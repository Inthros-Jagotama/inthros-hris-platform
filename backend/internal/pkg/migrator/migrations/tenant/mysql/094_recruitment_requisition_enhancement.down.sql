-- =============================================================================
-- Tenant Migration Down: 094_recruitment_requisition_enhancement (MySQL)
-- =============================================================================

SET @drop_opened_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'opened_at'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN opened_at',
  'DO 0'
);
PREPARE stmt FROM @drop_opened_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'position_id'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN position_id',
  'DO 0'
);
PREPARE stmt FROM @drop_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_priority = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'priority'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN priority',
  'DO 0'
);
PREPARE stmt FROM @drop_priority;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_requisition_number = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'requisition_number'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN requisition_number',
  'DO 0'
);
PREPARE stmt FROM @drop_requisition_number;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
