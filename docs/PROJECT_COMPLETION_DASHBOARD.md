# HRIS Platform — Project Completion Dashboard

**Generated:** 28 July 2026  
**Server:** ✅ Running (`http://localhost:8080`) — Status: `ok`

---

## 📊 Executive Summary

| Metric | Value |
|--------|-------|
| **Tenant Modules** | **16** complete (15 business + Settings with 18 reference CRUDs) |
| **Platform Modules** | **6** complete + **2** shared packages |
| **Total Go Files** | **224+** (155 source + 69 test) |
| **Total GORM Entities** | **121** (116 tenant + 5 platform) |
| **Total Test Functions** | **~1004+** |
| **Total OpenAPI Endpoints** | **684** |
| **Total OpenAPI Schemas** | **439** |
| **Total OpenAPI Tags** | **30** |
| **Module Type Filter** | ✅ **3 endpoints** (`/modules`, `/packages`, `/public/packages`) |
| **Bilingual Support** | ✅ **EN/ID** — Backend 80+ message pairs + Frontend 200+ locale keys, middleware auto-detect, field validation errors |
| **Frontend Phase 1** | ✅ **9/9 Platform Admin pages** — 100% complete with bilingual support |
| **Frontend Tenant (Phase 2)** | ✅ **20+ views** — Dashboard, Login, Profile, Organization, 18 Settings CRUDs |
| **Frontend Components** | **35+** (11 views, 10 form/action components, 3 composables, 2 utils, stores, services) |
| **Frontend Build** | ✅ **Clean** — zero warnings |
| **Migration Files** | **44 per dialect** (22 up + 22 down) |
| **Database Drivers** | PostgreSQL & MySQL |
| **Total Tenant Tables** | **148** (+ schema_migrations = 149) |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Modular Monolith                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐ │
│  │  Platform    │  │   Shared    │  │  16 Tenant Modules   │ │
│  │  Management  │  │   Kernel    │  │                      │ │
│  │              │  │             │  │  • Organization      │ │
│  │ • Company    │  │ • Config    │  │  • Employee          │ │
│  │ • Module     │  │ • Database  │  │  • Job Management    │ │
│  │ • License    │  │ • Auth/JWT  │  │  • Competency        │ │
│  │ • User       │  │ • Middleware│  │  • Employee Movement │ │
│  │ • Monitoring │  │ • Router    │  │  • Attendance        │ │
│  │              │  │ • Logger    │  │  • Approval          │ │
│  │              │  │ • Module    │  │  • Payroll           │ │
│  │              │  │   SDK       │  │  • Leave & Time Off  │ │
│  │              │  │ • Cache     │  │  • Performance       │ │
│  │              │  │ • Authz     │  │  • Recruitment (ATS) │ │
│  │              │  │ • Crypto    │  │  • Reimbursement     │ │
│  │              │  │ • Migrator  │  │  • Training & Dev    │ │
│  │              │  │             │  │  • Workforce Intell  │ │
│  │              │  │             │  │  • Career Intell     │ │
│  │  • Settings (Master) │ │
│  └─────────────┘  └─────────────┘  └──────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Module Completion Status

### Tenant Modules (16/16 — 100% Complete ✅)

| # | Module | Entities | Tests | OpenAPI Endpoints | Status | Completed |
|---|--------|:--------:|:-----:|:-----------------:|:------:|:---------:|
| 1 | Organization Management | 3 | 0 | 12 | ✅ Complete | 22 Jul 2026 |
| 2 | Employee Management | 9 | 35 | 29 | ✅ Complete | 22 Jul 2026 |
| 3 | Job Management | 18 | 67 | 88 | ✅ Complete | 22 Jul 2026 |
| 4 | Competency Management | 7 | 60 | 36 | ✅ Complete | 26 Jul 2026 |
| 5 | Employee Movement & Career | 2 | 62 | 15 | ✅ Complete | 26 Jul 2026 |
| 6 | Time & Attendance | 10 | 88 | 30 | ✅ Complete | 26 Jul 2026 |
| 7 | Approval Engine | 5 | 64 | 15 (+1) | ✅ Complete | 24 Jul 2026 |
| 8 | Payroll & Compensation | 21 | 39 | 47 | ✅ Complete | 24 Jul 2026 |
| 9 | Leave & Time Off | 6 | 39 | 23 | ✅ Complete | 26 Jul 2026 |
| 10 | Performance Management | 7 | 57 | 34 | ✅ Complete | 26 Jul 2026 |
| 11 | Recruitment & Onboarding (ATS) | 7 | 75 | 33 | ✅ Complete | 26 Jul 2026 |
| 12 | Reimbursement & Claim | 3 | 60 | 15 | ✅ Complete | 26 Jul 2026 |
| 13 | Training & Development | 7 | 31 | 35 | ✅ Complete | 26 Jul 2026 |
| 14 | **Workforce Intelligence & Strategic Planning** | **7** | **108** | **68** | ✅ Complete | **25 Jul 2026** |
| 15 | **Career Intelligence & Talent Management** | **4** | **65** | **19** | ✅ Complete | **26 Jul 2026** |
| 16 | **Settings (Master Data)** | **17** | **—** | **85** | ✅ Complete | **27 Jul 2026** |
| | **TOTAL** | **133** | **850** | **584** | **16/16 (100%)** | |

### Platform Modules (6/6 — 100% Complete ✅)

| # | Module | Entities | Tests | OpenAPI Endpoints | Status |
|---|--------|:--------:|:-----:|:-----------------:|:------:|
| 1 | Company & Tenant Lifecycle | 2 | — | 10 | ✅ Complete |
| 2 | Module Management | 2 | — | 7 | ✅ Complete |
| 3 | License Management | 1 | — | 4 | ✅ Complete |
| 4 | **Package Management** | **2** | **25** | **9** | **✅ Complete** |
| 5 | User & Auth (JWT) | 1 | — | 4 (+2 auth) | ✅ Complete |
| 6 | Monitoring | — | — | 3 | ✅ Complete |
| | **TOTAL** | **8** | **25** | **37** | **6/6 (100%)** |

### Shared Packages (7/7 — 100% Complete ✅)

| # | Package | Tests | Description |
|---|---------|:-----:|-------------|
| 1 | `internal/pkg/authz/` (RBAC) | **134** | Database-backed RBAC with 4 default roles, **98 permissions (24 resources)**, auto-reload, **31 assertions for setting resource** (V/C/U/D per role, ResourceFromPath, singularize) |
| 2 | `internal/pkg/cache/` | **51** | Two-tier distributed cache (local sync.Map + Redis) + Pub/Sub invalidation |
| 3 | `internal/pkg/config/` | — | Viper configuration loader (YAML + .env + env vars) |
| 4 | `internal/pkg/database/` | — | Multi-tenant DB connection manager with caching |
| 5 | `internal/pkg/crypto/` | ✅ | AES-256-GCM encryption for tenant credentials |
| 6 | `internal/pkg/migrator/` | ✅ | SQL migration engine (Up/Down/DownTo, embedded files, dual-dialect) |
| 7 | `internal/pkg/middleware/` | — | Auth, CORS, Logger, Recovery, Tenant resolver, **Localize** (Gin middleware stack) |
| **8** | **`internal/pkg/httputil/`** | **—** | **Bilingual response helpers** (SuccessJSON, CreatedJSON, ErrorJSON, NotFound) + locale message catalog (80+ EN/ID pairs) + custom Indonesian validators (NIK, NPWP, KK, Passport, SIM, No Rekening) |

---

## 🧪 Test Coverage Summary

| Module | Tests | % of Total | Key Test Areas |
|--------|:-----:|:----------:|----------------|
| **Shared: authz (RBAC)** | **134** | 14.6% | Enforcer DB loading, repository, service, handler — 24 resources, 98 permissions, **setting resource test coverage** (31 assertions across 7 test functions) |
| **Attendance** | **88** | 9.6% | Repository (37), service (25), handler (21) |
| **Recruitment (ATS)** | **75** | 8.2% | Repository (27), service (23), handler (16) |
| **Job Management** | **67** | 7.3% | Repository, service, handler |
| **Approval Engine** | **64** | 7.0% | Repository (25), service (25), handler (14) |
| **Employee Movement** | **62** | 6.8% | Repository (22), service (22), handler (14) |
| **Competency** | **60** | 6.6% | Service (25), repository (14), handler (15) |
| **Reimbursement** | **60** | 6.6% | Repository (18), service (20), handler (16) |
| **Performance** | **57** | 6.2% | Repository (14), service (24), handler (17) |
| **Shared: cache** | **51** | 5.6% | Unit (42), integration (8), benchmarks (31) |
| **Payroll** | **39** | 4.3% | Repository (13), service (21) |
| **Leave** | **39** | 4.3% | Repository (14), service (12), handler (12) |
| **Employee** | **35** | 3.8% | Repository, service, handler |
| **Training** | **31** | 3.3% | Repository (14), service (9), handler (8) |
| **Workforce Intelligence** | **108** | 10.8% | Repository (31), service (41), handler (36) |
| **Career Intelligence** | **65** | 6.5% | Repository (23), service (20), handler (22) |
| **Organization** | 0 | 0% | (basic CRUD — minimal tests) |
| | | | |
| **GRAND TOTAL** | **~1004** | **100%** | |

---

## 📖 Documentation Coverage

| Document | Description | Status |
|----------|-------------|:------:|
| `README.md` | Main project documentation (setup, API, testing, modules) | ✅ Complete (updated with module_type filter) |
| `ARCHITECTURE_DESIGN_v1.6_Updated.md` | Architecture design, module status, priority matrix | ✅ v16 updated (15 missing endpoints injected) |
| `docs/openapi-report.md` | OpenAPI comprehensive report (v15) | ✅ v15 — 684 endpoints, 439 schemas, 30 tags |
| `docs/go-module-architecture-report.md` | Go module architecture report (entities, services, tests) | ✅ Updated with Settings module (130 entities, 550 service methods, 1029 tests) |
| `docs/platform-architecture-design.md` | Platform architecture design | ✅ Complete |
| `docs/analisis-blueprint-vs-existing.md` | Gap analysis vs existing Laravel app | ✅ Complete |
| `docs/PROJECT_COMPLETION_DASHBOARD.md` | **This document** | ✅ **Updated — Phase 1 Frontend added** |
| `docs/frontend-development-plan.md` | Frontend Phase 1-4 development plan | ✅ **Phase 1 all 9 modules + Tenant 18 Settings CRUDs completed** |
| `docs/Phase-1-Completion-Report.md` | Phase 1 Frontend completion summary | ✅ **NEW — for presentation** |
| **TenantResolver Middleware (SaaS auto-detect)** | ✅ **Done (01 Aug 2026)** | `middleware.TenantResolver` — auto-determine company dari Host header/X-Tenant-ID untuk `/api/v1/tenant/**` (JWT menang → X-Tenant-ID UUID → X-Forwarded-Host → Host) + set response header `X-Tenant-ID`. FE tenant: state `company` di auth store disinkronkan otomatis dari response header (api.js interceptor kirim & baca `X-Tenant-ID`), company opsional di login. 7 unit test. Changelog ARCH v20 |
| **RBAC Permissions Lengkap (24 resource)** | ✅ **Done (01 Aug 2026)** | `authz/rbac.go` — `defaultResources()` + `loadDefaultPolicies()` + `seedDefaults()` kini mencakup **24 resource / 98 permissions** (sebelumnya 18/74): +`performance`, `recruitment`, `reimbursement`, `training`, `workforceintelligence`, `careerintelligence` (4 action tiap resource) — konsisten dengan tenant RBAC seed & `singularize()` map. Fix: resource tsb sebelumnya tidak pernah di-seed ke `rbac_permissions` → tidak muncul di menu RBAC platform-admin. Auto-assign saat restart hanya ke super_admin; role lain di-toggle manual via UI RBAC. +Rbac.vue sortOrder +6 resource. |
| **Rename BankProfileForm → BankAccountForm** | ✅ **Done (01 Aug 2026)** | FE tenant — komponen `employee/BankProfileForm.vue` di-rename jadi `BankAccountForm.vue` (import + tag `<BankAccountForm>` di `EmployeeForm.vue` step 9), label 'Bank Profile' → 'Bank Account' (locale `wizard_step_bank`/`tab_bank`) — konsistensi penamaan dengan label baru |

---

## 🗄️ Database & Migration

| Item | Count | Details |
|------|:----:|---------|
| **Tenant migration files (MySQL)** | **44** (22 up + 22 down) | `001_master_data` → `022_users` |
| **Tenant migration files (Postgres)** | **44** (22 up + 22 down) | Same as MySQL, dialect-adapted |
| **Platform migration files** | **8 DDL** (+6 down = 14) | Platform: 001–006, 007 RBAC, 012; seeders terpisah (2 file) |
| **Total tenant tables** | **148** | Across all 22 migrations |
| **Total with schema_migrations** | **149** | Auto-included by migrator engine |
| **Database drivers** | **2** | PostgreSQL 16+ & MySQL 8+ |
| **Migration engine** | ✅ | SQL-based, embedded, transaction-safe. Terbaru: `021_insurances` (tabel insurances resmi, sebelumnya AutoMigrate) + `022_users` (Level 2 Tenant RBAC identity) |
| **Tenant RBAC Seeder** | ✅ | `tenantseed.SeedTenantRBAC()` — auto-seed 68 permissions (17 resource × 4 action) + default roles Admin (full) & Employee (view-only) saat provisioning via CLI (`handleProvision`/`seed-data`) **dan API** (`company.Service.provisionTenant`); idempotent |

### Tenant Table Distribution

| Group | Tables | Group | Tables |
|-------|:-----:|-------|:-----:|
| **Payroll** | 16 | **Employee** | 14 |
| **Job Management** | 20 | **Master Data** | 12 |
| **Attendance** | 10 | **Competency** | 7 |
| **Performance** | 7 | **Recruitment** | 7 |
| **Training & Dev** | 7 | **Settings & Permissions** | 9 |
| **Leave** | 6 | **Approval** | 5 |
| **Organization** | 5 | **BPJS (seeded)** | 2 (1 setting + 13 rates) |
| **System** | 1 | **Tax** | 1 |
| **Employee Movement** | 2 | **Workforce Intelligence** | 7 |
| **Career Intelligence** | 4 | | |
| | | **Total (incl. 021 insurances, 022 users)** | **145** |

---

## 🔧 Infrastructure Status

| Component | Status | Details |
|-----------|:------:|---------|
| **API Server** | ✅ **Running** | `:8080` — Health check: `ok` |
| **OpenAPI Spec** | ✅ **Served** | `GET /openapi.json` — 684 endpoints |
| **Scalar UI** | ✅ **Served** | `GET /docs` — Interactive API docs with 684 endpoints |
| **RBAC Engine** | ✅ **Active** | 4 default roles, **98 permissions (24 resources)**, auto-reload |
| **On-Premise License Engine** | ✅ **Ready** | `internal/pkg/onpremise/` — RSA `.lic` (expires_at, allowed_modules, max_employees); CLI `licensectl` (gen-key/gen-lic); mode `on_premise` via `HRIS_LICENSE_DEPLOYMENT_MODE` (dormant di mode saas default); lister alternatif PlatformLicenseMiddleware. **`max_employees` di-enforce di `Service.Create()` → 403 `QUOTA_EXCEEDED`** (toast bilingual FE `employee.quota_exceeded`) |
| **Quota Audit (no bypass)** | ✅ **Audited** | Kuota terpusat di `Service.Create()` — satu-satunya pembuat Employee master. Payroll profiles / onboarding / employee-shift / sub-record TIDAK membuat Employee master (tidak perlu kuota). Frontend hanya 1 caller (`EmployeeForm.savePersonalData`). Jalur masa depan (batch import) otomatis kena kuota. *(Audit 31 Jul 2026)* |
| **Cache (Redis)** | 🔶 **Optional** | Redis required for distributed mode |
| **Docker Compose** | ✅ **Ready** | PostgreSQL, Redis, API, Asynqmon |
| **Dockerfile** | ✅ **Multi-stage** | Optimized Go build |
| **Connection Pool** | ✅ **Optimized** | Platform (10/5/1jam), Tenant (10/3/30mnt) |
| **Tenant Credentials** | ✅ **Encrypted** | AES-256-GCM via `internal/pkg/crypto/` |
| **Tenant Lifecycle** | ✅ **Managed** | Suspend/Activate/Terminate + connection cleanup |
| **Frontend Platform Admin** | ✅ **Built & Ready** | 9 views, 200+ locale keys, bilingual EN/ID |

---

## 🗺️ Phase Completion Roadmap

| Phase | Name | Modules | Status |
|:-----:|------|---------|:------:|
| **1** | **Foundation** | Platform, Organization, RBAC, Docker, Migration Engine, Tenant Lifecycle, **Platform Admin Frontend (9 pages)** | ✅ **100%** |
| **2** | **Core Modules** | Employee Management, Job Management | ✅ **100%** |
| **3** | **Payroll & Complex** | Payroll Engine, Competency, BPJS/PPh21 | ✅ **100%** |
| **4** | **Operations & Career** | Attendance, Leave, Employee Movement, Performance, Reimbursement, Recruitment, Training | ✅ **100%** |
| **5** | **Polish** | E2E Testing, Performance Optimization | ⬜ **Planned** |

---

---

## 🖥️ Phase 1 Frontend — Platform Admin (✅ 100%)

### Summary

| Modul | Halaman | Fitur Utama | Status |
|-------|---------|-------------|:------:|
| B.1 | Login | JWT auth, bilingual form, validation error per field | ✅ |
| B.2 | Dashboard | KPI cards (6), bar chart, system health, real-time polling 30s | ✅ |
| B.3 | Companies | Filter by status/package/search, license inline, provisioning progress, **Company Detail page + Rotate Credentials + CompanyActions.vue reusable** | ✅ |
| B.4 | Users | Filter by role, search, bulk actions (delete, change role), super_admin protection | ✅ |
| B.5 | Modules | Filter `?module_type=`, filter chips, depends_on tooltip, auto-slug | ✅ |
| B.6 | Licenses | Package integration, filter by package, copy key, expiry warnings | ✅ |
| B.7 | Monitoring | Auto-refresh toggle, pool utilization chart, alert thresholds | ✅ |
| B.8 | Packages | CRUD, publish/unpublish, dependency validation, module selector with Select All | ✅ |
| B.9 | RBAC | Roles CRUD, permission matrix grouped by resource, ToggleSwitch, system protection | ✅ |

### Bilingual Support

| Komponen | Status |
|----------|:------:|
| Language Store (Pinia) + localStorage persistence | ✅ |
| Axios interceptor `Accept-Language` header | ✅ |
| Language Switcher di HeaderBar (EN/ID) | ✅ |
| Composable `useI18n` → `t(key)` di semua view | ✅ |
| Locale files: 200+ keys di en.json & id.json | ✅ |
| Response handler bilingual (parse message EN/ID) | ✅ |
| Toast notification bilingual | ✅ |
| Validation error per field dengan locale support | ✅ |

### Shared Components

| Komponen | Penggunaan |
|----------|------------|
| **FormRow** | 9 views — form field wrapper + error display |
| **TextInput / SelectLabel / ToggleSwitch / etc.** | 9 form components reusable |
| **DatePicker** | Licenses, Companies (date fields) |
| **TextArea** | Packages, RBAC (description fields) |
| **useSlugify** | Packages, Modules, RBAC (auto-slug) |
| **useI18n / useNotification** | All views (bilingual text + toast) |
| **SkeletonTable / SkeletonCard** | All tenant views (loading skeleton) |
| **ConfirmDeleteDialog** | **20+ tenant views** — Custom dialog yang tetap terbuka selama API call, error tampil inline (menggantikan PrimeVue ConfirmDialog) |
| **CompanyActions** | **Platform Admin Companies (list & detail)** — Komponen reusable untuk Edit/Suspend/Activate/Terminate/Rotate + ConfirmDialog + Edit dialog + rotate password dialog (sekali lihat + copy); `mode="icons"` di list, `mode="buttons"` di header detail; emit `updated` → parent reload; tombol disembunyikan untuk company `terminated` |
| **formatDate.js** | **Utility** — Bilingual date formatting: `formatDate(value, locale)` → "30 July 2026" (EN) / "30 Juli 2026" (ID); `formatDateShort()` → "30/07/2026" |
| **primevueLocale.js** | **Utility** — PrimeVue locale configs EN/ID untuk DatePicker (nama bulan/hari bilingual, first day of week, tombol "Today"/"Hari Ini") |
| **DateInput bilingual** | **OrganizationSummary** — Calendar popup menampilkan nama bulan/hari sesuai bahasa aktif; bug fix: raw string `item.decree_date || null` (ganti `new Date(... + 'T00:00:00')`) |

### Tech Highlights
- Vue 3 Composition API + `<script setup>`
- PrimeVue 4: DataTable, Dialog, Tag, Chart, ToggleSwitch, ConfirmDialog
- PrimeIcons + Tailwind CSS 4
- Chart.js (bar chart Dashboard + line chart Monitoring)
- Pinia stores (auth + language)
- Slug pulse animation CSS (`slug-animation.css`)

---

## 📋 Remaining Work

### Phase 5a: Strategic Analytics Layer

| Task | Priority | Notes | Status |
|------|:--------:|-------|:------:|
| Workforce Intelligence & Strategic Planning | 🟢 High | 11 submodules, 7 entities, 70+ endpoints | ✅ **Done (25 Jul 2026)** |
| Career Intelligence & Talent Management | 🟢 High | 4 entities, 19 endpoints, 9-box grid, career paths, succession planning | ✅ **Done (26 Jul 2026)** |

### Phase 5b: Polish

| Task | Priority | Notes | Status |
|------|:--------:|-------|:------:|
| **Frontend Phase 1 — Platform Admin** | 🟢 High | 9 pages: Login, Dashboard, Companies, Users, Modules, Licenses, Packages, Monitoring, RBAC | ✅ **Done (26 Jul 2026)** |
| Frontend Phase 2 — Tenant Module Views | 🟡 Medium | Organization, Employee, Leave, Payroll, Attendance, dll. | ✅ **Partial (18 Settings CRUD views + Organization done)** |
| — Settings: Zones, Provinces, Regencies, Districts, Villages | 🟢 Easy | 5 geographic entities with parent-child hierarchy | ✅ **Done** |
| — Settings: Educations, Religions, MaritalStatuses, RelationshipTypes | 🟢 Easy | 4 simple reference CRUDs | ✅ **Done** |
| — Settings: Education Majors | 🟢 Easy | 1 reference CRUD — **seeder `seedEducationMajors` pakai kode 3 digit (001–020)**, 20 jurusan, UUID deterministik + idempotent (024_education_majors.sql) | ✅ **Done (1 Aug 2026)** |
| — Settings: Banks, EmploymentStatuses, Nationalities | 🟢 Easy | 3 simple reference CRUDs | ✅ **Done** |
| — Settings: JobFamilies, SalaryGrades | 🟢 Easy | 2 reference CRUDs (SalaryGrades: Grade, MinSalary, MaxSalary) | ✅ **Done** |
| — Settings: Insurances | 🟢 Easy | 1 reference CRUD (Code, Name, SortOrder) — Backend + Frontend + OpenAPI + **seeder kode 2 digit (01 BPJS Kesehatan, 02 BPJS Ketenagakerjaan)** | ✅ **Done** |
| — Settings: TER & PTKP | 🟢 Easy | 2 tax reference CRUDs (TER: Group, BrutoMin, BrutoMax, Rate; PTKP: Name, Group, Amount) | ✅ **Done** |
| — Organization Management | 🟡 Medium | TreeTable view, CRUD with parent selection, dark mode, bilingual (Positions CRUD ⏸️ postponed) | ✅ **Done** |
| — Employee Management (Wizard) | 🔴 Complex | Multi-step form (8/9 steps done — Step 9 Employment + Detail Page ⏸️ postponed) | ✅ **8/9 Steps Done — Remaining Postponed** |
| E2E Testing (Playwright) | 🟡 Medium | Phase 5 | ⬜ Planned |
| Performance Optimization | 🟢 High | Phase 5 | ⬜ Planned |
| CI/CD Pipeline | 🟡 Medium | Phase 5 | ⬜ Planned |

---

## 🔗 Related Documents

| Document | Path |
|----------|------|
| Main README | `./README.md` |
| Deployment Guide (SaaS & On-Premise) | `./docs/deployment-guide.md` |
| Architecture Design | `./ARCHITECTURE_DESIGN_v1.6_Updated.md` |
| OpenAPI Report (v15) | `./docs/openapi-report.md` |
| Go Module Architecture Report | `./docs/go-module-architecture-report.md` |
| Platform Architecture Design | `./docs/platform-architecture-design.md` |
| Blueprint vs Existing Analysis | `./docs/analisis-blueprint-vs-existing.md` |

---

*Dashboard generated from live server data and project documentation.*  
*Server: http://localhost:8080 | Scalar UI: http://localhost:8080/docs*
