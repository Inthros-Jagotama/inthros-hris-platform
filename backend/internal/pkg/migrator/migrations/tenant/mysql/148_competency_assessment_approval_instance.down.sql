-- 148_competency_assessment_approval_instance.down.sql

ALTER TABLE competency_event_targets
    DROP INDEX idx_comp_target_approval_instance,
    DROP INDEX idx_comp_target_status;

ALTER TABLE competency_event_targets
    DROP COLUMN approval_instance_id,
    DROP COLUMN status,
    DROP COLUMN finalized_at;
