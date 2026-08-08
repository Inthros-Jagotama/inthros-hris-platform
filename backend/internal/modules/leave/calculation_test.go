package leave

import "testing"

func TestCalculateLeaveDuration_FullDay_ExcludesWeekend(t *testing.T) {
	// 2026-01-15 Thu .. 2026-01-19 Mon, no holidays -> excludes Sat 17th/Sun 18th = 3 days
	days, details, err := CalculateLeaveDuration("2026-01-15", "2026-01-19", DurationFullDay, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if days != 3 {
		t.Errorf("expected 3 days, got %.2f", days)
	}
	if len(details) != 3 {
		t.Errorf("expected 3 detail rows, got %d", len(details))
	}
}

func TestCalculateLeaveDuration_FullDay_ExcludesHoliday(t *testing.T) {
	holidays := map[string]bool{"2026-01-16": true}
	days, _, err := CalculateLeaveDuration("2026-01-15", "2026-01-16", DurationFullDay, nil, nil, holidays)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if days != 1 {
		t.Errorf("expected 1 day (16th excluded as holiday), got %.2f", days)
	}
}

func TestCalculateLeaveDuration_FullDay_AllNonWorking_ReturnsError(t *testing.T) {
	// 2026-01-17 Sat .. 2026-01-18 Sun
	_, _, err := CalculateLeaveDuration("2026-01-17", "2026-01-18", DurationFullDay, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for range with no working days")
	}
}

func TestCalculateLeaveDuration_HalfDay(t *testing.T) {
	days, details, err := CalculateLeaveDuration("2026-01-15", "2026-01-15", DurationHalfDayAM, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if days != 0.5 {
		t.Errorf("expected 0.5 days, got %.2f", days)
	}
	if len(details) != 1 || details[0].DayFraction != 0.5 {
		t.Errorf("unexpected details: %+v", details)
	}
}

func TestCalculateLeaveDuration_HalfDay_OnWeekend_ReturnsError(t *testing.T) {
	_, _, err := CalculateLeaveDuration("2026-01-17", "2026-01-17", DurationHalfDayAM, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for half day on weekend")
	}
}

func TestCalculateLeaveDuration_Hourly(t *testing.T) {
	start := "09:00"
	end := "13:00"
	days, details, err := CalculateLeaveDuration("2026-01-15", "2026-01-15", DurationHourly, &start, &end, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if days != 0.5 {
		t.Errorf("expected 0.5 days (4h / 8h), got %.2f", days)
	}
	if len(details) != 1 {
		t.Errorf("expected 1 detail row, got %d", len(details))
	}
}

func TestCalculateLeaveDuration_Hourly_MissingTimes_ReturnsError(t *testing.T) {
	_, _, err := CalculateLeaveDuration("2026-01-15", "2026-01-15", DurationHourly, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error when start_time/end_time are missing")
	}
}

func TestCalculateLeaveDuration_EndBeforeStart_ReturnsError(t *testing.T) {
	_, _, err := CalculateLeaveDuration("2026-01-16", "2026-01-15", DurationFullDay, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error when end date is before start date")
	}
}
