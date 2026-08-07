-- 067_kpi_target_program_unit_weight.down.sql

SET @drop_program_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_program_items'
      AND column_name = 'unit_of_measurement'
  ),
  'ALTER TABLE performance_evaluation_program_items DROP COLUMN unit_of_measurement',
  'DO 0'
);
PREPARE stmt FROM @drop_program_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_program_weight = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_program_items'
      AND column_name = 'weight'
  ),
  'ALTER TABLE performance_evaluation_program_items DROP COLUMN weight',
  'DO 0'
);
PREPARE stmt FROM @drop_program_weight;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_detail_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'unit_of_measurement'
  ),
  'ALTER TABLE performance_evaluation_details DROP COLUMN unit_of_measurement',
  'DO 0'
);
PREPARE stmt FROM @drop_detail_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
