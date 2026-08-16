-- 150_rbac_submenu_permissions_extra.sql
-- Submenu permissions tambahan (tidak tercakup migration 149): resource submenu
-- yang dipakai card hub (Settings, Attendance, Recruitment, Career/WF Intel,
-- Leave Admin, Performance). ID deterministik sama dengan codeToUUID SeedTenantRBAC.

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d3f76b54-9dd5-571e-ac13-8af48a811ca9', 'attendance.corrections.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d3f76b54-9dd5-571e-ac13-8af48a811ca9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d3f76b54-9dd5-571e-ac13-8af48a811ca9', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c53796c0-4f54-5a09-8607-4500690818dd', 'attendance.corrections.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c53796c0-4f54-5a09-8607-4500690818dd', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('22ddc4b4-6555-5d77-8344-9a8495ad6502', 'attendance.corrections.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('22ddc4b4-6555-5d77-8344-9a8495ad6502', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e14d4927-171d-5212-9c00-497c83e79d74', 'attendance.corrections.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e14d4927-171d-5212-9c00-497c83e79d74', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c09b19a2-2059-5656-9646-72a3b0b5effb', 'attendance.admin.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c09b19a2-2059-5656-9646-72a3b0b5effb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c09b19a2-2059-5656-9646-72a3b0b5effb', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('68697b73-2f03-5a17-8da1-07b621dfb74d', 'attendance.admin.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('68697b73-2f03-5a17-8da1-07b621dfb74d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b6da54ed-998b-5e6d-96f7-41ce15c76e30', 'attendance.admin.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b6da54ed-998b-5e6d-96f7-41ce15c76e30', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('aee8d277-d31a-5148-af39-cd001d95655d', 'attendance.admin.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('aee8d277-d31a-5148-af39-cd001d95655d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d34be69e-5208-5dbc-aeab-d667e4b8676b', 'attendance.employee-shifts.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d34be69e-5208-5dbc-aeab-d667e4b8676b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d34be69e-5208-5dbc-aeab-d667e4b8676b', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5e014f91-d564-58a3-aed8-44ee2b9958ea', 'attendance.employee-shifts.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5e014f91-d564-58a3-aed8-44ee2b9958ea', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5a6595f7-84ed-513c-b817-87c24f3a4a41', 'attendance.employee-shifts.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5a6595f7-84ed-513c-b817-87c24f3a4a41', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('efe53475-6ca3-53f0-b309-7150ca230c11', 'attendance.employee-shifts.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('efe53475-6ca3-53f0-b309-7150ca230c11', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6f41363c-de6d-5356-b19f-49376c9bdb21', 'attendance.exempt-positions.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6f41363c-de6d-5356-b19f-49376c9bdb21', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6f41363c-de6d-5356-b19f-49376c9bdb21', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e9066995-ed62-52ec-b696-a5daa35c645d', 'attendance.exempt-positions.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e9066995-ed62-52ec-b696-a5daa35c645d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1b805d63-eeb4-5d01-932c-60ef54f3f142', 'attendance.exempt-positions.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1b805d63-eeb4-5d01-932c-60ef54f3f142', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9c2cdb3d-9ddd-5723-bdd3-3f631ab1d631', 'attendance.exempt-positions.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9c2cdb3d-9ddd-5723-bdd3-3f631ab1d631', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cc934216-17fb-5f9c-8717-56c3e97c3cec', 'attendance.sessions.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc934216-17fb-5f9c-8717-56c3e97c3cec', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc934216-17fb-5f9c-8717-56c3e97c3cec', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8b843684-8771-5f6f-b2f4-8f58840abe35', 'attendance.sessions.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8b843684-8771-5f6f-b2f4-8f58840abe35', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f9608099-02d3-580a-ad32-41c22142f150', 'attendance.sessions.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f9608099-02d3-580a-ad32-41c22142f150', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('15dba450-2e29-5c1e-9718-087f0622343f', 'attendance.sessions.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('15dba450-2e29-5c1e-9718-087f0622343f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('858b3454-a3e9-596c-bdbc-b5547cb3f67e', 'attendance.reports.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('858b3454-a3e9-596c-bdbc-b5547cb3f67e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('858b3454-a3e9-596c-bdbc-b5547cb3f67e', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0891bf18-b14c-5593-acdd-6f01d844e6b2', 'attendance.reports.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0891bf18-b14c-5593-acdd-6f01d844e6b2', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cca518a7-df35-5a3b-92e9-f0801d91ee19', 'attendance.reports.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cca518a7-df35-5a3b-92e9-f0801d91ee19', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ca3f0ab2-82cd-5354-9897-de4bc041d10e', 'attendance.reports.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ca3f0ab2-82cd-5354-9897-de4bc041d10e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('204c7f94-bb10-5070-9599-4e161a892b5a', 'recruitment.internal-candidates.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('204c7f94-bb10-5070-9599-4e161a892b5a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('204c7f94-bb10-5070-9599-4e161a892b5a', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('19b91cd9-a333-51f7-850c-e8c616cbd410', 'recruitment.internal-candidates.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('19b91cd9-a333-51f7-850c-e8c616cbd410', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e5ed6561-de01-5413-8143-e9b97a13a172', 'recruitment.internal-candidates.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e5ed6561-de01-5413-8143-e9b97a13a172', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b5a17a38-1512-5e6b-8b15-718c0ff4a622', 'recruitment.internal-candidates.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b5a17a38-1512-5e6b-8b15-718c0ff4a622', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3eaaabe5-fca4-5710-9a96-078ffbb1019a', 'recruitment.offers.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3eaaabe5-fca4-5710-9a96-078ffbb1019a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3eaaabe5-fca4-5710-9a96-078ffbb1019a', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ac9079f9-2717-5e8b-9384-3d37f46157db', 'recruitment.offers.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ac9079f9-2717-5e8b-9384-3d37f46157db', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cf6bc63c-d588-5eef-9562-3efb9be82015', 'recruitment.offers.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cf6bc63c-d588-5eef-9562-3efb9be82015', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4dac04bd-0b05-50fb-ad50-56c60b5bfd28', 'recruitment.offers.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4dac04bd-0b05-50fb-ad50-56c60b5bfd28', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1b829f66-0793-5e2a-add4-d8599ca32439', 'recruitment.assessments.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1b829f66-0793-5e2a-add4-d8599ca32439', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1b829f66-0793-5e2a-add4-d8599ca32439', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4792b3b6-0daf-5a93-b6d7-4cd8737d4848', 'recruitment.assessments.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4792b3b6-0daf-5a93-b6d7-4cd8737d4848', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4a66b455-c228-5f77-84a2-41bdfc6d4d9c', 'recruitment.assessments.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4a66b455-c228-5f77-84a2-41bdfc6d4d9c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('110cc565-956e-5d5c-a66f-b081f461af02', 'recruitment.assessments.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('110cc565-956e-5d5c-a66f-b081f461af02', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8077ff28-ef3e-5dcd-806e-1173671d48d4', 'workforceintelligence.candidate-search.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8077ff28-ef3e-5dcd-806e-1173671d48d4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8077ff28-ef3e-5dcd-806e-1173671d48d4', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d54371d6-ace7-5e5f-be06-074292534134', 'workforceintelligence.candidate-search.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d54371d6-ace7-5e5f-be06-074292534134', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f4714cfe-4030-5cca-a44c-07ea54d79fd4', 'workforceintelligence.candidate-search.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f4714cfe-4030-5cca-a44c-07ea54d79fd4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('33740e93-ff19-526c-99c6-714716de400f', 'workforceintelligence.candidate-search.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('33740e93-ff19-526c-99c6-714716de400f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5c14a8e4-8d55-58bd-a707-fa2ee82754d6', 'workforceintelligence.recruitment-analytics.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5c14a8e4-8d55-58bd-a707-fa2ee82754d6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5c14a8e4-8d55-58bd-a707-fa2ee82754d6', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('659dd4de-ea94-52a5-b4e6-027195775c8c', 'workforceintelligence.recruitment-analytics.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('659dd4de-ea94-52a5-b4e6-027195775c8c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cfb639e3-30d4-5e41-b228-19c85ab043ad', 'workforceintelligence.recruitment-analytics.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cfb639e3-30d4-5e41-b228-19c85ab043ad', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5af37aba-eb97-5647-ba42-b50d9347635e', 'workforceintelligence.recruitment-analytics.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5af37aba-eb97-5647-ba42-b50d9347635e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('275cb1ea-e027-5426-ae96-d982ccb13cab', 'workforceintelligence.quality-of-hire.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('275cb1ea-e027-5426-ae96-d982ccb13cab', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('275cb1ea-e027-5426-ae96-d982ccb13cab', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ab6bad52-2e8c-510b-8dee-f7c0dbf229ae', 'workforceintelligence.quality-of-hire.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ab6bad52-2e8c-510b-8dee-f7c0dbf229ae', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9487f69d-c0c0-5a65-9d88-9086b5a72b27', 'workforceintelligence.quality-of-hire.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9487f69d-c0c0-5a65-9d88-9086b5a72b27', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d7515849-33c5-5f62-9660-31c5ff881d0a', 'workforceintelligence.quality-of-hire.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d7515849-33c5-5f62-9660-31c5ff881d0a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('326e11fd-c5b7-53d8-ad7f-0f5cbf8e440f', 'workforceintelligence.training-analysis.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('326e11fd-c5b7-53d8-ad7f-0f5cbf8e440f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('326e11fd-c5b7-53d8-ad7f-0f5cbf8e440f', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('99dd0c9b-6451-577a-b791-fdea404901f6', 'workforceintelligence.training-analysis.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('99dd0c9b-6451-577a-b791-fdea404901f6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2117fdcd-5222-5cad-8d5f-6062986ab490', 'workforceintelligence.training-analysis.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2117fdcd-5222-5cad-8d5f-6062986ab490', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0ce5eed7-388f-5ba4-abfa-49baf979e907', 'workforceintelligence.training-analysis.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0ce5eed7-388f-5ba4-abfa-49baf979e907', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1873a269-d808-59f4-be39-a7ceefe15df9', 'careerintelligence.paths.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1873a269-d808-59f4-be39-a7ceefe15df9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1873a269-d808-59f4-be39-a7ceefe15df9', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cce58068-15ea-57b8-893a-09c655189924', 'careerintelligence.paths.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cce58068-15ea-57b8-893a-09c655189924', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7a312f25-4cb2-59c1-a43a-288d00e989f0', 'careerintelligence.paths.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7a312f25-4cb2-59c1-a43a-288d00e989f0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4d27cf5f-b45b-5b28-a31d-da468db7cb0c', 'careerintelligence.paths.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4d27cf5f-b45b-5b28-a31d-da468db7cb0c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('45f86032-8e9f-567f-9399-17322a1d7957', 'careerintelligence.successions.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('45f86032-8e9f-567f-9399-17322a1d7957', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('45f86032-8e9f-567f-9399-17322a1d7957', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b0b17ca1-7458-554b-8fce-f85a52c5b23d', 'careerintelligence.successions.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b0b17ca1-7458-554b-8fce-f85a52c5b23d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('089adaf3-28f7-5ad8-ae33-35cf1bcdc0f4', 'careerintelligence.successions.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('089adaf3-28f7-5ad8-ae33-35cf1bcdc0f4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c9eed7c8-43dc-5e6f-8914-4a7be51bd815', 'careerintelligence.successions.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c9eed7c8-43dc-5e6f-8914-4a7be51bd815', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('56a9a462-55e0-5460-aa36-3bf015870a72', 'careerintelligence.development.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56a9a462-55e0-5460-aa36-3bf015870a72', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56a9a462-55e0-5460-aa36-3bf015870a72', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6a7f51ba-0e76-515e-a9bf-9ca4b0d1315f', 'careerintelligence.development.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6a7f51ba-0e76-515e-a9bf-9ca4b0d1315f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8351ef89-d16e-55f1-9010-fe9f4607d097', 'careerintelligence.development.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8351ef89-d16e-55f1-9010-fe9f4607d097', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a2a951c1-e2fc-5950-a991-e1fc681d9223', 'careerintelligence.development.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a2a951c1-e2fc-5950-a991-e1fc681d9223', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e78b0adb-9fce-550e-950d-0201f1ba4d3b', 'setting.numbering.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e78b0adb-9fce-550e-950d-0201f1ba4d3b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e78b0adb-9fce-550e-950d-0201f1ba4d3b', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b6604cbd-c59a-5a4e-9677-86ddf51b9a55', 'setting.numbering.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b6604cbd-c59a-5a4e-9677-86ddf51b9a55', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('47185461-87c7-5aa6-8817-213f863fee03', 'setting.numbering.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('47185461-87c7-5aa6-8817-213f863fee03', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bb01216d-ecd4-5565-95b8-064724ce493b', 'setting.numbering.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bb01216d-ecd4-5565-95b8-064724ce493b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bf121f50-94c6-5ed9-9cb7-5668f1a00a86', 'setting.performance-perspectives.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bf121f50-94c6-5ed9-9cb7-5668f1a00a86', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bf121f50-94c6-5ed9-9cb7-5668f1a00a86', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('05879bfc-afcd-5594-ac8a-6d263ed78596', 'setting.performance-perspectives.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('05879bfc-afcd-5594-ac8a-6d263ed78596', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6aed3648-9c0e-5583-960d-3f4ece23ccb9', 'setting.performance-perspectives.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6aed3648-9c0e-5583-960d-3f4ece23ccb9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cc6c4fe8-e150-5366-9e13-36139aaf32ad', 'setting.performance-perspectives.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc6c4fe8-e150-5366-9e13-36139aaf32ad', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c833878f-d5f0-5164-90ad-cf8d71aa11ec', 'setting.performance-ratings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c833878f-d5f0-5164-90ad-cf8d71aa11ec', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c833878f-d5f0-5164-90ad-cf8d71aa11ec', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0bd7f55c-70da-590c-af79-b34d8f327d3a', 'setting.performance-ratings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0bd7f55c-70da-590c-af79-b34d8f327d3a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('457ae84d-4a33-57a8-9e65-33854775a55e', 'setting.performance-ratings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('457ae84d-4a33-57a8-9e65-33854775a55e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b7663c07-d1c8-5462-80e1-568d84b1169d', 'setting.performance-ratings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b7663c07-d1c8-5462-80e1-568d84b1169d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7ea5fab8-6471-57aa-b195-dbe163940163', 'setting.performance-formulas.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7ea5fab8-6471-57aa-b195-dbe163940163', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7ea5fab8-6471-57aa-b195-dbe163940163', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('42482b55-53b3-5290-9be6-f91a35077c88', 'setting.performance-formulas.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('42482b55-53b3-5290-9be6-f91a35077c88', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('29d7ae74-8d27-5532-a800-6360002535a9', 'setting.performance-formulas.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('29d7ae74-8d27-5532-a800-6360002535a9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fe5c34c8-7d05-5348-8f7a-04359245699a', 'setting.performance-formulas.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fe5c34c8-7d05-5348-8f7a-04359245699a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cb05f7be-319b-5ee6-8e15-055c03444133', 'setting.performance-components.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cb05f7be-319b-5ee6-8e15-055c03444133', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cb05f7be-319b-5ee6-8e15-055c03444133', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4b7ca604-4e03-55f2-818c-f884380180e9', 'setting.performance-components.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4b7ca604-4e03-55f2-818c-f884380180e9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3fdbec58-de7e-56d2-b612-d6a08c58e5eb', 'setting.performance-components.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3fdbec58-de7e-56d2-b612-d6a08c58e5eb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9a126ad2-0b8b-5567-b5c0-f9c6f0a7f52a', 'setting.performance-components.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9a126ad2-0b8b-5567-b5c0-f9c6f0a7f52a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('daf61bfd-191f-5cbf-b3ba-cec8ea92bcf2', 'setting.performance-scoring.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('daf61bfd-191f-5cbf-b3ba-cec8ea92bcf2', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('daf61bfd-191f-5cbf-b3ba-cec8ea92bcf2', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('19bf514f-87ef-5310-835d-fe336d6f4485', 'setting.performance-scoring.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('19bf514f-87ef-5310-835d-fe336d6f4485', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f58a032c-1548-5cf8-8bff-bd4cf99c1638', 'setting.performance-scoring.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f58a032c-1548-5cf8-8bff-bd4cf99c1638', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('95eff26e-9e66-5775-a169-c5f24fcb3252', 'setting.performance-scoring.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('95eff26e-9e66-5775-a169-c5f24fcb3252', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a22bc504-051a-5075-8804-b9223d70bb01', 'setting.payroll-periods.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a22bc504-051a-5075-8804-b9223d70bb01', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a22bc504-051a-5075-8804-b9223d70bb01', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e9e107e0-a4d8-551e-9e4a-b625d98f13cf', 'setting.payroll-periods.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e9e107e0-a4d8-551e-9e4a-b625d98f13cf', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a82809de-09fb-5033-9ae4-b9231a5264b8', 'setting.payroll-periods.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a82809de-09fb-5033-9ae4-b9231a5264b8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a5cabe84-ae23-555c-95fc-480ebc71d2ef', 'setting.payroll-periods.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a5cabe84-ae23-555c-95fc-480ebc71d2ef', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1bb8cc35-20ad-5613-a93e-645885b702e5', 'setting.salary-components.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1bb8cc35-20ad-5613-a93e-645885b702e5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1bb8cc35-20ad-5613-a93e-645885b702e5', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5430a4e2-eaba-518d-8d66-d400c018af40', 'setting.salary-components.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5430a4e2-eaba-518d-8d66-d400c018af40', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('91e65f07-0a5d-5dd3-a56c-41a7a8082b69', 'setting.salary-components.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('91e65f07-0a5d-5dd3-a56c-41a7a8082b69', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('483b9cdd-cb39-589b-96ab-2ea27fb6a89c', 'setting.salary-components.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('483b9cdd-cb39-589b-96ab-2ea27fb6a89c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0ca0a212-f80c-5819-8e22-b1fcd76fc8ed', 'setting.bpjs-settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0ca0a212-f80c-5819-8e22-b1fcd76fc8ed', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0ca0a212-f80c-5819-8e22-b1fcd76fc8ed', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ea4825f0-6571-507a-b5b2-41b9938d86c6', 'setting.bpjs-settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ea4825f0-6571-507a-b5b2-41b9938d86c6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5b6b1584-fc21-572e-b720-aebe1e28af41', 'setting.bpjs-settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5b6b1584-fc21-572e-b720-aebe1e28af41', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5a6eace9-1122-5bf8-8b65-ad6443b8bf40', 'setting.bpjs-settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5a6eace9-1122-5bf8-8b65-ad6443b8bf40', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('930d0a6d-dd64-5ea7-9bf8-c82acef61b5d', 'setting.pph21-settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('930d0a6d-dd64-5ea7-9bf8-c82acef61b5d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('930d0a6d-dd64-5ea7-9bf8-c82acef61b5d', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6bbcf755-85da-5cb2-a1e6-1a390af107ca', 'setting.pph21-settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6bbcf755-85da-5cb2-a1e6-1a390af107ca', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b96671c1-c02b-5694-86b2-d4b5a52c33f4', 'setting.pph21-settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b96671c1-c02b-5694-86b2-d4b5a52c33f4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7cbdd314-036e-5ad4-a325-60fc6610b08a', 'setting.pph21-settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7cbdd314-036e-5ad4-a325-60fc6610b08a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6700054f-7311-50e5-9e80-d742536d1727', 'setting.salary-structure.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6700054f-7311-50e5-9e80-d742536d1727', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6700054f-7311-50e5-9e80-d742536d1727', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6fabe8ca-6971-5578-bc14-06e5008bfad3', 'setting.salary-structure.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6fabe8ca-6971-5578-bc14-06e5008bfad3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d18afff6-de05-5498-89c1-0f6cfc0f8082', 'setting.salary-structure.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d18afff6-de05-5498-89c1-0f6cfc0f8082', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0a6e9701-02ba-587b-bb18-19d5f6cfefdd', 'setting.salary-structure.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0a6e9701-02ba-587b-bb18-19d5f6cfefdd', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a46e3bbb-5187-529a-8629-49153704355f', 'leave.accrual-policies.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a46e3bbb-5187-529a-8629-49153704355f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a46e3bbb-5187-529a-8629-49153704355f', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ebaaf572-0ade-54db-9512-aab1515293ba', 'leave.accrual-policies.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ebaaf572-0ade-54db-9512-aab1515293ba', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7b3b65e2-2281-554b-a1b9-7c873d5ba400', 'leave.accrual-policies.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7b3b65e2-2281-554b-a1b9-7c873d5ba400', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d0f56bff-64e5-52dd-9ab0-25c0e7f2dfb1', 'leave.accrual-policies.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d0f56bff-64e5-52dd-9ab0-25c0e7f2dfb1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('63a27150-bba7-53ac-831e-63ff3f5df92d', 'leave.reasons.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('63a27150-bba7-53ac-831e-63ff3f5df92d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('63a27150-bba7-53ac-831e-63ff3f5df92d', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fb9133d2-f5d0-51bb-a75b-085246ae9ad0', 'leave.reasons.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fb9133d2-f5d0-51bb-a75b-085246ae9ad0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f252b1d0-7b98-5518-be02-22edf049d2ec', 'leave.reasons.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f252b1d0-7b98-5518-be02-22edf049d2ec', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6a99584a-9da7-5208-bf42-4b061e04327c', 'leave.reasons.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6a99584a-9da7-5208-bf42-4b061e04327c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('28da771d-adb7-5245-a291-1f0d1e42e11b', 'performance.kpi.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('28da771d-adb7-5245-a291-1f0d1e42e11b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('28da771d-adb7-5245-a291-1f0d1e42e11b', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b3574559-4a50-5a47-b984-7a24c8df39ec', 'performance.kpi.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b3574559-4a50-5a47-b984-7a24c8df39ec', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e25d4354-c31a-55ad-8a78-5a5c16e3c569', 'performance.kpi.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e25d4354-c31a-55ad-8a78-5a5c16e3c569', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('42f8cb97-4f3d-501d-bc7e-36e12b86db1f', 'performance.kpi.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('42f8cb97-4f3d-501d-bc7e-36e12b86db1f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c904daa6-2fe0-5f8f-b7e2-5a053b2e7302', 'performance.okr.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c904daa6-2fe0-5f8f-b7e2-5a053b2e7302', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c904daa6-2fe0-5f8f-b7e2-5a053b2e7302', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('461b205b-7dd9-5727-ad38-87e04df66ae6', 'performance.okr.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('461b205b-7dd9-5727-ad38-87e04df66ae6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('91d2576c-43a2-51d7-9f98-74d20d4a9320', 'performance.okr.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('91d2576c-43a2-51d7-9f98-74d20d4a9320', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2862cf6c-0d71-5478-a4d4-274de1f8f023', 'performance.okr.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2862cf6c-0d71-5478-a4d4-274de1f8f023', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

