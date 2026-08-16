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
    uploaded_at           TIMESTAMP(6) NULL DEFAULT NULL,
    deleted_at            TIMESTAMP(6) NULL DEFAULT NULL,
    created_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at            TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_biztrav_doc_travel (business_travel_id),
    INDEX idx_biztrav_doc_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
