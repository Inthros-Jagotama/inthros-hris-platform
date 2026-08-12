-- =============================================================================
-- Tenant Migration Down: 093_recruitment_approval (PostgreSQL)
-- =============================================================================

ALTER TABLE job_requisitions
    DROP COLUMN IF EXISTS approval_instance_id;
