package repository

import (
	"context"
	"errors"

	"cbt-exam/internal/question/domain"
)

type QuestionRepository interface {
	// Basic CRUD
	Create(ctx context.Context, question *domain.Question) error
	GetByID(ctx context.Context, id string) (*domain.Question, error)
	List(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.Question, error)
	Update(ctx context.Context, question *domain.Question) error
	Delete(ctx context.Context, id string) error

	// Specific to exam questions
	GetExamQuestions(ctx context.Context, filter domain.QuestionFilter) ([]*domain.Question, error)
	CountExamQuestions(ctx context.Context, examID string) (int32, error)
}

// Errors
var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrInvalidQuestion  = errors.New("invalid question data")
	ErrExamNotFound     = errors.New("exam not found")
)
