-- =============================================================================
-- Tenant Migration: 075_attendance_allow_checkin_day_off
-- =============================================================================
-- Attendance: tambah setting untuk mengizinkan/memblokir check-in pada hari
-- libur (tanpa jadwal shift / IsDayOff). Default = 1 (diizinkan) sehingga
-- check-in & check-out tetap berjalan; admin bisa mematikannya dari halaman
-- Settings Absensi.

SET @add_allow_checkin = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_company_settings'
      AND column_name = 'allow_checkin_on_day_off'
  ),
  'DO 0',
  'ALTER TABLE attendance_company_settings ADD COLUMN allow_checkin_on_day_off TINYINT(1) NOT NULL DEFAULT 1 COMMENT ''Izinkan check-in pada hari libur/jadwal tidak aktif'''
);
PREPARE stmt FROM @add_allow_checkin;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
