-- 062_employeemovement_approval_instance.sql
-- Approval Module Phase 4: persist the approval_instance_id created for an
-- employee movement so its status can be updated via push-based callback
-- from the central approval module.

ALTER TABLE employee_movements
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_emp_mvmt_approval_instance ON employee_movements (approval_instance_id);
