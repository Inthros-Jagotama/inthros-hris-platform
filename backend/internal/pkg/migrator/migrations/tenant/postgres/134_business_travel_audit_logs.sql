-- 134_business_travel_audit_logs.sql
-- Business Travel Module: audit trail (§48 plan doc)

CREATE TABLE IF NOT EXISTS business_travel_audit_logs (
    id            CHAR(36)     NOT NULL PRIMARY KEY,
    entity_type   VARCHAR(50)  NOT NULL,
    entity_id     CHAR(36)     NOT NULL,
    action        VARCHAR(50)  NOT NULL,
    old_value     TEXT         NULL,
    new_value     TEXT         NULL,
    user_id       CHAR(36)     NULL,
    ip_address    VARCHAR(45)  NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_audit_entity ON business_travel_audit_logs (entity_type, entity_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_audit_action ON business_travel_audit_logs (action);

CREATE INDEX IF NOT EXISTS idx_biztrav_audit_created_at ON business_travel_audit_logs (created_at);
