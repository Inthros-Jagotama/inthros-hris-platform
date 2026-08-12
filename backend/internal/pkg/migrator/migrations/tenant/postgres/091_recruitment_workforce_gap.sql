-- =============================================================================
-- Tenant Migration: 091_recruitment_workforce_gap (PostgreSQL)
-- =============================================================================
-- S-1 Workforce Gap → Requisition
-- (docs/module-recruitment-strategic-layer-plan.md §4.1)
--
-- Recruitment membaca hiring need dari Workforce Intelligence melalui interface
-- narrow (wiring di cmd/server/main.go) — Recruitment TIDAK menghitung gap
-- sendiri. Kolom berikut hanya menyimpan referensi & alasan pembuatan
-- requisition dari sisi strategis:
--
--   reason_type        VARCHAR(30)  NULL  → alasan pembuatan requisition
--                        (NEW_POSITION | REPLACEMENT | EXPANSION | WORKFORCE_GAP)
--   workforce_gap_id   CHAR(36)     NULL  → referensi gap analysis Workforce
--                        Intelligence yang memicu requisition ini
--   workforce_plan_id  CHAR(36)     NULL  → referensi headcount plan (period)
--                        tempat gap teridentifikasi
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE job_requisitions
    ADD COLUMN IF NOT EXISTS reason_type VARCHAR(30) NULL,
    ADD COLUMN IF NOT EXISTS workforce_gap_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS workforce_plan_id CHAR(36) NULL;
