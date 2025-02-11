package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
	"github.com/ApesJs/cbt-exam/internal/session/domain"
	"github.com/ApesJs/cbt-exam/internal/session/repository"
)

type sessionService struct {
	repo repository.SessionRepository
	sessionv1.UnimplementedSessionServiceServer
}

func NewSessionService(repo repository.SessionRepository) sessionv1.SessionServiceServer {
	return &sessionService{
		repo: repo,
	}
}

func (s *sessionService) StartSession(ctx context.Context, req *sessionv1.StartSessionRequest) (*sessionv1.ExamSession, error) {
	// Validasi ujian aktif
	isActive, err := s.repo.IsExamActive(ctx, req.ExamId)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to check exam status: %v", err)
	}

	if !isActive {
		return nil, status.Error(codes.FailedPrecondition, "exam is not active")
	}

	// Validasi siswa tidak memiliki sesi aktif lain
	hasActive, err := s.repo.HasActiveSession(ctx, req.StudentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check active sessions: %v", err)
	}

	if hasActive {
		return nil, status.Error(codes.FailedPrecondition, "student already has an active session")
	}

	// Membuat sesi baru
	session := &domain.ExamSession{
		ExamID:    req.ExamId,
		StudentID: req.StudentId,
		Status:    domain.SessionStatusStarted,
		StartTime: time.Now(),
	}

	if err := s.repo.StartSession(ctx, session); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start session: %v", err)
	}

	return convertDomainToProto(session), nil
}

func (s *sessionService) GetSession(ctx context.Context, req *sessionv1.GetSessionRequest) (*sessionv1.ExamSession, error) {
	session, err := s.repo.GetSession(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	return convertDomainToProto(session), nil
}

func (s *sessionService) SubmitAnswer(ctx context.Context, req *sessionv1.SubmitAnswerRequest) (*sessionv1.SubmitAnswerResponse, error) {
	answer := domain.Answer{
		QuestionID:     req.QuestionId,
		SelectedChoice: req.SelectedChoice,
		AnsweredAt:     time.Now(),
	}

	err := s.repo.SubmitAnswer(ctx, req.SessionId, answer)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrSessionNotFound):
			return nil, status.Error(codes.NotFound, "session not found")
		case errors.Is(err, repository.ErrInvalidSessionState):
			return nil, status.Error(codes.FailedPrecondition, "session is not in valid state for answering")
		default:
			return nil, status.Errorf(codes.Internal, "failed to submit answer: %v", err)
		}
	}

	return &sessionv1.SubmitAnswerResponse{
		Success: true,
		Message: "Answer submitted successfully",
	}, nil
}

func (s *sessionService) FinishSession(ctx context.Context, req *sessionv1.FinishSessionRequest) (*sessionv1.ExamSession, error) {
	session, err := s.repo.GetSession(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	if session.Status == domain.SessionStatusFinished {
		return nil, status.Error(codes.FailedPrecondition, "session is already finished")
	}

	err = s.repo.FinishSession(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to finish session: %v", err)
	}

	session.Status = domain.SessionStatusFinished
	session.EndTime = time.Now()

	return convertDomainToProto(session), nil
}

func (s *sessionService) GetRemainingTime(ctx context.Context, req *sessionv1.GetRemainingTimeRequest) (*sessionv1.GetRemainingTimeResponse, error) {
	session, err := s.repo.GetSession(ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	// Mendapatkan durasi ujian dari ExamService (bisa menggunakan gRPC client)
	// Untuk sementara kita hardcode sebagai contoh
	const examDurationMinutes = 120

	remaining := session.CalculateRemainingTime(examDurationMinutes)

	return &sessionv1.GetRemainingTimeResponse{
		RemainingMinutes: remaining.Minutes,
		RemainingSeconds: remaining.Seconds,
	}, nil
}

// Helper function untuk konversi domain ke proto
func convertDomainToProto(session *domain.ExamSession) *sessionv1.ExamSession {
	protoSession := &sessionv1.ExamSession{
		Id:        session.ID,
		ExamId:    session.ExamID,
		StudentId: session.StudentID,
		Status:    convertStatusToProto(session.Status),
		StartTime: timestamppb.New(session.StartTime),
	}

	if !session.EndTime.IsZero() {
		protoSession.EndTime = timestamppb.New(session.EndTime)
	}

	for _, answer := range session.Answers {
		protoSession.Answers = append(protoSession.Answers, &sessionv1.Answer{
			QuestionId:     answer.QuestionID,
			SelectedChoice: answer.SelectedChoice,
			AnsweredAt:     timestamppb.New(answer.AnsweredAt),
		})
	}

	return protoSession
}

func convertStatusToProto(status domain.SessionStatus) sessionv1.SessionStatus {
	switch status {
	case domain.SessionStatusStarted:
		return sessionv1.SessionStatus_SESSION_STATUS_STARTED
	case domain.SessionStatusInProgress:
		return sessionv1.SessionStatus_SESSION_STATUS_IN_PROGRESS
	case domain.SessionStatusFinished:
		return sessionv1.SessionStatus_SESSION_STATUS_FINISHED
	case domain.SessionStatusTimeout:
		return sessionv1.SessionStatus_SESSION_STATUS_TIMEOUT
	default:
		return sessionv1.SessionStatus_SESSION_STATUS_UNSPECIFIED
	}
}
