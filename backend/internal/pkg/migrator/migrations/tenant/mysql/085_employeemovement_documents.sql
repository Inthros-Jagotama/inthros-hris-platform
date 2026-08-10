-- =============================================================================
-- Tenant Migration: 085_employeemovement_documents
-- =============================================================================
-- Employee Movement Enhancement (plan §12.15): Movement Documents.
--
-- Mendukung lebih dari satu dokumen per movement (selain decision letter
-- fields di employee_movements). File fisik di-upload lewat endpoint upload
-- generik (POST /api/v1/tenant/uploads), tabel ini hanya menyimpan metadata:
-- document_type (PROMOTION_SK, MUTATION_SK, DEMOTION_SK, RETIREMENT_LETTER,
-- OFFBOARDING_LETTER, OTHER), file_name, dan file_url publik.
--
-- uploaded_by menyimpan id user (CHAR(36)) yang meng-upload dokumen
-- (bisa NULL untuk data lama); created_at default ke waktu insert.

CREATE TABLE IF NOT EXISTS employee_movement_documents (
    id            CHAR(36) PRIMARY KEY,
    movement_id   CHAR(36) NOT NULL,
    document_type VARCHAR(30) NOT NULL COMMENT 'PROMOTION_SK, MUTATION_SK, DEMOTION_SK, RETIREMENT_LETTER, OFFBOARDING_LETTER, OTHER',
    file_name     VARCHAR(255) NOT NULL,
    file_url      VARCHAR(255) NOT NULL,
    uploaded_by   CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_emp_mvmt_doc_movement (movement_id),

    CONSTRAINT fk_empmvmt_doc_movement FOREIGN KEY (movement_id)
        REFERENCES employee_movements(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
