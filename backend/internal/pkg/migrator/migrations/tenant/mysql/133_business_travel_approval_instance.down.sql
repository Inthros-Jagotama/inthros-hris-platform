-- 133_business_travel_approval_instance.down.sql

SET @drop_biztrav_settle_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travel_settlements'
      AND INDEX_NAME = 'idx_biztrav_settle_approval_instance'
  ),
  'DROP INDEX idx_biztrav_settle_approval_instance ON business_travel_settlements',
  'DO 0'
);
PREPARE stmt FROM @drop_biztrav_settle_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_biztrav_settle_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travel_settlements'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE business_travel_settlements DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_biztrav_settle_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_biztrav_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travels'
      AND INDEX_NAME = 'idx_biztrav_approval_instance'
  ),
  'DROP INDEX idx_biztrav_approval_instance ON business_travels',
  'DO 0'
);
PREPARE stmt FROM @drop_biztrav_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_biztrav_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'business_travels'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE business_travels DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_biztrav_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
