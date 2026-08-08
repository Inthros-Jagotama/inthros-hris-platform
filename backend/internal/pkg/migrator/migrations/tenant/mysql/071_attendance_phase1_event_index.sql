-- =============================================================================
-- Tenant Migration: 071_attendance_phase1_event_index
-- =============================================================================
-- Attendance Module Phase 1 (docs/module-attendance-plan.md §51/§55):
--   attendance_events only had a single-column index on employee_id.
--   Daily queries (list an employee's events for a work date/range) need a
--   composite index covering employee_id + event_time so they don't fall
--   back to scanning every event the employee ever had.

SET @add_att_event_employee_time_index = IF(
  NOT EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'attendance_events'
      AND index_name = 'idx_att_event_employee_time'
  ),
  'ALTER TABLE attendance_events ADD INDEX idx_att_event_employee_time (employee_id, event_time_local)',
  'DO 0'
);
PREPARE stmt FROM @add_att_event_employee_time_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
