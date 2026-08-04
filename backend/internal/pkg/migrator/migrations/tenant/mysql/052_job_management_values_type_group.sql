-- Migration: 052_job_management_values_type_group
-- Menambah kolom type_group & description_group pada job_management_values
-- (pengelompokan tipe nilai jabatan beserta deskripsi grupnya — dipakai
-- frontend untuk menampilkan card/group tetap).

ALTER TABLE job_management_values
    ADD COLUMN type_group VARCHAR(255) NULL AFTER type,
    ADD COLUMN description_group TEXT NULL AFTER descriptions;
