-- Down Migration: 047_job_edu_exp_multi_major_family
-- Rollback: kembalikan kolom single education_major_id / job_family_id
-- (best-effort, diambil nilai pertama dari tabel pivot), lalu hapus tabel pivot.

-- 1) Tambah kolom single kembali
ALTER TABLE job_management_education_experiences ADD COLUMN IF NOT EXISTS education_major_id CHAR(36) NULL;
ALTER TABLE job_management_education_experiences ADD COLUMN IF NOT EXISTS job_family_id CHAR(36) NULL;

-- 2) Backfill single value dari tabel pivot (best-effort)
UPDATE job_management_education_experiences e
SET education_major_id = m.education_major_id
FROM job_management_majors m
WHERE m.job_management_education_experience_id = e.id;

UPDATE job_management_education_experiences e
SET job_family_id = jf.job_family_id
FROM job_management_job_family jf
WHERE jf.job_management_education_experience_id = e.id;

-- 3) Hapus tabel pivot
DROP TABLE IF EXISTS job_management_job_family;
DROP TABLE IF EXISTS job_management_majors;
