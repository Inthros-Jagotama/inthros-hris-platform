-- Down: kembalikan kolom deduct_bpjs_* dengan default lama (JHT & JP = true,
-- Kesehatan = false — aturan baku).
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_health_employee TINYINT(1) NOT NULL DEFAULT 0;
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_jht_employee TINYINT(1) NOT NULL DEFAULT 1;
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_jp_employee TINYINT(1) NOT NULL DEFAULT 1;
