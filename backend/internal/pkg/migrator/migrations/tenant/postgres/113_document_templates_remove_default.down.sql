ALTER TABLE document_templates
    ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX uq_document_templates_default_per_type
    ON document_templates (type)
    WHERE is_default = TRUE;
