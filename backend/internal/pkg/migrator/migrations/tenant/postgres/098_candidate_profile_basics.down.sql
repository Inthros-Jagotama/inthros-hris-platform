-- =============================================================================
-- Tenant Migration Down: 098_candidate_profile_basics (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_work_experiences;
DROP TABLE IF EXISTS candidate_educations;
ALTER TABLE candidates DROP COLUMN IF EXISTS candidate_number;
