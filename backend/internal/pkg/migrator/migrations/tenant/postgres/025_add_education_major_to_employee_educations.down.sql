-- Down Migration: 025_add_education_major_to_employee_educations
-- Hapus kolom education_major_id dari tabel employee_educations (idempotent).

DROP INDEX IF EXISTS idx_empedu_education_major;

ALTER TABLE employee_educations
    DROP COLUMN IF EXISTS education_major_id;
