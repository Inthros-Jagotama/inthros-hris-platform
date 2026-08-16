-- =============================================================================
-- 153_sensitive_field_view_permissions.sql
-- Sensitive Data Masking — permission per-field untuk melihat nilai asli
-- (bukan hasil masking). Default: hanya role Admin yang diberi akses;
-- role lain diatur manual lewat halaman RBAC.
-- ID deterministik sama persis dengan codeToUUID di SeedTenantRBAC
-- (uuid.NewSHA1), jadi aman: migrasi & re-seed tidak duplikat.
-- =============================================================================

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at) VALUES
    ('57032606-fdf3-5437-a06b-4f6b78bb585f', 'employee.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('6820f1d5-9dfd-54e2-a100-007e923a9434', 'employee.view_passport', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('24fc1ecd-deea-5c59-8b20-e4cd84480ef8', 'employee.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('e16e3dc0-d2b8-5ae8-a316-168300d857bd', 'employee.view_email', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('6fa7ca46-bb80-5eb2-8a4c-869a55e0d549', 'employee_family.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('5a86c981-34e9-55e5-a481-11870bb822c1', 'employee_bank_account.view_account_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('516d708f-2b4f-5a7a-8bfd-99a9d6bc39d0', 'employee_bank_account.view_account_name', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('3688b8ae-6c3f-5213-a05d-9e0bb767a829', 'emergency_contact.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id) VALUES
    ('57032606-fdf3-5437-a06b-4f6b78bb585f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('6820f1d5-9dfd-54e2-a100-007e923a9434', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('24fc1ecd-deea-5c59-8b20-e4cd84480ef8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('e16e3dc0-d2b8-5ae8-a316-168300d857bd', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('6fa7ca46-bb80-5eb2-8a4c-869a55e0d549', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('5a86c981-34e9-55e5-a481-11870bb822c1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('516d708f-2b4f-5a7a-8bfd-99a9d6bc39d0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('3688b8ae-6c3f-5213-a05d-9e0bb767a829', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');
