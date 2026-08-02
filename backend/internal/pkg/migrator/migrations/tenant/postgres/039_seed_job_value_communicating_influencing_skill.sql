-- Migration: 039_seed_job_value_communicating_influencing_skill
-- Menambahkan tipe 'communicating_influencing_skill' (Communicating & Influencing
-- Skill) ke tabel job_management_values — 3 level mengikuti
-- docs/seeder/JobManagementValuesTableSeeder.php (grup 'Communicating &
-- Influencing Skill', parent 'Psychological').
--
-- Catatan: teks note seeder (Berkomunikasi / Alasan / Perubahan Perilaku)
-- di-promote menjadi descriptions per level (keputusan user 02 Aug 2026).
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '11111111-0001-4000-8000-000000000001', 'communicating_influencing_skill', 1, 'Berkomunikasi', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '11111111-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '11111111-0002-4000-8000-000000000002', 'communicating_influencing_skill', 2, 'Alasan', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '11111111-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '11111111-0003-4000-8000-000000000003', 'communicating_influencing_skill', 3, 'Perubahan Perilaku', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '11111111-0003-4000-8000-000000000003');
