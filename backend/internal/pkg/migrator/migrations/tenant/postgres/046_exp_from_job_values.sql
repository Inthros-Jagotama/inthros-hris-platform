-- Migration: 046_exp_from_job_values
-- Ubah Pengalaman Kerja pada job_management_education_experiences:
--   SEBELUM: experience_range VARCHAR(50) — string hardcoded dropdown FE (0-2 Tahun, ...)
--   SESUDAH : experience_id CHAR(36) → job_management_values(id) WHERE type='experience'
--
-- Konsisten dengan Education di 045: opsi Pengalaman Kerja diambil dari
-- tabel job_management_values (type=experience, level 1-5).

DO $$
BEGIN
    -- 1) Tambah kolom experience_id (jika belum ada)
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'job_management_education_experiences'
          AND column_name = 'experience_id'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD COLUMN experience_id CHAR(36) NULL;
    END IF;

    -- 2) Backfill: cocokkan experience_range lama dengan descriptions di job_management_values (type=experience)
    UPDATE job_management_education_experiences e
    SET experience_id = v.id
    FROM job_management_values v
    WHERE v.type = 'experience'
      AND v.descriptions = e.experience_range
      AND e.experience_range IS NOT NULL
      AND e.experience_range <> '';

    -- 3) Drop kolom experience_range
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'job_management_education_experiences'
          AND column_name = 'experience_range'
    ) THEN
        ALTER TABLE job_management_education_experiences DROP COLUMN experience_range;
    END IF;

    -- 4) Index untuk experience_id
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = current_schema()
          AND tablename = 'job_management_education_experiences'
          AND indexname = 'idx_jmee_experience'
    ) THEN
        CREATE INDEX idx_jmee_experience ON job_management_education_experiences (experience_id);
    END IF;

    -- 5) FK experience_id → job_management_values(id)
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_jmee_experience'
    ) THEN
        ALTER TABLE job_management_education_experiences
            ADD CONSTRAINT fk_jmee_experience
            FOREIGN KEY (experience_id) REFERENCES job_management_values(id) ON DELETE SET NULL;
    END IF;
END $$;
