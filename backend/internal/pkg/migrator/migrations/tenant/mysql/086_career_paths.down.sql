-- =============================================================================
-- Tenant Migration Down: 086_career_paths (unified schema rollback)
-- =============================================================================
-- Mengembalikan career_paths ke bentuk edge career intelligence (018):
--   - hapus career_path_steps
--   - hapus kolom header EM (name, description, created_by, updated_by) + uk
--   - tambah kembali kolom edge CI (source/target/path_type/typical_tenure/
--     requirements/competencies/certifications) + index source/target.
-- Semua statement idempotent.

DROP TABLE IF EXISTS career_path_steps;

SET @drop_cp_name_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'uk_career_paths_name'),
  'ALTER TABLE career_paths DROP INDEX uk_career_paths_name',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_name_uk;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_name = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'name'),
  'ALTER TABLE career_paths DROP COLUMN name',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_description = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'description'),
  'ALTER TABLE career_paths DROP COLUMN description',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_description;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_created_by = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'created_by'),
  'ALTER TABLE career_paths DROP COLUMN created_by',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_updated_by = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'updated_by'),
  'ALTER TABLE career_paths DROP COLUMN updated_by',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_updated_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Kembalikan kolom edge CI.
SET @add_cp_source = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'source_title_id'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN source_title_id CHAR(36) NOT NULL AFTER id'
);
PREPARE stmt FROM @add_cp_source;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_target = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'target_title_id'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN target_title_id CHAR(36) NOT NULL AFTER source_title_id'
);
PREPARE stmt FROM @add_cp_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_path_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'path_type'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN path_type VARCHAR(30) NOT NULL AFTER target_title_id'
);
PREPARE stmt FROM @add_cp_path_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_typical_tenure = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'typical_tenure'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN typical_tenure INT NOT NULL DEFAULT 0 AFTER path_type'
);
PREPARE stmt FROM @add_cp_typical_tenure;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_requirements = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'requirements'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN requirements TEXT NULL AFTER typical_tenure'
);
PREPARE stmt FROM @add_cp_requirements;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_competencies = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'competencies'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN competencies TEXT NULL AFTER requirements'
);
PREPARE stmt FROM @add_cp_competencies;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_certifications = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'certifications'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN certifications TEXT NULL AFTER competencies'
);
PREPARE stmt FROM @add_cp_certifications;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_idx_source = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'idx_cp_source'),
  'DO 0',
  'CREATE INDEX idx_cp_source ON career_paths (source_title_id)'
);
PREPARE stmt FROM @add_cp_idx_source;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_idx_target = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'idx_cp_target'),
  'DO 0',
  'CREATE INDEX idx_cp_target ON career_paths (target_title_id)'
);
PREPARE stmt FROM @add_cp_idx_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
