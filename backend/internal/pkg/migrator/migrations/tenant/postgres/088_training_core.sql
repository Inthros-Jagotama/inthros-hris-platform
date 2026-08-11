-- =============================================================================
-- Tenant Migration: 088_training_core (PostgreSQL)
-- =============================================================================
-- Training & Development P0-BE (docs/module-training-development-plan.md §42 P0-BE):
--   - master provider & trainer + relasi session-trainer
--   - enhancement training_courses / training_sessions / training_participants / training_materials
--   - attendance detail (training_attendances)
--   - assessment + assessment results
--
-- Semua statement idempotent (IF NOT EXISTS / ADD COLUMN IF NOT EXISTS).
-- Kolom lama `external_vendor` (courses) dan `trainer_name` (sessions) TIDAK di-drop
-- (deprecate untuk kompatibilitas — plan §6).

-- -----------------------------------------------------------------------------
-- 1. training_providers — master penyelenggara (type: INTERNAL | EXTERNAL)
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
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at   TIMESTAMP,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trn_provider_code ON training_providers (code);
CREATE INDEX IF NOT EXISTS idx_trn_provider_type ON training_providers (type);

-- -----------------------------------------------------------------------------
-- 2. training_trainers — trainer internal (employee) / external (provider)
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
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_trainer_emp ON training_trainers (employee_id);
CREATE INDEX IF NOT EXISTS idx_trn_trainer_provider ON training_trainers (provider_id);

-- -----------------------------------------------------------------------------
-- 3. training_session_trainers — satu session bisa banyak trainer (MAIN/ASSISTANT)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_session_trainers (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    session_id CHAR(36)     NOT NULL,
    trainer_id CHAR(36)     NOT NULL,
    role       VARCHAR(20)  NOT NULL DEFAULT 'MAIN',
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_sess_trn_session ON training_session_trainers (session_id);
CREATE INDEX IF NOT EXISTS idx_trn_sess_trn_trainer ON training_session_trainers (trainer_id);

-- -----------------------------------------------------------------------------
-- 4. ALTER training_courses — course_type, delivery_type (preferred), is_mandatory
-- -----------------------------------------------------------------------------
ALTER TABLE training_courses ADD COLUMN IF NOT EXISTS course_type VARCHAR(20) NULL;
ALTER TABLE training_courses ADD COLUMN IF NOT EXISTS delivery_type VARCHAR(20) NULL;
ALTER TABLE training_courses ADD COLUMN IF NOT EXISTS is_mandatory BOOLEAN NOT NULL DEFAULT FALSE;

-- -----------------------------------------------------------------------------
-- 5. ALTER training_sessions — provider/delivery + datetime presisi + link
-- -----------------------------------------------------------------------------
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS provider_type VARCHAR(20) NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS delivery_mode VARCHAR(20) NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS provider_id CHAR(36) NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS start_datetime TIMESTAMP NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS end_datetime TIMESTAMP NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS meeting_url TEXT NULL;
ALTER TABLE training_sessions ADD COLUMN IF NOT EXISTS registration_deadline TIMESTAMP NULL;

CREATE INDEX IF NOT EXISTS idx_trn_sess_provider ON training_sessions (provider_id);

-- -----------------------------------------------------------------------------
-- 6. ALTER training_participants — enrollment (registration + completion)
-- -----------------------------------------------------------------------------
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS registration_status VARCHAR(20) NOT NULL DEFAULT 'REGISTERED';
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS registered_at TIMESTAMP NULL;
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP NULL;
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS completion_status VARCHAR(20) NOT NULL DEFAULT 'NOT_STARTED';
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS completion_date DATE NULL;
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS final_score DECIMAL(5,2) NULL;
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS passed BOOLEAN NULL;
ALTER TABLE training_participants ADD COLUMN IF NOT EXISTS remarks TEXT NULL;

-- Dedupe record aktif sebelum index unik (partial index mengabaikan soft-deleted
-- sehingga employee boleh diregistrasi ulang setelah record lama dihapus).
DELETE FROM training_participants a USING training_participants b
WHERE a.id > b.id
  AND a.session_id = b.session_id
  AND a.employee_id = b.employee_id
  AND a.deleted_at IS NULL
  AND b.deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_part_session_employee
    ON training_participants (session_id, employee_id) WHERE deleted_at IS NULL;

-- -----------------------------------------------------------------------------
-- 7. ALTER training_materials — deskripsi, wajib, ketersediaan
-- -----------------------------------------------------------------------------
ALTER TABLE training_materials ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE training_materials ADD COLUMN IF NOT EXISTS is_required BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE training_materials ADD COLUMN IF NOT EXISTS available_from TIMESTAMP NULL;

-- -----------------------------------------------------------------------------
-- 8. training_attendances — attendance detail per hari (multi-day session)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_attendances (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    participant_id  CHAR(36)     NOT NULL,
    attendance_date DATE         NOT NULL,
    check_in        TIMESTAMP    NULL,
    check_out       TIMESTAMP    NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'PRESENT',
    remarks         TEXT         NULL,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_att_part ON training_attendances (participant_id);

-- Satu baris per participant per hari (record aktif).
CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_att_part_date
    ON training_attendances (participant_id, attendance_date) WHERE deleted_at IS NULL;

-- -----------------------------------------------------------------------------
-- 9. training_assessments — definisi assessment per session
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_assessments (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    session_id    CHAR(36)     NOT NULL,
    name          VARCHAR(200) NOT NULL,
    type          VARCHAR(20)  NOT NULL DEFAULT 'OTHER',
    max_score     DECIMAL(8,2) NOT NULL DEFAULT 100.00,
    passing_score DECIMAL(8,2) NOT NULL DEFAULT 60.00,
    attempt_limit INT          NOT NULL DEFAULT 1,
    is_required   BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_assess_session ON training_assessments (session_id);

-- -----------------------------------------------------------------------------
-- 10. training_assessment_results — nilai per peserta per attempt
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_assessment_results (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    assessment_id  CHAR(36)     NOT NULL,
    participant_id CHAR(36)     NOT NULL,
    score          DECIMAL(8,2) NOT NULL DEFAULT 0.00,
    passed         BOOLEAN      NOT NULL DEFAULT FALSE,
    attempt        INT          NOT NULL DEFAULT 1,
    completed_at   TIMESTAMP    NULL,
    deleted_at     TIMESTAMP,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_assess_res_assessment ON training_assessment_results (assessment_id);
CREATE INDEX IF NOT EXISTS idx_trn_assess_res_part ON training_assessment_results (participant_id);

-- Maksimal satu result per (assessment, participant, attempt) untuk record aktif.
CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_assess_res_attempt
    ON training_assessment_results (assessment_id, participant_id, attempt) WHERE deleted_at IS NULL;
