# Approval Module — Central Workflow Engine Integration Plan

## Objective

Menjadikan modul `approval` (`backend/internal/modules/approval/`) sebagai satu-satunya sumber kebenaran untuk seluruh alur persetujuan di platform:

- Approval flow hanya bisa dikonfigurasi untuk module yang **disubscribe/diaktifkan** tenant (bukan free-text tanpa validasi)
- `SUPERVISOR` sebagai approver type benar-benar bisa di-resolve (saat ini masih stub)
- Semua module yang butuh approval (leave, reimbursement, employee movement, overtime, payroll, dst) merujuk ke satu approval engine yang sama — tidak lagi punya logic approve/reject sendiri-sendiri
- Approve/reject selalu lewat endpoint approval module (`POST /approval/instances/:id/action`), bukan endpoint masing-masing module

---

## Status Modul Saat Ini

Modul `approval` **sudah ada dan cukup lengkap**: model, DTO, repository, service (CRUD flow/step, `CreateInstance`, `SubmitAction` dengan mode ANY_ONE/ALL/N_OF_M, `ListMyPendingTasks`), handler, routes, dan `module.go` — sudah cocok dengan skema 5 tabel di migration `010_approval.sql` (mysql & postgres, identik field-per-field).

**Yang belum:**

| Gap | Lokasi | Dampak |
|-----|--------|--------|
| `SUPERVISOR` approver type belum di-resolve | `approval/service.go:689-694` — `resolveStepAssignees` return error stub | Flow dengan approver type SUPERVISOR tidak bisa jalan sama sekali |
| Tidak ada validasi module subscription pada flow | `approval_flows.module` adalah `VARCHAR` bebas | Tenant bisa bikin flow untuk module yang belum disubscribe |
| Tidak ada helper `IsModuleActive` | `platform/modulemgmt/` | Tidak ada cara cek module aktif dari modul lain |
| Hanya `payroll` yang terintegrasi, dan itu pun tidak sempurna | `payroll/service.go:913-960`, `main.go:411-494` | `approval_instance_id` tidak pernah disimpan (cuma di-log); sinkronisasi status manual (pull), bukan callback (push) |
| 4 module masih pakai status enum & approve/reject sendiri, tidak nyambung ke approval engine | `leave`, `reimbursement`, `employeemovement`, `attendance` (overtime) | Duplikasi logic, tidak konsisten, tidak auditable lewat satu tempat |

### Detail per module ad-hoc

| Module | Status enum | Approve/reject logic |
|--------|-------------|----------------------|
| `leave` | `LeaveStatus` (`leave/model.go:118-124`) | `SupervisorID`/`SupervisorActionAt`/`SupervisorNote` bolted-on di request; satu method `UpdateLeaveRequestStatus` (`service.go:469`) |
| `reimbursement` | `ReimbursementStatus` (`model.go:40-46`) | Switch-driven di `service.go:254-279` |
| `employeemovement` | `MovementStatus` (`model.go:25-29`) | `ApproveMovement` (`service.go:281`) |
| `attendance` (overtime) | `OvertimeStatus` (`model.go:298-303`) | Belum ada method approve/reject terpisah |

---

## Pola Integrasi yang Sudah Terbukti (payroll ↔ approval)

Dari `backend/cmd/server/main.go`:

```
main.go:411-417   → approvalRepo, approvalSvc (*approval.Service) dibuat sekali
main.go:420-422   → payrollApprovalAdapter{approvalSvc} — adapter membungkus approval.Service
                     agar cocok dengan interface sempit milik payroll
main.go:482       → approval.NewModuleWithService(dbManager, l, approvalSvc) — mount HTTP routes,
                     reuse instance service yang sama
main.go:487       → payroll.NewModule(dbManager, l, payrollApprovalEngine) — adapter di-inject
```

`payroll/service.go:20-28` mendefinisikan interface sempit sendiri (`ApprovalEngine` dengan `CreateApprovalInstance`/`GetApprovalInstanceStatus`) — payroll tidak import `approval.Service` langsung.

**Pola ini dipertahankan** (interface sempit + adapter + injection) karena sudah bersih dan terbukti jalan. Yang diperbaiki: instance ID harus disimpan, dan sinkronisasi status harus push (callback) bukan pull manual.

---

# Phase 1 — Approver Resolution Berbasis Hierarki Organisasi

## Objective
Membuat approver resolution mengacu ke hierarki Organization, dengan rule:

- **Atasan langsung** — approver = employee yang menempati parent Organization dari submitter (1 level ke atas)
- **Bisa pilih organization** — step bisa ditentukan approver-nya dari Organization tertentu yang dipilih eksplisit, tidak relatif ke submitter
- **Bisa pilih keduanya** — satu step bisa mewajibkan approval dari atasan hierarki DAN Organization yang dipilih eksplisit
- **Multiple level** — approval bisa berjenjang lebih dari 1 level (atasan langsung → atasan dari atasan → dst)
- **Di satu level bisa lebih dari satu yang approve** — satu step bisa punya lebih dari satu approver

Platform sudah punya cukup data untuk ini **tanpa field baru di `Employee`**: konvensi "Organization = Position = 1 Employee" sudah ada (`organization/model.go:19` `Organization.ParentID` untuk hierarki, `employee/model.go:279` `Employment.OrganizationID` untuk employee yang menempati suatu Organization). Atasan langsung = employee yang menempati parent Organization milik submitter.

## Desain

`approver_type` diperluas (bukan diganti — `ROLE`/`USER` tetap):

| `approver_type` | Resolusi |
|---|---|
| `SUPERVISOR` (redefinisi) | Naik `hierarchy_level` kali dari Organization submitter lewat `organizations.parent_id`, resolve ke employee yang menempati Organization hasil naik tersebut (`employments.organization_id`) |
| `ORGANIZATION` (baru) | Resolve ke employee yang menempati Organization-Organization yang dipilih eksplisit untuk step ini |
| `BOTH` (baru) | Gabungan (union) hasil resolusi `SUPERVISOR` + `ORGANIZATION` untuk step yang sama (dedup) |

**Perubahan skema** (migration baru, mysql + postgres):
- `approval_flow_steps.hierarchy_level INT NULL DEFAULT 1` — jumlah level naik (1 = atasan langsung, 2 = atasan dari atasan, dst); relevan untuk `SUPERVISOR`/`BOTH`
- Tabel baru `approval_flow_step_organizations (id, step_id FK, organization_id)` — satu step bisa punya banyak Organization target (`ORGANIZATION`/`BOTH`) — inilah yang mewujudkan "lebih dari satu approver di satu level". Digabung dengan `approval_mode`/`required_approvals` yang **sudah ada** (ANY_ONE/ALL/N_OF_M) di `approval_flow_steps`, tidak perlu perubahan di situ.

**"Multiple level"** cukup direpresentasikan sebagai beberapa `ApprovalFlowStep` (step_order berurutan), masing-masing resolve dari Organization submitter yang sama dengan `hierarchy_level` masing-masing — reuse mekanisme step sequencing yang sudah ada, tidak perlu konsep baru.

## Perubahan

**`approval` module:**
- `service.go` — implementasikan `resolveStepAssignees` sesuai tabel resolusi di atas (ganti stub `fmt.Errorf`). Perlu tahu Organization submitter saat ini — resolve dari `ApprovalInstance.CreatedBy` (user) → employee → `Employment.OrganizationID` yang aktif. `approval` tidak boleh import `employee` langsung (circular dependency, sama seperti alasan modul `performance` pakai raw query `db.Table("organizations")`) — pakai raw query ke tabel `employees`/`employments`/`organizations`.
- **Item terbuka yang perlu dikonfirmasi saat implementasi**: bagaimana `employees` terhubung ke `platform_users`/`ApprovalInstance.CreatedBy` (belum ditemukan FK `user_id` langsung di `Employee` dari riset awal) — perlu untuk dari "siapa yang submit" ke "employee/organization mana". Pastikan join ini sebelum menulis `resolveStepAssignees`.
- `dto.go` — `CreateStepRequest`/`UpdateStepRequest`: tambah `hierarchy_level *int` dan `organization_ids []string`, update binding `approver_type` jadi `oneof=SUPERVISOR ROLE USER ORGANIZATION BOTH`
- `repository.go` — CRUD untuk `approval_flow_step_organizations` mengikuti CRUD step (create/update/delete step harus sinkron dengan daftar organization)

---

# Phase 2 — Module-Subscription Awareness

## Objective
Approval flow hanya bisa dikonfigurasi untuk module yang aktif untuk tenant tersebut.

## Perubahan

**`modulemgmt`:**
- Tambah `IsModuleActive(ctx, companyID uuid.UUID, moduleSlug string) (bool, error)` di service, berbasis tabel `company_modules` (`modulemgmt/model.go:37-46`)

**`approval`:**
- Inject `IsModuleActive` (adapter pattern yang sama) ke `approval.Service`
- Validasi di `CreateFlow`/`UpdateFlow` — tolak jika `module` belum disubscribe tenant
- Endpoint baru `GET /approval/available-modules` — daftar module aktif tenant yang mendukung approval, dipakai frontend flow builder (module picker, bukan free-text)

---

# Phase 3 — Push-Based Status Sync

## Objective
Saat approval instance mencapai status final, module konsumen otomatis ter-update — tidak perlu ada pihak lain yang manual memanggil endpoint status module tersebut.

## Perubahan

**`approval/service.go`:**
- Tambah callback registry:
  ```go
  RegisterStatusHandler(module string, fn func(ctx context.Context, documentID uuid.UUID, status InstanceStatus, note string) error)
  ```
- Panggil handler yang terdaftar secara sinkron di dalam `SubmitAction` begitu instance mencapai status final (`APPROVED`/`REJECTED`/`CANCELLED`)

**Retrofit `payroll`:**
- Tambah kolom `approval_instance_id` di `payroll_runs` (migration baru)
- Simpan instance ID di `UpdatePayrollRunStatus` (saat ini cuma di-log, `service.go:933-937`)
- Daftarkan status handler payroll saat startup, hapus ketergantungan pada caller eksternal yang manual set `status=APPROVED`

---

# Phase 4 — Migrasi 4 Module Ad-hoc

Untuk masing-masing `leave`, `reimbursement`, `employeemovement`, `attendance` (overtime) — **satu per satu**, mulai dari `leave` (paling self-contained, cocok untuk validasi end-to-end resolusi hierarki dari Phase 1). Field ad-hoc `SupervisorID`/`SupervisorActionAt`/`SupervisorNote` (`leave/model.go:141-143`) jadi redundan begitu approval yang resolve supervisor-nya sendiri — bisa dihapus atau dipertahankan read-only untuk histori:

1. Tambah kolom + field `approval_instance_id` di model request
2. Saat submit: panggil `approvalEngine.CreateApprovalInstance(ctx, "<module>", requestID, flowID)` via adapter pattern yang sama seperti payroll
3. Daftarkan status handler (Phase 3) yang meng-update status field module tersebut saat approval resolve
4. **Approve/reject action pindah ke endpoint approval module** (`POST /approval/instances/:id/action`) — endpoint approve/reject milik module masing-masing dihapus
5. Nilai enum status yang sudah ada di masing-masing module **dipertahankan** (untuk kompatibilitas frontend) — yang berubah hanya mekanisme yang men-set nilainya

---

# Phase 5 — Dokumentasi & Frontend

- Update dokumen ini dengan status implementasi setiap phase
- Frontend:
  - Flow builder: module picker dibatasi ke module yang disubscribe (dari Phase 2), step builder dengan approver type Supervisor (+ hierarchy level)/Organization (multi-select)/Both/Role/User
  - "My Pending Tasks" inbox
  - Halaman submit di masing-masing module (leave, reimbursement, dst) tidak lagi menampilkan tombol approve/reject sendiri — pindah ke approval inbox

---

## File yang Kemungkinan Tersentuh

| Area | File |
|------|------|
| Approver resolution berbasis hierarki organisasi | Migration baru (mysql+postgres: kolom `hierarchy_level` + tabel `approval_flow_step_organizations`), `backend/internal/modules/approval/model.go`, `dto.go`, `repository.go`, `service.go` |
| Module subscription check | `backend/internal/platform/modulemgmt/service.go`, `backend/internal/modules/approval/service.go`, `handler.go`, `routes.go` |
| Status callback registry | `backend/internal/modules/approval/service.go` |
| Payroll retrofit | `backend/internal/modules/payroll/model.go`, `service.go`, migration baru, `backend/cmd/server/main.go` |
| Migrasi Leave/Reimbursement/EmployeeMovement/Attendance | `model.go`, `service.go`, `handler.go`, `routes.go` per module, migration per module, `backend/cmd/server/main.go` |

---

## Verifikasi

1. `cd backend && go build ./...` setelah setiap phase
2. `go test ./internal/modules/approval/... ./internal/modules/leave/... ./internal/modules/payroll/... -v`
3. Manual: buat approval flow untuk `leave` dengan step `SUPERVISOR` (`hierarchy_level=1`) → submit leave request → pastikan pending task muncul untuk employee yang menempati parent organization submitter → approve lewat `POST /approval/instances/:id/action` → pastikan status leave request otomatis jadi `APPROVED` tanpa panggilan kedua secara manual
4. Tambah step kedua dengan `approver_type=ORGANIZATION` menunjuk 2 Organization eksplisit, `approval_mode=ALL` → pastikan 2 pending task dibuat dan keduanya wajib approve sebelum step selesai
5. Tambah step dengan `approver_type=BOTH` → pastikan hasilnya gabungan approver dari hierarki + Organization eksplisit
6. Pastikan membuat flow untuk module yang belum disubscribe tenant ditolak (validasi Phase 2)

---

# Implementation Status

| Phase | Status | Completion Date | Notes |
|-------|--------|-----------------|-------|
| Phase 1 - Approver Resolution Berbasis Hierarki Organisasi | ⏳ Pending | - | |
| Phase 2 - Module-Subscription Awareness | ⏳ Pending | - | |
| Phase 3 - Push-Based Status Sync | ⏳ Pending | - | |
| Phase 4 - Migrasi 4 Module Ad-hoc | ⏳ Pending | - | |
| Phase 5 - Dokumentasi & Frontend | ⏳ Pending | - | |
