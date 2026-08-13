-- =============================================================================
-- Tenant Migration: 104_recruitment_interview_scorecard (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS interviewers (
    id           CHAR(36) PRIMARY KEY,
    interview_id CHAR(36) NOT NULL,
    employee_id  CHAR(36) NOT NULL,
    role         VARCHAR(50) NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_interviewer_int FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE,
    CONSTRAINT uq_interviewer_int_emp UNIQUE (interview_id, employee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_interviewer_int ON interviewers (interview_id);

CREATE TABLE IF NOT EXISTS interview_scorecard_items (
    id           CHAR(36) PRIMARY KEY,
    interview_id CHAR(36) NOT NULL,
    criterion    VARCHAR(255) NOT NULL,
    weight       DECIMAL(5,2) NOT NULL DEFAULT 0,
    score        DECIMAL(5,2) NULL,
    notes        TEXT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_scorecard_int FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_scorecard_int ON interview_scorecard_items (interview_id);
