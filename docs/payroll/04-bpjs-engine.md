# Payroll — BPJS Engine (🔶 Prioritas #3)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md). Skema config sudah selesai — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §34.

## 14. BPJS Configuration

> 🔶 **Status: konfigurasi ✅ ada (CRUD lengkap), kalkulator kontribusi ❌ belum ada.** Tabel aktual bernama `bpjs_settings` + `bpjs_rate_components` (bukan `bpjs_programs`/`bpjs_rules` seperti di bawah). `bpjs_rate_components` sudah mendukung effective dating, `rate_percent`, `min/max_base_amount`, kolom enum `bpjs_program` (HEALTH/JHT/JP/JKK/JKM/JKP) dan `paid_by` (EMPLOYEE/EMPLOYER) — jadi tarif memang tidak di-hard-code, sesuai prinsip di bawah. Tapi **tidak ada kode yang membaca tabel ini untuk menghasilkan angka kontribusi employee/employer** — CRUD-nya lengkap, konsumennya (payroll calculation) belum ada.

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

> 🔶 **Konfigurasi selesai, kalkulator belum.** Tabel `bpjs_settings`+`bpjs_rate_components` sudah effective-dated dan CRUD lengkap; item di bawah yang bercentang berarti "konfigurasinya ada", bukan "kontribusinya sudah terhitung otomatis" — itu masih menunggu [02-formula-engine.md](02-formula-engine.md) & [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md).

- [x] BPJS Kesehatan. — kolom `bpjs_program='HEALTH'` di `bpjs_rate_components`, config CRUD ada.
- [x] JHT. — `bpjs_program='JHT'`, config CRUD ada.
- [x] JKK. — `bpjs_program='JKK'`, config CRUD ada.
- [x] JKM. — `bpjs_program='JKM'`, config CRUD ada.
- [x] JP. — `bpjs_program='JP'`, config CRUD ada.
- [x] JKP. — `bpjs_program='JKP'`, config CRUD ada.
- [ ] Employee contribution. — **belum dihitung** (tidak ada kalkulator).
- [ ] Employer contribution. — **belum dihitung**.
- [ ] JKK risk level. — **tidak ada** tabel `bpjs_jkk_risk_levels`/relasi ke business sector seperti direncanakan §17; JKK hanya satu rate flat di `bpjs_rate_components`.
- [x] Effective-dated BPJS rules. — `bpjs_rate_components` sudah effective-dated.
- [ ] BPJS report. — belum ada (bergantung kalkulasi). Lihat juga [08-reporting-testing.md](08-reporting-testing.md).
