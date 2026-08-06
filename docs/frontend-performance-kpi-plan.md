# Frontend Performance Management - KPI Module Implementation Plan

## Overview

Implementasi frontend untuk module Performance Management dengan sub-module KPI (Key Performance Indicator). Mengikuti pola yang sudah ada pada module Job Management.

---

## Referensi Struktur (Job Management)

```
/job-management                    → JobManagement.vue (list page)
/job-management/values             → JobValuesIndex.vue (sub-menu)
/job-management/values/:type       → JobValuesForm.vue
/job-management/form?org_id=xxx    → job/JobManagementForm.vue (multi-section form)
```

---

## Struktur Performance Module

### Routes

```
/performance                       → PerformanceIndex.vue (dashboard/index dengan card menu)
/performance/kpi                   → kpi/KPIIndex.vue (list evaluasi KPI)
/performance/kpi/evaluation/:id    → kpi/KPIEvaluationForm.vue (form evaluasi)
/performance/kpi/templates         → kpi/KPITemplates.vue (list template KPI)
/performance/kpi/templates/:id     → kpi/KPITemplateForm.vue (form template)
```

### Settings Routes (Master Data)

```
/settings/performance-periods      → PerformancePeriodsView.vue
/settings/performance-perspectives → PerformancePerspectivesView.vue
/settings/performance-ratings      → PerformanceRatingsView.vue
/settings/performance-formulas     → PerformanceFormulasView.vue
```

---

## Files to Create/Modify

### 1. Router (`src/router/index.js`)

**Add routes:**

```javascript
// Performance Management - KPI
{
  path: 'performance',
  name: 'Performance',
  component: () => import('@/views/modules/performance/PerformanceIndex.vue'),
  meta: { title: 'Performance', titleKey: 'performance.title', descKey: 'performance.description', icon: 'pi pi-chart-line', module: 'performance' }
},
{
  path: 'performance/kpi',
  name: 'PerformanceKPI',
  component: () => import('@/views/modules/performance/kpi/KPIIndex.vue'),
  meta: { title: 'KPI', titleKey: 'kpi.title', descKey: 'kpi.description', icon: 'pi pi-chart-bar', module: 'performance' }
},
{
  path: 'performance/kpi/evaluation/:id',
  name: 'KPIEvaluation',
  component: () => import('@/views/modules/performance/kpi/KPIEvaluationForm.vue'),
  meta: { title: 'KPI Evaluation', titleKey: 'kpi.evaluation', descKey: 'kpi.evaluation_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/kpi', backLabelKey: 'nav.kpi' }
},
{
  path: 'performance/kpi/templates',
  name: 'KPITemplates',
  component: () => import('@/views/modules/performance/kpi/KPITemplates.vue'),
  meta: { title: 'KPI Templates', titleKey: 'kpi.templates', descKey: 'kpi.templates_desc', icon: 'pi pi-file', module: 'performance' }
},
{
  path: 'performance/kpi/templates/new',
  name: 'KPITemplateNew',
  component: () => import('@/views/modules/performance/kpi/KPITemplateForm.vue'),
  meta: { title: 'New KPI Template', titleKey: 'kpi.template_new', descKey: 'kpi.template_desc', icon: 'pi pi-plus', module: 'performance', backRoute: '/performance/kpi/templates', backLabelKey: 'kpi.templates' }
},
{
  path: 'performance/kpi/templates/:id/edit',
  name: 'KPITemplateEdit',
  component: () => import('@/views/modules/performance/kpi/KPITemplateForm.vue'),
  meta: { title: 'Edit KPI Template', titleKey: 'kpi.template_edit', descKey: 'kpi.template_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/kpi/templates', backLabelKey: 'kpi.templates' }
},

// Settings - Performance
{ path: 'settings/performance-periods', name: 'SettingsPerformancePeriods', component: () => import('@/views/settings/PerformancePeriodsView.vue'), meta: { title: 'Performance Periods', titleKey: 'performance_periods.title', descKey: 'performance_periods.description', icon: 'pi pi-calendar', module: 'setting' } },
{ path: 'settings/performance-perspectives', name: 'SettingsPerformancePerspectives', component: () => import('@/views/settings/PerformancePerspectivesView.vue'), meta: { title: 'BSC Perspectives', titleKey: 'performance_perspectives.title', descKey: 'performance_perspectives.description', icon: 'pi pi-th-large', module: 'setting' } },
{ path: 'settings/performance-ratings', name: 'SettingsPerformanceRatings', component: () => import('@/views/settings/PerformanceRatingsView.vue'), meta: { title: 'Performance Ratings', titleKey: 'performance_ratings.title', descKey: 'performance_ratings.description', icon: 'pi pi-star', module: 'setting' } },
{ path: 'settings/performance-formulas', name: 'SettingsPerformanceFormulas', component: () => import('@/views/settings/PerformanceFormulasView.vue'), meta: { title: 'KPI Formulas', titleKey: 'performance_formulas.title', descKey: 'performance_formulas.description', icon: 'pi pi-calculator', module: 'setting' } },
```

---

### 2. Sidebar Menu (`src/layouts/Sidebar.vue`)

**Modify `talentItems` to add submenu for Performance:**

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
      excludePaths: ['/performance/kpi'],
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
    }
  ]
}
```

---

### 3. Settings Index (`src/views/settings/SettingsIndex.vue`)

**Add new group for Performance settings:**

```javascript
{
  key: 'performance',
  icon: 'pi pi-chart-line',
  labelKey: 'settings.group_performance',
  items: [
    { path: '/settings/performance-periods', icon: 'pi pi-calendar', titleKey: 'settings.performance_periods', descKey: 'performance_periods.description', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' },
    { path: '/settings/performance-perspectives', icon: 'pi pi-th-large', titleKey: 'settings.performance_perspectives', descKey: 'performance_perspectives.description', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' },
    { path: '/settings/performance-ratings', icon: 'pi pi-star', titleKey: 'settings.performance_ratings', descKey: 'performance_ratings.description', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' },
    { path: '/settings/performance-formulas', icon: 'pi pi-calculator', titleKey: 'settings.performance_formulas', descKey: 'performance_formulas.description', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400' }
  ]
}
```

---

### 4. Views to Create

#### Directory Structure

```
src/views/modules/performance/
├── PerformanceIndex.vue           # Dashboard dengan card menu (KPI, Templates, dll)
└── kpi/
    ├── KPIIndex.vue               # List evaluasi KPI
    ├── KPIEvaluationForm.vue      # Form evaluasi dengan multi-section
    ├── KPITemplates.vue           # List template KPI
    ├── KPITemplateForm.vue        # Form template dengan indicators
    └── sections/
        ├── EvaluationInfoSection.vue
        ├── EvaluationDetailSection.vue
        ├── EvaluationProgressSection.vue
        └── EvaluationApprovalSection.vue

src/views/settings/
├── PerformancePeriodsView.vue     # CRUD Performance Periods
├── PerformancePerspectivesView.vue # CRUD BSC Perspectives
├── PerformanceRatingsView.vue     # CRUD Ratings
└── PerformanceFormulasView.vue    # CRUD KPI Formulas
```

---

### 5. Localization (`src/locales/en.json` & `id.json`)

**Add keys:**

```json
{
  "nav": {
    "performance": "Performance",
    "performance_dashboard": "Dashboard",
    "kpi": "KPI"
  },
  "performance": {
    "title": "Performance Management",
    "description": "Manage employee performance evaluation and KPI"
  },
  "kpi": {
    "title": "Key Performance Indicators",
    "description": "Manage KPI evaluations",
    "evaluation": "KPI Evaluation",
    "evaluation_desc": "Input and review KPI achievement",
    "templates": "KPI Templates",
    "templates_desc": "Manage KPI indicator templates",
    "template_new": "New Template",
    "template_edit": "Edit Template",
    "indicators": "Indicators",
    "add_indicator": "Add Indicator",
    "weight": "Weight",
    "target": "Target",
    "actual": "Actual",
    "achievement": "Achievement",
    "score": "Score",
    "final_score": "Final Score",
    "status": "Status",
    "submit": "Submit for Review",
    "approve": "Approve",
    "reject": "Reject",
    "complete": "Complete"
  },
  "performance_periods": {
    "title": "Performance Periods",
    "description": "Manage evaluation periods (Annual, Semester, Quarterly)"
  },
  "performance_perspectives": {
    "title": "BSC Perspectives",
    "description": "Balanced Scorecard perspectives (Financial, Customer, Internal, Learning)"
  },
  "performance_ratings": {
    "title": "Performance Ratings",
    "description": "Rating scales (Outstanding, Excellent, Good, Fair, Poor)"
  },
  "performance_formulas": {
    "title": "KPI Formulas",
    "description": "Calculation formulas (Higher Better, Lower Better, etc)"
  },
  "settings": {
    "group_performance": "Performance",
    "performance_periods": "Performance Periods",
    "performance_perspectives": "BSC Perspectives",
    "performance_ratings": "Performance Ratings",
    "performance_formulas": "KPI Formulas"
  }
}
```

---

## API Endpoints Used

### Performance Module

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/tenant/performance/periods` | List periods |
| POST | `/api/v1/tenant/performance/periods` | Create period |
| GET | `/api/v1/tenant/performance/periods/:id` | Get period |
| PUT | `/api/v1/tenant/performance/periods/:id` | Update period |
| DELETE | `/api/v1/tenant/performance/periods/:id` | Delete period |
| GET | `/api/v1/tenant/performance/perspectives` | List perspectives |
| GET | `/api/v1/tenant/performance/templates` | List templates |
| POST | `/api/v1/tenant/performance/templates` | Create template |
| GET | `/api/v1/tenant/performance/templates/:id` | Get template |
| PUT | `/api/v1/tenant/performance/templates/:id` | Update template |
| DELETE | `/api/v1/tenant/performance/templates/:id` | Delete template |
| GET | `/api/v1/tenant/performance/indicators` | List indicators |
| POST | `/api/v1/tenant/performance/indicators` | Create indicator |
| GET | `/api/v1/tenant/performance/evaluations` | List evaluations |
| POST | `/api/v1/tenant/performance/evaluations/snapshot` | Create evaluation with snapshot |
| GET | `/api/v1/tenant/performance/evaluations/:id/full` | Get evaluation with details |
| PUT | `/api/v1/tenant/performance/evaluation-details/:id/actual` | Update actual value |
| POST | `/api/v1/tenant/performance/evaluations/:id/recalculate` | Recalculate score |
| POST | `/api/v1/tenant/performance/evaluations/:id/submit` | Submit evaluation |
| POST | `/api/v1/tenant/performance/evaluations/:id/approve` | Approve evaluation |
| POST | `/api/v1/tenant/performance/evaluations/:id/reject` | Reject evaluation |
| POST | `/api/v1/tenant/performance/evaluations/:id/complete` | Complete evaluation |
| GET | `/api/v1/tenant/performance/ratings` | List ratings |
| GET | `/api/v1/tenant/performance/indicator-formulas` | List formulas |
| GET | `/api/v1/tenant/performance/dashboard/employee/:id` | Employee dashboard |
| GET | `/api/v1/tenant/performance/dashboard/manager/:id` | Manager dashboard |
| GET | `/api/v1/tenant/performance/dashboard/hr` | HR dashboard |

---

## Implementation Sequence

### Phase 1: Settings (Master Data)
1. Create `PerformancePeriodsView.vue` - CRUD periods
2. Create `PerformancePerspectivesView.vue` - CRUD perspectives
3. Create `PerformanceRatingsView.vue` - CRUD ratings
4. Create `PerformanceFormulasView.vue` - CRUD formulas
5. Update `SettingsIndex.vue` - Add performance group
6. Update `router/index.js` - Add settings routes
7. Update locales

### Phase 2: KPI Templates
1. Create `PerformanceIndex.vue` - Module index with card menu
2. Create `kpi/KPITemplates.vue` - List templates
3. Create `kpi/KPITemplateForm.vue` - Form with indicators management
4. Update `router/index.js` - Add template routes

### Phase 3: KPI Evaluation
1. Create `kpi/KPIIndex.vue` - List evaluations
2. Create `kpi/KPIEvaluationForm.vue` - Multi-section form
3. Create section components for evaluation
4. Update `router/index.js` - Add evaluation routes

### Phase 4: Menu & Navigation
1. Update `Sidebar.vue` - Add Performance submenu
2. Test navigation flow
3. Final locales update

---

## Component Details

### PerformanceIndex.vue
- Card grid menu similar to SettingsIndex.vue
- Cards: KPI Evaluations, KPI Templates, Dashboard
- Quick stats from HR dashboard API

### KPIIndex.vue
- DataTable with evaluations
- Filters: Period, Status, Organization
- Actions: View, Edit, Submit, Approve

### KPIEvaluationForm.vue
- Multi-section form similar to JobManagementForm.vue
- Sections:
  - Info (Employee, Period, Template)
  - KPI Details (Indicators with actual input)
  - Progress (Achievement chart)
  - Approval (Status workflow)
- Auto-calculate achievement & score

### KPITemplateForm.vue
- Template info (Name, Organization, Period)
- Indicators DataTable with inline add/edit
- Drag-drop reorder indicators
- Weight validation (total = 100%)

---

## Notes

- Ikuti pattern komponen yang sudah ada (DataTable, Form, Dialog)
- Gunakan composable untuk logic yang reusable
- Semua text menggunakan i18n keys
- Dark mode support required
- Responsive design (mobile-friendly)
