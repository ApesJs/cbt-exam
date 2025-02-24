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
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	"github.com/ApesJs/cbt-exam/internal/scoring/repository/postgres"
	"github.com/ApesJs/cbt-exam/internal/scoring/service"
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
			cfg.ScoringPort = portNum
			log.Printf("Using port from environment: %d", portNum)
		} else {
			log.Printf("Invalid PORT environment variable: %s", portStr)
		}
	}

	// Initialize PostgreSQL connection with retry mechanism
	var db *sql.DB
	maxRetries := 5
	retryDelay := time.Second * 3

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			log.Fatalf("Failed to connect to database after %d attempts", maxRetries)
		}

		// Test database connection
		if err := db.Ping(); err != nil {
			log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			log.Fatalf("Failed to ping database after %d attempts", maxRetries)
		}

		break
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Configure database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Initialize repository
	repo := postgres.NewPostgresRepository(db)

	// Initialize service
	svc := service.NewScoringService(repo)

	// Initialize gRPC server with options
	serverOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.MaxRecvMsgSize(4 * 1024 * 1024), // 4MB
	}
	server := grpc.NewServer(serverOptions...)
	scoringv1.RegisterScoringServiceServer(server, svc)

	// Enable reflection for development tools
	reflection.Register(server)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ScoringPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Channel untuk notifikasi server shutdown
	shutdownComplete := make(chan struct{})
	serverError := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting scoring service on port %d", cfg.ScoringPort)
		if err := server.Serve(lis); err != nil {
			serverError <- fmt.Errorf("failed to serve: %v", err)
			return
		}
		close(shutdownComplete)
	}()

	// Monitor channels untuk shutdown atau error
	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal. Initiating graceful shutdown...")

		// Create shutdown timeout context
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Graceful shutdown dengan timeout
		shutdownChan := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(shutdownChan)
		}()

		select {
		case <-shutdownChan:
			log.Println("Server shutdown gracefully")
		case <-shutdownCtx.Done():
			log.Println("Server shutdown timeout exceeded, forcing shutdown")
			server.Stop()
		}

	case err := <-serverError:
		log.Printf("Server error: %v", err)
	}

	// Final cleanup
	log.Println("Cleaning up resources...")

	// Close database connections
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	// Wait for all connections to close
	select {
	case <-shutdownComplete:
		log.Println("Server shutdown complete")
	case <-time.After(5 * time.Second):
		log.Println("Server shutdown timed out")
	}

	log.Println("Service stopped")
}
