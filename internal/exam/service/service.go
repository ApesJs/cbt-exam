package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	examv1 "cbt-exam/api/proto/exam/v1"
	"cbt-exam/internal/exam/domain"
	"cbt-exam/internal/exam/repository"
)

type examService struct {
	repo repository.ExamRepository
	examv1.UnimplementedExamServiceServer
}

func NewExamService(repo repository.ExamRepository) examv1.ExamServiceServer {
	return &examService{
		repo: repo,
	}
}

func (s *examService) CreateExam(ctx context.Context, req *examv1.CreateExamRequest) (*examv1.Exam, error) {
	exam := &domain.Exam{
		Title:          req.Title,
		Subject:        req.Subject,
		DurationMins:   req.DurationMinutes,
		TotalQuestions: req.TotalQuestions,
		IsRandom:       req.IsRandom,
		TeacherID:      req.TeacherId,
		ClassIDs:       req.ClassIds,
		Status:         domain.ExamStateCreated,
	}

	if err := s.repo.Create(ctx, exam); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create exam: %v", err)
	}

	return convertDomainToProto(exam), nil
}

func (s *examService) GetExam(ctx context.Context, req *examv1.GetExamRequest) (*examv1.Exam, error) {
	exam, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get exam: %v", err)
	}

	return convertDomainToProto(exam), nil
}

func (s *examService) ListExams(ctx context.Context, req *examv1.ListExamsRequest) (*examv1.ListExamsResponse, error) {
	exams, err := s.repo.List(ctx, req.TeacherId, req.PageSize, int32(len(req.PageToken)))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list exams: %v", err)
	}

	var protoExams []*examv1.Exam
	for _, exam := range exams {
		protoExams = append(protoExams, convertDomainToProto(exam))
	}

	// Simple pagination - using the last exam's ID as the next page token
	var nextPageToken string
	if len(exams) == int(req.PageSize) {
		nextPageToken = exams[len(exams)-1].ID
	}

	return &examv1.ListExamsResponse{
		Exams:         protoExams,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *examService) UpdateExam(ctx context.Context, req *examv1.UpdateExamRequest) (*examv1.Exam, error) {
	exam := &domain.Exam{
		ID:             req.Id,
		Title:          req.Exam.Title,
		Subject:        req.Exam.Subject,
		DurationMins:   req.Exam.DurationMinutes,
		TotalQuestions: req.Exam.TotalQuestions,
		IsRandom:       req.Exam.IsRandom,
		TeacherID:      req.Exam.TeacherId,
		ClassIDs:       req.Exam.ClassIds,
	}

	if err := s.repo.Update(ctx, exam); err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update exam: %v", err)
	}

	return convertDomainToProto(exam), nil
}

func (s *examService) DeleteExam(ctx context.Context, req *examv1.DeleteExamRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete exam: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *examService) ActivateExam(ctx context.Context, req *examv1.ActivateExamRequest) (*examv1.Exam, error) {
	exam, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get exam: %v", err)
	}

	if exam.Status != domain.ExamStateCreated {
		return nil, status.Error(codes.FailedPrecondition, "exam can only be activated when in CREATED state")
	}

	exam.Status = domain.ExamStateActive
	exam.StartTime = time.Now()

	if err := s.repo.UpdateStatus(ctx, req.Id, domain.ExamStateActive); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to activate exam: %v", err)
	}

	return convertDomainToProto(exam), nil
}

func (s *examService) DeactivateExam(ctx context.Context, req *examv1.DeactivateExamRequest) (*examv1.Exam, error) {
	exam, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get exam: %v", err)
	}

	if exam.Status != domain.ExamStateActive {
		return nil, status.Error(codes.FailedPrecondition, "exam can only be deactivated when in ACTIVE state")
	}

	exam.Status = domain.ExamStateFinished
	exam.EndTime = time.Now()

	if err := s.repo.UpdateStatus(ctx, req.Id, domain.ExamStateFinished); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to deactivate exam: %v", err)
	}

	return convertDomainToProto(exam), nil
}

func (s *examService) GetExamStatus(ctx context.Context, req *examv1.GetExamStatusRequest) (*examv1.ExamStatus, error) {
	getStatus, err := s.repo.GetStatus(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrExamNotFound) {
			return nil, status.Error(codes.NotFound, "exam not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get exam status: %v", err)
	}

	return convertStatusToProto(getStatus), nil
}

// Helper functions to convert between domain and proto models
func convertDomainToProto(exam *domain.Exam) *examv1.Exam {
	return &examv1.Exam{
		Id:              exam.ID,
		Title:           exam.Title,
		Subject:         exam.Subject,
		DurationMinutes: exam.DurationMins,
		TotalQuestions:  exam.TotalQuestions,
		IsRandom:        exam.IsRandom,
		TeacherId:       exam.TeacherID,
		ClassIds:        exam.ClassIDs,
		Status: &examv1.ExamStatus{
			State: convertExamState(exam.Status),
		},
		StartTime: timestamppb.New(exam.StartTime),
		EndTime:   timestamppb.New(exam.EndTime),
		CreatedAt: timestamppb.New(exam.CreatedAt),
		UpdatedAt: timestamppb.New(exam.UpdatedAt),
	}
}

func convertStatusToProto(status *domain.ExamStatus) *examv1.ExamStatus {
	var studentStatuses []*examv1.StudentStatus
	for _, s := range status.StudentStatuses {
		studentStatuses = append(studentStatuses, &examv1.StudentStatus{
			StudentId:   s.StudentID,
			StudentName: s.StudentName,
			ClassId:     s.ClassID,
			State:       convertStudentState(s.State),
			StartTime:   timestamppb.New(s.StartTime),
			EndTime:     timestamppb.New(s.EndTime),
		})
	}

	return &examv1.ExamStatus{
		Id:               status.ID,
		ExamId:           status.ExamID,
		State:            convertExamState(status.State),
		TotalStudents:    status.TotalStudents,
		StudentsStarted:  status.StudentsStarted,
		StudentsFinished: status.StudentsFinished,
		StudentStatuses:  studentStatuses,
	}
}

func convertExamState(state domain.ExamState) examv1.ExamState {
	switch state {
	case domain.ExamStateCreated:
		return examv1.ExamState_EXAM_STATE_CREATED
	case domain.ExamStateActive:
		return examv1.ExamState_EXAM_STATE_ACTIVE
	case domain.ExamStateFinished:
		return examv1.ExamState_EXAM_STATE_FINISHED
	default:
		return examv1.ExamState_EXAM_STATE_UNSPECIFIED
	}
}

func convertStudentState(state domain.ExamStudentState) examv1.ExamStudentState {
	switch state {
	case domain.ExamStudentStateNotStarted:
		return examv1.ExamStudentState_EXAM_STUDENT_STATE_NOT_STARTED
	case domain.ExamStudentStateInProgress:
		return examv1.ExamStudentState_EXAM_STUDENT_STATE_IN_PROGRESS
	case domain.ExamStudentStateFinished:
		return examv1.ExamStudentState_EXAM_STUDENT_STATE_FINISHED
	default:
		return examv1.ExamStudentState_EXAM_STUDENT_STATE_UNSPECIFIED
	}
}
