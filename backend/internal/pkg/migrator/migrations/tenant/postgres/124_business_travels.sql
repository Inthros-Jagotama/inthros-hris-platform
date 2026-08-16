-- 124_business_travels.sql
-- Business Travel Module: core travel request table
-- Request tidak menentukan funding method (lihat docs/module-attendance-business-travel-development-plan.md §4)

CREATE TABLE IF NOT EXISTS business_travels (
    id               CHAR(36)     NOT NULL PRIMARY KEY,
    request_number   VARCHAR(50)  NOT NULL,
    requester_id     CHAR(36)     NOT NULL,
    title            VARCHAR(200) NOT NULL,
    purpose          VARCHAR(500) NULL,
    description      TEXT         NULL,
    start_date       DATE         NOT NULL,
    end_date         DATE         NOT NULL,
    origin           VARCHAR(200) NULL,
    status           VARCHAR(30)  NOT NULL DEFAULT 'DRAFT',
    approval_status  VARCHAR(30)  NOT NULL DEFAULT 'DRAFT',
    created_by       CHAR(36)     NULL,
    deleted_at       TIMESTAMP    NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_biztrav_request_number ON business_travels (request_number);

CREATE INDEX IF NOT EXISTS idx_biztrav_requester ON business_travels (requester_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_status ON business_travels (status);

CREATE INDEX IF NOT EXISTS idx_biztrav_deleted_at ON business_travels (deleted_at);
