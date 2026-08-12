-- =============================================================================
-- Tenant Migration: 093_recruitment_approval (MySQL)
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
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN approval_instance_id CHAR(36) NULL AFTER approved_by'
);
PREPARE stmt FROM @add_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
