# 📄 Enterprise HRIS Blueprint & Architecture Specification
**Codebase:** Single Codebase (Go Engine + Vue 3 / PrimeVue 4)  
**Deployment Modes:** Subscription SaaS (Multi-Tenant) & On-Premise (Dedicated)  
**Access Model:** Hybrid Access Control (Platform License + Tenant DB RBAC)  
**Date:** July 2026  

---

## 📐 1. Executive Architecture Summary

Sistem HRIS dirancang menggunakan **Shared Application Engine + Multi-Tenant Database Isolation**. Aplikasi backend Go dan frontend Vue 3 berjalan menggunakan **satu codebase utama** yang dapat beradaptasi berdasarkan *environment configuration*.

┌─────────────────────────────────────────────────────────────────────────────┐
│                   SINGLE CODEBASE ENGINE (GO + VUE 3 / PRIMEVUE)            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [ MODE 1: SUBSCRIPTION SAAS ]              [ MODE 2: ON-PREMISE ]          │
│  • Single Vue 3 Frontend                    • Dedicated Vue 3 Frontend      │
│  • Shared Go Engine Cluster                 • Dedicated Go Engine Instance  │
│  • Multi-Tenant DBs (1 DB per Tenant)       • Single Local DB (Dedicated)   │
│  • License: Live via Platform Central DB    • License: Offline RSA (.lic)   │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                         HYBRID ACCESS CONTROL MODEL                         │
│  • Platform Level (Central DB / RSA): Lisensi Modul/Paket Langganan          │
│  • Tenant Level (Tenant DB): Identity (Users, Roles, Rules/Permissions)      │
└─────────────────────────────────────────────────────────────────────────────┘

---

## 🔐 2. Hybrid Access Control & Data Isolation Design

Sistem memisahkan tanggung jawab otorisasi menjadi dua tingkatan (*Separation of Concerns*):

### Level 1: Platform License Level (Modul & Paket)
* **Tanggung Jawab:** Platform Admin (Penyedia SaaS) / File Lisensi RSA.
* **Fungsi:** Mengontrol modul apa saja yang dibeli/diaktifkan oleh tenant (contoh: `attendance`, `payroll`, `leave`).
* **Mekanisme SaaS:** Disimpan di Platform Central DB dan di-cache pada Redis (`platform:tenant:{tenant_id}:modules`).
* **Mekanisme On-Premise:** Dibaca dari file lisensi terenkripsi `.lic` via RSA Public Key.

### Level 2: Tenant RBAC Level (Users, Roles, & Rules)
* **Tanggung Jawab:** Company Admin (Klien Tenant).
* **Fungsi:** Mengatur user internal, perolehan role (`HR Manager`, `Supervisor`), dan aturan *permission* (`payroll:view_slip`, `leave:approve`).
* **Mekanisme:** Disimpan penuh dan terisolasi di dalam **Database Tenant** masing-masing.

---

## ⚡ 3. Multi-Tenant Login & Multi-Account Handling

Untuk menangani skenario di mana satu alamat email yang sama terdaftar di beberapa tenant:

1. **Subdomain Isolation (Utama):** 
   * Requests dikirim via URL tenant (misal: `company-a.hris.com`).
   * Middleware Go (`TenantResolver`) membaca subdomain dan mengarahkan koneksi langsung ke **Database Tenant A**.
2. **Central Tenant Picker (Portal Login Terpusat):**
   * Jika user login via `app.hris.com`, sistem mengecek tabel mapping ringan di Central DB `global_user_tenants`.
   * Jika email berada di >1 tenant, aplikasi menampilkan *Tenant Selection Screen*.
3. **Contextual JWT Token:**
   * JWT Token menyimpan `tenant_id`, `user_id` spesifik tenant tersebut, dan array `permissions` aktif.

---

## 🛠️ 4. Action Plan & Technical Checklist

### Phase 1: Database & Migrator Refactoring (P0 - Immediate)
- [x] **Skema RBAC Tenant:** DDL migration SQL tersedia di `backend/internal/pkg/migrator/migrations/tenant/postgres/` (dan `mysql/`):
  - `users` → `022_users.sql` (+ `022_users.down.sql`)
  - `roles`, `permissions` → `011_settings.sql`
  - `role_permissions`, `user_roles` → `role_has_permissions`, `model_has_roles`, `model_has_permissions` (pola Spatie) di `011_settings.sql`
- [x] **Resolusi Migration Manual:** Skema tabel `insurances` dipindahkan ke file DDL migrasi resmi (`021_insurances.sql` + down, MySQL & PostgreSQL). Sebelumnya hanya dibuat via GORM AutoMigrate.
- [x] **Automated Seeder:** Fungsi Go `tenantseed.SeedTenantRBAC()` (`backend/internal/pkg/tenantseed/seed_rbac.go`) + `tenantseed.SeedTenantMasterData()` (`seed_data.go`) auto-seed master data `permissions` bawaan sistem (16 resource × 4 action = 64 permission, format `resource.action`) dan default roles (`Admin` full access, `Employee` view-only) saat provisioning tenant. Terintegrasi di **3 jalur**: `handleProvision` (CLI), command `seed-data`, dan `company.Service.provisionTenant()` (API — tenant dibuat via API juga auto-seed) — idempotent via deterministik UUID. *(Catatan: tidak via Asynq Worker karena worker belum ada — seed berjalan sinkron saat provisioning.)*

### Phase 2: Go Backend Auth & Middleware (P0 - Immediate)
- [x] **Dynamic DB Manager:** Konfigurasi `SetMaxOpenConns` dan `SetMaxIdleConns` pada `TenantDBManager` agar koneksi database hemat memori. *(Sudah diimplementasikan di `backend/internal/pkg/database/manager.go` — `SetMaxOpenConns`/`SetMaxIdleConns`/`SetConnMaxLifetime`/`SetConnMaxIdleTime` untuk platform & tenant + endpoint `/monitoring/pool` (PoolStats).)*
- [x] **2-Layer Guard Middleware:** *(Kedua layer selesai — lihat sub-item.)*
  - [x] `UserAuthMiddleware`: Validasi JWT & verifikasi permission internal user dari DB Tenant. *(Sudah: `middleware.AuthJWT` (JWT + claims), `middleware.TenantRequired` (propagate company_id/user_id ke context), `authz.NewMiddleware` (RBAC permission check).)*
  - [x] `PlatformLicenseMiddleware`: Validasi lisensi modul tenant via Redis (SaaS). *(Sudah: `middleware.LicenseMiddleware` — path→module slug mapping, cache Redis `hris:platform:license:modules:{companyID}` TTL 5 menit, fallback `modulemgmt.ListCompanyModules` via `companyModuleListerAdapter`, 403 `MODULE_NOT_LICENSED` jika modul tidak aktif; di-wire di router tenant antara TenantRequired & RBAC; cache di-invalidate saat subscribe/unsubscribe + saat Activate/DeactivateModule di modulemgmt service (titik choke semua perubahan modul). On-Premise `.lic` via `onpremise` reader — lihat Phase 3.)*
- [x] **JWT Payload Update:** Tambahkan field `tenant_id` dan `permissions` pada JWT claims. *(Sudah: `auth.Claims` memiliki `permissions` + `company_id` — konteks tenant diwakili `company_id`, bukan field bernama `tenant_id`; konsisten dengan `TenantRequired` yang routing via `company_id`.)*

### Phase 3: On-Premise License Engine (P1 - Short Term)
- [x] **RSA License Generator:** Buat utility CLI di Platform Admin untuk meng-generate file `.lic` terenkripsi (berisi `expires_at`, `allowed_modules`, `max_employees`). *(Sudah: `cmd/licensectl` — `gen-key` (RSA 2048 PKCS#8/PKIX) & `gen-lic` (payload JSON + RSA-SHA256 signature); package `internal/pkg/onpremise` (`GenerateKeyPair`, `SignLicense`, `WriteLicenseFile`). Contoh: `licensectl gen-key --out private.pem --pub public.pem` lalu `licensectl gen-lic --priv private.pem --company-id <uuid> --company "PT X" --expires 2027-12-31 --modules organization,employee,payroll --max-employees 500`.)*
- [x] **RSA License Reader:** Buat modul validator RSA di Go backend yang membaca file `.lic` lokal jika `DEPLOYMENT_MODE=ON_PREMISE`. *(Sudah: `onpremise.ReadLicenseFile`/`VerifyLicense` (verifikasi RSA-SHA256 + cek expired/tampered); config `license.deployment_mode` (`saas` default / `on_premise`, via `HRIS_LICENSE_DEPLOYMENT_MODE`); di-wire di main.go sebagai lister alternatif PlatformLicenseMiddleware (`onPremiseLister` — allowed_modules dari .lic menggantikan company_modules DB).)*
- [x] **Enforce `max_employees` di runtime (on-premise):** saat `POST /employees` di mode `on_premise`, jumlah employee saat ini dicek terhadap batas `.lic` dan ditolak (403 `QUOTA_EXCEEDED`) jika tercapai. *(Sudah: `employee.EmployeeQuotaChecker` + `ErrQuotaExceeded` (`internal/modules/employee/quota.go`); `Service.checkQuota` dipanggil di awal `Create`; `repo.CountEmployees` (soft-delete aware); di-wire via `newEmployeeModule` di main.go (`onPremiseQuotaChecker` dari `lic.MaxEmployees`); unit test 5 kasus di `service_test.go`.)*
- [x] **Audit jalur create employee (tidak ada bypass):** kuota terpusat di `Service.Create()` — **satu-satunya** tempat `&Employee{}` dibuat di jalur production (non-test) (`POST /employees`). Jalur lain terbukti TIDAK membuat Employee master sehingga tidak perlu kuota: payroll profiles (`CreateEmployeePayrollProfile`/`Bank`/`Bpjs`/`Tax` — sub-record employee existing), recruitment (`CreateEmployeeOnboarding` — record onboarding, module tidak impor employee module), attendance (`CreateEmployeeShift` — assignment), dan semua POST sub-record di module employee (addresses/families/educations/documents/dsb). Frontend hanya **1 caller** (`EmployeeForm.vue` `savePersonalData`) dengan toast bilingual `employee.quota_exceeded` + helper `getErrorCode()`. *(Audit 31 Jul 2026 — jalur masa depan seperti batch import otomatis kena kuota selama lewat `Service.Create()`; sebaiknya jangan tambahkan kuota ke sub-record karena alamat/keluarga tidak mengubah headcount.)*

### Phase 4: Vue 3 + PrimeVue Frontend Interceptors (P1 - Short Term)
- [x] **Axios Interceptor:** Injeksi header `X-Tenant-ID` secara otomatis dari subdomain URL pada setiap outgoing API request. *(Tidak diterapkan — **superseded**: backend tidak membaca header `X-Tenant-ID`; routing tenant dilakukan via `company_id` dari JWT claims (`AuthJWT` → `TenantRequired` → `GetCompanyID`). Menambahkan header ini hanya akan jadi dead code. Alternatif aktif: `Authorization: Bearer <token>` yang berisi `company_id` di claims, sudah di-inject oleh interceptor `api.js`.)*
- [ ] **Dynamic Menu Rendering (Pinia Store):**
  - [x] **Level 1 Filter:** Sembunyikan navigasi modul jika tidak ada dalam lisensi tenant. *(Sudah: store `activeModules` + `filterByModule()` di `Sidebar.vue` menyembunyikan item yang modulnya tidak aktif.)*
  - [x] **Level 2 Filter:** Sembunyikan tombol/sub-menu jika user tidak memiliki permission role yang sesuai. *(Sudah: `hasPermission()` di `stores/auth.js` — decode JWT payload `permissions` (array `resource.action`) tanpa dependency, dukung wildcard `*` / `resource.*`; diterapkan di `Sidebar.vue` (`filterByModule` kini juga cek `permission` tiap menu — Level 1 modul + Level 2 permission) dan aksi tombol CRUD di `Employees.vue` (`employee.create/update/delete`).)*

---

## 📦 5. On-Premise Distribution Template (`docker-compose.yml`)

```yaml
version: '3.8'

services:
  hris-backend:
    image: [registry.yourdomain.com/hris-backend:v1.6.0](https://registry.yourdomain.com/hris-backend:v1.6.0)
    restart: always
    environment:
      - DEPLOYMENT_MODE=ON_PREMISE
      - DB_HOST=hris-db
      - DB_NAME=hris_enterprise
    volumes:
      - ./license.lic:/etc/hris/license.lic:ro
    ports:
      - "8080:8080"
    depends_on:
      - hris-db

  hris-frontend:
    image: [registry.yourdomain.com/hris-frontend-vue:v1.6.0](https://registry.yourdomain.com/hris-frontend-vue:v1.6.0)
    restart: always
    ports:
      - "80:80"

  hris-db:
    image: postgres:16-alpine
    restart: always
    environment:
      - POSTGRES_DB=hris_enterprise
      - POSTGRES_USER=hris_admin
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata: