-- 069_okr_two_phase.down.sql

SET @drop_kr_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'kr_submitted_at'
  ),
  'ALTER TABLE okr_evaluations DROP COLUMN kr_submitted_at',
  'DO 0'
);
PREPARE stmt FROM @drop_kr_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_assessment_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'assessment_approval_instance_id'
  ),
  'ALTER TABLE okr_evaluations DROP COLUMN assessment_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_assessment_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_kr_approval_instance = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_evaluations'
      AND column_name = 'kr_approval_instance_id'
  ),
  'ALTER TABLE okr_evaluations DROP COLUMN kr_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_kr_approval_instance;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
