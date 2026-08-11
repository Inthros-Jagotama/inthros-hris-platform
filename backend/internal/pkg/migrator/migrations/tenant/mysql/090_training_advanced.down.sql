-- =============================================================================
-- Tenant Migration Down: 090_training_advanced (MySQL)
-- =============================================================================

DROP INDEX idx_trn_cert_certification ON training_certificates;

SET @db = DATABASE();

SET @drop_certification_id = (
    SELECT IF(
        COUNT(*) > 0,
        'ALTER TABLE training_certificates DROP COLUMN certification_id',
        'SELECT 1'
    )
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'training_certificates'
      AND COLUMN_NAME = 'certification_id'
);
PREPARE stmt FROM @drop_certification_id; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @drop_cert_file_url = (
    SELECT IF(
        COUNT(*) > 0,
        'ALTER TABLE training_certificates DROP COLUMN certificate_file_url',
        'SELECT 1'
    )
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'training_certificates'
      AND COLUMN_NAME = 'certificate_file_url'
);
PREPARE stmt FROM @drop_cert_file_url; EXECUTE stmt; DEALLOCATE PREPARE stmt;

DROP TABLE IF EXISTS training_certifications;
DROP TABLE IF EXISTS training_effectiveness_assessments;
DROP TABLE IF EXISTS training_evaluation_answers;
DROP TABLE IF EXISTS training_evaluation_questions;
DROP TABLE IF EXISTS training_evaluation_forms;
