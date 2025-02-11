package postgres

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/ApesJs/cbt-exam/internal/session/domain"
	"github.com/ApesJs/cbt-exam/internal/session/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) repository.SessionRepository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) StartSession(ctx context.Context, session *domain.ExamSession) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify exam is active
	var examActive bool
	err = tx.QueryRowContext(ctx,
		"SELECT status = 'ACTIVE' FROM exams WHERE id = $1",
		session.ExamID,
	).Scan(&examActive)

	if err == sql.ErrNoRows {
		return repository.ErrExamNotFound
	}
	if err != nil {
		return errors.Wrap(err, "failed to check exam status")
	}
	if !examActive {
		return repository.ErrExamNotActive
	}

	// Check for existing active session
	var activeCount int
	err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM exam_sessions 
         WHERE student_id = $1 AND status IN ('STARTED', 'IN_PROGRESS')`,
		session.StudentID,
	).Scan(&activeCount)
	if err != nil {
		return errors.Wrap(err, "failed to check active sessions")
	}
	if activeCount > 0 {
		return repository.ErrDuplicateSession
	}

	// Insert new session
	query := `
        INSERT INTO exam_sessions (exam_id, student_id, status, start_time)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		session.ExamID,
		session.StudentID,
		session.Status,
		session.StartTime,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to create session")
	}

	return tx.Commit()
}

func (r *postgresRepository) GetSession(ctx context.Context, id string) (*domain.ExamSession, error) {
	session := &domain.ExamSession{}

	query := `
        SELECT id, exam_id, student_id, status, start_time, end_time, 
               created_at, updated_at
        FROM exam_sessions 
        WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.ExamID,
		&session.StudentID,
		&session.Status,
		&session.StartTime,
		&session.EndTime,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrSessionNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	// Get answers
	answers, err := r.GetSessionAnswers(ctx, id)
	if err != nil {
		return nil, err
	}
	session.Answers = answers

	return session, nil
}

func (r *postgresRepository) UpdateSessionStatus(ctx context.Context, id string, status domain.SessionStatus) error {
	query := `
        UPDATE exam_sessions 
        SET status = $1, updated_at = CURRENT_TIMESTAMP,
            end_time = CASE 
                WHEN $1 IN ('FINISHED', 'TIMEOUT') THEN CURRENT_TIMESTAMP 
                ELSE end_time 
            END
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return errors.Wrap(err, "failed to update session status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrSessionNotFound
	}

	return nil
}

func (r *postgresRepository) FinishSession(ctx context.Context, id string) error {
	return r.UpdateSessionStatus(ctx, id, domain.SessionStatusFinished)
}

func (r *postgresRepository) SubmitAnswer(ctx context.Context, sessionID string, answer domain.Answer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify session is active
	var status domain.SessionStatus
	err = tx.QueryRowContext(ctx,
		"SELECT status FROM exam_sessions WHERE id = $1",
		sessionID,
	).Scan(&status)

	if err == sql.ErrNoRows {
		return repository.ErrSessionNotFound
	}
	if err != nil {
		return errors.Wrap(err, "failed to check session status")
	}

	if status != domain.SessionStatusStarted && status != domain.SessionStatusInProgress {
		return repository.ErrInvalidSessionState
	}

	// Update or insert answer
	query := `
        INSERT INTO session_answers (session_id, question_id, selected_choice, answered_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (session_id, question_id) 
        DO UPDATE SET selected_choice = EXCLUDED.selected_choice,
                      answered_at = EXCLUDED.answered_at`

	_, err = tx.ExecContext(ctx, query,
		sessionID,
		answer.QuestionID,
		answer.SelectedChoice,
		answer.AnsweredAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to submit answer")
	}

	// Update session status to in progress if it was just started
	if status == domain.SessionStatusStarted {
		err = r.UpdateSessionStatus(ctx, sessionID, domain.SessionStatusInProgress)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetSessionAnswers(ctx context.Context, sessionID string) ([]domain.Answer, error) {
	query := `
        SELECT question_id, selected_choice, answered_at
        FROM session_answers
        WHERE session_id = $1
        ORDER BY answered_at`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session answers")
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var answer domain.Answer
		err := rows.Scan(
			&answer.QuestionID,
			&answer.SelectedChoice,
			&answer.AnsweredAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan answer")
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *postgresRepository) IsExamActive(ctx context.Context, examID string) (bool, error) {
	var active bool
	err := r.db.QueryRowContext(ctx,
		"SELECT status = 'ACTIVE' FROM exams WHERE id = $1",
		examID,
	).Scan(&active)

	if err == sql.ErrNoRows {
		return false, repository.ErrExamNotFound
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to check exam status")
	}

	return active, nil
}

func (r *postgresRepository) HasActiveSession(ctx context.Context, studentID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM exam_sessions 
         WHERE student_id = $1 AND status IN ('STARTED', 'IN_PROGRESS')`,
		studentID,
	).Scan(&count)

	if err != nil {
		return false, errors.Wrap(err, "failed to check active sessions")
	}

	return count > 0, nil
}
