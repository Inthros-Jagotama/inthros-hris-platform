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
    requires_receipt    TINYINT(1)   NOT NULL DEFAULT 1,
    reimbursable        TINYINT(1)   NOT NULL DEFAULT 1,
    payroll_treatment   VARCHAR(50)  NULL,
    account_code        VARCHAR(50)  NULL,
    active              TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at          TIMESTAMP(6) NULL DEFAULT NULL,
    created_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at          TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uq_biztrav_expcat_code (code),
    INDEX idx_biztrav_expcat_active (active),
    INDEX idx_biztrav_expcat_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    deleted_at             TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at             TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at             TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_expplan_travel (business_travel_id),
    INDEX idx_biztrav_expplan_category (expense_category_id),
    INDEX idx_biztrav_expplan_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
