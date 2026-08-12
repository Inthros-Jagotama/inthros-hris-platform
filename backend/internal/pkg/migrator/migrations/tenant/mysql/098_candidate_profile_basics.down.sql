-- =============================================================================
-- Tenant Migration Down: 098_candidate_profile_basics (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_work_experiences;
DROP TABLE IF EXISTS candidate_educations;

SET @drop_candidate_number = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'candidate_number'
  ),
  'ALTER TABLE candidates DROP COLUMN candidate_number',
  'DO 0'
);
PREPARE stmt FROM @drop_candidate_number;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
