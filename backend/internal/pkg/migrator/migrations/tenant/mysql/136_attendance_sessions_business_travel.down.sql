-- 136_attendance_sessions_business_travel.down.sql

SET @drop_att_session_biztrav_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_sessions'
      AND INDEX_NAME = 'idx_att_session_biztrav'
  ),
  'DROP INDEX idx_att_session_biztrav ON attendance_sessions',
  'DO 0'
);
PREPARE stmt FROM @drop_att_session_biztrav_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_att_session_biztrav_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_sessions'
      AND COLUMN_NAME = 'business_travel_id'
  ),
  'ALTER TABLE attendance_sessions DROP COLUMN business_travel_id',
  'DO 0'
);
PREPARE stmt FROM @drop_att_session_biztrav_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
