-- 136_attendance_sessions_business_travel.sql
-- Business Travel becomes an Attendance source (§37 plan doc): dedicated
-- column on attendance_sessions, mirroring the leave_request_id/
-- leave_fraction pattern (migration 004) rather than a generic
-- source_type/source_id design, per Service.ApplyApprovedLeave's existing
-- convention (docs/module-attendance-business-travel-development-plan.md §54).

SET @add_att_session_biztrav_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_sessions'
      AND COLUMN_NAME = 'business_travel_id'
  ),
  'DO 0',
  'ALTER TABLE attendance_sessions ADD COLUMN business_travel_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_att_session_biztrav_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_att_session_biztrav_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_sessions'
      AND INDEX_NAME = 'idx_att_session_biztrav'
  ),
  'DO 0',
  'CREATE INDEX idx_att_session_biztrav ON attendance_sessions (business_travel_id)'
);
PREPARE stmt FROM @add_att_session_biztrav_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
