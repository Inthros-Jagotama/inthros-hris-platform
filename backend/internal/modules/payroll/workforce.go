package payroll

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
)

// workforceNonWorkedStatuses — status session attendance yang TIDAK dihitung
// sebagai hari hadir (sinkron dengan const di modul attendance).
var workforceNonWorkedStatuses = map[string]bool{
	"ABSENT":  true,
	"DAY_OFF": true,
	"EXEMPT":  true,
	"LEAVE":   true,
}

// workforceSummary merangkum hasil FINAL Workforce Management untuk satu
// employee dalam rentang periode. Payroll tidak menghitung ulang attendance —
// hanya mengonsumsi hasil akhir (docs/payroll/06 §20-22).
type workforceSummary struct {
	WorkingDays     float64 // hari kerja terjadwal (weekday) dalam rentang aktif
	WorkedDays      float64 // hari hadir (session non-absent/off/exempt/leave)
	AbsenceDays     float64 // weekday tanpa kehadiran & tanpa cuti disetujui
	PaidLeaveDays   float64 // hari cuti berbayar (detail is_paid=true, APPROVED_FINAL)
	UnpaidLeaveDays float64 // hari cuti tidak berbayar (is_paid=false)
	OvertimeHours   float64 // total jam lembur dari session (overtime_minutes)
}

// BuiltInValue mengembalikan nilai variabel built-in formula yang bersumber
// dari workforce. Variabel lain (GROSS, dst.) tetap di-handle di tempat lain.
func (w *workforceSummary) BuiltInValue(name string) (float64, bool) {
	if w == nil {
		return 0, false
	}
	switch name {
	case "WORKING_DAYS":
		return w.WorkingDays, true
	case "WORKED_DAYS":
		return w.WorkedDays, true
	case "ABSENCE_DAYS":
		return w.AbsenceDays, true
	case "UNPAID_LEAVE_DAYS":
		return w.UnpaidLeaveDays, true
	case "OVERTIME_HOURS":
		return w.OvertimeHours, true
	}
	return 0, false
}

// loadWorkforceSummary mengambil summary kehadiran/cuti/lembur employee untuk
// rentang periode [periodStart, periodEnd]. activeFrom/activeTo membatasi jendela
// "hari kerja" (mis. join/resign tengah bulan) — dipakai menghitung
// WORKING_DAYS & absence supaya hari sebelum join / setelah resign tidak
// dianggap alpa.
func (s *Service) loadWorkforceSummary(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd, activeFrom, activeTo string) (*workforceSummary, error) {
	db, err := s.repo.getDB(ctx)
	if err != nil {
		return nil, err
	}
	w := &workforceSummary{}

	start, err1 := parseFlexibleDate(periodStart)
	end, err2 := parseFlexibleDate(periodEnd)
	if err1 != nil || err2 != nil || start.After(end) {
		return w, nil
	}
	// Normalisasi ke YYYY-MM-DD supaya perbandingan BETWEEN konsisten
	// (kolom date bisa terbaca sebagai timestamp penuh di SQLite).
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")
	workStart := start
	if t, err := parseFlexibleDate(activeFrom); err == nil && t.After(workStart) {
		workStart = t
	}
	workEnd := end
	if activeTo != "" {
		if t, err := parseFlexibleDate(activeTo); err == nil && t.Before(workEnd) {
			workEnd = t
		}
	}
	if workStart.After(workEnd) {
		return w, nil
	}
	w.WorkingDays = weekdayCount(workStart, workEnd)

	// Sessions kehadiran (per employee per hari).
	var sessions []AttendanceSessionRead
	if err := db.Where(
		"employee_id = ? AND work_date BETWEEN ? AND ? AND deleted_at IS NULL",
		employeeID, startStr, endStr,
	).Find(&sessions).Error; err != nil {
		return nil, err
	}
	sessionDates := map[string]bool{}
	for _, sess := range sessions {
		// Kolom date bisa terbaca sebagai timestamp penuh (mis. SQLite) —
		// normalisasi ke YYYY-MM-DD untuk pencocokan per hari.
		if t, err := parseFlexibleDate(sess.WorkDate); err == nil {
			sessionDates[t.Format("2006-01-02")] = true
		}
		if !workforceNonWorkedStatuses[sess.Status] {
			w.WorkedDays++
		}
		w.OvertimeHours += float64(sess.OvertimeMinutes) / 60
	}
	w.OvertimeHours = math.Round(w.OvertimeHours*100) / 100

	// Detail cuti yang APPROVED_FINAL (is_paid + day_fraction).
	approvedLeave := db.Model(&LeaveRequestRead{}).
		Select("id").
		Where("employee_id = ? AND status = ? AND deleted_at IS NULL", employeeID, "APPROVED_FINAL")
	var details []LeaveRequestDetailRead
	if err := db.Where(
		"employee_id = ? AND leave_date BETWEEN ? AND ? AND leave_request_id IN (?)",
		employeeID, startStr, endStr, approvedLeave,
	).Find(&details).Error; err != nil {
		return nil, err
	}
	leaveDates := map[string]bool{}
	for _, d := range details {
		if t, err := parseFlexibleDate(d.LeaveDate); err == nil {
			leaveDates[t.Format("2006-01-02")] = true
		}
		if d.IsPaid {
			w.PaidLeaveDays += d.DayFraction
		} else {
			w.UnpaidLeaveDays += d.DayFraction
		}
	}

	// Absence = weekday dalam jendela kerja tanpa session & tanpa cuti disetujui.
	for d := workStart; !d.After(workEnd); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		ds := d.Format("2006-01-02")
		if sessionDates[ds] || leaveDates[ds] {
			continue
		}
		w.AbsenceDays++
	}
	return w, nil
}

// weekdayCount menghitung jumlah hari Senin-Jumat dalam rentang [start, end].
func weekdayCount(start, end time.Time) float64 {
	if start.After(end) {
		return 0
	}
	n := 0.0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if w := d.Weekday(); w != time.Saturday && w != time.Sunday {
			n++
		}
	}
	return n
}

// eligibleCalendarDays menghitung hari kalender efektif inklusif [start, end].
func eligibleCalendarDays(start, end time.Time) float64 {
	if start.After(end) {
		return 0
	}
	return float64(end.Sub(start).Hours()/24) + 1
}

// clampFactor membatasi faktor prorasi ke [0, 1].
func clampFactor(f float64) float64 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}
