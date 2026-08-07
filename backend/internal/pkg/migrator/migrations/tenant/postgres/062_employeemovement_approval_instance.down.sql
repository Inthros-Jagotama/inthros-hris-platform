-- 062_employeemovement_approval_instance.down.sql

DROP INDEX IF EXISTS idx_emp_mvmt_approval_instance;

ALTER TABLE employee_movements
    DROP COLUMN IF EXISTS approval_instance_id;
