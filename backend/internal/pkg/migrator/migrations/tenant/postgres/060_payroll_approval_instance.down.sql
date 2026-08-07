-- 060_payroll_approval_instance.down.sql

ALTER TABLE payroll_runs
    DROP COLUMN IF EXISTS approval_instance_id;
