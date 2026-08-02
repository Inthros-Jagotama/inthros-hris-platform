-- Down Migration: 045_edu_exp_education_from_job_values
-- Kembalikan relasi Pendidikan ke master educations(id).

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education'
    ) THEN
        ALTER TABLE job_management_education_experiences
            DROP CONSTRAINT fk_jmee_education;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_education
            FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL;
    END IF;
END $$;
