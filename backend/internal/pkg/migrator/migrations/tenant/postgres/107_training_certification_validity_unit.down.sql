-- =============================================================================
-- Tenant Migration Down: 107_training_certification_validity_unit (PostgreSQL)
-- =============================================================================

ALTER TABLE training_certifications
    DROP COLUMN IF EXISTS validity_period_unit;
