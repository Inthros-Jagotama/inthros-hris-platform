-- =============================================================================
-- Tenant Migration: 084_employeemovement_audit
-- =============================================================================
-- Employee Movement Enhancement (plan §12.6): Movement Audit Trail.
--
-- Mencatat seluruh perubahan lifecycle movement (CREATED, UPDATED, SUBMITTED,
-- APPROVED, REJECTED, CANCELLED, EXECUTED) lengkap dengan status lama/baru dan
-- snapshot JSON (old_data/new_data) sehingga transaksi HR dapat diaudit.

CREATE TABLE IF NOT EXISTS employee_movement_audits (
    id          CHAR(36) PRIMARY KEY,
    movement_id CHAR(36) NOT NULL,
    action      VARCHAR(30) NOT NULL,
    old_status  VARCHAR(20) NULL,
    new_status  VARCHAR(20) NULL,
    old_data    JSON NULL,
    new_data    JSON NULL,
    reason      TEXT NULL,
    acted_by    CHAR(36) NULL,
    acted_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_empmvmt_audit_movement FOREIGN KEY (movement_id)
        REFERENCES employee_movements(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_emp_mvmt_audit_movement ON employee_movement_audits (movement_id);
CREATE INDEX IF NOT EXISTS idx_emp_mvmt_audit_action ON employee_movement_audits (action);
