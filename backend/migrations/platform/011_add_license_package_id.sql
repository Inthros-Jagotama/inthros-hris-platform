-- Migration: 011_add_license_package_id
-- Database: Platform (MySQL)
-- Menambahkan relasi license ke package (opsional)

ALTER TABLE licenses
    ADD COLUMN package_id   CHAR(36)     NULL COMMENT 'Optional: associated package from package management' AFTER status,
    ADD COLUMN package_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Snapshot of package name at time of license creation' AFTER package_id,
    ADD INDEX idx_licenses_package (package_id),
    ADD CONSTRAINT fk_licenses_package FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE SET NULL;
