SET @drop_allow_checkin = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_company_settings'
      AND column_name = 'allow_checkin_on_day_off'
  ),
  'ALTER TABLE attendance_company_settings DROP COLUMN allow_checkin_on_day_off',
  'DO 0'
);
PREPARE stmt FROM @drop_allow_checkin;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
