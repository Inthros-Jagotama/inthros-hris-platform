# Payroll — Payroll Run, Calculation Order & Snapshot (❌ Prioritas #2)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md). Setelah ini selesai, konsumsi hasilnya di [04-bpjs-engine.md](04-bpjs-engine.md) dan [05-pph21-engine.md](05-pph21-engine.md).

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

Flow (desain target):

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
> Transisi `DRAFT → CALCULATED` hanya mengubah kolom `status` — **tidak ada kode yang mengisi `payroll_run_employees`/`payroll_run_items`, menghitung BPJS/PPh21, atau memverifikasi bahwa perhitungan benar-benar terjadi**. Siapa pun bisa memanggil endpoint status-update ini tanpa data payroll run pernah dihitung. Ini adalah gap implementasi paling kritis di seluruh modul.

Integrasi Approval sudah berfungsi nyata — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §26.

---

## 12. Calculation Order

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

## 13. Payroll Snapshot

> ❌ **Status: BELUM DIIMPLEMENTASIKAN.** Tabel `payroll_run_employees` dan `payroll_run_items` sudah ada di migration `007_payroll_run.sql` dengan kolom yang mendekati desain di bawah (component_id, calculation_type, base_amount, dll.), tapi **tidak pernah diisi oleh kode apa pun** — tidak ada service method yang melakukan insert ke kedua tabel ini. Snapshot secara konsep sudah disiapkan skemanya, tapi belum ada satu baris data pun yang pernah dibuat karena §12 (calculation) belum berjalan.

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

> Item Formula Engine murni ada di [02-formula-engine.md](02-formula-engine.md). Berikut item yang menjadi tanggung jawab Payroll Run execution:

- [ ] Build proration engine. — lihat juga [06-proration-workforce.md](06-proration-workforce.md) §19.
- [ ] Build calculation snapshot. — skema tabel (`payroll_run_employees`/`payroll_run_items`) sudah ada, tinggal diisi.
- [ ] Build calculation audit. — skema tabel (`pph21_calculation_logs`) sudah ada, tinggal diisi.
- [ ] Build recalculation mechanism.
