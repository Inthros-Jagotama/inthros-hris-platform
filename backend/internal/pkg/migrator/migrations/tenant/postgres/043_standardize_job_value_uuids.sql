-- Migration: 043_standardize_job_value_uuids
-- Menstandarkan seluruh id di tabel job_management_values menjadi UUID v4 acak.
--
-- Latar belakang: migration seed 033-042 memakai UUID tetap/serial (mis.
-- aaaaaaaa-0001-4000-8000-000000000001) demi idempotensi WHERE NOT EXISTS.
-- User meminta seluruh baris memakai UUID acak yang proper. Migration ini
-- mengubah SEMUA id yang mengikuti pola tetap '-4000-8000-' menjadi UUID()
-- (Postgres) / gen_random_uuid() (Postgres) — berlaku untuk DB existing MAUPUN
-- fresh install (karena 043 berjalan setelah 042).
--
-- Idempotent: hanya baris dengan pola id tetap yang di-update; UUID v4 acak
-- (mis. 945453fe-92e8-44b0-a786-de540bb70fdc) tidak mengandung '-4000-8000-'
-- sehingga tidak tersentuh saat migration dijalankan ulang.

UPDATE job_management_values
SET id = gen_random_uuid()
WHERE id LIKE '%-4000-8000-%';
