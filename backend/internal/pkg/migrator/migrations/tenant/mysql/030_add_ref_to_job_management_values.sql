-- Migration: 030_add_ref_to_job_management_values
-- Menambahkan kolom ref_id & ref_type ke tabel job_management_values
-- (relasi polymorphic: ref_type = entitas yang direferensikan, ref_id = id record-nya).
--
-- CATATAN (MySQL): Migration ini dibutuhkan untuk tenant lama yang tabel
-- job_management_values-nya dibuat sebelum kolom ref_id/ref_type ada, jadi
-- penambahan dilakukan secara kondisional (idempotent).

SET @add_ref_id_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND column_name = 'ref_id'
  ),
  'DO 0',
  'ALTER TABLE job_management_values ADD COLUMN ref_id CHAR(36) NULL AFTER sort'
);
PREPARE stmt_add_ref_id FROM @add_ref_id_sql;
EXECUTE stmt_add_ref_id;
DEALLOCATE PREPARE stmt_add_ref_id;

SET @add_ref_type_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND column_name = 'ref_type'
  ),
  'DO 0',
  'ALTER TABLE job_management_values ADD COLUMN ref_type VARCHAR(100) NULL AFTER ref_id'
);
PREPARE stmt_add_ref_type FROM @add_ref_type_sql;
EXECUTE stmt_add_ref_type;
DEALLOCATE PREPARE stmt_add_ref_type;

-- Index untuk lookup relasi (idempotent)
SET @idx_sql = IF(
  EXISTS(
    SELECT 1 FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'job_management_values'
      AND index_name = 'idx_jmv_ref'
  ),
  'DO 0',
  'ALTER TABLE job_management_values ADD INDEX idx_jmv_ref (ref_id)'
);
PREPARE stmt_idx FROM @idx_sql;
EXECUTE stmt_idx;
DEALLOCATE PREPARE stmt_idx;
