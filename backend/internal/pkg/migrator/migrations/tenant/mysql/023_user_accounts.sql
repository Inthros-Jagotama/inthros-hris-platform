-- =============================================================================
-- Tenant Migration: 023_user_accounts
-- =============================================================================
-- Akun login untuk employee (tabel users dibuat di 022_users.sql).
-- employee_accounts menghubungkan employee ↔ user tenant dan menyimpan
-- setup token untuk alur "set password via link email" (sekali pakai).

CREATE TABLE IF NOT EXISTS employee_accounts (
    id                  CHAR(36) PRIMARY KEY,
    company_id          CHAR(36) NOT NULL,
    employee_id         CHAR(36) NOT NULL,
    user_id             CHAR(36) NOT NULL,
    email               VARCHAR(255) NOT NULL,
    setup_token         VARCHAR(255) NULL,
    setup_token_expires TIMESTAMP NULL,
    created_by          CHAR(36) NULL,
    updated_by          CHAR(36) NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_employee_accounts_employee (employee_id),
    UNIQUE KEY uk_employee_accounts_email (email),
    INDEX idx_employee_accounts_setup_token (setup_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
