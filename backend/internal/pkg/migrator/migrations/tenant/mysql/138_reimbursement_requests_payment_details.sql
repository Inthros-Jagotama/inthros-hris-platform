-- 138_reimbursement_requests_payment_details.sql
-- Payment recorded directly in the reimbursement module (no payroll linkage —
-- product decision 2026-08-16): method, reference, note captured at PAY time.

ALTER TABLE reimbursement_requests
    ADD COLUMN payment_method    VARCHAR(50)  NULL AFTER paid_amount,
    ADD COLUMN payment_reference VARCHAR(200) NULL AFTER payment_method,
    ADD COLUMN payment_note      VARCHAR(500) NULL AFTER payment_reference;
