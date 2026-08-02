-- Migration: 038_seed_job_value_asset_authority
-- Menambahkan tipe 'asset_authority' (Wewenang Asset) ke tabel
-- job_management_values — 6 level mengikuti docs/seeder/JobManagementValuesTableSeeder.php
-- (grup 'Wewenang Asset').
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0001-4000-8000-000000000001', 'asset_authority', 1, 'Tidak ada Aset', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0002-4000-8000-000000000002', 'asset_authority', 2, 'Menggunakan Aset', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0003-4000-8000-000000000003', 'asset_authority', 3, 'Mengelola Aset', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0004-4000-8000-000000000004', 'asset_authority', 4, 'Memeriksa Aset', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0005-4000-8000-000000000005', 'asset_authority', 5, 'Memverifikasi Aset', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ffffffff-0006-4000-8000-000000000006', 'asset_authority', 6, 'Menyetujui Aset', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ffffffff-0006-4000-8000-000000000006');
