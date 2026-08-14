DROP TABLE IF EXISTS generated_documents;
DROP TABLE IF EXISTS document_template_audits;
DROP TABLE IF EXISTS document_template_versions;

DELETE FROM document_templates WHERE code IN ('CONTRACT_AGREEMENT_DEFAULT', 'MOVEMENT_SK_DEFAULT');

DROP INDEX IF EXISTS uq_document_templates_default_per_type;
DROP INDEX IF EXISTS uq_document_templates_active_per_type;
DROP INDEX IF EXISTS idx_document_templates_type_status;
DROP INDEX IF EXISTS uk_document_templates_code;

ALTER TABLE document_templates
    DROP COLUMN deleted_at,
    DROP COLUMN is_default,
    DROP COLUMN status,
    DROP COLUMN active_version_id,
    DROP COLUMN description,
    DROP COLUMN code;
