-- 016_training.sql
-- Training & Development Management Module
-- Tabel untuk end-to-end training management: course catalog, session scheduling,
-- participant registration, attendance, materials, evaluations, and certificates

-- =========================================================================
-- Training Categories (Kategori Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_categories (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    code        VARCHAR(20)  NOT NULL,
    name        VARCHAR(150) NOT NULL,
    description VARCHAR(500) NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trn_cat_code ON training_categories (code);

-- =========================================================================
-- Training Courses (Master Data Kursus)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_courses (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    category_id     CHAR(36)     NOT NULL,
    code            VARCHAR(20)  NOT NULL,
    name            VARCHAR(200) NOT NULL,
    description     TEXT         NULL,
    duration_hour   DECIMAL(8,2) NULL,
    min_score       DECIMAL(5,2) NULL,
    cost            DECIMAL(18,2) NULL,
    is_certified    BOOLEAN      NOT NULL DEFAULT FALSE,
    external_vendor VARCHAR(200) NULL,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at      TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trn_course_code ON training_courses (code);

CREATE INDEX IF NOT EXISTS idx_trn_course_cat ON training_courses (category_id);

-- =========================================================================
-- Training Sessions (Sesi/Kelas Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_sessions (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    course_id     CHAR(36)     NOT NULL,
    session_code  VARCHAR(20)  NOT NULL,
    trainer_name  VARCHAR(200) NOT NULL,
    location      VARCHAR(255) NULL,
    start_date    DATE         NOT NULL,
    end_date      DATE         NOT NULL,
    max_quota     INT          NOT NULL DEFAULT 30,
    status        VARCHAR(20)  NOT NULL DEFAULT 'SCHEDULED',
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_sess_course ON training_sessions (course_id);

CREATE INDEX IF NOT EXISTS idx_trn_sess_status ON training_sessions (status);

-- =========================================================================
-- Training Participants (Peserta Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_participants (
    id                 CHAR(36)     NOT NULL PRIMARY KEY,
    session_id         CHAR(36)     NOT NULL,
    employee_id        CHAR(36)     NOT NULL,
    attendance_status  VARCHAR(20)  NOT NULL DEFAULT 'PRESENT',
    score              DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    completed_at       DATE         NULL,
    deleted_at         TIMESTAMP,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_part_sess ON training_participants (session_id);

CREATE INDEX IF NOT EXISTS idx_trn_part_emp ON training_participants (employee_id);

-- =========================================================================
-- Training Materials (Materi Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_materials (
    id         CHAR(36)     NOT NULL PRIMARY KEY,
    session_id CHAR(36)     NOT NULL,
    title      VARCHAR(200) NOT NULL,
    file_url   TEXT         NULL,
    file_type  VARCHAR(50)  NULL,
    sort_order SMALLINT     NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_mat_sess ON training_materials (session_id);

-- =========================================================================
-- Training Evaluations (Evaluasi Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_evaluations (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    session_id  CHAR(36)     NOT NULL,
    employee_id CHAR(36)     NOT NULL,
    rating      SMALLINT     NOT NULL,
    feedback    TEXT         NULL,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_eval_sess ON training_evaluations (session_id);

CREATE INDEX IF NOT EXISTS idx_trn_eval_emp ON training_evaluations (employee_id);

-- =========================================================================
-- Training Certificates (Sertifikat Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_certificates (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    participant_id CHAR(36)     NOT NULL,
    certificate_no VARCHAR(50)  NOT NULL,
    issued_date    DATE         NOT NULL,
    expiry_date    DATE         NULL,
    deleted_at     TIMESTAMP,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_trn_cert_no ON training_certificates (certificate_no);

CREATE INDEX IF NOT EXISTS idx_trn_cert_part ON training_certificates (participant_id);
