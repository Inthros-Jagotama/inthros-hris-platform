-- Migration: 035_seed_job_value_authority
-- Mengisi tipe 'authority' (Wewenang) di tabel job_management_values — 8 level
-- dengan note 'Memiliki Wewenang', mengikuti docs/seeder/JobManagementValuesTableSeeder.php
-- (grup 'Wewenang').
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0001-4000-8000-000000000001', 'authority', 1, 'Beroperasi dalam instruksi langsung dan rinci dengan pengawasan yang sangat ketat dan berkelanjutan.', 'Memiliki Wewenang', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0002-4000-8000-000000000002', 'authority', 2, 'Tunduk pada instruksi dan pekerjaan yang ditetapkan, rutinitas, di bawah pengawasan ketat.', 'Memiliki Wewenang', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0003-4000-8000-000000000003', 'authority', 3, 'Beroperasi sesuai praktik dan prosedur standar, instruksi kerja umum, dengan pengawasan kemajuan dan hasil.', 'Memiliki Wewenang', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0004-4000-8000-000000000004', 'authority', 4, 'Beroperasi dalam praktik dan prosedur yang tercakup dalam preseden atau kebijakan yang jelas dan peninjauan hasil akhir.', 'Memiliki Wewenang', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0005-4000-8000-000000000005', 'authority', 5, 'Tunduk pada praktik dan prosedur luas yang tercakup dalam preseden fungsional dan kebijakan serta arahan manajerial.', 'Memiliki Wewenang', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0006-4000-8000-000000000006', 'authority', 6, 'Tunduk pada arahan umum dan tujuan kebijakan yang ditetapkan secara luas.', 'Memiliki Wewenang', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0006-4000-8000-000000000006');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0007-4000-8000-000000000007', 'authority', 7, 'Hanya tunduk pada panduan keseluruhan mengenai tujuan organisasi dan orientasi kebijakan strategis.', 'Memiliki Wewenang', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0007-4000-8000-000000000007');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'cccccccc-0008-4000-8000-000000000008', 'authority', 8, 'Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum terhadap tren bisnis.', 'Memiliki Wewenang', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'cccccccc-0008-4000-8000-000000000008');
