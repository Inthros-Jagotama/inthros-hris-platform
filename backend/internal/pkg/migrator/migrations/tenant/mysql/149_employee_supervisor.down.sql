-- =============================================================================
-- Tenant Migration (down): 149_employee_supervisor (MySQL)
-- =============================================================================

ALTER TABLE employees DROP FOREIGN KEY fk_employees_supervisor;

DROP INDEX idx_employees_supervisor ON employees;

ALTER TABLE employees DROP COLUMN supervisor_id;
