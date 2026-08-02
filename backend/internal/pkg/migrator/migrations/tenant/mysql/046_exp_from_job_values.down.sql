-- Down Migration: 046_exp_from_job_values
-- Kembalikan kolom experience_range VARCHAR(50) dan hapus experience_id + FK.

-- 1) Drop FK (idempotent)
SET @drop_fk = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_experience'
  ),
  'ALTER TABLE job_management_education_experiences DROP FOREIGN KEY fk_jmee_experience',
  'DO 0'
);
PREPARE stmt FROM @drop_fk;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2) Drop index (idempotent)
SET @drop_idx = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND index_name = 'idx_jmee_experience'
  ),
  'ALTER TABLE job_management_education_experiences DROP INDEX idx_jmee_experience',
  'DO 0'
);
PREPARE stmt2 FROM @drop_idx;
EXECUTE stmt2;
DEALLOCATE PREPARE stmt2;

-- 3) Drop kolom experience_id
SET @drop_col = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'experience_id'
  ),
  'ALTER TABLE job_management_education_experiences DROP COLUMN experience_id',
  'DO 0'
);
PREPARE stmt3 FROM @drop_col;
EXECUTE stmt3;
DEALLOCATE PREPARE stmt3;

-- 4) Tambah kembali kolom experience_range VARCHAR(50)
SET @add_exp = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'experience_range'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN experience_range VARCHAR(50) NULL'
);
PREPARE stmt4 FROM @add_exp;
EXECUTE stmt4;
DEALLOCATE PREPARE stmt4;
