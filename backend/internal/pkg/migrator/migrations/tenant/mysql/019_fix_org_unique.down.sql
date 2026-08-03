-- Down Migration: 019_fix_org_unique
-- =============================================================================
-- Mengembalikan constraint UNIQUE pada full_code seperti semula.
--
-- CATATAN (MySQL): 8.0 tidak mendukung `DROP INDEX IF EXISTS`, jadi drop/add
-- dilakukan secara kondisional lewat information_schema + prepared statement.
--
-- PENTING: index `uk_orgs_summary_code` dipakai oleh FK `fk_orgs_summary`
-- (kolom organization_summary_id) sebagai supporting index. Karena itu FK
-- harus di-drop TERLEBIH DAHULU sebelum index-nya, lalu dibuat ulang setelah
-- index `full_code` dikembalikan.

-- 1. Hapus FK fk_orgs_summary terlebih dahulu (jika masih ada)
SET @drop_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'organizations'
      AND CONSTRAINT_NAME = 'fk_orgs_summary'
  ),
  'ALTER TABLE organizations DROP FOREIGN KEY fk_orgs_summary',
  'DO 0'
);
PREPARE stmt_drop_fk FROM @drop_fk_sql;
EXECUTE stmt_drop_fk;
DEALLOCATE PREPARE stmt_drop_fk;

-- 2. Hapus composite unique key jika masih ada
SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'organizations'
      AND index_name = 'uk_orgs_summary_code'
  ),
  'ALTER TABLE organizations DROP INDEX `uk_orgs_summary_code`',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;

-- 3. Kembalikan UNIQUE constraint pada full_code jika belum ada
SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'organizations'
      AND index_name = 'full_code'
  ),
  'DO 0',
  'ALTER TABLE organizations ADD UNIQUE INDEX `full_code` (`full_code`)'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;

-- 4. Kembalikan FK fk_orgs_summary (MySQL otomatis membuat supporting index
--    pada organization_summary_id jika belum ada, sama seperti schema asli)
SET @add_fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'organizations'
      AND CONSTRAINT_NAME = 'fk_orgs_summary'
  ),
  'DO 0',
  'ALTER TABLE organizations ADD CONSTRAINT fk_orgs_summary FOREIGN KEY (organization_summary_id) REFERENCES organization_summaries(id) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk_sql;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;
