-- Migration: 047_job_edu_exp_multi_major_family
-- Dukungan multiple Jurusan & Bidang Pekerjaan pada job_management_education_experiences:
--
--   SEBELUM : education_major_id CHAR(36) (single, FK -> education_majors)
--             job_family_id      CHAR(36) (single, FK -> job_families)
--   SESUDAH : tabel pivot job_management_majors      (edu_exp <-> education_majors)
--             tabel pivot job_management_job_family  (edu_exp <-> job_families)
--
-- Data lama di-backfill ke tabel pivot sebelum kolom single dihapus.
-- Semua statement idempotent (aman dijalankan ulang).

-- 1) Tabel pivot: Jurusan
CREATE TABLE IF NOT EXISTS job_management_majors (
    id CHAR(36) NOT NULL,
    job_management_education_experience_id CHAR(36) NOT NULL,
    education_major_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_jmm_edu_exp (job_management_education_experience_id),
    INDEX idx_jmm_education_major (education_major_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2) Tabel pivot: Bidang Pekerjaan
CREATE TABLE IF NOT EXISTS job_management_job_family (
    id CHAR(36) NOT NULL,
    job_management_education_experience_id CHAR(36) NOT NULL,
    job_family_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_jmjf_edu_exp (job_management_education_experience_id),
    INDEX idx_jmjf_job_family (job_family_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3) Backfill data lama ke tabel pivot
INSERT INTO job_management_majors (id, job_management_education_experience_id, education_major_id)
SELECT UUID(), id, education_major_id
FROM job_management_education_experiences
WHERE education_major_id IS NOT NULL AND education_major_id <> '';

INSERT INTO job_management_job_family (id, job_management_education_experience_id, job_family_id)
SELECT UUID(), id, job_family_id
FROM job_management_education_experiences
WHERE job_family_id IS NOT NULL AND job_family_id <> '';

-- 4) Drop FK lama di job_management_education_experiences (idempotent)
SET @drop_fk_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education_major'
  ),
  'ALTER TABLE job_management_education_experiences DROP FOREIGN KEY fk_jmee_education_major',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_major;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_fk_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_job_family'
  ),
  'ALTER TABLE job_management_education_experiences DROP FOREIGN KEY fk_jmee_job_family',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_jf;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 5) Drop index lama (idempotent)
SET @drop_idx_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_education_major'
  ),
  'ALTER TABLE job_management_education_experiences DROP INDEX idx_jmee_education_major',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_major;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_idx_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_job_family'
  ),
  'ALTER TABLE job_management_education_experiences DROP INDEX idx_jmee_job_family',
  'DO 0'
);
PREPARE stmt FROM @drop_idx_jf;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 6) Drop kolom single lama (idempotent)
SET @drop_col_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_major_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN education_major_id',
  'DO 0'
);
PREPARE stmt FROM @drop_col_major;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_col_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_family_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN job_family_id',
  'DO 0'
);
PREPARE stmt FROM @drop_col_jf;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 7) FK tabel pivot -> parent & master (idempotent)
SET @fk_jmm_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_majors'
      AND CONSTRAINT_NAME = 'fk_jmm_edu_exp'
  ),
  'DO 0',
  'ALTER TABLE job_management_majors ADD CONSTRAINT fk_jmm_edu_exp FOREIGN KEY (job_management_education_experience_id) REFERENCES job_management_education_experiences(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_jmm_exp;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @fk_jmm_master = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_majors'
      AND CONSTRAINT_NAME = 'fk_jmm_education_major'
  ),
  'DO 0',
  'ALTER TABLE job_management_majors ADD CONSTRAINT fk_jmm_education_major FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_jmm_master;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @fk_jmjf_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_job_family'
      AND CONSTRAINT_NAME = 'fk_jmjf_edu_exp'
  ),
  'DO 0',
  'ALTER TABLE job_management_job_family ADD CONSTRAINT fk_jmjf_edu_exp FOREIGN KEY (job_management_education_experience_id) REFERENCES job_management_education_experiences(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_jmjf_exp;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @fk_jmjf_master = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_job_family'
      AND CONSTRAINT_NAME = 'fk_jmjf_job_family'
  ),
  'DO 0',
  'ALTER TABLE job_management_job_family ADD CONSTRAINT fk_jmjf_job_family FOREIGN KEY (job_family_id) REFERENCES job_families(id) ON DELETE CASCADE'
);
PREPARE stmt FROM @fk_jmjf_master;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
