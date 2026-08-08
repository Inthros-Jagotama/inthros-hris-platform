-- =============================================================================
-- Tenant Migration: 071_attendance_phase1_event_index
-- =============================================================================
-- Attendance Module Phase 1 (docs/module-attendance-plan.md §51/§55):
--   attendance_events only had a single-column index on employee_id.
--   Daily queries (list an employee's events for a work date/range) need a
--   composite index covering employee_id + event_time so they don't fall
--   back to scanning every event the employee ever had.

CREATE INDEX IF NOT EXISTS idx_att_event_employee_time ON attendance_events (employee_id, event_time_local);
