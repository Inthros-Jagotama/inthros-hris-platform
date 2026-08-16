-- =============================================================================
-- Tenant Migration: 149_employee_supervisor (MySQL)
-- =============================================================================
-- Menambahkan kolom supervisor_id (reports_to) pada tabel employees.
-- Dipakai Manager Assessment: daftar bawahan = employees.supervisor_id = manager.

ALTER TABLE employees
    ADD COLUMN supervisor_id CHAR(36) NULL AFTER status;

CREATE INDEX idx_employees_supervisor ON employees (supervisor_id);

-- Self-referencing FK: supervisor_id -> employees.id
ALTER TABLE employees
    ADD CONSTRAINT fk_employees_supervisor
    FOREIGN KEY (supervisor_id) REFERENCES employees (id) ON DELETE SET NULL;
