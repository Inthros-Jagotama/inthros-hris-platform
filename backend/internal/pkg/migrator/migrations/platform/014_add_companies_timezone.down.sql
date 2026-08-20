-- Down Migration: 014_add_companies_timezone
-- Database: Platform (Cross-Dialect)
--
-- Menghapus kolom timezone yang ditambahkan pada migration 014_add_companies_timezone.

ALTER TABLE companies
    DROP COLUMN timezone;
