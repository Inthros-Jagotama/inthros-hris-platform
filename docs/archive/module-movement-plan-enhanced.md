# Employee Movement & Career Management — Development Plan

> ✅ **Diarsipkan (2026-08-21)**: dokumen ini sangat basi terhadap kode aktual — diverifikasi ulang dan modul ini **~90%+ selesai**, jauh melebihi klaim dokumen. Step 12-13 (§9) yang ditandai belum selesai ternyata **sudah selesai** (FE detail dialog + status tag `rejected` di `EmployeeMovements.vue`, deep-link `MOVEMENT_*` di `Notifications.vue`). Seluruh "Enhancement Plan" §12-20 (P0/P1/P2, sebelumnya tanpa status sama sekali) **sudah diimplementasikan sebagian besar**: `ExecuteMovementTx` transaksi DB nyata + position-conflict + effective-date overlap check, migration 083 (snapshot fields persisted), migration 084 (`employee_movement_audits` + `recordAudit`), migration 086 (`career_paths`/`career_path_steps`), `eligibility_test.go` (promotion eligibility), `ProcessContractExpiration` + scheduler harian (§12.13), plus migration 085 (documents) & 087 (cancellation) yang bahkan tidak tercatat di dokumen ini. Hanya nama service persis di §16 (`MovementConflictService`, dst) yang tidak dipakai — logic-nya ada di method `Service`/`Repository` langsung, bukan gap fungsional. Diarsipkan sebagai selesai.
>
> 📅 Versi plan: 2026-08-10 · Status: **IMPLEMENTASI BERJALAN — langkah 3/12 selesai** (backend existing ✅ + 3 langkah baru ✅, FE placeholder ❌)
> ✅ **Keputusan bisnis sudah dikonfirmasi user (2026-08-10)** — lihat §11.
> 🔎 Berdasarkan struktur tabel `012_employee_movement.sql` (mysql + postgres) dan `062_employeemovement_approval_instance.sql`, serta audit modul `backend/internal/modules/employeemovement` dan `frontend/tenant/src/views/modules/EmployeeMovements.vue`.
> 📊 **Progres implementasi (per 2026-08-10):** ✅ 1) migration + enum `rejected` (082) · ✅ 2) G-1 ExecuteMovement transaksi employment · ✅ 3) G-3 auto-resolve flow · ✅ 4) G-4 enriched responses · ✅ 5) G-2 notifikasi `MOVEMENT_*` · ✅ 6) G-5 hapus approve manual · ✅ 7) G-6 contract extension count · ✅ 8) G-7 validasi per tipe · ✅ 9) G-8 slug/route disamakan · ✅ 10) FE halaman Movements + filter backend · ✅ 11) FE halaman Contracts (daftar enriched + filter + create/edit + upload dokumen) + filter backend. **Berikutnya:** 12) FE detail/deep-link/badge · 13) test & verifikasi.

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
| 12 | FE: aksi submit/execute/cancel + detail + deep-link notifikasi + badge `rejected` | FE | 5 |
| 13 | Test: unit/service + FE build + verifikasi manual E2E | — | semua |

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

## 12.8 P1 — Career Timeline

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

## 12.13 P1 — Contract Expiry Management

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

## 12.15 P1 — Movement Documents

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
| 1 | Transactional Execute Movement | P0 | BE |
| 2 | Position Conflict Detection | P0 | BE |
| 3 | Effective Date Conflict | P0 | BE |
| 4 | Movement Snapshot | P0 | DB/BE |
| 5 | Movement Audit Trail | P0 | DB/BE |
| 6 | Movement Documents | P1 | DB/BE/FE |
| 7 | Career Timeline | P1 | BE/FE |
| 8 | Contract Expiry | P1 | BE/Job/FE |
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
