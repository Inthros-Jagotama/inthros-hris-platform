-- Rollback 120: kembalikan data master & hapus kolom metode kalkulasi.

-- 120.3 ters grup C: kembalikan rate (geser turun 0,25%) & hapus baris 0-6,6jt.
UPDATE ters SET rate = rate - 0.25 WHERE "group" = 'C' AND bruto_min > 0;
DELETE FROM ters WHERE "group" = 'C' AND bruto_min <= 0 AND bruto_max >= 6600000;

-- 120.2 ptkps: kembalikan kategori lama (kondisi sebelum koreksi).
UPDATE ptkps SET "group" = 'B' WHERE name LIKE '%(K/0)%';
UPDATE ptkps SET "group" = 'A' WHERE name LIKE '%(TK/2)%';
UPDATE ptkps SET "group" = 'A' WHERE name LIKE '%(TK/3)%';
UPDATE ptkps SET "group" = 'B' WHERE name LIKE '%(K/3)%';

-- 120.1 pph21_calculation_logs: hapus kolom metode kalkulasi.
ALTER TABLE pph21_calculation_logs
    DROP COLUMN IF EXISTS calculation_method;
