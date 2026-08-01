-- Down Migration: 027_add_fk_empins_insurance
-- Menghapus DB-level FK constraint fk_empins_insurance (idempotent).

SET @drop_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_insurances'
      AND CONSTRAINT_NAME = 'fk_empins_insurance'
  ),
  'ALTER TABLE employee_insurances DROP FOREIGN KEY fk_empins_insurance',
  'DO 0'
);
PREPARE stmt_drop_fk FROM @drop_fk_sql;
EXECUTE stmt_drop_fk;
DEALLOCATE PREPARE stmt_drop_fk;
