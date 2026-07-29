-- Down Migration: 019_fix_org_unique
-- =============================================================================
-- Mengembalikan constraint UNIQUE pada full_code seperti semula.

-- 1. Hapus composite unique key
ALTER TABLE organizations DROP INDEX IF EXISTS `uk_orgs_summary_code`;

-- 2. Kembalikan UNIQUE constraint pada full_code (MySQL auto-names it as column name)
ALTER TABLE organizations ADD UNIQUE INDEX `full_code` (`full_code`);
