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
| **Platform & Tenant Management** | 6 platform entities | ✅ | ✅ **Completed** — Company, License, Module Management, Users, Monitoring |
| **Organization Management** | 5 GORM entities | ❌ | ✅ **Completed** — Organization, Zone, Position, Level, Summary |
| **Employee Management** | 9+ entities | ✅ 25+ tests | ✅ **Completed** — Personal, kontak, alamat, keluarga, pendidikan, dokumen, dll |
| **Job Management** | 18 GORM entities | ✅ **74 tests** | ✅ **Completed** — Titles, Values, Responsibilities, Authorities, Working Activities/Risks |
| **Payroll & Compensation Engine** | **21 GORM entities** | ✅ **34 tests** | ✅ **Completed (24 Juli 2026)** — BPJS, PPh21, Payroll Run, Employee Profiles |
| **Employee Movement & Career** | 8 entities | ✅ | ✅ **Completed** — Promosi/Demosi, PKWT, Mutation, Pensiun |
| **Competency Management** | 7 entities (DB only) | ❌ | 🗄️ **DB Schema Only** — Go module belum diimplementasi |
| **Time & Attendance** | **10 GORM entities** | ✅ **83 tests** | ✅ **Completed (26 Juli 2026)** — Settings, Shifts, Employee Shifts, Locations, Events, Sessions, Overtime, Exempt Positions |
| **Leave & Time Off** | **6 GORM entities** | ✅ **38 tests** | ✅ **Completed (26 Juli 2026)** — Leave Types, Accrual Policies, Reasons, Requests (23 endpoints), Details, Balances |
| **Approval Engine** | **5 GORM entities** | ✅ **67 tests** | ✅ **Completed (24 Juli 2026)** — Flows, Steps, Instances, Actions, Tasks |

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
* **Status:** ✅ **Sudah diimplementasikan** — Direktori migrasi dipisah: `migrations/tenant/mysql/` (22 file MySQL-optimized) dan `migrations/tenant/postgres/` (22 file PostgreSQL-optimized). Go code menggunakan `TenantRootPath(driver)` untuk memilih dialect yang sesuai secara otomatis saat provisioning.

### 2.5. Sinkronisasi Cache Terdistribusi ✅
* **Masalah:** Pembaruan *Feature Flags* atau *Permissions* di Redis pada satu instance Go server perlu dikonsumsi secara konsisten oleh instance lainnya pada lingkungan *multi-node deployment*.
* **Status:** ✅ **Sudah diimplementasikan** — Distributed cache dengan two-tier architecture:
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
* [ ] **Competency Management:** 🗄️ **DB Schema Only** — Tabel database tersedia (DDL migration 008_competency.sql — 7 tabel), Go module belum diimplementasikan.

### 3.2. Modul Operasional & Siklus Karier (Planned / Phase 2 Roadmap)
Untuk melengkapi cakupan *Full-Suite HRIS*, modul-modul operasional berikut masuk dalam skala prioritas pengembangan tahap berikutnya:

* [ ] **Employee Movement & Career Management:**
  * **Promosi & Demosi:** Perubahan jenjang jabatan, posisi, dan penyesuaian kelas gaji.
  * **Perpanjangan Kontrak (PKWT):** Pengelolaan masa berlaku kerja, peringatan jatuh tempo, dan adendum kontrak.
  * **Pensiun & Offboarding/PHK:** Pengelolaan masa purna bakti, pemutusan hubungan kerja, dan integrasi perhitungan pesangon/uang jasa pada modul Payroll.
* [x] **Time & Attendance:** Perekaman presensi, penjadwalan *shift*, lembur (*overtime*), dan kalkulasi keterlambatan.
* [x] **Leave & Time Off:** Pengajuan cuti, sakit, izin, manajemen kuota cuti tahunan, dan *multi-level approval*.  
  * **6 GORM entities** — LeaveType, LeaveAccrualPolicy, LeaveReason, LeaveRequest, LeaveRequestDetail, EmployeeLeaveBalance. Full CRUD untuk 5 sub-entities. **38 unit tests** (14 repo + 12 service + 12 handler). **23 endpoints**. |
* [ ] **Performance Management:** Penilaian kinerja (KPI, OKR, review 360) yang terintegrasi langsung dengan modul *Job Management* dan *Competency*.
* [ ] **Reimbursement & Claim:** Pengajuan dan verifikasi klaim kesehatan maupun operasional dinas.
* [ ] **Recruitment & Onboarding (ATS):** Manajemen kandidat, alur seleksi, hingga otomatisasi pendaftaran karyawan baru (*onboarding*).

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

---
*Document Version: 1.6-Updated (v3)*  
*Status: Approved for Architecture Enhancement*

### Changelog

| Versi | Tanggal | Perubahan |
|-------|---------|-----------|
| v5 | 26 Juli 2026 | ✅ Leave & Time Off selesai: 6 GORM entities, full CRUD, 38 unit tests, 23 endpoints |
| v4 | 26 Juli 2026 | ✅ Time & Attendance selesai: 10 GORM entities, full CRUD, 83 unit tests, 30 endpoints, OpenAPI docs |
| v3 | 24 Juli 2026 | ✅ Payroll & Compensation Engine selesai: 21 GORM entities, full CRUD, 34 unit tests, BPJS & PPh21 Indonesia |
| v2 | 22 Juli 2026 | ✅ Job Management selesai: 18 GORM entities, 74 unit tests, OpenAPI docs. Distribusi cache & migration fixes |
| v1 | 22 Juli 2026 | Initial architecture enhancement: Tenant lifecycle, crypto, connection pool, SQL dialect, distributed cache |
