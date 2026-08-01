-- Down Migration: 028_drop_name_from_employee_insurances
-- Mengembalikan kolom name (idempotent).

SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'name'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD COLUMN name VARCHAR(100) NULL AFTER number'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;
