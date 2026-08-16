-- =============================================================================
-- 154_sensitive_field_settings_permissions.sql
-- Permission khusus untuk halaman "Sensitive Field Settings" (toggle
-- encrypt-at-rest per field, cakupan tenant-wide). Sebelumnya endpoint
-- GET/PUT /employees/settings/sensitive-fields hanya digating oleh
-- employee.view / employee.update — terlalu longgar: siapa pun yang boleh
-- melihat/mengubah data karyawan jadi bisa mengubah setelan enkripsi.
-- Default: hanya role Admin. Role lain diatur manual lewat halaman RBAC.
-- ID deterministik sama persis dengan codeToUUID di SeedTenantRBAC
-- (uuid.NewSHA1), jadi aman: migrasi & re-seed tidak duplikat.
-- =============================================================================

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at) VALUES
    ('3236f7e6-df60-53d3-af56-58a86defa2dc', 'setting.sensitive-fields.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('4dcc317f-940c-5418-9291-c22df5f82d65', 'setting.sensitive-fields.manage', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id) VALUES
    ('3236f7e6-df60-53d3-af56-58a86defa2dc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345'),
    ('4dcc317f-940c-5418-9291-c22df5f82d65', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');
