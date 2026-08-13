-- =============================================================================
-- Tenant Migration: 105_recruitment_requisition_requirements (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS job_requisition_requirements (
    id               CHAR(36) PRIMARY KEY,
    requisition_id   CHAR(36) NOT NULL,
    requirement_type VARCHAR(50) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    description      TEXT NULL,
    minimum_value    DECIMAL(10,2) NULL,
    maximum_value    DECIMAL(10,2) NULL,
    is_required      TINYINT(1) NOT NULL DEFAULT 1,
    sort_order       INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reqmt_req FOREIGN KEY (requisition_id) REFERENCES job_requisitions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_reqmt_req ON job_requisition_requirements (requisition_id);

CREATE TABLE IF NOT EXISTS job_requisition_competencies (
    id             CHAR(36) PRIMARY KEY,
    requisition_id CHAR(36) NOT NULL,
    competency_id  CHAR(36) NOT NULL,
    required_level SMALLINT NULL,
    is_required    TINYINT(1) NOT NULL DEFAULT 1,
    weight         DECIMAL(5,2) NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reqcomp_req FOREIGN KEY (requisition_id) REFERENCES job_requisitions(id) ON DELETE CASCADE,
    CONSTRAINT fk_reqcomp_comp FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE,
    CONSTRAINT uq_reqcomp_req_comp UNIQUE (requisition_id, competency_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_reqcomp_req ON job_requisition_competencies (requisition_id);
CREATE INDEX idx_reqcomp_comp ON job_requisition_competencies (competency_id);
