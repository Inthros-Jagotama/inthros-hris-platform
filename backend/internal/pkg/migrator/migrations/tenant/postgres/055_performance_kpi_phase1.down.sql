-- 055_performance_kpi_phase1.down.sql
-- Rollback Phase 1 KPI Enhancement

-- Drop FK constraints first
ALTER TABLE performance_templates DROP CONSTRAINT IF EXISTS fk_perf_tmpl_period;
ALTER TABLE performance_evaluation_details DROP CONSTRAINT IF EXISTS fk_perf_detail_indicator;

-- Drop indexes
DROP INDEX IF EXISTS idx_perf_tmpl_period;
DROP INDEX IF EXISTS idx_perf_tmpl_effective;
DROP INDEX IF EXISTS idx_perf_ind_code;
DROP INDEX IF EXISTS idx_perf_eval_rating;
DROP INDEX IF EXISTS idx_perf_detail_indicator;

-- Remove columns from performance_templates
ALTER TABLE performance_templates
    DROP COLUMN IF EXISTS period_id,
    DROP COLUMN IF EXISTS effective_date,
    DROP COLUMN IF EXISTS expired_date;

-- Remove columns from performance_indicators
ALTER TABLE performance_indicators
    DROP COLUMN IF EXISTS code,
    DROP COLUMN IF EXISTS formula_type,
    DROP COLUMN IF EXISTS minimum_score,
    DROP COLUMN IF EXISTS maximum_score,
    DROP COLUMN IF EXISTS target_type,
    DROP COLUMN IF EXISTS is_required;

-- Remove columns from performance_evaluations
ALTER TABLE performance_evaluations
    DROP COLUMN IF EXISTS rating_id,
    DROP COLUMN IF EXISTS submitted_at,
    DROP COLUMN IF EXISTS approved_at;

-- Remove columns from performance_evaluation_details
ALTER TABLE performance_evaluation_details
    DROP COLUMN IF EXISTS indicator_id,
    DROP COLUMN IF EXISTS indicator_name,
    DROP COLUMN IF EXISTS target,
    DROP COLUMN IF EXISTS actual,
    DROP COLUMN IF EXISTS achievement,
    DROP COLUMN IF EXISTS remarks;
