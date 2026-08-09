# Employee Movement & Career Management — Development Plan

> 📅 Versi plan: 2026-08-10 · Status: **IMPLEMENTASI DIMULAI — langkah 1/12 selesai** (backend existing ✅, FE placeholder ❌)
> ✅ **Keputusan bisnis sudah dikonfirmasi user (2026-08-10)** — lihat §11.
> 🔎 Berdasarkan struktur tabel `012_employee_movement.sql` (mysql + postgres) dan `062_employeemovement_approval_instance.sql`, serta audit modul `backend/internal/modules/employeemovement` dan `frontend/tenant/src/views/modules/EmployeeMovements.vue`.

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
- **Repository**: CRUD + `ApproveMovement`/`ExecuteMovement`/`CancelMovement`/`ExtendContract`.
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

## G-5 🟡 ENDPOINT APPROVE MANUAL MASIH ADA

`POST /movements/:id/approve` (manual, tanpa engine) tetap eksis sebagai jalur paralel approval. Dengan Central Approval yang sudah jalan, jalur manual berisiko "dua pintu" yang tidak sinkron.

**Keputusan user §11.5: HAPUS.** `POST /movements/:id/approve` dihapus dari handler/routes; satu-satunya jalur approval = `submit` → Central Approval. (Service method `ApproveMovement` boleh dihapus atau dibiarkan tidak terpakai — lebih baik dihapus bersama test-nya.)

## G-6 🟡 CONTRACT EXTENSION COUNT HARDCODED

`service.go CreateContract` men-set `contract.ExtensionCount = 1` dengan komentar "dihitung manual oleh caller untuk extension > 1" — perpanjangan berantai tidak menghitung dengan benar.

**Rencana:** `ExtendContract` menghitung `extension_count = previous.extension_count + 1`, dan kontrak sebelumnya di-set `status = extended`.

## G-7 🟡 VALIDASI BISNIS PER TIPE MOVEMENT

Belum ada validasi "tipe X wajib field Y":
- `mutation` → wajib `to_organization_id` (dan/atau `to_position_id`)
- `promotion`/`demotion` → wajib `to_position_id`
- `status_change` → wajib `to_employment_status_id`
- `contract_extension` → wajib merujuk kontrak aktif
- `offboarding`/`retirement` → boleh tanpa `to_*`

**Rencana:** validasi service-level + `binding` DTO, kembalikan pesan field error (pola `getValidationErrors`).

## G-8 🟡 KONSISTENSI ROUTE/MENU & MODULE SLUG

- Backend menu (module.go): `/admin/career/movements`, `/admin/career/contracts`
- FE router saat ini: `/employee-movements` (satu halaman, module slug `employee-movement`)
- Approval module slug: `employeemovement` (tanpa tanda hubung)

**Keputusan user §11.6: DUA HALAMAN TERPISAH** mengikuti menu server:
- `/admin/career/movements` → halaman Movements (`EmployeeMovements.vue`)
- `/admin/career/contracts` → halaman Contracts (`EmployeeContracts.vue` baru)
- Router FE + sidebar disesuaikan; module slug disamakan (`employeemovement` atau `employee-movement` — pilih satu, samakan dengan permission & filter approval).

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
| 2 | G-1 ExecuteMovement transaksi employment + adapter (termasuk employee `is_active=false` utk offboarding/retirement) | BE | 1 |
| 3 | G-3 auto-resolve flow (`GetActiveFlowIDForModule`) | BE | approval engine |
| 4 | G-4 enriched responses (nama employee/org/posisi/status) | BE | — |
| 5 | G-2 notifikasi `MOVEMENT_*` (i18n + wiring; REJECTED → status `rejected`) | BE | 1 |
| 6 | G-5 hapus endpoint approve manual + service `ApproveMovement` + test | BE | — |
| 7 | G-6 contract extension count + G-7 validasi per tipe | BE | — |
| 8 | FE: halaman Movements (`/admin/career/movements`) + locale lengkap | FE | 2-7 |
| 9 | FE: halaman Contracts terpisah (`/admin/career/contracts`) + upload dokumen | FE | 7 |
| 10 | FE: aksi submit/execute/cancel + detail + deep-link notifikasi + badge `rejected` | FE | 5 |
| 11 | G-8 samakan slug module & route FE dengan menu server | BE/FE | — |
| 12 | Test: unit/service + FE build + verifikasi manual E2E | — | semua |

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
