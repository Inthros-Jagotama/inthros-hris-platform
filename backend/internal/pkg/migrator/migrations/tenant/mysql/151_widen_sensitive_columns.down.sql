-- 151_widen_sensitive_columns.down.sql
-- NB: only safe to run down if no encrypted (longer) values were written
-- while the wider columns were in place — this truncates on rollback.

ALTER TABLE employees
    MODIFY COLUMN nik VARCHAR(16) NULL,
    MODIFY COLUMN passport VARCHAR(50) NULL;

ALTER TABLE employee_families
    MODIFY COLUMN nik VARCHAR(16) NULL;

ALTER TABLE employee_bank_accounts
    MODIFY COLUMN account_number VARCHAR(50) NOT NULL;

ALTER TABLE emergency_contacts
    MODIFY COLUMN phone_number VARCHAR(50) NOT NULL;
