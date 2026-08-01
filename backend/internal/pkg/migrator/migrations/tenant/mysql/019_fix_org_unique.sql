-- =============================================================================
-- Tenant Migration: 019_fix_org_unique
-- =============================================================================
-- Mengubah constraint UNIQUE pada full_code menjadi composite unique
-- (organization_summary_id, full_code) agar setiap summary bisa memiliki
-- full_code yang sama tanpa konflik.
--
-- CATATAN (MySQL): 8.0 tidak mendukung `DROP INDEX IF EXISTS`, jadi drop/add
-- dilakukan secara kondisional lewat information_schema + prepared statement.
-- Aman untuk tenant baru (uk_orgs_summary_code sudah dibuat migration 002,
-- sehingga jadi no-op) maupun tenant lama (masih punya unique index 'full_code').

-- 1. Hapus unique index lama 'full_code' jika masih ada
SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'organizations'
      AND index_name = 'full_code'
  ),
  'ALTER TABLE organizations DROP INDEX `full_code`',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;

-- 2. Tambahkan composite unique key jika belum ada
SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'organizations'
      AND index_name = 'uk_orgs_summary_code'
  ),
  'DO 0',
  'ALTER TABLE organizations ADD UNIQUE KEY `uk_orgs_summary_code` (`organization_summary_id`, `full_code`)'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;
