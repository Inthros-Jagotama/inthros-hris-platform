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
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_trn_cat_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    deleted_at      TIMESTAMP(6) NULL DEFAULT NULL,
    created_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_trn_course_code (code),
    INDEX idx_trn_course_cat (category_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    deleted_at    TIMESTAMP(6) NULL DEFAULT NULL,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_sess_course (course_id),
    INDEX idx_trn_sess_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    deleted_at         TIMESTAMP(6) NULL DEFAULT NULL,
    created_at         TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at         TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_part_sess (session_id),
    INDEX idx_trn_part_emp (employee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    deleted_at TIMESTAMP(6) NULL DEFAULT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_mat_sess (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Training Evaluations (Evaluasi Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_evaluations (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    session_id  CHAR(36)     NOT NULL,
    employee_id CHAR(36)     NOT NULL,
    rating      SMALLINT     NOT NULL,
    feedback    TEXT         NULL,
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_eval_sess (session_id),
    INDEX idx_trn_eval_emp (employee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Training Certificates (Sertifikat Pelatihan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS training_certificates (
    id             CHAR(36)     NOT NULL PRIMARY KEY,
    participant_id CHAR(36)     NOT NULL,
    certificate_no VARCHAR(50)  NOT NULL,
    issued_date    DATE         NOT NULL,
    expiry_date    DATE         NULL,
    deleted_at     TIMESTAMP(6) NULL DEFAULT NULL,
    created_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at     TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_trn_cert_no (certificate_no),
    INDEX idx_trn_cert_part (participant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
