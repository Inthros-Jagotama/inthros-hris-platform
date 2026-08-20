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

	// Events are loaded up-front so the day-off branches below can reflect
	// what actually happened instead of unconditionally short-circuiting.
	events, err := s.repo.FindEventsForWorkDate(ctx, employeeID, workDate)
	if err != nil {
		return fmt.Errorf("failed to load events: %w", err)
	}
	checkin, checkout := selectCheckinCheckout(events)

	// Company setting: apakah check-in/check-out boleh dilakukan pada hari
	// libur (tanpa jadwal shift / IsDayOff). Default = diizinkan bila belum
	// ada baris setting (migration 075 default 1).
	allowOnDayOff := true
	if setting, err := s.repo.FindCompanySetting(ctx); err == nil {
		allowOnDayOff = setting.AllowCheckinOnDayOff
	}
	hasEvents := checkin != nil || checkout != nil

	if shiftErr != nil {
		// No shift assignment for this date - §25 treats an employee with no
		// working schedule as a day off. Saat check-in di hari libur
		// diizinkan dan karyawan benar-benar mencatat event, cerminkan event
		// tersebut (mis. MISSING_CHECKOUT) alih-alih memaksa DAY_OFF.
		session.ShiftID = nil
		if !allowOnDayOff || !hasEvents {
			session.Status = SessionStatusDayOff
			return s.repo.UpsertSession(ctx, session)
		}
		return s.repo.UpsertSession(ctx, finishSessionFromEvents(session, checkin, checkout))
	}
	if empShift.IsDayOff != nil && *empShift.IsDayOff {
		session.ShiftID = &empShift.AttendanceShiftID
		if !allowOnDayOff || !hasEvents {
			session.Status = SessionStatusDayOff
			return s.repo.UpsertSession(ctx, session)
		}
		// Diizinkan + ada event → lanjut ke perhitungan normal di bawah.
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

	// attendance_company_shifts.check_in_time/check_out_time carry no
	// timezone of their own - anchor them to the employee's resolved
	// organization timezone (company default, or zone override) when
	// resolvable, so lateness/early-leave reflect the org's official zone
	// rather than whatever offset the client device chose to report.
	// Falls back to the client-embedded EventTimeLocal offset when
	// org/timezone data is incomplete (must never block session
	// generation), and to UTC when there are no events yet at all.
	loc := time.UTC
	if empLoc := s.resolveEmployeeTimezone(ctx, employeeID); empLoc != nil {
		loc = empLoc
	} else {
		switch {
		case checkin != nil:
			loc = checkin.EventTimeLocal.Location()
		case checkout != nil:
			loc = checkout.EventTimeLocal.Location()
		}
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

// finishSessionFromEvents sets a session's status/event-IDs/work-minutes from
// its raw events when there is no shift to compute planned times against
// (day-off with check-in allowed, no shift assignment). Lateness/early-leave
// stay 0 — there are no planned times to compare with.
func finishSessionFromEvents(session *AttendanceSession, checkin, checkout *AttendanceEvent) *AttendanceSession {
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
		session.WorkMinutes = workMinutesBetween(checkin.EventTimeLocal, checkout.EventTimeLocal, session.BreakMinutes)
	case checkin != nil:
		session.CheckinEventID = &checkin.ID
		session.Status = SessionStatusMissingCheckOut
	case checkout != nil:
		session.CheckoutEventID = &checkout.ID
		session.Status = SessionStatusMissingCheckIn
	default:
		session.Status = SessionStatusOpen
	}
	return session
}

// ApplyApprovedLeave lets the Leave module push an approved leave onto this
// day's session (§26/§27/§50) - implements leave.AttendanceSessionUpdater,
// wired in main.go the same push-based way ApprovalEngine callbacks work.
// It only overwrites the session status to LEAVE when the day isn't already
// CLOSED (a real, completed attendance day): a half-day leave coexisting
// with half-day attendance just records LeaveFraction/LeaveRequestID
// without discarding real work data, matching §27's worked example.
func (s *Service) ApplyApprovedLeave(ctx context.Context, employeeID uuid.UUID, workDate string, leaveRequestID uuid.UUID, dayFraction float64) error {
	workDate = normalizeWorkDate(workDate)
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, employeeID, workDate)
	if err != nil {
		session = &AttendanceSession{EmployeeID: employeeID, WorkDate: workDate}
	}
	session.LeaveRequestID = &leaveRequestID
	fraction := dayFraction
	session.LeaveFraction = &fraction
	if session.Status != SessionStatusClosed {
		session.Status = SessionStatusLeave
	}
	return s.repo.UpsertSession(ctx, session)
}

// ApplyApprovedBusinessTravel lets Business Travel push an APPROVED travel
// onto a day's session (§37 plan doc), mirroring ApplyApprovedLeave's push
// pattern rather than a generic source_type/source_id design. Same
// CLOSED-day protection: a real, completed attendance day isn't overwritten,
// so a half day of actual work coexisting with travel just records
// BusinessTravelID without discarding it.
func (s *Service) ApplyApprovedBusinessTravel(ctx context.Context, employeeID uuid.UUID, workDate string, businessTravelID uuid.UUID) error {
	workDate = normalizeWorkDate(workDate)
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, employeeID, workDate)
	if err != nil {
		session = &AttendanceSession{EmployeeID: employeeID, WorkDate: workDate}
	}
	session.BusinessTravelID = &businessTravelID
	if session.Status != SessionStatusClosed {
		session.Status = SessionStatusBusinessTravel
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
