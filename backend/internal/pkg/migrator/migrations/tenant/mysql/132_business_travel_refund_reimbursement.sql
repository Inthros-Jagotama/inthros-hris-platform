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
    deleted_at               TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at               TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at               TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_refund_travel (business_travel_id),
    INDEX idx_biztrav_refund_settlement (settlement_id),
    INDEX idx_biztrav_refund_status (status),
    INDEX idx_biztrav_refund_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
    requested_at              TIMESTAMP(6)  NULL DEFAULT NULL,
    approved_at               TIMESTAMP(6)  NULL DEFAULT NULL,
    paid_at                   TIMESTAMP(6)  NULL DEFAULT NULL,
    payment_reference         VARCHAR(100)  NULL,
    paid_by                   CHAR(36)      NULL,
    notes                     VARCHAR(500)  NULL,
    deleted_at                 TIMESTAMP(6)  NULL DEFAULT NULL,
    created_at                 TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at                 TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_reimb_travel (business_travel_id),
    INDEX idx_biztrav_reimb_settlement (settlement_id),
    INDEX idx_biztrav_reimb_status (status),
    INDEX idx_biztrav_reimb_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
