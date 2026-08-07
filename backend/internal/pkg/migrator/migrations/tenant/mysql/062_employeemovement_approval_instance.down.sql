-- 062_employeemovement_approval_instance.down.sql

SET @drop_mvmt_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND INDEX_NAME = 'idx_emp_mvmt_approval_instance'
  ),
  'DROP INDEX idx_emp_mvmt_approval_instance ON employee_movements',
  'DO 0'
);
PREPARE stmt FROM @drop_mvmt_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE employee_movements DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
