-- =============================================================================
-- 142_competency_score_employee_unique.sql
-- Competency 360 Module — Score per Employee per Event (plan generik §15).
--
-- `competency_scores` saat ini UNIQUE (organization_id) — satu skor per
-- posisi. Untuk 360, skor adalah hasil per EMPLOYEE per assessment event:
-- satu posisi dapat dihuni banyak employee dan setiap event punya skor
-- sendiri. Unique menjadi (competency_event_id, employee_id).
--
-- ⚠️ BREAKING pada constraint lama: sebelum dieksekusi di production, wajib
-- cek data existing yang akan melanggar (duplikat event+employee). §34.1:
-- tabel ini belum punya jalur pengisian aktif — kemungkinan besar aman.
-- =============================================================================

ALTER TABLE competency_scores
    DROP CONSTRAINT IF EXISTS uk_comp_score_org;

-- (competency_event_id, employee_id) keduanya nullable di skema lama —
-- postgres memperlakukan NULL sebagai distinct sehingga baris lama tanpa
-- event/employee tetap bisa hidup berdampingan.
CREATE UNIQUE INDEX IF NOT EXISTS uk_comp_score_event_employee
    ON competency_scores (competency_event_id, employee_id);

-- Index untuk lookup per employee tetap dipertahankan dari migration 008
-- (idx_comp_score_employee), jadi tidak perlu dibuat ulang di sini.
