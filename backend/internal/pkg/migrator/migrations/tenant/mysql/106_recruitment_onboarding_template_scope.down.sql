-- =============================================================================
-- Tenant Migration Down: 106_recruitment_onboarding_template_scope (MySQL)
-- =============================================================================

SET @drop_employment_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'employment_type'
  ),
  'ALTER TABLE onboarding_task_templates DROP COLUMN employment_type',
  'DO 0'
);
PREPARE stmt FROM @drop_employment_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'position_id'
  ),
  'ALTER TABLE onboarding_task_templates DROP COLUMN position_id',
  'DO 0'
);
PREPARE stmt FROM @drop_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_organization_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'organization_id'
  ),
  'ALTER TABLE onboarding_task_templates DROP COLUMN organization_id',
  'DO 0'
);
PREPARE stmt FROM @drop_organization_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
