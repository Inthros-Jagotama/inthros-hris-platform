-- =============================================================================
-- Tenant Migration: 077_performance_template_created_by
-- =============================================================================
-- Performance: tambah kolom `created_by` di performance_templates (KPI) dan
-- okr_templates (OKR) untuk menegakkan ownership — hanya pembuat template yang
-- boleh edit/hapus. NULL untuk template legacy; diisi dari user yang login
-- saat create.

ALTER TABLE performance_templates ADD COLUMN IF NOT EXISTS created_by CHAR(36) NULL;
ALTER TABLE okr_templates ADD COLUMN IF NOT EXISTS created_by CHAR(36) NULL;
