-- 060_payroll_approval_instance.sql
-- Approval Module Phase 3: persist the approval_instance_id created for a
-- payroll run so its status can be updated via push-based callback instead
-- of relying on a separate manual call.

ALTER TABLE payroll_runs
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;
