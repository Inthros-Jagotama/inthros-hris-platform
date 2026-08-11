-- =============================================================================
-- Tenant Migration: 090_training_advanced (PostgreSQL)
-- =============================================================================
-- Training & Development P2-BE (docs/module-training-development-plan.md §42 P2-BE):
--   - training_evaluation_forms + training_evaluation_questions + training_evaluation_answers
--   - training_effectiveness_assessments (before/after score, 30/60/90 hari)
--   - training_certifications (master) + ALTER training_certificates
--     (certification_id, certificate_file_url)
--
-- Semua statement idempotent (IF NOT EXISTS).

-- -----------------------------------------------------------------------------
-- 1. training_evaluation_forms — form evaluasi per session
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_forms (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    session_id  CHAR(36)     NOT NULL,
    name        VARCHAR(200) NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_eval_form_session ON training_evaluation_forms (session_id);

-- -----------------------------------------------------------------------------
-- 2. training_evaluation_questions
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_questions (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    form_id       CHAR(36)     NOT NULL,
    question      TEXT         NOT NULL,
    question_type VARCHAR(20)  NOT NULL DEFAULT 'RATING',
    sort_order    INT          NOT NULL DEFAULT 0,
    is_required   BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_eval_question_form ON training_evaluation_questions (form_id);

-- -----------------------------------------------------------------------------
-- 3. training_evaluation_answers
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_answers (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    question_id    CHAR(36)     NOT NULL,
    participant_id CHAR(36)     NOT NULL,
    answer         TEXT         NULL,
    deleted_at     TIMESTAMP,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_eval_answer_question ON training_evaluation_answers (question_id);
CREATE INDEX IF NOT EXISTS idx_trn_eval_answer_participant ON training_evaluation_answers (participant_id);
CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_eval_answer_q_p ON training_evaluation_answers (question_id, participant_id);

-- -----------------------------------------------------------------------------
-- 4. training_effectiveness_assessments — dampak pasca training
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_effectiveness_assessments (
    id                   CHAR(36)     NOT NULL PRIMARY KEY,
    participant_id       CHAR(36)     NOT NULL,
    assessment_date      DATE         NOT NULL,
    assessor_employee_id CHAR(36)     NULL,
    before_score         DECIMAL(5,2) NULL,
    after_score          DECIMAL(5,2) NULL,
    effectiveness_score  DECIMAL(5,2) NULL,
    remarks              TEXT         NULL,
    deleted_at           TIMESTAMP,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_effect_participant ON training_effectiveness_assessments (participant_id);

-- -----------------------------------------------------------------------------
-- 5. training_certifications — master sertifikasi
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_certifications (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    code                  VARCHAR(30)  NOT NULL,
    name                  VARCHAR(200) NOT NULL,
    issuing_body          VARCHAR(200) NULL,
    validity_period_month INT          NULL,
    renewal_required      BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active             BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_certif_code ON training_certifications (code);

-- -----------------------------------------------------------------------------
-- 6. ALTER training_certificates — tambah relasi certification + file
-- -----------------------------------------------------------------------------
ALTER TABLE training_certificates
    ADD COLUMN IF NOT EXISTS certification_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS certificate_file_url TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_trn_cert_certification ON training_certificates (certification_id);
