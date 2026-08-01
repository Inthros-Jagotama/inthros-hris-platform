-- Migration: 026_replace_category_with_insurance_id
-- Mengganti kolom category (ENUM 'BPJS'/'Non BPJS') pada tabel employee_insurances
-- dengan kolom insurance_id CHAR(36) yang ber-relasi ke tabel master insurances.
--
-- CATATAN (MySQL): Kondisional/idempotent — aman untuk tenant lama maupun baru.

-- 1. Tambah kolom insurance_id (jika belum ada)
SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'insurance_id'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD COLUMN insurance_id CHAR(36) NULL AFTER employee_id'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;

-- 2. Index untuk lookup relasi (idempotent)
SET @idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND index_name = 'idx_empins_insurance'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD INDEX idx_empins_insurance (insurance_id)'
);
PREPARE stmt_idx FROM @idx_sql;
EXECUTE stmt_idx;
DEALLOCATE PREPARE stmt_idx;

-- 3. Backfill data lama: category 'BPJS' → insurance '01' (BPJS Kesehatan)
UPDATE employee_insurances ei
LEFT JOIN insurances i ON i.code = '01' AND i.deleted_at IS NULL
SET ei.insurance_id = i.id
WHERE ei.insurance_id IS NULL AND ei.category = 'BPJS';

-- 4. Hapus kolom category (digantikan insurance_id)
SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_insurances'
      AND column_name = 'category'
  ),
  'ALTER TABLE employee_insurances DROP COLUMN category',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;

-- 5. DB-level FK: relasi literal ke tabel insurances (idempotent)
SET @fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_insurances'
      AND CONSTRAINT_NAME = 'fk_empins_insurance'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD CONSTRAINT fk_empins_insurance FOREIGN KEY (insurance_id) REFERENCES insurances(id) ON DELETE SET NULL'
);
PREPARE stmt_fk FROM @fk_sql;
EXECUTE stmt_fk;
DEALLOCATE PREPARE stmt_fk;
