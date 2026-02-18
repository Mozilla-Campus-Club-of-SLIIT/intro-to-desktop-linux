package main

import (
	"log"
	"net"
	"os"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/auth"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/leaderboard"
	pb "github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // Import for gRPC reflection (optional, but good for debugging)
)

const (
	grpcPort = "50051"
)

func main() {
	// Initialize database connection
	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		dbAddr = "localhost:6379" // Default for local development
	}
	db, err := engine.NewDB(dbAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to Valkey at %s", dbAddr)

	// Initialize authentication service
	authService := auth.NewAuthService(db)

	// Initialize leaderboard gRPC server implementation
	leaderboardServer := leaderboard.NewLeaderboardServer(db, authService)

	// Create a TCP listener for gRPC
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server listening on port %s", grpcPort)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register our leaderboard service implementation with the gRPC server
	pb.RegisterDetachmentServiceServer(grpcServer, leaderboardServer)

	// Register reflection service on gRPC server.
	// This is optional but useful for debugging with tools like grpcurl.
	reflection.Register(grpcServer)

	// Start serving gRPC requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
