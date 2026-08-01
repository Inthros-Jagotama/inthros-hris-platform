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
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP,

    CONSTRAINT idx_insurance_code UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_insurances_deleted_at ON insurances (deleted_at);
