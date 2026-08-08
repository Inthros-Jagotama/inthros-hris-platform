package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const workDateLayout = "2006-01-02"

// normalizeWorkDate handles the same DB round-trip quirk documented for
// Leave (see leave/balance.go): some drivers (e.g. glebarez/sqlite in tests)
// return a `type:date` column as a full RFC3339 timestamp instead of plain
// "2006-01-02". Reformats to workDateLayout so it matches values that were
// never round-tripped.
func normalizeWorkDate(workDate string) string {
	if t, err := time.Parse(workDateLayout, workDate); err == nil {
		return t.Format(workDateLayout)
	}
	if t, err := time.Parse(time.RFC3339, workDate); err == nil {
		return t.Format(workDateLayout)
	}
	return workDate
}

// recalculateSession is the session calculation engine (§19-24, §43):
// resolves the employee's shift for workDate, reads that window's raw
// events, and creates/updates the corresponding AttendanceSession with
// lateness/early-leave/work-minutes and a status.
//
// It is triggered synchronously after each event is created (see CreateEvent)
// rather than by a scheduled job - ProcessDailyAttendance/DetectMissingAttendance
// (§44-45) remain unimplemented, so a session only exists once at least one
// event has been recorded for that work date. Proactively marking days
// ABSENT before any event happens is that scheduled job's responsibility,
// not this function's.
//
// Exempt-position handling (§28) is intentionally not resolved here - it
// needs the employee's organization_id, which requires a cross-module read
// into the employee module that doesn't exist yet anywhere in this
// codebase (same category of gap as Leave's deferred employee-active
// check). Day-of-week mask interpretation is also not applied, for the
// same reason documented in Phase 3: no bit-to-weekday convention exists.
func (s *Service) recalculateSession(ctx context.Context, employeeID uuid.UUID, workDate string) error {
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, employeeID, workDate)
	if err != nil {
		session = &AttendanceSession{EmployeeID: employeeID, WorkDate: workDate}
	}

	empShift, shiftErr := s.repo.FindEmployeeShiftForDate(ctx, employeeID, workDate)
	if shiftErr != nil {
		// No shift assignment for this date - §25 treats an employee with no
		// working schedule as a day off.
		session.Status = SessionStatusDayOff
		session.ShiftID = nil
		return s.repo.UpsertSession(ctx, session)
	}
	if empShift.IsDayOff != nil && *empShift.IsDayOff {
		session.Status = SessionStatusDayOff
		session.ShiftID = &empShift.AttendanceShiftID
		return s.repo.UpsertSession(ctx, session)
	}

	shift, err := s.repo.FindShiftByID(ctx, empShift.AttendanceShiftID)
	if err != nil {
		return fmt.Errorf("failed to load shift: %w", err)
	}
	session.ShiftID = &shift.ID

	workDateT, err := time.Parse(workDateLayout, workDate)
	if err != nil {
		return fmt.Errorf("invalid work date: %w", err)
	}

	events, err := s.repo.FindEventsForWorkDate(ctx, employeeID, workDate)
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}
	checkin, checkout := selectCheckinCheckout(events)

	// attendance_company_shifts.check_in_time/check_out_time carry no
	// timezone of their own - anchor them to whichever event's local
	// timestamp we have, since event_time_local's offset is this codebase's
	// only source of "local" (there's no company/tenant timezone setting -
	// see docs/module-attendance-plan.md §3, deferred in Phase 1). Falling
	// back to UTC when there are no events yet is harmless: planned times
	// aren't meaningfully comparable to anything in that case anyway.
	loc := time.UTC
	switch {
	case checkin != nil:
		loc = checkin.EventTimeLocal.Location()
	case checkout != nil:
		loc = checkout.EventTimeLocal.Location()
	}

	plannedStart, err := combineDateAndTime(workDateT, shift.CheckInTime, loc)
	if err != nil {
		return fmt.Errorf("invalid shift check_in_time: %w", err)
	}
	endDate := workDateT
	if shift.IsCrossMidnight {
		endDate = workDateT.AddDate(0, 0, 1)
	}
	plannedEnd, err := combineDateAndTime(endDate, shift.CheckOutTime, loc)
	if err != nil {
		return fmt.Errorf("invalid shift check_out_time: %w", err)
	}
	session.PlannedStartLocal = &plannedStart
	session.PlannedEndLocal = &plannedEnd

	lateToleranceMinutes := 0
	if setting, err := s.repo.FindCompanySetting(ctx); err == nil && setting.LateToleranceMinutes != nil {
		lateToleranceMinutes = *setting.LateToleranceMinutes
	}

	session.LatenessMinutes = 0
	session.EarlyLeaveMinutes = 0
	session.WorkMinutes = 0
	session.CheckinEventID = nil
	session.CheckoutEventID = nil

	switch {
	case checkin != nil && checkout != nil:
		session.CheckinEventID = &checkin.ID
		session.CheckoutEventID = &checkout.ID
		session.Status = SessionStatusClosed
		session.LatenessMinutes = minutesLate(plannedStart, checkin.EventTimeLocal, lateToleranceMinutes)
		session.EarlyLeaveMinutes = minutesEarly(plannedEnd, checkout.EventTimeLocal)
		session.WorkMinutes = workMinutesBetween(checkin.EventTimeLocal, checkout.EventTimeLocal, session.BreakMinutes)
	case checkin != nil:
		session.CheckinEventID = &checkin.ID
		session.Status = SessionStatusMissingCheckOut
		session.LatenessMinutes = minutesLate(plannedStart, checkin.EventTimeLocal, lateToleranceMinutes)
	case checkout != nil:
		session.CheckoutEventID = &checkout.ID
		session.Status = SessionStatusMissingCheckIn
		session.EarlyLeaveMinutes = minutesEarly(plannedEnd, checkout.EventTimeLocal)
	default:
		session.Status = SessionStatusOpen
	}

	return s.repo.UpsertSession(ctx, session)
}

// selectCheckinCheckout picks the first CHECKIN and the first CHECKOUT that
// follows it from a work date's events (already ordered by event_time_local
// ascending) - covering the cross-midnight case where the checkout's local
// calendar date is the day after the checkin's.
func selectCheckinCheckout(events []AttendanceEvent) (checkin, checkout *AttendanceEvent) {
	for i := range events {
		e := &events[i]
		if e.EventType == EventTypeCheckIn && checkin == nil {
			checkin = e
			continue
		}
		if e.EventType == EventTypeCheckOut && checkin != nil && checkout == nil {
			checkout = e
		}
	}
	return checkin, checkout
}

// minutesLate returns lateness after tolerance, floored at 0 (§21).
func minutesLate(plannedStart, actual time.Time, toleranceMinutes int) int {
	if !actual.After(plannedStart) {
		return 0
	}
	raw := int(actual.Sub(plannedStart).Minutes())
	late := raw - toleranceMinutes
	if late < 0 {
		return 0
	}
	return late
}

// minutesEarly returns how many minutes before plannedEnd the checkout
// happened, floored at 0 (§22).
func minutesEarly(plannedEnd, actual time.Time) int {
	if !actual.Before(plannedEnd) {
		return 0
	}
	return int(plannedEnd.Sub(actual).Minutes())
}

// workMinutesBetween returns checkout-checkin minus breakMinutes, floored at
// 0 (§23).
func workMinutesBetween(checkin, checkout time.Time, breakMinutes int) int {
	total := int(checkout.Sub(checkin).Minutes()) - breakMinutes
	if total < 0 {
		return 0
	}
	return total
}

// combineDateAndTime combines a calendar date with a "HH:MM" or "HH:MM:SS"
// time-of-day string (as stored in attendance_company_shifts) into a single
// time.Time on that date, in the given location.
func combineDateAndTime(date time.Time, hhmm string, loc *time.Location) (time.Time, error) {
	layout := "15:04"
	if len(hhmm) > 5 {
		layout = "15:04:05"
	}
	t, err := time.Parse(layout, hhmm)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, loc), nil
}
