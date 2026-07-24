package attendance

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Company Settings
// =========================================================================

func (s *Service) UpsertCompanySetting(ctx context.Context, req CreateCompanySettingRequest) (*CompanySettingResponse, error) {
	setting := &AttendanceCompanySetting{}
	if req.IsLocationRequired != nil {
		setting.IsLocationRequired = *req.IsLocationRequired
	}
	if req.IsFaceRequired != nil {
		setting.IsFaceRequired = *req.IsFaceRequired
	}
	if req.IsOvertimeEnabled != nil {
		setting.IsOvertimeEnabled = *req.IsOvertimeEnabled
	}
	setting.Latitude = req.Latitude
	setting.Longitude = req.Longitude
	setting.MaxDistanceMeter = req.MaxDistanceMeter
	setting.LateToleranceMinutes = req.LateToleranceMinutes
	setting.OvertimeMinMinutes = req.OvertimeMinMinutes

	if err := s.repo.UpsertCompanySetting(ctx, setting); err != nil {
		return nil, err
	}

	s.logger.Info("Company attendance setting updated", zap.String("id", setting.ID.String()))
	return settingToResponse(setting), nil
}

func (s *Service) GetCompanySetting(ctx context.Context) (*CompanySettingResponse, error) {
	setting, err := s.repo.FindCompanySetting(ctx)
	if err != nil {
		return nil, err
	}
	return settingToResponse(setting), nil
}

// =========================================================================
// Company Shifts
// =========================================================================

func (s *Service) CreateShift(ctx context.Context, req CreateCompanyShiftRequest) (*CompanyShiftResponse, error) {
	shift := &AttendanceCompanyShift{
		ShiftName:    req.ShiftName,
		CheckInTime:  req.CheckInTime,
		CheckOutTime: req.CheckOutTime,
	}
	if req.IsCrossMidnight != nil {
		shift.IsCrossMidnight = *req.IsCrossMidnight
	}

	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) GetShiftByID(ctx context.Context, id string) (*CompanyShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	shift, err := s.repo.FindShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) ListShifts(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	shifts, total, err := s.repo.ListShifts(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]CompanyShiftResponse, 0, len(shifts))
	for _, shift := range shifts {
		responses = append(responses, *shiftToResponse(&shift))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateShift(ctx context.Context, id string, req UpdateCompanyShiftRequest) (*CompanyShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	shift, err := s.repo.FindShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.ShiftName != nil {
		shift.ShiftName = *req.ShiftName
	}
	if req.CheckInTime != nil {
		shift.CheckInTime = *req.CheckInTime
	}
	if req.CheckOutTime != nil {
		shift.CheckOutTime = *req.CheckOutTime
	}
	if req.IsCrossMidnight != nil {
		shift.IsCrossMidnight = *req.IsCrossMidnight
	}
	if err := s.repo.UpdateShift(ctx, shift); err != nil {
		return nil, err
	}
	return shiftToResponse(shift), nil
}

func (s *Service) DeleteShift(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteShift(ctx, uid)
}

// =========================================================================
// Employee Shifts
// =========================================================================

func (s *Service) CreateEmployeeShift(ctx context.Context, req CreateEmployeeShiftRequest) (*EmployeeShiftResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	shiftID, err := uuid.Parse(req.AttendanceShiftID)
	if err != nil {
		return nil, fmt.Errorf("invalid attendance_shift_id: %w", err)
	}

	es := &AttendanceEmployeeShift{
		EmployeeID:        empID,
		AttendanceShiftID: shiftID,
		EffectiveDateFrom: req.EffectiveDateFrom,
		EffectiveDateTo:   req.EffectiveDateTo,
		DaysOfWeekMask:    req.DaysOfWeekMask,
		IsDayOff:          req.IsDayOff,
	}
	if err := s.repo.CreateEmployeeShift(ctx, es); err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) GetEmployeeShiftByID(ctx context.Context, id string) (*EmployeeShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	es, err := s.repo.FindEmployeeShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) ListEmployeeShifts(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	shifts, total, err := s.repo.ListEmployeeShifts(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EmployeeShiftResponse, 0, len(shifts))
	for _, es := range shifts {
		responses = append(responses, *employeeShiftToResponse(&es))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateEmployeeShift(ctx context.Context, id string, req UpdateEmployeeShiftRequest) (*EmployeeShiftResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	es, err := s.repo.FindEmployeeShiftByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.AttendanceShiftID != nil {
		sid, err := uuid.Parse(*req.AttendanceShiftID)
		if err != nil {
			return nil, fmt.Errorf("invalid attendance_shift_id: %w", err)
		}
		es.AttendanceShiftID = sid
	}
	if req.EffectiveDateFrom != nil {
		es.EffectiveDateFrom = *req.EffectiveDateFrom
	}
	if req.EffectiveDateTo != nil {
		es.EffectiveDateTo = req.EffectiveDateTo
	}
	if req.DaysOfWeekMask != nil {
		es.DaysOfWeekMask = req.DaysOfWeekMask
	}
	if req.IsDayOff != nil {
		es.IsDayOff = req.IsDayOff
	}
	if err := s.repo.UpdateEmployeeShift(ctx, es); err != nil {
		return nil, err
	}
	return employeeShiftToResponse(es), nil
}

func (s *Service) DeleteEmployeeShift(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeShift(ctx, uid)
}

// =========================================================================
// Locations
// =========================================================================

func (s *Service) CreateLocation(ctx context.Context, req CreateLocationRequest) (*LocationResponse, error) {
	loc := &AttendanceLocation{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		RadiusM:   50,
	}
	if req.RadiusM != nil && *req.RadiusM > 0 {
		loc.RadiusM = *req.RadiusM
	}
	if err := s.repo.CreateLocation(ctx, loc); err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) GetLocationByID(ctx context.Context, id string) (*LocationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	loc, err := s.repo.FindLocationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) ListLocations(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	locs, total, err := s.repo.ListLocations(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]LocationResponse, 0, len(locs))
	for _, loc := range locs {
		responses = append(responses, *locationToResponse(&loc))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateLocation(ctx context.Context, id string, req UpdateLocationRequest) (*LocationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	loc, err := s.repo.FindLocationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		loc.Name = *req.Name
	}
	if req.Latitude != nil {
		loc.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		loc.Longitude = *req.Longitude
	}
	if req.RadiusM != nil && *req.RadiusM > 0 {
		loc.RadiusM = *req.RadiusM
	}
	if err := s.repo.UpdateLocation(ctx, loc); err != nil {
		return nil, err
	}
	return locationToResponse(loc), nil
}

func (s *Service) DeleteLocation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteLocation(ctx, uid)
}

// =========================================================================
// Events (Check-in / Check-out)
// =========================================================================

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest) (*EventResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}

	utcTime, err := time.Parse(time.RFC3339, req.EventTimeUTC)
	if err != nil {
		return nil, fmt.Errorf("invalid event_time_utc: %w", err)
	}
	localTime, err := time.Parse(time.RFC3339, req.EventTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid event_time_local: %w", err)
	}

	event := &AttendanceEvent{
		EmployeeID:       empID,
		EventType:        EventType(req.EventType),
		EventTimeUTC:     utcTime,
		EventTimeLocal:   localTime,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		AccuracyM:        req.AccuracyM,
		LocationProvider: req.LocationProvider,
		ValidationStatus: ValidationPending,
	}
	if req.DeviceID != nil && *req.DeviceID != "" {
		did, err := uuid.Parse(*req.DeviceID)
		if err == nil {
			event.DeviceID = &did
		}
	}

	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, err
	}
	return eventToResponse(event), nil
}

func (s *Service) GetEventByID(ctx context.Context, id string) (*EventResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	event, err := s.repo.FindEventByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return eventToResponse(event), nil
}

func (s *Service) ListEvents(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	events, total, err := s.repo.ListEvents(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, *eventToResponse(&event))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Sessions
// =========================================================================

func (s *Service) GetSession(ctx context.Context, employeeID, workDate string) (*SessionResponse, error) {
	empUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	session, err := s.repo.FindSessionByEmployeeAndDate(ctx, empUUID, workDate)
	if err != nil {
		return nil, err
	}
	return sessionToResponse(session), nil
}

func (s *Service) ListSessions(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	sessions, total, err := s.repo.ListSessions(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, *sessionToResponse(&session))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Overtime Requests
// =========================================================================

func (s *Service) CreateOvertimeRequest(ctx context.Context, req CreateOvertimeRequest) (*OvertimeResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time_local: %w", err)
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTimeLocal)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time_local: %w", err)
	}

	overtime := &AttendanceOvertimeRequest{
		EmployeeID:       empID,
		WorkDate:         req.WorkDate,
		StartTimeLocal:   startTime,
		EndTimeLocal:     endTime,
		RequestedMinutes: req.RequestedMinutes,
		Reason:           strPtr(req.Reason),
		Status:           OvertimeSubmitted,
	}
	if err := s.repo.CreateOvertimeRequest(ctx, overtime); err != nil {
		return nil, err
	}
	return overtimeToResponse(overtime), nil
}

func (s *Service) GetOvertimeRequestByID(ctx context.Context, id string) (*OvertimeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOvertimeRequestByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return overtimeToResponse(o), nil
}

func (s *Service) ListOvertimeRequests(ctx context.Context, employeeID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, err := uuid.Parse(*employeeID)
		if err != nil {
			return nil, fmt.Errorf("invalid employee_id: %w", err)
		}
		empUUID = &uid
	}
	requests, total, err := s.repo.ListOvertimeRequests(ctx, empUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]OvertimeResponse, 0, len(requests))
	for _, o := range requests {
		responses = append(responses, *overtimeToResponse(&o))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Exempt Positions
// =========================================================================

func (s *Service) CreateExemptPosition(ctx context.Context, req CreateExemptPositionRequest) (*ExemptPositionResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	p := &AttendanceExemptPosition{
		OrganizationID: orgID,
		IsExempt:       true,
		Note:           strPtr(req.Note),
	}
	if req.IsExempt != nil {
		p.IsExempt = *req.IsExempt
	}
	if err := s.repo.CreateExemptPosition(ctx, p); err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) GetExemptPositionByID(ctx context.Context, id string) (*ExemptPositionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindExemptPositionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) ListExemptPositions(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	positions, total, err := s.repo.ListExemptPositions(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ExemptPositionResponse, 0, len(positions))
	for _, p := range positions {
		responses = append(responses, *exemptPositionToResponse(&p))
	}
	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateExemptPosition(ctx context.Context, id string, req UpdateExemptPositionRequest) (*ExemptPositionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindExemptPositionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.IsExempt != nil {
		p.IsExempt = *req.IsExempt
	}
	if req.Note != nil {
		p.Note = req.Note
	}
	if err := s.repo.UpdateExemptPosition(ctx, p); err != nil {
		return nil, err
	}
	return exemptPositionToResponse(p), nil
}

func (s *Service) DeleteExemptPosition(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteExemptPosition(ctx, uid)
}

// =========================================================================
// Helpers
// =========================================================================

func calcTotalPages(total int64, perPage int) int {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	if pages < 1 {
		return 1
	}
	return pages
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func settingToResponse(s *AttendanceCompanySetting) *CompanySettingResponse {
	return &CompanySettingResponse{
		ID:                   s.ID.String(),
		Latitude:             s.Latitude,
		Longitude:            s.Longitude,
		IsLocationRequired:   s.IsLocationRequired,
		IsFaceRequired:       s.IsFaceRequired,
		IsOvertimeEnabled:    s.IsOvertimeEnabled,
		MaxDistanceMeter:     s.MaxDistanceMeter,
		LateToleranceMinutes: s.LateToleranceMinutes,
		OvertimeMinMinutes:   s.OvertimeMinMinutes,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}

func shiftToResponse(s *AttendanceCompanyShift) *CompanyShiftResponse {
	return &CompanyShiftResponse{
		ID:              s.ID.String(),
		ShiftName:       s.ShiftName,
		CheckInTime:     s.CheckInTime,
		CheckOutTime:    s.CheckOutTime,
		IsCrossMidnight: s.IsCrossMidnight,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

func employeeShiftToResponse(es *AttendanceEmployeeShift) *EmployeeShiftResponse {
	resp := &EmployeeShiftResponse{
		ID:                es.ID.String(),
		EmployeeID:        es.EmployeeID.String(),
		AttendanceShiftID: es.AttendanceShiftID.String(),
		EffectiveDateFrom: es.EffectiveDateFrom,
		EffectiveDateTo:   es.EffectiveDateTo,
		DaysOfWeekMask:    es.DaysOfWeekMask,
		IsDayOff:          es.IsDayOff,
		CreatedAt:         es.CreatedAt,
		UpdatedAt:         es.UpdatedAt,
	}
	return resp
}

func locationToResponse(l *AttendanceLocation) *LocationResponse {
	return &LocationResponse{
		ID:        l.ID.String(),
		Name:      l.Name,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		RadiusM:   l.RadiusM,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}

func eventToResponse(e *AttendanceEvent) *EventResponse {
	resp := &EventResponse{
		ID:               e.ID.String(),
		EmployeeID:       e.EmployeeID.String(),
		EventType:        string(e.EventType),
		EventTimeUTC:     e.EventTimeUTC,
		EventTimeLocal:   e.EventTimeLocal,
		Latitude:         e.Latitude,
		Longitude:        e.Longitude,
		AccuracyM:        e.AccuracyM,
		LocationProvider: e.LocationProvider,
		IsInGeofence:     e.IsInGeofence,
		ValidationStatus: string(e.ValidationStatus),
		ValidationNote:   e.ValidationNote,
		CreatedAt:        e.CreatedAt,
	}
	if e.DeviceID != nil {
		did := e.DeviceID.String()
		resp.DeviceID = &did
	}
	if e.ValidatedLocationID != nil {
		vlid := e.ValidatedLocationID.String()
		resp.ValidatedLocationID = &vlid
	}
	if e.DistanceM != nil {
		resp.DistanceM = e.DistanceM
	}
	return resp
}

func sessionToResponse(s *AttendanceSession) *SessionResponse {
	resp := &SessionResponse{
		ID:                s.ID.String(),
		EmployeeID:        s.EmployeeID.String(),
		WorkDate:          s.WorkDate,
		IsOvertimeDay:     s.IsOvertimeDay,
		Status:            string(s.Status),
		LatenessMinutes:   s.LatenessMinutes,
		EarlyLeaveMinutes: s.EarlyLeaveMinutes,
		WorkMinutes:       s.WorkMinutes,
		BreakMinutes:      s.BreakMinutes,
		OvertimeMinutes:   s.OvertimeMinutes,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
	if s.ShiftID != nil {
		sid := s.ShiftID.String()
		resp.ShiftID = &sid
	}
	if s.CheckinEventID != nil {
		ceid := s.CheckinEventID.String()
		resp.CheckinEventID = &ceid
	}
	if s.CheckoutEventID != nil {
		ceid := s.CheckoutEventID.String()
		resp.CheckoutEventID = &ceid
	}
	if s.LeaveRequestID != nil {
		lrID := s.LeaveRequestID.String()
		resp.LeaveRequestID = &lrID
	}
	if s.PlannedStartLocal != nil {
		resp.PlannedStartLocal = s.PlannedStartLocal
	}
	if s.PlannedEndLocal != nil {
		resp.PlannedEndLocal = s.PlannedEndLocal
	}
	return resp
}

func overtimeToResponse(o *AttendanceOvertimeRequest) *OvertimeResponse {
	resp := &OvertimeResponse{
		ID:               o.ID.String(),
		EmployeeID:       o.EmployeeID.String(),
		WorkDate:         o.WorkDate,
		StartTimeLocal:   o.StartTimeLocal,
		EndTimeLocal:     o.EndTimeLocal,
		RequestedMinutes: o.RequestedMinutes,
		Reason:           o.Reason,
		Status:           string(o.Status),
		ApprovedAt:       o.ApprovedAt,
		ApprovalNote:     o.ApprovalNote,
		CreatedAt:        o.CreatedAt,
	}
	if o.ApprovedBy != nil {
		ab := o.ApprovedBy.String()
		resp.ApprovedBy = &ab
	}
	return resp
}

func exemptPositionToResponse(p *AttendanceExemptPosition) *ExemptPositionResponse {
	return &ExemptPositionResponse{
		ID:             p.ID.String(),
		OrganizationID: p.OrganizationID.String(),
		IsExempt:       p.IsExempt,
		Note:           p.Note,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
