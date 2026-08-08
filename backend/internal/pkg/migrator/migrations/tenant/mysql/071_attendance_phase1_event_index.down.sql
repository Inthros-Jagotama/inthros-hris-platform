SET @drop_att_event_employee_time_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_events'
      AND index_name = 'idx_att_event_employee_time'
  ),
  'ALTER TABLE attendance_events DROP INDEX idx_att_event_employee_time',
  'DO 0'
);
PREPARE stmt FROM @drop_att_event_employee_time_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
