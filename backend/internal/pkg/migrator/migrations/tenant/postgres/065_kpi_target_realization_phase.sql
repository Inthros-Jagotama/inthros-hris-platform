-- 065_kpi_target_realization_phase.sql
-- KPI Enhancement Phase 2: two-phase target/realization submission.
-- Status flow: DRAFT -> TARGET_SUBMITTED -> TARGET_APPROVED -> SUBMITTED ->
-- APPROVED -> COMPLETED (status column itself is unchanged, plain varchar).
-- Adds timestamps for the new target-approval checkpoint, mirroring the
-- existing submitted_at/approved_at columns for the realization checkpoint.

ALTER TABLE performance_evaluations
    ADD COLUMN IF NOT EXISTS target_submitted_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS target_approved_at TIMESTAMP NULL;
