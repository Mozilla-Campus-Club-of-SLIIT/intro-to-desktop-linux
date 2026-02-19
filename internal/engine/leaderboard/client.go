package leaderboard

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	pb "github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const (
	// grpcServerAddress is the address of the gRPC server.
	grpcServerAddress = "mozlive-grpc.mrbhanuka.dev:50051"
)

// MsgLeaderboardUpdate is the message sent to the TUI when the leaderboard state changes.
type MsgLeaderboardUpdate []Entry

// MsgLeaderboardError is the message sent to the TUI when a stream error occurs.
type MsgLeaderboardError error

// Entry represents a single record in the leaderboard.
type Entry struct {
	Username string
	Score    float64
}

var (
	globalClient *LeaderboardClient
	clientOnce   sync.Once
)

// GetClient returns the singleton leaderboard client.
func GetClient(ctx context.Context) (*LeaderboardClient, error) {
	var err error
	clientOnce.Do(func() {
		globalClient, err = NewLeaderboardClient()
		if err == nil {
			// Start the single long-lived background stream
			go globalClient.runBackground()
		}
	})
	return globalClient, err
}

// LeaderboardClient provides methods to interact with the leaderboard gRPC service and maintains state.
type LeaderboardClient struct {
	client pb.DetachmentServiceClient
	conn   *grpc.ClientConn

	mu                 sync.RWMutex
	currentLeaderboard []Entry
	leaderboardErr     error

	subscribers []chan MsgLeaderboardUpdate
	errSubs     []chan MsgLeaderboardError
}

// NewLeaderboardClient creates a new LeaderboardClient and establishes a gRPC connection.
func NewLeaderboardClient() (*LeaderboardClient, error) {
	// Adjusted keepalive to be less aggressive to avoid "too_many_pings" (ENHANCE_YOUR_CALM)
	// Server likely has a MinTime policy.
	kacp := keepalive.ClientParameters{
		Time:                60 * time.Second, // Ping every 60 seconds
		Timeout:             20 * time.Second, // Wait 20 seconds for response
		PermitWithoutStream: false,            // Only ping if there is an active stream
	}

	conn, err := grpc.NewClient(
		grpcServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	client := pb.NewDetachmentServiceClient(conn)
	return &LeaderboardClient{
		client:      client,
		conn:        conn,
		subscribers: make([]chan MsgLeaderboardUpdate, 0),
		errSubs:     make([]chan MsgLeaderboardError, 0),
	}, nil
}

// Close closes the gRPC client connection.
func (lc *LeaderboardClient) Close() error {
	if lc.conn != nil {
		return lc.conn.Close()
	}
	return nil
}

// GetEntries returns a copy of the current leaderboard entries.
func (lc *LeaderboardClient) GetEntries() []Entry {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	entries := make([]Entry, len(lc.currentLeaderboard))
	copy(entries, lc.currentLeaderboard)
	return entries
}

// Start establishes a subscription to the single background leaderboard stream.
func (lc *LeaderboardClient) Start(ctx context.Context) (<-chan MsgLeaderboardUpdate, <-chan MsgLeaderboardError) {
	updateCh := make(chan MsgLeaderboardUpdate, 1)
	errCh := make(chan MsgLeaderboardError, 1)

	lc.mu.Lock()
	lc.subscribers = append(lc.subscribers, updateCh)
	lc.errSubs = append(lc.errSubs, errCh)

	// Immediately send current state if available
	if len(lc.currentLeaderboard) > 0 {
		entries := make([]Entry, len(lc.currentLeaderboard))
		copy(entries, lc.currentLeaderboard)
		updateCh <- MsgLeaderboardUpdate(entries)
	}
	if lc.leaderboardErr != nil {
		errCh <- MsgLeaderboardError(lc.leaderboardErr)
	}
	lc.mu.Unlock()

	go func() {
		<-ctx.Done()
		lc.mu.Lock()
		defer lc.mu.Unlock()

		for i, sub := range lc.subscribers {
			if sub == updateCh {
				lc.subscribers = append(lc.subscribers[:i], lc.subscribers[i+1:]...)
				break
			}
		}
		for i, sub := range lc.errSubs {
			if sub == errCh {
				lc.errSubs = append(lc.errSubs[:i], lc.errSubs[i+1:]...)
				break
			}
		}
	}()

	return updateCh, errCh
}

// runBackground manages the single gRPC stream in the background.
func (lc *LeaderboardClient) runBackground() {
	bgCtx := context.Background()
	authClient := auth.GetClient()

	for {
		userID, err := authClient.GetUserID()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		accessToken, err := authClient.GetAccessToken()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		req := &pb.GetLeaderboardRequest{
			UserId:      userID,
			AccessToken: accessToken,
		}

		stream, err := lc.client.GetLeaderboard(bgCtx, req)
		if err != nil {
			lc.broadcastError(fmt.Errorf("failed to open leaderboard stream: %w", err))
			time.Sleep(10 * time.Second)
			continue
		}

		lc.updateState(nil, nil)

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				if status.Code(err) != codes.Canceled {
					lc.broadcastError(fmt.Errorf("leaderboard stream lost: %w", err))
				}
				break
			}

			var entries []Entry
			for _, record := range resp.Records {
				entries = append(entries, Entry{
					Username: record.Username,
					Score:    record.Score,
				})
			}

			lc.sortEntries(entries)
			if lc.updateState(entries, nil) {
				lc.broadcastUpdate(entries)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (lc *LeaderboardClient) broadcastUpdate(entries []Entry) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	msg := MsgLeaderboardUpdate(entries)
	for _, sub := range lc.subscribers {
		select {
		case sub <- msg:
		default:
		}
	}
}

func (lc *LeaderboardClient) broadcastError(err error) {
	if !lc.updateState(nil, err) {
		return
	}
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	msg := MsgLeaderboardError(err)
	for _, sub := range lc.errSubs {
		select {
		case sub <- msg:
		default:
		}
	}
}

func (lc *LeaderboardClient) updateState(entries []Entry, err error) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	changed := false

	// Handle error update
	if err != nil {
		if lc.leaderboardErr == nil || lc.leaderboardErr.Error() != err.Error() {
			lc.leaderboardErr = err
			changed = true
		}
	} else if lc.leaderboardErr != nil {
		lc.leaderboardErr = nil
		changed = true
	}

	// Handle entries update
	if entries != nil {
		isEqual := len(entries) == len(lc.currentLeaderboard)
		if isEqual {
			for i := range entries {
				if entries[i].Username != lc.currentLeaderboard[i].Username || entries[i].Score != lc.currentLeaderboard[i].Score {
					isEqual = false
					break
				}
			}
		}

		if !isEqual {
			lc.currentLeaderboard = entries
			changed = true
		}
	}

	return changed
}

func (lc *LeaderboardClient) sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Username < entries[j].Username
	})
}
