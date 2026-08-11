-- =============================================================================
-- Tenant Migration Down: 090_training_advanced (PostgreSQL)
-- =============================================================================

DROP INDEX IF EXISTS idx_trn_cert_certification;
ALTER TABLE training_certificates
    DROP COLUMN IF EXISTS certification_id,
    DROP COLUMN IF EXISTS certificate_file_url;

DROP TABLE IF EXISTS training_certifications;
DROP TABLE IF EXISTS training_effectiveness_assessments;
DROP TABLE IF EXISTS training_evaluation_answers;
DROP TABLE IF EXISTS training_evaluation_questions;
DROP TABLE IF EXISTS training_evaluation_forms;
