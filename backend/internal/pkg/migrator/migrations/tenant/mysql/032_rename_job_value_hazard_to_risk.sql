-- Migration: 032_rename_job_value_hazard_to_risk
-- Mengganti nilai type 'hazard' menjadi 'risk' pada tabel job_management_values
-- (konsistensi istilah: halaman Job Value Mapping kini menampilkan "Risk", dan
-- nilai type yang tersimpan mengikuti — filter ?type=risk di FE/API).
--
-- Idempotent: UPDATE dengan kondisi WHERE type = 'hazard' aman dijalankan ulang.

UPDATE job_management_values
SET type = 'risk'
WHERE type = 'hazard';
