# Payroll — Payroll Run, Calculation Order & Snapshot (✅ Selesai — 2026-08-14)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md). Hasil snapshot (`payroll_run_employees`/`payroll_run_items`) siap dikonsumsi [04-bpjs-engine.md](04-bpjs-engine.md) dan [05-pph21-engine.md](05-pph21-engine.md).

## 10. Payroll Period

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

## 11. Payroll Run

Payroll Run digunakan untuk menjalankan calculation terhadap employee dalam satu period.

> ✅ **Status: implementasi kalkulasi selesai (2026-08-14).** `Service.CalculatePayrollRun` (file baru `calculation.go`) sekarang benar-benar mengisi `payroll_run_employees` + `payroll_run_items`: pilih employee (dari daftar pre-selected atau auto-select dari employee payroll profile aktif), resolve salary structure (grade default + override employee + adjustment periode), evaluasi formula/reference via Formula Engine ([02-formula-engine.md](02-formula-engine.md)), hitung agregat (earning/deduction/employer contribution/net/company cost), dan simpan snapshot dalam transaksi. Transisi `DRAFT → CALCULATED` di `UpdatePayrollRunStatus` kini menjalankan kalkulasi sungguhan — gap paling kritis yang tercatat di bawah sudah ditutup. Endpoint baru: `POST /runs/:id/calculate`, `GET /runs/:id/employees`, `GET /runs/:id/items`.

Flow (desain target — langkah hitung kini berjalan):

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

> ✅ **Aktual (2026-08-14):** flow di atas kini berjalan untuk langkah struktur/komponen. `DRAFT → CALCULATED` memanggil `CalculatePayrollRun` (isi snapshot) lalu lanjut ke approval bila ApprovalEngine di-set (status `CALCULATED`), atau langsung `REVIEWED` bila tanpa approval engine:
>
> ```text
> DRAFT → [CalculatePayrollRun: isi payroll_run_employees + payroll_run_items]
>        → CALCULATED (approval) / REVIEWED (tanpa approval)
> CALCULATED → APPROVED
> APPROVED → LOCKED
> ```
>
> **Update 2026-08-16:** integrasi Approval & Notifikasi dilengkapi: (1) `flow_id` **auto-resolve** — saat `DRAFT → CALCULATED` tanpa `flow_id` eksplisit, service memanggil `ApprovalEngine.GetActiveFlowIDForModule("payroll")` (pola sama KPI/requisition), jadi run otomatis masuk alur approval selama ada flow aktif untuk modul payroll; (2) **notifikasi hasil** — saat instance approval mencapai final, `PAYROLL_APPROVED` / `PAYROLL_REJECTED` dikirim ke pembuat run (`notifyRunOutcome`, best-effort), dengan entri katalog bilingual di `notification/i18n.go`; (3) UI payroll (`Payroll.vue` + `PayrollRunDetail.vue`) kini memakai `PUT /status {status: CALCULATED}` untuk tombol Hitung, menampilkan badge "Menunggu Persetujuan" saat `approval_instance_id` ter-set, dan tombol "Periksa Status Persetujuan" yang memanggil `GET /runs/:id/approval`.
>
> **Update 2026-08-14/15:** BPJS kontribusi ([04-bpjs-engine.md](04-bpjs-engine.md)) dan PPh21 ([05-pph21-engine.md](05-pph21-engine.md)) sudah diintegrasikan ke kalkulasi run sebagai item `source_group='STATUTORY'` — BPJS dihitung setelah struktur (dasar upah `is_bpjs_base` dibaca dari nilai terhitung), PPh21 setelah BPJS (iuran BPJS employee dipakai sebagai pengurang) plus baris `pph21_calculation_logs`. Input Attendance/Leave/Overtime juga sudah masuk sebagai variabel built-in formula via [06-proration-workforce.md](06-proration-workforce.md) (`loadWorkforceSummary` → WORKING_DAYS/WORKED_DAYS/ABSENCE_DAYS/UNPAID_LEAVE_DAYS/OVERTIME_HOURS), dan prorasi join/resign configurable via `payroll_runs.proration_method`. Snapshot struktur yang sudah berjalan tidak berubah.

Integrasi Approval sudah berfungsi nyata — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §26.

---

## 12. Calculation Order

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-14).** `evaluateComponents` di `calculation.go` menjalankan urutan perhitungan: komponen FIXED/MANUAL dihitung lebih dulu (urut `display_order`), lalu FORMULA/PERCENTAGE/REFERENCE di-resolve secara iteratif dengan dependency resolver (baca nilai komponen lain yang sudah terhitung; siklus dependensi ditolak via `DetectCycles`). Agregat (GROSS = total earning, NET = gross − deduction, EMPLOYER_TOTAL_COST = gross + employer contribution) dihitung setelah seluruh item tersedia. Urutan 1-14 di bawah adalah panduan konseptual; eksekusi aktual berbasis dependency + display_order.

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

## 13. Payroll Snapshot

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-14).** `persistRunSnapshot` di `calculation.go` menulis `payroll_run_employees` + `payroll_run_items` dalam satu transaksi (hapus snapshot lama dulu → isi ulang, sehingga aman untuk recalculation). Kolom snapshot tambahan (`calculation_type`, `base_amount`, `rate`, `formula`, `formula_result`) ditambahkan via migration `116_payroll_run_snapshot` agar payroll item menyimpan jejak audit yang lengkap — lihat daftar di bawah. `calculated_amount` = kolom `amount`, `employee_amount`/`employer_amount` dibedakan via `paid_by` yang sudah ada.

Payroll yang sudah dihitung harus menyimpan snapshot.

Jangan bergantung pada konfigurasi live.

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

## Phase 2 — Payroll Engine (checklist, bagian Run/Snapshot)

> ✅ **Payroll Run execution + snapshot selesai (2026-08-14).** Item Formula Engine murni ada di [02-formula-engine.md](02-formula-engine.md).

- [x] Build proration engine. — `calculator/proration.go` (CALENDAR_DAYS/WORKING_DAYS/FIXED_30_DAYS/ATTENDANCE_DAYS) + faktor join tengah bulan di `calculateEmployee`. Integrasi penuh dengan workforce (attendance/leave/overtime) masih menunggu [06-proration-workforce.md](06-proration-workforce.md) §19.
- [x] Build calculation snapshot. — `payroll_run_employees`/`payroll_run_items` kini diisi oleh `CalculatePayrollRun`; kolom audit ditambah migration `116_payroll_run_snapshot`.
- [x] Build calculation audit. — snapshot item menyimpan `calculation_type`, `base_amount`, `formula`, `formula_result` sehingga kalkulasi bisa dijelaskan kembali; `pph21_calculation_logs` tetap dikhususkan untuk jejak PPh21 di [05-pph21-engine.md](05-pph21-engine.md).
- [x] Build recalculation mechanism. — `CalculatePayrollRun` bisa dipanggil ulang pada status DRAFT/CALCULATED; snapshot lama dihapus & diganti bersih (dicegah pada status REVIEWED/APPROVED/LOCKED/CANCELLED).
