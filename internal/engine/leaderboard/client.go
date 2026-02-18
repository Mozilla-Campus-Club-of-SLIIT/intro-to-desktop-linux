package leaderboard

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"time"

	pb "github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// grpcServerAddress is the address of the gRPC server.
	grpcServerAddress = "mozlive-grpc.mrbhanuka.dev:443" // From deploy/compose.yaml
)

// Entry represents a single record in the leaderboard.
type Entry struct {
	Username string
	Score    float64 // Changed to float64 to match protobuf definition
}

// LeaderboardClient provides methods to interact with the leaderboard gRPC service.
type LeaderboardClient struct {
	client pb.DetachmentServiceClient
	conn   *grpc.ClientConn
}

// NewLeaderboardClient creates a new LeaderboardClient and establishes a gRPC connection.
func NewLeaderboardClient(ctx context.Context) (*LeaderboardClient, error) {
	conn, err := grpc.DialContext(ctx, grpcServerAddress, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
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

// GetLeaderboardStream establishes a streaming connection to the leaderboard service.
// It returns a channel for leaderboard updates and an error channel.
func (lc *LeaderboardClient) GetLeaderboardStream(ctx context.Context, userID, accessToken string) (<-chan []Entry, <-chan error) {
	updateCh := make(chan []Entry)
	errCh := make(chan error, 1) // Buffered to prevent goroutine leak if error happens before reader is ready

	go func() {
		defer close(updateCh)
		defer close(errCh)

		req := &pb.GetLeaderboardRequest{
			UserId:      userID,
			AccessToken: accessToken,
		}

		stream, err := lc.client.GetLeaderboard(ctx, req)
		if err != nil {
			errCh <- fmt.Errorf("failed to open leaderboard stream: %w", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					log.Println("Leaderboard stream closed by server.")
					return
				}
				if err != nil {
					errCh <- fmt.Errorf("failed to receive leaderboard update: %w", err)
					return
				}

				var currentLeaderboard []Entry
				for _, record := range resp.Records {
					currentLeaderboard = append(currentLeaderboard, Entry{
						Username: record.Username,
						Score:    record.Score,
					})
				}

				// Re-sort the leaderboard just in case, or for client-specific ordering needs
				sort.Slice(currentLeaderboard, func(i, j int) bool {
					return currentLeaderboard[i].Score > currentLeaderboard[j].Score
				})

				select {
				case updateCh <- currentLeaderboard:
					// Successfully sent update
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case <-time.After(5 * time.Second): // Timeout for sending update
					log.Println("Timeout sending leaderboard update, client not consuming fast enough.")
					// Depending on requirements, could return error or drop update
				}
			}
		}
	}()

	return updateCh, errCh
}
