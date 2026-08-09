-- =============================================================================
-- Tenant Migration: 081_attendance_overtime_two_flow_backfill
-- =============================================================================
-- Attendance Module §32b (docs/module-attendance-plan.md §32b.8):
-- di alur baru, status APPROVED hanya tercapai SETELAH isian aktual disetujui.
-- Request yang sudah APPROVED sebelum migration 080 (alur lama, tanpa isian
-- aktual) di-backfill ke WAITING_ACTUAL agar karyawan bisa mengisi jam aktual
-- dan menyelesaikan alurnya secara konsisten.

UPDATE attendance_overtime_requests
SET status = 'WAITING_ACTUAL'
WHERE status = 'APPROVED'
  AND actual_submitted_at IS NULL;
