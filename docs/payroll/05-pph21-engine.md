# Payroll — PPh 21 Engine (🔶 Prioritas #4)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md). Skema config sudah selesai — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §34.

## 18. PPh 21

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15).** Kalkulator PPh21 ada — `backend/internal/modules/payroll/pph21.go` (`calculatePph21`) menghitung PPh21 bulanan metode `REGULAR_GROSS_ANNUALIZED`: gross taxable (komponen `is_taxable`), biaya jabatan (min(rate% × gross, max bulanan)), iuran BPJS employee yang boleh dikurangkan (flag `deduct_bpjs_health/jht/jp_employee`), annualisasi, PTKP (dari `pph21_ptkp_rates` via status di `employee_tax_profiles`), PKP (dibulatkan ke bawah ke `pkp_rounding_unit`), bracket progresif (`pph21_tax_brackets`), multiplier non-NPWP (`non_npwp_multiplier_percent`), lalu pajak bulanan. Hasilnya item potongan `source_group=STATUTORY` (komponen `pph21_component_id` dari setting) **plus** baris `pph21_calculation_logs` berisi jejak rinci (`gross_monthly` s/d `pph21_monthly` + `formula_json` dengan breakdown per bracket) — tabel yang selama ini kosong kini terisi, terintegrasi dengan `CalculatePayrollRun` ([03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)) dan dijalankan setelah BPJS (karena butuh iuran BPJS sebagai pengurang). Log lama dihapus saat recalculation.

PPh 21 dibuat sebagai statutory rule engine terpisah dari BPJS.

Komponen:

```text
Tax Profile
Tax Status
Taxable Income
Tax Rule
Tax Calculation
Tax Adjustment
```

Struktur:

```text
PPh21 Engine
├── Employee Tax Profile
├── Taxable Components
├── Deductible Components
├── Tax Rule
├── Tax Calculation
└── Tax Result
```

Rule harus mendukung effective dating karena regulasi pajak dapat berubah.

---

## Phase 5 — Tax (checklist)

> ✅ **Kalkulator sudah diimplementasikan (2026-08-15).** Item di bawah yang bercentang ✅ berarti sudah berfungsi penuh (konfigurasi + kalkulasi).

- [x] Tax profile. — via `employee_tax_profiles` (bukan `tax_profiles` generik), dibaca effective-dated oleh `FindActiveEmployeeTaxProfileByEmployeeID`.
- [x] Tax status. — PTKP status (`TK/0`, `K/1`, dst.) tersimpan di `employee_tax_profiles.ptkp_status`; dipetakan ke nominal tahunan via `pph21_ptkp_rates` (engine memakai status ini; status tanpa rate → PTKP 0 + log warning).
- [x] Taxable component mapping. — mapping eksplisit via flag `salary_components.is_taxable`; gross taxable = total nilai komponen `is_taxable=true` (dibaca dari hasil kalkulasi run).
- [x] PPh 21 rule engine. — `calculatePph21` di `pph21.go` (metode `REGULAR_GROSS_ANNUALIZED`; tarif bracket & PTKP dari konfigurasi effective-dated, bukan hard-code).
- [x] Tax calculation. — **SELESAI** — `pph21_calculation_logs` kini terisi per employee per run (gross/occ/bpjs-deductible/net/annual/ptkp/pkp/tax/monthly + `formula_json`), unique per run employee, dihapus & ditulis ulang saat recalc.
- [ ] Tax adjustment. — belum ada (penyesuaian pajak manual per employee/period di luar kalkulasi otomatis).
- [ ] Tax report. — belum ada (bergantung kalkulasi). Lihat juga [08-reporting-testing.md](08-reporting-testing.md).

PPh 21 harus menjadi engine terpisah dari BPJS agar perubahan regulasi dapat dikelola secara independen.
