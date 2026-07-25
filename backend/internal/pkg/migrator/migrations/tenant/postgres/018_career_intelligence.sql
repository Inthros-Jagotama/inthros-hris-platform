-- 018_career_intelligence.sql
-- Career Intelligence & Talent Management Module
-- Tabel untuk 9-box talent mapping, career interests, career paths, dan succession planning

-- =========================================================================
-- Career Talent Maps (9-Box Talent Mapping Grid)
-- =========================================================================
CREATE TABLE IF NOT EXISTS career_talent_maps (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id       CHAR(36)     NOT NULL,
    period            CHAR(7)      NOT NULL,
    performance       VARCHAR(20)  NOT NULL,
    potential         VARCHAR(20)  NOT NULL,
    grid_position     VARCHAR(30)  NOT NULL,
    notes             TEXT         NULL,
    assessor_id       CHAR(36)     NULL,
    assessed_at       DATE         NOT NULL,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    INDEX idx_ctm_employee (employee_id),
    INDEX idx_ctm_period (period),
    INDEX idx_ctm_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Career Interests (Employee Career Aspirations)
-- =========================================================================
CREATE TABLE IF NOT EXISTS career_interests (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id       CHAR(36)     NOT NULL,
    interest_type     VARCHAR(50)  NOT NULL,
    target_position   VARCHAR(100) NULL,
    target_department VARCHAR(100) NULL,
    motivation        TEXT         NULL,
    readiness_level   VARCHAR(20)  NULL,
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    recorded_at       DATE         NOT NULL,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    INDEX idx_ci_employee (employee_id),
    INDEX idx_ci_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Career Paths (Career Progression Routes)
-- =========================================================================
CREATE TABLE IF NOT EXISTS career_paths (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    source_title_id   CHAR(36)     NOT NULL,
    target_title_id   CHAR(36)     NOT NULL,
    path_type         VARCHAR(30)  NOT NULL,
    typical_tenure    INT          NOT NULL DEFAULT 0,
    requirements      TEXT         NULL,
    competencies      TEXT         NULL,
    certifications    TEXT         NULL,
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    INDEX idx_cp_source (source_title_id),
    INDEX idx_cp_target (target_title_id),
    INDEX idx_cp_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Career Succession Plans (Succession Planning for Key Positions)
-- =========================================================================
CREATE TABLE IF NOT EXISTS career_succession_plans (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    position_id       CHAR(36)     NOT NULL,
    successor_id      CHAR(36)     NOT NULL,
    readiness_level   VARCHAR(20)  NOT NULL,
    priority_order    INT          NOT NULL DEFAULT 1,
    target_date       DATE         NULL,
    development_plan  TEXT         NULL,
    notes             TEXT         NULL,
    status            VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE',
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    INDEX idx_csp_position (position_id),
    INDEX idx_csp_successor (successor_id),
    INDEX idx_csp_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
