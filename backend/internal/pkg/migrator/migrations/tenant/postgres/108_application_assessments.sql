-- =============================================================================
-- Tenant Migration: 108_application_assessments (PostgreSQL)
-- =============================================================================
-- Penilaian kandidat per aplikasi (G-12 susulan — "Penilaian Kandidat"):
--   education_match / experience_match  → penilai menandai pendidikan/pengalaman
--                                         kandidat sesuai requirement atau tidak
--   education_note / experience_note    → catatan penilai
--   competency_levels (JSONB)           → [{competency_id, level}] — level
--                                         kompetensi kandidat per requirement
--   score + breakdown (JSONB)           → skor akhir (Pendidikan 20%,
--                                         Pengalaman 30%, Kompetensi 50%)
-- Satu baris per application (UNIQUE application_id) — upsert saat dinilai ulang.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS application_assessments (
    id               CHAR(36) PRIMARY KEY,
    application_id   CHAR(36) NOT NULL,
    education_match  BOOLEAN NULL,
    education_note   TEXT NULL,
    experience_match BOOLEAN NULL,
    experience_note  TEXT NULL,
    competency_levels JSONB NULL,
    score            NUMERIC(5,2) NULL,
    breakdown        JSONB NULL,
    assessed_by      CHAR(36) NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_appassess_app FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT uq_appassess_app UNIQUE (application_id)
);

CREATE INDEX IF NOT EXISTS idx_appassess_app ON application_assessments (application_id);
