-- =============================================================================
-- Tenant Migration: 096_recruitment_employee_handoff.down (PostgreSQL)
-- =============================================================================

ALTER TABLE candidates
    DROP COLUMN IF EXISTS employee_id,
    DROP COLUMN IF EXISTS candidate_type;

ALTER TABLE employee
    DROP COLUMN IF EXISTS recruited_from_application_id;
