-- Migration: 029_relate_job_edu_exp_to_masters
-- Relasikan job_management_education_experiences ke master:
--   education_id        → educations(id)          (Pendidikan)
--   education_major_id  → education_majors(id)    (Jurusan)
--   job_family_id       → job_families(id)        (Bidang Pekerjaan)
--   experience_range    VARCHAR(50) NULL          (Pengalaman Kerja — hardcoded dropdown FE)
--
-- Kolom lama (job_management_value_education_id / job_management_value_experience_id)
-- yang menunjuk ke job_management_values tanpa FK dihapus (idempotent).

ALTER TABLE job_management_education_experiences
    ADD COLUMN IF NOT EXISTS education_id       CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS education_major_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS job_family_id      CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS experience_range   VARCHAR(50) NULL;

CREATE INDEX IF NOT EXISTS idx_jmee_education
    ON job_management_education_experiences (education_id);
CREATE INDEX IF NOT EXISTS idx_jmee_education_major
    ON job_management_education_experiences (education_major_id);
CREATE INDEX IF NOT EXISTS idx_jmee_job_family
    ON job_management_education_experiences (job_family_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_education
            FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education_major'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_education_major
            FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_job_family'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_job_family
            FOREIGN KEY (job_family_id) REFERENCES job_families(id) ON DELETE SET NULL;
    END IF;
END $$;

ALTER TABLE job_management_education_experiences
    DROP COLUMN IF EXISTS job_management_value_education_id,
    DROP COLUMN IF EXISTS job_management_value_experience_id;
