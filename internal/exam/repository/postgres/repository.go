package postgres

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"cbt-exam/internal/exam/domain"
	"cbt-exam/internal/exam/repository"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) repository.ExamRepository {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(ctx context.Context, exam *domain.Exam) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert exam
	query := `
        INSERT INTO exams (title, subject, duration_mins, total_questions, is_random, teacher_id, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		exam.Title,
		exam.Subject,
		exam.DurationMins,
		exam.TotalQuestions,
		exam.IsRandom,
		exam.TeacherID,
		exam.Status,
	).Scan(&exam.ID, &exam.CreatedAt, &exam.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to insert exam")
	}

	// Insert exam classes
	if len(exam.ClassIDs) > 0 {
		classQuery := `
            INSERT INTO exam_classes (exam_id, class_id)
            VALUES ($1, unnest($2::uuid[]))`

		_, err = tx.ExecContext(ctx, classQuery, exam.ID, pq.Array(exam.ClassIDs))
		if err != nil {
			return errors.Wrap(err, "failed to insert exam classes")
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.Exam, error) {
	exam := &domain.Exam{}

	query := `
        SELECT e.id, e.title, e.subject, e.duration_mins, e.total_questions, 
               e.is_random, e.teacher_id, e.status, e.start_time, e.end_time,
               e.created_at, e.updated_at,
               array_agg(ec.class_id) as class_ids
        FROM exams e
        LEFT JOIN exam_classes ec ON e.id = ec.exam_id
        WHERE e.id = $1
        GROUP BY e.id`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exam.ID,
		&exam.Title,
		&exam.Subject,
		&exam.DurationMins,
		&exam.TotalQuestions,
		&exam.IsRandom,
		&exam.TeacherID,
		&exam.Status,
		&exam.StartTime,
		&exam.EndTime,
		&exam.CreatedAt,
		&exam.UpdatedAt,
		pq.Array(&exam.ClassIDs),
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrExamNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exam")
	}

	return exam, nil
}

func (r *postgresRepository) List(ctx context.Context, teacherID string, limit int32, offset int32) ([]*domain.Exam, error) {
	query := `
        SELECT e.id, e.title, e.subject, e.duration_mins, e.total_questions, 
               e.is_random, e.teacher_id, e.status, e.start_time, e.end_time,
               e.created_at, e.updated_at,
               array_agg(ec.class_id) as class_ids
        FROM exams e
        LEFT JOIN exam_classes ec ON e.id = ec.exam_id
        WHERE e.teacher_id = $1
        GROUP BY e.id
        ORDER BY e.created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, teacherID, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list exams")
	}
	defer rows.Close()

	var exams []*domain.Exam
	for rows.Next() {
		exam := &domain.Exam{}
		err := rows.Scan(
			&exam.ID,
			&exam.Title,
			&exam.Subject,
			&exam.DurationMins,
			&exam.TotalQuestions,
			&exam.IsRandom,
			&exam.TeacherID,
			&exam.Status,
			&exam.StartTime,
			&exam.EndTime,
			&exam.CreatedAt,
			&exam.UpdatedAt,
			pq.Array(&exam.ClassIDs),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan exam row")
		}
		exams = append(exams, exam)
	}

	return exams, nil
}

func (r *postgresRepository) Update(ctx context.Context, exam *domain.Exam) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        UPDATE exams 
        SET title = $1, subject = $2, duration_mins = $3, total_questions = $4,
            is_random = $5, status = $6, updated_at = CURRENT_TIMESTAMP
        WHERE id = $7
        RETURNING updated_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		exam.Title,
		exam.Subject,
		exam.DurationMins,
		exam.TotalQuestions,
		exam.IsRandom,
		exam.Status,
		exam.ID,
	).Scan(&exam.UpdatedAt)

	if err == sql.ErrNoRows {
		return repository.ErrExamNotFound
	}
	if err != nil {
		return errors.Wrap(err, "failed to update exam")
	}

	// Update exam classes
	_, err = tx.ExecContext(ctx, "DELETE FROM exam_classes WHERE exam_id = $1", exam.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete old exam classes")
	}

	if len(exam.ClassIDs) > 0 {
		classQuery := `
            INSERT INTO exam_classes (exam_id, class_id)
            VALUES ($1, unnest($2::uuid[]))`

		_, err = tx.ExecContext(ctx, classQuery, exam.ID, pq.Array(exam.ClassIDs))
		if err != nil {
			return errors.Wrap(err, "failed to insert new exam classes")
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM exams WHERE id = $1", id)
	if err != nil {
		return errors.Wrap(err, "failed to delete exam")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrExamNotFound
	}

	return nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, examID string, status domain.ExamState) error {
	query := `
        UPDATE exams 
        SET status = $1, updated_at = CURRENT_TIMESTAMP,
            start_time = CASE 
                WHEN $1 = 'ACTIVE' THEN CURRENT_TIMESTAMP 
                ELSE start_time 
            END,
            end_time = CASE 
                WHEN $1 = 'FINISHED' THEN CURRENT_TIMESTAMP 
                ELSE end_time 
            END
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, examID)
	if err != nil {
		return errors.Wrap(err, "failed to update exam status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrExamNotFound
	}

	return nil
}

func (r *postgresRepository) GetStatus(ctx context.Context, examID string) (*domain.ExamStatus, error) {
	status := &domain.ExamStatus{
		ExamID: examID,
	}

	// Get exam status and counts
	query := `
        SELECT e.status,
               COUNT(DISTINCT es.student_id) as total_students,
               COUNT(CASE WHEN es.state != 'NOT_STARTED' THEN 1 END) as students_started,
               COUNT(CASE WHEN es.state = 'FINISHED' THEN 1 END) as students_finished
        FROM exams e
        LEFT JOIN exam_student_status es ON e.id = es.exam_id
        WHERE e.id = $1
        GROUP BY e.id, e.status`

	err := r.db.QueryRowContext(ctx, query, examID).Scan(
		&status.State,
		&status.TotalStudents,
		&status.StudentsStarted,
		&status.StudentsFinished,
	)

	if err == sql.ErrNoRows {
		return nil, repository.ErrExamNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exam status")
	}

	// Get student statuses
	studentQuery := `
        SELECT student_id, student_name, class_id, state, start_time, end_time
        FROM exam_student_status
        WHERE exam_id = $1
        ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, studentQuery, examID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get student statuses")
	}
	defer rows.Close()

	for rows.Next() {
		var studentStatus domain.StudentStatus
		err := rows.Scan(
			&studentStatus.StudentID,
			&studentStatus.StudentName,
			&studentStatus.ClassID,
			&studentStatus.State,
			&studentStatus.StartTime,
			&studentStatus.EndTime,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan student status")
		}
		status.StudentStatuses = append(status.StudentStatuses, studentStatus)
	}

	return status, nil
}

func (r *postgresRepository) UpdateStudentStatus(ctx context.Context, examID string, studentStatus *domain.StudentStatus) error {
	query := `
        INSERT INTO exam_student_status (
            exam_id, student_id, student_name, class_id, state, 
            start_time, end_time
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (exam_id, student_id) DO UPDATE
        SET state = EXCLUDED.state,
            start_time = COALESCE(exam_student_status.start_time, EXCLUDED.start_time),
            end_time = EXCLUDED.end_time,
            updated_at = CURRENT_TIMESTAMP`

	result, err := r.db.ExecContext(
		ctx,
		query,
		examID,
		studentStatus.StudentID,
		studentStatus.StudentName,
		studentStatus.ClassID,
		studentStatus.State,
		studentStatus.StartTime,
		studentStatus.EndTime,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update student status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrExamNotFound
	}

	return nil
}
