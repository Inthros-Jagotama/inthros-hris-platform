-- =============================================================================
-- Tenant Migration: 086_career_paths — UNIFIED CAREER PATH SCHEMA
-- =============================================================================
-- Employee Movement Enhancement (plan §12.9) + UNIFIKASI dengan Career
-- Intelligence (keputusan 2026-08-10): career_paths menjadi SATU sumber
-- kebenaran untuk kedua modul.
--
-- Sebelumnya (migration 018, career intelligence): career_paths = tabel EDGE
--   (source_title_id → target_title_id, path_type, typical_tenure,
--    requirements, competencies, certifications).
-- Sekarang (086): career_paths = HEADER jenjang (name/description/is_active/
--    created_by/updated_by + soft delete), dan career_path_steps = langkah.
--    Edge CI direpresentasikan sebagai path 2-langkah:
--      step 1 = source (sequence 1)
--      step 2 = target (sequence 2) — atribut CI (path_type, typical_tenure,
--               competencies, certifications) disimpan di step target.
--
-- Karena career_paths SUDAH ada (dibuat migration 018), 086 TIDAK membuat
-- ulang tabel — ia ALTER: menambah kolom EM, menghapus kolom edge CI lama
-- (dipindah ke career_path_steps), lalu membuat career_path_steps.
-- Semua statement idempotent (aman dijalankan ulang / pada tenant yang belum
-- punya kolom 018 — dibuat di migration yang sama saat provisioning pertama).

-- -----------------------------------------------------------------------------
-- 1. Tambah kolom header EM (idempotent) ke career_paths yang sudah ada.
-- -----------------------------------------------------------------------------

SET @add_cp_name = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'name'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN name VARCHAR(100) NULL AFTER id'
);
PREPARE stmt FROM @add_cp_name;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_description = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'description'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN description TEXT NULL AFTER name'
);
PREPARE stmt FROM @add_cp_description;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_created_by = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'created_by'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN created_by CHAR(36) NULL AFTER is_active'
);
PREPARE stmt FROM @add_cp_created_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_cp_updated_by = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'updated_by'),
  'DO 0',
  'ALTER TABLE career_paths ADD COLUMN updated_by CHAR(36) NULL AFTER created_by'
);
PREPARE stmt FROM @add_cp_updated_by;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 2. Unique index name (idempotent) — name diisi service (EM: nama jenjang;
--    CI: nama otomatis "PROMOTION: Source → Target").
-- -----------------------------------------------------------------------------

SET @add_cp_name_uk = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'uk_career_paths_name'),
  'DO 0',
  'ALTER TABLE career_paths ADD UNIQUE KEY uk_career_paths_name (name)'
);
PREPARE stmt FROM @add_cp_name_uk;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 3. Hapus kolom edge CI lama (idempotent) — dipindah ke career_path_steps.
--    Index source/target dihapus lebih dulu (mengacu kolom tsb).
-- -----------------------------------------------------------------------------

SET @drop_cp_idx_source = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'idx_cp_source'),
  'DROP INDEX idx_cp_source ON career_paths',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_idx_source;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_idx_target = IF(
  EXISTS(SELECT 1 FROM information_schema.STATISTICS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND INDEX_NAME = 'idx_cp_target'),
  'DROP INDEX idx_cp_target ON career_paths',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_idx_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_source = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'source_title_id'),
  'ALTER TABLE career_paths DROP COLUMN source_title_id',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_source;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_target = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'target_title_id'),
  'ALTER TABLE career_paths DROP COLUMN target_title_id',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_target;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_path_type = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'path_type'),
  'ALTER TABLE career_paths DROP COLUMN path_type',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_path_type;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_typical_tenure = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'typical_tenure'),
  'ALTER TABLE career_paths DROP COLUMN typical_tenure',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_typical_tenure;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_requirements = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'requirements'),
  'ALTER TABLE career_paths DROP COLUMN requirements',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_requirements;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_competencies = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'competencies'),
  'ALTER TABLE career_paths DROP COLUMN competencies',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_competencies;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_cp_certifications = IF(
  EXISTS(SELECT 1 FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'career_paths' AND COLUMN_NAME = 'certifications'),
  'ALTER TABLE career_paths DROP COLUMN certifications',
  'DO 0'
);
PREPARE stmt FROM @drop_cp_certifications;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- -----------------------------------------------------------------------------
-- 4. career_path_steps — langkah terpadu (EM + atribut CI opsional di step).
--    Kebijakan FK: career_path_id CASCADE; position_id TANPA FK (validasi
--    eksistensi di service layer — pola employee_movements.*_position_id).
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS career_path_steps (
    id                     CHAR(36) PRIMARY KEY,
    career_path_id         CHAR(36) NOT NULL,
    position_id            CHAR(36) NOT NULL,
    sequence               INT NOT NULL,
    -- EM: syarat per langkah.
    minimum_service_months INT NULL,
    requirements           TEXT NULL,
    -- CI: atribut edge (disimpan pada step target, sequence terakhir).
    path_type              VARCHAR(30) NULL,
    typical_tenure         INT NULL,
    competencies           TEXT NULL,
    certifications         TEXT NULL,
    created_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_career_path_steps_position (position_id),

    UNIQUE KEY uk_career_path_steps_sequence (career_path_id, sequence),
    UNIQUE KEY uk_career_path_steps_position (career_path_id, position_id),

    CONSTRAINT fk_career_path_steps_path FOREIGN KEY (career_path_id)
        REFERENCES career_paths(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
