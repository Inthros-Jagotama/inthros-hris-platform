-- 063_attendance_overtime_approval_instance.sql
-- Approval Module Phase 4: persist the approval_instance_id created for an
-- overtime request so its status can be updated via push-based callback
-- from the central approval module.

ALTER TABLE attendance_overtime_requests
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_att_overtime_approval_instance ON attendance_overtime_requests (approval_instance_id);
