-- Down Migration: 051_job_management_scores_is_complete
ALTER TABLE job_management_scores
    DROP COLUMN completed_at,
    DROP COLUMN is_complete;
