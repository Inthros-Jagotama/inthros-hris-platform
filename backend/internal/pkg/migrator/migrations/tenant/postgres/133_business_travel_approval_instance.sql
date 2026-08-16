-- 133_business_travel_approval_instance.sql
-- Approval Module integration: persist approval_instance_id for both the
-- Travel Approval flow and the Settlement Approval flow (dua alur terpisah,
-- lihat docs/module-attendance-business-travel-development-plan.md §54.3).

ALTER TABLE business_travels
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_biztrav_approval_instance ON business_travels (approval_instance_id);

ALTER TABLE business_travel_settlements
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_biztrav_settle_approval_instance ON business_travel_settlements (approval_instance_id);
