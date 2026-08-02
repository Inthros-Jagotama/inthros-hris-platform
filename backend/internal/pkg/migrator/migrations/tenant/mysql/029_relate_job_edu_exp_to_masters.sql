-- Migration: 029_relate_job_edu_exp_to_masters
-- Relasikan job_management_education_experiences ke master:
--   education_id        → educations(id)          (Pendidikan: SMA, D3, S1, ...)
--   education_major_id  → education_majors(id)    (Jurusan)
--   job_family_id       → job_families(id)        (Bidang Pekerjaan)
--   experience_range    VARCHAR(50) NULL          (Pengalaman Kerja — hardcoded dropdown FE)
--
-- Kolom lama (job_management_value_education_id / job_management_value_experience_id)
-- yang menunjuk ke job_management_values tanpa FK dihapus (idempotent).

-- 1) Tambah kolom FK baru (idempotent)
SET @add_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN education_id CHAR(36) NULL AFTER full_code'
);
PREPARE stmt_add_edu FROM @add_edu;
EXECUTE stmt_add_edu;
DEALLOCATE PREPARE stmt_add_edu;

SET @add_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_major_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN education_major_id CHAR(36) NULL AFTER education_id'
);
PREPARE stmt_add_major FROM @add_major;
EXECUTE stmt_add_major;
DEALLOCATE PREPARE stmt_add_major;

SET @add_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_family_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN job_family_id CHAR(36) NULL AFTER education_major_id'
);
PREPARE stmt_add_jf FROM @add_jf;
EXECUTE stmt_add_jf;
DEALLOCATE PREPARE stmt_add_jf;

SET @add_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'experience_range'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN experience_range VARCHAR(50) NULL AFTER job_family_id'
);
PREPARE stmt_add_exp FROM @add_exp;
EXECUTE stmt_add_exp;
DEALLOCATE PREPARE stmt_add_exp;

-- 2) Index untuk lookup relasi (idempotent)
SET @idx_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_education'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD INDEX idx_jmee_education (education_id)'
);
PREPARE stmt_idx_edu FROM @idx_edu;
EXECUTE stmt_idx_edu;
DEALLOCATE PREPARE stmt_idx_edu;

SET @idx_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_education_major'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD INDEX idx_jmee_education_major (education_major_id)'
);
PREPARE stmt_idx_major FROM @idx_major;
EXECUTE stmt_idx_major;
DEALLOCATE PREPARE stmt_idx_major;

SET @idx_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_job_family'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD INDEX idx_jmee_job_family (job_family_id)'
);
PREPARE stmt_idx_jf FROM @idx_jf;
EXECUTE stmt_idx_jf;
DEALLOCATE PREPARE stmt_idx_jf;

-- 3) FK constraints (idempotent)
SET @fk_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL'
);
PREPARE stmt_fk_edu FROM @fk_edu;
EXECUTE stmt_fk_edu;
DEALLOCATE PREPARE stmt_fk_edu;

SET @fk_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education_major'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_education_major FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE SET NULL'
);
PREPARE stmt_fk_major FROM @fk_major;
EXECUTE stmt_fk_major;
DEALLOCATE PREPARE stmt_fk_major;

SET @fk_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_job_family'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_job_family FOREIGN KEY (job_family_id) REFERENCES job_families(id) ON DELETE SET NULL'
);
PREPARE stmt_fk_jf FROM @fk_jf;
EXECUTE stmt_fk_jf;
DEALLOCATE PREPARE stmt_fk_jf;

-- 4) Hapus kolom lama (idempotent)
SET @drop_old_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_management_value_education_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN job_management_value_education_id',
  'DO 0'
);
PREPARE stmt_drop_old_edu FROM @drop_old_edu;
EXECUTE stmt_drop_old_edu;
DEALLOCATE PREPARE stmt_drop_old_edu;

SET @drop_old_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_management_value_experience_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN job_management_value_experience_id',
  'DO 0'
);
PREPARE stmt_drop_old_exp FROM @drop_old_exp;
EXECUTE stmt_drop_old_exp;
DEALLOCATE PREPARE stmt_drop_old_exp;
