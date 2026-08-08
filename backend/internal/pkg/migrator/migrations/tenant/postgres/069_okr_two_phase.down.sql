-- 069_okr_two_phase.down.sql

ALTER TABLE okr_evaluations
    DROP COLUMN IF EXISTS kr_submitted_at,
    DROP COLUMN IF EXISTS assessment_approval_instance_id,
    DROP COLUMN IF EXISTS kr_approval_instance_id;
