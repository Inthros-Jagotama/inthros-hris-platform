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
    ALTER COLUMN deleted_at TYPE TIMESTAMP USING NULL;

CREATE INDEX IF NOT EXISTS idx_accrual_deleted_at ON leave_accrual_policies (deleted_at);

-- ---------------------------------------------------------------------------
-- 2. Leave Balance Transactions (Ledger)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_balance_transactions (
    id                CHAR(36) PRIMARY KEY,
    employee_id       CHAR(36) NOT NULL,
    leave_type_id     CHAR(36) NOT NULL,
    balance_id        CHAR(36) NOT NULL,
    transaction_type  VARCHAR(20) NOT NULL,
    reference_type    VARCHAR(50) NULL,
    reference_id      CHAR(36) NULL,
    amount            DECIMAL(6, 2) NOT NULL,
    balance_before    DECIMAL(6, 2) NOT NULL,
    balance_after     DECIMAL(6, 2) NOT NULL,
    note              VARCHAR(255) NULL,
    created_by        CHAR(36) NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_leave_bal_txn_employee FOREIGN KEY (employee_id)   REFERENCES employees(id)               ON DELETE CASCADE,
    CONSTRAINT fk_leave_bal_txn_type     FOREIGN KEY (leave_type_id) REFERENCES leave_types(id)              ON DELETE CASCADE,
    CONSTRAINT fk_leave_bal_txn_balance  FOREIGN KEY (balance_id)    REFERENCES employee_leave_balances(id)  ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_leave_bal_txn_employee_type ON leave_balance_transactions (employee_id, leave_type_id, created_at);
CREATE INDEX IF NOT EXISTS idx_leave_bal_txn_balance ON leave_balance_transactions (balance_id);
CREATE INDEX IF NOT EXISTS idx_leave_bal_txn_reference ON leave_balance_transactions (reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_leave_bal_txn_type ON leave_balance_transactions (transaction_type);
