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
-- performance disederhanakan jadi 3 submenu: settings (templates,
-- indicators, periods, perspectives), operational (kpi, okr, evaluations),
-- & report (kosong, disiapkan untuk nanti — view only).
-- recruitment disederhanakan jadi 2 submenu: pipeline (requisitions,
-- applications, candidates, internal-candidates, assessments, offers),
-- onboarding (tetap). interviews dihapus (tidak dipakai).
-- training disederhanakan jadi 4 submenu (maksimal): settings (courses,
-- categories, providers, trainers), operations (sessions, participants,
-- planning, requests, needs), records (certificates, history), &
-- reports (view only; kept plural to avoid colliding with the existing
-- module-level action "report.view").

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

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5e0ab851-d7ee-5d66-8aa9-6411001e0581', 'performance.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5e0ab851-d7ee-5d66-8aa9-6411001e0581', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5e0ab851-d7ee-5d66-8aa9-6411001e0581', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('dcb055fc-0ea8-5688-87f8-44aa6333fc6b', 'performance.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('dcb055fc-0ea8-5688-87f8-44aa6333fc6b', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('56ad343a-2691-577c-8410-5db45d05c7f9', 'performance.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('56ad343a-2691-577c-8410-5db45d05c7f9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('df16dcce-765c-55b8-8e47-3637899cf680', 'performance.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('df16dcce-765c-55b8-8e47-3637899cf680', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('bf51f4af-691e-51c9-a613-7f3050e7a67c', 'performance.operational.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bf51f4af-691e-51c9-a613-7f3050e7a67c', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('bf51f4af-691e-51c9-a613-7f3050e7a67c', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('47d024ba-0fea-5f7f-afad-3264b6b4e1c6', 'performance.operational.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('47d024ba-0fea-5f7f-afad-3264b6b4e1c6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('6906a961-b805-5950-bdac-df3dd4292c77', 'performance.operational.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('6906a961-b805-5950-bdac-df3dd4292c77', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('4dccf52d-c457-53b3-997d-ef6ed91f3f2f', 'performance.operational.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('4dccf52d-c457-53b3-997d-ef6ed91f3f2f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('5b8f0eea-df9e-5ad6-94fc-5d938734f0ac', 'performance.report.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5b8f0eea-df9e-5ad6-94fc-5d938734f0ac', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('5b8f0eea-df9e-5ad6-94fc-5d938734f0ac', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('aa780452-d497-5ce5-9f6b-7b7b0d9a5527', 'recruitment.pipeline.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('aa780452-d497-5ce5-9f6b-7b7b0d9a5527', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('aa780452-d497-5ce5-9f6b-7b7b0d9a5527', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('088d9aef-b462-55fa-acb0-ca4e1335bad7', 'recruitment.pipeline.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('088d9aef-b462-55fa-acb0-ca4e1335bad7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ff6fb33c-a7c2-5c0c-9fca-85134faeb0b3', 'recruitment.pipeline.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ff6fb33c-a7c2-5c0c-9fca-85134faeb0b3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7d2365b5-145d-54df-b8bd-702fc5e45f5e', 'recruitment.pipeline.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7d2365b5-145d-54df-b8bd-702fc5e45f5e', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('687f578a-e8fc-5833-a546-b466b41926ac', 'training.settings.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('687f578a-e8fc-5833-a546-b466b41926ac', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('687f578a-e8fc-5833-a546-b466b41926ac', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('802810c7-70bd-5c10-a1c8-7e577023f9b7', 'training.settings.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('802810c7-70bd-5c10-a1c8-7e577023f9b7', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('b9dbb900-69f4-554c-8b10-712027cdd552', 'training.settings.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('b9dbb900-69f4-554c-8b10-712027cdd552', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0378521f-64bf-53a6-80ca-3f2886f48b42', 'training.settings.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0378521f-64bf-53a6-80ca-3f2886f48b42', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('ed4714c3-aadd-5e38-be92-7988005c0f55', 'training.operations.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ed4714c3-aadd-5e38-be92-7988005c0f55', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('ed4714c3-aadd-5e38-be92-7988005c0f55', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1547155c-003e-53ed-9759-efbf52f4f444', 'training.operations.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1547155c-003e-53ed-9759-efbf52f4f444', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('7e85a747-cd2d-5084-9a5b-1b51c2f244a3', 'training.operations.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('7e85a747-cd2d-5084-9a5b-1b51c2f244a3', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e641da29-f2e4-5f4f-b5f0-d696733a7e95', 'training.operations.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e641da29-f2e4-5f4f-b5f0-d696733a7e95', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('1549f78b-891d-567d-9332-9a73056e8ba6', 'training.records.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1549f78b-891d-567d-9332-9a73056e8ba6', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('1549f78b-891d-567d-9332-9a73056e8ba6', '3e562937-d5a1-543a-b1c8-af2f447500a4');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('69b8d4b0-7434-54f5-b451-c31ee029c58f', 'training.records.create', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('69b8d4b0-7434-54f5-b451-c31ee029c58f', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('e38969df-6c4d-50ef-a447-7d322874c9d9', 'training.records.update', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('e38969df-6c4d-50ef-a447-7d322874c9d9', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('20cf019e-1fcf-54ab-abd5-de9163dbbeba', 'training.records.delete', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('20cf019e-1fcf-54ab-abd5-de9163dbbeba', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', 'training.reports.view', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', 'ea1dcc10-3eb8-52c4-bd7c-bfb43e56d345');

INSERT IGNORE INTO role_has_permissions (permission_id, role_id)
VALUES ('0a99150f-5693-5702-8934-2d35c5459eab', '3e562937-d5a1-543a-b1c8-af2f447500a4');
