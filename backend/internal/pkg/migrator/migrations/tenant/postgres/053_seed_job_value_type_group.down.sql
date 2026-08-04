-- Down Migration: 053_seed_job_value_type_group
-- Mengosongkan kembali type_group & description_group yang diisi migration 053.

UPDATE job_management_values SET
    type_group = NULL,
    description_group = NULL
WHERE type_group IS NOT NULL OR description_group IS NOT NULL;
