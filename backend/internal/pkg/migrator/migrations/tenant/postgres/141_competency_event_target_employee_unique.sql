-- =============================================================================
-- 141_competency_event_target_employee_unique.sql
-- Competency 360 Module — Employee-centric Assessment Subject (plan generik §8).
--
-- `competency_event_targets` saat ini UNIQUE (competency_event_id,
-- organization_id) — satu target per posisi. Untuk 360, subject assessment
-- adalah EMPLOYEE: satu posisi dapat dihuni banyak employee, jadi unique
-- menjadi (competency_event_id, employee_id). organization_id dipertahankan
-- sebagai snapshot/reference posisi saat assessment.
--
-- ⚠️ BREAKING pada constraint lama: sebelum dieksekusi di production, wajib
-- cek data existing yang akan melanggar (duplikat employee per event).
-- §34.1: tabel ini belum punya jalur pengisian aktif di kode — kemungkinan
-- besar aman, tapi diverifikasi saat eksekusi (di luar scope implementasi ini).
-- =============================================================================

-- Drop unique constraint lama (event, organization).
ALTER TABLE competency_event_targets
    DROP CONSTRAINT IF EXISTS uk_comp_event_target;

-- Tambahkan unique baru (event, employee). employee_id tetap nullable untuk
-- backward compatibility data lama yang belum mengisi employee; unique index
-- postgres memperlakukan NULL sebagai distinct sehingga tidak menghalangi.
CREATE UNIQUE INDEX IF NOT EXISTS uk_comp_event_target_employee
    ON competency_event_targets (competency_event_id, employee_id);
