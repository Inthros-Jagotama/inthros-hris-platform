-- =============================================================================
-- Tenant Migration: 101_candidate_consents (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_consents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    action         VARCHAR(20) NOT NULL,
    notes          TEXT NULL,
    changed_by     CHAR(36) NULL,
    changed_at     BIGINT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_consent_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_consent_candidate ON candidate_consents (candidate_id);
