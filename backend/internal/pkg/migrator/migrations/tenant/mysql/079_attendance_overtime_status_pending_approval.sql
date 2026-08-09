-- =============================================================================
-- Tenant Migration: 079_attendance_overtime_status_pending_approval
-- =============================================================================
-- Attendance: kolom `status` di attendance_overtime_requests awalnya adalah
-- ENUM('SUBMITTED','APPROVED','REJECTED'). Sejak CreateOvertimeRequest
-- auto-resolve active flow (via GetActiveFlowIDForModule, pola sama dgn Leave),
-- request lembur yang masuk ke Central Approval Module disimpan dengan status
-- `PENDING_APPROVAL` — nilai yang tidak ada di enum, sehingga MySQL menolak
-- dengan "Data truncated for column 'status'".
--
-- Migration ini menambahkan `PENDING_APPROVAL` ke enum. Idempotent: jika nilai
-- sudah ada (misal tenant di-provision ulang setelah migration diterapkan),
-- statement di-skip.

SET @alter_overtime_status = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND COLUMN_NAME = 'status'
      AND COLUMN_TYPE LIKE '%PENDING_APPROVAL%'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests MODIFY COLUMN status ENUM(''SUBMITTED'',''PENDING_APPROVAL'',''APPROVED'',''REJECTED'') NOT NULL DEFAULT ''SUBMITTED'''
);
PREPARE stmt FROM @alter_overtime_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
