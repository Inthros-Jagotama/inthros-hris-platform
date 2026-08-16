-- 126_business_travel_activities_schedules.sql
-- Business Travel Module: activities/agenda & schedule/transportation

-- =========================================================================
-- Activities
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_activities (
    id                  CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id  CHAR(36)     NOT NULL,
    activity_date       DATE         NOT NULL,
    start_time          TIME         NULL,
    end_time            TIME         NULL,
    title               VARCHAR(200) NOT NULL,
    description         TEXT         NULL,
    location             VARCHAR(200) NULL,
    organizer           VARCHAR(150) NULL,
    notes               VARCHAR(500) NULL,
    deleted_at          TIMESTAMP    NULL,
    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_act_travel ON business_travel_activities (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_act_deleted_at ON business_travel_activities (deleted_at);

-- =========================================================================
-- Schedules / Transportation
-- =========================================================================
CREATE TABLE IF NOT EXISTS business_travel_schedules (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)     NOT NULL,
    schedule_type         VARCHAR(20)  NOT NULL DEFAULT 'DEPARTURE',
    departure_datetime    TIMESTAMP    NULL,
    arrival_datetime      TIMESTAMP    NULL,
    origin                VARCHAR(200) NULL,
    destination           VARCHAR(200) NULL,
    transportation_type   VARCHAR(30)  NOT NULL DEFAULT 'OTHER',
    provider              VARCHAR(150) NULL,
    booking_reference     VARCHAR(100) NULL,
    notes                 VARCHAR(500) NULL,
    deleted_at            TIMESTAMP    NULL,
    created_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_sched_travel ON business_travel_schedules (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_sched_deleted_at ON business_travel_schedules (deleted_at);
