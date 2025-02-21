package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
)

type ServiceClient struct {
	examClient     examv1.ExamServiceClient
	questionClient questionv1.QuestionServiceClient
	sessionClient  sessionv1.SessionServiceClient
	scoringClient  scoringv1.ScoringServiceClient
}

func NewServiceClient(examPort, questionPort, sessionPort, scoringPort int) (*ServiceClient, error) {
	// Connect to ExamService
	examConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", examPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to exam service: %v", err)
	}

	// Connect to QuestionService
	questionConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", questionPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to question service: %v", err)
	}

	// Connect to SessionService
	sessionConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", sessionPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session service: %v", err)
	}

	// Connect to ScoringService
	scoringConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", scoringPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to scoring service: %v", err)
	}

	return &ServiceClient{
		examClient:     examv1.NewExamServiceClient(examConn),
		questionClient: questionv1.NewQuestionServiceClient(questionConn),
		sessionClient:  sessionv1.NewSessionServiceClient(sessionConn),
		scoringClient:  scoringv1.NewScoringServiceClient(scoringConn),
	}, nil
}

// ExamService methods
func (c *ServiceClient) GetExam(ctx context.Context, examID string) (*examv1.Exam, error) {
	return c.examClient.GetExam(ctx, &examv1.GetExamRequest{
		Id: examID,
	})
}

func (c *ServiceClient) CreateExam(ctx context.Context, req *examv1.CreateExamRequest) (*examv1.Exam, error) {
	return c.examClient.CreateExam(ctx, req)
}

func (c *ServiceClient) ListExams(ctx context.Context, req *examv1.ListExamsRequest) (*examv1.ListExamsResponse, error) {
	return c.examClient.ListExams(ctx, req)
}

func (c *ServiceClient) UpdateExam(ctx context.Context, req *examv1.UpdateExamRequest) (*examv1.Exam, error) {
	return c.examClient.UpdateExam(ctx, req)
}

func (c *ServiceClient) DeleteExam(ctx context.Context, examID string) error {
	_, err := c.examClient.DeleteExam(ctx, &examv1.DeleteExamRequest{
		Id: examID,
	})
	return err
}

func (c *ServiceClient) ActivateExam(ctx context.Context, req *examv1.ActivateExamRequest) (*examv1.Exam, error) {
	return c.examClient.ActivateExam(ctx, req)
}

func (c *ServiceClient) DeactivateExam(ctx context.Context, req *examv1.DeactivateExamRequest) (*examv1.Exam, error) {
	return c.examClient.DeactivateExam(ctx, req)
}

func (c *ServiceClient) IsExamActive(ctx context.Context, examID string) (bool, error) {
	exam, err := c.GetExam(ctx, examID)
	if err != nil {
		return false, err
	}
	// Check if exam status exists and its state is ACTIVE
	return exam != nil && exam.Status != nil && exam.Status.State == examv1.ExamState_EXAM_STATE_ACTIVE, nil
}

// QuestionService methods
func (c *ServiceClient) GetExamQuestions(ctx context.Context, req *questionv1.GetExamQuestionsRequest) (*questionv1.GetExamQuestionsResponse, error) {
	return c.questionClient.GetExamQuestions(ctx, req)
}

func (c *ServiceClient) GetSession(ctx context.Context, req *sessionv1.GetSessionRequest) (*sessionv1.ExamSession, error) {
	return c.sessionClient.GetSession(ctx, req)
}

func (c *ServiceClient) FinishSession(ctx context.Context, req *sessionv1.FinishSessionRequest) (*sessionv1.ExamSession, error) {
	return c.sessionClient.FinishSession(ctx, req)
}

func (c *ServiceClient) StartSession(ctx context.Context, req *sessionv1.StartSessionRequest) (*sessionv1.ExamSession, error) {
	return c.sessionClient.StartSession(ctx, req)
}

func (c *ServiceClient) SubmitAnswer(ctx context.Context, req *sessionv1.SubmitAnswerRequest) (*sessionv1.SubmitAnswerResponse, error) {
	return c.sessionClient.SubmitAnswer(ctx, req)
}

func (c *ServiceClient) GetRemainingTime(ctx context.Context, req *sessionv1.GetRemainingTimeRequest) (*sessionv1.GetRemainingTimeResponse, error) {
	return c.sessionClient.GetRemainingTime(ctx, req)
}

func (c *ServiceClient) GetSessionAnswers(ctx context.Context, sessionID string) ([]*sessionv1.Answer, error) {
	session, err := c.GetSession(ctx, &sessionv1.GetSessionRequest{
		Id: sessionID,
	})
	if err != nil {
		return nil, err
	}
	return session.Answers, nil
}

// ScoringService methods
func (c *ServiceClient) CalculateExamScore(ctx context.Context, sessionID string) (*scoringv1.ExamScore, error) {
	return c.scoringClient.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})
}

// QuestionService methods
func (c *ServiceClient) CreateQuestion(ctx context.Context, req *questionv1.CreateQuestionRequest) (*questionv1.Question, error) {
	return c.questionClient.CreateQuestion(ctx, req)
}

func (c *ServiceClient) GetQuestion(ctx context.Context, req *questionv1.GetQuestionRequest) (*questionv1.Question, error) {
	return c.questionClient.GetQuestion(ctx, req)
}

func (c *ServiceClient) ListQuestions(ctx context.Context, req *questionv1.ListQuestionsRequest) (*questionv1.ListQuestionsResponse, error) {
	return c.questionClient.ListQuestions(ctx, req)
}

func (c *ServiceClient) UpdateQuestion(ctx context.Context, req *questionv1.UpdateQuestionRequest) (*questionv1.Question, error) {
	return c.questionClient.UpdateQuestion(ctx, req)
}

func (c *ServiceClient) DeleteQuestion(ctx context.Context, req *questionv1.DeleteQuestionRequest) error {
	_, err := c.questionClient.DeleteQuestion(ctx, req)
	return err
}

// ScoringService methods
func (c *ServiceClient) CalculateScore(ctx context.Context, req *scoringv1.CalculateScoreRequest) (*scoringv1.ExamScore, error) {
	return c.scoringClient.CalculateScore(ctx, req)
}

func (c *ServiceClient) GetScore(ctx context.Context, req *scoringv1.GetScoreRequest) (*scoringv1.ExamScore, error) {
	return c.scoringClient.GetScore(ctx, req)
}

func (c *ServiceClient) ListScores(ctx context.Context, req *scoringv1.ListScoresRequest) (*scoringv1.ListScoresResponse, error) {
	return c.scoringClient.ListScores(ctx, req)
}
