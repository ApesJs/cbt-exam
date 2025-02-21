package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	"github.com/ApesJs/cbt-exam/internal/scoring/domain"
	"github.com/ApesJs/cbt-exam/internal/scoring/repository"
	"github.com/ApesJs/cbt-exam/internal/scoring/service"
)

// MockScoringRepository adalah struct mock untuk repository.ScoringRepository
type MockScoringRepository struct {
	mock.Mock
}

// Implementasi metode-metode dari repository.ScoringRepository
func (m *MockScoringRepository) CreateScore(ctx context.Context, score *domain.ExamScore) error {
	args := m.Called(ctx, score)
	return args.Error(0)
}

func (m *MockScoringRepository) GetScore(ctx context.Context, id string) (*domain.ExamScore, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExamScore), args.Error(1)
}

func (m *MockScoringRepository) GetScoreByExamAndStudent(ctx context.Context, examID, studentID string) (*domain.ExamScore, error) {
	args := m.Called(ctx, examID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExamScore), args.Error(1)
}

func (m *MockScoringRepository) ListScores(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.ExamScore, error) {
	args := m.Called(ctx, examID, limit, offset)
	return args.Get(0).([]*domain.ExamScore), args.Error(1)
}

func (m *MockScoringRepository) GetCorrectAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]domain.Answer), args.Error(1)
}

func (m *MockScoringRepository) GetStudentAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]domain.Answer), args.Error(1)
}

// TestCalculateScore menguji kalkulasi skor ujian
func TestCalculateScore(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockScoringRepository)

	// Buat service dengan mock repository
	service := NewScoringService(mockRepo)

	// Setup data
	sessionID := "session-789"
	examID := "exam-123"
	studentID := "student-456"

	// Setup mock answers yang akan dikembalikan
	mockAnswers := []domain.Answer{
		{
			QuestionID:    "question-1",
			CorrectAnswer: "A",
			StudentAnswer: "A", // Benar
		},
		{
			QuestionID:    "question-2",
			CorrectAnswer: "B",
			StudentAnswer: "B", // Benar
		},
		{
			QuestionID:    "question-3",
			CorrectAnswer: "C",
			StudentAnswer: "D", // Salah
		},
		{
			QuestionID:    "question-4",
			CorrectAnswer: "D",
			StudentAnswer: "", // Tidak dijawab
		},
		{
			QuestionID:    "question-5",
			CorrectAnswer: "A",
			StudentAnswer: "B", // Salah
		},
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("GetCorrectAnswers", mock.Anything, sessionID).Return(mockAnswers, nil)
	mockRepo.On("CreateScore", mock.Anything, mock.MatchedBy(func(score *domain.ExamScore) bool {
		return score.SessionID == sessionID &&
			score.TotalQuestions == int32(len(mockAnswers)) &&
			score.CorrectAnswers == int32(2) && // 2 jawaban benar
			score.UnansweredCount == int32(1) && // 1 tidak dijawab
			score.WrongAnswers == int32(2) // 2 jawaban salah
	})).Return(nil).Run(func(args mock.Arguments) {
		score := args.Get(1).(*domain.ExamScore)
		score.ID = "score-123"
		score.ExamID = examID
		score.StudentID = studentID
		score.CreatedAt = time.Now()
	})

	// Panggil metode CalculateScore
	ctx := context.Background()
	resp, err := service.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "score-123", resp.Id)
	assert.Equal(t, sessionID, resp.SessionId)
	assert.Equal(t, examID, resp.ExamId)
	assert.Equal(t, studentID, resp.StudentId)
	assert.Equal(t, int32(5), resp.TotalQuestions)
	assert.Equal(t, int32(2), resp.CorrectAnswers)
	assert.Equal(t, int32(2), resp.WrongAnswers)
	assert.Equal(t, int32(1), resp.Unanswered)
	assert.Equal(t, float32(40.0), resp.Score) // 2/5 * 100 = 40%

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestCalculateScoreSessionNotFound menguji kalkulasi skor untuk sesi yang tidak ditemukan
func TestCalculateScoreSessionNotFound(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockScoringRepository)

	// Buat service dengan mock repository
	service := NewScoringService(mockRepo)

	// Setup data
	sessionID := "non-existent-session"

	// Setup ekspektasi untuk mock
	mockRepo.On("GetCorrectAnswers", mock.Anything, sessionID).Return([]domain.Answer{}, repository.ErrSessionNotFound)

	// Panggil metode CalculateScore
	ctx := context.Background()
	resp, err := service.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})

	// Assertions - karena sesi tidak ditemukan, harusnya error
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestGetScore menguji pengambilan skor ujian
func TestGetScore(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockScoringRepository)

	// Buat service dengan mock repository
	service := NewScoringService(mockRepo)

	// Setup data
	scoreID := "score-123"
	sessionID := "session-789"
	examID := "exam-123"
	studentID := "student-456"

	// Setup mock score yang akan dikembalikan
	mockScore := &domain.ExamScore{
		ID:              scoreID,
		ExamID:          examID,
		SessionID:       sessionID,
		StudentID:       studentID,
		TotalQuestions:  5,
		CorrectAnswers:  3,
		WrongAnswers:    1,
		UnansweredCount: 1,
		Score:           60.0, // 3/5 * 100 = 60%
		CreatedAt:       time.Now(),
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("GetScore", mock.Anything, scoreID).Return(mockScore, nil)

	// Panggil metode GetScore
	ctx := context.Background()
	resp, err := service.GetScore(ctx, &scoringv1.GetScoreRequest{
		Id: scoreID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, scoreID, resp.Id)
	assert.Equal(t, sessionID, resp.SessionId)
	assert.Equal(t, examID, resp.ExamId)
	assert.Equal(t, studentID, resp.StudentId)
	assert.Equal(t, int32(5), resp.TotalQuestions)
	assert.Equal(t, int32(3), resp.CorrectAnswers)
	assert.Equal(t, int32(1), resp.WrongAnswers)
	assert.Equal(t, int32(1), resp.Unanswered)
	assert.Equal(t, float32(60.0), resp.Score)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// TestListScores menguji pengambilan daftar skor ujian
func TestListScores(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockScoringRepository)

	// Buat service dengan mock repository
	service := NewScoringService(mockRepo)

	// Setup data
	examID := "exam-123"
	limit := int32(10)
	offset := int32(0)

	// Setup mock scores yang akan dikembalikan
	mockScores := []*domain.ExamScore{
		{
			ID:              "score-1",
			ExamID:          examID,
			SessionID:       "session-1",
			StudentID:       "student-1",
			TotalQuestions:  5,
			CorrectAnswers:  5,
			WrongAnswers:    0,
			UnansweredCount: 0,
			Score:           100.0,
			CreatedAt:       time.Now(),
		},
		{
			ID:              "score-2",
			ExamID:          examID,
			SessionID:       "session-2",
			StudentID:       "student-2",
			TotalQuestions:  5,
			CorrectAnswers:  4,
			WrongAnswers:    1,
			UnansweredCount: 0,
			Score:           80.0,
			CreatedAt:       time.Now(),
		},
	}

	// Setup ekspektasi untuk mock
	mockRepo.On("ListScores", mock.Anything, examID, limit, offset).Return(mockScores, nil)

	// Panggil metode ListScores
	ctx := context.Background()
	resp, err := service.ListScores(ctx, &scoringv1.ListScoresRequest{
		ExamId:    examID,
		PageSize:  limit,
		PageToken: "",
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(mockScores), len(resp.Scores))
	assert.Equal(t, mockScores[0].ID, resp.Scores[0].Id)
	assert.Equal(t, mockScores[0].SessionID, resp.Scores[0].SessionId)
	assert.Equal(t, mockScores[0].StudentID, resp.Scores[0].StudentId)
	assert.Equal(t, mockScores[0].Score, resp.Scores[0].Score)

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
}

// Helper function untuk membuat ScoringService
func NewScoringService(repo repository.ScoringRepository) scoringv1.ScoringServiceServer {
	return service.NewScoringService(repo)
}
