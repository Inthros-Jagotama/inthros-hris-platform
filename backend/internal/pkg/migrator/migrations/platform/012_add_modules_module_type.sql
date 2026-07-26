-- Migration: 012_add_modules_module_type
-- Database: Platform (Cross-Dialect)
-- Menambahkan kolom module_type untuk membedakan modul platform dan tenant

ALTER TABLE modules
    ADD COLUMN module_type VARCHAR(20) NOT NULL DEFAULT 'tenant';
