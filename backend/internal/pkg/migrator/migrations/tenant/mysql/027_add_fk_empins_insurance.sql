-- Migration: 027_add_fk_empins_insurance
-- Menambahkan DB-level FK constraint fk_empins_insurance pada tabel
-- employee_insurances.insurance_id → insurances.id.
--
-- Dibutuhkan untuk tenant yang sudah ter-record migration 026 (sebelum
-- statement FK ditambahkan ke 026) — idempotent, aman dijalankan ulang.

SET @fk_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_insurances'
      AND CONSTRAINT_NAME = 'fk_empins_insurance'
  ),
  'DO 0',
  'ALTER TABLE employee_insurances ADD CONSTRAINT fk_empins_insurance FOREIGN KEY (insurance_id) REFERENCES insurances(id) ON DELETE SET NULL'
);
PREPARE stmt_fk FROM @fk_sql;
EXECUTE stmt_fk;
DEALLOCATE PREPARE stmt_fk;
