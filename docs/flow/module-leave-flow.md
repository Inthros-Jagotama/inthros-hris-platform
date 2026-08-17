# Alur Pengisian Leave / Cuti (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Leave** dari setup master data sampai
request cuti selesai (approval + pemotongan saldo) — pola runbook seperti
[`module-payroll-user-flow.md`](module-payroll-user-flow.md) & [`module-reimbursement-flow.md`](module-reimbursement-flow.md).

- Plan pengembangan: [`../module-leave-plan.md`](../module-leave-plan.md) — 🟡 sebagian (inti selesai; accrual/adjustment/carry-forward/expiry & team calendar belum)
- Lokasi kode: `backend/internal/modules/leave/` · `frontend/tenant/src/views/modules/leave/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Leave)

---

## 1. Ringkasan Alur End-to-End

```
SETUP (sekali)                     PENGAJUAN (per request)                     SALDO
┌──────────────────┐    ┌─────────────────────────────────────────────────┐   ┌────────────┐
│ Leave Types      │    │ DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED  │   │ auto-deduct│
│ Accrual Policies │───▶│   │                        │                     │──▶│ + ledger   │
│ Leave Reasons    │    │   └── CANCELLED            └── REJECTED_FINAL     │   │ (usage)    │
└──────────────────┘    └─────────────────────────────────────────────────┘   └────────────┘
```

- **Status request:** `DRAFT → SUBMITTED → PENDING_APPROVAL → APPROVED_FINAL` · terminal: `REJECTED_FINAL`, `CANCELLED`
- **Approval** ditindaklanjuti di halaman **Approvals** generik (modul `leave`); saat `APPROVED_FINAL`
  saldo otomatis dikurangi (`applyLeaveUsage`) dan tercatat di ledger `leave_balance_transactions`.
- **`requested_days` dihitung server-side** (working-day calc: weekend + company holiday di-exclude)
  — klien cukup kirim tanggal & `duration_mode`.

---

## 2. Entitas Utama

| Entitas | Tabel | Deskripsi |
|---|---|---|
| Leave Type | `leave_types` | Jenis cuti (tahunan, sakit, dll.) + opsi `requires_attachment` |
| Accrual Policy | `leave_accrual_policies` | Aturan quota per type (masih master data — **engine accrual belum dibangun**) |
| Leave Reason | `leave_reasons` | Alasan cuti (opsional) |
| Leave Request | `leave_requests` | Pengajuan cuti + detail per tanggal (`leave_request_details`) |
| Balance | `employee_leave_balances` | Quota/used/remaining per employee+type+tahun |
| Balance Ledger | `leave_balance_transactions` | Riwayat USAGE / REVERSAL |

---

## 3. TAHAP 1 — SETUP (dikerjakan sekali)

### A. Leave Types

Menu **Leave → Types** (`leave/types`, `LeaveTypes.vue`).

- CRUD jenis cuti: kode/nama, kuota, `requires_attachment` (bila aktif → attachment wajib saat request).
- Endpoint: `POST/GET /types`, `GET/PUT/DELETE /types/:id`.

### B. Accrual Policies

Menu **Leave → Accrual Policies** (`leave/accrual-policies`, `LeaveAccrualPolicies.vue`).

- CRUD aturan accrual per type: pilih Leave Type + `effective_from`/`effective_to`.
- ⚠️ Saat ini **master data saja** — engine yang me-generate quota dari policy belum dibangun
  (lihat plan Phase 6 lanjutan).

### C. Leave Reasons

Menu **Leave → Reasons** (`leave/reasons`, `LeaveReasons.vue`).

- CRUD alasan cuti opsional: `POST/GET /reasons`, `GET/PUT/DELETE /reasons/:id`.
- Tanpa pagination server-side (array polos).

---

## 4. TAHAP 2 — PENGAJUAN (setiap request)

### A. Buat Request (DRAFT)

- Dari **halaman Leave** (`/leave`, `Leave.vue`) → tombol **New Request** → dialog:
  - Leave Type (dropdown dari `/types`) · rentang tanggal (`request_start_date`/`request_end_date`)
  - `duration_mode`: **FULL_DAY** / **HALF_DAY** / **HOURLY** (HOURLY → isi `start_time`/`end_time` `HH:mm`)
  - Leave Reason (opsional) + note · `attachment_url` bila type `requires_attachment`
  - `requested_days` **tidak diisi klien** — backend menghitung (working-day calc).
- Validasi backend: leave type aktif, attachment wajib bila disyaratkan, overlap tanggal.

### B. Submit (SUBMITTED → PENDING_APPROVAL)

- Tombol **Submit** → `PUT /requests/:id/status` `{"status":"SUBMITTED"}`:
  - Membuat **instance approval** di Central Approval (modul `leave`); flow aktif di-auto-resolve.
  - Status jadi `SUBMITTED`/`PENDING_APPROVAL` (tergantung state machine).

### C. Approval (via halaman Approvals)

- Approver bertindak di **Approvals** generik → keputusan dipropagasi via callback
  `HandleApprovalStatusChange` (hanya diproses bila masih `PENDING_APPROVAL`):
  - `APPROVED` → `APPROVED_FINAL` + **auto-deduct saldo** (`applyLeaveUsage`) + ledger USAGE
  - `REJECTED` → `REJECTED_FINAL`
- Notifikasi ke employee: `LEAVE_APPROVED` / `LEAVE_REJECTED`.

### D. Batal (CANCELLED)

- Employee membatalkan request miliknya (`DRAFT`/`SUBMITTED`/`PENDING_APPROVAL`) via tombol Cancel
  → `PUT /requests/:id/status` `{"status":"CANCELLED"}`.
- Bila pembatalan terjadi setelah `APPROVED_FINAL` (mis. HR) → **saldo dikembalikan**
  (`reverseLeaveUsage`) + ledger **REVERSAL**.
- Notifikasi `LEAVE_CANCELLED`.

---

## 5. Ringkasan Status & Transisi

| Status | Makna | Transisi masuk |
|---|---|---|
| `DRAFT` | Baru dibuat | create |
| `SUBMITTED` | Diajukan | submit |
| `PENDING_APPROVAL` | Menunggu approval | submit (ke approval) |
| `APPROVED_FINAL` | Disetujui + saldo terpotong | callback approval APPROVED |
| `REJECTED_FINAL` | Ditolak | callback approval REJECTED |
| `CANCELLED` | Dibatalkan | cancel (reversal saldo bila dari APPROVED_FINAL) |

> Endpoint transisi tunggal: `PUT /requests/:id/status` (generik — belum ada endpoint `submit`/`cancel` khusus).

---

## 6. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Instance approval modul `leave`; approver bertindak di halaman Approvals |
| **Notification** | `LEAVE_APPROVED` / `LEAVE_REJECTED` / `LEAVE_CANCELLED` ke employee |
| **Setting (Holiday)** | `HolidayProvider` — company holiday di-exclude dari perhitungan working-day |
| **Attendance** | `AttendanceSessionUpdater` — integrasi sesi hadir (backend) |
| **Employee / Organization** | 🚫 belum ada cross-module read (blocker: validasi employee aktif, team calendar, manager/HR dashboard) |

---

## 7. Peta Halaman UI

| Menu | Halaman | Isi |
|---|---|---|
| Leave (dashboard) | `Leave.vue` | Balance cards per type (`quota/used/remaining`) + list request sendiri + tombol Cancel + dialog New Request + kalender bulan berjalan (`GET /calendar`) + kartu menu admin |
| Leave → Types | `LeaveTypes.vue` | CRUD jenis cuti |
| Leave → Accrual Policies | `LeaveAccrualPolicies.vue` | CRUD aturan accrual |
| Leave → Reasons | `LeaveReasons.vue` | CRUD alasan cuti |
| Approvals (generik) | `Approvals.vue` | Approve/Reject request cuti |

> `leave/admin` (halaman kartu terpisah) sudah **di-redirect ke `/leave`** — kartu menu admin kini tampil di dashboard.

---

## 8. Endpoint API Utama

Semua di bawah `/api/v1/tenant/leave/`.

| Area | Endpoint |
|---|---|
| Master | `POST/GET /types`, `GET/PUT/DELETE /types/:id` · `POST/GET /accrual-policies`, `GET/PUT/DELETE /accrual-policies/:id` · `POST/GET /reasons`, `GET/PUT/DELETE /reasons/:id` |
| Request | `POST/GET /requests`, `GET/PUT/DELETE /requests/:id`, `PUT /requests/:id/status`, `GET /requests/:id/details` |
| Balance | `GET /balances`, `GET /balances/employees/:employeeId/types/:leaveTypeId` |
| Calendar | `GET /calendar?employee_id=&from=&to=` (employee calendar) |
| Reports | `GET /reports/usage?from=&to=` (tenant-wide) · `GET /reports/on-leave-today` |

---

## 9. Catatan Penting

- **Draft vs Ajukan**: request DRAFT belum masuk approval/notifikasi sampai Submit.
- **Saldo dipotong otomatis** saat `APPROVED_FINAL` dan **dikembalikan** saat keluar dari status itu (ledger USAGE/REVERSAL).
- **Accrual engine belum ada** — quota di-set dari mana pun harus lewat jalur lain (belum ada endpoint adjustment HR).
- **Team/Organization Calendar & Manager/HR Dashboard** belum ada (backend-blocked cross-module read).
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
