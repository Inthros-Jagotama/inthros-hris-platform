-- =============================================================================
-- Tenant Migration: 098_candidate_profile_basics (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE; CREATE TABLE IF NOT EXISTS.

SET @add_candidate_number = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'candidate_number'
  ),
  'DO 0',
  'ALTER TABLE candidates ADD COLUMN candidate_number VARCHAR(50) NULL'
);
PREPARE stmt FROM @add_candidate_number;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_edu_candidate ON candidate_educations (candidate_id);

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_exp_candidate ON candidate_work_experiences (candidate_id);
