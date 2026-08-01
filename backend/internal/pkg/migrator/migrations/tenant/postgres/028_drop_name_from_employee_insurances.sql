-- Migration: 028_drop_name_from_employee_insurances
-- Menghapus kolom name dari tabel employee_insurances.
-- Nama asuransi kini diambil dari relasi ke tabel master insurances.
--
-- CATATAN (PostgreSQL): Idempotent via DROP COLUMN IF EXISTS.

ALTER TABLE employee_insurances
    DROP COLUMN IF EXISTS name;
