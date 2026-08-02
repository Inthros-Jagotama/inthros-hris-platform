-- Down Migration: 034_seed_job_value_subordinate
-- Menghapus HANYA 5 row tipe 'subordinate' yang di-seed (by UUID tetap) agar
-- row subordinate buatan user (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  'bbbbbbbb-0001-4000-8000-000000000001',
  'bbbbbbbb-0002-4000-8000-000000000002',
  'bbbbbbbb-0003-4000-8000-000000000003',
  'bbbbbbbb-0004-4000-8000-000000000004',
  'bbbbbbbb-0005-4000-8000-000000000005'
);
