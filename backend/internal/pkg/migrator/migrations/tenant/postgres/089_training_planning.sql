-- =============================================================================
-- Tenant Migration: 089_training_planning (PostgreSQL)
-- =============================================================================
-- Training & Development P1-BE (docs/module-training-development-plan.md §42 P1-BE):
--   - training_plans + training_plan_items
--   - training_needs
--   - training_requests (approval via Central Approval Engine, module `training_request`)
--   - course objectives / competencies / prerequisites (sub-resource course)
--   - training_mandatories (mandatory training compliance)
--   - training_session_costs + training_documents
--
-- Semua statement idempotent (IF NOT EXISTS).

-- -----------------------------------------------------------------------------
-- 1. training_plans — rencana pelatihan tahunan
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_plans (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    code        VARCHAR(30)  NOT NULL,
    name        VARCHAR(200) NOT NULL,
    year        INT          NOT NULL,
    description TEXT         NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'DRAFT',
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_trn_plan_code ON training_plans (code);
CREATE INDEX IF NOT EXISTS idx_trn_plan_year ON training_plans (year);

-- -----------------------------------------------------------------------------
-- 2. training_plan_items
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_plan_items (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    training_plan_id    CHAR(36)     NOT NULL,
    course_id           CHAR(36)     NOT NULL,
    target_date         DATE         NULL,
    target_participants INT          NULL,
    estimated_cost      DECIMAL(14,2) NULL,
    priority            VARCHAR(20)  NOT NULL DEFAULT 'MEDIUM',
    deleted_at          TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_plan_item_plan ON training_plan_items (training_plan_id);
CREATE INDEX IF NOT EXISTS idx_trn_plan_item_course ON training_plan_items (course_id);

-- -----------------------------------------------------------------------------
-- 3. training_needs — kebutuhan training (operasional, bukan Intelligence)
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
    deleted_at       TIMESTAMP,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_need_employee ON training_needs (employee_id);
CREATE INDEX IF NOT EXISTS idx_trn_need_org ON training_needs (organization_id);
CREATE INDEX IF NOT EXISTS idx_trn_need_course ON training_needs (course_id);

-- -----------------------------------------------------------------------------
-- 4. training_requests — pengajuan training via Central Approval
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
    approved_at           TIMESTAMP    NULL,
    rejected_at           TIMESTAMP    NULL,
    supervisor_note       TEXT         NULL,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_req_employee ON training_requests (employee_id);
CREATE INDEX IF NOT EXISTS idx_trn_req_course ON training_requests (course_id);
CREATE INDEX IF NOT EXISTS idx_trn_req_status ON training_requests (status);
CREATE INDEX IF NOT EXISTS idx_trn_req_approval_instance ON training_requests (approval_instance_id);

-- -----------------------------------------------------------------------------
-- 5. training_course_objectives — learning objectives per course
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_objectives (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    course_id   CHAR(36)     NOT NULL,
    objective   TEXT         NOT NULL,
    sort_order  INT          NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_obj_course ON training_course_objectives (course_id);

-- -----------------------------------------------------------------------------
-- 6. training_course_competencies — relasi course ↔ competencies.id
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_competencies (
    id            CHAR(36)    NOT NULL PRIMARY KEY,
    course_id     CHAR(36)    NOT NULL,
    competency_id CHAR(36)    NOT NULL,
    target_level  INT         NULL,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_course_comp_course ON training_course_competencies (course_id);
CREATE INDEX IF NOT EXISTS idx_trn_course_comp_comp ON training_course_competencies (competency_id);

-- -----------------------------------------------------------------------------
-- 7. training_course_prerequisites
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_course_prerequisites (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    course_id           CHAR(36)     NOT NULL,
    prerequisite_type   VARCHAR(20)  NOT NULL DEFAULT 'COURSE',
    prerequisite_id     CHAR(36)     NULL,
    is_required         BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at          TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_preq_course ON training_course_prerequisites (course_id);

-- -----------------------------------------------------------------------------
-- 8. training_mandatories — mandatory training by target
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_mandatories (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    course_id             CHAR(36)     NOT NULL,
    organization_id       CHAR(36)     NULL,
    position_id           CHAR(36)     NULL,
    employment_status_id  CHAR(36)     NULL,
    due_days              INT          NULL,
    validity_period_month INT          NULL,
    is_active             BOOLEAN      NOT NULL DEFAULT TRUE,
    deleted_at            TIMESTAMP,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_mand_course ON training_mandatories (course_id);
CREATE INDEX IF NOT EXISTS idx_trn_mand_org ON training_mandatories (organization_id);

-- -----------------------------------------------------------------------------
-- 9. training_session_costs — biaya aktual per session
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_session_costs (
    id           CHAR(36)      NOT NULL PRIMARY KEY,
    session_id   CHAR(36)      NOT NULL,
    cost_type    VARCHAR(30)   NOT NULL DEFAULT 'OTHER',
    description  TEXT          NULL,
    amount       DECIMAL(14,2) NOT NULL DEFAULT 0.00,
    deleted_at   TIMESTAMP,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_cost_session ON training_session_costs (session_id);

-- -----------------------------------------------------------------------------
-- 10. training_documents — dokumen pendukung session
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_documents (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    session_id    CHAR(36)     NOT NULL,
    document_type VARCHAR(30)  NOT NULL DEFAULT 'OTHER',
    file_name     VARCHAR(255) NULL,
    file_url      TEXT         NULL,
    uploaded_by   CHAR(36)     NULL,
    deleted_at    TIMESTAMP,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_trn_doc_session ON training_documents (session_id);
