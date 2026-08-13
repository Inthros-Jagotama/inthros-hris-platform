-- =============================================================================
-- Tenant Migration: 106_recruitment_onboarding_template_scope (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_organization_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'organization_id'
  ),
  'DO 0',
  'ALTER TABLE onboarding_task_templates ADD COLUMN organization_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_organization_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'position_id'
  ),
  'DO 0',
  'ALTER TABLE onboarding_task_templates ADD COLUMN position_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_employment_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'onboarding_task_templates'
      AND COLUMN_NAME = 'employment_type'
  ),
  'DO 0',
  'ALTER TABLE onboarding_task_templates ADD COLUMN employment_type VARCHAR(50) NULL'
);
PREPARE stmt FROM @add_employment_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
