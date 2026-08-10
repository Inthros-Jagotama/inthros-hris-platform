-- =============================================================================
-- Tenant Migration: 087_employeemovement_cancellation
-- =============================================================================
-- Employee Movement Enhancement (plan §12.16): Movement Cancellation Approval.
--
-- Movement yang sudah `approved` TIDAK boleh dibatalkan langsung oleh HR —
-- pembatalan harus melalui Cancellation Request yang diproses Central Approval
-- Module (module slug `employeemovement_cancellation`). Selama permintaan
-- pembatalan berjalan, movement berada di status `cancellation_pending`.
--
-- Kolom baru: cancellation_approval_instance_id CHAR(36) NULL — approval
-- instance yang dibuat untuk cancellation request (terpisah dari
-- approval_instance_id milik submission).
-- Index baru: idx_emp_mvmt_cancellation_instance.

ALTER TABLE employee_movements
    ADD COLUMN IF NOT EXISTS cancellation_approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_emp_mvmt_cancellation_instance ON employee_movements (cancellation_approval_instance_id);
