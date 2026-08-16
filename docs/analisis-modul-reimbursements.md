# Analisis Modul Reimbursement

> Dibuat: 2026-08-16
> Sumber: eksplorasi kode backend (`backend/internal/modules/reimbursement`), migrasi SQL, routing frontend, dan dokumentasi terkait di `docs/`.

## 1. Ringkasan

Modul Reimbursement menangani pengajuan klaim penggantian biaya (expense claim) oleh karyawan — mis. transport, medis, dsb. — beserta workflow persetujuannya. Backend (model, service, API, RBAC, integrasi approval) **sudah lengkap dan teruji**, namun **frontend masih berupa placeholder "Coming soon"** dan dua integrasi lintas modul (pembayaran via payroll, notifikasi outcome) **belum diimplementasikan**.

Status singkat:

| Layer | Status |
|---|---|
| Database & migrasi | ✅ Selesai |
| Model / Repository / Service / Handler / Routes | ✅ Selesai, ada unit & integration test |
| Integrasi Approval Engine (central) | ✅ Selesai |
| RBAC permissions | ✅ Selesai |
| Integrasi payroll (pembayaran) | ❌ Belum ada |
| Notifikasi outcome (approved/rejected/paid) | ❌ Belum ada (direncanakan di Notification Phase 5) |
| Upload lampiran bukti (receipt) | ⚠️ Parsial — hanya field URL string, tidak ada endpoint upload khusus |
| Frontend UI | ❌ Placeholder saja, belum ada fitur nyata |

## 2. Lokasi Kode

### Backend — `backend/internal/modules/reimbursement/`

| File | Isi |
|---|---|
| `model.go` | Model GORM: `ReimbursementType`, `ReimbursementRequest`, `ReimbursementItem`, enum `ReimbursementStatus` |
| `dto.go` | DTO request/response untuk ketiga entitas + `PaginatedResponse` |
| `repository.go` | Data access (CRUD + agregasi `SumReimbursementItems`) |
| `service.go` | Business logic: state machine status, integrasi approval engine |
| `handler.go` | HTTP handler (Gin) |
| `routes.go` | Registrasi route |
| `module.go` | Definisi modul: menu, permission, `Migrate`, `Seed` |
| `*_test.go` | Test suite (repository, service, handler, approval integration) — ±60 test |

Terhubung di `backend/cmd/server/main.go`:
- Import baris 65.
- Baris 991–999: service dibuat lebih awal (bukan di dalam `NewModule`) agar handler status-approval bisa didaftarkan sebelum modul dibungkus:
```go
reimbursementResolver := reimbursement.NewTenantDBResolver(dbManager)
reimbursementRepo := reimbursement.NewRepository(reimbursementResolver)
reimbursementSvc := reimbursement.NewService(reimbursementRepo, l.Named("reimbursement"))
reimbursementSvc.SetApprovalEngine(sharedApprovalEngine)
approvalSvc.RegisterStatusHandler("reimbursement", func(ctx context.Context, documentID uuid.UUID, status approval.InstanceStatus, note string) error {
    return reimbursementSvc.HandleApprovalStatusChange(ctx, documentID, string(status), note)
})
```
- Baris 1287: `Module: reimbursement.NewModuleWithService(l, reimbursementSvc)` didaftarkan ke daftar modul.

### Migrasi Database
- `backend/internal/pkg/migrator/migrations/tenant/postgres/013_reimbursement.sql` (+ `.down.sql`) — skema dasar (3 tabel).
- `backend/internal/pkg/migrator/migrations/tenant/postgres/061_reimbursement_approval_instance.sql` (+ `.down.sql`) — menambah kolom `approval_instance_id`.
- Migrasi setara untuk MySQL tersedia di path yang sama dengan `tenant/mysql/...`.

### Frontend — `frontend/tenant/src`
- Route: `router/index.js` baris 436–441 → path `reimbursements`, komponen `@/views/modules/reimbursement/Reimbursements.vue`.
- Sidebar: `layouts/Sidebar.vue` baris 395 & 436–441, dikelompokkan di bawah "Finance — Payroll & Reimbursement", digate oleh `hasModule('reimbursement')` + permission `reimbursement.view`.
- i18n: string sudah ada di `en.json` / `id.json` (`reimbursement.title`, `.description`, `.types`, `.requests`, label nav) meski halaman belum berfungsi.
- **Isi lengkap `Reimbursements.vue` saat ini** (placeholder):
```html
<template>
  <div class="flex items-center justify-center h-full">
    <div class="text-center">
      <i class="pi pi-spinner text-3xl text-emerald-500 mb-3"></i>
      <h2 class="text-lg font-semibold text-gray-700">Reimbursements Module</h2>
      <p class="text-sm text-gray-400 mt-1">Coming soon</p>
    </div>
  </div>
</template>
```
Tidak ada Pinia store, composable, atau file API-call untuk reimbursement di frontend sama sekali. Satu-satunya jalur "nyata" untuk melihat pengajuan reimbursement saat ini adalah lewat inbox Approvals generik (`views/modules/approval/Approvals.vue`), yang menangani reimbursement hanya sebagai salah satu tipe dokumen umum, bukan UI khusus.

## 3. Struktur Tabel

### `reimbursement_types`
| Kolom | Tipe | Keterangan |
|---|---|---|
| id | CHAR(36) PK | |
| code | VARCHAR(50) | default `''` |
| name | VARCHAR(150) NOT NULL | |
| description | VARCHAR(500) | default `''` |
| is_active | SMALLINT | default `1` |
| deleted_at | TIMESTAMP | soft delete |
| created_at / updated_at | TIMESTAMP | |

Index: `idx_reimb_type_active(is_active)`, `idx_reimb_type_deleted_at`.

### `reimbursement_requests`
| Kolom | Tipe | Keterangan |
|---|---|---|
| id | CHAR(36) PK | |
| employee_id | CHAR(36) NOT NULL | referensi logis ke tabel employee |
| request_type_id | CHAR(36) NOT NULL | referensi logis ke `reimbursement_types.id` |
| title | VARCHAR(200) NOT NULL | |
| description | TEXT | |
| total_amount | DECIMAL(18,2) | default `0.00`, dihitung dari `reimbursement_items` |
| currency | VARCHAR(3) | default `'IDR'` |
| status | VARCHAR(50) | default `'DRAFT'` — lihat state machine §4 |
| supervisor_id | CHAR(36) NULL | |
| supervisor_action_at | TIMESTAMP NULL | |
| supervisor_note | VARCHAR(500) NULL | |
| hr_id | CHAR(36) NULL | |
| hr_action_at | TIMESTAMP NULL | |
| hr_note | VARCHAR(500) NULL | |
| paid_at | TIMESTAMP NULL | |
| paid_amount | DECIMAL(18,2) NULL | |
| submitted_at / approved_at / rejected_at / cancelled_at | TIMESTAMP NULL | |
| approval_instance_id | CHAR(36) NULL | ditambahkan migrasi 061 |
| deleted_at, created_at, updated_at | TIMESTAMP | |

Index: `employee_id`, `request_type_id`, `status`, `deleted_at`, `approval_instance_id`.

> Catatan konsistensi: SQL migration mendeklarasikan kolom timestamp aksi (`supervisor_action_at`, `hr_action_at`, `paid_at`, `submitted_at`, `approved_at`, `rejected_at`, `cancelled_at`) sebagai `TIMESTAMP`, tetapi model GORM (`model.go`) mengannotasikannya sebagai `int64`/`bigint` (Unix-nanosecond epoch, `json:"-"`, dikonversi ke `*time.Time` hanya di layer DTO lewat `unixNanoToTimePtr`). ini adalah inkonsistensi tipe antara `model.go` dan file `.sql` — perlu diperhatikan karena modul juga menjalankan `AutoMigrate` GORM di `Migrate()`, sehingga tipe kolom efektif tergantung path mana yang dieksekusi terakhir.

Tidak ada FK constraint eksplisit di SQL (konsisten dengan konvensi codebase ini: referensi antar-tabel dijaga di level aplikasi, bukan DB constraint).

### `reimbursement_items`
| Kolom | Tipe | Keterangan |
|---|---|---|
| id | CHAR(36) PK | |
| reimbursement_request_id | CHAR(36) NOT NULL | referensi ke `reimbursement_requests.id` |
| expense_date | DATE NOT NULL | |
| expense_type | VARCHAR(100) NOT NULL | |
| description | VARCHAR(500) NULL | |
| amount | DECIMAL(18,2) NOT NULL | |
| receipt_url | TEXT NULL | URL bukti/struk (lihat §5 soal upload) |
| deleted_at, created_at, updated_at | TIMESTAMP | |

Index: `reimbursement_request_id`, `deleted_at`.

### Status Enum
`DRAFT`, `SUBMITTED`, `PENDING_APPROVAL`, `APPROVED`, `REJECTED`, `PAID`, `CANCELLED` (`model.go` baris 42–50).

## 4. API Endpoints

Semua route berada di bawah `/reimbursements` dalam grup API tenant (`routes.go`):

| Method | Path | Fungsi |
|---|---|---|
| POST | `/reimbursements/types` | Buat tipe klaim (mis. Transport, Medical) |
| GET | `/reimbursements/types` | List tipe (paginated) |
| GET | `/reimbursements/types/:id` | Detail tipe |
| PUT | `/reimbursements/types/:id` | Update tipe |
| DELETE | `/reimbursements/types/:id` | Soft-delete tipe |
| POST | `/reimbursements/requests` | Buat pengajuan (draft) — karyawan mengajukan untuk diri sendiri via `user_id` di context |
| GET | `/reimbursements/requests` | List/filter berdasarkan `employee_id`, `status`, paginated |
| GET | `/reimbursements/requests/:id` | Detail pengajuan (preload items) |
| PUT | `/reimbursements/requests/:id` | Update field (title/description/currency/supervisor) — hanya saat `DRAFT` |
| PUT | `/reimbursements/requests/:id/status` | Transisi status (submit/approve/reject/pay/cancel) |
| DELETE | `/reimbursements/requests/:id` | Soft-delete pengajuan |
| GET | `/reimbursements/requests/:id/items` | List item biaya |
| POST | `/reimbursements/requests/:id/items` | Tambah item biaya (hanya saat `DRAFT`) |
| PUT | `/reimbursements/requests/:id/items/:itemId` | Edit item (hanya saat `DRAFT`) |
| DELETE | `/reimbursements/requests/:id/items/:itemId` | Hapus item (hanya saat `DRAFT`) |

**Permissions (RBAC)**: `reimbursement.view`, `reimbursement.create`, `reimbursement.update`, `reimbursement.delete`, `reimbursement.approve` — sudah di-seed sebagai bagian dari default RBAC set. `DependsOn: ["employee"]`.

## 5. Alur Bisnis (Flow)

### 5.1 State Machine (`service.go`, fungsi `UpdateReimbursementRequestStatus`)

```
DRAFT ──submit──► SUBMITTED / PENDING_APPROVAL ──approve──► APPROVED ──pay──► PAID
  │                        │
  └──────cancel────────────┴──reject──► REJECTED
                            (cancel juga bisa dari state manapun)
```

- **`DRAFT → SUBMITTED`**: hanya valid dari `DRAFT`. Sistem otomatis menghitung `total_amount` dari penjumlahan `reimbursement_items`, dan mencatat `submitted_at`.
  - **Routing approval**: jika `ApprovalEngine` terpasang (selalu terpasang via `sharedApprovalEngine`) dan `flow_id` disertakan di body request, service memanggil `approvalEngine.CreateApprovalInstance(ctx, "reimbursement", requestID, flowID)`. Jika berhasil, `approval_instance_id` disimpan dan status menjadi `PENDING_APPROVAL`. Jika gagal karena routing/assignee tidak ditemukan (`approval.RoutingError`), submit gagal total (HTTP 400 bilingual). Error approval lain hanya di-log sebagai warning dan status tetap jadi `SUBMITTED` biasa (fallback best-effort tanpa approval).
  - Jika tidak ada engine/flow, status langsung `SUBMITTED` (jalur lama ad-hoc supervisor/HR — field `supervisor_id`/`hr_id` tersedia di tabel tapi tidak ada endpoint approval terpisah selain endpoint status generik ini).
- **`SUBMITTED`/`PENDING_APPROVAL → APPROVED`**: mencatat `approved_at`.
- **`* → REJECTED`**: mencatat `rejected_at` (tidak ada guard status sumber).
- **`* → CANCELLED`**: mencatat `cancelled_at` (tidak ada guard status sumber).
- **`APPROVED → PAID`**: hanya valid dari `APPROVED`. Mencatat `paid_at`; `paid_amount` = nilai yang dikirim atau default ke `total_amount`.
  - **Penting**: status `PAID` ini **bukan integrasi otomatis dengan payroll**. Tidak ada referensi ke modul `reimbursement` di dalam `backend/internal/modules/payroll`, dan tidak ada payslip/disbursement line item yang dibuat otomatis. Menandai `PAID` murni flag manual lewat endpoint status yang sama (presumably dilakukan oleh HR/Finance secara manual setelah transfer dilakukan di luar sistem).

### 5.2 Callback Approval (push-based)

`HandleApprovalStatusChange` didaftarkan sebagai status handler modul `"reimbursement"` ke Approval Engine terpusat. Saat modul approval me-resolve instance ke `APPROVED`/`REJECTED`/`CANCELLED` (misalnya dari UI Approvals generik), callback ini meng-update `ReimbursementRequest` — hanya jika status saat ini masih `PENDING_APPROVAL` (idempotent/no-op di luar itu). Inilah mekanisme yang membuat aksi approve/reject dari halaman Approvals generik (`views/modules/approval/Approvals.vue`) tercermin balik ke record reimbursement.

### 5.3 Lampiran Bukti (Receipt)

`ReimbursementItem.ReceiptURL` hanya menyimpan string URL per item biaya. Tidak ditemukan endpoint upload file khusus di dalam modul ini — URL diasumsikan dihasilkan oleh layanan upload file bersama/generik di luar modul, lalu dikirim saat create/update item.

### 5.4 Notifikasi

Tidak ditemukan integrasi notifikasi langsung di modul reimbursement. `docs/module-notification-plan.md` secara eksplisit mencantumkan reimbursement sebagai bagian dari **Phase 5 (belum dikerjakan)**: notifikasi outcome (approved/rejected/paid) untuk payroll/employee-movement/reimbursement + recruitment masih pending.

### 5.5 Guardrail Lain
- Item biaya hanya bisa ditambah/edit/hapus selama status request masih `DRAFT`.
- Field request (title/description/currency/supervisor) hanya bisa diedit selama `DRAFT`.

## 6. Status Implementasi Detail

| Aspek | Status | Catatan |
|---|---|---|
| Skema DB & migrasi (Postgres + MySQL) | ✅ | 3 tabel, migrasi 013 & 061 |
| Model/Repository/Service/Handler/Routes | ✅ | Ada test lengkap (~60 test: repo 18, service 20, handler 16) |
| Integrasi Approval Engine terpusat | ✅ | Pola sama dengan payroll/leave/KPI |
| RBAC permission | ✅ | 5 permission ter-seed |
| Integrasi pembayaran via payroll | ❌ | `PAID` hanya flag manual, tidak ada linkage ke payslip/payroll run |
| Notifikasi status (approved/rejected/paid) | ❌ | Ditunda ke Notification Phase 5 |
| Upload lampiran bukti | ⚠️ | Hanya field URL string, tanpa endpoint upload khusus di modul ini |
| Frontend UI (list/create/edit/detail/approve/report) | ❌ | Placeholder "Coming soon"; route, sidebar, definisi menu (termasuk rencana sub-menu Requests/Types/Reports di `module.go`), dan string i18n sudah ada, tapi belum ada UI fungsional, store, atau composable |
| OpenAPI docs | ✅ | Terdokumentasi di `backend/internal/pkg/docs/openapi.json` |

**Catatan penting**: `docs/project-completion-dashboard.md` menandai "Reimbursement & Claim" sebagai "✅ Complete" — namun itu metrik **backend saja** (3 tabel, 60 test, 15 endpoint). Frontend belum diimplementasikan sama sekali, dan dua integrasi lintas modul (payroll payout, notifikasi) masih tertunda.

Tidak ditemukan file `docs/module-reimbursement-plan.md` tersendiri (berbeda dengan modul leave, movement, notification, career-intelligence, recruitment yang masing-masing punya `module-*-plan.md`) — mengindikasikan reimbursement belum punya dokumen perencanaan build-out frontend.

## 7. Rekomendasi Langkah Berikutnya

1. **Bangun frontend Reimbursement** — minimal: list pengajuan (dengan filter status/employee), form create/edit draft + item biaya, detail view, tombol submit/approve/reject/pay sesuai permission, dan halaman kelola `reimbursement_types`. Perlu store/composable baru karena belum ada sama sekali.
2. **Integrasi payroll**: tentukan apakah reimbursement `PAID` perlu otomatis masuk sebagai komponen payslip/off-cycle payment, atau tetap manual by design — perlu keputusan produk.
3. **Notifikasi outcome**: implementasikan sesuai rencana Phase 5 di `module-notification-plan.md` (submitted, approved, rejected, paid).
4. **Selaraskan tipe kolom timestamp** antara `model.go` (int64 unix-nano) dan migrasi SQL (`TIMESTAMP`) untuk menghindari ambiguitas AutoMigrate vs SQL migration.
5. **Upload bukti**: pastikan/pastikan ada mekanisme upload file bersama yang dipakai frontend nanti untuk mengisi `receipt_url`, atau tambahkan endpoint upload khusus jika belum ada di modul file/shared.
