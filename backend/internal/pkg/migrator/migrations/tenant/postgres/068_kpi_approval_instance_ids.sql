-- 068_kpi_approval_instance_ids.sql
-- Route KPI target/realization approval through the central approval
-- module: two separate approval instances per evaluation (target approval
-- and realization approval are independent checkpoints, potentially with
-- different flows/approvers), modules "performance_kpi_target" and
-- "performance_kpi_realization".

ALTER TABLE performance_evaluations
    ADD COLUMN IF NOT EXISTS target_approval_instance_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS realization_approval_instance_id CHAR(36) NULL;
