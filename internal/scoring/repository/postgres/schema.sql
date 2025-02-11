CREATE TABLE exam_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES exam_sessions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    total_questions INTEGER NOT NULL,
    correct_answers INTEGER NOT NULL,
    wrong_answers INTEGER NOT NULL,
    unanswered INTEGER NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (exam_id, student_id)
);

-- Indeks untuk mempercepat query
CREATE INDEX idx_score_exam ON exam_scores(exam_id);
CREATE INDEX idx_score_student ON exam_scores(student_id);
CREATE INDEX idx_score_session ON exam_scores(session_id);
CREATE INDEX idx_score_created ON exam_scores(created_at);