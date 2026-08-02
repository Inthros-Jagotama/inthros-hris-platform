-- Down Migration: 046_exp_from_job_values
-- Kembalikan kolom experience_range VARCHAR(50) dan hapus experience_id + FK.

DO $$
BEGIN
    -- 1) Drop FK
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_experience'
    ) THEN
        ALTER TABLE job_management_education_experiences DROP CONSTRAINT fk_jmee_experience;
    END IF;

    -- 2) Drop index
    IF EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'job_management_education_experiences'
          AND indexname = 'idx_jmee_experience'
    ) THEN
        DROP INDEX idx_jmee_experience;
    END IF;

    -- 3) Drop kolom experience_id
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'job_management_education_experiences'
          AND column_name = 'experience_id'
    ) THEN
        ALTER TABLE job_management_education_experiences DROP COLUMN experience_id;
    END IF;

    -- 4) Tambah kembali kolom experience_range VARCHAR(50)
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'job_management_education_experiences'
          AND column_name = 'experience_range'
    ) THEN
        ALTER TABLE job_management_education_experiences ADD COLUMN experience_range VARCHAR(50);
    END IF;
END $$;
