-- 063_attendance_overtime_approval_instance.down.sql

SET @drop_overtime_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND INDEX_NAME = 'idx_att_overtime_approval_instance'
  ),
  'DROP INDEX idx_att_overtime_approval_instance ON attendance_overtime_requests',
  'DO 0'
);
PREPARE stmt FROM @drop_overtime_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
