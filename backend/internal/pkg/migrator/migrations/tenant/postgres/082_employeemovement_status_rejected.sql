-- =============================================================================
-- Tenant Migration: 082_employeemovement_status_rejected
-- =============================================================================
-- Postgres no-op: kolom `status` di employee_movements sudah VARCHAR(255) sejak
-- 012_employee_movement.sql, sehingga nilai 'rejected' (8 karakter) sudah
-- tertampung tanpa perubahan skema.
--
-- Migration ini dibuat sebagai pasangan dari versi MySQL (yang memverifikasi
-- kolom bukan ENUM) agar kedua dialect tetap sinkron nomornya.

SELECT 1;
