Saya sudah pelajari struktur `005_leave.sql`. Struktur dasarnya sudah mencakup **Leave Type, Accrual Policy, Reason, Request, Request Detail, dan Balance**. Dari struktur tersebut, saya akan mempertahankan pola yang ada dan menambahkan kebutuhan development di atasnya, terutama karena Anda sudah memiliki **Approval Module terpusat**.

Berikut plan Markdown yang bisa langsung Anda simpan sebagai `performance-management-leave-plan.md` atau `leave-management-plan.md`.

> ⚠️ **Status vs. Plan ini**: dokumen ini awalnya ditulis seolah modul Leave belum ada (greenfield). Setelah dicek ulang terhadap kode aktual, **backend modul Leave sudah cukup jauh diimplementasikan** — termasuk integrasi ke Central Approval Module — sementara beberapa bagian rencana (cancellation flow sebagai sub-flow tersendiri, accrual engine dari `LeaveAccrualPolicy`, carry forward, expiry, adjustment HR) **memang belum ada** — balance ledger/usage/reversal (Phase 6), working-day calculation (Phase 3), employee calendar (Phase 7), dan notification (Phase 8) sudah terimplementasi sejak 2026-08-08/09 — dan **frontend awalnya masih placeholder "Coming soon"** (kini FE-1 & FE-2 selesai, lihat baris 7 di bawah). Lihat section **"Implementation Status"** di bagian bawah dokumen untuk status per-fase yang sudah diverifikasi terhadap kode, dan catatan blockquote (`>`) yang disisipkan di beberapa section untuk koreksi spesifik.
>
> ✅ **Update 2026-08-09**: Frontend Implementation Plan FE-1 (My Leave Dashboard) dan FE-2 (Admin Configuration) **sudah diimplementasikan** — lihat status di section "Frontend Implementation Plan" di bagian bawah dokumen.

# Leave Management Module Development Plan

## 1. Objective

Membangun modul **Leave Management** untuk mengelola seluruh proses cuti/izin karyawan mulai dari konfigurasi jenis cuti, perhitungan kuota, pengajuan, approval, saldo cuti, sampai histori penggunaan.

Modul harus terintegrasi dengan:

* Employee Management
* Organization Management
* Attendance
* Work Schedule / Shift
* Approval Module
* Notification
* Payroll
* Performance Management

Prinsip utama:

* Seluruh Primary Key menggunakan UUID.
* Employee menjadi pemilik Leave Request.
* Approval menggunakan **Central Approval Module**.
* Modul Leave tidak membuat engine approval sendiri.
* Saldo cuti harus dapat dilacak dan diaudit.
* Histori saldo tidak boleh hilang ketika terjadi adjustment.
* Leave Request menggunakan snapshot pada saat transaksi jika diperlukan.
* Perhitungan hari cuti harus memperhatikan kalender kerja dan hari libur.

---

# 2. Existing Database Structure

Berdasarkan `005_leave.sql`, tabel yang sudah tersedia:

```text
leave_types
    │
    ├── leave_accrual_policies
    │
    └── leave_requests
            │
            └── leave_request_details

employee_leave_balances

leave_reasons
```

Struktur saat ini:

```text
leave_types
├── leave_accrual_policies
│
└── leave_requests
    └── leave_request_details

leave_reasons

employee_leave_balances
```

---

# 3. Existing Tables

## 3.1 leave_types

Master jenis cuti/izin.

Field yang sudah tersedia:

* id
* code
* name
* description
* is_paid
* requires_attachment
* allow_half_day
* default_quota_days
* quota_period
* counts_against_quota
* allow_hourly
* is_active
* deleted_at
* created_at
* updated_at

Contoh:

```text
ANNUAL
Annual Leave

SICK
Sick Leave

MATERNITY
Maternity Leave

UNPAID
Unpaid Leave
```

### Enhancement

Pertimbangkan penambahan:

```text
requires_approval
minimum_notice_days
maximum_consecutive_days
allow_backdate
allow_future_date
gender
employment_type
```

Jika aturan tersebut diperlukan oleh perusahaan.

---

# 4. Leave Accrual Policy

## Existing Table

```text
leave_accrual_policies
```

Digunakan untuk menentukan hak cuti berdasarkan masa kerja dan periode berlaku.

Contoh:

```text
Base Quota       : 12 hari
Extra Every      : 2 tahun
Extra Days       : 1 hari
Max Extra        : 6 hari
```

### Enhancement

Tambahkan jika dibutuhkan:

```text
accrual_type
accrual_frequency
prorate_on_join
prorate_on_exit
carry_forward
max_carry_forward_days
expiry_month
```

---

# 5. Leave Reasons

## Existing Table

```text
leave_reasons
```

Digunakan untuk alasan pengajuan cuti.

Contoh:

```text
Personal
Family
Medical
Travel
Other
```

Tidak perlu membuat tabel baru jika kebutuhan masih sesuai struktur existing.

---

# 6. Leave Request

## Existing Table

```text
leave_requests
```

Merupakan header transaksi pengajuan cuti.

Struktur saat ini sudah mendukung:

* Employee
* Leave Type
* Start Date
* End Date
* Duration Mode
* Requested Days
* Reason
* Attachment
* Status
* Approval Instance
* Start Time
* End Time
* Submission
* Approval
* Rejection
* Cancellation

> ✅ **Sudah diimplementasikan.** `leave_requests` (migration `005_leave.sql`) sudah punya seluruh field di atas, termasuk `approval_instance_id` (sudah terhubung ke Central Approval Module, lihat Section 7). Field lama `supervisor_id`/`supervisor_action_at`/`supervisor_note`/`hr_id`/`hr_action_at`/`hr_note` masih ada di skema & model Go, tapi **sudah tidak dipakai sebagai penggerak workflow** — hanya `*_note` yang masih ditulis sebagai catatan bebas. Rekomendasi: hentikan penulisan field-field ini di kode baru (tanpa migration breaking untuk drop kolom), bukan dibiarkan ambigu.
>
> `requested_days` saat ini **diterima langsung dari payload client**, belum divalidasi/dihitung ulang di server terhadap kalender kerja (lihat Section 13 — gap nyata, bukan sekadar enhancement opsional).

---

# 7. Approval Integration

> ✅ **Sudah diimplementasikan sepenuhnya**, bukan lagi rencana. `backend/internal/modules/leave/service.go` sudah punya `ApprovalEngine` interface (`CreateApprovalInstance`, `GetApprovalInstanceStatus`), `SetApprovalEngine`, dan `HandleApprovalStatusChange` (push-callback yang hanya bertindak jika status masih `PENDING_APPROVAL`) — pola yang identik dengan payroll/performance. Di `cmd/server/main.go`, `leaveSvc.SetApprovalEngine(sharedApprovalEngine)` dan `approvalSvc.RegisterStatusHandler("leave", ...)` sudah didaftarkan, memakai adapter `sharedApprovalEngine` yang sama dengan payroll/reimbursement. Ada test coverage di `approval_integration_test.go`.
>
> Leave memakai **satu module slug polos `"leave"`** — dikonfirmasi tidak ada entry `"leave"` di `subscriptionModuleAliases`/`subscriptionModuleSubslots` (`approval/service.go`), berbeda dari pola dua-checkpoint KPI/OKR (`performance_kpi_target`/`performance_kpi_realization`, `okr_key_result`/`okr_assessment`). Artinya Leave hanya butuh **satu** approval checkpoint per request — bukan keputusan desain yang masih terbuka, tapi sudah demikian adanya di kode.

## Important

Field berikut sudah tersedia:

```text
approval_instance_id
```

Artinya Leave sudah siap menggunakan Central Approval Module.

Leave **tidak boleh membuat tabel approval supervisor/HR sendiri** sebagai mekanisme utama.

Field lama:

```text
supervisor_id
supervisor_action_at
supervisor_note

hr_id
hr_action_at
hr_note
```

perlu dievaluasi kembali.

### Recommendation

Jika Central Approval sudah menjadi standar seluruh HRIS, field tersebut sebaiknya **tidak digunakan sebagai workflow approval utama**.

Workflow menjadi:

```text
Leave Request
      │
      ▼
Approval Instance
      │
      ▼
Approval Flow
      │
      ▼
Approval Step
      │
      ▼
Approval Task
      │
      ▼
Approval Action
      │
      ▼
Leave Request Status
```

Dengan demikian:

```text
Leave
Overtime
KPI
OKR
Employee Movement
Recruitment
Reimbursement
```

semuanya menggunakan engine approval yang sama.

---

# 8. Leave Request Details

## Existing Table

```text
leave_request_details
```

Digunakan untuk menyimpan detail setiap tanggal cuti.

Contoh:

```text
Request

01 Jan - 03 Jan

↓

leave_request_details

01 Jan
02 Jan
03 Jan
```

Keuntungan:

* Perhitungan saldo lebih mudah.
* Validasi overlap lebih mudah.
* Integrasi Attendance lebih mudah.
* Histori per tanggal tersedia.

---

# 9. Leave Balance

## Existing Table

```text
employee_leave_balances
```

Struktur saat ini:

```text
employee_id
leave_type_id
period_year

quota_days
used_days
remaining_days
```

Unique:

```text
employee_id
leave_type_id
period_year
```

Ini sudah tepat untuk balance tahunan.

> ✅ Tabel `employee_leave_balances` sudah ada persis seperti ini (unique `employee_id, leave_type_id, period_year`), dengan `GetLeaveBalance`/`ListLeaveBalances` (read-only).
>
> ✅ **Gap P0 sudah teratasi (Phase 6, 2026-08-08)** — `leave/balance.go` kini menulis/mengurangi `used_days`/`remaining_days` saat leave di-approve (`applyLeaveUsage`, transisi ke `APPROVED_FINAL`) dan mengembalikannya saat dibatalkan/keluar dari `APPROVED_FINAL` (`reverseLeaveUsage`), keduanya tercatat di `leave_balance_transactions` (ledger) lewat `writeLeaveBalanceTransaction`. Detail di catatan Phase 6 di bawah. Yang masih berupa proposal: accrual dari `LeaveAccrualPolicy`, adjustment HR, carry forward, expiry.

---

# 10. Leave Balance Enhancement

> ✅ **Sudah diimplementasikan (Phase 6, 2026-08-08)** — tabel `leave_balance_transactions` + model + repository ada sejak Phase 1, dan sejak Phase 6 ledger benar-benar ditulis: setiap USAGE/REVERSAL lewat `writeLeaveBalanceTransaction` dengan `balance_before`/`balance_after` (lihat catatan Phase 6 di bawah). Entri ACCRUAL/ADJUSTMENT/CARRY_FORWARD/EXPIRY masih kosong karena fitur-fitur tersebut (accrual engine, adjustment, carry forward, expiry) tetap proposal.

Disarankan menambahkan mekanisme **Balance Ledger**.

## New Table

```text
leave_balance_transactions
```

### Purpose

Menyimpan seluruh perubahan saldo.

Field:

| Field            | Type         | Description                                                      |
| ---------------- | ------------ | ---------------------------------------------------------------- |
| id               | uuid         | Primary Key                                                      |
| employee_id      | uuid         | Employee                                                         |
| leave_type_id    | uuid         | Leave Type                                                       |
| balance_id       | uuid         | Balance                                                          |
| transaction_type | varchar      | ACCRUAL / USAGE / ADJUSTMENT / REVERSAL / CARRY_FORWARD / EXPIRY |
| reference_type   | varchar      | Source Module                                                    |
| reference_id     | uuid         | Source ID                                                        |
| amount           | decimal(6,2) | Transaction Amount                                               |
| balance_before   | decimal(6,2) | Previous Balance                                                 |
| balance_after    | decimal(6,2) | New Balance                                                      |
| note             | varchar      | Description                                                      |
| created_by       | uuid         | User                                                             |
| created_at       | timestamp    |                                                                  |

Contoh:

```text
Opening Balance
+12

Leave Request
-2

Adjustment
+1

Final Balance
11
```

---

# 11. Balance Calculation

Jangan menjadikan:

```text
remaining_days
```

sebagai satu-satunya sumber kebenaran.

Gunakan:

```text
Balance

=

Opening / Accrual

+

Adjustment

+

Carry Forward

-

Usage

-

Expiry
```

`employee_leave_balances` menjadi current balance/cache.

Sedangkan:

```text
leave_balance_transactions
```

menjadi histori/ledger.

---

# 12. Leave Request Validation

Sebelum request dapat disubmit:

### Employee Validation

* Employee aktif.
* Employee memiliki Organization aktif.
* Employee memiliki leave balance jika jenis cuti membutuhkan quota.

### Leave Type Validation

* Leave Type aktif.
* Employee memenuhi eligibility.
* Attachment tersedia jika diwajibkan.

### Date Validation

* Start Date <= End Date.
* Tidak boleh overlap dengan leave yang aktif.
* Tidak boleh menggunakan tanggal yang sudah approved.
* Sesuai aturan backdate.
* Sesuai minimum notice jika diterapkan.

### Balance Validation

```text
Requested Days <= Available Balance
```

Untuk leave yang:

```text
counts_against_quota = 0
```

saldo tidak dikurangi.

---

# 13. Working Day Calculation

Requested Days tidak boleh selalu dihitung:

```text
End Date - Start Date + 1
```

Perhitungan harus mempertimbangkan:

* Working Calendar
* Weekend
* Public Holiday
* Employee Schedule
* Shift
* Half Day
* Hourly Leave

Contoh:

```text
Monday
Tuesday
Wednesday Holiday
Thursday
Friday
```

Request:

```text
Monday - Friday
```

Maka:

```text
Requested Days = 4
```

bukan 5.

---

# 14. Half Day

Untuk:

```text
HALF_DAY_AM
HALF_DAY_PM
```

nilai:

```text
day_fraction = 0.50
```

Contoh:

```text
Monday      0.50
Tuesday     1.00
Wednesday   1.00
```

Total:

```text
2.50 days
```

---

# 15. Hourly Leave

Untuk:

```text
duration_mode = HOURLY
```

gunakan:

```text
start_time
end_time
```

Perhitungan:

```text
Duration Hours
=
end_time - start_time
```

Konversi ke quota harus mengikuti policy perusahaan.

Contoh:

```text
8 hours = 1 day
```

atau dapat menggunakan konfigurasi:

```text
hours_per_day
```

---

# 16. Leave Request Status

> ✅ Status machine di bawah ini **sudah diimplementasikan persis seperti tertulis** — konstanta `LeaveStatusDraft/Submitted/PendingApproval/ApprovedFinal/RejectedFinal/Cancelled` di `leave/service.go` cocok 1:1. Tidak perlu perubahan, hanya konfirmasi.

Status existing:

```text
DRAFT
SUBMITTED
PENDING_APPROVAL
APPROVED_FINAL
REJECTED_FINAL
CANCELLED
```

Flow:

```text
DRAFT
  │
  ▼
SUBMITTED
  │
  ▼
PENDING_APPROVAL
  │
  ├───────────────┐
  ▼               ▼
APPROVED_FINAL   REJECTED_FINAL
  │
  ▼
CANCELLED (jika diperbolehkan)
```

---

# 17. Approval Flow Configuration

Leave menggunakan Central Approval Module.

Contoh:

```text
Leave Request
      │
      ▼
Approval Flow
      │
      ├── Supervisor
      │
      └── HR
```

Namun konfigurasi tersebut **tidak disimpan di Leave Module**.

Disimpan pada:

```text
Approval Module
```

Leave hanya:

```text
create approval instance
```

dan menerima hasil approval.

---

# 18. Cancellation

> 🔶 **Sub-flow cancellation tetap belum ada**, tapi catatan asli perlu dikoreksi: reversal saldo **sudah ada** sejak Phase 6 — `reverseLeaveUsage` (`leave/balance.go`) mengembalikan hari yang terpakai + menulis ledger REVERSAL saat request keluar dari `APPROVED_FINAL` via `UpdateLeaveRequestStatus` (mis. cancel oleh HR setelah approve; jalur `HandleApprovalStatusChange` CANCELLED hanya terjadi dari `PENDING_APPROVAL` yang saldonya belum pernah dikurangi, jadi tidak butuh reversal). Yang masih belum ada: sub-flow cancellation-request dengan approval tersendiri, endpoint `cancel` khusus, dan tabel `leave_cancellation_requests` (§19 tetap proposal).

Employee dapat membatalkan request berdasarkan policy.

Contoh:

```text
DRAFT
SUBMITTED
PENDING_APPROVAL
```

dapat dibatalkan.

Sedangkan:

```text
APPROVED_FINAL
```

mungkin membutuhkan approval cancellation.

Flow:

```text
Approved Leave
      │
      ▼
Cancellation Request
      │
      ▼
Approval Module
      │
      ▼
Approved
      │
      ▼
Reverse Balance
```

Jika leave sudah mengurangi saldo:

```text
Usage -2
```

maka cancellation menghasilkan:

```text
Reversal +2
```

---

# 19. New Table - Leave Cancellation

Jika cancellation membutuhkan approval khusus, buat:

```text
leave_cancellation_requests
```

Field:

| Field                | Description         |
| -------------------- | ------------------- |
| id                   | UUID                |
| leave_request_id     | Leave Request       |
| employee_id          | Employee            |
| reason               | Cancellation Reason |
| status               | Status              |
| approval_instance_id | Approval Instance   |
| submitted_at         | Submit Time         |
| approved_at          | Approval Time       |
| rejected_at          | Rejection Time      |
| created_at           |                     |
| updated_at           |                     |

Jika kebutuhan cancellation masih sederhana, cancellation dapat tetap berada pada `leave_requests`.

---

# 20. Leave Calendar

Tambahkan fitur Calendar.

Menampilkan:

```text
Employee Leave Calendar
Team Leave Calendar
Organization Leave Calendar
Company Leave Calendar
```

Contoh:

```text
January 2027

Mon Tue Wed Thu Fri

Asep      ███
Budi          ██
Citra             ███
```

Manager dapat melihat benturan jadwal cuti team.

---

# 21. Team Leave Validation

Sistem dapat mencegah terlalu banyak anggota organisasi mengambil cuti pada waktu yang sama.

Contoh policy:

```text
Organization:
IT Support

Minimum Available:
70%
```

Jika 4 dari 5 anggota mengajukan cuti:

```text
Available = 20%

↓

Warning / Reject
```

Policy ini dapat menjadi enhancement jika belum diperlukan pada tahap awal.

---

# 22. Leave Eligibility

Beberapa jenis leave memiliki persyaratan.

Contoh:

```text
Maternity Leave
```

hanya untuk employee tertentu.

Atau:

```text
Annual Leave
```

baru tersedia setelah masa kerja tertentu.

Jika eligibility semakin kompleks, tambahkan:

```text
leave_eligibility_rules
```

Contoh:

```text
leave_type_id
employment_type
minimum_service_days
gender
marital_status
```

Gunakan hanya jika business rule memang diperlukan.

---

# 23. Notification

Integrasi dengan Notification Module.

### Employee

* Leave submitted
* Leave approved
* Leave rejected
* Leave cancelled

### Approver

* New leave request
* Reminder approval
* Approval overdue

### HR

* Leave approved
* Leave rejected
* Balance adjustment

---

# 24. Attendance Integration

Leave yang sudah:

```text
APPROVED_FINAL
```

dapat menjadi data Attendance.

Contoh:

```text
Leave
01 Jan
02 Jan

↓

Attendance

01 Jan = ON_LEAVE
02 Jan = ON_LEAVE
```

Leave yang:

```text
DRAFT
PENDING_APPROVAL
REJECTED
```

tidak boleh mempengaruhi attendance final.

---

# 25. Payroll Integration

Leave type memiliki:

```text
is_paid
```

Sehingga payroll dapat mengetahui:

```text
Paid Leave
Unpaid Leave
```

Contoh:

```text
Unpaid Leave
3 days

↓

Payroll Deduction
```

Perhitungan nominal sebaiknya tetap berada di Payroll Module.

Leave hanya menyediakan:

```text
leave_days
is_paid
employee
period
```

---

# 26. Leave Balance Adjustment

HR dapat melakukan adjustment.

Contoh:

```text
Employee
Annual Leave

Current = 8

Adjustment = +2

New Balance = 10
```

Adjustment harus menghasilkan ledger:

```text
ADJUSTMENT
+2
```

dan tidak boleh mengubah balance tanpa histori.

---

# 27. Carry Forward

Jika policy mengizinkan:

```text
2026 Remaining
=
5 days
```

Carry Forward:

```text
2027 Opening
=
5 days
```

Ledger:

```text
CARRY_FORWARD
+5
```

---

# 28. Expiry

Jika carry forward memiliki masa berlaku:

```text
Carry Forward
5 days

Expiry
31 March 2027
```

Jika tidak digunakan:

```text
EXPIRY
-remaining
```

Saldo harus tercatat pada ledger.

---

# 29. Menu Structure

> ✅ **Frontend sudah dibangun (FE-1 & FE-2, 2026-08-09).** `frontend/tenant/src/views/modules/Leave.vue` sudah menjadi My Leave dashboard (balance cards, list request sendiri, dialog New Request, kalender bulan berjalan), plus `LeaveAdmin.vue` (index kartu), `LeaveTypes.vue`, `LeaveAccrualPolicies.vue`, `LeaveReasons.vue` (CRUD Dialog inline) — route `leave`, `leave/admin`, `leave/types`, `leave/accrual-policies`, `leave/reasons`. Struktur menu di bawah ini adalah target penuh; bagian yang belum ada halamannya (Dashboard Manager/HR, Team Calendar, Reports, Adjustment) tetap backend-blocked — lihat "Eksplisit di luar cakupan rencana FE ini" di Frontend Implementation Plan.

```text
Leave Management

├── Dashboard
│
├── My Leave
│   ├── Request Leave
│   ├── My Requests
│   └── My Balance
│
├── Team Leave
│   ├── Team Calendar
│   ├── Pending Requests
│   └── Team Balance
│
├── Leave Calendar
│
├── Leave Balance
│
├── Leave Adjustment
│
├── Reports
│
└── Settings
    ├── Leave Types
    ├── Leave Reasons
    └── Accrual Policies
```

---

# 30. Dashboard

## Employee

Menampilkan:

```text
Annual Leave
Remaining : 8 days

Used      : 4 days
Quota     : 12 days
```

dan:

* Upcoming Leave
* Pending Request
* Recent History

---

## Manager

Menampilkan:

* Pending Approval
* Team Leave Calendar
* Team Leave Balance
* Upcoming Leave
* Leave Conflict

---

## HR

Menampilkan:

* Total Leave
* Leave Usage
* Balance Distribution
* Leave by Organization
* Leave by Type
* Unpaid Leave
* Expired Balance
* Adjustment History

---

# 31. API Plan

> ⚠️ Endpoint di bawah adalah **rencana**, bukan kondisi aktual. Endpoint yang **sudah ada** di `leave/routes.go` hari ini:
> ```http
> POST/GET      /api/v1/tenant/leave/types
> GET/PUT/DELETE /api/v1/tenant/leave/types/:id
> POST/GET      /api/v1/tenant/leave/accrual-policies
> GET/PUT/DELETE /api/v1/tenant/leave/accrual-policies/:id
> POST/GET      /api/v1/tenant/leave/reasons
> GET/PUT/DELETE /api/v1/tenant/leave/reasons/:id
> POST/GET      /api/v1/tenant/leave/requests
> GET           /api/v1/tenant/leave/requests/:id
> PUT           /api/v1/tenant/leave/requests/:id/status   ← satu-satunya endpoint pengubah status (generik, dipakai untuk submit/approve/reject/cancel sekaligus — belum ada endpoint khusus per aksi)
> DELETE        /api/v1/tenant/leave/requests/:id
> GET           /api/v1/tenant/leave/requests/:id/details
> GET           /api/v1/tenant/leave/balances
> GET           /api/v1/tenant/leave/balances/employees/:employeeId/types/:leaveTypeId
> ```
> Endpoint `submit`, `cancel`, `balances/{employee}/adjust`, dan `team-calendar` **belum ada** — perlu dibuat sebagai endpoint khusus (bukan terus menumpuk di satu endpoint status generik), terutama karena `submit`/`cancel`/`adjust` masing-masing butuh side-effect berbeda (validasi saldo, penulisan ledger, reversal) yang saat ini tidak dilakukan sama sekali oleh `PUT .../status`. **Update 2026-08-09**: `GET /leave/calendar?employee_id=&from=&to=` (Employee Calendar) sudah diimplementasikan — lihat Section 7/Phase 7.

## Leave Types

```http
GET    /api/v1/tenant/leave/types
POST   /api/v1/tenant/leave/types
GET    /api/v1/tenant/leave/types/{id}
PUT    /api/v1/tenant/leave/types/{id}
DELETE /api/v1/tenant/leave/types/{id}
```

## Leave Requests

```http
GET    /api/v1/tenant/leave/requests
POST   /api/v1/tenant/leave/requests
GET    /api/v1/tenant/leave/requests/{id}
PUT    /api/v1/tenant/leave/requests/{id}
POST   /api/v1/tenant/leave/requests/{id}/submit
POST   /api/v1/tenant/leave/requests/{id}/cancel
```

## Balance

```http
GET /api/v1/tenant/leave/balances
GET /api/v1/tenant/leave/balances/{employee}
POST /api/v1/tenant/leave/balances/{employee}/adjust
```

## Calendar

```http
GET /api/v1/tenant/leave/calendar
GET /api/v1/tenant/leave/team-calendar
```

---

# 32. Service Layer

Recommended services:

```text
LeaveRequestService
LeaveBalanceService
LeaveCalculationService
LeaveAccrualService
LeaveEligibilityService
LeaveValidationService
LeaveApprovalService
LeaveCancellationService
LeaveCalendarService
```

---

# 33. Core Business Flow

## Request Leave

```text
Employee
   │
   ▼
Select Leave Type
   │
   ▼
Select Date
   │
   ▼
Calculate Working Days
   │
   ▼
Validate Eligibility
   │
   ▼
Validate Balance
   │
   ▼
Create Leave Request
   │
   ▼
Create Request Details
   │
   ▼
Submit
   │
   ▼
Create Approval Instance
   │
   ▼
Central Approval Module
```

---

# 34. Approval Result

## Approved

```text
Approval Approved
        │
        ▼
Leave Request APPROVED_FINAL
        │
        ▼
Deduct Balance
        │
        ▼
Create USAGE Ledger
        │
        ▼
Generate Attendance Leave
```

## Rejected

```text
Approval Rejected
        │
        ▼
Leave Request REJECTED_FINAL
        │
        ▼
No Balance Deduction
```

---

# 35. Balance Transaction Flow

```text
Leave Approved
      │
      ▼
Balance Before
      │
      ▼
Usage Transaction
      │
      ▼
Balance After
```

Contoh:

```text
Before = 12

Usage = -2

After = 10
```

---

# 36. Seeder

Seeder yang direkomendasikan:

```text
LeaveTypeSeeder
LeaveReasonSeeder
LeaveAccrualPolicySeeder
```

Contoh Leave Types:

```text
ANNUAL
SICK
MATERNITY
PATERNITY
UNPAID
```

Contoh Reasons:

```text
Personal
Medical
Family
Travel
Other
```

Seeder harus mengikuti pola Seeder yang sudah digunakan pada project.

Seluruh ID menggunakan UUID.

---

# 37. Testing Plan

## Unit Test

### Leave Calculation

* Full day
* Half day
* Hourly
* Weekend
* Holiday
* Multiple dates

### Balance

* Accrual
* Usage
* Adjustment
* Reversal
* Carry Forward
* Expiry

### Eligibility

* Eligible employee
* Ineligible employee
* Minimum service period

---

# 38. Feature Test

## Create Leave

```text
Employee
→ Create Request
→ Calculate Duration
→ Submit
```

## Approval

```text
Submit
→ Approval Instance
→ Approver
→ Approve
→ Balance Deducted
```

## Rejection

```text
Submit
→ Reject
→ No Balance Deduction
```

## Cancellation

```text
Approved
→ Cancellation
→ Approval
→ Reversal Balance
```

---

# 39. Important Constraints

### Leave Request

```text
start_date <= end_date
```

### Leave Detail

```text
UNIQUE(
    leave_request_id,
    leave_date
)
```

### Balance

```text
UNIQUE(
    employee_id,
    leave_type_id,
    period_year
)
```

### Approval

Setiap Leave Request yang membutuhkan approval hanya boleh memiliki satu active approval instance.

---

# 40. Recommended Database Final Structure

```text
leave_types
    │
    ├── leave_accrual_policies
    │
    └── leave_requests
            │
            ├── leave_request_details
            │
            └── approval_instance
                    │
                    └── Central Approval Module

leave_reasons

employee_leave_balances
        │
        └── leave_balance_transactions

leave_cancellation_requests (optional)
```

---

# 41. Development Phases

## Phase 1 - Database Enhancement

* Review existing Leave tables.
* Fix data type inconsistencies. ✅ `leave_accrual_policies.deleted_at` was `INT NULL` (mismatched with the Go model's `gorm.DeletedAt`) — fixed to `TIMESTAMP NULL` in migration `070_leave_phase1_db_enhancement`.
* Add missing indexes. ✅ `idx_accrual_deleted_at` added (was the only soft-deletable Leave table without one).
* Add `leave_balance_transactions`. ✅ Table + `LeaveBalanceTransaction` model + `CreateLeaveBalanceTransaction`/`ListLeaveBalanceTransactions` repository methods added. **Kini sudah ditulis** sejak Phase 6 — setiap USAGE/REVERSAL lewat `writeLeaveBalanceTransaction` (catatan asli "nothing writes to it yet" hanya benar saat Phase 1 selesai, sebelum Phase 6 berjalan).
* Add cancellation table if required. ⏳ Deferred — no cancellation flow exists yet (Phase 18/19), revisit when that phase starts.
* Add eligibility tables if required. ⏳ Deferred — no eligibility rules needed yet per current business requirements (Section 22).

---

## Phase 2 - Master Data

* Leave Types ✅
* Leave Reasons ✅
* Accrual Policies ✅
* Leave Eligibility ⏳ Sengaja tidak dibangun — belum ada business rule konkret yang membutuhkannya (lihat §22: "gunakan hanya jika business rule memang diperlukan"). Revisit saat ada requirement nyata (mis. Maternity Leave gender-restricted, minimum service period untuk Annual Leave).

---

## Phase 3 - Leave Calculation Engine

* Working day calculation ✅ `CalculateLeaveDuration` (`leave/calculation.go`) — excludes Sat/Sun and holiday dates, wired into `CreateLeaveRequest` so `requested_days` is now server-computed, not client-trusted.
* Half day ✅ `DurationHalfDayAM/PM` → 0.5 fraction, rejected on non-working days, rejected if `LeaveType.AllowHalfDay = false`.
* Hourly leave ✅ `DurationHourly` → `(end_time - start_time) / 8h`, rejected if `LeaveType.AllowHourly = false`. `8h/day` is a hardcoded default — no per-company `hours_per_day` config exists yet, flagged as a future config point.
* Holiday handling ✅ New `HolidayProvider` interface (`leave/service.go`) + adapter wiring a second `setting.Service` instance in `main.go`, reading `setting.CompanyHoliday` via new `ListHolidayDatesInRange`. Missing/erroring holiday data degrades to "no holidays known" rather than failing the request.
* Shift handling ⏳ **Intentionally not implemented.** Attendance's `AttendanceEmployeeShift.DaysOfWeekMask` has no documented bit-to-weekday interpretation anywhere in the codebase (grep confirmed: the field is only ever passed through, never decoded) — inventing that convention unilaterally inside Leave would be a cross-module design decision, not a Phase 3 calculation task. Standard Sat/Sun + company holidays is the baseline until that convention is established elsewhere.
* Quota calculation 🔶 Duration-mode permission checks against `LeaveType.AllowHalfDay`/`AllowHourly` are done. Balance-quota validation (`requested_days <= available balance`, §12) is **not** part of this phase — that belongs to Phase 6 (Leave Balance) since it needs the balance-deduction logic that doesn't exist yet.
* `LeaveRequestDetail` rows ✅ now actually created per working date on submit (previously the repository methods existed but were never called from `CreateLeaveRequest`).

---

## Phase 4 - Leave Request

* Create Request ✅
* Draft ✅ (status exists, not separately exercised by a dedicated endpoint — creation goes straight to `SUBMITTED`)
* Submit ✅ (via `CreateLeaveRequest`; no separate draft→submit transition endpoint yet — see §31 gap)
* Validation 🔶 Implemented: leave type active check, required-attachment check, date-overlap check against the employee's own non-final requests (`CountOverlappingLeaveRequests`). **Not** implemented: employee/organization-active check (needs a cross-module `Employee` read no other Leave code currently does — deferred, not designed yet), backdate/minimum-notice rules (no `allow_backdate`/`minimum_notice_days` fields exist on `LeaveType` yet, §3), balance-quota check (`requested_days <= available balance` — belongs to Phase 6, needs balance-deduction logic first)
* Attachment ✅ Enforced when `LeaveType.RequiresAttachment = true`
* Request Detail ✅ (done in Phase 3 — `LeaveRequestDetail` rows created per working date on submit)

---

## Phase 5 - Approval Integration

Integrate dengan Central Approval Module.

* Create Approval Instance
* Resolve Approver
* Approval Task
* Approval Action
* Approval Result
* Notification

---

## Phase 6 - Leave Balance

* Accrual ⏳ Not implemented — no code seeds an ACCRUAL ledger entry or `employee_leave_balances.quota_days` from `LeaveAccrualPolicy`; a balance row is only lazily created (quota seeded from `LeaveType.DefaultQuotaDays`) the first time a USAGE deduction needs one.
* Usage ✅ `applyLeaveUsage` (`leave/balance.go`) deducts `RequestedDays` from `employee_leave_balances` and writes a USAGE ledger row, triggered from both `HandleApprovalStatusChange` (push-callback path) and `UpdateLeaveRequestStatus` (manual generic status endpoint) whenever a request transitions into `APPROVED_FINAL`. Skipped entirely when `LeaveType.CountsAgainstQuota = false` (§12).
* Adjustment ⏳ Not implemented — no HR-initiated balance adjustment endpoint/flow exists yet (§26 stays a proposal).
* Reversal ✅ `reverseLeaveUsage` (`leave/balance.go`) restores the deducted days and writes a REVERSAL ledger row, triggered from `UpdateLeaveRequestStatus` whenever a request transitions **away from** `APPROVED_FINAL` (e.g. cancelled by HR after approval). Note: `HandleApprovalStatusChange`'s own CANCELLED branch only fires from `PENDING_APPROVAL`, where the balance was never deducted, so no reversal is needed there — full post-approval cancellation as its own workflow (§18/§19) is still not built.
* Carry Forward ⏳ Not implemented (§27 stays a proposal).
* Expiry ⏳ Not implemented (§28 stays a proposal).
* Ledger ✅ Every USAGE/REVERSAL write goes through `writeLeaveBalanceTransaction`, recording `balance_before`/`balance_after` on `leave_balance_transactions` (table existed since Phase 1, now actually written to).

---

## Phase 7 - Calendar & Attendance 🔶 Sebagian (2026-08-09)

* Employee Calendar ✅ Diimplementasikan — lihat catatan di bawah
* Team Calendar ⏳ Deferred — butuh cross-module employee/organization read ("siapa bawahan siapa") yang tidak ada interface-nya di manapun di codebase ini, kategori gap yang sama dengan Manager/HR Dashboard di Attendance (Phase 10)
* Organization Calendar ⏳ Deferred — sama seperti Team Calendar
* Attendance Integration ✅ Sudah selesai sejak sebelumnya (bukan pekerjaan baru Phase 7 ini) — `leave.AttendanceSessionUpdater`/`SetAttendanceSessionUpdater`, diwire di `main.go` via `leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)`, lihat catatan Attendance backend Phase 9 (`docs/module-attendance-plan.md`)

> ✅ **Employee Calendar diimplementasikan.** `GET /api/v1/tenant/leave/calendar?employee_id=&from=&to=` (`leave/handler.go` `GetEmployeeCalendar`, `leave/service.go` `Service.GetEmployeeCalendar`) — mengembalikan `[]CalendarEntryResponse` (leave_request_id, leave_date, day_fraction, leave_type_id, status) dari `leave_request_details` JOIN `leave_requests` (`Repository.FindCalendarEntriesForEmployeeInRange`), pola yang identik dengan `attendance.Service.GetEmployeeCalendar`. Request berstatus `REJECTED_FINAL` sengaja dikecualikan dari hasil kalender — tanggal yang ditolak tidak pernah benar-benar jadi leave day, konsisten dengan prinsip §26 ("Pending Leave tidak boleh dianggap sebagai approved leave"). `CANCELLED` tetap ditampilkan supaya karyawan tahu kenapa tanggal itu kembali kosong.
>
> **Team/Organization Calendar sengaja tidak dibangun** — keduanya butuh mengetahui struktur organisasi/siapa bawahan siapa, sebuah cross-module read ke employee/organization yang tidak ada interface-nya di manapun di codebase ini (kategori gap yang identik dengan Manager/HR Dashboard Attendance, Phase 10). Employee Calendar tidak butuh itu karena `employee_id` sudah diberikan langsung oleh caller (sama seperti alasan Attendance's employee-level calendar tidak butuh cross-module read).
>
> Test: `calendar_test.go` — entries dalam rentang tanggal, dan request yang REJECTED_FINAL dikecualikan dari hasil.
>
> Bug driver ditemukan & diperbaiki saat implementasi: kolom `leave_date` (`type:date`) round-trip sebagai RFC3339 penuh di sqlite test driver (quirk yang sama seperti `submitted_at`/`work_date` yang sudah didokumentasikan di `balance.go`/Attendance `session.go`) — ditambahkan `normalizeLeaveDate` (`calculation.go`) mengikuti pola `normalizeWorkDate` Attendance.

---

## Phase 8 - Notification 🔶 Sebagian (2026-08-09)

* Request Submitted ⏳ Sengaja tidak dibangun — pemohon adalah aktor yang melakukan submit itu sendiri, feedback langsung dari response API sudah cukup; self-notification tidak menambah nilai
* Approval Required ⏳ Sengaja tidak dibangun — ini tanggung jawab modul Approval (pemberitahuan ke approver soal task baru), bukan Leave; Leave hanya membuat approval instance, tidak tahu/tidak seharusnya tahu siapa approver-nya
* Approved ✅ Sudah ada sejak sebelumnya (`notifyLeaveOutcome`, `LEAVE_APPROVED`) — kini juga terpicu dari `UpdateLeaveRequestStatus` (endpoint status generik/manual), tidak hanya dari `HandleApprovalStatusChange` (push-callback)
* Rejected ✅ Sama seperti Approved (`LEAVE_REJECTED`), kini juga dari kedua jalur
* Cancelled ✅ **Baru diimplementasikan** — `LEAVE_CANCELLED` ditambahkan ke `notifyLeaveOutcome`, terpicu dari kedua jalur (`HandleApprovalStatusChange` status CANCELLED, dan `UpdateLeaveRequestStatus` transisi ke CANCELLED)
* Reminder ⏳ Sengaja tidak dibangun — butuh scheduled job (cek approval yang overdue), tidak ada infrastruktur cron di codebase ini, kategori gap yang sama dengan Missing Checkout Reminder di Attendance (Phase 12)
* Balance Notification ⏳ Sengaja tidak dibangun — tidak ada fitur balance adjustment (§26 tetap proposal, Phase 6), jadi tidak ada event untuk dinotifikasi

> ✅ **Gap nyata yang ditemukan & diperbaiki**: sebelum perubahan ini, `notifyLeaveOutcome` hanya pernah dipanggil dari `HandleApprovalStatusChange` (jalur push-callback Central Approval Module) — transisi status lewat `UpdateLeaveRequestStatus` (endpoint generik `PUT /requests/:id/status`, dipakai HR untuk override manual) sama sekali tidak memicu notifikasi apapun, padahal efek balance-nya (`applyBalanceEffectOnStatusChange`) sudah konsisten di kedua jalur sejak Phase 6. Ditambahkan pemanggilan `notifyLeaveOutcome` di `UpdateLeaveRequestStatus` juga (hanya saat status benar-benar berubah, ke `APPROVED_FINAL`/`REJECTED_FINAL`/`CANCELLED`), supaya notifikasi konsisten dengan jalur mana pun yang men-drive transisi — sejalan dengan §34 yang menyebut "Approval Result" sebagai efek dari perubahan status leave request, bukan efek dari endpoint tertentu.
>
> Test: `notifier_integration_test.go` — notify saat CANCELLED via push-callback, notify saat APPROVED via endpoint status manual, dan tidak ada notify untuk transisi status yang sebenarnya no-op (status lama == status baru).

---

## Phase 9 - Dashboard & Reports 🔶 Sebagian (2026-08-09)

* Employee Dashboard ✅ Tidak butuh endpoint baru — sudah bisa disusun penuh dari endpoint yang ada (`GET /balances?employee_id=` untuk quota/used/remaining, `GET /requests?employee_id=` untuk upcoming/pending/recent history), murni komposisi di FE
* Manager Dashboard ⏳ Deferred — butuh cross-module employee/organization read ("siapa bawahan siapa") yang tidak ada interface-nya di manapun di codebase ini, kategori gap yang sama dengan Manager/HR Dashboard Attendance (Phase 10) dan Team/Organization Calendar Leave (Phase 7)
* HR Dashboard ⏳ Deferred — sama seperti Manager Dashboard
* Leave Usage Report ✅ Diimplementasikan — lihat catatan di bawah
* Balance Report ✅ Tidak butuh endpoint baru — `GET /balances` tanpa filter `employee_id` sudah tenant-wide sejak awal
* Organization Report ⏳ Deferred — sama seperti Manager/HR Dashboard, butuh cross-module organization read
* Leave Type Report ✅ Tidak butuh endpoint terpisah — data yang sama dari Leave Usage Report, dikelompokkan per `leave_type_id` di sisi consumer (pola yang sama dipakai Attendance untuk Late/Early Leave/Missing Attendance "reports", Phase 11)

> ✅ **Leave Usage Report diimplementasikan.** `GET /api/v1/tenant/leave/reports/usage?from=&to=` (`leave/handler.go` `GetLeaveUsageReport`, `leave/service.go` `Service.GetLeaveUsageReport`) — mengembalikan `[]LeaveRequestResponse` tenant-wide (semua employee) untuk request yang date range-nya overlap dengan `[from, to]` (`Repository.FindLeaveRequestsInRange`), pola yang identik dengan `attendance.Service.GetAttendanceReport` (Attendance Phase 11). **Request `REJECTED_FINAL`/`CANCELLED` sengaja TIDAK dikecualikan** di sini — beda dari Employee Calendar (Phase 7) yang mengecualikan `REJECTED_FINAL` karena kalender menunjukkan hari yang benar-benar jadi leave day, sedangkan usage report justru harus menunjukkan gambaran lengkap termasuk apa yang ditolak/dibatalkan.
>
> **Manager Dashboard, HR Dashboard, dan Organization Report sengaja tidak dibangun** — ketiganya butuh mengetahui struktur organisasi/siapa bawahan siapa, cross-module read ke employee/organization yang tidak ada interface-nya di manapun di codebase ini (kategori gap yang sama persis dengan Attendance Phase 10 dan Leave Phase 7). **Employee Dashboard dan Balance Report tidak butuh pekerjaan backend baru** — keduanya sudah bisa dipenuhi dari endpoint existing (`GET /balances`, `GET /requests`), murni pekerjaan FE komposisi, bukan gap backend.
>
> Test: `report_test.go` — request tenant-wide dalam rentang tanggal (lintas employee), dan request `REJECTED_FINAL` tetap muncul di hasil (beda perilaku dari kalender). Fixture ditulis lewat repository langsung (`seedLeaveRequest`) untuk menghindari quirk driver sqlite test yang sama seperti `calendar_test.go`.

---

## Phase 10 - Testing

* Unit Tests
* Feature Tests
* Approval Tests
* Balance Tests
* Calculation Tests
* Integration Tests

---

# 42. Priority

| Feature                | Priority |
| ---------------------- | -------- |
| Leave Type             | P0       |
| Leave Request          | P0       |
| Leave Calculation      | P0       |
| Leave Balance          | P0       |
| Approval Integration   | P0       |
| Request Detail         | P0       |
| Balance Ledger         | P0       |
| Leave Calendar         | P1       |
| Notification           | P1       |
| Attendance Integration | P1       |
| Leave Adjustment       | P1       |
| Carry Forward          | P1       |
| Payroll Integration    | P2       |
| Eligibility Rules      | P2       |
| Team Capacity Rule     | P2       |
| Advanced Analytics     | P3       |

---

# 43. Final Architecture

```text
                    ┌──────────────────────┐
                    │   Leave Management   │
                    └──────────┬───────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
        Leave Master      Leave Request     Leave Balance
              │                │                │
              │                ▼                ▼
              │        Approval Module      Balance Ledger
              │                │
              │                ▼
              │           Notification
              │
              ▼
       Leave Calculation
              │
              ├──────────► Attendance
              │
              └──────────► Payroll
```

# Implementation Status

Diverifikasi langsung terhadap kode per 2026-08-09.

| Phase (§41) | Status | Catatan |
|---|---|---|
| Phase 1 - Database Enhancement | ✅ Selesai (2026-08-08) | Migration `070_leave_phase1_db_enhancement` (mysql+postgres): fix `leave_accrual_policies.deleted_at` (INT → TIMESTAMP) + index, tambah tabel `leave_balance_transactions` + model `LeaveBalanceTransaction` + repository `CreateLeaveBalanceTransaction`/`ListLeaveBalanceTransactions`. Tabel cancellation (§19) dan eligibility (§22) sengaja ditunda — tidak dibutuhkan sampai fitur terkait (Phase 6 lanjutan/§18, §22) mulai dikerjakan |
| Phase 2 - Master Data | ✅ Selesai (2026-08-08) | Leave Types, Accrual Policies, Leave Reasons — CRUD lengkap di `leave/service.go`, `leave/handler.go`, `leave/routes.go`. Leave Eligibility sengaja tidak dibangun — tidak ada business rule konkret yang membutuhkannya saat ini (§22), revisit saat requirement muncul |
| Phase 3 - Leave Calculation Engine | ✅ Selesai (2026-08-08) | `leave/calculation.go` (`CalculateLeaveDuration`): full-day/half-day/hourly calculation, weekend + company-holiday exclusion via new `HolidayProvider` interface (adapter: `setting.Service.ListHolidayDatesInRange`), `LeaveRequestDetail` rows now actually persisted per date. Shift/`DaysOfWeekMask` handling and balance-quota validation intentionally deferred — see §41 |
| Phase 4 - Leave Request | 🔶 Sebagian (2026-08-08) | Create/Submit/List/Get/Delete/Details sudah ada, `requested_days` dihitung server-side (Phase 3). Validasi baru: leave type aktif, attachment wajib, overlap tanggal (`CountOverlappingLeaveRequests`). Belum ada: employee/organization aktif (butuh cross-module read ke `employee`, belum ada pola/interface untuk ini), backdate/minimum-notice (field belum ada di `LeaveType`), balance-quota check (nunggu Phase 6) |
| Phase 5 - Approval Integration | ✅ Selesai | `ApprovalEngine`, `SetApprovalEngine`, `HandleApprovalStatusChange`, wiring `main.go`, module slug tunggal `"leave"`, test coverage di `approval_integration_test.go` — lihat Section 7 |
| Phase 6 - Leave Balance | 🔶 Sebagian (2026-08-08) | **Usage + Reversal + Ledger selesai**: `leave/balance.go` (`applyLeaveUsage`/`reverseLeaveUsage`) — deduct saldo saat status masuk `APPROVED_FINAL`, reverse saat keluar dari `APPROVED_FINAL` (mis. cancel setelah approve), keduanya menulis ke `leave_balance_transactions`. Wired di `HandleApprovalStatusChange` dan `UpdateLeaveRequestStatus`. **Belum ada**: Accrual (seed quota dari `LeaveAccrualPolicy`), Adjustment (endpoint HR), Carry Forward, Expiry |
| Phase 7 - Calendar & Attendance | 🔶 Sebagian (2026-08-09) | Employee Calendar selesai (`GET /leave/calendar`). Attendance Integration sudah selesai sejak sebelumnya (`leave.AttendanceSessionUpdater`). Team/Organization Calendar sengaja ditunda — butuh cross-module employee/organization read yang belum ada |
| Phase 8 - Notification | 🔶 Sebagian (2026-08-09) | Approved/Rejected/Cancelled notification lewat `notifyLeaveOutcome`, kini konsisten di kedua jalur (`HandleApprovalStatusChange` push-callback dan `UpdateLeaveRequestStatus` manual). Request Submitted/Approval Required sengaja tidak dibangun (bukan tanggung jawab Leave), Reminder butuh scheduled job (belum ada infra), Balance Notification butuh fitur adjustment (belum ada) |
| Phase 9 - Dashboard & Reports | 🔶 Sebagian (2026-08-09) | Leave Usage Report selesai (`GET /leave/reports/usage`, tenant-wide). Employee Dashboard & Balance Report tidak butuh endpoint baru (sudah bisa disusun dari `GET /balances`/`GET /requests`). Manager/HR Dashboard & Organization Report sengaja ditunda — butuh cross-module employee/organization read yang belum ada |
| Phase 10 - Testing | 🔶 Sebagian | Test approval-integration sudah ada; test kalkulasi/balance/cancellation belum ada karena fiturnya sendiri belum ada |

**Frontend**: ✅ **FE-1 & FE-2 selesai (2026-08-09)** — `Leave.vue` (My Leave dashboard: balance cards + list request + dialog New Request + kalender bulan berjalan), `LeaveAdmin.vue` (index kartu), `LeaveTypes.vue`, `LeaveAccrualPolicies.vue`, `LeaveReasons.vue` (CRUD Dialog inline). Bilingual EN/ID di `locales/en.json` & `locales/id.json`, `npm run build` bersih, diverifikasi manual di browser tanpa console error. Team/Manager/HR Dashboard & Reports tetap backend-blocked (butuh cross-module employee/organization read yang belum ada).

**Rekomendasi urutan lanjutan** (sesuai prioritas P0 di §42, disesuaikan dengan gap yang sudah dikonfirmasi):
1. ~~Leave Balance auto-deduct on approve + reversal on cancel~~ ✅ Selesai (2026-08-08, Phase 6) — `applyLeaveUsage`/`reverseLeaveUsage` + ledger `leave_balance_transactions`.
2. ~~Working-day calculation server-side~~ ✅ Selesai (2026-08-08, Phase 3) — `CalculateLeaveDuration` wired ke `CreateLeaveRequest`, `requested_days` dihitung server-side bukan dipercaya dari client.
3. Endpoint `submit`/`cancel` khusus (pisah dari `PUT .../status` generik) — masih belum ada; side-effect balance/ledger sudah aman di jalur `UpdateLeaveRequestStatus`, endpoint khusus tinggal menyederhanakan UX/audit.
4. ~~Frontend dasar: My Leave (request + list + balance)~~ ✅ Selesai (2026-08-09) — FE-1 & FE-2 di Frontend Implementation Plan di bawah.
5. ~~Balance ledger table (§10)~~ ✅ Selesai (2026-08-08, Phase 6) — `leave_balance_transactions` sudah ditulis untuk USAGE/REVERSAL sejak auto-deduct berjalan.

---

# Frontend Implementation Plan ✅ Selesai (FE-1 & FE-2, 2026-08-09)

Ditambahkan 2026-08-09. Backend Leave sudah cukup matang (lihat Implementation Status di atas: Master Data, Calculation Engine, Approval Integration, Balance Usage/Reversal/Ledger semuanya selesai) — sisi frontend awalnya 0%: `frontend/tenant/src/views/modules/Leave.vue` hanya placeholder satu baris ("Leave Module — Coming soon"), sama seperti Attendance sebelum FE-nya dibangun. Section ini mengikuti pola yang sama persis dengan Frontend Implementation Plan Attendance (`docs/module-attendance-plan.md`) — konvensi FE sudah baku, bukan didesain ulang dari nol.

> ✅ **FE-1 (My Leave Dashboard) dan FE-2 (Admin Configuration) sudah diimplementasikan per 2026-08-09.** 5 halaman baru/diisi: `Leave.vue` (My Leave dashboard: balance cards per leave type, list request sendiri + cancel, dialog New Request, kalender bulan berjalan via `GET /leave/calendar`), `LeaveAdmin.vue` (index kartu, pola `AttendanceAdmin.vue`), `LeaveTypes.vue`, `LeaveAccrualPolicies.vue`, `LeaveReasons.vue` (CRUD Dialog inline, pola compact `NationalitiesView.vue`). Route terdaftar: `leave`, `leave/admin`, `leave/types`, `leave/accrual-policies`, `leave/reasons` (`meta.module: 'leave'`). Bilingual EN/ID di `locales/en.json`/`locales/id.json`, `npm run build` bersih (2.2s, zero warning), diverifikasi manual di browser (login tenant `andi.wijaya@test.local`): halaman Leave & Leave Admin render tanpa console error, kolom tabel sesuai rencana. Yang tetap di luar cakupan (lihat "Eksplisit di luar cakupan rencana FE ini" di bawah): Team/Organization Calendar, Manager/HR Dashboard, Leave Reports, Balance Adjustment UI, Attendance Integration UI, notification bell wiring — semuanya blocked oleh gap backend/modul lain, bukan keputusan FE.

## FE-1. Ringkasan & Prinsip

* **Pola FE mengikuti persis apa yang sudah divalidasi di Attendance FE** (`docs/module-attendance-plan.md`, sudah selesai FE-1 s.d. FE-5): reuse `Employees.vue` untuk list+pagination, pola compact `NationalitiesView.vue` (DataTable + `Dialog` inline) untuk CRUD entitas sederhana, `services/api.js`/`services/responseHandler.js` langsung tanpa service-layer/Pinia store baru per modul, `useAuth().hasPermission(slug)` untuk gating tombol.
* **Reuse langsung composable/util yang sudah dibuat untuk Attendance FE** — keduanya generik, tidak spesifik Attendance:
  - `composables/useMyEmployee.js` (`GET /user-accounts/me`, cache module-level) — dibutuhkan Leave persis dengan alasan yang sama seperti Attendance: `CreateLeaveRequest`/`ListLeaveRequests`/`ListLeaveBalances` semuanya butuh `employee_id` yang harus resolve dari user yang login, bukan endpoint "my-request" tersendiri.
  - `utils/localTime.js` — **tidak dibutuhkan Leave** kecuali untuk `duration_mode = HOURLY` (`start_time`/`end_time`, lihat dto.go `CreateLeaveRequest`); field ini string bebas format tanpa binding RFC3339 eksplisit di backend (beda dari Attendance overtime), jadi cukup `HH:mm` biasa — dicatat di sini supaya tidak salah asumsi format saat implementasi.
* **Approval TIDAK dibangun ulang di Leave**, sama seperti Attendance. Bedanya: Leave punya permission `leave.approve` terpisah di `Info().Permissions` (`leave/module.go`) — Attendance tidak. Namun tidak ada endpoint approve/reject khusus di `leave/routes.go` (hanya `PUT /requests/:id/status` generik yang di-drive oleh Central Approval Module lewat push-callback, lihat Section 7) — jadi permission ini kemungkinan besar dipakai untuk gating siapa yang boleh jadi approver di flow config Approval Module, bukan untuk tombol approve di FE Leave. FE Leave tetap murni link-out ke `/approvals`, tidak mencoba membuat tombol approve sendiri.
* **Balance ditampilkan read-only** — tidak ada endpoint adjustment (§26 tetap proposal, dikonfirmasi di Implementation Status Phase 6), jadi FE tidak boleh menampilkan form "Adjust Balance" yang akan gagal karena endpoint-nya tidak ada.
* **Cancellation lewat endpoint status generik** — `PUT /requests/:id/status` dengan `status: "CANCELLED"` adalah satu-satunya mekanisme yang ada (§18 dikonfirmasi belum ada sub-flow cancellation khusus/reversal approval). FE cukup memanggil endpoint ini untuk request yang masih `DRAFT`/`SUBMITTED`/`PENDING_APPROVAL` milik sendiri — bukan membangun UI cancellation-request-dengan-approval yang backend-nya tidak ada.

## FE-2. Halaman & Routing

Mengikuti pola nomenklatur route Attendance (`attendance/...` → `leave/...`), sibling di bawah path `leave/...`, `meta.module: 'attendance'` → `meta.module: 'leave'`, dst.

| Halaman | Route | Endpoint backend | Permission | Status |
|---|---|---|---|---|
| `Leave.vue` (My Leave dashboard: balance cards + list request saya + tombol "New Request" + kalender bulan berjalan) | `leave` | `GET /balances?employee_id=`, `GET /requests?employee_id=`, `GET /calendar?employee_id=&from=&to=` | `leave.view` | ✅ Selesai |
| `LeaveRequestForm.vue` (create request, Dialog atau halaman terpisah — lihat FE-3) | inline di `Leave.vue` | `POST /requests` | `leave.create` | ✅ Selesai (Dialog inline) |
| `LeaveAdmin.vue` (index kartu, pola sama `AttendanceAdmin.vue`) | `leave/admin` | - | `leave.update` | ✅ Selesai |
| `LeaveTypes.vue` + Dialog | `leave/types` | `POST/GET /types`, `GET/PUT/DELETE /types/:id` | `leave.view`/`create`/`update`/`delete` | ✅ Selesai |
| `LeaveAccrualPolicies.vue` + Dialog | `leave/accrual-policies` | `POST/GET /accrual-policies`, `GET/PUT/DELETE /accrual-policies/:id` | sama pola di atas | ✅ Selesai |
| `LeaveReasons.vue` + Dialog | `leave/reasons` | `POST/GET /reasons`, `GET/PUT/DELETE /reasons/:id` | sama pola di atas | ✅ Selesai |

Catatan: **tidak ada halaman `TeamLeave`/`LeaveDashboard`(Manager/HR)/`LeaveReports`** — masih backend-blocked (Implementation Status Phase 7/9). **Update 2026-08-09**: `GET /leave/calendar?employee_id=&from=&to=` (Employee Calendar) sudah diimplementasikan backend-side dan **sudah dipakai di `Leave.vue`** (section "My Leave This Month", pola list tanggal + Tag seperti `Attendance.vue`). Team/Organization Calendar tetap backend-blocked.

## FE-3. Development Phases (FE)

**Phase FE-1 — My Leave Dashboard ✅ Selesai (2026-08-09)**
`Leave.vue`: kartu balance per leave type (`quota_days`/`used_days`/`remaining_days` dari `GET /balances?employee_id=`), list request milik sendiri berstatus apapun (`GET /requests?employee_id=`, tabel dengan kolom leave_type/tanggal/status/requested_days + tombol Cancel untuk `DRAFT`/`SUBMITTED`/`PENDING_APPROVAL` via `PUT /requests/:id/status` `CANCELLED`), tombol "New Request" membuka Dialog form (leave type dropdown, date range via `DateInput` x2, `duration_mode` select incl. `HOURLY` → input `start_time`/`end_time` `HH:mm`, `leave_reason_id` dropdown dari `GET /reasons`, textarea note, input `attachment_url` jika `LeaveType.RequiresAttachment` — lihat catatan attachment di FE-4). `requested_days` **tidak dihitung di FE** — backend sudah menghitungnya server-side (`CalculateLeaveDuration`, Leave backend Phase 3), FE cukup kirim `request_start_date`/`request_end_date`/`duration_mode` dan tampilkan `requested_days` dari response. Tambahan sesuai update FE-2: section "My Leave This Month" memakai `GET /leave/calendar` (Employee Calendar).

> Catatan implementasi: `employee_id` di-resolve dari `composables/useMyEmployee.js` (`GET /user-accounts/me`) — persis alasan yang sama seperti Attendance FE. Gating tombol pakai `useAuth().hasPermission`: Admin (`leave.update`), New Request (`leave.create`). Setelah create/cancel, balance + requests + calendar di-reload bersamaan.

**Phase FE-2 — Admin Configuration ✅ Selesai (2026-08-09)**
`LeaveAdmin.vue` (index kartu, pola `AttendanceAdmin.vue`) + `LeaveTypes.vue`, `LeaveAccrualPolicies.vue`, `LeaveReasons.vue` — CRUD standar Dialog inline (pola compact `NationalitiesView.vue`), field sesuai `CreateLeaveTypeRequest`/`CreateAccrualPolicyRequest`/`CreateLeaveReasonRequest` (`leave/dto.go`). `LeaveAccrualPolicies.vue` memakai dropdown Leave Type (fetch dari `/types`), field `effective_from`/`effective_to` via `DateInput`. `LeaveReasons.vue` tanpa pagination server-side — `GET /reasons` mengembalikan array polos (`{success, data}`) bukan amplop paginated. `LeaveTypes.vue`/`LeaveAccrualPolicies.vue` pakai `DataTable` `lazy` + pagination server-side (`page`/`per_page`/`total`).

**Eksplisit di luar cakupan rencana FE ini** (semuanya backend-blocked, bukan keputusan FE):
* **Team Calendar / Organization Calendar** — tidak ada endpoint (`GET /team-calendar` di §31 API Plan tetap proposal murni). **Employee Calendar** (`GET /leave/calendar`) sudah tersedia sejak 2026-08-09 — lihat catatan update di FE-2 di atas, bukan lagi backend-blocked.
* **Manager Dashboard, HR Dashboard** (§30 backend Phase 9) — tidak ada endpoint, butuh cross-module employee/organization read. **Update 2026-08-09**: `GET /leave/reports/usage?from=&to=` (Leave Usage Report, tenant-wide) sudah diimplementasikan backend-side — halaman FE Report belum dibangun pada penulisan catatan ini, tapi tidak lagi backend-blocked.
* **Balance Adjustment UI** — tidak ada endpoint `POST /balances/{employee}/adjust` (§26 tetap proposal).
* **Attendance Integration UI** — tidak ada yang perlu ditampilkan di FE Leave; integrasi ini sudah selesai di sisi backend (Leave backend Phase 9 — `leave.AttendanceSessionUpdater`) dan hasilnya muncul di FE **Attendance**, bukan FE Leave.
* **Notification bell wiring** — sama seperti catatan Attendance FE, di luar cakupan modul ini (scope `docs/module-notification-plan.md`). Catatan: Leave justru sudah jadi Notifier consumer **pertama** di backend (`leave.Notifier`/`SetNotifier`, `docs/module-notification-plan.md` Phase 4 — mendahului Attendance's Phase 5 rollout), tapi notifikasi itu baru benar-benar terlihat pengguna setelah bell-nya jadi dropdown fungsional, yang tetap di luar cakupan FE plan ini.

## FE-4. Catatan Teknis

* **Response envelope & tabel/pagination**: identik dengan Attendance FE — `httputil.SuccessJSON`/paginated response yang sama, `DataTable` `lazy` + server pagination, `SkeletonTable`, debounce search — semua meniru `Employees.vue`.
* **Attachment upload** (`attachment_url` di `CreateLeaveRequest`) — ✅ **Investigasi selesai saat implementasi FE-1**: tidak ada endpoint upload generik di codebase ini (hanya upload employee-scoped: `PUT /employees/:id/photo` dan `POST/PUT /employees/:id/documents/upload`, keduanya terikat resource employee, tidak cocok untuk lampiran leave). Sesuai disiplin "jangan reinvent" — tanpa membangun endpoint upload baru di luar scope FE — dialog New Request memakai input teks `attachment_url` (ditampilkan hanya jika tipe cuti terpilih `requires_attachment`). Backend tetap menerima string URL apa pun; jika upload file generik dibutuhkan, itu endpoint backend baru di luar rencana FE ini.
* **`duration_mode = HOURLY`**: `start_time`/`end_time` di `CreateLeaveRequest` adalah string bebas (bukan RFC3339 seperti Attendance's `event_time_local`/overtime) — cukup input teks/time-picker `HH:mm`, tidak perlu `utils/localTime.js`.
* **Validasi overlap & balance**: backend sudah menolak overlap tanggal (`CountOverlappingLeaveRequests`) dan (begitu Phase 6 balance-quota check ada) saldo tidak cukup — FE cukup menampilkan error validasi apa adanya lewat `getValidationErrors`/`isValidationError`, tidak mengecek ulang di client.

1. Leave Module bertanggung jawab terhadap **business rule cuti**.
2. Approval dilakukan oleh **Central Approval Module**.
3. Leave tidak memiliki approval engine sendiri.
4. `employee_leave_balances` menyimpan current balance.
5. `leave_balance_transactions` menyimpan histori perubahan saldo.
6. Setiap perubahan saldo harus menghasilkan ledger transaction.
7. Approved Leave baru boleh mengurangi saldo.
8. Rejected Leave tidak mengurangi saldo.
9. Cancellation setelah deduction harus menghasilkan reversal.
10. Perhitungan hari harus mempertimbangkan working calendar.
11. Seluruh ID menggunakan UUID.
12. Histori transaksi tidak boleh dihapus secara fisik.
13. Configuration dibuat melalui Settings.
14. Employee hanya dapat mengajukan Leave sesuai eligibility.
15. Organization/Manager dapat melihat leave bawahannya berdasarkan struktur Organization.
16. Approval flow dikonfigurasi secara terpusat sehingga dapat digunakan oleh seluruh modul HRIS.
17. Leave Type, Accrual Policy, dan Eligibility tidak boleh mencampurkan business logic dengan UI.
18. Current balance dan transaction ledger harus konsisten melalui database transaction.
19. Setiap approval dan perubahan saldo harus dapat diaudit.
20. Struktur Leave harus dapat digunakan kembali oleh Attendance dan Payroll.
