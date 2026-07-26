-- Migration: 008_create_package_modules
-- Database: Platform (MySQL)
-- Tabel relasi many-to-many antara package dan module

CREATE TABLE IF NOT EXISTS package_modules (
    package_id    CHAR(36) NOT NULL,
    module_id     CHAR(36) NOT NULL,
    is_mandatory  TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Modul wajib yang tidak bisa dinonaktifkan',
    sort_order    INT NOT NULL DEFAULT 0,
    module_name   VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Nama modul saat dimasukkan ke paket (snapshot)',
    module_slug   VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'Slug modul saat dimasukkan ke paket (snapshot)',
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (package_id, module_id),

    INDEX idx_pm_package (package_id),
    INDEX idx_pm_module (module_id),
    INDEX idx_pm_mandatory (is_mandatory),

    CONSTRAINT fk_pm_package FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE,
    CONSTRAINT fk_pm_module  FOREIGN KEY (module_id)  REFERENCES modules(id)  ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
