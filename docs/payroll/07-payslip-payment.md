# Payroll — Payslip & Payment (✅ Prioritas #7-8 — DIIMPLEMENTASIKAN 2026-08-15)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Keduanya bergantung pada calculation engine ([02-formula-engine.md](02-formula-engine.md), [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)) sudah berjalan.

## 28. Payslip

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15).** `payslip.go` menyediakan `GeneratePayslips` (satu payslip per run employee dari snapshot `payroll_run_employees` + `payroll_run_items`, nomor `SLP-<period>-<seq>`, regenerasi aman), `PublishPayslip`/`CancelPayslip` (transisi DRAFT→PUBLISHED→CANCELLED), list/get, dan `GetPayslipHTML` (render server-side dengan rincian earnings/deductions/employer contributions dari snapshot). Endpoint: `POST /runs/:id/payslips`, `GET /runs/:id/payslips`, `GET /payslips/:id`, `GET /payslips/:id/html`, `POST /payslips/:id/publish`, `POST /payslips/:id/cancel`. Kolom `payroll_payslips.total_employer_contribution` ditambahkan (migration 118) untuk bagian "Employer Contribution" pada payslip. PDF: konversi HTML→PDF bisa memanfaatkan engine DOCX LibreOffice di modul documenttemplate atau tooling FE — belum termasuk di sini.

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

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15).** `payment.go` + tabel `payroll_payments` (migration 118) — `CreatePaymentBatch` membuat satu payment per run employee (nominal = net amount, status PENDING) dengan **snapshot rekening** dari `employee_bank_profiles` (bank profile utama aktif pada tanggal periode), sehingga perubahan rekening setelah run final tidak mengubah batch yang sudah dibuat. Employee tanpa bank profile aktif dilewati (dihitung & dilaporkan). Transisi status divalidasi (PENDING→PROCESSING→PAID, FAILED, REVERSED), export CSV bank transfer file tersedia. Endpoint: `POST /runs/:id/payments`, `GET /runs/:id/payments`, `GET /runs/:id/payments/export`, `GET /payments/:id`, `POST /payments/:id/status`.

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

> ✅ **Payment batch sudah diimplementasikan (2026-08-15)** — tabel `payroll_payments` + service `payment.go`.

- [x] Employee bank account. — `employee_bank_profiles` dikonsumsi `CreatePaymentBatch` via `FindActivePrimaryBankProfileByEmployeeID` (snapshot rekening disalin ke `payroll_payments`).
- [x] Payment batch. — satu `payroll_payments` per run employee (net amount, PENDING); regenerasi batch aman (batch lama dihapus).
- [x] Bank transfer file. — `GET /runs/:id/payments/export` menghasilkan CSV (employee, bank, rekening, nominal, status, reference).
- [x] Payment status. — PENDING/PROCESSING/PAID/FAILED/REVERSED dengan transisi divalidasi + timestamp + reference/failed_reason.
- [ ] Payment reconciliation. — belum ada (rekonsiliasi vs file bank balik / jurnal).
- [x] Payment reversal. — transisi PAID→REVERSED (dengan timestamp `reversed_at`).

## Phase 8 — Payslip (checklist, bagian payslip)

> ✅ **Generator payslip sudah ada (2026-08-15)** — `payslip.go` + endpoint di atas. Bagian reporting dari Phase 8 ada di [08-reporting-testing.md](08-reporting-testing.md).

- [x] Payslip HTML. — `GET /payslips/:id/html` (server-render: earnings, deductions, employer contributions, net, company cost).
- [ ] Payslip PDF. — belum (konversi HTML→PDF bisa memakai engine DOCX LibreOffice / tooling FE — follow-up).
- [ ] Employee portal. — endpoint sudah ada (`GET /payslips/:id` + `/html`), UI/portal employee belum dibangun (frontend payroll admin selesai 2026-08-15 — lihat [01-master-data-selesai.md](01-master-data-selesai.md)).
