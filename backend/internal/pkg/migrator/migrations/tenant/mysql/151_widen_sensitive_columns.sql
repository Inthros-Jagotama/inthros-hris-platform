-- =============================================================================
-- 151_widen_sensitive_columns.sql
-- Sensitive Data Masking — perbesar kolom yang berpotensi diisi ciphertext
-- (AES-256-GCM hex-encoded lebih panjang dari plaintext aslinya).
-- Dijalankan terlepas dari status toggle enkripsi saat ini, supaya
-- mengaktifkan enkripsi nanti tidak perlu migrasi skema tambahan.
-- =============================================================================

ALTER TABLE employees
    MODIFY COLUMN nik VARCHAR(255) NULL,
    MODIFY COLUMN passport VARCHAR(255) NULL;
-- phone_number and email already varchar(255), no change needed.

ALTER TABLE employee_families
    MODIFY COLUMN nik VARCHAR(255) NULL;

ALTER TABLE employee_bank_accounts
    MODIFY COLUMN account_number VARCHAR(255) NOT NULL;
-- account_name already varchar(255), no change needed.

-- emergency_contacts.phone_number already varchar(50) -> widen for ciphertext.
ALTER TABLE emergency_contacts
    MODIFY COLUMN phone_number VARCHAR(255) NOT NULL;
