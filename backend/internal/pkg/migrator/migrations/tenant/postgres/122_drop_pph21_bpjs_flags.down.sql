-- Down: kembalikan kolom deduct_bpjs_* dengan default lama (JHT & JP = true,
-- Kesehatan = false — aturan baku).
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_health_employee BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_jht_employee BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE pph21_settings ADD COLUMN deduct_bpjs_jp_employee BOOLEAN NOT NULL DEFAULT TRUE;
