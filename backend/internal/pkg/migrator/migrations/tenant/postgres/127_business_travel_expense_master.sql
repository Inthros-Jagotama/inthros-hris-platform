-- 127_business_travel_expense_master.sql
-- Business Travel Module: expense category master & expense plans (estimasi biaya)
-- Catatan: expense_plans TIDAK punya kolom funding_method — funding baru
-- ditentukan setelah travel approved (lihat §11 plan doc).

-- =========================================================================
-- Expense Categories (master)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_expense_categories (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    code                VARCHAR(50)  NOT NULL,
    name                VARCHAR(150) NOT NULL,
    description         VARCHAR(500) NULL,
    requires_receipt    SMALLINT     NOT NULL DEFAULT 1,
    reimbursable        SMALLINT     NOT NULL DEFAULT 1,
    payroll_treatment   VARCHAR(50)  NULL,
    account_code        VARCHAR(50)  NULL,
    active              SMALLINT     NOT NULL DEFAULT 1,
    deleted_at          TIMESTAMP    NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_biztrav_expcat_code ON business_travel_expense_categories (code);

CREATE INDEX IF NOT EXISTS idx_biztrav_expcat_active ON business_travel_expense_categories (active);

CREATE INDEX IF NOT EXISTS idx_biztrav_expcat_deleted_at ON business_travel_expense_categories (deleted_at);

-- =========================================================================
-- Expense Plans (estimasi biaya saat request)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_expense_plans (
    id                    CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)      NOT NULL,
    participant_id        CHAR(36)      NULL,
    expense_category_id   CHAR(36)      NOT NULL,
    description           VARCHAR(300)  NULL,
    quantity               DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit                   VARCHAR(30)   NULL,
    estimated_amount       DECIMAL(18,2) NOT NULL DEFAULT 0,
    notes                  VARCHAR(500)  NULL,
    deleted_at             TIMESTAMP     NULL,
    created_at             TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_expplan_travel ON business_travel_expense_plans (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_expplan_category ON business_travel_expense_plans (expense_category_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_expplan_deleted_at ON business_travel_expense_plans (deleted_at);
