-- 125_business_travel_participants_destinations.sql
-- Business Travel Module: participants & destinations

-- =========================================================================
-- Participants
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_participants (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id  CHAR(36)     NOT NULL,
    participant_type    VARCHAR(20)  NOT NULL DEFAULT 'EMPLOYEE',
    employee_id         CHAR(36)     NULL,
    name                VARCHAR(150) NULL,
    organization        VARCHAR(150) NULL,
    position            VARCHAR(150) NULL,
    identity_number     VARCHAR(50)  NULL,
    email               VARCHAR(150) NULL,
    phone               VARCHAR(30)  NULL,
    role                VARCHAR(30)  NOT NULL DEFAULT 'MEMBER',
    notes               VARCHAR(500) NULL,
    deleted_at          TIMESTAMP    NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_part_travel ON business_travel_participants (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_part_employee ON business_travel_participants (employee_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_part_deleted_at ON business_travel_participants (deleted_at);

-- =========================================================================
-- Destinations
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_destinations (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id  CHAR(36)     NOT NULL,
    sequence            INT          NOT NULL DEFAULT 1,
    country             VARCHAR(100) NULL,
    province            VARCHAR(100) NULL,
    city                VARCHAR(100) NULL,
    location             VARCHAR(200) NULL,
    arrival_date        DATE         NULL,
    departure_date      DATE         NULL,
    purpose             VARCHAR(300) NULL,
    notes               VARCHAR(500) NULL,
    deleted_at          TIMESTAMP    NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_dest_travel ON business_travel_destinations (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_dest_deleted_at ON business_travel_destinations (deleted_at);
