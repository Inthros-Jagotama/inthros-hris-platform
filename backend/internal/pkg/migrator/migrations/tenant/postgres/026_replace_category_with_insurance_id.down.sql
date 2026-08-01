-- Down Migration: 026_replace_category_with_insurance_id
-- Mengembalikan kolom category dan menghapus insurance_id (idempotent).

-- 1. Hapus FK constraint fk_empins_insurance
ALTER TABLE employee_insurances
    DROP CONSTRAINT IF EXISTS fk_empins_insurance;

-- 2. Hapus index insurance_id
DROP INDEX IF EXISTS idx_empins_insurance;

-- 3. Kembalikan kolom category
ALTER TABLE employee_insurances
    ADD COLUMN IF NOT EXISTS category VARCHAR(255) NULL;

-- 4. Backfill: insurance '01' (BPJS Kesehatan) → category 'BPJS'
UPDATE employee_insurances ei
SET category = 'BPJS'
FROM insurances i
WHERE i.id = ei.insurance_id
  AND ei.category IS NULL
  AND i.code = '01';

-- 5. Hapus kolom insurance_id
ALTER TABLE employee_insurances
    DROP COLUMN IF EXISTS insurance_id;
