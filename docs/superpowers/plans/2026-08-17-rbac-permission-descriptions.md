# RBAC Permission Descriptions (Bilingual) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a short bilingual (ID/EN) description under every row of the RBAC permissions grid, and fix the stale `rbac.submenus.*` locale keys to match the current post-consolidation submenu structure.

**Architecture:** Frontend-only i18n change. Two new/fixed locale namespaces (`rbac.submenus.*` rewritten, `rbac.descriptions.*` added) in `frontend/tenant/src/locales/{en,id}.json`, plus a new `descriptionLabel()` helper and template line in `RolePermissions.vue`, mirroring the existing `moduleLabel()`/`submenuLabel()` pattern exactly.

**Tech Stack:** Vue 3, `@/composables/useI18n`, JSON locale files.

**Spec:** `docs/superpowers/specs/2026-08-17-rbac-permission-descriptions-design.md`

## Global Constraints

- No backend, database, or API changes — locale JSON + one Vue template/script edit only.
- One short description per row (resource+submenu, or resource alone for module-level rows) — no per-action descriptions.
- Descriptions are one short sentence each, roughly 3-8 words.
- `rbac.submenus.<resource>` blocks must exactly match the current submenu list per resource (table in spec section 2) — no stale/extra submenu keys.
- Resources with zero submenus get no `rbac.submenus.<resource>` entry, but do get `rbac.descriptions.<resource>._module`.
- Missing `rbac.descriptions.*` key must render as empty string in the UI, never as the literal key.

---

## Current State (read before editing)

`tenantRBACResources()` in `backend/internal/pkg/tenantseed/seed_rbac.go` defines these 18 resources and their current submenus (this is the authoritative structure — matches spec section 2):

| Resource | Submenus |
|---|---|
| organization | (none) |
| employee | (none) |
| jobmanagement | setting, assessment |
| competency | settings, assessment, report |
| employeemovement | (none) |
| useraccount | (none) |
| attendance | settings, operations, report |
| approval | settings, operations |
| payroll | (none) |
| leave | settings, operations |
| performance | settings, operational, report |
| recruitment | pipeline, onboarding |
| reimbursement | requests, types, reports |
| training | settings, operations, records, reports |
| workforceintelligence | (none) |
| careerintelligence | (none) |
| setting | (none) |
| rbac | roles |
| notification | (none) |

`frontend/tenant/src/locales/en.json` and `id.json` currently have, inside the top-level `rbac` object (both files, same line range ~3753-3936):
- `rbac.submenus.*` (lines 3767-3897) — **stale**: still keyed on pre-consolidation submenu names (e.g. `organization.tree/zones/job-families/positions`, `attendance.dashboard/shifts/schedules/events/overtime/locations/business-travel/settings`, `payroll.runs/periods/salary-components/...`, `training.courses/categories/providers/...`, etc.) — none of these match the current submenu table above except a few incidental overlaps.
- `rbac.modules.*` (lines 3899-3916 en, 3899-3916 id) — **not stale**, already keyed by resource only, one label per resource. Do not touch. Note it's also missing `workforceintelligence` and `careerintelligence` keys — out of scope for this plan (spec doesn't ask for it; `moduleLabel()` already falls back to the raw slug when a key is missing, so this is a pre-existing, harmless gap).

`frontend/tenant/src/views/settings/RolePermissions.vue` currently has (around lines 119-129):
```js
function moduleLabel(resource) {
  const key = `rbac.modules.${resource}`
  return t(key) !== key ? t(key) : resource
}

function submenuLabel(resource, submenu) {
  const key = `rbac.submenus.${resource}.${submenu}`
  return t(key) !== key ? t(key) : submenu
}
```
and the submenu column template (lines 56-66):
```html
<Column field="submenu" :header="t('rbac.submenu')" style="width:220px">
  <template #body="{ data }">
    <span
      class="text-gray-800 dark:text-gray-100"
      :class="data.submenu ? 'font-normal' : 'font-semibold'"
    >
      {{ data.submenu ? submenuLabel(data.resource, data.submenu) : t('rbac.module_level') }}
    </span>
  </template>
</Column>
```

---

## Task 1: Rewrite `rbac.submenus.*` to match current submenu structure

**Files:**
- Modify: `frontend/tenant/src/locales/en.json:3767-3897` (the `"submenus": { ... }` block inside `rbac`)
- Modify: `frontend/tenant/src/locales/id.json:3767-3897` (same block)

**Interfaces:**
- Consumes: nothing (pure data)
- Produces: `rbac.submenus.<resource>.<submenu>` keys used by `submenuLabel()` in `RolePermissions.vue` (Task 3, unchanged function) and by Task 2's description work as the reference for which submenus exist.

- [ ] **Step 1: Replace the `"submenus"` object in `en.json`**

Replace the entire block from `"submenus": {` (line 3767) through its matching closing `},` (line 3898) with:

```json
    "submenus": {
      "jobmanagement": {
        "setting": "Settings",
        "assessment": "Assessment"
      },
      "competency": {
        "settings": "Settings",
        "assessment": "Assessment",
        "report": "Reports"
      },
      "attendance": {
        "settings": "Settings",
        "operations": "Operations",
        "report": "Reports"
      },
      "approval": {
        "settings": "Settings",
        "operations": "Operations"
      },
      "leave": {
        "settings": "Settings",
        "operations": "Operations"
      },
      "performance": {
        "settings": "Settings",
        "operational": "Operational",
        "report": "Reports"
      },
      "recruitment": {
        "pipeline": "Pipeline",
        "onboarding": "Onboarding"
      },
      "reimbursement": {
        "requests": "Requests",
        "types": "Types",
        "reports": "Reports"
      },
      "training": {
        "settings": "Settings",
        "operations": "Operations",
        "records": "Records",
        "reports": "Reports"
      },
      "rbac": {
        "roles": "Roles"
      }
    },
```

- [ ] **Step 2: Replace the `"submenus"` object in `id.json`**

Replace the same block (lines 3767-3898 in `id.json`) with:

```json
    "submenus": {
      "jobmanagement": {
        "setting": "Pengaturan",
        "assessment": "Assessment"
      },
      "competency": {
        "settings": "Pengaturan",
        "assessment": "Assessment",
        "report": "Laporan"
      },
      "attendance": {
        "settings": "Pengaturan",
        "operations": "Operasional",
        "report": "Laporan"
      },
      "approval": {
        "settings": "Pengaturan",
        "operations": "Operasional"
      },
      "leave": {
        "settings": "Pengaturan",
        "operations": "Operasional"
      },
      "performance": {
        "settings": "Pengaturan",
        "operational": "Operasional",
        "report": "Laporan"
      },
      "recruitment": {
        "pipeline": "Pipeline",
        "onboarding": "Onboarding"
      },
      "reimbursement": {
        "requests": "Permohonan",
        "types": "Tipe",
        "reports": "Laporan"
      },
      "training": {
        "settings": "Pengaturan",
        "operations": "Operasional",
        "records": "Rekam Data",
        "reports": "Laporan"
      },
      "rbac": {
        "roles": "Roles"
      }
    },
```

- [ ] **Step 3: Verify JSON validity**

Run: `node -e "JSON.parse(require('fs').readFileSync('frontend/tenant/src/locales/en.json','utf8')); JSON.parse(require('fs').readFileSync('frontend/tenant/src/locales/id.json','utf8')); console.log('OK')"`
Expected: `OK` printed, no parse error.

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/locales/en.json frontend/tenant/src/locales/id.json
git commit -m "fix(rbac): rewrite rbac.submenus.* locale keys to match current submenu structure"
```

---

## Task 2: Add `rbac.descriptions.*` locale keys

**Files:**
- Modify: `frontend/tenant/src/locales/en.json` (add new `"descriptions"` object inside `rbac`, immediately after the `"submenus"` block edited in Task 1)
- Modify: `frontend/tenant/src/locales/id.json` (same)

**Interfaces:**
- Consumes: nothing (pure data)
- Produces: `rbac.descriptions.<resource>.<submenu|_module>` keys, consumed by `descriptionLabel()` added in Task 3.

- [ ] **Step 1: Insert `"descriptions"` object into `en.json`**

Immediately after the `"submenus": { ... },` block from Task 1 (so right before the existing `"modules": {` key), insert:

```json
    "descriptions": {
      "organization": { "_module": "Company org structure, zones, positions" },
      "employee": { "_module": "Employee records and profiles" },
      "jobmanagement": {
        "_module": "Job titles, values, competency mapping",
        "setting": "Job attribute configuration",
        "assessment": "Job scoring and evaluation"
      },
      "competency": {
        "_module": "Competency framework and scoring",
        "settings": "Competency and rating scale setup",
        "assessment": "Assessment templates and events",
        "report": "Competency result reports"
      },
      "employeemovement": { "_module": "Transfers, promotions, contract changes" },
      "useraccount": { "_module": "User login accounts for employees" },
      "attendance": {
        "_module": "Attendance and time tracking",
        "settings": "Shift, location, working hours configuration",
        "operations": "Daily attendance activity",
        "report": "Attendance summary reports"
      },
      "approval": {
        "_module": "Approval flows and requests",
        "settings": "Approval flow configuration",
        "operations": "Pending and processed approvals"
      },
      "payroll": { "_module": "Payroll processing and runs" },
      "leave": {
        "_module": "Leave requests and balances",
        "settings": "Leave type and policy configuration",
        "operations": "Leave requests and balances"
      },
      "performance": {
        "_module": "Performance evaluation and KPIs",
        "settings": "Evaluation template configuration",
        "operational": "KPI, OKR, and evaluation activity",
        "report": "Performance result reports"
      },
      "recruitment": {
        "_module": "Hiring and candidate management",
        "pipeline": "Requisitions, candidates, interviews",
        "onboarding": "New hire onboarding"
      },
      "reimbursement": {
        "_module": "Expense reimbursement claims",
        "requests": "Employee reimbursement claims",
        "types": "Reimbursement category configuration",
        "reports": "Reimbursement summary reports"
      },
      "training": {
        "_module": "Training courses and certifications",
        "settings": "Course, category, provider setup",
        "operations": "Sessions, participants, attendance",
        "records": "Certificates and training history",
        "reports": "Training summary reports"
      },
      "workforceintelligence": { "_module": "Workforce analytics and insights" },
      "careerintelligence": { "_module": "Career path and succession insights" },
      "setting": { "_module": "System-wide reference data configuration" },
      "rbac": {
        "_module": "Roles and permission management",
        "roles": "Role definitions and assignments"
      },
      "notification": { "_module": "System notifications" }
    },
```

- [ ] **Step 2: Insert `"descriptions"` object into `id.json`**

At the same insertion point in `id.json` (immediately after Task 1's `"submenus"` block, before `"modules": {`), insert:

```json
    "descriptions": {
      "organization": { "_module": "Struktur organisasi, zona, posisi" },
      "employee": { "_module": "Data dan profil karyawan" },
      "jobmanagement": {
        "_module": "Jabatan, nilai, pemetaan kompetensi",
        "setting": "Konfigurasi atribut jabatan",
        "assessment": "Penilaian dan skor jabatan"
      },
      "competency": {
        "_module": "Kerangka dan penilaian kompetensi",
        "settings": "Konfigurasi kompetensi dan skala",
        "assessment": "Template dan event assessment",
        "report": "Laporan hasil kompetensi"
      },
      "employeemovement": { "_module": "Mutasi, promosi, perubahan kontrak" },
      "useraccount": { "_module": "Akun login untuk karyawan" },
      "attendance": {
        "_module": "Absensi dan pelacakan waktu kerja",
        "settings": "Konfigurasi shift, lokasi, jam kerja",
        "operations": "Aktivitas presensi harian",
        "report": "Laporan ringkasan absensi"
      },
      "approval": {
        "_module": "Alur persetujuan dan permohonan",
        "settings": "Konfigurasi alur persetujuan",
        "operations": "Persetujuan tertunda dan selesai"
      },
      "payroll": { "_module": "Proses dan menjalankan penggajian" },
      "leave": {
        "_module": "Permohonan dan saldo cuti",
        "settings": "Konfigurasi tipe dan kebijakan cuti",
        "operations": "Permohonan dan saldo cuti"
      },
      "performance": {
        "_module": "Evaluasi kinerja dan KPI",
        "settings": "Konfigurasi template evaluasi",
        "operational": "Aktivitas KPI, OKR, evaluasi",
        "report": "Laporan hasil kinerja"
      },
      "recruitment": {
        "_module": "Rekrutmen dan manajemen kandidat",
        "pipeline": "Requisisi, kandidat, wawancara",
        "onboarding": "Onboarding karyawan baru"
      },
      "reimbursement": {
        "_module": "Klaim reimbursemen biaya",
        "requests": "Klaim reimbursemen karyawan",
        "types": "Konfigurasi kategori reimbursemen",
        "reports": "Laporan ringkasan reimbursemen"
      },
      "training": {
        "_module": "Kursus pelatihan dan sertifikasi",
        "settings": "Konfigurasi kursus, kategori, penyedia",
        "operations": "Sesi, peserta, kehadiran",
        "records": "Sertifikat dan riwayat pelatihan",
        "reports": "Laporan ringkasan pelatihan"
      },
      "workforceintelligence": { "_module": "Analitik dan wawasan tenaga kerja" },
      "careerintelligence": { "_module": "Wawasan jenjang karir dan suksesi" },
      "setting": { "_module": "Konfigurasi data referensi sistem" },
      "rbac": {
        "_module": "Manajemen role dan permission",
        "roles": "Definisi dan penugasan role"
      },
      "notification": { "_module": "Notifikasi sistem" }
    },
```

- [ ] **Step 3: Verify JSON validity**

Run: `node -e "JSON.parse(require('fs').readFileSync('frontend/tenant/src/locales/en.json','utf8')); JSON.parse(require('fs').readFileSync('frontend/tenant/src/locales/id.json','utf8')); console.log('OK')"`
Expected: `OK` printed, no parse error.

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/locales/en.json frontend/tenant/src/locales/id.json
git commit -m "feat(rbac): add bilingual rbac.descriptions.* locale keys"
```

---

## Task 3: Render descriptions in `RolePermissions.vue`

**Files:**
- Modify: `frontend/tenant/src/views/settings/RolePermissions.vue:56-66` (submenu column template)
- Modify: `frontend/tenant/src/views/settings/RolePermissions.vue:119-129` (add `descriptionLabel()` next to `moduleLabel()`/`submenuLabel()`)

**Interfaces:**
- Consumes: `rbac.descriptions.<resource>.<submenu|_module>` keys from Task 2; `t()` from `useI18n()` (already imported in this file).
- Produces: `descriptionLabel(resource, submenu)` — returns the description string, or `''` if the key doesn't resolve. `submenu` may be `''`/falsy for module-level rows (matches how `data.submenu` is used elsewhere in this file, e.g. line 63).

- [ ] **Step 1: Add the `descriptionLabel()` helper**

In `frontend/tenant/src/views/settings/RolePermissions.vue`, right after the existing `submenuLabel()` function (currently lines 126-129):

```js
// Label bilingual submenu (rbac.submenus.<resource>.<submenu>) — fallback ke slug.
function submenuLabel(resource, submenu) {
  const key = `rbac.submenus.${resource}.${submenu}`
  return t(key) !== key ? t(key) : submenu
}

// Deskripsi singkat bilingual per baris (rbac.descriptions.<resource>.<submenu|_module>) — kosong jika tidak ada.
function descriptionLabel(resource, submenu) {
  const key = `rbac.descriptions.${resource}.${submenu || '_module'}`
  return t(key) !== key ? t(key) : ''
}
```

- [ ] **Step 2: Add the description line to the submenu column template**

Replace the submenu `<Column>` block (currently lines 56-66):

```html
      <!-- Kolom submenu: baris level-module ("Umum") diikuti tiap submenu -->
      <Column field="submenu" :header="t('rbac.submenu')" style="width:220px">
        <template #body="{ data }">
          <span
            class="text-gray-800 dark:text-gray-100"
            :class="data.submenu ? 'font-normal' : 'font-semibold'"
          >
            {{ data.submenu ? submenuLabel(data.resource, data.submenu) : t('rbac.module_level') }}
          </span>
        </template>
      </Column>
```

with:

```html
      <!-- Kolom submenu: baris level-module ("Umum") diikuti tiap submenu -->
      <Column field="submenu" :header="t('rbac.submenu')" style="width:220px">
        <template #body="{ data }">
          <span
            class="text-gray-800 dark:text-gray-100"
            :class="data.submenu ? 'font-normal' : 'font-semibold'"
          >
            {{ data.submenu ? submenuLabel(data.resource, data.submenu) : t('rbac.module_level') }}
          </span>
          <span class="text-xs text-gray-400 dark:text-gray-500 block">
            {{ descriptionLabel(data.resource, data.submenu) }}
          </span>
        </template>
      </Column>
```

- [ ] **Step 3: Manual verification**

Run: `cd frontend/tenant && npm run dev` (or the repo's existing dev-server command), open `/settings/rbac`, pick a role, open its permissions page (`/settings/rbac/roles/:id/permissions`). Confirm:
- Every row (module-level and submenu) shows a small muted description line under its label.
- Switch the app language to Indonesian and back to English (existing language switcher) — descriptions change accordingly.
- Rows for resources/submenus not yet covered (there should be none, since Task 2 covers all 18 resources) show no literal `rbac.descriptions...` text.

Stop the dev server when done.

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/views/settings/RolePermissions.vue
git commit -m "feat(rbac): show bilingual description under each permission row"
```

---

## Self-Review Notes

- Spec section 1 (locale key structure) → Task 2.
- Spec section 2 (fix `rbac.submenus.*`) → Task 1, table matches spec's table exactly.
- Spec section 3 (`RolePermissions.vue` change, incl. empty-string fallback behavior) → Task 3.
- Spec Testing section (manual only, both languages) → Task 3 Step 3.
- No backend/DB files touched anywhere in this plan, matching the spec's non-goals.
