# HRIS Platform — Project Completion Dashboard

> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md)  
> **Terkait:** [`go-module-architecture-report.md`](go-module-architecture-report.md) · [`openapi-report.md`](openapi-report.md) · [`frontend-development-plan.md`](frontend-development-plan.md)

**Generated:** 28 July 2026  
**Server:** ✅ Running (`http://localhost:8080`) — Status: `ok`

---

## 📊 Executive Summary

| Metric | Value |
|--------|-------|
| **Tenant Modules** | **16** complete (15 business + Settings with 18 reference CRUDs) *(+3 modul tambahan terdaftar di `main.go`: rbac, useraccount & notification — dihitung terpisah; total 20 folder modul tenant, lihat `go-module-architecture-report.md`)* |
| **Platform Modules** | **6** complete + **2** shared packages |
| **Total Go Files** | **346** (231 source + 115 test) |
| **Total GORM Entities** | **229** (221 tenant + 8 platform, dari `go-module-architecture-report.md`) |
| **Total Test Functions** | **~2077+** |
| **Total OpenAPI Endpoints** | **1166** |
| **Total OpenAPI Schemas** | **799** |
| **Total OpenAPI Tags** | **38** |
| **Module Type Filter** | ✅ **3 endpoints** (`/modules`, `/packages`, `/public/packages`) |
| **Bilingual Support** | ✅ **EN/ID** — Backend 80+ message pairs + Frontend 200+ locale keys, middleware auto-detect, field validation errors |
| **Frontend Phase 1** | ✅ **9/9 Platform Admin pages** — 100% complete with bilingual support |
| **Frontend Tenant (Phase 2)** | ✅ **20+ views** — Dashboard, Login, Profile, Organization, 19 Settings CRUDs |
| **Frontend Components** | **35+** (11 views, 10 form/action components, 3 composables, 2 utils, stores, services) |
| **Frontend Build** | ✅ **Clean** — zero warnings |
| **Migration Files** | **150 per dialect** (75 up + 75 down) |
| **Database Drivers** | PostgreSQL & MySQL |
| **Total Tenant Tables** | **199** (+ schema_migrations = 200) |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Modular Monolith                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐ │
│  │  Platform    │  │   Shared    │  │  19 Tenant Modules   │ │
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
│  │              │  │             │  │  • Notification      │ │
│  │              │  │             │  │  • RBAC              │ │
│  │              │  │             │  │  • User Account      │ │
│  │              │  │             │  │  • Settings (Master) │ │
│  └─────────────┘  └─────────────┘  └──────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Module Completion Status

### Tenant Modules (16/16 — 100% Complete ✅)

| # | Module | Entities | Tests | OpenAPI Endpoints | Status | Completed |
|---|--------|:--------:|:-----:|:-----------------:|:------:|:---------:|
| 1 | Organization Management | 3 | 0 | 12 | ✅ Complete | 22 Jul 2026 |
| 2 | Employee Management | 9 | 35 | 40 | ✅ Complete | 22 Jul 2026 |
| 3 | Job Management | 18 | 67 | 97 | ✅ Complete | 22 Jul 2026 *(05 Agu: +8 endpoint tree/clusters/details & auto-recalc)* |
| 4 | Competency Management | 7 | 60 | 67 | ✅ Complete | 26 Jul 2026 |
| 5 | Employee Movement & Career | 2 | 62 | 29 | ✅ Complete | 26 Jul 2026 |
| 6 | Time & Attendance | 10 | 88 | 92 | ✅ Complete | 26 Jul 2026 |
| 7 | Approval Engine | 5 | 64 | 18 | ✅ Complete | 24 Jul 2026 |
| 8 | Payroll & Compensation | 21 | 39 | 83 | ✅ Complete | 24 Jul 2026 |
| 9 | Leave & Time Off | 6 | 39 | 26 | ✅ Complete | 26 Jul 2026 |
| 10 | Performance Management | 7 | 57 | 147 | ✅ Complete | 26 Jul 2026 |
| 11 | Recruitment & Onboarding (ATS) | 7 | 75 | 101 | ✅ Complete | 26 Jul 2026 |
| 12 | Reimbursement & Claim | 3 | 75 | 15 | ✅ Complete | 16 Agu 2026 |
| 13 | Training & Development | 7 | 31 | 123 | ✅ Complete | 26 Jul 2026 |
| 14 | **Workforce Intelligence & Strategic Planning** | **7** | **108** | **70** | ✅ Complete | **25 Jul 2026** |
| 15 | **Career Intelligence & Talent Management** | **4** | **65** | **23** | ✅ Complete | **26 Jul 2026** |
| 16 | **Settings (Master Data)** | **17** | **—** | **107** | ✅ Complete | **27 Jul 2026** |
| | **TOTAL** | **133** | **850** | **1050** | **16/16 (100%)** | |

### Platform Modules (6/6 — 100% Complete ✅)

| # | Module | Entities | Tests | OpenAPI Endpoints | Status |
|---|--------|:--------:|:-----:|:-----------------:|:------:|
| 1 | Company & Tenant Lifecycle | 2 | — | 11 | ✅ Complete |
| 2 | Module Management | 2 | — | 7 | ✅ Complete |
| 3 | License Management | 1 | — | 5 | ✅ Complete |
| 4 | **Package Management** | **2** | **25** | **8** | **✅ Complete** |
| 5 | User & Auth (JWT) | 1 | — | 6 (+2 auth) | ✅ Complete |
| 6 | Monitoring | — | — | 5 | ✅ Complete |
| 7 | RBAC (platform) | — | — | 10 | ✅ Complete |
| 8 | Login & Refresh | — | — | 2 | ✅ Complete |
| | **TOTAL** | **8** | **25** | **54** | **8/8 (100%)** |

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

## 📋 Module Development Plan Status

Ringkasan status implementasi per dokumen plan modul — diverifikasi terhadap kode & masing-masing file plan per **12 Agu 2026**. Detail lengkap ada di file plan-nya.

| Plan | Status | Sisa Pekerjaan |
|------|:------:|----------------|
| [`performance-management-kpi-plan.md`](performance-management-kpi-plan.md) | ✅ **Selesai** | KPI Phase 1–4 (2026-08-06) + Scoring Phase 5 (07 Agu) + Seeder Phase 6 + Approval Integration & Two-Phase Target/Realization Phase 8 + Program Component/Bottom-Up Subordinate Scoring & OKR Self-Assessment Phase 9 (08 Agu) — backend & FE selesai; plan OKR & frontend KPI/OKR sudah di-archive |
| `module-training-development-plan.md` *(di-archive: `docs/archive/`)* | ✅ **100% Selesai** | P0–P2 BE+FE selesai (2026-08-11): migration 088/089/090, CRUD lengkap, Central Approval, sertifikat, reports/history, 14 view FE — tanpa TODO tersisa |
| `module-movement-plan.md` *(di-archive: `docs/archive/`)* | ✅ **Selesai** | 13/13 langkah + enhancement §12 P0–P2 ✅ (backend + FE Movements/Contracts/detail/reports/dashboard). Sisa: eksekusi E2E manual (menunggu environment tenant + akun `employeemovement.*`) + 4 enhancement opsional (bukan blocker) — **di-archive 19 Agu 2026** |
| [`module-notification-plan.md`](module-notification-plan.md) | 🟡 **Sebagian** | Phase 1–4 ✅ + FE-1/2 ✅ + **outcome notification reimbursement ✅ (16 Agu 2026)**: `SetNotifier` di-wire + REIMBURSEMENT_APPROVED/REJECTED/PAID + i18n EN/ID. Belum: Phase 5 — outcome notification payroll/employeemovement/recruitment; FE-3 deep-link per `reference_type` |
| [`module-career-intelligence-plan.md`](module-career-intelligence-plan.md) | 🟡 **Sebagian** | Backend ✅ penuh (talent maps, interests, career paths, gap analysis, succession plans). Belum: halaman FE Talent Maps/Interests/Gap Analysis/Succession Plans (baru career paths) |
| [`module-leave-plan.md`](module-leave-plan.md) | 🟡 **Sebagian** | Phase 1–3 & 5 ✅; FE-1/2 ✅. Belum (Phase 4/6–10): accrual engine, adjustment HR, carry forward, expiry, sub-flow cancellation, team calendar, manager/HR dashboard |
| [`module-attendance-plan.md`](archive/module-attendance-plan.md) | 🟡 **Sebagian** | Mayoritas ✅ (geofence, session calc, corrections, overtime 2-alur, leave integration, notif) + FE-1–5 ✅. Belum: Absent/Exempt detection (butuh scheduled job), Manager/HR Dashboard & Team Calendar, Payroll integration (Phase 13), WRONG_CHECKIN/CHECKOUT otomatis |
| [`career-intelligence-training-enhancement-plan.md`](career-intelligence-training-enhancement-plan.md) | 📋 **Proposal** | 11 fitur enhancement (training profile, gap, recommendation, readiness, dll.) — belum dieksekusi |
| `module-recruitment-development-plan.md` *(di-archive: `docs/archive/`)* | ✅ **100% Selesai** | Blueprint Integrated Recruitment **operasional** — **G-1 s.d. G-12 ✅ (12–13 Agu 2026)**: G-1 requisition → Central Approval (migration 093); G-2 requisition enhancement (094: requisition_number/priority/position_id/opened_at); G-3 offer management (095 tabel `job_offers`, workflow approval modul `recruitment_offer`); G-4 → Employee/Movement (096: `employee.recruited_from_application_id` + `candidates.candidate_type`/`employee_id`); G-5 stage history; G-6 candidate enhancement (6 sub-tabel); G-7 screening & assessment; G-8 multi-interviewer & scorecard; G-9 requirement/kompetensi & match score; G-10 onboarding template scope; G-11 analytics; G-12 FE penuh (9 view + hub) — **101 endpoint**, **272 test** (module lolos), migrasi 093–106. Strategic layer **dipisah** ke [`module-recruitment-strategic-layer-plan.md`](archive/module-recruitment-strategic-layer-plan.md) (S-1 s.d. S-7 ✅); alur pengguna: [`flow/module-recruitment-flow.md`](flow/module-recruitment-flow.md) |
| `module-recruitment-strategic-layer-plan.md` *(di-archive: `docs/archive/`)* | ✅ **Selesai** | Strategic layer ↔ Recruitment (S-1 s.d. S-7 ✅ 12 Agu 2026) — **di-archive 19 Agu 2026** |
| [`workforce-intelligence-training-enhancement-plan.md`](workforce-intelligence-training-enhancement-plan.md) | 📋 **Proposal** | 13 fitur enhancement (training demand, coverage, cost, readiness, risk, dll.) — belum dieksekusi |

---

## 🧪 Test Coverage Summary

| Module | Tests | % of Total | Key Test Areas |
|--------|:-----:|:----------:|----------------|
| **Shared: authz (RBAC)** | **134** | 14.6% | Enforcer DB loading, repository, service, handler — 24 resources, 98 permissions, **setting resource test coverage** (31 assertions across 7 test functions) |
| **Attendance** | **88** | 9.6% | Repository (37), service (25), handler (21) |
| **Recruitment (ATS)** | **128** | 11.8% | Repository (27), service (73), handler (28) — G-1..G-4: approval requisition + requisition enhancement + offer management + handoff Employee/Movement |
| **Job Management** | **67** | 7.3% | Repository, service, handler |
| **Approval Engine** | **64** | 7.0% | Repository (25), service (25), handler (14) |
| **Employee Movement** | **62** | 6.8% | Repository (22), service (22), handler (14) |
| **Competency** | **60** | 6.6% | Service (25), repository (14), handler (15) |
| **Reimbursement** | **75** | 8.1% | Repository (18), service (20), handler (16) + approval auto-resolve & notifier integration (16 Agu 2026) |
| **Performance** | **57** | 6.2% | Repository (14), service (24), handler (17) |
| **Shared: cache** | **51** | 5.6% | Unit (42), integration (8), benchmarks (31) |
| **Payroll** | **39** | 4.3% | Repository (13), service (21) |
| **Leave** | **39** | 4.3% | Repository (14), service (12), handler (12) |
| **Employee** | **41** | 3.8% | Repository, service, handler — G-4: recruited_from_application_id |
| **Training** | **31** | 3.3% | Repository (14), service (9), handler (8) |
| **Workforce Intelligence** | **108** | 10.8% | Repository (31), service (41), handler (36) |
| **Career Intelligence** | **65** | 6.5% | Repository (23), service (20), handler (22) |
| **Organization** | 0 | 0% | (basic CRUD — minimal tests) |
| | | | |
| **GRAND TOTAL** | **~1600** | **100%** | *(termasuk shared packages & middleware — per-modul lihat kolom Tests di tabel atas)* |

> ℹ️ **Angka resmi terkini** (entities, service/repo/handler methods, route regs, test per modul) ada di [`go-module-architecture-report.md`](go-module-architecture-report.md) — di-generate dari analisis statis kode (`python scripts/generate_go_module_report.py`). Baris per-modul di tabel atas adalah snapshot per tanggal completion; jika ada selisih, gunakan report arsitektur sebagai sumber kebenaran.

---

## 📖 Documentation Coverage

| Document | Description | Status |
|----------|-------------|:------:|
| `README.md` | Main project documentation (setup, API, testing, modules) | ✅ Complete (updated with module_type filter) |
| `docs/platform-architecture-design.md` | Architecture design (modular monolith, multi-tenant) | ✅ Complete |
| `docs/openapi-report.md` | OpenAPI comprehensive report (report v23, spec 1.6.3) | ✅ v23 — 1166 endpoints, 701 paths, 799 schemas, 38 tags |
| `docs/go-module-architecture-report.md` | Go module architecture report (entities, services, tests) | ✅ 229 entities, 1265 service methods, 2077 tests |
| `docs/platform-architecture-design.md` | Platform architecture design | ✅ Complete |
| `docs/analisis-blueprint-vs-existing.md` | Gap analysis vs existing Laravel app | ✅ Complete |
| `docs/project-completion-dashboard.md` | **This document** | ✅ **Updated — Phase 1 Frontend added** |
| `docs/frontend-development-plan.md` | Frontend Phase 1-4 development plan | ✅ **Phase 1 all 9 modules + Tenant 19 Settings CRUDs completed** |
| **TenantResolver Middleware (SaaS auto-detect)** | ✅ **Done (01 Aug 2026)** | `middleware.TenantResolver` — auto-determine company dari Host header/X-Tenant-ID untuk `/api/v1/tenant/**` (JWT menang → X-Tenant-ID UUID → X-Forwarded-Host → Host) + set response header `X-Tenant-ID`. FE tenant: state `company` di auth store disinkronkan otomatis dari response header (api.js interceptor kirim & baca `X-Tenant-ID`), company opsional di login. 7 unit test. Changelog ARCH v20 |
| **RBAC Permissions Lengkap (24 resource)** | ✅ **Done (01 Aug 2026)** | `authz/rbac.go` — `defaultResources()` + `loadDefaultPolicies()` + `seedDefaults()` kini mencakup **24 resource / 98 permissions** (sebelumnya 18/74): +`performance`, `recruitment`, `reimbursement`, `training`, `workforceintelligence`, `careerintelligence` (4 action tiap resource) — konsisten dengan tenant RBAC seed & `singularize()` map. Fix: resource tsb sebelumnya tidak pernah di-seed ke `rbac_permissions` → tidak muncul di menu RBAC platform-admin. Auto-assign saat restart hanya ke super_admin; role lain di-toggle manual via UI RBAC. +Rbac.vue sortOrder +6 resource. |
| **Rename BankProfileForm → BankAccountForm** | ✅ **Done (01 Aug 2026)** | FE tenant — komponen `employee/BankProfileForm.vue` di-rename jadi `BankAccountForm.vue` (import + tag `<BankAccountForm>` di `EmployeeForm.vue` step 9), label 'Bank Profile' → 'Bank Account' (locale `wizard_step_bank`/`tab_bank`) — konsistensi penamaan dengan label baru |
| **Company Holidays CRUD (Settings)** | ✅ **Done (01 Aug 2026)** | `company_holidays` (Hari Libur Perusahaan) — Backend: model/repo/service/handler/routes/module + permissions `setting.company_holiday.*`; OpenAPI +5 endpoints +3 schemas (373 paths, 689 endpoints, 442 schemas, report v16); FE: CompanyHolidaysView.vue (DateInput + ToggleSwitch + ConfirmDeleteDialog) + route + SettingsIndex card + locale EN/ID. **Catatan:** tabel tanpa updated_at/deleted_at → Update `Updates(map)` + Delete hard-delete; model IsActive tanpa tag `default:true` (gotcha GORM); fix bilingual: label field pakai key `desc` terpisah dari page description |
| **Company Holidays — Calendar View + Badge Fix** | ✅ **Done (01 Aug 2026)** | FE tenant CompanyHolidaysView.vue — Toggle Table/Calendar + DatePicker inline (primevue ^4.5.5) + slot `#date` badge rose penuh di tanggal libur + panel legend bulan berjalan + klik tanggal → edit (ada libur) / tambah dgn tanggal terisi (belum ada); `viewYear/viewMonth` sinkron via `@month-change`/`@year-change`. **Bug fix:** slot `#date` v4 menerima object `{day, month 0-indexed, year}` (bukan Date) → `date.day` + `toDateKey()` dual-mode (fix `date.getDate is not a function`); normalisasi `normDateKey()` (`slice(0,10)`) di `holidayMap`/`monthHolidays` agar badge muncul apa pun format `holiday_date` dari API |
| **Org Tree drag & drop (pindah parent)** | ✅ **Done (02 Aug 2026)** | FE tenant `modules/Organizations.vue` mode Tree — PrimeVue Tree `:draggable-nodes`/`:droppable-nodes`/`validate-drop` + `@node-drop` (`dragNode`/`dropNode`/`dropPosition`/`accept`); parent baru dari `dropPosition` (0 → dropNode; ±1 → parent dropNode; dropNode null → root); guard anti-siklus self/descendant (toast + reload, tanpa ghost state); `PUT /organizations/:id { parent_id }` + reload; `moving` ref sebagai guard concurrency. **Bug fix:** kode lama pakai props `draggable`/`droppable` + `event.node`/`event.dropIndex` yang TIDAK ada di PrimeVue 4.5.5 (API asli: `draggableNodes`/`droppableNodes`/`validateDrop` + `dragNode`/`dropPosition`/`accept`) → drag tidak pernah aktif. Build FE ✅ (tenant dev :5174) |
| **Job Value Mapping — seed tipe activity/subordinate/authority/impact_unauthorized/authority_unauthorized/asset_authority/communicating_influencing_skill/thinking_environment/thinking_chalenge/psychological/technical/managerial** | ✅ **Done (02 Aug 2026)** | `job_management_values` — migration **033** (`activity`/Aktifitas Fisik, 5 level, UUID `aaaa...`), **034** (`subordinate`/Total Bawahan, 5 level, UUID `bbbb...`), **035** (`authority`/Wewenang, 8 level note 'Memiliki Wewenang', UUID `cccc...`), **036** (`impact_unauthorized`/Dampak Tanpa Wewenang Keuangan, 6 level PAnciliary–Berpengaruh, UUID `dddd...`), **037** (`authority_unauthorized`/Otoritas Tanpa Wewenang Keuangan, 8 level note 'Tidak memiliki Wewenang', UUID `eeee...`), **038** (`asset_authority`/Wewenang Asset, 6 level Tidak ada Aset–Menyetujui Aset, UUID `ffff...`), **039** (`communicating_influencing_skill`/Communicating & Influencing Skill, 3 level note Berkomunikasi–Perubahan Perilaku, UUID `1111...`), **040** (`thinking_environment`/Lingkungan Berpikir, 8 level + `thinking_chalenge`/Tantangan Berpikir, 5 level — grup 'Problem Solving & Decision Making', UUID `2222...`; teks note seeder di-promote menjadi descriptions per level, pola sama dgn 039), **041** (`kecerdasan`/Intelligence 5 level + `innovation_creativity` 8 + `self_confidence` 8 + `flexibility` 8 + `tenacity` 8 + `continuous_learning` 8 = 45 row — grup 'Psychological', UUID `3333...`; teks note seeder di-promote menjadi descriptions per level), **042** (16 tipe grup 'Technical' x 8 level + 6 tipe grup 'Managerial' x 5 level = 158 row, UUID `4444...`; descriptions Primary→Unique Authority utk Technical, Task→Total Managerial utk Managerial — pola sama 039-041) — data persis `JobManagementValuesTableSeeder.php` legacy (kini di `docs/tmp/inthros-web/database/seeders/` — folder diabaikan), idempotent `WHERE NOT EXISTS` per id, down scoped ke UUID seed. FE: pola **card menu** (`JobValuesIndex.vue` mengikuti SettingsIndex) + page title per-tipe di layout (AppLayout special case JobValuesType) + util bersama `@/utils/jobValues.js` (single source of truth label/desc/valid types — dipakai AppLayout, HeaderBar, Index, Section, Form). OpenAPI deskripsi filter `?type=` diperbarui. Total endpoint tidak berubah. Catatan: 6 row authority lama (dari grup 'Wewenang Asset') telah hilang dari DB — di-seed ulang 8 level grup 'Wewenang' sesuai keputusan user; tipe `impact` (Memiliki Wewenang Keuangan) & `authority` sudah ada, 036/037 menambah varian 'Tidak Memiliki Wewenang Keuangan' sebagai tipe baru, 038 menambah 'Wewenang Asset', 039 menambah 'Communicating & Influencing Skill', 040 menambah 2 tipe grup 'Problem Solving & Decision Making', 041 menambah 6 tipe grup 'Psychological' |
| **Job Value UUID standardization (migration 043)** | ✅ **Done (02 Aug 2026)** | `job_management_values` — semua **306 baris distandarkan ke UUID v4 acak proper** (kolom char(36)). Latar: seed 033-042 memakai UUID tetap/serial (`aaaaaaaa-0001-4000-8000-...`) demi idempotensi `WHERE NOT EXISTS`; user minta standar UUID proper. Migration **043** (mysql + postgres + down) mengubah semua id pola tetap `-4000-8000-` → UUID v4 (MySQL: ekspresi `RANDOM_BYTES`-based v4; Postgres: `gen_random_uuid()`), berlaku untuk DB existing MAUPUN fresh install (043 berjalan setelah 042). Down migration mengembalikan id tetap seed dengan mencocokkan (type, level, descriptions, sort) — kombinasi unik per row. Tanpa FK & child tables kosong (assets/financials/potency_comp/relationships = 0 row) → aman tanpa dangling reference. **Bug fix penting:** versi awal memakai `MOD(ABS(RAND()*100), 4)` yang saat RAND ≈ 0.9999 menghasilkan posisi `SUBSTR` 5 → segmen ke-4 hanya 3 karakter (UUID 35-char, 31 row malformed) — diperbaiki ke `FLOOR(RAND()*4)` (integer 0-3, 30x test valid). |
| **Icon fix Job Value Mapping (PrimeIcons 8.0.0)** | ✅ **Done (02 Aug 2026)** | `JobValuesIndex.vue` memakai 4 ikon yang **tidak ada** di PrimeIcons 8.0.0 → ikon tidak render (kosong) untuk kartu: `pi-brain` (Lingkungan Berpikir + Kecerdasan) → diganti `pi-wave-pulse` / `pi-microchip`; `pi-puzzle` (Tantangan Berpikir) → `pi-flag`; `pi-target` (Achievement Orientation) → `pi-bullseye`; `pi-handshake` (Building Partnership) → `pi-link`. Audit semua 37 ikon unik di JobValuesIndex vs `primeicons.css` → **0 missing tersisa**. Build FE ✅ + browser test ✅ (7 grup card render, 0 console errors) |
| **Seed education & experience (migration 049)** | ✅ **Done (04 Aug 2026)** | `job_management_values` — tipe `education` (5 level: SMP, SMA, D3, S1, S2/S3) & `experience` (5 level: 0–2, 3–5, 6–8, 9–11, >12 th) — MySQL & Postgres (up/down, idempoten by UUID v4). Menutup gap referensi pendidikan yang tidak ada di seed 033–042. |
| **Seed cash/impact/environment/risk/relationship/frequency/asset (migration 050)** | ✅ **Done (04 Aug 2026)** | `job_management_values` — 7 tipe referensi yang belum ter-seed: `cash` (5 level Jumlah Uang), `impact` (6 level Dampak dengan wewenang keuangan), `environment` (5 Lingkungan Kerja), `risk` (5 Resiko/Bahaya), `relationship` (5 Lingkup Hubungan Kerja), `frequency` (5 Frekuensi), `asset` (8 Nilai Asset) — level & deskripsi dari seeder legacy `JobManagementValuesTableSeeder.php`, MySQL & Postgres (up/down, idempoten by UUID v4). |
| **Job Score Calculator Go + recalc otomatis** | ✅ **Done (05 Aug 2026)** | `backend/internal/modules/jobmanagement/calculator.go` — port penuh `JobValueCalculator.php`: `CalculateForOrganization(s)` + batch `CalculateForOrganizationIDs`, mapping 5 tabel (MAP_DEFAULT/EXTENDED/LINEAR_5/LINEAR_8/COMMUNICATION `{1:1,2:3,3:6}`), 10 komponen (edu-exp, potentials, competencies, problem solving, financial, asset, subordinate, scope, activity, risk), `isResultComplete`. **Hook recalc otomatis** di Create/Update/Delete tiap section yang memengaruhi skor → `RecalculateJobScore` → `UpsertJobScore`. Bobot (`weight`) tidak mengubah skor (keputusan user — rumus legacy). |
| **is_complete & completed_at di job_management_scores (migration 051)** | ✅ **Done (05 Aug 2026)** | MySQL (`TINYINT(1) NOT NULL DEFAULT 0` + `TIMESTAMP NULL`) & Postgres (`SMALLINT` + `TIMESTAMP`), up/down. Diisi `scoreFromResult` dari `JobScoreResult.IsComplete` (`completed_at` hanya saat complete). Diekspos ke API → badge Lengkap/Belum Lengkap di FE. |
| **type_group & description_group job_management_values (migration 052/053)** | ✅ **Done (05 Aug 2026)** | **052** tambah kolom `type_group` & `description_group` (MySQL + Postgres, up/down); **053** seed mengisi keduanya untuk semua tipe (contoh `type_group=psychological`, `type=kecerdasan` → `description_group=Kecerdasan`) — idempoten by UUID. Dipakai form potensi & tree endpoint. |
| **Job value clusters (migration 054) + endpoint mapping** | ✅ **Done (05 Aug 2026)** | Tabel `job_management_value_clusters` untuk mapping tipe technical/managerial → cluster kompetensi. Endpoint `GET/PUT /api/v1/tenant/job-management/values/clusters/:type` (list & simpan mapping) — dipakai setting mapping Job Value submenu Technical & Managerial + card Kompetensi Teknis/Manajerial (filter cluster ≠ Core/Manajerial dari table `competencies`, option dari `job_management_values` type `technical`/`managerial`). |
| **Endpoint tree job management values** | ✅ **Done (05 Aug 2026)** | `GET /api/v1/tenant/job-management/values/tree` — hierarki `type_group → daftar tipe (label = description_group) → options per tipe (level + deskripsi)`, urutan grup tetap (education…frequency). Dipakai form potensi: filter type_group → multi-select tipe psikologis → tabel isian per tipe. |
| **Halaman score: breakdown dari field components** | ✅ **Done (05 Aug 2026)** | `JobScoreSection.vue` ditulis ulang — breakdown 10 komponen dari field `components` (JSON nested `job_management_scores`), bukan `sub_component_points`; menampilkan poin (level → poin) + skor per komponen, urut sesuai navigasi; `competencies` memakai `base_score` (hindari dobel hitung); footer total with/without financial. **Ringkasan skor sticky** (`JobScoreSummary.vue`) di form utama — selalu terlihat saat pindah navigasi. |
| **Daftar job management: kolom score + search, tanpa code/order/level** | ✅ **Done (05 Aug 2026)** | `JobManagement.vue` — kolom `level` & `sort_order` (order) & `code` dihapus; kolom **Score** (`job_value_without_financial`), **With Financial** (Yes/No dari `has_financial_authority`), **Status Complete** (badge is_complete) ditambahkan (data score di-loop pagination dari `GET /scores` lalu di-merge via organization_id); **pencarian** (input 🔍 debounce 350ms, reset page 1). |
| **Search param di API organizations** | ✅ **Done (05 Aug 2026)** | `backend/internal/modules/organization/` — handler/service/repo `List`/`FindAll` menerima `?search=` → filter `LIKE %…%` pada `organizations.code`/`full_code`/`nomenclature` (prefix tabel `organizations.` untuk hindari error 1052 ambiguous saat JOIN `organization_summaries`). Dipakai pencarian daftar job management. |
| **Potensi & kompetensi: card self-contained + composable bersama** | ✅ **Done (05 Aug 2026)** | JobPotencySection dipecah per card: `PsychologicalPotencyCard`, `SkillPotencyCard` (Communication & Influencing), `ProblemSolvingPotencyCard`, `TechnicalPotencyCard` (filter competencies cluster ≠ Core/Manajerial + option type `technical` + bobot % per kompetensi), `ManagerialPotencyCard` (bobot mengikuti technical 100−x). Logika simpan/hapus/hydrate diekstrak ke composable `usePotencyLevels` + komponen tabel bersama; delete per-row; bobot simpan on-blur; default bobot technical 100/jumlah kompetensi. |

---

## 🗄️ Database & Migration

| Item | Count | Details |
|------|:----:|---------|
| **Tenant migration files (MySQL)** | **150** (75 up + 75 down) | `001_master_data` → `076_notification_params` |
| **Tenant migration files (Postgres)** | **150** (75 up + 75 down) | Same as MySQL, dialect-adapted (penomoran 001–076; nomor 044 tidak ada) |
| **Platform migration files** | **9 DDL** (+7 down = 16) | Platform: 001–007, 012, 013; seeders terpisah (2 file) |
| **Total tenant tables** | **199** | Across all 89 migrations |
| **Total with schema_migrations** | **176** | Auto-included by migrator engine |
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
| **OpenAPI Spec** | ✅ **Served** | `GET /openapi.json` — 1166 endpoints |
| **Scalar UI** | ✅ **Served** | `GET /docs` — Interactive API docs with 1166 endpoints |
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
| Career Intelligence & Talent Management | 🟢 High | 5 entities, 23 endpoints, 9-box grid, career paths, succession planning | ✅ **Done (26 Jul 2026)** |

### Phase 5b: Polish

| Task | Priority | Notes | Status |
|------|:--------:|-------|:------:|
| **Frontend Phase 1 — Platform Admin** | 🟢 High | 9 pages: Login, Dashboard, Companies, Users, Modules, Licenses, Packages, Monitoring, RBAC | ✅ **Done (26 Jul 2026)** |
| Frontend Phase 2 — Tenant Module Views | 🟡 Medium | Organization, Employee, Leave, Payroll, Attendance, dll. | ✅ **Partial (19 Settings CRUD views + Organization done + Leave FE-1/FE-2 + Attendance FE-1 s.d. FE-5)** |
| — Settings: Zones, Provinces, Regencies, Districts, Villages | 🟢 Easy | 5 geographic entities with parent-child hierarchy | ✅ **Done** |
| — Settings: Educations, Religions, MaritalStatuses, RelationshipTypes | 🟢 Easy | 4 simple reference CRUDs | ✅ **Done** |
| — Settings: Education Majors | 🟢 Easy | 1 reference CRUD — **seeder `seedEducationMajors` pakai kode 3 digit (001–020)**, 20 jurusan, UUID deterministik + idempotent (024_education_majors.sql) | ✅ **Done (1 Aug 2026)** |
| — Settings: Banks, EmploymentStatuses, Nationalities | 🟢 Easy | 3 simple reference CRUDs | ✅ **Done** |
| — Settings: JobFamilies, SalaryGrades | 🟢 Easy | 2 reference CRUDs (SalaryGrades: Grade, MinSalary, MaxSalary) | ✅ **Done** |
| — Settings: Insurances | 🟢 Easy | 1 reference CRUD (Code, Name, SortOrder) — Backend + Frontend + OpenAPI + **seeder kode 2 digit (01 BPJS Kesehatan, 02 BPJS Ketenagakerjaan)** | ✅ **Done** |
| — Settings: Company Holidays | 🟢 Easy | 1 reference CRUD (HolidayDate unique, Name, Description, IsActive) — Backend + Frontend (CompanyHolidaysView.vue) + OpenAPI (report v16, 373 paths) | ✅ **Done (1 Aug 2026)** |
| — Settings: TER & PTKP | 🟢 Easy | 2 tax reference CRUDs (TER: Group, BrutoMin, BrutoMax, Rate; PTKP: Name, Group, Amount) | ✅ **Done** |
| — Organization Management | 🟡 Medium | TreeTable view (semua level), CRUD with parent selection, **Tree view drag & drop (pindah parent)**, dark mode, bilingual (Positions CRUD ⏸️ postponed) | ✅ **Done** |
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
| Architecture Design | `./docs/platform-architecture-design.md` |
| OpenAPI Report | `./docs/openapi-report.md` |
| Go Module Architecture Report | `./docs/go-module-architecture-report.md` |
| Platform Architecture Design | `./docs/platform-architecture-design.md` |
| Blueprint vs Existing Analysis | `./docs/analisis-blueprint-vs-existing.md` |

---

*Dashboard generated from live server data and project documentation.*  
*Server: http://localhost:8080 | Scalar UI: http://localhost:8080/docs*
