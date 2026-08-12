-- =============================================================================
-- Tenant Migration: 093_recruitment_approval (PostgreSQL)
-- =============================================================================
-- G-1 Module Approval Integration (requisition)
-- (docs/module-recruitment-development-plan.md §G-1)
--
-- Job requisition dirutekan melalui Central Approval:
--
--   DRAFT → SUBMITTED → (Module Approval) → APPROVED/REJECTED → OPEN
--
-- Kolom berikut menyimpan instance approval yang sedang/berhasil memproses
-- requisition. Status keputusan (APPROVED/REJECTED) didorong kembali dari
-- modul approval via status handler push-callback (pola employeemovement) —
-- Approval module tetap source of truth untuk proses persetujuan.
--
--   approval_instance_id  CHAR(36)  NULL  → id ApprovalInstance (Central
--                        Approval) untuk requisition ini.
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE job_requisitions
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;
