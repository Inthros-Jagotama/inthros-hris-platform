# HRIS Platform — Deployment Guide

> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md)  
> **Terkait:** [`platform-architecture-design.md`](platform-architecture-design.md) · [`api/api-usage-guide.md`](api/api-usage-guide.md) · [`openapi-report.md`](openapi-report.md)

Panduan lengkap deployment untuk **dua mode distribusi**:

| Mode | Karakteristik | Lisensi Modul | Cocok Untuk |
|---|---|---|---|
| **Subscription SaaS** (multi-tenant) | Satu instance melayani banyak tenant; tiap tenant punya database sendiri | Per-company via tabel `company_modules` di platform DB | Vendor SaaS, hosting multi-customer |
| **On-Premise** (dedicated) | Satu instance untuk satu klien; data 100% di infrastruktur klien | Via file `.lic` RSA yang ditandatangani vendor | Enterprise dengan kebijakan data on-premise |

Kedua mode memakai **binari yang sama** (`cmd/server`) — perbedaan hanya pada nilai konfigurasi `license.deployment_mode`.

---

## 1. Arsitektur Singkat

```
┌────────────────────────────────────────────────────────────────┐
│                        HRIS Platform Server                     │
│  cmd/server (Go, modular monolith)                              │
│                                                                  │
│  Platform DB ──► hris_platform (companies, users, licenses,      │
│                  modules, company_modules, RBAC)                 │
│                                                                  │
│  Tenant DB(s) ─► SATU per tenant (hris_<slug>)                   │
│                  - SaaS: dibuat otomatis saat provisioning       │
│                  - On-Premise: satu tenant = satu DB             │
│                                                                  │
│  Redis ────────► Distributed cache (2-tier + Pub/Sub invalidate) │
└────────────────────────────────────────────────────────────────┘
```

### Perbedaan Kunci Antara Mode

| Aspek | SaaS | On-Premise |
|---|---|---|
| `license.deployment_mode` | `saas` (default) | `on_premise` |
| Sumber lisensi modul | `company_modules` DB (via `modulemgmt`) | `.lic` RSA (`allowed_modules`) |
| Jumlah tenant | Banyak (multi-tenant) | Satu (single-tenant dedicated) |
| Provisioning tenant | Otomatis via API/CLI (buat DB + migrasi + seed) | Tidak diperlukan (DB dibuat saat setup awal) |
| Batas employee (`max_employees`) | Tidak ada (unlimited) | Wajib diisi di `.lic`, di-enforce runtime → 403 `QUOTA_EXCEEDED` |
| Pembuatan lisensi | Auto saat company signup / subscribe package | Via CLI `licensectl` oleh vendor |
| Enkripsi kredensial tenant DB | ✅ AES-256-GCM (`encryption_key`) | Tidak relevan (single DB) |

---

## 2. Prasyarat (Kedua Mode)

- Go 1.24+ (untuk build dari source) **atau** Docker
- Database: **MySQL 8.0+** atau **PostgreSQL 16+** (driver dikonfigurasi via `HRIS_DATABASE_DRIVER`)
- Redis 7+ (wajib untuk distributed cache)
- (Opsional) Mailpit untuk dev email, Asynqmon untuk queue monitoring

### Build Binary

```bash
cd backend

# Server utama
make build                 # → ./bin/server

# CLI installer (provisioning, seed, migrate, encrypt)
make build-installer       # → ./bin/installer

# CLI lisensi on-premise (vendor side)
go build -o bin/licensectl ./cmd/licensectl
```

### Setup Infrastruktur (Docker)

```bash
cd docker

# PostgreSQL + PgBouncer + Redis + API
docker compose --profile postgres up -d

# ATAU MySQL
docker compose --profile mysql up -d
```

---

## 3. Mode A — Subscription SaaS (Multi-Tenant)

### 3.1 Konfigurasi

```yaml
# config/config.yaml (atau via env)
license:
  deployment_mode: "saas"          # default — tidak wajib ditulis

database:
  driver: "postgres"               # atau "mysql"
  platform_host: "localhost"
  platform_port: 5432
  platform_db: "hris_platform"
  platform_user: "hris"
  platform_password: "hris_secret"

  tenant_host: "localhost"
  tenant_port: 5432
  tenant_super_user: "hris"
  tenant_super_password: "hris_secret"
```

> ⚠️ **Kredensial tenant DB**: di production wajib set `encryption_key` (AES-256-GCM, 32-byte hex = 64 karakter) agar kredensial tenant terenkripsi di platform DB. Generate: `installer encrypt-passwords`.

### 3.2 Menjalankan Server

```bash
cd backend
# Env wajib production:
export HRIS_JWT_SECRET="<strong-random-secret>"
export HRIS_DATABASE_PLATFORM_PASSWORD="..."
export HRIS_DATABASE_TENANT_SUPER_PASSWORD="..."

./bin/server --config ./config/config.yaml
# → migrasi platform otomatis dijalankan saat startup
```

### 3.3 Provisioning Tenant Baru

Ada **dua cara**:

**A. Via API** (disarankan — integrasi dengan flow signup):
```bash
curl -X POST http://localhost:8080/api/v1/platform/companies \
  -H "Authorization: Bearer <super_admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PT Contoh Jaya",
    "admin_name": "Budi",
    "admin_email": "admin@contoh.com",
    "admin_password": "admin123"
  }'
```
Alur otomatis: buat company → buat tenant DB → jalankan migrasi tenant → **seed master data** (religion, banks, region BPS ~90k rows, salary grades, PTKP, TER, BPJS) → **seed RBAC** (roles Admin/Employee + 64 permissions).

**B. Via CLI installer:**
```bash
./bin/installer provision --company=<company_id> --config ./config/config.yaml
./bin/installer seed-data --company=<company_id> --config ./config/config.yaml   # jika perlu seed ulang
./bin/installer seed-modules --config ./config/config.yaml                        # daftarkan module ke platform
```

### 3.4 Aktivasi Modul per Tenant

Modul tenant diaktifkan melalui salah satu:

1. **Package subscribe** (tenant side): `POST /api/v1/tenant/packages/:id/subscribe` — otomatis mengaktifkan module package & membuat license
2. **Module Management** (platform admin): aktifkan/nonaktifkan `company_modules` per company

Enforcement di runtime: `PlatformLicenseMiddleware` membaca `company_modules` (via cache Redis, key per company) → modul nonaktif ditolak dengan **403 `MODULE_NOT_LICENSED`**.

### 3.5 Operasional SaaS

- **Suspension/Activation/Termination**: dikelola di Platform Admin (Companies) — suspend menonaktifkan login tenant; terminate menghapus koneksi tenant
- **Pembaruan cache**: setiap aktivasi/nonaktivasi modul otomatis meng-invalidate cache lisensi company via Redis Pub/Sub
- **Connection pool**: platform 10/5/1jam; per-tenant 10/3/30mnt/5mnt idle → 50 tenant × 10 = 500 koneksi (rekomendasi PgBouncer untuk 1000+ tenant)

### 3.6 Alur Signup Company via API (Request/Response Lengkap)

Endpoint: **`POST /api/v1/platform/companies`** (super_admin auth)

**Request:**

```json
{
  "name": "PT Contoh Jaya",
  "npwp": "0123456789012345",
  "nib": "8123456789012",
  "address": "Jl. Sudirman No. 1, Jakarta",
  "email": "contact@contohjaya.com",
  "phone": "02112345678",
  "admin_name": "Budi Santoso",
  "admin_email": "admin@contohjaya.com",
  "admin_password": "rahasia123",
  "package_id": "<uuid-package>"   // opsional — auto-create license + aktivasi modul package
}
```

| Field | Wajib | Validasi |
|---|---|---|
| `name` | ✅ | min 3, max 255 |
| `admin_name` | ✅ | min 1 |
| `admin_email` | ✅ | format email |
| `admin_password` | ✅ | min 6 |
| `npwp` | ❌ | len 16 |
| `nib` | ❌ | max 25 |
| `email` | ❌ | format email |
| `phone` | ❌ | max 20 |
| `package_id` | ❌ | UUID package valid |

**Alur internal otomatis (service `Create`):**
1. Generate `slug` dari name (jika duplikat → tambah suffix 8-char random)
2. Simpan company (status `ACTIVE`)
3. **Provision tenant DB**: buat database `hris_<slug>` → jalankan SQL migrations (dialect sesuai driver) → **seed master data** (religion, bank, region BPS ~90k rows, salary grades, PTKP, TER, BPJS) → **seed RBAC** (roles Admin/Employee + 64 permissions)
   - Jika provisioning gagal → company otomatis diset `SUSPENDED` (log error)
4. Buat **admin user** role `company_admin` (bcrypt hash)
5. Jika `package_id` ada → auto-create **license** + aktivasi modul package

**Response 201 (success):**

```json
{
  "success": true,
  "data": {
    "id": "uuid-company",
    "name": "PT Contoh Jaya",
    "slug": "pt-contoh-jaya",
    "status": "ACTIVE",
    "admin_user": { "id": "uuid-user", "name": "Budi Santoso", "email": "admin@contohjaya.com", "role": "company_admin" },
    "license_info": { "id": "uuid-lic", "license_key": "xxx", "plan_type": "pro", "package_id": "uuid-package" },
    "created_at": "2026-07-31T10:00:00Z",
    "updated_at": "2026-07-31T10:00:00Z"
  }
}
```

> ℹ️ `admin_user` dan `license_info` hanya muncul jika pembuatannya berhasil (opsional). **`provisioning_info` tidak disertakan di response 201** — status provisioning (provisioned, is_active, db_name) tersedia via `GET /api/v1/platform/companies/:id`.

> ⚠️ Jika provisioning tenant gagal: response tetap 201 dengan `status: "SUSPENDED"` — periksa log server, lalu jalankan `installer provision --company=<id>` untuk retry (atau lihat status via `GET /companies/:id`).

### 3.7 Subscribe / Unsubscribe Package (Tenant Side)

**Subscribe** — `POST /api/v1/tenant/packages/:id/subscribe` (company_admin auth):

1. Validasi package ada & berstatus `published` (jika tidak → 400 `VALIDATION_ERROR`)
2. Auto-create **license** (plan `pro`, masa aktif 1 tahun dari hari ini) via `LicenseCreator`
3. **Aktivasi modul package** (`ActivatePackageModules`) → insert/update `company_modules`
4. **Invalidasi cache lisensi** company (Redis) agar middleware langsung mengenali modul baru

**Response 201:**

```json
{
  "success": true,
  "message": "package.subscribed",
  "data": {
    "license_id": "uuid-license",
    "license_key": "xxxx",
    "plan_type": "pro",
    "package_id": "uuid-package",
    "package_name": "Enterprise HR",
    "activated_modules": ["Employee Management", "Payroll", "Leave"]
  }
}
```

**Unsubscribe** — `POST /api/v1/tenant/packages/:id/unsubscribe`:

1. **Deaktivasi semua modul package** (`DeactivatePackageModules`) → set `company_modules.enabled = false`
2. **Invalidasi cache lisensi** company
3. **Suspend lisensi aktif** yang terkait package (`status = suspended`)

**Response:** `{ "success": true, "message": "package.unsubscribed" }`

> ℹ️ Endpoint `GET /api/v1/tenant/packages` (daftar package published) dan `GET /api/v1/public/packages` (tanpa auth) tersedia untuk browsing sebelum subscribe.

### 3.8 Suspend / Activate / Terminate Tenant

Endpoint platform (semua super_admin auth):

| Operasi | Endpoint | Prasyarat Status | Efek |
|---|---|---|---|
| **Suspend** | `POST /api/v1/platform/companies/:id/suspend` | `ACTIVE` | Deactivate tenant connection (`is_active=false`) + status → `SUSPENDED`; user tenant **tidak bisa login** |
| **Activate** | `POST /api/v1/platform/companies/:id/activate` | `SUSPENDED` | Reactivate tenant connection + status → `ACTIVE`; login tenant normal kembali |
| **Terminate** | `POST /api/v1/platform/companies/:id/terminate` | `ACTIVE` / `SUSPENDED` (bukan `TERMINATED`) | **DROP tenant database** + remove `tenant_connections` record + status → `TERMINATED`; **permanen, data hilang** |

Response error bila prasyarat status tidak terpenuhi: **409** dengan kode `SUSPEND_FAILED` / `ACTIVATE_FAILED` / `TERMINATE_FAILED`.

Response sukses: `{ "success": true, "message": "company.suspended|activated|terminated", "data": { ...CompanyResponse } }`.

> ⚠️ **Terminate tidak bisa di-undo.** Backup dulu via `POST /api/v1/platform/companies/:id/backup` (atau snapshot DB manual) sebelum terminate. Untuk retensi data audit, pertimbangkan hanya *suspend* + *soft delete* (`DELETE /:id`) yang menonaktifkan koneksi tanpa drop DB.

### 3.9 Rotasi Kredensial DB Tenant

Kredensial DB tenant disimpan di tabel **`tenant_connections`** (platform DB) dan **selalu dienkripsi AES-256-GCM** (wajib set `encryption_key` di production) — lihat `SaveTenantConnection` → `crypto.EncryptString`.

#### A. Via API (disarankan)

Endpoint: **`POST /api/v1/platform/companies/:id/rotate-credentials`** (super_admin auth) — backend menangani semuanya: `ALTER USER` di server DB + update `tenant_connections` terenkripsi + close cache koneksi tenant (reconnect otomatis memakai kredensial baru).

**Request:**

```json
{
  "new_password": "<opsional-min-8-max-128>"
}
```

| Field | Wajib | Keterangan |
|---|---|---|
| `new_password` | ❌ | Jika kosong/dilewati, backend **meng-generate password acak kuat** (24 char, charset aman SQL tanpa quote/backslash) dan mengembalikannya sekali di response |

**Alur internal (`Service.RotateCredentials` → `database.Manager.RotateTenantCredentials`):**
1. Validasi company ID (UUID valid, bukan `TERMINATED`)
2. Generate password acak kuat jika `new_password` kosong
3. Ambil `tenant_connections` (dekripsi password saat ini) → connect sebagai superuser (DSN tanpa database spesifik)
4. **`ALTER USER` sesuai dialect**: MySQL `ALTER USER 'user'@'%' IDENTIFIED BY '...'`; PostgreSQL `ALTER USER "user" WITH PASSWORD '...'`
5. Update `tenant_connections.password` terenkripsi AES-256-GCM (hanya setelah ALTER sukses)
6. Tutup koneksi cache tenant → koneksi berikutnya memakai kredensial baru

**Response 200:**

```json
{
  "success": true,
  "data": {
    "company_id": "uuid-company",
    "rotated": true,
    "new_password": "<hanya-jika-auto-generated>"
  }
}
```

> ⚠️ `new_password` di response hanya muncul jika backend yang generate (field kosong di request). Jika Anda kirim password sendiri, backend tidak mengembalikannya.

**Error handling:**

| Kondisi | Kode |
|---|---|
| `new_password` < 8 atau > 128 char | 400 `VALIDATION_ERROR` |
| Password mengandung `'` atau `\` (raw SQL literal) | 409 `ROTATE_FAILED` (ditolak demi cegah SQL injection) |
| Company `TERMINATED` / ID invalid | 409 `ROTATE_FAILED` |
| Company belum ter-provision (tidak ada record `tenant_connections`) | 409 `ROTATE_FAILED` — rotasi butuh record koneksi yang sudah ada |
| User tenant adalah superuser/root | ⚠️ Hanya warning log — lihat catatan di bawah |

Contoh:
```bash
curl -X POST http://localhost:8080/api/v1/platform/companies/<uuid>/rotate-credentials \
  -H "Authorization: Bearer <super_admin_token>" \
  -H "Content-Type: application/json" \
  -d '{}'   # auto-generate password baru
```

#### B. Manual (fallback / on-premise single DB)

Jika endpoint tidak tersedia (versi lama) atau hanya ingin rotasi satu akun DB:

1. **Rotasi password di level DB server**:
   ```sql
   -- PostgreSQL
   ALTER USER "hris_<slug>" WITH PASSWORD '<new-password>';
   -- MySQL
   ALTER USER 'hris_<slug>'@'%' IDENTIFIED BY '<new-password>';
   FLUSH PRIVILEGES;
   ```
2. **Perbarui record `tenant_connections`** dengan password baru — wajib terenkripsi AES-256-GCM (`encryption_key` sama). Gunakan `database.Manager.SaveTenantConnection` / `crypto.EncryptString` (jangan plaintext langsung).
3. **Restart server** agar pool koneksi tenant memakai kredensial baru.
4. **Verifikasi**: login tenant + akses endpoint (mis. `GET /api/v1/tenant/employees`) harus sukses.

#### Catatan Penting

- **Migrasi legacy plaintext**: jika ada kredensial tersimpan plaintext (sebelum fitur enkripsi):
  ```bash
  HRIS_ENCRYPTION_KEY=<64-char-hex> ./bin/installer encrypt-passwords --config ./config/config.yaml
  ```
  Idempotent — mendeteksi `LooksEncrypted` dan melewatkan yang sudah terenkripsi dengan kunci sama.
- **⚠️ Jangan rotasi akun superuser/root**: `company.Service.Create` memakai `root` untuk development — merotasi root akan memutus provisioning & drop tenant selanjutnya. Rotasi API hanya cocok untuk **dedicated DB user per tenant** (praktik terbaik production).
- **Asumsi host MySQL `'%'`**: `ALTER USER 'user'@'%'` gagal jika akun dibuat untuk host spesifik (`localhost`) — pastikan akun tenant dibuat dengan host `'%'`.

---

## 4. Mode B — On-Premise (Dedicated)

### 4.1 Konsep

Vendor menandatangani file `.lic` (RSA-SHA256) berisi:

```json
{
  "company_id": "uuid",
  "company_name": "PT Contoh Jaya",
  "expires_at": "2027-12-31T00:00:00Z",
  "allowed_modules": ["organization", "employee", "payroll"],
  "max_employees": 500
}
```

File `.lic` + public key dideploy di sisi klien. Server memverifikasi signature & expiry saat startup (**fail-closed**: jika invalid/expired → server **gagal start**). Daftar modul dari `.lic` menggantikan `company_modules` DB; `max_employees` di-enforce saat create employee.

### 4.2 Langkah 1 — Vendor: Generate Keypair & Lisensi

```bash
cd backend

# 1) Generate RSA keypair (private disimpan vendor, public diberikan ke klien)
go run ./cmd/licensectl gen-key --bits 2048 --out private.pem --pub public.pem

# 2) Generate .lic (semua flag wajib)
go run ./cmd/licensectl gen-lic \
  --priv private.pem \
  --out license.lic \
  --company-id 00000000-0000-0000-0000-000000000001 \
  --company "PT Contoh Jaya" \
  --expires 2027-12-31 \
  --modules organization,employee,payroll \
  --max-employees 500
```

**File yang dikirim ke klien:**
- `license.lic` (payload + signature)
- `public.pem` (public key untuk verifikasi)

> 🔐 **Jangan pernah** mengirim `private.pem` ke klien — itu kunci untuk menandatangani lisensi.

### 4.3 Langkah 2 — Klien: Konfigurasi & Deploy

```yaml
# config/config.yaml
license:
  deployment_mode: "on_premise"
  license_file: "/etc/hris/license.lic"
  public_key_file: "/etc/hris/public.pem"
```

Atau via environment (produksi disarankan):

```bash
export HRIS_LICENSE_DEPLOYMENT_MODE=on_premise
export HRIS_LICENSE_LICENSE_FILE=/etc/hris/license.lic
export HRIS_LICENSE_PUBLIC_KEY_FILE=/etc/hris/public.pem
```

### 4.4 Langkah 3 — Klien: Setup Database Tunggal

On-premise = satu tenant. Buat database tenant + jalankan migrasi + seed sekali saja:

```bash
# 1) Buat company & tenant DB via API super_admin (atau via installer)
./bin/installer provision --company=<company_id> --config ./config/config.yaml

# 2) Seed master data + RBAC (sekali saat setup)
./bin/installer seed-data --company=<company_id> --config ./config/config.yaml
```

> Untuk keperluan audit, seed idempotent — aman dijalankan ulang.

### 4.5 Langkah 4 — Jalankan Server

```bash
./bin/server --config ./config/config.yaml
```

Startup log yang diharapkan:
```
INFO On-premise license loaded  {"company": "PT Contoh Jaya", "expires_at": "2027-12-31T00:00:00Z", "max_employees": 500}
```

### 4.6 Perilaku Enforcement di On-Premise

| Skema | Perilaku |
|---|---|
| Modul di luar `allowed_modules` | **403 `MODULE_NOT_LICENSED`** (PlatformLicenseMiddleware) |
| Create employee saat jumlah = `max_employees` | **403 `QUOTA_EXCEEDED`** (EmployeeQuotaChecker → `Service.Create()`) |
| Toast frontend | Bilingual `employee.quota_exceeded` (EN/ID) |
| `.lic` expired / signature invalid / key salah | **Server gagal start** (fail-closed, `l.Fatal`) |
| `max_employees` di `.lic` | Enforced di satu-satunya jalur create (`Service.Create()`) — tidak ada bypass |

### 4.7 Perpanjangan Lisensi

Vendor membuat `.lic` baru dengan `--expires` lebih lama, klien mengganti file `license.lic` lalu restart server.

---

## 5. Perbandingan Lengkap

| Item | SaaS | On-Premise |
|---|---|---|
| `deployment_mode` | `saas` | `on_premise` |
| Infrastruktur | Vendor (cloud) | Klien |
| Tenant isolation | DB per tenant | Single tenant |
| Lisensi modul | `company_modules` DB + package | File `.lic` RSA |
| Batas employee | Unlimited | `max_employees` (enforced) |
| Provisioning | Otomatis per signup | Sekali saat setup |
| Signing lisensi | Tidak ada | `licensectl gen-lic` (vendor) |
| Gagal verifikasi lisensi | N/A (DB-based) | Fail-closed startup |
| Aktivasi modul baru | Instant (via subscribe/management) | Butuh `.lic` baru + restart |
| Kredensial tenant DB | Per-tenant (AES-256 terenkripsi) | Single DB |
| Cocok untuk | Banyak customer, churn tinggi | Enterprise, kepatuhan data |

---

## 6. Production Checklist (Kedua Mode)

- [ ] `HRIS_JWT_SECRET` diganti dari default
- [ ] `server.mode: "release"`
- [ ] `logger.level: "info"` (atau `warn` di production)
- [ ] `encryption_key` diset (SaaS wajib; on-premise opsional)
- [ ] Redis tidak memakai password kosong
- [ ] Database credentials bukan root/super
- [ ] CORS `allowed_origins` dibatasi (bukan `*`) jika frontend domain tetap
- [ ] Backup: platform DB + semua tenant DB (atau single DB on-premise)
- [ ] Monitoring: `/api/v1/platform/monitoring/*` (health, pool stats, cache stats)

---

## 7. Troubleshooting

| Gejala | Penyebab | Solusi |
|---|---|---|
| `Fatal: On-premise license invalid` | `.lic` expired/tampered/key salah | Regenerate `.lic`; cek tanggal server |
| `Fatal: On-premise public key not found` | Path public key salah | Set `HRIS_LICENSE_PUBLIC_KEY_FILE` benar |
| 403 `MODULE_NOT_LICENSED` | Modul tidak aktif (SaaS) / tidak di `.lic` (on-premise) | SaaS: aktivasi via module management; on-premise: buat `.lic` baru |
| 403 `QUOTA_EXCEEDED` | Jumlah employee = `max_employees` | Perbarui `.lic` dengan `--max-employees` lebih besar |
| `Error 1045 (28000) Access denied` | Kredensial DB salah | Cek `HRIS_DATABASE_*` / `.env` |
| Data region tidak muncul saat provisioning | Seed gagal diam-diam | Jalankan `installer seed-data --company=<id>` (sekarang hard-fail) |
| Redis connection refused | Redis mati / salah host | `docker compose up -d redis`; cek `HRIS_REDIS_HOST` |

---

## 8. Referensi

| Sumber | Path |
|---|---|
| Config template | `backend/config/config.yaml` |
| Env template | `backend/.env.example` |
| Makefile | `backend/Makefile` |
| Docker compose | `docker/docker-compose.yml` |
| On-premise engine | `backend/internal/pkg/onpremise/` |
| Licensectl CLI | `backend/cmd/licensectl/` |
| Installer CLI | `backend/cmd/installer/` |
| OpenAPI Report | `docs/openapi-report.md` |
| Project Dashboard | `docs/project-completion-dashboard.md` |

---

*Dokumen dibuat: 31 Juli 2026 — disinkronkan dengan implementasi mode SaaS (company_modules + package subscribe) dan mode On-Premise (`.lic` RSA + max_employees enforcement).*
