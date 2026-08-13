-- =============================================================================
-- Tenant Migration: 101_candidate_consents (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 3b: candidate_consents (LAST of the originally-planned G-6
-- sub-tables — closes out education/experience/skills/certifications/
-- documents/consents; only candidates.status/source_id remain open, both
-- explicitly skipped by earlier design decisions, not deferred-to-later)
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-consents-design.md)
--
-- candidate_consents — append-only audit log of data-processing consent
-- (GDPR-style). action GRANTED|REVOKED enforced at Gin binding layer, not
-- DB constraint (module convention). No updated_at — rows are never
-- updated, only inserted. changed_by = staff user who recorded the entry
-- (this system has no public candidate-facing portal, so consent is
-- documented by HR/recruiter, not self-service).
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_consents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    action         VARCHAR(20) NOT NULL,
    notes          TEXT NULL,
    changed_by     CHAR(36) NULL,
    changed_at     BIGINT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_consent_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_consent_candidate ON candidate_consents (candidate_id);
