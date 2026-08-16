-- 132_business_travel_refund_reimbursement.sql
-- Business Travel Module: refund (§35) & reimbursement (§36) hasil settlement
-- Catatan: business_travel_reimbursements di sini bersifat internal ke module ini.
-- Integrasi ke module Reimbursement standalone hanya dilakukan JIKA tenant
-- subscribe module tersebut (lihat docs/module-attendance-business-travel-development-plan.md §54.7).

-- =========================================================================
-- Refunds
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_refunds (
    id                    CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)      NOT NULL,
    settlement_id          CHAR(36)      NULL,
    participant_id         CHAR(36)      NULL,
    refund_amount           DECIMAL(18,2) NOT NULL DEFAULT 0,
    refund_date             DATE          NULL,
    refund_reference        VARCHAR(100)  NULL,
    refunded_by             CHAR(36)      NULL,
    refund_document         TEXT          NULL,
    status                  VARCHAR(30)   NOT NULL DEFAULT 'PENDING',
    notes                   VARCHAR(500)  NULL,
    deleted_at               TIMESTAMP     NULL,
    created_at               TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_refund_travel ON business_travel_refunds (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_refund_settlement ON business_travel_refunds (settlement_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_refund_status ON business_travel_refunds (status);

CREATE INDEX IF NOT EXISTS idx_biztrav_refund_deleted_at ON business_travel_refunds (deleted_at);

-- =========================================================================
-- Reimbursements (additional claim hasil settlement)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_reimbursements (
    id                    CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)      NOT NULL,
    participant_id         CHAR(36)      NULL,
    settlement_id           CHAR(36)      NULL,
    amount                   DECIMAL(18,2) NOT NULL DEFAULT 0,
    status                   VARCHAR(30)   NOT NULL DEFAULT 'REQUESTED',
    requested_at              TIMESTAMP     NULL,
    approved_at               TIMESTAMP     NULL,
    paid_at                   TIMESTAMP     NULL,
    payment_reference         VARCHAR(100)  NULL,
    paid_by                   CHAR(36)      NULL,
    notes                     VARCHAR(500)  NULL,
    deleted_at                 TIMESTAMP     NULL,
    created_at                 TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_reimb_travel ON business_travel_reimbursements (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_reimb_settlement ON business_travel_reimbursements (settlement_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_reimb_status ON business_travel_reimbursements (status);

CREATE INDEX IF NOT EXISTS idx_biztrav_reimb_deleted_at ON business_travel_reimbursements (deleted_at);
