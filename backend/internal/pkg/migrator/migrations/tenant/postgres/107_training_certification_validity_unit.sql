-- =============================================================================
-- Tenant Migration: 107_training_certification_validity_unit (PostgreSQL)
-- =============================================================================
-- Menambah kolom validity_period_unit pada training_certifications agar masa
-- berlaku dapat dinyatakan dalam tahun ('year') ATAU bulan ('month').
-- Nilai default 'month' menjaga kompatibilitas data lama.

ALTER TABLE training_certifications
    ADD COLUMN IF NOT EXISTS validity_period_unit VARCHAR(10) NOT NULL DEFAULT 'month';
