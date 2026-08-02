-- Down Migration: 035_seed_job_value_authority
-- Menghapus HANYA 8 row tipe 'authority' yang di-seed (by UUID tetap) agar
-- row authority buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'cccccccc-0001-4000-8000-000000000001',
  'cccccccc-0002-4000-8000-000000000002',
  'cccccccc-0003-4000-8000-000000000003',
  'cccccccc-0004-4000-8000-000000000004',
  'cccccccc-0005-4000-8000-000000000005',
  'cccccccc-0006-4000-8000-000000000006',
  'cccccccc-0007-4000-8000-000000000007',
  'cccccccc-0008-4000-8000-000000000008'
);
