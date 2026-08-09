-- =============================================================================
-- Tenant Migration: 080_attendance_overtime_two_flow
-- =============================================================================
-- Attendance Module §32b (docs/module-attendance-plan.md) — Overtime dua alur:
--   SELF (ajukan sendiri → approval #1 → isian aktual → approval #2) dan
--   ASSIGNED (atasan tugaskan bawahan → notifikasi → isian aktual → approval).
--   Kedua alur bisa dibatalkan sebelum isian aktual.
--
-- Migration ini menambahkan kolom pendukung alur tersebut ke
-- attendance_overtime_requests dan memperluas ENUM status dengan
-- WAITING_ACTUAL, ACTUAL_SUBMITTED, dan CANCELLED. Semua statement idempotent
-- (di-skip jika kolom/nilai sudah ada).

SET @add_flow_type = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'flow_type'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN flow_type VARCHAR(20) NOT NULL DEFAULT ''SELF'' AFTER approval_instance_id'
);
PREPARE stmt FROM @add_flow_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_assigned_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'assigned_by'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN assigned_by CHAR(36) NULL AFTER flow_type'
);
PREPARE stmt FROM @add_assigned_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_assigned_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'assigned_at'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN assigned_at DATETIME(6) NULL AFTER assigned_by'
);
PREPARE stmt FROM @add_assigned_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_start_time_local = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_start_time_local'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_start_time_local DATETIME(6) NULL AFTER assigned_at'
);
PREPARE stmt FROM @add_actual_start_time_local;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_end_time_local = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_end_time_local'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_end_time_local DATETIME(6) NULL AFTER actual_start_time_local'
);
PREPARE stmt FROM @add_actual_end_time_local;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_note = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_note'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_note VARCHAR(500) NULL AFTER actual_end_time_local'
);
PREPARE stmt FROM @add_actual_note;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_attachment_url = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'attachment_url'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN attachment_url VARCHAR(500) NULL AFTER actual_note'
);
PREPARE stmt FROM @add_attachment_url;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_approval_instance_id CHAR(36) NULL AFTER attachment_url'
);
PREPARE stmt FROM @add_actual_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_submitted_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_submitted_at'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_submitted_at DATETIME(6) NULL AFTER actual_approval_instance_id'
);
PREPARE stmt FROM @add_actual_submitted_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_approved_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approved_by'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_approved_by CHAR(36) NULL AFTER actual_submitted_at'
);
PREPARE stmt FROM @add_actual_approved_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_actual_approved_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_approved_at'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_approved_at DATETIME(6) NULL AFTER actual_approved_by'
);
PREPARE stmt FROM @add_actual_approved_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cancelled_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'cancelled_by'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN cancelled_by CHAR(36) NULL AFTER actual_approved_at'
);
PREPARE stmt FROM @add_cancelled_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cancelled_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'cancelled_at'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN cancelled_at DATETIME(6) NULL AFTER cancelled_by'
);
PREPARE stmt FROM @add_cancelled_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Index untuk lookup dispatch callback approval instance kedua
SET @add_actual_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND INDEX_NAME = 'idx_att_overtime_actual_approval_instance'
  ),
  'DO 0',
  'CREATE INDEX idx_att_overtime_actual_approval_instance ON attendance_overtime_requests (actual_approval_instance_id)'
);
PREPARE stmt FROM @add_actual_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Perluas ENUM status: tambah WAITING_ACTUAL, ACTUAL_SUBMITTED, CANCELLED
SET @alter_overtime_status = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'attendance_overtime_requests'
      AND COLUMN_NAME = 'status'
      AND COLUMN_TYPE LIKE '%WAITING_ACTUAL%'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests MODIFY COLUMN status ENUM(''SUBMITTED'',''PENDING_APPROVAL'',''APPROVED'',''REJECTED'',''WAITING_ACTUAL'',''ACTUAL_SUBMITTED'',''CANCELLED'') NOT NULL DEFAULT ''SUBMITTED'''
);
PREPARE stmt FROM @alter_overtime_status;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
