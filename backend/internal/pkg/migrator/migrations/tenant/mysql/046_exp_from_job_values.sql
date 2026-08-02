-- Migration: 046_exp_from_job_values
-- Ubah Pengalaman Kerja pada job_management_education_experiences:
--   SEBELUM: experience_range VARCHAR(50) — string hardcoded dropdown FE (0-2 Tahun, ...)
--   SESUDAH : experience_id CHAR(36) → job_management_values(id) WHERE type='experience'
--
-- Sesuai keputusan user (konsisten dengan Education di 045): opsi Pengalaman Kerja
-- diambil dari tabel job_management_values (type=experience, level 1-5:
-- 0-2 Tahun, 3-5 Tahun, 6-8 Tahun, 9-11 Tahun, > 12 Tahun).
--
-- Idempotent: backfill data lama experience_range → experience_id (match descriptions),
-- lalu drop kolom experience_range, tambah index + FK ke job_management_values.

-- 1) Tambah kolom experience_id (idempotent)
SET @add_exp_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'experience_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN experience_id CHAR(36) NULL AFTER experience_range'
);
PREPARE stmt_add_exp_id FROM @add_exp_id;
EXECUTE stmt_add_exp_id;
DEALLOCATE PREPARE stmt_add_exp_id;

-- 2) Backfill: cocokkan experience_range lama dengan descriptions di job_management_values (type=experience)
UPDATE job_management_education_experiences e
JOIN job_management_values v ON v.type = 'experience' AND v.descriptions = e.experience_range
SET e.experience_id = v.id
WHERE e.experience_range IS NOT NULL AND e.experience_range <> '';

-- 3) Drop kolom experience_range (idempotent — tabel masih baru/kosong)
SET @drop_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'experience_range'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN experience_range',
  'DO 0'
);
PREPARE stmt_drop_exp FROM @drop_exp;
EXECUTE stmt_drop_exp;
DEALLOCATE PREPARE stmt_drop_exp;

-- 4) Index untuk experience_id (idempotent)
SET @add_idx = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_experience'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD INDEX idx_jmee_experience (experience_id)'
);
PREPARE stmt_add_idx FROM @add_idx;
EXECUTE stmt_add_idx;
DEALLOCATE PREPARE stmt_add_idx;

-- 5) FK experience_id → job_management_values(id) (idempotent)
SET @add_fk = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_experience'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_experience FOREIGN KEY (experience_id) REFERENCES job_management_values(id) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;
