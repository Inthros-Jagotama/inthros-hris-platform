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
-- NULL untuk template legacy; saat otorisasi di-fallback ke organization_id.

ALTER TABLE performance_templates ADD COLUMN IF NOT EXISTS created_by_org_id UUID NULL;
ALTER TABLE okr_templates ADD COLUMN IF NOT EXISTS created_by_org_id UUID NULL;
