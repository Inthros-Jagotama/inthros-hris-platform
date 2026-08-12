-- =============================================================================
-- Tenant Migration Down: 093_recruitment_approval (MySQL)
-- =============================================================================

SET @drop_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE job_requisitions DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
