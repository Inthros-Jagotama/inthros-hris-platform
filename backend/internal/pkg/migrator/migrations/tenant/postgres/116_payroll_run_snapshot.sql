-- =============================================================================
-- Tenant Migration: 116_payroll_run_snapshot
-- =============================================================================
-- Perluas payroll_run_items dengan field snapshot sesuai docs/payroll/
-- 03-payroll-run-snapshot.md §13: payroll item menyimpan calculation_type,
-- base_amount, rate, formula, dan formula_result agar histori kalkulasi bisa
-- dijelaskan kembali (audit) walau konfigurasi live berubah.
-- calculated_amount = kolom amount yang sudah ada; employee/employer amount
-- dibedakan via kolom paid_by yang sudah ada.

ALTER TABLE payroll_run_items
    ADD COLUMN calculation_type VARCHAR(255) NOT NULL DEFAULT 'FIXED',
    ADD COLUMN base_amount       DECIMAL(18, 2) NOT NULL DEFAULT 0,
    ADD COLUMN rate              DECIMAL(8, 4) NULL,
    ADD COLUMN formula           TEXT NULL,
    ADD COLUMN formula_result    DECIMAL(18, 2) NULL;
