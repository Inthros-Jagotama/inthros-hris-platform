-- Migration: 047_job_edu_exp_multi_major_family
-- Dukungan multiple Jurusan & Bidang Pekerjaan pada job_management_education_experiences:
--
--   SEBELUM : education_major_id CHAR(36) (single, FK -> education_majors)
--             job_family_id      CHAR(36) (single, FK -> job_families)
--   SESUDAH : tabel pivot job_management_majors      (edu_exp <-> education_majors)
--             tabel pivot job_management_job_family  (edu_exp <-> job_families)
--
-- Data lama di-backfill ke tabel pivot sebelum kolom single dihapus.
-- Statement memakai IF EXISTS / IF NOT EXISTS (idempotent).

-- 1) Tabel pivot: Jurusan
CREATE TABLE IF NOT EXISTS job_management_majors (
    id CHAR(36) PRIMARY KEY,
    job_management_education_experience_id CHAR(36) NOT NULL,
    education_major_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_jmm_edu_exp ON job_management_majors (job_management_education_experience_id);
CREATE INDEX IF NOT EXISTS idx_jmm_education_major ON job_management_majors (education_major_id);

-- 2) Tabel pivot: Bidang Pekerjaan
CREATE TABLE IF NOT EXISTS job_management_job_family (
    id CHAR(36) PRIMARY KEY,
    job_management_education_experience_id CHAR(36) NOT NULL,
    job_family_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_jmjf_edu_exp ON job_management_job_family (job_management_education_experience_id);
CREATE INDEX IF NOT EXISTS idx_jmjf_job_family ON job_management_job_family (job_family_id);

-- 3) Backfill data lama ke tabel pivot
INSERT INTO job_management_majors (id, job_management_education_experience_id, education_major_id)
SELECT gen_random_uuid(), id, education_major_id
FROM job_management_education_experiences
WHERE education_major_id IS NOT NULL AND education_major_id <> '';

INSERT INTO job_management_job_family (id, job_management_education_experience_id, job_family_id)
SELECT gen_random_uuid(), id, job_family_id
FROM job_management_education_experiences
WHERE job_family_id IS NOT NULL AND job_family_id <> '';

-- 4) Drop FK, index, dan kolom single lama
ALTER TABLE job_management_education_experiences DROP CONSTRAINT IF EXISTS fk_jmee_education_major;
ALTER TABLE job_management_education_experiences DROP CONSTRAINT IF EXISTS fk_jmee_job_family;
DROP INDEX IF EXISTS idx_jmee_education_major;
DROP INDEX IF EXISTS idx_jmee_job_family;
ALTER TABLE job_management_education_experiences DROP COLUMN IF EXISTS education_major_id;
ALTER TABLE job_management_education_experiences DROP COLUMN IF EXISTS job_family_id;

-- 5) FK tabel pivot -> parent & master (idempotent via DO block)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmm_edu_exp') THEN
        ALTER TABLE job_management_majors
            ADD CONSTRAINT fk_jmm_edu_exp
            FOREIGN KEY (job_management_education_experience_id)
            REFERENCES job_management_education_experiences(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmm_education_major') THEN
        ALTER TABLE job_management_majors
            ADD CONSTRAINT fk_jmm_education_major
            FOREIGN KEY (education_major_id)
            REFERENCES education_majors(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmjf_edu_exp') THEN
        ALTER TABLE job_management_job_family
            ADD CONSTRAINT fk_jmjf_edu_exp
            FOREIGN KEY (job_management_education_experience_id)
            REFERENCES job_management_education_experiences(id) ON DELETE CASCADE;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmjf_job_family') THEN
        ALTER TABLE job_management_job_family
            ADD CONSTRAINT fk_jmjf_job_family
            FOREIGN KEY (job_family_id)
            REFERENCES job_families(id) ON DELETE CASCADE;
    END IF;
END $$;
