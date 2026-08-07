-- 065_kpi_target_realization_phase.down.sql

ALTER TABLE performance_evaluations
    DROP COLUMN IF EXISTS target_approved_at,
    DROP COLUMN IF EXISTS target_submitted_at;
