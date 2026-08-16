-- 137_reimbursement_requests_fix_employee_id.sql
-- Repair misattributed reimbursement requests: rows created before the
-- employee_id fix (2026-08-16) stored the JWT user-account UUID in
-- employee_id, so they never showed up under the employee filter
-- (employee_accounts maps user_id -> employee_id, two distinct UUIDs).
-- Backfill the real employee UUID wherever a mapping exists.

UPDATE reimbursement_requests rr
JOIN employee_accounts ea ON ea.user_id = rr.employee_id
SET rr.employee_id = ea.employee_id;
