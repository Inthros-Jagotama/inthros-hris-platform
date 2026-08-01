-- Migration: 026_replace_category_with_insurance_id
-- Mengganti kolom category pada tabel employee_insurances dengan kolom
-- insurance_id CHAR(36) yang ber-relasi ke tabel master insurances.
--
-- CATATAN (PostgreSQL): Idempotent via ADD COLUMN IF NOT EXISTS.

ALTER TABLE employee_insurances
    ADD COLUMN IF NOT EXISTS insurance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_empins_insurance
    ON employee_insurances (insurance_id);

-- Backfill data lama: category 'BPJS' → insurance '01' (BPJS Kesehatan)
UPDATE employee_insurances ei
SET insurance_id = i.id
FROM insurances i
WHERE i.code = '01'
  AND ei.insurance_id IS NULL
  AND ei.category = 'BPJS';

ALTER TABLE employee_insurances
    DROP COLUMN IF EXISTS category;

-- DB-level FK: relasi literal ke tabel insurances (idempotent)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_empins_insurance'
    ) THEN
        ALTER TABLE employee_insurances
            ADD CONSTRAINT fk_empins_insurance
            FOREIGN KEY (insurance_id) REFERENCES insurances(id)
            ON DELETE SET NULL;
    END IF;
END $$;
