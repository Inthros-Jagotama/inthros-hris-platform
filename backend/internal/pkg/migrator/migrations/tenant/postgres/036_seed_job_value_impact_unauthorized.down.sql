-- Down Migration: 036_seed_job_value_impact_unauthorized
-- Menghapus HANYA 6 row tipe 'impact_unauthorized' yang di-seed (by UUID tetap)
-- agar row buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'dddddddd-0001-4000-8000-000000000001',
  'dddddddd-0002-4000-8000-000000000002',
  'dddddddd-0003-4000-8000-000000000003',
  'dddddddd-0004-4000-8000-000000000004',
  'dddddddd-0005-4000-8000-000000000005',
  'dddddddd-0006-4000-8000-000000000006'
);
