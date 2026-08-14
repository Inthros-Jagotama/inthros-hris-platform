# Payroll — Proration & Workforce Integration (✅ Prioritas #5-6 — DIIMPLEMENTASIKAN 2026-08-15)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Dibutuhkan oleh [02-formula-engine.md](02-formula-engine.md)/[03-payroll-run-snapshot.md](03-payroll-run-snapshot.md) untuk kasus join/resign tengah bulan, dan sebagai input kalkulasi (Attendance/Leave/Overtime).

## 19. Proration

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

## 20. Integrasi Attendance

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

## 21. Integrasi Overtime

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

## 22. Integrasi Leave

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

## 23. Integrasi Employee Movement & Career

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

## 24. Integrasi Recruitment

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

## 25. Integrasi Performance

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

## Phase 4 — Workforce Integration (checklist)

> ✅ **Integrasi selesai (2026-08-15).** Payroll kini mengonsumsi HASIL FINAL Workforce Management (bukan menghitung ulang, sesuai prinsip §20-22): `loadWorkforceSummary` di `workforce.go` membaca `attendance_sessions` (hari hadir, status, `overtime_minutes`), `leave_request_details` + `leave_requests` (cuti APPROVED_FINAL, `is_paid`, `day_fraction`), lalu mengisi variabel built-in formula `WORKING_DAYS`/`WORKED_DAYS`/`ABSENCE_DAYS`/`UNPAID_LEAVE_DAYS`/`OVERTIME_HOURS` (read-only, tidak ada tulis-balik ke modul workforce). Prorasi join **dan resign** tengah bulan didukung, metode configurable per run via kolom `payroll_runs.proration_method` (migration 117): CALENDAR_DAYS (default) | WORKING_DAYS | FIXED_30_DAYS | ATTENDANCE_DAYS. Movement/promotion (position→grading effective-dated) sudah didukung sejak [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md).

- [x] Attendance integration. — `attendance_sessions` dibaca per employee per periode → `WORKED_DAYS`/`ABSENCE_DAYS` (status ABSENT/DAY_OFF/EXEMPT/LEAVE tidak dihitung hadir; weekday tanpa session & tanpa cuti = alpa).
- [x] Leave integration. — `leave_request_details` (join `leave_requests` status APPROVED_FINAL) → `UNPAID_LEAVE_DAYS`/`PAID_LEAVE_DAYS` via `is_paid` + `day_fraction`.
- [x] Overtime integration. — `attendance_sessions.overtime_minutes` → `OVERTIME_HOURS` (contoh formula `OVERTIME_HOURS * OVERTIME_RATE`).
- [ ] Shift integration. — hari kerja dihitung sebagai weekday (Senin-Jumat) default; integrasi `attendance_employee_shifts` (days_of_week_mask) belum — catatan di plan §19.
- [x] Employee movement integration. — employment effective-dated (join/promosi/transfer → position→grading) sudah dipakai sejak [03](03-payroll-run-snapshot.md); kini overlap-periode (`FindEmploymentByEmployeeIDForPeriod`) sehingga resign tengah bulan tetap dihitung terprorasi.
- [x] Absence deduction. — formula komponen bisa memakai `ABSENCE_DAYS` (contoh `BASIC / WORKING_DAYS * ABSENCE_DAYS`).
- [x] Overtime earning. — formula komponen bisa memakai `OVERTIME_HOURS` (contoh `OVERTIME_HOURS * 15000`).
