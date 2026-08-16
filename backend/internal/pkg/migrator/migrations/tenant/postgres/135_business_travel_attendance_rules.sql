-- 135_business_travel_attendance_rules.sql
-- Business Travel Module: konfigurasi attendance rule (§37 plan doc).
-- business_travel_id NULL berarti rule berlaku sebagai default company-wide,
-- diisi jika override khusus per perjalanan.

CREATE TABLE IF NOT EXISTS business_travel_attendance_rules (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)     NULL,
    rule_type             VARCHAR(30)  NOT NULL DEFAULT 'FULL_DAY',
    description           VARCHAR(300) NULL,
    is_default             SMALLINT     NOT NULL DEFAULT 0,
    active                 SMALLINT     NOT NULL DEFAULT 1,
    deleted_at              TIMESTAMP    NULL,
    created_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_attrule_travel ON business_travel_attendance_rules (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_attrule_active ON business_travel_attendance_rules (active);

CREATE INDEX IF NOT EXISTS idx_biztrav_attrule_deleted_at ON business_travel_attendance_rules (deleted_at);
