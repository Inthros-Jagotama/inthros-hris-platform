-- 061_reimbursement_approval_instance.sql
-- Approval Module Phase 4: persist the approval_instance_id created for a
-- reimbursement request so its status can be updated via push-based
-- callback from the central approval module.

ALTER TABLE reimbursement_requests
    ADD COLUMN IF NOT EXISTS approval_instance_id CHAR(36) NULL;

CREATE INDEX IF NOT EXISTS idx_reimb_req_approval_instance ON reimbursement_requests (approval_instance_id);
