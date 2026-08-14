ALTER TABLE document_templates
    ADD CONSTRAINT fk_doctpl_active_version FOREIGN KEY (active_version_id) REFERENCES document_template_versions(id) ON DELETE SET NULL;
