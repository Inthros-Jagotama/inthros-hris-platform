-- =============================================================================
-- Tenant Migration: 096_recruitment_employee_handoff.down (MySQL)
-- =============================================================================

SET @drop_candidate_employee_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'employee_id'
  ),
  'ALTER TABLE candidates DROP COLUMN employee_id',
  'DO 0'
);
PREPARE stmt FROM @drop_candidate_employee_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_candidate_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'candidate_type'
  ),
  'ALTER TABLE candidates DROP COLUMN candidate_type',
  'DO 0'
);
PREPARE stmt FROM @drop_candidate_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_recruited_from_application_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee'
      AND COLUMN_NAME = 'recruited_from_application_id'
  ),
  'ALTER TABLE employee DROP COLUMN recruited_from_application_id',
  'DO 0'
);
PREPARE stmt FROM @drop_recruited_from_application_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
