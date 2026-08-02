-- Migration: 030_add_ref_to_job_management_values
-- Menambahkan kolom ref_id & ref_type ke tabel job_management_values
-- (relasi polymorphic: ref_type = entitas yang direferensikan, ref_id = id record-nya).
--
-- CATATAN (PostgreSQL): Migration ini hanya dibutuhkan untuk tenant lama,
-- jadi gunakan ADD COLUMN IF NOT EXISTS (idempotent).

ALTER TABLE job_management_values
    ADD COLUMN IF NOT EXISTS ref_id CHAR(36) NULL;

ALTER TABLE job_management_values
    ADD COLUMN IF NOT EXISTS ref_type VARCHAR(100) NULL;

CREATE INDEX IF NOT EXISTS idx_jmv_ref
    ON job_management_values (ref_id);
