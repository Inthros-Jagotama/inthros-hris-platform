-- =============================================================================
-- Tenant Migration: 109_document_numbering_settings (PostgreSQL)
-- =============================================================================
-- Pengaturan penomoran dokumen otomatis (SK/surat keputusan) untuk mutasi
-- karyawan (employee_movement) dan kontrak karyawan (employee_contract).
-- format_template mendukung placeholder seperti {sequence:3}, {month_roman},
-- {year}. last_sequence & last_reset_key dipakai untuk melacak nomor urut
-- terakhir per periode reset (reset_period).
-- Idempotent: CREATE TABLE IF NOT EXISTS.
--
-- See mysql version for full column documentation.

CREATE TABLE IF NOT EXISTS document_numbering_settings (
    id              CHAR(36) PRIMARY KEY,
    document_type   VARCHAR(50) NOT NULL,
    format_template VARCHAR(255) NOT NULL,
    reset_period    VARCHAR(20) NOT NULL DEFAULT 'yearly',
    last_sequence   INT NOT NULL DEFAULT 0,
    last_reset_key  VARCHAR(16) NOT NULL DEFAULT '',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_doc_numbering_type UNIQUE (document_type)
);

INSERT INTO document_numbering_settings (id, document_type, format_template, reset_period, last_sequence, last_reset_key)
VALUES
    (gen_random_uuid(), 'employee_movement', 'SK/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, ''),
    (gen_random_uuid(), 'employee_contract', 'CTR/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, '');
