-- 064_seed_performance_components.sql
-- KPI Enhancement: seed the 3 fixed Performance Components (Indikator/KPI,
-- Program, Kinerja Bawahan/Subordinate). The scoring engine
-- (CalculateEvaluationComponentScoring) only knows how to auto-calculate
-- these exact codes (KPI, PROGRAM, SUBORDINATE) — the module is locked to
-- exactly these 3 rows (create/delete routes removed at the API level).
--
-- Idempotent: inserted only if a row with that code doesn't already exist.

INSERT INTO performance_components (id, code, name, description, sort_order, is_active, created_at, updated_at)
SELECT 'a1b2c3d4-0001-4000-8000-000000000001', 'KPI', 'Indikator', 'Indikator Kinerja Utama (KPI) sesuai template posisi jabatan.', 1, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM performance_components WHERE code = 'KPI');

INSERT INTO performance_components (id, code, name, description, sort_order, is_active, created_at, updated_at)
SELECT 'a1b2c3d4-0002-4000-8000-000000000002', 'PROGRAM', 'Program', 'Program kerja yang diajukan sendiri oleh karyawan pada setiap evaluasi.', 2, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM performance_components WHERE code = 'PROGRAM');

INSERT INTO performance_components (id, code, name, description, sort_order, is_active, created_at, updated_at)
SELECT 'a1b2c3d4-0003-4000-8000-000000000003', 'SUBORDINATE', 'Kinerja Bawahan', 'Dihitung otomatis dari rata-rata skor akhir evaluasi karyawan yang melapor ke posisi ini.', 3, true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM performance_components WHERE code = 'SUBORDINATE');
