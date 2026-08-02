-- Migration: 036_seed_job_value_impact_unauthorized
-- Menambahkan tipe 'impact_unauthorized' (Dampak pada Hasil Akhir — Tidak
-- Memiliki Wewenang Keuangan) ke tabel job_management_values — 6 level mengikuti
-- docs/seeder/JobManagementValuesTableSeeder.php (grup 'Dampak pada Hasil Akhir
-- (Tidak Memiliki Wewenang Keuangan)').
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0001-4000-8000-000000000001', 'impact_unauthorized', 1, 'PAnciliary : Penyediaan jasa insidentil untuk penggunaan lainnya', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0002-4000-8000-000000000002', 'impact_unauthorized', 2, 'Suportif : Penyediaan layanan dukungan yang biasanya bersifat informasi dan pencatatan dalam suatu departemen. Atau Pengoperasian atau pemeliharaan sederhana peralatan atau mesin sekunder atau pendukung.', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0003-4000-8000-000000000003', 'impact_unauthorized', 3, 'Operasional : Pengoperasian proses dan/atau peralatan yang berhubungan langsung dengan rantai nilai inti bisnis ATAU pengoperasian/pemeliharaan peralatan atau sistem penting yang sangat kompleks dan merupakan inti bisnis', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0004-4000-8000-000000000004', 'impact_unauthorized', 4, 'Analitik : Penyediaan layanan khusus yang biasanya bersifat analitis, diagnostik, dan konsultasi.ATAU pengoperasian/pemeliharaan peralatan atau sistem penting dan sangat kompleks yang merupakan inti bisnis', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0005-4000-8000-000000000005', 'impact_unauthorized', 5, 'Pemandu : Memimpin area aktivitas yang dapat diidentifikasi seperti tim karyawan atau proyek sehari-hari dalam parameter yang ditentukan dengan baik, ATAU Memberikan nasihat dan bimbingan dalam bidang keahlian pada tingkat pengembangan kerangka kerja/kebijakan.', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dddddddd-0006-4000-8000-000000000006', 'impact_unauthorized', 6, 'Berpengaruh : Memberikan hasil dari operasi yang terdiri dari beberapa tim (terkait) yang melakukan beragam aktivitas.ATAU memastikan penyampaian program strategis yang efektif. ATAU Memimpin penyediaan kebijakan dan kerangka fungsional yang memungkinkan kinerja organisasi.', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dddddddd-0006-4000-8000-000000000006');
