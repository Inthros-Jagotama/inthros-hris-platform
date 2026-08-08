Saya sudah pelajari struktur `005_leave.sql`. Struktur dasarnya sudah mencakup **Leave Type, Accrual Policy, Reason, Request, Request Detail, dan Balance**. Dari struktur tersebut, saya akan mempertahankan pola yang ada dan menambahkan kebutuhan development di atasnya, terutama karena Anda sudah memiliki **Approval Module terpusat**.

Berikut plan Markdown yang bisa langsung Anda simpan sebagai `performance-management-leave-plan.md` atau `leave-management-plan.md`.

> ⚠️ **Status vs. Plan ini**: dokumen ini awalnya ditulis seolah modul Leave belum ada (greenfield). Setelah dicek ulang terhadap kode aktual, **backend modul Leave sudah cukup jauh diimplementasikan** — termasuk integrasi ke Central Approval Module — sementara beberapa bagian rencana (balance ledger, cancellation flow, working-day calculation, calendar, accrual engine, notification) **memang belum ada**, dan **frontend masih placeholder "Coming soon"**. Lihat section **"Implementation Status"** di bagian bawah dokumen untuk status per-fase yang sudah diverifikasi terhadap kode, dan catatan blockquote (`>`) yang disisipkan di beberapa section untuk koreksi spesifik.

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
> ❌ **Gap nyata, prioritas P0**: tidak ada satupun kode yang menulis/mengurangi `used_days`/`remaining_days` saat leave di-approve, atau mengembalikannya saat dibatalkan. Approval flow saat ini "selesai" tanpa pernah menyentuh saldo yang seharusnya ia lindungi. Section 10 (Balance Ledger) di bawah masih berupa proposal — belum ada implementasinya sama sekali, bukan cuma "enhancement" opsional.

---

# 10. Leave Balance Enhancement

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

> ❌ **Belum ada sama sekali** di luar nilai status `CANCELLED` yang bisa di-set lewat endpoint generik `PUT /requests/:id/status`. Tidak ada sub-flow request/approval cancellation, tidak ada reversal saldo (karena saldo sendiri belum pernah dikurangi — lihat Section 9/10), tidak ada endpoint `cancel` khusus. Desain di bawah ini masih sepenuhnya proposal.

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

> ❌ **Frontend belum dibangun.** `frontend/tenant/src/views/modules/Leave.vue` saat ini hanya stub placeholder ("Leave Module — Coming soon"). Tidak ada halaman request list, request form, balance view, calendar, atau settings — seluruh struktur menu di bawah ini masih 0% dikerjakan.

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
> Endpoint `submit`, `cancel`, `balances/{employee}/adjust`, dan seluruh Calendar di bawah **belum ada** — perlu dibuat sebagai endpoint khusus (bukan terus menumpuk di satu endpoint status generik), terutama karena `submit`/`cancel`/`adjust` masing-masing butuh side-effect berbeda (validasi saldo, penulisan ledger, reversal) yang saat ini tidak dilakukan sama sekali oleh `PUT .../status`.

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
* Add `leave_balance_transactions`. ✅ Table + `LeaveBalanceTransaction` model + `CreateLeaveBalanceTransaction`/`ListLeaveBalanceTransactions` repository methods added. Nothing writes to it yet — that's Phase 6 (accrual/usage/adjustment business logic).
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

* Create Request
* Draft
* Submit
* Validation
* Attachment
* Request Detail

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

* Accrual
* Usage
* Adjustment
* Reversal
* Carry Forward
* Expiry
* Ledger

---

## Phase 7 - Calendar & Attendance

* Employee Calendar
* Team Calendar
* Organization Calendar
* Attendance Integration

---

## Phase 8 - Notification

* Request Submitted
* Approval Required
* Approved
* Rejected
* Cancelled
* Reminder
* Balance Notification

---

## Phase 9 - Dashboard & Reports

* Employee Dashboard
* Manager Dashboard
* HR Dashboard
* Leave Usage Report
* Balance Report
* Organization Report
* Leave Type Report

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

Diverifikasi langsung terhadap kode per 2026-08-08.

| Phase (§41) | Status | Catatan |
|---|---|---|
| Phase 1 - Database Enhancement | ✅ Selesai (2026-08-08) | Migration `070_leave_phase1_db_enhancement` (mysql+postgres): fix `leave_accrual_policies.deleted_at` (INT → TIMESTAMP) + index, tambah tabel `leave_balance_transactions` + model `LeaveBalanceTransaction` + repository `CreateLeaveBalanceTransaction`/`ListLeaveBalanceTransactions`. Tabel cancellation (§19) dan eligibility (§22) sengaja ditunda — tidak dibutuhkan sampai fitur terkait (Phase 6 lanjutan/§18, §22) mulai dikerjakan |
| Phase 2 - Master Data | ✅ Selesai (2026-08-08) | Leave Types, Accrual Policies, Leave Reasons — CRUD lengkap di `leave/service.go`, `leave/handler.go`, `leave/routes.go`. Leave Eligibility sengaja tidak dibangun — tidak ada business rule konkret yang membutuhkannya saat ini (§22), revisit saat requirement muncul |
| Phase 3 - Leave Calculation Engine | ✅ Selesai (2026-08-08) | `leave/calculation.go` (`CalculateLeaveDuration`): full-day/half-day/hourly calculation, weekend + company-holiday exclusion via new `HolidayProvider` interface (adapter: `setting.Service.ListHolidayDatesInRange`), `LeaveRequestDetail` rows now actually persisted per date. Shift/`DaysOfWeekMask` handling and balance-quota validation intentionally deferred — see §41 |
| Phase 4 - Leave Request | 🔶 Sebagian | Create/Submit/List/Get/Delete/Details sudah ada, `requested_days` kini dihitung server-side (Phase 3). Validasi lain (eligibility, overlap tanggal, balance quota, backdate) di §12 **belum diimplementasikan** |
| Phase 5 - Approval Integration | ✅ Selesai | `ApprovalEngine`, `SetApprovalEngine`, `HandleApprovalStatusChange`, wiring `main.go`, module slug tunggal `"leave"`, test coverage di `approval_integration_test.go` — lihat Section 7 |
| Phase 6 - Leave Balance | 🔶 Sebagian (2026-08-08) | Tabel ledger `leave_balance_transactions` + repository sudah ada (Phase 1), tapi **belum ada logic** accrual/usage/adjustment/reversal/carry-forward/expiry yang benar-benar menulis ke sana atau ke `employee_leave_balances` — lihat Section 9/10 |
| Phase 7 - Calendar & Attendance | ❌ Belum ada | Tidak ada endpoint calendar, tidak ada integrasi Attendance |
| Phase 8 - Notification | ❌ Belum ada | Tidak ditemukan pemanggilan Notification module dari modul Leave |
| Phase 9 - Dashboard & Reports | ❌ Belum ada | Tidak ada endpoint/handler dashboard atau report |
| Phase 10 - Testing | 🔶 Sebagian | Test approval-integration sudah ada; test kalkulasi/balance/cancellation belum ada karena fiturnya sendiri belum ada |

**Frontend**: ❌ belum dimulai — `Leave.vue` hanya placeholder "Coming soon", tidak ada halaman di bawah `views/modules/leave/`.

**Rekomendasi urutan lanjutan** (sesuai prioritas P0 di §42, disesuaikan dengan gap yang sudah dikonfirmasi):
1. Leave Balance auto-deduct on approve + reversal on cancel (saat ini approval "selesai" tanpa efek ke saldo sama sekali — paling kritis).
2. Working-day calculation server-side (saat ini `requested_days` tidak divalidasi, celah untuk data salah/manipulasi).
3. Endpoint `submit`/`cancel` khusus (pisah dari `PUT .../status` generik) supaya side-effect balance/ledger bisa dipasang dengan aman.
4. Frontend dasar: My Leave (request + list + balance) — modul backend sudah cukup matang untuk mulai dari sisi UI.
5. Balance ledger table (§10) menyusul begitu titik (1) mulai berjalan, supaya auto-deduct langsung punya jejak audit sejak awal alih-alih ditambah belakangan.

---

# 44. Design Principles

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
