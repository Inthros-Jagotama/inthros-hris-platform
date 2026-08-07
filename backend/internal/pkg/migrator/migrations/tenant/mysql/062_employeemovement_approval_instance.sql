-- 062_employeemovement_approval_instance.sql
-- Approval Module Phase 4: persist the approval_instance_id created for an
-- employee movement so its status can be updated via push-based callback
-- from the central approval module.

SET @add_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN approval_instance_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_mvmt_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND INDEX_NAME = 'idx_emp_mvmt_approval_instance'
  ),
  'DO 0',
  'CREATE INDEX idx_emp_mvmt_approval_instance ON employee_movements (approval_instance_id)'
);
PREPARE stmt FROM @add_mvmt_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
