-- Migration: 020_add_passport_to_employees
-- Menambahkan kolom passport ke tabel employees

ALTER TABLE employees
    ADD COLUMN passport VARCHAR(50) NULL;
