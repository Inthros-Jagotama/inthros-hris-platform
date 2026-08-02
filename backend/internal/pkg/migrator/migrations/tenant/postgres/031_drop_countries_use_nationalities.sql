-- Migration: 031_drop_countries_use_nationalities
-- Deduplikasi master data: tabel `countries` dihapus, hanya `nationalities` yang dipakai.
-- FK lama employees.nationality_id → countries(id) di-retarget menjadi
-- employees.nationality_id → nationalities(code).
--
-- CATATAN (PostgreSQL): statement dibuat idempotent.

-- 1. Drop FK lama employees.nationality_id → countries(id)
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employees_nationality;

-- 1b. Lebarkan employees.nationality_id CHAR(2) → VARCHAR(20)
-- (nationalities.code varchar(20) berisi kode ISO 2-3 char seperti "US" dan "LNY";
--  CHAR(2) sebelumnya mengikuti countries.id dan tidak muat kode 3 char.)
ALTER TABLE employees ALTER COLUMN nationality_id TYPE VARCHAR(20);

-- 2. UNIQUE index pada nationalities.code (target FK baru)
CREATE UNIQUE INDEX IF NOT EXISTS uq_nationalities_code ON nationalities (code);

-- 3. FK baru employees.nationality_id → nationalities(code)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_employees_nationality'
    ) THEN
        ALTER TABLE employees
            ADD CONSTRAINT fk_employees_nationality
            FOREIGN KEY (nationality_id) REFERENCES nationalities(code) ON DELETE SET NULL;
    END IF;
END $$;

-- 4. Drop tabel countries (referensi data sudah tidak dipakai)
DROP TABLE IF EXISTS countries;
