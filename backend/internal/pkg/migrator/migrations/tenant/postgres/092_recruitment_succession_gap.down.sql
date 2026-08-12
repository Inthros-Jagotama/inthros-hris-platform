-- =============================================================================
-- Tenant Migration Down: 092_recruitment_succession_gap (PostgreSQL)
-- =============================================================================

ALTER TABLE job_requisitions
    DROP COLUMN IF EXISTS succession_position_id;
