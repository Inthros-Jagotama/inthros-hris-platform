package payroll

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/payroll/calculator"
)

// ItemCategory — kategori baris payroll_run_items (sinkron dengan migration 007).
const (
	ItemCategoryEmployeeEarning      = "EMPLOYEE_EARNING"
	ItemCategoryEmployeeDeduction    = "EMPLOYEE_DEDUCTION"
	ItemCategoryEmployerContribution = "EMPLOYER_CONTRIBUTION"
	ItemCategoryInformation          = "INFORMATION"
)

// SourceGroup — asal baris payroll_run_items.
const (
	SourceGroupStructure  = "STRUCTURE"
	SourceGroupAdjustment = "ADJUSTMENT"
	SourceGroupStatutory  = "STATUTORY"
)

// runEmployeeContext mengumpulkan data yang dibutuhkan kalkulasi satu employee.
type runEmployeeContext struct {
	employee    EmployeeRead
	employment  *EmploymentRead
	position    *PositionRead
	grading     *GradingRead
	gradeComps  []SalaryGradeComponent    // default per grade
	empComps    []SalaryEmployeeComponent // override per employee
	adjustments []SalaryEmployeeAdjustment
}

// componentInput adalah komponen yang siap dihitung untuk satu employee:
// gabungan grade default + override employee (override menang) + adjustment.
type componentInput struct {
	component     SalaryComponent
	baseAmount    float64 // dari grade/employee structure
	adjustmentAmt float64 // penyesuaian sekali-jalan periode
	sourceGroup   string
	sourceTable   string
	sourceID      *uuid.UUID
	displayOrder  int
}

// employeeCalcResult adalah hasil kalkulasi satu employee.
type employeeCalcResult struct {
	runEmployee PayrollRunEmployee
	items       []PayrollRunItem
	pph21Log    *Pph21CalculationLog // jejak kalkulasi PPh21 (jika ada)
}

// CalculatePayrollRun menjalankan perhitungan untuk sebuah payroll run dan
// menyimpan snapshot ke payroll_run_employees + payroll_run_items.
//
// Bisa dipanggil berulang (recalculation): snapshot lama dihapus lalu diisi
// ulang. Status run harus DRAFT atau CALCULATED.
func (s *Service) CalculatePayrollRun(ctx context.Context, runID string) (*PayrollRunResponse, error) {
	uid, err := uuid.Parse(runID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	run, err := s.repo.FindPayrollRunByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if run.Status != "DRAFT" && run.Status != "CALCULATED" {
		return nil, fmt.Errorf("payroll run tidak bisa dihitung ulang pada status %s (hanya DRAFT/CALCULATED)", run.Status)
	}

	period, err := s.repo.FindPayrollPeriodByID(ctx, run.PayrollPeriodID)
	if err != nil {
		return nil, err
	}

	// 1. Tentukan daftar employee yang dihitung. Run tanpa employee tetap bisa
	// dihitung (snapshot kosong, total 0) — mis. off-cycle run yang belum diisi.
	employees, err := s.selectRunEmployees(ctx, run, period)
	if err != nil {
		return nil, err
	}
	if len(employees) == 0 {
		s.logger.Warn("Payroll run calculated without employees",
			zap.String("run_id", run.ID.String()),
			zap.String("run_code", run.RunCode),
		)
		if err := s.persistRunSnapshot(ctx, run.ID, nil); err != nil {
			return nil, err
		}
		run.TotalEmployees = 0
		run.TotalEarning = 0
		run.TotalDeduction = 0
		run.TotalEmployerContribution = 0
		run.TotalNet = 0
		run.TotalCompanyCost = 0
		now := time.Now()
		run.CalculatedAt = &now
		if err := s.repo.UpdatePayrollRun(ctx, run); err != nil {
			return nil, err
		}
		response := toPayrollRunResponse(run)
		return &response, nil
	}

	// 2. Validasi & siapkan komponen yang bisa direferensikan formula.
	components, err := s.loadAllActiveComponents(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]employeeCalcResult, 0, len(employees))
	for _, emp := range employees {
		res, err := s.calculateEmployee(ctx, run, period, emp, components)
		if err != nil {
			return nil, fmt.Errorf("kalkulasi gagal untuk %s (%s): %w", emp.Name, emp.EmployeeID, err)
		}
		results = append(results, res)
	}

	// 4. Simpan snapshot dalam satu transaksi: hapus lama, tulis baru.
	if err := s.persistRunSnapshot(ctx, run.ID, results); err != nil {
		return nil, err
	}

	// 5. Update totals run.
	run.TotalEmployees = len(results)
	run.TotalEarning = 0
	run.TotalDeduction = 0
	run.TotalEmployerContribution = 0
	run.TotalNet = 0
	run.TotalCompanyCost = 0
	for _, res := range results {
		run.TotalEarning += res.runEmployee.TotalEarning
		run.TotalDeduction += res.runEmployee.TotalDeduction
		run.TotalEmployerContribution += res.runEmployee.TotalEmployerContribution
		run.TotalNet += res.runEmployee.NetAmount
		run.TotalCompanyCost += res.runEmployee.TotalCompanyCost
	}
	now := time.Now()
	run.CalculatedAt = &now
	if err := s.repo.UpdatePayrollRun(ctx, run); err != nil {
		return nil, err
	}

	s.logger.Info("Payroll run calculated",
		zap.String("run_id", run.ID.String()),
		zap.Int("employees", run.TotalEmployees),
		zap.Float64("total_net", run.TotalNet),
	)
	response := toPayrollRunResponse(run)
	return &response, nil
}

// selectRunEmployees menentukan employee yang dihitung:
//   - Jika run sudah punya payroll_run_employees (status != EXCLUDED) → dipakai.
//   - Jika belum → ambil seluruh employee active yang punya employee payroll
//     profile aktif yang mencakup periode (asOfDate).
func (s *Service) selectRunEmployees(ctx context.Context, run *PayrollRun, period *PayrollPeriod) ([]EmployeeRead, error) {
	existing, err := s.repo.FindPayrollRunEmployeesByRunID(ctx, run.ID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		ids := make([]uuid.UUID, 0, len(existing))
		for _, e := range existing {
			if e.Status == "EXCLUDED" {
				continue
			}
			ids = append(ids, e.EmployeeID)
		}
		if len(ids) == 0 {
			return nil, nil
		}
		return s.repo.FindEmployeesByIDs(ctx, ids)
	}

	// Auto-select: employee dengan payroll profile aktif pada asOfDate.
	profiles, err := s.findActivePayrollProfiles(ctx, period.AsOfDate)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.EmployeeID)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	return s.repo.FindEmployeesByIDs(ctx, ids)
}

// findActivePayrollProfiles mengambil profile payroll yang berlaku pada tanggal
// tertentu (is_payroll_active + status ACTIVE + rentang effective).
func (s *Service) findActivePayrollProfiles(ctx context.Context, asOfDate string) ([]EmployeePayrollProfile, error) {
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []EmployeePayrollProfile
	if err := db.Where(
		"is_payroll_active = ? AND status = ? AND effective_start_date <= ? AND (effective_end_date IS NULL OR effective_end_date >= ?)",
		true, "ACTIVE", asOfDate, asOfDate,
	).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// loadAllActiveComponents mengambil seluruh salary component ACTIVE.
func (s *Service) loadAllActiveComponents(ctx context.Context) ([]SalaryComponent, error) {
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []SalaryComponent
	if err := db.Where("status = ?", "ACTIVE").Order("display_order ASC, code ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// validateComponentFormulas memastikan seluruh formula komponen valid secara
// sintaks dan tidak ada circular dependency antar komponen yang direferensikan.
func (s *Service) validateComponentFormulas(components []SalaryComponent) error {
	engine := s.formulaEngine
	deps := map[string][]string{}
	for i := range components {
		c := &components[i]
		if c.Formula == nil || strings.TrimSpace(*c.Formula) == "" {
			continue
		}
		if err := engine.Validate(*c.Formula); err != nil {
			return &ValidationError{Message: fmt.Sprintf("formula komponen %s tidak valid: %v", c.Code, err)}
		}
		vars, err := engine.ReferencedVariables(*c.Formula)
		if err != nil {
			return &ValidationError{Message: fmt.Sprintf("gagal membaca referensi formula %s: %v", c.Code, err)}
		}
		deps[c.Code] = vars
	}
	cycles := calculator.DetectCycles(deps)
	if len(cycles) > 0 {
		msgs := make([]string, 0, len(cycles))
		for _, cyc := range cycles {
			msgs = append(msgs, cyc.String())
		}
		return &ValidationError{Message: fmt.Sprintf("terdapat circular dependency antar komponen: %s", strings.Join(msgs, "; "))}
	}
	return nil
}

// calculateEmployee menghitung seluruh komponen satu employee dan menghasilkan
// snapshot run employee + items.
func (s *Service) calculateEmployee(ctx context.Context, run *PayrollRun, period *PayrollPeriod, emp EmployeeRead, components []SalaryComponent) (employeeCalcResult, error) {
	ec := runEmployeeContext{employee: emp}
	asOf := period.AsOfDate

	// Resolusi employment → position → grading. Pakai overlap dengan periode
	// (bukan asOf akhir periode) agar employee join/resign tengah bulan tetap
	// punya struktur untuk diprorasi.
	employment, err := s.repo.FindEmploymentByEmployeeIDForPeriod(ctx, emp.ID, period.StartDate, period.EndDate)
	if err != nil {
		return employeeCalcResult{}, err
	}
	ec.employment = employment
	if employment != nil && employment.PositionID != nil {
		pos, err := s.repo.FindPositionByID(ctx, *employment.PositionID)
		if err != nil {
			return employeeCalcResult{}, err
		}
		ec.position = pos
		if pos.GradingID != nil {
			grading, err := s.repo.FindGradingByID(ctx, *pos.GradingID)
			if err != nil {
				return employeeCalcResult{}, err
			}
			ec.grading = grading
		}
	}

	// Load struktur gaji.
	if ec.grading != nil {
		gradeComps, err := s.repo.FindAllSalaryGradeComponentsByGradingID(ctx, ec.grading.ID, asOf)
		if err != nil {
			return employeeCalcResult{}, err
		}
		ec.gradeComps = gradeComps
	}
	empComps, err := s.repo.FindAllSalaryEmployeeComponentsByEmployeeID(ctx, emp.ID, asOf)
	if err != nil {
		return employeeCalcResult{}, err
	}
	ec.empComps = empComps
	adjustments, err := s.repo.FindAllSalaryEmployeeAdjustmentsByPeriod(ctx, emp.ID, period.PeriodYear, period.PeriodMonth)
	if err != nil {
		return employeeCalcResult{}, err
	}
	ec.adjustments = adjustments

	// Gabungkan jadi komponen siap hitung.
	inputs := s.buildComponentInputs(ec, components)

	// Summary workforce (attendance/leave/overtime) — input variabel built-in
	// formula (WORKING_DAYS, WORKED_DAYS, ABSENCE_DAYS, UNPAID_LEAVE_DAYS,
	// OVERTIME_HOURS). Jendela aktif dibatasi join/resign tengah bulan.
	activeFrom, activeTo := period.StartDate, period.EndDate
	if employment != nil {
		if employment.EffectiveDate != "" && employment.EffectiveDate > activeFrom {
			activeFrom = employment.EffectiveDate
		}
		if employment.EffectiveEndDate != nil && *employment.EffectiveEndDate != "" && *employment.EffectiveEndDate < activeTo {
			activeTo = *employment.EffectiveEndDate
		}
	}
	workforce, err := s.loadWorkforceSummary(ctx, emp.ID, period.StartDate, period.EndDate, activeFrom, activeTo)
	if err != nil {
		return employeeCalcResult{}, err
	}

	// Proration: faktor join/resign tengah bulan sesuai metode run.
	prorateFactor := 1.0
	if employment != nil && employment.EffectiveDate != "" {
		endDate := ""
		if employment.EffectiveEndDate != nil {
			endDate = *employment.EffectiveEndDate
		}
		factor, err := s.employmentProrationFactor(employment.EffectiveDate, endDate, period.StartDate, period.EndDate, run.ProrationMethod, workforce)
		if err != nil {
			return employeeCalcResult{}, err
		}
		prorateFactor = factor
	}

	// Evaluasi formula dengan dependency resolver. values berisi nilai terhitung
	// per komponen (code → amount), dipakai dasar upah BPJS.
	items, values, err := s.evaluateComponents(ec, inputs, prorateFactor, workforce)
	if err != nil {
		return employeeCalcResult{}, err
	}

	// Kontribusi BPJS (statutory) — setelah struktur selesai supaya dasar upah
	// (is_bpjs_base) tersedia di values.
	bpjsItems, err := s.calculateBpjsContributions(ctx, period, ec, components, values)
	if err != nil {
		return employeeCalcResult{}, err
	}
	items = append(items, bpjsItems...)

	// PPh21 (statutory) — setelah BPJS karena butuh iuran BPJS yang boleh
	// dikurangkan sebagai pengurang penghasilan.
	pph21Item, pph21Log, err := s.calculatePph21(ctx, run, period, ec, components, values, bpjsItems)
	if err != nil {
		return employeeCalcResult{}, err
	}
	if pph21Item != nil {
		items = append(items, *pph21Item)
	}

	// Hitung agregat.
	var totalEarning, totalDeduction, totalEmployerContribution float64
	for _, item := range items {
		switch item.ItemCategory {
		case ItemCategoryEmployeeEarning:
			totalEarning += item.Amount
		case ItemCategoryEmployeeDeduction:
			totalDeduction += item.Amount
		case ItemCategoryEmployerContribution:
			totalEmployerContribution += item.Amount
		}
	}
	netAmount := totalEarning - totalDeduction
	totalCompanyCost := totalEarning + totalEmployerContribution

	runEmployee := PayrollRunEmployee{
		PayrollRunID:              run.ID,
		EmployeeID:                emp.ID,
		EmployeeCode:              emp.EmployeeID,
		EmployeeName:              emp.Name,
		TotalEarning:              totalEarning,
		TotalDeduction:            totalDeduction,
		TotalEmployerContribution: totalEmployerContribution,
		NetAmount:                 netAmount,
		TotalCompanyCost:          totalCompanyCost,
		Status:                    "CALCULATED",
	}
	if ec.employment != nil {
		runEmployee.EmploymentID = &ec.employment.ID
	}
	if ec.position != nil {
		runEmployee.PositionID = &ec.position.ID
		runEmployee.PositionTitle = stringPtr(ec.position.Title)
	}
	if ec.grading != nil {
		runEmployee.GradingID = &ec.grading.ID
		runEmployee.GradingName = stringPtr(ec.grading.Name)
	}

	return employeeCalcResult{runEmployee: runEmployee, items: items, pph21Log: pph21Log}, nil
}

// employmentProrationFactor menghitung faktor prorasi employee yang join
// dan/atau resign tengah bulan sesuai metode run (docs/payroll/06 §19):
//   - CALENDAR_DAYS: hari kalender efektif / total hari kalender periode
//   - FIXED_30_DAYS: hari kalender efektif / 30
//   - WORKING_DAYS:  hari kerja efektif / total hari kerja (weekday)
//   - ATTENDANCE_DAYS: hari hadir aktual / hari kerja (dari workforce summary)
func (s *Service) employmentProrationFactor(effectiveDate, effectiveEndDate, periodStart, periodEnd, method string, workforce *workforceSummary) (float64, error) {
	join, err := parseFlexibleDate(effectiveDate)
	if err != nil {
		return 1.0, nil // tanggal tidak terparse → tanpa prorasi (bukan blocker)
	}
	start, err := parseFlexibleDate(periodStart)
	if err != nil {
		return 1.0, nil
	}
	end, err := parseFlexibleDate(periodEnd)
	if err != nil {
		return 1.0, nil
	}
	if end.Before(start) {
		return 1.0, nil
	}

	// Jendela aktif: dari join (atau awal periode) sampai resign (atau akhir).
	eligibleStart := join
	if start.After(eligibleStart) {
		eligibleStart = start
	}
	eligibleEnd := end
	if effectiveEndDate != "" {
		if t, err := parseFlexibleDate(effectiveEndDate); err == nil && t.Before(eligibleEnd) {
			eligibleEnd = t
		}
	}
	if eligibleStart.After(eligibleEnd) {
		return 0, nil // join setelah periode berakhir / resign sebelum join
	}

	switch calculator.ProrationMethod(method) {
	case calculator.ProrationWorkingDays:
		total := weekdayCount(start, end)
		if total <= 0 {
			return 1.0, nil
		}
		return clampFactor(weekdayCount(eligibleStart, eligibleEnd) / total), nil
	case calculator.ProrationFixed30Days:
		return clampFactor(eligibleCalendarDays(eligibleStart, eligibleEnd) / 30), nil
	case calculator.ProrationAttendanceDays:
		if workforce != nil && workforce.WorkingDays > 0 {
			return clampFactor(workforce.WorkedDays / workforce.WorkingDays), nil
		}
		return 1.0, nil
	default: // CALENDAR_DAYS
		total := calculator.TotalCalendarDays(start, end)
		if total <= 0 {
			return 1.0, nil
		}
		return clampFactor(eligibleCalendarDays(eligibleStart, eligibleEnd) / total), nil
	}
}

// parseFlexibleDate mem-parse tanggal yang bisa berupa "YYYY-MM-DD" atau
// timestamp penuh (mis. "2026-01-16T00:00:00Z" dari SQLite).
func parseFlexibleDate(s string) (time.Time, error) {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) >= 10 {
		trimmed = trimmed[:10]
	}
	return time.Parse("2006-01-02", trimmed)
}

// buildComponentInputs menggabungkan grade default + override employee
// (override menang) + adjustment periode (tambahan sekali-jalan).
func (s *Service) buildComponentInputs(ec runEmployeeContext, components []SalaryComponent) []componentInput {
	compByID := map[uuid.UUID]SalaryComponent{}
	for _, c := range components {
		compByID[c.ID] = c
	}

	gradeByComp := map[uuid.UUID]SalaryGradeComponent{}
	for _, gc := range ec.gradeComps {
		gradeByComp[gc.SalaryComponentID] = gc
	}
	empByComp := map[uuid.UUID]SalaryEmployeeComponent{}
	for _, ec2 := range ec.empComps {
		empByComp[ec2.SalaryComponentID] = ec2
	}
	adjByComp := map[uuid.UUID][]SalaryEmployeeAdjustment{}
	for _, a := range ec.adjustments {
		adjByComp[a.SalaryComponentID] = append(adjByComp[a.SalaryComponentID], a)
	}

	// Map compID → input. Grade default diisi dulu, lalu override employee
	// MENGGANTI nilai default untuk komponen yang sama (override menang).
	inputsByComp := map[uuid.UUID]*componentInput{}
	getInput := func(compID uuid.UUID, sourceGroup, sourceTable string) *componentInput {
		if in, ok := inputsByComp[compID]; ok {
			return in
		}
		comp, ok := compByID[compID]
		if !ok {
			return nil // komponen sudah non-aktif → dilewati
		}
		in := &componentInput{
			component:    comp,
			sourceGroup:  sourceGroup,
			sourceTable:  sourceTable,
			displayOrder: comp.DisplayOrder,
		}
		inputsByComp[compID] = in
		return in
	}

	// 1. Grade default.
	for _, gc := range ec.gradeComps {
		in := getInput(gc.SalaryComponentID, SourceGroupStructure, "salary_grade_components")
		if in == nil {
			continue
		}
		in.baseAmount = gc.Amount
		for _, a := range adjByComp[gc.SalaryComponentID] {
			in.adjustmentAmt += a.Amount
			in.sourceID = &a.ID
		}
	}

	// 2. Override employee — menang atas grade default.
	for _, ec2 := range ec.empComps {
		in := getInput(ec2.SalaryComponentID, SourceGroupStructure, "salary_employee_components")
		if in == nil {
			continue
		}
		in.baseAmount = ec2.Amount
		for _, a := range adjByComp[ec2.SalaryComponentID] {
			in.adjustmentAmt += a.Amount
			in.sourceID = &a.ID
		}
	}

	// 3. Komponen yang hanya muncul sebagai adjustment (tidak di struktur).
	for compID, adjs := range adjByComp {
		in := getInput(compID, SourceGroupAdjustment, "salary_employee_adjustments")
		if in == nil {
			continue
		}
		for _, a := range adjs {
			in.adjustmentAmt += a.Amount
			in.sourceID = &a.ID
		}
	}

	inputs := make([]componentInput, 0, len(inputsByComp))
	for _, in := range inputsByComp {
		inputs = append(inputs, *in)
	}
	sort.SliceStable(inputs, func(i, j int) bool {
		return inputs[i].displayOrder < inputs[j].displayOrder
	})
	return inputs
}

// evaluateComponents menghitung nilai tiap komponen. Komponen FORMULA/REFERENCE
// dievaluasi dengan resolver yang membaca hasil komponen lain (dependency
// resolver) + variabel built-in. Iterasi berhenti saat tidak ada progress lagi.
func (s *Service) evaluateComponents(ec runEmployeeContext, inputs []componentInput, prorateFactor float64, workforce *workforceSummary) ([]PayrollRunItem, map[string]float64, error) {
	// Map component ID → code untuk resolusi REFERENCE (referensi antar komponen).
	codeByID := map[uuid.UUID]string{}
	for _, in := range inputs {
		codeByID[in.component.ID] = in.component.Code
	}

	values := map[string]float64{} // code → amount terhitung
	var items []PayrollRunItem

	// Pass 1: komponen FIXED/MANUAL langsung dihitung.
	pending := make([]componentInput, 0, len(inputs))
	for _, in := range inputs {
		if in.component.CalculationType == CalculationTypeFixed || in.component.CalculationType == CalculationTypeManual {
			amt := in.baseAmount + in.adjustmentAmt
			if in.component.IsProratable && prorateFactor < 1 {
				amt = calculator.ProrateToWholeNumber(amt * prorateFactor)
			}
			values[in.component.Code] = amt
			items = append(items, s.buildRunItem(ec, in, amt))
		} else {
			pending = append(pending, in)
		}
	}

	// Pass 2..n: FORMULA/PERCENTAGE/REFERENCE — iterasi sampai stabil.
	remaining := pending
	for len(remaining) > 0 {
		progress := false
		var next []componentInput
		for _, in := range remaining {
			amt, ok, err := s.tryEvaluateComponent(in, values, codeByID, workforce)
			if err != nil {
				return nil, nil, err
			}
			if !ok {
				next = append(next, in)
				continue
			}
			progress = true
			values[in.component.Code] = amt
			items = append(items, s.buildRunItem(ec, in, amt))
		}
		if !progress {
			if len(next) > 0 {
				var missing []string
				for _, in := range next {
					missing = append(missing, in.component.Code)
				}
				return nil, nil, fmt.Errorf("tidak bisa menyelesaikan dependensi komponen: %s (periksa formula/reference)", strings.Join(missing, ", "))
			}
			break
		}
		remaining = next
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ComponentCode < items[j].ComponentCode
	})
	return items, values, nil
}

// tryEvaluateComponent mencoba menghitung komponen FORMULA/PERCENTAGE/REFERENCE.
// Return ok=false jika dependensi belum tersedia.
func (s *Service) tryEvaluateComponent(in componentInput, values map[string]float64, codeByID map[uuid.UUID]string, workforce *workforceSummary) (float64, bool, error) {
	switch in.component.CalculationType {
	case CalculationTypeReference:
		if in.component.ReferenceComponentID == nil {
			return 0, false, fmt.Errorf("komponen %s bertipe REFERENCE tanpa reference_component_id", in.component.Code)
		}
		targetCode, known := codeByID[*in.component.ReferenceComponentID]
		if !known {
			return 0, false, fmt.Errorf("komponen %s mereferensikan komponen yang tidak ada di struktur (id %s)", in.component.Code, in.component.ReferenceComponentID.String())
		}
		targetVal, ok := values[targetCode]
		if !ok {
			return 0, false, nil // target belum terhitung → tunggu iterasi berikutnya
		}
		return targetVal, true, nil
	default:
		if in.component.Formula == nil || strings.TrimSpace(*in.component.Formula) == "" {
			return 0, false, fmt.Errorf("komponen %s bertipe %s tanpa formula", in.component.Code, in.component.CalculationType)
		}
		return s.evaluateFormulaWithValues(in, values, workforce)
	}
}

// evaluateFormulaWithValues mengevaluasi formula dengan resolver yang membaca
// nilai komponen terhitung + variabel built-in. Return ok=false bila ada
// variabel non-built-in yang belum tersedia (dependency belum resolve).
func (s *Service) evaluateFormulaWithValues(in componentInput, values map[string]float64, workforce *workforceSummary) (float64, bool, error) {
	engine := s.formulaEngine
	formula := *in.component.Formula
	vars, err := engine.ReferencedVariables(formula)
	if err != nil {
		return 0, false, err
	}
	for _, v := range vars {
		if engine.Registry().IsBuiltIn(v) {
			continue
		}
		if _, ok := values[v]; !ok {
			return 0, false, nil // dependency belum tersedia
		}
	}
	amt, err := engine.Evaluate(formula, func(name string) (float64, bool) {
		if v, ok := values[name]; ok {
			return v, true
		}
		// Variabel built-in workforce (attendance/leave/overtime) — diisi dari
		// summary hasil final Workforce Management.
		if v, ok := workforce.BuiltInValue(name); ok {
			return v, true
		}
		// Variabel built-in lain yang belum dihitung dianggap 0 (mis. GROSS) —
		// dipakai untuk formula agregat seperti NET_SALARY.
		if engine.Registry().IsBuiltIn(name) {
			return 0, true
		}
		return 0, false
	})
	if err != nil {
		return 0, false, err
	}
	return amt, true, nil
}

// buildRunItem membangun snapshot item payroll.
func (s *Service) buildRunItem(ec runEmployeeContext, in componentInput, amount float64) PayrollRunItem {
	item := PayrollRunItem{
		EmployeeID:        ec.employee.ID,
		SalaryComponentID: in.component.ID,
		ComponentCode:     in.component.Code,
		ComponentName:     in.component.Name,
		ComponentType:     in.component.ComponentType,
		CalculationType:   in.component.CalculationType,
		Amount:            amount,
		BaseAmount:        in.baseAmount + in.adjustmentAmt,
		CurrencyCode:      "IDR",
		SourceGroup:       in.sourceGroup,
		PrintOnPayslip:    in.component.PrintOnSalaryStructure,
	}
	if in.component.Formula != nil {
		item.Formula = in.component.Formula
		item.FormulaResult = &amount
	}
	if in.sourceTable != "" {
		item.SourceTable = &in.sourceTable
	}
	item.SourceID = in.sourceID

	switch in.component.ComponentType {
	case "EARNING":
		item.ItemCategory = ItemCategoryEmployeeEarning
		item.PaidBy = "EMPLOYER"
		item.AffectsGrossPay = true
		item.AffectsNetPay = true
		item.AffectsCompanyCost = true
	case "DEDUCTION":
		item.ItemCategory = ItemCategoryEmployeeDeduction
		item.PaidBy = "EMPLOYEE"
		item.AffectsNetPay = true
	case "EMPLOYER_CONTRIBUTION":
		item.ItemCategory = ItemCategoryEmployerContribution
		item.PaidBy = "EMPLOYER"
		item.AffectsCompanyCost = true
	default:
		item.ItemCategory = ItemCategoryInformation
		item.PaidBy = "NONE"
	}
	return item
}

// persistRunSnapshot menyimpan snapshot dalam satu transaksi: hapus item lama,
// hapus employee lama, tulis employee baru, tulis item baru (relasi FK).
func (s *Service) persistRunSnapshot(ctx context.Context, runID uuid.UUID, results []employeeCalcResult) error {
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// Hapus snapshot lama (log PPh21 & item dulu karena FK ke employees).
		if err := tx.Where("payroll_run_id = ?", runID).Delete(&Pph21CalculationLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("payroll_run_id = ?", runID).Delete(&PayrollRunItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("payroll_run_id = ?", runID).Delete(&PayrollRunEmployee{}).Error; err != nil {
			return err
		}

		// Tulis employee baru — simpan ID hasil insert untuk relasi item & log.
		for i := range results {
			emp := results[i].runEmployee
			if err := tx.Create(&emp).Error; err != nil {
				return err
			}
			results[i].runEmployee = emp
			for j := range results[i].items {
				results[i].items[j].PayrollRunID = runID
				results[i].items[j].PayrollRunEmployeeID = emp.ID
			}
			if err := tx.CreateInBatches(results[i].items, 100).Error; err != nil {
				return err
			}

			// Tulis log PPh21 — tautkan ke employee & item yang baru dibuat.
			if results[i].pph21Log != nil {
				log := *results[i].pph21Log
				log.PayrollRunID = runID
				log.PayrollRunEmployeeID = emp.ID
				for j := range results[i].items {
					if results[i].items[j].SourceType != nil && *results[i].items[j].SourceType == "PPH21" {
						id := results[i].items[j].ID
						log.PayrollRunItemID = &id
						break
					}
				}
				if err := tx.Create(&log).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// stringPtr helper untuk *string.
func stringPtr(s string) *string { return &s }
