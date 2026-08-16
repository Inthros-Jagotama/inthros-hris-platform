-- 135_business_travel_attendance_rules.sql
-- Business Travel Module: konfigurasi attendance rule (§37 plan doc).
-- business_travel_id NULL berarti rule berlaku sebagai default company-wide,
-- diisi jika override khusus per perjalanan.

CREATE TABLE IF NOT EXISTS business_travel_attendance_rules (
    id                    CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id    CHAR(36)     NULL,
    rule_type             VARCHAR(30)  NOT NULL DEFAULT 'FULL_DAY',
    description           VARCHAR(300) NULL,
    is_default             TINYINT(1)   NOT NULL DEFAULT 0,
    active                 TINYINT(1)   NOT NULL DEFAULT 1,
    deleted_at              TIMESTAMP(6) NULL DEFAULT NULL,
    created_at              TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at              TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_attrule_travel (business_travel_id),
    INDEX idx_biztrav_attrule_active (active),
    INDEX idx_biztrav_attrule_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
