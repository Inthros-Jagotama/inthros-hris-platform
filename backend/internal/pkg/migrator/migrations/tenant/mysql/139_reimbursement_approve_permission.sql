-- 139_reimbursement_approve_permission.sql
-- Seed RBAC yang hilang: permission `reimbursement.approve` + link ke role Admin.
-- Seed asli (seed_rbac.go) hanya meng-seed view/create/update/delete untuk
-- reimbursement, sehingga gate FE `hasPermission('reimbursement.approve')`
-- (tombol Pay, mode admin melihat semua request) tidak pernah terpenuhi.
-- ID deterministik sama persis dengan yang dihasilkan codeToUUID di
-- SeedTenantRBAC (uuid.NewSHA1), jadi aman: migrasi & re-seed tidak duplikat.

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6e081adf-9ec0-53dc-acc7-3fe7af35dec4', 'reimbursement.approve', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6e081adf-9ec0-53dc-acc7-3fe7af35dec4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');
