-- =============================================================================
-- Tenant Migration: 099_candidate_skills_certifications (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_skills (
    id            CHAR(36) PRIMARY KEY,
    candidate_id  CHAR(36) NOT NULL,
    competency_id CHAR(36) NOT NULL,
    level         SMALLINT NULL,
    notes         TEXT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_skill_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_skill_competency FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_skill_candidate ON candidate_skills (candidate_id);
CREATE INDEX idx_cand_skill_competency ON candidate_skills (competency_id);

CREATE TABLE IF NOT EXISTS candidate_certifications (
    id                   CHAR(36) PRIMARY KEY,
    candidate_id         CHAR(36) NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    issuing_organization VARCHAR(255) NULL,
    issue_date           DATE NULL,
    expiry_date          DATE NULL,
    credential_url       TEXT NULL,
    notes                TEXT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_cert_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_cert_candidate ON candidate_certifications (candidate_id);
