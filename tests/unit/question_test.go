package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
	"github.com/ApesJs/cbt-exam/internal/question/domain"
	"github.com/ApesJs/cbt-exam/internal/question/repository"
	"github.com/ApesJs/cbt-exam/internal/question/service"
)

type ServiceClientInterface interface {
	GetExam(ctx context.Context, examID string) (*examv1.Exam, error)
	IsExamActive(ctx context.Context, examID string) (bool, error)
}

// MockQuestionRepository adalah struct mock untuk repository.QuestionRepository
type MockQuestionRepository struct {
	mock.Mock
}

// Implementasi metode-metode dari repository.QuestionRepository
func (m *MockQuestionRepository) Create(ctx context.Context, question *domain.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockQuestionRepository) GetByID(ctx context.Context, id string) (*domain.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Question), args.Error(1)
}

func (m *MockQuestionRepository) List(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.Question, error) {
	args := m.Called(ctx, examID, limit, offset)
	return args.Get(0).([]*domain.Question), args.Error(1)
}

func (m *MockQuestionRepository) Update(ctx context.Context, question *domain.Question) error {
	args := m.Called(ctx, question)
	return args.Error(0)
}

func (m *MockQuestionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuestionRepository) GetExamQuestions(ctx context.Context, filter domain.QuestionFilter) ([]*domain.Question, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Question), args.Error(1)
}

func (m *MockQuestionRepository) CountExamQuestions(ctx context.Context, examID string) (int32, error) {
	args := m.Called(ctx, examID)
	return args.Get(0).(int32), args.Error(1)
}

// MockServiceClient adalah mock dari client.ServiceClient
type MockServiceClient struct {
	mock.Mock
}

// Implementasi metode-metode yang dibutuhkan dari client.ServiceClient
func (m *MockServiceClient) GetExam(ctx context.Context, examID string) (*examv1.Exam, error) {
	args := m.Called(ctx, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*examv1.Exam), args.Error(1)
}

func (m *MockServiceClient) IsExamActive(ctx context.Context, examID string) (bool, error) {
	args := m.Called(ctx, examID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockServiceClient) CreateExam(ctx context.Context, req *examv1.CreateExamRequest) (*examv1.Exam, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*examv1.Exam), args.Error(1)
}

func (m *MockServiceClient) ListExams(ctx context.Context, req *examv1.ListExamsRequest) (*examv1.ListExamsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*examv1.ListExamsResponse), args.Error(1)
}

func (m *MockServiceClient) UpdateExam(ctx context.Context, req *examv1.UpdateExamRequest) (*examv1.Exam, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*examv1.Exam), args.Error(1)
}

func (m *MockServiceClient) DeleteExam(ctx context.Context, examID string) error {
	args := m.Called(ctx, examID)
	return args.Error(0)
}

func (m *MockServiceClient) ActivateExam(ctx context.Context, req *examv1.ActivateExamRequest) (*examv1.Exam, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*examv1.Exam), args.Error(1)
}

func (m *MockServiceClient) DeactivateExam(ctx context.Context, req *examv1.DeactivateExamRequest) (*examv1.Exam, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*examv1.Exam), args.Error(1)
}

func (m *MockServiceClient) GetExamQuestions(ctx context.Context, req *questionv1.GetExamQuestionsRequest) (*questionv1.GetExamQuestionsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*questionv1.GetExamQuestionsResponse), args.Error(1)
}

func (m *MockServiceClient) GetSession(ctx context.Context, req *sessionv1.GetSessionRequest) (*sessionv1.ExamSession, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*sessionv1.ExamSession), args.Error(1)
}

func (m *MockServiceClient) StartSession(ctx context.Context, req *sessionv1.StartSessionRequest) (*sessionv1.ExamSession, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*sessionv1.ExamSession), args.Error(1)
}

func (m *MockServiceClient) FinishSession(ctx context.Context, req *sessionv1.FinishSessionRequest) (*sessionv1.ExamSession, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*sessionv1.ExamSession), args.Error(1)
}

func (m *MockServiceClient) SubmitAnswer(ctx context.Context, req *sessionv1.SubmitAnswerRequest) (*sessionv1.SubmitAnswerResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*sessionv1.SubmitAnswerResponse), args.Error(1)
}

func (m *MockServiceClient) GetRemainingTime(ctx context.Context, req *sessionv1.GetRemainingTimeRequest) (*sessionv1.GetRemainingTimeResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*sessionv1.GetRemainingTimeResponse), args.Error(1)
}

func (m *MockServiceClient) GetSessionAnswers(ctx context.Context, sessionID string) ([]*sessionv1.Answer, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]*sessionv1.Answer), args.Error(1)
}

func (m *MockServiceClient) CalculateExamScore(ctx context.Context, sessionID string) (*scoringv1.ExamScore, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*scoringv1.ExamScore), args.Error(1)
}

func (m *MockServiceClient) CreateQuestion(ctx context.Context, req *questionv1.CreateQuestionRequest) (*questionv1.Question, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*questionv1.Question), args.Error(1)
}

func (m *MockServiceClient) GetQuestion(ctx context.Context, req *questionv1.GetQuestionRequest) (*questionv1.Question, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*questionv1.Question), args.Error(1)
}

func (m *MockServiceClient) ListQuestions(ctx context.Context, req *questionv1.ListQuestionsRequest) (*questionv1.ListQuestionsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*questionv1.ListQuestionsResponse), args.Error(1)
}

func (m *MockServiceClient) UpdateQuestion(ctx context.Context, req *questionv1.UpdateQuestionRequest) (*questionv1.Question, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*questionv1.Question), args.Error(1)
}

func (m *MockServiceClient) DeleteQuestion(ctx context.Context, req *questionv1.DeleteQuestionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockServiceClient) CalculateScore(ctx context.Context, req *scoringv1.CalculateScoreRequest) (*scoringv1.ExamScore, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*scoringv1.ExamScore), args.Error(1)
}

func (m *MockServiceClient) GetScore(ctx context.Context, req *scoringv1.GetScoreRequest) (*scoringv1.ExamScore, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*scoringv1.ExamScore), args.Error(1)
}

func (m *MockServiceClient) ListScores(ctx context.Context, req *scoringv1.ListScoresRequest) (*scoringv1.ListScoresResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*scoringv1.ListScoresResponse), args.Error(1)
}

// TestCreateQuestion menguji pembuatan pertanyaan baru
func TestCreateQuestion(t *testing.T) {
	// Setup mock repository dan client
	mockRepo := new(MockQuestionRepository)
	mockClient := new(MockServiceClient)

	// Buat service dengan mock repository dan client
	service := NewQuestionService(mockRepo, mockClient)

	// Setup data ujian yang akan dikembalikan oleh mock client
	examID := "exam-123"
	mockExam := &examv1.Exam{
		Id: examID,
		Status: &examv1.ExamStatus{
			State: examv1.ExamState_EXAM_STATE_CREATED,
		},
	}

	// Setup ekspektasi untuk mock client
	mockClient.On("GetExam", mock.Anything, examID).Return(mockExam, nil)

	// Buat request untuk CreateQuestion
	req := &questionv1.CreateQuestionRequest{
		ExamId:       examID,
		QuestionText: "Berapakah hasil dari 2 + 2?",
		Choices: []*questionv1.Choice{
			{Text: "3"},
			{Text: "4"},
			{Text: "5"},
			{Text: "6"},
		},
		CorrectAnswer: "B", // Jawaban B adalah "4"
	}

	// Setup ekspektasi untuk mock repository
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(question *domain.Question) bool {
		return question.ExamID == req.ExamId &&
			question.QuestionText == req.QuestionText &&
			question.CorrectAnswer == req.CorrectAnswer &&
			len(question.Choices) == len(req.Choices)
	})).Return(nil).Run(func(args mock.Arguments) {
		question := args.Get(1).(*domain.Question)
		question.ID = "question-123"
		question.CreatedAt = time.Now()
		question.UpdatedAt = time.Now()

		// Set ID untuk choices
		for i := range question.Choices {
			question.Choices[i].ID = "choice-" + string('A'+rune(i))
		}
	})

	// Panggil metode CreateQuestion
	ctx := context.Background()
	resp, err := service.CreateQuestion(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "question-123", resp.Id)
	assert.Equal(t, req.ExamId, resp.ExamId)
	assert.Equal(t, req.QuestionText, resp.QuestionText)
	assert.Equal(t, req.CorrectAnswer, resp.CorrectAnswer)
	assert.Equal(t, len(req.Choices), len(resp.Choices))

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

// TestCreateQuestionExamActive menguji pembuatan pertanyaan untuk ujian yang sudah aktif
func TestCreateQuestionExamActive(t *testing.T) {
	// Setup mock repository dan client
	mockRepo := new(MockQuestionRepository)
	mockClient := new(MockServiceClient)

	// Buat service dengan mock repository dan client
	service := NewQuestionService(mockRepo, mockClient)

	// Setup data ujian yang akan dikembalikan oleh mock client
	examID := "exam-123"
	mockExam := &examv1.Exam{
		Id: examID,
		Status: &examv1.ExamStatus{
			State: examv1.ExamState_EXAM_STATE_ACTIVE,
		},
	}

	// Setup ekspektasi untuk mock client
	mockClient.On("GetExam", mock.Anything, examID).Return(mockExam, nil)

	// Buat request untuk CreateQuestion
	req := &questionv1.CreateQuestionRequest{
		ExamId:       examID,
		QuestionText: "Berapakah hasil dari 2 + 2?",
		Choices: []*questionv1.Choice{
			{Text: "3"},
			{Text: "4"},
			{Text: "5"},
			{Text: "6"},
		},
		CorrectAnswer: "B", // Jawaban B adalah "4"
	}

	// Panggil metode CreateQuestion
	ctx := context.Background()
	resp, err := service.CreateQuestion(ctx, req)

	// Assertions - karena ujian sudah aktif, harusnya error
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockClient.AssertExpectations(t)
}

// TestGetExamQuestions menguji pengambilan pertanyaan untuk ujian
func TestGetExamQuestions(t *testing.T) {
	// Setup mock repository dan client
	mockRepo := new(MockQuestionRepository)
	mockClient := new(MockServiceClient)

	// Buat service dengan mock repository dan client
	service := NewQuestionService(mockRepo, mockClient)

	// Setup data
	examID := "exam-123"
	questionFilter := domain.QuestionFilter{
		ExamID:    examID,
		Randomize: true,
		Limit:     10,
	}

	// Setup mock questions yang akan dikembalikan
	mockQuestions := []*domain.Question{
		{
			ID:            "question-1",
			ExamID:        examID,
			QuestionText:  "Pertanyaan 1",
			CorrectAnswer: "A",
			Choices: []domain.Choice{
				{ID: "choice-A", Text: "Pilihan A"},
				{ID: "choice-B", Text: "Pilihan B"},
				{ID: "choice-C", Text: "Pilihan C"},
				{ID: "choice-D", Text: "Pilihan D"},
			},
		},
		{
			ID:            "question-2",
			ExamID:        examID,
			QuestionText:  "Pertanyaan 2",
			CorrectAnswer: "B",
			Choices: []domain.Choice{
				{ID: "choice-A", Text: "Pilihan A"},
				{ID: "choice-B", Text: "Pilihan B"},
				{ID: "choice-C", Text: "Pilihan C"},
				{ID: "choice-D", Text: "Pilihan D"},
			},
		},
	}

	// Setup ekspektasi untuk mock
	mockClient.On("IsExamActive", mock.Anything, examID).Return(true, nil)
	mockRepo.On("GetExamQuestions", mock.Anything, questionFilter).Return(mockQuestions, nil)

	// Panggil metode GetExamQuestions
	ctx := context.Background()
	resp, err := service.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: true,
		Limit:     10,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(mockQuestions), len(resp.Questions))
	assert.Equal(t, mockQuestions[0].ID, resp.Questions[0].Id)
	assert.Equal(t, mockQuestions[0].QuestionText, resp.Questions[0].QuestionText)
	assert.Equal(t, mockQuestions[0].CorrectAnswer, resp.Questions[0].CorrectAnswer)
	assert.Equal(t, len(mockQuestions[0].Choices), len(resp.Questions[0].Choices))

	// Verifikasi bahwa ekspektasi terpenuhi
	mockRepo.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

// TestGetExamQuestionsExamNotActive menguji pengambilan pertanyaan untuk ujian yang tidak aktif
func TestGetExamQuestionsExamNotActive(t *testing.T) {
	// Setup mock repository dan client
	mockRepo := new(MockQuestionRepository)
	mockClient := new(MockServiceClient)

	// Buat service dengan mock repository dan client
	service := NewQuestionService(mockRepo, mockClient)

	// Setup data
	examID := "exam-123"

	// Setup ekspektasi untuk mock
	mockClient.On("IsExamActive", mock.Anything, examID).Return(false, nil)

	// Panggil metode GetExamQuestions
	ctx := context.Background()
	resp, err := service.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: true,
		Limit:     10,
	})

	// Assertions - karena ujian tidak aktif, harusnya error
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Verifikasi bahwa ekspektasi terpenuhi
	mockClient.AssertExpectations(t)
}

// Helper function untuk membuat QuestionService
func NewQuestionService(repo repository.QuestionRepository, client ServiceClientInterface) questionv1.QuestionServiceServer {
	return service.NewQuestionService(repo, client.(interface{}))
}
