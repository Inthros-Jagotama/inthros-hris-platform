# Payroll — PPh 21 Engine (🔶 Prioritas #4)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Bergantung pada [02-formula-engine.md](02-formula-engine.md) dan [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md). Skema config sudah selesai — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §34.

## 18. PPh 21

> 🔶 **Status: konfigurasi ✅ ada (CRUD lengkap), kalkulator pajak ❌ belum ada.** Tabel aktual: `pph21_settings` + `pph21_ptkp_rates` + `pph21_tax_brackets` (bukan `tax_profiles`/`tax_rules`/`tax_rule_brackets` generik). Tabel `pph21_calculation_logs` (migration `007_payroll_run.sql`) sudah menyediakan kolom lengkap untuk jejak kalkulasi (`gross_monthly`, `ptkp_annual`, `pkp_annual`, `pph21_monthly`, `formula_json`, dll.) — tapi **tidak ada satu baris pun yang pernah ditulis ke tabel ini**, karena tidak ada kode kalkulator PPh21 yang berjalan. Sama seperti BPJS: config-nya siap dan effective-dated, konsumennya (payroll calculation) belum dibangun.

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

> 🔶 **Konfigurasi selesai (nama tabel beda dari rencana), kalkulator belum.** Tabel aktual: `pph21_settings`+`pph21_ptkp_rates`+`pph21_tax_brackets`, bukan `tax_profiles`/`tax_rules`/`tax_rule_brackets` generik.

- [x] Tax profile. — via `employee_tax_profiles` (bukan `tax_profiles` generik).
- [ ] Tax status. — belum diverifikasi ulang di audit ini apakah PTKP status (TK/0, K/1, dst.) sudah ada di `employee_tax_profiles`.
- [ ] Taxable component mapping. — tidak ditemukan mapping "component mana yang taxable" secara eksplisit.
- [ ] PPh 21 rule engine. — **belum ada**, ini bagian dari [02-formula-engine.md](02-formula-engine.md).
- [ ] Tax calculation. — **belum ada** — `pph21_calculation_logs` tabel kosong, tidak pernah ditulisi.
- [ ] Tax adjustment.
- [ ] Tax report. — belum ada (bergantung kalkulasi). Lihat juga [08-reporting-testing.md](08-reporting-testing.md).

PPh 21 harus menjadi engine terpisah dari BPJS agar perubahan regulasi dapat dikelola secara independen.
