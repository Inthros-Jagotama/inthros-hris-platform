-- =============================================================================
-- Tenant Migration: 103_recruitment_assessment (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS recruitment_assessments (
    id             CHAR(36) PRIMARY KEY,
    requisition_id CHAR(36) NULL,
    name           VARCHAR(255) NOT NULL,
    type           VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    scheduled_at   BIGINT NOT NULL DEFAULT 0,
    location       VARCHAR(255) NULL,
    meeting_link   TEXT NULL,
    notes          TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_assess_req FOREIGN KEY (requisition_id) REFERENCES job_requisitions(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_assess_req ON recruitment_assessments (requisition_id);

CREATE TABLE IF NOT EXISTS assessment_participants (
    id             CHAR(36) PRIMARY KEY,
    assessment_id  CHAR(36) NOT NULL,
    application_id CHAR(36) NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'INVITED',
    score          DECIMAL(5,2) NULL,
    result         VARCHAR(10) NULL,
    recommendation TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_partic_assess FOREIGN KEY (assessment_id) REFERENCES recruitment_assessments(id) ON DELETE CASCADE,
    CONSTRAINT fk_partic_app FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT uq_partic_assess_app UNIQUE (assessment_id, application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_partic_assess ON assessment_participants (assessment_id);
CREATE INDEX idx_partic_app ON assessment_participants (application_id);
