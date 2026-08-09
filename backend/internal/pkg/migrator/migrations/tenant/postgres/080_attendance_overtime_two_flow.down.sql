-- =============================================================================
-- Tenant Migration (down): 080_attendance_overtime_two_flow
-- =============================================================================
-- Kembalikan attendance_overtime_requests ke struktur sebelum dua-alur:
-- drop index + 13 kolom pendukung. (Status Postgres tetap VARCHAR — tidak ada
-- enum yang perlu di-revert.)

DROP INDEX IF EXISTS idx_att_overtime_actual_approval_instance;

ALTER TABLE attendance_overtime_requests
    DROP COLUMN IF EXISTS flow_type,
    DROP COLUMN IF EXISTS assigned_by,
    DROP COLUMN IF EXISTS assigned_at,
    DROP COLUMN IF EXISTS actual_start_time_local,
    DROP COLUMN IF EXISTS actual_end_time_local,
    DROP COLUMN IF EXISTS actual_note,
    DROP COLUMN IF EXISTS attachment_url,
    DROP COLUMN IF EXISTS actual_approval_instance_id,
    DROP COLUMN IF EXISTS actual_submitted_at,
    DROP COLUMN IF EXISTS actual_approved_by,
    DROP COLUMN IF EXISTS actual_approved_at,
    DROP COLUMN IF EXISTS cancelled_by,
    DROP COLUMN IF EXISTS cancelled_at;
