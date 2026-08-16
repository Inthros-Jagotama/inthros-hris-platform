-- 148_competency_assessment_approval_instance.down.sql

DROP INDEX IF EXISTS idx_comp_target_approval_instance;

DROP INDEX IF EXISTS idx_comp_target_status;

ALTER TABLE competency_event_targets
    DROP COLUMN IF EXISTS approval_instance_id,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS finalized_at;
