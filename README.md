# HRIS Platform

**Human Resource Information System — Enterprise Grade**

Modular monolith backend untuk platform HRIS enterprise dengan arsitektur multi-tenant, dibangun menggunakan **Go** + **Gin** + **GORM**.

---

## 📋 Daftar Isi

- [Arsitektur](#arsitektur)
- [Tech Stack](#tech-stack)
- [Struktur Proyek](#struktur-proyek)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Konfigurasi](#konfigurasi)
- [Environment Variables](#environment-variables)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Docker](#docker)
- [Testing](#testing)
- [Module SDK](#module-sdk)
- [Database Migration](#database-migration)
- [Tenant Provisioning](#tenant-provisioning-end-to-end-verified-)
- [Deployment (SaaS & On-Premise)](#deployment-saas-on-premise)
- [Roadmap](#roadmap)
- [Dokumentasi Lainnya](#dokumentasi-lainnya)

---

## 🏗️ Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    Go Modular Monolith                       │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐ │
│  │  Platform    │  │   Shared    │  │   Tenant Modules     │ │
│  │  Management  │  │   Kernel    │  │                      │ │
│  │              │  │             │  │  ┌──────────────────┐│ │
│  │ • Company    │  │ • Config    │  │  │ Organization     ││ │
│  │ • Module     │  │ • Database  │  │  │ Employee         ││ │
│  │ • License    │  │ • Auth/JWT  │  │  │ Attendance       ││ │
│  │ • User       │  │ • Middleware│  │  │ Leave & Time Off ││ │
│  │ • Feature    │  │ • Router   │  │  │ Payroll          ││ │
│  │   Flag       │  │ • Logger   │  │  │ Competency       ││ │
│  │              │  │ • Module   │  │  │ Job Management   ││ │
│  │              │  │   SDK      │  │  │ Approval         ││ │
│  │              │  │             │  │  │ Reimbursement    ││ │
│  └─────────────┘  └─────────────┘  │  └──────────────────┘│ │
│                                     └──────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
         │                    │                      │
         ▼                    ▼                      ▼
  ┌─────────────┐   ┌──────────────┐      ┌──────────────────┐
  │ Platform DB  │   │  Tenant 1 DB │      │   Tenant 2 DB    │
  │ (PostgreSQL) │   │  (PostgreSQL │      │   (MySQL)        │
  │  or MySQL)   │   │   or MySQL)  │      │    or MySQL)     │
  └─────────────┘   └──────────────┘      └──────────────────┘
```

### Design Principles

| Prinsip | Deskripsi |
|---|---|
| **Modular Monolith** | Satu binary Go, multiple internal modules dengan batasan tegas |
| **Multi-Tenant** | Database per tenant — isolasi data total antar company |
| **Module SDK** | Setiap modul mengikuti kontrak interface yang sama |
| **Multi-Driver DB** | Support PostgreSQL & MySQL, bisa berbeda per tenant |
| **API-First** | REST API dengan JWT authentication |
| **Security by Design** | JWT + RBAC + Audit Log + Rate Limiting |

---

## 🔬 Gap Analysis & Rekomendasi Arsitektur (Production Readiness)

Berdasarkan tinjauan arsitektur teknis, berikut area kritis yang harus diterapkan untuk menjamin keandalan sistem pada skala produksi:

### 1. Tenant Lifecycle & Resource Cleanup ✅
- **Masalah:** Penanganan status tenant (*Suspend*, *Soft Delete*, *Terminate*) berisiko menyisakan koneksi TCP/GORM yang menggantung di memori aplikasi.
- **Solusi:** ✅ **Sudah diimplementasikan** — `CloseTenantConnection(companyID)` pada `database.Manager` sudah ada dan dipanggil di `DeactivateTenantConnection()`, `RemoveTenantConnection()`, `DropTenantDB()`, dan `CloseAll()`.

### 2. Keamanan Kredensial Database Tenant ✅
- **Masalah:** Kredensial koneksi database tenant disimpan dalam bentuk *plain text*.
- **Status:** ✅ **Sudah diimplementasikan**
  - Package `internal/pkg/crypto/crypto.go` — AES-256-GCM encrypt/decrypt
  - `SaveTenantConnection()` encrypt password, `FindTenantConnection()` decrypt password
  - CLI `encrypt-passwords` untuk migrasi data legacy
  - Key 32-byte hex via env `HRIS_ENCRYPTION_KEY`

### 3. Optimasi Connection Pooling ✅
- **Masalah:** Alokasi `SetMaxOpenConns(100)` per tenant berisiko menghabiskan `max_connections` host database.
- **Status:** ✅ **Sudah diimplementasikan** — Platform dan tenant memiliki pool terpisah:

| Pool | max_open | max_idle | max_lifetime | max_idle_time |
|------|:--------:|:--------:|:------------:|:-------------:|
| Platform (single DB) | 10 | 5 | 1 jam | — |
| Tenant (per DB) | 10 | 3 | 30 menit | 5 menit |
| PgBouncer (optional) | 10/tenant | — | transaction mode | — |

**Pool math:** 50 tenant × 10 open = 500 koneksi (sebelumnya 1.250).
**PoolStats():** Endpoint `GET /monitoring/pool` untuk inspeksi real-time.

### 4. Penanganan Dialek SQL Migrasi ✅
- **Masalah:** Eksekusi file `.sql` mentah bermasalah untuk dual-driver (PostgreSQL vs MySQL) karena perbedaan sintaksis DDL.
- **Status:** ✅ **Sudah diimplementasikan** — Migrasi dipisah per dialect: `internal/pkg/migrator/migrations/tenant/mysql/` (86 file) dan `internal/pkg/migrator/migrations/tenant/postgres/` (86 file). Go code menggunakan `TenantRootPath(driver)` untuk seleksi otomatis saat provisioning. (File migrasi di-embed ke binary via `//go:embed migrations` — direktori filesystem `backend/migrations/` sudah dihapus; satu-satunya sumber adalah embed tree `internal/pkg/migrator/migrations/`.)

### 5. Sinkronisasi Cache Terdistribusi ✅
- **Masalah:** Pembaruan *Feature Flags* atau *Permissions* di Redis perlu dikonsumsi konsisten oleh semua instance server.
- **Status:** ✅ **Sudah diimplementasikan**
  - Two-tier cache: local `sync.Map` + shared `go-redis/v9`
  - Pub/Sub invalidation via channel `hris:cache:invalidate`
  - API: `Get`, `Set`, `SetJSON`, `Invalidate`, `InvalidatePrefix`
  - Monitoring: `GET /health` mencakup status Redis cache
  - Config: `cache.default_ttl`, `cache.key_prefix`
  - **Testing:** ✅ **81 test functions** — 42 unit tests (cache + PubSub) + 8 integration tests (full lifecycle, cross-instance, TTL, concurrent) + 31 benchmarks (throughput, latency, data size). Semua menggunakan `miniredis` — mock Redis in-memory tanpa dependensi eksternal.

### 6. Enkripsi & Masking Data Karyawan Sensitif ✅
- **Masalah:** Field karyawan sensitif (NIK, no. paspor, no. telepon, email, NIK keluarga, no. rekening bank, kontak darurat) tersimpan plain text di tenant DB dan tampil apa adanya ke siapa pun yang bisa mengakses data karyawan.
- **Status:** ✅ **Sudah diimplementasikan**
  - Encrypt-on-write (AES-256-GCM, `internal/pkg/crypto/`) untuk 8 field di 4 model (`employee`, `employee_family`, `employee_bank_account`, `emergency_contact`) — bisa di-toggle per field lewat halaman **Pengaturan Data Sensitif** (`setting.sensitive-fields.*`)
  - Registry field sensitif: `internal/modules/employee/sensitive_field_registry.go`
  - Masking tampilan (`internal/pkg/mask.PartialMask`, length-tiered) untuk pengguna tanpa permission lihat nilai asli; permission per-field `<resource>.view_<field>` (mis. `employee.view_nik`) — default hanya role Admin
  - Endpoint: `GET/PUT /api/v1/tenant/employees/settings/sensitive-fields`
  - Migrasi tenant 150–154 (tabel `sensitive_field_settings`, pelebaran kolom untuk ciphertext, permission RBAC)
  - Spec & plan: `docs/superpowers/specs/2026-08-16-sensitive-data-masking-design.md`, `docs/superpowers/plans/2026-08-16-sensitive-data-masking.md`

### Prioritas Eksekusi

| Area | Komponen | Prioritas | Action Item Utama |
| :--- | :--- | :---: | :--- |
| **Security** | Tenant Credentials | ✅ Done | AES-256-GCM encrypt/decrypt via `internal/pkg/crypto/`, CLI `encrypt-passwords` |
| **Security** | Sensitive Employee Fields | ✅ Done | Encrypt-on-write + role-based view masking untuk NIK/telepon/paspor/email/rekening bank dkk., toggle per field via halaman Pengaturan Data Sensitif |
| **Database** | SQL Dialect | ✅ Done | Migrasi dipisah per dialect: `mysql/` dan `postgres/`, dipilih otomatis via `TenantRootPath(driver)` |
| **Performance** | Connection Pool | ✅ Done | Platform pool (10/5/1jam) & Tenant pool (10/3/30mnt/5mnt) terpisah. PgBouncer + PoolStats() |
| **Architecture** | Cache Sync | ✅ Done | Distributed cache (local sync.Map + Redis) + Pub/Sub invalidation via `internal/pkg/cache/` |

---

## ⚙️ Tech Stack

| Layer | Teknologi |
|---|---|
| **Bahasa** | Go 1.22+ / Vue 3 (Composition API) |
| **HTTP Framework** | Gin (backend) / Axios (frontend) |
| **Frontend UI** | PrimeVue 4 + Tailwind CSS 4 + PrimeIcons |
| **Frontend State** | Pinia + Vue Router 4 |
| **ORM** | GORM (PostgreSQL & MySQL) |
| **Database** | PostgreSQL 16+ / MySQL 8+ |
| **Cache** | Redis 7+ |
| **Queue** | Asynq (Redis-based) |
| **Auth** | JWT (golang-jwt) |
| **Config** | Viper + godotenv |
| **Logging** | Zap (structured) |
| **Tracing** | OpenTelemetry (optional) |
| **Container** | Docker + Docker Compose |

---

## 📁 Struktur Proyek

```
hris-platform/
├── backend/                          # Go backend
│   ├── cmd/
│   │   ├── server/main.go            # Entry point API server
│   │   └── installer/main.go         # CLI installer (provisioning)
│   ├── internal/
│   │   ├── platform/                 # Platform Management Module
│   │   │   └── company/              #   Company CRUD + TenantConnection
│   │   ├── modules/                  # Tenant Modules
│   │   │   ├── organization/         #   Organization tree CRUD
│   │   │   ├── employee/             #   Employee CRUD + 8 sub-modules
│   │   │   ├── jobmanagement/        #   Job Management (18 entities)
│   │   │   ├── competency/           #   Competency Management (7 entities)
│   │   │   ├── employeemovement/     #   Employee Movement & Career Management (6 entities, 29 endpoints, 150 tests)
│   │   │   ├── attendance/           #   Time & Attendance (11 entities, 98 endpoints, 146 tests)
│   │   │   ├── approval/             #   Approval Engine (6 entities, 18 endpoints, 105 tests)
│   │   │   ├── payroll/              #   Payroll & Compensation Engine (21 entities, 83 endpoints)
│   │   │   ├── leave/                #   Leave & Time Off (7 entities, 26 endpoints, 79 tests)
│   │   │   ├── performance/          #   Performance Management (17 entities, 147 endpoints, 125 tests)
│   │   │   ├── recruitment/          #   Recruitment & Onboarding ATS (7 entities, 101 endpoints, 75 tests)
│   │   │   ├── training/            #   Training & Development (28 entities, 123 endpoints, 71 tests)
│   │   │   ├── careerintelligence/   #   Career Intelligence & Talent Management (5 entities, 23 endpoints, 66 tests)
│   │   │   ├── workforceintelligence/#   Workforce Intelligence & Strategic Planning (7 entities, 70 endpoints, 112 tests)
│   │   │   ├── reimbursement/        #   Reimbursement & Claim (3 entities, 15 endpoints, 67 tests)
│   │   │   └── setting/              #   Settings — 19 reference CRUDs (Zones, Provinces, Regencies, Districts, Villages, Educations, Education Majors, Religions, MaritalStatuses, RelationshipTypes, Banks, EmploymentStatuses, Nationalities, JobFamilies, SalaryGrades, TER, PTKP, Insurances, Company Holidays)
│   │   └── pkg/                      # Shared Kernel│       │   ├── config/               # Viper configuration loader
│       │   ├── database/             # Multi-tenant DB manager
│       │   ├── driver/               # Shared DB driver type
│       │   ├── auth/                 # JWT generation & validation
│       │   ├── authz/                # RBAC enforcer (database-backed) + CRUD API
│       │   ├── middleware/           # Gin middleware
│       │   ├── router/               # Router setup & module registration
│       │   ├── logger/               # Zap logger
│       │   ├── cache/                # Distributed two-tier cache + Pub/Sub
│       │   └── module/               # Module SDK interface
│   ├── config/
│   │   └── config.yaml               # Base configuration
│   └── internal/
│       └── pkg/
│           └── migrator/             # SQL Migration Engine
│               ├── migrator.go       #  Core: Up(), Down(), DownTo()
│               ├── embed.go          #  //go:embed migrations/
│               └── migrations/
│                   ├── platform/     #  Platform DDL (8 files)
│                   ├── seeders/      #  Seed data (2 file — 001 + down)
│                   └── tenant/       #  Tenant migrations (86 file per dialect — mysql & postgres)
│   ├── docker/
│   │   └── Dockerfile
│   ├── .env.example                  # Environment template
│   ├── Makefile
│   ├── go.mod
│   └── go.sum
├── frontend/                         # Vue 3 + PrimeVue (Platform Admin & Tenant)
├── docker/
│   └── docker-compose.yml            # Full infra compose
└── docs/
    ├── api/
    │   └── api-usage-guide.md        # Panduan penggunaan API (curl, auth, contoh)
    ├── platform-architecture-design.md  # Dokumen arsitektur utama
    ├── deployment-guide.md           # Panduan deployment SaaS & on-premise
    ├── openapi-report.md             # Laporan OpenAPI (generated via make docs)
    ├── panduan-uiux-hris-enterprise.md # Standar & prompt UI/UX
    ├── frontend-development-plan.md  # Roadmap implementasi frontend
    ├── ...                           # Lihat bagian "Dokumentasi Lainnya"
    ├── archive/                      # 🚫 DIABAIKAN — plan lama (jangan dibaca)
    ├── backlog/                      # 🚫 DIABAIKAN — ide/backlog
    ├── seeder/                       # 🚫 DIABAIKAN — file sementara
    └── tmp/                          # 🚫 DIABAIKAN — temp & node_modules
```

> ⚠️ **Folder yang diabaikan di `docs/`:** `archive/`, `backlog/`, `seeder/`, `tmp/` — jangan dibaca/diindeks (detail: [`docs/README.md`](docs/README.md#-folder-yang-diabaikan)).

---

## ✅ Prerequisites

- **Go** 1.22+
- **PostgreSQL** 16+ atau **MySQL** 8+
- **Redis** 7+
- **Docker** & Docker Compose (optional, untuk development environment)

---

## 🚀 Quick Start

### 1. Clone & Setup

```bash
git clone <repo-url> hris-platform
cd hris-platform/backend

# Copy environment
cp .env.example .env

# Edit .env sesuai environment lokal
nano .env

# Download dependencies
go mod download
```

> ⚠️ **Jangan commit file `.env` ke repository!** File ini berisi secret/kredensial lokal. Sudah seharusnya di-*gitignore*.

### 2. Database Setup

Buat database platform:

```bash
# PostgreSQL
createdb hris_platform

# atau via psql
psql -U postgres -c "CREATE DATABASE hris_platform;"
```

### 3. Run Server

```bash
# Development mode
make run

# atau langsung
go run ./cmd/server --config ./config/config.yaml
```

Server akan berjalan di `http://localhost:8080`.

### 4. Verify

```bash
curl http://localhost:8080/healthz
# {"status":"ok","service":"hris-platform"}
```

---

## 🔧 Konfigurasi

### Priority Loading

Konfigurasi di-load dengan urutan prioritas (low → high):

```
1. Default values (hardcoded)
2. config/config.yaml
3. .env file (via godotenv)
4. OS environment variables (HRIS_* prefix)
```

### Configuration File

```yaml
# backend/config/config.yaml
server:
  port: "8080"
  mode: "debug"               # debug | release | test

database:
  driver: "postgres"           # postgres | mysql
  platform_host: "localhost"
  platform_port: 5432
  platform_db: "hris_platform"
  # ... lihat config.yaml untuk lengkapnya
```

### Berganti ke MySQL

```yaml
database:
  driver: "mysql"
  platform_host: "localhost"
  platform_port: 3306
```

---

## 🌐 Environment Variables

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp backend/.env.example backend/.env
```

| Variable | Default | Deskripsi |
|---|---|---|
| `HRIS_SERVER_PORT` | `8080` | Port server |
| `HRIS_SERVER_MODE` | `debug` | Mode: debug/release/test |
| `HRIS_DATABASE_DRIVER` | `postgres` | Driver: postgres/mysql |
| `HRIS_DATABASE_PLATFORM_HOST` | `localhost` | Host platform DB |
| `HRIS_DATABASE_PLATFORM_PORT` | `5432` | Port platform DB |
| `HRIS_DATABASE_PLATFORM_DB` | `hris_platform` | Nama platform DB |
| `HRIS_REDIS_HOST` | `localhost` | Host Redis |
| `HRIS_REDIS_PORT` | `6379` | Port Redis |
| `HRIS_JWT_SECRET` | - | Secret key untuk JWT |
| `HRIS_JWT_ACCESS_TOKEN_TTL` | `15` | Access token expiry (menit) |
| `HRIS_JWT_REFRESH_TOKEN_TTL` | `24` | Refresh token expiry (jam) |
| `HRIS_LOGGER_LEVEL` | `info` | Log level |
| `HRIS_CORS_ALLOWED_ORIGINS` | `*` | CORS origins |

> Lihat [`.env.example`](backend/.env.example) untuk daftar lengkap.

---

## 📖 API Documentation

### Interactive Docs (Scalar UI)

Dokumentasi API interaktif tersedia setelah server berjalan:

```
http://localhost:8080/docs          # Scalar UI (interaktif)
http://localhost:8080/openapi.json  # OpenAPI 3.0 spec (JSON)
```

### Health Check

```bash
GET  /healthz          # Server health
GET  /readyz           # Readiness check
```

---

### 🔒 RBAC (Role-Based Access Control)

Platform menggunakan **database-backed RBAC** dengan **4 default role** dalam hierarki:

```
super_admin  (full access — platform + tenant)
  └── company_admin  (platform view + tenant full)
        └── manager  (tenant view/create/update)
              └── employee  (tenant view-only)
```

**Arsitektur Database-Backed:**

Role, permission, dan role-permission assignments disimpan di tabel platform:
- `rbac_roles` — daftar role dengan hierarki (parent_id)
- `rbac_permissions` — daftar permission (resource + action)
- `rbac_role_permissions` — many-to-many assignment

Saat startup, `NewEnforcerFromDB()`:
1. **Seed defaults** jika tabel kosong (4 role + 74+ permission dari 14 resource)
2. **Load ke memori** sebagai map[role]map[resource]actions
3. Siap digunakan untuk Check() — sub-milidetik, tanpa query DB per request

Setelah perubahan role/permission via API, enforcer di-**Reload()** otomatis — tanpa restart server.

**Inheritance behavior:**
- Jika suatu role **tidak memiliki** policy untuk resource → inherit dari parent
- Jika suatu role **memiliki** policy (dengan daftar action terbatas) → action yang tidak tercantum = intentional deny
- Contoh: `employee` punya `organization: view` → `create` akan **denied**, tidak inherit dari `manager`

#### Default Permission Matrix

| Resource | super_admin | company_admin | manager | employee |
|----------|:-----------:|:-------------:|:-------:|:--------:|
| `company` | ✅ * | ✅ view | ✅ view (via inherit) | ✅ view (via inherit) |
| `user` | ✅ * | ✅ view | ✅ view (via inherit) | ✅ view (via inherit) |
| `module` | ✅ * | ❌ | ✅ view (via inherit) | ✅ view (via inherit) |
| `license` | ✅ * | ✅ view | ✅ view (via inherit) | ✅ view (via inherit) |
| `monitoring` | ✅ * | ❌ | ❌ | ❌ |
| `organization` | ✅ * | ✅ * | ✅ V/C/U | ✅ V |
| `employee` | ✅ * | ✅ * | ✅ V/C/U | ✅ V |
| `attendance` | ✅ * | ✅ * | ✅ V | ✅ V |
| `leave` | ✅ * | ✅ * | ✅ V/C | ✅ V |
| `payroll` | ✅ * | ✅ * | ❌ | ✅ V |
| `competency` | ✅ * | ✅ * | ✅ V/C/U | ✅ V |
| `jobmanagement` | ✅ * | ✅ * | ✅ V/C/U | ❌ |
| `approval` | ✅ * | ✅ * | ✅ V/C/U | ❌ |
| `setting` | ✅ * | ✅ * | ✅ V/C/U | ✅ V |

> * = all actions, V = view, C = create, U = update

**Permission format:** `resource.action` (contoh: `company.create`, `user.view`)

Endpoint yang tidak memiliki akses akan mengembalikan **403 Forbidden**:

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to perform this action",
    "details": {
      "role": "company_admin",
      "resource": "company",
      "action": "create"
    }
  }
}
```

#### RBAC Management API (super_admin only)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/platform/rbac/roles` | List all roles with permissions |
| `POST` | `/platform/rbac/roles` | Create a new role |
| `GET` | `/platform/rbac/roles/:id` | Get role detail |
| `PUT` | `/platform/rbac/roles/:id` | Update role |
| `DELETE` | `/platform/rbac/roles/:id` | Delete role (non-system only) |
| `GET` | `/platform/rbac/permissions` | List all permissions |
| `POST` | `/platform/rbac/permissions` | Create a new permission |
| `DELETE` | `/platform/rbac/permissions/:id` | Delete permission (non-system only) |
| `POST` | `/platform/rbac/roles/:id/permissions` | Assign permission to role |
| `DELETE` | `/platform/rbac/roles/:id/permissions/:permissionId` | Revoke permission from role |

Semua perubahan role/permission akan otomatis me-reload enforcer (`Service.Sync()` → `Enforcer.Reload()`).

---

### Platform API (`/api/v1/platform`)

> 💡 **Auth required**: ✅ = JWT + RBAC, 🔓 = Public

#### 🔐 Auth

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `POST` | `/login` | Login platform admin | 🔓 Public | - |
| `POST` | `/refresh` | Refresh access token | 🔓 Public | - |

#### 👥 Users

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/users` | List platform users | ✅ | `user.view` (super_admin, company_admin) |
| `POST` | `/users` | Create platform user | ✅ | `user.create` (super_admin only) |
| `GET` | `/users/:id` | Get user detail | ✅ | `user.view` (super_admin, company_admin) |
| `PUT` | `/users/:id` | Update user | ✅ | `user.update` (super_admin only) |

#### 🏢 Companies — Tenant Lifecycle Management

| Method | Endpoint | Deskripsi | Auth | RBAC | Cleanup? |
|---|---|---|---|---|---|
| `POST` | `/companies` | Create company + provision tenant DB 🔄 | ✅ | `company.create` (super_admin only) | ❌ |
| `GET` | `/companies` | List all companies | ✅ | `company.view` (super_admin, company_admin) | ❌ |
| `GET` | `/companies/:id` | Get company detail | ✅ | `company.view` (super_admin, company_admin) | ❌ |
| `PUT` | `/companies/:id` | Update company info | ✅ | `company.update` (super_admin only) | ❌ |
| `DELETE` | `/companies/:id` | **Soft delete** — deactivate connection + deleted_at | ✅ | `company.delete` (super_admin only) | ✅ Deactivate |
| `POST` | `/companies/:id/suspend` | **Suspend** — deactivate connection, set status suspended | ✅ | `company.suspend` (super_admin only) | ✅ Deactivate |
| `POST` | `/companies/:id/activate` | **Activate** — reactivate connection, set status active | ✅ | `company.activate` (super_admin only) | ✅ Reactivate |
| `POST` | `/companies/:id/terminate` | **Terminate** — drop database + remove connection, set status terminated | ✅ | `company.terminate` (super_admin only) | ✅ Drop DB |
| `POST` | `/companies/:id/backup` | Trigger tenant backup (Phase 2) | ✅ | `company.backup` (super_admin only) | ❌ |
| `POST` | `/companies/:id/restore` | Trigger tenant restore (Phase 2) | ✅ | `company.restore` (super_admin only) | ❌ |

**Lifecycle Actions vs Database Cleanup:**

| Action | Company Status | Tenant Connection | Tenant Database |
|--------|---------------|-------------------|-----------------|
| `Create` | `active` | `is_active = true` | ✅ Created + migrated (199 tables) |
| `Suspend` | `suspended` | `is_active = false` + cache cleared | ✅ Data preserved |
| `Activate` | `active` | `is_active = true` + cache cleared | ✅ Reconnected |
| `Soft Delete` | (hidden via `deleted_at`) | `is_active = false` + cache cleared | ✅ Data preserved |
| `Terminate` | `terminated` | ❌ Removed from DB | ❌ Dropped |

#### 🧩 Modules

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/modules?module_type=` | List all modules (filter: `?module_type=platform` atau `tenant`) | ✅ | `module.view` (super_admin only) |
| `POST` | `/modules` | Register new module | ✅ | `module.create` (super_admin only) |
| `GET` | `/modules/:id` | Get module detail | ✅ | `module.view` (super_admin only) |
| `PUT` | `/modules/:id` | Update module | ✅ | `module.update` (super_admin only) |
| `GET` | `/modules/:id/companies` | List companies using module | ✅ | `module.view` (super_admin only) |
| `POST` | `/modules/:id/activate` | Activate module for company | ✅ | `module.activate` (super_admin only) |
| `POST` | `/modules/:id/deactivate` | Deactivate module for company | ✅ | `module.deactivate` (super_admin only) |

#### 🔑 Licenses

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/licenses` | List all licenses | ✅ | `license.view` (super_admin, company_admin) |
| `POST` | `/licenses` | Create license for company | ✅ | `license.create` (super_admin only) |
| `GET` | `/licenses/:id` | Get license detail | ✅ | `license.view` (super_admin, company_admin) |
| `PUT` | `/licenses/:id` | Update license | ✅ | `license.update` (super_admin only) |
| `DELETE` | `/licenses/:id` | Soft-delete license (deactivates modules) | ✅ | `license.delete` (super_admin only) |

#### 📊 Monitoring

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/monitoring/health` | Platform health + DB status | ✅ | `monitoring.view` (super_admin only) |
| `GET` | `/monitoring/tenants` | List tenant connections health | ✅ | `monitoring.view` (super_admin only) |
| `GET` | `/monitoring/tenants/:id` | Get tenant health detail | ✅ | `monitoring.view` (super_admin only) |

#### 📦 Package Management

| Method | Endpoint | Deskripsi | Auth | RBAC |
|---|---|---|---|---|
| `GET` | `/packages?module_type=` | List all packages (filter: `?module_type=platform` atau `tenant`) | ✅ | `package.view` (super_admin) |
| `POST` | `/packages` | Create package with modules | ✅ | `package.create` (super_admin only) |
| `GET` | `/packages/:id` | Get package detail | ✅ | `package.view` (super_admin) |
| `PUT` | `/packages/:id` | Update package | ✅ | `package.update` (super_admin only) |
| `DELETE` | `/packages/:id` | Soft-delete package | ✅ | `package.delete` (super_admin only) |
| `POST` | `/packages/:id/publish` | Publish package (validate dependencies) | ✅ | `package.update` (super_admin only) |
| `POST` | `/packages/:id/unpublish` | Unpublish package | ✅ | `package.update` (super_admin only) |
| `GET` | `/packages/:id/validate` | Validate module dependencies | ✅ | `package.view` (super_admin) |

**Public endpoint (no auth):**

| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/api/v1/public/packages?module_type=` | List published packages (filter: `?module_type=platform` atau `tenant`) |

**Package-License Integration:**
- Setiap package dapat di-assign ke license melalui field `package_id`
- Saat create company dengan `package_id`, license auto-created
- Package memiliki validasi dependensi modul (`depends_on`) — semua dependensi wajib disertakan dalam package yang sama

### Tenant API (`/api/v1/tenant`)

Semua endpoint tenant memerlukan **JWT Bearer Token** di header:
```
Authorization: Bearer <access_token>
```

| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/organizations` | List organizations |
| `POST` | `/organizations` | Create organization |
| `GET` | `/organizations/:id` | Get organization |
| `PUT` | `/organizations/:id` | Update organization |
| `DELETE` | `/organizations/:id` | Delete organization |
| `GET` | `/organizations?tree=true` | Get organization tree |
| `GET` | `/employees` | List employees (pagination) |
| `POST` | `/employees` | Create employee |
| `GET` | `/employees/:id` | Get employee with all sub-modules |
| `PUT` | `/employees/:id` | Update employee |
| `DELETE` | `/employees/:id` | Delete employee (hard delete) |

**Competency Management (8.1 — Competencies)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competencies` | List competencies (pagination) |
| `POST` | `/competencies` | Create competency |
| `GET` | `/competencies/:id` | Get competency |
| `PUT` | `/competencies/:id` | Update competency |
| `DELETE` | `/competencies/:id` | Delete competency |

**Competency Management (8.2 — Competence Values — legacy)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competence-values` | List competence values |
| `POST` | `/competence-values` | Create competence value |
| `GET` | `/competence-values/:id` | Get competence value |
| `PUT` | `/competence-values/:id` | Update competence value |
| `DELETE` | `/competence-values/:id` | Delete competence value |

**Competency Management (8.3 — Competency Values — structured)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competency-values` | List competency values |
| `POST` | `/competency-values` | Create competency value |
| `GET` | `/competency-values/:id` | Get competency value |
| `PUT` | `/competency-values/:id` | Update competency value |
| `DELETE` | `/competency-values/:id` | Delete competency value |

**Competency Management (8.4 — Competency Events)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competency-events` | List competency events |
| `POST` | `/competency-events` | Create competency event |
| `GET` | `/competency-events/:id` | Get competency event |
| `PUT` | `/competency-events/:id` | Update competency event |
| `DELETE` | `/competency-events/:id` | Delete competency event |

**Competency Management (8.5 — Competency Event Targets)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competency-event-targets` | List event targets |
| `POST` | `/competency-event-targets` | Create event target |
| `GET` | `/competency-event-targets/:id` | Get event target |
| `PUT` | `/competency-event-targets/:id` | Update event target |
| `DELETE` | `/competency-event-targets/:id` | Delete event target |

**Competency Management (8.6 — Competency Scores)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competency-scores` | List competency scores |
| `POST` | `/competency-scores` | Create competency score |
| `GET` | `/competency-scores/:id` | Get competency score (with details) |
| `PUT` | `/competency-scores/:id` | Update competency score |
| `DELETE` | `/competency-scores/:id` | Delete competency score |

**Competency Management (8.7 — Competency Score Details)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/competency-score-details` | List score details by score ID |
| `POST` | `/competency-score-details` | Create score detail |
| `GET` | `/competency-score-details/:id` | Get score detail |
| `PUT` | `/competency-score-details/:id` | Update score detail |
| `DELETE` | `/competency-score-details/:id` | Delete score detail |
| `POST` | `/employees/:id/addresses` | Create address (type: MAIN/DOMICILE) |
| `PUT` | `/employees/:id/addresses/:addressId` | Update address |
| `DELETE` | `/employees/:id/addresses/:addressId` | Delete address |
| `POST` | `/employees/:id/emergency-contacts` | Create emergency contact |
| `PUT` | `/employees/:id/emergency-contacts/:contactId` | Update emergency contact |
| `DELETE` | `/employees/:id/emergency-contacts/:contactId` | Delete emergency contact |
| `POST` | `/employees/:id/families` | Create family member |
| `PUT` | `/employees/:id/families/:familyId` | Update family member |
| `DELETE` | `/employees/:id/families/:familyId` | Delete family member |
| `POST` | `/employees/:id/educations` | Create education record |
| `PUT` | `/employees/:id/educations/:educationId` | Update education |
| `DELETE` | `/employees/:id/educations/:educationId` | Delete education |
| `POST` | `/employees/:id/experiences` | Create work experience |
| `PUT` | `/employees/:id/experiences/:experienceId` | Update experience |
| `DELETE` | `/employees/:id/experiences/:experienceId` | Delete experience |
| `POST` | `/employees/:id/documents` | Create document |
| `PUT` | `/employees/:id/documents/:documentId` | Update document |
| `DELETE` | `/employees/:id/documents/:documentId` | Delete document |
| `POST` | `/employees/:id/insurances` | Create insurance (BPJS) |
| `PUT` | `/employees/:id/insurances/:insuranceId` | Update insurance |
| `DELETE` | `/employees/:id/insurances/:insuranceId` | Delete insurance |
| `POST` | `/employees/:id/employments` | Create employment record |
| `PUT` | `/employees/:id/employments/:employmentId` | Update employment |
| `DELETE` | `/employees/:id/employments/:employmentId` | Delete employment |

**Job Management — Titles (9.1)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/titles` | List job titles |
| `POST` | `/job-management/titles` | Create job title |
| `GET` | `/job-management/titles/:id` | Get job title with subs |
| `PUT` | `/job-management/titles/:id` | Update job title |
| `DELETE` | `/job-management/titles/:id` | Delete job title |
| `POST` | `/job-management/titles/:id/subs` | Create job title sub |
| `GET` | `/job-management/titles/:id/subs` | List subs under a title |
| `GET` | `/job-management/titles/:id/subs/:subId` | Get job title sub |
| `PUT` | `/job-management/titles/:id/subs/:subId` | Update job title sub |
| `DELETE` | `/job-management/titles/:id/subs/:subId` | Delete job title sub |

**Job Management — Values (9.3)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/values` | List job values |
| `POST` | `/job-management/values` | Create job value |
| `GET` | `/job-management/values/tree` | Get job values tree (type_group → types → level options) |
| `GET` | `/job-management/values/clusters/:type` | List cluster mapping (`technical`/`managerial`) |
| `PUT` | `/job-management/values/clusters/:type` | Update cluster mapping |
| `GET` | `/job-management/values/:id` | Get job value |
| `PUT` | `/job-management/values/:id` | Update job value |
| `DELETE` | `/job-management/values/:id` | Delete job value |

**Job Management — Objectives (9.4)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/objectives` | List job objectives |
| `POST` | `/job-management/objectives` | Create job objective |
| `GET` | `/job-management/objectives/:id` | Get job objective |
| `PUT` | `/job-management/objectives/:id` | Update job objective |
| `DELETE` | `/job-management/objectives/:id` | Delete job objective |

**Job Management — Identifications (9.5)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/identifications` | List job identifications |
| `POST` | `/job-management/identifications` | Create job identification |
| `GET` | `/job-management/identifications/:id` | Get job identification |
| `PUT` | `/job-management/identifications/:id` | Update job identification |
| `DELETE` | `/job-management/identifications/:id` | Delete job identification |

**Job Management — Responsibilities (9.6)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/responsibilities` | List job responsibilities |
| `POST` | `/job-management/responsibilities` | Create job responsibility |
| `GET` | `/job-management/responsibilities/:id` | Get job responsibility |
| `PUT` | `/job-management/responsibilities/:id` | Update job responsibility |
| `DELETE` | `/job-management/responsibilities/:id` | Delete job responsibility |

**Job Management — Education Experiences (9.7)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/education-experiences` | List education experiences |
| `POST` | `/job-management/education-experiences` | Create education experience |
| `GET` | `/job-management/education-experiences/:id` | Get education experience |
| `PUT` | `/job-management/education-experiences/:id` | Update education experience |
| `DELETE` | `/job-management/education-experiences/:id` | Delete education experience |

**Job Management — HR Authorities (9.8)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/hr-authorities` | List HR authorities |
| `POST` | `/job-management/hr-authorities` | Create HR authority |
| `GET` | `/job-management/hr-authorities/:id` | Get HR authority |
| `PUT` | `/job-management/hr-authorities/:id` | Update HR authority |
| `DELETE` | `/job-management/hr-authorities/:id` | Delete HR authority |

**Job Management — Operational Authorities (9.9)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/operational-authorities` | List operational authorities |
| `POST` | `/job-management/operational-authorities` | Create operational authority |
| `GET` | `/job-management/operational-authorities/:id` | Get operational authority |
| `PUT` | `/job-management/operational-authorities/:id` | Update operational authority |
| `DELETE` | `/job-management/operational-authorities/:id` | Delete operational authority |

**Job Management — Working Activities (9.10)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/working-activities` | List working activities |
| `POST` | `/job-management/working-activities` | Create working activity |
| `GET` | `/job-management/working-activities/:id` | Get working activity |
| `PUT` | `/job-management/working-activities/:id` | Update working activity |
| `DELETE` | `/job-management/working-activities/:id` | Delete working activity |

**Job Management — Working Risks (9.11)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/working-risks` | List working risks |
| `POST` | `/job-management/working-risks` | Create working risk |
| `GET` | `/job-management/working-risks/:id` | Get working risk |
| `PUT` | `/job-management/working-risks/:id` | Update working risk |
| `DELETE` | `/job-management/working-risks/:id` | Delete working risk |

**Job Management — Relationships (9.12)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/relationships` | List job relationships |
| `POST` | `/job-management/relationships` | Create job relationship |
| `GET` | `/job-management/relationships/:id` | Get job relationship |
| `PUT` | `/job-management/relationships/:id` | Update job relationship |
| `DELETE` | `/job-management/relationships/:id` | Delete job relationship |
| `POST` | `/job-management/relationships/:id/details` | Create relationship detail |
| `GET` | `/job-management/relationships/:id/details` | List relationship details |
| `GET` | `/job-management/relationships/:id/details/:detailId` | Get relationship detail |
| `PUT` | `/job-management/relationships/:id/details/:detailId` | Update relationship detail |
| `DELETE` | `/job-management/relationships/:id/details/:detailId` | Delete relationship detail |

**Job Management — Subordinate Controls (9.13)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/subordinate-controls` | List subordinate controls |
| `POST` | `/job-management/subordinate-controls` | Create subordinate control |
| `GET` | `/job-management/subordinate-controls/:id` | Get subordinate control |
| `PUT` | `/job-management/subordinate-controls/:id` | Update subordinate control |
| `DELETE` | `/job-management/subordinate-controls/:id` | Delete subordinate control |

**Job Management — Assets (9.14)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/assets` | List job assets |
| `POST` | `/job-management/assets` | Create job asset |
| `GET` | `/job-management/assets/:id` | Get job asset |
| `PUT` | `/job-management/assets/:id` | Update job asset |
| `DELETE` | `/job-management/assets/:id` | Delete job asset |

**Job Management — Financials (9.15)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/financials` | List job financials |
| `POST` | `/job-management/financials` | Create job financial |
| `GET` | `/job-management/financials/:id` | Get job financial |
| `PUT` | `/job-management/financials/:id` | Update job financial |
| `DELETE` | `/job-management/financials/:id` | Delete job financial |

**Job Management — Potency Competencies (9.16)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/potency-competencies` | List potency competencies |
| `POST` | `/job-management/potency-competencies` | Create potency competency |
| `GET` | `/job-management/potency-competencies/:id` | Get potency competency |
| `PUT` | `/job-management/potency-competencies/:id` | Update potency competency |
| `DELETE` | `/job-management/potency-competencies/:id` | Delete potency competency |

**Job Management — Scores (9.17)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/job-management/scores` | List all scores |
| `GET` | `/job-management/scores/org/:orgId` | Get score by organization |
| `PUT` | `/job-management/scores/org/:orgId` | Upsert/recalculate score for organization (body kosong = hitung ulang otomatis) |

**Job Management — Competency Groups (9.18)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/job-management/competency-groups` | Create competency group |
| `GET` | `/job-management/competency-groups` | List competency groups |
| `GET` | `/job-management/competency-groups/:id` | Get competency group |
| `PUT` | `/job-management/competency-groups/:id` | Update competency group |
| `DELETE` | `/job-management/competency-groups/:id` | Delete competency group |

**Employee Movement & Career Management — Movements**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/employee-movements/movements` | List movements (pagination) |
| `POST` | `/employee-movements/movements` | Create movement (promotion, demotion, mutation, etc.) |
| `GET` | `/employee-movements/movements/:id` | Get movement by ID |
| `PUT` | `/employee-movements/movements/:id` | Update movement (draft only) |
| `DELETE` | `/employee-movements/movements/:id` | Delete movement (draft only) |
| `POST` | `/employee-movements/movements/:id/approve` | Approve movement (draft → approved) |
| `POST` | `/employee-movements/movements/:id/execute` | Execute movement (approved → executed) |
| `POST` | `/employee-movements/movements/:id/cancel` | Cancel movement (draft or approved only) |
| `GET` | `/employee-movements/employees/:employeeId/movements` | List movements by employee |

**Employee Movement & Career Management — Contracts**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/employee-movements/contracts` | List contracts (pagination) |
| `POST` | `/employee-movements/contracts` | Create contract (PKWT/PKWTT/Daily) |
| `GET` | `/employee-movements/contracts/:id` | Get contract by ID |
| `PUT` | `/employee-movements/contracts/:id` | Update contract |
| `DELETE` | `/employee-movements/contracts/:id` | Delete contract |
| `GET` | `/employee-movements/employees/:employeeId/contracts` | List contracts by employee |

**Reimbursement & Claim — Reimbursement Types**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/reimbursements/types` | List reimbursement types (pagination) |
| `POST` | `/reimbursements/types` | Create reimbursement type |
| `GET` | `/reimbursements/types/:id` | Get reimbursement type by ID |
| `PUT` | `/reimbursements/types/:id` | Update reimbursement type |
| `DELETE` | `/reimbursements/types/:id` | Delete reimbursement type |

**Reimbursement & Claim — Reimbursement Requests**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/reimbursements/requests` | List reimbursement requests (pagination) |
| `POST` | `/reimbursements/requests` | Create reimbursement request |
| `GET` | `/reimbursements/requests/:id` | Get reimbursement request by ID |
| `PUT` | `/reimbursements/requests/:id` | Update reimbursement request |
| `PUT` | `/reimbursements/requests/:id/status` | Update request status (submit/approve/reject/pay/cancel) |
| `DELETE` | `/reimbursements/requests/:id` | Delete reimbursement request |

**Reimbursement & Claim — Reimbursement Items**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/reimbursements/requests/:requestId/items` | List items attached to a request |
| `POST` | `/reimbursements/requests/:requestId/items` | Add a new item to a reimbursement request |
| `PUT` | `/reimbursements/requests/:requestId/items/:itemId` | Update a reimbursement item |
| `DELETE` | `/reimbursements/requests/:requestId/items/:itemId` | Delete a reimbursement item |

**Performance Management — KPI & OKR**

Semua endpoint performance berada di `/api/v1/tenant/performance/` dan terbagi dalam **3 kelompok**: **Master Data** (dipakai bersama KPI & OKR), **KPI** (`/performance/kpi/*`, 58 endpoints), dan **OKR** (`/performance/okr/*`, 42 endpoints).

**Master Data (dipakai KPI & OKR)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` / `GET` | `/performance/periods` | Create / List performance periods (MONTHLY/QUARTERLY/SEMESTER/ANNUAL) |
| `GET` / `PUT` / `DELETE` | `/performance/periods/:id` | Get / Update / Delete performance period |
| `POST` / `GET` | `/performance/ratings` | Create / List rating scales (contoh: A=90-100) |
| `GET` / `PUT` / `DELETE` | `/performance/ratings/:id` | Get / Update / Delete rating |
| `POST` / `GET` | `/performance/indicator-formulas` | Create / List KPI score formulas (MANUAL/HIGHER_BETTER/LOWER_BETTER/RANGE) |
| `GET` / `PUT` / `DELETE` | `/performance/indicator-formulas/:id` | Get / Update / Delete formula |
| `GET` | `/performance/logs` | List performance audit logs |
| `GET` | `/performance/logs/:id` | Get audit log by ID |

**KPI — Balanced Scorecard (`/performance/kpi/*`)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` / `GET` | `/performance/kpi/perspectives` | Create / List BSC perspectives |
| `GET` / `PUT` / `DELETE` | `/performance/kpi/perspectives/:id` | Get / Update / Delete perspective |
| `POST` / `GET` | `/performance/kpi/templates` | Create / List KPI templates |
| `GET` / `PUT` / `DELETE` | `/performance/kpi/templates/:id` | Get / Update / Delete template |
| `POST` / `GET` | `/performance/kpi/indicators` | Create / List KPI indicators |
| `GET` / `PUT` / `DELETE` | `/performance/kpi/indicators/:id` | Get / Update / Delete indicator |
| `POST` | `/performance/kpi/evaluations` | Create KPI evaluation |
| `GET` | `/performance/kpi/evaluations` | List evaluations (filter employee/organization/period/status) |
| `GET` / `PUT` / `DELETE` | `/performance/kpi/evaluations/:id` | Get / Update / Delete evaluation |
| `POST` | `/performance/kpi/evaluations/snapshot` | Create evaluation + snapshot KPI dari template (nilai target terkunci) |
| `GET` | `/performance/kpi/evaluations/:id/full` | Get evaluation lengkap (details, targets, comments, attachments) |
| `PUT` | `/performance/kpi/evaluations/:id/status` | Update status workflow (DRAFT → PLAN_SUBMITTED → PLAN_APPROVED → ACTUAL_SUBMITTED → ACTUAL_APPROVED → COMPLETED) |
| `POST` | `/performance/kpi/evaluations/:id/submit` | Submit evaluation |
| `POST` | `/performance/kpi/evaluations/:id/approve` | Approve evaluation |
| `POST` | `/performance/kpi/evaluations/:id/reject` | Reject evaluation |
| `POST` | `/performance/kpi/evaluations/:id/complete` | Complete evaluation |
| `POST` | `/performance/kpi/evaluations/:id/recalculate` | Recalculate score |
| `PUT` | `/performance/kpi/evaluations/:id/actuals` | Bulk input actual values |
| `POST` / `GET` | `/performance/kpi/targets` | Create / List KPI targets |
| `POST` | `/performance/kpi/progress` | Catat progres realisasi KPI |
| `POST` | `/performance/kpi/comments` | Tambah komentar/review |
| `POST` | `/performance/kpi/attachments` | Lampirkan file bukti |
| `GET` | `/performance/kpi/dashboard/employee/:employeeId` | Dashboard KPI employee |
| `GET` | `/performance/kpi/dashboard/manager/:managerId` | Dashboard KPI manager |
| `GET` | `/performance/kpi/dashboard/hr` | Dashboard KPI HR |

**OKR — Objectives & Key Results (`/performance/okr/*`)**
| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` / `GET` | `/performance/okr/templates` | Create / List OKR templates |
| `GET` / `PUT` / `DELETE` | `/performance/okr/templates/:id` | Get / Update / Delete template |
| `POST` | `/performance/okr/templates/:id/duplicate` | Duplicate template (beserta objectives & key results) |
| `POST` | `/performance/okr/objectives` | Create objective |
| `GET` / `PUT` / `DELETE` | `/performance/okr/objectives/:id` | Get / Update / Delete objective |
| `POST` | `/performance/okr/key-results` | Create key result (target, unit, formula, bobot) |
| `GET` / `PUT` / `DELETE` | `/performance/okr/key-results/:id` | Get / Update / Delete key result |
| `POST` / `GET` | `/performance/okr/evaluations` | Create (snapshot dari template) / List OKR evaluations |
| `GET` / `PUT` / `DELETE` | `/performance/okr/evaluations/:id` | Get / Update / Delete evaluation |
| `POST` | `/performance/okr/evaluations/:id/submit` | Submit (→ SUBMITTED) |
| `POST` | `/performance/okr/evaluations/:id/approve` | Approve (→ APPROVED) |
| `POST` | `/performance/okr/evaluations/:id/reject` | Reject (→ REJECTED) |
| `POST` | `/performance/okr/evaluations/:id/complete` | Complete (→ COMPLETED) |
| `POST` | `/performance/okr/evaluations/:id/recalculate` | Recalculate score |
| `PUT` | `/performance/okr/evaluations/:id/actuals` | Bulk input actual values |
| `POST` | `/performance/okr/progress` | Catat progres/check-in |
| `POST` | `/performance/okr/comments` | Tambah komentar (mendukung reply) |
| `POST` | `/performance/okr/attachments` | Lampirkan file bukti |
| `GET` | `/performance/okr/dashboard/hr` | Dashboard OKR HR |

> 🔍 Daftar lengkap 117 endpoint (17 master data + 58 KPI + 42 OKR) dengan skema request/response: lihat [`docs/openapi-report.md`](docs/openapi-report.md).

**Settings — Master Data / Reference CRUD**

Semua endpoint settings berada di `/api/v1/tenant/settings/`. Masing-masing entity memiliki 5 CRUD endpoint standar (`GET`, `POST`, `GET/:id`, `PUT/:id`, `DELETE/:id`):

| Entity | Endpoint | Fields | PK Type | Frontend View |
|--------|----------|--------|:-------:|:-------------:|
| Zones | `/settings/zones` | Code, Name, Region, IsActive, SortOrder | `UUID` | `ZonesView.vue` |
| Provinces | `/settings/provinces` | Code, Name | `CHAR(2)` — BPS code | `ProvincesView.vue` |
| Regencies | `/settings/regencies` | Code, Name, ProvinceID (parent) | `CHAR(4)` — BPS code | `RegenciesView.vue` |
| Districts | `/settings/districts` | Code, Name, RegencyID (parent) | `CHAR(6)` — BPS code | `DistrictsView.vue` |
| Villages | `/settings/villages` | Code, Name, DistrictID (parent) | `CHAR(10)` — BPS code | `VillagesView.vue` |
| Educations | `/settings/educations` | Code, Name, SortOrder | `UUID` | `EducationsView.vue` |
| Education Majors | `/settings/education-majors` | Code, Name, SortOrder | `UUID` | `EducationMajorsView.vue` |
| Religions | `/settings/religions` | Code, Name, SortOrder | `UUID` | `ReligionsView.vue` |
| Marital Statuses | `/settings/marital-statuses` | Code, Name, SortOrder | `UUID` | `MaritalStatusesView.vue` |
| Relationship Types | `/settings/relationship-types` | Code, Name, SortOrder | `UUID` | `RelationshipTypesView.vue` |
| Banks | `/settings/banks` | Code, Name, SortOrder | `UUID` | `BanksView.vue` |
| Employment Statuses | `/settings/employment-statuses` | Code, Name, SortOrder | `UUID` | `EmploymentStatusesView.vue` |
| Nationalities | `/settings/nationalities` | Code, Name, SortOrder | `UUID` | `NationalitiesView.vue` |
| Job Families | `/settings/job-families` | Code, Name, SortOrder | `UUID` | `JobFamiliesView.vue` |
| Salary Grades | `/settings/salary-grades` | Code, Name, Grade, MinSalary, MaxSalary, SortOrder | `UUID` | `SalaryGradesView.vue` |
| Ters | `/settings/ters` | Group, BrutoMin, BrutoMax, Rate, SortOrder | `UUID` | `TersView.vue` |
| Ptkps | `/settings/ptkps` | Status, DependentCount, Amount, SortOrder | `UUID` | `PtkpsView.vue` |
| Insurances | `/settings/insurances` | Code, Name, SortOrder | `UUID` | `InsurancesView.vue` |
| Company Holidays | `/settings/company-holidays` | HolidayDate (unique), Name, Description, IsActive | `UUID` | `CompanyHolidaysView.vue` |

> **Query parameters:** Semua endpoint GET mendukung `?page=&per_page=&search=&sort_by=&sort_order=` — standar server-side pagination & sorting. Untuk zona, ada tambahan `?active_only=true`.

**Tenant: Packages — Browse & Subscribe**

| Method | Endpoint | Deskripsi |
|---|---|---|
| `GET` | `/packages` | List published packages (browse available packages) |
| `POST` | `/packages/:id/subscribe` | Subscribe to package — create license + auto-activate modules |
| `POST` | `/packages/:id/unsubscribe` | Unsubscribe from package — deactivate modules + suspend license |

**Subscribe response (201):**
```json
{
  "success": true,
  "data": {
    "license_id": "uuid",
    "license_key": "xxx",
    "plan_type": "pro",
    "package_id": "uuid",
    "package_name": "HR Pro Package",
    "activated_modules": ["Employee Management", "Payroll"]
  },
  "message": "Successfully subscribed to package"
}
```

### Response Format

```json
// Success
{
    "success": true,
    "data": { ... },
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 100
    }
}

// Error (general)
{
    "success": false,
    "error": {
        "code": "NOT_FOUND",
        "message": "Company not found"
    }
}

// Error (validation — field-level)
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Validation failed",
        "fields": {
            "name": ["This field is required"],
            "email": ["Must be a valid email address", "Already exists"],
            "admin_name": ["This field is required"]
        }
    }
}
```

**Frontend field-level validation display** (platform-admin):
Gunakan `getValidationErrors()` dari `services/responseHandler.js` untuk mengekstrak field errors:
```js
import { getValidationErrors } from '@/services/responseHandler'

const errors = ref({})
try {
  await api.post('/endpoint', data)
} catch (e) {
  errors.value = getValidationErrors(e)
  // errors.value = { name: "This field is required", email: "Must be valid, Already exists" }
  // Array values otomatis di-implode dengan ', '
}
```

Tampilkan error per-field di `FormRow`:
```vue
<FormRow label="Email" :errors="errors?.email">
  <TextInput v-model="form.email" :class="{ 'p-invalid': errors?.email }" />
</FormRow>
```

### Bilingual Support (Bahasa Indonesia & English)

API mendukung dua bahasa untuk pesan response.

- **English** (default): Jika header `Accept-Language` tidak disertakan atau bernilai selain `id`
- **Bahasa Indonesia**: Kirim header `Accept-Language: id`

Header ini memengaruhi semua pesan response, termasuk:
- Pesan sukses (created, updated, deleted)
- Pesan error (not found, internal error)
- Pesan validasi (field-level errors)

**Contoh:**
```bash
# Response in English (default)
curl http://localhost:8080/api/v1/platform/companies/not-found
# {"success":false,"error":{"code":"NOT_FOUND","message":"Not found"}}

# Response in Bahasa Indonesia
curl -H "Accept-Language: id" http://localhost:8080/api/v1/platform/companies/not-found
# {"success":false,"error":{"code":"NOT_FOUND","message":"Tidak ditemukan"}}
```

### Custom Validators for Indonesian Data Formats

Endpoint tenant (prefixed `/api/v1/tenant`) mendukung validasi built-in untuk format data Indonesia:

| Tag | Format | Contoh | Keterangan |
|-----|--------|--------|------------|
| `nik` | 16 digit angka | `3273010101900001` | NIK KTP |
| `npwp` | 15-16 digit angka | `0123456789012345` | NPWP |
| `npwp_format` | XX.XXX.XXX.X-XXX.XXX | `01.234.567.8-901.234` | NPWP terformat |
| `kk` | 16 digit angka | `1234567890123456` | Kartu Keluarga |
| `phone_id` | +628/08xx (7-11 digit) | `08123456789` | No. Telepon Indonesia |
| `postal_code` | 5 digit angka | `12345` | Kode Pos |
| `date_id` | YYYY-MM-DD | `2026-12-31` | Tanggal |
| `passport` | 1 huruf + 8 angka | `A12345678` | Passport |
| `sim` | 12 digit angka | `123456789012` | SIM |
| `no_rekening` | 8-20 digit angka | `1234567890` | No. Rekening Bank |

**Contoh penggunaan di DTO:**
```go
type CreateEmployeeRequest struct {
    Name        string  `json:"name" binding:"required"`
    NIK         string  `json:"nik" binding:"omitempty,nik"`
    NPWP        *string `json:"npwp" binding:"omitempty,npwp"`
    Phone       *string `json:"phone" binding:"omitempty,phone_id"`
    Passport    *string `json:"passport" binding:"omitempty,passport"`
    NoRekening  *string `json:"no_rekening" binding:"omitempty,no_rekening"`
}
```

Error response untuk validasi field:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validasi gagal",
    "fields": {
      "nik": ["Format NIK tidak valid, harus 16 digit angka"],
      "npwp": ["Format NPWP tidak valid, harus 15-16 digit angka"]
    }
  }
}
```

---

## 🖥️ Frontend — Platform Admin (Phase 1 ✅)

Platform Admin frontend dibangun dengan **Vue 3 + PrimeVue 4 + Tailwind CSS 4** dan mencakup **9 halaman** (Platform Admin) + **19+ views** (Tenant) = **28+ total views** dengan **32+ komponen reusable** yang terintegrasi penuh dengan backend API.

### Fitur Umum

- 🌐 **Bilingual EN/ID** — Language switcher di HeaderBar, semua label/title/toast mendukung 2 bahasa
- ✅ **Field-Level Validation** — Error per-field dari API ditampilkan langsung di bawah input
- 🔍 **Filter & Search** — Setiap halaman memiliki filter chips, search input, dan sortable columns
- ♻️ **Real-time Polling** — Dashboard & Monitoring dengan auto-refresh 30s
- 📦 **Shared Components** — FormRow, TextInput, SelectLabel, ToggleSwitch, DatePicker, TextArea, dll reusable
- 📅 **Bilingual Date Formatting** — `formatDate(value, locale)` → "30 July 2026" (EN) / "30 Juli 2026" (ID)
- 🗓️ **PrimeVue DatePicker Locale** — Calendar popup menampilkan nama bulan/hari sesuai bahasa aktif + tombol "Today"/"Hari Ini"
- 🌳 **Organization Tree drag & drop** — Pindah parent organisasi via drag & drop di Tree view (`Organizations.vue`): PrimeVue Tree `draggable-nodes`/`droppable-nodes`/`validate-drop` + `@node-drop` (`dragNode`/`dropNode`/`dropPosition`/`accept`), guard anti-siklus (self/descendant), lalu `PUT /organizations/:id { parent_id }`; parent juga bisa diganti via dropdown di dialog edit
- 🔗 **Axios Interceptors** — Auto JWT + Accept-Language header + 401 auto-refresh

### Halaman

| Modul | Halaman | Fitur Unggulan |
|-------|---------|----------------|
| B.1 | **Login** (`Login.vue`) | JWT auth, bilingual form, error per field |
| B.2 | **Dashboard** (`Dashboard.vue`) | 6 KPI cards, bar chart (company/month), system health, polling 30s |
| B.3 | **Companies** (`Companies.vue`) | Filter by status/package, search, license inline, **provisioning progress** (DB column + tooltip) |
| B.4 | **Users** (`Users.vue`) | Filter by role, search, **bulk actions** (Change Role + Delete), super_admin protection |
| B.5 | **Modules** (`Modules.vue`) | Filter `?module_type=`, filter chips, depends_on tooltip, auto-slug |
| B.6 | **Licenses** (`Licenses.vue`) | Package integration, filter by package, copy key, expiry warnings |
| B.7 | **Monitoring** (`Monitoring.vue`) | Auto-refresh toggle, **pool utilization chart** (Line), **alert thresholds** |
| B.8 | **Packages** (`Packages.vue`) | CRUD + publish/unpublish, **dependency validation**, module selector with Select All |
| B.9 | **RBAC** (`Rbac.vue`) | Roles CRUD, **permission matrix** grouped by module/resource, system protection |

### Bilingual Architecture

```
[Language Switcher HeaderBar]
    → Pinia Store (language.js) — persist to localStorage
    → Axios Interceptor — header Accept-Language: en|id
    → Backend Localize middleware — response otomatis sesuai bahasa
    → Response Handler — parse message + validation errors bilingual
```

- **200+ locale keys** di `en.json` & `id.json` untuk static UI text
- **9+ bilingual composables** (`useI18n`, `useNotification`, `useSlugify`)
- **2 bilingual utils** (`formatDate.js` — date formatting EN/ID; `primevueLocale.js` — PrimeVue DatePicker locale)
- **Semua view** menggunakan `t('key')` untuk label, header, filter, tooltip
- **DateInput bilingual** — Calendar popup menampilkan nama bulan Bahasa Indonesia (Januari, Februari, ...) saat bahasa ID aktif

### Struktur Frontend

```
frontend/platform-admin/src/
├── App.vue                    # Root — Toast + ConfirmDialog
├── main.js                    # PrimeVue, Tailwind, Router, Pinia
├── assets/css/                # Global CSS (slug-animation.css)
├── components/                # 11 reusable form/components (CompanyActions, Skeleton, dll)
├── composables/               # useI18n, useSkeletonPage, useSlugify
├── layouts/                   # AppLayout (bilingual title), HeaderBar (theme/lang/Live/user menu), Sidebar (8 menu collapsible)
├── locales/                   # en.json, id.json (391 keys)
├── router/index.js            # 11 routes (login + 10 halaman) + auth guard
├── services/                  # api.js (axios) + responseHandler.js
├── stores/                    # auth.js + language.js + theme.js (Pinia)
└── views/                     # 9 halaman
```

---

## 💻 Development

## 🔧 CLI Installer Commands

CLI Installer (`go run ./cmd/installer`) menyediakan perintah untuk provisioning tenant, migrasi, dan seeding data.

### Prerequisites

```bash
cd backend
go mod download
```

### Perintah Dasar

| Command | Deskripsi |
|---------|-----------|
| `provision` | Provision tenant baru (create DB + migrate + seed data) |
| `migrate` | Jalankan migrasi pending untuk tenant |
| `seed-data` | Seed master reference data ke tenant DB |
| `seed-modules` | Register modul platform & tenant ke database |
| `encrypt-passwords` | Encrypt legacy credentials |

### Seeders — Master Reference Data

Berikut data master yang di-seed ke setiap tenant baru (dan bisa di-re-seed untuk tenant existing):

#### Reference Data (Simple CRUD)

| Tabel | Records | Keterangan |
|-------|:-------:|------------|
| `religions` | 7 | Islam, Kristen, Katolik, Hindu, Budha, Konghucu, Lainnya |
| `educations` | 12 | TK s.d. S3 |
| `marital_statuses` | 4 | Belum Kawin, Kawin, Cerai Hidup, Cerai Mati |
| `relationship_types` | 10 | Suami, Istri, Anak, Orang Tua, dll |
| `employment_statuses` | 6 | Tetap, Kontrak, Magang, Harian Lepas, Outsourcing, Pensiunan |
| `banks` | 25 | Bank Mandiri, BCA, BNI, BRI, BSI, dll |
| `nationalities` | 241 | Afghanistan, Albania, dll (ISO 3166-1 alpha-2) |
| `job_families` | 15 | Finance, HR, IT, Marketing, Sales, dll |

#### Data Wilayah Indonesia (Bulk)

| Tabel | Records | PK Type | Sumber Data |
|-------|:-------:|:-------:|:-----------:|
| `provinces` | **38** | `CHAR(2)` — BPS code (e.g. `11`) | Kemendagri |
| `regencies` | **188** | `CHAR(4)` — BPS code (e.g. `1101`) | Kemendagri |
| `districts` | **7.230** | `CHAR(6)` — BPS code (e.g. `110101`) | Embedded SQL (`tenantseed/seeddata/002_seed_districts.sql`) |
| `villages` | **83.441** | `CHAR(10)` — BPS code (e.g. `1101012001`) | Embedded SQL (`tenantseed/seeddata/003_seed_villages.sql`) |

> ℹ️ **ID = Code:** Tabel wilayah menggunakan kode BPS (Badan Pusat Statistik) sebagai primary key, bukan UUID. Hal ini karena data wilayah merupakan referensi tetap yang di-seed dari dataset Kemendagri — foreign keys (`province_id`, `regency_id`, `district_id`) menggunakan kode BPS parent yang sesuai.

#### Payroll & Pajak Data

| Tabel | Records | Keterangan |
|-------|:-------:|------------|
| `salary_grades` | **12** | Golongan IA s.d. IVC (Rp1,5jt–200jt) |
| `ptkps` | **12** | Status PTKP TK/0–K/I/3 (Rp54jt–126jt) |
| `pph21_ptkp_rates` | **12** | PTKP rates 2024 |
| `pph21_tax_brackets` | **5** | Bracket PPh21: 5%–35% (UU HPP) |
| `ters` | **124** | TER Group A=44, B=40, C=40 (rate 0%–34%) |
| `bpjs_settings` | **1** | Default BPJS 2024 (health max Rp12jt) |
| `bpjs_rate_components` | **13** | HEALTH(2), JHT(2), JP(2), JKK(5), JKM(1), JKP(1) |

**Total:** ~90.948 records (termasuk 83.441 villages + 7.230 districts)

### Contoh Penggunaan

#### 1. Provision Tenant Baru (Auto-Seed)

```bash
cd backend

# Provision + auto-migrate + auto-seed (satu perintah)
go run ./cmd/installer provision --company=<uuid> --config=./config/config.yaml
```

Ini akan:
1. Create database `hris_tenant_<uuid[:8]>`
2. Save TenantConnection
3. Jalankan 145+ migration
4. ✅ **Auto-seed semua data master** (~90.948 records)

#### 2. Seed Data ke Tenant Existing

```bash
# Seed semua data master (idempotent — record yang sudah ada di-skip)
go run ./cmd/installer seed-data --company=<uuid> --config=./config/config.yaml
```

#### 3. Jalankan Migrasi Tenant

```bash
go run ./cmd/installer migrate --company=<uuid> --config=./config/config.yaml
```

#### 4. Register Modules

```bash
# Daftarkan module platform & tenant ke database
make seed-modules
# atau:
go run ./cmd/installer seed-modules --config=./config/config.yaml
```

#### 5. Encrypt Legacy Passwords

```bash
export HRIS_ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

go run ./cmd/installer encrypt-passwords --config=./config/config.yaml
```

### Makefile Shortcuts

| Perintah | Fungsi |
|---------|--------|
| `make seed-modules` | Register modules ke platform database |
| `make seed` | Jalankan platform seeding (via server) |
| `make migrate` | Jalankan platform migrations |
| `make build-installer` | Build CLI installer binary (ke `bin/installer`) |
| `make docs` | Verifikasi OpenAPI + regenerate `docs/openapi-report.md` |
| `make db-docs` | Regenerate `docs/database-schema.md` dari migrasi SQL |
| `make check-db-docs` | Verifikasi `docs/database-schema.md` sinkron dengan migrasi SQL postgres **&** mysql (tanpa menimpa) |
| `make arch-report` | Regenerate `docs/go-module-architecture-report.md` dari analisis statis kode Go |
| `make check-arch-report` | Verifikasi `docs/go-module-architecture-report.md` sinkron dengan kode Go (tanpa menimpa) |

> **Catatan:** `make seed` (server) berbeda dari `go run ./cmd/installer seed-data` (installer). Gunakan `make seed` untuk platform-level dan installer `seed-data --company=<uuid>` untuk tenant-level master data.

### Idempotency

Semua seeder bersifat **idempotent** — jika dijalankan ulang, record yang sudah ada akan di-skip (berdasarkan UUID deterministik).

```
Contoh output:
  religions                  inserted=0     skipped=7     5ms
  districts                  inserted=0     skipped=7293  2.3s
  salary_grades              inserted=0     skipped=12    8ms
  bpjs_settings              inserted=0     skipped=1     5ms
  ters                       inserted=0     skipped=124   15ms
```

## 💻 Development

### Makefile Commands

Jalankan dari direktori `backend/`:

```bash
make build            # Build binary
make run              # Run server
make run-hot          # Run with hot reload (air)
make test             # Run all tests
make race-test        # Run all tests with race detector in Docker (golang:1.24)
make db-race-up       # Start MySQL + PostgreSQL containers for the migrator integration test
make db-race-down     # Stop & remove the race-test DB containers
make lint             # Run linter
make vet              # Run go vet
make coverage         # Run tests with coverage report
make docs             # Verify OpenAPI + regenerate docs/openapi-report.md
make db-docs          # Regenerate docs/database-schema.md from SQL migrations
make check-db-docs    # Verify database-schema.md in sync with postgres & mysql migrations (no overwrite)
make arch-report      # Regenerate docs/go-module-architecture-report.md from Go source
make check-arch-report  # Verify go-module-architecture-report.md in sync with Go source (no overwrite)
make docker           # Build Docker image
make tidy             # Tidy dependencies
make clean            # Clean build artifacts
make help             # Show all commands
```

### Menambah Module Baru

Setiap modul harus mengimplementasikan `module.Module` interface:

```go
type Module interface {
    Info() ModuleInfo
    RegisterRoutes(router *gin.RouterGroup)
    Migrate(db *gorm.DB) error
    Seed(db *gorm.DB) error
    Permissions() []string
}
```

Template struktur module:

```text
internal/modules/{nama_module}/
├── handler.go
├── service.go
├── repository.go
├── model.go
├── dto.go
├── routes.go
└── module.go
```

---

## 🐳 Docker

### Development Environment

```bash
# Start semua services
make docker-compose-up

# Services:
# - PostgreSQL :5432
# - Redis      :6379
# - API       :8080
# - Asynqmon  :8081

# Stop
make docker-compose-down
```

### Build Docker Image

```bash
make docker
docker run -p 8080:8080 hris-platform:latest
```

---

## 🧪 Testing

### Test Coverage

| Package | Unit Tests | Integration Tests | Benchmarks | Total |
|---------|:----------:|:-----------------:|:----------:|:----:|
| `internal/pkg/cache/` | 42 (cache + PubSub) | 8 (full lifecycle, 2-instance Pub/Sub, TTL, concurrent, data types) | 31 (set/get/invalidate/PubSub/concurrent/data size) | **81** |
| `internal/modules/competency/` | 54 (service 25 + repository 14 + handler 15) | — | — | **54** |
| `internal/modules/payroll/` | **34** (repository 13 + service 21) | — | — | **34** |
| `internal/pkg/authz/` | **80+** (enforcer 26 + repository 22 + service 20 + handler 12) | — | — | **80+** |
| `internal/modules/approval/` | **67** (repository 25 + service 25 + handler 14) | — | — | **67** |
| `internal/modules/employeemovement/` | **58** (service 22 + repository 22 + handler 14) | — | — | **58** |
| `internal/modules/attendance/` | **83** (repository 37 + service 25 + handler 21) | — | — | **83** |
| `internal/modules/leave/` | **38** (repository 14 + service 12 + handler 12) | — | — | **38** |
| `internal/modules/performance/` | **55** (repository 14 + service 24 + handler 17) | — | — | **55** |
| `internal/modules/recruitment/` | **66** (repository 27 + service 23 + handler 16) | — | — | **66** |
| `internal/modules/reimbursement/` | **48** (repository 18 + service 20 + handler 16) | — | — | **48** |
| `internal/modules/training/` | **31** (repository 14 + service 9 + handler 8) | — | — | **31** |
| `internal/modules/workforceintelligence/` | **108** (repository 31 + service 41 + handler 36) | — | — | **108** |
| `internal/modules/careerintelligence/` | **65** (repository 23 + service 20 + handler 22) | — | — | **65** |

### Run Tests

```bash
# Semua test
make test

# Cache package specific
go test ./internal/pkg/cache/ -v

# Benchmarks
go test ./internal/pkg/cache/ -bench=. -benchmem

# Integration tests
go test ./internal/pkg/cache/ -run TestCacheIntegration -v

# Dengan coverage
make coverage     # hasil: coverage.html

# Short test (tanpa integration)
make test-short
```

### Race Detector (`make race-test`)

Menjalankan full suite dengan **race detector** untuk memastikan tidak ada *data race*:

```bash
# Cara standar — otomatis menyalakan MySQL + PostgreSQL via Docker, lalu
# menjalankan seluruh suite dengan -race
make race-test
```

**Kenapa Docker?** Race detector Go memerlukan cgo + kompiler C (`CGO_ENABLED=1`). Di Windows tanpa gcc terpasang, `make test` / `make test-verbose` (yang juga memakai `-race`) tidak bisa berjalan — image resmi `golang:1.24` sudah menyertakan gcc, sehingga `-race` berfungsi di dalam container.

**Alur kerja `make race-test`:**

1. **`db-race-up`** (dependency otomatis) — menjalankan container MySQL + PostgreSQL (`--network host`) dengan kredensial yang sesuai DSN *hardcoded* di `internal/pkg/migrator/migrator_integration_test.go`:
   - MySQL: `root@` tanpa password, port `3306`
   - PostgreSQL: `postgres:password`, port `5432`
   - Database: `hris_migrate_test`
   - Idempotent — container yang sudah ada dipakai ulang; jika rusak, dibuat ulang otomatis.
2. **Test race** — `go test -race -count=1 -p 4 -timeout 15m ./...` di container `golang:1.24`, dengan volume cache (`go-mod-cache` & `go-build-cache`) agar run berikutnya cepat.

**Perintah terkait:**

```bash
make race-test                              # Full suite + race detector (DB otomatis)
make race-test PKG=./internal/modules/leave/...   # Hanya paket tertentu
make db-race-up                             # Nyalakan DB saja (idempotent)
make db-race-down                           # Hentikan & hapus container DB
```

> ⚠️ **Catatan port:** `db-race-up` mengikat port `3306`/`5432` di dalam Docker VM. Jika service `docker compose --profile mysql/postgres up` sedang berjalan, hentikan dulu karena akan terjadi konflik port.

### Cache Test Details

Semua cache tests menggunakan **miniredis** — pure Go Redis server untuk testing, tanpa perlu Redis asli berjalan.

**Unit tests (cache_test.go + pubsub_test.go):**
- Cache: Set/Get/Miss/SetJSON/Invalidate/InvalidatePrefix/LocalCache/TTL/Ping/Close/Concurrent/ErrorPath
- PubSub: Publish/SingleKey/MultipleKeys/HandleMessage/CrossInstance/Concurrent/ErrorPath

**Integration tests (cache_integration_test.go):**
- Full lifecycle: New → Set → Get → SetJSON → Get (JSON) → Invalidate → Ping → Close
- Two-instance Pub/Sub invalidation (cross-instance propagation)
- TTL expiry (Redis + local cache wall clock)
- Concurrent workflow (20 goroutines mixed Set/Get/SetJSON/Invalidate)
- Cache miss scenarios (non-existent, invalidated, empty value, idempotent)
- Large payload (50KB byte-per-byte validation)
- Complex data types (nested JSON struct)
- Cross-instance prefix invalidation

**Benchmarks (bench_test.go):**
- Set/Get latency (hot cache ~21ns, cold cache ~28µs, miss ~27µs)
- JSON serialization overhead
- Invalidate/InvalidatePrefix throughput
- Concurrent operations (RunParallel)
- Data size comparison (64B → 1MB)
- PubSub publish/handle throughput

---

## 🧩 Module SDK

Setiap modul di HRIS Platform harus mengikuti kontrak Module SDK:

1. **ModuleInfo** — identitas modul (nama, slug, versi, dependensi)
2. **RegisterRoutes** — daftarkan semua endpoint HTTP
3. **Migrate** — migrasi database modul
4. **Seed** — data awal / seeder
5. **Permissions** — daftar permission yang dibutuhkan

```yaml
# manifest.yaml (contoh)
name: "Organization Management"
slug: "organization"
version: "1.0.0"
depends_on: []
permissions:
  - "organization.view"
  - "organization.create"
menus:
  - name: "Organization"
    icon: "building"
    route: "/admin/organizations"
```

---

## 🔄 Database Migration

Platform menggunakan **Go SQL Migration Runner** (`internal/pkg/migrator/`) dengan dua mode:

### SQL Migrator (file `.sql` embedded)

File SQL migration di-embed ke binary Go via `//go:embed` — tidak perlu file eksternal saat runtime.

**File Convention:**
```text
001_create_companies.sql          ← Up migration (DDL)
001_create_companies.down.sql     ← Down migration (rollback) — optional
```

**Startup Flow:**
```text
[1] SQL Migrator (platform/*.sql)  → CREATE TABLE, indexes, FKs
[2] GORM AutoMigrate               → Sync Go struct columns
[3] SQL Seeders (seeders/*.sql)    → INSERT initial data
[4] Module Seeders                 → Business logic seed data
```

### CLI Usage

```bash
# Normal startup (run all pending migrations + seeders)
go run ./cmd/server --config ./config/config.yaml

# Rollback ALL applied migrations (down files required)
go run ./cmd/server --config ./config/config.yaml --migrate-down

# Rollback to specific version (exclusive)
# Contoh: applied [001,002,003,004,005,006]
# --migrate-to 004 → rollback 006, 005 (004 tetap applied)
go run ./cmd/server --config ./config/config.yaml --migrate-to 004
```

> ⚠️ **--migrate-down** dan **--migrate-to** bersifat mutually exclusive.
> Setelah selesai, server akan exit (tidak start).

### Tracking Table

Setiap migration yang sukses dicatat di tabel `schema_migrations`:

```sql
CREATE TABLE schema_migrations (
    version     VARCHAR(14) PRIMARY KEY,   -- '001', '002'
    name        VARCHAR(255) NOT NULL,      -- 'create_companies'
    applied_at  TIMESTAMP NOT NULL,         -- waktu eksekusi
    checksum    VARCHAR(64) NOT NULL,
    file_path   VARCHAR(500) NOT NULL
);
```

### Migration Files

- `internal/pkg/migrator/migrations/platform/` — 8 platform DDL files (001–006 + RBAC + 012, termasuk roles, permissions, role_permissions)
- `internal/pkg/migrator/migrations/seeders/` — 2 file (1 seeder `001_seed_super_admin.sql` + down)
- `internal/pkg/tenantseed/seeddata/` — 2 embedded bulk SQL files (tenant: districts & villages), di-embed via `//go:embed seeddata` → CWD-independent
- `internal/pkg/migrator/migrations/tenant/` — Tenant migrations (86 file per dialect — mysql & postgres)

### Transaction Safety

- **Up**: SQL content + insert `schema_migrations` record dalam 1 database transaction
- **Down**: Down SQL + delete `schema_migrations` record dalam 1 database transaction
- Jika salah satu gagal → entire transaction rollback

---

## 🏗️ Tenant Provisioning (End-to-End Verified ✅)

Platform memiliki **Tenant Provisioning Engine** yang secara otomatis membuat database dan menjalankan migrasi saat company baru dibuat.

### Provisioning Flow

```
POST /api/v1/platform/companies
    ↓
1. Generate slug dari company name
2. Cek duplikasi slug (tambahkan UUID suffix jika perlu)
3. Save company ke platform DB (status: active)
4. provisionTenant():
   ├── a. Generate db_name = hris_{slug}
   ├── b. Connect sebagai superuser (root@localhost)
   ├── c. Buat database tenant (CREATE DATABASE IF NOT EXISTS)
   ├── d. Simpan TenantConnection ke platform DB (ID = companyID)
   ├── e. Connect ke tenant DB via GORM    └── f. Jalankan 89 tenant SQL migrations (199 tables)
5. Jika provisioning berhasil → company status: active
6. Jika provisioning gagal → company status: suspended (data tetap tersimpan)
```

### Tenant Migration Files (89 files → 199 tables)

| File | Isi |
|------|-----|
| `001_master_data.sql` | Master tables (religions, educations, provinces, dll) |
| `002_organization.sql` | Organization structure, zones, job families, positions |
| `003_employee.sql` | Employees, employments, families, educations, documents |
| `004_attendance.sql` | Attendance settings, shifts, events, overtime |
| `005_leave.sql` | Leave types, requests, accrual policies, company holidays |
| `006_payroll_structure.sql` | Salary components, grades, payroll periods |
| `007_payroll_run.sql` | Payroll run, run items, run employees, payslips |
| `008_competency.sql` | Competencies, values, events, scores + FK dari migration 002 |
| `009_job_management.sql` | Job titling, values, objectives, responsibilities |
| `010_approval.sql` | Approval flows, steps, instances, tasks |
| `011_settings.sql` | Settings & RBAC: features, roles, permissions, holidays, document templates |
| `012_employee_movement.sql` | Employee movements & contracts |
| `013_reimbursement.sql` | Reimbursement types, requests, items |
| `014_performance.sql` | Performance periods, perspectives, templates, indicators, evaluations, targets |
| `015_recruitment.sql` | Job requisitions, candidates, applications, interviews, onboarding tasks |
| `016_training.sql` | Training categories, courses, sessions, participants, materials, evaluations, certificates |
| `017_workforce_intelligence.sql` | Workforce planning headcounts, forecasts, KPIs, analytics cache, scenarios, risk indicators, health scores |
| `018_career_intelligence.sql` | Career talent maps (9-box grid), career interests, career paths, succession plans |
| `019_fix_org_unique.sql` | Composite unique (organization_summary_id, full_code) |
| `020_add_passport_to_employees.sql` | Add passport column to employees (idempotent) |
| `021_insurances.sql` | Insurance master (BPJS Kesehatan, BPJS Ketenagakerjaan, dll) |
| `022_users.sql` | Users table (Level 2 Tenant RBAC) |
| `023_user_accounts.sql` | Employee user accounts (login access for employees) |
| `024_education_majors.sql` | Education majors master (jurusan pendidikan) |

> **Catatan:** Total 200 tabel termasuk `schema_migrations` (199 tabel tenant + schema_migrations, auto-created oleh migrator engine).

### Daftar Lengkap 199 Tabel Tenant (+ schema_migrations = 200)

**Approval (5):**
`approval_actions`, `approval_flow_steps`, `approval_flows`, `approval_instances`, `approval_tasks`

**Attendance (10):**
`attendance_company_settings`, `attendance_company_shifts`, `attendance_device_captures`,
`attendance_employee_shifts`, `attendance_events`, `attendance_exempt_positions`,
`attendance_face_captures`, `attendance_locations`, `attendance_overtime_requests`, `attendance_sessions`

**BPJS (3):**
`bpjs_rate_components`, `bpjs_settings`, `insurances`

**Competency (7):**
`competence_values`, `competencies`, `competency_event_targets`, `competency_events`,
`competency_score_details`, `competency_scores`, `competency_values`

**Employee (17):**
`emergency_contacts`, `employee_addresses`, `employee_bank_accounts`, `employee_bank_profiles`, `employee_bpjs_profiles`,
`employee_contracts`, `employee_documents`, `employee_educations`, `employee_experiences`, `employee_families`,
`employee_insurances`, `employee_leave_balances`, `employee_movements`, `employee_payroll_profiles`, `employee_tax_profiles`,
`employees`, `employments`

**Job Management (20):**
`job_families`, `job_family_competencies`, `job_management_assets`, `job_management_competency_groups`,
`job_management_education_experiences`, `job_management_financials`, `job_management_hr_authorities`,
`job_management_identifications`, `job_management_objectives`, `job_management_operational_authorities`,
`job_management_potency_competencies`, `job_management_relationships`, `job_management_responsibilities`,
`job_management_scores`, `job_management_subordinate_controls`, `job_management_title_subs`,
`job_management_titles`, `job_management_values`, `job_management_working_activities`,
`job_management_working_risks`

**Leave (6):**
`company_holidays`, `leave_accrual_policies`, `leave_reasons`, `leave_request_details`, `leave_requests`, `leave_types`

**Master Data (13):**
`banks`, `districts`, `document_templates`, `educations`, `employment_statuses`,
`gradings`, `marital_statuses`, `nationalities`, `provinces`, `regencies`, `relationship_types`, `religions`, `villages`

**Organization (5):**
`organization_levels`, `organization_summaries`, `organizations`, `positions`, `zones`

**Payroll (17):**
`payroll_payslips`, `payroll_periods`, `payroll_profile_change_logs`, `payroll_run_employees`,
`payroll_run_items`, `payroll_runs`, `pph21_calculation_logs`, `pph21_ptkp_rates`,
`pph21_settings`, `pph21_tax_brackets`, `ptkps`, `salary_change_logs`,
`salary_components`, `salary_employee_adjustments`, `salary_employee_components`, `salary_grade_components`,
`salary_grades`

**Performance (7):**
`performance_periods`, `performance_perspectives`, `performance_templates`, `performance_indicators`,
`performance_evaluations`, `performance_evaluation_details`, `performance_targets`

**Reimbursement (3):**
`reimbursement_items`, `reimbursement_requests`, `reimbursement_types`

**Recruitment (7):**
`job_requisitions`, `candidates`, `job_applications`, `interviews`,
`onboarding_task_templates`, `employee_onboardings`, `onboarding_task_items`

**Workforce Intelligence (7):**
`workforce_planning_headcounts`, `workforce_forecasts`, `workforce_kpis`,
`workforce_analytics_cache`, `workforce_scenarios`, `workforce_risk_indicators`,
`workforce_health_scores`

**Career Intelligence (4):**
`career_talent_maps`, `career_interests`, `career_paths`, `career_succession_plans`

**Training & Development (7):**
`training_categories`, `training_courses`, `training_sessions`, `training_participants`,
`training_materials`, `training_evaluations`, `training_certificates`

**Settings & Permissions (8):**
`feature_permission`, `features`, `model_has_permissions`, `model_has_roles`,
`permissions`, `role_has_permissions`, `roles`, `users`

**System (1):**
`schema_migrations`

**Tax (1):**
`ters`

### Provisioning Test Results ✅

**Test Date:** 22 Juli 2026  
**Environment:** Development — Laragon MySQL 8.0  
**Driver:** MySQL (`multiStatements=true`, `parseTime=True`)

| Item | Status | Detail |
|------|--------|--------|
| Company status | ✅ **active** | API mengembalikan `status: "active"` |
| Tenant database | ✅ Created | `hris_final-provision-test` |
| Tenant connection | ✅ Saved | Record di `tenant_connections` tersimpan |
| Migrations | ✅ **89 files** | 001 → 090 sukses semua (nomor 044 tidak ada) |
| Total tables | ✅ **199 tables** | Setiap migrasi menciptakan tabel sesuai DDL |
| Server log | ✅ Clean | "Tenant provisioning completed successfully" |

#### API Test Response

```json
{
  "success": true,
  "data": {
    "id": "0b77f721-cee5-46f0-bc75-45219fc2316d",
    "name": "Final Provision Test",
    "slug": "final-provision-test",
    "email": "final@test.com",
    "phone": "021777777",
    "status": "active"
  }
}
```

### Issues Resolved During Development

| Issue | Root Cause | Fix |
|-------|-----------|-----|
| `TenantConnection` ID field | Struct Go tanpa `ID` → INSERT gagal (PK CHAR(36) tanpa default) | Tambah field ID — reuse companyID sebagai PK |
| `ssl_mode` column mismatch | GORM tag `sslmode` (tanpa underscore) tapi DDL `ssl_mode` | GORM tag: `sslmode` → `ssl_mode` |
| Multi-statement SQL gagal | MySQL driver default blokir multi-statement | Tambah `multiStatements=true` ke DSN |
| Access denied tenant DB | Kredensial pakai nama DB sebagai user (user tidak ada) | Development: gunakan `root`/`""` |
| FK dependency cross-file | `002_organization.sql` punya FK ke tabel `competencies` (migration 008) | Pindah FK ke migration 008 via ALTER TABLE |

---

## ✅ Module Completion Status

### Modul Inti (Completed ✅)

| Modul | Status | Detail |
|-------|:------:|--------|
| **Platform & Tenant Management** | ✅ **Completed** | Provisioning DB multi-tenant, isolasi database, switching context tenant, lifecycle management (Suspend/Activate/Terminate) |
| **Organization Management** | ✅ **Completed** | Multi-Company Architecture, Dynamic Department Hierarchy (Adjacency List), Location & Geofencing Zones, Organization Summary |
| **Employee Management** | ✅ **Completed** | Data personal, kontak, alamat, keluarga, pendidikan, dokumen, riwayat kerja, rekening/pajak, 8 sub-modules |
| **Job Management** | ✅ **Completed** | 18 GORM entities: Titles, Subs, Values, Objectives, Identifications, Responsibilities, Education Experiences, HR/Operational Authorities, Working Activities/Risks, Relationships, Subordinate Controls, Assets, Financials, Potency Competencies, Scores, Competency Groups |
| **Competency Management** | ✅ **Completed** | 7 GORM entities: Competencies, CompetenceValues (legacy), CompetencyValues (structured), CompetencyEvents, CompetencyEventTargets, CompetencyScores, CompetencyScoreDetails |
| **RBAC Management (Database-Backed)** | ✅ **Completed** | 4 default roles with hierarchy, 13 seeded resources (70+ permissions), CRUD API (10 endpoints), enforcer auto-reload, **80+ unit tests** |
| **Employee Movement & Career Management** | ✅ **Completed** | 2 entities (EmployeeMovement, EmployeeContract) with 8 movement types, contract extension chain, 3-step approval flow (draft→approved→executed), **58 unit tests**, 29 OpenAPI endpoints |
| **Approval Engine** | ✅ **Completed (24 Juli 2026)** | **5 GORM entities**: ApprovalFlow, ApprovalFlowStep, ApprovalInstance, ApprovalAction, ApprovalTask. Multi-step workflow engine dengan flexible approval flows. **67 unit tests** (25 repository + 25 service + 14 handler + 3 helpers). **18 endpoints** — 3 approval modes (ANY_ONE, ALL, N_OF_M), 3 approver types (USER, ROLE, SUPERVISOR), instance tracking & task management.
  - **Flows**: CRUD approval flows, soft-delete, module-based filtering
  - **Steps**: Ordered step management per flow (auto-increment step_order)
  - **Instances**: Create approval instances for any module/document, with status workflow PENDING→APPROVED/REJECTED/CANCELLED
  - **Actions**: Submit approve/reject actions, automatic next-step advancement or full approval
  - **Tasks**: Task management per approver with assignee resolution |
| **Payroll & Compensation Engine** | ✅ **Completed (24 Juli 2026)** | **21 GORM entities**: SalaryComponent, PayrollPeriod, PayrollRun, BpjsSetting, BpjsRateComponent, Pph21Setting, Pph21PtkpRate, Pph21TaxBracket, EmployeePayrollProfile, EmployeeBankProfile, EmployeeBpjsProfile, EmployeeTaxProfile, PayrollRunEmployee, PayrollRunItem, PayrollPayslip, Pph21CalculationLog, dan lainnya. Full CRUD (Create/Read/Update/Delete) untuk seluruh entity. **34 unit tests** (13 repository + 21 service).
  - **BPJS Indonesia**: BPJS Settings & Rate Components (Kesehatan & Ketenagakerjaan)
  - **PPh21**: Settings, Tax Brackets, PTKP Rates, Calculation Logs
  - **Payroll Run**: Periods, Runs (status workflow DRAFT→CALCULATED→REVIEWED→APPROVED→LOCKED), Run Employees, Items, Payslips
  - **OpenAPI**: **83 endpoints** terintegrasi — periods, runs, payslips, payments + export, reports (summary/detail/bpjs/tax/bank), salary grade & employee components, BPJS/PPh21 settings & brackets |
| **Performance Management** | ✅ **Completed (31 Juli 2026)** | **15 GORM entities** (7 KPI + 8 OKR): PerformancePeriod, PerformancePerspective, PerformanceTemplate, PerformanceIndicator, PerformanceEvaluation, PerformanceEvaluationDetail, PerformanceTarget, OKRTemplate, OKRObjective, OKRKeyResult, OKREvaluation, OKREvaluationDetail, OKRProgress, OKRComment, OKRAttachment. KPI (Balanced Scorecard) & OKR framework dengan period-based evaluations, status workflow (DRAFT → PLAN_SUBMITTED → PLAN_APPROVED → ACTUAL_SUBMITTED → ACTUAL_APPROVED → COMPLETED), dan dashboard per role (employee/manager/HR). **55 unit tests** (14 repository + 24 service + 17 handler). **147 endpoints** — 17 master data (periods, ratings, indicator-formulas, logs) + 80 KPI (`/performance/kpi/*`) + 50 OKR (`/performance/okr/*`). Full OpenAPI documentation terintegrasi. |
| **Recruitment & Onboarding ATS** | ✅ **Completed (31 Juli 2026)** | **7 GORM entities**: Requisition, Candidate, Application, Interview, InterviewResult, OnboardingTask, OnboardingChecklist. End-to-end hiring pipeline: job req → candidate sourcing → application → interview → offer → onboarding. **66 unit tests** (27 repository + 23 service + 16 handler). **101 endpoints** across 7 resource groups (+ requirements, competencies, offers workflow, consents, screenings, assessments, scorecards). Full ATS workflow with status tracking (OPEN→IN_PROGRESS→FILLED→CANCELLED). |
| **Time & Attendance** | ✅ **Completed (26 Juli 2026)** | **10 GORM entities**: Company Settings, Company Shifts, Employee Shifts, Locations (Geofence), Device Captures, Face Captures, Events (Check-in/out), Sessions (Daily Work), Overtime Requests, Exempt Positions. Full CRUD untuk 8 sub-entities. **83 unit tests** (37 repo + 25 service + 21 handler). **92 endpoints** (termasuk Business Travel, overtime & stats). OpenAPI docs enhanced — versi 1.6.3 |
| **Leave & Time Off** | ✅ **Completed (26 Juli 2026)** | **6 GORM entities**: LeaveTypes, LeaveAccrualPolicies, LeaveReasons, LeaveRequests, LeaveRequestDetails, EmployeeLeaveBalances. Full CRUD untuk 5 sub-entities. **38 unit tests** (14 repo + 12 service + 12 handler). **26 endpoints**. |
| **Reimbursement & Claim** | ✅ **Completed (30 Juli 2026)** | **3 GORM entities**: ReimbursementType, ReimbursementRequest, ReimbursementItem. Multi-step status flow (DRAFT->SUBMITTED->APPROVED->PAID). **48 unit tests** (18 repo + 20 service + 16 handler). **15 OpenAPI endpoints**. Approval integration via Approval Engine. |
| **Training & Development Management** | ✅ **Completed (1 Agustus 2026)** | **7 GORM entities**: TrainingCategory, TrainingCourse, TrainingSession, TrainingParticipant, TrainingMaterial, TrainingEvaluation, TrainingCertificate. Full CRUD untuk 7 resource groups dengan validasi (quota checks, FK existence, session status). **31 unit tests** (14 repository + 9 service + 8 handler). **123 endpoints** — P0 core (categories, courses, sessions, participants, materials, evaluations, certificates), P1 planning & governance (plans, needs, requests, mandatories), P2 advanced (evaluation forms, effectiveness, certifications, history, reports). Migration files (016_training) untuk MySQL & Postgres. |
| **Workforce Intelligence & Strategic Planning** | ✅ **Completed (25 Juli 2026)** | **7 GORM entities**: WorkforcePlanningHeadcount, WorkforceForecast, WorkforceKPI, WorkforceAnalyticsCache, WorkforceScenario, WorkforceRiskIndicator, WorkforceHealthScore. **108 unit tests** (31 repo + 41 service + 36 handler). **70 endpoints** — Planning, Analytics, Risk, Executive, Scenarios, People Analytics, Cost, Capacity, Health, KPI, quality-of-hire & projections. Migration files (017_workforce_intelligence) untuk MySQL & Postgres. |
| **Career Intelligence & Talent Management** | ✅ **Completed (1 September 2026)** | **4 GORM entities**: CareerTalentMap, CareerInterest, CareerPath, CareerSuccessionPlan. **65 unit tests** (23 repo + 20 service + 22 handler). **23 endpoints** — 9-box talent mapping, career interests, career paths + gap analysis, succession plans + gaps. OpenAPI docs (23 endpoints, 24th tag). Migration files (018_career_intelligence) untuk MySQL & Postgres. |

### Modul Operasional & Siklus Karier (Planned 🗓️)

| Modul | Prioritas | Scope |
|-------|:---------:|-------|
| **Organization History, Versioning & Cloning** | 🟢 High | Change Capture, Full Structure Cloning DRAFT, Version Audit Trail |

---

## ✅ Completed Work

### 📄 Documentation

| # | Item | File |
|---|------|------|
| ✅ | Analisis blueprint v3 vs existing Laravel app | `docs/analisis-blueprint-vs-existing.md` |
| ✅ | Platform architecture design (modular monolith, multi-tenant) | `docs/platform-architecture-design.md` |
| ✅ | Project completion dashboard (25 modules — 19 tenant + 6 platform, 1547 tests, 199 tables) | `docs/project-completion-dashboard.md` |
| ✅ | OpenAPI comprehensive report (1166 endpoints, 701 paths, 799 schemas, 38 tags) | `docs/openapi-report.md` |
| ✅ | Go module architecture report (199 entities, 961 service methods, 1547 tests) | `docs/go-module-architecture-report.md` |
| ✅ | Environment variables template | `backend/.env.example` |
| ✅ | Build & development Makefile | `backend/Makefile` |
| ✅ | README utama proyek | `README.md` |

### ⚙️ Shared Kernel (Backend Infrastructure)

| # | Package | Files | Deskripsi |
|---|---------|-------|-----------|
| ✅ | `internal/pkg/config/` | `config.go` | Viper configuration loader (YAML + .env + env vars) |
| ✅ | `internal/pkg/database/` | `manager.go` | Multi-tenant DB connection manager with caching |
| ✅ | `internal/pkg/driver/` | `driver.go` | Shared DB driver type (PostgreSQL / MySQL) |
| ✅ | `internal/pkg/auth/` | `jwt.go` | JWT token generation & validation (access + refresh) |
| ✅ | `internal/pkg/middleware/` | `auth.go`, `cors.go`, `logger.go`, `recovery.go`, `tenant.go` | Gin middleware stack (auth, CORS, logging, recovery, tenant resolver) |
| ✅ | `internal/pkg/router/` | `router.go` | Router setup & module registration |
| ✅ | `internal/pkg/logger/` | `logger.go` | Zap structured logger |
| ✅ | `internal/pkg/module/` | `module.go` | Module SDK interface & registration |
| ✅ | `internal/pkg/cache/` | `cache.go`, `pubsub.go` | Distributed two-tier cache (local sync.Map + Redis) + Pub/Sub invalidation |
| ✅ | `internal/pkg/cache/` (tests) | `cache_test.go`, `pubsub_test.go`, `bench_test.go`, `cache_integration_test.go` | **42 unit tests + 8 integration tests + 31 benchmarks** — test coverage with miniredis |
| ✅ | `internal/pkg/authz/` | `rbac.go`, `model.go`, `repository.go`, `service.go`, `handler.go`, `routes.go` | Database-backed RBAC: 4 roles with hierarchy, 70+ seeded permissions, enforcer with auto-reload |
| ✅ | `internal/pkg/authz/` (tests) | `helpers_test.go`, `rbac_test.go`, `repository_test.go`, `service_test.go`, `handler_test.go` | **80+ unit tests** — enforcer DB loading (26) + repository (22) + service (20) + handler (12) with SQLite in-memory |
| ✅ | `internal/modules/competency/` | `model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`, `module.go` | Competency Management — 7 entities full CRUD + auto-migrate + Module SDK |
| ✅ | `internal/modules/competency/` (tests) | `helpers_test.go`, `service_test.go`, `repository_test.go`, `handler_test.go` | **54 unit tests** — service (25) + repository (14) + handler (15) tests with SQLite in-memory |
| ✅ | `internal/modules/employeemovement/` | `model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `routes.go`, `module.go` | Employee Movement & Career Management — 2 entities (movements + contracts), 8 movement types, 3-step approval flow, contract chain |
| ✅ | `internal/modules/employeemovement/` (tests) | `helpers_test.go`, `service_test.go`, `repository_test.go`, `handler_test.go` | **58 unit tests** — service (22) + repository (22) + handler (14) tests with SQLite in-memory |

### 🏢 Platform Module — Company

| # | Item | Files |
|---|------|-------|
| ✅ | Company model (UUID, status lifecycle) | `internal/platform/company/model.go` |
| ✅ | Company CRUD (Create, Read, Update, Delete) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Company lifecycle (Suspend, Activate, Backup, Restore) | `handler.go`, `service.go` |
| ✅ | Company request/response DTOs | `dto.go` |
| ✅ | Company routes registration | `routes.go` |
| ✅ | TenantConnection model (multi-tenant DB config) | `tenant_connection.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |

### 👤 Platform Module — Users & Auth

| # | Item | Files |
|---|------|-------|
| ✅ | PlatformUser model (UUID, bcrypt password, roles) | `internal/platform/user/model.go` |
| ✅ | JWT Authentication (Login, Refresh Token) | `service.go`, `handler.go` |
| ✅ | User CRUD (Create, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs | `dto.go` |
| ✅ | Routes registration | `routes.go` |
| ✅ | Auto-seed super admin (development) | `module.go` (via `EnsureSeed`) |

### 🧩 Platform Module — Module Management

| # | Item | Files |
|---|------|-------|
| ✅ | PlatformModule model (registered modules) | `internal/platform/modulemgmt/model.go` |
| ✅ | CompanyModule model (company-module association) | `internal/platform/modulemgmt/model.go` |
| ✅ | Module CRUD (Register, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Module activation/deactivation for company | `handler.go`, `service.go` |
| ✅ | Company-module listing | `repository.go` (via JOIN query) |
| ✅ | Routes registration | `routes.go` |

### 🔑 Platform Module — License Management

| # | Item | Files |
|---|------|-------|
| ✅ | License model (plan types, limits, dates) | `internal/platform/license/model.go` |
| ✅ | License CRUD (Create, Read, Update, List) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Plan-based auto calculation (max employees, modules) | `service.go` |
| ✅ | License key generation (via xid) | `service.go` |
| ✅ | Routes registration | `routes.go` |

### 📊 Platform Module — Monitoring

| # | Item | Files |
|---|------|-------|
| ✅ | Platform health check (DB connection status) | `internal/platform/monitoring/handler.go` |
| ✅ | Tenant connection health listing | `internal/platform/monitoring/handler.go` |
| ✅ | Tenant connection health detail | `internal/platform/monitoring/handler.go` |
| ✅ | Routes registration | `internal/platform/monitoring/routes.go` |

### 🏛️ Tenant Module — Organization

| # | Item | Files |
|---|------|-------|
| ✅ | Organization model (tree hierarchy with parent_id) | `internal/modules/organization/model.go` |
| ✅ | Organization CRUD + Tree view | `handler.go`, `service.go`, `repository.go` |
| ✅ | Organization request/response DTOs | `dto.go` |
| ✅ | Organization routes registration | `routes.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |

### 📦 API Documentation

| # | Item | File |
|---|------|------|
| ✅ | OpenAPI 3.0 JSON specification (**1166 endpoints**, 701 paths, 799 schemas, 38 tags) | `internal/pkg/docs/openapi.json` |
| ✅ | Scalar UI served at `/docs` (interactive documentation) | `internal/pkg/docs/scalar.go` |
| ✅ | OpenAPI spec served at `/openapi.json` | `internal/pkg/docs/scalar.go` |

### 🐳 Infrastructure

| # | Item | File |
|---|------|------|
| ✅ | Multi-stage Dockerfile | `backend/docker/Dockerfile` |
| ✅ | Docker Compose (PostgreSQL, Redis, API, Asynqmon) | `docker/docker-compose.yml` |
| ✅ | PostgreSQL init script | `docker/init-db.sql` |
| ✅ | CLI Installer stub for tenant provisioning | `backend/cmd/installer/main.go` |
| ✅ | API server entry point | `backend/cmd/server/main.go` |
| ✅ | Server config YAML | `backend/config/config.yaml` |
| ✅ | Go module dependencies | `go.mod`, `go.sum` |

### 🗄️ Database Support

| # | Item | Detail |
|---|------|--------|
| ✅ | PostgreSQL driver | `gorm.io/driver/postgres v1.5.9` |
| ✅ | MySQL driver | `gorm.io/driver/mysql v1.6.0` |
| ✅ | Multi-driver DSN generation | PostgreSQL & MySQL format in config |
| ✅ | Per-tenant driver configuration | `driver` field in `tenant_connections` table |
| ✅ | Shared driver type | `internal/pkg/driver/driver.go` |

### 📦 Go Dependencies

```
github.com/gin-gonic/gin v1.10.0          # HTTP framework
github.com/golang-jwt/jwt/v5 v5.2.1       # JWT auth
github.com/google/uuid v1.6.0             # UUID generation
github.com/gosimple/slug v1.15.0          # URL slug
github.com/joho/godotenv v1.5.1           # .env file loader
github.com/spf13/viper v1.19.0            # Configuration
go.uber.org/zap v1.27.0                   # Structured logging
gorm.io/driver/mysql v1.6.0               # MySQL driver
gorm.io/driver/postgres v1.5.9            # PostgreSQL driver
gorm.io/gorm v1.30.0                      # ORM
```

### ✅ Tenant Provisioning

| # | Item | Detail |
|---|------|--------|
| ✅ | Provisioning Engine | Database creation + TenantConnection save |
| ✅ | Tenant SQL Migrations | 178 files per dialect (89 up + 89 down) → 199 tables |
| ✅ | Multi-statement MySQL support | `multiStatements=true` di DSN |
| ✅ | Error handling / graceful failure | Company status = `suspended` jika provisioning gagal |
| ✅ | End-to-end test | Company active ✅, 199 tables ✅, MySQL |

### ✅ Tenant Lifecycle Management

| # | Item | Detail |
|---|------|--------|
| ✅ | `Suspend` endpoint | Deactivate connection + status suspended |
| ✅ | `Activate` endpoint | Reactivate connection + status active |
| ✅ | `Terminate` endpoint | Drop database + remove connection + status terminated |
| ✅ | Soft delete cleanup | `DELETE /:id` now deactivates connection too |
| ✅ | Cached connection cleanup | `CloseTenantConnection` removes cached GORM connections |
| ✅ | End-to-end test | Suspend ✅ → Activate ✅ → Terminate ✅ |

### 🏛️ Tenant Module — Employee ✅

| # | Item | Files |
|---|------|-------|
| ✅ | Employee model (9 GORM models with UUID hooks) | `internal/modules/employee/model.go` |
| ✅ | Employee CRUD (Create, Read, Update, Delete) | `handler.go`, `service.go`, `repository.go` |
| ✅ | 8 sub-modules (Addresses, Emergency Contacts, Families, Educations, Experiences, Documents, Insurances, Employments) | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs with validation (oneof, required, email, max) | `dto.go` |
| ✅ | Paginated List response | `service.go` |
| ✅ | Routes registration (40 endpoints) | `routes.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |
| ✅ | AutoMigrate for 9 models during tenant provisioning | `module.go` |
| ✅ | End-to-end tested: Create, List, Get, Update, Delete + all sub-modules | API verified |

### 🏛️ Tenant Module — Job Management ✅

| # | Item | Files |
|---|------|-------|
| ✅ | Job Management model (18 GORM entities with UUID hooks) | `internal/modules/jobmanagement/model.go` |
| ✅ | Full CRUD for all 18 entities | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs with validation for all 18 entities | `dto.go` |
| ✅ | Paginated List responses for list endpoints | `service.go` |
| ✅ | Routes registration (97 endpoints) | `routes.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |
| ✅ | SQLite-integrated unit tests (74 tests) | `*_test.go` |
| ✅ | OpenAPI 3.0 documentation (23 schemas + 97 endpoints) | `internal/pkg/docs/openapi.json` |
| ✅ | RBAC permission: `jobmanagement.*` for company_admin | `internal/pkg/authz/rbac.go` |

### 🏛️ Tenant Module — Competency Management ✅

| # | Item | Files |
|---|------|-------|
| ✅ | Competency model (7 GORM entities with UUID hooks) | `internal/modules/competency/model.go` |
| ✅ | Full CRUD for all 7 entities: Competencies, CompetenceValues (legacy), CompetencyValues (structured), CompetencyEvents, CompetencyEventTargets, CompetencyScores, CompetencyScoreDetails | `handler.go`, `service.go`, `repository.go` |
| ✅ | Request/Response DTOs with validation (required, oneof, max) | `dto.go` |
| ✅ | Paginated List responses for all list endpoints | `service.go` |
| ✅ | Routes registration (67 endpoints) | `routes.go` |
| ✅ | Module registration (Module SDK compliance) | `module.go` |
| ✅ | Context-driven tenant DB resolver | `module.go` |
| ✅ | AutoMigrate for 7 models during tenant provisioning | `module.go` |
| ✅ | SQLite-integrated unit tests (54 tests: 25 service + 14 repository + 15 handler) — covers all 7 entities with CRUD + edge cases (invalid UUID, not found, pagination, validation errors) | `*_test.go` |

### Build Status

```bash
$ go vet ./...    # ✅ Lulus
$ go build ./...  # ✅ Berhasil
```

---

## 🚀 Deployment (SaaS & On-Premise)

Platform mendukung **dua mode distribusi** dengan binary yang sama (`cmd/server`) — perbedaan hanya pada konfigurasi `license.deployment_mode`:

| Mode | `deployment_mode` | Lisensi Modul | Batas Employee |
|---|---|---|---|
| **Subscription SaaS** (multi-tenant) | `saas` (default) | Per-company via `company_modules` DB + package subscribe | Unlimited |
| **On-Premise** (dedicated) | `on_premise` | File `.lic` RSA (`allowed_modules`) | `max_employees` di-enforce → 403 `QUOTA_EXCEEDED` |

Ringkasan cepat:

```bash
# SaaS (default) — jalankan server
./bin/server --config ./config/config.yaml

# On-Premise — vendor generate lisensi, klien deploy
cd backend && go run ./cmd/licensectl gen-key --out private.pem --pub public.pem
go run ./cmd/licensectl gen-lic --priv private.pem --out license.lic \
  --company-id <uuid> --company "PT X" --expires 2027-12-31 \
  --modules organization,employee,payroll --max-employees 500

export HRIS_LICENSE_DEPLOYMENT_MODE=on_premise
export HRIS_LICENSE_LICENSE_FILE=/etc/hris/license.lic
export HRIS_LICENSE_PUBLIC_KEY_FILE=/etc/hris/public.pem
./bin/server --config ./config/config.yaml
```

> 📖 **Dokumentasi lengkap**: [`docs/deployment-guide.md`](docs/deployment-guide.md) — arsitektur kedua mode, langkah instalasi SaaS & on-premise, perbandingan, production checklist, dan troubleshooting.

---

## 🗺️ Roadmap

### Phase 1: Foundation ✅
- ✅ Platform Architecture Design
- ✅ Core Shared Packages (config, database, auth, middleware, router)
- ✅ Platform Module — Company Management (CRUD + TenantConnection)
- ✅ Tenant Module — Organization (tree CRUD)
- ✅ Docker & CI/CD setup
- ✅ Multi-database driver support (PostgreSQL + MySQL)
- ✅ Environment configuration (.env)
- ✅ Tenant Provisioning Engine (12 migrations → 108 tables, end-to-end verified)
- ✅ RBAC Authorization Engine (role hierarchy, resource-action policy)
- ✅ SQL Migration Runner (Up/Down/DownTo rollback, embedded SQL files)
- ✅ Tenant Lifecycle Management (Suspend/Activate/Terminate + connection cleanup)

### Phase 2: Tenant Core Modules ✅
- ✅ Employee Module (full CRUD + 8 sub-modules: addresses, emergency contacts, families, educations, experiences, documents, insurances, employments)

### Phase 3: Payroll & Complex ✅
- ✅ **Job Management Module** (16+ models, 18 GORM entities, full CRUD + scoring) -- [Selesai 22 Juli 2026]
- ✅ **Competency Management** (7 entities: Competencies, Values, Events, Targets, Scores, Score Details) — [Selesai 23 Juli 2026]
- ✅ **Payroll & Compensation Engine** (21 GORM entities, full CRUD, 34 unit tests) — [Selesai 24 Juli 2026]

### Phase 4: Operations & Career ✅
- ✅ **Employee Movement & Career Management** (Promosi/Demosi, PKWT, Pensiun/PHK) — [Selesai 23 Juli 2026]
- ✅ **Time & Attendance** (presensi, shift, lembur) — [Selesai 26 Juli 2026]
- ✅ **Leave & Time Off** (cuti, izin, sakit, multi-level approval) — [Selesai 26 Juli 2026]
- ✅ **Performance Management** (KPI, OKR, review 360) — [Selesai 31 Juli 2026]
- ✅ **Reimbursement & Claim** — [Selesai 30 Juli 2026]
- ✅ **Recruitment & Onboarding (ATS)** — [Selesai 31 Juli 2026]
- ✅ **Training & Development Management** (Training categories, courses, sessions, participants, materials, evaluations, certificates) — [Selesai 1 Agustus 2026]

### Production Readiness 🎯
- ✅ AES-256-GCM encryption untuk kredensial tenant DB (`internal/pkg/crypto/`, CLI `encrypt-passwords`)
- ✅ Enkripsi & masking data karyawan sensitif (NIK, telepon, paspor, email, rekening bank, kontak darurat) — encrypt-on-write + role-based view masking, toggle per field via halaman Pengaturan Data Sensitif
- ✅ **Connection Pool optimization** (Platform: 10/5/1jam, Tenant: 10/3/30mnt/5mnt idle→close, PoolStats(), PgBouncer support)
- ✅ SQL dialect separation (PostgreSQL vs MySQL migrations) — 84 file per dialect (42 up + 42 down), auto-select via `TenantRootPath(driver)`
- ✅ Redis Pub/Sub untuk distributed cache invalidation (`internal/pkg/cache/` — two-tier + Pub/Sub)
- ⬜ Frontend Implementation (Vue 3 + PrimeVue)

### Phase 5: Polish
- ⬜ E2E Testing (Playwright)
- ⬜ Performance Optimization

---

## 📚 Dokumentasi Lainnya

| Dokumen | Deskripsi |
|---|---|
| [`docs/README.md`](docs/README.md) | **Index pusat** seluruh dokumentasi (navigasi cepat ke semua dokumen) |
| [`docs/api/api-usage-guide.md`](docs/api/api-usage-guide.md) | Panduan penggunaan API — menjalankan server, auth, format response, contoh curl, error codes |
| [`docs/database-schema.md`](docs/database-schema.md) | Struktur database & ERD — Platform DB (11 tabel) + Tenant DB (173 tabel), relasi FK, konvensi kolom |
| [`docs/platform-architecture-design.md`](docs/platform-architecture-design.md) | Architecture design document lengkap (satu-satunya dokumen arsitektur) |
| [`docs/deployment-guide.md`](docs/deployment-guide.md) | Panduan deployment lengkap: Subscription SaaS (multi-tenant) & On-Premise (dedicated `.lic` RSA) |
| [`docs/openapi-report.md`](docs/openapi-report.md) | OpenAPI comprehensive report (report v23, spec 1.6.3 — 1166 endpoints, 701 paths, 799 schemas, 38 tags) |
| [`docs/go-module-architecture-report.md`](docs/go-module-architecture-report.md) | Go module architecture report (137 entities, 622 service methods, 1265 tests) |
| [`docs/project-completion-dashboard.md`](docs/project-completion-dashboard.md) | Project completion dashboard (25 modules, 1547 tests, 199 tables) |
| [`docs/panduan-uiux-hris-enterprise.md`](docs/panduan-uiux-hris-enterprise.md) | Standar UI/UX enterprise (modal-first, high-density, prompt AI, warna badge) |
| [`docs/frontend-development-plan.md`](docs/frontend-development-plan.md) | Roadmap implementasi frontend Platform Admin & Tenant |
| [`docs/job-management-score-analysis.md`](docs/job-management-score-analysis.md) | Analisa perhitungan Job Management Score (dirujuk `calculator.go`) |
| [`docs/analisis-blueprint-vs-existing.md`](docs/analisis-blueprint-vs-existing.md) | Analisis blueprint vs existing Laravel app |
| [`backend/.env.example`](backend/.env.example) | Template environment variables |
| [`backend/Makefile`](backend/Makefile) | Build & development commands |

---

## 🛡️ Security

- **JWT Authentication** — Access + Refresh token
- **Casbin RBAC** — Role-based access control
- **Multi-Tenant Isolation** — Database per tenant
- **Input Validation** — via go-playground/validator
- **SQL Injection Prevention** — GORM parameterized queries
- **Audit Log** — Semua mutasi tercatat

---

## 📄 Lisensi

Proprietary — All rights reserved.

---

*Dokumen ini diperbarui pada: 26 Juli 2026*
