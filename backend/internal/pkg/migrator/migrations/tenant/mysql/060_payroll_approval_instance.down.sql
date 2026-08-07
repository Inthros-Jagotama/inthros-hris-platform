-- 060_payroll_approval_instance.down.sql

SET @drop_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'payroll_runs'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE payroll_runs DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
