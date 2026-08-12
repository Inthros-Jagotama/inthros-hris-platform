-- =============================================================================
-- Tenant Migration: 095_recruitment_offer (MySQL)
-- =============================================================================
-- G-3 Offer Management (entity baru)
-- (docs/module-recruitment-development-plan.md §G-3)
--
-- Tabel job_offers — penawaran kerja ke kandidat setelah seleksi:
--
--   id                    CHAR(36) PK
--   application_id        CHAR(36) NN   → job_applications (kandidat + requisition)
--   offer_number          VARCHAR(50)    → nomor offer (auto-generated OFF-YYYYMM-XXXXXXXX)
--   employment_type       VARCHAR(50)    → FULL_TIME | PART_TIME | CONTRACT | INTERNSHIP ...
--   salary                DECIMAL(15,2)  → gaji pokok
--   allowances            DECIMAL(15,2)  → tunjangan
--   benefits              TEXT           → teks bebas benefit
--   start_date            VARCHAR(10)    → tanggal mulai (YYYY-MM-DD)
--   expiry_date           VARCHAR(10)    → tanggal kedaluwarsa offer
--   status                VARCHAR(30) NN DEFAULT 'DRAFT'
--     DRAFT → PENDING_APPROVAL → APPROVED → SENT → ACCEPTED/REJECTED/EXPIRED
--     WITHDRAWN (dari DRAFT/APPROVED/SENT)
--   sent_at               BIGINT NULL    → unix nano saat offer dikirim
--   accepted_at           BIGINT NULL    → unix nano saat kandidat menerima
--   rejected_at           BIGINT NULL    → unix nano saat kandidat menolak
--   approval_instance_id  CHAR(36) NULL  → instance Central Approval (G-1 pattern)
--   created_at / updated_at
--
-- Index: application_id (idx_offer_app), status (idx_offer_status).
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS job_offers (
    id                    CHAR(36) PRIMARY KEY,
    application_id        CHAR(36) NOT NULL,
    offer_number          VARCHAR(50) NULL,
    employment_type       VARCHAR(50) NULL,
    salary                DECIMAL(15,2) NOT NULL DEFAULT 0,
    allowances            DECIMAL(15,2) NOT NULL DEFAULT 0,
    benefits              TEXT NULL,
    start_date            VARCHAR(10) NULL,
    expiry_date           VARCHAR(10) NULL,
    status                VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    sent_at               BIGINT NULL,
    accepted_at           BIGINT NULL,
    rejected_at           BIGINT NULL,
    approval_instance_id  CHAR(36) NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_offers_application FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_offer_app ON job_offers (application_id);
CREATE INDEX idx_offer_status ON job_offers (status);
