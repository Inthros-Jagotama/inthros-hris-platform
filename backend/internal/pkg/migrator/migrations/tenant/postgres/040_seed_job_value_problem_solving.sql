-- Migration: 040_seed_job_value_problem_solving
-- Menambahkan tipe 'thinking_environment' (8 level) dan 'thinking_chalenge'
-- (5 level) — kedua tipe di bawah grup 'Problem Solving & Decision Making'
-- mengikuti docs/seeder/JobManagementValuesTableSeeder.php (parent
-- 'Problem Solving & Decision Making', codeMap: Thinking Environment →
-- thinking_environment, Thinking Chalenge → thinking_chalenge).
--
-- Catatan: teks note seeder (Berulang-ulang, Bermotif, Variabel, Adaptif,
-- Belum dipetakan, dst.) di-promote menjadi descriptions per level — pola
-- sama dengan 039 communicating_influencing_skill (keputusan user).
--
-- Idempotent: setiap row di-INSERT hanya jika id (UUID tetap) belum ada.

-- ── thinking_environment (8 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0001-4000-8000-000000000001', 'thinking_environment', 1, 'Berulang-ulang', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0001-4000-8000-000000000001');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0002-4000-8000-000000000002', 'thinking_environment', 2, 'Bermotif', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0002-4000-8000-000000000002');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0003-4000-8000-000000000003', 'thinking_environment', 3, 'Variabel', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0003-4000-8000-000000000003');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0004-4000-8000-000000000004', 'thinking_environment', 4, 'Adaptif', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0004-4000-8000-000000000004');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0005-4000-8000-000000000005', 'thinking_environment', 5, 'Belum dipetakan', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0005-4000-8000-000000000005');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0006-4000-8000-000000000006', 'thinking_environment', 6, 'Didefinisikan secara luas', '', 6, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0006-4000-8000-000000000006');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0007-4000-8000-000000000007', 'thinking_environment', 7, 'Didefinisikan Secara Umum', '', 7, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0007-4000-8000-000000000007');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0008-4000-8000-000000000008', 'thinking_environment', 8, 'Didefinisikan Secara Abstrak', '', 8, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0008-4000-8000-000000000008');

-- ── thinking_chalenge (5 level) ──
INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0009-4000-8000-000000000009', 'thinking_chalenge', 1, 'Berulang-ulang', '', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0009-4000-8000-000000000009');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0010-4000-8000-000000000010', 'thinking_chalenge', 2, 'Bermotif', '', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0010-4000-8000-000000000010');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0011-4000-8000-000000000011', 'thinking_chalenge', 3, 'Variabel', '', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0011-4000-8000-000000000011');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0012-4000-8000-000000000012', 'thinking_chalenge', 4, 'Adaptif', '', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0012-4000-8000-000000000012');

INSERT INTO job_management_values (id, type, level, descriptions, note, sort, created_at, updated_at)
SELECT '22222222-0013-4000-8000-000000000013', 'thinking_chalenge', 5, 'Belum dipetakan', '', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '22222222-0013-4000-8000-000000000013');
