-- Down Migration: 052_job_management_values_type_group
ALTER TABLE job_management_values
    DROP COLUMN IF EXISTS description_group,
    DROP COLUMN IF EXISTS type_group;
