-- =============================================================================
-- Tenant Migration: 070_leave_phase1_db_enhancement
-- =============================================================================
-- Leave Module Phase 1 (docs/module-leave-plan.md):
--   1. Fix leave_accrual_policies.deleted_at — was mistakenly typed INT
--      instead of TIMESTAMP (mismatched with the Go model's gorm.DeletedAt),
--      and had no index unlike every other soft-deletable table in this
--      module.
--   2. Add leave_balance_transactions — the balance ledger table, so future
--      accrual/usage/adjustment/reversal/carry-forward/expiry writes have an
--      auditable history instead of only ever overwriting the current
--      balance row in employee_leave_balances.

-- ---------------------------------------------------------------------------
-- 1. Fix leave_accrual_policies.deleted_at type + add missing index
-- ---------------------------------------------------------------------------
ALTER TABLE leave_accrual_policies
    MODIFY COLUMN deleted_at TIMESTAMP NULL;

SET @add_accrual_deleted_at_index = IF(
  NOT EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'leave_accrual_policies'
      AND index_name = 'idx_accrual_deleted_at'
  ),
  'ALTER TABLE leave_accrual_policies ADD INDEX idx_accrual_deleted_at (deleted_at)',
  'DO 0'
);
PREPARE stmt FROM @add_accrual_deleted_at_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- 2. Leave Balance Transactions (Ledger)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_balance_transactions (
    id                CHAR(36) PRIMARY KEY,
    employee_id       CHAR(36) NOT NULL,
    leave_type_id     CHAR(36) NOT NULL,
    balance_id        CHAR(36) NOT NULL,
    transaction_type  ENUM('ACCRUAL', 'USAGE', 'ADJUSTMENT', 'REVERSAL', 'CARRY_FORWARD', 'EXPIRY') NOT NULL,
    reference_type    VARCHAR(50) NULL,
    reference_id      CHAR(36) NULL,
    amount            DECIMAL(6, 2) NOT NULL,
    balance_before    DECIMAL(6, 2) NOT NULL,
    balance_after     DECIMAL(6, 2) NOT NULL,
    note              VARCHAR(255) NULL,
    created_by        CHAR(36) NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_leave_bal_txn_employee_type (employee_id, leave_type_id, created_at),
    INDEX idx_leave_bal_txn_balance (balance_id),
    INDEX idx_leave_bal_txn_reference (reference_type, reference_id),
    INDEX idx_leave_bal_txn_type (transaction_type),

    CONSTRAINT fk_leave_bal_txn_employee FOREIGN KEY (employee_id)   REFERENCES employees(id)               ON DELETE CASCADE,
    CONSTRAINT fk_leave_bal_txn_type     FOREIGN KEY (leave_type_id) REFERENCES leave_types(id)              ON DELETE CASCADE,
    CONSTRAINT fk_leave_bal_txn_balance  FOREIGN KEY (balance_id)    REFERENCES employee_leave_balances(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
