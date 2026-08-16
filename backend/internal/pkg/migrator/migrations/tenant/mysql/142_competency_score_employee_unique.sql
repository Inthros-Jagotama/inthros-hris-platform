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

-- Index uk_comp_score_org dipakai FK fk_comp_score_org (MySQL menolak drop
-- index yang masih dipakai constraint). Drop FK dulu, ganti index, re-add FK.
ALTER TABLE competency_scores
    DROP FOREIGN KEY fk_comp_score_org,
    DROP INDEX uk_comp_score_org,
    ADD UNIQUE KEY uk_comp_score_event_employee (competency_event_id, employee_id),
    ADD INDEX idx_comp_score_org (organization_id);

ALTER TABLE competency_scores
    ADD CONSTRAINT fk_comp_score_org FOREIGN KEY (organization_id)
        REFERENCES organizations(id) ON DELETE CASCADE;
