-- =============================================================================
-- 151_sensitive_field_settings.sql (postgres)
-- Sensitive Data Masking — tabel setting toggle enkripsi per field.
-- =============================================================================

CREATE TABLE sensitive_field_settings (
    id CHAR(36) PRIMARY KEY,
    field_key VARCHAR(100) NOT NULL UNIQUE,
    is_encryption_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_by CHAR(36) NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO sensitive_field_settings (id, field_key, is_encryption_enabled, updated_at) VALUES
    ('a3f1b2c4-0001-5f1a-9c1e-000000000001', 'employee.nik', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000002', 'employee.passport', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000003', 'employee.phone_number', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000004', 'employee.email', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000005', 'employee_family.nik', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000006', 'employee_bank_account.account_number', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000007', 'employee_bank_account.account_name', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000008', 'emergency_contact.phone_number', FALSE, CURRENT_TIMESTAMP)
ON CONFLICT (field_key) DO NOTHING;
