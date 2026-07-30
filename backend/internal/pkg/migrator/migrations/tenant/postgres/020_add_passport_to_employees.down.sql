-- Down Migration: 020_add_passport_to_employees
ALTER TABLE employees
    DROP COLUMN passport;
