package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	// Initialize repository
	repo := postgres.NewPostgresRepository(db)

	// Initialize service
	svc := service.NewScoringService(repo)

	// Initialize gRPC server
	server := grpc.NewServer()
	scoringv1.RegisterScoringServiceServer(server, svc)

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

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		server.GracefulStop()
	}()

	log.Printf("Starting scoring service on port %d", cfg.Port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
