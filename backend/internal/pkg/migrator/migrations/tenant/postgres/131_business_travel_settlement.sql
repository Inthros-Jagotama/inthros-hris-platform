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
    submitted_at                 TIMESTAMP     NULL,
    approved_at                  TIMESTAMP     NULL,
    settled_at                   TIMESTAMP     NULL,
    notes                        VARCHAR(500)  NULL,
    deleted_at                   TIMESTAMP     NULL,
    created_at                   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_settle_travel ON business_travel_settlements (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_settle_participant ON business_travel_settlements (participant_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_settle_status ON business_travel_settlements (status);

CREATE INDEX IF NOT EXISTS idx_biztrav_settle_deleted_at ON business_travel_settlements (deleted_at);

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
    deleted_at                         TIMESTAMP     NULL,
    created_at                         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                         TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_settleitem_settlement ON business_travel_settlement_items (business_travel_settlement_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_settleitem_expense ON business_travel_settlement_items (expense_id);
