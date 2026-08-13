-- =============================================================================
-- Tenant Migration Down: 106_recruitment_onboarding_template_scope (PostgreSQL)
-- =============================================================================

ALTER TABLE onboarding_task_templates
    DROP COLUMN IF EXISTS employment_type,
    DROP COLUMN IF EXISTS position_id,
    DROP COLUMN IF EXISTS organization_id;
