-- 068_kpi_approval_instance_ids.down.sql

ALTER TABLE performance_evaluations
    DROP COLUMN IF EXISTS realization_approval_instance_id,
    DROP COLUMN IF EXISTS target_approval_instance_id;
