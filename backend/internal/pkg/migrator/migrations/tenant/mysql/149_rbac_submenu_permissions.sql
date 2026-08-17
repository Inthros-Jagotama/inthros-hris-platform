-- 149_rbac_submenu_permissions.sql
-- Seed permission level-submenu ("resource.submenu.action") untuk tenant existing.
-- ID deterministik sama persis dengan codeToUUID di SeedTenantRBAC, jadi aman
-- dijalankan berulang (ON CONFLICT / INSERT IGNORE).
-- Admin mendapat semua action; Employee hanya view.
-- Catatan: organization, setting, employee, & employeemovement TIDAK punya
-- permission level-submenu —
-- cukup module-level (organization.view/create/update/delete, dst).
-- payroll juga TIDAK punya submenu — halaman konfigurasi payroll (periods,
-- salary-components, bpjs-settings, pph21-settings, salary-structure)
-- sudah ke-gate lewat menu Settings global (setting.view).
-- jobmanagement disederhanakan jadi 2 submenu: setting (titles, objectives,
-- identifications, responsibilities, authorities, working-conditions,
-- competencies) & assessment (values, scores).
-- competency disederhanakan jadi 3 submenu: settings (competencies, values,
-- indicators, templates, events, raters), assessment (my-assessments,
-- manager-assessments, scores), & report (results, reports — view only).
-- attendance disederhanakan jadi 3 submenu: settings (shifts, employee-shifts,
-- locations, exempt-positions, settings, admin), operations (dashboard,
-- schedules, events, sessions, overtime, corrections, business-travel),
-- & report (reports — view only).
-- approval disederhanakan jadi 2 submenu: settings (flows), operations
-- (tasks, instances).
-- leave disederhanakan jadi 2 submenu: settings (settings, types,
-- accrual-policies, reasons), operations (dashboard, requests, balances).

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
VALUES ('2a40f9ea-06a3-52a4-95b6-b09dad649f37', 'attendance.operations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2a40f9ea-06a3-52a4-95b6-b09dad649f37', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2a40f9ea-06a3-52a4-95b6-b09dad649f37', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('86dc438a-696c-5d81-b03b-ffdd8f60a481', 'attendance.operations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('86dc438a-696c-5d81-b03b-ffdd8f60a481', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bd7ccd79-1978-5820-9901-b2b552ec731b', 'attendance.operations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bd7ccd79-1978-5820-9901-b2b552ec731b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('49b4bc4d-42aa-5166-b66d-72e91147ec6d', 'attendance.operations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('49b4bc4d-42aa-5166-b66d-72e91147ec6d', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a61c3b32-c8a1-5dd5-a552-cca4b7e5d779', 'attendance.report.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a61c3b32-c8a1-5dd5-a552-cca4b7e5d779', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a61c3b32-c8a1-5dd5-a552-cca4b7e5d779', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('de15d7fc-db5b-5a4e-8dd7-174d1ce6dc87', 'approval.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('de15d7fc-db5b-5a4e-8dd7-174d1ce6dc87', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('de15d7fc-db5b-5a4e-8dd7-174d1ce6dc87', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a0b7b2d9-56a1-5c25-927f-0ecdbbe5f18f', 'approval.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a0b7b2d9-56a1-5c25-927f-0ecdbbe5f18f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('f82a7893-7b72-5324-84f2-d86bb2786bcb', 'approval.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('f82a7893-7b72-5324-84f2-d86bb2786bcb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('43a619b9-d420-5765-b872-a8c76ee0cd3f', 'approval.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('43a619b9-d420-5765-b872-a8c76ee0cd3f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ef03cbb7-50ef-5905-83f2-41d7263e1469', 'approval.operations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ef03cbb7-50ef-5905-83f2-41d7263e1469', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ef03cbb7-50ef-5905-83f2-41d7263e1469', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('182c2097-5b6f-5c50-9360-556ce4a26b7c', 'approval.operations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('182c2097-5b6f-5c50-9360-556ce4a26b7c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('2ff128ca-807b-50b5-b2d2-512bbeaf83a5', 'approval.operations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('2ff128ca-807b-50b5-b2d2-512bbeaf83a5', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('364011da-ce28-5a1b-b209-93c83ae03308', 'approval.operations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('364011da-ce28-5a1b-b209-93c83ae03308', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('17622ba1-8e7e-5510-94a5-15ca11099960', 'leave.operations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('17622ba1-8e7e-5510-94a5-15ca11099960', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('17622ba1-8e7e-5510-94a5-15ca11099960', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('51add501-6150-58db-a7e4-46480ee68acb', 'leave.operations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('51add501-6150-58db-a7e4-46480ee68acb', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('c63f59e1-e202-550e-b52b-61331878d4c8', 'leave.operations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('c63f59e1-e202-550e-b52b-61331878d4c8', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('a1c67c35-b973-56f5-b405-bcd9a8aeb1cc', 'leave.operations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('a1c67c35-b973-56f5-b405-bcd9a8aeb1cc', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');
