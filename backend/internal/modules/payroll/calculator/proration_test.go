package calculator

import (
	"math"
	"testing"
	"time"
)

func TestProrate(t *testing.T) {
	// 10.000.000 / 22 * 12 (contoh di docs/payroll/06 §19)
	got, err := Prorate(10000000, 12, 22, ProrationWorkingDays)
	if err != nil {
		t.Fatalf("Prorate: %v", err)
	}
	if math.Abs(got-5454545.454545) > 0.01 {
		t.Errorf("expected ~5454545.45, got %v", got)
	}
}

func TestProrateCalendarDays(t *testing.T) {
	// 1 bulan penuh 31 hari → jumlah tidak berubah.
	got, err := Prorate(3100000, 31, 31, ProrationCalendarDays)
	if err != nil {
		t.Fatalf("Prorate: %v", err)
	}
	if got != 3100000 {
		t.Errorf("expected full amount 3100000, got %v", got)
	}

	// Bergabung di tengah bulan: 15 hari dari 30.
	got2, err := Prorate(3000000, 15, 30, ProrationCalendarDays)
	if err != nil {
		t.Fatalf("Prorate: %v", err)
	}
	if got2 != 1500000 {
		t.Errorf("expected 1500000, got %v", got2)
	}
}

func TestProrateFixed30Days(t *testing.T) {
	got, err := ProrateFixed30Days(3000000, 15)
	if err != nil {
		t.Fatalf("ProrateFixed30Days: %v", err)
	}
	if got != 1500000 {
		t.Errorf("expected 1500000, got %v", got)
	}
}

func TestProrateInvalid(t *testing.T) {
	if _, err := Prorate(1000, 10, 0, ProrationCalendarDays); err == nil {
		t.Error("expected error for zero total days")
	}
	if _, err := Prorate(1000, 10, 30, "WEIRD"); err == nil {
		t.Error("expected error for unknown method")
	}
	if IsValidProrationMethod("WEIRD") {
		t.Error("WEIRD should not be a valid proration method")
	}
}

func TestProrateClampsEligible(t *testing.T) {
	// eligible > total → dibatasi ke total (full amount).
	got, err := Prorate(1000000, 40, 30, ProrationCalendarDays)
	if err != nil {
		t.Fatalf("Prorate: %v", err)
	}
	if got != 1000000 {
		t.Errorf("expected clamped full amount 1000000, got %v", got)
	}
}

func TestEligibleCalendarDays(t *testing.T) {
	periodStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)

	// Bergabung 1 Jan → 31 hari penuh.
	if got := EligibleCalendarDays(periodStart, periodStart, periodEnd); got != 31 {
		t.Errorf("expected 31 days, got %v", got)
	}

	// Bergabung 16 Jan → 16 hari (16..31 inklusif).
	join := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)
	if got := EligibleCalendarDays(join, periodStart, periodEnd); got != 16 {
		t.Errorf("expected 16 days, got %v", got)
	}

	// Bergabung setelah period berakhir → 0.
	joinLate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if got := EligibleCalendarDays(joinLate, periodStart, periodEnd); got != 0 {
		t.Errorf("expected 0 days, got %v", got)
	}
}

func TestProrateToWholeNumber(t *testing.T) {
	if got := ProrateToWholeNumber(5454545.45); got != 5454545 {
		t.Errorf("expected 5454545, got %v", got)
	}
}
