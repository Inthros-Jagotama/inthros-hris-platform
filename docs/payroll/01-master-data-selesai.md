# Payroll — Master Data & Approval (✅ Selesai — BACKEND SAJA)

> Bagian-bagian di file ini **sudah diimplementasikan di backend** per audit 2026-08-12, re-verifikasi 2026-08-14. Ref: [00-overview.md](00-overview.md) untuk prinsip & roadmap keseluruhan.
>
> 🚫 **Frontend: BELUM ADA SAMA SEKALI (0%).** `frontend/tenant/src/views/modules/payroll/Payroll.vue` hanya berisi placeholder 1 baris ("Payroll Module — Coming soon"). Route `/payroll` (`router/index.js:228-231`) dan menu sidebar (`layouts/Sidebar.vue:394,437-438`, digate permission `payroll.view`) sudah terdaftar dan mengarah ke halaman placeholder itu, tapi **tidak ada UI apa pun** untuk salary component, salary structure, payroll period/run, BPJS, PPh21, atau payslip — dan tidak ada payroll API service file di frontend sama sekali. Semua "✅ selesai" di file ini murni backend/API; belum bisa dipakai user tanpa membangun FE dari nol.

## 4. Salary Component

Salary component adalah unit dasar Payroll.

### Jenis component

#### Earnings

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

#### Employee Deduction

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

#### Employer Contribution

```text
BPJS_KES_ER
JHT_ER
JKK_ER
JKM_ER
JP_ER
JKP_ER
```

### Atribut component

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

> ✅ Tabel aktual: `salary_components` — CRUD lengkap.

---

## 5. Calculation Type

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

> ⚠️ Kolom `calculation_type`/`formula` sudah ada di skema `salary_components`, tapi eksekusi formula-nya sendiri **belum ada** — lihat [02-formula-engine.md](02-formula-engine.md).

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

## 7. Salary Structure

> ✅ **Status: selesai**, tapi via nama tabel berbeda dari rencana awal — lihat §34 di bawah untuk skema aktual.

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

Tabel rencana awal (tidak dipakai, diganti — lihat §34):

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

## 8. Payroll Policy

> ❌ **Tidak ada** entity `PayrollPolicy` terpusat — konsep "policy" tersebar di `bpjs_settings`/`pph21_settings`/`employee_payroll_profiles` tanpa satu tabel payroll-policy terpusat. Bagian ini masih rencana, belum jadi prioritas revisi karena bukan blocker calculation engine.

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

## 26. Approval

> ✅ **Status: sudah diimplementasikan dan berfungsi nyata** — satu-satunya bagian non-CRUD di modul ini yang benar-benar jalan. `service.go` mendefinisikan interface `ApprovalEngine` (line 24, di-inject via `SetApprovalEngine`), `UpdatePayrollRunStatus` memanggil `approvalEngine.CreateApprovalInstance(ctx, "payroll", ...)` saat transisi DRAFT→CALCULATED (line 932) dan menyimpan `approval_instance_id` (kolom ditambahkan migration `060_payroll_approval_instance.sql`). `CheckPayrollRunApproval`/`HandleApprovalStatusChange` (line 996, 1051) menangani sinkronisasi status balik dari modul Approval (push-callback), endpoint `GET /runs/:id/approval` tersedia.

Payroll harus menggunakan existing Approval Module.

Flow (desain — implementasi aktual hanya punya fallback manual "REVIEWED" ketika tidak ada `ApprovalEngine`/`FlowID` yang di-set, lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)):

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

## 34. Database Structure (skema AKTUAL)

> ⚠️ 21 tabel di bawah adalah skema **AKTUAL** (migration `006_payroll_structure` + `007_payroll_run` + `060_payroll_approval_instance`, postgres+mysql, sudah live), **bukan rencana**. Berbeda dari desain awal dokumen ini (mis. tidak ada `salary_structures`/`bpjs_programs`/`tax_profiles` generik) — penamaan di bawah adalah keputusan arsitektur yang sudah tertanam di 21 GORM struct (`backend/internal/modules/payroll/model.go`). Referensi lengkap kolom: `docs/database-schema.md` §Payroll.

Master data & konfigurasi gaji (migration `006_payroll_structure`):

```text
salary_components            -- master komponen gaji (bukan salary_component_rules)
salary_grade_components       -- default komponen per grade, effective-dated (pengganti salary_structures)
salary_employee_components    -- override per employee, effective-dated (pengganti employee_salary_structures)
salary_change_logs            -- audit perubahan salary_employee_components
salary_employee_adjustments   -- penyesuaian sekali-jalan per periode (pengganti payroll_adjustments/bonuses/deductions)

payroll_periods                -- status hanya OPEN/CLOSED
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

Payroll run & processing (migration `007_payroll_run`, 6 tabel; `060_payroll_approval_instance` **bukan tabel baru** — hanya `ALTER TABLE payroll_runs ADD COLUMN approval_instance_id`) — lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md) untuk status masing-masing tabel:

```text
payroll_runs                   -- + approval_instance_id (ditambahkan via migration 060); status DRAFT/CALCULATED/REVIEWED/APPROVED/LOCKED/CANCELLED
payroll_run_employees          -- skema ada, TIDAK PERNAH DIISI kode apa pun
payroll_run_items              -- skema ada, TIDAK PERNAH DIISI kode apa pun
payroll_payslips               -- skema ada; TIDAK ADA endpoint/service untuk generate
pph21_calculation_logs         -- skema ada (formula_json dll.); TIDAK PERNAH DITULIS kode apa pun
payroll_profile_change_logs    -- audit perubahan profile (employee_payroll_profiles dll.), bukan audit kalkulasi
```

**Tidak ada tabel apa pun** untuk: `payroll_payments` (payment/bank transfer batch — belum ada sama sekali), `payroll_audit_logs` run-level (hanya ada log perubahan config seperti di atas).

Seluruh primary key menggunakan UUID (`CHAR(36)`), konsisten dengan modul lain di repo ini.

---

## 35. Relasi Utama

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

## Phase 1 — Foundation (Development Phase checklist)

> ✅ Sebagian besar Phase 1 sudah selesai (audit 2026-08-12).

- [x] Create Payroll module. — `backend/internal/modules/payroll/`, module.go, `IsCore: true`.
- [x] Create payroll permissions. — terdaftar di module.go (lihat §32 di [00-overview.md](00-overview.md); verifikasi granular belum di-cross-check ulang di audit ini).
- [x] Create payroll period. — tabel `payroll_periods`, tapi status hanya OPEN/CLOSED, bukan 10-state yang direncanakan.
- [ ] Create payroll policy. — **tidak ada** entity `PayrollPolicy` terpusat (lihat §8 di atas).
- [x] Create salary component. — tabel `salary_components`, CRUD lengkap.
- [x] Create salary structure. — via `salary_grade_components` (default per grade) + `salary_employee_components` (override per employee).
- [x] Create employee salary assignment. — `salary_employee_components`, effective-dated.
- [x] Implement effective dating. — dipakai konsisten di `salary_grade_components`, `salary_employee_components`, `bpjs_rate_components`, `pph21_*`.
- [ ] Implement basic payroll calculation. — **BELUM**, lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md).

## Phase 6 — Approval

> ✅ **Sebagian besar sudah selesai dan nyata berfungsi** — satu-satunya phase non-config yang punya implementasi jalan.

- [x] Integrate existing Approval Module. — `ApprovalEngine` interface + `approval_instance_id` (migration 060) + push-callback `HandleApprovalStatusChange`.
- [ ] Payroll review. — implementasi aktual hanya fallback manual "REVIEWED" (bukan tahap terpisah eksplisit).
- [x] HR approval. — via Approval Module generik (flow ditentukan di sana, bukan hard-code payroll).
- [x] Finance approval. — via Approval Module generik.
- [x] Management approval. — via Approval Module generik.
- [x] Finalization. — status `LOCKED` sebagai state akhir.
