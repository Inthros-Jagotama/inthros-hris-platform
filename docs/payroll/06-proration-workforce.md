# Payroll — Proration & Workforce Integration (❌ Prioritas #5-6)

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

> ❌ **Belum ada integrasi apa pun** — payroll module tidak memanggil Attendance/Leave/Overtime module manapun (tidak ditemukan di survei kode).

- [ ] Attendance integration.
- [ ] Leave integration.
- [ ] Overtime integration.
- [ ] Shift integration.
- [ ] Employee movement integration.
- [ ] Absence deduction.
- [ ] Overtime earning.
