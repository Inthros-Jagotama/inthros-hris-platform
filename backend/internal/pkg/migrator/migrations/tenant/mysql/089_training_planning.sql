-- =============================================================================
-- Tenant Migration: 089_training_planning (MySQL)
-- =============================================================================
-- Training & Development P1-BE (docs/module-training-development-plan.md §42 P1-BE):
--   - training_plans + training_plan_items
--   - training_needs
--   - training_requests (approval via Central Approval Engine, module `training_request`)
--   - course objectives / competencies / prerequisites (sub-resource course)
--   - training_mandatories (mandatory training compliance)
--   - training_session_costs + training_documents
--
-- Semua statement idempotent (CREATE TABLE IF NOT EXISTS).

-- -----------------------------------------------------------------------------
-- 1. training_plans
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_plans (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    code        VARCHAR(30)  NOT NULL,
    name        VARCHAR(200) NOT NULL,
    year        INT          NOT NULL,
    description TEXT         NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX uk_trn_plan_code (code),
    INDEX idx_trn_plan_year (year)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 2. training_plan_items
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_plan_items (
    id                  CHAR(36)      NOT NULL PRIMARY KEY,
    training_plan_id    CHAR(36)      NOT NULL,
    course_id           CHAR(36)      NOT NULL,
    target_date         DATE          NULL,
    target_participants INT           NULL,
    estimated_cost      DECIMAL(14,2) NULL,
    priority            VARCHAR(20)   NOT NULL DEFAULT 'MEDIUM',
    deleted_at          TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at          TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at          TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_plan_item_plan (training_plan_id),
    INDEX idx_trn_plan_item_course (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 3. training_needs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_needs (
    id               CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id      CHAR(36)     NULL,
    organization_id  CHAR(36)     NULL,
    position_id      CHAR(36)     NULL,
    course_id        CHAR(36)     NULL,
    reason           TEXT         NULL,
    priority         VARCHAR(20)  NOT NULL DEFAULT 'MEDIUM',
    source_type      VARCHAR(30)  NOT NULL DEFAULT 'MANUAL',
    source_id        CHAR(36)     NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'OPEN',
    deleted_at       TIMESTAMP(6) NULL DEFAULT NULL,
    created_at       TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at       TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_need_employee (employee_id),
    INDEX idx_trn_need_org (organization_id),
    INDEX idx_trn_need_course (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 4. training_requests
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_requests (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id           CHAR(36)     NOT NULL,
    course_id             CHAR(36)     NOT NULL,
    session_id            CHAR(36)     NULL,
    requested_date        DATE         NOT NULL,
    reason                TEXT         NULL,
    priority              VARCHAR(20)  NOT NULL DEFAULT 'MEDIUM',
    status                VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    approval_instance_id  CHAR(36)     NULL,
    approved_at           TIMESTAMP(6) NULL,
    rejected_at           TIMESTAMP(6) NULL,
    supervisor_note       TEXT         NULL,
    deleted_at            TIMESTAMP(6) NULL DEFAULT NULL,
    created_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_req_employee (employee_id),
    INDEX idx_trn_req_course (course_id),
    INDEX idx_trn_req_status (status),
    INDEX idx_trn_req_approval_instance (approval_instance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 5. training_course_objectives
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_objectives (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    course_id   CHAR(36)     NOT NULL,
    objective   TEXT         NOT NULL,
    sort_order  INT          NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_obj_course (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 6. training_course_competencies
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_competencies (
    id            CHAR(36)    NOT NULL PRIMARY KEY,
    course_id     CHAR(36)    NOT NULL,
    competency_id CHAR(36)    NOT NULL,
    target_level  INT         NULL,
    deleted_at    TIMESTAMP(6) NULL DEFAULT NULL,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_course_comp_course (course_id),
    INDEX idx_trn_course_comp_comp (competency_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 7. training_course_prerequisites
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_prerequisites (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    course_id           CHAR(36)     NOT NULL,
    prerequisite_type   VARCHAR(20)  NOT NULL DEFAULT 'COURSE',
    prerequisite_id     CHAR(36)     NULL,
    is_required         TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at          TIMESTAMP(6) NULL DEFAULT NULL,
    created_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_preq_course (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 8. training_mandatories
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_mandatories (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    course_id             CHAR(36)     NOT NULL,
    organization_id       CHAR(36)     NULL,
    position_id           CHAR(36)     NULL,
    employment_status_id  CHAR(36)     NULL,
    due_days              INT          NULL,
    validity_period_month INT          NULL,
    is_active             TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at            TIMESTAMP(6) NULL DEFAULT NULL,
    created_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_mand_course (course_id),
    INDEX idx_trn_mand_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 9. training_session_costs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_session_costs (
    id           CHAR(36)      NOT NULL PRIMARY KEY,
    session_id   CHAR(36)      NOT NULL,
    cost_type    VARCHAR(30)   NOT NULL DEFAULT 'OTHER',
    description  TEXT          NULL,
    amount       DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    deleted_at   TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at   TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at   TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_cost_session (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------------------------
-- 10. training_documents
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_documents (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    session_id    CHAR(36)     NOT NULL,
    document_type VARCHAR(30)  NOT NULL DEFAULT 'OTHER',
    file_name     VARCHAR(255) NULL,
    file_url      TEXT         NULL,
    uploaded_by   CHAR(36)     NULL,
    deleted_at    TIMESTAMP(6) NULL DEFAULT NULL,
    created_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at    TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_trn_doc_session (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
