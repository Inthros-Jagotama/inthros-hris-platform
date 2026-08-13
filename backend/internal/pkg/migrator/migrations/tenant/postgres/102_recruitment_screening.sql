-- =============================================================================
-- Tenant Migration: 102_recruitment_screening (PostgreSQL)
-- =============================================================================
-- G-7 sub-project 1: application_screenings — pipeline screening record,
-- many-per-application (mirrors interviews cardinality). Purely a
-- supporting record: it does NOT auto-transition job_applications.status
-- or write to job_application_stage_histories (G-5) — recruiter still
-- moves the application status manually via the existing status endpoint.
-- (docs/module-recruitment-development-plan.md §G-7)
--
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
);

CREATE INDEX IF NOT EXISTS idx_screen_app ON application_screenings (application_id);
