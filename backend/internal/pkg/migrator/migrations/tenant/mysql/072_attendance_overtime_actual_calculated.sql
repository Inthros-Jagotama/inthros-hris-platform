-- =============================================================================
-- Tenant Migration: 072_attendance_overtime_actual_calculated
-- =============================================================================
-- Attendance Module Phase 7 (docs/module-attendance-plan.md §31):
--   requested_minutes alone was being treated as the final overtime value.
--   Add actual_minutes (derived from the session's real checkout vs.
--   planned checkout) and calculated_minutes (actual, capped by the
--   approved/requested minutes) so the two are tracked separately.

SET @add_actual_minutes = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'actual_minutes'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN actual_minutes INT NULL AFTER requested_minutes'
);
PREPARE stmt FROM @add_actual_minutes;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_calculated_minutes = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_overtime_requests'
      AND column_name = 'calculated_minutes'
  ),
  'DO 0',
  'ALTER TABLE attendance_overtime_requests ADD COLUMN calculated_minutes INT NULL AFTER actual_minutes'
);
PREPARE stmt FROM @add_calculated_minutes;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
