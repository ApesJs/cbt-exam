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
	examDomain "github.com/ApesJs/cbt-exam/internal/exam/domain"
	examRepo "github.com/ApesJs/cbt-exam/internal/exam/repository"
	questionDomain "github.com/ApesJs/cbt-exam/internal/question/domain"
	scoringDomain "github.com/ApesJs/cbt-exam/internal/scoring/domain"
	scoringRepo "github.com/ApesJs/cbt-exam/internal/scoring/repository"
	sessionDomain "github.com/ApesJs/cbt-exam/internal/session/domain"
)

// IntegrationTestSuite adalah struct yang berisi semua mock repositories dan service instances
// yang diperlukan untuk tes integrasi
type IntegrationTestSuite struct {
	MockExamRepo     *MockExamRepository
	MockQuestionRepo *MockQuestionRepository
	MockSessionRepo  *MockSessionRepository
	MockScoringRepo  *MockScoringRepository
	MockClient       *MockServiceClient

	ExamService     examv1.ExamServiceServer
	QuestionService questionv1.QuestionServiceServer
	SessionService  sessionv1.SessionServiceServer
	ScoringService  scoringv1.ScoringServiceServer
}

// SetupIntegrationTest membuat dan menginisialisasi IntegrationTestSuite
func SetupIntegrationTest() *IntegrationTestSuite {
	suite := &IntegrationTestSuite{
		MockExamRepo:     new(MockExamRepository),
		MockQuestionRepo: new(MockQuestionRepository),
		MockSessionRepo:  new(MockSessionRepository),
		MockScoringRepo:  new(MockScoringRepository),
		MockClient:       new(MockServiceClient),
	}

	suite.ExamService = NewExamService(suite.MockExamRepo)
	suite.QuestionService = NewQuestionService(suite.MockQuestionRepo, suite.MockClient)
	suite.SessionService = NewSessionService(suite.MockSessionRepo)
	suite.ScoringService = NewScoringService(suite.MockScoringRepo)

	return suite
}

// TestFullExamWorkflow menguji alur kerja lengkap dari ujian, mulai dari pembuatan ujian
// hingga penghitungan skor
func TestFullExamWorkflow(t *testing.T) {
	// Setup integration test suite
	suite := SetupIntegrationTest()

	// Context untuk semua operasi
	ctx := context.Background()

	// Step 1: Buat ujian baru
	teacherID := "teacher-123"
	classIDs := []string{"class-1", "class-2"}
	examTitle := "Ujian Matematika"
	examSubject := "Matematika"
	durationMinutes := int32(60)
	totalQuestions := int32(5)

	// Setup ekspektasi untuk mock exam repository
	mockExam := &examv1.Exam{
		Id:              "exam-123",
		Title:           examTitle,
		Subject:         examSubject,
		DurationMinutes: durationMinutes,
		TotalQuestions:  totalQuestions,
		IsRandom:        true,
		TeacherId:       teacherID,
		ClassIds:        classIDs,
		Status: &examv1.ExamStatus{
			State: examv1.ExamState_EXAM_STATE_CREATED,
		},
	}

	suite.MockExamRepo.On("Create", mock.Anything, mock.MatchedBy(func(exam *examDomain.Exam) bool {
		return exam.Title == examTitle &&
			exam.Subject == examSubject &&
			exam.DurationMins == durationMinutes
	})).Return(nil).Run(func(args mock.Arguments) {
		exam := args.Get(1).(*examDomain.Exam)
		exam.ID = "exam-123"
		exam.CreatedAt = time.Now()
		exam.UpdatedAt = time.Now()
	})

	// Panggil CreateExam
	examResp, err := suite.ExamService.CreateExam(ctx, &examv1.CreateExamRequest{
		Title:           examTitle,
		Subject:         examSubject,
		DurationMinutes: durationMinutes,
		TotalQuestions:  totalQuestions,
		IsRandom:        true,
		TeacherId:       teacherID,
		ClassIds:        classIDs,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, examResp)
	assert.Equal(t, "exam-123", examResp.Id)

	// Step 2: Buat pertanyaan untuk ujian
	// Setup ekspektasi untuk mock client dan question repository
	suite.MockClient.On("GetExam", mock.Anything, "exam-123").Return(mockExam, nil)

	for i := 1; i <= 5; i++ {
		questionID := "question-" + string(rune('0'+i))
		questionText := "Pertanyaan " + string(rune('0'+i)) + "?"
		correctAnswer := string(rune('A' + (i-1)%4))

		choices := []*questionv1.Choice{
			{Text: "Jawaban A"},
			{Text: "Jawaban B"},
			{Text: "Jawaban C"},
			{Text: "Jawaban D"},
		}

		suite.MockQuestionRepo.On("Create", mock.Anything, mock.MatchedBy(func(q *questionDomain.Question) bool {
			return q.ExamID == "exam-123" &&
				q.QuestionText == questionText &&
				q.CorrectAnswer == correctAnswer
		})).Return(nil).Run(func(args mock.Arguments) {
			q := args.Get(1).(*questionDomain.Question)
			q.ID = questionID
			q.CreatedAt = time.Now()
			q.UpdatedAt = time.Now()

			// Setup choices
			for j := range q.Choices {
				q.Choices[j].ID = "choice-" + string(rune('A'+j))
			}
		})

		// Panggil CreateQuestion
		questionResp, err := suite.QuestionService.CreateQuestion(ctx, &questionv1.CreateQuestionRequest{
			ExamId:        "exam-123",
			QuestionText:  questionText,
			CorrectAnswer: correctAnswer,
			Choices:       choices,
		})

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, questionResp)
		assert.Equal(t, questionID, questionResp.Id)
	}

	// Step 3: Aktifkan ujian
	suite.MockExamRepo.On("GetByID", mock.Anything, "exam-123").Return(&examDomain.Exam{
		ID:             "exam-123",
		Title:          examTitle,
		Subject:        examSubject,
		DurationMins:   durationMinutes,
		TotalQuestions: totalQuestions,
		IsRandom:       true,
		TeacherID:      teacherID,
		ClassIDs:       classIDs,
		Status:         examDomain.ExamStateCreated,
	}, nil)

	suite.MockExamRepo.On("UpdateStatus", mock.Anything, "exam-123", examDomain.ExamStateActive).Return(nil)

	// Panggil ActivateExam
	activateResp, err := suite.ExamService.ActivateExam(ctx, &examv1.ActivateExamRequest{
		Id:       "exam-123",
		ClassIds: classIDs,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, activateResp)
	assert.Equal(t, "exam-123", activateResp.Id)

	// Step 4: Mulai sesi ujian untuk seorang siswa
	studentID := "student-456"
	sessionID := "session-789"

	suite.MockSessionRepo.On("IsExamActive", mock.Anything, "exam-123").Return(true, nil)
	suite.MockSessionRepo.On("HasActiveSession", mock.Anything, studentID).Return(false, nil)
	suite.MockSessionRepo.On("StartSession", mock.Anything, mock.MatchedBy(func(session *sessionDomain.ExamSession) bool {
		return session.ExamID == "exam-123" &&
			session.StudentID == studentID
	})).Return(nil).Run(func(args mock.Arguments) {
		session := args.Get(1).(*sessionDomain.ExamSession)
		session.ID = sessionID
		session.CreatedAt = time.Now()
		session.UpdatedAt = time.Now()
	})

	// Panggil StartSession
	sessionResp, err := suite.SessionService.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    "exam-123",
		StudentId: studentID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, sessionResp)
	assert.Equal(t, sessionID, sessionResp.Id)

	// Step 5: Ambil pertanyaan ujian
	mockQuestions := []*questionDomain.Question{
		{
			ID:            "question-1",
			ExamID:        "exam-123",
			QuestionText:  "Pertanyaan 1?",
			CorrectAnswer: "A",
			Choices: []questionDomain.Choice{
				{ID: "choice-A", Text: "Jawaban A"},
				{ID: "choice-B", Text: "Jawaban B"},
				{ID: "choice-C", Text: "Jawaban C"},
				{ID: "choice-D", Text: "Jawaban D"},
			},
		},
		{
			ID:            "question-2",
			ExamID:        "exam-123",
			QuestionText:  "Pertanyaan 2?",
			CorrectAnswer: "B",
			Choices: []questionDomain.Choice{
				{ID: "choice-A", Text: "Jawaban A"},
				{ID: "choice-B", Text: "Jawaban B"},
				{ID: "choice-C", Text: "Jawaban C"},
				{ID: "choice-D", Text: "Jawaban D"},
			},
		},
		{
			ID:            "question-3",
			ExamID:        "exam-123",
			QuestionText:  "Pertanyaan 3?",
			CorrectAnswer: "C",
			Choices: []questionDomain.Choice{
				{ID: "choice-A", Text: "Jawaban A"},
				{ID: "choice-B", Text: "Jawaban B"},
				{ID: "choice-C", Text: "Jawaban C"},
				{ID: "choice-D", Text: "Jawaban D"},
			},
		},
		{
			ID:            "question-4",
			ExamID:        "exam-123",
			QuestionText:  "Pertanyaan 4?",
			CorrectAnswer: "D",
			Choices: []questionDomain.Choice{
				{ID: "choice-A", Text: "Jawaban A"},
				{ID: "choice-B", Text: "Jawaban B"},
				{ID: "choice-C", Text: "Jawaban C"},
				{ID: "choice-D", Text: "Jawaban D"},
			},
		},
		{
			ID:            "question-5",
			ExamID:        "exam-123",
			QuestionText:  "Pertanyaan 5?",
			CorrectAnswer: "A",
			Choices: []questionDomain.Choice{
				{ID: "choice-A", Text: "Jawaban A"},
				{ID: "choice-B", Text: "Jawaban B"},
				{ID: "choice-C", Text: "Jawaban C"},
				{ID: "choice-D", Text: "Jawaban D"},
			},
		},
	}

	suite.MockClient.On("IsExamActive", mock.Anything, "exam-123").Return(true, nil)
	suite.MockQuestionRepo.On("GetExamQuestions", mock.Anything, mock.MatchedBy(func(filter questionDomain.QuestionFilter) bool {
		return filter.ExamID == "exam-123" &&
			filter.Randomize == true &&
			filter.Limit == totalQuestions
	})).Return(mockQuestions, nil)

	// Panggil GetExamQuestions
	questionsResp, err := suite.QuestionService.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    "exam-123",
		Randomize: true,
		Limit:     totalQuestions,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, questionsResp)
	assert.Equal(t, len(mockQuestions), len(questionsResp.Questions))

	// Step 6: Kirim jawaban untuk pertanyaan
	suite.MockSessionRepo.On("SubmitAnswer", mock.Anything, sessionID, mock.MatchedBy(func(answer sessionDomain.Answer) bool {
		return answer.QuestionID == "question-1" &&
			answer.SelectedChoice == "A" // Jawaban benar
	})).Return(nil)

	suite.MockSessionRepo.On("SubmitAnswer", mock.Anything, sessionID, mock.MatchedBy(func(answer sessionDomain.Answer) bool {
		return answer.QuestionID == "question-2" &&
			answer.SelectedChoice == "C" // Jawaban salah
	})).Return(nil)

	// Jawaban untuk question-1
	answerResp1, err := suite.SessionService.SubmitAnswer(ctx, &sessionv1.SubmitAnswerRequest{
		SessionId:      sessionID,
		QuestionId:     "question-1",
		SelectedChoice: "A", // Jawaban benar
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, answerResp1)
	assert.True(t, answerResp1.Success)

	// Jawaban untuk question-2
	answerResp2, err := suite.SessionService.SubmitAnswer(ctx, &sessionv1.SubmitAnswerRequest{
		SessionId:      sessionID,
		QuestionId:     "question-2",
		SelectedChoice: "C", // Jawaban salah
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, answerResp2)
	assert.True(t, answerResp2.Success)

	// Step 7: Selesaikan sesi ujian
	mockSession := &sessionDomain.ExamSession{
		ID:        sessionID,
		ExamID:    "exam-123",
		StudentID: studentID,
		Status:    sessionDomain.SessionStatusInProgress,
		StartTime: time.Now().Add(-30 * time.Minute),
		Answers: []sessionDomain.Answer{
			{QuestionID: "question-1", SelectedChoice: "A"},
			{QuestionID: "question-2", SelectedChoice: "C"},
		},
	}

	suite.MockSessionRepo.On("GetSession", mock.Anything, sessionID).Return(mockSession, nil)
	suite.MockSessionRepo.On("FinishSession", mock.Anything, sessionID).Return(nil)

	// Panggil FinishSession
	finishResp, err := suite.SessionService.FinishSession(ctx, &sessionv1.FinishSessionRequest{
		Id: sessionID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, finishResp)
	assert.Equal(t, sessionID, finishResp.Id)
	assert.Equal(t, sessionv1.SessionStatus_SESSION_STATUS_FINISHED, finishResp.Status)

	// Step 8: Hitung skor ujian
	mockAnswers := []scoringDomain.Answer{
		{
			QuestionID:    "question-1",
			CorrectAnswer: "A",
			StudentAnswer: "A", // Benar
		},
		{
			QuestionID:    "question-2",
			CorrectAnswer: "B",
			StudentAnswer: "C", // Salah
		},
		{
			QuestionID:    "question-3",
			CorrectAnswer: "C",
			StudentAnswer: "", // Tidak dijawab
		},
		{
			QuestionID:    "question-4",
			CorrectAnswer: "D",
			StudentAnswer: "", // Tidak dijawab
		},
		{
			QuestionID:    "question-5",
			CorrectAnswer: "A",
			StudentAnswer: "", // Tidak dijawab
		},
	}

	scoreID := "score-123"

	suite.MockScoringRepo.On("GetCorrectAnswers", mock.Anything, sessionID).Return(mockAnswers, nil)
	suite.MockScoringRepo.On("CreateScore", mock.Anything, mock.MatchedBy(func(score *scoringDomain.ExamScore) bool {
		return score.SessionID == sessionID &&
			score.TotalQuestions == int32(len(mockAnswers)) &&
			score.CorrectAnswers == int32(1) && // 1 jawaban benar
			score.WrongAnswers == int32(1) && // 1 jawaban salah
			score.UnansweredCount == int32(3) // 3 tidak dijawab
	})).Return(nil).Run(func(args mock.Arguments) {
		score := args.Get(1).(*scoringDomain.ExamScore)
		score.ID = scoreID
		score.ExamID = "exam-123"
		score.StudentID = studentID
		score.CreatedAt = time.Now()
	})

	// Panggil CalculateScore
	scoreResp, err := suite.ScoringService.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, scoreResp)
	assert.Equal(t, scoreID, scoreResp.Id)
	assert.Equal(t, sessionID, scoreResp.SessionId)
	assert.Equal(t, "exam-123", scoreResp.ExamId)
	assert.Equal(t, studentID, scoreResp.StudentId)
	assert.Equal(t, int32(5), scoreResp.TotalQuestions)
	assert.Equal(t, int32(1), scoreResp.CorrectAnswers)
	assert.Equal(t, int32(1), scoreResp.WrongAnswers)
	assert.Equal(t, int32(3), scoreResp.Unanswered)
	assert.Equal(t, float32(20.0), scoreResp.Score) // 1/5 * 100 = 20%

	// Verifikasi bahwa semua ekspektasi terpenuhi
	suite.MockExamRepo.AssertExpectations(t)
	suite.MockQuestionRepo.AssertExpectations(t)
	suite.MockSessionRepo.AssertExpectations(t)
	suite.MockScoringRepo.AssertExpectations(t)
	suite.MockClient.AssertExpectations(t)
}

// TestExamWorkflowErrors menguji alur kerja ujian dengan berbagai skenario error
func TestExamWorkflowErrors(t *testing.T) {
	// Setup integration test suite
	suite := SetupIntegrationTest()

	// Context untuk semua operasi
	ctx := context.Background()

	// Skenario 1: Mencoba mengaktifkan ujian yang tidak ada
	examID := "non-existent-exam"
	classIDs := []string{"class-1", "class-2"}

	suite.MockExamRepo.On("GetByID", mock.Anything, examID).Return(nil, examRepo.ErrExamNotFound)

	// Panggil ActivateExam
	activateResp, err := suite.ExamService.ActivateExam(ctx, &examv1.ActivateExamRequest{
		Id:       examID,
		ClassIds: classIDs,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, activateResp)

	// Verifikasi status error
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	// Skenario 2: Mencoba memulai sesi ujian untuk ujian yang tidak aktif
	examID = "inactive-exam"
	studentID := "student-123"

	suite.MockSessionRepo.On("IsExamActive", mock.Anything, examID).Return(false, nil)

	// Panggil StartSession
	sessionResp, err := suite.SessionService.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    examID,
		StudentId: studentID,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, sessionResp)

	// Verifikasi status error
	st, ok = status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Skenario 3: Mencoba mengambil pertanyaan untuk ujian yang tidak aktif
	suite.MockClient.On("IsExamActive", mock.Anything, examID).Return(false, nil)

	// Panggil GetExamQuestions
	questionsResp, err := suite.QuestionService.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: true,
		Limit:     10,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, questionsResp)

	// Verifikasi status error
	st, ok = status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	// Skenario 4: Mencoba menghitung skor untuk sesi yang tidak ada
	sessionID := "non-existent-session"

	suite.MockScoringRepo.On("GetCorrectAnswers", mock.Anything, sessionID).Return([]scoringDomain.Answer{}, scoringRepo.ErrSessionNotFound)

	// Panggil CalculateScore
	scoreResp, err := suite.ScoringService.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, scoreResp)

	// Verifikasi status error
	st, ok = status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	// Verifikasi bahwa semua ekspektasi terpenuhi
	suite.MockExamRepo.AssertExpectations(t)
	suite.MockSessionRepo.AssertExpectations(t)
	suite.MockScoringRepo.AssertExpectations(t)
	suite.MockClient.AssertExpectations(t)
}
