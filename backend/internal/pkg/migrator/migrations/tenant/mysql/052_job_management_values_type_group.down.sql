-- Down Migration: 052_job_management_values_type_group
ALTER TABLE job_management_values
    DROP COLUMN description_group,
    DROP COLUMN type_group;
