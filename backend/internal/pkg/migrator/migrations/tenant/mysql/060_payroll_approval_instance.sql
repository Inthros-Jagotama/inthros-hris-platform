-- 060_payroll_approval_instance.sql
-- Approval Module Phase 3: persist the approval_instance_id created for a
-- payroll run so its status can be updated via push-based callback instead
-- of relying on a separate manual call.

SET @add_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'payroll_runs'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE payroll_runs ADD COLUMN approval_instance_id CHAR(36) NULL AFTER locked_at'
);
PREPARE stmt FROM @add_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
