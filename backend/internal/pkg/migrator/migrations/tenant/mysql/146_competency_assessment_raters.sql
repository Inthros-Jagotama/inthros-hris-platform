-- =============================================================================
-- 146_competency_assessment_raters.sql
-- Competency 360 Module — Rater Assignment (plan generik §9).
-- Menentukan siapa yang menilai employee subject (self/superior/peer/
-- subordinate/other) pada sebuah assessment target.
-- =============================================================================

CREATE TABLE IF NOT EXISTS competency_assessment_raters (
    id                       CHAR(36) PRIMARY KEY,
    competency_event_target_id CHAR(36) NOT NULL,
    rater_employee_id        CHAR(36) NOT NULL,
    rater_type               VARCHAR(20) NOT NULL,
    weight                   DECIMAL(6,2) NOT NULL DEFAULT 0,
    status                   VARCHAR(20) NOT NULL DEFAULT 'assigned',
    assigned_at              TIMESTAMP NULL,
    submitted_at             TIMESTAMP NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_comp_rater_target_employee (competency_event_target_id, rater_employee_id),
    INDEX idx_comp_rater_employee (rater_employee_id),
    INDEX idx_comp_rater_type (rater_type),
    INDEX idx_comp_rater_status (status),
    CONSTRAINT fk_comp_rater_target   FOREIGN KEY (competency_event_target_id) REFERENCES competency_event_targets(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_rater_employee FOREIGN KEY (rater_employee_id)         REFERENCES employees(id)                 ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
