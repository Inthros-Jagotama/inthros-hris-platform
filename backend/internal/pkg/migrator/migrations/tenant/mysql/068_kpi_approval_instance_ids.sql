-- 068_kpi_approval_instance_ids.sql
-- Route KPI target/realization approval through the central approval
-- module: two separate approval instances per evaluation (target approval
-- and realization approval are independent checkpoints, potentially with
-- different flows/approvers), modules "performance_kpi_target" and
-- "performance_kpi_realization".

SET @add_target_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'target_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN target_approval_instance_id CHAR(36) NULL AFTER target_approved_at'
);
PREPARE stmt FROM @add_target_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_realization_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluations'
      AND column_name = 'realization_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluations ADD COLUMN realization_approval_instance_id CHAR(36) NULL AFTER target_approval_instance_id'
);
PREPARE stmt FROM @add_realization_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
