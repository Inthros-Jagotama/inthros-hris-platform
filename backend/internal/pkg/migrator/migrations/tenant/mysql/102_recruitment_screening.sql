-- =============================================================================
-- Tenant Migration: 102_recruitment_screening (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS application_screenings (
    id             CHAR(36) PRIMARY KEY,
    application_id CHAR(36) NOT NULL,
    screened_by    CHAR(36) NULL,
    screened_at    BIGINT NOT NULL DEFAULT 0,
    score          DECIMAL(5,2) NULL,
    result         VARCHAR(10) NOT NULL DEFAULT 'HOLD',
    notes          TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_screen_app FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_screen_app ON application_screenings (application_id);
