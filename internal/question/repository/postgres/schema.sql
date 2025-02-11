CREATE TABLE questions (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           exam_id UUID NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
                           question_text TEXT NOT NULL,
                           correct_answer VARCHAR(1) NOT NULL,
                           created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE choices (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
                         text TEXT NOT NULL,
                         created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_question_exam ON questions(exam_id);
CREATE INDEX idx_choices_question ON choices(question_id);