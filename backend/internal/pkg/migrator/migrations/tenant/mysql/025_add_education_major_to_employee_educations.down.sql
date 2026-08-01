-- Down Migration: 025_add_education_major_to_employee_educations
-- Hapus kolom education_major_id dari tabel employee_educations (kondisional, idempotent).

SET @drop_idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_educations'
      AND index_name = 'idx_empedu_education_major'
  ),
  'ALTER TABLE employee_educations DROP INDEX idx_empedu_education_major',
  'DO 0'
);
PREPARE stmt_drop_idx FROM @drop_idx_sql;
EXECUTE stmt_drop_idx;
DEALLOCATE PREPARE stmt_drop_idx;

SET @drop_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'employee_educations'
      AND column_name = 'education_major_id'
  ),
  'ALTER TABLE employee_educations DROP COLUMN education_major_id',
  'DO 0'
);
PREPARE stmt_drop FROM @drop_sql;
EXECUTE stmt_drop;
DEALLOCATE PREPARE stmt_drop;
