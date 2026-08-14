# Payroll — Master Data & Approval (🔶 Backend CRUD jalan, skema belum lengkap; FE 0%)

> Bagian-bagian di file ini **sudah diimplementasikan sebagian di backend** per audit 2026-08-12, re-verifikasi skema 2026-08-14. Ref: [00-overview.md](00-overview.md) untuk prinsip & roadmap keseluruhan.
>
> 🚫 **Frontend: BELUM ADA SAMA SEKALI (0%).** `frontend/tenant/src/views/modules/payroll/Payroll.vue` hanya berisi placeholder 1 baris ("Payroll Module — Coming soon"). Route `/payroll` (`router/index.js:228-231`) dan menu sidebar (`layouts/Sidebar.vue:394,437-438`, digate permission `payroll.view`) sudah terdaftar dan mengarah ke halaman placeholder itu, tapi **tidak ada UI apa pun** untuk salary component, salary structure, payroll period/run, BPJS, PPh21, atau payslip — dan tidak ada payroll API service file di frontend sama sekali.
>
> ⚠️ **Skema belum 100% sesuai desain (re-verifikasi 2026-08-14, update 2026-08-14).** CRUD dasar (create/list/update) memang jalan, tapi ada gap struktural nyata dibanding requirement §4-9 — lihat **§Gap Analysis & Rencana Perbaikan** di bagian bawah file ini. **Sudah diperbaiki (2026-08-14):** validasi `calculation_type` enum di Go + kolom `formula`/`reference_component_id` (migration `115_payroll_formula_engine`). **Masih gap:** `salary_components` tidak punya `company_id`/effective-dating sendiri, assignment scope cuma 2 dari 5 level yang direncanakan (Grade + Employee saja, tanpa Organization/Position/Employment-Type/Company), dan tidak ada kolom `version` di tabel mana pun (hanya effective_from/to tanpa overlap guard; `pph21_settings` bahkan cuma bisa punya 1 row aktif tanpa histori).

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

> 🔶 Tabel aktual: `salary_components` — CRUD lengkap, **tapi kolom belum lengkap**: tidak ada `company_id` (tidak ada scoping per-company di level component), tidak ada `effective_from`/`effective_to` sendiri (validity window baru ada di child table `salary_grade_components`/`salary_employee_components`, bukan di component-nya), `active` diimplementasikan sebagai `status` VARCHAR enum (bukan boolean — gap kosmetik, bukan blocker). Lihat §Gap Analysis di bawah.

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

> ✅ **Selesai (2026-08-14):** kolom `calculation_type` kini divalidasi enum di Go (`FIXED|PERCENTAGE|FORMULA|REFERENCE|MANUAL`, lihat `validateSalaryComponentCalculation` di `service.go`) dan di binding DTO. Kolom `formula` (TEXT) dan `reference_component_id` (self-FK) **sudah ada** di `salary_components` via migration `115_payroll_formula_engine` — jadi FORMULA/PERCENTAGE/REFERENCE sekarang punya tempat menyimpan datanya, dan Formula Engine (`calculator/`) siap meng-parse/mengevaluasinya. Endpoint bantu: `POST /payroll/formula/validate`, `GET /payroll/formula/variables`.

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

> ❌ **Koreksi (re-verifikasi 2026-08-14, saat brainstorming Payroll Run): BUKAN "selesai" — tidak ada CRUD API sama sekali.** `SalaryGradeComponent`/`SalaryEmployeeComponent` di `model.go` hanya GORM struct (tabelnya dibuat migration `006_payroll_structure`), tapi **tidak ada satu pun method di `repository.go`, `service.go`, `dto.go`, `handler.go`, atau `routes.go`** untuk keduanya — grep `GradeComponent`/`EmployeeComponent` di file-file itu nihil. Tidak bisa create/read/update/delete assignment grade atau employee override lewat API mana pun. Yang benar-benar CRUD lengkap hanya `salary_components` (§4) — master data komponennya saja, bukan assignment-nya ke grade/employee. Klaim sebelumnya salah, harus diperbaiki sebelum dianggap prasyarat siap untuk Payroll Run.

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

> ❌ **Aktual (re-verifikasi 2026-08-14): hanya 2 level yang benar-benar ada.** `salary_grade_components` = level **Grade** (FK `grading_id`). `salary_employee_components` = level **Employee** (FK `employee_id` wajib, plus kolom `employment_id`/`position_id`/`grading_id` yang sifatnya kontekstual/denormalisasi pada row employee — bukan scope assignment independen dengan resolusi prioritas sendiri). **Tidak ada** tabel/kolom untuk level Organization, Position (sebagai scope berdiri sendiri), Employment Type, atau Company/default-policy. Rantai prioritas 5-level di atas **tidak punya skema pendukung sama sekali** — lihat §Gap Analysis di bawah untuk rencana perbaikan.

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
payroll_run_employees          -- ✅ DIISI oleh CalculatePayrollRun (2026-08-14, lihat 03)
payroll_run_items              -- ✅ DIISI oleh CalculatePayrollRun (2026-08-14, + kolom audit via migration 116)
payroll_payslips               -- ✅ DIISI oleh GeneratePayslips (2026-08-15, lihat 07): per run employee + total_employer_contribution (migration 118)
payroll_payments               -- ✅ tabel baru migration 118: payment batch per run employee (snapshot rekening bank)
pph21_calculation_logs         -- ✅ DIISI oleh kalkulator PPh21 (2026-08-15, lihat 05): 1 baris per run employee + formula_json
payroll_profile_change_logs    -- audit perubahan profile (employee_payroll_profiles dll.), bukan audit kalkulasi
```

**Tidak ada tabel apa pun** untuk: `payroll_audit_logs` run-level (hanya ada log perubahan config + snapshot + `pph21_calculation_logs` sebagai jejak kalkulasi).

Seluruh primary key menggunakan UUID (`CHAR(36)`), konsisten dengan modul lain di repo ini.

> ❌ **Gap versioning (re-verifikasi 2026-08-14):** tidak ada kolom `version` di tabel mana pun (`salary_components`, `bpjs_settings`, `bpjs_rate_components`, `pph21_settings`, `pph21_ptkp_rates`, `pph21_tax_brackets`). Effective-dating (`effective_start_date`/`effective_end_date`) dipakai sebagai gantinya, tapi under-enforced: unique constraint hanya di `(code, effective_start_date)` — tidak ada exclusion constraint yang mencegah dua row dengan rentang tanggal overlap untuk scope yang sama. Lebih parah lagi, `pph21_settings` punya `UNIQUE(setting_code)` **tanpa effective dating sama sekali** — hanya bisa ada 1 row aktif, tidak ada histori versi PPh21 sama sekali walau regulasi pajak berubah setiap tahun. Lihat §Gap Analysis di bawah.

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
- [ ] Create salary structure. — **koreksi**: tabel `salary_grade_components`/`salary_employee_components` ada, tapi tidak ada CRUD API sama sekali (lihat §7).
- [ ] Create employee salary assignment. — **koreksi**: `salary_employee_components` tidak punya repository/service/handler apa pun.
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

---

## Gap Analysis (re-verifikasi 2026-08-14)

Ringkasan hasil audit skema aktual (`006_payroll_structure.sql` + `model.go`) dibanding requirement §4-9 di atas:

| # | Requirement | Kondisi Aktual | Verdict |
|---|---|---|---|
| 1 | `salary_components` perlu `company_id` + `effective_from/to` sendiri | Tidak ada — component tidak punya kolom tenant-scope maupun validity window sendiri (baru ada di child table) | ❌ GAP |
| 2 | `type` (EARNING/DEDUCTION/EMPLOYER_CONTRIBUTION) | String bebas (`component_type VARCHAR(255)`, tanpa CHECK) — bisa menampung semua kode yang direncanakan | ✅ OK |
| 3 | `calculation_type` (FIXED/PERCENTAGE/FORMULA/REFERENCE/MANUAL) | ✅ **SELESAI (2026-08-14):** validasi enum di service (`validateSalaryComponentCalculation`) + binding DTO; kolom `formula` + `reference_component_id` ditambahkan migration `115_payroll_formula_engine`; Formula Engine `calculator/` meng-parse/mengevaluasi keduanya | ✅ OK |
| 4 | Assignment scope 5-level (Company/Organization/Position/Grade/Employment Type/Employee) dengan prioritas resolusi | Hanya 2 level nyata: Grade (`salary_grade_components`) dan Employee (`salary_employee_components`). `position_id`/`employment_id` di tabel employee-level cuma kolom kontekstual, bukan scope assignment independen | ❌ GAP (terbesar) |
| 5 | `version` field untuk resolusi rule historis | Tidak ada kolom `version` di tabel mana pun — hanya `effective_start/end_date` tanpa exclusion constraint anti-overlap. `pph21_settings` malah cuma 1 row aktif (unique constraint tanpa histori sama sekali) | ❌ GAP |

**Kesimpulan**: CRUD dasar (create/list/update) berfungsi dan cukup untuk "component/structure exists in DB", tapi **belum memenuhi arsitektur configuration-driven payroll engine** yang jadi prinsip inti modul ini (lihat [00-overview.md](00-overview.md) §Prinsip Utama poin 1-2). Formula/reference calculation type dan multi-level assignment scope adalah prasyarat struktural untuk [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md) — bukan pekerjaan terpisah yang bisa ditunda.

---

## Rencana Perbaikan Gap (Backend)

Urutan disusun agar tidak menghalangi start Formula Engine ([02-formula-engine.md](02-formula-engine.md)) — migration ini jadi prasyarat, sebaiknya dikerjakan sebagai bagian awal sub-project Formula Engine, bukan sub-project terpisah:

- [ ] Migration: tambah `company_id`, `effective_start_date`, `effective_end_date` ke `salary_components`.
- [x] Migration: tambah `formula TEXT NULL` dan `reference_component_id CHAR(36) NULL` (self-FK ke `salary_components`) untuk mendukung calculation_type FORMULA/REFERENCE. — **SELESAI** migration `115_payroll_formula_engine` (postgres+mysql, up/down).
- [x] Tambah validasi `calculation_type` di service layer (Go const/enum: `FIXED|PERCENTAGE|FORMULA|REFERENCE|MANUAL`) — CHECK constraint di DB opsional, minimal validasi di `service.go`. — **SELESAI** `validateSalaryComponentCalculation` + `ValidationError` → HTTP 400.
- [ ] Migration: tambah `organization_id CHAR(36) NULL` dan `employment_type VARCHAR(50) NULL` ke `salary_grade_components` (atau tabel scope baru `salary_scope_components`) untuk mendukung 3 level assignment yang belum ada (Organization, Position sebagai scope berdiri sendiri, Employment Type). Perlu keputusan desain: extend tabel existing vs tabel baru — evaluasi saat brainstorming sub-project ini.
- [ ] Implementasi priority-resolution logic (Employee > Position > Organization > Company > Default) di service layer saat query salary component aktif untuk seorang employee.
- [ ] Migration: tambah `version INT NOT NULL DEFAULT 1` ke `bpjs_rate_components`, `pph21_ptkp_rates`, `pph21_tax_brackets`, `salary_components` — atau alternatif lebih murah: tambah exclusion constraint (Postgres `EXCLUDE USING gist` on daterange, atau app-level overlap check) tanpa kolom version baru. Keputusan desain saat brainstorming.
- [ ] Migration: ubah `pph21_settings` dari `UNIQUE(setting_code)` single-row menjadi effective-dated (tambah `effective_start_date`/`effective_end_date`, ubah unique constraint jadi `(setting_code, effective_start_date)`) agar histori perubahan regulasi PPh21 tahunan bisa disimpan.

> ⚠️ Item-item ini **harus dikerjakan sebelum atau bersamaan** dengan Task 1 di sub-project Formula Engine ([02-formula-engine.md](02-formula-engine.md)), karena Formula Engine butuh kolom `formula`/`reference_component_id` untuk punya sesuatu yang bisa di-parse.

---

## Rencana Implementasi Frontend

> ✅ **FE payroll selesai (2026-08-15)** untuk area inti. Berikut progress per item dari rencana awal:

- [x] Ganti `frontend/tenant/src/views/modules/payroll/Payroll.vue` placeholder dengan layout index (list Payroll Run + create + Calculate + aksi status + dashboard).
- [x] API service file `frontend/tenant/src/services/api.js` (pola modul lain) untuk memanggil endpoint payroll yang sudah ada di backend (48 handler, lihat [00-overview.md](00-overview.md)).
- [x] **Salary Component**: list + form create/edit (CRUD `salary_components`) — `SalaryComponentsView.vue`.
- [ ] **Salary Structure**: 2 view — grade-level default component (`salary_grade_components`) dan employee-level override (`salary_employee_components`), dengan tabel effective-dating yang jelas ke user. *(belum dibangun — backend sudah ada)*
- [x] **Payroll Period**: list + form create/edit, toggle OPEN/CLOSED — `PayrollPeriodsView.vue`.
- [x] **BPJS Settings & Rate Components**: form konfigurasi tarif per program (HEALTH/JHT/JP/JKK/JKM/JKP) dengan effective dating — `BpjsSettingsView.vue` (+ endpoint baru `GET /bpjs/rate-components`).
- [x] **PPh21 Settings, PTKP Rates, Tax Brackets**: form konfigurasi pajak — `Pph21SettingsView.vue`.
- [x] **Employee Payroll/Bank/BPJS/Tax Profiles**: form per-employee — `PayrollProfilesView.vue` (list + create + delete; update endpoint belum ada di backend).
- [x] **Payroll Run**: list + detail + status-transition UI (DRAFT→CALCULATED→APPROVED→LOCKED), termasuk tombol trigger approval (`GET /runs/:id/approval`) — `Payroll.vue` + `PayrollRunDetail.vue` (tab Employees/Items/Payslips/Payments/Reports).
- [x] i18n: string payroll ditambahkan ke `frontend/tenant/src/locales/en.json`/`id.json`.
- [ ] Update `docs/openapi-report.md` referensi kalau ada penyesuaian kontrak API saat FE dibangun.

> ⚠️ Sisa pekerjaan FE yang masih terbuka: **Salary Structure** (grade-level/employee-level override), **payslip PDF**, **employee portal** (self-service payslip), **payroll journal**, **payment reconciliation**.
