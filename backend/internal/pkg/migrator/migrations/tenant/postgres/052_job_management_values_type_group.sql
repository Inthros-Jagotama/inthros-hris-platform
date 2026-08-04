-- Migration: 052_job_management_values_type_group
-- Menambah kolom type_group & description_group pada job_management_values
-- (pengelompokan tipe nilai jabatan beserta deskripsi grupnya — dipakai
-- frontend untuk menampilkan card/group tetap).
-- Catatan: Postgres tidak mendukung AFTER — kolom ditambahkan di akhir.

ALTER TABLE job_management_values
    ADD COLUMN type_group VARCHAR(255),
    ADD COLUMN description_group TEXT;
