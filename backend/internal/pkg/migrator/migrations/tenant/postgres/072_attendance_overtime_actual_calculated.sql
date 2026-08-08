-- =============================================================================
-- Tenant Migration: 072_attendance_overtime_actual_calculated
-- =============================================================================
-- Attendance Module Phase 7 (docs/module-attendance-plan.md §31):
--   requested_minutes alone was being treated as the final overtime value.
--   Add actual_minutes (derived from the session's real checkout vs.
--   planned checkout) and calculated_minutes (actual, capped by the
--   approved/requested minutes) so the two are tracked separately.

ALTER TABLE attendance_overtime_requests ADD COLUMN IF NOT EXISTS actual_minutes INT NULL;
ALTER TABLE attendance_overtime_requests ADD COLUMN IF NOT EXISTS calculated_minutes INT NULL;
