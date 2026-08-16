-- 129_business_travel_expenses.sql
-- Business Travel Module: actual expenses & expense documents
-- funding_method_id di sini boleh berbeda dengan funding awal travel (lihat §21 plan doc).

-- =========================================================================
-- Expenses (actual)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_expenses (
    id                     CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id     CHAR(36)      NOT NULL,
    participant_id         CHAR(36)      NULL,
    expense_category_id    CHAR(36)      NOT NULL,
    expense_date           DATE          NOT NULL,
    description            VARCHAR(300)  NULL,
    quantity                DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit                    VARCHAR(30)   NULL,
    amount                  DECIMAL(18,2) NOT NULL DEFAULT 0,
    funding_method_id       CHAR(36)      NULL,
    vendor                  VARCHAR(150)  NULL,
    receipt_number          VARCHAR(100)  NULL,
    status                  VARCHAR(30)   NOT NULL DEFAULT 'DRAFT',
    notes                   VARCHAR(500)  NULL,
    deleted_at              TIMESTAMP     NULL,
    created_at              TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_exp_travel ON business_travel_expenses (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_exp_category ON business_travel_expenses (expense_category_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_exp_funding_method ON business_travel_expenses (funding_method_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_exp_deleted_at ON business_travel_expenses (deleted_at);

-- =========================================================================
-- Expense Documents (bukti pengeluaran)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_expense_documents (
    id                            CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_expense_id    CHAR(36)     NOT NULL,
    document_type                 VARCHAR(30)  NOT NULL DEFAULT 'RECEIPT',
    file_name                     VARCHAR(255) NOT NULL,
    file_path                     TEXT         NOT NULL,
    mime_type                     VARCHAR(100) NULL,
    file_size                     BIGINT       NULL,
    uploaded_by                   CHAR(36)     NULL,
    uploaded_at                   TIMESTAMP    NULL,
    deleted_at                    TIMESTAMP    NULL,
    created_at                    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_expdoc_expense ON business_travel_expense_documents (business_travel_expense_id);
