-- =============================================================================
-- Tenant Migration (down): 080_attendance_overtime_two_flow
-- =============================================================================
-- Kembalikan attendance_overtime_requests ke struktur sebelum dua-alur:
-- drop 13 kolom pendukung + index, revert ENUM status ke nilai lama.
-- Baris yang berstatus WAITING_ACTUAL/ACTUAL_SUBMITTED/CANCELLED akan gagal
-- saat rollback enum — operator harus menormalkan/menghapus baris tsb dulu.

SET @revert_overtime_status = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND COLUMN_NAME = 'status'
      AND COLUMN_TYPE NOT LIKE '%WAITING_ACTUAL%'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests MODIFY COLUMN status ENUM(''SUBMITTED'',''PENDING_APPROVAL'',''APPROVED'',''REJECTED'') NOT NULL DEFAULT ''SUBMITTED'''
);
PREPARE stmt FROM @revert_overtime_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND INDEX_NAME = 'idx_att_overtime_actual_approval_instance'
  ),
  'ALTER TABLE attendance_overtime_requests DROP INDEX idx_att_overtime_actual_approval_instance',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cancelled_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'cancelled_at'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN cancelled_at',
  'DO 0'
);
PREPARE stmt FROM @drop_cancelled_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cancelled_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'cancelled_by'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN cancelled_by',
  'DO 0'
);
PREPARE stmt FROM @drop_cancelled_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approved_at'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_approved_at',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_approved_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approved_by'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_approved_by',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_approved_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_submitted_at'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_submitted_at',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approval_instance_id'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_attachment_url = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'attachment_url'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN attachment_url',
  'DO 0'
);
PREPARE stmt FROM @drop_attachment_url;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_note = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_note'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_note',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_note;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_end_time_local = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_end_time_local'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_end_time_local',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_end_time_local;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_actual_start_time_local = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_start_time_local'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN actual_start_time_local',
  'DO 0'
);
PREPARE stmt FROM @drop_actual_start_time_local;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_assigned_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'assigned_at'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN assigned_at',
  'DO 0'
);
PREPARE stmt FROM @drop_assigned_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_assigned_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'assigned_by'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN assigned_by',
  'DO 0'
);
PREPARE stmt FROM @drop_assigned_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_flow_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'flow_type'
  ),
  'ALTER TABLE attendance_overtime_requests DROP COLUMN flow_type',
  'DO 0'
);
PREPARE stmt FROM @drop_flow_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
