-- 069_okr_two_phase.sql
-- Two-phase OKR flow: employee proposes Key Results (DRAFT -> KR_SUBMITTED ->
-- KR_APPROVED, "OKR Active") before self-assessment (KR_APPROVED -> SUBMITTED
-- -> COMPLETED). Two independent approval-module checkpoints per evaluation,
-- modules "okr_key_result" and "okr_assessment" — mirrors
-- 068_kpi_approval_instance_ids.sql for KPI.

ALTER TABLE okr_evaluations
    ADD COLUMN IF NOT EXISTS kr_approval_instance_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS assessment_approval_instance_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS kr_submitted_at TIMESTAMP NULL;
