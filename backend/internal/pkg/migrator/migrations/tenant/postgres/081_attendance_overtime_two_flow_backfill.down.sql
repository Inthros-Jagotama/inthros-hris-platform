-- =============================================================================
-- Tenant Migration (down): 081_attendance_overtime_two_flow_backfill
-- =============================================================================
-- Perkiraan rollback: request yang di-backfill (WAITING_ACTUAL tanpa isian
-- aktual) dikembalikan ke APPROVED. Request yang sudah mengisi aktual
-- (actual_submitted_at terisi) tidak disentuh.

UPDATE attendance_overtime_requests
SET status = 'APPROVED'
WHERE status = 'WAITING_ACTUAL'
  AND actual_submitted_at IS NULL;
