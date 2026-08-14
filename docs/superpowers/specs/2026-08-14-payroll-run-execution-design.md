# Payroll Run Execution + Snapshot — Design

> Ref: `docs/payroll/03-payroll-run-snapshot.md` (priority #2 in the payroll roadmap, `docs/payroll/00-overview.md`). Depends on the already-implemented Formula Engine (`docs/payroll/02-formula-engine.md`, `backend/internal/modules/payroll/calculator/`).

## Goal

Make `UpdatePayrollRunStatus`'s DRAFT→CALCULATED transition actually calculate payroll instead of just flipping a status column. Populate `payroll_run_employees`/`payroll_run_items` with real computed amounts, using the Formula Engine to evaluate `salary_components.formula` for PERCENTAGE/FORMULA/REFERENCE types.

## Explicitly out of scope (deferred to later sub-projects)

- **BPJS/PPh21 actual contribution/tax calculation** — `docs/payroll/04-bpjs-engine.md` / `05-pph21-engine.md` (priority #3/#4). Any salary component whose formula references `TAXABLE_INCOME` or similar statutory-derived variables will simply evaluate those as 0 for now — not a blocker, since no BPJS/PPh21 component is currently assigned via `salary_grade_components`/`salary_employee_components` in this codebase (those config tables are separate: `bpjs_settings`, `pph21_settings`).
- **Attendance/Leave/Overtime/Proration integration** — `docs/payroll/06-proration-workforce.md` (priority #5-6). Workforce-derived variables (`WORKING_DAYS`, `WORKED_DAYS`, `ABSENCE_DAYS`, `UNPAID_LEAVE_DAYS`, `OVERTIME_HOURS`, `OVERTIME_RATE`) default to `0` in this sub-project. No proration of BASIC salary for mid-month joiners/leavers.
- **Full CRUD for `salary_grade_components`/`salary_employee_components`** — only read access needed for calculation is built here (see Prerequisite below). Create/update/delete admin API for these two tables is a separate, later piece of work.
- **Payroll Policy resolution beyond Grade/Employee** — Organization/Position/Employment-Type/Company scope levels don't exist in the schema (see `docs/payroll/01-master-data-selesai.md` §7 gap note) and are not built here.
- **Recalculation history/versioning** — recalculating an already-CALCULATED run (via DRAFT→CALCULATED again) overwrites the previous snapshot rows. No audit trail of prior calculation attempts in this sub-project (see `docs/payroll/08-reporting-testing.md` for future audit trail work).
- **Payslip generation, payment, reporting** — unaffected by this work, remain separate sub-projects.

## Prerequisite: read access to salary structure assignment

`SalaryGradeComponent`/`SalaryEmployeeComponent` currently have GORM structs (`model.go`) but zero repository/service methods (confirmed by grep — no matches in `repository.go`, `service.go`, `dto.go`, `handler.go`). Calculation needs to read this data, so this sub-project adds:

- `Repository.FindActiveSalaryGradeComponents(ctx, gradingID uuid.UUID, asOfDate string) ([]SalaryGradeComponent, error)` — `WHERE grading_id = ? AND status = 'ACTIVE' AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)`.
- `Repository.FindActiveSalaryEmployeeComponents(ctx, employeeID uuid.UUID, asOfDate string) ([]SalaryEmployeeComponent, error)` — same effective-date filter, `WHERE employee_id = ?`.

No handler/route/DTO for these — internal to the service layer for this sub-project. (Admin CRUD for grade/employee component assignment is flagged as a follow-up, not built here.)

## Data Model Changes

None. All target tables (`payroll_run_employees`, `payroll_run_items`, `payroll_runs` aggregate columns) already exist from migration `007_payroll_run.sql`.

## Calculation Flow

Triggered inside `Service.UpdatePayrollRunStatus` when `req.Status == "CALCULATED" && pr.Status == "DRAFT"`, before the existing approval-instance creation logic. Wrapped in a single DB transaction (existing `s.repo` transaction pattern used elsewhere in this module — check `service.go` for the established helper, e.g. `r.getDB(ctx)` / GORM `.Transaction()`).

```text
1. Load PayrollPeriod (for AsOfDate, used as the effective-date filter for all component/profile lookups).
2. Load all EmployeePayrollProfile where IsPayrollActive = true and effective at period.AsOfDate.
3. DELETE existing payroll_run_employees/payroll_run_items for this PayrollRunID (overwrite semantics).
4. For each employee profile:
   a. Load active SalaryGradeComponent rows for the employee's grading (if any) — this seeds default components.
   b. Load active SalaryEmployeeComponent rows for the employee — these override/add to (a) by salary_component_id.
   c. Merge into a single component list (employee-level entry wins over grade-level entry for the same salary_component_id).
   d. If the merged list is empty: insert one PayrollRunEmployee row with Status="EXCLUDED", TotalEarning/Deduction/etc = 0, no PayrollRunItem rows. Continue to next employee — do NOT fail the whole run.
   e. Otherwise, for each component in the merged list, load the SalaryComponent (code, calculation_type, formula, reference_component_id).
   f. Build a dependency graph (component code -> referenced variable names, via Engine.ReferencedVariables on FORMULA/PERCENTAGE/REFERENCE types) and topologically order evaluation; use calculator.DetectCycles first — a cycle aborts the ENTIRE payroll run with a descriptive error naming the run and the cycle path (fail loudly, matches the approval.RoutingError policy already in this file).
   g. Evaluate each component in dependency order via a PayrollContext resolver (see below), producing a float64 amount per component.
   h. Compute aggregate totals per employee: TOTAL_EARNINGS (sum of EARNING-type items), TOTAL_DEDUCTIONS (sum of DEDUCTION-type items), TOTAL_EMPLOYER_CONTRIBUTION (sum of EMPLOYER_CONTRIBUTION-type items), NET_SALARY = TOTAL_EARNINGS - TOTAL_DEDUCTIONS, EMPLOYER_TOTAL_COST = TOTAL_EARNINGS + TOTAL_EMPLOYER_CONTRIBUTION.
   i. Insert one PayrollRunEmployee row (Status="CALCULATED") with these totals, and one PayrollRunItem row per evaluated component (component_code, component_name, component_type, amount, paid_by, affects_gross_pay/net_pay/company_cost flags derived from component_type).
5. Update payroll_runs aggregate columns (TotalEmployees, TotalEarning, TotalDeduction, TotalEmployerContribution, TotalNet, TotalCompanyCost) as sums across all PayrollRunEmployee rows just inserted (including EXCLUDED employees contributing 0).
6. Commit transaction. Only after successful commit does the existing approval-instance-creation logic run (unchanged).
```

## `PayrollContext` (new: `calculator/context.go`)

Implements the existing `VariableResolver` interface (defined in `evaluator.go`). Two-tier resolution:

1. **Built-in aggregate variables** (GROSS, TOTAL_EARNINGS, TOTAL_DEDUCTIONS, TOTAL_EMPLOYEE_DEDUCTION, TOTAL_EMPLOYER_CONTRIBUTION, NET_SALARY, EMPLOYER_TOTAL_COST) — computed from the running sums of already-evaluated components at resolution time (only valid to reference in a component evaluated AFTER all its dependency components, guaranteed by the topological order from step 4f above).
2. **Workforce variables** (WORKING_DAYS, WORKED_DAYS, ABSENCE_DAYS, UNPAID_LEAVE_DAYS, OVERTIME_HOURS, OVERTIME_RATE) — hardcoded to `0` in this sub-project (documented limitation, see Out of Scope).
3. **Component-reference variables** (any other name) — looked up in a `map[string]float64` of already-evaluated component amounts for the current employee, keyed by normalized component code. Unresolved reference (dependency not in the merged component list) is an evaluation error, aborting the run per the fail-loudly policy.

## Error Handling

- Circular dependency in an employee's component formulas → abort entire payroll run (transaction rolled back), error message includes run code + cycle path from `calculator.DetectCycles`.
- Formula references a component code not assigned to the employee (grade or override) → abort entire payroll run, error names the employee and missing component code.
- Employee with zero assigned components → not an error; EXCLUDED row, run continues.
- Formula parse/evaluate error (malformed expression — shouldn't happen since `salary_components` validates on create/update, but defensive) → abort entire payroll run.
- DB failure mid-loop → transaction rollback, `payroll_runs.status` stays DRAFT (no partial CALCULATED state visible).

## Testing Plan

- `calculator/context_test.go`: `PayrollContext` resolves built-in aggregates correctly given a pre-populated component-amount map; resolves workforce variables to 0; errors on unresolved component reference.
- `service_test.go` additions (following existing test patterns in this file):
  - Single employee, single FIXED component → correct PayrollRunEmployee/PayrollRunItem amounts.
  - PERCENTAGE component referencing another component (`JHT_EMP = BPJS_WAGE * 2%`) → correct evaluation order and result.
  - Employee with grade-level default + employee-level override for the same component code → override wins.
  - Employee with zero components → EXCLUDED row, run still succeeds for other employees.
  - Circular dependency between two components → `UpdatePayrollRunStatus` returns error, run stays DRAFT, no rows written.
  - Recalculate (DRAFT→CALCULATED→ back to DRAFT via direct status update → CALCULATED again) → old snapshot rows overwritten, not duplicated.
  - `payroll_runs` aggregate totals match sum of `payroll_run_employees` rows.

## Files Touched

- `backend/internal/modules/payroll/repository.go` — `FindActiveSalaryGradeComponents`, `FindActiveSalaryEmployeeComponents`.
- `backend/internal/modules/payroll/calculator/context.go` (new) + `context_test.go` (new) — `PayrollContext`.
- `backend/internal/modules/payroll/service.go` — `calculatePayrollRun` (new private method), called from `UpdatePayrollRunStatus`; `service_test.go` additions.
- `docs/payroll/03-payroll-run-snapshot.md` — update status blockquotes (§11-13) and Phase 2 checklist once implemented.
- `docs/payroll/01-master-data-selesai.md` — note the new read-only repository methods for grade/employee components (still flag full CRUD as outstanding).
- `docs/database-schema.md` — no schema change, so no update needed here.
