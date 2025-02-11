package domain

import (
	"time"
)

type ExamScore struct {
	ID              string    `json:"id"`
	ExamID          string    `json:"exam_id"`
	SessionID       string    `json:"session_id"`
	StudentID       string    `json:"student_id"`
	TotalQuestions  int32     `json:"total_questions"`
	CorrectAnswers  int32     `json:"correct_answers"`
	WrongAnswers    int32     `json:"wrong_answers"`
	UnansweredCount int32     `json:"unanswered"`
	Score           float32   `json:"score"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CalculateScore menghitung nilai berdasarkan jawaban benar
func (s *ExamScore) CalculateScore() {
	if s.TotalQuestions > 0 {
		s.Score = float32(s.CorrectAnswers) / float32(s.TotalQuestions) * 100
	}
	s.WrongAnswers = s.TotalQuestions - s.CorrectAnswers - s.UnansweredCount
}

type Answer struct {
	QuestionID    string `json:"question_id"`
	CorrectAnswer string `json:"correct_answer"`
	StudentAnswer string `json:"student_answer"`
}
