-- =============================================================================
-- Tenant Migration (down): 149_employee_supervisor
-- =============================================================================

ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employees_supervisor;

DROP INDEX IF EXISTS idx_employees_supervisor;

ALTER TABLE employees DROP COLUMN IF EXISTS supervisor_id;
