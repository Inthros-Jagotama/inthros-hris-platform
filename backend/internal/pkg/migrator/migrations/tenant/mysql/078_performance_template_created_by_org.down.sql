SET @drop_kpi_created_by_org = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_templates'
      AND column_name = 'created_by_org_id'
  ),
  'ALTER TABLE performance_templates DROP COLUMN created_by_org_id',
  'DO 0'
);
PREPARE stmt FROM @drop_kpi_created_by_org;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_okr_created_by_org = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'okr_templates'
      AND column_name = 'created_by_org_id'
  ),
  'ALTER TABLE okr_templates DROP COLUMN created_by_org_id',
  'DO 0'
);
PREPARE stmt FROM @drop_okr_created_by_org;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
