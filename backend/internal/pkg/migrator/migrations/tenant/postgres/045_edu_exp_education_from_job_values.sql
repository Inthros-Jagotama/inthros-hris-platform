-- Migration: 045_edu_exp_education_from_job_values
-- Ubah relasi Pendidikan (education_id) pada job_management_education_experiences:
--   SEBELUM: education_id → educations(id)  (master settings)
--   SESUDAH : education_id → job_management_values(id) WHERE type='education'
--
-- Sesuai keputusan user: opsi Pendidikan diambil dari tabel job_management_values
-- (type=education, level 1-5), bukan master educations.

DO $$
BEGIN
    -- 1) Drop FK lama menuju educations (jika ada)
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education'
    ) THEN
        ALTER TABLE job_management_education_experiences
            DROP CONSTRAINT fk_jmee_education;
    END IF;

    -- 2) Tambah ulang FK menuju job_management_values(id) (jika belum ada)
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_education'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_education
            FOREIGN KEY (education_id) REFERENCES job_management_values(id) ON DELETE SET NULL;
    END IF;
END $$;
