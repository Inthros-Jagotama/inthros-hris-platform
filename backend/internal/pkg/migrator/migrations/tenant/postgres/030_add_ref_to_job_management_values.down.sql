-- Down Migration: 030_add_ref_to_job_management_values
-- Hapus kolom ref_id & ref_type dari tabel job_management_values.

DROP INDEX IF EXISTS idx_jmv_ref;

ALTER TABLE job_management_values
    DROP COLUMN IF EXISTS ref_type;

ALTER TABLE job_management_values
    DROP COLUMN IF EXISTS ref_id;
