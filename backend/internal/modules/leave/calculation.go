package leave

import (
	"fmt"
	"time"
)

// defaultHoursPerDay converts hourly leave into a fraction of a day when no
// per-company "hours_per_day" configuration exists yet (see §15 of
// docs/module-leave-plan.md — flagged as a future config point).
const defaultHoursPerDay = 8.0

const dateLayout = "2006-01-02"
const timeLayout = "15:04"

// normalizeLeaveDate handles the same DB round-trip quirk noted in
// balance.go's requestYear: some drivers (e.g. glebarez/sqlite in tests)
// return a `type:date` column as a full RFC3339 timestamp instead of plain
// "2006-01-02". Reformats to dateLayout so it matches values that were
// never round-tripped.
func normalizeLeaveDate(leaveDate string) string {
	if t, err := time.Parse(dateLayout, leaveDate); err == nil {
		return t.Format(dateLayout)
	}
	if t, err := time.Parse(time.RFC3339, leaveDate); err == nil {
		return t.Format(dateLayout)
	}
	return leaveDate
}

// LeaveDayDetail is one row of the per-date breakdown produced by
// CalculateLeaveDuration, used to populate LeaveRequestDetail.
type LeaveDayDetail struct {
	Date        string
	DayFraction float64
}

// CalculateLeaveDuration computes the authoritative RequestedDays for a leave
// request server-side (§13-15 of the plan), instead of trusting the raw
// client-supplied value. Weekends (Saturday/Sunday) and dates present in
// holidayDates are treated as non-working and excluded.
//
// Per-employee shift/working-calendar overrides (attendance module's
// DaysOfWeekMask) are intentionally NOT considered here — that bitmask has no
// documented weekday interpretation anywhere in the codebase yet, so
// inventing one inside the Leave module would be a cross-module design
// decision, not a Phase 3 calculation-engine concern. Standard Sat/Sun +
// company holidays is the deferred-but-documented baseline until that
// convention exists.
func CalculateLeaveDuration(startDate, endDate string, mode DurationMode, startTime, endTime *string, holidayDates map[string]bool) (float64, []LeaveDayDetail, error) {
	start, err := time.Parse(dateLayout, startDate)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid request_start_date: %w", err)
	}
	end, err := time.Parse(dateLayout, endDate)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid request_end_date: %w", err)
	}
	if end.Before(start) {
		return 0, nil, fmt.Errorf("request_end_date must not be before request_start_date")
	}

	switch mode {
	case DurationHalfDayAM, DurationHalfDayPM:
		if !start.Equal(end) {
			return 0, nil, fmt.Errorf("half day leave must have the same start and end date")
		}
		if isNonWorkingDay(start, holidayDates) {
			return 0, nil, fmt.Errorf("%s is not a working day", startDate)
		}
		return 0.5, []LeaveDayDetail{{Date: startDate, DayFraction: 0.5}}, nil

	case DurationHourly:
		if !start.Equal(end) {
			return 0, nil, fmt.Errorf("hourly leave must have the same start and end date")
		}
		if startTime == nil || endTime == nil {
			return 0, nil, fmt.Errorf("start_time and end_time are required for hourly leave")
		}
		if isNonWorkingDay(start, holidayDates) {
			return 0, nil, fmt.Errorf("%s is not a working day", startDate)
		}
		st, err := time.Parse(timeLayout, *startTime)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid start_time: %w", err)
		}
		et, err := time.Parse(timeLayout, *endTime)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid end_time: %w", err)
		}
		if !et.After(st) {
			return 0, nil, fmt.Errorf("end_time must be after start_time")
		}
		hours := et.Sub(st).Hours()
		days := hours / defaultHoursPerDay
		return days, []LeaveDayDetail{{Date: startDate, DayFraction: days}}, nil

	default: // DurationFullDay
		var total float64
		var details []LeaveDayDetail
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			if isNonWorkingDay(d, holidayDates) {
				continue
			}
			details = append(details, LeaveDayDetail{Date: d.Format(dateLayout), DayFraction: 1.0})
			total += 1.0
		}
		if len(details) == 0 {
			return 0, nil, fmt.Errorf("requested date range contains no working days")
		}
		return total, details, nil
	}
}

func isNonWorkingDay(d time.Time, holidayDates map[string]bool) bool {
	if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		return true
	}
	return holidayDates[d.Format(dateLayout)]
}
