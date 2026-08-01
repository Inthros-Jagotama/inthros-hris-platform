-- Migration: 025_add_education_major_to_employee_educations
-- Menambahkan kolom education_major_id ke tabel employee_educations
-- (relasi ke tabel master education_majors dari module setting).
--
-- CATATAN (PostgreSQL): Migration ini hanya dibutuhkan untuk tenant lama,
-- jadi gunakan ADD COLUMN IF NOT EXISTS (idempotent).

ALTER TABLE employee_educations
    ADD COLUMN IF NOT EXISTS education_major_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_empedu_education_major
    ON employee_educations (education_major_id);
