-- =============================================================================
-- Tenant Migration: 094_recruitment_requisition_enhancement (MySQL)
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
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE.

SET @add_requisition_number = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'requisition_number'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN requisition_number VARCHAR(50) NULL AFTER title'
);
PREPARE stmt FROM @add_requisition_number;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_priority = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'priority'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN priority VARCHAR(10) NOT NULL DEFAULT ''MEDIUM'' AFTER requisition_number'
);
PREPARE stmt FROM @add_priority;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_position_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'position_id'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN position_id CHAR(36) NULL AFTER priority'
);
PREPARE stmt FROM @add_position_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @add_opened_at = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'job_requisitions'
      AND COLUMN_NAME = 'opened_at'
  ),
  'DO 0',
  'ALTER TABLE job_requisitions ADD COLUMN opened_at BIGINT NULL AFTER closed_at'
);
PREPARE stmt FROM @add_opened_at;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
