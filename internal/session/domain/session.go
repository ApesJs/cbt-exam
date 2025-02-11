package domain

import (
	"time"
)

type SessionStatus string

const (
	SessionStatusStarted    SessionStatus = "STARTED"
	SessionStatusInProgress SessionStatus = "IN_PROGRESS"
	SessionStatusFinished   SessionStatus = "FINISHED"
	SessionStatusTimeout    SessionStatus = "TIMEOUT"
)

type ExamSession struct {
	ID        string        `json:"id"`
	ExamID    string        `json:"exam_id"`
	StudentID string        `json:"student_id"`
	Status    SessionStatus `json:"status"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Answers   []Answer      `json:"answers"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Answer struct {
	QuestionID     string    `json:"question_id"`
	SelectedChoice string    `json:"selected_choice"`
	AnsweredAt     time.Time `json:"answered_at"`
}

type RemainingTime struct {
	Minutes int32 `json:"remaining_minutes"`
	Seconds int32 `json:"remaining_seconds"`
}

// Helper function untuk menghitung sisa waktu
func (s *ExamSession) CalculateRemainingTime(durationMinutes int32) *RemainingTime {
	if s.Status == SessionStatusFinished || s.Status == SessionStatusTimeout {
		return &RemainingTime{
			Minutes: 0,
			Seconds: 0,
		}
	}

	endTime := s.StartTime.Add(time.Duration(durationMinutes) * time.Minute)
	remaining := time.Until(endTime)

	if remaining < 0 {
		return &RemainingTime{
			Minutes: 0,
			Seconds: 0,
		}
	}

	minutes := int32(remaining.Minutes())
	seconds := int32(remaining.Seconds()) % 60

	return &RemainingTime{
		Minutes: minutes,
		Seconds: seconds,
	}
}
