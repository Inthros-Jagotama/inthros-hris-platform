# Plan Pengembangan Business Travel / Perjalanan Dinas

## 1. Tujuan

Menambahkan fitur **Business Travel / Perjalanan Dinas** pada HRIS sebagai bagian dari module Attendance, dengan integrasi ke:

- Existing Approval Module
- Payroll
- Finance/Accounting
- Document Management
- Employee/Organization

Prinsip utama:

> **Pengajuan perjalanan dinas tidak menentukan metode pembiayaan.**

Metode pembiayaan dan nominal deposit ditentukan **setelah perjalanan dinas disetujui**, oleh pihak yang memiliki kewenangan mengelola pembiayaan.

---

# 2. Business Flow

Flow utama:

```text
REQUEST PERJALANAN DINAS
        |
        v
APPROVAL
(existing Approval Module)
        |
        +---- REJECTED
        |
        +---- APPROVED
                 |
                 v
         PEMBIAYAAN / FUNDING
                 |
        +--------+---------+
        |        |         |
     DEPOSIT  REIMBURSE  OTHER
        |
        v
Isi nominal deposit
Upload bukti transfer jika ada
        |
        v
       PERJALANAN
        |
        v
    PENYELESAIAN
        |
        +-- Actual Expense
        +-- Bukti Perjalanan
        +-- Bukti Pengeluaran
        |
        v
      APPROVAL
(existing Approval Module)
        |
        v
      SETTLEMENT
        |
   +----+----------------+
   |                     |
Deposit > Actual     Deposit < Actual
   |                     |
 Refund              Reimbursement
   |                     |
   +----------+----------+
              |
              v
            CLOSED
```

---

# 3. Prinsip Desain

Pisahkan lifecycle:

```text
Travel
    !=
Funding
    !=
Expense
    !=
Settlement
    !=
Reimbursement
```

Hubungan:

```text
Business Travel
      |
      +-- Participants
      +-- Destinations
      +-- Activities
      +-- Schedule
      |
      +-- Funding
      |     |
      |     +-- Deposit
      |     +-- Reimbursement
      |     +-- Company Paid
      |     +-- Other
      |
      +-- Actual Expenses
      |
      +-- Settlement
            |
            +-- Refund
            +-- Additional Reimbursement
```

---

# 4. Request Perjalanan Dinas

Pada tahap request, user **tidak memilih funding method**.

Request hanya berisi kebutuhan perjalanan:

```text
Judul
Tujuan
Kegiatan
Tanggal Mulai
Tanggal Selesai
Asal
Destinasi
Peserta
Estimasi kebutuhan/biaya
Catatan
```

### Tidak boleh diisi pada tahap request

```text
Funding Method
Deposit Amount
Transfer Reference
Payment Date
```

Karena informasi tersebut baru ditentukan setelah approval.

---

# 5. Data Business Travel

Tabel:

```text
business_travels
```

Field:

```text
id
company_id
request_number
requester_id
title
purpose
description
start_date
end_date
origin
status
approval_status
created_at
updated_at
```

Semua primary key menggunakan UUID.

---

# 6. Status Business Travel

Status Travel dipisahkan dari status Funding dan Settlement.

```text
DRAFT
SUBMITTED
APPROVED
REJECTED
IN_PROGRESS
COMPLETED
CLOSED
CANCELLED
```

Flow:

```text
DRAFT
  ↓
SUBMITTED
  ↓
APPROVED
  ↓
IN_PROGRESS
  ↓
COMPLETED
  ↓
CLOSED
```

Approval menggunakan **existing Approval Module**.

Jangan membuat approval engine baru khusus Business Travel.

---

# 7. Peserta Perjalanan

Satu perjalanan dapat memiliki banyak peserta.

Peserta dapat berupa:

```text
EMPLOYEE
NON_EMPLOYEE
```

Tabel:

```text
business_travel_participants
```

Field:

```text
id
business_travel_id
participant_type
employee_id
name
organization
position
identity_number
email
phone
role
notes
```

Untuk employee:

```text
participant_type = EMPLOYEE
employee_id = UUID
```

Untuk pihak eksternal:

```text
participant_type = NON_EMPLOYEE
employee_id = NULL
name = ...
organization = ...
position = ...
```

Role:

```text
LEADER
MEMBER
DRIVER
NARASUMBER
CLIENT
CONSULTANT
OTHER
```

---

# 8. Tujuan Perjalanan

Satu perjalanan dapat memiliki lebih dari satu tujuan.

Tabel:

```text
business_travel_destinations
```

Field:

```text
id
business_travel_id
sequence
country
province
city
location
arrival_date
departure_date
purpose
notes
```

Contoh:

```text
Jakarta
   ↓
Bandung
   ↓
Cimahi
```

---

# 9. Kegiatan / Agenda

Tabel:

```text
business_travel_activities
```

Field:

```text
id
business_travel_id
activity_date
start_time
end_time
title
description
location
organizer
notes
```

---

# 10. Schedule dan Transportasi

Tabel:

```text
business_travel_schedules
```

Field:

```text
id
business_travel_id
schedule_type
departure_datetime
arrival_datetime
origin
destination
transportation_type
provider
booking_reference
notes
```

Jenis:

```text
DEPARTURE
RETURN
TRANSFER
OTHER
```

Transportasi:

```text
AIRPLANE
TRAIN
BUS
CAR
COMPANY_CAR
RENTAL_CAR
MOTORCYCLE
OTHER
```

---

# 11. Estimasi Biaya

Estimasi biaya dapat dimasukkan saat request untuk memberikan gambaran kebutuhan anggaran.

Namun:

> **Estimasi biaya bukan funding method dan bukan komitmen pembayaran.**

Tabel:

```text
business_travel_expense_plans
```

Field:

```text
id
business_travel_id
participant_id
expense_category_id
description
quantity
unit
estimated_amount
notes
```

Tidak ada:

```text
funding_method
```

pada tabel ini karena funding baru ditentukan setelah approval.

---

# 12. Expense Category

Buat master:

```text
business_travel_expense_categories
```

Contoh:

```text
TRANSPORTATION
ACCOMMODATION
MEAL
TOLL
PARKING
AIRPORT_TAX
TAXI
CAR_RENTAL
COMMUNICATION
REPRESENTATION
DAILY_ALLOWANCE
OTHER
```

Field:

```text
id
company_id
code
name
description
requires_receipt
reimbursable
payroll_treatment
account_code
active
```

---

# 13. Approval Perjalanan

Gunakan existing Approval Module.

Flow:

```text
DRAFT
   ↓
SUBMITTED
   ↓
Existing Approval Workflow
   |
   +---- REJECTED
   |
   +---- APPROVED
```

Contoh approver:

```text
Requester
   ↓
Supervisor
   ↓
Manager
   ↓
HR / Authorized Approver
```

Konfigurasi approval mengikuti module approval yang sudah ada.

---

# 14. Funding Setelah Approval

Setelah Travel:

```text
status = APPROVED
```

baru dapat dilakukan proses:

```text
Funding
```

Funding diisi oleh pihak lain yang memiliki kewenangan.

Contoh:

```text
Requester:
Employee A

Approved By:
Manager

Funding By:
Finance/Admin

Funding:
DEPOSIT

Amount:
Rp 5.000.000
```

---

# 15. Funding Actor

Penting untuk memisahkan siapa yang:

```text
requester
approved_by
funded_by
```

Jangan mengasumsikan requester adalah orang yang memberikan deposit.

Tambahkan informasi:

```text
created_by
funded_by
funded_at
```

Jika approval menggunakan existing Approval Module, approval actor dapat diperoleh dari module tersebut.

---

# 16. Funding Method

Funding method baru dipilih setelah approval.

Master:

```text
business_travel_funding_methods
```

Default:

```text
DEPOSIT
REIMBURSEMENT
COMPANY_PAID
OTHER
```

Contoh:

```text
DEPOSIT
```

berarti perusahaan memberikan uang sebelum perjalanan.

```text
REIMBURSEMENT
```

berarti employee membayar terlebih dahulu dan meminta penggantian setelah perjalanan.

```text
COMPANY_PAID
```

berarti perusahaan membayar langsung.

```text
OTHER
```

untuk mekanisme pembiayaan lain yang dapat dikonfigurasi.

---

# 17. Funding Table

Gunakan tabel terpisah:

```text
business_travel_fundings
```

Field:

```text
id
business_travel_id
funding_method_id
participant_id
amount
funding_date
payment_method
payment_reference
funded_by
status
notes
created_at
updated_at
```

Status:

```text
PENDING
PROCESSING
FUNDED
CANCELLED
REVERSED
```

---

# 18. Bukti Transfer Funding

Buat tabel:

```text
business_travel_funding_documents
```

Field:

```text
id
business_travel_funding_id
document_type
file_name
file_path
mime_type
file_size
uploaded_by
uploaded_at
```

Jenis:

```text
TRANSFER_RECEIPT
PAYMENT_PROOF
OTHER
```

Upload bukti transfer dapat dibuat:

```text
REQUIRED
OPTIONAL
```

berdasarkan konfigurasi.

---

# 19. Deposit / Advance

Deposit adalah **advance**, bukan actual expense.

Contoh:

```text
Funding Method:
DEPOSIT

Amount:
Rp 5.000.000
```

Maka:

```text
Advance = Rp 5.000.000
```

Belum boleh dianggap sebagai:

```text
Travel Expense = Rp 5.000.000
```

Actual expense baru dicatat setelah perjalanan.

---

# 20. Deposit Per Participant

Jika banyak employee mendapatkan deposit:

```text
Employee A
Deposit = 5.000.000

Employee B
Deposit = 3.000.000

Employee C
Deposit = 2.000.000
```

Maka funding harus dapat dikaitkan ke:

```text
participant_id
```

Jika deposit berlaku untuk seluruh perjalanan, `participant_id` dapat dibuat nullable sesuai konfigurasi.

---

# 21. Actual Expense

Setelah perjalanan, peserta mengisi actual expense.

Tabel:

```text
business_travel_expenses
```

Field:

```text
id
business_travel_id
participant_id
expense_category_id
expense_date
description
quantity
unit
amount
funding_method_id
vendor
receipt_number
status
notes
```

Funding method pada actual expense boleh berbeda dengan funding awal.

Contoh:

```text
Travel Funding:
DEPOSIT
```

Tetapi actual:

```text
Hotel       → DEPOSIT
Meal        → DEPOSIT
Parking     → REIMBURSEMENT
```

Ini memungkinkan kondisi nyata.

---

# 22. Bukti Pengeluaran

Tabel:

```text
business_travel_expense_documents
```

Jenis:

```text
RECEIPT
INVOICE
TICKET
HOTEL_BILL
BOARDING_PASS
TOLL_RECEIPT
OTHER
```

Setiap expense dapat memiliki satu atau banyak dokumen.

---

# 23. Bukti Perjalanan Dinas

Selain bukti expense, simpan dokumen perjalanan.

Tabel:

```text
business_travel_documents
```

Jenis:

```text
TRAVEL_ORDER
INVITATION
TICKET
BOARDING_PASS
HOTEL
MEETING_DOCUMENT
ATTENDANCE_PROOF
PHOTO
TRAVEL_REPORT
OTHER
```

---

# 24. Penyelesaian / Settlement

Setelah perjalanan selesai:

```text
IN_PROGRESS
     ↓
COMPLETED
     ↓
SETTLEMENT
```

Peserta/user mengisi:

```text
Actual Expense
Bukti Pengeluaran
Bukti Perjalanan
Catatan
```

Kemudian sistem menghitung:

```text
Total Advance
Total Actual Expense
Refund
Additional Reimbursement
```

---

# 25. Settlement Table

Tabel:

```text
business_travel_settlements
```

Field:

```text
id
business_travel_id
participant_id
total_advance
total_actual_expense
total_company_paid
total_reimbursement
total_refund
balance
status
submitted_at
approved_at
settled_at
notes
```

Status:

```text
PENDING
SUBMITTED
APPROVED
REFUND_REQUIRED
REIMBURSEMENT_REQUIRED
BALANCED
REFUND_PROCESSING
REIMBURSEMENT_PROCESSING
SETTLED
REJECTED
```

---

# 26. Approval Penyelesaian

Settlement menggunakan **existing Approval Module**.

Flow:

```text
Settlement
    ↓
SUBMITTED
    ↓
Existing Approval Module
    |
    +---- REJECTED
    |
    +---- APPROVED
            ↓
         SETTLEMENT
```

Setelah settlement disetujui, sistem menentukan apakah:

```text
Refund
Reimbursement
Balanced
```

---

# 27. Settlement Calculation

Formula dasar:

```text
Actual Expense
-
Advance
=
Settlement Difference
```

Namun expense harus dipisahkan berdasarkan funding.

Gunakan:

```text
DEPOSIT
REIMBURSEMENT
COMPANY_PAID
OTHER
```

---

# 28. Scenario 1 — Reimbursement Murni

Tidak ada deposit.

```text
Advance = 0

Actual Expense = 2.500.000
```

Hasil:

```text
Reimbursement = 2.500.000
```

Flow:

```text
Travel Approved
      ↓
Funding = REIMBURSEMENT
      ↓
Travel
      ↓
Settlement
      ↓
Approval
      ↓
Reimbursement
      ↓
Closed
```

---

# 29. Scenario 2 — Deposit Sama dengan Actual

```text
Deposit        = 5.000.000
Actual Expense = 5.000.000
```

Hasil:

```text
Refund = 0
Reimbursement = 0
Balance = 0
```

Status:

```text
BALANCED
```

Kemudian:

```text
CLOSED
```

---

# 30. Scenario 3 — Deposit Lebih Besar

```text
Deposit        = 5.000.000
Actual Expense = 4.000.000
```

Hasil:

```text
Refund = 1.000.000
```

Flow:

```text
Settlement Approved
       ↓
REFUND_REQUIRED
       ↓
Employee Refund
       ↓
REFUND_CONFIRMED
       ↓
CLOSED
```

---

# 31. Scenario 4 — Deposit Lebih Kecil

```text
Deposit        = 5.000.000
Actual Expense = 6.500.000
```

Hasil:

```text
Additional Reimbursement = 1.500.000
```

Flow:

```text
Settlement Approved
       ↓
REIMBURSEMENT_REQUIRED
       ↓
Reimbursement Process
       ↓
REIMBURSEMENT_PAID
       ↓
CLOSED
```

---

# 32. Scenario 5 — Deposit + Actual Reimbursement

Kasus:

```text
Initial Deposit:
Rp 5.000.000

Actual:
Hotel       Rp 2.500.000
Transport   Rp 2.000.000
Meal        Rp 1.500.000
Parking     Rp   500.000
-------------------------
Total       Rp 6.500.000
```

Hasil:

```text
Deposit             Rp 5.000.000
Actual              Rp 6.500.000
Additional Claim    Rp 1.500.000
```

Sistem harus menangani:

```text
Advance Settlement
+
Additional Reimbursement
```

dalam satu Business Travel.

Tidak perlu membuat perjalanan baru.

---

# 33. Scenario 6 — Mixed Funding

Contoh:

```text
Hotel       → COMPANY_PAID
Transport   → DEPOSIT
Meal        → DEPOSIT
Parking     → REIMBURSEMENT
```

Sistem harus dapat memisahkan:

```text
Company Paid Expense
Deposit Expense
Reimbursement Expense
```

Sehingga settlement tetap akurat.

---

# 34. Company Paid

Company Paid tidak menjadi hutang kepada employee.

Contoh:

```text
Hotel
Amount = 2.500.000
Funding = COMPANY_PAID
```

Maka:

```text
Company Travel Cost = 2.500.000
Employee Payable = 0
```

---

# 35. Refund

Jika:

```text
Actual < Deposit
```

maka:

```text
Refund Required
```

Sistem mencatat:

```text
refund_amount
refund_date
refund_reference
refunded_by
refund_document
```

Tabel dapat berupa:

```text
business_travel_refunds
```

---

# 36. Reimbursement

Jika:

```text
Actual > Deposit
```

atau funding awal:

```text
REIMBURSEMENT
```

maka sistem membuat claim.

Tabel:

```text
business_travel_reimbursements
```

Field:

```text
id
business_travel_id
participant_id
settlement_id
amount
status
requested_at
approved_at
paid_at
payment_reference
paid_by
notes
```

Status:

```text
REQUESTED
APPROVED
PROCESSING
PAID
REJECTED
CANCELLED
```

---

# 37. Attendance Integration

Business Travel menjadi source Attendance.

Gunakan reference:

```text
attendance.source_type = BUSINESS_TRAVEL
attendance.source_id = business_travel_id
```

Contoh:

```text
20 Aug → BUSINESS_TRAVEL
21 Aug → BUSINESS_TRAVEL
22 Aug → BUSINESS_TRAVEL
```

Konfigurasi:

```text
business_travel_attendance_rules
```

Untuk:

```text
FULL_DAY
HALF_DAY
TRAVEL_DAY
NON_WORKING_DAY
```

---

# 38. Payroll Integration

Business Travel dapat terintegrasi dengan Payroll untuk komponen tertentu.

Contoh:

```text
Daily Allowance
Travel Allowance
Travel Incentive
```

Tetapi:

```text
Reimbursement
Refund
Advance
```

tidak otomatis dianggap salary.

Gunakan:

```text
payroll_treatment
```

pada expense category.

---

# 39. Accounting / Finance Integration

Transaksi yang perlu didukung:

```text
Advance
Travel Expense
Reimbursement Payable
Refund
Company Paid
```

Contoh:

### Advance

```text
Employee Advance       Dr
    Cash/Bank              Cr
```

### Actual Expense menggunakan Advance

```text
Travel Expense         Dr
    Employee Advance       Cr
```

### Actual lebih besar

```text
Travel Expense         Dr
    Employee Payable       Cr
```

### Actual lebih kecil

```text
Travel Expense         Dr
Cash/Bank               Dr
    Employee Advance       Cr
```

Accounting mapping harus configurable.

---

# 40. Separation of Duties

Sistem harus mendukung pemisahan:

```text
Requester
    ↓
Approver
    ↓
Funder / Finance
    ↓
Settlement Submitter
    ↓
Settlement Approver
    ↓
Reimbursement Payer
```

Tidak harus selalu orang berbeda, tetapi sistem tidak boleh mengasumsikan bahwa semua proses dilakukan requester.

---

# 41. Permissions

```text
business_travel.view
business_travel.create
business_travel.update
business_travel.submit
business_travel.cancel

business_travel.approve
business_travel.reject

business_travel.funding.view
business_travel.funding.create
business_travel.funding.update
business_travel.funding.confirm
business_travel.funding.upload

business_travel.expense.view
business_travel.expense.create
business_travel.expense.update
business_travel.expense.delete

business_travel.settlement.view
business_travel.settlement.create
business_travel.settlement.submit
business_travel.settlement.approve

business_travel.refund.process
business_travel.reimbursement.process

business_travel.report
business_travel.export
```

---

# 42. Database Structure

Tabel utama:

```text
business_travels

business_travel_destinations
business_travel_activities
business_travel_participants
business_travel_schedules

business_travel_expense_categories
business_travel_expense_plans

business_travel_funding_methods
business_travel_fundings
business_travel_funding_documents

business_travel_expenses
business_travel_expense_documents

business_travel_documents

business_travel_settlements
business_travel_settlement_items

business_travel_refunds
business_travel_reimbursements

business_travel_audit_logs
```

Approval tidak perlu membuat tabel workflow baru jika existing Approval Module sudah mendukung generic approval reference.

---

# 43. Relasi Database

```text
business_travels
    |
    +-- participants
    |
    +-- destinations
    |
    +-- activities
    |
    +-- schedules
    |
    +-- expense_plans
    |
    +-- fundings
    |      |
    |      +-- funding_documents
    |
    +-- expenses
    |      |
    |      +-- expense_documents
    |
    +-- documents
    |
    +-- settlements
           |
           +-- settlement_items
           |
           +-- refunds
           |
           +-- reimbursements
```

---

# 44. API

## Travel

```text
GET    /api/business-travels
POST   /api/business-travels
GET    /api/business-travels/{id}
PUT    /api/business-travels/{id}
DELETE /api/business-travels/{id}

POST   /api/business-travels/{id}/submit
POST   /api/business-travels/{id}/cancel
```

## Funding

```text
GET    /api/business-travels/{id}/fundings
POST   /api/business-travels/{id}/fundings
PUT    /api/business-travels/{id}/fundings/{fundingId}
POST   /api/business-travels/{id}/fundings/{fundingId}/confirm
POST   /api/business-travels/{id}/fundings/{fundingId}/documents
```

## Expenses

```text
GET    /api/business-travels/{id}/expenses
POST   /api/business-travels/{id}/expenses
PUT    /api/business-travels/{id}/expenses/{expenseId}
DELETE /api/business-travels/{id}/expenses/{expenseId}
POST   /api/business-travels/{id}/expenses/{expenseId}/documents
```

## Settlement

```text
GET    /api/business-travels/{id}/settlement
POST   /api/business-travels/{id}/settlement
POST   /api/business-travels/{id}/settlement/submit
```

## Refund

```text
POST   /api/business-travels/{id}/refund
POST   /api/business-travels/{id}/refund/confirm
```

## Reimbursement

```text
GET    /api/business-travels/{id}/reimbursements
POST   /api/business-travels/{id}/reimbursements/{reimbursementId}/process
POST   /api/business-travels/{id}/reimbursements/{reimbursementId}/paid
```

---

# 45. UI Flow

## A. Request

Form:

```text
General Information
├── Title
├── Purpose
├── Start Date
├── End Date
├── Origin
└── Destination

Participants
├── Employee
└── Non Employee

Activities

Schedule

Estimated Expense

Documents
```

Tidak ada:

```text
Funding Method
Deposit
Transfer
```

---

## B. Setelah Approved

Pada detail Travel muncul:

```text
Funding
```

Contoh:

```text
Funding Method:
[ Deposit ]

Amount:
[ Rp 5.000.000 ]

Funding Date:
[ 20-08-2026 ]

Payment Reference:
[ TRX-12345 ]

Transfer Proof:
[ Upload ]
```

---

## C. Settlement

```text
Settlement

Total Advance             Rp 5.000.000
Total Actual Expense      Rp 6.500.000
Company Paid              Rp 1.000.000
Eligible Expense          Rp 6.500.000

Additional Reimbursement  Rp 1.500.000
```

Atau:

```text
Total Advance             Rp 5.000.000
Total Actual Expense      Rp 4.000.000

Refund Required           Rp 1.000.000
```

---

# 46. Dashboard

```text
Business Travel

Pending Approval             8
Approved / Upcoming         12
Active Travel                6
Settlement Pending           9

Funding Pending              4
Outstanding Advance         15

Pending Reimbursement
Rp 35.000.000

Pending Refund
Rp 12.000.000
```

---

# 47. Reporting

## Travel Report

```text
Travel Number
Requester
Destination
Start Date
End Date
Status
```

## Funding Report

```text
Travel
Participant
Funding Method
Amount
Funding Date
Funded By
Status
```

## Advance Report

```text
Travel
Employee
Advance
Actual Expense
Remaining
Settlement Status
```

## Reimbursement Report

```text
Travel
Employee
Claim
Approved
Paid
Payment Date
```

## Refund Report

```text
Travel
Employee
Advance
Actual Expense
Refund
Refund Date
Status
```

## Travel Cost Report

```text
Travel
Company Paid
Advance
Reimbursement
Total Actual Cost
```

---

# 48. Audit Trail

Tabel:

```text
business_travel_audit_logs
```

Event:

```text
TRAVEL_CREATED
TRAVEL_UPDATED
TRAVEL_SUBMITTED
TRAVEL_APPROVED
TRAVEL_REJECTED

FUNDING_CREATED
FUNDING_UPDATED
FUNDING_CONFIRMED
FUNDING_DOCUMENT_UPLOADED

EXPENSE_CREATED
EXPENSE_UPDATED
EXPENSE_DELETED

SETTLEMENT_CREATED
SETTLEMENT_SUBMITTED
SETTLEMENT_APPROVED

REFUND_REQUIRED
REFUND_CONFIRMED

REIMBURSEMENT_REQUIRED
REIMBURSEMENT_APPROVED
REIMBURSEMENT_PAID

TRAVEL_CLOSED
```

Simpan:

```text
entity_type
entity_id
action
old_value
new_value
user_id
timestamp
ip_address
```

---

# 49. Development Phase

## Phase 1 — Foundation

- [x] Create Business Travel module — **keputusan: tidak jadi module Go terpisah, digabung ke module `attendance` existing** (lihat §54.1).
- [x] Migration — 12 migrasi (`124`–`135`, postgres+mysql, up+down = 48 file) di `backend/internal/pkg/migrator/migrations/tenant/{postgres,mysql}/`. Belum dijalankan (di luar scope task ini).
- [x] Model/entity — `backend/internal/modules/attendance/model_businesstravel.go` (package `attendance`), build-verified tanpa collision nama tipe.
- [x] Repository/service — `repository_businesstravel.go` (CRUD travel/participant/destination) + `service_businesstravel.go` (Create/Get/List/Update/Submit/Cancel + approval wiring), build & `go vet` bersih.
- [x] UUID — semua PK `CHAR(36)` + `uuid.UUID` dengan `BeforeCreate` hook, mengikuti pola reimbursement.
- [x] Company scope — sengaja **tidak** pakai kolom `company_id` (arsitektur DB-per-tenant, lihat §54.2).
- [x] Permission — semua slug `business_travel.*` (§41) ditambahkan ke `attModule.Info().Permissions()` di `module.go`. Catatan: codebase ini tidak punya middleware permission per-route di module manapun (dicek di seluruh `internal/modules`) — RBAC ditegakkan di level module/menu yang lebih kasar, jadi deklarasi ini konsisten dengan pola existing, bukan gap baru.
- [x] Travel CRUD — `handler_businesstravel.go` + route `/attendance/business-travels` (POST/GET/GET:id/PUT).
- [x] Travel status — `TravelStatus` enum (DRAFT→SUBMITTED→APPROVED/REJECTED→...→CLOSED/CANCELLED) di model, transisi DRAFT→SUBMITTED→CANCELLED sudah diimplementasi di service.

## Phase 2 — Travel Information

- [x] Participants — dibuat inline saat `CreateBusinessTravel` (`repo.CreateParticipant`), belum ada endpoint tambah/edit participant terpisah setelah travel dibuat.
- [x] Employee participant — `participant_type = EMPLOYEE` + `employee_id`.
- [x] Non-employee participant — `participant_type = NON_EMPLOYEE` + field manual (name/organization/dst).
- [x] Destinations — dibuat inline saat `CreateBusinessTravel` (`repo.CreateDestination`), belum ada endpoint tambah/edit terpisah.
- [x] Activities — `POST/GET /attendance/business-travels/:id/activities`.
- [x] Schedule — `POST/GET /attendance/business-travels/:id/schedules`.
- [x] Transportation — bagian dari Schedule (`transportation_type` field), sama endpoint di atas.

## Phase 3 — Approval Integration

- [x] Integrate existing Approval Module — reuse `Service.approvalEngine` (field sama dengan overtime/correction), module slug baru `"business_travel"` agar flow terpisah dari `"attendance"`.
- [x] Submit request — `Service.SubmitBusinessTravel` (DRAFT→SUBMITTED, `CreateApprovalInstance`/`GetActiveFlowIDForModule`, `approval.RoutingError` fail loudly seperti overtime).
- [x] Approval callback/status synchronization — `Service.HandleBusinessTravelApprovalStatusChange` (APPROVED/REJECTED), diregister via `approvalSvc.RegisterStatusHandler("business_travel", ...)` di `backend/cmd/server/main.go` (titik sama dengan wiring `attendanceSvc.SetApprovalEngine`). Approval end-to-end sudah nyambung secara kode; belum diverifikasi dengan test/manual run.
- [x] Rejection — status REJECTED di-set oleh callback di atas.
- [ ] Approval history — belum ada endpoint untuk menampilkan riwayat approval (bisa query langsung ke module approval by instance ID untuk sementara).

## Phase 4 — Funding

- [x] Funding method master — CRUD `POST/GET /attendance/business-travel-funding-methods` (code/name bebas diisi admin; DEPOSIT/REIMBURSEMENT/COMPANY_PAID/OTHER dari §16 belum di-seed otomatis — codebase ini tidak punya pola auto-seed master data di Module.Seed() manapun, jadi harus dibuat manual via endpoint ini, konsisten dengan reimbursement_types dkk).
- [x] Funding transaction — `POST/GET /attendance/business-travels/:id/fundings`, `PUT .../fundings/:fundingId`. Digating `ErrBusinessTravelNotApproved` jika travel belum APPROVED/IN_PROGRESS/COMPLETED (Rule 1, §52).
- [x] Deposit / [x] Reimbursement funding / [x] Company Paid / [x] Other — semua ditangani generik lewat `funding_method_id`, bukan tipe terpisah di kode (sesuai §16: master funding method configurable).
- [x] Funding actor — `FundedBy` diisi dari `authctx.GetUserID(ctx)` saat create, terpisah dari `RequesterID`/`CreatedBy` travel (Rule 2, §52).
- [x] Payment reference — field `payment_reference` di create/update.
- [x] Transfer proof — `POST .../fundings/:fundingId/documents`, menyimpan `file_path`/URL yang didapat client dari endpoint upload generik (§54.4) — module ini tidak menangani upload file mentah.
- [x] Funding confirmation — `POST .../fundings/:fundingId/confirm` (PENDING/PROCESSING → FUNDED).

## Phase 5 — Actual Expense

- [x] Expense category — CRUD `POST/GET /attendance/business-travel-expense-categories` (sama pola dengan funding methods, master diisi manual — codebase tidak punya auto-seed).
- [x] Actual expense — `POST/GET /attendance/business-travels/:id/expenses`, `PUT/DELETE .../expenses/:expenseId`. Digating `ErrBusinessTravelNotApproved` sama seperti funding (hanya bisa dicatat setelah travel APPROVED/IN_PROGRESS/COMPLETED).
- [x] Funding method per expense — `funding_method_id` opsional per expense, independen dari funding awal travel (mixed funding §33).
- [x] Expense receipt — `POST .../expenses/:expenseId/documents`, pola sama dengan funding transfer proof (URL dari endpoint upload generik, §54.4).
- [x] Expense validation — kategori & travel divalidasi exist saat create; update/delete ditolak (`ErrExpenseInvalidState`) jika expense sudah berstatus APPROVED.
- [x] Expense per participant — `participant_id` opsional per expense.

## Phase 6 — Settlement

- [x] Settlement form — `POST /attendance/business-travels/:id/settlements` (digating `ErrBusinessTravelNotCompleted`, travel harus COMPLETED — §24).
- [x] Calculate total advance — sum funding `FUNDED` dengan method code `DEPOSIT` (funding lain seperti REIMBURSEMENT/COMPANY_PAID/OTHER tidak dihitung sebagai advance).
- [x] Calculate actual expense — sum semua expense (difilter participant jika diisi).
- [x] Calculate company paid — sum expense dengan `funding_method_id` ber-code `COMPANY_PAID`, dikeluarkan dari rekonsiliasi advance (§34, Rule: Company Paid bukan hutang ke employee).
- [x] Calculate refund — `diff = (totalActual - totalCompanyPaid) - totalAdvance`; diff < 0 → `TotalRefund = -diff` (§30 Scenario 3).
- [x] Calculate additional reimbursement — diff > 0 → `TotalReimbursement = diff` (§31 Scenario 4). diff == 0 → BALANCED (§29 Scenario 2). Formula sudah menangani Scenario 1 (reimbursement murni, advance=0) dan Scenario 6 (mixed funding) karena company-paid & advance dipisah per item, bukan diasumsikan satu metode per travel.
- [x] Settlement per participant — `participant_id` opsional di `CreateSettlementRequest`; kosong = gabungan seluruh peserta.
- [x] Settlement approval — `POST .../settlements/:settlementId/submit`, module slug approval terpisah `"business_travel_settlement"` (bukan `"business_travel"` milik travel), diregister di `main.go`. Hasil akhir (BALANCED/REFUND_REQUIRED/REIMBURSEMENT_REQUIRED) baru ditentukan & disimpan permanen saat approval APPROVED (`HandleSettlementApprovalStatusChange`), bukan saat create — nilai di CreateSettlement adalah proyeksi awal.

## Phase 7 — Refund

- [x] Refund calculation — sudah dihitung di Phase 6 (`Settlement.TotalRefund`/`Balance`).
- [x] Refund transaction — record `business_travel_refunds` **dibuat otomatis** oleh `HandleSettlementApprovalStatusChange` saat settlement APPROVED dan `TotalRefund > 0`, status awal `PENDING`.
- [ ] Refund proof — `refund_document` diisi sebagai string URL langsung di `ConfirmRefundRequest` (bukan endpoint upload multi-dokumen terpisah seperti funding/expense — refund cuma satu bukti transfer balik, dianggap cukup satu field).
- [x] Refund confirmation — `POST /attendance/business-travels/:id/refunds/:refundId/confirm` (PENDING → CONFIRMED, `refunded_by`/`refund_date`/`refund_reference` terisi otomatis).
- [x] Close settlement — `Service.maybeSettleAndCloseTravel`: begitu refund CONFIRMED (atau reimbursement PAID), settlement terkait → `SETTLED`, lalu travel → `CLOSED` otomatis jika **semua** settlement milik travel tsb sudah BALANCED/SETTLED (mendukung settlement per participant, §33).

## Phase 8 — Reimbursement

- [x] Generate reimbursement claim — record `business_travel_reimbursements` **dibuat otomatis** oleh `HandleSettlementApprovalStatusChange` saat settlement APPROVED dan `TotalReimbursement > 0`, status awal `REQUESTED`.
- [x] Reimbursement approval/status — `POST .../reimbursements/:reimbursementId/{approve,process,pay}` (REQUESTED→APPROVED→PROCESSING→PAID). **Gerbang §54.7 diterapkan di `ProcessTravelReimbursement`**: cek `moduleChecker.IsModuleActive(companyID, "reimbursement")` via `Service.SetModuleChecker` (adapter `approvalModuleCheckerAdapter` yang sama dipakai `approval.Service`, diwire di `main.go`). Catatan jujur: module Reimbursement standalone belum punya API publik untuk menerima claim dari luar, jadi baik subscribed maupun tidak, claim tetap diproses internal — pengecekan subscription sudah ada dan di-log sebagai hint, tapi push cross-module belum diimplementasikan (todo lanjutan jika module Reimbursement menyediakan endpoint ingest).
- [x] Payment / [x] Payment reference — `POST .../pay` (PROCESSING→PAID, `paid_by`/`paid_at`/`payment_reference`).
- [ ] Payment proof — belum ada field/endpoint upload bukti pembayaran reimbursement (model `TravelReimbursement` tidak punya kolom dokumen — beda dari Funding/Expense yang punya tabel `*_documents` terpisah).
- [x] Close travel — sama seperti Phase 7, `maybeSettleAndCloseTravel` dipanggil setelah `PayTravelReimbursement` sukses.

## Phase 9 — Attendance

- [x] Business Travel attendance source — **keputusan desain: dedicated column, bukan generic `source_type`/`source_id`**. Plan doc §37 awalnya mengusulkan `attendance.source_type = BUSINESS_TRAVEL` generik, tapi codebase sudah punya pola established untuk kasus yang sama persis (Leave → Attendance, migration 004 `leave_request_id`/`leave_fraction` + `Service.ApplyApprovedLeave`, dipush via `leaveSvc.SetAttendanceSessionUpdater(attendanceSvc)`). Business Travel mengikuti pola itu: kolom `business_travel_id` baru di `attendance_sessions` (migration `136`, postgres+mysql), status session baru `BUSINESS_TRAVEL`, method `Service.ApplyApprovedBusinessTravel` (mirror `ApplyApprovedLeave` persis termasuk proteksi CLOSED-day). Dipanggil otomatis dari `HandleBusinessTravelApprovalStatusChange` saat travel APPROVED (`pushBusinessTravelAttendance`, iterasi semua participant `EMPLOYEE` × tanggal `StartDate..EndDate`) — tidak perlu wiring updater terpisah seperti leave karena sudah satu Service yang sama.
- [x] Travel day — session di-mark `BUSINESS_TRAVEL` untuk setiap hari dalam rentang `start_date..end_date`.
- [ ] Full day / [ ] Half day / [ ] Weekend — belum ada pembedaan granular; saat ini semua hari travel diperlakukan sama (full day). Tabel `business_travel_attendance_rules` (migration `135`) sudah ada di skema tapi belum dipakai untuk membedakan FULL_DAY/HALF_DAY/TRAVEL_DAY/NON_WORKING_DAY — masih best-effort binary (travel day vs bukan).
- [x] Attendance integration — sesi tidak ditimpa jika sudah `CLOSED` (hari kerja nyata yang sudah selesai), sama seperti proteksi leave.

## Phase 10 — Payroll / Accounting

- [x] Payroll treatment — **keputusan desain penting, baca §54.8**: `business_travel_expense_categories.payroll_treatment` (§12) di-repurpose untuk menyimpan **UUID `SalaryComponent`** milik module Payroll secara langsung, bukan label teks bebas seperti tersirat di §12. Alasan: payroll module tidak punya lookup component-by-code, hanya by-ID (`FindSalaryComponentByID`), jadi tidak ada cara resolve dari kode/label ke component tanpa menambah endpoint baru di module Payroll (di luar scope perubahan ini, risiko lebih tinggi karena Payroll adalah module produksi dengan kalkulasi pajak yang sensitif).
- [x] Travel allowance / [x] Daily allowance — ditangani generik: expense apapun dengan kategori yang `payroll_treatment`-nya diisi UUID component valid akan didorong ke payroll sebagai `SalaryEmployeeAdjustment` satu-kali, terlepas dari nama kategorinya (admin yang menentukan kategori mana yang "Daily Allowance" dsb dengan mengisi `payroll_treatment`).
- [ ] Accounting mapping — **belum diimplementasikan.** §39 minta accounting mapping configurable (Advance/Travel Expense/Reimbursement Payable/Refund/Company Paid ke akun COA) — tidak ada module Accounting/GL di codebase ini untuk diintegrasikan, jadi bagian ini di-skip sepenuhnya, bukan sekadar tertunda.
- [ ] Advance accounting / [ ] Reimbursement payable / [ ] Refund / [ ] Company Paid (accounting) — sama seperti di atas, tidak ada sink accounting untuk didorong. Data mentahnya sendiri (advance dari funding, actual expense, refund, reimbursement) sudah lengkap tercatat di tabel-tabel Phase 4-8 dan bisa jadi sumber laporan/ekspor manual ke sistem akuntansi eksternal.

**Implementasi (§54.8):** `pushBusinessTravelPayrollAdjustments` dipanggil dari `HandleSettlementApprovalStatusChange` saat settlement APPROVED — iterasi `SettlementItem` ber-`ItemType=ACTUAL`, ambil `Expense` → `ExpenseCategory.PayrollTreatment` (parse sebagai UUID component) → `Expense.ParticipantID` → `Participant.EmployeeID`, lalu push `SalaryEmployeeAdjustment(source_type="BUSINESS_TRAVEL", status="APPROVED")` via `Service.payrollAdjuster` (interface `PayrollAdjustmentProvider`, diwire di `main.go` dengan adapter `attendancePayrollAdjusterAdapter` yang membungkus `*payroll.Repository` — **tidak ada kode baru ditambahkan ke module `payroll`**, murni reuse `CreateSalaryEmployeeAdjustment` yang sudah ada). Item non-payroll (kategori tanpa `payroll_treatment`, atau expense tanpa participant employee) dilewati diam-diam. Best-effort: gagal push tidak menggagalkan approval settlement.

## Phase 11 — Documents

- [ ] Travel order.
- [ ] Ticket.
- [ ] Boarding pass.
- [ ] Hotel.
- [ ] Receipt.
- [ ] Transfer proof.
- [ ] Travel report.
- [ ] Activity evidence.

## Phase 12 — Reporting

- [ ] Travel report.
- [ ] Funding report.
- [ ] Advance report.
- [ ] Expense report.
- [ ] Settlement report.
- [ ] Refund report.
- [ ] Reimbursement report.
- [ ] Company travel cost.

---

# 50. Testing

## Request

- [ ] Create travel.
- [ ] Submit travel.
- [ ] Validate mandatory fields.
- [ ] Ensure funding cannot be selected during request.

## Approval

- [ ] Submit to existing approval module.
- [ ] Approved.
- [ ] Rejected.
- [ ] Verify status synchronization.

## Funding

- [ ] Create deposit after approval.
- [ ] Create reimbursement funding.
- [ ] Create company paid.
- [ ] Add funding amount.
- [ ] Upload transfer proof.
- [ ] Confirm funding.
- [ ] Verify requester cannot perform unauthorized funding action.

## Settlement

### Deposit Balanced

```text
Deposit = 5M
Actual = 5M

Expected:
Balance = 0
```

### Deposit + Refund

```text
Deposit = 5M
Actual = 4M

Expected:
Refund = 1M
```

### Deposit + Reimbursement

```text
Deposit = 5M
Actual = 6.5M

Expected:
Additional Reimbursement = 1.5M
```

### Reimbursement Only

```text
Deposit = 0
Actual = 2.5M

Expected:
Reimbursement = 2.5M
```

### Mixed Funding

```text
Hotel       = COMPANY_PAID
Transport   = DEPOSIT
Meal        = DEPOSIT
Parking     = REIMBURSEMENT
```

Expected:

```text
Company Paid
+
Advance Settlement
+
Reimbursement
```

dihitung secara terpisah.

---

# 51. Target Architecture

```text
                         BUSINESS TRAVEL
                                |
          +---------------------+---------------------+
          |                     |                     |
    Travel Request         Participants          Activities
          |                     |                     |
          +---------------------+---------------------+
                                |
                         EXISTING APPROVAL
                                |
                         +------+------+
                         |             |
                     REJECTED       APPROVED
                                       |
                                       v
                                  FUNDING
                                       |
                         +-------------+-------------+
                         |             |             |
                      DEPOSIT     REIMBURSEMENT  COMPANY PAID
                         |
                    Deposit Amount
                         |
                    Transfer Proof
                         |
                         v
                      TRAVEL
                         |
                         v
                    COMPLETED
                         |
                         v
                    SETTLEMENT
                         |
                 +-------+--------+
                 |                |
             Actual <          Actual >
             Advance            Advance
                 |                |
               REFUND       REIMBURSEMENT
                 |                |
                 +-------+--------+
                         |
                  EXISTING APPROVAL
                         |
                         v
                       CLOSED
```

---

# 52. Prinsip Bisnis Final

## Rule 1 — Request Tidak Menentukan Pembiayaan

```text
Request
    ↓
Approval
    ↓
Funding
```

Bukan:

```text
Request + Funding
```

---

## Rule 2 — Funding Diisi Pihak Berwenang

Requester tidak otomatis menjadi funder.

```text
Requester
    ≠
Funder
```

kecuali memang memiliki permission.

---

## Rule 3 — Deposit Adalah Advance

```text
Deposit != Expense
```

Deposit hanya menjadi uang muka yang nantinya direkonsiliasi.

---

## Rule 4 — Actual Expense Dicatat Terpisah

```text
Advance
   +
Actual Expense
   ↓
Settlement
```

---

## Rule 5 — Settlement Menentukan Hasil Akhir

```text
Actual < Advance
    → Refund

Actual = Advance
    → Balanced

Actual > Advance
    → Additional Reimbursement
```

---

## Rule 6 — Approval Settlement Menggunakan Module Existing

Tidak membuat approval engine baru.

```text
Travel Approval
      ↓
Existing Approval Module

Settlement Approval
      ↓
Existing Approval Module
```

---

## Rule 7 — Reimbursement Setelah Settlement

Jika hasil settlement membutuhkan reimbursement:

```text
Settlement Approved
       ↓
Reimbursement Required
       ↓
Reimbursement Paid
       ↓
Closed
```

---

# 53. Kesimpulan

Arsitektur yang direkomendasikan:

```text
1. Request perjalanan dinas
       ↓
2. Approval
       ↓
3. Funding oleh pihak berwenang
       ↓
4. Perjalanan
       ↓
5. Penyelesaian + bukti
       ↓
6. Approval penyelesaian
       ↓
7. Settlement
       ↓
8. Refund / Reimbursement / Balanced
       ↓
9. Closed
```

Desain ini memastikan sistem dapat menangani:

- Reimbursement tanpa deposit.
- Deposit.
- Deposit dengan sisa uang dan refund.
- Deposit yang ternyata kurang dan membutuhkan reimbursement.
- Company Paid.
- Kombinasi beberapa metode pembiayaan.
- Banyak employee.
- Employee dan non-employee.
- Funding dilakukan oleh pihak lain.
- Bukti transfer.
- Bukti perjalanan.
- Bukti pengeluaran.
- Settlement per participant.
- Approval yang menggunakan module existing.

**Core design:** `Business Travel → Approval → Funding → Travel → Settlement → Approval → Refund/Reimbursement → Closed`.

---

# 54. Rencana Implementasi BE & FE (Sesuai Struktur Codebase Existing)

Berdasarkan struktur module `reimbursement`, `attendance`, dan `approval` yang sudah ada di codebase, berikut rencana implementasi konkret.

## 54.1 Struktur Module Backend

**Keputusan:** Business Travel **tidak** dibuat sebagai module Go terpisah — masuk ke dalam module **`attendance`** yang sudah ada (`backend/internal/modules/attendance/`), karena secara bisnis business travel adalah salah satu source attendance (§37) dan bagian dari module Attendance sesuai judul plan doc ini. Tidak ada module baru yang perlu diregistrasi ke module registry/`main.go` — cukup extend module `attendance` yang sudah terdaftar.

File ditambahkan sebagai file baru dengan `package attendance`, mengikuti pola split-file yang sudah dipakai attendance module (`model.go`, `handler.go`, `service.go`, `repository.go`, `session.go`, `geofence.go`, dst — satu package, banyak file per domain):

```text
backend/internal/modules/attendance/
├── model_businesstravel.go        # GORM entities (semua tabel di section 42) — SUDAH DIBUAT
├── dto_businesstravel.go           # request/response DTO
├── repository_businesstravel.go    # data access via *gorm.DB per-tenant
├── service_businesstravel.go        # business logic, status/state machine, approval-engine wiring
├── handler_businesstravel.go        # Gin HTTP handlers
├── routes_businesstravel.go          # route registration (ditambahkan ke routes.go existing)
└── *_businesstravel_test.go          # unit + approval_integration_test.go
```

Konsekuensi penamaan: karena satu package dengan attendance yang sudah punya banyak tipe & fungsi, semua identifier baru diberi prefix implisit lewat nama tipe (`BusinessTravel`, `TravelStatus`, dst) agar tidak bentrok — sudah diverifikasi tidak ada collision nama tipe dengan file attendance lain saat `model_businesstravel.go` dibuat (`go build ./internal/modules/attendance/...` sukses).

Permission tetap pakai namespace `business_travel.*` (section 41) meski secara kode menyatu dengan module `attendance` — pemisahan permission tidak terikat pemisahan package Go.

## 54.2 Migrasi Database

Ikuti konvensi migrator existing: file SQL bernomor urut, dua dialect (`mysql` & `postgres`), embedded via `//go:embed`. Migrasi terbaru saat ini berhenti di `123_pph21_settings_calculation_method.sql`, sehingga migrasi business travel dimulai dari **124**.

```text
backend/internal/pkg/migrator/migrations/tenant/postgres/124_business_travels.sql
backend/internal/pkg/migrator/migrations/tenant/postgres/124_business_travels.down.sql
backend/internal/pkg/migrator/migrations/tenant/mysql/124_business_travels.sql
backend/internal/pkg/migrator/migrations/tenant/mysql/124_business_travels.down.sql
125_business_travel_participants_destinations.sql   (+ mysql)
126_business_travel_activities_schedules.sql          (+ mysql)
127_business_travel_expense_master.sql                 (+ mysql)   -- expense_categories, expense_plans
128_business_travel_funding.sql                          (+ mysql)  -- funding_methods, fundings, funding_documents
129_business_travel_expenses.sql                          (+ mysql) -- expenses, expense_documents
130_business_travel_documents.sql                          (+ mysql)
131_business_travel_settlement.sql                          (+ mysql) -- settlements, settlement_items
132_business_travel_refund_reimbursement.sql                 (+ mysql)
133_business_travel_approval_instance.sql                      (+ mysql) -- approval_instance_id column, mengikuti pola 061_reimbursement_approval_instance.sql
134_business_travel_audit_logs.sql                              (+ mysql)
135_business_travel_attendance_rules.sql                          (+ mysql)
```

Semua primary key `UUID`. Tambahkan `.down.sql` untuk setiap migrasi agar rollback konsisten dengan pola existing. Nomor final harus dicek ulang saat implementasi karena migration lain mungkin sudah ditambahkan di antara waktu penulisan plan ini dan eksekusi — cek `ls backend/internal/pkg/migrator/migrations/tenant/postgres | sort | tail`.

Karena arsitektur multi-tenant di sini adalah **DB-per-tenant** (bukan kolom `company_id` per baris), `module.go` cukup menyediakan `TenantDBFunc` yang resolve `*gorm.DB` dari `company_id` di context (pola sama seperti `reimbursement/module.go:23-30`). Kolom `company_id` yang disebut di section 5 (`business_travels.company_id`) **tidak diperlukan** mengikuti konvensi existing — hapus dari model kecuali ada kebutuhan cross-tenant reporting khusus.

## 54.3 Integrasi Approval Module

Ikuti pola reimbursement persis (section 42/section 3 dokumen ini sudah sejalan):

1. Tambahkan kolom `approval_instance_id` pada `business_travels` dan `business_travel_settlements` (dua alur approval terpisah: Travel Approval & Settlement Approval).
2. Definisikan interface lokal di `service.go`:
   ```go
   type ApprovalEngine interface {
       CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
       GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
   }
   ```
3. Saat submit travel / submit settlement, panggil `CreateApprovalInstance(ctx, "business_travel", documentID, flowID)` (gunakan module slug berbeda untuk settlement, misal `"business_travel_settlement"`, atau flow berbeda dalam module yang sama — tentukan saat desain flow approval).
4. Daftarkan callback status di `main.go`:
   ```go
   approvalSvc.RegisterStatusHandler("business_travel", businessTravelSvc.HandleApprovalStatusChange)
   approvalSvc.RegisterStatusHandler("business_travel_settlement", businessTravelSvc.HandleSettlementApprovalStatusChange)
   ```
5. Handle `approval.RoutingError` secara graceful (mis. tidak ada approver terkonfigurasi) sesuai pola reimbursement.
6. **Jangan** membuat approval engine/tabel workflow baru — sesuai Rule 6 (section 52).

## 54.4 Upload / Dokumen

Gunakan endpoint upload generik existing: `POST /api/v1/tenant/uploads` (`backend/internal/pkg/upload/handler.go`), bukan endpoint upload khusus per module. Field-field dokumen (`business_travel_funding_documents.file_path`, `business_travel_expense_documents`, `business_travel_documents`) cukup menyimpan URL yang dikembalikan endpoint tersebut, mengikuti pola `receipt_url` di `reimbursement_items`. Catatan dari `docs/analisis-modul-reimbursements.md`: reimbursement backend sudah siap tapi belum dipakai penuh dari frontend — pastikan business travel FE benar-benar memanggil endpoint upload ini end-to-end (jangan mengulang gap yang sama).

## 54.5 Routing & Permission

Routes didaftarkan di `routes.go` module, permission slugs sesuai section 41, contoh path:

```text
POST   /api/v1/tenant/business-travels
GET    /api/v1/tenant/business-travels
GET    /api/v1/tenant/business-travels/:id
PUT    /api/v1/tenant/business-travels/:id
POST   /api/v1/tenant/business-travels/:id/submit
POST   /api/v1/tenant/business-travels/:id/fundings
POST   /api/v1/tenant/business-travels/:id/settlement
...
```

(prefix `/api/v1/tenant/` mengikuti konvensi endpoint reimbursement/attendance, sesuaikan dengan prefix aktual di `routes.go` module lain saat implementasi — verifikasi via grep sebelum menulis kode).

Tidak ada module baru untuk diregistrasi di `backend/cmd/server/main.go` (lihat §54.1) — cukup tambahkan wiring approval untuk business travel service di titik yang sama dengan wiring `attendanceSvc` yang sudah ada, agar approval callback terdaftar sebelum server start (`approvalSvc.RegisterStatusHandler("business_travel", attendanceSvc.HandleBusinessTravelApprovalStatusChange)`).

## 54.6 Struktur Frontend

**Status: diimplementasikan (versi ringkas)**, di `frontend/tenant/src/views/modules/attendance/business-travel/`:

```text
frontend/tenant/src/views/modules/attendance/business-travel/
├── BusinessTravelList.vue     # daftar + filter status + dialog create (tanpa field funding)
└── BusinessTravelDetail.vue    # halaman tab tunggal: Info, Funding, Expenses, Settlement, Refund & Reimbursement
```

Deviasi dari rencana awal (2 file, bukan 9), alasan pragmatis:

- **Form digabung ke List** sebagai dialog create, bukan `BusinessTravelForm.vue` terpisah — request perjalanan dinas cukup sederhana (title/purpose/dates/origin/description) untuk satu dialog, edit-while-DRAFT belum ada UI-nya (baru create).
- **Funding/Expenses/Settlement/Reimbursements digabung jadi tab di dalam satu `BusinessTravelDetail.vue`** (pola manual tab-button, bukan `TabView` PrimeVue — mengikuti `PayrollRunDetail.vue`/`ApplicationDetail.vue`), bukan file terpisah — konteks travel yang sama dipakai lintas tab (participants untuk funding/expense per-participant, dsb), lebih mudah dikelola dalam satu komponen daripada mem-props-drilling antar file.
- **Participants/Destinations ditampilkan read-only** di tab Info (diisi saat create), belum ada UI tambah/edit setelah travel dibuat meski endpoint backend untuk activities/schedules sudah dipakai (ada dialog add).
- **`BusinessTravelDashboard.vue`/`BusinessTravelReports.vue` (§46-47) belum dibuat** — Phase 12 backend (reporting endpoints) juga belum ada.
- **Master data (funding method, expense category) tidak punya halaman admin terpisah** — quick-add dialog inline di dalam form Funding/Expense (`POST /business-travel-funding-methods`, `POST /business-travel-expense-categories`) karena belum ada tempat lain untuk membuatnya.
- **Upload dokumen (transfer proof, receipt) sudah diwire di FE** — tombol paperclip per baris Funding/Expense di `BusinessTravelDetail.vue`, two-step (upload ke endpoint generik → attach URL ke `.../fundings/:fundingId/documents` atau `.../expenses/:expenseId/documents`), pola sama dengan `AttendanceOvertime.vue`'s `attachment_url`.

Wiring aktual:

- Routes: `/attendance/business-travel` (list) & `/attendance/business-travel/:id` (detail) di `frontend/tenant/src/router/index.js`, `module: 'attendance'` (bukan module terpisah, konsisten dengan §54.1).
- Menu: card baru di `Attendance.vue` (pola sama dengan card Overtime/Corrections), **bukan** entri Sidebar terpisah — Business Travel diakses lewat halaman Attendance seperti Overtime/Corrections, bukan top-level module di sidebar.
- i18n: section `"business_travel": {...}` baru ditambahkan ke `en.json` & `id.json` (bukan nested di dalam `"attendance"`).
- Permission: **belum digating di FE** (tombol/menu tidak dicek `hasPermission('business_travel.*')` sama sekali) — konsisten dengan catatan §49 Phase 1 bahwa codebase ini tidak punya pola enforcement permission granular per-action, hanya module-level.
- Verifikasi: `npm run build` sukses tanpa error. **Belum ditest di browser nyata** (belum ada environment jalan) — ini murni verifikasi kompilasi, bukan verifikasi fungsional UI.

## 54.7 Integrasi dengan Module Reimbursement — Bergantung Subscription

Business Travel punya tabel `business_travel_reimbursements` sendiri (section 36) untuk mencatat claim additional reimbursement hasil settlement. Namun payout/pemrosesan pembayaran reimbursement idealnya memakai module **Reimbursement** yang sudah ada (`backend/internal/modules/reimbursement/`) jika tenant berlangganan module tersebut — jangan duplikasi logic payout.

Aturan integrasi:

```text
IF company subscribed ke module "reimbursement"
    → business_travel_reimbursements di-push/link ke Reimbursement module
      (mis. auto-create ReimbursementRequest, atau tampilkan link/status dari module tsb)
ELSE
    → business_travel_reimbursements diproses & ditutup mandiri di dalam
      Business Travel module (tanpa dependency ke module reimbursement)
```

Implementasi mengikuti pola `ModuleSubscriptionChecker` yang sudah dipakai module **Approval** (`backend/internal/modules/approval/service.go:15-24`, `ensureModuleSubscribed` di `service.go:266`):

```go
type ModuleSubscriptionChecker interface {
    IsModuleActive(companyID, moduleSlug string) (bool, error)
    ListActiveModules(companyID string) ([]string, error)
}
```

- Business Travel service mendefinisikan interface sempit yang sama (atau reuse tipe yang sama jika sudah diekspor), lalu di-set via `SetModuleChecker(mc)`.
- Sebelum membuat/mem-push claim reimbursement, panggil `moduleChecker.IsModuleActive(companyID, "reimbursement")`.
- Jika `true`: buat/panggil integrasi ke module reimbursement (misal via service call langsung antar-module, atau lewat interface adapter serupa `ApprovalEngine`/`Notifier` — **jangan** import package `reimbursement` langsung dari `businesstravel` kalau bisa dihindari; pakai narrow interface + adapter di `main.go`, sesuai catatan di `approval/service.go:15-20`).
- Jika `false`: business travel tetap berjalan penuh, hanya saja reimbursement diproses & ditutup secara mandiri di dalam module business travel (status `PAID` di-set langsung oleh Finance/Admin lewat endpoint business travel sendiri, tanpa keterikatan ke module lain).
- Wiring `SetModuleChecker` dilakukan di `main.go`, sama seperti `approvalSvc.SetModuleChecker(...)` — cek dulu adapter `ModuleSubscriptionChecker` yang membungkus `modulemgmt.Service` apakah sudah reusable/exported, supaya tidak membuat adapter duplikat.
- Refund (section 35) **tidak** melalui pengecekan ini karena refund adalah uang kembali ke perusahaan, bukan pembayaran ke employee — tidak relevan dengan module reimbursement.

Tambahkan test khusus (`approval_integration_test.go`-style, mis. `reimbursement_subscription_test.go`) yang memverifikasi kedua cabang: subscribed → terintegrasi, tidak subscribed → mandiri, meniru pola `approval/module_subscription_test.go`.

## 54.8 Integrasi Payroll (§38 plan doc)

Business Travel tidak mengimpor package `payroll` secara langsung — narrow-interface-plus-adapter, pola sama seperti `ApprovalEngine`/`Notifier`/`ModuleSubscriptionChecker`:

```go
type PayrollAdjustmentProvider interface {
    CreateAdjustment(ctx context.Context, employeeID, salaryComponentID uuid.UUID, periodYear, periodMonth int, amount float64, sourceType, reason string) error
}
```

Diwire di `main.go` via `Service.SetPayrollAdjuster(attendancePayrollAdjusterAdapter{repo: payrollRepo})`, di mana adapter itu murni memanggil `payrollRepo.CreateSalaryEmployeeAdjustment` yang sudah ada — **tidak ada kode baru di module `payroll`**.

**Keputusan desain — `payroll_treatment` di-repurpose jadi UUID component**: field `business_travel_expense_categories.payroll_treatment` (VARCHAR di migration `127`) sekarang diisi UUID `SalaryComponent` milik Payroll secara langsung (bukan label bebas seperti "TRAVEL_ALLOWANCE" yang tersirat di §12/§38). Alasan: Payroll module hanya punya `FindSalaryComponentByID`, tidak ada lookup by-code — menambah endpoint lookup baru ke module Payroll (module produksi dengan kalkulasi pajak sensitif) dianggap risiko yang tidak sepadan untuk fitur ini. Admin yang mengelola expense category harus tahu UUID component target saat setup (tidak ada dropdown picker di FE untuk ini — **gap FE yang belum ditutup**).

Trigger: `pushBusinessTravelPayrollAdjustments` dipanggil dari `HandleSettlementApprovalStatusChange` saat settlement APPROVED (bukan saat expense dibuat, karena expense bisa direvisi/dihapus sebelum settlement final). Untuk setiap `SettlementItem` ber-`ItemType=ACTUAL`: resolve `Expense` → `ExpenseCategory.PayrollTreatment` (harus valid UUID) → `Expense.ParticipantID` → `Participant.EmployeeID`. Item yang tidak bisa diresolusi (kategori tanpa payroll_treatment, atau expense tanpa participant employee) dilewati diam-diam — ini yang memastikan Refund/Reimbursement/Advance **tidak** otomatis dianggap salary (Rule di §38), karena hanya actual expense item yang eksplisit dikonfigurasi payroll-eligible yang didorong.

Periode payroll (`period_year`/`period_month`) diambil dari `travel.EndDate`, bukan tanggal settlement disetujui — asumsinya expense atribusinya ke bulan perjalanan terjadi, bukan bulan settlement diproses (bisa beda bulan). `SourceType` diisi `"BUSINESS_TRAVEL"` dan `Status` langsung `"APPROVED"` (bukan `"DRAFT"`) supaya otomatis terhitung payroll run berikutnya tanpa approval manual tambahan — settlement sendiri sudah melalui approval module.

**Belum ada test khusus** untuk alur ini (`payroll_integration_test.go`-style) — perlu ditambahkan sebelum dianggap production-ready.

## 54.9 Urutan Kerja yang Disarankan

Urutan ini menyelaraskan Phase 1–12 (section 49) dengan struktur BE/FE di atas agar tiap tahap punya deliverable yang bisa dites end-to-end:

1. **Migrations 124–135** (semua tabel sekaligus, karena banyak FK saling terkait) + `model.go` lengkap.
2. **Module skeleton**: `module.go`, `routes.go`, registrasi kosong di `main.go` (tanpa approval dulu) → smoke test module ter-load.
3. **Travel CRUD + status** (Phase 1–2): repository/service/handler/routes untuk `business_travels`, participants, destinations, activities, schedules. FE: `BusinessTravelList.vue`, `BusinessTravelForm.vue`, `BusinessTravelDetail.vue`.
4. **Approval integration** (Phase 3): submit → approval instance, callback handler, sinkronisasi status. Tes dengan `approval_integration_test.go` seperti reimbursement.
5. **Funding** (Phase 4): funding methods master + seed default (`DEPOSIT`, `REIMBURSEMENT`, `COMPANY_PAID`, `OTHER`), funding transaction, upload transfer proof. FE: `BusinessTravelFunding.vue`.
6. **Expense** (Phase 5): expense category master + seed, actual expense CRUD, expense document upload. FE: `BusinessTravelExpenses.vue`.
7. **Settlement** (Phase 6–8): kalkulasi advance/actual/refund/reimbursement, settlement approval (approval instance kedua), refund & reimbursement flow. FE: `BusinessTravelSettlement.vue`, `BusinessTravelReimbursements.vue`.
8. **Attendance integration** (Phase 9): isi `attendance.source_type = BUSINESS_TRAVEL` saat travel APPROVED/IN_PROGRESS, sesuai `business_travel_attendance_rules`.
9. **Payroll/Accounting** (Phase 10): `payroll_treatment` pada expense category dipakai payroll engine existing; accounting mapping configurable (tabel/interface baru jika belum ada — cek dulu apakah module `payroll` sudah punya generic accounting-mapping yang bisa dipakai ulang).
10. **Documents, Dashboard, Reporting** (Phase 11–12): `BusinessTravelDashboard.vue`, `BusinessTravelReports.vue`.

Sebelum mulai coding, disarankan membuat dokumen analisis khusus (`docs/analisis-modul-business-travel.md`) dengan pola sama seperti `docs/analisis-modul-reimbursements.md`, memetakan file & baris exact di `main.go` untuk wiring, supaya tim implementasi punya referensi presisi seperti yang sudah ada untuk reimbursement.
