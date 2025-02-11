package domain

import (
	"time"
)

type ExamState string
type ExamStudentState string

const (
	ExamStateCreated  ExamState = "CREATED"
	ExamStateActive   ExamState = "ACTIVE"
	ExamStateFinished ExamState = "FINISHED"

	ExamStudentStateNotStarted ExamStudentState = "NOT_STARTED"
	ExamStudentStateInProgress ExamStudentState = "IN_PROGRESS"
	ExamStudentStateFinished   ExamStudentState = "FINISHED"
)

type Exam struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Subject        string    `json:"subject"`
	DurationMins   int32     `json:"duration_minutes"`
	TotalQuestions int32     `json:"total_questions"`
	IsRandom       bool      `json:"is_random"`
	TeacherID      string    `json:"teacher_id"`
	ClassIDs       []string  `json:"class_ids"`
	Status         ExamState `json:"status"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StudentStatus struct {
	StudentID   string           `json:"student_id"`
	StudentName string           `json:"student_name"`
	ClassID     string           `json:"class_id"`
	State       ExamStudentState `json:"state"`
	StartTime   time.Time        `json:"start_time"`
	EndTime     time.Time        `json:"end_time"`
}

type ExamStatus struct {
	ID               string          `json:"id"`
	ExamID           string          `json:"exam_id"`
	State            ExamState       `json:"state"`
	TotalStudents    int32           `json:"total_students"`
	StudentsStarted  int32           `json:"students_started"`
	StudentsFinished int32           `json:"students_finished"`
	StudentStatuses  []StudentStatus `json:"student_statuses"`
}
