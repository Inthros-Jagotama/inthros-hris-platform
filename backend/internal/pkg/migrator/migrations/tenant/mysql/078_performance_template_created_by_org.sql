-- =============================================================================
-- Tenant Migration: 078_performance_template_created_by_org
-- =============================================================================
-- Performance: tambah kolom `created_by_org_id` di performance_templates (KPI)
-- dan okr_templates (OKR). Kolom ini menyimpan ORGANISASI pembuat saat template
-- pertama kali dibuat — otorisasi edit/hapus membandingkan organisasi user yang
-- sedang login dengan kolom ini (bukan organisasi template). Jika karyawan
-- sudah pindah organisasi, dia tidak lagi bisa mengubah/menghapus template
-- meskipun dia yang membuatnya.
--
-- NULL untuk template yang dibuat sebelum migration ini (legacy). Saat otorisasi,
-- NULL di-fallback ke organization_id (template). Backfill tidak dilakukan karena
-- organisasi pembuat legacy tidak diketahui — fallback adalah pendekatan terbaik.

SET @add_kpi_created_by_org = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'created_by_org_id'
  ),
  'DO 0',
  'ALTER TABLE performance_templates ADD COLUMN created_by_org_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_kpi_created_by_org;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_okr_created_by_org = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_templates'
      AND column_name = 'created_by_org_id'
  ),
  'DO 0',
  'ALTER TABLE okr_templates ADD COLUMN created_by_org_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_okr_created_by_org;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
