package service

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	"github.com/ApesJs/cbt-exam/internal/scoring/domain"
	"github.com/ApesJs/cbt-exam/internal/scoring/repository"
)

type scoringService struct {
	repo repository.ScoringRepository
	scoringv1.UnimplementedScoringServiceServer
}

func NewScoringService(repo repository.ScoringRepository) scoringv1.ScoringServiceServer {
	return &scoringService{
		repo: repo,
	}
}

func (s *scoringService) CalculateScore(ctx context.Context, req *scoringv1.CalculateScoreRequest) (*scoringv1.ExamScore, error) {
	// Get all answers from the session (both correct answers and student answers)
	answers, err := s.repo.GetCorrectAnswers(ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get answers: %v", err)
	}

	// Calculate score
	score := &domain.ExamScore{
		SessionID:      req.SessionId,
		TotalQuestions: int32(len(answers)),
	}

	// Count correct answers and unanswered questions
	for _, answer := range answers {
		if answer.StudentAnswer == "" {
			score.UnansweredCount++
		} else if answer.StudentAnswer == answer.CorrectAnswer {
			score.CorrectAnswers++
		}
	}

	// Calculate score percentage
	score.CalculateScore()

	// Save the score
	if err := s.repo.CreateScore(ctx, score); err != nil {
		if errors.Is(err, repository.ErrDuplicateScore) {
			return nil, status.Error(codes.AlreadyExists, "score already exists for this exam and student")
		}
		return nil, status.Errorf(codes.Internal, "failed to save score: %v", err)
	}

	return convertDomainToProto(score), nil
}

func (s *scoringService) GetScore(ctx context.Context, req *scoringv1.GetScoreRequest) (*scoringv1.ExamScore, error) {
	score, err := s.repo.GetScore(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrScoreNotFound) {
			return nil, status.Error(codes.NotFound, "score not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get score: %v", err)
	}

	return convertDomainToProto(score), nil
}

func (s *scoringService) ListScores(ctx context.Context, req *scoringv1.ListScoresRequest) (*scoringv1.ListScoresResponse, error) {
	scores, err := s.repo.ListScores(ctx, req.ExamId, req.PageSize, int32(len(req.PageToken)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list scores: %v", err)
	}

	var protoScores []*scoringv1.ExamScore
	for _, score := range scores {
		protoScores = append(protoScores, convertDomainToProto(score))
	}

	// Simple pagination using last score's ID as next page token
	var nextPageToken string
	if len(scores) == int(req.PageSize) {
		nextPageToken = scores[len(scores)-1].ID
	}

	return &scoringv1.ListScoresResponse{
		Scores:        protoScores,
		NextPageToken: nextPageToken,
	}, nil
}

// Helper function untuk konversi domain ke proto
func convertDomainToProto(score *domain.ExamScore) *scoringv1.ExamScore {
	return &scoringv1.ExamScore{
		Id:             score.ID,
		ExamId:         score.ExamID,
		SessionId:      score.SessionID,
		StudentId:      score.StudentID,
		TotalQuestions: score.TotalQuestions,
		CorrectAnswers: score.CorrectAnswers,
		WrongAnswers:   score.WrongAnswers,
		Unanswered:     score.UnansweredCount,
		Score:          score.Score,
		CreatedAt:      timestamppb.New(score.CreatedAt),
	}
}
