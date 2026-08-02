-- Migration: 045_edu_exp_education_from_job_values
-- Ubah relasi Pendidikan (education_id) pada job_management_education_experiences:
--   SEBELUM: education_id → educations(id)  (master settings)
--   SESUDAH : education_id → job_management_values(id) WHERE type='education'
--
-- Sesuai keputusan user: opsi Pendidikan diambil dari tabel job_management_values
-- (type=education, level 1-5: Sekolah Menengah Pertama → Strata 3), bukan master educations.
-- Backend juga diubah: relasi model Education → JobValue (scope type='education').
--
-- Idempotent: drop FK jika ada (tanpa error bila belum ada), lalu tambah ulang
-- menunjuk ke job_management_values. Data lama kosong (tabel baru) → aman.

-- 1) Drop FK lama menuju educations (idempotent — error apa pun diabaikan via dynamic SQL)
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

-- 2) Tambah ulang FK menuju job_management_values(id) (idempotent)
SET @add_fk = IF(
  EXISTS(
    SELECT 1 FROM information_schema.TABLE_CONSTRAINTS
    WHERE CONSTRAINT_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_management_education_experiences'
      AND CONSTRAINT_NAME = 'fk_jmee_education'
  ),
  'DO 0',
  'ALTER TABLE job_management_education_experiences ADD CONSTRAINT fk_jmee_education FOREIGN KEY (education_id) REFERENCES job_management_values(id) ON DELETE SET NULL'
);
PREPARE stmt_add_fk FROM @add_fk;
EXECUTE stmt_add_fk;
DEALLOCATE PREPARE stmt_add_fk;
