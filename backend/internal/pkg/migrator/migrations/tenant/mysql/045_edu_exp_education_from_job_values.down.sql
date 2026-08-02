-- Down Migration: 045_edu_exp_education_from_job_values
-- Kembalikan relasi Pendidikan ke master educations(id).

SET @drop_fk = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education'
  ),
  'ALTER TABLE job_management_education_experiences DROP FOREIGN KEY fk_jmee_education',
  'DO 0'
);
PREPARE stmt_drop_fk FROM @drop_fk;
EXECUTE stmt_drop_fk;
DEALLOCATE PREPARE stmt_drop_fk;

SET @add_fk = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;
