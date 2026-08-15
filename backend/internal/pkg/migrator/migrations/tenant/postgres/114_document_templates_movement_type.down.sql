DROP INDEX IF EXISTS idx_document_templates_type_movement_status;

DROP INDEX IF EXISTS uq_document_templates_active_per_type;
CREATE UNIQUE INDEX uq_document_templates_active_per_type
    ON document_templates (type)
    WHERE status = 'ACTIVE';

ALTER TABLE document_templates
    DROP COLUMN movement_type;
