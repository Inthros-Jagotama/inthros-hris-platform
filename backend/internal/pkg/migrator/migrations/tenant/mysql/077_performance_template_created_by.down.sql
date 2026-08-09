SET @drop_kpi_created_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'created_by'
  ),
  'ALTER TABLE performance_templates DROP COLUMN created_by',
  'DO 0'
);
PREPARE stmt FROM @drop_kpi_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_okr_created_by = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_templates'
      AND column_name = 'created_by'
  ),
  'ALTER TABLE okr_templates DROP COLUMN created_by',
  'DO 0'
);
PREPARE stmt FROM @drop_okr_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
