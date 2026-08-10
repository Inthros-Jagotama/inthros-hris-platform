-- =============================================================================
-- Tenant Migration Down: 087_employeemovement_cancellation
-- =============================================================================

SET @drop_mvmt_cancel_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND INDEX_NAME = 'idx_emp_mvmt_cancellation_instance'
  ),
  'DROP INDEX idx_emp_mvmt_cancellation_instance ON employee_movements',
  'DO 0'
);
PREPARE stmt FROM @drop_mvmt_cancel_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cancel_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'cancellation_approval_instance_id'
  ),
  'ALTER TABLE employee_movements DROP COLUMN cancellation_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_cancel_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
