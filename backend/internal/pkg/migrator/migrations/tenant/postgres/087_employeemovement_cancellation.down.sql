-- =============================================================================
-- Tenant Migration Down: 087_employeemovement_cancellation
-- =============================================================================

DROP INDEX IF EXISTS idx_emp_mvmt_cancellation_instance;
ALTER TABLE employee_movements
    DROP COLUMN IF EXISTS cancellation_approval_instance_id;
