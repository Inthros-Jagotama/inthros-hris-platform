-- Down Migration: 037_seed_job_value_authority_unauthorized
-- Menghapus HANYA 8 row tipe 'authority_unauthorized' yang di-seed (by UUID tetap)
-- agar row buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'eeeeeeee-0001-4000-8000-000000000001',
  'eeeeeeee-0002-4000-8000-000000000002',
  'eeeeeeee-0003-4000-8000-000000000003',
  'eeeeeeee-0004-4000-8000-000000000004',
  'eeeeeeee-0005-4000-8000-000000000005',
  'eeeeeeee-0006-4000-8000-000000000006',
  'eeeeeeee-0007-4000-8000-000000000007',
  'eeeeeeee-0008-4000-8000-000000000008'
);
