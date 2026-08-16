-- =============================================================================
-- Tenant Migration: 123_pph21_settings_calculation_method (MySQL DOWN)
-- =============================================================================
-- Kembalikan kolom calculation_method ke ENUM('REGULAR_GROSS_ANNUALIZED').
-- Catatan: hanya aman jika tidak ada baris yang menyimpan nilai 'TER'.

ALTER TABLE pph21_settings
    MODIFY COLUMN calculation_method ENUM('REGULAR_GROSS_ANNUALIZED') NOT NULL DEFAULT 'REGULAR_GROSS_ANNUALIZED';
