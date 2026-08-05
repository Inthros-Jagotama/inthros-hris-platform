-- Migration: 042_seed_job_value_technical_managerial
-- Menambahkan level kompetensi Technical (8 level) dan Managerial (5 level)
-- ke tabel job_management_values.
--
-- Catatan: kolom descriptions = nama level; kolom note = deskripsi lengkap per level
-- (technical 8 level, managerial 5 level) sesuai ketentuan user.
-- type memakai nilai literal 'technical' / 'managerial' (bukan slug per kompetensi)
-- karena kompetensi individual kini bersumber dari tabel competencies.
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID) belum ada.
--
-- CATATAN untuk tenant yang sudah apply 042 versi lama (158 row slugs):
-- jalankan down 042 dulu (menghapus row slugs via type_group), lalu up ini,
-- agar tidak terjadi duplikat (171 row).

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'b7c2f94b-13bd-4bc8-aa7f-83a839400d35', 'technical', 1, 'Primary', 'Membaca, menulis, melakukan perhitungan sederhana.', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'b7c2f94b-13bd-4bc8-aa7f-83a839400d35');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '0bb49f73-c068-419d-ada5-a986d6966005', 'technical', 2, 'Elementary Vocational', 'Sudah menggunakan teknologi sederhana.', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '0bb49f73-c068-419d-ada5-a986d6966005');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '158b213b-3262-4099-9446-55f7916f05aa', 'technical', 3, 'Vocational', 'Menggunakan kaidah-kaidah sederhana.', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '158b213b-3262-4099-9446-55f7916f05aa');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '39c0d3ab-09fe-40d8-a272-040cc238af3a', 'technical', 4, 'Advanced Vocational', 'Pengetahuan/keahlian/pengalaman selevel akademi dengan penekanan pada pengetahuan praktis, pekerjaan teknikal dan supervisi.', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '39c0d3ab-09fe-40d8-a272-040cc238af3a');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '359b3f8b-ad31-4036-bf6a-f53168d3687f', 'technical', 5, 'Basic Professional', 'Pengetahuan/keahlian bidang teknis, ilmiah, spesialis, konseptual dengan keterlibatan dalam praktik dan prosedur.', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '359b3f8b-ad31-4036-bf6a-f53168d3687f');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '9a6cfdad-1e43-48cb-9b3b-6421b23985af', 'technical', 6, 'Seasoned Professional', 'Pengetahuan/keahlian dengan pengalaman yang luas.', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '9a6cfdad-1e43-48cb-9b3b-6421b23985af');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'aeaafba0-8478-4384-82c0-36dcc91f6eeb', 'technical', 7, 'Professional Mastery', 'Pemahaman mendalam konsep, prinsip, dan praktik melalui pengalaman luas dalam bisnis serta mampu menerjemahkan visi misi perusahaan.', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'aeaafba0-8478-4384-82c0-36dcc91f6eeb');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '61116efa-98ec-4f5a-8cfb-e7b263e95b9d', 'technical', 8, 'Unique Authority', 'Otoritas unik di bidangnya.', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '61116efa-98ec-4f5a-8cfb-e7b263e95b9d');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'b2ad32dc-b067-4c8d-94c6-c56b7a85503b', 'managerial', 1, 'Task', 'Melibatkan perencanaan harian dan sikap proaktif.', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'b2ad32dc-b067-4c8d-94c6-c56b7a85503b');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '546baaef-b125-401d-901f-644f0809495e', 'managerial', 2, 'Supervisory', 'Planning, organizing, dan controlling mingguan atau bulanan, mengelola tim kecil serta koordinasi terbatas.', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '546baaef-b125-401d-901f-644f0809495e');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'e87e33f4-f746-48d6-b8e1-bf78572b23ca', 'managerial', 3, 'Managerial', 'Planning, organizing, dan controlling tahunan atau lebih, termasuk budgeting dan koordinasi lintas fungsi.', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'e87e33f4-f746-48d6-b8e1-bf78572b23ca');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '0034107e-e72b-4425-a269-fd520018bd67', 'managerial', 4, 'Diverse Managerial', 'Integrasi berbagai fungsi, perencanaan jangka panjang, serta kemampuan menyelesaikan konflik (conflict resolution).', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '0034107e-e72b-4425-a269-fd520018bd67');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT 'f5169e21-d099-4ca6-9495-9c5e7ee923ac', 'managerial', 5, 'Total Managerial', 'Bertanggung jawab terhadap skala dan visi perusahaan, penyusunan kebijakan organisasi, serta koordinasi dengan pihak eksternal.', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'f5169e21-d099-4ca6-9495-9c5e7ee923ac');
