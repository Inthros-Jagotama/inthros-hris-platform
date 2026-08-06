# Frontend Performance Management - OKR Module Implementation Plan

## Overview

Implementasi frontend untuk module Performance Management dengan sub-module OKR (Objective & Key Results). Mengikuti pola yang sudah ada pada module KPI.

---

## Referensi Struktur (KPI Module)

```
/performance/kpi                   → KPIIndex.vue (list evaluasi)
/performance/kpi/evaluation/:id    → KPIEvaluationDetail.vue
/performance/kpi/templates         → KPITemplates.vue
/performance/kpi/templates/:id     → KPITemplateForm.vue
```

---

## Struktur OKR Module

### Routes

```
/performance/okr                        → okr/OKRIndex.vue (list evaluasi OKR)
/performance/okr/evaluation/:id         → okr/OKREvaluationDetail.vue (detail evaluasi)
/performance/okr/templates              → okr/OKRTemplates.vue (list template OKR)
/performance/okr/templates/new          → okr/OKRTemplateForm.vue (form template baru)
/performance/okr/templates/:id/edit     → okr/OKRTemplateForm.vue (edit template)
```

---

## Files to Create/Modify

### 1. Router (`src/router/index.js`)

**Add routes:**

```javascript
// Performance Management - OKR
{
  path: 'performance/okr',
  name: 'PerformanceOKR',
  component: () => import('@/views/modules/performance/okr/OKRIndex.vue'),
  meta: { title: 'OKR', titleKey: 'okr.title', descKey: 'okr.description', icon: 'pi pi-bullseye', module: 'performance' }
},
{
  path: 'performance/okr/evaluation/:id',
  name: 'OKREvaluation',
  component: () => import('@/views/modules/performance/okr/OKREvaluationDetail.vue'),
  meta: { title: 'OKR Evaluation', titleKey: 'okr.evaluation', descKey: 'okr.evaluation_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/okr', backLabelKey: 'nav.okr' }
},
{
  path: 'performance/okr/templates',
  name: 'OKRTemplates',
  component: () => import('@/views/modules/performance/okr/OKRTemplates.vue'),
  meta: { title: 'OKR Templates', titleKey: 'okr.templates', descKey: 'okr.templates_desc', icon: 'pi pi-file', module: 'performance' }
},
{
  path: 'performance/okr/templates/new',
  name: 'OKRTemplateNew',
  component: () => import('@/views/modules/performance/okr/OKRTemplateForm.vue'),
  meta: { title: 'New OKR Template', titleKey: 'okr.template_new', descKey: 'okr.template_desc', icon: 'pi pi-plus', module: 'performance', backRoute: '/performance/okr/templates', backLabelKey: 'okr.templates' }
},
{
  path: 'performance/okr/templates/:id/edit',
  name: 'OKRTemplateEdit',
  component: () => import('@/views/modules/performance/okr/OKRTemplateForm.vue'),
  meta: { title: 'Edit OKR Template', titleKey: 'okr.template_edit', descKey: 'okr.template_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/okr/templates', backLabelKey: 'okr.templates' }
},
```

---

### 2. Sidebar Menu (`src/layouts/Sidebar.vue`)

**Update Performance children to include OKR:**

```javascript
{
  key: 'performance',
  label: t('nav.performance'),
  icon: 'pi pi-chart-line',
  command: () => router.push('/performance'),
  moduleSlug: 'performance',
  permission: 'performance.view',
  children: [
    {
      key: 'performance_dashboard',
      label: t('nav.performance_dashboard'),
      icon: 'pi pi-home',
      command: () => router.push('/performance'),
      path: '/performance',
      excludePaths: ['/performance/kpi', '/performance/okr'],
      moduleSlug: 'performance',
      permission: 'performance.view'
    },
    {
      key: 'kpi',
      label: t('nav.kpi'),
      icon: 'pi pi-chart-bar',
      command: () => router.push('/performance/kpi'),
      path: '/performance/kpi',
      moduleSlug: 'performance',
      permission: 'performance.view'
    },
    {
      key: 'okr',
      label: t('nav.okr'),
      icon: 'pi pi-bullseye',
      command: () => router.push('/performance/okr'),
      path: '/performance/okr',
      moduleSlug: 'performance',
      permission: 'performance.view'
    }
  ]
}
```

---

### 3. PerformanceIndex.vue

**Update card menu to include OKR:**

```javascript
const menuItems = [
  {
    key: 'kpi',
    icon: 'pi pi-chart-bar',
    titleKey: 'performance.kpi_evaluations',
    descKey: 'performance.kpi_evaluations_desc',
    path: '/performance/kpi',
    tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400'
  },
  {
    key: 'okr',
    icon: 'pi pi-bullseye',
    titleKey: 'performance.okr_evaluations',
    descKey: 'performance.okr_evaluations_desc',
    path: '/performance/okr',
    tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400'
  },
  {
    key: 'kpi_templates',
    icon: 'pi pi-file',
    titleKey: 'performance.kpi_templates',
    descKey: 'performance.kpi_templates_desc',
    path: '/performance/kpi/templates',
    tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
  },
  {
    key: 'okr_templates',
    icon: 'pi pi-file-edit',
    titleKey: 'performance.okr_templates',
    descKey: 'performance.okr_templates_desc',
    path: '/performance/okr/templates',
    tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400'
  }
]
```

---

### 4. Views to Create

#### Directory Structure

```
src/views/modules/performance/okr/
├── OKRIndex.vue                   # List evaluasi OKR
├── OKREvaluationDetail.vue        # Detail evaluasi dengan objectives & key results
├── OKRTemplates.vue               # List template OKR
└── OKRTemplateForm.vue            # Form template dengan objectives & key results
```

---

### 5. Localization (`src/locales/en.json` & `id.json`)

**Add keys for English:**

```json
{
  "nav": {
    "okr": "OKR"
  },
  "okr": {
    "title": "Objective & Key Results",
    "description": "Manage OKR evaluations",
    "evaluation": "OKR Evaluation",
    "evaluation_desc": "View and manage OKR achievement",
    "templates": "OKR Templates",
    "templates_desc": "Manage OKR objective templates",
    "template_new": "New Template",
    "template_edit": "Edit Template",
    "template_desc": "Define objectives and key results",
    "objectives": "Objectives",
    "key_results": "Key Results",
    "add_objective": "Add Objective",
    "add_key_result": "Add Key Result",
    "objective_title": "Objective Title",
    "key_result_title": "Key Result Title",
    "weight": "Weight",
    "target": "Target",
    "actual": "Actual",
    "achievement": "Achievement",
    "score": "Score",
    "final_score": "Final Score",
    "status": "Status",
    "progress": "Progress",
    "check_in": "Check-in",
    "add_progress": "Add Progress",
    "submit": "Submit for Review",
    "approve": "Approve",
    "reject": "Reject",
    "complete": "Complete",
    "duplicate": "Duplicate Template",
    "target_type": "Target Type",
    "formula_type": "Formula Type",
    "unit": "Unit",
    "no_objectives": "No objectives defined",
    "no_key_results": "No key results defined",
    "objective_weight_total": "Total objective weight must be 100%",
    "key_result_weight_total": "Total key result weight must be 100%"
  },
  "performance": {
    "okr_evaluations": "OKR Evaluations",
    "okr_evaluations_desc": "Manage employee OKR evaluations",
    "okr_templates": "OKR Templates",
    "okr_templates_desc": "Manage OKR objective templates"
  }
}
```

**Add keys for Indonesian:**

```json
{
  "nav": {
    "okr": "OKR"
  },
  "okr": {
    "title": "Objective & Key Results",
    "description": "Kelola evaluasi OKR",
    "evaluation": "Evaluasi OKR",
    "evaluation_desc": "Lihat dan kelola pencapaian OKR",
    "templates": "Template OKR",
    "templates_desc": "Kelola template objective OKR",
    "template_new": "Template Baru",
    "template_edit": "Edit Template",
    "template_desc": "Definisikan objective dan key results",
    "objectives": "Objective",
    "key_results": "Key Results",
    "add_objective": "Tambah Objective",
    "add_key_result": "Tambah Key Result",
    "objective_title": "Judul Objective",
    "key_result_title": "Judul Key Result",
    "weight": "Bobot",
    "target": "Target",
    "actual": "Aktual",
    "achievement": "Pencapaian",
    "score": "Skor",
    "final_score": "Skor Akhir",
    "status": "Status",
    "progress": "Progress",
    "check_in": "Check-in",
    "add_progress": "Tambah Progress",
    "submit": "Ajukan Review",
    "approve": "Setujui",
    "reject": "Tolak",
    "complete": "Selesai",
    "duplicate": "Duplikat Template",
    "target_type": "Tipe Target",
    "formula_type": "Tipe Formula",
    "unit": "Satuan",
    "no_objectives": "Belum ada objective",
    "no_key_results": "Belum ada key results",
    "objective_weight_total": "Total bobot objective harus 100%",
    "key_result_weight_total": "Total bobot key result harus 100%"
  },
  "performance": {
    "okr_evaluations": "Evaluasi OKR",
    "okr_evaluations_desc": "Kelola evaluasi OKR karyawan",
    "okr_templates": "Template OKR",
    "okr_templates_desc": "Kelola template objective OKR"
  }
}
```

---

## API Endpoints Used

### OKR Module (`/performance/okr/*`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/okr/templates` | List templates |
| POST | `/api/v1/tenant/performance/okr/templates` | Create template |
| GET | `/api/v1/tenant/performance/okr/templates/:id` | Get template with objectives |
| PUT | `/api/v1/tenant/performance/okr/templates/:id` | Update template |
| DELETE | `/api/v1/tenant/performance/okr/templates/:id` | Delete template |
| POST | `/api/v1/tenant/performance/okr/templates/:id/duplicate` | Duplicate template |
| GET | `/api/v1/tenant/performance/okr/templates/:id/objectives` | List objectives |
| POST | `/api/v1/tenant/performance/okr/objectives` | Create objective |
| GET | `/api/v1/tenant/performance/okr/objectives/:id` | Get objective with key results |
| PUT | `/api/v1/tenant/performance/okr/objectives/:id` | Update objective |
| DELETE | `/api/v1/tenant/performance/okr/objectives/:id` | Delete objective |
| GET | `/api/v1/tenant/performance/okr/objectives/:id/key-results` | List key results |
| POST | `/api/v1/tenant/performance/okr/key-results` | Create key result |
| GET | `/api/v1/tenant/performance/okr/key-results/:id` | Get key result |
| PUT | `/api/v1/tenant/performance/okr/key-results/:id` | Update key result |
| DELETE | `/api/v1/tenant/performance/okr/key-results/:id` | Delete key result |
| GET | `/api/v1/tenant/performance/okr/evaluations` | List evaluations |
| POST | `/api/v1/tenant/performance/okr/evaluations` | Create evaluation with snapshot |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id` | Get evaluation |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id/details` | Get evaluation with details |
| PUT | `/api/v1/tenant/performance/okr/evaluations/:id` | Update evaluation |
| DELETE | `/api/v1/tenant/performance/okr/evaluations/:id` | Delete evaluation |
| PUT | `/api/v1/tenant/performance/okr/evaluations/:id/actuals` | Bulk update actuals |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/recalculate` | Recalculate score |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/submit` | Submit evaluation |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/approve` | Approve evaluation |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/reject` | Reject evaluation |
| POST | `/api/v1/tenant/performance/okr/evaluations/:id/complete` | Complete evaluation |
| PUT | `/api/v1/tenant/performance/okr/evaluation-details/:id` | Update detail actual |
| GET | `/api/v1/tenant/performance/okr/evaluation-details/:id/progress` | List progress |
| POST | `/api/v1/tenant/performance/okr/progress` | Create progress (check-in) |
| PUT | `/api/v1/tenant/performance/okr/progress/:id` | Update progress |
| DELETE | `/api/v1/tenant/performance/okr/progress/:id` | Delete progress |
| GET | `/api/v1/tenant/performance/okr/evaluations/:id/comments` | List comments |
| POST | `/api/v1/tenant/performance/okr/comments` | Create comment |
| PUT | `/api/v1/tenant/performance/okr/comments/:id` | Update comment |
| DELETE | `/api/v1/tenant/performance/okr/comments/:id` | Delete comment |
| GET | `/api/v1/tenant/performance/okr/evaluation-details/:id/attachments` | List attachments |
| POST | `/api/v1/tenant/performance/okr/attachments` | Upload attachment |
| DELETE | `/api/v1/tenant/performance/okr/attachments/:id` | Delete attachment |
| GET | `/api/v1/tenant/performance/okr/dashboard/hr` | HR dashboard |

---

## Implementation Sequence

### Phase 1: Routes & Navigation
1. Update `router/index.js` - Add OKR routes
2. Update `Sidebar.vue` - Add OKR menu item under Performance
3. Update `PerformanceIndex.vue` - Add OKR cards
4. Update locales (en/id)

### Phase 2: OKR Templates
1. Create `okr/OKRTemplates.vue` - List templates with DataTable
2. Create `okr/OKRTemplateForm.vue` - Form with nested objectives & key results
3. Support CRUD for objectives and key results
4. Duplicate template functionality
5. Weight validation (objectives = 100%, key results per objective = 100%)

### Phase 3: OKR Evaluation
1. Create `okr/OKRIndex.vue` - List evaluations with filters
2. Create `okr/OKREvaluationDetail.vue` - Detail view with:
   - Header card (employee, period, status, final score)
   - Objectives grouped with their key results
   - Actual value input for DRAFT status
   - Progress check-in panel
   - Workflow actions (Submit, Approve, Reject, Complete)
3. Support comments and attachments

### Phase 4: Progress Check-in
1. Add progress tracking per key result
2. Progress history timeline
3. Achievement chart/visualization

---

## Component Details

### OKRIndex.vue
- DataTable with evaluations
- Filters: Period, Status, Organization, Employee
- Columns: Employee, Period, Status, Final Score, Actions
- Actions: View, Edit (if DRAFT), Delete

### OKRTemplates.vue
- DataTable with templates
- Filters: Organization, Period, Status
- Columns: Name, Organization, Objective Count, Status, Actions
- Actions: View, Edit, Duplicate, Delete

### OKRTemplateForm.vue
- Template info section (Name, Organization, Period, Status)
- Objectives panel with accordion/expandable rows
- Each objective has:
  - Title, Description, Weight
  - Key Results sub-table
- Each key result has:
  - Title, Target Type, Target Value, Unit, Formula, Weight
- Add/Edit/Delete objectives and key results inline
- Weight validation

### OKREvaluationDetail.vue
- Header card with:
  - Employee info
  - Period info
  - Template name
  - Status badge
  - Final score with rating color
- Objectives section:
  - Card per objective with title, weight, achievement
  - Expand to show key results table
  - Key results columns: Title, Target, Actual, Achievement, Score
  - Inline actual value input (editable in DRAFT)
- Progress panel (collapsible):
  - Add progress button
  - Progress history list per key result
- Actions footer:
  - Submit (DRAFT → SUBMITTED)
  - Approve/Reject (SUBMITTED → APPROVED/DRAFT)
  - Complete (APPROVED → COMPLETED)
  - Recalculate score

---

## Notes

- Ikuti pattern komponen yang sudah ada pada KPI
- Gunakan composable untuk logic yang reusable
- Semua text menggunakan i18n keys
- Dark mode support required
- Responsive design (mobile-friendly)
- OKR menggunakan master data yang sama dengan KPI:
  - Performance Periods
  - Performance Ratings
  - Performance Indicator Formulas

---

## Implementation Status

| Phase | Status | Completion Date | Notes |
|-------|--------|-----------------|-------|
| Phase 1 - Routes & Navigation | ⏳ Pending | - | |
| Phase 2 - OKR Templates | ⏳ Pending | - | |
| Phase 3 - OKR Evaluation | ⏳ Pending | - | |
| Phase 4 - Progress Check-in | ⏳ Pending | - | |
