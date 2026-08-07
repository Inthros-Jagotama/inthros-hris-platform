-- 065_kpi_target_realization_phase.down.sql

SET @drop_target_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_approved_at'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN target_approved_at',
  'DO 0'
);
PREPARE stmt FROM @drop_target_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_target_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_submitted_at'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN target_submitted_at',
  'DO 0'
);
PREPARE stmt FROM @drop_target_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
