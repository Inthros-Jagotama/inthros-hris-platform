-- Down Migration: 029_relate_job_edu_exp_to_masters
-- Mengembalikan kolom lama job_management_value_* dan menghapus kolom FK baru.

-- 1) Hapus FK constraints (idempotent)
SET @drop_fk_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education'
  ),
  'ALTER TABLE job_management_education_experiences DROP FOREIGN KEY fk_jmee_education',
  'DO 0'
);
PREPARE stmt FROM @drop_fk_edu;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

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

-- 2) Restore kolom lama (idempotent)
SET @add_old_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_management_value_education_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN job_management_value_education_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_old_edu;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_old_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_management_value_experience_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN job_management_value_experience_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_old_exp;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3) Drop kolom FK baru (idempotent)
SET @drop_edu = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN education_id',
  'DO 0'
);
PREPARE stmt FROM @drop_edu;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_major_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN education_major_id',
  'DO 0'
);
PREPARE stmt FROM @drop_major;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_family_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN job_family_id',
  'DO 0'
);
PREPARE stmt FROM @drop_jf;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

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
PREPARE stmt FROM @drop_exp;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
