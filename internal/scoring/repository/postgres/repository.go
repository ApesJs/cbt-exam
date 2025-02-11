package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/ApesJs/cbt-exam/internal/scoring/domain"
	"github.com/ApesJs/cbt-exam/internal/scoring/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) repository.ScoringRepository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) CreateScore(ctx context.Context, score *domain.ExamScore) error {
	query := `
        INSERT INTO exam_scores (
            exam_id, session_id, student_id, total_questions,
            correct_answers, wrong_answers, unanswered, score
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		score.ExamID,
		score.SessionID,
		score.StudentID,
		score.TotalQuestions,
		score.CorrectAnswers,
		score.WrongAnswers,
		score.UnansweredCount,
		score.Score,
	).Scan(&score.ID, &score.CreatedAt, &score.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return repository.ErrDuplicateScore
			}
		}
		return errors.Wrap(err, "failed to create score")
	}

	return nil
}

func (r *postgresRepository) GetScore(ctx context.Context, id string) (*domain.ExamScore, error) {
	score := &domain.ExamScore{}

	query := `
        SELECT id, exam_id, session_id, student_id, total_questions,
               correct_answers, wrong_answers, unanswered, score,
               created_at, updated_at
        FROM exam_scores
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&score.ID,
		&score.ExamID,
		&score.SessionID,
		&score.StudentID,
		&score.TotalQuestions,
		&score.CorrectAnswers,
		&score.WrongAnswers,
		&score.UnansweredCount,
		&score.Score,
		&score.CreatedAt,
		&score.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrScoreNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get score")
	}

	return score, nil
}

func (r *postgresRepository) GetScoreByExamAndStudent(ctx context.Context, examID, studentID string) (*domain.ExamScore, error) {
	score := &domain.ExamScore{}

	query := `
        SELECT id, exam_id, session_id, student_id, total_questions,
               correct_answers, wrong_answers, unanswered, score,
               created_at, updated_at
        FROM exam_scores
        WHERE exam_id = $1 AND student_id = $2`

	err := r.db.QueryRowContext(ctx, query, examID, studentID).Scan(
		&score.ID,
		&score.ExamID,
		&score.SessionID,
		&score.StudentID,
		&score.TotalQuestions,
		&score.CorrectAnswers,
		&score.WrongAnswers,
		&score.UnansweredCount,
		&score.Score,
		&score.CreatedAt,
		&score.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrScoreNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get score")
	}

	return score, nil
}

func (r *postgresRepository) ListScores(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.ExamScore, error) {
	query := `
        SELECT id, exam_id, session_id, student_id, total_questions,
               correct_answers, wrong_answers, unanswered, score,
               created_at, updated_at
        FROM exam_scores
        WHERE exam_id = $1
        ORDER BY score DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, examID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list scores")
	}
	defer rows.Close()

	var scores []*domain.ExamScore
	for rows.Next() {
		score := &domain.ExamScore{}
		err := rows.Scan(
			&score.ID,
			&score.ExamID,
			&score.SessionID,
			&score.StudentID,
			&score.TotalQuestions,
			&score.CorrectAnswers,
			&score.WrongAnswers,
			&score.UnansweredCount,
			&score.Score,
			&score.CreatedAt,
			&score.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan score")
		}
		scores = append(scores, score)
	}

	return scores, nil
}

func (r *postgresRepository) GetCorrectAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	query := `
        SELECT q.id, q.correct_answer, sa.selected_choice
        FROM exam_sessions es
        JOIN questions q ON q.exam_id = es.exam_id
        LEFT JOIN session_answers sa ON sa.question_id = q.id AND sa.session_id = es.id
        WHERE es.id = $1`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get answers")
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var answer domain.Answer
		err := rows.Scan(
			&answer.QuestionID,
			&answer.CorrectAnswer,
			&answer.StudentAnswer,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan answer")
		}
		answers = append(answers, answer)
	}

	if len(answers) == 0 {
		return nil, repository.ErrSessionNotFound
	}

	return answers, nil
}

func (r *postgresRepository) GetStudentAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	query := `
        SELECT sa.question_id, q.correct_answer, sa.selected_choice
        FROM session_answers sa
        JOIN questions q ON q.id = sa.question_id
        WHERE sa.session_id = $1`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get student answers")
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var answer domain.Answer
		err := rows.Scan(
			&answer.QuestionID,
			&answer.CorrectAnswer,
			&answer.StudentAnswer,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan answer")
		}
		answers = append(answers, answer)
	}

	if len(answers) == 0 {
		return nil, repository.ErrSessionNotFound
	}

	return answers, nil
}
