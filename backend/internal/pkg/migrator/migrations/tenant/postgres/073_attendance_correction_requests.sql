-- =============================================================================
-- Tenant Migration: 073_attendance_correction_requests
-- =============================================================================
-- Attendance Module Phase 8 (docs/module-attendance-plan.md §16/§33-34):
--   No correction workflow existed at all - a wrong/missing check-in or
--   check-out had no way to be fixed without mutating attendance_events,
--   which breaks §15's "raw event is immutable" principle. This table lets
--   employees/HR request a correction that, once approved through the
--   Central Approval Module, is applied to the session (not the raw event)
--   and triggers a recalculation.

CREATE TABLE IF NOT EXISTS attendance_correction_requests (
    id                       CHAR(36) PRIMARY KEY,
    employee_id              CHAR(36) NOT NULL,
    attendance_session_id    CHAR(36) NOT NULL,
    correction_type          VARCHAR(50) NOT NULL,
    requested_checkin        TIMESTAMP(6) NULL,
    requested_checkout       TIMESTAMP(6) NULL,
    reason                   VARCHAR(255) NOT NULL,
    status                   VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED',
    approval_instance_id     CHAR(36) NULL,
    created_by               CHAR(36) NULL,
    approved_at              TIMESTAMP(6) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_att_correction_employee FOREIGN KEY (employee_id)           REFERENCES employees(id)            ON DELETE CASCADE,
    CONSTRAINT fk_att_correction_session  FOREIGN KEY (attendance_session_id) REFERENCES attendance_sessions(id)  ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_att_correction_employee ON attendance_correction_requests (employee_id);
CREATE INDEX IF NOT EXISTS idx_att_correction_session ON attendance_correction_requests (attendance_session_id);
CREATE INDEX IF NOT EXISTS idx_att_correction_status ON attendance_correction_requests (status);
CREATE INDEX IF NOT EXISTS idx_att_correction_approval_instance ON attendance_correction_requests (approval_instance_id);
