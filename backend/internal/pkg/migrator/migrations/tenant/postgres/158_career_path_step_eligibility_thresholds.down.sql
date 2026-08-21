-- 158_career_path_step_eligibility_thresholds.down.sql (postgres)
ALTER TABLE career_path_steps
    DROP COLUMN IF EXISTS min_performance_score,
    DROP COLUMN IF EXISTS min_competency_score,
    DROP COLUMN IF EXISTS min_okr_score;
