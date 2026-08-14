# Payroll Module — Development Plan

> 📅 Revisi audit: 2026-08-12 (sinkron dengan implementasi aktual) · Status: **master data & CRUD ✅ selesai; calculation engine (inti modul ini) ❌ belum ada sama sekali**.
> ✅ **Fakta aktual (audit 2026-08-12):** modul payroll **bukan greenfield** — `backend/internal/modules/payroll/` sudah berisi ±6.289 baris kode (model.go 671, service.go 1563, repository.go 804, handler.go 603, dto.go 951, routes.go 116), 21 GORM entity, 21 tabel (migration `006_payroll_structure` + `007_payroll_run` + `060_payroll_approval_instance`), 48 handler function, 71 repository function, 43 test (repository_test.go 17 + service_test.go 26). Seluruhnya adalah **CRUD master data + status-transition sederhana** — bukan calculation engine.
> 🔎 **Sumber:** audit `backend/internal/modules/payroll/` (model.go/service.go/handler.go/repository.go/routes.go/module.go) + migration `006_payroll_structure.sql`, `007_payroll_run.sql`, `060_payroll_approval_instance.sql` (postgres+mysql) + `docs/database-schema.md` §Payroll + `docs/project-completion-dashboard.md` (baris 81, 153, 259) + `docs/openapi-report.md` (baris 115, 731+).
> 🚫 **Yang TIDAK ADA sama sekali** (meski `project-completion-dashboard.md` mengklaim modul ini "✅ Complete" — klaim itu **menyesatkan**, hanya benar untuk layer CRUD/master-data): Formula Engine (parser/variable registry/dependency resolver/circular-dependency detection — §6), kalkulator kontribusi BPJS aktual (§14-17), kalkulator PPh21 aktual (§18 — tabel `pph21_calculation_logs` ada tapi tidak pernah ditulisi kode manapun), eksekusi "Calculate" pada payroll run (status berpindah manual via API tanpa verifikasi ada perhitungan yang jalan — §11-13), generate payslip PDF/HTML (tabel `payroll_payslips` ada tapi tidak ada endpoint/service untuk membuatnya — §28), payment/bank transfer batch (§27 — tidak ada tabel maupun kode sama sekali), payroll-run-level audit trail (§31 — yang ada cuma `salary_change_logs`/`payroll_profile_change_logs` untuk perubahan config, bukan audit kalkulasi).
> ✅ **Yang SUDAH ADA dan berfungsi nyata** (bukan cuma CRUD kosong): integrasi Approval Module untuk payroll run (§26) — `ApprovalEngine` interface + `approval_instance_id` (migration 060) + push-callback `HandleApprovalStatusChange`, ini implementasi asli yang bisa dipakai, bukan stub.
> 🔧 **Catatan konsistensi docs:** `project-completion-dashboard.md` baris 81 mencatat "Payroll & Compensation | 21 entities | 39 tests | 47 endpoints | ✅ Complete" — setelah audit ini, baris tersebut perlu diperjelas menjadi "master data & run-status scaffolding selesai; calculation engine belum" supaya tidak menyesatkan pembaca lain.

## 1. Tujuan

Module Payroll menjadi bagian dari HRIS dan bertanggung jawab untuk:

- Menghitung gaji karyawan secara periodik.
- Mendukung komponen gaji yang dapat dikonfigurasi secara dinamis.
- Mendukung earning, deduction, dan employer contribution.
- Menghitung BPJS Kesehatan.
- Menghitung BPJS Ketenagakerjaan: JHT, JKK, JKM, JP, dan JKP.
- Mendukung PPh 21 melalui rule engine terpisah.
- Mendukung perubahan aturan berdasarkan effective date.
- Menyimpan snapshot hasil payroll agar histori tidak berubah.
- Mendukung approval payroll.
- Menghasilkan payslip.
- Mendukung payment/bank transfer.
- Menyediakan reporting dan audit trail.
- Terintegrasi dengan Employee, Organization, Recruitment, Workforce Management, Performance, Approval, dan Accounting.

---

# 2. Prinsip Arsitektur

Payroll menggunakan pendekatan **configuration-driven payroll engine**.

```text
Employee
   |
   +-- Employment
   +-- Salary Structure
   +-- Payroll Profile
           |
           v
     Payroll Engine
           |
     +-----+------+
     |            |
  Earnings    Deductions
     |            |
     +-----+------+
           |
     Statutory Rules
     +-- BPJS Kesehatan
     +-- BPJS Ketenagakerjaan
     +-- PPh 21
           |
           v
      Payroll Result
       +---------+
       |         |
    Payslip   Accounting
```

## Prinsip utama

1. Payroll Engine tidak hard-code tarif BPJS.
2. Salary component dikonfigurasi melalui database.
3. Formula dihitung melalui Formula Engine yang aman.
4. Semua rule mendukung `effective_from` dan `effective_to`.
5. Payroll yang sudah finalized menggunakan snapshot konfigurasi dan hasil perhitungan.
6. Approval menggunakan Approval Module yang sudah ada.
7. Semua ID menggunakan UUID.
8. Semua konfigurasi yang bersifat company-specific memiliki `company_id`.
9. Data payroll harus memiliki audit trail.
10. Payroll Engine harus dapat digunakan kembali untuk company dengan kebijakan berbeda.

---

# 3. Struktur Module

```text
Payroll
├── Configuration
│   ├── Payroll Policy
│   ├── Payroll Calendar
│   ├── Salary Component
│   ├── Salary Component Rule
│   └── Salary Structure
│
├── Payroll Engine
│   ├── Formula Engine
│   ├── Calculation Context
│   ├── Dependency Resolver
│   ├── Proration Engine
│   └── Snapshot
│
├── Statutory
│   ├── BPJS Kesehatan
│   ├── JHT
│   ├── JKK
│   ├── JKM
│   ├── JP
│   ├── JKP
│   └── PPh 21
│
├── Processing
│   ├── Payroll Period
│   ├── Payroll Run
│   ├── Payroll Calculation
│   └── Payroll Finalization
│
├── Approval
├── Payment
├── Payslip
├── Reporting
└── Audit
```

---

# 4. Salary Component

Salary component adalah unit dasar Payroll.

## Jenis component

### Earnings

```text
BASIC
POSITION_ALLOWANCE
TRANSPORT
MEAL
COMMUNICATION
HOUSING
PERFORMANCE_ALLOWANCE
OVERTIME
BONUS
THR
OTHER_INCOME
```

### Employee Deduction

```text
BPJS_KES_EMP
JHT_EMP
JP_EMP
PPH21
LOAN
KASBON
ABSENCE
LATE
OTHER_DEDUCTION
```

### Employer Contribution

```text
BPJS_KES_ER
JHT_ER
JKK_ER
JKM_ER
JP_ER
JKP_ER
```

## Atribut component

Minimal:

```text
id
company_id
code
name
description
type
calculation_type
formula
taxable
bpjs_basis
proratable
recurring
sequence
active
effective_from
effective_to
created_at
updated_at
```

---

# 5. Calculation Type

Component minimal mendukung:

```text
FIXED
PERCENTAGE
FORMULA
REFERENCE
MANUAL
```

Contoh:

### Fixed

```text
TRANSPORT = 500000
```

### Percentage

```text
JHT_EMP = BPJS_WAGE * 2%
```

### Formula

```text
OVERTIME_HOURS * OVERTIME_RATE
```

### Reference

Mengambil nilai dari component lain.

```text
PERFORMANCE_BONUS
    -> PERFORMANCE_RESULT
```

### Manual

Nilai dimasukkan oleh payroll processor.

```text
BONUS_SPECIAL = 2500000
```

---

# 6. Formula Engine

> ❌ **Status: BELUM DIIMPLEMENTASIKAN.** Tidak ada parser, evaluator, variable registry, dependency resolver, atau circular-dependency detection di mana pun dalam `backend/internal/modules/payroll/`. Satu-satunya jejak konsep "formula" adalah kolom `Pph21CalculationLog.FormulaJSON *string` (`model.go:629`) — kolom JSON kosong yang dimaksudkan untuk menyimpan jejak perhitungan, tapi tidak pernah ditulisi oleh kode apa pun. Seluruh bagian ini murni desain target.

Jangan menggunakan `eval()`.

Formula Engine harus memiliki parser dan evaluator sendiri.

## Contoh formula

```text
BASIC + POSITION_ALLOWANCE
```

```text
BPJS_WAGE * 2%
```

```text
OVERTIME_HOURS * OVERTIME_RATE
```

```text
GROSS - TOTAL_EMPLOYEE_DEDUCTION
```

## Variable registry

Contoh variable:

```text
BASIC
GROSS
BPJS_WAGE
TAXABLE_INCOME
WORKING_DAYS
WORKED_DAYS
ABSENCE_DAYS
OVERTIME_HOURS
OVERTIME_RATE
TOTAL_EARNINGS
TOTAL_DEDUCTIONS
```

## Formula Engine harus mendukung

- Arithmetic operation.
- Percentage.
- Parentheses.
- Component reference.
- Variable reference.
- Conditional rule.
- Rounding.
- Dependency validation.
- Circular dependency detection.
- Formula validation sebelum payroll dijalankan.

---

# 7. Salary Structure

Salary Structure mengelompokkan component yang berlaku untuk employee.

Contoh:

```text
Staff Regular
├── Basic Salary
├── Position Allowance
├── Transport
├── Meal
├── BPJS Kesehatan
├── JHT
├── JP
└── PPh 21
```

Tabel:

```text
salary_structures
salary_structure_components
employee_salary_structures
```

Assignment dapat berdasarkan:

```text
Company
Organization
Position
Employee Grade
Employment Type
Employee
```

Prioritas:

```text
Employee
   ↓
Position
   ↓
Organization
   ↓
Company
   ↓
Default Payroll Policy
```

---

# 8. Payroll Policy

Payroll Policy merupakan konfigurasi utama perusahaan.

Contoh:

```text
PT ABC - Monthly Payroll
```

Policy mengatur:

```text
Payroll Calendar
Salary Components
Salary Rules
BPJS Rules
Tax Rules
Proration
Rounding
Overtime
Absence
Cutoff
Payment Date
```

Satu company dapat memiliki beberapa policy:

```text
Regular Employee
Executive
Contract Employee
Daily Worker
```

---

# 9. Effective Dating dan Versioning

Semua rule penting harus memiliki:

```text
effective_from
effective_to
version
```

Contoh:

```text
JHT Employee

Version 1
2026-01-01
2026-12-31
2%

Version 2
2027-01-01
2.5%
```

Payroll Agustus 2026 tetap menggunakan Version 1 walaupun Version 2 dibuat kemudian.

Konfigurasi yang membutuhkan versioning:

- Salary Component.
- Salary Component Rule.
- Salary Structure.
- Payroll Policy.
- BPJS Rule.
- Tax Rule.
- Overtime Rule.
- Proration Rule.

---

# 10. Payroll Period

Tabel:

```text
payroll_periods
```

Contoh:

```text
August 2026

period_start : 2026-08-01
period_end   : 2026-08-31
cutoff_date  : 2026-08-20
payment_date : 2026-08-25
```

Status (rencana awal — **belum diimplementasikan seperti ini**):

```text
DRAFT
CALCULATING
CALCULATED
REVIEW
APPROVAL
APPROVED
PROCESSING
PAID
CLOSED
CANCELLED
```

> ⚠️ **Aktual:** `PayrollPeriod.Status` (`model.go:176`) hanya mendukung **`OPEN`/`CLOSED`** — jauh lebih sederhana dari 10-state di atas. Status detail (DRAFT/CALCULATED/REVIEWED/APPROVED/LOCKED/CANCELLED) ternyata ada di level `PayrollRun`, bukan `PayrollPeriod` — lihat §11 di bawah untuk state machine yang benar-benar berjalan.

---

# 11. Payroll Run

Payroll Run digunakan untuk menjalankan calculation terhadap employee dalam satu period.

Flow:

```text
Create Period
      |
      v
Select Employees
      |
      v
Load Employee Data
      |
      v
Load Salary Structure
      |
      v
Load Attendance
      |
      v
Load Leave
      |
      v
Load Overtime
      |
      v
Load Bonus/Adjustment
      |
      v
Calculate Earnings
      |
      v
Calculate Employer Contribution
      |
      v
Calculate Employee Deduction
      |
      v
Calculate Tax
      |
      v
Calculate Net Salary
      |
      v
Validation
      |
      v
Review
      |
      v
Approval
      |
      v
Finalize
```

> ⚠️ **Aktual (`service.go:915-991`, `UpdatePayrollRunStatus`):** flow di atas **belum berjalan** — hanya kerangka status-transition manual yang ada, tanpa langkah "Load Attendance/Leave/Overtime" atau "Calculate *" mana pun benar-benar dieksekusi:
>
> ```text
> DRAFT → CALCULATED   (kalau ApprovalEngine + FlowID di-set: buat approval instance;
>                        kalau tidak, langsung lompat ke REVIEWED)
> CALCULATED → APPROVED
> APPROVED → LOCKED
> ```
>
> Transisi `DRAFT → CALCULATED` hanya mengubah kolom `status` — **tidak ada kode yang mengisi `payroll_run_employees`/`payroll_run_items`, menghitung BPJS/PPh21, atau memverifikasi bahwa perhitungan benar-benar terjadi**. Siapa pun bisa memanggil endpoint status-update ini tanpa data payroll run pernah dihitung. Ini adalah gap implementasi paling kritis di seluruh modul — lihat §12-§13.

---

# 12. Calculation Order

> ❌ **Status: BELUM DIIMPLEMENTASIKAN.** Tidak ada kode di `backend/internal/modules/payroll/` yang menjalankan urutan perhitungan berikut. Bagian ini murni rencana/desain target.

Urutan perhitungan harus eksplisit.

```text
1. Basic Salary
2. Fixed Allowance
3. Variable Allowance
4. Overtime
5. Bonus
6. Gross Salary
7. BPJS Wage
8. Employer Contribution
9. Taxable Income
10. Employee BPJS Deduction
11. PPh 21
12. Other Deduction
13. Net Salary
14. Employer Total Cost
```

Formula dasar:

```text
GROSS =
    TOTAL_EARNINGS
```

```text
TOTAL_EMPLOYEE_DEDUCTION =
    BPJS_EMPLOYEE
    + PPH21
    + OTHER_DEDUCTION
```

```text
NET_SALARY =
    GROSS
    - TOTAL_EMPLOYEE_DEDUCTION
```

```text
EMPLOYER_TOTAL_COST =
    GROSS
    + TOTAL_EMPLOYER_CONTRIBUTION
```

---

# 13. Payroll Snapshot

> ❌ **Status: BELUM DIIMPLEMENTASIKAN.** Tabel `payroll_run_employees` dan `payroll_run_items` sudah ada di migration `007_payroll_run.sql` dengan kolom yang mendekati desain di bawah (component_id, calculation_type, base_amount, dll.), tapi **tidak pernah diisi oleh kode apa pun** — tidak ada service method yang melakukan insert ke kedua tabel ini. Snapshot secara konsep sudah disiapkan skemanya, tapi belum ada satu baris data pun yang pernah dibuat karena §12 (calculation) belum berjalan.

Payroll yang sudah dihitung harus menyimpan snapshot.

Jangan bergantung pada konfigurasi live.

Contoh:

```text
payrolls
payroll_items
```

Payroll item menyimpan:

```text
component_id
component_code
component_name
calculation_type
base_amount
rate
calculated_amount
employee_amount
employer_amount
formula
formula_result
```

Contoh:

```text
JHT_EMP

Base:
10,000,000

Rate:
2%

Formula:
BPJS_WAGE * 2%

Result:
200,000
```

Tujuan:

- Histori tidak berubah.
- Audit mudah.
- Payslip konsisten.
- Payroll dapat dijelaskan kembali.
- Perubahan konfigurasi tidak mengubah payroll lama.

---

# 14. BPJS Configuration

> 🔶 **Status: konfigurasi ✅ ada (CRUD lengkap), kalkulator kontribusi ❌ belum ada.** Tabel aktual bernama `bpjs_settings` + `bpjs_rate_components` (bukan `bpjs_programs`/`bpjs_rules` seperti di bawah) — lihat §34 untuk skema sebenarnya. `bpjs_rate_components` sudah mendukung effective dating, `rate_percent`, `min/max_base_amount`, kolom enum `bpjs_program` (HEALTH/JHT/JP/JKK/JKM/JKP) dan `paid_by` (EMPLOYEE/EMPLOYER) — jadi tarif memang tidak di-hard-code, sesuai prinsip di bawah. Tapi **tidak ada kode yang membaca tabel ini untuk menghasilkan angka kontribusi employee/employer** — CRUD-nya lengkap, konsumennya (payroll calculation) belum ada.

Gunakan tabel:

```text
bpjs_programs
bpjs_rules
```

Program:

```text
BPJS_KES
JHT
JKK
JKM
JP
JKP
```

Rule minimal:

```text
program_id
employee_percentage
employer_percentage
calculation_basis
minimum_wage
maximum_wage
fixed_amount
effective_from
effective_to
```

Tarif tidak boleh ditanam langsung dalam source code.

Default rule harus dapat diperbarui berdasarkan regulasi terbaru.

---

# 15. BPJS Kesehatan

Konfigurasi mendukung:

```text
Employee Contribution
Employer Contribution
Wage Basis
Minimum Wage
Maximum Wage
Effective Date
```

Contoh konfigurasi awal:

```text
Employee : 1%
Employer : 4%
Total    : 5%
```

Catatan: nilai tersebut merupakan default konfigurasi dan harus dapat diperbarui jika regulasi berubah.

---

# 16. BPJS Ketenagakerjaan

Program:

```text
JHT
JKK
JKM
JP
JKP
```

Contoh struktur rule:

```text
JHT
├── Employee
└── Employer

JKK
└── Employer

JKM
└── Employer

JP
├── Employee
└── Employer

JKP
└── Rule khusus
```

Default rule awal dapat di-seed berdasarkan ketentuan BPJS yang berlaku, tetapi tetap harus dikelola melalui konfigurasi dan effective dating.

---

# 17. JKK Risk Level

JKK membutuhkan konfigurasi risk level.

Tabel:

```text
bpjs_jkk_risk_levels
```

Contoh:

```text
VERY_LOW
LOW
MEDIUM
HIGH
VERY_HIGH
```

Relasi:

```text
Company
   ↓
Business Sector
   ↓
Risk Level
   ↓
JKK Rate
```

Tarif JKK tidak boleh di-hard-code.

---

# 18. PPh 21

> 🔶 **Status: konfigurasi ✅ ada (CRUD lengkap), kalkulator pajak ❌ belum ada.** Tabel aktual: `pph21_settings` + `pph21_ptkp_rates` + `pph21_tax_brackets` (bukan `tax_profiles`/`tax_rules`/`tax_rule_brackets` generik). Tabel `pph21_calculation_logs` (migration `007_payroll_run.sql`) sudah menyediakan kolom lengkap untuk jejak kalkulasi (`gross_monthly`, `ptkp_annual`, `pkp_annual`, `pph21_monthly`, `formula_json`, dll.) — tapi **tidak ada satu baris pun yang pernah ditulis ke tabel ini**, karena tidak ada kode kalkulator PPh21 yang berjalan. Sama seperti BPJS: config-nya siap dan effective-dated, konsumennya (payroll calculation) belum dibangun.

PPh 21 dibuat sebagai statutory rule engine terpisah dari BPJS.

Komponen:

```text
Tax Profile
Tax Status
Taxable Income
Tax Rule
Tax Calculation
Tax Adjustment
```

Struktur:

```text
PPh21 Engine
├── Employee Tax Profile
├── Taxable Components
├── Deductible Components
├── Tax Rule
├── Tax Calculation
└── Tax Result
```

Rule harus mendukung effective dating karena regulasi pajak dapat berubah.

---

# 19. Proration

Payroll harus mendukung:

- New Employee.
- Resignation.
- Join tengah bulan.
- Salary change.
- Promotion.
- Unpaid leave.
- Absence.
- Transfer.
- Employment termination.

Metode proration:

```text
CALENDAR_DAYS
WORKING_DAYS
FIXED_30_DAYS
ATTENDANCE_DAYS
```

Contoh:

```text
Basic Salary = 10,000,000
Working Days = 22
Eligible Days = 12

Prorated Salary =
10,000,000 / 22 * 12
```

Metode harus configurable berdasarkan Payroll Policy.

---

# 20. Integrasi Attendance

```text
Attendance
    |
    +-- Working Days
    +-- Absence
    +-- Late
    +-- Attendance Adjustment
             |
             v
          Payroll
```

Payroll tidak menghitung attendance dari awal.

Payroll menggunakan hasil final dari Workforce Management.

---

# 21. Integrasi Overtime

```text
Overtime
   |
   +-- Overtime Hours
   +-- Overtime Rate
   +-- Overtime Approval
          |
          v
       Payroll
```

Contoh:

```text
OVERTIME =
OVERTIME_HOURS * OVERTIME_RATE
```

---

# 22. Integrasi Leave

Payroll mengambil informasi:

```text
Paid Leave
Unpaid Leave
Half Day
Special Leave
```

Contoh:

```text
Unpaid Leave
    ↓
Absence/Unpaid Deduction
    ↓
Payroll
```

---

# 23. Integrasi Employee Movement & Career

Promotion:

```text
Promotion
    ↓
Position Change
    ↓
Salary Structure Change
    ↓
Effective Date
    ↓
Payroll
```

Salary change harus memiliki effective date sehingga payroll memilih struktur yang berlaku pada period tersebut.

---

# 24. Integrasi Recruitment

Flow:

```text
Recruitment
    ↓
Employee
    ↓
Employment
    ↓
Salary Structure
    ↓
Payroll Profile
```

Data gaji tidak perlu dimasukkan ulang.

---

# 25. Integrasi Performance

Performance dapat menghasilkan:

```text
Bonus
Performance Allowance
Incentive
```

Contoh:

```text
Performance Score = 92

Bonus Rule:
92 - 100
    -> 10% Basic Salary
```

Performance Engine tetap bertanggung jawab atas scoring.

Payroll hanya menerima nilai final:

```text
PERFORMANCE_BONUS = X
```

---

# 26. Approval

> ✅ **Status: sudah diimplementasikan dan berfungsi nyata** — satu-satunya bagian non-CRUD di modul ini yang benar-benar jalan. `service.go` mendefinisikan interface `ApprovalEngine` (line 24, di-inject via `SetApprovalEngine`), `UpdatePayrollRunStatus` memanggil `approvalEngine.CreateApprovalInstance(ctx, "payroll", ...)` saat transisi DRAFT→CALCULATED (line 932) dan menyimpan `approval_instance_id` (kolom ditambahkan migration `060_payroll_approval_instance.sql`). `CheckPayrollRunApproval`/`HandleApprovalStatusChange` (line 996, 1051) menangani sinkronisasi status balik dari modul Approval (push-callback), endpoint `GET /runs/:id/approval` tersedia. Flow di bawah ini mencerminkan desain yang sudah terpakai, bukan sekadar rencana — tapi catatan: alur di bawah masih menyebut "Payroll Review" sebagai langkah terpisah sebelum HR Approval, sedangkan implementasi aktual hanya punya fallback manual "REVIEWED" ketika tidak ada `ApprovalEngine`/`FlowID` yang di-set (lihat §11).

Payroll harus menggunakan existing Approval Module.

Flow:

```text
Calculate
   ↓
Payroll Review
   ↓
HR Approval
   ↓
Finance Approval
   ↓
Management Approval
   ↓
Finalize
```

Jangan membuat Approval Engine khusus Payroll.

---

# 27. Payment

> ❌ **Status: BELUM DIIMPLEMENTASIKAN — tidak ada tabel, struct, atau kode sama sekali.** Tidak ada `payroll_payments` atau padanannya di skema aktual (§34). `EmployeeBankProfile` (tabel `employee_bank_profiles`) sudah menyimpan rekening bank per employee untuk keperluan payroll, tapi tidak ada kode yang mengonsumsinya untuk menghasilkan payment/disbursement record. Seluruh bagian ini murni desain target, prioritas paling akhir karena bergantung pada calculation engine (§12-13) yang juga belum ada.

Tabel:

```text
payroll_payments
```

Data:

```text
employee_id
bank_id
account_number
account_name
amount
payment_date
payment_status
reference
```

Status:

```text
PENDING
PROCESSING
PAID
FAILED
REVERSED
```

Nomor rekening untuk payment batch harus menggunakan snapshot agar perubahan rekening setelah payroll finalized tidak mengubah batch sebelumnya.

---

# 28. Payslip

> 🔶 **Status: skema tabel ✅ ada, generator ❌ belum ada.** Tabel `payroll_payslips` sudah ada (migration `007_payroll_run.sql`) dengan kolom `payslip_number`, `total_earning/deduction`, `net_amount`, `status DRAFT/PUBLISHED/CANCELLED`, `generated_at/published_at/cancelled_at` — tapi **tidak ada endpoint `/payslips`, tidak ada service method `CreatePayslip`/`GeneratePayslip`, dan tidak ada kode generate PDF/HTML apa pun** di seluruh modul. Fungsional ini stub murni: tabelnya siap, tidak ada yang mengisinya.

Payslip dibuat berdasarkan finalized payroll.

Isi minimal:

```text
Company
Employee
Employee Number
Position
Payroll Period

Earnings
--------
Basic Salary
Allowance
Overtime
Bonus

Deductions
----------
BPJS Kesehatan
JHT
JP
PPh 21
Other Deduction

Employer Contribution
---------------------
BPJS Kesehatan
JHT
JKK
JKM
JP

Net Salary
```

Format:

```text
PDF
HTML
Employee Portal
```

---

# 29. Reporting

## Payroll Summary

```text
Total Employee
Gross Salary
Employee Deduction
Employer Contribution
Net Salary
Total Payroll Cost
```

## BPJS Report

```text
Employee
BPJS Number
Wage Basis
Employee Contribution
Employer Contribution
Total Contribution
```

## Tax Report

```text
Employee
Taxable Income
PPh 21
```

## Bank Transfer

```text
Employee
Bank
Account
Amount
```

## Payroll Detail

```text
Employee
Component
Base
Rate
Employee Amount
Employer Amount
```

---

# 30. Dashboard

Contoh:

```text
Payroll August 2026

Employees               1,245
Gross Salary          Rp 8.2 M
Employee Deduction    Rp 1.1 M
Employer Contribution Rp 0.8 M
Net Payroll           Rp 7.1 M
Total Employer Cost   Rp 9.0 M
```

Status:

```text
Draft
Calculated
Review
Approval
Approved
Paid
Closed
```

---

# 31. Audit Trail

> 🔶 **Status: audit config-level ✅ ada, audit kalkulasi ❌ belum ada.** `salary_change_logs` (perubahan `salary_employee_components`) dan `payroll_profile_change_logs` (perubahan `employee_payroll_profiles` dll.) sudah ada dan berfungsi — tapi keduanya mengaudit **perubahan konfigurasi**, bukan kalkulasi payroll. Tidak ada `payroll_calculation_logs`/`payroll_audit_logs` seperti di bawah; yang paling dekat adalah tabel `pph21_calculation_logs` yang skemanya sudah siap tapi tidak pernah ditulisi (lihat §18) — begitu calculation engine (§12-13) dibangun, tabel itu perlu diisi sebagai bagian dari audit kalkulasi, bukan tabel generik baru seperti direncanakan di bawah.

Tabel:

```text
payroll_calculation_logs
```

Simpan:

```text
employee_id
payroll_id
component_id
action
old_value
new_value
formula
calculation_context
user_id
created_at
```

Audit diperlukan untuk:

- Perubahan payroll.
- Manual adjustment.
- Recalculation.
- Approval.
- Finalization.
- Payment.
- Reversal.

---

# 32. Security

Permission granular:

```text
payroll.view
payroll.create
payroll.calculate
payroll.review
payroll.approve
payroll.finalize
payroll.pay
payroll.view_salary
payroll.view_payslip
payroll.export
payroll.manage_component
payroll.manage_rule
payroll.manage_policy
```

Akses nominal gaji tidak boleh otomatis diberikan kepada seluruh user HR.

---

# 33. Multi-Company

Semua konfigurasi company-specific menggunakan:

```text
company_id
```

Contoh scope:

```text
GLOBAL
COUNTRY
COMPANY
```

Contoh:

```text
BPJS Rule
    Scope: Indonesia

Payroll Policy
    Scope: Company

Salary Structure
    Scope: Company

Salary Component
    Scope: Company
```

---

# 34. Database Structure

> ⚠️ **Bagian ini sudah diimplementasikan — 21 tabel di bawah adalah skema AKTUAL** (migration `006_payroll_structure` + `007_payroll_run` + `060_payroll_approval_instance`, postgres+mysql, sudah live), **bukan lagi rencana**. Skema ini berbeda dari desain awal yang diusulkan di dokumen ini (mis. tidak ada `salary_structures`/`bpjs_programs`/`tax_profiles` generik) — penamaan di bawah adalah keputusan arsitektur yang sudah tertanam di 21 GORM struct (`backend/internal/modules/payroll/model.go`) dan butuh rewrite (bukan sekadar rename) kalau mau diselaraskan ke desain awal. Referensi lengkap kolom: `docs/database-schema.md` §Payroll.

Master data & konfigurasi gaji (migration `006_payroll_structure`):

```text
salary_components            -- master komponen gaji (bukan salary_component_rules)
salary_grade_components       -- default komponen per grade, effective-dated (pengganti salary_structures)
salary_employee_components    -- override per employee, effective-dated (pengganti employee_salary_structures)
salary_change_logs            -- audit perubahan salary_employee_components
salary_employee_adjustments   -- penyesuaian sekali-jalan per periode (pengganti payroll_adjustments/bonuses/deductions)

payroll_periods                -- status hanya OPEN/CLOSED (bukan 10-state di §10)
employee_payroll_profiles
employee_bank_profiles
employee_bpjs_profiles
employee_tax_profiles

bpjs_settings                  -- pengganti bpjs_programs (flat, bukan per-program)
bpjs_rate_components           -- pengganti bpjs_rules/bpjs_jkk_risk_levels (kolom bpjs_program enum: HEALTH/JHT/JP/JKK/JKM/JKP, paid_by EMPLOYEE/EMPLOYER)
pph21_settings                 -- pengganti tax_profiles/tax_rules
pph21_ptkp_rates
pph21_tax_brackets             -- pengganti tax_rule_brackets
```

Payroll run & processing (migration `007_payroll_run` + `060_payroll_approval_instance`):

```text
payroll_runs                   -- + approval_instance_id (migration 060); status DRAFT/CALCULATED/REVIEWED/APPROVED/LOCKED/CANCELLED (§10-11)
payroll_run_employees          -- skema ada, TIDAK PERNAH DIISI kode apa pun (§13)
payroll_run_items              -- skema ada, TIDAK PERNAH DIISI kode apa pun (§13)
payroll_payslips               -- skema ada; TIDAK ADA endpoint/service untuk generate (§28)
pph21_calculation_logs         -- skema ada (formula_json dll.); TIDAK PERNAH DITULIS kode apa pun (§18, §31)
payroll_profile_change_logs    -- audit perubahan profile (employee_payroll_profiles dll.), bukan audit kalkulasi
```

**Tidak ada tabel apa pun** untuk: `payroll_payments` (§27, payment/bank transfer batch — belum ada sama sekali), `payroll_audit_logs` run-level (§31 — hanya ada log perubahan config seperti di atas).

Seluruh primary key menggunakan UUID (`CHAR(36)`), konsisten dengan modul lain di repo ini.

---

# 35. Relasi Utama

```text
Company
   |
   +-- Payroll Policy
   |      |
   |      +-- Salary Components
   |      +-- BPJS Rules
   |      +-- Tax Rules
   |
   +-- Employee
          |
          +-- Employment
          +-- Salary Structure
          +-- Tax Profile
          |
          v
     Payroll Period
          |
          v
      Payroll Run
          |
          v
       Payroll
          |
     +----+----+
     |         |
 Payroll Items Employer Cost
     |
     v
  Payslip
```

---

# 36. Development Phase

## Phase 1 — Foundation

> ✅ Sebagian besar Phase 1 sudah selesai (audit 2026-08-12) — lihat detail per item.

- [x] Create Payroll module. — `backend/internal/modules/payroll/`, module.go, `IsCore: true`.
- [x] Create payroll permissions. — terdaftar di module.go (lihat §32 untuk daftar; verifikasi granular belum di-cross-check ulang di audit ini).
- [x] Create payroll period. — tabel `payroll_periods`, tapi status hanya OPEN/CLOSED (§10), bukan 10-state yang direncanakan.
- [ ] Create payroll policy. — **tidak ada** entity `PayrollPolicy`; konsep "policy" tersebar di `bpjs_settings`/`pph21_settings`/`employee_payroll_profiles` tanpa satu tabel payroll-policy terpusat.
- [x] Create salary component. — tabel `salary_components`, CRUD lengkap.
- [x] Create salary structure. — via `salary_grade_components` (default per grade) + `salary_employee_components` (override per employee), bukan `salary_structures`/`salary_structure_components` seperti direncanakan §7, tapi konsep effective-dated grade→employee override sudah berfungsi.
- [x] Create employee salary assignment. — `salary_employee_components`, effective-dated.
- [x] Implement effective dating. — dipakai konsisten di `salary_grade_components`, `salary_employee_components`, `bpjs_rate_components`, `pph21_*`.
- [ ] Implement basic payroll calculation. — **BELUM** — lihat §12, tidak ada kode kalkulasi apa pun berjalan, ini gap terbesar seluruh modul.

Output:

```text
Basic Salary
+ Allowance
- Deduction
= Net Salary
```

---

## Phase 2 — Payroll Engine

> ❌ **Belum ada satu item pun** — ini prioritas #1 pekerjaan berikutnya (lihat §39). `payroll_run_employees`/`payroll_run_items`/`pph21_calculation_logs` sudah punya skema tabel (calculation snapshot & audit), tapi menunggu engine ini untuk diisi.

- [ ] Build Formula Engine.
- [ ] Build variable registry.
- [ ] Build component dependency resolver.
- [ ] Detect circular dependency.
- [ ] Build calculation context.
- [ ] Build rounding engine.
- [ ] Build proration engine.
- [ ] Build calculation snapshot. — skema tabel (`payroll_run_employees`/`payroll_run_items`) sudah ada, tinggal diisi.
- [ ] Build calculation audit. — skema tabel (`pph21_calculation_logs`) sudah ada, tinggal diisi.
- [ ] Build recalculation mechanism.

---

## Phase 3 — BPJS

> 🔶 **Konfigurasi selesai, kalkulator belum.** Tabel `bpjs_settings`+`bpjs_rate_components` sudah effective-dated dan CRUD lengkap; item di bawah yang bercentang berarti "konfigurasinya ada", bukan "kontribusinya sudah terhitung otomatis" — itu masih Phase 2.

- [x] BPJS Kesehatan. — kolom `bpjs_program='HEALTH'` di `bpjs_rate_components`, config CRUD ada.
- [x] JHT. — `bpjs_program='JHT'`, config CRUD ada.
- [x] JKK. — `bpjs_program='JKK'`, config CRUD ada.
- [x] JKM. — `bpjs_program='JKM'`, config CRUD ada.
- [x] JP. — `bpjs_program='JP'`, config CRUD ada.
- [x] JKP. — `bpjs_program='JKP'`, config CRUD ada.
- [ ] Employee contribution. — **belum dihitung** (tidak ada kalkulator).
- [ ] Employer contribution. — **belum dihitung**.
- [ ] JKK risk level. — **tidak ada** tabel `bpjs_jkk_risk_levels`/relasi ke business sector seperti direncanakan §17; JKK hanya satu rate flat di `bpjs_rate_components`.
- [x] Effective-dated BPJS rules. — `bpjs_rate_components` sudah effective-dated.
- [ ] BPJS report. — belum ada (bergantung kalkulasi).

---

## Phase 4 — Workforce Integration

> ❌ **Belum ada integrasi apa pun** — payroll module tidak memanggil Attendance/Leave/Overtime module manapun (tidak ditemukan di survei kode).

- [ ] Attendance integration.
- [ ] Leave integration.
- [ ] Overtime integration.
- [ ] Shift integration.
- [ ] Employee movement integration.
- [ ] Absence deduction.
- [ ] Overtime earning.

---

## Phase 5 — Tax

> 🔶 **Konfigurasi selesai (nama tabel beda dari rencana), kalkulator belum.** Tabel aktual: `pph21_settings`+`pph21_ptkp_rates`+`pph21_tax_brackets` (§18), bukan `tax_profiles`/`tax_rules`/`tax_rule_brackets` generik.

- [x] Tax profile. — via `employee_tax_profiles` (bukan `tax_profiles` generik).
- [ ] Tax status. — belum diverifikasi ulang di audit ini apakah PTKP status (TK/0, K/1, dst.) sudah ada di `employee_tax_profiles`.
- [ ] Taxable component mapping. — tidak ditemukan mapping "component mana yang taxable" secara eksplisit.
- [ ] PPh 21 rule engine. — **belum ada**, ini bagian dari Phase 2.
- [ ] Tax calculation. — **belum ada** — `pph21_calculation_logs` tabel kosong, tidak pernah ditulisi.
- [ ] Tax adjustment.
- [ ] Tax report. — belum ada (bergantung kalkulasi).

PPh 21 harus menjadi engine terpisah dari BPJS agar perubahan regulasi dapat dikelola secara independen.

---

## Phase 6 — Approval

> ✅ **Sebagian besar sudah selesai dan nyata berfungsi** (§26) — satu-satunya phase non-config yang punya implementasi jalan.

- [x] Integrate existing Approval Module. — `ApprovalEngine` interface + `approval_instance_id` (migration 060) + push-callback `HandleApprovalStatusChange`.
- [ ] Payroll review. — implementasi aktual hanya fallback manual "REVIEWED" (bukan tahap terpisah eksplisit).
- [x] HR approval. — via Approval Module generik (flow ditentukan di sana, bukan hard-code payroll).
- [x] Finance approval. — via Approval Module generik.
- [x] Management approval. — via Approval Module generik.
- [x] Finalization. — status `LOCKED` sebagai state akhir.

---

## Phase 7 — Payment

> ❌ **Belum ada satu item pun** — tidak ada tabel maupun kode sama sekali (§27). `employee_bank_profiles` sudah ada tapi belum dikonsumsi.

- [ ] Employee bank account. — data sudah ada di `employee_bank_profiles`, tapi belum dikonsumsi proses payment.
- [ ] Payment batch.
- [ ] Bank transfer file.
- [ ] Payment status.
- [ ] Payment reconciliation.
- [ ] Payment reversal.

---

## Phase 8 — Payslip & Reporting

> 🔶 **Tabel payslip ada (§28), generator dan seluruh reporting belum ada.**

- [ ] Payslip HTML. — tabel `payroll_payslips` ada, generator tidak ada.
- [ ] Payslip PDF. — tidak ada.
- [ ] Employee portal. — tidak diverifikasi; tidak ditemukan endpoint publik payslip di survei.
- [ ] Payroll summary.
- [ ] Payroll detail.
- [ ] BPJS report.
- [ ] Tax report.
- [ ] Bank transfer report.
- [ ] Payroll journal.

---

# 37. Testing Strategy

> 🔶 **Status aktual: 43 test ada (`repository_test.go` 17 + `service_test.go` 26), tapi semuanya level CRUD/repository** — testing untuk unit test kalkulasi (basic salary, BPJS, tax, dsb.) dan scenario test di bawah **belum ada karena kode yang diuji belum ada** (§12, §14-18). Golden dataset regression test juga belum dibuat.

Payroll harus memiliki test lebih ketat dibanding module biasa.

## Unit Test

- [ ] Basic salary.
- [ ] Allowance.
- [ ] Deduction.
- [ ] Percentage.
- [ ] Formula.
- [ ] Proration.
- [ ] Overtime.
- [ ] BPJS.
- [ ] Tax.
- [ ] Bonus.
- [ ] Net salary.
- [ ] Employer cost.

## Scenario Test

- [ ] Employee baru.
- [ ] Employee resign.
- [ ] Join tengah bulan.
- [ ] Promotion.
- [ ] Salary increase.
- [ ] Unpaid leave.
- [ ] Overtime.
- [ ] Bonus.
- [ ] BPJS.
- [ ] Tax.
- [ ] Manual adjustment.

## Regression Test

Buat golden dataset:

```text
Employee A
Salary             = X
Expected Gross     = X
Expected BPJS      = X
Expected Tax       = X
Expected Net       = X
Expected Employer  = X
```

Setiap perubahan Payroll Engine harus melewati golden dataset.

---

# 38. Integrasi HRIS

Struktur HRIS:

```text
HRIS
|
+-- Organization
|
+-- Employee
|
+-- Recruitment
|
+-- Employee Movement & Career
|
+-- Workforce Management
|   +-- Attendance
|   +-- Shift
|   +-- Leave
|   +-- Overtime
|
+-- Performance
|
+-- Payroll
|   +-- Configuration
|   +-- Salary Component
|   +-- Salary Structure
|   +-- Payroll Policy
|   +-- Payroll Engine
|   +-- BPJS
|   +-- Tax
|   +-- Processing
|   +-- Approval
|   +-- Payment
|   +-- Payslip
|   +-- Reporting
|
+-- Approval
|
+-- Accounting
```

---

# 39. Prioritas Implementasi

> ⚠️ **Urutan di bawah ini rencana awal — sudah tidak relevan lagi.** Item 1, 2 (via grade/employee override), 5, 12 sudah selesai (audit 2026-08-12). Prioritas nyata sekarang adalah **Formula Engine → Payroll Run execution/Snapshot → BPJS Engine → PPh 21 Engine**, karena keempatnya saling bergantung dan tanpa itu tidak ada satu pun payroll yang benar-benar bisa dihitung — lihat urutan revisi di bawah.

Urutan prioritas (rencana awal — lihat status di §36):

```text
1. Salary Component        ✅ selesai
2. Salary Structure        ✅ selesai (via salary_grade_components/salary_employee_components)
3. Payroll Policy          ❌ tidak ada entity terpusat (§36 Phase 1)
4. Formula Engine          ❌ belum ada — lihat prioritas revisi di bawah
5. Payroll Period          ✅ selesai (status OPEN/CLOSED, lebih sederhana dari rencana)
6. Payroll Run             🔶 CRUD+status-transition ada, "Calculate" tidak pernah dieksekusi
7. Payroll Snapshot        ❌ skema tabel ada, tidak pernah diisi
8. Proration                ❌ belum ada
9. BPJS Engine              🔶 config ada, kalkulator belum
10. Workforce Integration   ❌ belum ada
11. PPh 21 Engine           🔶 config ada, kalkulator belum
12. Approval                 ✅ selesai dan berfungsi nyata
13. Payment                 ❌ belum ada sama sekali
14. Payslip                 🔶 tabel ada, generator belum
15. Reporting                ❌ belum ada
16. Accounting Integration   ❌ belum ada
```

**Prioritas revisi (2026-08-12), berdasarkan apa yang sudah ada:**

```text
1. Formula Engine        (blocker untuk semua kalkulasi — §6, §12)
2. Payroll Run execution + Snapshot (isi payroll_run_employees/payroll_run_items — §11-13)
3. BPJS Engine            (baca bpjs_rate_components, hasilkan kontribusi — §14-17)
4. PPh 21 Engine          (baca pph21_*, isi pph21_calculation_logs — §18)
5. Proration              (dibutuhkan Formula/Payroll Run untuk kasus join/resign tengah bulan — §19)
6. Workforce Integration  (Attendance/Leave/Overtime sebagai input kalkulasi — §20-22)
7. Payslip generator      (tabel sudah siap, tinggal endpoint+PDF/HTML — §28)
8. Payment/bank transfer  (belum ada sama sekali, prioritas paling akhir — §27)
9. Reporting + Accounting Integration
```

Jangan memulai dari tabel payroll sederhana lalu menambahkan fitur satu per satu — **bagian ini sudah terjadi** (tabel/CRUD sudah lengkap). Yang tersisa dan harus benar sejak awal justru komponen intinya:

```text
Formula Engine
        +
Payroll Run Execution (baca salary_grade_components/salary_employee_components yang sudah ada)
        +
Payroll Snapshot (isi payroll_run_employees/payroll_run_items yang sudah ada)
```

---

# 40. Target Arsitektur Akhir

Target akhir:

```text
                    HRIS
                     |
              Payroll Policy
                     |
             +-------+-------+
             |               |
       Salary Structure   Statutory Rules
             |               |
             |        +------+------+
             |        |             |
             |       BPJS          Tax
             |        |             |
             +--------+-------------+
                      |
                 Payroll Engine
                      |
        +-------------+-------------+
        |             |             |
     Earnings     Deductions    Employer Cost
        |             |             |
        +-------------+-------------+
                      |
                 Payroll Snapshot
                      |
             +--------+--------+
             |                 |
          Payslip          Payment
             |
             v
          Employee

Payroll
   |
   +--> Approval
   |
   +--> Accounting
   |
   +--> Reporting
```

## Prinsip akhir

**Payroll Engine tidak boleh mengetahui bisnis secara hard-code.**

Engine hanya menjalankan:

```text
Component
+
Rule
+
Formula
+
Context
+
Effective Date
=
Payroll Result
```

Dengan arsitektur ini, penambahan komponen baru seperti `meal allowance`, `transport`, `loan deduction`, `performance bonus`, atau aturan statutory baru dapat dilakukan melalui konfigurasi tanpa harus mengubah core Payroll Engine.
