-- 061_reimbursement_approval_instance.down.sql

SET @drop_reimb_approval_index = IF(
  EXISTS(
    SELECT 1 FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'reimbursement_requests'
      AND INDEX_NAME = 'idx_reimb_req_approval_instance'
  ),
  'DROP INDEX idx_reimb_req_approval_instance ON reimbursement_requests',
  'DO 0'
);
PREPARE stmt FROM @drop_reimb_approval_index;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @drop_approval_instance_id = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'reimbursement_requests'
      AND COLUMN_NAME = 'approval_instance_id'
  ),
  'ALTER TABLE reimbursement_requests DROP COLUMN approval_instance_id',
  'DO 0'
);
PREPARE stmt FROM @drop_approval_instance_id;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
