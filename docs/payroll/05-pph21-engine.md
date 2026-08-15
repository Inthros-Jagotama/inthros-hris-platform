# Payroll — PPh 21 Engine (🔶 Prioritas #4)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md). Skema config sudah selesai — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §34.

## 18. PPh 21

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-15), disempurnakan (2026-08-15), metode TER (2026-08-15).** Kalkulator PPh21 ada — `backend/internal/modules/payroll/pph21.go` (`calculatePph21`) mendukung **dua metode** (dipilih via `pph21_settings.calculation_method`, default `TER`):
> - **TER (Tarif Efektif Rata-rata, PP 58/2023)** — Januari–November: pajak = **bruto bulanan × tarif TER** langsung (tanpa pengurangan). Tarif dibaca dari tabel `ters` (grup kategori A/B/C + rentang `bruto_min`/`bruto_max`), kategori ditentukan dari status PTKP (`TK/0`, `TK/1`, `K/0` → A; `TK/2`, `TK/3`, `K/1`, `K/2` → B; `K/3` → C). Desember: metode normal (annualized) − **YTD potongan Jan–Nov** (`SumPph21YtdByEmployeeAndYear`).
> - **REGULAR_GROSS_ANNUALIZED** — metode lama: gross taxable (komponen `is_taxable`), biaya jabatan (min(rate% × gross, max bulanan)), iuran BPJS employee yang boleh dikurangkan (flag `deduct_bpjs_health/jht/jp_employee`), **pengurang non-BPJS seperti iuran pensiun (komponen bertanda `is_pph21_deductible`, bisa lebih dari satu — migration 119)**, annualisasi, PTKP (dari tabel `ptkps` via status di `employee_tax_profiles` — satu sumber kebenaran, migration 121), PKP (dibulatkan ke bawah ke `pkp_rounding_unit`), bracket progresif (`pph21_tax_brackets`), multiplier non-NPWP (`non_npwp_multiplier_percent`), lalu pajak bulanan.
>
> Hasilnya item potongan `source_group=STATUTORY` — **komponen wadah ditentukan oleh flag `is_pph21_component` di komponen gaji (sumber kebenaran tunggal; `pph21_component_id` di setting dihapus)** — **plus** baris `pph21_calculation_logs` berisi jejak rinci (`gross_monthly` s/d `pph21_monthly` + `pension_deductible_monthly` + `calculation_method`/`ter_group`/`ter_rate` + `formula_json` dengan breakdown per bracket). Terintegrasi dengan `CalculatePayrollRun` ([03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)) dan dijalankan setelah BPJS (karena butuh iuran BPJS sebagai pengurang). Log lama dihapus saat recalculation.

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
- [x] Tax status. — PTKP status (`TK/0`, `K/1`, dst.) tersimpan di `employee_tax_profiles.ptkp_status`; dipetakan ke nominal tahunan via tabel **`ptkps`** (kolom `code` = status ternormalisasi mis. `TK0`, `ptkp` = nominal tahunan, `group` = kategori TER). **Satu sumber kebenaran** — `pph21_ptkp_rates` (tabel duplikat) dihapus di migration 121. Status tanpa baris ptkps → PTKP 0 + log warning; kategori TER diambil dari `group`, fallback ke mapping hardcode.
- [x] Taxable component mapping. — mapping eksplisit via flag `salary_components.is_taxable`; gross taxable = total nilai komponen `is_taxable=true` (dibaca dari hasil kalkulasi run).
- [x] Deductible component mapping. — komponen pengurang ditandai flag `salary_components.is_pph21_deductible` (mis. iuran pensiun; bisa banyak) + iuran BPJS employee via flag setting; dijumlahkan sebagai pengurang bruto (`pension_deductible_monthly` di log).
- [x] Output component (wadah hasil). — komponen bertanda `salary_components.is_pph21_component` dipakai sebagai baris potongan hasil pajak (migration 119; `pph21_component_id` di `pph21_settings` dihapus).
- [x] PPh 21 rule engine. — `calculatePph21` di `pph21.go` (metode `TER` default atau `REGULAR_GROSS_ANNUALIZED`; tarif TER dari tabel `ters`, tarif bracket & PTKP dari konfigurasi effective-dated, bukan hard-code).
- [x] Tax calculation. — **SELESAI** — `pph21_calculation_logs` kini terisi per employee per run (gross/occ/bpjs-deductible/net/annual/ptkp/pkp/tax/monthly + `formula_json`), unique per run employee, dihapus & ditulis ulang saat recalc.
- [x] Metode TER (PP 58/2023). — **SELESAI (2026-08-15)** — default `calculation_method=TER`; Jan–Nov pakai tarif dari tabel `ters` (seed 3 kategori A/B/C sesuai PER-2/PJ/2024; seed grup C & `ptkps.group` diperbaiki di migration 120); Desember = metode normal − YTD Jan–Nov. Kategori TER dibaca dari kolom `group` tabel `ptkps` (konsolidasi migration 121).
- [ ] Tax adjustment. — belum ada (penyesuaian pajak manual per employee/period di luar kalkulasi otomatis).
- [ ] Tax report. — belum ada (bergantung kalkulasi). Lihat juga [08-reporting-testing.md](08-reporting-testing.md).

PPh 21 harus menjadi engine terpisah dari BPJS agar perubahan regulasi dapat dikelola secara independen.
