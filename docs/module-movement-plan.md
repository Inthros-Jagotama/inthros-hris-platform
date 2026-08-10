# Employee Movement & Career Management — Development Plan

> 📅 Versi plan: 2026-08-10 · Status: **IMPLEMENTASI BERJALAN — langkah 13/13 selesai** + **enhancement §12 P0 (1–5) ✅ + P1 (6–8) ✅** (backend ✅ + FE halaman Movements/Contracts/detail ✅ + checklist E2E dibuat ✅ — eksekusi manual menunggu environment tenant)
> ✅ **Keputusan bisnis sudah dikonfirmasi user (2026-08-10)** — lihat §11.
> 🔎 Berdasarkan struktur tabel `012_employee_movement.sql` (mysql + postgres) dan `062_employeemovement_approval_instance.sql`, serta audit modul `backend/internal/modules/employeemovement` dan `frontend/tenant/src/views/modules/EmployeeMovements.vue`.
> 📊 **Progres implementasi (per 2026-08-10):** ✅ 1) migration + enum `rejected` (082) · ✅ 2) G-1 ExecuteMovement transaksi employment · ✅ 3) G-3 auto-resolve flow · ✅ 4) G-4 enriched responses · ✅ 5) G-2 notifikasi `MOVEMENT_*` (termasuk deep-link FE) · ✅ 6) G-5 hapus approve manual · ✅ 7) G-6 contract extension count · ✅ 8) G-7 validasi per tipe · ✅ 9) G-8 slug/route disamakan · ✅ 10) FE halaman Movements + filter backend (termasuk badge `rejected`) · ✅ 11) FE halaman Contracts (daftar enriched + filter + create/edit + upload dokumen) + filter backend · ✅ 12) FE detail dialog movement + aksi dari detail · ✅ 13) checklist verifikasi E2E manual dibuat (`docs/module-movement-e2e-checklist.md`; unit/service & build sudah PASS di tiap langkah) · ✅ **Enhancement §12 P0 1–3** (transactional ExecuteMovement + position conflict + effective-date conflict — lihat log §3.17) · ✅ **Enhancement §12 P0 4** (Movement Snapshot — lihat log §3.18) · ✅ **Enhancement §12 P0 5** (Movement Audit Trail — lihat log §3.19) · ✅ **Enhancement §12 P1 6** (Movement Documents — lihat log §3.20) · ✅ **Enhancement §12 P1 7** (Career Timeline — lihat log §3.21) · ✅ **Enhancement §12 P1 8** (Contract Expiry Management — lihat log §3.22). **Eksekusi E2E manual:** menunggu environment tenant + akun ber-permission `employeemovement.*`.

---

# 1. Objective

Membangun modul **Employee Movement & Career Management**: riwayat pergerakan karyawan (promosi, demosi, mutasi/rotasi, perpanjangan kontrak, perubahan status, pensiun, offboarding) dan pengelolaan kontrak kerja (PKWT/PKWTT/harian).

- **Backend**: sudah diimplementasikan penuh (CRUD + alur approval + eksekusi) — plan ini fokus pada **gap fungsional** yang belum ada.
- **Frontend**: masih placeholder ("Coming soon") — plan ini mencakup pembangunan halaman penuh.

---

# 2. Existing Database Structure

Sumber: `backend/internal/pkg/migrator/migrations/tenant/mysql/012_employee_movement.sql` (+ postgres) dan `062_employeemovement_approval_instance.sql` (mysql + postgres).

## 2.1 Tabel `employee_movements` (riwayat pergerakan)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `employee_id` | CHAR(36) NN → `employees(id)` CASCADE | karyawan |
| `movement_type` | VARCHAR(50) NN | `promotion, demotion, mutation, contract_extension, status_change, retirement, offboarding, other` |
| `from_employment_id` / `to_employment_id` | CHAR(36) NULL → `employments(id)` SET NULL | employment sebelum/sesudah |
| `from_organization_id` / `to_organization_id` | CHAR(36) NULL → `organizations(id)` SET NULL | |
| `from_position_id` / `to_position_id` | CHAR(36) NULL → `positions(id)` SET NULL | |
| `from_employment_status_id` / `to_employment_status_id` | CHAR(36) NULL → `employment_statuses(id)` SET NULL | |
| `reason`, `notes` | TEXT NULL | |
| `decision_letter_number` | VARCHAR(50) NN | nomor SK |
| `decision_letter_date` | DATE NN | tanggal SK |
| `effective_date` | DATE NN | tanggal berlaku |
| `status` | VARCHAR(20) DEFAULT `draft` | **target enum final**: `draft, pending_approval, approved, rejected, executed, cancelled` (perlu migration tambah `rejected` — keputusan §11.4) |
| `approved_by` / `approved_at` | CHAR(36) / TIMESTAMP NULL | |
| `executed_by` / `executed_at` | CHAR(36) / TIMESTAMP NULL | |
| `approval_instance_id` | CHAR(36) NULL (migration 062) | instance Central Approval |
| `created_by` / `updated_by` / timestamps | | audit |

Index: employee, type, status, effective_date, from/to org, from/to position, approval_instance (062).

## 2.2 Tabel `employee_contracts` (PKWT & perjanjian kerja)

| Kolom | Tipe | Keterangan |
|---|---|---|
| `id` | CHAR(36) PK | |
| `employee_id` | CHAR(36) NN → `employees(id)` CASCADE | |
| `contract_number` | VARCHAR(50) NN | |
| `contract_type` | VARCHAR(20) NN | `pkwt, pkwtt, daily, other` |
| `start_date` | DATE NN | |
| `end_date` | DATE NULL | |
| `extension_count` | INT DEFAULT 0 | jumlah perpanjangan |
| `previous_contract_id` | CHAR(36) NULL → `employee_contracts(id)` SET NULL | rantai perpanjangan |
| `decision_letter_number` | VARCHAR(50) NULL | |
| `notes` | TEXT NULL | |
| `document_url` | VARCHAR(255) NULL | lampiran (dapat memakai upload generik `POST /uploads`) |
| `status` | VARCHAR(20) DEFAULT `active` | `active, expired, extended, terminated` |
| audit (`created_by`/`updated_by`/timestamps) | | |

Index: employee, status, type, end_date.

---

# 3. Status Aktual

## 3.1 Backend — ✅ SUDAH IMPLEMENTASI (per 2026-08-10)

Modul `backend/internal/modules/employeemovement/` (±3.200 baris) sudah lengkap:

- **Model**: `EmployeeMovement` + `EmployeeContract` (enum movement/contract type & status, `TableName`, `BeforeCreate`) — migration 012 + 062.
- **DTO**: `Create/UpdateMovementRequest`, `MovementResponse`, `SubmitMovementRequest`, `Create/UpdateContractRequest`, `ContractResponse`, pagination.
- **Service**: CRUD movement, `SubmitMovement` (routing ke Central Approval), `HandleApprovalStatusChange` (push-callback), `Approve/Execute/CancelMovement`, CRUD kontrak + alur `ExtendContract`.
- **Repository**: CRUD + `ExecuteMovement`/`CancelMovement`/`ExtendContract` (`ApproveMovement` sudah dihapus sejak G-5, 2026-08-10).
- **Handler + Routes**: 
  - `POST/GET /employee-movements/movements`, `GET/PUT/DELETE /movements/:id`, `POST /movements/:id/submit|execute|cancel` (approve manual **dihapus** §11.5), `GET /employees/:employeeId/movements`
  - `POST/GET /employee-movements/contracts`, `GET/PUT/DELETE /contracts/:id`, `GET /employees/:employeeId/contracts`
- **Approval**: interface `ApprovalEngine` (narrow, pola payroll/leave/reimbursement), di-wire di `cmd/server/main.go` (`SetApprovalEngine` + `RegisterStatusHandler("employeemovement", ...)`).
- **Menu (server)**: "Career" → Movements & Contracts (routes `/admin/career/*`).
- **Test**: `handler_test.go`, `service_test.go`, `repository_test.go`, `approval_integration_test.go`.

## 3.2 Frontend — ❌ BELUM (placeholder)

- `frontend/tenant/src/views/modules/EmployeeMovements.vue` = **1 baris placeholder** ("Coming soon"), tidak ada komponen/API call.
- Router sudah terdaftar: `/employee-movements` (meta module `employee-movement`).
- Locale hanya 4 key: `employee_movement.title/description/movements/contracts`.
- `Approvals.vue` sudah punya deep-link `case 'employeemovement'` → `GET /employee-movements/movements/:id` (siap dipakai).

## 3.3 Selisih migration vs model

| Item | Migration 012 | Model Go |
|---|---|---|
| `status` enum | `draft, approved, executed, cancelled` | `draft, pending_approval, approved, executed, cancelled` (model lebih kaya) |
| `status` target | — | **+ `rejected`** (keputusan §11.4 → migration + enum) |
| `approval_instance_id` | ditambah migration 062 | ada |
| `extension_count` | ada | ada |

Status model lebih kaya dari komentar migration (wajar karena migration 012 tidak diubah setelah 062). Gap aktual: **`rejected` belum ada di enum model maupun migration** — perlu ditambahkan.

## 3.4 Log Implementasi — Langkah 1 ✅ (migration 082 + enum `rejected`)

> Implementasi selesai **2026-08-10**. Referensi untuk langkah-langkah berikut.

### 3.4.1 Fakta desain

Kolom `status` di `employee_movements` sudah `VARCHAR(20)` sejak migration 012 di **kedua dialect** (bukan ENUM) → nilai `rejected` (8 karakter) tertampung **tanpa perubahan skema**. Migration 082 dibuat sebagai pasangan sinkron nomor + verifikasi anti-drift (pola yang sama dengan migration 079 untuk overtime `PENDING_APPROVAL`).

### 3.4.2 File yang dibuat

```
backend/internal/pkg/migrator/migrations/tenant/mysql/082_employeemovement_status_rejected.sql        (up)
backend/internal/pkg/migrator/migrations/tenant/mysql/082_employeemovement_status_rejected.down.sql   (down)
backend/internal/pkg/migrator/migrations/tenant/postgres/082_employeemovement_status_rejected.sql      (up)
backend/internal/pkg/migrator/migrations/tenant/postgres/082_employeemovement_status_rejected.down.sql (down)
```

**MySQL up** — idempotent via `information_schema`:
- Jika kolom masih `VARCHAR` → no-op (`DO 0`) — kondisi normal.
- Jika kolom sudah drift jadi `ENUM` → `MODIFY COLUMN status ENUM('draft','pending_approval','approved','rejected','executed','cancelled') DEFAULT 'draft'` (nilai `rejected` disisipkan otomatis).

**MySQL down** — guard **anti data-truncation**:
- Konversi ke ENUM (tanpa `rejected`) hanya dijalankan jika (a) kolom bukan VARCHAR **dan** (b) **tidak ada baris berstatus `rejected`** — jika ada baris tsb, rollback di-skip dan operator harus menormalkan/menghapus datanya dulu (komentar peringatan ada di file).

**Postgres up/down** — no-op `SELECT 1` (VARCHAR sudah menampung `rejected`).

### 3.4.3 Perubahan kode

- `backend/internal/modules/employeemovement/model.go` — konstanta baru:
  ```go
  MovementStatusRejected MovementStatus = "rejected" // keputusan plan §11.4: status terpisah utk approval ditolak
  ```
- Belum ada perubahan di service/handler/DTO — pemakaian `rejected` menyusul di langkah 5 (G-2 notifikasi + `HandleApprovalStatusChange` REJECTED → `rejected`) dan langkah 10 (badge FE).

### 3.4.4 Validasi

- `go build ./...` ✅ · `go test ./internal/modules/employeemovement/` ✅
- Migration 082 applied di tenant MySQL (`df687f34-…`): tercatat di `schema_migrations`, kolom tetap `varchar(20)` ✅
- `TestMigratorIntegration/MySQL` **PASS** (mencakup up + down 082 di database bersih) ✅
- Catatan: Postgres lokal tidak aktif saat validasi — no-op aman, diverifikasi saat server Postgres tersedia.

## 3.5 Log Implementasi — Langkah 2 ✅ (G-1 ExecuteMovement transaksi employment)

> Implementasi selesai **2026-08-10**. Referensi untuk langkah-langkah berikut.

### 3.5.1 Interface narrow + adapter

- `employeemovement.CareerExecutor` (interface di `service.go`): `FindCurrentEmployment`, `CloseEmployment`, `CreateEmployment`, `SetEmployeeInactive` + data holder `CareerEmployment` (berisi `ID` yang hanya terisi saat `FindCurrentEmployment`, serta `to_*` fields) — **tanpa import modul employee** (hanya `uuid`).
- `employeeCareerAdapter` di `cmd/server/main.go` membungkus **instance `employee.Repository` terpisah** (pola sama dengan `setting.Service` untuk holiday leave) dan di-wire via `employeeMovementSvc.SetCareerExecutor(...)` sebelum module di-mount.

### 3.5.2 Perubahan kode

| File | Perubahan |
|---|---|
| `employeemovement/service.go` | Interface `CareerExecutor` + `SetCareerExecutor`; `ExecuteMovement` = fetch movement → validasi `approved` → HR data change per tipe → `repo.ExecuteMovement` (status `executed`) |
| `employeemovement/repository.go` | `ExecuteMovement(ctx, id, executedBy, toEmploymentID *uuid.UUID)` — `to_employment_id` di-persist hanya jika tidak nil |
| `employee/repository.go` | +`FindActiveEmploymentByEmployeeID` (end_date NULL, terbaru), +`CloseEmployment` (guard `effective_end_date IS NULL`), +`SetEmployeeStatus` (update status saja, hindari Preload berat) |
| `cmd/server/main.go` | +`employeeCareerAdapter` (4 method) + wiring `SetCareerExecutor` |
| `service_test.go`, `repository_test.go`, `helpers_test.go` | fake executor + 2 integration test (promotion & offboarding, repo employee nyata + adapter, SQLite) + `dayBefore` test + helper `ptrUUID`/`uuidPtr`/`dateStartsWith`/`testLogger` |

### 3.5.3 Perilaku per tipe movement

| Tipe | Perubahan HR data saat execute |
|---|---|
| `promotion`, `demotion`, `mutation`, `status_change`, `other` | Tutup employment aktif (`end = effective_date - 1`) → buat employment baru (`to_*` + SK + effective_date) → `to_employment_id` di-persist di movement |
| `offboarding`, `retirement` | Tutup employment aktif, **tanpa** employment baru; employee di-set `status = inactive` |
| `contract_extension` | Tanpa perubahan employment (hanya status movement) |

`effective_date` boleh masa depan (§11.2): employment baru disimpan dengan tanggal tsb, yang lama aktif sampai sehari sebelumnya. `dayBefore` menangani format plain date (`YYYY-MM-DD`) maupun RFC3339 (driver MySQL/SQLite).

### 3.5.4 Catatan & tradeoff

- **Non-transaksional lintas modul** (pola leave→attendance): jika langkah tengah gagal, movement tetap `approved` (retry oleh HR). Konsekuensi & mitigasi didokumentasikan di catatan §G-1.
- `CloseEmployment` tidak akan menimpa employment yang sudah tertutup (guard `effective_end_date IS NULL`).
- `other` diperlakukan sebagai perubahan kepegawaian (buat employment baru) — keputusan yang bisa ditinjau ulang.

### 3.5.5 Validasi

- `go build ./...` ✅ · `go vet` (cmd/server, employeemovement, employee) ✅
- `go test` employeemovement / employee / approval — semua **PASS** ✅
- Integration test (repo nyata + adapter, SQLite): promotion → employment lama `effective_end_date = 2026-07-31`, employment baru terisi `to_*` fields, `to_employment_id` persist di movement; offboarding → employment tertutup tanpa yang baru + employee `inactive` ✅
- E2E via API tidak dieksekusi penuh (akun admin tenant untuk permission `employeemovement.*` belum tersedia saat validasi) — alur diverifikasi via integration test setara.

## 3.6 Log Implementasi — Langkah 3 ✅ (G-3 auto-resolve flow)

> Implementasi selesai **2026-08-10**. Referensi untuk langkah-langkah berikut.

### 3.6.1 Perubahan

- `employeemovement/service.go`:
  - Interface `ApprovalEngine` + method **`GetActiveFlowIDForModule(ctx, module string) (string, error)`** (pola sama leave/attendance).
  - `SubmitMovement` di-refactor: bila client tidak mengirim `flow_id`, flow aktif module `"employeemovement"` di-auto-resolve. Prioritas: `flow_id` eksplisit → auto-resolve → error bila keduanya kosong (`approval flow not configured: provide flow_id or activate an approval flow for module employeemovement`).
- `approval_integration_test.go`: fake engine + `GetActiveFlowIDForModule` (resolvedFlowID/resolvedFlowErr/resolveCalls); test lama `NoFlowID_ReturnsError` diganti 3 test baru: **NoFlowID_AutoResolves** (flow ter-resolve dipakai), **NoFlowResolved_ReturnsError** (flow kosong → error), **ExplicitFlowID_BeatsAutoResolve** (flow_id eksplisit menang, auto-resolve tidak dipanggil).
- `cmd/server/main.go`: **tanpa perubahan** — `sharedApprovalEngine` (payrollApprovalAdapter) sudah mengimplementasi `GetActiveFlowIDForModule` → otomatis memenuhi interface baru.

### 3.6.2 Validasi

- `go build ./...` ✅ · `go test ./internal/modules/employeemovement/` ✅

## 3.7 Log Implementasi — Langkah 6 ✅ (G-5 hapus approve manual)

> Implementasi selesai **2026-08-10**. Keputusan §11.5: satu pintu approval via `submit` → Central Approval.

### 3.7.1 Yang dihapus

| File | Penghapusan |
|---|---|
| `routes.go` | Route `POST /movements/:id/approve` dihapus (+komentar: approval hanya lewat submit) |
| `handler.go` | Handler `ApproveMovement` dihapus |
| `service.go` | Service `ApproveMovement` dihapus; komentar interface `ApprovalEngine` & `SubmitMovement` diperbarui (bukan lagi "manual fallback") |
| `repository.go` | Repository `ApproveMovement` dihapus |
| `dto.go` | Komentar `SubmitMovementRequest` diperbarui |
| `module.go` | Permission `employeemovement.approve` dihapus dari daftar (view/create/update/delete/execute) |
| `openapi.json` | Endpoint `/api/v1/tenant/employee-movements/movements/{id}/approve` dihapus (JSON tetap valid) |
| `handler_test.go`, `service_test.go`, `repository_test.go` | Test approve manual dihapus (handler/service/repo) |

### 3.7.2 Verifikasi

- Tidak ada referensi `ApproveMovement` / `/approve` tersisa di modul & seluruh backend ✅
- Seed RBAC tenant (`tenantseed/seed_rbac.go`) memang tidak pernah memuat `employeemovement.approve` → konsisten ✅
- FE belum punya halaman movement → tidak ada referensi frontend yang perlu dihapus ✅
- `go build ./...` ✅ · `go test` employeemovement + approval **PASS** ✅

## 3.8 Log Implementasi — Langkah 5 ✅ (G-2 notifikasi `MOVEMENT_*` + REJECTED → `rejected`)

> Implementasi selesai **2026-08-10**. Pola sama dengan notifikasi attendance/leave.

### 3.8.1 Perubahan

| File | Perubahan |
|---|---|
| `notification/i18n.go` | +4 tipe bilingual: `MOVEMENT_SUBMITTED`, `MOVEMENT_APPROVED`, `MOVEMENT_REJECTED`, `MOVEMENT_EXECUTED` (EN/ID) |
| `employeemovement/service.go` | +interface `Notifier` + `SetNotifier` + helper `notifyMovementOutcome` (best-effort, resolve user via `FindUserIDByEmployeeID`); `HandleApprovalStatusChange`: REJECTED → **`rejected`** (keputusan §11.4, bukan lagi `cancelled`) + notify APPROVED/REJECTED; `ExecuteMovement` sukses → notify `MOVEMENT_EXECUTED` |
| `employeemovement/repository.go` | +`FindUserIDByEmployeeID` (query `employee_accounts`, pola attendance) |
| `cmd/server/main.go` | `employeeMovementSvc.SetNotifier(notificationSvc)` (satu baris) |
| `frontend/.../Notifications.vue` | Deep-link: `reference_type == 'employeemovement'` atau `type` berawalan `MOVEMENT_` → `/admin/career/movements` (pola sama overtime/leave) |
| Test | `fakeNotifier`; `HandleApprovalStatusChange_Rejected` → expect `rejected`; +`Approved_NotifiesEmployee` (seed `employee_accounts`), +`ExecuteMovement_NotifiesEmployee` |

### 3.8.2 Catatan

- `MOVEMENT_SUBMITTED` → approver sudah otomatis dapat task-assigned push dari Central Approval (tidak perlu kirim manual dari modul) — konsisten dengan pola modul lain.
- Notifikasi best-effort: employee tanpa user account → di-skip (tidak menggagalkan approval/eksekusi).

### 3.8.3 Validasi

- `go build ./...` ✅ · `go test` employeemovement + notification **PASS** ✅
- `npm run build` (tenant FE) ✅

---

# 4. Gap Analysis & Backend Enhancement Plan

> ⚠️ Prioritas diurutkan berdasarkan dampak bisnis.

## G-1 🔴 EKSEKUSI TIDAK MENYENTUH EMPLOYMENT (gap terbesar)

`ExecuteMovement` (repository.go:238) **hanya mengubah `status` movement → `executed`** — tidak membuat employment baru, tidak menutup employment lama, dan tidak meng-update organisasi/posisi karyawan. Padahal inti bisnis promosi/mutasi/status-change adalah **perubahan data kepegawaian riil** (`employments`, lihat `employee.Employment` model: `organization_id`, `position_id`, `employment_status_id`, `effective_date`, `effective_end_date`).

**Rencana (keputusan user §11.1-11.3 sudah masuk):**
1. `ExecuteMovement` menjadi transaksi: buat `to_employment_id` (employment baru dengan nilai `to_*` + `effective_date`), tutup employment aktif lama (`effective_end_date = effective_date - 1`), simpan `to_employment_id` di movement.
2. Perlu akses repository `employments` — pola: interface narrow + adapter di `main.go` (seperti `AttendanceSessionUpdater` leave, `ApprovalEngine`), **bukan** import langsung modul employee.
3. **`effective_date` boleh di masa depan** (keputusan §11.2): employment baru disimpan dengan `effective_date` tersebut, employment lama tetap aktif sampai `effective_date - 1`.
4. Tipe khusus: `contract_extension` → hanya update kontrak (lihat G-6); `offboarding`/`retirement` → tutup employment aktif **dan tandai employee non-aktif** (`is_active = false` — keputusan §11.3).
5. Eksekusi tetap **manual oleh HR** via tombol Execute (keputusan §11.1); scheduler otomatis = fase berikutnya (opsional).

> ⚠️ **Catatan implementasi (2026-08-10)**: `ExecuteMovement` memakai pola adapter **non-transaksional lintas modul** (sama dengan leave→attendance) — `CreateEmployment`/`CloseEmployment`/`SetEmployeeInactive` berjalan di koneksi terpisah dari update status movement. Jika langkah tengah gagal, movement tetap `approved` sehingga HR bisa retry; konsekuensi: bila gagal setelah employment baru dibuat, ada employment baru tanpa movement `executed` (retry bisa membuat employment ganda — mitigasi: operator cek sebelum retry). `CloseEmployment` dilengkapi guard `effective_end_date IS NULL` agar tidak menimpa employment yang sudah tertutup. Tipe `other` diperlakukan sebagai perubahan kepegawaian (buat employment baru) — keputusan ini bisa ditinjau ulang di §11.

## G-2 🟠 TIDAK ADA NOTIFIKASI

`internal/modules/notification/i18n.go` tidak punya tipe `MOVEMENT_*` sama sekali (bandingkan leave/overtime/performance yang punya katalog). HR/approver/employee tidak mendapat notifikasi saat submit/approve/reject/execute.

**Rencana:** tambah tipe notifikasi (pola i18n + `Notifier`):
- `MOVEMENT_SUBMITTED` → approver (task-assigned otomatis dari engine pusat, seperti `OVERTIME_ACTUAL_SUBMITTED`)
- `MOVEMENT_APPROVED` / `MOVEMENT_REJECTED` → pengaju/employee (status `rejected` baru — keputusan §11.4)
- `MOVEMENT_EXECUTED` → employee
- `MOVEMENT_CANCELLED` → pengaju (bila perlu)

Deep-link FE: `reference_type = 'employeemovement'` → `/admin/career/movements` atau `/admin/career/contracts` (keputusan §11.6; pola `Notifications.vue`).

## G-3 🟠 SUBMIT BUTUH `flow_id` MANUAL (tanpa auto-resolve)

`SubmitMovement` mengharuskan client mengirim `flow_id`; tanpa itu → error "approval engine or flow_id not configured". Modul attendance/leave sudah punya pola **auto-resolve active flow** (`GetActiveFlowIDForModule`).

**Rencana:** perluas interface `ApprovalEngine` + `GetActiveFlowIDForModule(ctx, module)` (ikuti persis implementasi attendance `CreateOvertimeRequest`), lalu `SubmitMovement` memakai flow aktif module `"employeemovement"` bila `flow_id` tidak dikirim. FE tidak perlu dropdown flow.

## G-4 🟡 RESPONS TIDAK ENRICHED (hanya UUID)

`MovementResponse`/`ContractResponse` hanya membawa UUID (`employee_id`, `to_organization_id`, `to_position_id`, dll.) tanpa nama — FE harus resolve satu-satu (banyak request). Pola yang sudah ada: `ListOvertimeRequests` mengisi `employee_name`/`organization_name` via JOIN.

**Rencana:** tambah field nama di respons list/detail (batch JOIN):
- `employee_name`, `employee_code`
- `from/to_organization_name`, `from/to_position_name`, `from/to_employment_status_name`
- `contract` → `employee_name`, `employee_code`, `previous_contract_number`

### 3.9 Log Implementasi Langkah 4 (G-4 enriched responses) — 2026-08-10

**Tujuan:** respons list/detail movement & contract membawa nama display (bukan hanya UUID) agar FE tidak resolve satu-satu.

**Perubahan:**
- `dto.go` — `MovementResponse` + `EmployeeName`, `EmployeeCode`, `From/ToOrganizationName`, `From/ToPositionName`, `From/ToEmploymentStatusName`; `ContractResponse` + `EmployeeName`, `EmployeeCode`, `PreviousContractNumber`.
- `repository.go` — method batch resolver (raw table query, tanpa import modul employee/organization — pola `attendance.GetEmployeeInfoByIDs`): `GetEmployeeInfoByIDs` (employees → name + employee_id/employee_code), `resolveNamesByIDs` (helper generik), `GetOrganizationNamesByIDs` (nomenclature), `GetPositionNamesByIDs` (title), `GetEmploymentStatusNamesByIDs` (name), `GetContractNumbersByIDs` (employee_contracts → contract_number).
- `service.go` — helper `enrichMovementResponses` / `enrichContractResponses` (batch collect id per tabel → satu query per tabel → map); dipanggil di semua jalur respons: `Get/List/ListByEmployee/Update/Create/Submit` movement & `Get/List/ListByEmployee/Update/Create` contract. Enrichment best-effort (gagal resolve → warn log, tidak error).
- `enrichment_test.go` (baru) — 4 test: list & get movement enriched (employee/org/posisi/status), contract enriched (employee + previous_contract_number), dan no-refs → nama kosong tanpa error.

**Validasi:** `go build` ✅ · `go test ./internal/modules/employeemovement/` **PASS** ✅ (termasuk 4 test enrichment baru).

## G-5 🟡 ENDPOINT APPROVE MANUAL MASIH ADA

`POST /movements/:id/approve` (manual, tanpa engine) tetap eksis sebagai jalur paralel approval. Dengan Central Approval yang sudah jalan, jalur manual berisiko "dua pintu" yang tidak sinkron.

**Keputusan user §11.5: HAPUS.** `POST /movements/:id/approve` dihapus dari handler/routes; satu-satunya jalur approval = `submit` → Central Approval. (Service method `ApproveMovement` boleh dihapus atau dibiarkan tidak terpakai — lebih baik dihapus bersama test-nya.)

## G-6 🟡 CONTRACT EXTENSION COUNT HARDCODED

`service.go CreateContract` men-set `contract.ExtensionCount = 1` dengan komentar "dihitung manual oleh caller untuk extension > 1" — perpanjangan berantai tidak menghitung dengan benar.

**Rencana:** `ExtendContract` menghitung `extension_count = previous.extension_count + 1`, dan kontrak sebelumnya di-set `status = extended`.

### 3.10 Log Implementasi Langkah 7 (G-6 contract extension count) — 2026-08-10

**Tujuan:** perpanjangan kontrak berantai menghitung `extension_count` dengan benar (sebelumnya hardcoded `1`).

**Perubahan:**
- `repository.go ExtendContract` — di dalam transaksi: load previous contract (`tx.First`) → set `status = extended` → set `newContract.ExtensionCount = previous.ExtensionCount + 1` → create. Previous yang tidak ada → error + rollback (mencegah orphan contract).
- `service.go CreateContract` — hapus hardcoded `contract.ExtensionCount = 1`; count kini diturunkan dari rantai kontrak sebelumnya.
- Test: `TestRepo_ExtendContract_ChainCount` (2x extension → count 1 lalu 2, kedua previous berstatus `extended`), `TestRepo_ExtendContract_MissingPrevious` (rollback, 0 kontrak persist), `TestService_CreateContract_WithPrevious_ChainCount` (via service, response membawa count berantai).

**Validasi:** `go build` ✅ · `go vet` ✅ · test employeemovement/employee/approval **PASS** ✅.

## G-7 🟡 VALIDASI BISNIS PER TIPE MOVEMENT

Belum ada validasi "tipe X wajib field Y":
- `mutation` → wajib `to_organization_id` (dan/atau `to_position_id`)
- `promotion`/`demotion` → wajib `to_position_id`
- `status_change` → wajib `to_employment_status_id`
- `contract_extension` → wajib merujuk kontrak aktif
- `offboarding`/`retirement` → boleh tanpa `to_*`

**Rencana:** validasi service-level + `binding` DTO, kembalikan pesan field error (pola `getValidationErrors`).

### 3.11 Log Implementasi Langkah 8 (G-7 validasi per tipe) — 2026-08-10

**Tujuan:** menegakkan "tipe X wajib field Y" di service level dengan pesan field error yang bisa dipetakan handler ke 400.

**Perubahan:**
- `service.go` — tipe error baru `MovementValidationError`; `validateMovementFields(type, toOrg, toPos, toStatus, hasActiveContract)` menegakkan: `mutation` → wajib `to_organization_id` ATAU `to_position_id`; `promotion`/`demotion` → wajib `to_position_id`; `status_change` → wajib `to_employment_status_id`; `contract_extension` → wajib ada kontrak aktif; `offboarding`/`retirement` → tanpa validasi (boleh tanpa to_*). Dipanggil via `validateMovementCreate` / `validateMovementUpdate` (re-validasi saat movement_type berubah pada update).
- `repository.go` — `HasActiveContractByEmployeeID` (status = active) untuk dukungan contract_extension.
- `handler.go` — `movementErrStatus` memetakan `MovementValidationError` → 400 `VALIDATION_ERROR` (bukan 500).
- Test: 6 test validasi baru (mutation org/pos, promotion, status_change, contract_extension dengan/sans kontrak aktif, offboarding lolos tanpa to_*, update re-validasi). Helper `createTestMovement` beralih ke `MovementTypeOther` (tanpa kewajiban to_*); test lama promotion dilengkapi `to_position_id`.

**Validasi:** `go build` ✅ · `go vet` ✅ · test employeemovement/employee/approval **PASS** ✅.

## G-8 🟡 KONSISTENSI ROUTE/MENU & MODULE SLUG

- Backend menu (module.go): `/admin/career/movements`, `/admin/career/contracts`
- FE router saat ini: `/employee-movements` (satu halaman, module slug `employee-movement`)
- Approval module slug: `employeemovement` (tanpa tanda hubung)

**Keputusan user §11.6: DUA HALAMAN TERPISAH** mengikuti menu server:
- `/admin/career/movements` → halaman Movements (`EmployeeMovements.vue`)
- `/admin/career/contracts` → halaman Contracts (`EmployeeContracts.vue` baru)
- Router FE + sidebar disesuaikan; module slug disamakan (`employeemovement` atau `employee-movement` — pilih satu, samakan dengan permission & filter approval).

### 3.12 Log Implementasi Langkah 11 (G-8 slug/route konsisten) — 2026-08-10

**Tujuan:** menyamakan route FE dengan menu server & slug module.

**Perubahan (backend sudah benar — tanpa perubahan):** `ModuleSlug = "employeemovement"`, menu `/admin/career/movements` & `/admin/career/contracts` (module.go). Approval module slug juga `employeemovement` (konsisten, G-3 auto-resolve memakainya).

**FE:**
- `router/index.js` — route lama `employee-movements` (slug `employee-movement`) diganti dua route: `/admin/career/movements` (EmployeeMovements.vue, title `employee_movement.movements`) & `/admin/career/contracts` (EmployeeContracts.vue baru); slug `employeemovement` + backRoute `/admin/career`.
- `layouts/Sidebar.vue` — item `movement` diarahkan ke `/admin/career/movements` (label kini `employee_movement.movements`) + item baru `contracts` → `/admin/career/contracts` (ikon `pi-file-edit`).
- `views/modules/EmployeeContracts.vue` (baru) — placeholder mengikuti pola EmployeeMovements.vue.
- Locale `employee_movement.contracts_coming_soon` (en/id). `nav.movement` jadi dead key (Sidebar sudah tidak memakainya).

**Validasi:** `npm run build` FE **PASS** ✅ · JSON locale valid.

### 3.13 Log Implementasi Langkah 10 (FE halaman Movements) — 2026-08-10

**Tujuan:** halaman `/admin/career/movements` (EmployeeMovements.vue) — daftar enriched, filter, form create per tipe, aksi submit/execute/cancel/delete.

**Backend (penunjang filter list):**
- `ListMovements` (repo/service/handler) + 3 parameter opsional: `movement_type`, `status`, `search` (decision letter number / alasan, LIKE). Test filter baru di repository_test.go (`TestRepository_ListMovements_Filters*`).

**FE (`EmployeeMovements.vue` — rewrite penuh, sebelumnya placeholder):**
- Toolbar: total records, filter type (8 tipe), filter status (6 status), pencarian, tombol reset, tombol **Add Movement**.
- Tabel lazy pagination (SkeletonTable saat load): kolom employee (nama + employee_code), movement_type (Tag berwarna), to_position/to_organization, to_employment_status, decision_letter_number, effective_date (formatDate bilingual), status (Tag), aksi.
- Aksi per baris sesuai status: **Submit** (draft → pending_approval), **Execute** (approved → executed), **Cancel** (pending_approval/approved), **Delete** (draft), semuanya via ConfirmActionDialog/ConfirmDeleteDialog.
- Dialog create: pilih employee + movement_type; field `to_*` tampil kondisional sesuai tipe (mutation → to_organization (+to_position), promotion/demotion → to_position, status_change → to_employment_status); decision_letter_number/date + effective_date (DateInput) + reason. Validasi frontend mengikuti aturan G-7; error VALIDATION_ERROR dari backend ditampilkan per-field via getValidationErrors.
- Referensi dropdown: employees (`per_page=500`), organizations **`active_only=true`** (summary aktif — sesuai instruksi user: *position ambil dari organization yang organization summary active*), employment-statuses (`settings/employment-statuses`).
- Locale `employee_movement.*` lengkap (en/id): label tipe & status, konfirmasi aksi, hint per tipe, field_required.

**Validasi:** `go build` + `go vet` + test employeemovement **PASS** ✅ · `npm run build` FE **PASS** ✅ · JSON locale valid ✅

### 3.14 Log Implementasi Langkah 11 (FE halaman Contracts) — 2026-08-10

**Tujuan:** halaman `/admin/career/contracts` (EmployeeContracts.vue, sebelumnya placeholder) — daftar enriched, filter, create/edit, upload dokumen, delete.

**Backend (penunjang filter list):** `ListContracts` (repo/service/handler) + 2 parameter opsional: `status`, `search` (contract_number / nama employee / employee_id, parameterized LIKE + escape wildcard, JOIN employees). Test baru `TestRepo_ListContracts_Filters` (filter status, search by number, kombinasi tak-cocok) + seed minimal tabel employees.

**FE (`EmployeeContracts.vue` — rewrite penuh):**
- Toolbar: total records, filter status (active/expired/extended/terminated), pencarian, reset, tombol **New Contract**.
- Tabel lazy pagination + SkeletonTable: employee (nama + kode), contract_number, contract_type (Tag), start_date/end_date (formatDate bilingual), extension_count (badge xN), status (Tag), dokumen (link lampiran), aksi edit/delete.
- Dialog create/edit: employee (create only, disabled saat edit), contract_number, contract_type (pkwt/pkwtt/daily/other), start_date (+ minDate utk end_date), end_date, previous_contract_id (opsional — opsi = kontrak milik employee terpilih, minus dirinya), decision_letter_number, notes, upload dokumen, status (edit only).
- Upload: `POST /api/v1/tenant/uploads` (FormData `file`) → `data.url` disimpan ke `document_url` (pola AttendanceOvertime).
- Delete via ConfirmDeleteDialog; error VALIDATION_ERROR per-field via getValidationErrors.
- Locale: key `contracts_coming_soon` **dihapus**, diganti keys contract lengkap (en/id).

**Validasi:** `go build` + `go vet` + test employeemovement **PASS** ✅ · `npm run build` FE **PASS** ✅ · JSON locale valid ✅ (tanpa dead key `contracts_coming_soon`)

## 3.15 Log Implementasi Langkah 12 (FE detail dialog movement) — 2026-08-10

**Tujuan:** menyelesaikan langkah 12: detail movement + deep-link notifikasi + badge `rejected`.

### 3.15.1 Audit — sebagian sudah ada

Saat dicek ulang, dua dari tiga item langkah 12 **sudah terimplementasi pada langkah sebelumnya**:

- **Badge status `rejected`** ✅ — `statusSeverity` di `EmployeeMovements.vue` sudah memetakan `rejected → danger` + locale `status_rejected` (masuk bersamaan langkah 10).
- **Deep-link notifikasi** ✅ — `Notifications.vue` sudah punya case `reference_type == 'employeemovement'` ATAU `type` berawalan `MOVEMENT_` → `/admin/career/movements` (masuk bersamaan langkah 5/G-2). `Approvals.vue` juga sudah punya `case 'employeemovement'` → `GET /employee-movements/movements/:id`.

### 3.15.2 Perubahan FE

- `frontend/tenant/src/views/modules/EmployeeMovements.vue`:
  - Tombol **View (detail)** (`pi pi-eye`) di kolom aksi setiap baris → membuka **Dialog Detail Movement** (`detailVisible`/`detailItem`, lebar 680px):
    - Header: Tag tipe movement + Tag status (severity warna sama dengan tabel).
    - Ringkasan karyawan (nama + kode) & nomor SK (mono, break-all).
    - Section **Dari → Ke**: box `from_*` (abu-abu) dan box `to_*` (hijau) berisi nama organization/position/employment_status dari respons enriched (G-4); box hanya dirender bila ada data (`hasAnyField`).
    - Tanggal: decision_letter_date, effective_date (formatDate bilingual), created_at (formatDateTime).
    - Alasan & catatan (bila ada).
    - Riwayat approval: `approved_at` / `executed_at` (bila ada) — `approved_by`/`executed_by` UUID tidak ditampilkan (tidak human-readable).
  - **Aksi dari dalam dialog detail** mengikuti status: `draft` → Submit + Delete; `pending_approval`/`approved` → Cancel; `approved` → Execute; tombol Close. Aksi membuka ConfirmActionDialog/ConfirmDeleteDialog yang sama dengan tabel (reuse `actionTarget`).
  - Import `ViewLabel.vue` (komponen display key-value yang sudah dipakai `CompanyDetail.vue`).
- Locale `employee_movement.*` (en/id): +6 key — `detail_title`, `from`, `to`, `approved_at`, `executed_at`, `created_at`.

### 3.15.3 Validasi

- `npm run build` (tenant FE) **PASS** ✅ · JSON locale en/id valid ✅ · tidak ada dead key baru ✅

## 3.16 Log Implementasi Langkah 13 (checklist verifikasi E2E manual) — 2026-08-10

**Tujuan:** menyediakan checklist verifikasi E2E manual sesuai Testing Plan §10 — alur lengkap movement (draft → submit → approve → execute → cek employment), offboarding/retirement, CRUD kontrak (termasuk extension chain), detail/deep-link, dan regresi.

**Deliverable:** `docs/module-movement-e2e-checklist.md` (baru) — berisi:

- **§0 Prasyarat**: menjalankan backend (`make run`) + FE (`npm run dev` :5174), akun tenant ber-permission `employeemovement.*` + `approval.*`, **approval flow aktif untuk module `employeemovement`** (tanpa ini submit gagal, G-3), data pendukung, dan tabel data yang dipakai.
- **§1 Skenario A (promotion, inti G-1)**: 10 langkah A1–A10 + edge negatif (submit tanpa flow, execute sebelum approved, double execute, approve manual 404). Verifikasi employment via `GET /employees/:id` → `employments[]` (lama `effective_end_date = effective_date − 1`, baru berisi `to_*` + SK + `effective_date`).
- **§2 Skenario B (offboarding/retirement)**: employee `status = inactive`, tanpa employment baru (§11.3).
- **§3 Skenario C (CRUD kontrak)**: create PKWT + upload, filter, edit, **extension chain** (G-6: `extension_count` berantai + previous `extended`), delete.
- **§4 Skenario D (langkah 12)**: detail dialog (badge `rejected`), aksi dari dialog, deep-link notifikasi `MOVEMENT_*` → `/admin/career/movements`, deep-link approval.
- **§5 Regresi & lintas bahasa**: switch EN/ID, filter kombinasi, build & console bersih, catatan permission FE (tombol aksi belum di-gate `hasPermission` — batasan ada di backend authz).
- **§6–8**: kriteria penerimaan (acceptance plan §10), bukti yang disimpan, dan referensi endpoint (movement/contract/approval/employee/upload/notification).

**Endpoint diverifikasi terhadap kode:** `POST /approval/instances/:id/actions` body `{"action":"APPROVE"|"REJECT"}` (`approval/dto.go` `SubmitActionRequest`), `GET /approval/active-flow?module=...`, `GET /employees/:id` → `employments[]` (`employee/dto.go` `EmploymentResponse`), `GET /notifications` (type `MOVEMENT_*`).

**Status:** checklist **siap dieksekusi**; eksekusi manual memerlukan environment tenant berjalan + akun ber-permission `employeemovement.*` (catatan §3.5.5: akun tsb belum tersedia saat validasi sebelumnya) — jadi bagian ini **menunggu environment**, bukan perubahan kode.

## 3.17 Log Implementasi — Enhancement §12 P0 1–3 (transactional execute + conflict detection) 2026-08-10

**Tujuan:** menutup tiga gap P0 enhancement plan (§12.2–§12.4): eksekusi atomic, conflict posisi, dan conflict tanggal efektif.

### 3.17.1 Perubahan

| File | Perubahan |
|---|---|
| `employeemovement/service.go` | +`MovementConflictError`; interface `CareerExecutor` **tx-aware** (semua method menerima `*gorm.DB`); `validateMovement` + **position conflict check saat create/update** (§12.3 soft-check); `ExecuteMovement` ditulis ulang → `repo.ExecuteMovementTx` dengan closure: conflict checks (§12.3 posisi + §12.4 tanggal) → `FindCurrentEmployment` → `CloseEmployment` → `CreateEmployment`/`SetEmployeeInactive` → return `toEmploymentID`; `contract_extension` tetap jalur non-tx (tanpa HR data change) |
| `employeemovement/repository.go` | +`ExecuteMovementTx(ctx, id, executedBy, hrChanges func(tx) (*uuid.UUID, error))` — BEGIN → reload movement (guard `approved`) → jalankan HR changes → update `executed` + `to_employment_id` → COMMIT; error → ROLLBACK. +`PositionConflict(ctx, tx, posID, excludeEmp, effectiveDate)` (posisi terisi employee lain saat effective date) +`EmploymentEffectiveDateConflict(ctx, tx, empID, effectiveDate)` (employment terbuka yang mulai ≥ effective date). Keduanya menerima `tx` (nil → koneksi biasa) |
| `employee/repository.go` | 4 metode di-refactor ke helper bersama + varian **Tx**: `FindActiveEmploymentByEmployeeIDTx`, `CloseEmploymentTx`, `CreateEmploymentTx`, `SetEmployeeStatusTx` — dipakai adapter agar operasi berjalan di transaksi modul employeemovement (satu koneksi, satu unit kerja) |
| `cmd/server/main.go` | `employeeCareerAdapter` tx-aware (meneruskan `tx` ke varian Tx repository) |
| `employeemovement/handler.go` | `movementErrStatus` + `MovementConflictError` → **409 `CONFLICT_ERROR`** (create/update/execute); execute membedakan konflik dari `EXECUTE_FAILED` |

### 3.17.2 Aturan conflict yang diimplementasikan

- **Posisi (§12.3)** — `position_id = to_position_id AND employee_id <> employee AND (effective_end_date IS NULL OR effective_end_date >= effective_date)` → terisi → ditolak. Dicek saat draft (create/update) dan **diulang atomik di execute**.
- **Tanggal efektif (§12.4)** — untuk employee yang sama: ada employment `effective_end_date IS NULL` dengan `effective_date >= movement.effective_date` → overlap (mencakup backdate ke employment aktif maupun tabrakan dengan employment future-dated hasil eksekusi movement lain) → ditolak. Future-dated movement berurutan (chain) tetap diizinkan: eksekusi A (09-01) lalu B (08-15) → B ditolak; sebaliknya B dulu lalu A → A menutup B (08-31) lalu membuat A (09-01) — tanpa overlap.
- **Guard offboarding/retirement** — cek tanggal efektif juga dijalankan untuk tipe deactivation: mencegah `CloseEmployment` menutup employment future-dated menjadi periode invalid (`effective_end_date` sebelum `effective_date` sendiri) sementara employment aktif sebenarnya tetap terbuka (hasil review kode saat implementasi).

### 3.17.3 Tradeoff & catatan

- Conflict check posisi saat **create** bersifat **blocking** (strict, sesuai acceptance §12.3 "validation dilakukan saat create/submit dan diulang saat execute"). Konsekuensi: skenario swap dua karyawan harus dieksekusi berurutan (buat+execute yang pertama dulu) — belum ada dukungan swap eksplisit (exception §12.3 di-defer).
- Tipe `mutation` yang hanya mengisi `to_organization_id` (tanpa posisi) tidak bisa dicek occupancy posisi (P1 §12.12 — konsistensi org+posisi belum masuk scope ini).
- `PositionConflict`/`EmploymentEffectiveDateConflict` memakai raw query ke tabel `employments` (pola enrichment G-4 yang sudah ada), bukan import modul employee.
- Tidak ada row-locking `FOR UPDATE` (agar kompatibel driver SQLite di test) — anti double-execute dijaga oleh guard `WHERE status = 'approved'` pada update dalam transaksi.

### 3.17.4 Validasi

- `go build ./...` ✅ · `go vet` (employeemovement, employee, cmd/server) ✅
- `go test` employeemovement / employee / approval — semua **PASS** ✅
- Test baru (7): create posisi terisi → `MovementConflictError`; create posisi dibebaskan sebelum effective date → lolos; execute posisi terisi → 409 + movement tetap `approved` tanpa perubahan HR; execute tanggal efektif overlap employment future-dated → ditolak; execute backdate → ditolak; **offboarding tabrakan employment future-dated → ditolak** (guard §12.4); **rollback**: gagal create employment → movement tetap `approved` (bukan sebagian `executed`).

---

# 3.18 Log Implementasi — Enhancement P0 item 4: Movement Snapshot (2026-08-10)

## 3.18.1 Yang dikerjakan

Migration + model + service snapshot sehingga histori movement tidak berubah ketika master data Organization/Position/EmploymentStatus diubah namanya (plan §12.5):

| File | Perubahan |
|---|---|
| `migrations/tenant/{mysql,postgres}/083_employeemovement_snapshot.{down.}sql` | Migration up/down: 6 kolom baru `from_organization_name`, `from_position_name`, `from_employment_status_name`, `to_organization_name`, `to_position_name`, `to_employment_status_name` (VARCHAR(255), nullable) di `employee_movements` |
| `employeemovement/model.go` | `EmployeeMovement` + 6 field snapshot (`gorm` column mapped; nullable string) |
| `employeemovement/dto.go` | `ToResponse()` mengisi 6 field `*Name` langsung dari snapshot row |
| `employeemovement/service.go` | Helper `fillMovementSnapshot` — resolve nama (batch query `GetOrganizationNamesByIDs`/`GetPositionNamesByIDs`/`GetEmploymentStatusNamesByIDs`) dari `from_*/to_*_id` lalu persist ke row; dipanggil di `CreateMovement` & `UpdateMovement`. `enrichMovementResponses` kini **snapshot-aware**: hanya mengisi nama yang masih kosong (fallback untuk movement lama tanpa snapshot), tidak menimpa nilai snapshot |
| `employeemovement/enrichment_test.go` | +2 test: snapshot ter-persist saat create (cek DB row & response); snapshot tetap saat master position di-rename setelahnya |

## 3.18.2 Keputusan desain

1. **Snapshot diisi saat create/update** (bukan saat submit/execute) — plan §12.5 "diisi saat movement dibuat/submit"; update draft juga menyegarkan karena `to_*` boleh berubah sebelum submit.
2. **Best-effort resolution**: bila id referensi tidak ter-resolve (data master terhapus), nama snapshot dibiarkan kosong dan enrichment G-4 mengisi dari live data sebagai fallback — create/update tidak pernah gagal karena nama tidak ketemu.
3. **Enrichment tidak menimpa snapshot** — setelah migration, daftar/detail tetap memakai nama snapshot (kondisi master saat movement dibuat). Movement lama (snapshot NULL) tetap menampilkan nama live via enrichment.
4. **Foreign key `from_*/to_*_id` tetap disimpan** — relasi & navigasi tidak berubah (plan §12.5).
5. `from_employment_type_name`/`to_employment_type_name` (opsional §14) **tidak diimplementasikan** — employment type belum ada di model; dapat ditambahkan belakangan bila kebutuhan muncul.

## 3.18.3 Validasi

- `go build ./...` ✅ · `go vet ./internal/modules/employeemovement/...` ✅ · `gofmt` bersih ✅
- `go test ./internal/modules/employeemovement/...` — semua **PASS** ✅
- Test baru (2): `TestService_CreateMovement_SnapshotPersisted` (snapshot tersimpan di row + tampil di respons) · `TestService_Snapshot_ImmutableOnMasterRename` (rename title posisi setelah create → nama snapshot tetap `Software Engineer`).

---

# 3.19 Log Implementasi — Enhancement P0 item 5: Movement Audit Trail (2026-08-10)

## 3.19.1 Yang dikerjakan

Mencatat seluruh perubahan lifecycle movement sehingga transaksi HR dapat diaudit (plan §12.6):

| File | Perubahan |
|---|---|
| `migrations/tenant/{mysql,postgres}/084_employeemovement_audit.{down.}sql` | Tabel baru `employee_movement_audits`: `id`, `movement_id` (FK → `employee_movements` ON DELETE CASCADE, index), `action` (index), `old_status`/`new_status` (VARCHAR(20)), `old_data`/`new_data` (JSON snapshot movement sebelum/sesudah), `reason`, `acted_by` (CHAR(36), nullable utk aksi sistem/callback), `acted_at` (default CURRENT_TIMESTAMP) |
| `employeemovement/model.go` | +`EmployeeMovementAudit` + enum `MovementAuditAction` (CREATED, UPDATED, SUBMITTED, APPROVED, REJECTED, CANCELLED, EXECUTED); `BeforeCreate` set ID + `acted_at` |
| `employeemovement/dto.go` | +`MovementAuditResponse` + `PaginatedMovementAuditResponse` + `ToResponse` |
| `employeemovement/repository.go` | +`CreateAudit` (best-effort, tidak menggagalkan operasi utama) + `ListAuditsByMovementID` (acted_at DESC, pagination) |
| `employeemovement/service.go` | +helper `recordAudit(ctx, movementID, action, oldStatus, newStatus, oldData, newData, reason, actedBy)` + `movementAuditJSON` (marshal movement → JSON snapshot) + `statusPtr`. Dipanggil di: **CreateMovement** (CREATED, new snapshot) · **UpdateMovement** (UPDATED, old+new snapshot) · **SubmitMovement** (SUBMITTED, draft → pending_approval) · **HandleApprovalStatusChange** (APPROVED/REJECTED/CANCELLED, snapshot sebelum + note reject sbg reason; acted_by nil krn aksi dari push-callback) · **ExecuteMovement** (EXECUTED, reload state akhir + acted_by executor; kedua jalur: contract_extension non-tx & transaksi employment) · **CancelMovement** (CANCELLED, newData mencerminkan status akhir) |
| `employeemovement/handler.go` + `routes.go` | +`ListMovementAudits` → **`GET /api/v1/tenant/employee-movements/movements/:id/audits`** (page/per_page) |
| `employeemovement/module.go` | `Migrate` AutoMigrate + `&EmployeeMovementAudit{}` |
| `employeemovement/helpers_test.go` | AutoMigrate test DB + audit model |
| `employeemovement/audit_test.go` (baru) | 8 test: lifecycle per aksi (create/update/submit/approve/reject/cancel/execute) + pagination (newest-first) + handler GET audits |

## 3.19.2 Keputusan desain

1. **Audit terpisah dari notifikasi**: `recordAudit` best-effort (kegagalan hanya warn log) — mengikuti pola `organization.captureHistory`, bukan menggagalkan operasi movement.
2. **Snapshot JSON utuh**: `old_data`/`new_data` berisi marshal penuh `EmployeeMovement` (bukan hanya diff) sehingga histori dapat direkonstruksi meskipun row diubah setelahnya.
3. **acted_by nil untuk push-callback** — aksi APPROVED/REJECTED/CANCELLED datang dari Central Approval (bukan user langsung); CREATED/UPDATED/SUBMITTED/CANCELLED-user memakai `authctx.GetUserID`, EXECUTED memakai id executor dari request.
4. **FK ON DELETE CASCADE** — audit otomatis terhapus bila movement dihapus (konsisten dengan `employee_id` CASCADE pada movement).
5. **Endpoint** mengikuti pola route modul: `GET /movements/:id/audits` (bukan `/employee-movements/{id}/audits` di plan §15 yang memakai shorthand path modul).

## 3.19.3 Validasi

- `go build ./...` ✅ · `go vet ./internal/modules/employeemovement/...` ✅ · `gofmt` bersih ✅
- `go test ./internal/modules/employeemovement/...` — semua **PASS** ✅
- Test baru (8): lifecycle CREATED/UPDATED/SUBMITTED/APPROVED/REJECTED/CANCELLED/EXECUTED (status transition + snapshot + reason reject + acted_by executor) · pagination (total/total_pages + newest-first) · handler `GET /audits` 200.

---

# 3.20 Log — Enhancement §12 P1 (6) Movement Documents

## 3.20.1 Yang dikerjakan

Mendukung lebih dari satu dokumen per movement (plan §12.15) — selain decision letter fields yang sudah ada:

| File | Perubahan |
|---|---|
| `migrations/tenant/{mysql,postgres}/085_employeemovement_documents.{down.}sql` | Tabel baru `employee_movement_documents`: `id`, `movement_id` (FK → `employee_movements` ON DELETE CASCADE, index), `document_type` (VARCHAR(30): PROMOTION_SK, MUTATION_SK, DEMOTION_SK, RETIREMENT_LETTER, OFFBOARDING_LETTER, OTHER), `file_name`, `file_url` (VARCHAR(255)), `uploaded_by` (CHAR(36), nullable), `created_at` (default CURRENT_TIMESTAMP) |
| `employeemovement/model.go` | +`EmployeeMovementDocument` + enum `MovementDocumentType`; `BeforeCreate` set ID |
| `employeemovement/dto.go` | +`CreateMovementDocumentRequest` (document_type oneof + file_name/file_url required; file_url wajib diawali `/`) + `MovementDocumentResponse` + `PaginatedMovementDocumentResponse` + `ToResponse` |
| `employeemovement/repository.go` | +`CreateMovementDocument` + `ListDocumentsByMovementID` (created_at DESC, pagination) + `FindMovementDocumentByID` + `DeleteMovementDocument` (RowsAffected 0 → not found) |
| `employeemovement/service.go` | +`CreateMovementDocument` (validasi movement ada + persist metadata + `uploaded_by` dari authctx) + `ListMovementDocuments` (pagination) + `DeleteMovementDocument` |
| `employeemovement/handler.go` + `routes.go` | +3 endpoint: `GET`/`POST /movements/:id/documents` + `DELETE /movements/:id/documents/:documentId` |
| `employeemovement/module.go` | `Migrate` AutoMigrate + `&EmployeeMovementDocument{}` |
| `employeemovement/helpers_test.go` | AutoMigrate test DB + document model |
| `employeemovement/document_test.go` (baru) | 6 test: create (metadata + movement not found) · list (newest-first + pagination) · delete (sisa dokumen lain tetap ada + not found) · handler POST/GET/DELETE end-to-end |

## 3.20.2 Keputusan desain

1. **File fisik vs metadata dipisah**: upload file lewat endpoint upload generik `POST /api/v1/tenant/uploads` (sudah ada — validasi ekstensi + ukuran, menyimpan ke `{uploadDir}/attachments`, mengembalikan URL publik `/uploads/attachments/{uuid}{ext}`). Tabel baru hanya menyimpan metadata (`file_url`), konsisten dengan pola `employee_contracts.document_url` — tidak ada duplikasi multipart handling di module.
2. **`file_url` wajib diawali `/`** (binding `startswith=/`) karena URL publik selalu path relatif hasil upload generik.
3. **Validasi movement ada di service** sebelum insert metadata — selain jadi FK guard, memberi error yang jelas (bukan FK violation cryptic).
4. **Belum di-audit (plan §12.6 tidak mencakup aksi dokumen)** — lifecycle audit tetap untuk perubahan status movement saja, sesuai plan.
5. **CASCADE hapus**: dokumen ikut terhapus bila movement dihapus (FK ON DELETE CASCADE).

## 3.20.3 Validasi

- `go build ./...` ✅ · `go vet ./internal/modules/employeemovement/...` ✅ · `gofmt` bersih ✅
- `go test ./internal/modules/employeemovement/...` — semua **PASS** ✅ (deterministik — dijalankan 3x)
- Test baru (6): create metadata + movement not found · list newest-first + pagination (total/total_pages) · delete + not found + dokumen lain tetap ada · handler `POST` 201 → `GET` 200 → `DELETE` 200.
- Catatan: test ordering pakai jeda 2ms antar insert karena SQLite menyimpan `created_at` mikro-detik (dua insert cepat bisa berbagi timestamp → urutan jatuh ke tie-break UUID acak).

---

# 3.21 Log — Enhancement §12 P1 (7) Career Timeline

## 3.21.1 Yang dikerjakan

Career timeline read model (plan §12.8) — **tanpa tabel duplikasi** `employee_career_history`:

| File | Perubahan |
|---|---|
| `employeemovement/repository.go` | +3 query non-paginated utk timeline: `FindEmploymentsByEmployeeID` (table `employments`, effective_date ASC), `FindExecutedMovementsByEmployeeID` (hanya status `executed`), `FindAllContractsByEmployeeID` (start_date ASC) |
| `employeemovement/dto.go` | +`CareerHistoryResponse`/`CareerHistoryData` + `CareerPositionInfo` (current position) + `CareerTimelineEntry` (date, event_type, title, description, movement_type/contract_type, referensi id sumber) |
| `employeemovement/service.go` | +`GetCareerHistory`: baca 3 sumber → resolve nama org/posisi/status employment via batch query (G-4; movement pakai snapshot §12.5) → bangun event JOINED (employment pertama) + MOVEMENT (executed) + CONTRACT → urut kronologis ASC (tanggal sama: JOINED→MOVEMENT→CONTRACT) → hitung current position (employment terbuka terakhir) · +helper `movementFromToLabel` (label "dari → ke" dari snapshot) + `careerEventPriority` + `currentPosition` + `normalizeDate` (kolom DATE bisa dikembalikan driver sbg DATETIME/RFC3339) |
| `employeemovement/handler.go` + `routes.go` | +`GetCareerHistory` → **`GET /api/v1/tenant/employee-movements/employees/:employeeId/career-history`** |
| `employeemovement/career_test.go` (baru) | 4 test: timeline lengkap (JOINED→CONTRACT→MOVEMENT + snapshot from→to + current position org/posisi/status) · hanya movement executed yg masuk · current position prefer employment terbuka · handler 200 |

## 3.21.2 Keputusan desain

1. **Read model, bukan tabel baru**: timeline dibentuk dari `employee_movements` + `employments` + `employee_contracts` (keputusan §12.8 / §13 — tidak membuat `employee_career_history`).
2. **Movement = sumber transaksi**: hanya status `executed` masuk timeline; draft/pending bukan histori nyata.
3. **JOINED dari employment pertama** (effective_date terawal); current position = employment **terbuka** (effective_end_date NULL) dengan tanggal efektif terbesar, fallback employment terakhir.
4. **Movement memakai snapshot names** (§12.5) — histori tidak berubah walau master data diganti nama; employment di-resolve via batch query (G-4).
5. **normalizeDate** dipakai utk semua tanggal respons (kolom DATE dapat dikembalikan driver sbg DATETIME/RFC3339 — sama seperti yg sudah diantisipasi `dayBefore`).

## 3.21.3 Validasi

- `go build ./...` ✅ · `go vet ./internal/modules/employeemovement/...` ✅ · `gofmt` bersih ✅
- `go test ./internal/modules/employeemovement/...` — semua **PASS** ✅ (deterministik — dijalankan 4x)
- Test baru (4): JOINED→CONTRACT→MOVEMENT terurut + deskripsi from→to snapshot + referensi id sumber · hanya executed movements · current position prefer employment terbuka · handler `GET /career-history` 200.
- Perbaikan determinisme: helper `createAuditedMovement` kini memberi jeda 2ms setelah create (pola sama document_test) karena SQLite menyimpan `acted_at` mikro-detik — dua aksi cepat bisa berbagi timestamp dan urutan audit jatuh ke tie-break UUID acak.

---

# 3.22 Log — Enhancement §12 P1 (8) Contract Expiry Management

## 3.22.1 Yang dikerjakan

Scheduled process `ProcessContractExpiration` (plan §12.13) — mark expired otomatis + reminder ke employee & HR:

| File | Perubahan |
|---|---|
| `employeemovement/repository.go` | +`FindContractsExpiringOn` (status=active, end_date = target) + `FindContractsExpiredBefore` (status=active, end_date < tanggal) + `MarkContractsExpired` (batch by ids, guard status active) + `FindUserIDsWithPermission` (resolve user HR via permissions + role_has_permissions + model_has_roles/model_has_permissions — UNION, cross-dialect MySQL/Postgres) |
| `employeemovement/service.go` | +`ProcessContractExpiration(ctx)`: (1) mark expired kontrak active yang end_date < hari ini → notif `CONTRACT_EXPIRED`; (2) reminder H-30/H-14/H-7/H-1 → notif `CONTRACT_EXPIRING` (params: nomor kontrak + tanggal akhir). +helper `notifyContractEvent` (employee pemilik via akun terhubung + seluruh user HR via permission `employeemovement.view`; best-effort) + `addDays` · const `contractExpiryReminderDays` = {30,14,7,1} + `contractExpiryHRPermission` |
| `notification/i18n.go` | +2 entri catalog: `CONTRACT_EXPIRING` + `CONTRACT_EXPIRED` (EN/ID, %s placeholder nomor kontrak + tanggal) |
| `cmd/server/main.go` | +`runContractExpirationScheduler` (goroutine + `time.Ticker` 24 jam — **tanpa dependency cron baru**, keputusan user) + `runContractExpirationPass` (iterate company status=active di platform DB, panggil `ProcessContractExpiration` per tenant dgn company_id di context) |
| `employeemovement/expiry_test.go` (baru) | 3 test: mark expired (past→expired, future tetap active, notif CONTRACT_EXPIRED) · reminder H-30/H-7 dapat notif, H-15 tidak (bukan jadwal) · tanpa notifier tetap jalan |

## 3.22.2 Keputusan desain

1. **Scheduler tanpa dependency baru**: goroutine + `time.Ticker` harian (24 jam) di main.go — keputusan dikonfirmasi user (codebase belum punya infra cron). Pass pertama langsung dijalankan saat server start.
2. **Reminder berbasis tanggal tepat**: kontrak di-remind hanya pada hari H-30/H-14/H-7/H-1 (bukan rentang) — mencegah notifikasi berulang tiap hari; kontrak di luar jadwal tidak dapat notif.
3. **Penerima notifikasi**: employee pemilik kontrak (via `employee_accounts` → user) + seluruh user HR (semua yang punya permission `employeemovement.view`). Ini memenuhi "Notification dikirim kepada HR" (plan §12.13) tanpa menambah relasi manager.
4. **Best-effort per kontrak**: kegagalan notifikasi satu kontrak hanya di-log, tidak menggagalkan sisa proses. `MarkContractsExpired` guard ulang status active di query.
5. **`ProcessContractExpiration` murni per-tenant** (company_id dari context) dan tanpa tahu jadwal — mudah diuji; scheduler yang menentukan kapan dipanggil.
6. **Penerima di-dedup**: employee yang sekaligus HR (memegang permission `employeemovement.view`) hanya menerima satu notifikasi (set user id, bukan dua jalur terpisah) — perbaikan dari hasil review.
7. **Keterbatasan yang disengaja**: (a) reminder berbasis tanggal TEPAT (`end_date = hari ini + N`) — bila server down di H-30, reminder H-30 terlewat dan tidak menyusul (range-based + guard "sudah di-remind" bisa ditambahkan kemudian); (b) kontrak dengan `end_date = hari ini` masih berstatus active sepanjang hari kedaluwarsanya dan baru dipindah ke expired sehari setelahnya (interpretasi end date = hari terakhir kontrak berlaku); (c) `time.Now()` memakai zona waktu server — konsisten selama server TZ disetel ke TZ tenant.

## 3.22.3 Validasi

- `go build ./...` ✅ · `go vet ./internal/modules/employeemovement/... ./internal/modules/notification/...` ✅ · `gofmt` bersih ✅
- `go test ./internal/modules/employeemovement/... ./internal/modules/notification/...` — semua **PASS** ✅
- Test baru (3): mark expired + notif CONTRACT_EXPIRED · reminder H-30/H-7 ✓ & H-15 ✗ · tanpa notifier tetap mark expired.

---

# 5. API Plan

## 5.1 Endpoint Existing (sudah jalan)

```http
POST  /api/v1/tenant/employee-movements/movements              # buat draft
GET   /api/v1/tenant/employee-movements/movements              # list semua (pagination)
GET   /api/v1/tenant/employee-movements/movements/:id          # detail
PUT   /api/v1/tenant/employee-movements/movements/:id          # update (hanya draft)
DELETE /api/v1/tenant/employee-movements/movements/:id         # hapus (hanya draft)
POST  /api/v1/tenant/employee-movements/movements/:id/submit   # kirim approval (flow_id)
POST  /api/v1/tenant/employee-movements/movements/:id/approve  # manual (G-5) — ❌ DIHAPUS (§11.5)
POST  /api/v1/tenant/employee-movements/movements/:id/execute  # eksekusi (G-1)
POST  /api/v1/tenant/employee-movements/movements/:id/cancel   # batal
GET   /api/v1/tenant/employee-movements/employees/:employeeId/movements
POST  /api/v1/tenant/employee-movements/contracts              # buat kontrak (extend bila previous_contract_id)
GET   /api/v1/tenant/employee-movements/contracts
GET   /api/v1/tenant/employee-movements/contracts/:id
PUT   /api/v1/tenant/employee-movements/contracts/:id
DELETE /api/v1/tenant/employee-movements/contracts/:id
GET   /api/v1/tenant/employee-movements/employees/:employeeId/contracts
```

## 5.2 Perubahan yang Direncanakan

- `GET/POST /movements` — respons + nama-nama (G-4), filter `movement_type`/`status`/`employee_id` + search (G-4).
- `POST /movements/:id/submit` — `flow_id` opsional (auto-resolve G-3).
- `POST /movements/:id/execute` — transaksi employment (G-1) + employee non-aktif utk offboarding/retirement (§11.3).
- `status` model + migration: tambah `rejected` (§11.4); handler approval set status `rejected` saat REJECTED (bukan `cancelled`).
- `POST /uploads` (generik, sudah ada) — dipakai `contracts.document_url` (G-2/G-6).
- Notifikasi `MOVEMENT_*` (G-2).

---

# 6. Frontend Plan

> Status saat ini: `EmployeeMovements.vue` = placeholder 1 baris. Seluruh poin di bawah **belum ada**.

## 6.1 Halaman & Navigasi

> **Keputusan §11.6: DUA halaman terpisah** (mengikuti menu server), bukan satu halaman tab.

- **`EmployeeMovements.vue`** → halaman Movements di route `/admin/career/movements`.
- **`EmployeeContracts.vue`** (baru) → halaman Contracts di route `/admin/career/contracts`.
- Router FE + sidebar disesuaikan dengan menu server (pola `AdminLayout`/sidebar yang sudah ada).
- Header halaman: judul + deskripsi + tombol **+ Movement** / **+ Contract**.
- Konsisten dengan pola halaman lain: `SkeletonTable` saat load, `DataTable` lazy pagination, filter + search, tombol aksi per baris, `ConfirmDeleteDialog` untuk hapus/batal.

## 6.2 Tab Movements

- **Kolom tabel**: employee (nama + kode), movement_type (badge), dari → ke (organisasi/posisi), effective_date, decision_letter_number, status (badge: draft/pending_approval/approved/**rejected**/executed/cancelled), dibuat tanggal.
- **Filter**: tipe movement, status, search nama karyawan.
- **Dialog form (create/update)**: 
  - dropdown employee (`GET /employees?per_page=500`)
  - dropdown movement_type (bilingual)
  - field dari → ke: organization (`GET /organizations` tree), position, employment_status — field "dari" terisi otomatis dari employment aktif karyawan saat tipe dipilih (default), bisa dikoreksi
  - decision_letter_number/date, effective_date (DateInput), reason, notes
  - validasi dinamis per tipe (G-7) — tampilkan error field
- **Aksi per baris** (mengikuti status):
  - `draft`: Edit, **Submit** (→ approval), Delete
  - `pending_approval`: Cancel
  - `approved`: **Execute**
  - aksi memakai `ConfirmDeleteDialog` (seragam) untuk konfirmasi
- **Detail drawer/dialog**: ringkasan movement + tombol aksi.

## 6.3 Halaman Contracts (terpisah, §11.6)

- **Kolom**: employee, contract_number, contract_type (badge), start/end date, extension_count, status (badge: active/expired/extended/terminated), link dokumen.
- **Dialog form**: contract_number, contract_type, start_date, end_date, decision_letter_number, notes, **upload dokumen** (`POST /uploads` → document_url), opsi **perpanjangan** (pilih kontrak sebelumnya → rantai extension_count, G-6).
- **Aksi**: edit (aktif), terminate (dengan konfirmasi), hapus.
- Filter: tipe, status.

## 6.4 Integrasi & Konsistensi

- **Approval deep-link**: `Approvals.vue` sudah punya `case 'employeemovement'` → detail movement di modal approval (renderer khusus seperti attendance/overtime bila perlu).
- **Notifikasi**: deep-link `MOVEMENT_*` → `/employee-movements` (pola `Notifications.vue`).
- **Locale**: lengkapi `employee_movement.*` bilingual (en/id): status, tipe, label form, pesan.
- **Permission**: `employeemovement.*` (create/update/delete/approve/execute) untuk tampil/sembunyikan aksi (pola `hasPermission`).

---

# 7. Notification Plan

| Type | Penerima | Trigger |
|---|---|---|
| `MOVEMENT_SUBMITTED` | approver | `SubmitMovement` → instance approval dibuat |
| `MOVEMENT_APPROVED` | pengaju / employee | instance APPROVED (push-callback) |
| `MOVEMENT_REJECTED` | pengaju / employee | instance REJECTED |
| `MOVEMENT_EXECUTED` | employee | `ExecuteMovement` sukses |

i18n en/id di `internal/modules/notification/i18n.go` + FE deep-link. (Pola sama: `OVERTIME_*`, `LEAVE_*`.)

---

# 8. Approval Integration (sudah ada, catatan)

- `SubmitMovement` → `approvalEngine.CreateApprovalInstance("employeemovement", id, flowID)`.
- `RegisterStatusHandler("employeemovement", ...)` di `main.go` → `HandleApprovalStatusChange`:
  - `APPROVED` → movement `approved` (+`approved_at`)
  - `REJECTED` → movement **`rejected`** (status baru §11.4 — bukan lagi `cancelled`)
  - `CANCELLED` → movement `cancelled`
- Setelah approve, eksekusi tetap manual oleh HR via `execute` (§11.1).

---

# 9. Urutan Implementasi

| # | Item | Area | Ketergantungan |
|---|---|---|---|
| 1 | **Migration + enum status `rejected`** (mysql+postgres) + model `MovementStatusRejected` — ✅ **SELESAI (2026-08-10, migration 082)** | BE | — |
| 2 | G-1 ExecuteMovement transaksi employment + adapter (termasuk employee `is_active=false` utk offboarding/retirement) — ✅ **SELESAI (2026-08-10)** | BE | 1 |
| 3 | G-3 auto-resolve flow (`GetActiveFlowIDForModule`) — ✅ **SELESAI (2026-08-10)** | BE | approval engine |
| 4 | G-4 enriched responses (nama employee/org/posisi/status) — ✅ **SELESAI (2026-08-10)** | BE | — |
| 5 | G-2 notifikasi `MOVEMENT_*` (i18n + wiring; REJECTED → status `rejected`) — ✅ **SELESAI (2026-08-10)** | BE | 1 |
| 6 | G-5 hapus endpoint approve manual + service `ApproveMovement` + test — ✅ **SELESAI (2026-08-10)** | BE | — |
| 7 | G-6 contract extension count — ✅ **SELESAI (2026-08-10)** | BE | — |
| 8 | G-7 validasi bisnis per tipe — ✅ **SELESAI (2026-08-10)** | BE | — |
| 9 | G-8 samakan slug module & route FE dengan menu server — ✅ **SELESAI (2026-08-10)** | BE/FE | — |
| 10 | FE: halaman Movements (`/admin/career/movements`) — daftar enriched + filter (type/status/search) + form create per tipe + aksi submit/execute/cancel + delete — ✅ **SELESAI (2026-08-10)** | FE | 2-7, 9 |
| 11 | FE: halaman Contracts terpisah (`/admin/career/contracts`) — daftar enriched + filter status/search + dialog create/edit + upload dokumen + delete — ✅ **SELESAI (2026-08-10)** | FE | 7 |
| 12 | FE: aksi submit/execute/cancel + detail + deep-link notifikasi + badge `rejected` — ✅ **SELESAI (2026-08-10; deep-link & badge sudah masuk langkah 5/10, detail dialog baru)** | FE | 5 |
| 13 | Test: unit/service + FE build + verifikasi manual E2E — ✅ **unit/service + build SELESAI** (PASS di tiap langkah); ✅ **checklist E2E dibuat** (`docs/module-movement-e2e-checklist.md`); ⏳ **eksekusi manual** menunggu environment tenant | — | semua |

---

# 10. Testing Plan

- **Backend**: unit/service transisi state movement (`draft → pending_approval → approved → executed`, `approved → rejected`, `draft → cancelled`); `ExecuteMovement` benar-benar insert employment baru + tutup lama (`effective_date - 1`); offboarding/retirement → employee `is_active=false`; approve manual sudah tidak ada; contract extend (extension_count berantai, previous → extended); validasi per tipe; auto-resolve flow.
- **Integration**: instance approval + push-callback + notifikasi `MOVEMENT_*` (termasuk REJECTED → status `rejected`); eksekusi promosi mengubah employment employee; eksekusi manual oleh HR (§11.1).
- **Frontend**: build bersih, verifikasi manual alur HR (buat draft → submit → approve → execute → cek employment employee) di halaman Movements, dan CRUD kontrak di halaman Contracts terpisah.

---

# 11. Keputusan Bisnis — ✅ DIKONFIRMASI (2026-08-10)

| # | Pertanyaan | Keputusan | Dampak Implementasi |
|---|---|---|---|
| 1 | Eksekusi otomatis vs manual | **Manual oleh HR** (tombol Execute; scheduler opsional di fase berikutnya) | `POST /movements/:id/execute` dipanggil HR; tidak ada job scheduler di fase ini |
| 2 | Employment di masa depan | **Ya, boleh** — employment baru aktif sejak `effective_date`, yang lama ditutup `effective_date - 1` | `ExecuteMovement` set `effective_date` apa adanya (boleh > hari ini); validasi end-date |
| 3 | Offboarding/retirement | **Tandai employee non-aktif** (`is_active = false`) + tutup employment aktif | Adapter `employees` perlu update `is_active`; tidak membuat employment baru |
| 4 | Status `rejected` | **Tambahkan status terpisah `rejected`** | Migration + enum model `MovementStatusRejected`; `HandleApprovalStatusChange` REJECTED → status `rejected`; FE badge/color baru |
| 5 | Approve manual (G-5) | **Hapus** — paksa lewat `submit` (satu pintu approval) | Hapus `POST /movements/:id/approve` + `ApproveMovement` service + test-nya |
| 6 | Navigasi (G-8) | **Dua halaman terpisah** seperti menu server | `EmployeeMovements.vue` (route `/admin/career/movements`) + `EmployeeContracts.vue` baru (route `/admin/career/contracts`); samakan module slug |

---

# 12. Enhancement Plan — Employee Movement & Career Management

> Bagian ini merupakan enhancement lanjutan terhadap implementasi existing. Tidak menggantikan keputusan bisnis pada §11. Enhancement diprioritaskan pada integritas histori employment, validasi movement, career management, dan integrasi dengan modul HRIS lainnya.

## 12.1 Prinsip Enhancement

1. `employee_movements` menjadi sumber histori transaksi perubahan status/posisi employee.
2. `employments` tetap menjadi sumber current employment.
3. Movement tidak boleh mengubah histori lama secara destruktif.
4. Setiap execution harus atomic/transactional.
5. Future-dated movement tetap diperbolehkan sesuai keputusan §11.2.
6. Karena model organisasi adalah **organization = position** dan satu position hanya ditempati satu employee, movement harus melakukan conflict detection.
7. Approval tetap menggunakan Central Approval Module.
8. Performance, competency, dan career path menjadi sumber informasi/eligibility; movement tetap menjadi transaksi eksekusi.
9. Seluruh ID menggunakan UUID/CHAR(36) sesuai pola database existing.
10. Dokumen, audit, dan histori harus tetap tersedia setelah movement dieksekusi.

---

## 12.2 P0 — Transactional Execute Movement

> ✅ **SELESAI (2026-08-10)** — `ExecuteMovement` kini berjalan dalam satu DB transaction (`Repository.ExecuteMovementTx`): conflict detection + close old employment + create new employment + update movement `executed` + commit; kegagalan di langkah mana pun → ROLLBACK (employment lama utuh, movement tetap `approved`, retry oleh HR). Lihat log §3.17.

### Problem

`ExecuteMovement` mengubah employment lama dan membuat employment baru. Proses tersebut harus atomic agar tidak terjadi kondisi sebagian berhasil.

### Target Flow

```text
Execute Movement
      ↓
BEGIN TRANSACTION
      ↓
Validate current state
      ↓
Validate position / organization conflict
      ↓
Close old employment
      ↓
Create new employment
      ↓
Update movement = executed
      ↓
Write audit
      ↓
COMMIT
```

Jika salah satu proses gagal:

```text
ROLLBACK
```

### Acceptance Criteria

- Tidak ada employment setengah berubah.
- Movement hanya menjadi `executed` setelah seluruh perubahan berhasil.
- Jika insert employment gagal, employment lama tetap utuh.
- Audit hanya dibuat setelah transaksi berhasil atau menggunakan transactional audit.

---

## 12.3 P0 — Position / Organization Conflict Detection

> ✅ **SELESAI (2026-08-10)** — `PositionConflict` dicek saat create/update (draft) dan **diulang atomik di dalam transaksi execute**: target position tidak boleh terisi employment terbuka employee lain pada effective date. Konflik → error `CONFLICT_ERROR` (409) `MovementConflictError`. Lihat log §3.17.

Dengan konsep bisnis:

```text
Organization = Position
1 Position = 1 Employee
```

maka sebelum movement dieksekusi harus dilakukan validasi target.

### Validation

```text
Employee A
    ↓
Target Position B
    ↓
Check active employment
    ↓
Position B occupied?
```

Jika sudah ditempati:

```text
❌ Position already occupied
```

### Exception

Conflict dapat dilewati hanya jika movement tersebut secara eksplisit memang melakukan perpindahan employee lama terlebih dahulu dalam satu transaksi bisnis yang valid.

### Acceptance Criteria

- Tidak ada dua employee aktif pada position yang sama.
- Validation dilakukan saat create/submit dan diulang saat execute.
- Execute tidak boleh bergantung hanya pada validation frontend.

---

## 12.4 P0 — Effective Date Conflict Detection

> ✅ **SELESAI (2026-08-10)** — `EmploymentEffectiveDateConflict` (di dalam transaksi execute): employee tidak boleh memiliki employment terbuka yang mulai pada/ setelah effective date — mencegah backdate ke employment aktif & tabrakan dengan employment future-dated dari eksekusi sebelumnya. Konflik → 409 `CONFLICT_ERROR`. Lihat log §3.17.

Movement future-dated harus dapat digunakan, tetapi histori employment tidak boleh overlap.

Contoh:

```text
Movement A
Promotion
Effective: 2026-09-01
```

kemudian:

```text
Movement B
Mutation
Effective: 2026-08-15
```

Sistem harus mendeteksi apakah kedua movement menghasilkan employment period yang conflict.

### Rule

Untuk employee yang sama:

```text
employment_from <= employment_to
```

dan tidak boleh terdapat dua active employment period yang overlap.

### Acceptance Criteria

- Future movement diperbolehkan.
- Movement dengan effective date conflict ditolak.
- Recalculation/validation menggunakan tanggal efektif, bukan tanggal approval.

---

## 12.5 P0 — Movement Snapshot

Saat ini movement menyimpan foreign key `from_*` dan `to_*`. Karena nama Organization, Position, dan Employment Status dapat berubah, histori movement perlu memiliki snapshot.

### Recommended Fields

Tambahkan pada `employee_movements`:

```text
from_organization_name
from_position_name
from_employment_status_name

to_organization_name
to_position_name
to_employment_status_name
```

Jika diperlukan:

```text
from_employment_type_name
to_employment_type_name
```

### Tujuan

Jika Position berubah nama pada tahun berikutnya, histori movement tetap menampilkan nama saat transaksi dibuat.

### Acceptance Criteria

- Snapshot diisi saat movement dibuat/submit.
- Snapshot tidak berubah ketika master data berubah.
- Foreign key tetap disimpan untuk relasi dan navigasi.

> ✅ **SELESAI (2026-08-10)** — migration `083_employeemovement_snapshot` (mysql + postgres, up + down) menambah 6 kolom `from_*/to_*_name` di `employee_movements`. `fillMovementSnapshot` (service) meresolve nama Organization/Position/EmploymentStatus saat create/update dan mempersist-nya bersama movement; `ToResponse` membaca snapshot dari row; enrichment G-4 kini hanya mengisi nama kosong (fallback untuk movement lama tanpa snapshot) — master data yang diubah setelahnya tidak menulis ulang histori. Lihat log §3.18.

---

## 12.6 P0 — Movement Audit Trail

Employee Movement adalah transaksi HR yang penting dan harus dapat diaudit.

### Recommended Table

```text
employee_movement_audits
```

### Fields

| Field | Description |
|---|---|
| `id` | UUID |
| `movement_id` | Employee movement |
| `action` | Action yang dilakukan |
| `old_status` | Status sebelumnya |
| `new_status` | Status setelah action |
| `old_data` | JSON snapshot sebelum perubahan |
| `new_data` | JSON snapshot setelah perubahan |
| `reason` | Alasan |
| `acted_by` | User |
| `acted_at` | Timestamp |

### Actions

```text
CREATED
UPDATED
SUBMITTED
APPROVED
REJECTED
CANCELLED
EXECUTED
```

### Acceptance Criteria

Semua perubahan lifecycle movement tercatat.

> ✅ **SELESAI (2026-08-10)** — tabel `employee_movement_audits` (migration `084_employeemovement_audit`, mysql + postgres) + `recordAudit` di service mencatat CREATED / UPDATED / SUBMITTED / APPROVED / REJECTED / CANCELLED / EXECUTED dengan `old/new_status` + snapshot JSON `old/new_data`. Endpoint `GET /movements/:id/audits`. Lihat log §3.19.

---

## 12.7 P1 — Future-Dated Movement Processing

Keputusan §11.2 tetap dipertahankan:

> Employment masa depan diperbolehkan dan execution dilakukan manual oleh HR.

Enhancement yang disarankan adalah scheduler **opsional**, bukan menggantikan manual execution pada fase existing.

### Opsi A — Manual

```text
APPROVED
   ↓
HR Execute
   ↓
EXECUTED
```

### Opsi B — Scheduled Execution

```text
APPROVED
   ↓
effective_date reached
   ↓
ProcessEffectiveMovementsJob
   ↓
EXECUTED
```

Implementasi scheduler sebaiknya menjadi feature flag/configuration jika nanti diperlukan.

---

## 12.8 P1 — Career Timeline — ✅ **SELESAI (2026-08-10)** — lihat log §3.21

Tambahkan career timeline pada Employee Detail.

### Endpoint

```http
GET /api/v1/tenant/employees/{employeeId}/career-history
```

### Timeline

```text
2024
  Joined
  Staff IT

2025
  Mutation
  IT → Finance

2026
  Promotion
  Staff → Supervisor
```

### Data Source

Tidak perlu membuat `employee_career_history` jika informasi sudah dapat dibentuk dari:

```text
employee_movements
        +
employments
        +
employee_contracts
```

Movement menjadi sumber transaksi, sedangkan API career-history menjadi read model/query khusus.

---

## 12.9 P1 — Career Path

Karena module mencakup Career Management, tambahkan konfigurasi career path.

### New Table

```text
career_paths
career_path_steps
```

### Structure

```text
career_paths
    │
    └── career_path_steps
             │
             ├── position_id
             ├── sequence
             ├── minimum_service_months
             └── requirements
```

### Example

```text
Staff
  ↓
Senior Staff
  ↓
Supervisor
  ↓
Manager
  ↓
Senior Manager
```

### Catatan

Career Path adalah **planning/configuration**, bukan movement transaction.

---

## 12.10 P1 — Promotion Eligibility

Promotion tidak cukup hanya memvalidasi `to_position_id`.

Tambahkan eligibility rule jika business requirement sudah ditetapkan.

### Contoh

```text
Minimum Service       >= 24 months
Performance Score     >= 80
Competency Level      >= 3
Required Training     = completed
```

### Flow

```text
Create Promotion
       ↓
Check Eligibility
       ↓
Eligible?
   ┌───┴───┐
  Yes      No
   ↓        ↓
Submit    Reject
```

### New Table — Optional

```text
career_path_requirements
```

atau dapat dikembangkan sebagai konfigurasi `career_path_steps` jika rule masih sederhana.

---

## 12.11 P1 — Performance Integration

Promotion dapat menggunakan hasil Performance Management.

Contoh:

```text
Employee
   │
   ├── KPI Score       87
   ├── OKR Score       91
   ├── Competency      85
   │
   ▼
Promotion Eligibility
   │
   ▼
Eligible
```

### Integration Principle

Employee Movement **tidak menghitung KPI/OKR**.

Movement hanya membaca hasil final dari:

```text
Performance Management
Competency Management
```

sebagai input eligibility/recommendation.

---

## 12.12 P1 — Mutation Enhancement

Mutation harus mendukung perpindahan lengkap:

```text
FROM
Organization A
Position A

        ↓

TO
Organization B
Position B
```

Karena:

```text
Organization = Position
```

maka validasi target wajib memeriksa pasangan:

```text
organization_id
position_id
```

### Acceptance Criteria

- Organization dan Position harus konsisten.
- Target position harus tersedia.
- Tidak boleh membuat dua active employment pada target position.

---

## 12.13 P1 — Contract Expiry Management — ✅ **SELESAI (2026-08-10)** — lihat log §3.22

Existing `employee_contracts` sudah mendukung `previous_contract_id`, `extension_count`, status, dan document URL.

Tambahkan scheduled process:

```text
ProcessContractExpiration
```

### Reminder

```text
30 days before
14 days before
7 days before
1 day before
```

### Status Transition

```text
ACTIVE
  ↓ end_date reached
EXPIRED
```

### Acceptance Criteria

- Contract expired otomatis dapat terdeteksi.
- Notification dikirim kepada HR.
- Employee/manager dapat menerima reminder sesuai permission.

---

## 12.14 P1 — Contract Extension Chain

Pertahankan chain existing:

```text
Contract A
    ↓
extension
Contract B
    ↓
extension
Contract C
```

Dengan:

```text
previous_contract_id
extension_count
```

### Enhancement

Pastikan setiap extension:

- Memiliki contract number baru jika policy mengharuskan.
- Memiliki effective date baru.
- Menyimpan previous contract.
- Tidak mengubah histori contract sebelumnya.
- Dapat ditelusuri dari contract terbaru ke awal chain.

---

## 12.15 P1 — Movement Documents — ✅ **SELESAI (2026-08-10)** — lihat log §3.20

Saat ini movement memiliki informasi SK seperti:

```text
decision_letter_number
decision_letter_date
```

Untuk mendukung lebih dari satu dokumen, tambahkan:

```text
employee_movement_documents
```

### Fields

| Field | Description |
|---|---|
| `id` | UUID |
| `movement_id` | Movement |
| `document_type` | Jenis dokumen |
| `file_name` | Nama file |
| `file_url` | Lokasi file |
| `uploaded_by` | User |
| `created_at` | Timestamp |

### Document Type Example

```text
PROMOTION_SK
MUTATION_SK
DEMOTION_SK
RETIREMENT_LETTER
OFFBOARDING_LETTER
OTHER
```

---

## 12.16 P1 — Movement Cancellation After Approval

Movement yang sudah approved tidak sebaiknya dapat dibatalkan secara langsung tanpa audit/approval tambahan jika policy HR mensyaratkannya.

Recommended future flow:

```text
APPROVED
   ↓
Cancellation Request
   ↓
Central Approval Module
   ↓
CANCELLED
```

Namun keputusan §11.1 tetap berlaku untuk execution:

```text
Approved → HR Execute
```

---

## 12.17 P2 — Movement Reporting

Tambahkan report:

### Movement Report

```text
Promotion
Demotion
Mutation
Contract Extension
Status Change
Retirement
Offboarding
```

Filter:

```text
Period
Organization
Position
Employee
Movement Type
Status
```

### Career History Report

```text
Employee
Join Date
Position History
Organization History
Promotion History
Mutation History
Status History
```

### Contract Report

```text
Active
Expiring
Expired
Extended
Terminated
```

---

## 12.18 P2 — Dashboard

### HR Dashboard

```text
Employee Movement
-------------------------
Promotion             12
Mutation              20
Demotion               2
Status Change          8
Retirement             3
Offboarding             5

Pending Approval       10
Effective This Month    8
```

### Contract

```text
Active                150
Expiring < 30 days     12
Expired                 5
```

---

# 13. Recommended Database Enhancement

## Existing Tables

```text
employee_movements
employee_contracts
```

## P0 / P1 New Tables

```text
employee_movement_audits
employee_movement_documents
```

## Career Management

```text
career_paths
career_path_steps
```

## Optional

```text
career_path_requirements
```

Tidak direkomendasikan membuat:

```text
employee_career_history
```

karena career history dapat dibentuk dari existing transactional data.

---

# 14. Recommended Field Enhancement

## `employee_movements`

Tambahkan snapshot fields:

```text
from_organization_name
from_position_name
from_employment_status_name

to_organization_name
to_position_name
to_employment_status_name
```

Jika dibutuhkan:

```text
from_employment_type_name
to_employment_type_name
```

### Indexes

Pastikan tersedia index untuk:

```text
employee_id
movement_type
status
effective_date
from_organization_id
to_organization_id
from_position_id
to_position_id
```

---

# 15. Enhancement API Plan

## Career History

```http
GET /api/v1/tenant/employees/{employeeId}/career-history
```

## Eligibility

```http
GET /api/v1/tenant/employees/{employeeId}/movement-eligibility
GET /api/v1/tenant/employees/{employeeId}/promotion-eligibility
```

## Documents

```http
GET    /api/v1/tenant/employee-movements/{id}/documents
POST   /api/v1/tenant/employee-movements/{id}/documents
DELETE /api/v1/tenant/employee-movements/{id}/documents/{documentId}
```

## Audit

```http
GET /api/v1/tenant/employee-movements/{id}/audits
```

## Career Paths

```http
GET    /api/v1/tenant/career-paths
POST   /api/v1/tenant/career-paths
GET    /api/v1/tenant/career-paths/{id}
PUT    /api/v1/tenant/career-paths/{id}
DELETE /api/v1/tenant/career-paths/{id}
```

---

# 16. Enhancement Service Layer

Recommended services:

```text
MovementExecutionService
MovementValidationService
MovementConflictService
MovementSnapshotService
MovementAuditService
MovementDocumentService
CareerHistoryService
CareerPathService
CareerEligibilityService
ContractExpirationService
```

### Responsibility

```text
MovementValidationService
    ↓
Validation business rules

MovementConflictService
    ↓
Position / employment overlap

MovementExecutionService
    ↓
Atomic employment changes

MovementSnapshotService
    ↓
Historical snapshot

MovementAuditService
    ↓
Lifecycle audit

CareerEligibilityService
    ↓
Promotion / career eligibility
```

---

# 17. Enhancement Testing Plan

## 17.1 Transaction Test

```text
Execute Movement
→ Old Employment Closed
→ New Employment Created
→ Movement Executed
```

Simulasikan failure pada new employment:

```text
Old Employment remains unchanged
Movement remains approved
```

---

## 17.2 Position Conflict Test

```text
Position A occupied by Employee A

Employee B
→ Movement to Position A

Expected:
REJECTED / VALIDATION ERROR
```

---

## 17.3 Future Date Test

```text
Effective Date > Today
```

Expected:

```text
Approved
but employment only effective from effective_date
```

---

## 17.4 Overlap Test

Test:

```text
Movement A = 01 Jan
Movement B = 15 Jan
```

Pastikan employment periods tidak overlap secara invalid.

---

## 17.5 Snapshot Test

1. Create movement.
2. Snapshot position name = `Staff`.
3. Rename position menjadi `Senior Staff`.
4. Movement history tetap menampilkan `Staff`.

---

## 17.6 Career History Test

Pastikan timeline menampilkan:

```text
Join
→ Promotion
→ Mutation
→ Promotion
→ Current Position
```

secara kronologis.

---

## 17.7 Contract Test

Test:

```text
Contract A
→ Extension
→ Contract B
→ Extension
→ Contract C
```

Pastikan:

```text
A.previous = null
B.previous = A
C.previous = B
```

dan:

```text
extension_count
A = 0
B = 1
C = 2
```

---

# 18. Enhancement Development Order

| # | Item | Priority | Area |
|---|---|---|---|
| 1 | Transactional Execute Movement — ✅ **SELESAI (2026-08-10)** | P0 | BE |
| 2 | Position Conflict Detection — ✅ **SELESAI (2026-08-10)** | P0 | BE |
| 3 | Effective Date Conflict — ✅ **SELESAI (2026-08-10)** | P0 | BE |
| 4 | Movement Snapshot — ✅ **SELESAI (2026-08-10)** | P0 | DB/BE |
| 5 | Movement Audit Trail — ✅ **SELESAI (2026-08-10)** | P0 | DB/BE |
| 6 | Movement Documents — ✅ **SELESAI (2026-08-10)** | P1 | DB/BE/FE |
| 7 | Career Timeline — ✅ **SELESAI (2026-08-10)** | P1 | BE/FE |
| 8 | Contract Expiry — ✅ **SELESAI (2026-08-10)** | P1 | BE/Job/FE |
| 9 | Performance Integration | P1 | BE |
| 10 | Career Path | P1 | DB/BE/FE |
| 11 | Promotion Eligibility | P1/P2 | BE |
| 12 | Movement Cancellation Approval | P2 | Approval/BE |
| 13 | Reports | P2 | BE/FE |
| 14 | Dashboard | P2 | BE/FE |

---

# 19. Final Target Architecture

```text
                         Employee Movement
                                │
              ┌─────────────────┼─────────────────┐
              │                 │                 │
              ▼                 ▼                 ▼
         Validation          Approval          Documents
              │                 │                 │
              │                 ▼                 │
              │        Central Approval          │
              │                 │                 │
              └─────────────────┼─────────────────┘
                                ▼
                         Execute Movement
                                │
                         DB Transaction
                                │
                 ┌──────────────┴──────────────┐
                 ▼                             ▼
        Old Employment                  New Employment
                 │                             │
                 └──────────────┬──────────────┘
                                ▼
                       Employee Career History
                                │
            ┌───────────────────┼───────────────────┐
            ▼                   ▼                   ▼
       Performance          Competency         Career Path
          KPI/OKR                                  │
            │                                      ▼
            └──────────────────────────── Promotion
                                                   │
                                                   ▼
                                             Employee Movement
```

---

# 20. Final Design Principles

1. `employee_movements` adalah **transaction/history**, bukan current employee state.
2. `employments` adalah sumber current/future employment state.
3. Movement execution wajib atomic.
4. Target Position harus divalidasi karena satu Position hanya boleh memiliki satu employee aktif sesuai konsep HRIS.
5. Future-dated movement diperbolehkan.
6. Approval tetap menggunakan Central Approval Module.
7. `rejected` berbeda dari `cancelled`.
8. Movement yang approved tetap memerlukan execution sesuai keputusan §11.
9. Snapshot diperlukan untuk menjaga histori terhadap perubahan master Organization/Position/Status.
10. Audit trail wajib tersedia untuk lifecycle movement.
11. Career history tidak perlu memiliki tabel duplikasi jika dapat dibentuk dari movement + employment.
12. Career Path merupakan konfigurasi/perencanaan karier, bukan transaksi movement.
13. Performance KPI/OKR dapat menjadi input promotion eligibility, bukan dihitung di Movement Module.
14. Contract extension harus mempertahankan chain melalui `previous_contract_id`.
15. Contract expiration harus dapat dimonitor dan diberi notification.
16. Dokumen movement harus mendukung lebih dari satu dokumen.
17. Semua enhancement tetap mengikuti pola UUID yang sudah digunakan project.
18. Semua perubahan database harus memiliki migration MySQL dan PostgreSQL jika project mempertahankan dua driver tersebut.
19. Business rules berada di service layer, bukan controller/frontend.
20. Frontend hanya menampilkan action yang sesuai permission dan state movement.
