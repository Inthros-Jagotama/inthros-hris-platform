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
    correction_type          ENUM('MISSING_CHECKIN', 'MISSING_CHECKOUT', 'WRONG_CHECKIN', 'WRONG_CHECKOUT') NOT NULL,
    requested_checkin        DATETIME(6) NULL,
    requested_checkout       DATETIME(6) NULL,
    reason                   VARCHAR(255) NOT NULL,
    status                   VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED',
    approval_instance_id     CHAR(36) NULL,
    created_by               CHAR(36) NULL,
    approved_at              DATETIME(6) NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_att_correction_employee (employee_id),
    INDEX idx_att_correction_session (attendance_session_id),
    INDEX idx_att_correction_status (status),
    INDEX idx_att_correction_approval_instance (approval_instance_id),

    CONSTRAINT fk_att_correction_employee FOREIGN KEY (employee_id)           REFERENCES employees(id)            ON DELETE CASCADE,
    CONSTRAINT fk_att_correction_session  FOREIGN KEY (attendance_session_id) REFERENCES attendance_sessions(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
