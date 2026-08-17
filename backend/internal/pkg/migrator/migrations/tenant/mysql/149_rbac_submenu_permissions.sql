-- 149_rbac_submenu_permissions.sql
-- Seed permission level-submenu ("resource.submenu.action") untuk tenant existing.
-- ID deterministik sama persis dengan codeToUUID di SeedTenantRBAC, jadi aman
-- dijalankan berulang (ON CONFLICT / INSERT IGNORE).
-- Admin mendapat semua action; Employee hanya view.
-- Catatan: organization, setting, employee, & employeemovement TIDAK punya
-- permission level-submenu —
-- cukup module-level (organization.view/create/update/delete, dst).
-- jobmanagement disederhanakan jadi 2 submenu: setting (titles, objectives,
-- identifications, responsibilities, authorities, working-conditions,
-- competencies) & assessment (values, scores).
-- competency disederhanakan jadi 3 submenu: settings (competencies, values,
-- indicators, templates, events, raters), assessment (my-assessments,
-- manager-assessments, scores), & report (results, reports — view only).

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('89a2d0a0-e8de-5ef1-a2e2-691ddd21a503', 'attendance.dashboard.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('89a2d0a0-e8de-5ef1-a2e2-691ddd21a503', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('89a2d0a0-e8de-5ef1-a2e2-691ddd21a503', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('68fa6bb5-bda7-50b9-93f5-7f3401777094', 'attendance.dashboard.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('68fa6bb5-bda7-50b9-93f5-7f3401777094', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fffa8a75-5945-5dfa-935a-37d0d5429cb2', 'attendance.dashboard.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fffa8a75-5945-5dfa-935a-37d0d5429cb2', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f4638f7f-3324-5de2-818e-037759ea90c5', 'attendance.dashboard.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f4638f7f-3324-5de2-818e-037759ea90c5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('81ca4559-35c3-5d88-aceb-a82bb44c1721', 'attendance.shifts.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('81ca4559-35c3-5d88-aceb-a82bb44c1721', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('81ca4559-35c3-5d88-aceb-a82bb44c1721', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4cb4ea07-658f-501e-8a62-e5892f493e57', 'attendance.shifts.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4cb4ea07-658f-501e-8a62-e5892f493e57', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f1359e90-712b-5eaa-bb6e-a4997369e2a0', 'attendance.shifts.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f1359e90-712b-5eaa-bb6e-a4997369e2a0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e3204f30-d9fd-527f-b19a-6912ce24a21f', 'attendance.shifts.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e3204f30-d9fd-527f-b19a-6912ce24a21f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('79a83f86-d6e1-58d4-810a-0f17d2298dac', 'attendance.schedules.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('79a83f86-d6e1-58d4-810a-0f17d2298dac', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('79a83f86-d6e1-58d4-810a-0f17d2298dac', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('56f5f7f6-6142-587c-87c7-9029a5e8a565', 'attendance.schedules.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56f5f7f6-6142-587c-87c7-9029a5e8a565', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f4c59610-0886-59a0-bfc9-41483719da4f', 'attendance.schedules.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f4c59610-0886-59a0-bfc9-41483719da4f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0c36958a-e9d4-5174-b048-56b02da775b3', 'attendance.schedules.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0c36958a-e9d4-5174-b048-56b02da775b3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('af76ae9c-d7fa-5506-9812-f3796f4a48db', 'attendance.events.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('af76ae9c-d7fa-5506-9812-f3796f4a48db', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('af76ae9c-d7fa-5506-9812-f3796f4a48db', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3c224edb-6800-5168-8f04-e31bc14be468', 'attendance.events.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3c224edb-6800-5168-8f04-e31bc14be468', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bc94dd67-7c2a-56cd-8f81-198cd103b059', 'attendance.events.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bc94dd67-7c2a-56cd-8f81-198cd103b059', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2b2220fd-c0ab-54e4-b219-34f613b2f0c1', 'attendance.events.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2b2220fd-c0ab-54e4-b219-34f613b2f0c1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('659f01c7-a5f1-52c3-b157-bf9db85b7d2e', 'attendance.overtime.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('659f01c7-a5f1-52c3-b157-bf9db85b7d2e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('659f01c7-a5f1-52c3-b157-bf9db85b7d2e', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3a21b1bc-cea8-5018-97b0-b3a8a6d7192b', 'attendance.overtime.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3a21b1bc-cea8-5018-97b0-b3a8a6d7192b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ed9b57d9-73f5-5f87-b408-332ed688d1f6', 'attendance.overtime.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ed9b57d9-73f5-5f87-b408-332ed688d1f6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('48558a00-469e-5d49-829d-bf16d41318ef', 'attendance.overtime.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('48558a00-469e-5d49-829d-bf16d41318ef', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('94441d07-b611-5c6a-b32b-698892599a2e', 'attendance.locations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('94441d07-b611-5c6a-b32b-698892599a2e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('94441d07-b611-5c6a-b32b-698892599a2e', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8f1c83fb-b2a4-5ff9-a64d-aba55855762d', 'attendance.locations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8f1c83fb-b2a4-5ff9-a64d-aba55855762d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c3d699fd-bdb4-5601-a124-18aa1f53fab0', 'attendance.locations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c3d699fd-bdb4-5601-a124-18aa1f53fab0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d31288fd-5923-53ef-8654-59d540cf3868', 'attendance.locations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d31288fd-5923-53ef-8654-59d540cf3868', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7fedcc32-fbd2-5f16-8b0e-5d5379f0340e', 'attendance.business-travel.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7fedcc32-fbd2-5f16-8b0e-5d5379f0340e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7fedcc32-fbd2-5f16-8b0e-5d5379f0340e', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2bb49c2a-5e5c-540b-8bb3-241dfb2d7e69', 'attendance.business-travel.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2bb49c2a-5e5c-540b-8bb3-241dfb2d7e69', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d247849b-a555-5147-9169-44991de65d44', 'attendance.business-travel.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d247849b-a555-5147-9169-44991de65d44', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5b216e6d-816c-515d-86e9-f5b34ba7ffeb', 'attendance.business-travel.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5b216e6d-816c-515d-86e9-f5b34ba7ffeb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b85971d9-9ba0-5f51-a87f-da8d94e4bf31', 'attendance.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b85971d9-9ba0-5f51-a87f-da8d94e4bf31', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b85971d9-9ba0-5f51-a87f-da8d94e4bf31', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('85b5ece0-f963-52e6-aa2f-d6a36b5d8c55', 'attendance.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('85b5ece0-f963-52e6-aa2f-d6a36b5d8c55', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d4ab5d8d-2983-50e6-a624-5b6e824cee5f', 'attendance.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d4ab5d8d-2983-50e6-a624-5b6e824cee5f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('66786d93-2f19-5409-b442-56bc565fe263', 'attendance.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('66786d93-2f19-5409-b442-56bc565fe263', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f8198f71-ec65-5c9b-96dc-e65534125407', 'approval.tasks.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f8198f71-ec65-5c9b-96dc-e65534125407', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f8198f71-ec65-5c9b-96dc-e65534125407', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5b1c464f-ca91-5fa9-9c75-7f4db13daffc', 'approval.tasks.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5b1c464f-ca91-5fa9-9c75-7f4db13daffc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b0816bbc-98ee-5029-9466-4e2a2ed986e1', 'approval.tasks.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b0816bbc-98ee-5029-9466-4e2a2ed986e1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a6b9bdeb-42ce-51f3-8afa-3ae98f062f83', 'approval.tasks.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a6b9bdeb-42ce-51f3-8afa-3ae98f062f83', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('645fa664-9dd0-552c-8740-f523523ea1dd', 'approval.flows.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('645fa664-9dd0-552c-8740-f523523ea1dd', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('645fa664-9dd0-552c-8740-f523523ea1dd', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('52eecef9-305b-5235-a0d9-2203f84aff0a', 'approval.flows.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('52eecef9-305b-5235-a0d9-2203f84aff0a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('aee62640-1f21-5f32-995a-1e704394723e', 'approval.flows.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('aee62640-1f21-5f32-995a-1e704394723e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d39e9cbe-aeec-52e5-86ca-14ca840949d2', 'approval.flows.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d39e9cbe-aeec-52e5-86ca-14ca840949d2', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2afebae2-d68c-5457-9a0e-a1a3f79fe150', 'approval.instances.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2afebae2-d68c-5457-9a0e-a1a3f79fe150', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2afebae2-d68c-5457-9a0e-a1a3f79fe150', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('de87b5fa-c2e5-5ed6-9020-db262deba0ee', 'approval.instances.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('de87b5fa-c2e5-5ed6-9020-db262deba0ee', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('04799bf5-ae57-5d6d-b13e-cf3fae20bc0d', 'approval.instances.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('04799bf5-ae57-5d6d-b13e-cf3fae20bc0d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('af9bd05e-f475-5b62-a658-011652a7e347', 'approval.instances.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('af9bd05e-f475-5b62-a658-011652a7e347', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('265c4113-b9eb-5155-8fc5-b73e1cd0f530', 'payroll.runs.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('265c4113-b9eb-5155-8fc5-b73e1cd0f530', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('265c4113-b9eb-5155-8fc5-b73e1cd0f530', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('16f46b3a-f999-5435-95a2-d49691a8e9ab', 'payroll.runs.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('16f46b3a-f999-5435-95a2-d49691a8e9ab', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('64d811fe-76d3-5786-9611-34b9fa07293e', 'payroll.runs.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('64d811fe-76d3-5786-9611-34b9fa07293e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e653acff-bb3c-5217-9302-71a908066633', 'payroll.runs.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e653acff-bb3c-5217-9302-71a908066633', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('037676fb-fa9a-5bc9-a56c-a409db6d451a', 'payroll.periods.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('037676fb-fa9a-5bc9-a56c-a409db6d451a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('037676fb-fa9a-5bc9-a56c-a409db6d451a', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('31e5f3ac-961b-5351-aa8d-80613932ba0f', 'payroll.periods.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('31e5f3ac-961b-5351-aa8d-80613932ba0f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e5c02c82-9b5e-5d70-a67b-a6bc337b7eea', 'payroll.periods.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e5c02c82-9b5e-5d70-a67b-a6bc337b7eea', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ac9f3b64-21fb-5285-b62e-58c2ac585932', 'payroll.periods.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ac9f3b64-21fb-5285-b62e-58c2ac585932', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b8ae71e3-e98a-5fc9-a9d2-931e2ee57d64', 'payroll.salary-components.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b8ae71e3-e98a-5fc9-a9d2-931e2ee57d64', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b8ae71e3-e98a-5fc9-a9d2-931e2ee57d64', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('56db9d93-a63f-56ed-942f-1872118d3611', 'payroll.salary-components.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56db9d93-a63f-56ed-942f-1872118d3611', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4b887706-b787-51c3-9478-b874263dbc0e', 'payroll.salary-components.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4b887706-b787-51c3-9478-b874263dbc0e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('24e304c6-39f9-58da-b541-3a61ea172a0a', 'payroll.salary-components.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('24e304c6-39f9-58da-b541-3a61ea172a0a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2c7c0cd3-cf37-5d24-9939-31e775d6daa7', 'payroll.profiles.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2c7c0cd3-cf37-5d24-9939-31e775d6daa7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2c7c0cd3-cf37-5d24-9939-31e775d6daa7', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e3e274c3-a600-528a-914a-77009dd7b870', 'payroll.profiles.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e3e274c3-a600-528a-914a-77009dd7b870', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e01357b9-d44d-5a53-a8ea-b71c67423205', 'payroll.profiles.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e01357b9-d44d-5a53-a8ea-b71c67423205', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('85b67dbb-af54-5a99-88c5-6919cdd13bf2', 'payroll.profiles.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('85b67dbb-af54-5a99-88c5-6919cdd13bf2', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1a09b93b-229f-557d-a317-6a1cf2294bbc', 'payroll.bpjs-settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1a09b93b-229f-557d-a317-6a1cf2294bbc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1a09b93b-229f-557d-a317-6a1cf2294bbc', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4a34eb2d-8d75-5da9-8068-682da88d46c5', 'payroll.bpjs-settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4a34eb2d-8d75-5da9-8068-682da88d46c5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ddb3fec9-9430-527a-87e8-972c5af8641e', 'payroll.bpjs-settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ddb3fec9-9430-527a-87e8-972c5af8641e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b1eea49c-23eb-5550-8271-5fe312acaa20', 'payroll.bpjs-settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b1eea49c-23eb-5550-8271-5fe312acaa20', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2ae98b82-f8f6-52d0-ad66-333f6b0f7b2b', 'payroll.pph21-settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2ae98b82-f8f6-52d0-ad66-333f6b0f7b2b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2ae98b82-f8f6-52d0-ad66-333f6b0f7b2b', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0b2b73b6-8d70-5a51-ba7f-83bc09e9274e', 'payroll.pph21-settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0b2b73b6-8d70-5a51-ba7f-83bc09e9274e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6d8e2d09-731f-5bce-a35d-f56fd5d80692', 'payroll.pph21-settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6d8e2d09-731f-5bce-a35d-f56fd5d80692', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c8c5f8d2-5302-5634-92dd-f3711fa71299', 'payroll.pph21-settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c8c5f8d2-5302-5634-92dd-f3711fa71299', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c249d3ca-6d08-55a8-8744-5e6de91fc33f', 'payroll.ptkp-rates.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c249d3ca-6d08-55a8-8744-5e6de91fc33f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c249d3ca-6d08-55a8-8744-5e6de91fc33f', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('388cd215-9637-584f-8891-486193a66f01', 'payroll.ptkp-rates.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('388cd215-9637-584f-8891-486193a66f01', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7ef86ca7-ac63-5936-9e1c-0a76bbc2947a', 'payroll.ptkp-rates.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7ef86ca7-ac63-5936-9e1c-0a76bbc2947a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c4da3b1e-39f6-5ff5-b2ba-63aa922dbcf1', 'payroll.ptkp-rates.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c4da3b1e-39f6-5ff5-b2ba-63aa922dbcf1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('97fdcf92-070d-5fe3-b711-4987455381bf', 'payroll.tax-brackets.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('97fdcf92-070d-5fe3-b711-4987455381bf', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('97fdcf92-070d-5fe3-b711-4987455381bf', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('097a93f9-21ac-572c-bbde-7733bfa8ab5f', 'payroll.tax-brackets.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('097a93f9-21ac-572c-bbde-7733bfa8ab5f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0010d207-0b8d-520b-b286-02665d156913', 'payroll.tax-brackets.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0010d207-0b8d-520b-b286-02665d156913', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d8ef3c94-503a-58d6-8691-4580cd6abfa9', 'payroll.tax-brackets.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d8ef3c94-503a-58d6-8691-4580cd6abfa9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('04976b88-eed2-580d-a244-4e3bb1dddbe8', 'leave.dashboard.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('04976b88-eed2-580d-a244-4e3bb1dddbe8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('04976b88-eed2-580d-a244-4e3bb1dddbe8', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('eafe152e-65fd-52da-8af3-5aef7d6acb91', 'leave.dashboard.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('eafe152e-65fd-52da-8af3-5aef7d6acb91', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1db4e988-6cf6-5997-8e60-625700403e20', 'leave.dashboard.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1db4e988-6cf6-5997-8e60-625700403e20', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2c0f3ff8-2b3d-56fc-ac92-28b9ffb54a53', 'leave.dashboard.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2c0f3ff8-2b3d-56fc-ac92-28b9ffb54a53', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a3a01d1e-2b94-55a6-b1bd-2ba5f7bc1bc6', 'leave.requests.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a3a01d1e-2b94-55a6-b1bd-2ba5f7bc1bc6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a3a01d1e-2b94-55a6-b1bd-2ba5f7bc1bc6', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('297c9d5c-2cd8-53b3-9617-919322459e60', 'leave.requests.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('297c9d5c-2cd8-53b3-9617-919322459e60', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('716f7ea6-0e93-50f8-83d0-202de0c230b8', 'leave.requests.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('716f7ea6-0e93-50f8-83d0-202de0c230b8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e180bcac-8a70-5758-8f2e-e2cbc1d0b93b', 'leave.requests.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e180bcac-8a70-5758-8f2e-e2cbc1d0b93b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('96654182-bc67-594b-8467-e2a79ff28c17', 'leave.types.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('96654182-bc67-594b-8467-e2a79ff28c17', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('96654182-bc67-594b-8467-e2a79ff28c17', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('36f4e28d-cffa-534e-afac-5c65c0d98588', 'leave.types.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('36f4e28d-cffa-534e-afac-5c65c0d98588', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('01a2bbb5-403d-50aa-9a7f-496bccd7e465', 'leave.types.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('01a2bbb5-403d-50aa-9a7f-496bccd7e465', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('14cca6ea-9c5d-5b0e-83f3-6aa70af1f974', 'leave.types.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('14cca6ea-9c5d-5b0e-83f3-6aa70af1f974', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ff5c48d0-f69a-5b9d-9aad-0acfc95b6874', 'leave.balances.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ff5c48d0-f69a-5b9d-9aad-0acfc95b6874', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ff5c48d0-f69a-5b9d-9aad-0acfc95b6874', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9767c6cf-7fe1-5a70-a94e-bdb956c3bcae', 'leave.balances.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9767c6cf-7fe1-5a70-a94e-bdb956c3bcae', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('139a5f68-7222-5b6c-913e-f43925e647c7', 'leave.balances.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('139a5f68-7222-5b6c-913e-f43925e647c7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0231cf23-d482-530e-a4e7-d00b4687c173', 'leave.balances.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0231cf23-d482-530e-a4e7-d00b4687c173', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('207d3e01-43b1-594a-9acb-fe4783d50309', 'leave.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('207d3e01-43b1-594a-9acb-fe4783d50309', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('207d3e01-43b1-594a-9acb-fe4783d50309', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5d77618b-8c6b-57ab-b6eb-b85fdf9da153', 'leave.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5d77618b-8c6b-57ab-b6eb-b85fdf9da153', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('08e4707e-f44b-555d-9709-5971266d253f', 'leave.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('08e4707e-f44b-555d-9709-5971266d253f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d57c992d-0881-52ca-96d6-3a27cb0cb440', 'leave.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d57c992d-0881-52ca-96d6-3a27cb0cb440', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('08e9d74b-476c-5c9c-8c69-8094fb9c0a17', 'performance.evaluations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('08e9d74b-476c-5c9c-8c69-8094fb9c0a17', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('08e9d74b-476c-5c9c-8c69-8094fb9c0a17', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('960992ab-4cc4-5fdb-8c6e-cda605c1c34f', 'performance.evaluations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('960992ab-4cc4-5fdb-8c6e-cda605c1c34f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('265932c2-e660-5244-b269-496df073dd61', 'performance.evaluations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('265932c2-e660-5244-b269-496df073dd61', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0c1fc082-1d47-5a27-afc9-83f3997dbcb5', 'performance.evaluations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0c1fc082-1d47-5a27-afc9-83f3997dbcb5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2c0c7c23-244a-518d-a4b7-f6c00c7b6937', 'performance.templates.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2c0c7c23-244a-518d-a4b7-f6c00c7b6937', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2c0c7c23-244a-518d-a4b7-f6c00c7b6937', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('26716d43-3286-5e27-9f70-075deccda901', 'performance.templates.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('26716d43-3286-5e27-9f70-075deccda901', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('19a4fae6-50a4-5978-9f99-0a30b8fb4ba0', 'performance.templates.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('19a4fae6-50a4-5978-9f99-0a30b8fb4ba0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('66ca4635-d793-58e6-a935-d604396ed053', 'performance.templates.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('66ca4635-d793-58e6-a935-d604396ed053', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('52bdc408-0f28-5f57-b345-2e8b68b3b41b', 'performance.indicators.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('52bdc408-0f28-5f57-b345-2e8b68b3b41b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('52bdc408-0f28-5f57-b345-2e8b68b3b41b', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9ebf583a-585e-51fb-8000-c947c07d3606', 'performance.indicators.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9ebf583a-585e-51fb-8000-c947c07d3606', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('70eddcf6-792a-5cfe-a21f-5747c0e045fa', 'performance.indicators.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('70eddcf6-792a-5cfe-a21f-5747c0e045fa', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b15df6f9-d5ee-55d2-adf9-df0b6ca7e0c8', 'performance.indicators.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b15df6f9-d5ee-55d2-adf9-df0b6ca7e0c8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('60b04a7e-be9c-50fb-b37a-d6e637a96581', 'performance.periods.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('60b04a7e-be9c-50fb-b37a-d6e637a96581', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('60b04a7e-be9c-50fb-b37a-d6e637a96581', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ddd46246-e8fe-571c-a18f-9aaed45f5288', 'performance.periods.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ddd46246-e8fe-571c-a18f-9aaed45f5288', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9f1fae6d-58e4-5906-a3f8-4dc4027e3901', 'performance.periods.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9f1fae6d-58e4-5906-a3f8-4dc4027e3901', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b93cc631-2272-52e2-b499-382a8b4b5969', 'performance.periods.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b93cc631-2272-52e2-b499-382a8b4b5969', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0329587c-ccfd-5bc2-83e7-8adc6288ebd3', 'performance.perspectives.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0329587c-ccfd-5bc2-83e7-8adc6288ebd3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0329587c-ccfd-5bc2-83e7-8adc6288ebd3', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b0ee799c-d3ce-536d-aece-78a04a4d9ad5', 'performance.perspectives.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b0ee799c-d3ce-536d-aece-78a04a4d9ad5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7c533642-b3bf-5a03-ab1c-084bb5ce056f', 'performance.perspectives.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7c533642-b3bf-5a03-ab1c-084bb5ce056f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fb5c76f1-d845-5a8b-ae7c-03173f0fc59d', 'performance.perspectives.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fb5c76f1-d845-5a8b-ae7c-03173f0fc59d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f0fabdc6-cb16-52ce-9afc-c8b6372e9aba', 'recruitment.requisitions.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f0fabdc6-cb16-52ce-9afc-c8b6372e9aba', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f0fabdc6-cb16-52ce-9afc-c8b6372e9aba', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f880ea50-83f5-5f61-8593-c99320a2cb67', 'recruitment.requisitions.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f880ea50-83f5-5f61-8593-c99320a2cb67', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('23db20bf-c9a4-56f5-9387-07eb4a562b20', 'recruitment.requisitions.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('23db20bf-c9a4-56f5-9387-07eb4a562b20', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5c140fa7-c270-5f0e-b8b4-da55820f6aeb', 'recruitment.requisitions.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5c140fa7-c270-5f0e-b8b4-da55820f6aeb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3d60d36e-4c5e-53b8-af3c-46a2412fb6fc', 'recruitment.candidates.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3d60d36e-4c5e-53b8-af3c-46a2412fb6fc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3d60d36e-4c5e-53b8-af3c-46a2412fb6fc', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('424ae719-dd2f-59a9-99b3-58b291e486aa', 'recruitment.candidates.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('424ae719-dd2f-59a9-99b3-58b291e486aa', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9d9a2637-cfbd-5a9d-b520-0624b4185425', 'recruitment.candidates.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9d9a2637-cfbd-5a9d-b520-0624b4185425', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9aa57ba4-09c1-5f78-96fc-d5505f82e2c0', 'recruitment.candidates.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9aa57ba4-09c1-5f78-96fc-d5505f82e2c0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cc6460b5-8305-5da9-a5cc-e430001a5c0a', 'recruitment.applications.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc6460b5-8305-5da9-a5cc-e430001a5c0a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc6460b5-8305-5da9-a5cc-e430001a5c0a', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('82c725e6-b0bd-50eb-958b-bda4341c18e1', 'recruitment.applications.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('82c725e6-b0bd-50eb-958b-bda4341c18e1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d83a289e-5546-5b15-91d3-03fa0d5a1a8d', 'recruitment.applications.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d83a289e-5546-5b15-91d3-03fa0d5a1a8d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7b053a04-c63a-55c2-92a0-9b464f562377', 'recruitment.applications.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7b053a04-c63a-55c2-92a0-9b464f562377', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('676b46d1-372a-5420-ac06-26e7847076cf', 'recruitment.interviews.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('676b46d1-372a-5420-ac06-26e7847076cf', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('676b46d1-372a-5420-ac06-26e7847076cf', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('16d4dccf-839e-5a66-ae4c-7102f1f371b4', 'recruitment.interviews.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('16d4dccf-839e-5a66-ae4c-7102f1f371b4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2af43e23-4a4a-50cb-8e5c-8a1948d9e115', 'recruitment.interviews.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2af43e23-4a4a-50cb-8e5c-8a1948d9e115', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fb0c758c-2ee8-5c02-8b52-3d4ac37b6f45', 'recruitment.interviews.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fb0c758c-2ee8-5c02-8b52-3d4ac37b6f45', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1759f095-97e0-547a-9c68-a0ea7b88bb8d', 'recruitment.onboarding.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1759f095-97e0-547a-9c68-a0ea7b88bb8d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1759f095-97e0-547a-9c68-a0ea7b88bb8d', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f9f365aa-f92f-55fa-bb0e-e9bed50159fc', 'recruitment.onboarding.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f9f365aa-f92f-55fa-bb0e-e9bed50159fc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6bbe5881-705b-5d71-9bcf-8bc734c56e32', 'recruitment.onboarding.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6bbe5881-705b-5d71-9bcf-8bc734c56e32', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b695b9ed-7320-59f4-9961-d1e30b6bb574', 'recruitment.onboarding.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b695b9ed-7320-59f4-9961-d1e30b6bb574', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c20d4fca-6628-509c-a316-a11cfea680a4', 'reimbursement.requests.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c20d4fca-6628-509c-a316-a11cfea680a4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c20d4fca-6628-509c-a316-a11cfea680a4', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ac7a833a-de70-5b2c-8441-cea32ede0ff4', 'reimbursement.requests.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ac7a833a-de70-5b2c-8441-cea32ede0ff4', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('76e29558-1670-5ddc-9e3c-1bd60433cc04', 'reimbursement.requests.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('76e29558-1670-5ddc-9e3c-1bd60433cc04', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3323d68c-247f-593b-a7b0-4c086cd8c68d', 'reimbursement.requests.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3323d68c-247f-593b-a7b0-4c086cd8c68d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('83f901cb-a00d-5a1f-9b7f-8ec82f2771a7', 'reimbursement.types.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('83f901cb-a00d-5a1f-9b7f-8ec82f2771a7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('83f901cb-a00d-5a1f-9b7f-8ec82f2771a7', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ec1792bd-1b7b-5c7f-b147-c008d68db7a7', 'reimbursement.types.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ec1792bd-1b7b-5c7f-b147-c008d68db7a7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f7056c45-687c-5f4a-a1ba-d542bac556a3', 'reimbursement.types.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f7056c45-687c-5f4a-a1ba-d542bac556a3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c04406cb-2c39-5a01-a19e-588003396a19', 'reimbursement.types.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c04406cb-2c39-5a01-a19e-588003396a19', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('503f2491-b172-5e87-a511-b3c7d5b4e7b8', 'reimbursement.reports.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('503f2491-b172-5e87-a511-b3c7d5b4e7b8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('503f2491-b172-5e87-a511-b3c7d5b4e7b8', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('db595c2b-a129-5d48-ba71-078f3b4755cc', 'reimbursement.reports.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('db595c2b-a129-5d48-ba71-078f3b4755cc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ea171e41-8760-5ec8-8b72-c1ca8da20342', 'reimbursement.reports.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ea171e41-8760-5ec8-8b72-c1ca8da20342', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d36ec1a6-1d12-50f9-a792-64fc60065263', 'reimbursement.reports.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d36ec1a6-1d12-50f9-a792-64fc60065263', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('30e138b7-de50-55a6-9173-0ef15c3d6617', 'training.courses.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('30e138b7-de50-55a6-9173-0ef15c3d6617', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('30e138b7-de50-55a6-9173-0ef15c3d6617', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('767e6719-de5f-5eb4-8e09-0e909f43b8ab', 'training.courses.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('767e6719-de5f-5eb4-8e09-0e909f43b8ab', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('415855ee-2fd7-568f-9471-5f7cfc91e1e8', 'training.courses.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('415855ee-2fd7-568f-9471-5f7cfc91e1e8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d150bfcf-84b9-5dc5-a567-776d96aea35f', 'training.courses.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d150bfcf-84b9-5dc5-a567-776d96aea35f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ab361091-9b49-52ab-a745-929ca34fc997', 'training.categories.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ab361091-9b49-52ab-a745-929ca34fc997', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ab361091-9b49-52ab-a745-929ca34fc997', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7966895d-f85d-5a91-a20e-4e0bd5d2dd61', 'training.categories.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7966895d-f85d-5a91-a20e-4e0bd5d2dd61', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('abde48ee-a825-54ba-b2e4-66f611baab62', 'training.categories.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('abde48ee-a825-54ba-b2e4-66f611baab62', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c8c96ed3-6d2f-5323-9ba4-1582116e59b6', 'training.categories.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c8c96ed3-6d2f-5323-9ba4-1582116e59b6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0cc27862-fb2d-5274-8c8a-2b963224579c', 'training.providers.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0cc27862-fb2d-5274-8c8a-2b963224579c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0cc27862-fb2d-5274-8c8a-2b963224579c', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a4aa09f5-3360-5241-88c1-5871c4f508f3', 'training.providers.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a4aa09f5-3360-5241-88c1-5871c4f508f3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c1313d65-3c57-5aaf-b6ca-9bc936a958c3', 'training.providers.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c1313d65-3c57-5aaf-b6ca-9bc936a958c3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8357e918-1002-5fc8-ac27-780087c25ba3', 'training.providers.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8357e918-1002-5fc8-ac27-780087c25ba3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8f4d8784-4f46-5416-9986-d2d8f9b34a1c', 'training.trainers.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8f4d8784-4f46-5416-9986-d2d8f9b34a1c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8f4d8784-4f46-5416-9986-d2d8f9b34a1c', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5f3d7c09-dc1e-56b2-a960-a8876474e324', 'training.trainers.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5f3d7c09-dc1e-56b2-a960-a8876474e324', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('75800331-432c-5430-8174-bc22339dca54', 'training.trainers.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('75800331-432c-5430-8174-bc22339dca54', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('abd52b25-c377-597c-b754-bbcc7610b8f0', 'training.trainers.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('abd52b25-c377-597c-b754-bbcc7610b8f0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ba2230f4-1bb8-5f9a-8531-542b41077b83', 'training.sessions.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ba2230f4-1bb8-5f9a-8531-542b41077b83', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ba2230f4-1bb8-5f9a-8531-542b41077b83', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('19b3c7cc-2781-589b-b24d-d83df075bde8', 'training.sessions.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('19b3c7cc-2781-589b-b24d-d83df075bde8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0f7c574a-fd1c-50ea-a5e7-0d3b678b0a8e', 'training.sessions.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0f7c574a-fd1c-50ea-a5e7-0d3b678b0a8e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('69e15a7e-ef1f-57a9-b5ef-989f39e33411', 'training.sessions.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('69e15a7e-ef1f-57a9-b5ef-989f39e33411', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('98cedd7b-3ae9-5d18-97fe-0d86cba998c9', 'training.participants.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('98cedd7b-3ae9-5d18-97fe-0d86cba998c9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('98cedd7b-3ae9-5d18-97fe-0d86cba998c9', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('fd432ed0-9a10-51fb-bb9f-869a32d5aa93', 'training.participants.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('fd432ed0-9a10-51fb-bb9f-869a32d5aa93', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9efaa95d-5cbc-52c2-8acd-492f473ed6c7', 'training.participants.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9efaa95d-5cbc-52c2-8acd-492f473ed6c7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e1ab24fa-65e0-5586-9d60-0b7df44ae99a', 'training.participants.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e1ab24fa-65e0-5586-9d60-0b7df44ae99a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4c31673f-7659-56ac-a315-23f57e0220fd', 'training.planning.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4c31673f-7659-56ac-a315-23f57e0220fd', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4c31673f-7659-56ac-a315-23f57e0220fd', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('987d998a-8db4-51f3-a91c-aee8c5592909', 'training.planning.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('987d998a-8db4-51f3-a91c-aee8c5592909', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('49ec0dff-c5ea-5ad7-b12c-4a9e9d148e7c', 'training.planning.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('49ec0dff-c5ea-5ad7-b12c-4a9e9d148e7c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('082a557c-3037-5762-a313-918960d21560', 'training.planning.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('082a557c-3037-5762-a313-918960d21560', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('3d368e5e-a99e-5643-869c-9bc693fd8362', 'training.requests.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3d368e5e-a99e-5643-869c-9bc693fd8362', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('3d368e5e-a99e-5643-869c-9bc693fd8362', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f52c5e1c-e27f-5702-b26f-26be4e7d41c0', 'training.requests.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f52c5e1c-e27f-5702-b26f-26be4e7d41c0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('8793684f-1d59-56cf-9905-e03f21d6c008', 'training.requests.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('8793684f-1d59-56cf-9905-e03f21d6c008', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('cc14e0f3-7e37-526d-a70e-c736479abf41', 'training.requests.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('cc14e0f3-7e37-526d-a70e-c736479abf41', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ac5624ad-b89c-54ae-8eb1-eb6604752a32', 'training.needs.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ac5624ad-b89c-54ae-8eb1-eb6604752a32', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ac5624ad-b89c-54ae-8eb1-eb6604752a32', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c3aefd20-c54f-5614-a6c5-7bfe7b5d7203', 'training.needs.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c3aefd20-c54f-5614-a6c5-7bfe7b5d7203', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('40967a83-e249-5ab2-9089-527fa0263708', 'training.needs.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('40967a83-e249-5ab2-9089-527fa0263708', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6c20409e-5850-5a37-a9ce-d456324b8bbb', 'training.needs.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6c20409e-5850-5a37-a9ce-d456324b8bbb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('56808200-48f3-5e4b-ab88-7585684d3d31', 'training.certificates.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56808200-48f3-5e4b-ab88-7585684d3d31', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56808200-48f3-5e4b-ab88-7585684d3d31', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1fd82008-de16-51d6-89dc-2e9df4084402', 'training.certificates.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1fd82008-de16-51d6-89dc-2e9df4084402', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a109c5f3-fd1e-5256-a2f8-72288ee28e28', 'training.certificates.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a109c5f3-fd1e-5256-a2f8-72288ee28e28', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ec1d7fe5-888a-5ad2-a089-268887a8edcf', 'training.certificates.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ec1d7fe5-888a-5ad2-a089-268887a8edcf', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6a56f4da-af1a-5a67-b691-ee651981f4e0', 'training.history.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6a56f4da-af1a-5a67-b691-ee651981f4e0', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6a56f4da-af1a-5a67-b691-ee651981f4e0', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('776b9c3b-d56e-5a98-93f6-4bf98a3c711a', 'training.history.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('776b9c3b-d56e-5a98-93f6-4bf98a3c711a', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f848c54e-573c-50ed-9f01-50651d99ab92', 'training.history.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f848c54e-573c-50ed-9f01-50651d99ab92', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2ca3dfcf-40d0-5110-b6f4-9f799a40c1e3', 'training.history.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2ca3dfcf-40d0-5110-b6f4-9f799a40c1e3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', 'training.reports.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('10972ea0-c8c9-5df3-91c9-07bf874d59e1', 'training.reports.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('10972ea0-c8c9-5df3-91c9-07bf874d59e1', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('d09d7b99-cd66-5f5f-b2c0-8092924fd83f', 'training.reports.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('d09d7b99-cd66-5f5f-b2c0-8092924fd83f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6b140aea-d616-5461-8fa4-130115d6739d', 'training.reports.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6b140aea-d616-5461-8fa4-130115d6739d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b4271c32-c0fb-50ef-8f5e-cdb98151ae96', 'rbac.roles.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b4271c32-c0fb-50ef-8f5e-cdb98151ae96', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b4271c32-c0fb-50ef-8f5e-cdb98151ae96', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('534e842f-30c8-5b90-9b0e-4411488aa471', 'rbac.roles.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('534e842f-30c8-5b90-9b0e-4411488aa471', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1302d874-847e-5bae-8a65-f491dc5e7eb6', 'rbac.roles.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1302d874-847e-5bae-8a65-f491dc5e7eb6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a16e766d-9f94-52ce-a83a-9c00d0093e06', 'rbac.roles.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a16e766d-9f94-52ce-a83a-9c00d0093e06', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c59b2317-1f08-5a4e-9034-7bce4c501f43', 'jobmanagement.setting.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c59b2317-1f08-5a4e-9034-7bce4c501f43', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c59b2317-1f08-5a4e-9034-7bce4c501f43', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ec82462b-00e4-516e-baca-677360931ddf', 'jobmanagement.setting.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ec82462b-00e4-516e-baca-677360931ddf', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4642a76a-c29c-56d6-90f6-e56ac506f1eb', 'jobmanagement.setting.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4642a76a-c29c-56d6-90f6-e56ac506f1eb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('97ea7d47-d665-5bda-8ad5-11329f878f80', 'jobmanagement.setting.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('97ea7d47-d665-5bda-8ad5-11329f878f80', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1e34a394-b875-545e-9bcd-ae0d0e6fa386', 'jobmanagement.assessment.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1e34a394-b875-545e-9bcd-ae0d0e6fa386', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1e34a394-b875-545e-9bcd-ae0d0e6fa386', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ff80361a-85e6-58fe-894b-4813ecd90526', 'jobmanagement.assessment.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ff80361a-85e6-58fe-894b-4813ecd90526', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('9825b818-8668-571a-aee0-180c9979e36d', 'jobmanagement.assessment.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('9825b818-8668-571a-aee0-180c9979e36d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bfb9a601-2941-5b73-93eb-75775cf5a3fe', 'jobmanagement.assessment.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bfb9a601-2941-5b73-93eb-75775cf5a3fe', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('eb66b714-ef2e-5f12-b8a3-e64aa91634a8', 'competency.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('eb66b714-ef2e-5f12-b8a3-e64aa91634a8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('eb66b714-ef2e-5f12-b8a3-e64aa91634a8', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('89c278fd-7a5c-52d9-9369-a3f58ec41ad3', 'competency.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('89c278fd-7a5c-52d9-9369-a3f58ec41ad3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5765f5de-6b72-5045-975b-8423ba0dbb37', 'competency.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5765f5de-6b72-5045-975b-8423ba0dbb37', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('18afb851-661d-5c05-9c5a-664a3be5347f', 'competency.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('18afb851-661d-5c05-9c5a-664a3be5347f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('94fe83aa-cf27-58fc-b8b5-7abc03f1563f', 'competency.assessment.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('94fe83aa-cf27-58fc-b8b5-7abc03f1563f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('94fe83aa-cf27-58fc-b8b5-7abc03f1563f', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e582546c-edd7-59e5-b584-af67173757ce', 'competency.assessment.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e582546c-edd7-59e5-b584-af67173757ce', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6b95da8e-0260-539d-b3fb-0495fdf4b044', 'competency.assessment.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6b95da8e-0260-539d-b3fb-0495fdf4b044', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a53de392-b016-5eea-8e4e-277f72e8c432', 'competency.assessment.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a53de392-b016-5eea-8e4e-277f72e8c432', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0029d707-f241-57e9-91d1-3dc6d44f24c9', 'competency.report.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0029d707-f241-57e9-91d1-3dc6d44f24c9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0029d707-f241-57e9-91d1-3dc6d44f24c9', '3e562937-d5a1-543a-b1c8-af2f447500a4');
