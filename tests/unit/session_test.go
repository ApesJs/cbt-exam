package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
	"github.com/ApesJs/cbt-exam/internal/session/domain"
	"github.com/ApesJs/cbt-exam/internal/session/repository"
	"github.com/ApesJs/cbt-exam/internal/session/service"
)

// MockSessionRepository adalah struct mock untuk repository.SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Implementasi metode-metode dari repository.SessionRepository
func (m *MockSessionRepository) StartSession(ctx context.Context, session *domain.ExamSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSession(ctx context.Context, id string) (*domain.ExamSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExamSession), args.Error(1)
}

func (m *MockSessionRepository) UpdateSessionStatus(ctx context.Context, id string, status domain.SessionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSessionRepository) FinishSession(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionRepository) SubmitAnswer(ctx context.Context, sessionID string, answer domain.Answer) error {
	args := m.Called(ctx, sessionID, answer)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSessionAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]domain.Answer), args.Error(1)
}

func (m *MockSessionRepository) IsExamActive(ctx context.Context, examID string) (bool, error) {
	args := m.Called(ctx, examID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockSessionRepository) HasActiveSession(ctx context.Context, studentID string) (bool, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).(bool), args.Error(1)
}

// TestStartSession menguji pembuatan sesi ujian baru
func TestStartSession(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	examID := "exam-123"
	studentID := "student-456"

	// Setup ekspektasi untuk mock
	mockRepo.On("IsExamActive", mock.Anything, examID).Return(true, nil)
	mockRepo.On("HasActiveSession", mock.Anything, studentID).Return(false, nil)
	mockRepo.On("StartSession", mock.Anything, mock.MatchedBy(func(session *domain.ExamSession) bool {
		return session.ExamID == examID &&
			session.StudentID == studentID &&
			session.Status == domain.SessionStatusStarted
	})).Return(nil).Run(func(args mock.Arguments) {
		session := args.Get(1).(*domain.ExamSession)
		session.ID = "session-789"
		session.CreatedAt = time.Now()
		session.UpdatedAt = time.Now()
	})

	// Panggil metode StartSession
	ctx := context.Background()
	resp, err := service.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    examID,
		StudentId: studentID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "session-789", resp.Id)
	assert.Equal(t, examID, resp.ExamId)
	assert.Equal(t, studentID, resp.StudentId)
	assert.Equal(t, sessionv1.SessionStatus_SESSION_STATUS_STARTED, resp.Status)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestStartSessionExamNotActive menguji pembuatan sesi ujian untuk ujian yang tidak aktif
func TestStartSessionExamNotActive(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	examID := "exam-123"
	studentID := "student-456"

	// Setup ekspektasi untuk mock
	mockRepo.On("IsExamActive", mock.Anything, examID).Return(false, nil)

	// Panggil metode StartSession
	ctx := context.Background()
	resp, err := service.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    examID,
		StudentId: studentID,
	})

	// Assertions - karena ujian tidak aktif, harusnya error
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestStartSessionAlreadyHasActiveSession menguji pembuatan sesi ujian untuk siswa yang sudah memiliki sesi aktif
func TestStartSessionAlreadyHasActiveSession(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	examID := "exam-123"
	studentID := "student-456"

	// Setup ekspektasi untuk mock
	mockRepo.On("IsExamActive", mock.Anything, examID).Return(true, nil)
	mockRepo.On("HasActiveSession", mock.Anything, studentID).Return(true, nil)

	// Panggil metode StartSession
	ctx := context.Background()
	resp, err := service.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    examID,
		StudentId: studentID,
	})

	// Assertions - karena siswa sudah memiliki sesi aktif, harusnya error
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestSubmitAnswer menguji pengiriman jawaban untuk pertanyaan ujian
func TestSubmitAnswer(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	sessionID := "session-789"
	questionID := "question-123"
	selectedChoice := "B"

	// Setup ekspektasi untuk mock
	mockRepo.On("SubmitAnswer", mock.Anything, sessionID, mock.MatchedBy(func(answer domain.Answer) bool {
		return answer.QuestionID == questionID &&
			answer.SelectedChoice == selectedChoice
	})).Return(nil)

	// Panggil metode SubmitAnswer
	ctx := context.Background()
	resp, err := service.SubmitAnswer(ctx, &sessionv1.SubmitAnswerRequest{
		SessionId:      sessionID,
		QuestionId:     questionID,
		SelectedChoice: selectedChoice,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "successfully")

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestFinishSession menguji penyelesaian sesi ujian
func TestFinishSession(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	sessionID := "session-789"
	examID := "exam-123"
	studentID := "student-456"
	startTime := time.Now().Add(-30 * time.Minute)

	// Setup mock session yang akan dikembalikan
	mockSession := &domain.ExamSession{
		ID:        sessionID,
		ExamID:    examID,
		StudentID: studentID,
		Status:    domain.SessionStatusInProgress,
		StartTime: startTime,
		Answers:   []domain.Answer{},
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("GetSession", mock.Anything, sessionID).Return(mockSession, nil)
	mockRepo.On("FinishSession", mock.Anything, sessionID).Return(nil)

	// Panggil metode FinishSession
	ctx := context.Background()
	resp, err := service.FinishSession(ctx, &sessionv1.FinishSessionRequest{
		Id: sessionID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, sessionID, resp.Id)
	assert.Equal(t, examID, resp.ExamId)
	assert.Equal(t, studentID, resp.StudentId)
	assert.Equal(t, sessionv1.SessionStatus_SESSION_STATUS_FINISHED, resp.Status)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestGetRemainingTime menguji pengambilan sisa waktu ujian
func TestGetRemainingTime(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockSessionRepository)

	// Buat service dengan mock repository
	service := NewSessionService(mockRepo)

	// Setup data
	sessionID := "session-789"
	examID := "exam-123"
	studentID := "student-456"
	startTime := time.Now().Add(-30 * time.Minute) // 30 menit yang lalu

	// Setup mock session yang akan dikembalikan
	mockSession := &domain.ExamSession{
		ID:        sessionID,
		ExamID:    examID,
		StudentID: studentID,
		Status:    domain.SessionStatusInProgress,
		StartTime: startTime,
		Answers:   []domain.Answer{},
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("GetSession", mock.Anything, sessionID).Return(mockSession, nil)

	// Panggil metode GetRemainingTime
	ctx := context.Background()
	resp, err := service.GetRemainingTime(ctx, &sessionv1.GetRemainingTimeRequest{
		SessionId: sessionID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Karena durasi ujian diatur sebagai konstanta dalam service (misalnya 120 menit)
	// dan waktu mulai adalah 30 menit yang lalu, maka sisa waktu seharusnya sekitar 90 menit
	// tetapi karena kita menggunakan waktu sebenarnya, kita hanya bisa memastikan bahwa
	// sisa waktu lebih besar dari 0 dan kurang dari atau sama dengan 90 menit
	assert.GreaterOrEqual(t, resp.RemainingMinutes, int32(0))
	assert.LessOrEqual(t, resp.RemainingMinutes, int32(90))
	assert.GreaterOrEqual(t, resp.RemainingSeconds, int32(0))
	assert.Less(t, resp.RemainingSeconds, int32(60))

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// Helper function untuk membuat SessionService
func NewSessionService(repo repository.SessionRepository) sessionv1.SessionServiceServer {
	return service.NewSessionService(repo)
}
