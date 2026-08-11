-- =============================================================================
-- Tenant Migration: 088_training_core (MySQL)
-- =============================================================================
-- Training & Development P0-BE (docs/module-training-development-plan.md §42 P0-BE):
--   - master provider & trainer + relasi session-trainer
--   - enhancement training_courses / training_sessions / training_participants / training_materials
--   - attendance detail (training_attendances)
--   - assessment + assessment results
--
-- Semua statement idempotent (CREATE TABLE IF NOT EXISTS + guard
-- information_schema + PREPARE/EXECUTE untuk ALTER). Kolom lama
-- `external_vendor` (courses) & `trainer_name` (sessions) TIDAK di-drop (deprecate).

-- -----------------------------------------------------------------------------
-- 1. training_providers
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_providers (
    id           CHAR(36)     NOT NULL PRIMARY KEY,
    code         VARCHAR(20)  NOT NULL,
    name         VARCHAR(200) NOT NULL,
    type         VARCHAR(20)  NOT NULL DEFAULT 'EXTERNAL',
    contact_name VARCHAR(150) NULL,
    email        VARCHAR(150) NULL,
    phone        VARCHAR(50)  NULL,
    address      TEXT         NULL,
    website      VARCHAR(200) NULL,
    is_active    TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at   TIMESTAMP(6) NULL DEFAULT NULL,
    created_at   TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at   TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_trn_provider_code (code),
    INDEX idx_trn_provider_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 2. training_trainers
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_trainers (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    type        VARCHAR(20)  NOT NULL,
    employee_id CHAR(36)     NULL,
    provider_id CHAR(36)     NULL,
    name        VARCHAR(200) NOT NULL,
    email       VARCHAR(150) NULL,
    phone       VARCHAR(50)  NULL,
    bio         TEXT         NULL,
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_trainer_emp (employee_id),
    INDEX idx_trn_trainer_provider (provider_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 3. training_session_trainers
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_session_trainers (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    session_id CHAR(36)     NOT NULL,
    trainer_id CHAR(36)     NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'MAIN',
    deleted_at TIMESTAMP(6) NULL DEFAULT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_sess_trn_session (session_id),
    INDEX idx_trn_sess_trn_trainer (trainer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 4. ALTER training_courses
-- -----------------------------------------------------------------------------
SET @add_course_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='course_type'),
  'DO 0',
  'ALTER TABLE training_courses ADD COLUMN course_type VARCHAR(20) NULL AFTER external_vendor'
);
PREPARE stmt FROM @add_course_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_delivery_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='delivery_type'),
  'DO 0',
  'ALTER TABLE training_courses ADD COLUMN delivery_type VARCHAR(20) NULL AFTER course_type'
);
PREPARE stmt FROM @add_delivery_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_is_mandatory = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_courses' AND COLUMN_NAME='is_mandatory'),
  'DO 0',
  'ALTER TABLE training_courses ADD COLUMN is_mandatory TINYINT(1) NOT NULL DEFAULT 0 AFTER delivery_type'
);
PREPARE stmt FROM @add_is_mandatory; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 5. ALTER training_sessions
-- -----------------------------------------------------------------------------
SET @add_provider_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='provider_type'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN provider_type VARCHAR(20) NULL AFTER trainer_name'
);
PREPARE stmt FROM @add_provider_type; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_delivery_mode = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='delivery_mode'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN delivery_mode VARCHAR(20) NULL AFTER provider_type'
);
PREPARE stmt FROM @add_delivery_mode; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_provider_id = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='provider_id'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN provider_id CHAR(36) NULL AFTER delivery_mode'
);
PREPARE stmt FROM @add_provider_id; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_start_datetime = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='start_datetime'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN start_datetime TIMESTAMP(6) NULL AFTER provider_id'
);
PREPARE stmt FROM @add_start_datetime; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_end_datetime = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='end_datetime'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN end_datetime TIMESTAMP(6) NULL AFTER start_datetime'
);
PREPARE stmt FROM @add_end_datetime; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_meeting_url = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='meeting_url'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN meeting_url TEXT NULL AFTER end_datetime'
);
PREPARE stmt FROM @add_meeting_url; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_registration_deadline = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND COLUMN_NAME='registration_deadline'),
  'DO 0',
  'ALTER TABLE training_sessions ADD COLUMN registration_deadline TIMESTAMP(6) NULL AFTER meeting_url'
);
PREPARE stmt FROM @add_registration_deadline; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_sess_provider_index = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_sessions' AND INDEX_NAME='idx_trn_sess_provider'),
  'DO 0',
  'CREATE INDEX idx_trn_sess_provider ON training_sessions (provider_id)'
);
PREPARE stmt FROM @add_sess_provider_index; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 6. ALTER training_participants
-- -----------------------------------------------------------------------------
SET @add_reg_status = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='registration_status'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN registration_status VARCHAR(20) NOT NULL DEFAULT ''REGISTERED'' AFTER employee_id'
);
PREPARE stmt FROM @add_reg_status; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_registered_at = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='registered_at'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN registered_at TIMESTAMP(6) NULL AFTER registration_status'
);
PREPARE stmt FROM @add_registered_at; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_approved_at = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='approved_at'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN approved_at TIMESTAMP(6) NULL AFTER registered_at'
);
PREPARE stmt FROM @add_approved_at; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_completion_status = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='completion_status'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN completion_status VARCHAR(20) NOT NULL DEFAULT ''NOT_STARTED'' AFTER attendance_status'
);
PREPARE stmt FROM @add_completion_status; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_completion_date = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='completion_date'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN completion_date DATE NULL AFTER completion_status'
);
PREPARE stmt FROM @add_completion_date; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_final_score = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='final_score'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN final_score DECIMAL(5,2) NULL AFTER completion_date'
);
PREPARE stmt FROM @add_final_score; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_passed = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='passed'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN passed TINYINT(1) NULL AFTER final_score'
);
PREPARE stmt FROM @add_passed; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_remarks = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='remarks'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN remarks TEXT NULL AFTER passed'
);
PREPARE stmt FROM @add_remarks; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Dedupe record aktif sebelum index unik (MySQL tidak punya partial index —
-- pakai generated column yang NULL untuk record soft-deleted, sehingga unique
-- hanya berlaku untuk (session_id, employee_id) aktif).
DELETE p1 FROM training_participants p1
INNER JOIN training_participants p2
  ON p1.session_id = p2.session_id AND p1.employee_id = p2.employee_id
WHERE p1.id > p2.id AND p1.deleted_at IS NULL AND p2.deleted_at IS NULL;

SET @add_active_key = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND COLUMN_NAME='active_session_employee'),
  'DO 0',
  'ALTER TABLE training_participants ADD COLUMN active_session_employee CHAR(73) GENERATED ALWAYS AS (IF(deleted_at IS NULL, CONCAT(session_id, ''|'', employee_id), NULL)) STORED'
);
PREPARE stmt FROM @add_active_key; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_part_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_participants' AND INDEX_NAME='uk_trn_part_active_key'),
  'DO 0',
  'CREATE UNIQUE INDEX uk_trn_part_active_key ON training_participants (active_session_employee)'
);
PREPARE stmt FROM @add_part_uk; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 7. ALTER training_materials
-- -----------------------------------------------------------------------------
SET @add_mat_description = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='description'),
  'DO 0',
  'ALTER TABLE training_materials ADD COLUMN description TEXT NULL AFTER title'
);
PREPARE stmt FROM @add_mat_description; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_mat_is_required = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='is_required'),
  'DO 0',
  'ALTER TABLE training_materials ADD COLUMN is_required TINYINT(1) NOT NULL DEFAULT 0 AFTER description'
);
PREPARE stmt FROM @add_mat_is_required; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_mat_available_from = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_materials' AND COLUMN_NAME='available_from'),
  'DO 0',
  'ALTER TABLE training_materials ADD COLUMN available_from TIMESTAMP(6) NULL AFTER is_required'
);
PREPARE stmt FROM @add_mat_available_from; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 8. training_attendances
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_attendances (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    participant_id  CHAR(36)     NOT NULL,
    attendance_date DATE         NOT NULL,
    check_in        TIMESTAMP(6) NULL,
    check_out       TIMESTAMP(6) NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'PRESENT',
    remarks         TEXT         NULL,
    deleted_at      TIMESTAMP(6) NULL DEFAULT NULL,
    created_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_att_part (participant_id),
    INDEX idx_trn_att_date (attendance_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Unique (participant_id, attendance_date) untuk record aktif via generated column.
SET @add_att_active_key = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_attendances' AND COLUMN_NAME='att_active_key'),
  'DO 0',
  'ALTER TABLE training_attendances ADD COLUMN att_active_key CHAR(73) GENERATED ALWAYS AS (IF(deleted_at IS NULL, CONCAT(participant_id, ''|'', attendance_date), NULL)) STORED'
);
PREPARE stmt FROM @add_att_active_key; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_att_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_attendances' AND INDEX_NAME='uk_trn_att_active_key'),
  'DO 0',
  'CREATE UNIQUE INDEX uk_trn_att_active_key ON training_attendances (att_active_key)'
);
PREPARE stmt FROM @add_att_uk; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 9. training_assessments
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_assessments (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    session_id    CHAR(36)     NOT NULL,
    name          VARCHAR(200) NOT NULL,
    type          VARCHAR(20)  NOT NULL DEFAULT 'OTHER',
    max_score     DECIMAL(8,2) NOT NULL DEFAULT 100.00,
    passing_score DECIMAL(8,2) NOT NULL DEFAULT 60.00,
    attempt_limit INT          NOT NULL DEFAULT 1,
    is_required   TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at    TIMESTAMP(6) NULL DEFAULT NULL,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_assess_session (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 10. training_assessment_results
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_assessment_results (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    assessment_id  CHAR(36)     NOT NULL,
    participant_id CHAR(36)     NOT NULL,
    score          DECIMAL(8,2) NOT NULL DEFAULT 0.00,
    passed         TINYINT(1)   NOT NULL DEFAULT 0,
    attempt        INT          NOT NULL DEFAULT 1,
    completed_at   TIMESTAMP(6) NULL,
    deleted_at     TIMESTAMP(6) NULL DEFAULT NULL,
    created_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_assess_res_assessment (assessment_id),
    INDEX idx_trn_assess_res_part (participant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Unique (assessment_id, participant_id, attempt) untuk record aktif.
SET @add_res_active_key = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_assessment_results' AND COLUMN_NAME='res_active_key'),
  'DO 0',
  'ALTER TABLE training_assessment_results ADD COLUMN res_active_key CHAR(110) GENERATED ALWAYS AS (IF(deleted_at IS NULL, CONCAT(assessment_id, ''|'', participant_id, ''|'', attempt), NULL)) STORED'
);
PREPARE stmt FROM @add_res_active_key; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @add_res_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='training_assessment_results' AND INDEX_NAME='uk_trn_assess_res_active_key'),
  'DO 0',
  'CREATE UNIQUE INDEX uk_trn_assess_res_active_key ON training_assessment_results (res_active_key)'
);
PREPARE stmt FROM @add_res_uk; EXECUTE stmt; DEALLOCATE PREPARE stmt;
