-- =============================================================================
-- Tenant Migration: 104_recruitment_interview_scorecard (PostgreSQL)
-- =============================================================================
-- G-8: interviewers (multi-interviewer per interview; interviews.interviewer_id
-- existing tetap sebagai pewawancara utama — backward compatible, tidak
-- breaking existing data/FE) + interview_scorecard_items (kriteria berbobot
-- bebas per interview, tanpa master — mis. "Technical Skill" 30%). Endpoint
-- POST /interviews/:id/complete menghitung weighted average scorecard items
-- ke interviews.score existing (bukan kolom baru).
-- (docs/module-recruitment-development-plan.md §G-8)
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS interviewers (
    id           CHAR(36) PRIMARY KEY,
    interview_id CHAR(36) NOT NULL,
    employee_id  CHAR(36) NOT NULL,
    role         VARCHAR(50) NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_interviewer_int FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE,
    CONSTRAINT uq_interviewer_int_emp UNIQUE (interview_id, employee_id)
);

CREATE INDEX IF NOT EXISTS idx_interviewer_int ON interviewers (interview_id);

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
);

CREATE INDEX IF NOT EXISTS idx_scorecard_int ON interview_scorecard_items (interview_id);
