-- =============================================================================
-- Tenant Migration Down: 107_training_certification_validity_unit (MySQL)
-- =============================================================================

SET @drop_validity_period_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'training_certifications'
      AND COLUMN_NAME = 'validity_period_unit'
  ),
  'ALTER TABLE training_certifications DROP COLUMN validity_period_unit',
  'DO 0'
);
PREPARE stmt FROM @drop_validity_period_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
