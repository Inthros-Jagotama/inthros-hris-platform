-- Down Migration: 051_job_management_scores_is_complete
ALTER TABLE job_management_scores
    DROP COLUMN IF EXISTS completed_at,
    DROP COLUMN IF EXISTS is_complete;
