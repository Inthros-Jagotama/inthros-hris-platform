-- =============================================================================
-- Tenant Migration Down: 091_recruitment_workforce_gap (MySQL)
-- =============================================================================

SET @drop_workforce_plan_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'workforce_plan_id'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN workforce_plan_id',
  'DO 0'
);
PREPARE stmt FROM @drop_workforce_plan_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_workforce_gap_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'workforce_gap_id'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN workforce_gap_id',
  'DO 0'
);
PREPARE stmt FROM @drop_workforce_gap_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_reason_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'reason_type'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN reason_type',
  'DO 0'
);
PREPARE stmt FROM @drop_reason_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
