-- 128_business_travel_funding.sql
-- Business Travel Module: funding method master, funding transactions & funding documents
-- Funding baru dapat dibuat setelah business_travels.status = APPROVED (lihat §14-17 plan doc).

-- =========================================================================
-- Funding Methods (master)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_funding_methods (
    id           CHAR(36)     NOT NULL PRIMARY KEY,
    code         VARCHAR(50)  NOT NULL,
    name         VARCHAR(100) NOT NULL,
    description  VARCHAR(300) NULL,
    active       SMALLINT     NOT NULL DEFAULT 1,
    deleted_at   TIMESTAMP    NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_biztrav_fundmethod_code ON business_travel_funding_methods (code);

-- =========================================================================
-- Fundings
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_fundings (
    id                   CHAR(36)      NOT NULL PRIMARY KEY,
    business_travel_id   CHAR(36)      NOT NULL,
    funding_method_id    CHAR(36)      NOT NULL,
    participant_id       CHAR(36)      NULL,
    amount               DECIMAL(18,2) NOT NULL DEFAULT 0,
    funding_date         DATE          NULL,
    payment_method       VARCHAR(50)   NULL,
    payment_reference    VARCHAR(100)  NULL,
    funded_by            CHAR(36)      NULL,
    status               VARCHAR(30)   NOT NULL DEFAULT 'PENDING',
    notes                VARCHAR(500)  NULL,
    deleted_at           TIMESTAMP     NULL,
    created_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_funding_travel ON business_travel_fundings (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_funding_method ON business_travel_fundings (funding_method_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_funding_status ON business_travel_fundings (status);

CREATE INDEX IF NOT EXISTS idx_biztrav_funding_deleted_at ON business_travel_fundings (deleted_at);

-- =========================================================================
-- Funding Documents (bukti transfer)
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_funding_documents (
    id                            CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_funding_id    CHAR(36)     NOT NULL,
    document_type                 VARCHAR(30)  NOT NULL DEFAULT 'OTHER',
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

CREATE INDEX IF NOT EXISTS idx_biztrav_funddoc_funding ON business_travel_funding_documents (business_travel_funding_id);
