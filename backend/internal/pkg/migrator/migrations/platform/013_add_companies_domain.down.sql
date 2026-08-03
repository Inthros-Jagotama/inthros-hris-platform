-- Down Migration: 013_add_companies_domain
-- Database: Platform (Cross-Dialect)
--
-- CATATAN: index idx_companies_subdomain & idx_companies_domain (single-column)
-- otomatis ikut ter-drop bersama kolomnya di kedua dialect (MySQL & PostgreSQL),
-- jadi tidak perlu DROP INDEX eksplisit (yang sintaksnya beda per dialect).

ALTER TABLE companies
    DROP COLUMN subdomain,
    DROP COLUMN domain;
