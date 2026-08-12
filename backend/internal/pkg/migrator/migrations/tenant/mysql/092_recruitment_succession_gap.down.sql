-- =============================================================================
-- Tenant Migration Down: 092_recruitment_succession_gap (MySQL)
-- =============================================================================

SET @drop_succession_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'succession_position_id'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN succession_position_id',
  'DO 0'
);
PREPARE stmt FROM @drop_succession_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
