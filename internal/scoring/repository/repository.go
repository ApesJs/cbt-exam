package repository

import (
	"context"
	"errors"

	"github.com/ApesJs/cbt-exam/internal/scoring/domain"
)

type ScoringRepository interface {
	// Score operations
	CreateScore(ctx context.Context, score *domain.ExamScore) error
	GetScore(ctx context.Context, id string) (*domain.ExamScore, error)
	GetScoreByExamAndStudent(ctx context.Context, examID, studentID string) (*domain.ExamScore, error)
	ListScores(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.ExamScore, error)

	// Answer validation
	GetCorrectAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error)
	GetStudentAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error)
}

// Errors
var (
	ErrScoreNotFound   = errors.New("score not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrExamNotFound    = errors.New("exam not found")
	ErrDuplicateScore  = errors.New("score already exists for this exam and student")
)
