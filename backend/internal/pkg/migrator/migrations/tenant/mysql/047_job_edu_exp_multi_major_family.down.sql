-- Down Migration: 047_job_edu_exp_multi_major_family
-- Rollback: kembalikan kolom single education_major_id / job_family_id
-- (best-effort, diambil nilai pertama dari tabel pivot), lalu hapus tabel pivot.

-- 1) Tambah kolom single kembali (idempotent)
SET @add_col_major = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'education_major_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN education_major_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_col_major;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_col_jf = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_education_experiences'
      AND column_name = 'job_family_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD COLUMN job_family_id CHAR(36) NULL'
);
PREPARE stmt FROM @add_col_jf;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2) Backfill single value dari tabel pivot (best-effort)
UPDATE job_management_education_experiences e
LEFT JOIN job_management_majors m ON m.job_management_education_experience_id = e.id
SET e.education_major_id = m.education_major_id;

UPDATE job_management_education_experiences e
LEFT JOIN job_management_job_family jf ON jf.job_management_education_experience_id = e.id
SET e.job_family_id = jf.job_family_id;

-- 3) Hapus tabel pivot
DROP TABLE IF EXISTS job_management_job_family;
DROP TABLE IF EXISTS job_management_majors;
