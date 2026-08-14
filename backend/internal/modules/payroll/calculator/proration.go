package calculator

import (
	"fmt"
	"math"
	"time"
)

// ProrationMethod menentukan metode prorasi gaji (docs/payroll/06 §19).
type ProrationMethod string

const (
	ProrationCalendarDays   ProrationMethod = "CALENDAR_DAYS"   // jumlah hari kalender
	ProrationWorkingDays    ProrationMethod = "WORKING_DAYS"    // jumlah hari kerja
	ProrationFixed30Days    ProrationMethod = "FIXED_30_DAYS"   // dibagi 30 (konvensi umum payroll)
	ProrationAttendanceDays ProrationMethod = "ATTENDANCE_DAYS" // berdasarkan kehadiran aktual
)

// IsValidProrationMethod memvalidasi metode prorasi.
func IsValidProrationMethod(method string) bool {
	switch ProrationMethod(method) {
	case ProrationCalendarDays, ProrationWorkingDays, ProrationFixed30Days, ProrationAttendanceDays:
		return true
	default:
		return false
	}
}

// Prorate menghitung jumlah yang diprorasi:
//
//	amount * (eligible / total)
//
// totalDays harus > 0. eligibleDays dibatasi ke [0, totalDays].
func Prorate(amount float64, eligibleDays, totalDays float64, method ProrationMethod) (float64, error) {
	if totalDays <= 0 {
		return 0, fmt.Errorf("total days harus lebih dari 0 untuk prorasi")
	}
	if eligibleDays < 0 {
		eligibleDays = 0
	}
	if eligibleDays > totalDays {
		eligibleDays = totalDays
	}
	switch method {
	case ProrationCalendarDays, ProrationWorkingDays, ProrationFixed30Days, ProrationAttendanceDays:
		return amount * (eligibleDays / totalDays), nil
	default:
		return 0, fmt.Errorf("metode prorasi tidak dikenal: %s", method)
	}
}

// EligibleCalendarDays menghitung jumlah hari kalender efektif antara joinDate
// (atau periodStart) hingga periodEnd, inklusif.
func EligibleCalendarDays(joinDate, periodStart, periodEnd time.Time) float64 {
	start := joinDate
	if periodStart.After(start) {
		start = periodStart
	}
	if start.After(periodEnd) {
		return 0
	}
	return float64(periodEnd.Sub(start).Hours()/24) + 1
}

// TotalCalendarDays menghitung jumlah hari kalender dalam rentang [start, end].
func TotalCalendarDays(start, end time.Time) float64 {
	if start.After(end) {
		return 0
	}
	return float64(end.Sub(start).Hours()/24) + 1
}

// ProrateFixed30Days menghitung prorasi dengan konvensi 30 hari:
// amount * (eligibleDays / 30), eligibleDays = hari kalender efektif.
func ProrateFixed30Days(amount float64, eligibleDays float64) (float64, error) {
	if eligibleDays < 0 {
		eligibleDays = 0
	}
	return amount * (eligibleDays / 30), nil
}

// ProrateToWholeNumber membulatkan hasil prorasi ke rupiah penuh (ke bawah),
// sesuai konvensi umum payroll Indonesia.
func ProrateToWholeNumber(value float64) float64 {
	return math.Floor(value)
}
