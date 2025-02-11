CREATE TYPE session_status AS ENUM ('STARTED', 'IN_PROGRESS', 'FINISHED', 'TIMEOUT');

CREATE TABLE exam_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    status session_status NOT NULL DEFAULT 'STARTED',
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (exam_id, student_id)
);

CREATE TABLE session_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES exam_sessions(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_choice VARCHAR(1) NOT NULL,
    answered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, question_id)
);

CREATE INDEX idx_session_exam ON exam_sessions(exam_id);
CREATE INDEX idx_session_student ON exam_sessions(student_id);
CREATE INDEX idx_session_status ON exam_sessions(status);
CREATE INDEX idx_answer_session ON session_answers(session_id);