-- Down Migration: 028_drop_name_from_employee_insurances
-- Mengembalikan kolom name (idempotent).

ALTER TABLE employee_insurances
    ADD COLUMN IF NOT EXISTS name VARCHAR(100) NULL;
