-- 018_career_intelligence.sql
-- Career Intelligence & Talent Management Module
-- Tabel untuk 9-box talent mapping, career interests, career paths, dan succession planning

-- =========================================================================
-- Career Talent Maps (9-Box Talent Mapping Grid)
-- =========================================================================
CREATE TABLE IF NOT EXISTS career_talent_maps (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id       CHAR(36)     NOT NULL COMMENT 'Employee being assessed',
    period            CHAR(7)      NOT NULL COMMENT 'Format: 2026-Q1',
    performance       VARCHAR(20)  NOT NULL COMMENT 'LOW / MEDIUM / HIGH',
    potential         VARCHAR(20)  NOT NULL COMMENT 'LOW / MEDIUM / HIGH',
    grid_position     VARCHAR(30)  NOT NULL COMMENT '9-BOX-1 through 9-BOX-9',
    notes             TEXT         NULL COMMENT 'Assessment notes',
    assessor_id       CHAR(36)     NULL COMMENT 'ID of the assessor (manager/HR)',
    assessed_at       DATE         NOT NULL COMMENT 'Date of assessment',
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
    employee_id       CHAR(36)     NOT NULL COMMENT 'Employee expressing interest',
    interest_type     VARCHAR(50)  NOT NULL COMMENT 'LEADERSHIP / SPECIALIST / INTERNATIONAL / ENTREPRENEUR',
    target_position   VARCHAR(100) NULL COMMENT 'Desired position title',
    target_department VARCHAR(100) NULL COMMENT 'Desired department',
    motivation        TEXT         NULL COMMENT 'Career motivation or reason',
    readiness_level   VARCHAR(20)  NULL COMMENT 'NOW / 1_YEAR / 2_3_YEARS / 3_PLUS',
    is_active         BOOLEAN      NOT NULL DEFAULT TRUE,
    recorded_at       DATE         NOT NULL COMMENT 'Date interest was recorded',
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
    source_title_id   CHAR(36)     NOT NULL COMMENT 'Source position title ID',
    target_title_id   CHAR(36)     NOT NULL COMMENT 'Target position title ID',
    path_type         VARCHAR(30)  NOT NULL COMMENT 'PROMOTION / LATERAL / DEMOTION / CROSSFUNCTIONAL',
    typical_tenure    INT          NOT NULL DEFAULT 0 COMMENT 'Typical tenure in months',
    requirements      TEXT         NULL COMMENT 'Requirements for this career path',
    competencies      TEXT         NULL COMMENT 'Required competencies (JSON list)',
    certifications    TEXT         NULL COMMENT 'Required certifications (JSON list)',
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
    position_id       CHAR(36)     NOT NULL COMMENT 'Key position planned for succession',
    successor_id      CHAR(36)     NOT NULL COMMENT 'Employee ID of potential successor',
    readiness_level   VARCHAR(20)  NOT NULL COMMENT 'READY_NOW / READY_1YR / READY_2YR / POTENTIAL',
    priority_order    INT          NOT NULL DEFAULT 1 COMMENT 'Priority among multiple successors',
    target_date       DATE         NULL COMMENT 'Target date for succession',
    development_plan  TEXT         NULL COMMENT 'Development plan for the successor',
    notes             TEXT         NULL COMMENT 'Additional notes',
    status            VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE / COMPLETED / REMOVED',
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    INDEX idx_csp_position (position_id),
    INDEX idx_csp_successor (successor_id),
    INDEX idx_csp_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
