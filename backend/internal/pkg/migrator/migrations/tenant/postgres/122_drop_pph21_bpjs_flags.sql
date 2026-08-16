-- ---------------------------------------------------------------------------
-- 122. Hapus kolom deduct_bpjs_* dari pph21_settings — aturan pengurang iuran
--      BPJS kini di-hardcode sesuai regulasi (JHT & JP boleh dikurangkan,
--      BPJS Kesehatan tidak). Lihat engine pph21.go.
-- ---------------------------------------------------------------------------

ALTER TABLE pph21_settings DROP COLUMN deduct_bpjs_health_employee;
ALTER TABLE pph21_settings DROP COLUMN deduct_bpjs_jht_employee;
ALTER TABLE pph21_settings DROP COLUMN deduct_bpjs_jp_employee;
