-- =============================================================================
-- Tenant Migration: 117_payroll_proration_method
-- =============================================================================
-- Metode prorasi gaji configurable per payroll run (docs/payroll/06 §19).
-- Nilai: CALENDAR_DAYS (default) | WORKING_DAYS | FIXED_30_DAYS | ATTENDANCE_DAYS

ALTER TABLE payroll_runs
    ADD COLUMN proration_method VARCHAR(255) NOT NULL DEFAULT 'CALENDAR_DAYS';
