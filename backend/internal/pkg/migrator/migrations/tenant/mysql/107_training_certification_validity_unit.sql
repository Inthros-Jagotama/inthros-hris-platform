-- =============================================================================
-- Tenant Migration: 107_training_certification_validity_unit (MySQL)
-- =============================================================================
-- Menambah kolom validity_period_unit pada training_certifications agar masa
-- berlaku dapat dinyatakan dalam tahun ('year') ATAU bulan ('month').
-- Nilai default 'month' menjaga kompatibilitas data lama.
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_validity_period_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'training_certifications'
      AND COLUMN_NAME = 'validity_period_unit'
  ),
  'DO 0',
  'ALTER TABLE training_certifications ADD COLUMN validity_period_unit VARCHAR(10) NOT NULL DEFAULT ''month'''
);
PREPARE stmt FROM @add_validity_period_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
