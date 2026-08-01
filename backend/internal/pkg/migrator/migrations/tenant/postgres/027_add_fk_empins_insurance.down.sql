-- Down Migration: 027_add_fk_empins_insurance
-- Menghapus DB-level FK constraint fk_empins_insurance (idempotent).

ALTER TABLE employee_insurances
    DROP CONSTRAINT IF EXISTS fk_empins_insurance;
