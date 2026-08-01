-- Migration: 020_add_passport_to_employees
-- Menambahkan kolom passport ke tabel employees.
--
-- CATATAN (MySQL): Migration 003_employee.sql sudah menyertakan kolom
-- `passport` di CREATE TABLE. Migration ini hanya dibutuhkan untuk tenant
-- lama yang tabel employees-nya dibuat sebelum kolom passport ada, jadi
-- penambahan dilakukan secara kondisional (idempotent).

SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND column_name = 'passport'
  ),
  'DO 0',
  'ALTER TABLE employees ADD COLUMN passport VARCHAR(50) NULL AFTER nationality_id'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;
