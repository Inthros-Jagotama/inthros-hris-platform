-- Down Migration: 019_fix_org_unique
-- =============================================================================
-- Mengembalikan constraint UNIQUE pada full_code seperti semula.

-- 1. Hapus composite unique index
DROP INDEX IF EXISTS uk_orgs_summary_code;

-- 2. Kembalikan UNIQUE constraint pada full_code
ALTER TABLE organizations ADD CONSTRAINT organizations_full_code_key UNIQUE (full_code);
