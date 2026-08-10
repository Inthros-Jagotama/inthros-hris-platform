-- =============================================================================
-- Tenant Migration: 086_career_paths — UNIFIED CAREER PATH SCHEMA (PostgreSQL)
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
-- career_paths SUDAH ada (migration 018) → 086 ALTER (tambah kolom EM, hapus
-- kolom edge CI lama) + CREATE career_path_steps. Semua idempotent.

-- -----------------------------------------------------------------------------
-- 1. Tambah kolom header EM.
-- -----------------------------------------------------------------------------
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS name VARCHAR(100) NULL;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS description TEXT NULL;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS created_by CHAR(36) NULL;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS updated_by CHAR(36) NULL;

-- -----------------------------------------------------------------------------
-- 2. Unique index name.
-- -----------------------------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS uk_career_paths_name ON career_paths (name);

-- -----------------------------------------------------------------------------
-- 3. Hapus kolom edge CI lama (dipindah ke career_path_steps).
-- -----------------------------------------------------------------------------
DROP INDEX IF EXISTS idx_cp_source;
DROP INDEX IF EXISTS idx_cp_target;
ALTER TABLE career_paths DROP COLUMN IF EXISTS source_title_id;
ALTER TABLE career_paths DROP COLUMN IF EXISTS target_title_id;
ALTER TABLE career_paths DROP COLUMN IF EXISTS path_type;
ALTER TABLE career_paths DROP COLUMN IF EXISTS typical_tenure;
ALTER TABLE career_paths DROP COLUMN IF EXISTS requirements;
ALTER TABLE career_paths DROP COLUMN IF EXISTS competencies;
ALTER TABLE career_paths DROP COLUMN IF EXISTS certifications;

-- -----------------------------------------------------------------------------
-- 4. career_path_steps — langkah terpadu (EM + atribut CI opsional di step).
--    position_id TANPA FK (validasi eksistensi di service layer).
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS career_path_steps (
    id                     CHAR(36) PRIMARY KEY,
    career_path_id         CHAR(36) NOT NULL,
    position_id            CHAR(36) NOT NULL,
    sequence               INT NOT NULL,
    minimum_service_months INT NULL,
    requirements           TEXT NULL,
    path_type              VARCHAR(30) NULL,
    typical_tenure         INT NULL,
    competencies           TEXT NULL,
    certifications         TEXT NULL,
    created_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_career_path_steps_path FOREIGN KEY (career_path_id)
        REFERENCES career_paths(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_career_path_steps_position ON career_path_steps (position_id);

CREATE UNIQUE INDEX IF NOT EXISTS uk_career_path_steps_sequence
    ON career_path_steps (career_path_id, sequence);

CREATE UNIQUE INDEX IF NOT EXISTS uk_career_path_steps_position
    ON career_path_steps (career_path_id, position_id);
