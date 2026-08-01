-- Down Migration: 020_add_passport_to_employees
-- Hapus kolom passport dari tabel employees (kondisional, idempotent).

SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND column_name = 'passport'
  ),
  'ALTER TABLE employees DROP COLUMN passport',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;
