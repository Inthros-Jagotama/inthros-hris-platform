# Payroll — Formula Engine (❌ Prioritas #1)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Ini adalah **blocker untuk semua kalkulasi lain** (BPJS, PPh21, Payroll Run) — lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md), [04-bpjs-engine.md](04-bpjs-engine.md), [05-pph21-engine.md](05-pph21-engine.md).

## 6. Formula Engine

> ❌ **Status: BELUM DIIMPLEMENTASIKAN.** Tidak ada parser, evaluator, variable registry, dependency resolver, atau circular-dependency detection di mana pun dalam `backend/internal/modules/payroll/`. Satu-satunya jejak konsep "formula" adalah kolom `Pph21CalculationLog.FormulaJSON *string` (`model.go:629`) — kolom JSON kosong yang dimaksudkan untuk menyimpan jejak perhitungan, tapi tidak pernah ditulisi oleh kode apa pun. Seluruh bagian ini murni desain target.

Jangan menggunakan `eval()`.

Formula Engine harus memiliki parser dan evaluator sendiri.

### Contoh formula

```text
BASIC + POSITION_ALLOWANCE
```

```text
BPJS_WAGE * 2%
```

```text
OVERTIME_HOURS * OVERTIME_RATE
```

```text
GROSS - TOTAL_EMPLOYEE_DEDUCTION
```

### Variable registry

Contoh variable:

```text
BASIC
GROSS
BPJS_WAGE
TAXABLE_INCOME
WORKING_DAYS
WORKED_DAYS
ABSENCE_DAYS
OVERTIME_HOURS
OVERTIME_RATE
TOTAL_EARNINGS
TOTAL_DEDUCTIONS
```

### Formula Engine harus mendukung

- Arithmetic operation.
- Percentage.
- Parentheses.
- Component reference.
- Variable reference.
- Conditional rule.
- Rounding.
- Dependency validation.
- Circular dependency detection.
- Formula validation sebelum payroll dijalankan.

---

## Phase 2 — Payroll Engine (checklist, bagian Formula Engine)

> ❌ **Belum ada satu item pun** — ini prioritas #1 pekerjaan berikutnya. Item snapshot/audit selengkapnya ada di [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md), item ini fokus ke Formula Engine itu sendiri.

- [ ] Build Formula Engine.
- [ ] Build variable registry.
- [ ] Build component dependency resolver.
- [ ] Detect circular dependency.
- [ ] Build calculation context.
- [ ] Build rounding engine.

Referensi kolom `calculation_type`/`formula` pada `salary_components` sudah ada di skema aktual — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §5.
