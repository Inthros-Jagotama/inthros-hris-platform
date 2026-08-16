-- 133_business_travel_approval_instance.sql
-- Approval Module integration: persist approval_instance_id for both the
-- Travel Approval flow and the Settlement Approval flow (dua alur terpisah,
-- lihat docs/module-attendance-business-travel-development-plan.md §54.3).

SET @add_biztrav_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travels'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE business_travels ADD COLUMN approval_instance_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_biztrav_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_biztrav_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travels'
      AND INDEX_NAME = 'idx_biztrav_approval_instance'
  ),
  'DO 0',
  'CREATE INDEX idx_biztrav_approval_instance ON business_travels (approval_instance_id)'
);
PREPARE stmt FROM @add_biztrav_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_biztrav_settle_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travel_settlements'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE business_travel_settlements ADD COLUMN approval_instance_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_biztrav_settle_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_biztrav_settle_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travel_settlements'
      AND INDEX_NAME = 'idx_biztrav_settle_approval_instance'
  ),
  'DO 0',
  'CREATE INDEX idx_biztrav_settle_approval_instance ON business_travel_settlements (approval_instance_id)'
);
PREPARE stmt FROM @add_biztrav_settle_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
