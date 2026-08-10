-- =============================================================================
-- Tenant Migration: 084_employeemovement_audit
-- =============================================================================
-- Employee Movement Enhancement (plan §12.6): Movement Audit Trail.
--
-- Mencatat seluruh perubahan lifecycle movement (CREATED, UPDATED, SUBMITTED,
-- APPROVED, REJECTED, CANCELLED, EXECUTED) lengkap dengan status lama/baru dan
-- snapshot JSON (old_data/new_data) sehingga transaksi HR dapat diaudit.
--
-- acted_by menyimpan id user (CHAR(36)) yang melakukan aksi (bisa NULL untuk
-- aksi sistem/kron job); acted_at default ke waktu insert.

CREATE TABLE IF NOT EXISTS employee_movement_audits (
    id          CHAR(36) PRIMARY KEY,
    movement_id CHAR(36) NOT NULL,
    action      VARCHAR(30) NOT NULL COMMENT 'CREATED, UPDATED, SUBMITTED, APPROVED, REJECTED, CANCELLED, EXECUTED',
    old_status  VARCHAR(20) NULL,
    new_status  VARCHAR(20) NULL,
    old_data    JSON NULL,
    new_data    JSON NULL,
    reason      TEXT NULL,
    acted_by    CHAR(36) NULL,
    acted_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_emp_mvmt_audit_movement (movement_id),
    INDEX idx_emp_mvmt_audit_action (action),

    CONSTRAINT fk_empmvmt_audit_movement FOREIGN KEY (movement_id)
        REFERENCES employee_movements(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
