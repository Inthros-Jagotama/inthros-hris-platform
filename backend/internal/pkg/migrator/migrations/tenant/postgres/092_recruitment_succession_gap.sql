-- =============================================================================
-- Tenant Migration: 092_recruitment_succession_gap (PostgreSQL)
-- =============================================================================
-- S-5 Succession Planning → Fallback External Recruitment
-- (docs/module-recruitment-strategic-layer-plan.md §4.5)
--
-- CI menandai posisi kunci tanpa successor siap (succession gap); Recruitment
-- menjadi fallback external: requisition dengan reason_type=SUCCESSION_GAP
-- menautkan ke posisi kunci tsb via succession_position_id. Recruitment TIDAK
-- menghitung readiness succession sendiri — ia membaca hasil CI melalui
-- interface narrow (wiring di cmd/server/main.go).
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE job_requisitions
    ADD COLUMN IF NOT EXISTS succession_position_id CHAR(36) NULL;
