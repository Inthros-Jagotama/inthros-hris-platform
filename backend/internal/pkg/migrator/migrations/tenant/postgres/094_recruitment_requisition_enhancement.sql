-- =============================================================================
-- Tenant Migration: 094_recruitment_requisition_enhancement (PostgreSQL)
-- =============================================================================
-- G-2 Job Requisition Enhancement (master position + operational fields)
-- (docs/module-recruitment-development-plan.md §G-2)
--
-- Field operasional baru untuk job_requisitions:
--
--   requisition_number   VARCHAR(50) NULL  → nomor requisition (auto-generated
--                        REQ-YYYYMM-XXXX; pola nomor dokumen TRN-* training)
--   priority             VARCHAR(10) NOT NULL DEFAULT 'MEDIUM'
--                        (LOW | MEDIUM | HIGH | URGENT)
--   position_id          CHAR(36)    NULL  → referensi master position
--                        (tabel positions — tanpa FK karena modul Organization
--                        tidak mengekspos CRUD position; referensi saja)
--   opened_at            BIGINT      NULL  → unix nano saat requisition
--                        menjadi OPEN (diset otomatis saat approval APPROVED
--                        atau transisi status manual ke OPEN)
--
-- Catatan: approval_status TIDAK ditambahkan — G-1 sudah meng-cover via status
-- requisition (SUBMITTED/OPEN/REJECTED) + approval_instance_id.
--
-- Idempotent: ADD COLUMN IF NOT EXISTS.

ALTER TABLE job_requisitions
    ADD COLUMN IF NOT EXISTS requisition_number VARCHAR(50) NULL,
    ADD COLUMN IF NOT EXISTS priority VARCHAR(10) NOT NULL DEFAULT 'MEDIUM',
    ADD COLUMN IF NOT EXISTS position_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS opened_at BIGINT NULL;
