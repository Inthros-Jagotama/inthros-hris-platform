-- =============================================================================
-- Tenant Migration: 083_employeemovement_snapshot
-- =============================================================================
-- Employee Movement Enhancement (plan §12.5): Movement Snapshot.
--
-- Menyimpan nama Organization / Position / Employment Status pada saat
-- movement dibuat, sehingga histori movement tidak berubah ketika master
-- data (nomenclature organisasi, title posisi, nama status) diubah namanya.
--
-- Kolom baru (semua nullable — movement lama tetap valid tanpa snapshot):
--   from_organization_name
--   from_position_name
--   from_employment_status_name
--   to_organization_name
--   to_position_name
--   to_employment_status_name
--
-- Foreign key (from_*/to_*_id) tetap dipertahankan untuk relasi & navigasi.

ALTER TABLE employee_movements
    ADD COLUMN IF NOT EXISTS from_organization_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS from_position_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS from_employment_status_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS to_organization_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS to_position_name VARCHAR(255) NULL,
    ADD COLUMN IF NOT EXISTS to_employment_status_name VARCHAR(255) NULL;
