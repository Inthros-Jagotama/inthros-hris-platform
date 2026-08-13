-- =============================================================================
-- Tenant Migration: 100_candidate_documents (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_documents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    document_type  VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    name           VARCHAR(255) NOT NULL,
    file_url       TEXT NOT NULL,
    notes          TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_doc_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_doc_candidate ON candidate_documents (candidate_id);
