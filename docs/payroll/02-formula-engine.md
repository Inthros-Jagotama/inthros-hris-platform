# Payroll — Formula Engine (✅ Selesai — 2026-08-14)

> Ref: [00-overview.md](00-overview.md) §Roadmap Prioritas. Ini adalah **blocker untuk semua kalkulasi lain** (BPJS, PPh21, Payroll Run) — lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md), [04-bpjs-engine.md](04-bpjs-engine.md), [05-pph21-engine.md](05-pph21-engine.md).

## 6. Formula Engine

> ✅ **Status: DIIMPLEMENTASIKAN (2026-08-14).** Package `backend/internal/modules/payroll/calculator/` sekarang berisi engine lengkap: lexer, parser recursive-descent, evaluator AST, variable registry built-in, dependency resolver (ekstraksi variabel yang direferensikan), circular-dependency detection, dan rounding engine. Dipakai service layer untuk validasi `calculation_type` FORMULA/PERCENTAGE/REFERENCE saat create/update salary component, plus endpoint `POST /payroll/formula/validate` dan `GET /payroll/formula/variables`. Prasyarat struktural (kolom `formula` + `reference_component_id` di `salary_components`, validasi enum `calculation_type`) ikut selesai via migration `115_payroll_formula_engine` — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §5 & Gap Analysis.

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

- [x] Arithmetic operation. — `+ - * /` dengan precedence & associativity standar.
- [x] Percentage. — postfix `%` (`BPJS_WAGE * 2%` = `BPJS_WAGE * 0.02`).
- [x] Parentheses. — grouping `( ... )`.
- [x] Component reference. — variabel non-built-in diperlakukan sebagai kode komponen lain.
- [x] Variable reference. — variabel built-in registry (GROSS, BPJS_WAGE, OVERTIME_HOURS, dst.).
- [ ] Conditional rule. — belum ada sintaks IF/ternary; bisa diekspresikan via komponen/atribut lain (belum diimplementasikan).
- [x] Rounding. — `calculator/rounding.go` (ROUND/CEIL/FLOOR/NONE + unit rounding).
- [x] Dependency validation. — `ValidateReferences` memvalidasi variabel vs registry + komponen tersedia.
- [x] Circular dependency detection. — `DetectCycles` (self-reference & siklus tidak langsung).
- [x] Formula validation sebelum payroll dijalankan. — validasi saat create/update salary component; validasi dependency saat payroll run memakai engine ini (lihat [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md)).

---

## Phase 2 — Payroll Engine (checklist, bagian Formula Engine)

> ✅ **Formula Engine selesai (2026-08-14).** Item snapshot/audit selengkapnya ada di [03-payroll-run-snapshot.md](03-payroll-run-snapshot.md), item ini fokus ke Formula Engine itu sendiri.

- [x] Build Formula Engine. — `calculator/` (lexer.go, parser.go, evaluator.go, engine.go).
- [x] Build variable registry. — `calculator/registry.go` (`DefaultRegistry`, `VariableMeta`).
- [x] Build component dependency resolver. — `Engine.ReferencedVariables` + `Engine.ValidateReferences` (validasi vs komponen tersedia).
- [x] Detect circular dependency. — `DetectCycles` di `engine.go` (DFS, deteksi self-reference & siklus tidak langsung, deduplikasi rotasi).
- [x] Build calculation context. — evaluator menerima `VariableResolver` (konteks nilai per variabel saat eksekusi payroll run).
- [x] Build rounding engine. — `calculator/rounding.go` (`ROUND`/`CEIL`/`FLOOR`/`NONE` + `RoundToUnit`).

> ✅ Kolom `calculation_type`/`formula`/`reference_component_id` pada `salary_components` kini benar-benar ada di skema (migration `115_payroll_formula_engine`) — lihat [01-master-data-selesai.md](01-master-data-selesai.md) §5 & Gap Analysis.
