-- =============================================================================
-- Tenant Migration: 021_insurances
-- =============================================================================
-- Tabel master asuransi (BPJS Kesehatan, BPJS Ketenagakerjaan, dll).
-- Sebelumnya dibuat via GORM AutoMigrate — kini diresmikan sebagai file DDL.
-- Struktur mengikuti model Insurance (setting module):
--   id, code (unique), name, sort_order, created_at, updated_at, deleted_at

CREATE TABLE IF NOT EXISTS insurances (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    UNIQUE KEY idx_insurance_code (code),
    INDEX idx_insurances_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
