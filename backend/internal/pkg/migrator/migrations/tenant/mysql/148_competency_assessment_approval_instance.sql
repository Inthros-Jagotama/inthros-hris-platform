-- =============================================================================
-- 148_competency_assessment_approval_instance.sql
-- Competency 360 Module — Approval Engine Integration (plan generik §13, §34.2).
-- Finalisasi hasil assessment diproses Central Approval: saat seluruh rater
-- submit (atau HR menutup assessment), dibuat approval instance ber-module
-- competency_360_assessment untuk document = competency_event_target.
-- Kolom approval_instance_id menyimpan instance tsb (pola migrasi 061/133);
-- kolom status menyimpan business-state assessment (draft/open/in_progress/
-- submitted/finalized) — approval status ditangani engine, bukan di sini.
-- =============================================================================

ALTER TABLE competency_event_targets
    ADD COLUMN approval_instance_id CHAR(36) NULL AFTER employee_id,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'draft' AFTER approval_instance_id,
    ADD COLUMN finalized_at TIMESTAMP NULL AFTER status;

ALTER TABLE competency_event_targets
    ADD INDEX idx_comp_target_approval_instance (approval_instance_id),
    ADD INDEX idx_comp_target_status (status);
