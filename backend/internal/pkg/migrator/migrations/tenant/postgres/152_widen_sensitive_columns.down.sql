-- 152_widen_sensitive_columns.down.sql
-- NB: only safe to run down if no encrypted (longer) values were written
-- while the wider columns were in place — this truncates on rollback.

ALTER TABLE employees
    ALTER COLUMN nik TYPE VARCHAR(16),
    ALTER COLUMN passport TYPE VARCHAR(50);

ALTER TABLE employee_families
    ALTER COLUMN nik TYPE VARCHAR(16);

ALTER TABLE employee_bank_accounts
    ALTER COLUMN account_number TYPE VARCHAR(50);

ALTER TABLE emergency_contacts
    ALTER COLUMN phone_number TYPE VARCHAR(50);
