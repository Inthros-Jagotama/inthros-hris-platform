-- =============================================================================
-- Tenant Migration: 097_recruitment_stage_history (PostgreSQL)
-- =============================================================================
-- G-5 Pipeline Stage History
-- (docs/module-recruitment-development-plan.md §G-5)
--
-- recruitment_stages (master, seeded dari 8 CandidateStatus existing):
--   id          CHAR(36) PK
--   code        VARCHAR(20) NN UNIQUE  → NEW | SCREENED | SHORTLISTED |
--                                        INTERVIEWED | OFFERED | ACCEPTED |
--                                        REJECTED | WITHDRAWN
--   name        VARCHAR(100) NN        → label tampilan
--   sort_order  INT NN DEFAULT 0
--
-- job_application_stage_histories (audit trail transisi status aplikasi):
--   id              CHAR(36) PK
--   application_id  CHAR(36) NN  → job_applications
--   from_stage_id   CHAR(36) NULL → recruitment_stages (NULL untuk histori
--                                  pertama saat aplikasi dibuat berstatus NEW)
--   to_stage_id     CHAR(36) NN  → recruitment_stages
--   changed_by      CHAR(36) NULL → user id aktor (NULL bila sistem)
--   notes           TEXT NULL
--   changed_at      BIGINT NN    → unix nano
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS recruitment_stages (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_recruitment_stages_code UNIQUE (code)
);

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
);

CREATE INDEX IF NOT EXISTS idx_ash_app ON job_application_stage_histories (application_id);
CREATE INDEX IF NOT EXISTS idx_ash_changed_at ON job_application_stage_histories (changed_at);
