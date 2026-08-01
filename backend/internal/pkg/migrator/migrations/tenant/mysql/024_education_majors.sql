-- =============================================================================
-- Tenant Migration: 024_education_majors
-- =============================================================================
-- Tabel master jurusan pendidikan (education major).
-- Sebelumnya dibuat via GORM AutoMigrate — kini diresmikan sebagai file DDL.
-- Struktur mengikuti model EducationMajor (setting module):
--   id, code (unique), name, sort_order, created_at, updated_at, deleted_at

CREATE TABLE IF NOT EXISTS education_majors (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(200) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    UNIQUE KEY idx_education_major_code (code),
    INDEX idx_education_majors_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
