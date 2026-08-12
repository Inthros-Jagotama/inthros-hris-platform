# HRIS Platform — Panduan Penggunaan API (API Usage Guide)

> 🔗 **Index dokumentasi:** [`docs/README.md`](../README.md)  
> **Terkait:** [`openapi-report.md`](../openapi-report.md) · [`deployment-guide.md`](../deployment-guide.md) · [`platform-architecture-design.md`](../platform-architecture-design.md)

Panduan praktis **cara menggunakan API** HRIS Platform: dari menjalankan server, autentikasi, format request/response, sampai contoh pemanggilan end-to-end (curl).

> 📖 Dokumen ini berfokus pada **cara pakai**. Untuk daftar lengkap seluruh 943 endpoint + skema, lihat:
> - [`docs/openapi-report.md`](../openapi-report.md) — laporan komprehensif per modul
> - `backend/internal/pkg/docs/openapi.json` — OpenAPI 3.0 spec (sumber kebenaran)

---

## Daftar Isi

1. [Menjalankan Server (dari Makefile)](#1-menjalankan-server-dari-makefile)
2. [Dokumentasi Interaktif](#2-dokumentasi-interaktif)
3. [Struktur URL API](#3-struktur-url-api)
4. [Autentikasi & Otorisasi](#4-autentikasi--otorisasi)
5. [Format Response](#5-format-response)
6. [Pagination, Filter & Sorting](#6-pagination-filter--sorting)
7. [Bilingual Support](#7-bilingual-support)
8. [Contoh Penggunaan (curl)](#8-contoh-penggunaan-curl)
9. [Error Codes](#9-error-codes)
10. [Maintenance Dokumen API (Makefile)](#10-maintenance-dokumen-api-makefile)

---

## 1. Menjalankan Server (dari Makefile)

Semua perintah dijalankan dari direktori `backend/`.

### 1.1 Prasyarat

| Kebutuhan | Versi | Keterangan |
|---|---|---|
| Go | 1.22+ | Build & run server |
| Database | PostgreSQL 16+ / MySQL 8+ | Driver diatur di config |
| Redis | 7+ | Distributed cache (wajib) |
| Python 3 | - | Untuk tooling docs (`check-openapi`, `docs`) |
| Docker (opsional) | - | Infrastruktur lokal via compose |

### 1.2 Setup Infrastruktur (Docker)

```bash
cd docker

# PostgreSQL + PgBouncer + Redis + Mailpit + Asynqmon + API
docker compose --profile postgres up -d

# ATAU MySQL
docker compose --profile mysql up -d
```

> Mailpit tersedia di `http://localhost:8025` (email dev), SMTP `localhost:1025`.
> Asynqmon di `http://localhost:8081` (queue monitoring).

### 1.3 Build & Run

```bash
cd backend

make build          # → ./bin/server (binary release)
make run            # go run ./cmd/server --config ./config/config.yaml
make run-hot        # hot reload via Air (auto-install jika belum ada)
```

Server aktif di **`http://localhost:8080`**. Verifikasi:

```bash
curl http://localhost:8080/healthz
# {"status":"ok","service":"hris-platform"}
```

### 1.4 Migrasi & Seed

```bash
make migrate        # Jalankan SQL migrations platform (otomatis juga saat startup)
make seed           # Jalankan seeders platform
make seed-modules   # Seed/update tabel modules dari definisi module
```

> ⚠️ **Migrasi tenant tidak berjalan saat startup** — tenant DB dibuat & dimigrasi otomatis saat provisioning company via API (`POST /api/v1/platform/companies`) atau CLI `installer`.

### 1.5 Konfigurasi

Config default: `backend/config/config.yaml`, dapat di-override via env `HRIS_*`:

| Env | Default | Deskripsi |
|---|---|---|
| `HRIS_SERVER_PORT` | `8080` | Port server |
| `HRIS_DATABASE_DRIVER` | `mysql` | `postgres` / `mysql` |
| `HRIS_DATABASE_PLATFORM_HOST` | `localhost` | Host platform DB |
| `HRIS_REDIS_HOST` | `localhost` | Host Redis |
| `HRIS_JWT_SECRET` | `change-this...` | **Wajib diganti di production** |
| `HRIS_ENCRYPTION_KEY` | - | AES-256-GCM 64-char hex (wajib production SaaS) |

Contoh menjalankan dengan driver PostgreSQL:

```bash
export HRIS_DATABASE_DRIVER=postgres
export HRIS_DATABASE_PLATFORM_PORT=5432
export HRIS_DATABASE_PLATFORM_USER=hris
export HRIS_DATABASE_PLATFORM_PASSWORD=hris_secret
make run
```

---

## 2. Dokumentasi Interaktif

Setelah server jalan, dokumentasi API tersedia di:

| URL | Deskripsi |
|---|---|
| `http://localhost:8080/docs` | **Scalar UI** — explore & try endpoint langsung dari browser |
| `http://localhost:8080/openapi.json` | OpenAPI 3.0 spec mentah (JSON) |
| `docs/openapi-report.md` | Laporan markdown statis (943 endpoint, 552 paths, 627 schemas, 33 tag) |

---

## 3. Struktur URL API

Semua API menggunakan prefix versi `/api/v1` dan dikelompokkan per konteks:

| Prefix | Konteks | Auth |
|---|---|---|
| `/api/v1/platform/*` | Platform admin (companies, users, modules, licenses, packages, RBAC, monitoring) | JWT + RBAC |
| `/api/v1/tenant/*` | Data tenant per-company (employees, organization, payroll, dsb.) | JWT + tenant + license + RBAC |
| `/api/v1/public/*` | Endpoint publik tanpa auth | ❌ Publik |
| `/api/v1/tenant/auth/*` | Login/refresh user tenant | ❌ Publik |
| `/healthz`, `/readyz` | Health check | ❌ Publik |
| `/docs`, `/openapi.json` | Dokumentasi | ❌ Publik |
| `/uploads/*` | File statis (upload) | ❌ Publik |

### Endpoint Publik (tanpa token)

| Method | Endpoint | Deskripsi |
|---|---|---|
| `POST` | `/api/v1/platform/login` | Login admin platform |
| `POST` | `/api/v1/platform/refresh` | Refresh access token platform |
| `POST` | `/api/v1/tenant/auth/login` | Login user tenant (employee/company_admin) |
| `POST` | `/api/v1/tenant/auth/refresh` | Refresh access token tenant |
| `GET` | `/api/v1/public/packages` | List package published (browsing) |
| `GET` | `/api/v1/public/companies/resolve` | Resolve company dari hostname/subdomain |
| `POST` | `/api/v1/public/account/setup-password` | Set password via link email (token) |
| `GET` | `/healthz` / `/readyz` | Health & readiness |

---

## 4. Autentikasi & Otorisasi

### 4.1 Alur Autentikasi

```
1. POST /login (email + password)
        │
        ▼
   { access_token, refresh_token, token_type: "Bearer", expires_in, user }
        │
        ▼
2. Setiap request terproteksi sertakan header:
   Authorization: Bearer <access_token>

3. Access token expired (15 menit default) →
   POST /refresh { "refresh_token": "..." } → dapat access_token baru
```

### 4.2 Login Platform Admin

```bash
curl -X POST http://localhost:8080/api/v1/platform/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@hris-platform.com",
    "password": "admin123"
  }'
```

> 🔑 **Default super admin (development):** `superadmin@hris-platform.com` / `admin123` — di-seed otomatis oleh `Service.EnsureSeed` jika belum ada.

**Response 200:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "uuid",
      "name": "Super Admin",
      "email": "superadmin@hris-platform.com",
      "role": "super_admin"
    }
  }
}
```

### 4.3 Login User Tenant

```bash
curl -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "company_slug": "pt-contoh-jaya",
    "email": "admin@contohjaya.com",
    "password": "rahasia123"
  }'
```

> Company diidentifikasi via `company_slug` ATAU `company_id`. Jika keduanya kosong, tenant di-resolve otomatis dari **Host header** (mode SaaS multi-subdomain).

### 4.4 Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/platform/refresh \
  -H "Content-Type: application/json" \
  -d '{ "refresh_token": "<refresh_token>" }'
```

### 4.5 Tenant Resolution (mode SaaS)

Tenant (company) pada request `/api/v1/tenant/*` ditentukan dengan prioritas:

1. **`company_id` dari JWT claims** — identitas menang, tidak bisa di-override
2. **Header `X-Tenant-ID`** — eksplisit (UUID valid), untuk dev/ops
3. **Host header** (`X-Forwarded-Host` → `Host`) → resolve via subdomain/domain company

Response menyertakan header `X-Tenant-ID` dengan company aktif.

### 4.6 RBAC (Role-Based Access Control)

4 role default dengan hierarki: `super_admin` → `company_admin` → `manager` → `employee`.

Format permission: **`resource.action`** (contoh: `company.create`, `employee.view`). Setiap endpoint tenant/platform diperiksa oleh middleware RBAC:

| Method | Endpoint | Diizinkan |
|---|---|---|
| `POST /api/v1/platform/companies` | Buat company | `super_admin` saja |
| `GET /api/v1/platform/users` | List users | `super_admin`, `company_admin` |
| `GET /api/v1/tenant/employees` | List employee | semua role (hak berbeda) |
| `POST /api/v1/tenant/employees` | Create employee | `manager` ke atas |

Akses ditolak → **403** dengan detail role/resource/action.

### 4.7 License Enforcement (modul tenant)

Modul yang tidak aktif untuk company ditolak oleh `PlatformLicenseMiddleware`:

```json
{
  "success": false,
  "error": {
    "code": "MODULE_NOT_LICENSED",
    "message": "..."
  }
}
```

On-premise: kuota employee juga di-enforce → **403 `QUOTA_EXCEEDED`**.

---

## 5. Format Response

### 5.1 Sukses

```json
// GET — list/detail
{
  "success": true,
  "data": { ... }
}

// POST/PUT — dengan pesan lokal
{
  "success": true,
  "data": { ... },
  "message": "Created successfully"
}

// DELETE — tanpa data
{
  "success": true,
  "message": "Deleted successfully"
}
```

### 5.2 Error (umum)

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Company not found"
  }
}
```

### 5.3 Error Validasi (per-field)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "name": ["This field is required"],
      "email": ["Must be a valid email address"],
      "admin_password": ["Minimum length is 6"]
    }
  }
}
```

### 5.4 List & Pagination

```json
{
  "success": true,
  "data": [ ... ],
  "page": 1,
  "per_page": 20,
  "total": 100,
  "total_pages": 5
}
```

---

## 6. Pagination, Filter & Sorting

Endpoint list (GET) mendukung query parameter standar:

| Param | Contoh | Deskripsi |
|---|---|---|
| `page` | `?page=2` | Halaman (default 1) |
| `per_page` | `?per_page=50` | Baris per halaman (max 100) |
| `search` | `?search=budi` | Pencarian teks (field tergantung modul) |
| `sort_by` | `?sort_by=name` | Kolom sorting |
| `sort_order` | `?sort_order=desc` | `asc` / `desc` |
| Filter khusus | `?module_type=tenant`, `?status=ACTIVE`, `?year=2026` | Bervariasi per modul |

Contoh:

```bash
curl "http://localhost:8080/api/v1/tenant/settings/zones?page=1&per_page=10&search=jabodetabek&sort_by=name&sort_order=asc" \
  -H "Authorization: Bearer <token>"
```

---

## 7. Bilingual Support

API mendukung dua bahasa pesan via header **`Accept-Language`**:

| Header | Bahasa |
|---|---|
| (tanpa header) | **English** (default) |
| `Accept-Language: id` | **Bahasa Indonesia** |

Berpengaruh pada: pesan sukses, pesan error, dan pesan validasi per-field.

```bash
curl http://localhost:8080/api/v1/tenant/employees/xxx \
  -H "Authorization: Bearer <token>" \
  -H "Accept-Language: id"
# → { "error": { "code": "NOT_FOUND", "message": "Tidak ditemukan" } }
```

**Validator data Indonesia** yang didukung di endpoint tenant:

| Tag | Format | Contoh |
|---|---|---|
| `nik` | 16 digit | `3273010101900001` |
| `npwp` | 15-16 digit | `0123456789012345` |
| `phone_id` | +628/08xx | `08123456789` |
| `postal_code` | 5 digit | `12345` |
| `date_id` | YYYY-MM-DD | `2026-12-31` |
| `passport` | 1 huruf + 8 digit | `A12345678` |
| `sim` | 12 digit | `123456789012` |
| `no_rekening` | 8-20 digit | `1234567890` |

---

## 8. Contoh Penggunaan (curl)

### 8.1 Alur End-to-End: Platform → Company → Tenant

**Step 1 — Login super admin:**

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/platform/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superadmin@hris-platform.com","password":"admin123"}' \
  | python -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")

echo $TOKEN
```

**Step 2 — Buat company (provision tenant DB otomatis):**

```bash
curl -X POST http://localhost:8080/api/v1/platform/companies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PT Contoh Jaya",
    "npwp": "0123456789012345",
    "address": "Jl. Sudirman No. 1, Jakarta",
    "email": "contact@contohjaya.com",
    "admin_name": "Budi Santoso",
    "admin_email": "admin@contohjaya.com",
    "admin_password": "rahasia123"
  }'
```

**Response 201** → `id` company, `slug` (`pt-contoh-jaya`), `admin_user`, dan `provisioning_info` (via `GET /companies/:id`).

**Step 3 — Login user tenant (company admin):**

```bash
TENANT_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/tenant/auth/login \
  -H "Content-Type: application/json" \
  -d '{"company_slug":"pt-contoh-jaya","email":"admin@contohjaya.com","password":"rahasia123"}' \
  | python -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")
```

**Step 4 — Buat organization:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/organizations \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Direktorat Operasional","parent_id":null}'
```

**Step 5 — List employees (pagination):**

```bash
curl "http://localhost:8080/api/v1/tenant/employees?page=1&per_page=20" \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

**Step 6 — Buat employee:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/employees \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "EMP-001",
    "full_name": "Ani Wijaya",
    "email": "ani@contohjaya.com",
    "nik": "3273010101900001",
    "phone": "08123456789",
    "organization_id": "<org-uuid>"
  }'
```

### 8.2 Endpoint Tenant Lainnya

| Modul | Contoh Endpoint |
|---|---|
| Settings / Master data | `GET /api/v1/tenant/settings/banks`, `.../religions`, `.../zones`, `.../company-holidays` |
| Job Management | `GET/POST /api/v1/tenant/job-management/titles`, `.../values/tree` |
| Training & Development | `GET/POST /api/v1/tenant/trainings/categories`, `.../courses`, `.../plans`, `.../needs`, `.../requests`, `.../providers`, `.../trainers`, `.../sessions`, `.../evaluation-forms`, `.../certifications`, `GET .../reports/participation`, `.../history` → lihat §8.8 |
| Career Intelligence | `GET/POST /api/v1/tenant/career-intelligence/paths`, `.../talent-maps`, `.../interests`, `.../successions`, `GET .../successions/gaps`, `GET .../paths/gap-analysis` → lihat §8.10 |
| Payroll | `GET /api/v1/tenant/payroll/...` |
| Leave | `GET/POST /api/v1/tenant/leave/types`, `.../requests`, `.../balances`, `GET /api/v1/tenant/leave/calendar`, `GET /api/v1/tenant/leave/reports/usage` → lihat §8.7 |
| Performance — Master Data | `GET/POST /api/v1/tenant/performance/periods`, `.../ratings`, `.../indicator-formulas`, `.../logs` |
| Performance — KPI | `GET/POST /api/v1/tenant/performance/kpi/templates`, `.../kpi/indicators`, `.../kpi/evaluations`, `.../kpi/dashboard/hr` |
| Performance — OKR | `GET/POST /api/v1/tenant/performance/okr/templates`, `.../okr/objectives`, `.../okr/key-results`, `.../okr/evaluations`, `.../okr/dashboard/hr` |
| Competency | `GET/POST /api/v1/tenant/competency/competencies`, `.../competence-values`, `.../event-targets`, `.../events`, `.../scores`, `GET .../scores/{scoreId}/details` |
| Workforce Intelligence | `GET /api/v1/tenant/workforce-intelligence/executive/summary`, `GET .../analytics/quality-of-hire` → lihat §8.9 |
| Package (subscribe) | `POST /api/v1/tenant/packages/:id/subscribe` |
| Company self-service | `GET/PUT /api/v1/tenant/companies/me` |
| Module aktif | `GET /api/v1/tenant/company-modules` |
| Organization Summaries | `GET/POST /api/v1/tenant/organization-summaries`, `GET .../stats`, `GET/PUT/DELETE .../{id}` |
| Approval Engine | `GET/POST /api/v1/tenant/approval/flows`, `GET /api/v1/tenant/approval/available-modules`, `POST /api/v1/tenant/approval/instances/:id/actions` |
| Employee Movement | `GET/POST /api/v1/tenant/employee-movements/movements`, `POST .../movements/:id/submit`, `.../:id/execute`, `GET .../employees/:employeeId/career-history`, `.../movement-eligibility`, `.../promotion-eligibility`, `GET .../reports/movements`, `.../reports/contracts`, `.../dashboard` |
| Attendance | `GET/POST /api/v1/tenant/attendance/shifts`, `.../locations`, `.../events`, `GET .../sessions`, `.../calendar`, `.../summary`, `.../reports/sessions`, `POST/GET .../corrections`, `POST .../overtime-requests`, `.../overtime-requests/assign`, `.../overtime-requests/:id/actual`, `GET .../overtime-requests/assignable-employees` |
| Notification | `GET /api/v1/tenant/notifications`, `.../unread-count`, `PATCH .../:id/read`, `POST .../read-all` |
| Reimbursements | `GET/POST /api/v1/tenant/reimbursements/types`, `.../requests`, `PUT .../requests/{id}/status`, `GET/POST .../requests/{requestId}/items` |
| User Accounts | `GET /api/v1/tenant/user-accounts/me`, `POST .../employees/:employeeId`, `GET .../employees/:employeeId`, `POST .../employees/:employeeId/resend` |

**Platform Admin (contoh endpoint):**

| Modul | Contoh Endpoint |
|---|---|
| Licenses | `GET/POST /api/v1/platform/licenses`, `GET/PUT/DELETE .../{id}` |
| Modules | `GET/POST /api/v1/platform/modules`, `GET/PUT .../{id}`, `POST .../{id}/activate`, `.../{id}/deactivate`, `GET .../{id}/companies` |
| Monitoring | `GET /api/v1/platform/monitoring/health`, `.../pool`, `.../seed-status`, `.../tenants`, `GET .../tenants/{id}` |
| RBAC | `GET/POST /api/v1/platform/rbac/roles`, `GET/PUT/DELETE .../roles/{id}`, `POST .../roles/{id}/permissions`, `DELETE .../roles/{id}/permissions/{permissionId}`, `GET/POST .../permissions`, `DELETE .../permissions/{id}` |
| Packages | `GET/POST /api/v1/platform/packages`, `GET/PUT/DELETE .../{id}`, `POST .../{id}/publish`, `.../{id}/unpublish`, `GET .../{id}/validate` |

> 🔍 Daftar lengkap per modul: lihat [`docs/openapi-report.md`](../openapi-report.md).

### 8.3 Contoh Penggunaan OKR (End-to-End)

Alur lengkap pengelolaan **OKR** dari membuat template sampai menyelesaikan evaluasi. Semua endpoint berada di `/api/v1/tenant/performance/okr/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Cek Konteks & Scope → Template OKR (Objectives) → Evaluasi (snapshot) → Proposal Key Results → Persetujuan KR → Input Actual → Workflow Assessment → Dashboard HR
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), serta UUID `organization_id`, `period_id` (performance period), dan `employee_id` yang sudah ada.

**Step 1 — Cek konteks OKR & template aktif untuk user saat ini (self-assessment):**

Sebelum membuat evaluasi, employee dapat mengecek konteksnya: apakah sudah memiliki posisi (employment) aktif, berada di organisasi mana, dan template OKR **aktif** apa saja yang tersedia untuk organisasi tersebut:

```bash
curl "http://localhost:8080/api/v1/tenant/performance/okr/my-context" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": {
#     "has_position": true,
#     "employee_id": "<employee-uuid>",
#     "organization_id": "<org-uuid>",
#     "organization_name": "PT Contoh Sejahtera",
#     "templates": [
#       { "id": "<template-uuid>", "name": "OKR Sales Team — Q3 2026", "status": 1, "period_code": "2026-Q3" }
#     ]
#   } }
```

> `has_position: false` berarti user belum memiliki posisi aktif — evaluasi OKR belum bisa dimulai. Self-assessment hanya memungkinkan setelah ada template **aktif** (`status = 1`) untuk organisasi posisi terakhir user. Mirip dengan `my-context` di modul KPI.

**Step 2 — Cek scope pembuatan objective (cascading top-down):**

Employee yang menyiapkan OKR untuk **bawahannya** dapat mengecek apakah dirinya berhak membuat objective (cascading) dan organisasi bawahan mana yang memenuhi syarat — objective hanya bisa dibuat untuk **anak organisasi langsung** (melewati organisasi kosong), dan hanya setelah menerima objective sendiri (kecuali puncak hierarki yang men-seed cascade):

```bash
curl "http://localhost:8080/api/v1/tenant/performance/okr/templates/objective-scope" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": {
#     "organization_id": "<org-uuid>",
#     "organization_name": "Direktorat Operasional",
#     "eligible": true,
#     "ineligible_reason_key": "",
#     "subordinate_organizations": [
#       { "id": "<child-org-uuid>", "name": "Divisi Penjualan" },
#       { "id": "<child-org-uuid>", "name": "Cabang Surabaya" }
#     ]
#   } }
```

> `eligible: false` berarti user belum memenuhi syarat — `ineligible_reason_key` memuat alasan (mis. `okr.objective_scope_ineligible_no_position` saat tidak punya posisi aktif). Organisasi di `subordinate_organizations` adalah **anak efektif** (turunan langsung melewati organisasi kosong).

**Step 3 — Buat OKR template:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/templates \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "name": "OKR Sales Team — Q3 2026",
    "description": "Target kuartal tim penjualan",
    "effective_date": "2026-07-01",
    "expired_date": "2026-09-30"
  }'
```

**Response 201:**
```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": "<template-uuid>",
    "organization_id": "<org-uuid>",
    "name": "OKR Sales Team — Q3 2026",
    "status": 0,
    "created_at": "2026-08-06T10:00:00Z"
  }
}
```

**Step 4 — Buat objective di dalam template:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/objectives \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "<template-uuid>",
    "title": "Meningkatkan pendapatan penjualan",
    "description": "Fokus ekspansi pasar & penjualan produk unggulan",
    "weight": 100
  }'
```

**Step 5 — (Opsional) Buat key result master data di template:**

Template-level key results **tidak lagi disalin** ke evaluasi — Key Results diusulkan oleh employee saat self-assessment (lihat Step 7). Endpoint master data ini tetap tersedia sebagai standar acuan/referensi:

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/key-results \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "objective_id": "<objective-uuid>",
    "title": "Meraih target penjualan Rp 2 Miliar",
    "target_type": "CURRENCY",
    "target_value": 2000000000,
    "unit": "IDR",
    "formula_type": "HIGHER_BETTER",
    "weight": 50,
    "is_required": true
  }'

# Verifikasi isi template (objectives yang tersedia untuk usulan KR)
curl "http://localhost:8080/api/v1/tenant/performance/okr/templates/<template-uuid>/objectives" \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

**Step 6 — Buat evaluasi OKR untuk employee:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "<employee-uuid>",
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "template_id": "<template-uuid>"
  }'
```

> Evaluasi dibuat dengan status `DRAFT` dan **hanya mereferensikan template** — objectives diterima via `template_id` (frontend membaca `GET /okr/templates/:id/objectives` untuk tahu objective mana yang harus diisi). **Key Results tidak disalin** dari template: employee mengusulkan Key Results-nya sendiri per objective saat fase DRAFT (lihat Step 7).

**Response 201:**
```json
{
  "success": true,
  "data": {
    "id": "<evaluation-uuid>",
    "employee_id": "<employee-uuid>",
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "template_id": "<template-uuid>",
    "status": "DRAFT",
    "details": []
  }
}
```

**Step 7 — Usulkan Key Results per objective (saat DRAFT):**

Karyawan mengusulkan Key Results-nya sendiri di bawah objective hasil snapshot, lengkap dengan target yang diajukan. Total bobot Key Result dalam satu objective maksimal 100%:

```bash
# Usulkan Key Result baru di bawah objective
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/key-results \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "evaluation_id": "<evaluation-uuid>",
    "objective_id": "<objective-uuid>",
    "objective_title": "Meningkatkan pendapatan penjualan",
    "objective_weight": 100,
    "title": "Skor CSAT ≥ 85",
    "target_type": "PERCENTAGE",
    "target_value": 85,
    "unit": "%",
    "formula_type": "HIGHER_BETTER",
    "weight": 100
  }'
# → { "success": true, "data": { "id": "<detail-uuid>", "key_result_title": "Skor CSAT ≥ 85", "target_value": 85, "target_type": "PERCENTAGE", ... } }

# Ubah target Key Result yang diajukan (sebelum di-submit)
curl -X PUT http://localhost:8080/api/v1/tenant/performance/okr/evaluation-key-results/<detail-uuid>/target \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "title": "Skor CSAT ≥ 88", "target_value": 88, "weight": 100 }'

# Hapus Key Result yang diajukan (masih DRAFT)
curl -X DELETE http://localhost:8080/api/v1/tenant/performance/okr/evaluation-key-results/<detail-uuid> \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

> Ulangi untuk setiap objective yang diterima. Setelah proposal selesai, cek `GET /api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/details`.

**Step 8 — Submit & persetujuan Key Results (checkpoint pertama):**

```bash
# Employee mengajukan proposal Key Results
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/submit-key-results \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → status: DRAFT → KR_SUBMITTED

# Atasan menyetujui proposal — evaluasi siap diisi actual
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/approve-key-results \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → status: KR_SUBMITTED → KR_APPROVED ("OKR Active")

# (Opsional) Menolak proposal — kembali ke DRAFT untuk revisi
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/reject-key-results \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "notes": "Target perlu direvisi agar lebih realistis" }'
```

> Checkpoint ini dirutekan melalui **approval engine** — module slug sub-checkpoint OKR dengan **fallback ke module induk `performance`** (satu flow cukup untuk semua checkpoint di bawahnya). Jika flow dikonfigurasi tapi tidak bisa di-resolve, submission **gagal** (hard-fail) — lihat `GET /api/v1/tenant/approval/active-flow?module=performance` di [8.5](#85-contoh-penggunaan-approval-engine--employee-movement). Tanpa approval engine, approve/reject manual tetap tersedia.

**Step 9 — Input nilai aktual (actual) secara massal:**

Ambil UUID detail dari `GET /api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>` lalu kirim nilai realisasi:

```bash
curl -X PUT http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/actuals \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "details": [
      { "id": "<detail-uuid>", "actual_value": 1850000000 }
    ]
  }'
```

**Step 10 — Hitung ulang skor (bisa diulang kapan saja):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/recalculate \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → achievement & final_score dihitung ulang dari actual vs target + formula tiap key result
```

**Step 11 — Workflow assessment (checkpoint kedua — persetujuan akhir):**

Setelah Key Results disetujui (KR_APPROVED), employee mengisi actual (Step 9) lalu mengajukan penilaian akhir. Persetujuan akhir **langsung menyelesaikan** evaluasi — tidak ada aksi "Complete" terpisah:

```bash
# Employee mengajukan penilaian (hanya bisa setelah KR_APPROVED)
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/submit \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → status: KR_APPROVED → SUBMITTED

# Atasan menyetujui penilaian → langsung COMPLETED (hasil terkunci)
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/approve \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → status: SUBMITTED → COMPLETED

# (Opsional) Menolak penilaian — kembali ke KR_APPROVED (KR tetap disetujui,
# employee cukup merevisi actual self-assessment)
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/reject \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "notes": "Actual perlu dilengkapi" }'
```

> Endpoint `POST /evaluations/:id/complete` masih tersedia sebagai jalur manual/legacy, namun alur normal memakai `approve` sebagai langkah penyelesaian.

**Step 12 — Lihat dashboard OKR HR (ringkasan seluruh evaluasi):**

```bash
curl "http://localhost:8080/api/v1/tenant/performance/okr/dashboard/hr?period_id=<period-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → total_evaluations, sebaran status (draft/submitted/approved/completed), skor & achievement rata-rata
```

**Catatan tambahan OKR:**

| Endpoint | Deskripsi |
|---|---|
| `GET /api/v1/tenant/performance/okr/templates/objective-scope` | Cek scope pembuatan objective cascading (`eligible` + daftar org bawahan) |
| `POST /api/v1/tenant/performance/okr/evaluations/:id/key-results` | Usulkan Key Results karyawan per objective (saat DRAFT) |
| `PUT /api/v1/tenant/performance/okr/evaluation-key-results/:id/target` | Ubah target Key Results yang diajukan |
| `DELETE /api/v1/tenant/performance/okr/evaluation-key-results/:id` | Hapus Key Results yang diajukan |
| `POST /api/v1/tenant/performance/okr/evaluations/:id/submit-key-results` | Ajukan proposal Key Results (→ `KR_SUBMITTED`) |
| `POST /api/v1/tenant/performance/okr/evaluations/:id/approve-key-results` | Setujui proposal Key Results (→ `KR_APPROVED`) |
| `POST /api/v1/tenant/performance/okr/evaluations/:id/reject-key-results` | Tolak proposal Key Results (→ `DRAFT`) |
| `POST /api/v1/tenant/performance/okr/templates/:id/duplicate` | Duplikat template (beserta objective & key results) sebagai template baru |
| `GET /api/v1/tenant/performance/okr/objectives/:id/key-results` | List key results dalam sebuah objective |
| `POST /api/v1/tenant/performance/okr/progress` | Catat progres/check-in per detail (riwayat tanggal) |
| `POST /api/v1/tenant/performance/okr/comments` | Komentar/review evaluasi (dukung reply via `parent_id`) |
| `POST /api/v1/tenant/performance/okr/attachments` | Lampirkan file bukti ke evaluation detail |

> 💡 Status evaluasi OKR **dua fase**: `DRAFT` → `KR_SUBMITTED` → `KR_APPROVED` → `SUBMITTED` → `COMPLETED` (atau `REJECTED` di tiap fase — reject KR kembali ke `DRAFT`, reject assessment kembali ke `KR_APPROVED`). Untuk KPI (BSC) gunakan prefix `/performance/kpi/*` — lihat [8.4](#84-contoh-penggunaan-kpi-bsc-end-to-end). Contoh approval flow & employee movement ada di [8.5](#85-contoh-penggunaan-approval-engine--employee-movement).

### 8.4 Contoh Penggunaan KPI (BSC) — End-to-End

Alur lengkap pengelolaan **KPI (Balanced Scorecard)** dari perspektif sampai dashboard & scoring komponen. Semua endpoint berada di `/api/v1/tenant/performance/kpi/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Perspektif BSC → Template KPI (scope org) → Indikator → Evaluasi (snapshot) → Input Actual → Workflow → Dashboard → Scoring (per evaluasi & batch)
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), serta UUID `organization_id`, `period_id` (performance period), dan `employee_id` yang sudah ada.

**Step 1 — Buat perspektif BSC** (Financial, Customer, Internal Process, Learning & Growth):

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/perspectives \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Financial",
    "description": "Perspektif keuangan",
    "sort_order": 1
  }'
```

> Ulangi untuk perspektif lain (`Customer`, `Internal Process`, `Learning & Growth`). Catat UUID tiap perspektif — dipakai saat membuat indikator.

**Step 2 — Buat template KPI:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/templates \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "name": "KPI Sales Team — Q3 2026",
    "description": "Target KPI tim penjualan",
    "effective_date": "2026-07-01",
    "expired_date": "2026-09-30"
  }'
```

**Response 201:**
```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": "<template-uuid>",
    "organization_id": "<org-uuid>",
    "name": "KPI Sales Team — Q3 2026",
    "status": "DRAFT",
    "created_at": "2026-08-06T10:00:00Z"
  }
}
```

**Step 3 — (Opsional) Lihat opsi organisasi untuk scoped template:**

Saat menyiapkan template KPI yang dibatasi pada organisasi turunan (sebelum mengisi `organization_id` pada Step 2), HR dapat mengambil daftar organisasi yang boleh dipilih — hanya organisasi **di bawah** organisasi milik user saat ini (hierarki `ParentID`):

```bash
curl "http://localhost:8080/api/v1/tenant/performance/kpi/templates/organization-scope" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [
#     { "id": "<child-org-uuid>", "name": "Divisi Penjualan" },
#     { "id": "<child-org-uuid>", "name": "Cabang Surabaya" }
#   ] }
```

> Organisasi milik user sendiri **tidak** disertakan — hanya turunannya. Array kosong jika user tidak memiliki organisasi turunan.

**Step 4 — Buat indikator KPI di dalam template** (linked ke perspektif):

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/indicators \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "performance_template_id": "<template-uuid>",
    "perspective_id": "<perspective-uuid>",
    "indicator_type": "MAXIMIZATION",
    "title": "Pendapatan penjualan",
    "weight": 50,
    "target_value": 2000000000,
    "unit_of_measurement": "IDR",
    "formula_type": "HIGHER_BETTER",
    "target_type": "CURRENCY",
    "is_required": true
  }'
```

> `indicator_type` hanya `MAXIMIZATION` / `MINIMIZATION`; `formula_type`: `MANUAL` / `HIGHER_BETTER` / `LOWER_BETTER` / `RANGE`. Ulangi untuk indikator lain (mis. `MINIMIZATION` untuk "Keluhan pelanggan").

**Step 5 — Buat evaluasi KPI untuk employee (snapshot dari template):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/snapshot \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "<employee-uuid>",
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "template_id": "<template-uuid>",
    "supervisor_id": "<supervisor-uuid>"
  }'
```

> Evaluasi dibuat dengan status `DRAFT` dan **menyalin (snapshot)** seluruh indikator dari template ke `evaluation_details` — perubahan template setelah ini tidak memengaruhi evaluasi yang sudah dibuat.

**Response 201:**
```json
{
  "success": true,
  "message": "Created successfully",
  "data": {
    "id": "<evaluation-uuid>",
    "employee_id": "<employee-uuid>",
    "organization_id": "<org-uuid>",
    "period_id": "<period-uuid>",
    "template_id": "<template-uuid>",
    "status": "DRAFT",
    "details": [
      { "id": "<detail-uuid>", "indicator_name": "Pendapatan penjualan", "target": 2000000000, "actual": 0 }
    ]
  }
}
```

**Step 6 — Input nilai aktual (actual) secara massal:**

Ambil UUID detail dari `GET /api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/full` lalu kirim nilai realisasi:

```bash
curl -X PUT http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/actuals \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "details": [
      { "detail_id": "<detail-uuid>", "actual": 1850000000, "remarks": "Realisasi akhir kuartal" }
    ]
  }'
```

**Step 7 — Hitung ulang skor (bisa diulang kapan saja):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/recalculate \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → achievement & final_score dihitung ulang dari actual vs target + formula tiap indikator
```

**Step 8 — Hitung ulang skor seluruh evaluasi dalam satu period (batch):**

Untuk re-kalkulasi massal (mis. setelah mengubah formula atau bobot komponen), jalankan batch scoring untuk **semua** evaluasi di sebuah period:

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/periods/<period-uuid>/recalculate-scoring \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": {
#     "period_id": "<period-uuid>",
#     "total": 42, "processed": 41, "failed": 1
#   } }
```

> `total` = jumlah evaluasi dalam period; `processed` = berhasil dihitung ulang; `failed` = gagal (mis. evaluasi tanpa komponen scoring aktif). Berbeda dengan `POST /evaluations/:id/calculate-scoring` yang hanya menghitung satu evaluasi.

**Step 9 — Workflow status (2-tahap: plan dulu, lalu actual):**

```bash
# 1) Employee mengajukan rencana KPI
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/submit \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → DRAFT → PLAN_SUBMITTED

# 2) Atasan menyetujui rencana
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/approve \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → PLAN_SUBMITTED → PLAN_APPROVED

# 3) Setelah periode berjalan, employee mengajukan realisasi (actual)
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/submit \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → PLAN_APPROVED → ACTUAL_SUBMITTED

# 4) Atasan menyetujui realisasi
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/approve \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → ACTUAL_SUBMITTED → ACTUAL_APPROVED

# 5) HR/atasan menyelesaikan evaluasi (hasil akhir terkunci)
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/complete \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → ACTUAL_APPROVED → COMPLETED

# (Opsional) Menolak dengan catatan — kembali ke status sebelumnya untuk revisi
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/reject \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "notes": "Actual perlu dilengkapi" }'
```

> Alternatif manual: `PUT /api/v1/tenant/performance/kpi/evaluations/<id>/status` dengan body `{ "status": "PLAN_APPROVED", "notes": "..." }` untuk transisi status langsung (nilai `status`: `DRAFT`, `PLAN_SUBMITTED`, `PLAN_APPROVED`, `ACTUAL_SUBMITTED`, `ACTUAL_APPROVED`, `COMPLETED`).

**Step 10 — Lihat dashboard KPI:**

```bash
# Dashboard employee (progress KPI sendiri)
curl "http://localhost:8080/api/v1/tenant/performance/kpi/dashboard/employee/<employee-uuid>?period_id=<period-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Dashboard manager (ringkasan tim + pending reviews)
curl "http://localhost:8080/api/v1/tenant/performance/kpi/dashboard/manager/<manager-uuid>?period_id=<period-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Dashboard HR (statistik penyelesaian, rating distribution, top/bottom performers)
curl "http://localhost:8080/api/v1/tenant/performance/kpi/dashboard/hr?period_id=<period-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

**Step 11 — Konfigurasi komponen scoring (Phase 5):**

Komponen scoring memecah nilai evaluasi menjadi beberapa komponen (mis. `KPI Target`, `Competency`, `Work Program`) dengan bobot yang bisa diatur per organisasi, lalu dihitung oleh scoring engine.

```bash
# 1) Buat komponen scoring (master data)
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/components \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "KPI_TARGET",
    "name": "KPI Target",
    "description": "Skor dari pencapaian target indikator KPI",
    "sort_order": 1,
    "is_active": true
  }'
# → ulangi untuk "Competency", "Work Program" — catat UUID tiap komponen

# 2) Aktifkan komponen + set bobot untuk organisasi (upsert, otomatis menambah/mengupdate)
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/organization-components \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "<org-uuid>",
    "component_id": "<component-uuid>",
    "weight": 70,
    "is_enabled": true,
    "sort_order": 1
  }'

# 3) Cek daftar komponen aktif + bobot untuk organisasi
curl "http://localhost:8080/api/v1/tenant/performance/kpi/organizations/<org-uuid>/components" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# 4) Jalankan scoring engine — hitung skor tiap komponen dari data terkait lalu simpan
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/calculate-scoring \
  -H "Authorization: Bearer $TENANT_TOKEN"

# 5) Lihat skor per komponen hasil scoring engine
curl "http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/components" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# 6) Isi skor manual untuk komponen yang tidak bisa dihitung otomatis
#    (mis. Work Program — wajib diisi reviewer; score berkisar 0-100)
curl -X PUT http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/components/<component-uuid> \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "score": 85 }'
```

> 💡 `weight` antar komponen per organisasi menentukan proporsi kontribusi tiap komponen ke nilai akhir evaluasi.

**Catatan tambahan KPI:**

| Endpoint | Deskripsi |
|---|---|
| `PUT /api/v1/tenant/performance/kpi/targets/:id` | Update target KPI (nilai target/aktual) |
| `POST /api/v1/tenant/performance/kpi/progress` | Catat progres realisasi per detail (riwayat tanggal) |
| `POST /api/v1/tenant/performance/kpi/comments` | Komentar/review evaluasi |
| `POST /api/v1/tenant/performance/kpi/attachments` | Lampirkan file bukti ke evaluation detail |
| `GET /api/v1/tenant/performance/kpi/evaluations/:id/full` | Detail evaluasi lengkap (details, targets, comments, attachments) |
| `GET /api/v1/tenant/performance/kpi/evaluations/:id/progress-summary` | Ringkasan progres keseluruhan evaluasi |
| `POST /api/v1/tenant/performance/kpi/components` | Master data komponen scoring (Phase 5) |
| `POST /api/v1/tenant/performance/kpi/organization-components` | Set bobot/aktifkan komponen per organisasi (upsert) |
| `POST /api/v1/tenant/performance/kpi/evaluations/:id/calculate-scoring` | Jalankan scoring engine (hitung skor per komponen) |

> 💡 Status evaluasi KPI (BSC): `DRAFT` → `PLAN_SUBMITTED` → `PLAN_APPROVED` → `ACTUAL_SUBMITTED` → `ACTUAL_APPROVED` → `COMPLETED`. Sama seperti OKR yang juga dua fase (`DRAFT → KR_SUBMITTED → KR_APPROVED → SUBMITTED → COMPLETED`) — KPI mewajibkan persetujuan rencana **dan** realisasi secara terpisah.

> 🔗 Ingin melihat alur pengajuan & persetujuan (approval flow, submit, approve/reject)? Lihat [8.5 Contoh Penggunaan: Approval Engine & Employee Movement](#85-contoh-penggunaan-approval-engine--employee-movement).

### 8.5 Contoh Penggunaan: Approval Engine & Employee Movement

Employee movement (promosi/mutasi/status change) bisa diajukan lewat **approval flow terpusat**: HR membuat draft movement → men-submit ke flow → approver menyetujui/menolak lewat approval engine → movement tereksekusi.

**Step 1 — Lihat module yang tersedia untuk approval flow:**

```bash
# Module yang aktif/disubscribe tenant — dipakai flow builder agar hanya menampilkan module valid
curl "http://localhost:8080/api/v1/tenant/approval/available-modules" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": ["leave", "reimbursement", "employeemovement", ...] }
```

**Step 2 — Buat approval flow untuk module movement:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/approval/flows \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module": "employeemovement",
    "name": "Movement Approval - 2 Levels",
    "version": 1
  }'
# → { "success": true, "data": { "id": "<flow-uuid>", "module": "employeemovement", ... } }
```

**Step 3 — (Opsional) Resolusi otomatis flow aktif per module:**

Alih-alih memilih `flow_id` secara manual, konsumen yang butuh auto-resolution (mis. submission KPI/OKR dua tahap) dapat meminta flow aktif untuk sebuah module:

```bash
curl "http://localhost:8080/api/v1/tenant/approval/active-flow?module=performance_kpi_target" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": {
#     "id": "<flow-uuid>", "module": "performance", "name": "Performance Approval", "version": 1, "is_active": true, "steps": [...]
#   } }
```

> Parameter `module` **wajib**. Jika tidak ada flow spesifik untuk sub-checkpoint (mis. `performance_kpi_target`), sistem otomatis **fallback** ke flow module induk (`performance`) hingga flow khusus dibuat — satu flow cukup untuk mencakup semua checkpoint di bawahnya.

**Step 4 — Buat draft movement (promosi):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/employee-movements/movements \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "<employee-uuid>",
    "movement_type": "promotion",
    "from_employment_id": "<from-employment-uuid>",
    "to_employment_id": "<to-employment-uuid>",
    "decision_letter_number": "SK/2026/0821",
    "decision_letter_date": "2026-08-07",
    "effective_date": "2026-09-01",
    "reason": "Kinerja melebihi target semester I"
  }'
# → { "success": true, "data": { "id": "<movement-uuid>", "status": "draft", ... } }
```

**Step 5 — Submit movement ke approval engine:**

```bash
# Hanya movement berstatus draft yang bisa di-submit; flow_id wajib diisi
curl -X POST http://localhost:8080/api/v1/tenant/employee-movements/movements/<movement-uuid>/submit \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "flow_id": "<flow-uuid>" }'
# → { "success": true, "data": { "id": "<movement-uuid>", "status": "pending_approval", "approval_instance_id": "<instance-uuid>", ... } }
```

**Step 6 — Approver menyetujui via approval engine:**

```bash
# (Opsional) Cek task pending untuk user approver
curl "http://localhost:8080/api/v1/tenant/approval/tasks/pending?page=1&per_page=20" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Setujui instance approval (action: APPROVE | REJECT)
curl -X POST http://localhost:8080/api/v1/tenant/approval/instances/<instance-uuid>/actions \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "action": "APPROVE", "note": "Disetujui" }'
```

**Step 7 — (Opsional) Eksekusi movement setelah disetujui:**

```bash
# Setelah status approved, HR mengeksekusi perpindahan (update employment efektif)
curl -X POST http://localhost:8080/api/v1/tenant/employee-movements/movements/<movement-uuid>/execute \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

> 💡 Status movement: `draft` → `pending_approval` → `approved` → `executed`. Jika approver **REJECT**, status menjadi `cancelled` (tidak ada status rejected khusus di EmployeeMovement). Endpoint `/movements/:id/approve` dan `/movements/:id/execute` tetap tersedia sebagai jalur manual tanpa approval engine.

---

### 8.6 Contoh Penggunaan: Attendance & Notifications

#### 8.6.1 Attendance — Kalender, Rekap & Koreksi

Kalender & rekap dipakai untuk tampilan dashboard kehadiran; koreksi dipakai karyawan untuk memperbaiki check-in/check-out yang salah atau tidak tercatat.

**Step 1 — Kalender kehadiran satu karyawan:**

```bash
curl "http://localhost:8080/api/v1/tenant/attendance/calendar?employee_id=<employee-uuid>&from=2026-07-01&to=2026-07-31" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [ { "id": "<session-uuid>", "work_date": "2026-07-01", "status": "completed", "work_minutes": 480, ... } ] }
```

> Query `employee_id`, `from`, dan `to` **wajib** diisi — jika kosong akan error `VALIDATION_ERROR`.

**Step 2 — Rekap kehadiran (summary) karyawan:**

```bash
curl "http://localhost:8080/api/v1/tenant/attendance/summary?employee_id=<employee-uuid>&from=2026-07-01&to=2026-07-31" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "employee_id": "<employee-uuid>", "total_sessions": 22, "present_days": 20,
#     "late_days": 2, "missing_checkin_days": 0, "missing_checkout_days": 1, "leave_days": 1,
#     "total_work_minutes": 10180, "total_overtime_minutes": 240 } }
```

**Step 3 — Laporan sesi semua karyawan (untuk HR):**

```bash
curl "http://localhost:8080/api/v1/tenant/attendance/reports/sessions?from=2026-07-01&to=2026-07-31" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [ { "id": "<session-uuid>", "employee_id": "<employee-uuid>", "work_date": "2026-07-01", ... } ] }
```

**Step 4 — Ajukan koreksi kehadiran:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/attendance/corrections \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "<employee-uuid>",
    "attendance_session_id": "<session-uuid>",
    "correction_type": "checkin_missing",
    "requested_checkin": "2026-07-15T08:05:00+07:00",
    "reason": "Lupa check-in saat masuk kantor",
    "flow_id": "<flow-uuid>"
  }'
# → 201 { "success": true, "data": { "id": "<correction-uuid>", "status": "pending", ... } }
```

**Step 5 — Cek status pengajuan koreksi:**

```bash
curl "http://localhost:8080/api/v1/tenant/attendance/corrections?employee_id=<employee-uuid>&page=1&per_page=20" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [ { "id": "<correction-uuid>", "status": "approved", ... } ], "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
```

> Jika `flow_id` diberikan, koreksi diproses lewat **Approval Engine** (`approval/instances/:id/actions`). Tanpa flow, status langsung jadi `approved`/`rejected` sesuai logika internal modul.

#### 8.6.2 Notification — Feed, Badge & Tandai Dibaca

Notifikasi in-app dikirim per **user_id** (bukan employee_id) dan otomatis dibuat oleh modul lain (mis. Approval saat ada instance baru).

**Step 1 — Lihat feed notifikasi (paginated, bisa filter belum dibaca):**

```bash
curl "http://localhost:8080/api/v1/tenant/notifications?is_read=false&page=1&per_page=20" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "data": [ { "id": "<notif-uuid>", "type": "approval",
#     "title": "Approval baru menunggu Anda", "body": "Permintaan cuti Andi menunggu persetujuan Anda.",
#     "reference_type": "approval_instance", "reference_id": "<instance-uuid>", "is_read": false } ],
#     "total": 5, "page": 1, "per_page": 20 } }
```

**Step 2 — Badge jumlah belum dibaca:**

```bash
curl "http://localhost:8080/api/v1/tenant/notifications/unread-count" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "unread_count": 3 } }
```

**Step 3 — Tandai satu notifikasi sudah dibaca:**

```bash
curl -X PATCH http://localhost:8080/api/v1/tenant/notifications/<notif-uuid>/read \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "success": true } }
```

**Step 4 — Tandai semua notifikasi sudah dibaca:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/notifications/read-all \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "success": true } }
```

> 💡 User hanya bisa membaca/menandai **notifikasi miliknya sendiri** — endpoint memakai identitas user dari token, bukan parameter.

---

### 8.7 Contoh Penggunaan: Leave — Kalender & Laporan Penggunaan

Kalender cuti dipakai untuk tampilan kalender per karyawan; laporan penggunaan dipakai HR untuk melihat semua permintaan cuti yang beririsan dengan rentang tanggal tertentu.

#### 8.7.1 Leave Calendar — Kalender Cuti Satu Karyawan

**Step 1 — Entri cuti harian karyawan dalam rentang tanggal:**

```bash
curl "http://localhost:8080/api/v1/tenant/leave/calendar?employee_id=<employee-uuid>&from=2026-07-01&to=2026-07-31" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [ { "leave_request_id": "<req-uuid>", "leave_date": "2026-07-02",
#     "day_fraction": 1.0, "leave_type_id": "<type-uuid>", "status": "approved" } ] }
```

> Query `employee_id`, `from`, dan `to` **wajib** diisi — jika kosong akan error `VALIDATION_ERROR`. Setiap entri mewakili **satu hari cuti** (`day_fraction` bisa 0.5 untuk cuti setengah hari).

#### 8.7.2 Leave Usage Report — Laporan Penggunaan Cuti (HR)

**Step 2 — Semua permintaan cuti yang beririsan dengan rentang tanggal (untuk HR):**

```bash
curl "http://localhost:8080/api/v1/tenant/leave/reports/usage?from=2026-07-01&to=2026-07-31" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": [ { "id": "<req-uuid>", "employee_id": "<employee-uuid>",
#     "leave_type_id": "<type-uuid>", "request_start_date": "2026-07-02", "request_end_date": "2026-07-03",
#     "requested_days": 2.0, "status": "approved", ... } ] }
```

> Query `from` dan `to` **wajib** diisi. Response non-paginated; bentuk item sama dengan `GET /api/v1/tenant/leave/requests` (`LeaveRequestResponse`). Jika perlu dikelompokkan per jenis cuti, cukup agregasi `leave_type_id` di sisi klien.

> 💡 Alur lengkap cuti (buat jenis → buat request → setujui via Approval Engine → cek balance) mengikuti pola yang sama dengan contoh di §8.3–§8.6; `calendar` & `reports/usage` adalah endpoint *read-only* untuk tampilan & laporan.

---

### 8.8 Contoh Penggunaan: Training & Development (End-to-End)

Alur lengkap pengelolaan **Training & Development**: dari perencanaan (plan + need) sampai operasional (session, attendance, assessment) dan evaluasi (form evaluasi, efektivitas, sertifikat). Semua endpoint berada di `/api/v1/tenant/trainings/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Training Plan (tahunan) → Training Need → Course & Category → Session → Attendance & Assessment → Evaluation Form → Efektivitas & Sertifikat → Laporan
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), serta UUID `organization_id` / `position_id` dan `employee_id` yang sudah ada. Modul Training harus aktif (license) untuk tenant.

**Step 1 — Buat rencana pelatihan tahunan (Training Plan):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/trainings/plans \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PLN-2026",
    "name": "Rencana Pelatihan 2026",
    "year": 2026,
    "description": "Program pengembangan kompetensi tahunan",
    "status": "ACTIVE"
  }'
# → 201 { "success": true, "data": { "id": "<plan-uuid>", "status": "ACTIVE", ... } }
```

**Step 2 — Catat kebutuhan pelatihan (Training Need):**

Kebutuhan bisa di-input manual atau berasal dari sumber lain (performance/competency/career/succession/workforce/onboarding).

```bash
curl -X POST http://localhost:8080/api/v1/tenant/trainings/needs \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "<org-uuid>",
    "course_id": "<course-uuid>",
    "reason": "Gap kompetensi tim penjualan",
    "priority": "HIGH",
    "source_type": "COMPETENCY",
    "status": "OPEN"
  }'
# → 201 { "success": true, "data": { "id": "<need-uuid>", "status": "OPEN", ... } }
```

> 💡 **Catatan `source_type` (enum lengkap):** `MANUAL`, `PERFORMANCE`, `COMPETENCY`, `CAREER`, `SUCCESSION`, `COMPLIANCE`, `WORKFORCE`, **`ONBOARDING`** (S-7). Nilai `ONBOARDING` **tidak perlu di-input manual** — dibuat otomatis oleh sistem saat onboarding bertransisi ke `COMPLETED` (handoff **Recruitment → Training**), dengan anti-duplikat per employee/course. Manual `source_type` umumnya dipakai untuk `MANUAL`, `PERFORMANCE`, `COMPETENCY`, atau `CAREER`.

**Step 3 — Buat course di dalam kategori:**

```bash
# Buat kategori dulu
curl -X POST http://localhost:8080/api/v1/tenant/trainings/categories \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "code": "SOFT", "name": "Soft Skill", "description": "Pengembangan soft skill" }'

# Lalu buat course
curl -X POST http://localhost:8080/api/v1/tenant/trainings/courses \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": "<category-uuid>",
    "code": "NEG-101",
    "name": "Negotiation Skills",
    "duration_hour": 16,
    "min_score": 70,
    "course_type": "SOFT_SKILL",
    "delivery_type": "IN_HOUSE",
    "is_certified": true
  }'
# → 201 { "success": true, "data": { "id": "<course-uuid>", ... } }

# (Opsional) Tambah objective & prasyarat course
curl -X POST http://localhost:8080/api/v1/tenant/trainings/courses/<course-uuid>/objectives \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "objective": "Mampu melakukan negosiasi kontrak", "sort_order": 1 }'

curl -X POST http://localhost:8080/api/v1/tenant/trainings/courses/<course-uuid>/prerequisites \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "prerequisite_type": "COMPETENCY", "is_required": true }'
```

**Step 4 — Buat session pelatihan (jadwal + trainer):**

```bash
# Daftarkan provider & trainer
curl -X POST http://localhost:8080/api/v1/tenant/trainings/providers \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "code": "VND-01", "name": "Training House ID", "type": "EXTERNAL" }'

curl -X POST http://localhost:8080/api/v1/tenant/trainings/trainers \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "type": "EXTERNAL", "provider_id": "<provider-uuid>", "name": "Bpk. Andi Trainer" }'

# Buat session
curl -X POST http://localhost:8080/api/v1/tenant/trainings/sessions \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "<course-uuid>",
    "session_code": "NEG-101-AUG",
    "trainer_name": "Bpk. Andi Trainer",
    "location": "Jakarta",
    "start_date": "2026-08-17",
    "end_date": "2026-08-19",
    "max_quota": 20,
    "delivery_mode": "ONSITE",
    "provider_type": "EXTERNAL"
  }'
# → 201 { "success": true, "data": { "id": "<session-uuid>", "status": "DRAFT", ... } }
```

**Step 5 — Daftarkan peserta & catat kehadiran:**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/trainings/participants \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "session_id": "<session-uuid>", "employee_id": "<employee-uuid>" }'

# Catat kehadiran per hari
curl -X POST http://localhost:8080/api/v1/tenant/trainings/sessions/<session-uuid>/attendance \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "participant_id": "<participant-uuid>", "attendance_date": "2026-08-17", "status": "PRESENT", "check_in": "08:00", "check_out": "17:00" }'

# Buat assessment & input nilai peserta
curl -X POST http://localhost:8080/api/v1/tenant/trainings/sessions/<session-uuid>/assessments \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "name": "Post-Test Negosiasi", "type": "POST_TEST", "max_score": 100, "passing_score": 70, "is_required": true }'

curl -X POST http://localhost:8080/api/v1/tenant/trainings/assessments/<assessment-uuid>/results \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "participant_id": "<participant-uuid>", "score": 85 }'
```

**Step 6 — Evaluasi pelatihan (form + jawaban + efektivitas):**

```bash
# Buat form evaluasi untuk session
curl -X POST http://localhost:8080/api/v1/tenant/trainings/evaluation-forms \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "session_id": "<session-uuid>", "name": "Evaluasi Pelatihan Negosiasi", "is_active": true }'

# Tambah pertanyaan
curl -X POST http://localhost:8080/api/v1/tenant/trainings/evaluation-forms/<form-uuid>/questions \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "question": "Seberapa relevan materi dengan pekerjaan?", "question_type": "RATING", "is_required": true }'

# Peserta mengisi jawaban
curl -X POST http://localhost:8080/api/v1/tenant/trainings/evaluation-forms/<form-uuid>/participants/<participant-uuid>/answers \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "answers": [ { "question_id": "<question-uuid>", "answer": "5" } ] }'

# Penilaian efektivitas (before/after score) oleh atasan
curl -X POST http://localhost:8080/api/v1/tenant/trainings/effectiveness \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": "<participant-uuid>",
    "assessment_date": "2026-09-01",
    "before_score": 60,
    "after_score": 85
  }'
```

**Step 7 — Terbitkan sertifikat & lihat laporan:**

```bash
# Generate sertifikat untuk peserta (isikan nomor + file URL hasil upload)
curl -X POST http://localhost:8080/api/v1/tenant/trainings/participants/<participant-uuid>/certificate \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "certification_id": "<certification-uuid>", "certificate_file_url": "/uploads/certs/neg-101-a.pdf", "expiry_date": "2028-08-19" }'

# Riwayat pelatihan karyawan
curl "http://localhost:8080/api/v1/tenant/trainings/history?employee_id=<employee-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Laporan partisipasi & biaya (HR)
curl "http://localhost:8080/api/v1/tenant/trainings/reports/participation?course_id=<course-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"
curl "http://localhost:8080/api/v1/tenant/trainings/reports/cost?year=2026" \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Kartu dashboard training
curl "http://localhost:8080/api/v1/tenant/trainings/reports/dashboard" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → { "success": true, "data": { "total_courses": 12, "total_sessions": 8, "completion_rate": 91.5, ... } }
```

**Catatan tambahan Training:**

| Endpoint | Deskripsi |
|---|---|
| `POST /api/v1/tenant/trainings/requests` (+ `/:id/submit`, `/:id/cancel`) | Permintaan pelatihan karyawan — diproses lewat **Central Approval** (module `training_request`) |
| `GET/POST /api/v1/tenant/trainings/mandatories` | Kebijakan pelatihan wajib per organisasi/posisi (compliance) |
| `GET /api/v1/tenant/trainings/plans/:id/items`, `POST .../items` | Item rencana (course target + estimasi biaya) dalam plan tahunan |
| `GET/POST /api/v1/tenant/trainings/sessions/:id/costs` | Komponen biaya per session (biaya trainer, venue, materi, dll) |
| `GET/POST /api/v1/tenant/trainings/sessions/:id/documents` | Dokumen session (proposal, quotation, invoice, kontrak) |
| `GET/POST /api/v1/tenant/trainings/certifications` | Master sertifikasi (badan penerbit, masa berlaku, renewal) |
| `GET /api/v1/tenant/trainings/reports/compliance` | Laporan kepatuhan pelatihan wajib per karyawan |

> 💡 Alur **Central Approval**: buat request → `POST /requests/:id/submit` (kirim ke approval engine) → approver bertindak via `/approval/instances/:id/actions` — lihat [8.5](#85-contoh-penggunaan-approval-engine--employee-movement). Untuk perpindahan karyawan (promosi/mutasi) lihat juga contoh employee movement di [8.5](#85-contoh-penggunaan-approval-engine--employee-movement), dan karier karyawan (career history, eligibility promosi) di endpoint `/api/v1/tenant/employee-movements/employees/:employeeId/*`.

---

### 8.9 Contoh Penggunaan: Workforce Intelligence — Quality of Hire

Alur penggunaan metrik **Quality of Hire (S-6)**: Workforce Intelligence **membaca** data operasional lintas modul (Recruitment/Training/Performance menyediakan data) dan menghitung skor komposit kualitas hire — dari interview score, onboarding completion (proxy probation), performance evaluasi terbaru, sampai retention. Endpoint berada di `/api/v1/tenant/workforce-intelligence/analytics/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Data operasional (interviews, onboarding, performance evaluations, employment) → Quality of Hire (overall + breakdown by source/requisition/organization)
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), serta data hire yang sudah ada: kandidat ACCEPTED (`job_applications`), interview dengan score (`interviews`), onboarding (`onboarding`), evaluasi performance, dan data employment. Modul Workforce Intelligence harus aktif (license) untuk tenant.

**Step 1 — Lihat metrik agregat Quality of Hire:**

```bash
curl -X GET http://localhost:8080/api/v1/tenant/workforce-intelligence/analytics/quality-of-hire \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → 200 { "success": true, "data": {
#     "overall_score": 78.5,
#     "hires_analyzed": 12,
#     "recruitment_match_score": 0,
#     "interview_score": 81.2,
#     "assessment_score": 0,
#     "onboarding_completion_rate": 75.0,
#     "performance_score": 79.0,
#     "retention_rate": 83.3,
#     "by_source": [
#       { "key": "JOB_PORTAL", "hires": 5, "score": 76.4 },
#       { "key": "REFERRAL", "hires": 4, "score": 84.1 },
#       { "key": "SOCIAL_MEDIA", "hires": 3, "score": 72.8 }
#     ],
#     "by_requisition": [
#       { "key": "<requisition-uuid>", "hires": 6, "score": 80.2 },
#       { "key": "<requisition-uuid>", "hires": 6, "score": 76.8 }
#     ],
#     "by_organization": [
#       { "key": "<org-uuid>", "hires": 8, "score": 79.6 },
#       { "key": "<org-uuid>", "hires": 4, "score": 76.3 }
#     ]
#   } }
```

**Step 2 — Interpretasi komponen skor:**

| Field | Sumber data | Catatan |
|---|---|---|
| `interview_score` | Rata-rata `interviews.score` per hire | Semakin tinggi semakin baik |
| `onboarding_completion_rate` | % hire dengan onboarding `COMPLETED` | Proxy masa percobaan |
| `performance_score` | Rata-rata `performance_evaluations.final_score` (evaluasi selesai) | Skor evaluasi terbaru |
| `retention_rate` | % hire dengan employment aktif (`effective_end_date IS NULL`) | Retensi karyawan baru |
| `recruitment_match_score` / `assessment_score` | **Placeholder `0`** | Data kompetensi kandidat (G-9) & assessment belum dikumpulkan |

**Step 3 — Bandingkan dengan Executive Summary untuk konteks lebih luas:**

```bash
curl -X GET http://localhost:8080/api/v1/tenant/workforce-intelligence/executive/summary \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

> 💡 **Catatan:** `overall_score` adalah rata-rata komponen yang **tersedia** (interview, onboarding, performance, retention) — komponen placeholder tidak ikut dihitung. Breakdown `by_source` membantu menilai channel rekrutmen mana yang menghasilkan hire berkualitas, sedangkan `by_requisition`/`by_organization` untuk analisis per kebutuhan/unit.

### 8.10 Contoh Penggunaan: Career Intelligence — Succession Gaps & Fallback External Recruitment

Alur penggunaan **Succession Gaps (S-5)**: Career Intelligence menandai **posisi kunci** (memiliki ≥1 succession plan ACTIVE) yang tidak punya successor siap (`READY_NOW`). Posisi seperti itu membutuhkan **fallback external recruitment** — requisition dibuat dengan `reason_type=SUCCESSION_GAP` + `succession_position_id`. Endpoint berada di `/api/v1/tenant/career-intelligence/successions/*` dan `/api/v1/tenant/recruitment/*`.

```
Succession plans (ACTIVE) → deteksi posisi kunci tanpa successor READY_NOW → gap (requires_external_recruitment=true) → Job Requisition fallback (reason_type=SUCCESSION_GAP)
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), succession plan ACTIVE untuk posisi kunci (`career-intelligence/successions`), serta data employee dengan readiness (`READY_NOW`/lainnya). Modul Career Intelligence & Recruitment harus aktif (license) untuk tenant.

**Step 1 — Lihat daftar succession gaps:**

```bash
curl -X GET http://localhost:8080/api/v1/tenant/career-intelligence/successions/gaps \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → 200 { "success": true, "data": [
#     {
#       "position_id": "<position-uuid>",
#       "position_title": "Direktur Operasional",
#       "organization_id": "<org-uuid>",
#       "successor_count": 2,
#       "ready_now_count": 0,
#       "has_ready_successor": false,
#       "requires_external_recruitment": true
#     },
#     {
#       "position_id": "<position-uuid>",
#       "position_title": "Head of Engineering",
#       "organization_id": "<org-uuid>",
#       "successor_count": 3,
#       "ready_now_count": 1,
#       "has_ready_successor": true,
#       "requires_external_recruitment": false
#     }
#   ] }
```

**Step 2 — Identifikasi posisi yang butuh fallback external recruitment:**

Filter hasil di mana `requires_external_recruitment: true` (tidak ada satu pun successor `READY_NOW`). Pada contoh di atas, **Direktur Operasional** adalah kandidat fallback.

**Step 3 — Buat Job Requisition fallback (reason_type=SUCCESSION_GAP):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/recruitment/requisitions \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": "<org-uuid>",
    "position_id": "<position-uuid>",
    "title": "Direktur Operasional",
    "reason_type": "SUCCESSION_GAP",
    "succession_position_id": "<position-uuid>",
    "slots_available": 1,
    "target_start_date": "2026-10-01"
  }'
# → 201 { "success": true, "data": { "id": "<requisition-uuid>", "reason_type": "SUCCESSION_GAP", ... } }
```

> 💡 **Catatan:** `succession_position_id` menautkan requisition ke posisi kunci (`positions.id`) yang gap-nya ditandai S-5. Nilai `reason_type` lain: `NEW_POSITION`, `REPLACEMENT`, `EXPANSION`, `WORKFORCE_GAP`. Setelah requisition disetujui, alur rekrutmen lanjut seperti biasa (candidate → interview → offer → onboarding) — dan Quality of Hire di §8.9 akan mengevaluasi hasilnya.

## 9. Error Codes

| Kode | HTTP | Arti |
|---|---|---|
| `UNAUTHORIZED` | 401 | Token tidak valid / tidak disertakan |
| `INVALID_CREDENTIALS` | 401 | Email/password salah (login tenant) |
| `INVALID_REFRESH_TOKEN` | 401 | Refresh token invalid/expired |
| `FORBIDDEN` | 403 | RBAC: tidak punya permission |
| `TENANT_REQUIRED` | 403 | Context tenant (`company_id`) tidak ada |
| `MODULE_NOT_LICENSED` | 403 | Modul tidak aktif untuk company |
| `QUOTA_EXCEEDED` | 403 | Kuota employee (on-premise) terlampaui |
| `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `VALIDATION_ERROR` | 400 | Validasi gagal (detail di `error.fields`) |
| `BAD_REQUEST` | 400 | Request tidak valid |
| `SUSPEND_FAILED` / `ACTIVATE_FAILED` / `TERMINATE_FAILED` | 409 | Prasyarat status company tidak terpenuhi |
| `ROTATE_FAILED` | 409 | Rotasi kredensial DB tenant gagal |
| `INTERNAL_ERROR` | 500 | Kesalahan server |

---

## 10. Maintenance Dokumen API (Makefile)

Dokumen API dijaga **tersinkronisasi dengan kode** lewat tooling berikut:

| Perintah | Fungsi |
|---|---|
| `make check-openapi` | Cek semua endpoint yang terdaftar di `routes.go` sudah terdokumentasi di `openapi.json` (jalankan `scripts/check_missing_openapi.py`) |
| `make docs` | Jalankan `check-openapi` lalu **regenerate** `docs/openapi-report.md` (jalankan `scripts/generate_openapi_report.py`) |
| `make test` | Test semua (verbose + race detector) |
| `make lint` / `make vet` | Golangci-lint / go vet |

**Workflow saat ada perubahan API:**

```bash
cd backend

# 1. Ubah kode (handler/service/routes)
# 2. Tambahkan endpoint baru ke openapi.json (backend/internal/pkg/docs/openapi.json)
# 3. Verifikasi tidak ada endpoint yang lupa didokumentasikan
make check-openapi

# 4. Regenerate laporan markdown
make docs
```

> `make check-openapi` **gagal (exit code 1)** jika ada endpoint missing ATAU phantom (dokumentasi yang tidak terdaftar di kode) — ini memaksa dokumen API selalu sinkron.

---

## Lampiran: Ringkasan Perintah Makefile

```text
=== Build ===            make build, make build-installer
=== Run ===              make run, make run-hot
=== Test ===             make test, test-verbose, test-short, test-pkg pkg=,
                         test-employee, test-organization, test-cache*,
                         bench-cache, coverage, cover-view, cover-func
=== Lint ===             make lint, make vet
=== API Docs ===         make check-openapi, make docs, make db-docs, make check-db-docs,
                         make arch-report, make check-arch-report
=== Database ===         make migrate, make seed, make seed-modules
=== Docker ===           make docker, make docker-compose-up, make docker-compose-down
=== Utilities ===        make tidy, make clean, make help
```

---

*Dokumen disusun dari analisis `backend/Makefile`, `backend/cmd/server/main.go`, `backend/internal/pkg/router/router.go`, dan konvensi response `backend/internal/pkg/httputil/response.go`.*
