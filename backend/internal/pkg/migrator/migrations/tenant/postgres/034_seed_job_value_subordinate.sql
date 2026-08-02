-- Migration: 034_seed_job_value_subordinate
-- Menambahkan tipe 'subordinate' (Total Bawahan / Total Subordinates) ke
-- tabel job_management_values — 5 level mengikuti docs/seeder/JobManagementValuesTableSeeder.php.
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'bbbbbbbb-0001-4000-8000-000000000001', 'subordinate', 1, 'Very Small', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'bbbbbbbb-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'bbbbbbbb-0002-4000-8000-000000000002', 'subordinate', 2, 'Small', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'bbbbbbbb-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'bbbbbbbb-0003-4000-8000-000000000003', 'subordinate', 3, 'Medium', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'bbbbbbbb-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'bbbbbbbb-0004-4000-8000-000000000004', 'subordinate', 4, 'Large', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'bbbbbbbb-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'bbbbbbbb-0005-4000-8000-000000000005', 'subordinate', 5, 'Very Large', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'bbbbbbbb-0005-4000-8000-000000000005');
