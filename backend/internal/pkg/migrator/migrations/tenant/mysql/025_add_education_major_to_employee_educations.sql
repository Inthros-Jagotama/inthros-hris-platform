-- Migration: 025_add_education_major_to_employee_educations
-- Menambahkan kolom education_major_id ke tabel employee_educations
-- (relasi ke tabel master education_majors dari module setting).
--
-- CATATAN (MySQL): Migration ini dibutuhkan untuk tenant lama yang tabel
-- employee_educations-nya dibuat sebelum kolom education_major_id ada, jadi
-- penambahan dilakukan secara kondisional (idempotent).

SET @add_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_educations'
      AND column_name = 'education_major_id'
  ),
  'DO 0',
  'ALTER TABLE employee_educations ADD COLUMN education_major_id CHAR(36) NULL AFTER education_id'
);
PREPARE stmt_add FROM @add_sql;
EXECUTE stmt_add;
DEALLOCATE PREPARE stmt_add;

-- Index untuk lookup relasi (idempotent)
SET @idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_educations'
      AND index_name = 'idx_empedu_education_major'
  ),
  'DO 0',
  'ALTER TABLE employee_educations ADD INDEX idx_empedu_education_major (education_major_id)'
);
PREPARE stmt_idx FROM @idx_sql;
EXECUTE stmt_idx;
DEALLOCATE PREPARE stmt_idx;
