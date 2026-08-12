-- =============================================================================
-- Tenant Migration: 097_recruitment_stage_history (MySQL)
-- =============================================================================
-- G-5 Pipeline Stage History — lihat versi postgres untuk penjelasan kolom.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS recruitment_stages (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_recruitment_stages_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS job_application_stage_histories (
    id              CHAR(36) PRIMARY KEY,
    application_id  CHAR(36) NOT NULL,
    from_stage_id   CHAR(36) NULL,
    to_stage_id     CHAR(36) NOT NULL,
    changed_by      CHAR(36) NULL,
    notes           TEXT NULL,
    changed_at      BIGINT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ash_application FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT fk_ash_from_stage FOREIGN KEY (from_stage_id) REFERENCES recruitment_stages(id),
    CONSTRAINT fk_ash_to_stage FOREIGN KEY (to_stage_id) REFERENCES recruitment_stages(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_ash_app ON job_application_stage_histories (application_id);
CREATE INDEX idx_ash_changed_at ON job_application_stage_histories (changed_at);
