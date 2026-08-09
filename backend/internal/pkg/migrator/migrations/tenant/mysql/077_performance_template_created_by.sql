-- =============================================================================
-- Tenant Migration: 077_performance_template_created_by
-- =============================================================================
-- Performance: tambah kolom `created_by` di performance_templates (KPI) dan
-- okr_templates (OKR) untuk menegakkan ownership — hanya pembuat template yang
-- boleh edit/hapus. NULL untuk template yang dibuat sebelum migration ini
-- (legacy, tanpa pemilik tercatat); diisi dari user yang login saat create.

SET @add_kpi_created_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'created_by'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD COLUMN created_by CHAR(36) NULL'
);
PREPARE stmt FROM @add_kpi_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_okr_created_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_templates'
      AND column_name = 'created_by'
  ),
  'DO 0',
  'ALTER TABLE okr_templates ADD COLUMN created_by CHAR(36) NULL'
);
PREPARE stmt FROM @add_okr_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
