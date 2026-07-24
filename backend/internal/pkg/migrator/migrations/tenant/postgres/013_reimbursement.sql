-- 013_reimbursement.sql
-- Reimbursement & Claim Module
-- Tabel untuk mengelola pengajuan reimbursement/klaim karyawan

-- =========================================================================
-- Reimbursement Types
-- =========================================================================
CREATE TABLE IF NOT EXISTS reimbursement_types (
    id          CHAR(36)     NOT NULL PRIMARY KEY,
    code        VARCHAR(50)  NOT NULL DEFAULT '',
    name        VARCHAR(150) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    is_active   TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_reimb_type_active (is_active),
    INDEX idx_reimb_type_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Reimbursement Requests
-- =========================================================================
CREATE TABLE IF NOT EXISTS reimbursement_requests (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    employee_id         CHAR(36)     NOT NULL,
    request_type_id     CHAR(36)     NOT NULL,
    title               VARCHAR(200) NOT NULL,
    description         TEXT         NULL,
    total_amount        DECIMAL(18,2) NOT NULL DEFAULT 0.00,
    currency            VARCHAR(3)   NOT NULL DEFAULT 'IDR',
    status              VARCHAR(50)  NOT NULL DEFAULT 'DRAFT',
    supervisor_id       CHAR(36)     NULL,
    supervisor_action_at TIMESTAMP(6) NULL DEFAULT NULL,
    supervisor_note     VARCHAR(500) NULL,
    hr_id               CHAR(36)     NULL,
    hr_action_at        TIMESTAMP(6) NULL DEFAULT NULL,
    hr_note             VARCHAR(500) NULL,
    paid_at             TIMESTAMP(6) NULL DEFAULT NULL,
    paid_amount         DECIMAL(18,2) NULL,
    submitted_at        TIMESTAMP(6) NULL DEFAULT NULL,
    approved_at         TIMESTAMP(6) NULL DEFAULT NULL,
    rejected_at         TIMESTAMP(6) NULL DEFAULT NULL,
    cancelled_at        TIMESTAMP(6) NULL DEFAULT NULL,
    deleted_at          TIMESTAMP(6) NULL DEFAULT NULL,
    created_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_reimb_req_employee (employee_id),
    INDEX idx_reimb_req_type (request_type_id),
    INDEX idx_reimb_req_status (status),
    INDEX idx_reimb_req_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Reimbursement Items (detail pengeluaran per pengajuan)
-- =========================================================================
CREATE TABLE IF NOT EXISTS reimbursement_items (
    id                          CHAR(36)     NOT NULL PRIMARY KEY,
    reimbursement_request_id    CHAR(36)     NOT NULL,
    expense_date                DATE         NOT NULL,
    expense_type                VARCHAR(100) NOT NULL,
    description                 VARCHAR(500) NULL,
    amount                      DECIMAL(18,2) NOT NULL,
    receipt_url                 TEXT         NULL,
    deleted_at                  TIMESTAMP(6) NULL DEFAULT NULL,
    created_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                  TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_reimb_item_request (reimbursement_request_id),
    INDEX idx_reimb_item_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
