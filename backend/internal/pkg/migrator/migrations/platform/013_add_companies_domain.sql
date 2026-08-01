-- Migration: 013_add_companies_domain
-- Database: Platform (Cross-Dialect)
-- Menambahkan kolom domain & subdomain untuk akses aplikasi SaaS per-company.
-- Setiap company punya alamat URL aplikasi sendiri (mis. <subdomain>.app.hris.com).
-- Digunakan tenant FE untuk resolve company dari hostname tanpa subdomain eksplisit.

ALTER TABLE companies
    ADD COLUMN subdomain VARCHAR(100) NULL,
    ADD COLUMN domain    VARCHAR(255) NULL;

CREATE UNIQUE INDEX idx_companies_subdomain ON companies (subdomain);
CREATE UNIQUE INDEX idx_companies_domain    ON companies (domain);
