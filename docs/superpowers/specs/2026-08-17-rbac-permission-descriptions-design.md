# RBAC Permission Descriptions (Bilingual) — Design

## Problem

The RBAC permissions page (`RolePermissions.vue`) shows a grid of
resource/submenu rows × action (view/create/update/delete) columns, but
each row is labeled only by a short slug-derived name (e.g. "Settings",
"Operations") with no explanation of what that permission actually
grants. Admins configuring roles have no way to tell what a given
row means without guessing from the label alone.

Separately, the existing `rbac.submenus.<resource>.<submenu>` locale
keys (in `frontend/tenant/src/locales/{en,id}.json`) are stale — they
still reference submenu names from before this session's RBAC
consolidation work (e.g. old attendance submenus like `dashboard`,
`shifts`, `corrections` instead of the current `settings`,
`operations`, `report`). Most current submenu rows fall back to raw
slugs because their locale key doesn't exist.

## Goals

- Add a short, bilingual (Indonesian + English) description under each
  resource/submenu row in the RBAC permissions grid.
- Fix `rbac.submenus.*` to match the current (post-consolidation)
  submenu structure across all resources.
- No backend, database, or API changes — this is a frontend-only,
  additive change to the existing i18n-based labeling pattern already
  used for `rbac.modules.*` and `rbac.submenus.*`.

## Non-goals

- Per-action (view/create/update/delete) descriptions — one short
  description per row (resource+submenu, or resource alone for
  module-level rows) is sufficient per this session's decision.
- Storing descriptions in the database — would require a schema
  migration and re-seeding every tenant for content that changes
  independently of tenant data; the existing `rbac.modules.*` /
  `rbac.submenus.*` precedent already keeps this kind of label in
  frontend i18n, not the DB.

## Approach

Add a new locale key namespace, `rbac.descriptions.<resource>.<submenu>`,
parallel to the existing `rbac.submenus.<resource>.<submenu>`
structure. For resources with no submenus (module-level only, e.g.
`organization`, `setting`, `employee`), use
`rbac.descriptions.<resource>._module`.

## 1. Locale key structure

For every resource currently defined in `tenantRBACResources()`
(`backend/internal/pkg/tenantseed/seed_rbac.go`), add:

- `rbac.descriptions.<resource>._module` — one short description for
  the module-level permission row, always present.
- `rbac.descriptions.<resource>.<submenu>` — one per submenu, for
  resources that have submenus.

Descriptions are one short sentence each (roughly 3-8 words), e.g.:
- `rbac.descriptions.attendance.settings` → ID: "Konfigurasi shift, lokasi, jam kerja" / EN: "Shift, location, working hours configuration"
- `rbac.descriptions.attendance.operations` → ID: "Aktivitas presensi harian" / EN: "Daily attendance activity"

## 2. Fix `rbac.submenus.*`

Rewrite the `rbac.submenus.<resource>` blocks in both `en.json` and
`id.json` to exactly match the current submenu list per resource (as
defined in `tenantRBACResources()` today, post-consolidation):

| Resource | Current submenus |
|---|---|
| jobmanagement | setting, assessment |
| competency | settings, assessment, report |
| attendance | settings, operations, report |
| approval | settings, operations |
| leave | settings, operations |
| performance | settings, operational, report |
| recruitment | pipeline, onboarding |
| reimbursement | requests, types, reports |
| training | settings, operations, records, reports |
| rbac | roles |

All other resources (organization, employee, jobmanagement's parent,
employeemovement, useraccount, payroll, workforceintelligence,
careerintelligence, setting, notification) have zero submenus — no
`rbac.submenus.<resource>` entry needed for them.

## 3. `RolePermissions.vue` change

Below each row's existing label (module label via `moduleLabel()`, or
submenu label via `submenuLabel()`), add a small muted text line:

```html
<span class="text-xs text-gray-400 dark:text-gray-500 block">
  {{ descriptionLabel(resource, submenu) }}
</span>
```

New helper `descriptionLabel(resource, submenu)`: looks up
`rbac.descriptions.<resource>.<submenu || '_module'>`, returns empty
string if the key doesn't resolve (Vue i18n's `t()` returns the key
itself when missing — guard against that by checking key existence, so
an unlabeled row shows nothing rather than a literal key string).

## Testing

- Manual: open the RBAC permissions page in both language settings (ID
  and EN), confirm every row shows a description, confirm submenu
  labels no longer fall back to raw slugs.
- No automated tests — this is static locale content + a template
  change with no business logic to unit test.
