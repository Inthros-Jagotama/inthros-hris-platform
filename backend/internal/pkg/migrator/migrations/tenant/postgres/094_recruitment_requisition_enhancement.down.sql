-- =============================================================================
-- Tenant Migration Down: 094_recruitment_requisition_enhancement (PostgreSQL)
-- =============================================================================

ALTER TABLE job_requisitions
    DROP COLUMN IF EXISTS opened_at,
    DROP COLUMN IF EXISTS position_id,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS requisition_number;
