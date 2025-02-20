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

func (c *ServiceClient) IsExamActive(ctx context.Context, examID string) (bool, error) {
	exam, err := c.GetExam(ctx, examID)
	if err != nil {
		return false, err
	}
	// Check if exam status exists and its state is ACTIVE
	return exam != nil && exam.Status != nil && exam.Status.State == examv1.ExamState_EXAM_STATE_ACTIVE, nil
}

// QuestionService methods
func (c *ServiceClient) GetExamQuestions(ctx context.Context, examID string, randomize bool, limit int32) ([]*questionv1.Question, error) {
	resp, err := c.questionClient.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: randomize,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}
	return resp.Questions, nil
}

// SessionService methods
func (c *ServiceClient) GetSession(ctx context.Context, sessionID string) (*sessionv1.ExamSession, error) {
	return c.sessionClient.GetSession(ctx, &sessionv1.GetSessionRequest{
		Id: sessionID,
	})
}

func (c *ServiceClient) FinishSession(ctx context.Context, sessionID string) error {
	_, err := c.sessionClient.FinishSession(ctx, &sessionv1.FinishSessionRequest{
		Id: sessionID,
	})
	return err
}

func (c *ServiceClient) GetSessionAnswers(ctx context.Context, sessionID string) ([]*sessionv1.Answer, error) {
	session, err := c.GetSession(ctx, sessionID)
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
