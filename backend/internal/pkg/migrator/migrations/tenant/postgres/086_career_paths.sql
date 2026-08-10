-- =============================================================================
-- Tenant Migration: 086_career_paths
-- =============================================================================
-- Employee Movement Enhancement (plan §12.9): Career Path.
--
-- Career Path adalah planning/configuration jenjang karier (bukan movement
-- transaction). Satu career_paths berisi deretan career_path_steps yang
-- menunjuk ke posisi (organisasi = posisi), diurutkan by sequence, dengan
-- syarat opsional minimum masa kerja (minimum_service_months) dan
-- requirements textual.
--
-- Kebijakan FK: career_path_id memakai FK CASCADE ke career_paths (kedua
-- tabel lahir di migration yang sama). position_id TIDAK memakai FK karena
-- mengikuti pola employee_movements.from_/to_position_id — validasi eksistensi
-- posisi dilakukan di service layer.

CREATE TABLE IF NOT EXISTS career_paths (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_by  CHAR(36),
    updated_by  CHAR(36),
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_career_paths_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS career_path_steps (
    id                     CHAR(36) PRIMARY KEY,
    career_path_id         CHAR(36) NOT NULL,
    position_id            CHAR(36) NOT NULL,
    sequence               INT NOT NULL,
    minimum_service_months INT,
    requirements           TEXT,
    created_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_career_path_steps_sequence UNIQUE (career_path_id, sequence),
    CONSTRAINT uk_career_path_steps_position UNIQUE (career_path_id, position_id),

    CONSTRAINT fk_career_path_steps_path FOREIGN KEY (career_path_id)
        REFERENCES career_paths(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_career_path_steps_position ON career_path_steps (position_id);
