-- =============================================================================
-- Tenant Migration: 019_fix_org_unique
-- =============================================================================
-- Mengubah constraint UNIQUE pada full_code menjadi composite unique
-- (organization_summary_id, full_code) agar setiap summary bisa memiliki
-- full_code yang sama tanpa konflik.

-- 1. Hapus constraint unique lama (PostgreSQL auto-names it as table_column_key)
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS organizations_full_code_key;

-- 2. Hapus index unique lama jika ada (fallback untuk naming yang berbeda)
DROP INDEX IF EXISTS idx_organizations_full_code;

-- 3. Tambahkan composite unique index
--    PostgreSQL tidak bisa membuat UNIQUE constraint dengan NULLable columns
--    secara langsung, jadi kita gunakan partial unique index untuk menangani
--    NULL pada organization_summary_id.
--    Root organization (tanpa summary) tetap bisa memiliki full_code global unique.
CREATE UNIQUE INDEX IF NOT EXISTS uk_orgs_summary_code ON organizations (COALESCE(organization_summary_id, 'none'), full_code);
