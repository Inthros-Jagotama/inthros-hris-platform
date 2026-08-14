DROP TABLE IF EXISTS payroll_payments;

ALTER TABLE payroll_payslips
    DROP COLUMN total_employer_contribution;
