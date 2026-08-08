-- 069_okr_two_phase.sql
-- Two-phase OKR flow: employee proposes Key Results (DRAFT -> KR_SUBMITTED ->
-- KR_APPROVED, "OKR Active") before self-assessment (KR_APPROVED -> SUBMITTED
-- -> COMPLETED). Two independent approval-module checkpoints per evaluation,
-- modules "okr_key_result" and "okr_assessment" — mirrors
-- 068_kpi_approval_instance_ids.sql for KPI.

SET @add_kr_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'kr_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE okr_evaluations ADD COLUMN kr_approval_instance_id CHAR(36) NULL AFTER approved_by'
);
PREPARE stmt FROM @add_kr_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_assessment_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'assessment_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE okr_evaluations ADD COLUMN assessment_approval_instance_id CHAR(36) NULL AFTER kr_approval_instance_id'
);
PREPARE stmt FROM @add_assessment_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_kr_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'kr_submitted_at'
  ),
  'DO 0',
  'ALTER TABLE okr_evaluations ADD COLUMN kr_submitted_at TIMESTAMP NULL AFTER assessment_approval_instance_id'
);
PREPARE stmt FROM @add_kr_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
