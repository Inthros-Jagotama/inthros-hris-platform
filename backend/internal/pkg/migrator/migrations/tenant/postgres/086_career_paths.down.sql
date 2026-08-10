-- =============================================================================
-- Tenant Migration Down: 086_career_paths (unified schema rollback, PostgreSQL)
-- =============================================================================
-- Mengembalikan career_paths ke bentuk edge career intelligence (018).

DROP TABLE IF EXISTS career_path_steps;

DROP INDEX IF EXISTS uk_career_paths_name;
ALTER TABLE career_paths DROP COLUMN IF EXISTS name;
ALTER TABLE career_paths DROP COLUMN IF EXISTS description;
ALTER TABLE career_paths DROP COLUMN IF EXISTS created_by;
ALTER TABLE career_paths DROP COLUMN IF EXISTS updated_by;

-- Kembalikan kolom edge CI.
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS source_title_id CHAR(36) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS target_title_id CHAR(36) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS path_type VARCHAR(30) NOT NULL DEFAULT 'PROMOTION';
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS typical_tenure INT NOT NULL DEFAULT 0;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS requirements TEXT NULL;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS competencies TEXT NULL;
ALTER TABLE career_paths ADD COLUMN IF NOT EXISTS certifications TEXT NULL;

CREATE INDEX IF NOT EXISTS idx_cp_source ON career_paths (source_title_id);
CREATE INDEX IF NOT EXISTS idx_cp_target ON career_paths (target_title_id);
