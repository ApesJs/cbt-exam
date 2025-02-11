package repository

import (
	"context"
	"errors"

	"github.com/ApesJs/cbt-exam/internal/session/domain"
)

type SessionRepository interface {
	// Session management
	StartSession(ctx context.Context, session *domain.ExamSession) error
	GetSession(ctx context.Context, id string) (*domain.ExamSession, error)
	UpdateSessionStatus(ctx context.Context, id string, status domain.SessionStatus) error
	FinishSession(ctx context.Context, id string) error

	// Answer management
	SubmitAnswer(ctx context.Context, sessionID string, answer domain.Answer) error
	GetSessionAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error)

	// Validations
	IsExamActive(ctx context.Context, examID string) (bool, error)
	HasActiveSession(ctx context.Context, studentID string) (bool, error)
}

// Errors
var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrExamNotFound        = errors.New("exam not found")
	ErrExamNotActive       = errors.New("exam is not active")
	ErrDuplicateSession    = errors.New("student already has an active session")
	ErrAnswerExists        = errors.New("answer already exists for this question")
	ErrInvalidSessionState = errors.New("invalid session state")
)
