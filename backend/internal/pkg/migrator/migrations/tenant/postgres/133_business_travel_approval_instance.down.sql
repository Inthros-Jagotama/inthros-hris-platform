-- 133_business_travel_approval_instance.down.sql

DROP INDEX IF EXISTS idx_biztrav_settle_approval_instance;

ALTER TABLE business_travel_settlements
    DROP COLUMN IF EXISTS approval_instance_id;

DROP INDEX IF EXISTS idx_biztrav_approval_instance;

ALTER TABLE business_travels
    DROP COLUMN IF EXISTS approval_instance_id;
