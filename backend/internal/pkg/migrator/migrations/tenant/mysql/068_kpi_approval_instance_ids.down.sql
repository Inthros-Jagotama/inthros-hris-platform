-- 068_kpi_approval_instance_ids.down.sql

SET @drop_realization_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'realization_approval_instance_id'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN realization_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_realization_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_target_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_approval_instance_id'
  ),
  'ALTER TABLE performance_evaluations DROP COLUMN target_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_target_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
