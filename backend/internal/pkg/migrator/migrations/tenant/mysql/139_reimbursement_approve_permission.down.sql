-- 139_reimbursement_approve_permission.down.sql

DELETE FROM role_has_permissions
WHERE permission_id = '6e081adf-9ec0-53dc-acc7-3fe7af35dec4';

DELETE FROM permissions
WHERE id = '6e081adf-9ec0-53dc-acc7-3fe7af35dec4';
