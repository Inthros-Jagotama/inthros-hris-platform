-- =============================================================================
-- 152_widen_sensitive_columns.sql
-- Sensitive Data Masking — perbesar kolom yang berpotensi diisi ciphertext
-- (AES-256-GCM hex-encoded lebih panjang dari plaintext aslinya).
-- Dijalankan terlepas dari status toggle enkripsi saat ini, supaya
-- mengaktifkan enkripsi nanti tidak perlu migrasi skema tambahan.
-- =============================================================================

ALTER TABLE employees
    ALTER COLUMN nik TYPE VARCHAR(255),
    ALTER COLUMN passport TYPE VARCHAR(255);

ALTER TABLE employee_families
    ALTER COLUMN nik TYPE VARCHAR(255);

ALTER TABLE employee_bank_accounts
    ALTER COLUMN account_number TYPE VARCHAR(255);

ALTER TABLE emergency_contacts
    ALTER COLUMN phone_number TYPE VARCHAR(255);
