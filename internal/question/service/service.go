package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	questionv1 "cbt-exam/api/proto/question/v1"
	"cbt-exam/internal/question/domain"
	"cbt-exam/internal/question/repository"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type questionService struct {
	repo repository.QuestionRepository
	questionv1.UnimplementedQuestionServiceServer
}

func NewQuestionService(repo repository.QuestionRepository) questionv1.QuestionServiceServer {
	return &questionService{
		repo: repo,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, req *questionv1.CreateQuestionRequest) (*questionv1.Question, error) {
	question := &domain.Question{
		ExamID:        req.ExamId,
		QuestionText:  req.QuestionText,
		CorrectAnswer: req.CorrectAnswer,
	}

	// Convert choices from proto to domain
	for _, c := range req.Choices {
		question.Choices = append(question.Choices, domain.Choice{
			Text: c.Text,
		})
	}

	if err := s.repo.Create(ctx, question); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create question: %v", err)
	}

	return convertDomainToProto(question), nil
}

func (s *questionService) GetQuestion(ctx context.Context, req *questionv1.GetQuestionRequest) (*questionv1.Question, error) {
	question, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrQuestionNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get question: %v", err)
	}

	return convertDomainToProto(question), nil
}

func (s *questionService) ListQuestions(ctx context.Context, req *questionv1.ListQuestionsRequest) (*questionv1.ListQuestionsResponse, error) {
	questions, err := s.repo.List(ctx, req.ExamId, req.PageSize, int32(len(req.PageToken)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list questions: %v", err)
	}

	var protoQuestions []*questionv1.Question
	for _, q := range questions {
		protoQuestions = append(protoQuestions, convertDomainToProto(q))
	}

	// Simple pagination using last question's ID as next page token
	var nextPageToken string
	if len(questions) == int(req.PageSize) {
		nextPageToken = questions[len(questions)-1].ID
	}

	return &questionv1.ListQuestionsResponse{
		Questions:     protoQuestions,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *questionService) UpdateQuestion(ctx context.Context, req *questionv1.UpdateQuestionRequest) (*questionv1.Question, error) {
	question := &domain.Question{
		ID:            req.Id,
		ExamID:        req.Question.ExamId,
		QuestionText:  req.Question.QuestionText,
		CorrectAnswer: req.Question.CorrectAnswer,
	}

	// Convert choices from proto to domain
	for _, c := range req.Question.Choices {
		question.Choices = append(question.Choices, domain.Choice{
			ID:   c.Id,
			Text: c.Text,
		})
	}

	if err := s.repo.Update(ctx, question); err != nil {
		if errors.Is(err, repository.ErrQuestionNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update question: %v", err)
	}

	return convertDomainToProto(question), nil
}

func (s *questionService) DeleteQuestion(ctx context.Context, req *questionv1.DeleteQuestionRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, repository.ErrQuestionNotFound) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete question: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *questionService) GetExamQuestions(ctx context.Context, req *questionv1.GetExamQuestionsRequest) (*questionv1.GetExamQuestionsResponse, error) {
	filter := domain.QuestionFilter{
		ExamID:    req.ExamId,
		Randomize: req.Randomize,
		Limit:     req.Limit,
	}

	questions, err := s.repo.GetExamQuestions(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exam questions: %v", err)
	}

	var protoQuestions []*questionv1.Question
	for _, q := range questions {
		protoQuestions = append(protoQuestions, convertDomainToProto(q))
	}

	return &questionv1.GetExamQuestionsResponse{
		Questions: protoQuestions,
	}, nil
}

// Helper functions to convert between domain and proto models
func convertDomainToProto(q *domain.Question) *questionv1.Question {
	protoQuestion := &questionv1.Question{
		Id:            q.ID,
		ExamId:        q.ExamID,
		QuestionText:  q.QuestionText,
		CorrectAnswer: q.CorrectAnswer,
		//CreatedAt:     timestamppb.New(q.CreatedAt),
	}

	for _, c := range q.Choices {
		protoQuestion.Choices = append(protoQuestion.Choices, &questionv1.Choice{
			Id:   c.ID,
			Text: c.Text,
		})
	}

	return protoQuestion
}
