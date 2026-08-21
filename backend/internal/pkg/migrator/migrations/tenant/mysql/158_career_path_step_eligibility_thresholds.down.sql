-- 158_career_path_step_eligibility_thresholds.down.sql (mysql)
ALTER TABLE career_path_steps
    DROP COLUMN min_performance_score,
    DROP COLUMN min_competency_score,
    DROP COLUMN min_okr_score;
