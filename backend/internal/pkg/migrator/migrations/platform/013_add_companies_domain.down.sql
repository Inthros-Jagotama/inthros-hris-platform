-- Down Migration: 013_add_companies_domain
-- Database: Platform (Cross-Dialect)

DROP INDEX idx_companies_subdomain ON companies;
DROP INDEX idx_companies_domain    ON companies;

ALTER TABLE companies
    DROP COLUMN subdomain,
    DROP COLUMN domain;
