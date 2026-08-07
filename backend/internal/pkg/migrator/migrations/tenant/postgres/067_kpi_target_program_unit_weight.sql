-- 067_kpi_target_program_unit_weight.sql
-- KPI Enhancement: employee fills unit_of_measurement on the indicator
-- target itself (template no longer authors it), and Program items gain
-- their own weight + unit_of_measurement (weight makes the Program
-- component score a proper weighted sum, mirroring indicators).

ALTER TABLE performance_evaluation_details
    ADD COLUMN IF NOT EXISTS unit_of_measurement VARCHAR(50) NULL;

ALTER TABLE performance_evaluation_program_items
    ADD COLUMN IF NOT EXISTS weight DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS unit_of_measurement VARCHAR(50) NULL;
