-- =============================================================================
-- Tenant Migration: 100_candidate_documents (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 3a: candidate_documents
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-documents-design.md)
--
-- candidate_documents — referensi dokumen kandidat (bukan binary). File
-- sesungguhnya diupload lewat endpoint generik POST /api/v1/tenant/uploads
-- (backend/internal/pkg/upload) yang mengembalikan URL; tabel ini hanya
-- menyimpan referensi URL tersebut. document_type enum di-enforce di layer
-- Gin binding (oneof=...), bukan DB constraint — pola sama dengan
-- CandidateType/OfferStatus di modul ini.
--
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
);

CREATE INDEX IF NOT EXISTS idx_cand_doc_candidate ON candidate_documents (candidate_id);
