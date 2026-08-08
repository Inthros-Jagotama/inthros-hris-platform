-- 070_leave_phase1_db_enhancement.down.sql

DROP TABLE IF EXISTS leave_balance_transactions;

SET @drop_accrual_deleted_at_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'leave_accrual_policies'
      AND index_name = 'idx_accrual_deleted_at'
  ),
  'ALTER TABLE leave_accrual_policies DROP INDEX idx_accrual_deleted_at',
  'DO 0'
);
PREPARE stmt FROM @drop_accrual_deleted_at_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

ALTER TABLE leave_accrual_policies
    MODIFY COLUMN deleted_at INT NULL;
