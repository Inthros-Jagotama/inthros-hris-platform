-- =============================================================================
-- Tenant Migration: 087_employeemovement_cancellation
-- =============================================================================
-- Employee Movement Enhancement (plan §12.16): Movement Cancellation Approval.
--
-- Movement yang sudah `approved` TIDAK boleh dibatalkan langsung oleh HR —
-- pembatalan harus melalui Cancellation Request yang diproses Central Approval
-- Module (module slug `employeemovement_cancellation`). Selama permintaan
-- pembatalan berjalan, movement berada di status `cancellation_pending`.
--
-- Kolom baru:
--   cancellation_approval_instance_id CHAR(36) NULL
--     → approval instance yang dibuat untuk cancellation request (terpisah dari
--       approval_instance_id milik submission), sehingga status callback Central
--       Approval dapat membedakan hasil approval submission vs pembatalan.
--
-- Index baru:
--   idx_emp_mvmt_cancellation_instance (cancellation_approval_instance_id)

SET @add_cancel_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND COLUMN_NAME = 'cancellation_approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE employee_movements ADD COLUMN cancellation_approval_instance_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_cancel_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_mvmt_cancel_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'employee_movements'
      AND INDEX_NAME = 'idx_emp_mvmt_cancellation_instance'
  ),
  'DO 0',
  'CREATE INDEX idx_emp_mvmt_cancellation_instance ON employee_movements (cancellation_approval_instance_id)'
);
PREPARE stmt FROM @add_mvmt_cancel_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
