-- 136_attendance_sessions_business_travel.down.sql

DROP INDEX IF EXISTS idx_att_session_biztrav;

ALTER TABLE attendance_sessions
    DROP COLUMN IF EXISTS business_travel_id;
