ALTER TABLE document_templates
    ADD COLUMN code VARCHAR(100) NULL,
    ADD COLUMN description TEXT NULL,
    ADD COLUMN active_version_id CHAR(36) NULL,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'INACTIVE',
    ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE UNIQUE INDEX uk_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type_status ON document_templates (type, status);

-- Only one ACTIVE template per document type (partial unique index, Postgres only)
CREATE UNIQUE INDEX uq_document_templates_active_per_type
    ON document_templates (type)
    WHERE status = 'ACTIVE';

-- Only one default (reference) template per document type
CREATE UNIQUE INDEX uq_document_templates_default_per_type
    ON document_templates (type)
    WHERE is_default = TRUE;

CREATE TABLE IF NOT EXISTS document_template_versions (
    id            CHAR(36) PRIMARY KEY,
    template_id   CHAR(36) NOT NULL REFERENCES document_templates(id) ON DELETE CASCADE,
    version       INT NOT NULL,
    content       TEXT NOT NULL,
    paper_size    VARCHAR(20) NOT NULL DEFAULT 'A4',
    orientation   VARCHAR(20) NOT NULL DEFAULT 'portrait',
    margin_top    INT NOT NULL DEFAULT 20,
    margin_right  INT NOT NULL DEFAULT 20,
    margin_bottom INT NOT NULL DEFAULT 20,
    margin_left   INT NOT NULL DEFAULT 20,
    created_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_doctplver_template_version UNIQUE (template_id, version)
);

CREATE TABLE IF NOT EXISTS document_template_audits (
    id          CHAR(36) PRIMARY KEY,
    template_id CHAR(36) NOT NULL REFERENCES document_templates(id) ON DELETE CASCADE,
    version_id  CHAR(36) NULL,
    action      VARCHAR(50) NOT NULL,
    actor_id    CHAR(36) NULL,
    payload     JSONB NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_doctplaudit_template ON document_template_audits (template_id);

CREATE TABLE IF NOT EXISTS generated_documents (
    id                   CHAR(36) PRIMARY KEY,
    template_id          CHAR(36) NOT NULL REFERENCES document_templates(id),
    template_version_id  CHAR(36) NOT NULL REFERENCES document_template_versions(id),
    document_type        VARCHAR(50) NOT NULL,
    reference_type       VARCHAR(50) NOT NULL,
    reference_id         CHAR(36) NOT NULL,
    file_name            VARCHAR(255) NOT NULL,
    file_path            VARCHAR(500) NOT NULL,
    mime_type            VARCHAR(100) NOT NULL DEFAULT 'application/pdf',
    generated_by         CHAR(36) NULL,
    generated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gendoc_reference ON generated_documents (reference_type, reference_id);

INSERT INTO document_templates (id, name, code, type, description, content, active_version_id, status, is_default, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'Perjanjian Kerja Waktu Tertentu (Default)', 'CONTRACT_AGREEMENT_DEFAULT', 'CONTRACT_AGREEMENT', 'Template referensi bawaan untuk Perjanjian Kontrak', '<h2>PERJANJIAN KERJA WAKTU TERTENTU</h2><p>Nomor: {{contract.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan: {{employee.position}}</p>', NULL, 'REFERENCE', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'SK Movement (Default)', 'MOVEMENT_SK_DEFAULT', 'MOVEMENT_SK', 'Template referensi bawaan untuk SK Movement', '<h2>SURAT KEPUTUSAN MUTASI/PROMOSI/DEMOSI</h2><p>Nomor: {{movement.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan Baru: {{movement.new_position}}</p>', NULL, 'REFERENCE', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;
