-- Down Migration: 033_seed_job_value_activity
-- Menghapus HANYA 5 row tipe 'activity' yang di-seed (by UUID tetap) agar
-- row activity buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'aaaaaaaa-0001-4000-8000-000000000001',
  'aaaaaaaa-0002-4000-8000-000000000002',
  'aaaaaaaa-0003-4000-8000-000000000003',
  'aaaaaaaa-0004-4000-8000-000000000004',
  'aaaaaaaa-0005-4000-8000-000000000005'
);
