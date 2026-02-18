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
	"google.golang.org/grpc/status"
)

const (
	// grpcServerAddress is the address of the gRPC server.
	grpcServerAddress = "mozlive-grpc.mrbhanuka.dev:50051"
)

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
		globalClient, err = NewLeaderboardClient(ctx)
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
}

// NewLeaderboardClient creates a new LeaderboardClient and establishes a gRPC connection.
func NewLeaderboardClient(ctx context.Context) (*LeaderboardClient, error) {
	conn, err := grpc.NewClient(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	client := pb.NewDetachmentServiceClient(conn)
	return &LeaderboardClient{
		client: client,
		conn:   conn,
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

	// Return a copy to avoid data race and external modification
	entries := make([]Entry, len(lc.currentLeaderboard))
	copy(entries, lc.currentLeaderboard)
	return entries
}

// GetError returns the last error encountered during streaming.
func (lc *LeaderboardClient) GetError() error {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.leaderboardErr
}

// Start establishes a streaming connection to the leaderboard service and maintains the internal state.
// It returns channels for updates and errors.
func (lc *LeaderboardClient) Start(ctx context.Context) (<-chan []Entry, <-chan error) {
	updateCh := make(chan []Entry)
	errCh := make(chan error, 1)

	go func() {
		defer close(updateCh)
		defer close(errCh)

		authClient := auth.GetClient()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			userID, err := authClient.GetUserID()
			if err != nil {
				lc.handleError(ctx, errCh, fmt.Errorf("failed to get user ID: %w", err))
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
					continue
				}
			}

			accessToken, err := authClient.GetAccessToken()
			if err != nil {
				lc.handleError(ctx, errCh, fmt.Errorf("failed to get access token: %w", err))
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
					continue
				}
			}

			req := &pb.GetLeaderboardRequest{
				UserId:      userID,
				AccessToken: accessToken,
			}

			stream, err := lc.client.GetLeaderboard(ctx, req)
			if err != nil {
				lc.handleError(ctx, errCh, fmt.Errorf("failed to open leaderboard stream: %w", err))
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}

			lc.updateState(nil, nil) // Clear error on successful connection

		recvLoop:
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break recvLoop
				}
				if err != nil {
					if status.Code(err) == codes.Canceled {
						return
					}
					lc.handleError(ctx, errCh, fmt.Errorf("failed to receive leaderboard update: %w", err))
					break recvLoop
				}

				var entries []Entry
				for _, record := range resp.Records {
					entries = append(entries, Entry{
						Username: record.Username,
						Score:    record.Score,
					})
				}

				lc.sortEntries(entries)
				lc.updateState(entries, nil)

				select {
				case updateCh <- entries:
				case <-ctx.Done():
					return
				default:
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()

	return updateCh, errCh
}

// handleError updates internal state and sends error to channel if possible.
func (lc *LeaderboardClient) handleError(ctx context.Context, errCh chan<- error, err error) {
	lc.updateState(nil, err)
	select {
	case errCh <- err:
	case <-ctx.Done():
	default:
	}
}

// updateState safely updates the internal leaderboard state.
func (lc *LeaderboardClient) updateState(entries []Entry, err error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if entries != nil {
		lc.currentLeaderboard = entries
	}
	lc.leaderboardErr = err
}

// sortEntries handles the sorting logic for the leaderboard.
func (lc *LeaderboardClient) sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		// Sort by score descending
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		// Then by username ascending for deterministic ordering
		return entries[i].Username < entries[j].Username
	})
}
