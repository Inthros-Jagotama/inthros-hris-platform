-- Migration: 041_seed_job_value_psychological
-- Menambahkan 6 tipe di bawah grup 'Psychological' ke tabel
-- job_management_values mengikuti docs/seeder/JobManagementValuesTableSeeder.php
-- (codeMap Psychological: Kecerdasan → kecerdasan, Innovation & Creativity →
-- innovation_creativity, Self Confidence → self_confidence, Flexibility →
-- flexibility, Tenacity → tenacity, Continuous Learning → continuous_learning).
--
-- Catatan: teks note seeder di-promote menjadi descriptions per level — pola
-- sama dengan 039/040 (keputusan user).
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

-- ── kecerdasan (5 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0001-4000-8000-000000000001', 'kecerdasan', 1, 'Kurang', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0002-4000-8000-000000000002', 'kecerdasan', 2, 'Cukup', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0003-4000-8000-000000000003', 'kecerdasan', 3, 'Rata-rata', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0004-4000-8000-000000000004', 'kecerdasan', 4, 'Diatas rata-rata', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0005-4000-8000-000000000005', 'kecerdasan', 5, 'Istimewa', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0005-4000-8000-000000000005');

-- ── innovation_creativity (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0006-4000-8000-000000000006', 'innovation_creativity', 1, 'Primary', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0006-4000-8000-000000000006');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0007-4000-8000-000000000007', 'innovation_creativity', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0007-4000-8000-000000000007');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0008-4000-8000-000000000008', 'innovation_creativity', 3, 'Vocational', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0008-4000-8000-000000000008');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0009-4000-8000-000000000009', 'innovation_creativity', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0009-4000-8000-000000000009');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0010-4000-8000-000000000010', 'innovation_creativity', 5, 'Basic Professional', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0010-4000-8000-000000000010');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0011-4000-8000-000000000011', 'innovation_creativity', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0011-4000-8000-000000000011');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0012-4000-8000-000000000012', 'innovation_creativity', 7, 'Professional Mastery', '', 7, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0012-4000-8000-000000000012');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0013-4000-8000-000000000013', 'innovation_creativity', 8, 'Unique Authority', '', 8, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0013-4000-8000-000000000013');

-- ── self_confidence (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0014-4000-8000-000000000014', 'self_confidence', 1, 'Primary', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0014-4000-8000-000000000014');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0015-4000-8000-000000000015', 'self_confidence', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0015-4000-8000-000000000015');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0016-4000-8000-000000000016', 'self_confidence', 3, 'Vocational', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0016-4000-8000-000000000016');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0017-4000-8000-000000000017', 'self_confidence', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0017-4000-8000-000000000017');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0018-4000-8000-000000000018', 'self_confidence', 5, 'Basic Professional', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0018-4000-8000-000000000018');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0019-4000-8000-000000000019', 'self_confidence', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0019-4000-8000-000000000019');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0020-4000-8000-000000000020', 'self_confidence', 7, 'Professional Mastery', '', 7, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0020-4000-8000-000000000020');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0021-4000-8000-000000000021', 'self_confidence', 8, 'Unique Authority', '', 8, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0021-4000-8000-000000000021');

-- ── flexibility (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0022-4000-8000-000000000022', 'flexibility', 1, 'Primary', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0022-4000-8000-000000000022');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0023-4000-8000-000000000023', 'flexibility', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0023-4000-8000-000000000023');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0024-4000-8000-000000000024', 'flexibility', 3, 'Vocational', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0024-4000-8000-000000000024');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0025-4000-8000-000000000025', 'flexibility', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0025-4000-8000-000000000025');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0026-4000-8000-000000000026', 'flexibility', 5, 'Basic Professional', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0026-4000-8000-000000000026');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0027-4000-8000-000000000027', 'flexibility', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0027-4000-8000-000000000027');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0028-4000-8000-000000000028', 'flexibility', 7, 'Professional Mastery', '', 7, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0028-4000-8000-000000000028');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0029-4000-8000-000000000029', 'flexibility', 8, 'Unique Authority', '', 8, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0029-4000-8000-000000000029');

-- ── tenacity (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0030-4000-8000-000000000030', 'tenacity', 1, 'Primary', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0030-4000-8000-000000000030');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0031-4000-8000-000000000031', 'tenacity', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0031-4000-8000-000000000031');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0032-4000-8000-000000000032', 'tenacity', 3, 'Vocational', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0032-4000-8000-000000000032');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0033-4000-8000-000000000033', 'tenacity', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0033-4000-8000-000000000033');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0034-4000-8000-000000000034', 'tenacity', 5, 'Basic Professional', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0034-4000-8000-000000000034');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0035-4000-8000-000000000035', 'tenacity', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0035-4000-8000-000000000035');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0036-4000-8000-000000000036', 'tenacity', 7, 'Professional Mastery', '', 7, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0036-4000-8000-000000000036');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0037-4000-8000-000000000037', 'tenacity', 8, 'Unique Authority', '', 8, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0037-4000-8000-000000000037');

-- ── continuous_learning (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0038-4000-8000-000000000038', 'continuous_learning', 1, 'Primary', '', 1, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0038-4000-8000-000000000038');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0039-4000-8000-000000000039', 'continuous_learning', 2, 'Elementary Vocational', '', 2, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0039-4000-8000-000000000039');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0040-4000-8000-000000000040', 'continuous_learning', 3, 'Vocational', '', 3, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0040-4000-8000-000000000040');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0041-4000-8000-000000000041', 'continuous_learning', 4, 'Advanced Vocational', '', 4, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0041-4000-8000-000000000041');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0042-4000-8000-000000000042', 'continuous_learning', 5, 'Basic Professional', '', 5, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0042-4000-8000-000000000042');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0043-4000-8000-000000000043', 'continuous_learning', 6, 'Seasoned Professional', '', 6, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0043-4000-8000-000000000043');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0044-4000-8000-000000000044', 'continuous_learning', 7, 'Professional Mastery', '', 7, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0044-4000-8000-000000000044');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '33333333-0045-4000-8000-000000000045', 'continuous_learning', 8, 'Unique Authority', '', 8, NOW(), NOW()
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '33333333-0045-4000-8000-000000000045');
