# Payroll Module — Overview & Arsitektur

> 📅 Revisi audit: 2026-08-12 (sinkron dengan implementasi aktual) · Status: **master data & CRUD ✅ selesai; calculation engine (inti modul ini) ❌ belum ada sama sekali**.
> ✅ **Fakta aktual (audit 2026-08-12):** modul payroll **bukan greenfield** — `backend/internal/modules/payroll/` sudah berisi ±6.289 baris kode (model.go 671, service.go 1563, repository.go 804, handler.go 603, dto.go 951, routes.go 116), 21 GORM entity, 21 tabel (migration `006_payroll_structure` + `007_payroll_run` + `060_payroll_approval_instance`), 48 handler function, 71 repository function, 43 test (repository_test.go 17 + service_test.go 26). Seluruhnya adalah **CRUD master data + status-transition sederhana** — bukan calculation engine.
> 🔎 **Sumber:** audit `backend/internal/modules/payroll/` (model.go/service.go/handler.go/repository.go/routes.go/module.go) + migration `006_payroll_structure.sql`, `007_payroll_run.sql`, `060_payroll_approval_instance.sql` (postgres+mysql) + `docs/database-schema.md` §Payroll + `docs/project-completion-dashboard.md` (baris 81, 153, 259) + `docs/openapi-report.md` (baris 115, 731+).
> 🚫 **Yang TIDAK ADA sama sekali** (meski `project-completion-dashboard.md` mengklaim modul ini "✅ Complete" — klaim itu **menyesatkan**, hanya benar untuk layer CRUD/master-data): Formula Engine, kalkulator kontribusi BPJS aktual, kalkulator PPh21 aktual, eksekusi "Calculate" pada payroll run, generate payslip PDF/HTML, payment/bank transfer batch, payroll-run-level audit trail.
> ✅ **Yang SUDAH ADA dan berfungsi nyata** (bukan cuma CRUD kosong): integrasi Approval Module untuk payroll run — lihat [01-master-data-selesai.md](01-master-data-selesai.md).
> 🔧 **Catatan konsistensi docs:** `project-completion-dashboard.md` baris 81 mencatat "Payroll & Compensation | 21 entities | 39 tests | 47 endpoints | ✅ Complete" — setelah audit ini, baris tersebut perlu diperjelas menjadi "master data & run-status scaffolding selesai; calculation engine belum" supaya tidak menyesatkan pembaca lain.
> 🚫 **Frontend (re-verifikasi 2026-08-14): 0% — belum ada implementasi sama sekali.** `frontend/tenant/src/views/modules/payroll/Payroll.vue` hanya placeholder "Coming soon". Route `/payroll` dan menu sidebar sudah terdaftar (digate `payroll.view`), tapi tidak ada UI, komponen, maupun API service file untuk payroll di frontend mana pun. Seluruh "selesai" di dokumen ini murni backend/API.

## Daftar file plan (split 2026-08-14)

Dokumen ini awalnya satu file besar (`docs/module-payroll-development-plan.md`, 1716 baris). Setelah audit, dipecah menjadi sub-plan berikut, mengikuti urutan prioritas implementasi revisi:

| # | File | Isi | Status |
|---|---|---|---|
| 00 | [00-overview.md](00-overview.md) | Tujuan, prinsip arsitektur, struktur module, security, multi-company, target arsitektur akhir, roadmap prioritas | referensi |
| 01 | [01-master-data-selesai.md](01-master-data-selesai.md) | Salary Component, Salary Structure, Payroll Policy, Database Structure, Approval — **CRUD jalan tapi skema belum lengkap (gap: formula/reference, scope assignment, versioning), frontend 0%** | 🔶 BE / 🚫 FE |
| 02 | [02-formula-engine.md](02-formula-engine.md) | Formula Engine, parser/evaluator, variable registry | ❌ prioritas #1 |
| 03 | [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md) | Payroll Period, Payroll Run, Calculation Order, Snapshot | ❌ prioritas #2 |
| 04 | [04-bpjs-engine.md](04-bpjs-engine.md) | BPJS Kesehatan/Ketenagakerjaan kalkulator | 🔶 prioritas #3 |
| 05 | [05-pph21-engine.md](05-pph21-engine.md) | PPh 21 rule engine & kalkulator | 🔶 prioritas #4 |
| 06 | [06-proration-workforce.md](06-proration-workforce.md) | Proration, integrasi Attendance/Leave/Overtime/Movement/Recruitment/Performance | ❌ prioritas #5-6 |
| 07 | [07-payslip-payment.md](07-payslip-payment.md) | Payslip generator, Payment/bank transfer batch | 🔶/❌ prioritas #7-8 |
| 08 | [08-reporting-testing.md](08-reporting-testing.md) | Reporting, Dashboard, Audit Trail, Testing Strategy | ❌ prioritas #9 |

File asli `docs/module-payroll-development-plan.md` tetap ada sebagai index pendek yang mengarahkan ke sini (lihat file itu).

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

## 2. Prinsip Arsitektur

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

### Prinsip utama

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

## 3. Struktur Module

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

## 9. Effective Dating dan Versioning (prinsip cross-cutting)

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

> ✅ Sudah dipakai konsisten di `salary_grade_components`, `salary_employee_components`, `bpjs_rate_components`, `pph21_*` — lihat [01-master-data-selesai.md](01-master-data-selesai.md).

---

## 32. Security

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

## 33. Multi-Company

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

## 38. Integrasi HRIS

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

## 39. Roadmap Prioritas

> ⚠️ **Urutan rencana awal sudah tidak relevan** — Salary Component, Salary Structure (via grade/employee override), Payroll Period, Approval sudah selesai (audit 2026-08-12). Prioritas nyata sekarang:

```text
1. Formula Engine          → 02-formula-engine.md          (blocker untuk semua kalkulasi)
2. Payroll Run + Snapshot  → 03-payroll-run-snapshot.md    (isi payroll_run_employees/payroll_run_items)
3. BPJS Engine             → 04-bpjs-engine.md              (baca bpjs_rate_components, hasilkan kontribusi)
4. PPh 21 Engine           → 05-pph21-engine.md             (baca pph21_*, isi pph21_calculation_logs)
5. Proration                → 06-proration-workforce.md     (dibutuhkan Formula/Payroll Run)
6. Workforce Integration    → 06-proration-workforce.md     (Attendance/Leave/Overtime sebagai input)
7. Payslip generator         → 07-payslip-payment.md         (tabel sudah siap)
8. Payment/bank transfer     → 07-payslip-payment.md         (belum ada sama sekali)
9. Reporting + Testing       → 08-reporting-testing.md
```

Yang tersisa dan harus benar sejak awal:

```text
Formula Engine
        +
Payroll Run Execution (baca salary_grade_components/salary_employee_components yang sudah ada)
        +
Payroll Snapshot (isi payroll_run_employees/payroll_run_items yang sudah ada)
```

---

## 40. Target Arsitektur Akhir

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

### Prinsip akhir

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
