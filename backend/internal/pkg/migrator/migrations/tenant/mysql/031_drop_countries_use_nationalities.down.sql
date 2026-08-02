-- Down Migration: 031_drop_countries_use_nationalities
-- Restore tabel `countries` dan kembalikan FK employees.nationality_id →
-- countries(id). (Kondisional / idempotent.)

-- 1. Drop FK baru employees.nationality_id → nationalities(code)
SET @drop_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND constraint_name = 'fk_employees_nationality'
      AND constraint_type = 'FOREIGN KEY'
  ),
  'ALTER TABLE employees DROP FOREIGN KEY fk_employees_nationality',
  'DO 0'
);
PREPARE stmt_drop_fk FROM @drop_fk_sql;
EXECUTE stmt_drop_fk;
DEALLOCATE PREPARE stmt_drop_fk;

-- 2. Hapus UNIQUE index pada nationalities.code
SET @drop_uq_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'nationalities'
      AND index_name = 'uq_nationalities_code'
  ),
  'ALTER TABLE nationalities DROP INDEX uq_nationalities_code',
  'DO 0'
);
PREPARE stmt_drop_uq FROM @drop_uq_sql;
EXECUTE stmt_drop_uq;
DEALLOCATE PREPARE stmt_drop_uq;

-- 2b. Kembalikan employees.nationality_id ke CHAR(2) (struktur asli 003_employee.sql)
ALTER TABLE employees MODIFY COLUMN nationality_id CHAR(2) NULL;

-- 3. Recreate tabel countries (struktur asli dari 001_master_data.sql)
CREATE TABLE IF NOT EXISTS countries (
    id          CHAR(2) PRIMARY KEY,
    code        VARCHAR(2) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    INDEX idx_countries_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Kembalikan FK lama → countries(id)
SET @add_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND constraint_name = 'fk_employees_nationality'
      AND constraint_type = 'FOREIGN KEY'
  ),
  'DO 0',
  'ALTER TABLE employees ADD CONSTRAINT fk_employees_nationality FOREIGN KEY (nationality_id) REFERENCES countries(id) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk_sql;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;
