package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	"github.com/ApesJs/cbt-exam/internal/question/repository/postgres"
	"github.com/ApesJs/cbt-exam/internal/question/service"
	"github.com/ApesJs/cbt-exam/pkg/client"
	"github.com/ApesJs/cbt-exam/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override port from environment variable if exists
	if portStr := os.Getenv("PORT"); portStr != "" {
		if portNum, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = portNum
			log.Printf("Using port from environment: %d", portNum)
		} else {
			log.Printf("Invalid PORT environment variable: %s", portStr)
		}
	}

	// Initialize service client
	pkgClient, err := client.NewServiceClient(
		cfg.ExamPort,
		cfg.QuestionPort,
		cfg.SessionPort,
		cfg.ScoringPort,
	)
	if err != nil {
		log.Fatalf("Failed to create service client: %v", err)
	}

	// Initialize PostgreSQL connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize repository
	repo := postgres.NewPostgresRepository(db)

	// Initialize service
	svc := service.NewQuestionService(repo, pkgClient)

	// Initialize gRPC server
	server := grpc.NewServer()
	questionv1.RegisterQuestionServiceServer(server, svc)

	// Enable reflection for development tools
	reflection.Register(server)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Channel untuk notifikasi server shutdown
	shutdownComplete := make(chan struct{})

	// Start server in a goroutine
	go func() {
		log.Printf("Starting question service on port %d", cfg.Port)
		if err := server.Serve(lis); err != nil {
			log.Printf("Server error: %v", err)
		}
		close(shutdownComplete)
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Received shutdown signal. Initiating graceful shutdown...")

	// Gracefully stop the server
	server.GracefulStop()

	// Wait for server to complete shutdown
	<-shutdownComplete
	log.Println("Server shutdown complete")

	// Additional cleanup
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
