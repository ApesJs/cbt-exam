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

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	"github.com/ApesJs/cbt-exam/internal/exam/repository/postgres"
	"github.com/ApesJs/cbt-exam/internal/exam/service"
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
	svc := service.NewExamService(repo)

	// Initialize gRPC server
	server := grpc.NewServer()
	examv1.RegisterExamServiceServer(server, svc)

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

	// Start server in a goroutine
	go func() {
		log.Printf("Starting exam service on port %d", cfg.Port)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Gracefully stop the server
	server.GracefulStop()
	log.Println("Server stopped")
}
