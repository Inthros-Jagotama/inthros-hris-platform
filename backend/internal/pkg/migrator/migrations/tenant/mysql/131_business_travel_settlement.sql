-- 131_business_travel_settlement.sql
-- Business Travel Module: settlement & settlement items (§25 plan doc)

-- =========================================================================
-- Settlements
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_settlements (
    id                        CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id        CHAR(36)      NOT NULL,
    participant_id             CHAR(36)      NULL,
    total_advance               DECIMAL(18,2) NOT NULL DEFAULT 0,
    total_actual_expense        DECIMAL(18,2) NOT NULL DEFAULT 0,
    total_company_paid          DECIMAL(18,2) NOT NULL DEFAULT 0,
    total_reimbursement         DECIMAL(18,2) NOT NULL DEFAULT 0,
    total_refund                DECIMAL(18,2) NOT NULL DEFAULT 0,
    balance                     DECIMAL(18,2) NOT NULL DEFAULT 0,
    status                      VARCHAR(30)   NOT NULL DEFAULT 'PENDING',
    submitted_at                 TIMESTAMP(6)  NULL DEFAULT NULL,
    approved_at                  TIMESTAMP(6)  NULL DEFAULT NULL,
    settled_at                   TIMESTAMP(6)  NULL DEFAULT NULL,
    notes                        VARCHAR(500)  NULL,
    deleted_at                   TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at                   TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                   TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_settle_travel (business_travel_id),
    INDEX idx_biztrav_settle_participant (participant_id),
    INDEX idx_biztrav_settle_status (status),
    INDEX idx_biztrav_settle_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Settlement Items (breakdown per expense/funding method)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_settlement_items (
    id                              CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_settlement_id    CHAR(36)      NOT NULL,
    expense_id                        CHAR(36)      NULL,
    funding_method_id                 CHAR(36)      NULL,
    item_type                         VARCHAR(30)   NOT NULL DEFAULT 'ACTUAL',
    category                          VARCHAR(100)  NULL,
    amount                             DECIMAL(18,2) NOT NULL DEFAULT 0,
    notes                              VARCHAR(500)  NULL,
    deleted_at                         TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at                         TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                         TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_settleitem_settlement (business_travel_settlement_id),
    INDEX idx_biztrav_settleitem_expense (expense_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
