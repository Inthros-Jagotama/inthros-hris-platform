-- 063_attendance_overtime_approval_instance.down.sql

DROP INDEX IF EXISTS idx_att_overtime_approval_instance;

ALTER TABLE attendance_overtime_requests
    DROP COLUMN IF EXISTS approval_instance_id;
