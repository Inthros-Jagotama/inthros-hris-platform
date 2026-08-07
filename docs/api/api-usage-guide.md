# HRIS Platform — Panduan Penggunaan API (API Usage Guide)

> 🔗 **Index dokumentasi:** [`docs/README.md`](../README.md)  
> **Terkait:** [`openapi-report.md`](../openapi-report.md) · [`deployment-guide.md`](../deployment-guide.md) · [`platform-architecture-design.md`](../platform-architecture-design.md)

Panduan praktis **cara menggunakan API** HRIS Platform: dari menjalankan server, autentikasi, format request/response, sampai contoh pemanggilan end-to-end (curl).

> 📖 Dokumen ini berfokus pada **cara pakai**. Untuk daftar lengkap seluruh 800 endpoint + skema, lihat:
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
| `docs/openapi-report.md` | Laporan markdown statis (800 endpoint, 450 paths, 497 schemas, 32 tag) |

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
| Payroll | `GET /api/v1/tenant/payroll/...` |
| Leave | `GET /api/v1/tenant/leave/...` |
| Performance — Master Data | `GET/POST /api/v1/tenant/performance/periods`, `.../ratings`, `.../indicator-formulas`, `.../logs` |
| Performance — KPI | `GET/POST /api/v1/tenant/performance/kpi/templates`, `.../kpi/indicators`, `.../kpi/evaluations`, `.../kpi/dashboard/hr` |
| Performance — OKR | `GET/POST /api/v1/tenant/performance/okr/templates`, `.../okr/objectives`, `.../okr/key-results`, `.../okr/evaluations`, `.../okr/dashboard/hr` |
| Workforce Intelligence | `GET /api/v1/tenant/workforce-intelligence/executive/summary` |
| Package (subscribe) | `POST /api/v1/tenant/packages/:id/subscribe` |
| Company self-service | `GET/PUT /api/v1/tenant/companies/me` |
| Module aktif | `GET /api/v1/tenant/company-modules` |

> 🔍 Daftar lengkap per modul: lihat [`docs/openapi-report.md`](../openapi-report.md).

### 8.3 Contoh Penggunaan OKR (End-to-End)

Alur lengkap pengelolaan **OKR** dari membuat template sampai menyelesaikan evaluasi. Semua endpoint berada di `/api/v1/tenant/performance/okr/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Template OKR → Objective → Key Result → Evaluasi (snapshot) → Input Actual → Workflow → Dashboard HR
```

> **Prasyarat:** `TENANT_TOKEN` dari login tenant (lihat Step 3 di 8.1), serta UUID `organization_id`, `period_id` (performance period), dan `employee_id` yang sudah ada.

**Step 1 — Buat OKR template:**

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

**Step 2 — Buat objective di dalam template:**

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

**Step 3 — Buat key result di dalam objective (target terukur):**

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
```

> Ulangi untuk key result lain (mis. `PERCENTAGE` untuk "Growth pelanggan baru 20%"). Verifikasi isi template:
> `GET /api/v1/tenant/performance/okr/templates/<template-uuid>/objectives`

**Step 4 — Buat evaluasi OKR untuk employee (snapshot dari template):**

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

> Evaluasi dibuat dengan status `DRAFT` dan **menyalin (snapshot)** seluruh objective & key result dari template ke `evaluation_details` — perubahan template setelah ini tidak memengaruhi evaluasi yang sudah dibuat.

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
    "details": [
      { "id": "<detail-uuid>", "key_result_title": "Meraih target penjualan Rp 2 Miliar", "target_value": 2000000000, "actual_value": 0 }
    ]
  }
}
```

**Step 5 — Input nilai aktual (actual) secara massal:**

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

**Step 6 — Hitung ulang skor (bisa diulang kapan saja):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/recalculate \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → achievement & final_score dihitung ulang dari actual vs target + formula tiap key result
```

**Step 7 — Workflow status (DRAFT → SUBMITTED → APPROVED → COMPLETED):**

```bash
# Employee mengajukan evaluasi
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/submit \
  -H "Authorization: Bearer $TENANT_TOKEN"

# Atasan menyetujui
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/approve \
  -H "Authorization: Bearer $TENANT_TOKEN"

# (Opsional) menolak dengan catatan — evaluasi kembali ke revisi
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/reject \
  -H "Authorization: Bearer $TENANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ "notes": "Actual perlu dilengkapi" }'

# Menyelesaikan evaluasi (hasil akhir terkunci)
curl -X POST http://localhost:8080/api/v1/tenant/performance/okr/evaluations/<evaluation-uuid>/complete \
  -H "Authorization: Bearer $TENANT_TOKEN"
```

**Step 8 — Lihat dashboard OKR HR (ringkasan seluruh evaluasi):**

```bash
curl "http://localhost:8080/api/v1/tenant/performance/okr/dashboard/hr?period_id=<period-uuid>" \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → total_evaluations, sebaran status (draft/submitted/approved/completed), skor & achievement rata-rata
```

**Catatan tambahan OKR:**

| Endpoint | Deskripsi |
|---|---|
| `POST /api/v1/tenant/performance/okr/templates/:id/duplicate` | Duplikat template (beserta objective & key results) sebagai template baru |
| `GET /api/v1/tenant/performance/okr/objectives/:id/key-results` | List key results dalam sebuah objective |
| `POST /api/v1/tenant/performance/okr/progress` | Catat progres/check-in per detail (riwayat tanggal) |
| `POST /api/v1/tenant/performance/okr/comments` | Komentar/review evaluasi (dukung reply via `parent_id`) |
| `POST /api/v1/tenant/performance/okr/attachments` | Lampirkan file bukti ke evaluation detail |

> 💡 Status evaluasi OKR: `DRAFT` → `SUBMITTED` → `APPROVED` → `COMPLETED` (atau `REJECTED`). Untuk KPI (BSC) gunakan prefix `/performance/kpi/*` — lihat [8.4](#84-contoh-penggunaan-kpi-bsc-end-to-end).

### 8.4 Contoh Penggunaan KPI (BSC) — End-to-End

Alur lengkap pengelolaan **KPI (Balanced Scorecard)** dari perspektif sampai dashboard & scoring komponen. Semua endpoint berada di `/api/v1/tenant/performance/kpi/*` dan memerlukan `Authorization: Bearer <tenant_token>`.

```
Perspektif BSC → Template KPI → Indikator → Evaluasi (snapshot) → Input Actual → Workflow → Dashboard → Scoring Komponen
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

**Step 3 — Buat indikator KPI di dalam template** (linked ke perspektif):

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

**Step 4 — Buat evaluasi KPI untuk employee (snapshot dari template):**

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

**Step 5 — Input nilai aktual (actual) secara massal:**

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

**Step 6 — Hitung ulang skor (bisa diulang kapan saja):**

```bash
curl -X POST http://localhost:8080/api/v1/tenant/performance/kpi/evaluations/<evaluation-uuid>/recalculate \
  -H "Authorization: Bearer $TENANT_TOKEN"
# → achievement & final_score dihitung ulang dari actual vs target + formula tiap indikator
```

**Step 7 — Workflow status (2-tahap: plan dulu, lalu actual):**

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

**Step 8 — Lihat dashboard KPI:**

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

**Step 9 — Konfigurasi komponen scoring (Phase 5):**

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

> 💡 Status evaluasi KPI (BSC): `DRAFT` → `PLAN_SUBMITTED` → `PLAN_APPROVED` → `ACTUAL_SUBMITTED` → `ACTUAL_APPROVED` → `COMPLETED`. Berbeda dari OKR yang satu-tahap (`DRAFT → SUBMITTED → APPROVED → COMPLETED`) — KPI mewajibkan persetujuan rencana **dan** realisasi secara terpisah.

---

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
=== API Docs ===         make check-openapi, make docs
=== Database ===         make migrate, make seed, make seed-modules
=== Docker ===           make docker, make docker-compose-up, make docker-compose-down
=== Utilities ===        make tidy, make clean, make help
```

---

*Dokumen disusun dari analisis `backend/Makefile`, `backend/cmd/server/main.go`, `backend/internal/pkg/router/router.go`, dan konvensi response `backend/internal/pkg/httputil/response.go`.*
