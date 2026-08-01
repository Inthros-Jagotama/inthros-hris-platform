-- Migration: 028_drop_name_from_employee_insurances
-- Menghapus kolom name dari tabel employee_insurances.
-- Nama asuransi kini diambil dari relasi ke tabel master insurances
-- (employee_insurances.insurance_id → insurances.id → insurances.name).
--
-- CATATAN (MySQL): Kondisional/idempotent.

SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'name'
  ),
  'ALTER TABLE employee_insurances DROP COLUMN name',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;
