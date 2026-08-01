-- Migration: 020_add_passport_to_employees
-- Menambahkan kolom passport ke tabel employees.
--
-- CATATAN (PostgreSQL): Migration 003_employee.sql sudah menyertakan kolom
-- `passport` di CREATE TABLE. Migration ini hanya dibutuhkan untuk tenant
-- lama, jadi gunakan ADD COLUMN IF NOT EXISTS (idempotent).

ALTER TABLE employees
    ADD COLUMN IF NOT EXISTS passport VARCHAR(50) NULL;
