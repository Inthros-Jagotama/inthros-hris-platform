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
    is_active   SMALLINT     NOT NULL DEFAULT 1,
    deleted_at  TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reimb_type_active ON reimbursement_types (is_active);

CREATE INDEX IF NOT EXISTS idx_reimb_type_deleted_at ON reimbursement_types (deleted_at);

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
    supervisor_action_at TIMESTAMP   NULL,
    supervisor_note     VARCHAR(500) NULL,
    hr_id               CHAR(36)     NULL,
    hr_action_at        TIMESTAMP    NULL,
    hr_note             VARCHAR(500) NULL,
    paid_at             TIMESTAMP    NULL,
    paid_amount         DECIMAL(18,2) NULL,
    submitted_at        TIMESTAMP    NULL,
    approved_at         TIMESTAMP    NULL,
    rejected_at         TIMESTAMP    NULL,
    cancelled_at        TIMESTAMP    NULL,
    deleted_at          TIMESTAMP    NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reimb_req_employee ON reimbursement_requests (employee_id);

CREATE INDEX IF NOT EXISTS idx_reimb_req_type ON reimbursement_requests (request_type_id);

CREATE INDEX IF NOT EXISTS idx_reimb_req_status ON reimbursement_requests (status);

CREATE INDEX IF NOT EXISTS idx_reimb_req_deleted_at ON reimbursement_requests (deleted_at);

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
    deleted_at                  TIMESTAMP    NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reimb_item_request ON reimbursement_items (reimbursement_request_id);

CREATE INDEX IF NOT EXISTS idx_reimb_item_deleted_at ON reimbursement_items (deleted_at);
