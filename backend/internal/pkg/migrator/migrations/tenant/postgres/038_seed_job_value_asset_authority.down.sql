-- Down Migration: 038_seed_job_value_asset_authority
-- Menghapus HANYA 6 row tipe 'asset_authority' yang di-seed (by UUID tetap) agar
-- row asset_authority buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'ffffffff-0001-4000-8000-000000000001',
  'ffffffff-0002-4000-8000-000000000002',
  'ffffffff-0003-4000-8000-000000000003',
  'ffffffff-0004-4000-8000-000000000004',
  'ffffffff-0005-4000-8000-000000000005',
  'ffffffff-0006-4000-8000-000000000006'
);
