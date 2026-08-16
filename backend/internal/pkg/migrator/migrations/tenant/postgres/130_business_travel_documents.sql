-- 130_business_travel_documents.sql
-- Business Travel Module: dokumen perjalanan (travel order, tiket, boarding pass, dst — §23 plan doc)

CREATE TABLE IF NOT EXISTS business_travel_documents (
    id                   CHAR(36)     NOT NULL PRIMARY KEY,
    business_travel_id   CHAR(36)     NOT NULL,
    document_type        VARCHAR(30)  NOT NULL DEFAULT 'OTHER',
    file_name             VARCHAR(255) NOT NULL,
    file_path             TEXT         NOT NULL,
    mime_type             VARCHAR(100) NULL,
    file_size             BIGINT       NULL,
    uploaded_by           CHAR(36)     NULL,
    uploaded_at           TIMESTAMP    NULL,
    deleted_at            TIMESTAMP    NULL,
    created_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biztrav_doc_travel ON business_travel_documents (business_travel_id);

CREATE INDEX IF NOT EXISTS idx_biztrav_doc_deleted_at ON business_travel_documents (deleted_at);
