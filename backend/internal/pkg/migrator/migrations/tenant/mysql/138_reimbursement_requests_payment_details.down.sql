-- 138_reimbursement_requests_payment_details.down.sql

ALTER TABLE reimbursement_requests
    DROP COLUMN IF EXISTS payment_method,
    DROP COLUMN IF EXISTS payment_reference,
    DROP COLUMN IF EXISTS payment_note;
