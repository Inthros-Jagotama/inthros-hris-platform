# Architecture Design Document (ADD)
**Blueprint Roadmap v1.6 & System Specification**

---

## 1. Executive Summary & Core Architecture

Sistem ini dirancang menggunakan pendekatan **Modular Monolith** berbasis bahasa pemrograman **Go (Golang)** dengan pola isolasi data **Database per Tenant**. Arsitektur ini menggabungkan fleksibilitas pemeliharaan kode berbasis modul dengan keamanan isolasi data tingkat tinggi antar tenant.

### Tech Stack Utama
* **Backend:** Go (Golang) dengan Web Framework (Gin/Echo) & GORM ORM
* **Frontend:** Next.js / Vue.js
* **Database:** PostgreSQL / MySQL (Multi-tenant isolation)
* **Caching & Message Bus:** Redis & Asynq (Distributed Task Queue)
* **Connection Pooler:** PgBouncer (Rekomendasi Production)

---

## 2. Module Implementation Status

### 2.0. Modul Go Backend — Status Implementasi

| Modul | Entities | Tests | Status |
|:---|---:|:---:|:---:|
| **Platform & Tenant Management** | 6 platform entities | ✅ | ✅ **Completed** — Company, License, Module Management (with `module_type` filter: platform/tenant), Users, Monitoring, Package Management (with `module_type` filter) |
| **Organization Management** | 5 GORM entities | ❌ | ✅ **Completed** — Organization, Zone, Position, Level, Summary |
| **Employee Management** | 9+ entities | ✅ 25+ tests | ✅ **Completed** — Personal, kontak, alamat, keluarga, pendidikan, dokumen, dll |
| **Job Management** | 18 GORM entities | ✅ **74 tests** | ✅ **Completed** — Titles, Values, Responsibilities, Authorities, Working Activities/Risks |
| **Competency Management** | 7 GORM entities | ✅ **54 tests** | ✅ **Completed (26 Jul 2026)** — Competencies, Values, Events, Targets, Scores, Details. **35 OpenAPI endpoints** |
| **Employee Movement & Career** | 2 GORM entities | ✅ **58 tests** | ✅ **Completed (26 Jul 2026)** — Movements (promotion/demotion/mutation/retirement/offboarding), Contracts (PKWT/PKWTT). **15 OpenAPI endpoints** |
| **Approval Engine** | **5 GORM entities** | ✅ **67 tests** | ✅ **Completed (24 Juli 2026)** — Flows, Steps, Instances, Actions, Tasks. **15 OpenAPI endpoints** |
| **Payroll & Compensation Engine** | **21 GORM entities** | ✅ **34 tests** | ✅ **Completed (24 Juli 2026)** — BPJS, PPh21, Payroll Run, Employee Profiles. **47 OpenAPI endpoints** |
| **Time & Attendance** | **10 GORM entities** | ✅ **83 tests** | ✅ **Completed (26 Juli 2026)** — Settings, Shifts, Employee Shifts, Locations, Events, Sessions, Overtime, Exempt Positions. **30 OpenAPI endpoints** |
| **Leave & Time Off** | **6 GORM entities** | ✅ **38 tests** | ✅ **Completed (26 Juli 2026)** — Leave Types, Accrual Policies, Reasons, Requests (23 endpoints), Details, Balances |
| **Performance Management** | **7 GORM entities** | ✅ **55 tests** | ✅ **Completed (26 Jul 2026)** — Periods, Perspectives, Templates, Indicators, Evaluations, Details, Targets. KPI/OKR/BSC framework. **34 OpenAPI endpoints** |
| **Recruitment & Onboarding ATS** | **7 GORM entities** | ✅ **66 tests** | ✅ **Completed (26 Jul 2026)** — Requisitions, Candidates, Applications, Interviews, Onboardings, Task Templates, Items. **33 OpenAPI endpoints** |
| **Reimbursement & Claim** | **3 GORM entities** | ✅ **48 tests** | ✅ **Completed (26 Jul 2026)** — Reimbursement Types, Requests, Items. **15 OpenAPI endpoints** |
| **Training & Development Management** | **7 GORM entities** | ✅ **31 tests** | ✅ **Completed (26 Jul 2026)** — Training categories, courses, sessions, participants, materials, evaluations, certificates. **35 OpenAPI endpoints** |
| **Workforce Intelligence & Strategic Planning** | **7 GORM entities** | ✅ **108 tests** | ✅ **Completed (25 Jul 2026)** — Strategic analytics layer with 11 submodules: Workforce Planning (headcounts, forecasts, gap analysis, projections), Intelligence (KPI dashboards), Analytics (9 dashboards: headcount, attendance, leave, overtime, payroll, performance, learning, recruitment, movement), Capacity (dashboard, utilization, forecast, bottlenecks), Cost (summary, payroll breakdown, per-employee, per-department, budget-vs-actual), Risk (dashboard, 4 risk monitors, indicators), Executive (8 widgets), Scenario Planning (7 operations), Organization Health (span of control, succession readiness), People Analytics (7 correlation analyses). **68 OpenAPI endpoints** |
| **Career Intelligence & Talent Management** | **4 GORM entities** | ✅ **65 tests** | ✅ **Completed (26 Jul 2026)** — Strategic talent analytics for 9-box talent mapping (CareerTalentMap), career interests tracking (CareerInterest), career path gap analysis (CareerPath), and succession planning (CareerSuccessionPlan). **19 OpenAPI endpoints** across 4 resource groups. Migration files 018 (4 tables). Full OpenAPI documentation (24th tag). |

---

## 2. Analisis Arsitektur & Rekomendasi Teknis (Gap Analysis)

Berdasarkan tinjauan arsitektur teknis, beberapa area kritis berikut harus diterapkan untuk menjamin keandalan sistem pada skala produksi (*production-ready*):

### 2.1. Tenant Lifecycle & Resource Cleanup ✅
* **Masalah:** Penanganan status tenant (*Suspend*, *Soft Delete*, *Terminate*) berisiko menyisakan koneksi TCP/GORM yang menggantung di memori aplikasi (`Manager.tenants map`).
* **Status:** ✅ **Sudah diimplementasikan** — `CloseTenantConnection(companyID string)` pada `database.Manager` sudah ada dan dipanggil di `DeactivateTenantConnection()`, `RemoveTenantConnection()`, `DropTenantDB()`, dan `CloseAll()`.

### 2.2. Keamanan Kredensial Database Tenant ✅
* **Masalah:** Kredensial koneksi database tenant pada skema Platform/Master tidak boleh disimpan dalam bentuk *plain text*.
* **Status:** ✅ **Sudah diimplementasikan** — Enkripsi AES-256-GCM pada layer repository.
  * Package: `internal/pkg/crypto/crypto.go` — `Encrypt()`/`Decrypt()` dengan 12-byte random nonce
  * `SaveTenantConnection()` — encrypt password sebelum INSERT ke `tenant_connections`
  * `FindTenantConnection()` — decrypt password setelah SELECT (dengan fallback untuk legacy plaintext)
  * CLI: `encrypt-passwords` untuk migrasi data legacy ke format terenkripsi
  * Key: 32-byte hex-encoded via env `HRIS_ENCRYPTION_KEY`

### 2.3. Optimasi Connection Pooling & Database Limit ✅
* **Masalah:** Alokasi koneksi maksimum yang terlalu besar per tenant (misal `SetMaxOpenConns(100)`) berisiko menghabiskan kuota `max_connections` pada host database jika jumlah tenant meningkat.
* **Status:** ✅ **Sudah diimplementasikan** — Platform dan tenant kini memiliki pool terpisah.

#### Arsitektur Pool
```
Platform Pool (single DB):
  max_open=10  max_idle=5  lifetime=1jam

Tenant Pool (per DB tenant):
  max_open=10  max_idle=3  lifetime=30mnt  idle_timeout=5mnt

PgBouncer (opsional, untuk production):
  mode=transaction  pool=10/tenant  max_client=500
```

#### Pool Math
| Skenario | Platform | 1 Tenant | 10 Tenants | 50 Tenants |
|----------|:--------:|:--------:|:----------:|:----------:|
| **Before** | 25 open | 25 open | 250 open | 1,250 open |
| **After** | 10 open | 10 open | 100 open | 500 open |
| **+PgBouncer** | — | — | 10 ke PG | 50 ke PG |

#### Implementasi
- **Pool terpisah**: `database.Config` sekarang memiliki field terpisah untuk platform (`MaxOpenConns`/`MaxIdleConns`) dan tenant (`TenantMaxOpenConns`/`TenantMaxIdleConns`/`TenantConnMaxLifetimeMs`/`TenantConnMaxIdleTimeMs`)
- **Idle eviction**: `SetConnMaxIdleTime(5m)` menutup koneksi idle yang tidak terpakai
- **PoolStats()**: Method baru untuk inspeksi pool secara real-time (lihat `GET /monitoring/pool`)
- **PgBouncer**: Service tersedia di `docker-compose.yml` dengan transaction mode, siap production

### 2.4. Penanganan Dialek SQL Migrasi ✅
* **Masalah:** Penggunaan eksekusi file `.sql` mentah memerlukan penanganan khusus jika sistem mendukung dual-driver (PostgreSQL & MySQL) karena perbedaan sintaksis DDL (seperti `UUID` vs `VARCHAR/CHAR(36)`).
* **Status:** ✅ **Sudah diimplementasikan** — Direktori migrasi dipisah: `migrations/tenant/mysql/` (30 file MySQL-optimized) dan `migrations/tenant/postgres/` (30 file PostgreSQL-optimized). Go code menggunakan `TenantRootPath(driver)` untuk memilih dialect yang sesuai secara otomatis saat provisioning.

### 2.5. Sinkronisasi Cache Terdistribusi ✅
* **Masalah:** Pembaruan *Feature Flags* atau *Permissions* di Redis pada satu instance Go server perlu dikonsumsi secara konsisten oleh instance lainnya pada lingkungan *multi-node deployment*.
* **Status:** ✅ **Sudah diimplementasikan** — Distributed cache dengan two-tier architecture:

### 2.6. Bilingual Support (Bahasa Indonesia & English) ✅
* **Masalah:** Response API perlu mendukung dua bahasa untuk menjangkau pengguna yang lebih luas (Indonesia dan internasional).
* **Status:** ✅ **Sudah diimplementasikan** — Bilingual messaging system dengan arsitektur:

#### Komponen Utama
1. **Middleware `Localize()`** (`internal/pkg/middleware/localize.go`) — auto-detect bahasa dari header `Accept-Language` dan simpan di `gin.Context` dengan key `"lang"`.
2. **Response Helpers** (`internal/pkg/httputil/response.go`) — `SuccessJSON`, `CreatedJSON`, `UpdatedJSON`, `DeletedJSON`, `ErrorJSON`, `NotFound`, dll. Semua helper otomatis membaca bahasa dari context tanpa perlu parsing ulang.
3. **Message Catalog** (`internal/pkg/httputil/locale.go`) — Map `localeMessages` berisi 80+ pasangan EN/ID untuk:
   - **Validation error** (required, email, min, max, uuid, dll)
   - **HTTP error** (not_found, unauthorized, forbidden, conflict, dll)
   - **Success action** (created, updated, deleted, approved, rejected, dll)
   - **Module-specific** (company, user, license, package, approval, employee, training, career, dll)
   - **Custom Indonesian validators** (NIK, NPWP, KK, Passport, SIM, No Rekening, Kode Pos, Phone ID)
4. **Validation Error Response** — Field-specific errors balik dalam array dengan key sesuai field:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": {
      "en": "Validation failed",
      "id": "Validasi gagal"
    },
    "errors": {
      "email": ["Must be a valid email address"],
      "nik": ["Format NIK tidak valid, harus 16 digit angka"]
    }
  }
}
```

#### Cara Kerja
1. Request masuk → `Localize()` middleware baca `Accept-Language` header
2. Bahasa disimpan di `gin.Context` sebagai `"lang"`
3. Semua response helper (`SuccessJSON`, `CreatedJSON`, `ErrorJSON`, `NotFound`, dll) panggil `tCtx(c, key)` untuk auto-translate
4. Jika key tidak ditemukan, fallback ke English

#### Penggunaan API
```bash
# English (default)
curl http://localhost:8080/api/v1/platform/companies/not-found

# Bahasa Indonesia
curl -H "Accept-Language: id" http://localhost:8080/api/v1/platform/companies/not-found
```
  * **Package:** `internal/pkg/cache/` — `cache.go` + `pubsub.go`
  * **Local cache:** `sync.Map` dengan TTL untuk akses sub-milidetik
  * **Shared store:** `go-redis/v9` untuk cache bersama antar instance
  * **Pub/Sub invalidation:** Semua instance subscribe ke channel `hris:cache:invalidate`. Saat satu instance mengubah data, semua instance lain langsung evict local cache.
  * **API:** `Get`, `Set`, `SetJSON`, `Invalidate`, `InvalidatePrefix`
  * **Monitoring:** `GET /health` mencakup status Redis cache via `Ping()`

---

## 3. Status & Roadmap Kelengkapan Modul HRIS

### 3.1. Modul Inti (Completed / Core Engine)
Modul-modul berikut telah selesai didesain dan diimplementasikan sebagai fondasi utama sistem:

* [x] **Platform & Tenant Management:** Provisioning DB multi-tenant, isolasi database, switching context tenant.
* [x] **Organization Management:** Organization Summary, Positions/Jabatan, Zone/Wilayah.
  * **Multi-Company Architecture:** Pengelolaan entitas Holding, Anak Perusahaan, dan Unit Bisnis.
  * **Dynamic Department Hierarchy:** Struktur pohon departemen/divisi (Adjacency List) berkedalaman fleksibel beserta pemetaan *Cost Center*.
  * **Location & Geofencing Zones:** Pengelolaan cabang, lokasi kerja, serta penentuan radius zonasi (*latitude/longitude*) untuk presensi.
  * **Organization Summary:** Aggregation metrics headcount per entitas/departemen.
  * 🗓️ **Organization History, Versioning & Cloning (Planned):**
    * **Partial & Total Change Capture:** Pencatatan delta perubahan skala kecil (audit log) dan reorganisasi skala besar.
    * **Full Structure Cloning:** Fitur *deep copy* seluruh pohon struktur ke versi *DRAFT* untuk simulasi reorganisasi.
    * **Version Audit Trail:** Perbandingan (*diff*) antar versi arsitektur organisasi secara historis.
* [x] **Employee Management:** Data personal, kontak, alamat, keluarga, pendidikan, dokumen, riwayat kerja, rekening/pajak, serta pengaturan akun.
* [x] **Job Management:** Deskripsi jabatan lengkap (*Responsibilities, Working Activities, Operational/HR Authorities, Working Risks, Title Subs, Assets, Financials*).
* [x] **Payroll & Compensation Engine:** **21 GORM entities** — SalaryComponent, PayrollPeriod, PayrollRun, EmployeePayrollProfile, EmployeeBankProfile, EmployeeBpjsProfile, EmployeeTaxProfile, BpjsSetting, BpjsRateComponent, Pph21Setting, Pph21PtkpRate, Pph21TaxBracket, PayrollRunEmployee, PayrollRunItem, PayrollPayslip, Pph21CalculationLog, dan lainnya. Full CRUD (Create/Read/Update/Delete) untuk seluruh entity. **34 unit tests** (13 repository + 21 service).
  * **BPJS Indonesia:** BPJS Settings, BPJS Rate Components (Kesehatan & Ketenagakerjaan), Employee BPJS Profiles.
  * **PPh21:** PPh21 Settings, Tax Brackets, PTKP Rates, Calculation Logs.
  * **Payroll Run:** Payroll Periods, Payroll Runs (dengan status workflow), Run Employees, Run Items, Payroll Payslips.
  * **Employee Payroll:** Employee Payroll Profiles, Bank Profiles, BPJS Profiles, Tax Profiles.
* [x] **Time & Attendance:** **10 GORM entities** — Company Settings, Company Shifts, Employee Shifts, Locations (Geofence), Device Captures, Face Captures, Events (Check-in/out), Sessions (Daily Work), Overtime Requests, Exempt Positions. Full CRUD untuk 8 sub-entities. **83 unit tests** (37 repository + 25 service + 21 handler).
* [x] **Competency Management:** ✅ **Completed (26 Jul 2026)** — 7 GORM entities, full CRUD. Competencies, Values (legacy + structured), Events, Targets, Scores, Score Details. **35 OpenAPI endpoints**, **54 unit tests**.
* [x] **Employee Movement & Career Management:** ✅ **Completed (26 Jul 2026)** — 2 GORM entities. 8 movement types (promotion, demotion, mutation, contract_extension, status_change, retirement, offboarding, other). 4 contract types (PKWT, PKWTT, daily, other). **15 OpenAPI endpoints**, **58 unit tests**.
* [x] **Reimbursement & Claim:** ✅ **Completed (26 Jul 2026)** — 3 GORM entities (ReimbursementType, ReimbursementRequest, ReimbursementItem). Multi-step approval flow (DRAFT->SUBMITTED->APPROVED->PAID). **15 OpenAPI endpoints**, **48 unit tests**.
* [x] **Training & Development Management:** ✅ **Completed (26 Jul 2026)** — 7 GORM entities (TrainingCategory, TrainingCourse, TrainingSession, TrainingParticipant, TrainingMaterial, TrainingEvaluation, TrainingCertificate). Full CRUD with validations. **35 OpenAPI endpoints**, **31 unit tests**.
* [x] **Workforce Intelligence & Strategic Planning:** ✅ **Completed (25 Jul 2026)** — Strategic analytics layer aggregating from all 14 tenant modules. 11 submodules: Workforce Planning (headcount planning, demand/supply/hiring forecasts, gap analysis, hiring/retirement/growth/budget projections), Workforce Intelligence (KPI dashboards with 13 metrics: total/active HC, attrition, turnover, retention, promotion, mobility, tenure, diversity, employment type), Workforce Analytics (9 dashboard endpoints: headcount, attendance, leave, overtime, payroll, performance, learning, recruitment, movement), Capacity & Utilization (dashboard, resource availability, utilization rate, forecast, bottleneck analysis), Workforce Cost Analytics (summary, payroll cost breakdown, per-employee, per-department, budget-vs-actual), Risk Dashboard (high turnover, retirement, contract expiry, high absenteeism risk monitors with level/score/threshold and recommendations), Executive Dashboard (8 widgets: summary, growth, cost trend, attrition trend, capacity, hiring progress, risk overview, health score), Forecast & Scenario Planning (7 CRUD + run/clone operations for growth/new branch/demand/budget/retirement simulations), Organization Health (span of control, manager ratio, succession readiness, composite health score), People Analytics (7 correlation analyses: training-vs-performance, overtime-vs-productivity, attendance-vs-performance, compensation-vs-turnover, source-vs-retention, career-progression, learning-effectiveness). **7 GORM entities, 68 OpenAPI endpoints, 108 unit tests.**
* [x] **Career Intelligence & Talent Management:** ✅ **Completed (26 Jul 2026)** — Strategic talent analytics module for 9-box talent mapping (CareerTalentMap), career interests tracking (CareerInterest), career path gap analysis (CareerPath), and succession planning (CareerSuccessionPlan). **4 GORM entities, 65 unit tests** (23 repository + 20 service + 22 handler). **19 OpenAPI endpoints** across 4 resource groups: Talent Maps (7), Career Interests (3), Career Paths (4), Succession Plans (5). Full OpenAPI documentation with 24 tags, migration files 018 (4 new tables).

### 3.2. Modul Operasional & Siklus Karier (Planned / Phase 2 Roadmap)
Untuk melengkapi cakupan *Full-Suite HRIS*, modul-modul operasional berikut masuk dalam skala prioritas pengembangan tahap berikutnya:

* [x] **Time & Attendance:** ✅ Perekaman presensi, penjadwalan *shift*, lembur (*overtime*), dan kalkulasi keterlambatan.
* [x] **Leave & Time Off:** ✅ Pengajuan cuti, sakit, izin, manajemen kuota cuti tahunan, dan *multi-level approval*.
* [x] **Performance Management:** ✅ **Completed (26 Jul 2026)** — KPI/OKR/BSC dengan 7 GORM entities, period-based evaluations, perspective scoring, dan 55 unit tests.
* [x] **Recruitment & Onboarding (ATS):** ✅ **Completed (26 Jul 2026)** — End-to-end ATS dengan 7 GORM entities, candidate pipeline, interview scheduling, onboarding workflows, dan 66 unit tests.
* [x] **Training & Development Management:** ✅ **Completed (26 Jul 2026)** — 7 GORM entities, training categories, courses, sessions, participants, materials, evaluations, certificates. 35 endpoints, 31 unit tests.
* [ ] **Organization History, Versioning & Cloning:** Change Capture, Full Structure Cloning, Version Audit Trail.

---

## 4. Matriks Prioritas Eksekusi

| Area | Komponen | Prioritas | Action Item Utama |
| :--- | :--- | :---: | :--- |
| **Security** | Tenant Credentials | ✅ Done | AES-256-GCM encrypt/decrypt via `internal/pkg/crypto/`, CLI `encrypt-passwords` untuk legacy. |
| **Database** | SQL Dialect | ✅ Done | Migrasi dipisah per dialect: `mysql/` dan `postgres/`, dipilih otomatis via `TenantRootPath(driver)`. |
| **Resource** | Lifecycle Tenant | ✅ Done | `CloseTenantConnection()` sudah terimplementasi dan terintegrasi di lifecycle management. |
| **Performance**| Connection Pool | ✅ Done | Platform pool (10/5/1jam) & Tenant pool (10/3/30mnt/5mnt) terpisah. `PoolStats()` + PgBouncer. |
| **Architecture**| Cache Sync | ✅ Done | Distributed cache (local sync.Map + Redis) + Pub/Sub invalidation via `internal/pkg/cache/`. |
| **Payroll**| Payroll & Compensation Engine | ✅ **Done** | 21 GORM entities, full CRUD, 34 unit tests. BPJS, PPh21, Payroll Run, Employee Payroll Profiles. |
| **Attendance**| Time & Attendance | ✅ **Done (26 Juli 2026)** | 10 GORM entities, full CRUD, 83 unit tests. Settings, Shifts, Locations, Events, Sessions, Overtime. |
| **Leave**| Leave & Time Off | ✅ **Done (26 Juli 2026)** | 6 GORM entities, full CRUD, 38 unit tests, 21 endpoints. Types, Accrual Policies, Reasons, Requests, Balances. |
| **Competency**| Competency Management | ✅ **Done (26 Jul 2026)** | 7 GORM entities, full CRUD, 54 unit tests. Competencies, Values, Events, Scores. **35 OpenAPI endpoints**. |
| **Movement**| Employee Movement & Career | ✅ **Done (26 Jul 2026)** | 2 GORM entities, 58 unit tests. 8 movement types, contract extension chain. **15 OpenAPI endpoints**. |
| **Reimbursement**| Reimbursement & Claim | ✅ **Done (26 Jul 2026)** | 3 GORM entities, 48 unit tests. Multi-step status flow. **15 OpenAPI endpoints**. |
| **Performance**| Performance Management | ✅ **Done (26 Jul 2026)** | 7 GORM entities, 55 unit tests, 34 endpoints. KPI/OKR/BSC framework with period-based evaluations. |
| **Recruitment**| Recruitment & Onboarding ATS | ✅ **Done (26 Jul 2026)** | 7 GORM entities, 66 unit tests, 33 endpoints. End-to-end ATS pipeline. |
| **Training**| Training & Development Management | ✅ **Done (26 Jul 2026)** | 7 GORM entities, 31 unit tests, 35 endpoints. Categories, courses, sessions, participants, materials, evaluations, certificates. |
| **Workforce Intelligence**| Workforce Intelligence & Strategic Planning | ✅ **Done (25 Jul 2026)** | 7 analytics entities, read-only analytics layer, 11 submodules (Planning, Intelligence, Analytics, Capacity, Cost, Risk, Executive, Scenarios, Health, People Analytics, AI Insights). 68 OpenAPI endpoints, 108 unit tests. |
| **Career Intelligence** | Career Intelligence & Talent Management | ✅ **Done (26 Jul 2026)** | 4 GORM entities (CareerTalentMap, CareerInterest, CareerPath, CareerSuccessionPlan), 65 unit tests, 19 OpenAPI endpoints. 9-box talent grid, career path gap analysis, succession planning. Migration 018 (4 tables). |

---
*Document Version: 1.6-Updated (v12)*  
*Status: Approved for Architecture Enhancement*

### Changelog

| Versi | Tanggal | Perubahan |
| v12 | 26 Jul 2026 | ✅ **Module Type Filter** ditambahkan. Kolom `module_type` (platform/tenant) di tabel modules. Filter `?module_type=` di 3 endpoints: GET /modules, GET /packages, GET /public/packages. Seed installer diperbaiki: slug `job-management`→`jobmanagement`, `employee-movement`→`employeemovement`. Semua `depends_on` 20 module telah diselaraskan dengan Info(). DTO dan model modulemgmt diperbarui. Migration 012_add_modules_module_type. |
| v11 | 26 Jul 2026 | ✅ **Career Intelligence & Talent Management** selesai. 4 GORM entities (CareerTalentMap, CareerInterest, CareerPath, CareerSuccessionPlan). 65 unit tests (23 repo + 20 service + 22 handler). 19 OpenAPI endpoints. Migration 018 (4 tables). OpenAPI full coverage (556 endpoints, 359 schemas, 26 tags). Go arch report: 116 tenant entities, 480 service methods, 1004 tests total. |
| v10 | 25 Jul 2026 | ✅ **Workforce Intelligence & Strategic Planning** selesai. 7 GORM entities (WorkforcePlanningHeadcount, WorkforceForecast, WorkforceKPI, WorkforceAnalyticsCache, WorkforceScenario, WorkforceRiskIndicator, WorkforceHealthScore). 68 OpenAPI endpoints, 108 unit tests (31 repository + 41 service + 36 handler). Migration files 001-017 (139 tables). OpenAPI full coverage (525 endpoints, 334 schemas, 23 tags). |
| v8 | 26 Jul 2026 | ✅ Training & Development Management (7 entities, 31 tests, 35 endpoints) selesai. OpenAPI full coverage (457 endpoints, 290 schemas, 22 tags). Migration files 001-016 (132 tables). |
| v7 | 26 Jul 2026 | ✅ Performance Management (7 entities, 55 tests, 34 endpoints) & Recruitment ATS (7 entities, 66 tests, 33 endpoints) selesai. OpenAPI full coverage (415 endpoints, 253 schemas, 21 tags). Migration files 001-015 (120 tables). |
| v6 | 26 Jul 2026 | ✅ Competency, Employee Movement, Reimbursement selesai. OpenAPI docs full coverage (348 endpoints, 208 schemas, 19 tags) |
| v5 | 26 Juli 2026 | ✅ Leave & Time Off selesai: 6 GORM entities, full CRUD, 38 unit tests, 23 endpoints |
| v4 | 26 Juli 2026 | ✅ Time & Attendance selesai: 10 GORM entities, full CRUD, 83 unit tests, 30 endpoints, OpenAPI docs |
| v3 | 24 Juli 2026 | ✅ Payroll & Compensation Engine selesai: 21 GORM entities, full CRUD, 34 unit tests, BPJS & PPh21 Indonesia |
| v2 | 22 Juli 2026 | ✅ Job Management selesai: 18 GORM entities, 74 unit tests, OpenAPI docs. Distribusi cache & migration fixes |
| v1 | 22 Juli 2026 | Initial architecture enhancement: Tenant lifecycle, crypto, connection pool, SQL dialect, distributed cache |
