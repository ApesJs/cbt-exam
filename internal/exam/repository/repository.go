package repository

import (
	"context"
	"errors"

	"cbt-exam/internal/exam/domain"
)

type ExamRepository interface {
	// Exam operations
	Create(ctx context.Context, exam *domain.Exam) error
	GetByID(ctx context.Context, id string) (*domain.Exam, error)
	List(ctx context.Context, teacherID string, limit int32, offset int32) ([]*domain.Exam, error)
	Update(ctx context.Context, exam *domain.Exam) error
	Delete(ctx context.Context, id string) error

	// Exam status operations
	UpdateStatus(ctx context.Context, examID string, status domain.ExamState) error
	GetStatus(ctx context.Context, examID string) (*domain.ExamStatus, error)
	UpdateStudentStatus(ctx context.Context, examID string, studentStatus *domain.StudentStatus) error
}

// Errors
var (
	ErrExamNotFound      = errors.New("exam not found")
	ErrExamAlreadyExists = errors.New("exam already exists")
	ErrInvalidExamState  = errors.New("invalid exam state")
)
