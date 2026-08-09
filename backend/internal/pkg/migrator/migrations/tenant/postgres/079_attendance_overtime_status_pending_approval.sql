-- =============================================================================
-- Tenant Migration: 079_attendance_overtime_status_pending_approval
-- =============================================================================
-- Postgres no-op: kolom `status` di attendance_overtime_requests sudah
-- VARCHAR(255) sejak awal (004_attendance.sql), sehingga nilai 'PENDING_APPROVAL'
-- (16 karakter) sudah tertampung tanpa perubahan skema.
--
-- Migration ini dibuat sebagai pasangan dari versi MySQL (yang mengubah
-- ENUM('SUBMITTED','APPROVED','REJECTED') menjadi menyertakan
-- 'PENDING_APPROVAL') agar kedua dialect tetap sinkron nomornya.

SELECT 1;
