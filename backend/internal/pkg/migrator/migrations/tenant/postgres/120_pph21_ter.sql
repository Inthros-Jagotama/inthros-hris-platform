-- =============================================================================
-- Tenant Migration: 120_pph21_ter
-- =============================================================================
-- Adopsi metode TER (PP 58/2023 / PER-2/PJ/2024) di engine PPh21:
--   1. pph21_calculation_logs: catat metode kalkulasi (TER vs REGULAR).
--   2. Perbaiki data master yang salah:
--      a. ptkps."group" (kategori TER) sesuai aturan resmi:
--         A = TK/0, TK/1, K/0 · B = TK/2, TK/3, K/1, K/2 · C = K/3 + K/I/*.
--         Seed lama menempatkan TK/2 & TK/3 di A, K/0 di B, K/3 di B.
--      b. ters grup C: baris pertama (0 s.d. 6.600.000 → 0%) hilang dan
--         seluruh rate bergeser 1 tingkat lebih rendah dari tabel resmi.
--         Koreksi: insert baris 0-6,6jt (0%) + naikkan rate baris lainnya 0,25%.

-- ---------------------------------------------------------------------------
-- 120.1 pph21_calculation_logs: kolom metode kalkulasi
-- ---------------------------------------------------------------------------
ALTER TABLE pph21_calculation_logs
    ADD COLUMN calculation_method VARCHAR(255) NOT NULL DEFAULT 'REGULAR_GROSS_ANNUALIZED';

-- ---------------------------------------------------------------------------
-- 120.2 ptkps: koreksi kategori TER sesuai aturan resmi
-- ---------------------------------------------------------------------------
UPDATE ptkps SET "group" = 'A' WHERE id IN (SELECT id FROM ptkps WHERE name LIKE '%(K/0)%');
UPDATE ptkps SET "group" = 'B' WHERE id IN (SELECT id FROM ptkps WHERE name LIKE '%(TK/2)%');
UPDATE ptkps SET "group" = 'B' WHERE id IN (SELECT id FROM ptkps WHERE name LIKE '%(TK/3)%');
UPDATE ptkps SET "group" = 'C' WHERE id IN (SELECT id FROM ptkps WHERE name LIKE '%(K/3)%');

-- ---------------------------------------------------------------------------
-- 120.3 ters grup C: koreksi tabel tarif efektif bulanan
-- ---------------------------------------------------------------------------
-- Baris pertama yang hilang: 0 s.d. 6.600.000 → 0%.
INSERT INTO ters (id, "group", bruto_min, bruto_max, rate, created_at, updated_at)
SELECT gen_random_uuid(), 'C', 0, 6600000, 0.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM ters WHERE "group" = 'C' AND bruto_min <= 0 AND bruto_max >= 6600000);

-- Rate seluruh baris C lainnya digeser naik 0,25% (koreksi pergeseran).
UPDATE ters SET rate = rate + 0.25 WHERE "group" = 'C' AND bruto_min > 0;
