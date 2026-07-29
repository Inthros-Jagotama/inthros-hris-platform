-- =============================================================================
-- Tenant Migration: 019_fix_org_unique
-- =============================================================================
-- Mengubah constraint UNIQUE pada full_code menjadi composite unique
-- (organization_summary_id, full_code) agar setiap summary bisa memiliki
-- full_code yang sama tanpa konflik.

-- 1. Hapus unique index lama pada full_code (MySQL auto-names it as column name)
ALTER TABLE organizations DROP INDEX IF EXISTS `full_code`;

-- 2. Tambahkan composite unique key (organization_summary_id, full_code)
--    agar full_code hanya unique dalam satu summary yang sama.
ALTER TABLE organizations ADD UNIQUE KEY `uk_orgs_summary_code` (`organization_summary_id`, `full_code`);
