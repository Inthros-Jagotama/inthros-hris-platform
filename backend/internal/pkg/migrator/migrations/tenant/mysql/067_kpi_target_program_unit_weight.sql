-- 067_kpi_target_program_unit_weight.sql
-- KPI Enhancement: employee fills unit_of_measurement on the indicator
-- target itself (template no longer authors it), and Program items gain
-- their own weight + unit_of_measurement (weight makes the Program
-- component score a proper weighted sum, mirroring indicators).

SET @add_detail_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_details'
      AND column_name = 'unit_of_measurement'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_details ADD COLUMN unit_of_measurement VARCHAR(50) NULL AFTER weight'
);
PREPARE stmt FROM @add_detail_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_program_weight = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_program_items'
      AND column_name = 'weight'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_program_items ADD COLUMN weight DECIMAL(5,2) NOT NULL DEFAULT 0.00 AFTER formula_type'
);
PREPARE stmt FROM @add_program_weight;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_program_unit = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'performance_evaluation_program_items'
      AND column_name = 'unit_of_measurement'
  ),
  'DO 0',
  'ALTER TABLE performance_evaluation_program_items ADD COLUMN unit_of_measurement VARCHAR(50) NULL AFTER weight'
);
PREPARE stmt FROM @add_program_unit;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
