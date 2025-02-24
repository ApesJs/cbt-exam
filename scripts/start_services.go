package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Service struct {
	name string
	port int
	dir  string
	file string
	cmd  *exec.Cmd
}

func main() {
	services := []Service{
		{"exam", 50051, "cmd/exam", "main.go", nil},
		{"question", 50052, "cmd/question", "main.go", nil},
		{"session", 50053, "cmd/session", "main.go", nil},
		{"scoring", 50054, "cmd/scoring", "main.go", nil},
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel untuk menangkap error dari service
	errorChan := make(chan error, len(services))

	// Jalankan setiap service
	for i := range services {
		svc := &services[i] // Gunakan pointer agar bisa mengupdate cmd
		fmt.Printf("Starting %s service on port %d\n", svc.name, svc.port)

		if runtime.GOOS == "windows" {
			svc.cmd = exec.Command("go", "run", fmt.Sprintf(".\\%s\\%s", svc.dir, svc.file))
		} else {
			svc.cmd = exec.Command("go", "run", fmt.Sprintf("./%s/%s", svc.dir, svc.file))
		}

		// Set environment variables
		svc.cmd.Env = append(os.Environ(),
			fmt.Sprintf("PORT=%d", svc.port),
			fmt.Sprintf("SERVICE_NAME=%s", svc.name),
		)

		svc.cmd.Stdout = os.Stdout
		svc.cmd.Stderr = os.Stderr

		// Start command
		if err := svc.cmd.Start(); err != nil {
			log.Printf("Failed to start %s service: %v", svc.name, err)
			continue
		}

		// Monitor service
		go func(s *Service) {
			if err := s.cmd.Wait(); err != nil {
				errorChan <- fmt.Errorf("%s service stopped with error: %v", s.name, err)
			}
		}(svc)

		// Tunggu sebentar sebelum menjalankan service berikutnya
		time.Sleep(3 * time.Second)
	}

	// Wait for interrupt signal
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		// Graceful shutdown
		for _, svc := range services {
			if svc.cmd != nil && svc.cmd.Process != nil {
				log.Printf("Stopping %s service...", svc.name)
				if err := svc.cmd.Process.Kill(); err != nil {
					log.Printf("Error stopping %s service: %v", svc.name, err)
				}
			}
		}
	case err := <-errorChan:
		log.Printf("Service error: %v", err)
	}

	log.Println("All services stopped")
}
