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
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_comp_rater_target_employee UNIQUE (competency_event_target_id, rater_employee_id),

    CONSTRAINT fk_comp_rater_target   FOREIGN KEY (competency_event_target_id) REFERENCES competency_event_targets(id) ON DELETE CASCADE,
    CONSTRAINT fk_comp_rater_employee FOREIGN KEY (rater_employee_id)         REFERENCES employees(id)                 ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comp_rater_employee ON competency_assessment_raters (rater_employee_id);

CREATE INDEX IF NOT EXISTS idx_comp_rater_type ON competency_assessment_raters (rater_type);

CREATE INDEX IF NOT EXISTS idx_comp_rater_status ON competency_assessment_raters (status);
