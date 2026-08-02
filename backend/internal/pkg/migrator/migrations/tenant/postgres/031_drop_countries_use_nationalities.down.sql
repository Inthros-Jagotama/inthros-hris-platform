-- Down Migration: 031_drop_countries_use_nationalities
-- Restore tabel `countries` dan kembalikan FK employees.nationality_id →
-- countries(id). (Idempotent.)

-- 1. Drop FK baru employees.nationality_id → nationalities(code)
ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employees_nationality;

-- 2. Hapus UNIQUE index pada nationalities.code
DROP INDEX IF EXISTS uq_nationalities_code;

-- 2b. Kembalikan employees.nationality_id ke CHAR(2) (struktur asli 003_employee.sql)
ALTER TABLE employees ALTER COLUMN nationality_id TYPE CHAR(2);

-- 3. Recreate tabel countries (struktur asli dari 001_master_data.sql)
CREATE TABLE IF NOT EXISTS countries (
    id          CHAR(2) PRIMARY KEY,
    code        VARCHAR(2) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);

-- 4. Kembalikan FK lama → countries(id)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_employees_nationality'
    ) THEN
        ALTER TABLE employees
            ADD CONSTRAINT fk_employees_nationality
            FOREIGN KEY (nationality_id) REFERENCES countries(id) ON DELETE SET NULL;
    END IF;
END $$;
