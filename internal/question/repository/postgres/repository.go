package postgres

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/ApesJs/cbt-exam/internal/question/domain"
	"github.com/ApesJs/cbt-exam/internal/question/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) repository.QuestionRepository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(ctx context.Context, question *domain.Question) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert question
	query := `
        INSERT INTO questions (exam_id, question_text, correct_answer)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		question.ExamID,
		question.QuestionText,
		question.CorrectAnswer,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to insert question")
	}

	// Insert choices
	if len(question.Choices) > 0 {
		choiceQuery := `
            INSERT INTO choices (question_id, text)
            VALUES ($1, $2)
            RETURNING id`

		for i := range question.Choices {
			err = tx.QueryRowContext(
				ctx,
				choiceQuery,
				question.ID,
				question.Choices[i].Text,
			).Scan(&question.Choices[i].ID)
			if err != nil {
				return errors.Wrap(err, "failed to insert choice")
			}
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.Question, error) {
	query := `
        SELECT q.id, q.exam_id, q.question_text, q.correct_answer, 
               q.created_at, q.updated_at
        FROM questions q
        WHERE q.id = $1`

	question := &domain.Question{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&question.ID,
		&question.ExamID,
		&question.QuestionText,
		&question.CorrectAnswer,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrQuestionNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question")
	}

	// Get choices
	choiceQuery := `
        SELECT id, text
        FROM choices
        WHERE question_id = $1
        ORDER BY id`

	rows, err := r.db.QueryContext(ctx, choiceQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get choices")
	}
	defer rows.Close()

	for rows.Next() {
		var choice domain.Choice
		err := rows.Scan(&choice.ID, &choice.Text)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan choice")
		}
		question.Choices = append(question.Choices, choice)
	}

	return question, nil
}

func (r *postgresRepository) List(ctx context.Context, examID string, limit int32, offset int32) ([]*domain.Question, error) {
	query := `
        SELECT q.id, q.exam_id, q.question_text, q.correct_answer, 
               q.created_at, q.updated_at
        FROM questions q
        WHERE q.exam_id = $1
        ORDER BY q.created_at
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, examID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list questions")
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		question := &domain.Question{}
		err := rows.Scan(
			&question.ID,
			&question.ExamID,
			&question.QuestionText,
			&question.CorrectAnswer,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan question")
		}
		questions = append(questions, question)
	}

	// Get choices for each question
	for _, q := range questions {
		choiceQuery := `
            SELECT id, text
            FROM choices
            WHERE question_id = $1
            ORDER BY id`

		rows, err := r.db.QueryContext(ctx, choiceQuery, q.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get choices")
		}
		defer rows.Close()

		for rows.Next() {
			var choice domain.Choice
			err := rows.Scan(&choice.ID, &choice.Text)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan choice")
			}
			q.Choices = append(q.Choices, choice)
		}
	}

	return questions, nil
}

func (r *postgresRepository) Update(ctx context.Context, question *domain.Question) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update question
	query := `
        UPDATE questions 
        SET question_text = $1, correct_answer = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
        RETURNING updated_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		question.QuestionText,
		question.CorrectAnswer,
		question.ID,
	).Scan(&question.UpdatedAt)

	if err == sql.ErrNoRows {
		return repository.ErrQuestionNotFound
	}
	if err != nil {
		return errors.Wrap(err, "failed to update question")
	}

	// Delete existing choices
	_, err = tx.ExecContext(ctx, "DELETE FROM choices WHERE question_id = $1", question.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete old choices")
	}

	// Insert new choices
	if len(question.Choices) > 0 {
		choiceQuery := `
            INSERT INTO choices (question_id, text)
            VALUES ($1, $2)
            RETURNING id`

		for i := range question.Choices {
			err = tx.QueryRowContext(
				ctx,
				choiceQuery,
				question.ID,
				question.Choices[i].Text,
			).Scan(&question.Choices[i].ID)
			if err != nil {
				return errors.Wrap(err, "failed to insert choice")
			}
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM questions WHERE id = $1", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete question")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrQuestionNotFound
	}

	return nil
}

func (r *postgresRepository) GetExamQuestions(ctx context.Context, filter domain.QuestionFilter) ([]*domain.Question, error) {
	var query string
	if filter.Randomize {
		query = `
            SELECT q.id, q.exam_id, q.question_text, q.correct_answer, 
                   q.created_at, q.updated_at
            FROM questions q
            WHERE q.exam_id = $1
            ORDER BY RANDOM()
            LIMIT $2`
	} else {
		query = `
            SELECT q.id, q.exam_id, q.question_text, q.correct_answer, 
                   q.created_at, q.updated_at
            FROM questions q
            WHERE q.exam_id = $1
            ORDER BY q.created_at
            LIMIT $2`
	}

	rows, err := r.db.QueryContext(ctx, query, filter.ExamID, filter.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exam questions")
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		question := &domain.Question{}
		err := rows.Scan(
			&question.ID,
			&question.ExamID,
			&question.QuestionText,
			&question.CorrectAnswer,
			&question.CreatedAt,
			&question.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan question")
		}
		questions = append(questions, question)
	}

	// Get choices for each question
	for _, q := range questions {
		choiceQuery := `
            SELECT id, text
            FROM choices
            WHERE question_id = $1
            ORDER BY id`

		rows, err := r.db.QueryContext(ctx, choiceQuery, q.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get choices")
		}
		defer rows.Close()

		for rows.Next() {
			var choice domain.Choice
			err := rows.Scan(&choice.ID, &choice.Text)
			if err != nil {
				return nil, errors.Wrap(err, "failed to scan choice")
			}
			q.Choices = append(q.Choices, choice)
		}
	}

	return questions, nil
}

func (r *postgresRepository) CountExamQuestions(ctx context.Context, examID string) (int32, error) {
	var count int32
	err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM questions WHERE exam_id = $1",
		examID,
	).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count questions")
	}
	return count, nil
}
