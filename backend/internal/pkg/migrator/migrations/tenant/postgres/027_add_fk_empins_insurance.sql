-- Migration: 027_add_fk_empins_insurance
-- Menambahkan DB-level FK constraint fk_empins_insurance pada tabel
-- employee_insurances.insurance_id → insurances.id.
--
-- Dibutuhkan untuk tenant yang sudah ter-record migration 026 (sebelum
-- statement FK ditambahkan ke 026) — idempotent, aman dijalankan ulang.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_empins_insurance'
    ) THEN
        ALTER TABLE employee_insurances
            ADD CONSTRAINT fk_empins_insurance
            FOREIGN KEY (insurance_id) REFERENCES insurances(id)
            ON DELETE SET NULL;
    END IF;
END $$;
