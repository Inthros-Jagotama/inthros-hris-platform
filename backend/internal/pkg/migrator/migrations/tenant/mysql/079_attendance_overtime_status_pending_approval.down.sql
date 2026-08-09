-- =============================================================================
-- Tenant Migration (down): 079_attendance_overtime_status_pending_approval
-- =============================================================================
-- Kembalikan ENUM status attendance_overtime_requests ke nilai awal
-- (SUBMITTED/APPROVED/REJECTED). Baris yang berstatus PENDING_APPROVAL akan
-- gagal saat rollback — operator harus menormalkan/menghapus baris tsb dulu.

SET @revert_overtime_status = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND COLUMN_NAME = 'status'
      AND COLUMN_TYPE NOT LIKE '%PENDING_APPROVAL%'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests MODIFY COLUMN status ENUM(''SUBMITTED'',''APPROVED'',''REJECTED'') NOT NULL DEFAULT ''SUBMITTED'''
);
PREPARE stmt FROM @revert_overtime_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
