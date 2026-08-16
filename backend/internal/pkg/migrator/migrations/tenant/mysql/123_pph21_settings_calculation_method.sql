-- =============================================================================
-- Tenant Migration: 123_pph21_settings_calculation_method (MySQL)
-- =============================================================================
-- Ubah kolom calculation_method di pph21_settings dari
-- ENUM('REGULAR_GROSS_ANNUALIZED') menjadi VARCHAR(255).
--
-- Latar belakang: migrasi 120 (pph21_ter) memperkenalkan metode TER pada engine
-- dan DTO (oneof=TER REGULAR_GROSS_ANNUALIZED), namun hanya mengubah tabel
-- pph21_calculation_logs. Kolom calculation_method di pph21_settings tetap ENUM
-- dengan satu nilai sehingga menyimpan/update 'TER' gagal dengan
-- MySQL Error 1265 "Data truncated for column 'calculation_method'".
-- Postgres sejak awal memakai VARCHAR(255), jadi perubahan ini menyamakan MySQL.

ALTER TABLE pph21_settings
    MODIFY COLUMN calculation_method VARCHAR(255) NOT NULL DEFAULT 'REGULAR_GROSS_ANNUALIZED';
