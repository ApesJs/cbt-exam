package domain

import (
	"time"
)

type Question struct {
	ID            string    `json:"id"`
	ExamID        string    `json:"exam_id"`
	QuestionText  string    `json:"question_text"`
	Choices       []Choice  `json:"choices"`
	CorrectAnswer string    `json:"correct_answer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Choice struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Untuk mendapatkan soal ujian dengan jumlah dan urutan tertentu
type QuestionFilter struct {
	ExamID    string
	Randomize bool
	Limit     int32
}
