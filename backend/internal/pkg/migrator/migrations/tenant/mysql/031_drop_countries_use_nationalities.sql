-- Migration: 031_drop_countries_use_nationalities
-- Deduplikasi master data: tabel `countries` dihapus, hanya `nationalities` yang dipakai.
--
-- LATAR BELAKANG: ada duplikasi antara countries (id=code ISO alpha-2) dan
-- nationalities. Keputusan: pertahankan nationalities saja.
--   - nationalities berisi 240 negara (ISO alpha-2) + "Lainnya / Other" (LNY),
--     dengan kolom `code` = ISO code (sama dengan nilai yang disimpan
--     employees.nationality_id).
--   - FK lama employees.nationality_id → countries(id) di-retarget menjadi
--     employees.nationality_id → nationalities(code).
--
-- CATATAN (MySQL): Semua statement dibuat kondisional (idempotent) agar aman
-- dijalankan ulang.

-- 1. Drop FK lama employees.nationality_id → countries(id)
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

-- 1b. Lebarkan employees.nationality_id CHAR(2) → VARCHAR(20)
-- (nationalities.code varchar(20) berisi kode ISO 2-3 char seperti "US" dan "LNY";
--  CHAR(2) sebelumnya mengikuti countries.id dan tidak muat kode 3 char.)
SET @modify_col_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND column_name = 'nationality_id'
      AND column_type = 'char(2)'
  ),
  'ALTER TABLE employees MODIFY COLUMN nationality_id VARCHAR(20) NULL',
  'DO 0'
);
PREPARE stmt_modify_col FROM @modify_col_sql;
EXECUTE stmt_modify_col;
DEALLOCATE PREPARE stmt_modify_col;

-- 2. UNIQUE index pada nationalities.code (target FK baru)
SET @add_uq_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'nationalities'
      AND index_name = 'uq_nationalities_code'
  ),
  'DO 0',
  'ALTER TABLE nationalities ADD UNIQUE INDEX uq_nationalities_code (code)'
);
PREPARE stmt_add_uq FROM @add_uq_sql;
EXECUTE stmt_add_uq;
DEALLOCATE PREPARE stmt_add_uq;

-- 3. FK baru employees.nationality_id → nationalities(code)
SET @add_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE table_schema = DATABASE()
      AND table_name = 'employees'
      AND constraint_name = 'fk_employees_nationality'
      AND constraint_type = 'FOREIGN KEY'
  ),
  'DO 0',
  'ALTER TABLE employees ADD CONSTRAINT fk_employees_nationality FOREIGN KEY (nationality_id) REFERENCES nationalities(code) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk_sql;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;

-- 4. Drop tabel countries (referensi data sudah tidak dipakai)
DROP TABLE IF EXISTS countries;
