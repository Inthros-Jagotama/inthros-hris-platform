-- Down Migration: 032_rename_job_value_hazard_to_risk
-- Mengembalikan nilai type 'risk' menjadi 'hazard' (rollback).

UPDATE job_management_values
SET type = 'hazard'
WHERE type = 'risk';
