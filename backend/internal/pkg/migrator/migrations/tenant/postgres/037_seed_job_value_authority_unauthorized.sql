-- Migration: 037_seed_job_value_authority_unauthorized
-- Menambahkan tipe 'authority_unauthorized' (Memiliki Wewenang / Tidak Memiliki
-- Wewenang Keuangan) ke tabel job_management_values — 8 level mengikuti
-- docs/seeder/JobManagementValuesTableSeeder.php (grup 'Memiliki Wewenang',
-- note 'Tidak memiliki Wewenang').
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0001-4000-8000-000000000001', 'authority_unauthorized', 1, 'Dikendalikan secara ketat : Beroperasi dalam instruksi langsung dan rinci dengan pengawasan yang sangat ketat dan berkelanjutan.', 'Tidak memiliki Wewenang', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0002-4000-8000-000000000002', 'authority_unauthorized', 2, 'Terkendali : Tunduk pada instruksi dan pekerjaan yang ditetapkan, rutinitas, di bawah pengawasan ketat.', 'Tidak memiliki Wewenang', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0003-4000-8000-000000000003', 'authority_unauthorized', 3, 'Terstandar : Beroperasi sesuai praktik dan prosedur standar, instruksi kerja umum dan pengawasan kemajuan dan hasil.', 'Tidak memiliki Wewenang', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0004-4000-8000-000000000004', 'authority_unauthorized', 4, 'Secara umum diatur : Beroperasi dalam praktik dan prosedur yang tercakup dalam preseden atau kebijakan yang jelas dan peninjauan hasil akhir.', 'Tidak memiliki Wewenang', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0005-4000-8000-000000000005', 'authority_unauthorized', 5, 'Terarah dengan Jelas; Tunduk pada praktik dan prosedur luas yang tercakup dalam preseden fungsional dan kebijakan serta arahan manajerial', 'Tidak memiliki Wewenang', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0006-4000-8000-000000000006', 'authority_unauthorized', 6, 'Dipandu : Hanya tunduk pada panduan keseluruhan mengenai tujuan organisasi secara luas dan orientasi kebijakan strategis.', 'Tidak memiliki Wewenang', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0006-4000-8000-000000000006');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0007-4000-8000-000000000007', 'authority_unauthorized', 7, 'Dipandu Secara Strategis : Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum', 'Tidak memiliki Wewenang', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0007-4000-8000-000000000007');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'eeeeeeee-0008-4000-8000-000000000008', 'authority_unauthorized', 8, 'Dipandu Secara Strategis : Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum dalam menanggapi tren bisnis.', 'Tidak memiliki Wewenang', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'eeeeeeee-0008-4000-8000-000000000008');
