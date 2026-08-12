-- =============================================================================
-- Tenant Migration: 099_candidate_skills_certifications (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 2: candidate_skills + candidate_certifications
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-skills-certifications-design.md)
--
-- candidate_skills — kandidat + skill (referensi competency.Competency
-- master existing, bukan Skill Master baru) + level proficiency 1-5.
-- competency_id NOT NULL (beda dari candidate_educations yang nullable) —
-- skill tanpa referensi competency tidak actionable untuk candidate
-- matching (G-9). ON DELETE CASCADE (satu-satunya opsi valid untuk FK
-- NOT NULL, konsisten dengan child table competency module sendiri).
--
-- candidate_certifications — tanpa master (tidak ada Certification Master
-- di sistem), field bebas + info penerbit/masa berlaku.
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_skills (
    id            CHAR(36) PRIMARY KEY,
    candidate_id  CHAR(36) NOT NULL,
    competency_id CHAR(36) NOT NULL,
    level         SMALLINT NULL,
    notes         TEXT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_skill_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_skill_competency FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_skill_candidate ON candidate_skills (candidate_id);
CREATE INDEX IF NOT EXISTS idx_cand_skill_competency ON candidate_skills (competency_id);

CREATE TABLE IF NOT EXISTS candidate_certifications (
    id                   CHAR(36) PRIMARY KEY,
    candidate_id         CHAR(36) NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    issuing_organization VARCHAR(255) NULL,
    issue_date           DATE NULL,
    expiry_date          DATE NULL,
    credential_url       TEXT NULL,
    notes                TEXT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_cert_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_cert_candidate ON candidate_certifications (candidate_id);
