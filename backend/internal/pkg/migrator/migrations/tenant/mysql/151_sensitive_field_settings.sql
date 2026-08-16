-- =============================================================================
-- 151_sensitive_field_settings.sql
-- Sensitive Data Masking — tabel setting toggle enkripsi per field.
-- Setiap baris mewakili satu field sensitif yang bisa di-enkripsi saat
-- ditulis (encrypt-on-write). Toggle ini independen dari permission
-- view per-field (lihat migration 153).
-- =============================================================================

CREATE TABLE sensitive_field_settings (
    id CHAR(36) PRIMARY KEY,
    field_key VARCHAR(100) NOT NULL,
    is_encryption_enabled TINYINT(1) NOT NULL DEFAULT 0,
    updated_by CHAR(36) NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_sensitive_field_settings_key (field_key)
);

INSERT IGNORE INTO sensitive_field_settings (id, field_key, is_encryption_enabled, updated_at) VALUES
    ('a3f1b2c4-0001-5f1a-9c1e-000000000001', 'employee.nik', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000002', 'employee.passport', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000003', 'employee.phone_number', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000004', 'employee.email', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000005', 'employee_family.nik', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000006', 'employee_bank_account.account_number', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000007', 'employee_bank_account.account_name', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000008', 'emergency_contact.phone_number', 0, CURRENT_TIMESTAMP);
