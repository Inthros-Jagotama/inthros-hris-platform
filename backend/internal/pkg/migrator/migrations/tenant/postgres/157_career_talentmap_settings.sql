-- =============================================================================
-- Tenant Migration: 157_career_talentmap_settings (PostgreSQL)
-- =============================================================================
-- Ambang batas skor (skala 0-100) untuk membanding otomatis skor Performance
-- Management (performance) dan Competency (potential) menjadi LOW/MEDIUM/HIGH
-- saat generate Talent Map (9-box grid). Singleton table (satu baris per
-- tenant, dibuat otomatis oleh service saat pertama kali dibaca — pola sama
-- seperti employee_id_format_settings / attendance_company_settings).
--
-- Band: score < low_max => LOW; low_max <= score < high_min => MEDIUM;
-- score >= high_min => HIGH. Validasi low_max < high_min dilakukan di service.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS career_talentmap_settings (
    id                    CHAR(36) PRIMARY KEY,
    performance_low_max   DECIMAL(5,2) NOT NULL DEFAULT 50.00,
    performance_high_min  DECIMAL(5,2) NOT NULL DEFAULT 80.00,
    potential_low_max     DECIMAL(5,2) NOT NULL DEFAULT 50.00,
    potential_high_min    DECIMAL(5,2) NOT NULL DEFAULT 80.00,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
