-- =============================================================================
-- Tenant Migration: 108_application_assessments (MySQL)
-- =============================================================================
-- Penilaian kandidat per aplikasi (G-12 susulan — "Penilaian Kandidat"):
--   education_match / experience_match  → penilai menandai pendidikan/pengalaman
--                                         kandidat sesuai requirement atau tidak
--   education_note / experience_note    → catatan penilai
--   competency_levels (JSON)            → [{competency_id, level}] — level
--                                         kompetensi kandidat per requirement
--   score + breakdown (JSON)            → skor akhir (Pendidikan 20%,
--                                         Pengalaman 30%, Kompetensi 50%)
-- Satu baris per application (UNIQUE application_id) — upsert saat dinilai ulang.
-- Idempotent: CREATE TABLE IF NOT EXISTS.
--
-- See postgres version for full column documentation.

CREATE TABLE IF NOT EXISTS application_assessments (
    id               CHAR(36) PRIMARY KEY,
    application_id   CHAR(36) NOT NULL,
    education_match  TINYINT(1) NULL,
    education_note   TEXT NULL,
    experience_match TINYINT(1) NULL,
    experience_note  TEXT NULL,
    competency_levels JSON NULL,
    score            DECIMAL(5,2) NULL,
    breakdown        JSON NULL,
    assessed_by      CHAR(36) NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_appassess_app FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT uq_appassess_app UNIQUE (application_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_appassess_app ON application_assessments (application_id);
