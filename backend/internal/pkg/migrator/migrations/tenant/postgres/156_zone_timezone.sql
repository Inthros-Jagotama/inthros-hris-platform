-- =============================================================================
-- Tenant Migration: 156_zone_timezone (PostgreSQL)
-- =============================================================================
-- Menambahkan kolom timezone pada tabel zones untuk menyimpan override timezone
-- per zona geografis. Kolom ini nullable karena zone dapat menggunakan timezone
-- default dari company jika tidak di-specify.
--
-- Timezone value berupa IANA timezone identifier (contoh: 'Asia/Jakarta',
-- 'America/New_York', dll). Digunakan untuk context khusus per lokasi/zone
-- dalam operasi datetime dan shift calculation.

ALTER TABLE zones
    ADD COLUMN timezone VARCHAR(64) NULL;
