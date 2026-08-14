# Payroll — Reporting, Dashboard, Audit Trail & Testing (❌ Prioritas #9)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada seluruh calculation engine ([02](02-formula-engine.md)-[07](07-payslip-payment.md)) sudah berjalan sebelum reporting punya data nyata untuk ditampilkan.

## 29. Reporting

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

> 🔶 **Status: audit config-level ✅ ada, audit kalkulasi ❌ belum ada.** `salary_change_logs` (perubahan `salary_employee_components`) dan `payroll_profile_change_logs` (perubahan `employee_payroll_profiles` dll.) sudah ada dan berfungsi — tapi keduanya mengaudit **perubahan konfigurasi**, bukan kalkulasi payroll. Tidak ada `payroll_calculation_logs`/`payroll_audit_logs` seperti di bawah; yang paling dekat adalah tabel `pph21_calculation_logs` yang skemanya sudah siap tapi tidak pernah ditulisi (lihat [05-pph21-engine.md](05-pph21-engine.md)) — begitu calculation engine dibangun, tabel itu perlu diisi sebagai bagian dari audit kalkulasi, bukan tabel generik baru seperti direncanakan di bawah.

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

- [ ] Payroll summary.
- [ ] Payroll detail.
- [ ] BPJS report.
- [ ] Tax report.
- [ ] Bank transfer report.
- [ ] Payroll journal.

---

## 37. Testing Strategy

> 🔶 **Status aktual: 43 test ada (`repository_test.go` 17 + `service_test.go` 26), tapi semuanya level CRUD/repository** — testing untuk unit test kalkulasi (basic salary, BPJS, tax, dsb.) dan scenario test di bawah **belum ada karena kode yang diuji belum ada**. Golden dataset regression test juga belum dibuat.

Payroll harus memiliki test lebih ketat dibanding module biasa.

### Unit Test

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

### Scenario Test

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
