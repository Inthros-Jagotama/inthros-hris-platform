-- Down Migration: 019_fix_org_unique
-- =============================================================================
-- Mengembalikan constraint UNIQUE pada full_code seperti semula.
--
-- CATATAN (MySQL): 8.0 tidak mendukung `DROP INDEX IF EXISTS`, jadi drop/add
-- dilakukan secara kondisional lewat information_schema + prepared statement.

-- 1. Hapus composite unique key jika masih ada
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

-- 2. Kembalikan UNIQUE constraint pada full_code jika belum ada
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
