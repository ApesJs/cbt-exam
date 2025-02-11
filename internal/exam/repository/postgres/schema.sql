CREATE TYPE exam_state AS ENUM ('CREATED', 'ACTIVE', 'FINISHED');
CREATE TYPE exam_student_state AS ENUM ('NOT_STARTED', 'IN_PROGRESS', 'FINISHED');

CREATE TABLE exams (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       title VARCHAR(255) NOT NULL,
                       subject VARCHAR(100) NOT NULL,
                       duration_mins INTEGER NOT NULL,
                       total_questions INTEGER NOT NULL,
                       is_random BOOLEAN DEFAULT false,
                       teacher_id UUID NOT NULL,
                       status exam_state NOT NULL DEFAULT 'CREATED',
                       start_time TIMESTAMP WITH TIME ZONE,
                       end_time TIMESTAMP WITH TIME ZONE,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exam_classes (
                              exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
                              class_id UUID NOT NULL,
                              PRIMARY KEY (exam_id, class_id)
);

CREATE TABLE exam_student_status (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     exam_id UUID REFERENCES exams(id) ON DELETE CASCADE,
                                     student_id UUID NOT NULL,
                                     student_name VARCHAR(255) NOT NULL,
                                     class_id UUID NOT NULL,
                                     state exam_student_state NOT NULL DEFAULT 'NOT_STARTED',
                                     start_time TIMESTAMP WITH TIME ZONE,
                                     end_time TIMESTAMP WITH TIME ZONE,
                                     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     UNIQUE (exam_id, student_id)
);

CREATE INDEX idx_exam_teacher ON exams(teacher_id);
CREATE INDEX idx_exam_status ON exams(status);
CREATE INDEX idx_student_status_exam ON exam_student_status(exam_id);
CREATE INDEX idx_student_status_student ON exam_student_status(student_id);