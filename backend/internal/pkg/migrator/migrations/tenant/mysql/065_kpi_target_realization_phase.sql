-- 065_kpi_target_realization_phase.sql
-- KPI Enhancement Phase 2: two-phase target/realization submission.
-- Status flow: DRAFT -> TARGET_SUBMITTED -> TARGET_APPROVED -> SUBMITTED ->
-- APPROVED -> COMPLETED (status column itself is unchanged, plain varchar).
-- Adds timestamps for the new target-approval checkpoint, mirroring the
-- existing submitted_at/approved_at columns for the realization checkpoint.

SET @add_target_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_submitted_at'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN target_submitted_at TIMESTAMP(6) NULL AFTER approved_at'
);
PREPARE stmt FROM @add_target_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_target_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_approved_at'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN target_approved_at TIMESTAMP(6) NULL AFTER target_submitted_at'
);
PREPARE stmt FROM @add_target_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
