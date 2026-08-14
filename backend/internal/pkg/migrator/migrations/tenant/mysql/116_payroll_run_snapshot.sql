-- =============================================================================
-- Tenant Migration: 116_payroll_run_snapshot
-- =============================================================================
-- Perluas payroll_run_items dengan field snapshot sesuai docs/payroll/
-- 03-payroll-run-snapshot.md §13.

ALTER TABLE payroll_run_items
    ADD COLUMN calculation_type VARCHAR(255) NOT NULL DEFAULT 'FIXED' AFTER component_type,
    ADD COLUMN base_amount       DECIMAL(18, 2) NOT NULL DEFAULT 0 AFTER amount,
    ADD COLUMN rate              DECIMAL(8, 4) NULL AFTER base_amount,
    ADD COLUMN formula           TEXT NULL AFTER rate,
    ADD COLUMN formula_result    DECIMAL(18, 2) NULL AFTER formula;
