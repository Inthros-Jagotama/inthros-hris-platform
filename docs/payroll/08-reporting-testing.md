# Payroll — Reporting, Dashboard, Audit Trail & Testing (✅ Prioritas #9 — DIIMPLEMENTASIKAN 2026-08-15)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada seluruh calculation engine ([02](02-formula-engine.md)-[07](07-payslip-payment.md)) sudah berjalan sebelum reporting punya data nyata untuk ditampilkan.

## 29. Reporting

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15).** `report.go` menyediakan laporan read-only dari snapshot run (tidak ada tulis-balik): `GetPayrollSummaryReport` (ringkasan run), `GetPayrollDetailReport` (per employee per komponen: base/rate/amount/kategori), `GetBpjsReport` (wage basis, kontribusi employee/employer, total + nomor kepesertaan dari `employee_bpjs_profiles`), `GetTaxReport` (dari `pph21_calculation_logs`: taxable income + PPh21), `GetBankTransferReport` (dari `payroll_payments`), dan `GetPayrollDashboard` (agregat run). Endpoint: `GET /runs/:id/reports/{summary,detail,bpjs,tax,bank}` + `GET /runs/:id/dashboard`.

### Payroll Summary

```text
Total Employee
Gross Salary
Employee Deduction
Employer Contribution
Net Salary
Total Payroll Cost
```

### BPJS Report

```text
Employee
BPJS Number
Wage Basis
Employee Contribution
Employer Contribution
Total Contribution
```

### Tax Report

```text
Employee
Taxable Income
PPh 21
```

### Bank Transfer

```text
Employee
Bank
Account
Amount
```

### Payroll Detail

```text
Employee
Component
Base
Rate
Employee Amount
Employer Amount
```

---

## 30. Dashboard

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15).** `GetPayrollDashboard` mengembalikan agregat run (total employee, gross, deduction, employer contribution, net, total employer cost + status) — data tersedia langsung di `payroll_runs` yang diisi `CalculatePayrollRun`. Frontend dashboard ada di tab Overview `PayrollRunDetail.vue` (frontend payroll selesai 2026-08-15).

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

## 31. Audit Trail

> 🔶 **Status: audit config-level ✅ ada, audit kalkulasi sebagian ✅ (update 2026-08-15).** `salary_change_logs` dan `payroll_profile_change_logs` mengaudit **perubahan konfigurasi**; audit kalkulasi kini ada lewat snapshot item (`payroll_run_items` menyimpan `calculation_type`/`base_amount`/`rate`/`formula`/`formula_result`) dan `pph21_calculation_logs` yang kini **terisi** per employee per run (lihat [05-pph21-engine.md](05-pph21-engine.md)). Yang masih belum ada: `payroll_audit_logs` run-level generik seperti di bawah — sejauh ini jejak run-level cukup dari `payroll_runs.*_at` + snapshot + log PPh21.

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

## Phase 8 — Reporting (checklist, bagian reporting)

> Bagian payslip generator dari Phase 8 ada di [07-payslip-payment.md](07-payslip-payment.md).

- [x] Payroll summary. — `GET /runs/:id/reports/summary` (total employee, gross, deduction, employer contribution, net, company cost).
- [x] Payroll detail. — `GET /runs/:id/reports/detail` (per employee per komponen: component, base, rate, amount, kategori).
- [x] BPJS report. — `GET /runs/:id/reports/bpjs` (wage basis, employee/employer/total contribution + nomor BPJS).
- [x] Tax report. — `GET /runs/:id/reports/tax` (taxable income + PPh21 dari log kalkulasi).
- [x] Bank transfer report. — `GET /runs/:id/reports/bank` (employee, bank, rekening, amount, status — dari payment batch).
- [ ] Payroll journal. — belum ada (jurnal akuntansi/posting ke Accounting module — di luar scope sub-plan ini).

---

## 37. Testing Strategy

> 🔶 **Status aktual (update 2026-08-15):** test CRUD/repository sudah lama ada; sejak kalkulator dibangun, test kalkulasi juga sudah ada — `calculation_test.go` (struktur/formula/reference/adjustment/proration/recalculation/status-transition), `bpjs_test.go`, `pph21_test.go`, `workforce_test.go` (resign mid-period, metode WORKING_DAYS/ATTENDANCE_DAYS, absence/overtime/unpaid-leave via formula) + `calculator/engine_test.go` + `calculator/proration_test.go`. Checklist unit & scenario di bawah yang sudah tertutup ditandai ✅; yang masih **belum ada**: golden dataset regression test dan scenario end-to-end terpisah (menunggu [08-reporting-testing.md](08-reporting-testing.md) itu sendiri).

Payroll harus memiliki test lebih ketat dibanding module biasa.

> ✅ **Golden dataset regression test ada** — `TestGoldenDatasetRegression` di `report_test.go`: input tetap (BASIC 10jt, BPJS 1%/4%/2%, PPh21 TK/0) harus selalu menghasilkan gross 10jt, BPJS employee 300rb, PPh21 240rb, deduction 540rb, net 9.460.000, employer 400rb, company cost 10.4jt — diverifikasi lewat snapshot run & payslip.

### Unit Test

- [x] Basic salary. — `TestCalculatePayrollRun_BasicStructure`.
- [x] Allowance. — `TestCalculatePayrollRun_BasicStructure` (TRANSPORT override).
- [x] Deduction. — `TestWorkforceAbsenceDeduction`/`TestWorkforceUnpaidLeave`.
- [x] Percentage. — `calculator/engine_test.go` (postfix `%`).
- [x] Formula. — `TestCalculatePayrollRun_FormulaComponent` + engine tests.
- [x] Proration. — `calculator/proration_test.go` + `TestProration*` (join, resign, WORKING_DAYS, ATTENDANCE_DAYS, FIXED_30_DAYS).
- [x] Overtime. — `TestWorkforceOvertimeEarning`.
- [x] BPJS. — `bpjs_test.go` (wage basis, caps, JKK risk class, fixed amount, no-profile).
- [x] Tax. — `pph21_test.go` (dasar, progresif, non-NPWP, log persist).
- [ ] Bonus. — belum ada test khusus (komponen FIXED/adjustment sudah dicakup `TestCalculatePayrollRun_Adjustment`).
- [x] Net salary. — semua test kalkulasi (net = earning − deduction).
- [x] Employer cost. — `TestBpjsHealthContributions` (company cost).

### Scenario Test

- [x] Employee baru. — `TestCalculatePayrollRun_BasicStructure` (auto-select payroll profile).
- [x] Employee resign. — `TestProrationResignMidPeriod`.
- [x] Join tengah bulan. — `TestCalculatePayrollRun_ProrationJoinMidPeriod` + `TestProrationWorkingDaysMethod`.
- [ ] Promotion. — resolusi position→grading effective-dated sudah dipakai kalkulasi (belum ada scenario test promosi khusus).
- [ ] Salary increase. — `TestCalculatePayrollRun_Recalculation` mengubah grade component (mendekati, tapi bukan scenario promosi/kenaikan terpisah).
- [x] Unpaid leave. — `TestWorkforceUnpaidLeave`.
- [x] Overtime. — `TestWorkforceOvertimeEarning`.
- [ ] Bonus. — belum ada scenario test khusus (adjustment sudah dicakup).
- [x] BPJS. — `TestBpjsHealthContributions` dst.
- [x] Tax. — `TestPph21Basic` dst.
- [x] Manual adjustment. — `TestCalculatePayrollRun_Adjustment`.

### Regression Test

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

> ✅ **SELESAI (2026-08-15):** `TestGoldenDatasetRegression` — input & ekspektasi persis seperti format di atas, diverifikasi terhadap snapshot run + payslip. Perubahan engine yang mengubah angka ini akan mem-fail test.
