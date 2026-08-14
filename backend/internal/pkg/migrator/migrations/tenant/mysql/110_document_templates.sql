ALTER TABLE document_templates
    ADD COLUMN code VARCHAR(100) NULL AFTER name,
    ADD COLUMN description TEXT NULL AFTER code,
    ADD COLUMN active_version_id CHAR(36) NULL AFTER content,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'INACTIVE' AFTER active_version_id,
    ADD COLUMN is_default TINYINT(1) NOT NULL DEFAULT 0 AFTER status,
    ADD COLUMN deleted_at TIMESTAMP NULL AFTER updated_at;

CREATE UNIQUE INDEX uk_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type_status ON document_templates (type, status);

CREATE TABLE IF NOT EXISTS document_template_versions (
    id            CHAR(36) PRIMARY KEY,
    template_id   CHAR(36) NOT NULL,
    version       INT NOT NULL,
    content       LONGTEXT NOT NULL,
    paper_size    VARCHAR(20) NOT NULL DEFAULT 'A4',
    orientation   VARCHAR(20) NOT NULL DEFAULT 'portrait',
    margin_top    INT NOT NULL DEFAULT 20,
    margin_right  INT NOT NULL DEFAULT 20,
    margin_bottom INT NOT NULL DEFAULT 20,
    margin_left   INT NOT NULL DEFAULT 20,
    created_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_doctplver_template FOREIGN KEY (template_id) REFERENCES document_templates(id) ON DELETE CASCADE,
    CONSTRAINT uq_doctplver_template_version UNIQUE (template_id, version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS document_template_audits (
    id          CHAR(36) PRIMARY KEY,
    template_id CHAR(36) NOT NULL,
    version_id  CHAR(36) NULL,
    action      VARCHAR(50) NOT NULL,
    actor_id    CHAR(36) NULL,
    payload     JSON NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_doctplaudit_template FOREIGN KEY (template_id) REFERENCES document_templates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_doctplaudit_template ON document_template_audits (template_id);

CREATE TABLE IF NOT EXISTS generated_documents (
    id                   CHAR(36) PRIMARY KEY,
    template_id          CHAR(36) NOT NULL,
    template_version_id  CHAR(36) NOT NULL,
    document_type        VARCHAR(50) NOT NULL,
    reference_type       VARCHAR(50) NOT NULL,
    reference_id         CHAR(36) NOT NULL,
    file_name            VARCHAR(255) NOT NULL,
    file_path            VARCHAR(500) NOT NULL,
    mime_type            VARCHAR(100) NOT NULL DEFAULT 'application/pdf',
    generated_by         CHAR(36) NULL,
    generated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_gendoc_template FOREIGN KEY (template_id) REFERENCES document_templates(id),
    CONSTRAINT fk_gendoc_version FOREIGN KEY (template_version_id) REFERENCES document_template_versions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_gendoc_reference ON generated_documents (reference_type, reference_id);

-- Seed: one REFERENCE (default) template per document type. Reference
-- templates store content directly and never use active_version_id.
INSERT INTO document_templates (id, name, code, type, description, content, active_version_id, status, is_default, is_active, created_at, updated_at)
VALUES
    (UUID(), 'Perjanjian Kerja Waktu Tertentu (Default)', 'CONTRACT_AGREEMENT_DEFAULT', 'CONTRACT_AGREEMENT', 'Template referensi bawaan untuk Perjanjian Kontrak', '<h2>PERJANJIAN KERJA WAKTU TERTENTU</h2><p>Nomor: {{contract.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan: {{employee.position}}</p>', NULL, 'REFERENCE', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (UUID(), 'SK Movement (Default)', 'MOVEMENT_SK_DEFAULT', 'MOVEMENT_SK', 'Template referensi bawaan untuk SK Movement', '<h2>SURAT KEPUTUSAN MUTASI/PROMOSI/DEMOSI</h2><p>Nomor: {{movement.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan Baru: {{movement.new_position}}</p>', NULL, 'REFERENCE', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE id = id;
