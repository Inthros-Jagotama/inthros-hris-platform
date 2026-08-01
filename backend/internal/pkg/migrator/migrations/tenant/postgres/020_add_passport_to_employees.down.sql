-- Down Migration: 020_add_passport_to_employees
-- Hapus kolom passport dari tabel employees (idempotent).

ALTER TABLE employees
    DROP COLUMN IF EXISTS passport;
