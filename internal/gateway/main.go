package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ApesJs/cbt-exam/internal/gateway/handler"
	"github.com/ApesJs/cbt-exam/pkg/client"
	"github.com/ApesJs/cbt-exam/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize service client
	serviceClient, err := client.NewServiceClient(
		cfg.ExamPort,
		cfg.QuestionPort,
		cfg.SessionPort,
		cfg.ScoringPort,
	)
	if err != nil {
		log.Fatalf("Failed to create service client: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Initialize handlers
	examHandler := handler.NewExamHandler(serviceClient)
	questionHandler := handler.NewQuestionHandler(serviceClient)
	sessionHandler := handler.NewSessionHandler(serviceClient)
	scoringHandler := handler.NewScoringHandler(serviceClient)

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		// Exam routes
		exam := v1.Group("/exams")
		{
			exam.POST("", examHandler.CreateExam)
			exam.GET("", examHandler.ListExams)
			exam.GET("/:id", examHandler.GetExam)
			exam.PUT("/:id", examHandler.UpdateExam)
			exam.DELETE("/:id", examHandler.DeleteExam)
			exam.POST("/:id/activate", examHandler.ActivateExam)
			exam.POST("/:id/deactivate", examHandler.DeactivateExam)
		}

		// Question routes
		question := v1.Group("/questions")
		{
			question.POST("", questionHandler.CreateQuestion)
			question.GET("", questionHandler.ListQuestions)
			question.GET("/:id", questionHandler.GetQuestion)
			question.PUT("/:id", questionHandler.UpdateQuestion)
			question.DELETE("/:id", questionHandler.DeleteQuestion)
			question.GET("/exam/:examId", questionHandler.GetExamQuestions)
		}

		// Session routes
		session := v1.Group("/sessions")
		{
			session.POST("", sessionHandler.StartSession)
			session.GET("/:id", sessionHandler.GetSession)
			session.POST("/:id/answer", sessionHandler.SubmitAnswer)
			session.POST("/:id/finish", sessionHandler.FinishSession)
			session.GET("/:id/time", sessionHandler.GetRemainingTime)
		}

		// Scoring routes
		score := v1.Group("/scores")
		{
			score.POST("/calculate", scoringHandler.CalculateScore)
			score.GET("/:id", scoringHandler.GetScore)
			score.GET("/exam/:examId", scoringHandler.ListScores)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give 5 seconds for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
