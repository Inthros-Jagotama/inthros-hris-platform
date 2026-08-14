# Payroll — Payslip & Payment (🔶/❌ Prioritas #7-8)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Keduanya bergantung pada calculation engine ([02-formula-engine.md](02-formula-engine.md), [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)) sudah berjalan.

## 28. Payslip

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

## 27. Payment

> ❌ **Status: BELUM DIIMPLEMENTASIKAN — tidak ada tabel, struct, atau kode sama sekali.** Tidak ada `payroll_payments` atau padanannya di skema aktual. `EmployeeBankProfile` (tabel `employee_bank_profiles`) sudah menyimpan rekening bank per employee untuk keperluan payroll, tapi tidak ada kode yang mengonsumsinya untuk menghasilkan payment/disbursement record. Seluruh bagian ini murni desain target, prioritas paling akhir karena bergantung pada calculation engine yang juga belum ada.

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

## Phase 7 — Payment (checklist)

> ❌ **Belum ada satu item pun** — tidak ada tabel maupun kode sama sekali. `employee_bank_profiles` sudah ada tapi belum dikonsumsi.

- [ ] Employee bank account. — data sudah ada di `employee_bank_profiles`, tapi belum dikonsumsi proses payment.
- [ ] Payment batch.
- [ ] Bank transfer file.
- [ ] Payment status.
- [ ] Payment reconciliation.
- [ ] Payment reversal.

## Phase 8 — Payslip (checklist, bagian payslip)

> 🔶 **Tabel payslip ada, generator belum ada.** Bagian reporting dari Phase 8 ada di [08-reporting-testing.md](08-reporting-testing.md).

- [ ] Payslip HTML. — tabel `payroll_payslips` ada, generator tidak ada.
- [ ] Payslip PDF. — tidak ada.
- [ ] Employee portal. — tidak diverifikasi; tidak ditemukan endpoint publik payslip di survei.
