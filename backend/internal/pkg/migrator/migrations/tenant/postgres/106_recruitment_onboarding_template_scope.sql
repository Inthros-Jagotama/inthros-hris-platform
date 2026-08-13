-- =============================================================================
-- Tenant Migration: 106_recruitment_onboarding_template_scope (PostgreSQL)
-- =============================================================================
-- G-10: onboarding_task_templates gets 3 nullable scope columns so a
-- template can be restricted to an organization/position/employment_type
-- (mis. Software Engineer -> Laptop, Repository Access, Security Training).
-- No FK on organization_id/position_id — same convention as
-- job_requisitions.position_id (Organization module doesn't expose
-- position CRUD to FK against). NULL on any field means "applies
-- universally" — CreateEmployeeOnboarding includes a template if each of
-- its scope fields is NULL or matches the target requisition.
-- (docs/module-recruitment-development-plan.md §G-10)
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE onboarding_task_templates
    ADD COLUMN IF NOT EXISTS organization_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS position_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS employment_type VARCHAR(50) NULL;
