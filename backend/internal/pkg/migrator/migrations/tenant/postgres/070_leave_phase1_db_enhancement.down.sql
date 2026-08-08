-- 070_leave_phase1_db_enhancement.down.sql

DROP TABLE IF EXISTS leave_balance_transactions;

DROP INDEX IF EXISTS idx_accrual_deleted_at;

ALTER TABLE leave_accrual_policies
    ALTER COLUMN deleted_at TYPE INT USING NULL;
