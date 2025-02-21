package unit

import (
	"context"
	"github.com/ApesJs/cbt-exam/internal/exam/service"
	"testing"
	"time"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	"github.com/ApesJs/cbt-exam/internal/exam/domain"
	"github.com/ApesJs/cbt-exam/internal/exam/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockExamRepository adalah struct mock untuk repository.ExamRepository
type MockExamRepository struct {
	mock.Mock
}

// Implementasi metode-metode dari repository.ExamRepository
func (m *MockExamRepository) Create(ctx context.Context, exam *domain.Exam) error {
	args := m.Called(ctx, exam)
	return args.Error(0)
}

func (m *MockExamRepository) GetByID(ctx context.Context, id string) (*domain.Exam, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Exam), args.Error(1)
}

func (m *MockExamRepository) List(ctx context.Context, teacherID string, limit int32, offset int32) ([]*domain.Exam, error) {
	args := m.Called(ctx, teacherID, limit, offset)
	return args.Get(0).([]*domain.Exam), args.Error(1)
}

func (m *MockExamRepository) Update(ctx context.Context, exam *domain.Exam) error {
	args := m.Called(ctx, exam)
	return args.Error(0)
}

func (m *MockExamRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExamRepository) UpdateStatus(ctx context.Context, examID string, status domain.ExamState) error {
	args := m.Called(ctx, examID, status)
	return args.Error(0)
}

func (m *MockExamRepository) GetStatus(ctx context.Context, examID string) (*domain.ExamStatus, error) {
	args := m.Called(ctx, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExamStatus), args.Error(1)
}

func (m *MockExamRepository) UpdateStudentStatus(ctx context.Context, examID string, studentStatus *domain.StudentStatus) error {
	args := m.Called(ctx, examID, studentStatus)
	return args.Error(0)
}

// TestCreateExam menguji pembuatan ujian baru
func TestCreateExam(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockExamRepository)

	// Buat service dengan mock repository
	service := NewExamService(mockRepo)

	// Buat request untuk CreateExam
	req := &examv1.CreateExamRequest{
		Title:           "Ujian Matematika",
		Subject:         "Matematika",
		DurationMinutes: 60,
		TotalQuestions:  20,
		IsRandom:        true,
		TeacherId:       "teacher-123",
		ClassIds:        []string{"class-1", "class-2"},
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(exam *domain.Exam) bool {
		return exam.Title == req.Title &&
			exam.Subject == req.Subject &&
			exam.DurationMins == req.DurationMinutes &&
			exam.TotalQuestions == req.TotalQuestions &&
			exam.IsRandom == req.IsRandom &&
			exam.TeacherID == req.TeacherId
	})).Return(nil).Run(func(args mock.Arguments) {
		exam := args.Get(1).(*domain.Exam)
		exam.ID = "exam-123"
		exam.CreatedAt = time.Now()
		exam.UpdatedAt = time.Now()
	})

	// Panggil metode CreateExam
	ctx := context.Background()
	resp, err := service.CreateExam(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "exam-123", resp.Id)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Subject, resp.Subject)
	assert.Equal(t, req.DurationMinutes, resp.DurationMinutes)
	assert.Equal(t, req.TotalQuestions, resp.TotalQuestions)
	assert.Equal(t, req.IsRandom, resp.IsRandom)
	assert.Equal(t, req.TeacherId, resp.TeacherId)
	assert.Equal(t, req.ClassIds, resp.ClassIds)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestGetExam menguji pengambilan data ujian
func TestGetExam(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockExamRepository)

	// Buat service dengan mock repository
	service := NewExamService(mockRepo)

	// Buat data exam yang akan dikembalikan oleh mock
	examID := "exam-123"
	mockExam := &domain.Exam{
		ID:             examID,
		Title:          "Ujian Matematika",
		Subject:        "Matematika",
		DurationMins:   60,
		TotalQuestions: 20,
		IsRandom:       true,
		TeacherID:      "teacher-123",
		ClassIDs:       []string{"class-1", "class-2"},
		Status:         domain.ExamStateCreated,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("GetByID", mock.Anything, examID).Return(mockExam, nil)

	// Panggil metode GetExam
	ctx := context.Background()
	resp, err := service.GetExam(ctx, &examv1.GetExamRequest{Id: examID})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, examID, resp.Id)
	assert.Equal(t, mockExam.Title, resp.Title)
	assert.Equal(t, mockExam.Subject, resp.Subject)
	assert.Equal(t, mockExam.DurationMins, resp.DurationMinutes)
	assert.Equal(t, mockExam.TotalQuestions, resp.TotalQuestions)
	assert.Equal(t, mockExam.IsRandom, resp.IsRandom)
	assert.Equal(t, mockExam.TeacherID, resp.TeacherId)
	assert.Equal(t, mockExam.ClassIDs, resp.ClassIds)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestGetExamNotFound menguji pengambilan data ujian yang tidak ditemukan
func TestGetExamNotFound(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockExamRepository)

	// Buat service dengan mock repository
	service := NewExamService(mockRepo)

	// Setup ekspektasi untuk mock
	examID := "non-existent-exam"
	mockRepo.On("GetByID", mock.Anything, examID).Return(nil, repository.ErrExamNotFound)

	// Panggil metode GetExam
	ctx := context.Background()
	resp, err := service.GetExam(ctx, &examv1.GetExamRequest{Id: examID})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// Helper function untuk membuat ExamService
func NewExamService(repo repository.ExamRepository) examv1.ExamServiceServer {
	return service.NewExamService(repo)
}
