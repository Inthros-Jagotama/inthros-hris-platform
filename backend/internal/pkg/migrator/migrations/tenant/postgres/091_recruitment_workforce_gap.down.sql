-- =============================================================================
-- Tenant Migration Down: 091_recruitment_workforce_gap (PostgreSQL)
-- =============================================================================

ALTER TABLE job_requisitions
    DROP COLUMN IF EXISTS workforce_plan_id,
    DROP COLUMN IF EXISTS workforce_gap_id,
    DROP COLUMN IF EXISTS reason_type;
