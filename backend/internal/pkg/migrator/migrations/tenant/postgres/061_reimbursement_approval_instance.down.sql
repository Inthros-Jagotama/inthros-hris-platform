-- 061_reimbursement_approval_instance.down.sql

DROP INDEX IF EXISTS idx_reimb_req_approval_instance;

ALTER TABLE reimbursement_requests
    DROP COLUMN IF EXISTS approval_instance_id;
