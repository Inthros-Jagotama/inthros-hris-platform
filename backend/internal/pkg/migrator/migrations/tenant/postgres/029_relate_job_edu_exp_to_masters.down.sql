-- Down Migration: 029_relate_job_edu_exp_to_masters (PostgreSQL)
-- Mengembalikan kolom lama job_management_value_* dan menghapus kolom FK baru.

ALTER TABLE job_management_education_experiences
    DROP CONSTRAINT IF EXISTS fk_jmee_education,
    DROP CONSTRAINT IF EXISTS fk_jmee_education_major,
    DROP CONSTRAINT IF EXISTS fk_jmee_job_family;

ALTER TABLE job_management_education_experiences
    ADD COLUMN IF NOT EXISTS job_management_value_education_id   CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS job_management_value_experience_id  CHAR(36) NULL;

ALTER TABLE job_management_education_experiences
    DROP COLUMN IF EXISTS education_id,
    DROP COLUMN IF EXISTS education_major_id,
    DROP COLUMN IF EXISTS job_family_id,
    DROP COLUMN IF EXISTS experience_range;
