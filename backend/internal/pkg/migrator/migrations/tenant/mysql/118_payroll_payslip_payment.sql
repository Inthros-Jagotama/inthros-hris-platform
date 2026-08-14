-- =============================================================================
-- Tenant Migration: 118_payroll_payslip_payment
-- =============================================================================
-- Payment batch (docs/payroll/07 §27) + kolom employer contribution di
-- payslip (bagian "Employer Contribution" pada payslip, §28).

-- ---------------------------------------------------------------------------
-- 118.1 payroll_payslips: tambah total employer contribution
-- ---------------------------------------------------------------------------
ALTER TABLE payroll_payslips
    ADD COLUMN total_employer_contribution DECIMAL(18, 2) NOT NULL DEFAULT 0 AFTER net_amount;

-- ---------------------------------------------------------------------------
-- 118.2 payroll_payments
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS payroll_payments (
    id                          CHAR(36) PRIMARY KEY,
    payroll_run_id              CHAR(36) NOT NULL,
    payroll_run_employee_id     CHAR(36) NOT NULL,
    employee_id                 CHAR(36) NOT NULL,
    employee_code               VARCHAR(50) NOT NULL,
    employee_name               VARCHAR(255) NOT NULL,
    amount                      DECIMAL(18, 2) NOT NULL DEFAULT 0,
    currency_code               CHAR(3) NOT NULL DEFAULT 'IDR',
    payment_date                DATE NOT NULL,
    employee_bank_profile_id    CHAR(36) NULL,
    bank_code                   VARCHAR(50) NULL,
    bank_name                   VARCHAR(150) NULL,
    bank_branch                 VARCHAR(150) NULL,
    bank_account_number         VARCHAR(100) NOT NULL,
    bank_account_holder_name    VARCHAR(255) NOT NULL,
    status                      VARCHAR(255) NOT NULL DEFAULT 'PENDING',
    reference                   VARCHAR(100) NULL,
    processed_at                TIMESTAMP NULL,
    paid_at                     TIMESTAMP NULL,
    failed_at                   TIMESTAMP NULL,
    failed_reason               VARCHAR(255) NULL,
    reversed_at                 TIMESTAMP NULL,
    created_by                  CHAR(36) NULL,
    updated_by                  CHAR(36) NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT uk_payment_run_employee UNIQUE (payroll_run_id, payroll_run_employee_id),

    CONSTRAINT fk_payment_run_employee FOREIGN KEY (payroll_run_employee_id) REFERENCES payroll_run_employees(id) ON DELETE CASCADE
);

CREATE INDEX idx_payment_employee ON payroll_payments (employee_id);
CREATE INDEX idx_payment_status ON payroll_payments (status);
