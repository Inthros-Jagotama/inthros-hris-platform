SET @drop_calculated_minutes = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'calculated_minutes'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN calculated_minutes',
  'DO 0'
);
PREPARE stmt FROM @drop_calculated_minutes;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_minutes = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_minutes'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_minutes',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_minutes;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
