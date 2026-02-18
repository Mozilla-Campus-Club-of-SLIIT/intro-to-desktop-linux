package leaderboard

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	pb "github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/pb"
	"github.com/valkey-io/valkey-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the DetachmentServiceServer interface for leaderboard operations.
type Server struct {
	pb.UnimplementedDetachmentServiceServer // Embed for forward compatibility
	db                                      *engine.DB
	authService                             *auth.AuthService
}

// NewLeaderboardServer creates a new Leaderboard server.
func NewLeaderboardServer(db *engine.DB, authService *auth.AuthService) *Server {
	return &Server{
		db:          db,
		authService: authService,
	}
}

// AddLeaderboardRecord adds a new record (username and score) to the leaderboard.
// This is a wrapper around db.UpdateScore.
func (s *Server) AddLeaderboardRecord(ctx context.Context, username string, score float64) error {
	return s.db.UpdateScore(ctx, username, score)
}

// GetLeaderboard implements the gRPC server-side streaming method.
func (s *Server) GetLeaderboard(req *pb.GetLeaderboardRequest, stream pb.DetachmentService_GetLeaderboardServer) error {
	ctx := stream.Context()

	// 1. Authenticate user
	verified, err := s.authService.VerifySession(ctx, req.GetUserId(), req.GetAccessToken())
	if err != nil {
		log.Printf("Authentication error for user %s: %v", req.GetUserId(), err)
		return status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !verified {
		return status.Errorf(codes.Unauthenticated, "unauthorized access")
	}

	// 2. Send initial full leaderboard
	initialLeaderboard, err := s.db.GetLeaderboard(ctx, -1) // -1 to get all
	if err != nil {
		log.Printf("Error getting initial leaderboard: %v", err)
		return status.Errorf(codes.Internal, "failed to retrieve initial leaderboard: %v", err)
	}

	resp := &pb.GetLeaderboardResponse{}
	for _, entry := range initialLeaderboard {
		resp.Records = append(resp.Records, &pb.LeaderboardRecord{
			Username: entry.Member,
			Score:    entry.Score,
		})
	}
	if err := stream.Send(resp); err != nil {
		log.Printf("Error sending initial leaderboard: %v", err)
		return status.Errorf(codes.Internal, "failed to send initial leaderboard: %v", err)
	}

	// 3. Keep connection open and send updates
	updateCh := make(chan *pb.GetLeaderboardResponse, 10) // Buffered channel

	go func() {
		// Create a new context for the subscription that is derived from the main stream context,
		// but use a context that can be cancelled if the gRPC stream context is cancelled.
		// This ensures the subscription Goroutine can clean up.
		subCtx, cancelSub := context.WithCancel(ctx)
		defer cancelSub()

		err := s.db.Subscribe(subCtx, []string{engine.ChannelLeaderboard}, func(msg valkey.PubSubMessage) {
			parts := strings.Split(msg.Message, ":")
			if len(parts) == 3 && parts[0] == "update" {
				_ = parts[1]
				_, parseErr := strconv.ParseFloat(parts[2], 64)
				if parseErr != nil {
					log.Printf("Error parsing score from leaderboard update message '%s': %v", msg.Message, parseErr)
					return
				}

				updatedLeaderboard, err := s.db.GetLeaderboard(subCtx, -1) // Use subCtx for DB call
				if err != nil {
					log.Printf("Error re-fetching leaderboard for update: %v", err)
					return
				}

				updateResp := &pb.GetLeaderboardResponse{}
				for _, entry := range updatedLeaderboard {
					updateResp.Records = append(updateResp.Records, &pb.LeaderboardRecord{
						Username: entry.Member,
						Score:    entry.Score,
					})
				}
				select {
				case updateCh <- updateResp:
				case <-subCtx.Done(): // Check if subscription context is done
					log.Println("Subscription context done, stopping leaderboard update send to channel")
					return
				default:
					log.Println("Leaderboard update channel is full, dropping update")
				}
			} else {
				log.Printf("Received malformed leaderboard update message: %s", msg.Message)
			}
		})
		if err != nil {
			log.Printf("Leaderboard PubSub subscription error: %v", err)
		}
		// Ensure the channel is closed when the subscription goroutine finishes
		close(updateCh)
	}()

	// Continuously send updates from the channel to the gRPC stream
	for {
		select {
		case update, ok := <-updateCh:
			if !ok {
				// Channel closed, subscription ended
				return nil
			}
			if err := stream.Send(update); err != nil {
				log.Printf("Error sending leaderboard update to stream: %v", err)
				return status.Errorf(codes.Internal, "failed to send leaderboard update: %v", err)
			}
		case <-ctx.Done(): // Original stream context
			log.Println("Client disconnected from leaderboard stream.")
			return nil
		}
	}
}

// ReportCommand is not implemented here, as this server is for leaderboard only.
// It must exist to satisfy the DetachmentServiceServer interface.
func (s *Server) ReportCommand(ctx context.Context, req *pb.ReportCommandRequest) (*pb.ReportCommandResponse, error) {
	// This method will be implemented in another server or combined.
	return nil, status.Errorf(codes.Unimplemented, "method ReportCommand not implemented")
}
