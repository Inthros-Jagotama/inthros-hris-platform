-- =============================================================================
-- Tenant Migration: 080_attendance_overtime_two_flow
-- =============================================================================
-- Attendance Module §32b (docs/module-attendance-plan.md) — Overtime dua alur:
--   SELF (ajukan sendiri → approval #1 → isian aktual → approval #2) dan
--   ASSIGNED (atasan tugaskan bawahan → notifikasi → isian aktual → approval).
--   Kedua alur bisa dibatalkan sebelum isian aktual.
--
-- Kolom status di Postgres sudah VARCHAR(255) sejak awal (004_attendance.sql),
-- jadi nilai WAITING_ACTUAL/ACTUAL_SUBMITTED/CANCELLED sudah tertampung tanpa
-- perubahan skema — cukup menambah kolom pendukung alur + index.

ALTER TABLE attendance_overtime_requests
    ADD COLUMN IF NOT EXISTS flow_type VARCHAR(20) NOT NULL DEFAULT 'SELF',
    ADD COLUMN IF NOT EXISTS assigned_by CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMP(6) NULL,
    ADD COLUMN IF NOT EXISTS actual_start_time_local TIMESTAMP(6) NULL,
    ADD COLUMN IF NOT EXISTS actual_end_time_local TIMESTAMP(6) NULL,
    ADD COLUMN IF NOT EXISTS actual_note VARCHAR(500) NULL,
    ADD COLUMN IF NOT EXISTS attachment_url VARCHAR(500) NULL,
    ADD COLUMN IF NOT EXISTS actual_approval_instance_id CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS actual_submitted_at TIMESTAMP(6) NULL,
    ADD COLUMN IF NOT EXISTS actual_approved_by CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS actual_approved_at TIMESTAMP(6) NULL,
    ADD COLUMN IF NOT EXISTS cancelled_by CHAR(36) NULL,
    ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMP(6) NULL;

CREATE INDEX IF NOT EXISTS idx_att_overtime_actual_approval_instance
    ON attendance_overtime_requests (actual_approval_instance_id);
