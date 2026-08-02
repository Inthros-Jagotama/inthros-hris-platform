-- Migration: 033_seed_job_value_activity
-- Menambahkan tipe 'activity' (Aktifitas Fisik / Physical Activity) ke
-- tabel job_management_values — 5 level mengikuti docs/seeder/JobManagementValuesTableSeeder.php.
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'aaaaaaaa-0001-4000-8000-000000000001', 'activity', 1, 'Banyak duduk sedikit bergerak', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aaaaaaaa-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'aaaaaaaa-0002-4000-8000-000000000002', 'activity', 2, 'Seimbang duduk dan berdiri', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aaaaaaaa-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'aaaaaaaa-0003-4000-8000-000000000003', 'activity', 3, 'Sedikit duduk, banyak berdiri dan berjalan', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aaaaaaaa-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'aaaaaaaa-0004-4000-8000-000000000004', 'activity', 4, 'Aktivitas fisik tinggi, gunakan organ dan indra', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aaaaaaaa-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'aaaaaaaa-0005-4000-8000-000000000005', 'activity', 5, 'Aktivitas sangat tinggi, melakukan pengawasan', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aaaaaaaa-0005-4000-8000-000000000005');
