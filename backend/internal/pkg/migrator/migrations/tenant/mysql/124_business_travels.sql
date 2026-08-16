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
    deleted_at       TIMESTAMP(6) NULL DEFAULT NULL,
    created_at       TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at       TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uq_biztrav_request_number (request_number),
    INDEX idx_biztrav_requester (requester_id),
    INDEX idx_biztrav_status (status),
    INDEX idx_biztrav_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
