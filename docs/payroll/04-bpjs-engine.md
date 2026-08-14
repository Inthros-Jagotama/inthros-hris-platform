# Payroll — BPJS Engine (🔶 Prioritas #3)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md). Skema config sudah selesai — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §34.

## 14. BPJS Configuration

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-14).** Kalkulator kontribusi BPJS sekarang ada — `backend/internal/modules/payroll/bpjs.go` (`calculateBpjsContributions`) membaca `bpjs_settings` + `bpjs_rate_components` yang ACTIVE & effective-dated pada tanggal periode plus profil BPJS employee (`employee_bpjs_profiles`), menghitung kontribusi employee (paid_by=EMPLOYEE → `EMPLOYEE_DEDUCTION`) dan employer (paid_by=EMPLOYER → `EMPLOYER_CONTRIBUTION`), lalu menghasilkan `payroll_run_items` dengan `source_group=STATUTORY` yang terintegrasi dengan `CalculatePayrollRun` ([03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)). Tarif tetap tidak di-hard-code: `rate_percent`/`fixed_amount`/`min_max_base_amount`/`jkk_risk_class` semuanya dari konfigurasi. Dasar upah = total komponen `is_bpjs_base`; cap program (kesehatan/pensiun) & cap per-rate diterapkan; rate JKK khusus risk class hanya dipakai employee dengan class yang sama; employee tanpa profil BPJS dilewati. Repository query effective-dated: `FindActiveBpjsSettingByDate`, `FindActiveBpjsRateComponentsBySettingID`, `FindActiveEmployeeBpjsProfileByEmployeeID`.

Tabel rencana awal (tidak dipakai, diganti — lihat §34 di [01-master-data-selesai.md](01-master-data-selesai.md)):

```text
bpjs_programs
bpjs_rules
```

Program:

```text
BPJS_KES
JHT
JKK
JKM
JP
JKP
```

Rule minimal:

```text
program_id
employee_percentage
employer_percentage
calculation_basis
minimum_wage
maximum_wage
fixed_amount
effective_from
effective_to
```

Tarif tidak boleh ditanam langsung dalam source code.

Default rule harus dapat diperbarui berdasarkan regulasi terbaru.

---

## 15. BPJS Kesehatan

Konfigurasi mendukung:

```text
Employee Contribution
Employer Contribution
Wage Basis
Minimum Wage
Maximum Wage
Effective Date
```

Contoh konfigurasi awal:

```text
Employee : 1%
Employer : 4%
Total    : 5%
```

Catatan: nilai tersebut merupakan default konfigurasi dan harus dapat diperbarui jika regulasi berubah.

---

## 16. BPJS Ketenagakerjaan

Program:

```text
JHT
JKK
JKM
JP
JKP
```

Contoh struktur rule:

```text
JHT
├── Employee
└── Employer

JKK
└── Employer

JKM
└── Employer

JP
├── Employee
└── Employer

JKP
└── Rule khusus
```

Default rule awal dapat di-seed berdasarkan ketentuan BPJS yang berlaku, tetapi tetap harus dikelola melalui konfigurasi dan effective dating.

---

## 17. JKK Risk Level

JKK membutuhkan konfigurasi risk level.

Tabel (rencana awal — **tidak ada** di skema aktual):

```text
bpjs_jkk_risk_levels
```

Contoh:

```text
VERY_LOW
LOW
MEDIUM
HIGH
VERY_HIGH
```

Relasi:

```text
Company
   ↓
Business Sector
   ↓
Risk Level
   ↓
JKK Rate
```

Tarif JKK tidak boleh di-hard-code.

> ⚠️ **Aktual:** JKK hanya satu rate flat di `bpjs_rate_components`, tidak ada relasi ke business sector/risk level seperti direncanakan di atas.

---

## Phase 3 — BPJS (checklist)

> ✅ **Kalkulator sudah diimplementasikan (2026-08-14).** `calculateBpjsContributions` di `bpjs.go` menghasilkan item BPJS per employee; item yang di bawah yang bercentang ✅ berarti sudah berfungsi penuh (konfigurasi + kalkulasi).

- [x] BPJS Kesehatan. — kolom `bpjs_program='HEALTH'` di `bpjs_rate_components`, config CRUD ada.
- [x] JHT. — `bpjs_program='JHT'`, config CRUD ada.
- [x] JKK. — `bpjs_program='JKK'`, config CRUD ada.
- [x] JKM. — `bpjs_program='JKM'`, config CRUD ada.
- [x] JP. — `bpjs_program='JP'`, config CRUD ada.
- [x] JKP. — `bpjs_program='JKP'`, config CRUD ada.
- [x] Employee contribution. — dihitung di `calculateBpjsContributions` (paid_by=EMPLOYEE → `EMPLOYEE_DEDUCTION`, mengurangi net pay).
- [x] Employer contribution. — dihitung (paid_by=EMPLOYER → `EMPLOYER_CONTRIBUTION`, menambah biaya perusahaan).
- [x] JKK risk level. — matching via kolom `jkk_risk_class` di `bpjs_rate_components` vs `employee_bpjs_profiles.jkk_risk_class` (rate khusus risk class hanya untuk employee dengan class yang sama; rate tanpa risk class berlaku untuk semua). Tabel `bpjs_jkk_risk_levels`/relasi business sector tetap tidak ada — keputusan desain: risk class disimpan langsung di profil employee.
- [x] Effective-dated BPJS rules. — `bpjs_rate_components` sudah effective-dated; kalkulator memilih rate yang berlaku pada tanggal periode.
- [ ] BPJS report. — belum ada (bergantung kalkulasi). Lihat juga [08-reporting-testing.md](08-reporting-testing.md).
