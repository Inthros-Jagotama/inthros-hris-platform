-- 067_kpi_target_program_unit_weight.down.sql

ALTER TABLE performance_evaluation_program_items
    DROP COLUMN IF EXISTS unit_of_measurement,
    DROP COLUMN IF EXISTS weight;

ALTER TABLE performance_evaluation_details
    DROP COLUMN IF EXISTS unit_of_measurement;
