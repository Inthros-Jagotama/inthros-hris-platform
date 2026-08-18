# Alur Pengisian Attendance & Business Travel (Runbook)

Dokumen ini menjelaskan **cara pakai / pengisian** modul **Attendance** — setup master, sesi
check-in/check-out (geofence), koreksi, dan lembur (overtime 2 alur) — **plus Business Travel /
Perjalanan Dinas (Bagian B)**, bagian dari module Attendance. Pola runbook seperti
[`module-payroll-user-flow.md`](module-payroll-user-flow.md).

- Plan pengembangan: [`module-attendance-plan.md`](../archive/module-attendance-plan.md) — 🟡 sebagian (inti ✅; Absent/Exempt detection, Manager/HR dashboard, payroll integration belum)
- Lokasi kode: `backend/internal/modules/attendance/` · `frontend/tenant/src/views/modules/attendance/`
- Daftar endpoint + contoh curl: [`../api/api-usage-guide.md`](../api/api-usage-guide.md) → §8.2 (tabel Attendance)

---

## 1. Ringkasan Alur End-to-End

```
SETUP (sekali)                    HARIAN                              KOREKSI & LEMBUR
┌──────────────────┐   ┌──────────────────────────┐   ┌─────────────────────────────────┐
│ Locations (GPS)  │   │ Sesi shift               │   │ Koreksi: SUBMITTED → PENDING_    │
│ Shifts           │──▶│  OPEN → CLOSED           │──▶│   APPROVAL → APPROVED/REJECTED   │
│ Events/Exempt    │   │  (MISSING_*, ABSENT,     │   │ Lembur alur SELF: SUBMITTED →     │
│ Settings         │   │   DAY_OFF, EXEMPT, LEAVE)│   │   PENDING_APPROVAL → APPROVED     │
└──────────────────┘   └──────────────────────────┘   │ Lembur alur ASSIGNED: WAITING_    │
                                                      │   ACTUAL → isi aktual → approval  │
                                                      └─────────────────────────────────┘
```

- **Session status:** `OPEN → CLOSED` · status lain otomatis: `MISSING_CHECKIN`, `MISSING_CHECKOUT`,
  `ABSENT`, `DAY_OFF`, `EXEMPT`, `LEAVE` (integrasi leave).
- **Koreksi:** `SUBMITTED → PENDING_APPROVAL → APPROVED/REJECTED`.
- **Lembur 2 alur:** alur **SELF** (request → approval) & alur **ASSIGNED** (penugasan →
  `WAITING_ACTUAL` → isian aktual → approval kedua).

---

## 2. Entitas Utama

| Entitas | Deskripsi |
|---|---|
| Attendance Location | Titik geofence (nama, koordinat, radius) untuk validasi check-in/out |
| Shift / Employee Shift | Jadwal kerja + penugasan shift per employee |
| Attendance Event | Hari libur/event (digunakan kalkulasi sesi) |
| Exempt Position | Posisi yang di-exempt dari absensi |
| Attendance Session | Sesi harian per employee (check-in/out + status) |
| Correction Request | Permintaan koreksi sesi (missing/wrong check-in/out) |
| Overtime Request | Lembur: alur SELF (request) / ASSIGNED (penugasan + isian aktual) |

---

## 3. TAHAP 1 — SETUP (dikerjakan sekali)

| Menu | Halaman | Isi |
|---|---|---|
| Locations | `AttendanceLocations.vue` | CRUD titik geofence (nama, koordinat, radius) |
| Shifts | `AttendanceShifts.vue` + `AttendanceEmployeeShifts.vue` | Master shift + penugasan per employee |
| Events | `AttendanceEvents.vue` | Hari libur/event |
| Exempt Positions | `AttendanceExemptPositions.vue` | Posisi exempt absensi |
| Settings | `AttendanceSettings.vue` | Konfigurasi modul |

Endpoint: `GET/POST /locations`, `GET/PUT/DELETE /locations/:id` · `/shifts`, `/shifts/:id` ·
`/employee-shifts`, `/employee-shifts/:id` · `/events`, `/events/:id` ·
`/exempt-positions`, `/exempt-positions/:id` · `/settings`.

---

## 4. TAHAP 2 — HARIAN (check-in / check-out)

- Employee check-in/check-out pada sesi shift miliknya (validasi **geofence** via `AttendanceLocations`).
- Sesi berjalan `OPEN` → `CLOSED`; status lain dihasilkan otomatis:
  `MISSING_CHECKIN`, `MISSING_CHECKOUT`, `ABSENT`, `DAY_OFF`, `EXEMPT`, `LEAVE`.
- Endpoint: `GET /sessions`, `GET /sessions/detail` · ringkasan `GET /summary`, `GET /stats/summary`.

> ⚠️ Deteksi `ABSENT`/`EXEMPT` otomatis butuh **scheduled job** (belum ada infra) — status ini
> saat ini bergantung jalur lain (integrasi leave/event).

---

## 5. TAHAP 3 — KOREKSI

- Employee/HR mengajukan koreksi sesi (mis. lupa check-in): `POST /corrections`.
- Alur: `SUBMITTED → PENDING_APPROVAL → APPROVED` / `REJECTED` (via Central Approval).
- Sesi asli **tidak pernah di-mutasi langsung** — koreksi tercatat terpisah.
- Endpoint: `GET/POST /corrections`, `GET/PUT/DELETE /corrections/:id` (+ status).

---

## 6. TAHAP 4 — LEMBUR (2 alur)

Lembur punya **dua alur** (§32b) — keduanya berujung pada **isian aktual + approval kedua**.
Satu request memiliki **dua instance approval** (request & aktual) dengan `document_id` yang sama.

### 6.1 Alur SELF (request oleh karyawan)

1. **Buat request** — `POST /overtime-requests`: `employee_id`, `work_date`, `start_time_local`,
   `end_time_local`, `requested_minutes`, `reason`. Status `SUBMITTED` → `PENDING_APPROVAL`
   (instance approval #1 dibuat, modul `attendance`, flow auto-resolve).
2. **Approval #1** (halaman Approvals) — `APPROVED` → **`WAITING_ACTUAL`** (rencana disetujui;
   notif `OVERTIME_APPROVED`) · `REJECTED` → `REJECTED` (notif `OVERTIME_REJECTED`).
3. **Isi aktual** — `PUT /overtime-requests/:id/actual`: `actual_start_time_local`,
   `actual_end_time_local` (harus setelah start), `actual_note`, `attachment_url` →
   **`ACTUAL_SUBMITTED`** (instance approval #2 dibuat; hanya pemilik request yang boleh mengisi).
4. **Approval #2** — `APPROVED` → `APPROVED` (**final**, notif `OVERTIME_ACTUAL_APPROVED`) ·
   `REJECTED` → `REJECTED` (notif `OVERTIME_ACTUAL_REJECTED`).

### 6.2 Alur ASSIGNED (penugasan oleh atasan)

1. **Pilih bawahan** — `GET /overtime-requests/assignable-employees` (bawahan efektif dari org tree).
2. **Tugaskan** — `POST /overtime-requests/assign`: `assigned_employee_id`, `work_date`,
   start/end time, `requested_minutes`, `reason` → langsung **`WAITING_ACTUAL`** tanpa approval
   penugasan; notif `OVERTIME_ASSIGNED` ke bawahan.
3. **Isi aktual** — karyawan mengisi aktual (langkah 3 alur SELF) → `ACTUAL_SUBMITTED`.
4. **Approval #2** — seperti alur SELF → `APPROVED` (final) / `REJECTED`.

### 6.3 Batal

- `POST /overtime-requests/:id/cancel` → `CANCELLED` (sebelum isian aktual — kedua alur).

**Ringkasan status lembur:** `SUBMITTED → PENDING_APPROVAL → WAITING_ACTUAL → ACTUAL_SUBMITTED → APPROVED` · `REJECTED` (di approval #1/#2) · `CANCELLED`

---

## 7. Ringkasan Status

| Entitas | Status |
|---|---|
| Session | `OPEN → CLOSED` · `MISSING_CHECKIN` / `MISSING_CHECKOUT` / `ABSENT` / `DAY_OFF` / `EXEMPT` / `LEAVE` |
| Correction | `SUBMITTED → PENDING_APPROVAL → APPROVED` / `REJECTED` |
| Overtime | `SUBMITTED → PENDING_APPROVAL → WAITING_ACTUAL → ACTUAL_SUBMITTED → APPROVED` · `REJECTED` (approval #1/#2) · `CANCELLED` (kedua alur SELF & ASSIGNED) |

---

## 8. Integrasi Lintas Modul

| Modul | Peran |
|---|---|
| **Central Approval** | Approval koreksi & lembur (kedua alur) |
| **Leave** | `AttendanceSessionUpdater` — sesi jadi `LEAVE` saat cuti |
| **Notification** | Notifikasi hasil approval (pola modul lain) |
| Payroll | 🚫 belum terintegrasi (Phase 13 plan — sisa) |
| Employee / Organization | 🚫 belum ada cross-module read (blocker Manager/HR dashboard & team calendar) |

---

## 9. Peta Halaman UI

| Menu | Halaman |
|---|---|
| Attendance (hub) | `Attendance.vue` |
| Sesi | `AttendanceSessions.vue` |
| Koreksi | `AttendanceCorrections.vue` |
| Lembur | `AttendanceOvertime.vue` |
| Shift | `AttendanceShifts.vue` / `AttendanceEmployeeShifts.vue` |
| Event | `AttendanceEvents.vue` |
| Exempt Positions | `AttendanceExemptPositions.vue` |
| Locations | `AttendanceLocations.vue` |
| Reports | `AttendanceReports.vue` |
| Settings | `AttendanceSettings.vue` |

---

## 10. Endpoint API Utama

Semua di bawah `/api/v1/tenant/attendance/`.

| Area | Endpoint |
|---|---|
| Setup | `GET/POST /locations`, `GET/PUT/DELETE /locations/:id` · `/shifts`, `/shifts/:id` · `/employee-shifts`, `/employee-shifts/:id` · `/events`, `/events/:id` · `/exempt-positions`, `/exempt-positions/:id` · `/settings` |
| Sesi | `GET /sessions`, `GET /sessions/detail`, `GET /calendar`, `GET /summary`, `GET /stats/summary`, `GET /stats/overtime-trend`, `GET /reports/sessions` |
| Koreksi | `GET/POST /corrections`, `GET/PUT/DELETE /corrections/:id` |
| Lembur | `GET/POST /overtime-requests`, `GET/PUT/DELETE /overtime-requests/:id`, `PUT /overtime-requests/:id/actual`, `POST /overtime-requests/:id/cancel`, `POST /overtime-requests/assign`, `GET /overtime-requests/assignable-employees` |

---

## 11. Catatan Penting

- **Geofence** divalidasi saat check-in/out (AttendanceLocations).
- **Koreksi tidak mengubah sesi asli** — tercatat terpisah + approval.
- **Lembur ASSIGNED** menunggu isian aktual sebelum approval kedua.
- **Belum ada**: scheduled job deteksi ABSENT/EXEMPT, Manager/HR Dashboard & Team Calendar,
  payroll integration (Phase 13), deteksi otomatis WRONG_CHECKIN/CHECKOUT.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.

---

# ══ BAGIAN B — BUSINESS TRAVEL / PERJALANAN DINAS ══

> Business Travel adalah bagian dari module **Attendance** — runbook pengisiannya digabung ke
> dokumen ini. Plan: [`module-attendance-business-travel-development-plan.md`](../archive/module-attendance-business-travel-development-plan.md)
> · UI: `frontend/tenant/src/views/modules/attendance/business-travel/`
> (`BusinessTravelList.vue` / `BusinessTravelDetail.vue`).

## 12. Ringkasan Alur (Business Travel)

```
SETUP (sekali)                    PENGAJUAN                    PELAKSANAAN                 SETTLEMENT
┌────────────────────┐   ┌────────────────────────┐   ┌──────────────────────────┐   ┌──────────────────────┐
│ Expense Categories │   │ DRAFT → SUBMITTED       │   │ IN_PROGRESS              │   │ Settlement: PENDING   │
│ Funding Methods    │──▶│  → APPROVED (approval)  │──▶│  ├─ expenses + bukti      │──▶│  → SUBMITTED         │
└────────────────────┘   │  → IN_PROGRESS          │   │  ├─ funding (PENDING →    │   │  → APPROVED/REJECTED │
                         └────────────────────────┘   │  │    PROCESSING → FUNDED) │   │  → refund/reimburse  │
                                                      │  └─ schedules/agenda       │   │  → CLOSED            │
                                                      └──────────────────────────┘   └──────────────────────┘
```

- **Travel status:** `DRAFT → SUBMITTED → APPROVED → IN_PROGRESS → COMPLETED → CLOSED` · terminal: `REJECTED`, `CANCELLED`
- **Funding:** `PENDING → PROCESSING → FUNDED` · terminal: `CANCELLED`, `REVERSED`
- **Expense:** `DRAFT → SUBMITTED → APPROVED/REJECTED` · **Settlement:** `PENDING → SUBMITTED → APPROVED/REJECTED`
- **Reimbursement:** `REQUESTED → APPROVED → PAID` (process → approve → pay)

---

## 13. Entitas Utama (Business Travel)

| Entitas | Deskripsi |
|---|---|
| Business Travel | Perjalanan dinas (judul, tanggal mulai-selesai, alasan, status) |
| Participant | Employee peserta perjalanan |
| Destination | Kota/tempat tujuan |
| Activity / Agenda | Kegiatan per hari |
| Schedule / Transport | Jadwal & transportasi |
| Expense + Bukti | Estimasi biaya + dokumen bukti per expense |
| Funding | Pendanaan (company paid / deposit / advance) |
| Settlement | Penyelesaian akhir (scenario refund/reimbursement) |
| Refund / Reimbursement | Pengembalian dana |

---

## 14. TAHAP 1 — SETUP Master (Business Travel)

| Menu | Endpoint |
|---|---|
| Expense Categories | `GET/POST /business-travel-expense-categories` |
| Funding Methods | `GET/POST /business-travel-funding-methods` |

---

## 15. TAHAP 2 — PENGAJUAN (Business Travel)

1. **Buat travel (DRAFT)** — `POST /business-travels`: judul, tanggal, alasan, `requested_by`.
2. **Lengkapi detail:**
   - Peserta: `POST /business-travels/:id/participants` (hapus via `.../participants/:participantId`)
   - Tujuan: `POST .../destinations` · Agenda: `POST .../activities` · Jadwal/transport: `POST .../schedules`
   - Estimasi biaya: `POST .../expenses` (+ kategori dari master)
3. **Submit** — `POST /business-travels/:id/submit` → status `SUBMITTED`, instance approval dibuat
   (Central Approval). Keputusan APPROVED → `APPROVED`; REJECTED → `REJECTED`.
4. **Batal** — `POST /business-travels/:id/cancel` → `CANCELLED`.

---

## 16. TAHAP 3 — PELAKSANAAN (IN_PROGRESS)

- Travel berjalan → status `IN_PROGRESS`.
- **Biaya aktual:** `POST /business-travels/:id/expenses` + upload bukti per expense
  (`POST .../expenses/:expenseId/documents`); expense `DRAFT → SUBMITTED → APPROVED/REJECTED`.
- **Funding** (per participant): `POST .../fundings` → `PENDING`; proses → `PROCESSING`;
  konfirmasi transfer → `POST .../fundings/:fundingId/confirm` → `FUNDED`. Terminal: `CANCELLED`, `REVERSED`.
- **Documents** per travel: `POST .../documents`.

---

## 17. TAHAP 4 — SETTLEMENT (penyelesaian)

1. **Buat settlement** — `POST /business-travels/:id/settlements` → `PENDING`.
2. **Submit** — `POST .../settlements/:settlementId/submit` → `SUBMITTED` + approval.
3. **Keputusan** — `APPROVED`/`REJECTED` (callback Central Approval; modul slug `business_travel`).
4. **Hasil akhir** tergantung scenario (reimbursement murni / deposit vs actual / mixed):
   - **Refund** (dana kembali ke company): `POST /business-travels/:id/refunds` → `POST .../refunds/:refundId/confirm`
   - **Reimbursement** (company bayar ke employee): `POST .../reimbursements` →
     `process` → `approve` → `pay` (`REQUESTED → APPROVED → PAID`)
5. Travel ditutup → `COMPLETED` → `CLOSED`.

> ⚠️ Saat settlement `APPROVED`, item actual expense yang dikonfigurasi payroll-eligible di-push
> ke payroll adjustment (`pushBusinessTravelPayrollAdjustments`, `SourceType=BUSINESS_TRAVEL`,
> status langsung `APPROVED`) — **test khusus alur ini belum ada** (lihat plan §54).

---

## 18. Ringkasan Status (Business Travel)

| Entitas | Status |
|---|---|
| Travel | `DRAFT → SUBMITTED → APPROVED → IN_PROGRESS → COMPLETED → CLOSED` · `REJECTED`, `CANCELLED` |
| Funding | `PENDING → PROCESSING → FUNDED` · `CANCELLED`, `REVERSED` |
| Expense | `DRAFT → SUBMITTED → APPROVED` / `REJECTED` |
| Settlement | `PENDING → SUBMITTED → APPROVED` / `REJECTED` |
| Refund | `PENDING → (confirm)` |
| Reimbursement | `REQUESTED → APPROVED → PAID` · `REJECTED` |

---

## 19. Integrasi Lintas Modul (Business Travel)

| Modul | Peran |
|---|---|
| **Central Approval** | Approval travel, settlement (dan alur terkait) |
| **Attendance** | Satu package — sesi/periode terkait perjalanan |
| **Payroll** | Push adjustment saat settlement APPROVED (actual expense payroll-eligible) — test khusus belum ada |
| **Document / Uploads** | Bukti per expense/travel |

---

## 20. Peta Halaman UI (Business Travel)

| Menu | Halaman |
|---|---|
| Attendance → Business Travel (list) | `BusinessTravelList.vue` |
| Attendance → Business Travel (detail) | `BusinessTravelDetail.vue` (travel + peserta/tujuan/agenda/schedule/expenses/funding/settlement) |

---

## 21. Endpoint API Utama (Business Travel)

Semua di bawah `/api/v1/tenant/attendance/`.

| Area | Endpoint |
|---|---|
| Master | `GET/POST /business-travel-expense-categories` · `GET/POST /business-travel-funding-methods` |
| Travel | `GET/POST /business-travels`, `GET/PUT/DELETE /business-travels/:id`, `POST /business-travels/:id/submit`, `POST /business-travels/:id/cancel` |
| Detail | `.../:id/participants(+/:participantId)`, `.../destinations(+/:id)`, `.../activities(+/:activityId)`, `.../schedules(+/:scheduleId)`, `.../documents(+/:documentId)` |
| Expense | `.../expenses(+/:expenseId)` + `.../expenses/:expenseId/documents` |
| Funding | `.../fundings(+/:fundingId)`, `POST .../fundings/:fundingId/confirm`, `.../fundings/:fundingId/documents` |
| Settlement | `POST/GET .../settlements(+/:id)`, `POST .../settlements/:settlementId/submit`, `GET/PUT /business-travel-settlements/:settlementId` |
| Refund / Reimburse | `POST .../refunds`, `POST .../refunds/:refundId/confirm` · `POST .../reimbursements(+/:reimbursementId/process)` + `/approve` + `/pay` |

---

## 22. Catatan Penting (Business Travel)

- **Settlement menentukan hasil akhir** — refund (dana kembali) vs reimbursement (dibayar ke employee),
  termasuk scenario mixed funding (deposit + reimbursement).
- **Funding wajib di-confirm** (`confirm`) sebelum dianggap `FUNDED`.
- **Payroll push** otomatis saat settlement APPROVED hanya untuk actual expense yang eksplisit
  payroll-eligible — Refund/Reimbursement/Advance tidak dianggap salary.
- **Server restart** diperlukan setelah perubahan backend agar migrasi & fitur baru aktif.
