-- =============================================================================
-- Tenant Migration: 149_employee_supervisor
-- =============================================================================
-- Menambahkan kolom supervisor_id (reports_to) pada tabel employees.
-- Menunjuk ke employee lain (atasan langsung) dalam tabel yang sama.
-- Dipakai Manager Assessment: daftar bawahan = employees.supervisor_id = manager.
-- Self-referencing FK sengaja TANPA ON DELETE SET NULL dari sisi FK constraint
-- agar hubungan atasan-bawahan tidak terputus diam-diam saat employee dihapus
-- (penghapusan employee ditangani aplikasi). Kolom tetap nullable.

ALTER TABLE employees
    ADD COLUMN supervisor_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_employees_supervisor ON employees (supervisor_id);

-- FK self-referencing: supervisor_id -> employees.id
ALTER TABLE employees
    ADD CONSTRAINT fk_employees_supervisor
    FOREIGN KEY (supervisor_id) REFERENCES employees (id) ON DELETE SET NULL;
