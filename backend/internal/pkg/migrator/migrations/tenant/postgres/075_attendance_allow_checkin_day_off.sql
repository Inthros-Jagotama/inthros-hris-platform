-- =============================================================================
-- Tenant Migration: 075_attendance_allow_checkin_day_off
-- =============================================================================
-- Attendance: tambah setting untuk mengizinkan/memblokir check-in pada hari
-- libur (tanpa jadwal shift / IsDayOff). Default = TRUE (diizinkan) sehingga
-- check-in & check-out tetap berjalan; admin bisa mematikannya dari halaman
-- Settings Absensi.

ALTER TABLE attendance_company_settings ADD COLUMN IF NOT EXISTS allow_checkin_on_day_off BOOLEAN NOT NULL DEFAULT TRUE;
