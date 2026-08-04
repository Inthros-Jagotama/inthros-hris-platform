-- Migration: 049_seed_job_value_education_experience
-- Menambahkan tipe 'education' (Pendidikan) dan 'experience' (Pengalaman Kerja) ke
-- tabel job_management_values — 5 level masing-masing.
--
-- Latar belakang: migration 045/046 mengubah relasi job_management_education_experiences
-- agar menunjuk ke job_management_values (type='education' / type='experience'),
-- namun belum ada seed untuk kedua tipe tsb — sehingga opsi dropdown di FE kosong
-- kecuali diisi manual. Seeder ini menutup gap tersebut.
--
-- Level mengikuti keputusan user (04 Agu 2026):
--   education  → level 1-5: Sekolah Menengah Pertama → Strata 3 (SMP, SMA, D3, S1, S2/S3)
--   experience → level 1-5: 0-2 Tahun, 3-5 Tahun, 6-8 Tahun, 9-11 Tahun, > 12 Tahun
--   (experience konsisten dengan migration 046 backfill yang mencocokkan descriptions).
--
-- Idempotent GANDA: row di-INSERT hanya jika (a) id (UUID v4 acak) belum ada DAN
-- (b) belum ada row dengan type+level yang sama. Pengaman (b) mencegah duplikat
-- pada tenant yang sudah membuat row manual untuk tipe-tipe ini via UI
-- (tabel tidak punya UNIQUE (type, level)).
--
-- Menggunakan UUID v4 acak (bukan pola '-4000-8000-') agar konsisten dengan
-- migration 043_standardize_job_value_uuids.

-- ── tipe 'education' (Pendidikan) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '129d8602-5eda-4095-b403-2bdcafa80ae0', 'education', 1, 'Sekolah Menengah Pertama (SMP)', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '129d8602-5eda-4095-b403-2bdcafa80ae0')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'education' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'df9019c2-e110-44c3-8815-82fdef3c700a', 'education', 2, 'Sekolah Menengah Atas (SMA)', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'df9019c2-e110-44c3-8815-82fdef3c700a')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'education' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '8f2ef00c-4f6a-4904-88ed-bcdf4917476b', 'education', 3, 'Diploma 3 (D3)', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '8f2ef00c-4f6a-4904-88ed-bcdf4917476b')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'education' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'c1791488-e91b-45cc-99a8-30ab0cb80951', 'education', 4, 'Strata 1 (S1)', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'c1791488-e91b-45cc-99a8-30ab0cb80951')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'education' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'ed02d3d8-8eb8-4c38-a026-931e4a3998ec', 'education', 5, 'Strata 2 / Strata 3 (S2/S3)', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'ed02d3d8-8eb8-4c38-a026-931e4a3998ec')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'education' AND level = 5);

-- ── tipe 'experience' (Pengalaman Kerja) — 5 level ──
INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4e3e2ed3-131d-4ebc-a198-03e77b097cd1', 'experience', 1, '0-2 Tahun', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4e3e2ed3-131d-4ebc-a198-03e77b097cd1')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'experience' AND level = 1);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '4ffa3804-d566-4986-9a53-5282e85506da', 'experience', 2, '3-5 Tahun', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '4ffa3804-d566-4986-9a53-5282e85506da')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'experience' AND level = 2);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'dab6918a-1c0e-48b7-bfe7-e28eff2dca1f', 'experience', 3, '6-8 Tahun', 3, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'dab6918a-1c0e-48b7-bfe7-e28eff2dca1f')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'experience' AND level = 3);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT '62cbaf52-e99a-41fa-8cc2-7e61eb3c4969', 'experience', 4, '9-11 Tahun', 4, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = '62cbaf52-e99a-41fa-8cc2-7e61eb3c4969')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'experience' AND level = 4);

INSERT INTO job_management_values (id, type, level, descriptions, sort, created_at, updated_at)
SELECT 'd89056a0-39da-45a2-a643-04598c630888', 'experience', 5, '> 12 Tahun', 5, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM job_management_values WHERE id = 'd89056a0-39da-45a2-a643-04598c630888')
  AND NOT EXISTS (SELECT 1 FROM job_management_values WHERE type = 'experience' AND level = 5);
