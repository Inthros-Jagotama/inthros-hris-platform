-- Migration: 042_seed_job_value_technical_managerial
-- Menambahkan grup 'Technical' (16 tipe x 8 level) dan 'Managerial' (6 tipe x 5 level)
-- ke tabel job_management_values mengikuti docs/seeder/JobManagementValuesTableSeeder.php.
--
-- Catatan: teks note seeder di-promote menjadi descriptions per level — pola sama
-- dengan 039/040/041 (keputusan user).
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0001-4000-8000-000000000001', 'competency_based_human_resources_management', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0002-4000-8000-000000000002', 'competency_based_human_resources_management', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0003-4000-8000-000000000003', 'competency_based_human_resources_management', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0004-4000-8000-000000000004', 'competency_based_human_resources_management', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0005-4000-8000-000000000005', 'competency_based_human_resources_management', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0006-4000-8000-000000000006', 'competency_based_human_resources_management', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0006-4000-8000-000000000006');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0007-4000-8000-000000000007', 'competency_based_human_resources_management', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0007-4000-8000-000000000007');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0008-4000-8000-000000000008', 'competency_based_human_resources_management', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0008-4000-8000-000000000008');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0009-4000-8000-000000000009', 'competency_development', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0009-4000-8000-000000000009');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0010-4000-8000-000000000010', 'competency_development', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0010-4000-8000-000000000010');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0011-4000-8000-000000000011', 'competency_development', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0011-4000-8000-000000000011');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0012-4000-8000-000000000012', 'competency_development', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0012-4000-8000-000000000012');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0013-4000-8000-000000000013', 'competency_development', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0013-4000-8000-000000000013');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0014-4000-8000-000000000014', 'competency_development', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0014-4000-8000-000000000014');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0015-4000-8000-000000000015', 'competency_development', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0015-4000-8000-000000000015');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0016-4000-8000-000000000016', 'competency_development', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0016-4000-8000-000000000016');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0017-4000-8000-000000000017', 'people_development', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0017-4000-8000-000000000017');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0018-4000-8000-000000000018', 'people_development', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0018-4000-8000-000000000018');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0019-4000-8000-000000000019', 'people_development', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0019-4000-8000-000000000019');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0020-4000-8000-000000000020', 'people_development', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0020-4000-8000-000000000020');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0021-4000-8000-000000000021', 'people_development', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0021-4000-8000-000000000021');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0022-4000-8000-000000000022', 'people_development', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0022-4000-8000-000000000022');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0023-4000-8000-000000000023', 'people_development', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0023-4000-8000-000000000023');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0024-4000-8000-000000000024', 'people_development', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0024-4000-8000-000000000024');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0025-4000-8000-000000000025', 'career_management', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0025-4000-8000-000000000025');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0026-4000-8000-000000000026', 'career_management', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0026-4000-8000-000000000026');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0027-4000-8000-000000000027', 'career_management', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0027-4000-8000-000000000027');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0028-4000-8000-000000000028', 'career_management', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0028-4000-8000-000000000028');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0029-4000-8000-000000000029', 'career_management', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0029-4000-8000-000000000029');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0030-4000-8000-000000000030', 'career_management', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0030-4000-8000-000000000030');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0031-4000-8000-000000000031', 'career_management', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0031-4000-8000-000000000031');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0032-4000-8000-000000000032', 'career_management', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0032-4000-8000-000000000032');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0033-4000-8000-000000000033', 'hr_assessment', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0033-4000-8000-000000000033');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0034-4000-8000-000000000034', 'hr_assessment', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0034-4000-8000-000000000034');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0035-4000-8000-000000000035', 'hr_assessment', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0035-4000-8000-000000000035');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0036-4000-8000-000000000036', 'hr_assessment', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0036-4000-8000-000000000036');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0037-4000-8000-000000000037', 'hr_assessment', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0037-4000-8000-000000000037');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0038-4000-8000-000000000038', 'hr_assessment', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0038-4000-8000-000000000038');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0039-4000-8000-000000000039', 'hr_assessment', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0039-4000-8000-000000000039');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0040-4000-8000-000000000040', 'hr_assessment', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0040-4000-8000-000000000040');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0041-4000-8000-000000000041', 'recruitement_selection', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0041-4000-8000-000000000041');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0042-4000-8000-000000000042', 'recruitement_selection', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0042-4000-8000-000000000042');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0043-4000-8000-000000000043', 'recruitement_selection', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0043-4000-8000-000000000043');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0044-4000-8000-000000000044', 'recruitement_selection', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0044-4000-8000-000000000044');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0045-4000-8000-000000000045', 'recruitement_selection', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0045-4000-8000-000000000045');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0046-4000-8000-000000000046', 'recruitement_selection', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0046-4000-8000-000000000046');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0047-4000-8000-000000000047', 'recruitement_selection', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0047-4000-8000-000000000047');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0048-4000-8000-000000000048', 'recruitement_selection', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0048-4000-8000-000000000048');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0049-4000-8000-000000000049', 'job_analysis_evaluation', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0049-4000-8000-000000000049');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0050-4000-8000-000000000050', 'job_analysis_evaluation', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0050-4000-8000-000000000050');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0051-4000-8000-000000000051', 'job_analysis_evaluation', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0051-4000-8000-000000000051');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0052-4000-8000-000000000052', 'job_analysis_evaluation', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0052-4000-8000-000000000052');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0053-4000-8000-000000000053', 'job_analysis_evaluation', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0053-4000-8000-000000000053');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0054-4000-8000-000000000054', 'job_analysis_evaluation', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0054-4000-8000-000000000054');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0055-4000-8000-000000000055', 'job_analysis_evaluation', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0055-4000-8000-000000000055');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0056-4000-8000-000000000056', 'job_analysis_evaluation', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0056-4000-8000-000000000056');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0057-4000-8000-000000000057', 'organizational_development', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0057-4000-8000-000000000057');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0058-4000-8000-000000000058', 'organizational_development', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0058-4000-8000-000000000058');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0059-4000-8000-000000000059', 'organizational_development', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0059-4000-8000-000000000059');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0060-4000-8000-000000000060', 'organizational_development', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0060-4000-8000-000000000060');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0061-4000-8000-000000000061', 'organizational_development', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0061-4000-8000-000000000061');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0062-4000-8000-000000000062', 'organizational_development', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0062-4000-8000-000000000062');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0063-4000-8000-000000000063', 'organizational_development', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0063-4000-8000-000000000063');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0064-4000-8000-000000000064', 'organizational_development', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0064-4000-8000-000000000064');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0065-4000-8000-000000000065', 'human_resources_information_system', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0065-4000-8000-000000000065');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0066-4000-8000-000000000066', 'human_resources_information_system', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0066-4000-8000-000000000066');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0067-4000-8000-000000000067', 'human_resources_information_system', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0067-4000-8000-000000000067');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0068-4000-8000-000000000068', 'human_resources_information_system', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0068-4000-8000-000000000068');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0069-4000-8000-000000000069', 'human_resources_information_system', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0069-4000-8000-000000000069');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0070-4000-8000-000000000070', 'human_resources_information_system', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0070-4000-8000-000000000070');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0071-4000-8000-000000000071', 'human_resources_information_system', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0071-4000-8000-000000000071');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0072-4000-8000-000000000072', 'human_resources_information_system', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0072-4000-8000-000000000072');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0073-4000-8000-000000000073', 'workload_analysis', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0073-4000-8000-000000000073');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0074-4000-8000-000000000074', 'workload_analysis', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0074-4000-8000-000000000074');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0075-4000-8000-000000000075', 'workload_analysis', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0075-4000-8000-000000000075');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0076-4000-8000-000000000076', 'workload_analysis', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0076-4000-8000-000000000076');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0077-4000-8000-000000000077', 'workload_analysis', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0077-4000-8000-000000000077');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0078-4000-8000-000000000078', 'workload_analysis', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0078-4000-8000-000000000078');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0079-4000-8000-000000000079', 'workload_analysis', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0079-4000-8000-000000000079');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0080-4000-8000-000000000080', 'workload_analysis', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0080-4000-8000-000000000080');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0081-4000-8000-000000000081', 'performance_apraisal', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0081-4000-8000-000000000081');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0082-4000-8000-000000000082', 'performance_apraisal', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0082-4000-8000-000000000082');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0083-4000-8000-000000000083', 'performance_apraisal', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0083-4000-8000-000000000083');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0084-4000-8000-000000000084', 'performance_apraisal', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0084-4000-8000-000000000084');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0085-4000-8000-000000000085', 'performance_apraisal', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0085-4000-8000-000000000085');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0086-4000-8000-000000000086', 'performance_apraisal', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0086-4000-8000-000000000086');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0087-4000-8000-000000000087', 'performance_apraisal', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0087-4000-8000-000000000087');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0088-4000-8000-000000000088', 'performance_apraisal', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0088-4000-8000-000000000088');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0089-4000-8000-000000000089', 'remuneration_manajemen', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0089-4000-8000-000000000089');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0090-4000-8000-000000000090', 'remuneration_manajemen', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0090-4000-8000-000000000090');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0091-4000-8000-000000000091', 'remuneration_manajemen', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0091-4000-8000-000000000091');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0092-4000-8000-000000000092', 'remuneration_manajemen', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0092-4000-8000-000000000092');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0093-4000-8000-000000000093', 'remuneration_manajemen', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0093-4000-8000-000000000093');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0094-4000-8000-000000000094', 'remuneration_manajemen', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0094-4000-8000-000000000094');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0095-4000-8000-000000000095', 'remuneration_manajemen', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0095-4000-8000-000000000095');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0096-4000-8000-000000000096', 'remuneration_manajemen', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0096-4000-8000-000000000096');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0097-4000-8000-000000000097', 'reward_punisment_management', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0097-4000-8000-000000000097');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0098-4000-8000-000000000098', 'reward_punisment_management', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0098-4000-8000-000000000098');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0099-4000-8000-000000000099', 'reward_punisment_management', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0099-4000-8000-000000000099');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0100-4000-8000-000000000100', 'reward_punisment_management', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0100-4000-8000-000000000100');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0101-4000-8000-000000000101', 'reward_punisment_management', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0101-4000-8000-000000000101');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0102-4000-8000-000000000102', 'reward_punisment_management', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0102-4000-8000-000000000102');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0103-4000-8000-000000000103', 'reward_punisment_management', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0103-4000-8000-000000000103');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0104-4000-8000-000000000104', 'reward_punisment_management', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0104-4000-8000-000000000104');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0105-4000-8000-000000000105', 'health_safety_environment', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0105-4000-8000-000000000105');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0106-4000-8000-000000000106', 'health_safety_environment', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0106-4000-8000-000000000106');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0107-4000-8000-000000000107', 'health_safety_environment', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0107-4000-8000-000000000107');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0108-4000-8000-000000000108', 'health_safety_environment', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0108-4000-8000-000000000108');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0109-4000-8000-000000000109', 'health_safety_environment', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0109-4000-8000-000000000109');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0110-4000-8000-000000000110', 'health_safety_environment', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0110-4000-8000-000000000110');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0111-4000-8000-000000000111', 'health_safety_environment', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0111-4000-8000-000000000111');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0112-4000-8000-000000000112', 'health_safety_environment', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0112-4000-8000-000000000112');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0113-4000-8000-000000000113', 'hubungan_industrial', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0113-4000-8000-000000000113');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0114-4000-8000-000000000114', 'hubungan_industrial', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0114-4000-8000-000000000114');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0115-4000-8000-000000000115', 'hubungan_industrial', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0115-4000-8000-000000000115');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0116-4000-8000-000000000116', 'hubungan_industrial', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0116-4000-8000-000000000116');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0117-4000-8000-000000000117', 'hubungan_industrial', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0117-4000-8000-000000000117');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0118-4000-8000-000000000118', 'hubungan_industrial', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0118-4000-8000-000000000118');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0119-4000-8000-000000000119', 'hubungan_industrial', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0119-4000-8000-000000000119');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0120-4000-8000-000000000120', 'hubungan_industrial', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0120-4000-8000-000000000120');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0121-4000-8000-000000000121', 'budgeting', 1, 'Primary', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0121-4000-8000-000000000121');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0122-4000-8000-000000000122', 'budgeting', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0122-4000-8000-000000000122');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0123-4000-8000-000000000123', 'budgeting', 3, 'Vocational', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0123-4000-8000-000000000123');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0124-4000-8000-000000000124', 'budgeting', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0124-4000-8000-000000000124');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0125-4000-8000-000000000125', 'budgeting', 5, 'Basic Professional', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0125-4000-8000-000000000125');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0126-4000-8000-000000000126', 'budgeting', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0126-4000-8000-000000000126');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0127-4000-8000-000000000127', 'budgeting', 7, 'Professional Mastery', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0127-4000-8000-000000000127');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0128-4000-8000-000000000128', 'budgeting', 8, 'Unique Authority', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0128-4000-8000-000000000128');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0129-4000-8000-000000000129', 'integrity', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0129-4000-8000-000000000129');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0130-4000-8000-000000000130', 'integrity', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0130-4000-8000-000000000130');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0131-4000-8000-000000000131', 'integrity', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0131-4000-8000-000000000131');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0132-4000-8000-000000000132', 'integrity', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0132-4000-8000-000000000132');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0133-4000-8000-000000000133', 'integrity', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0133-4000-8000-000000000133');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0134-4000-8000-000000000134', 'achievement_orientation', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0134-4000-8000-000000000134');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0135-4000-8000-000000000135', 'achievement_orientation', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0135-4000-8000-000000000135');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0136-4000-8000-000000000136', 'achievement_orientation', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0136-4000-8000-000000000136');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0137-4000-8000-000000000137', 'achievement_orientation', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0137-4000-8000-000000000137');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0138-4000-8000-000000000138', 'achievement_orientation', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0138-4000-8000-000000000138');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0139-4000-8000-000000000139', 'building_partnership', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0139-4000-8000-000000000139');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0140-4000-8000-000000000140', 'building_partnership', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0140-4000-8000-000000000140');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0141-4000-8000-000000000141', 'building_partnership', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0141-4000-8000-000000000141');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0142-4000-8000-000000000142', 'building_partnership', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0142-4000-8000-000000000142');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0143-4000-8000-000000000143', 'building_partnership', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0143-4000-8000-000000000143');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0144-4000-8000-000000000144', 'planning_organizing', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0144-4000-8000-000000000144');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0145-4000-8000-000000000145', 'planning_organizing', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0145-4000-8000-000000000145');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0146-4000-8000-000000000146', 'planning_organizing', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0146-4000-8000-000000000146');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0147-4000-8000-000000000147', 'planning_organizing', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0147-4000-8000-000000000147');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0148-4000-8000-000000000148', 'planning_organizing', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0148-4000-8000-000000000148');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0149-4000-8000-000000000149', 'leadership', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0149-4000-8000-000000000149');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0150-4000-8000-000000000150', 'leadership', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0150-4000-8000-000000000150');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0151-4000-8000-000000000151', 'leadership', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0151-4000-8000-000000000151');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0152-4000-8000-000000000152', 'leadership', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0152-4000-8000-000000000152');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0153-4000-8000-000000000153', 'leadership', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0153-4000-8000-000000000153');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0154-4000-8000-000000000154', 'developing_others', 1, 'Task', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0154-4000-8000-000000000154');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0155-4000-8000-000000000155', 'developing_others', 2, 'Supervisory', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0155-4000-8000-000000000155');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0156-4000-8000-000000000156', 'developing_others', 3, 'Managerial', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0156-4000-8000-000000000156');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0157-4000-8000-000000000157', 'developing_others', 4, 'Diverse Managerial', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0157-4000-8000-000000000157');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '44444444-0158-4000-8000-000000000158', 'developing_others', 5, 'Total Managerial', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '44444444-0158-4000-8000-000000000158');
