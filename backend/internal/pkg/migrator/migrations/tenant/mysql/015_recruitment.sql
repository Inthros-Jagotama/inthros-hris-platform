-- 015_recruitment.sql
-- Recruitment & Onboarding (ATS) Module
-- Tabel untuk end-to-end recruitment & applicant tracking system

-- =========================================================================
-- Job Requisitions (Lowongan Pekerjaan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS job_requisitions (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    organization_id   CHAR(36)     NOT NULL,
    title             VARCHAR(255) NOT NULL,
    department        VARCHAR(150) NULL,
    employment_type   VARCHAR(50)  NULL,
    location          VARCHAR(255) NULL,
    min_salary        DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    max_salary        DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    description       TEXT         NULL,
    requirements      TEXT         NULL,
    responsibilities  TEXT         NULL,
    slots_available   INT          NOT NULL DEFAULT 1,
    slots_filled      INT          NOT NULL DEFAULT 0,
    status            VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    requested_by      CHAR(36)     NULL,
    approved_by       CHAR(36)     NULL,
    target_start_date DATE         NULL,
    closed_at         BIGINT       NULL DEFAULT 0,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_req_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Candidates (Kandidat Pelamar)
-- =========================================================================
CREATE TABLE IF NOT EXISTS candidates (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    phone           VARCHAR(50)  NULL,
    address         TEXT         NULL,
    current_company VARCHAR(255) NULL,
    current_title   VARCHAR(255) NULL,
    resume_url      TEXT         NULL,
    portfolio_url   TEXT         NULL,
    linkedin_url    TEXT         NULL,
    source          VARCHAR(50)  NOT NULL DEFAULT 'direct',
    notes           TEXT         NULL,
    created_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_cand_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Job Applications (Lamaran Pekerjaan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS job_applications (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    requisition_id    CHAR(36)     NOT NULL,
    candidate_id      CHAR(36)     NOT NULL,
    status            VARCHAR(50)  NOT NULL DEFAULT 'NEW',
    applied_at        BIGINT       NOT NULL DEFAULT 0,
    screened_at       BIGINT       NULL DEFAULT 0,
    shortlisted_at    BIGINT       NULL DEFAULT 0,
    offered_at        BIGINT       NULL DEFAULT 0,
    accepted_at       BIGINT       NULL DEFAULT 0,
    rejected_at       BIGINT       NULL DEFAULT 0,
    withdrawn_at      BIGINT       NULL DEFAULT 0,
    rejection_reason  TEXT         NULL,
    notes             TEXT         NULL,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_app_req (requisition_id),
    INDEX idx_app_cand (candidate_id),
    INDEX idx_app_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Interviews (Wawancara)
-- =========================================================================
CREATE TABLE IF NOT EXISTS interviews (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    application_id    CHAR(36)     NOT NULL,
    interviewer_id    CHAR(36)     NOT NULL,
    stage             VARCHAR(50)  NOT NULL,
    scheduled_at      BIGINT       NOT NULL DEFAULT 0,
    duration_minutes  INT          NOT NULL DEFAULT 60,
    location          VARCHAR(255) NULL,
    meeting_link      TEXT         NULL,
    status            VARCHAR(20)  NOT NULL DEFAULT 'SCHEDULED',
    score             DECIMAL(5,2) NULL,
    feedback          TEXT         NULL,
    completed_at      BIGINT       NULL DEFAULT 0,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_int_app (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Onboarding Task Templates (Template Tugas Onboarding)
-- =========================================================================
CREATE TABLE IF NOT EXISTS onboarding_task_templates (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    description   TEXT         NULL,
    category      VARCHAR(50)  NULL,
    day_offset    INT          NOT NULL DEFAULT 0,
    assigned_role VARCHAR(50)  NULL,
    is_mandatory  TINYINT(1)   NOT NULL DEFAULT 1,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Employee Onboardings (Onboarding Karyawan Baru)
-- =========================================================================
CREATE TABLE IF NOT EXISTS employee_onboardings (
    id              CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id     CHAR(36)     NOT NULL,
    application_id  CHAR(36)     NOT NULL,
    start_date      DATE         NOT NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'PENDING',
    buddy_id        CHAR(36)     NULL,
    completed_at    BIGINT       NULL DEFAULT 0,
    notes           TEXT         NULL,
    created_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_onb_emp (employee_id),
    INDEX idx_onb_app (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Onboarding Task Items (Tugas spesifik dalam onboarding)
-- =========================================================================
CREATE TABLE IF NOT EXISTS onboarding_task_items (
    id                      CHAR(36)     NOT NULL PRIMARY KEY,
    employee_onboarding_id  CHAR(36)     NOT NULL,
    template_id             CHAR(36)     NULL,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT         NULL,
    assigned_to             CHAR(36)     NULL,
    due_date                BIGINT       NULL DEFAULT 0,
    is_completed            TINYINT(1)   NOT NULL DEFAULT 0,
    completed_at            BIGINT       NULL DEFAULT 0,
    notes                   TEXT         NULL,
    created_at              TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at              TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_onb_task_item (employee_onboarding_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
