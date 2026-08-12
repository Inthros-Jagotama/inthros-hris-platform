-- =============================================================================
-- Tenant Migration: 098_candidate_profile_basics (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 1: candidate_number + candidate_educations +
-- candidate_work_experiences
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-profile-basics-design.md)
--
-- candidates.candidate_number VARCHAR(50) NULL — auto-generated
-- CAND-YYYYMM-XXXXXXXX (pola requisition_number/offer_number).
--
-- candidate_educations — riwayat pendidikan kandidat, education_id/
-- education_major_id merujuk ke master setting.Education/EducationMajor
-- (pola employee_educations existing) dengan fallback teks bebas.
--
-- candidate_work_experiences — riwayat pekerjaan kandidat, tanpa master
-- (tidak ada master employer/job-title di sistem).
--
-- Idempotent: ADD COLUMN / CREATE TABLE IF NOT EXISTS.

ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS candidate_number VARCHAR(50) NULL;

CREATE TABLE IF NOT EXISTS candidate_educations (
    id                  CHAR(36) PRIMARY KEY,
    candidate_id        CHAR(36) NOT NULL,
    education_id        CHAR(36) NULL,
    institution_name    VARCHAR(255) NOT NULL,
    education_major_id  CHAR(36) NULL,
    major               VARCHAR(255) NULL,
    gpa                 DECIMAL(3,2) NULL,
    start_year          INT NULL,
    end_year            INT NULL,
    is_highest          BOOLEAN NOT NULL DEFAULT false,
    notes               TEXT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_edu_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_edu_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL,
    CONSTRAINT fk_cand_edu_major FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_cand_edu_candidate ON candidate_educations (candidate_id);

CREATE TABLE IF NOT EXISTS candidate_work_experiences (
    id               CHAR(36) PRIMARY KEY,
    candidate_id     CHAR(36) NOT NULL,
    company_name     VARCHAR(255) NOT NULL,
    job_title        VARCHAR(255) NOT NULL,
    employment_type  VARCHAR(50) NULL,
    start_date       DATE NOT NULL,
    end_date         DATE NULL,
    is_current       BOOLEAN NOT NULL DEFAULT false,
    description      TEXT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_exp_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_exp_candidate ON candidate_work_experiences (candidate_id);
