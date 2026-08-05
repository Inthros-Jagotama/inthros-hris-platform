-- 056_performance_kpi_phase2.down.sql
-- Rollback Phase 2 KPI Enhancement: Drop 6 new tables

-- Drop FK constraint from performance_evaluations to performance_ratings first
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_perf_eval_rating'
    ) THEN
        ALTER TABLE performance_evaluations DROP CONSTRAINT fk_perf_eval_rating;
    END IF;
END $$;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS performance_logs;
DROP TABLE IF EXISTS performance_indicator_formulas;
DROP TABLE IF EXISTS performance_ratings;
DROP TABLE IF EXISTS performance_attachments;
DROP TABLE IF EXISTS performance_comments;
DROP TABLE IF EXISTS performance_progress;
