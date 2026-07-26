-- Migration: 007_create_packages
-- Database: Platform (MySQL)
-- Tabel untuk menyimpan paket-modul (bundling modul tenant dengan harga)

CREATE TABLE IF NOT EXISTS packages (
    id          CHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NULL,
    price       DECIMAL(15,2) NOT NULL DEFAULT 0.00 COMMENT 'Harga paket dalam Rupiah',
    status      VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT 'draft | published | archived',
    is_public   TINYINT(1) NOT NULL DEFAULT 0 COMMENT '1 = tampil di halaman public',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,

    INDEX idx_packages_slug (slug),
    INDEX idx_packages_status (status),
    INDEX idx_packages_public (is_public),
    INDEX idx_packages_sort (sort_order),
    INDEX idx_packages_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
