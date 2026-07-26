-- Migration: 010_add_modules_depends_on
-- Database: Platform (MySQL)
-- Menambahkan kolom depends_on untuk menyimpan dependensi antar modul

ALTER TABLE modules
    ADD COLUMN depends_on TEXT NULL COMMENT 'Comma-separated list of module slugs that this module depends on' AFTER is_core;
