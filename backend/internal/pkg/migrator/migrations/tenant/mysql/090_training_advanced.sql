-- =============================================================================
-- Tenant Migration: 090_training_advanced (MySQL)
-- =============================================================================
-- Training & Development P2-BE (docs/module-training-development-plan.md §42 P2-BE):
--   - training_evaluation_forms + training_evaluation_questions + training_evaluation_answers
--   - training_effectiveness_assessments (before/after score, 30/60/90 hari)
--   - training_certifications (master) + ALTER training_certificates
--     (certification_id, certificate_file_url)
--
-- Idempotent: CREATE TABLE IF NOT EXISTS; ALTER via information_schema + PREPARE/EXECUTE.

-- -----------------------------------------------------------------------------
-- 1. training_evaluation_forms
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_forms (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    session_id  CHAR(36)     NOT NULL,
    name        VARCHAR(200) NOT NULL,
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_eval_form_session (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 2. training_evaluation_questions
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_questions (
    id            CHAR(36)    NOT NULL PRIMARY KEY,
    form_id       CHAR(36)    NOT NULL,
    question      TEXT        NOT NULL,
    question_type VARCHAR(20) NOT NULL DEFAULT 'RATING',
    sort_order    INT         NOT NULL DEFAULT 0,
    is_required   TINYINT(1)  NOT NULL DEFAULT 1,
    deleted_at    TIMESTAMP(6) NULL DEFAULT NULL,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_eval_question_form (form_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 3. training_evaluation_answers
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_evaluation_answers (
    id             CHAR(36)    NOT NULL PRIMARY KEY,
    question_id    CHAR(36)    NOT NULL,
    participant_id CHAR(36)    NOT NULL,
    answer         TEXT        NULL,
    deleted_at     TIMESTAMP(6) NULL DEFAULT NULL,
    created_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uk_trn_eval_answer_q_p (question_id, participant_id),
    INDEX idx_trn_eval_answer_question (question_id),
    INDEX idx_trn_eval_answer_participant (participant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 4. training_effectiveness_assessments
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
    deleted_at           TIMESTAMP(6) NULL DEFAULT NULL,
    created_at           TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at           TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_effect_participant (participant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 5. training_certifications — master sertifikasi
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_certifications (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    code                  VARCHAR(30)  NOT NULL,
    name                  VARCHAR(200) NOT NULL,
    issuing_body          VARCHAR(200) NULL,
    validity_period_month INT          NULL,
    renewal_required      TINYINT(1)   NOT NULL DEFAULT 0,
    is_active             TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at            TIMESTAMP(6) NULL DEFAULT NULL,
    created_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uk_trn_certif_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 6. ALTER training_certificates — tambah relasi certification + file (idempotent)
-- -----------------------------------------------------------------------------
SET @db = DATABASE();

SET @add_certification_id = (
    SELECT IF(
        COUNT(*) = 0,
        'ALTER TABLE training_certificates ADD COLUMN certification_id CHAR(36) NULL AFTER participant_id',
        'SELECT 1'
    )
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'training_certificates'
      AND COLUMN_NAME = 'certification_id'
);
PREPARE stmt FROM @add_certification_id; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_cert_file_url = (
    SELECT IF(
        COUNT(*) = 0,
        'ALTER TABLE training_certificates ADD COLUMN certificate_file_url TEXT NULL AFTER certificate_no',
        'SELECT 1'
    )
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'training_certificates'
      AND COLUMN_NAME = 'certificate_file_url'
);
PREPARE stmt FROM @add_cert_file_url; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_idx_cert_cert = (
    SELECT IF(
        COUNT(*) = 0,
        'ALTER TABLE training_certificates ADD INDEX idx_trn_cert_certification (certification_id)',
        'SELECT 1'
    )
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'training_certificates'
      AND INDEX_NAME = 'idx_trn_cert_certification'
);
PREPARE stmt FROM @add_idx_cert_cert; EXECUTE stmt; DEALLOCATE PREPARE stmt;
