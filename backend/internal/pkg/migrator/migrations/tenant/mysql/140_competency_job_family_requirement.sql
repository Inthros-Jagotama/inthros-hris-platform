-- =============================================================================
-- 140_competency_job_family_requirement.sql
-- Competency 360 Module — Job Family Competency Requirement (plan generik §2.1/§21).
-- `job_family_competencies` (migration 002) selama ini hanya (job_family_id,
-- competency_id) tanpa required level/weight dan belum dipakai kode apapun.
-- Tambahkan kolom requirement agar job family dapat menentukan level & bobot
-- competency yang dibutuhkan untuk sebuah posisi.
-- =============================================================================

ALTER TABLE job_family_competencies
    ADD COLUMN required_level SMALLINT NULL AFTER competency_id,
    ADD COLUMN weight         DECIMAL(6, 2) NOT NULL DEFAULT 1 AFTER required_level;
