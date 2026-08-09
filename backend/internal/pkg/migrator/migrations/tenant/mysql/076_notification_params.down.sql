SET @drop_params = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'notifications'
      AND column_name = 'params'
  ),
  'ALTER TABLE notifications DROP COLUMN params',
  'DO 0'
);
PREPARE stmt FROM @drop_params;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
