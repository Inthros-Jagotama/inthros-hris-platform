DROP INDEX idx_document_templates_type_movement_status ON document_templates;

ALTER TABLE document_templates
    DROP COLUMN movement_type;
