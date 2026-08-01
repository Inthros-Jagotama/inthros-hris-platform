-- Down Migration: 026_replace_category_with_insurance_id
-- Mengembalikan kolom category (ENUM 'BPJS'/'Non BPJS') dan menghapus insurance_id.

-- 1. Hapus FK constraint fk_empins_insurance (jika ada)
SET @drop_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_insurances'
      AND CONSTRAINT_NAME = 'fk_empins_insurance'
  ),
  'ALTER TABLE employee_insurances DROP FOREIGN KEY fk_empins_insurance',
  'DO 0'
);
PREPARE stmt_drop_fk FROM @drop_fk_sql;
EXECUTE stmt_drop_fk;
DEALLOCATE PREPARE stmt_drop_fk;

-- 2. Hapus index insurance_id (jika ada)
SET @drop_idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND index_name = 'idx_empins_insurance'
  ),
  'ALTER TABLE employee_insurances DROP INDEX idx_empins_insurance',
  'DO 0'
);
PREPARE stmt_drop_idx FROM @drop_idx_sql;
EXECUTE stmt_drop_idx;
DEALLOCATE PREPARE stmt_drop_idx;

-- 3. Kembalikan kolom category (jika belum ada)
SET @add_cat_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'category'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD COLUMN category ENUM(''BPJS'',''Non BPJS'') NULL AFTER employee_id'
);
PREPARE stmt_add_cat FROM @add_cat_sql;
EXECUTE stmt_add_cat;
DEALLOCATE PREPARE stmt_add_cat;

-- 4. Backfill: insurance '01' (BPJS Kesehatan) → category 'BPJS'
UPDATE employee_insurances ei
LEFT JOIN insurances i ON i.id = ei.insurance_id
SET ei.category = 'BPJS'
WHERE ei.category IS NULL AND ei.insurance_id IS NOT NULL AND i.code = '01';

-- 5. Hapus kolom insurance_id
SET @drop_ins_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'insurance_id'
  ),
  'ALTER TABLE employee_insurances DROP COLUMN insurance_id',
  'DO 0'
);
PREPARE stmt_drop_ins FROM @drop_ins_sql;
EXECUTE stmt_drop_ins;
DEALLOCATE PREPARE stmt_drop_ins;
