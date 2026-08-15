-- ---------------------------------------------------------------------------
-- 121. Konsolidasi PTKP: satu tabel `ptkps` (setting module) sebagai sumber
--      kebenaran tunggal untuk engine PPh21. Tabel pph21_ptkp_rates (payroll)
--      yang duplikat dihapus.
-- ---------------------------------------------------------------------------

-- 1. Kolom `code` di ptkps untuk lookup engine (normalisasi status, mis. "TK0").
ALTER TABLE ptkps ADD COLUMN code VARCHAR(20) NOT NULL DEFAULT '';

-- 2. Backfill code dari nama — pola "(TK/0)" di akhir nama → "TK0".
UPDATE ptkps
SET code = UPPER(REPLACE(SUBSTRING(name FROM '\\(([^)]+)\\)'), '/', ''))
WHERE code = '';

-- 3. Drop tabel duplikat pph21_ptkp_rates (digantikan ptkps).
DROP TABLE IF EXISTS pph21_ptkp_rates;
