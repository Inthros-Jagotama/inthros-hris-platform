package payroll

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// createSessionsForWeekdaysExcept membuat session CLOSED untuk seluruh weekday
// periode Jan 2026 kecuali tanggal yang di-exclude (format "2006-01-02").
func createSessionsForWeekdaysExcept(ctx context.Context, repo *Repository, employeeID uuid.UUID, exclude map[string]bool, overtimeMinutes int) {
	for d := 1; d <= 31; d++ {
		date := time.Date(2026, 1, d, 0, 0, 0, 0, time.UTC)
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}
		ds := date.Format("2006-01-02")
		if exclude[ds] {
			continue
		}
		createTestAttendanceSession(ctx, repo, employeeID, ds, "CLOSED", overtimeMinutes)
	}
}

func runCalc(t *testing.T, repo *Repository, svc *Service, ctx context.Context, runCode string, prorationMethod string) (*PayrollRunResponse, error) {
	t.Helper()
	period := createTestPayrollPeriod(ctx, repo)
	run, err := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID.String(),
		RunCode:         runCode,
		ProrationMethod: prorationMethod,
	})
	if err != nil {
		t.Fatalf("CreatePayrollRun: %v", err)
	}
	return svc.CalculatePayrollRun(ctx, run.ID)
}

// TestProrationResignMidPeriod: employee resign 16 Jan → eligible 16 dari 31
// hari kalender → BASIC 3.1jt × 16/31 = 1.6jt.
func TestProrationResignMidPeriod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-RESIGN", "Endang")
	emp := createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2020-01-01")
	// Resign tengah bulan.
	db, _ := repo.getDB(ctx)
	if err := db.Model(&EmploymentRead{}).Where("id = ?", emp.ID).Update("effective_end_date", "2026-01-16").Error; err != nil {
		t.Fatalf("set resign date: %v", err)
	}
	createTestPayrollProfile(ctx, repo, employee.ID, "2020-01-01")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 3100000)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-RESIGN", "")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "prorated earning", resp.TotalEarning, 1600000)
}

// TestProrationWorkingDaysMethod: join 16 Jan, metode WORKING_DAYS → 11/22
// hari kerja = 0.5 → BASIC 3.1jt × 0.5 = 1.55jt.
func TestProrationWorkingDaysMethod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading := createTestGrading(ctx, repo, "G-1", "Staff")
	position := createTestPosition(ctx, repo, "Staff HR", &grading.ID)
	employee := createTestEmployee(ctx, repo, "EMP-WD", "Fajar")
	createTestEmployment(ctx, repo, &employee.ID, &position.ID, "2026-01-16")
	createTestPayrollProfile(ctx, repo, employee.ID, "2026-01-16")

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 3100000)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-WORKINGDAYS", "WORKING_DAYS")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "working-days prorated earning", resp.TotalEarning, 1550000)
}

// TestProrationAttendanceDaysMethod: employee kerja 11 dari 22 hari kerja,
// metode ATTENDANCE_DAYS → faktor 0.5 → BASIC 3.1jt × 0.5 = 1.55jt.
func TestProrationAttendanceDaysMethod(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)
	// 11 session weekday (tanpa tanggal 8 & 15 Jan + beberapa lain) — cukup
	// 11 hari hadir dari 22 hari kerja.
	worked := 0
	for d := 1; d <= 31 && worked < 11; d++ {
		date := time.Date(2026, 1, d, 0, 0, 0, 0, time.UTC)
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}
		createTestAttendanceSession(ctx, repo, employee.ID, date.Format("2006-01-02"), "CLOSED", 0)
		worked++
	}

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 3100000)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-ATTDAYS", "ATTENDANCE_DAYS")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "attendance-days prorated earning", resp.TotalEarning, 1550000)
}

// TestWorkforceAbsenceDeduction: 2 hari absen (tanpa session) dari 22 hari
// kerja → formula BASIC / WORKING_DAYS * ABSENCE_DAYS = 22jt/22×2 = 2jt potong.
func TestWorkforceAbsenceDeduction(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 22000000)

	// Session semua weekday kecuali 8 & 15 Jan (2 hari absen).
	exclude := map[string]bool{"2026-01-08": true, "2026-01-15": true}
	createSessionsForWeekdaysExcept(ctx, repo, employee.ID, exclude, 0)

	// Potongan absence via formula.
	formula := "BASIC / WORKING_DAYS * ABSENCE_DAYS"
	absent := createTestComponent(ctx, repo, "ABSENT_DEDUCT", "Potongan Absen", "DEDUCTION", "FORMULA", 20)
	absent.Formula = &formula
	if err := repo.UpdateSalaryComponent(ctx, absent); err != nil {
		t.Fatalf("update absent component: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, absent.ID, 0)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-ABSENCE", "")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "earning", resp.TotalEarning, 22000000)
	assertClose(t, "absence deduction", resp.TotalDeduction, 2000000)
	assertClose(t, "net", resp.TotalNet, 20000000)
}

// TestWorkforceOvertimeEarning: 3 session × 120 menit lembur = 6 jam →
// formula OVERTIME_HOURS * 15000 = 90rb earning.
func TestWorkforceOvertimeEarning(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 22000000)

	// 3 session dengan lembur 120 menit (semua weekday lain hadir normal).
	createTestAttendanceSession(ctx, repo, employee.ID, "2026-01-06", "CLOSED", 120)
	createTestAttendanceSession(ctx, repo, employee.ID, "2026-01-07", "CLOSED", 120)
	createTestAttendanceSession(ctx, repo, employee.ID, "2026-01-08", "CLOSED", 120)
	createSessionsForWeekdaysExcept(ctx, repo, employee.ID, map[string]bool{"2026-01-06": true, "2026-01-07": true, "2026-01-08": true}, 0)

	formula := "OVERTIME_HOURS * 15000"
	ot := createTestComponent(ctx, repo, "OVERTIME_PAY", "Lembur", "EARNING", "FORMULA", 20)
	ot.Formula = &formula
	if err := repo.UpdateSalaryComponent(ctx, ot); err != nil {
		t.Fatalf("update overtime component: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, ot.ID, 0)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-OVERTIME", "")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "overtime earning", resp.TotalEarning, 22090000)
}

// TestWorkforceUnpaidLeave: 2 hari cuti tidak berbayar (APPROVED_FINAL,
// is_paid=false) → formula BASIC / WORKING_DAYS * UNPAID_LEAVE_DAYS = 2jt potong.
func TestWorkforceUnpaidLeave(t *testing.T) {
	_, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	grading, _, employee := setupCalcEnv(t, repo, ctx)

	basic := createTestComponent(ctx, repo, "BASIC", "Basic Salary", "EARNING", "FIXED", 10)
	createTestGradeComponent(ctx, repo, grading.ID, basic.ID, 22000000)

	// Hadir semua weekday (tidak ada absence).
	createSessionsForWeekdaysExcept(ctx, repo, employee.ID, nil, 0)

	// 2 hari cuti tidak berbayar yang disetujui.
	lr := createTestLeaveRequest(ctx, repo, employee.ID, "APPROVED_FINAL")
	createTestLeaveRequestDetail(ctx, repo, lr.ID, employee.ID, "2026-01-08", false, 1.0)
	createTestLeaveRequestDetail(ctx, repo, lr.ID, employee.ID, "2026-01-15", false, 1.0)

	formula := "BASIC / WORKING_DAYS * UNPAID_LEAVE_DAYS"
	ulp := createTestComponent(ctx, repo, "UNPAID_LEAVE_DEDUCT", "Potongan Cuti Tanpa Bayar", "DEDUCTION", "FORMULA", 20)
	ulp.Formula = &formula
	if err := repo.UpdateSalaryComponent(ctx, ulp); err != nil {
		t.Fatalf("update ulp component: %v", err)
	}
	createTestGradeComponent(ctx, repo, grading.ID, ulp.ID, 0)

	svc := NewService(repo, zap.NewNop())
	resp, err := runCalc(t, repo, svc, ctx, "RUN-UNPAID", "")
	if err != nil {
		t.Fatalf("CalculatePayrollRun: %v", err)
	}
	assertClose(t, "unpaid leave deduction", resp.TotalDeduction, 2000000)
	assertClose(t, "net", resp.TotalNet, 20000000)
}
