ALTER TABLE document_templates
    ADD COLUMN is_default TINYINT(1) NOT NULL DEFAULT 0 AFTER status;
