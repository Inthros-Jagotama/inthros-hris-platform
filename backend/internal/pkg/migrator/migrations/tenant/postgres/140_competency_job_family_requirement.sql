-- =============================================================================
-- 140_competency_job_family_requirement.sql
-- Competency 360 Module — Job Family Competency Requirement (plan generik §2.1/§21).
-- `job_family_competencies` (migration 002) selama ini hanya (job_family_id,
-- competency_id) tanpa required level/weight dan belum dipakai kode apapun.
-- Tambahkan kolom requirement agar job family dapat menentukan level & bobot
-- competency yang dibutuhkan untuk sebuah posisi.
-- =============================================================================

ALTER TABLE job_family_competencies
    ADD COLUMN IF NOT EXISTS required_level SMALLINT NULL,
    ADD COLUMN IF NOT EXISTS weight         DECIMAL(6, 2) NOT NULL DEFAULT 1;
