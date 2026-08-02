-- Down Migration: 030_add_ref_to_job_management_values
-- Hapus kolom ref_id & ref_type dari tabel job_management_values (kondisional, idempotent).

SET @drop_idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND index_name = 'idx_jmv_ref'
  ),
  'ALTER TABLE job_management_values DROP INDEX idx_jmv_ref',
  'DO 0'
);
PREPARE stmt_drop_idx FROM @drop_idx_sql;
EXECUTE stmt_drop_idx;
DEALLOCATE PREPARE stmt_drop_idx;

SET @drop_ref_type_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND column_name = 'ref_type'
  ),
  'ALTER TABLE job_management_values DROP COLUMN ref_type',
  'DO 0'
);
PREPARE stmt_drop_ref_type FROM @drop_ref_type_sql;
EXECUTE stmt_drop_ref_type;
DEALLOCATE PREPARE stmt_drop_ref_type;

SET @drop_ref_id_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND column_name = 'ref_id'
  ),
  'ALTER TABLE job_management_values DROP COLUMN ref_id',
  'DO 0'
);
PREPARE stmt_drop_ref_id FROM @drop_ref_id_sql;
EXECUTE stmt_drop_ref_id;
DEALLOCATE PREPARE stmt_drop_ref_id;
