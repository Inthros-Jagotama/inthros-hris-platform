-- Down Migration: 039_seed_job_value_communicating_influencing_skill
-- Menghapus HANYA 3 row tipe 'communicating_influencing_skill' yang di-seed
-- (by UUID tetap) agar row communicating_influencing_skill buatan user
-- (lewat UI) tidak ikut terhapus saat rollback.

DELETE FROM job_management_values
WHERE id IN (
  '11111111-0001-4000-8000-000000000001',
  '11111111-0002-4000-8000-000000000002',
  '11111111-0003-4000-8000-000000000003'
);
