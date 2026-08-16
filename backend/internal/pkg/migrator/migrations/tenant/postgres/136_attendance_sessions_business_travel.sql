-- 136_attendance_sessions_business_travel.sql
-- Business Travel becomes an Attendance source (§37 plan doc): dedicated
-- column on attendance_sessions, mirroring the leave_request_id/
-- leave_fraction pattern (migration 004) rather than a generic
-- source_type/source_id design, per Service.ApplyApprovedLeave's existing
-- convention (docs/module-attendance-business-travel-development-plan.md §54).

ALTER TABLE attendance_sessions
    ADD COLUMN IF NOT EXISTS business_travel_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_att_session_biztrav ON attendance_sessions (business_travel_id);
