-- Migration: 014_add_companies_timezone
-- Database: Platform (Cross-Dialect)
-- Menambahkan kolom timezone untuk menyimpan timezone perusahaan.
-- Timezone default adalah 'Asia/Jakarta' dan diperlukan untuk operasi datetime
-- yang melibatkan timezone awareness di aplikasi.

ALTER TABLE companies
    ADD COLUMN timezone VARCHAR(64) NOT NULL DEFAULT 'Asia/Jakarta' AFTER phone;
