package attendance

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

// =========================================================================
// Company Settings
// =========================================================================

func (r *Repository) UpsertCompanySetting(ctx context.Context, s *AttendanceCompanySetting) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	// Only one setting row per tenant — upsert
	var existing AttendanceCompanySetting
	if err := db.First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return db.Create(s).Error
		}
		return err
	}
	s.ID = existing.ID
	return db.Save(s).Error
}

func (r *Repository) FindCompanySetting(ctx context.Context) (*AttendanceCompanySetting, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s AttendanceCompanySetting
	if err := db.First(&s).Error; err != nil {
		return nil, fmt.Errorf("company setting not found: %w", err)
	}
	return &s, nil
}

// =========================================================================
// Company Shifts
// =========================================================================

func (r *Repository) CreateShift(ctx context.Context, s *AttendanceCompanyShift) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(s).Error
}

func (r *Repository) FindShiftByID(ctx context.Context, id uuid.UUID) (*AttendanceCompanyShift, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s AttendanceCompanyShift
	if err := db.First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("shift not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) ListShifts(ctx context.Context, page, perPage int) ([]AttendanceCompanyShift, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var shifts []AttendanceCompanyShift
	var total int64
	query := db.Model(&AttendanceCompanyShift{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("shift_name ASC").Find(&shifts).Error; err != nil {
		return nil, 0, err
	}
	return shifts, total, nil
}

func (r *Repository) UpdateShift(ctx context.Context, s *AttendanceCompanyShift) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

func (r *Repository) DeleteShift(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&AttendanceCompanyShift{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("shift not found")
	}
	return nil
}

// =========================================================================
// Employee Shifts
// =========================================================================

func (r *Repository) CreateEmployeeShift(ctx context.Context, es *AttendanceEmployeeShift) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(es).Error
}

func (r *Repository) FindEmployeeShiftByID(ctx context.Context, id uuid.UUID) (*AttendanceEmployeeShift, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var es AttendanceEmployeeShift
	if err := db.First(&es, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("employee shift not found: %w", err)
	}
	return &es, nil
}

func (r *Repository) ListEmployeeShifts(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]AttendanceEmployeeShift, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var shifts []AttendanceEmployeeShift
	var total int64
	query := db.Model(&AttendanceEmployeeShift{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("effective_date_from DESC").Find(&shifts).Error; err != nil {
		return nil, 0, err
	}
	return shifts, total, nil
}

// CountOverlappingEmployeeShifts counts the employee's other shift
// assignments whose [effective_date_from, effective_date_to] range overlaps
// [from, to] (a nil `to` means open-ended / ongoing). excludeID lets an
// update check overlap against everything except the row being updated.
func (r *Repository) CountOverlappingEmployeeShifts(ctx context.Context, employeeID uuid.UUID, from string, to *string, excludeID *uuid.UUID) (int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return 0, err
	}
	query := db.Model(&AttendanceEmployeeShift{}).Where("employee_id = ?", employeeID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if to != nil {
		query = query.Where("effective_date_from <= ?", *to)
	}
	query = query.Where("effective_date_to IS NULL OR effective_date_to >= ?", from)

	var count int64
	err = query.Count(&count).Error
	return count, err
}

func (r *Repository) UpdateEmployeeShift(ctx context.Context, es *AttendanceEmployeeShift) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(es).Error
}

func (r *Repository) DeleteEmployeeShift(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&AttendanceEmployeeShift{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee shift not found")
	}
	return nil
}

// =========================================================================
// Locations
// =========================================================================

func (r *Repository) CreateLocation(ctx context.Context, l *AttendanceLocation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(l).Error
}

func (r *Repository) FindLocationByID(ctx context.Context, id uuid.UUID) (*AttendanceLocation, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var l AttendanceLocation
	if err := db.First(&l, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return &l, nil
}

// FindAllLocations returns every registered geofence location, unpaginated,
// for use by event-time location validation (see geofence.go).
func (r *Repository) FindAllLocations(ctx context.Context) ([]AttendanceLocation, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var locs []AttendanceLocation
	if err := db.Find(&locs).Error; err != nil {
		return nil, err
	}
	return locs, nil
}

func (r *Repository) ListLocations(ctx context.Context, page, perPage int) ([]AttendanceLocation, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var locs []AttendanceLocation
	var total int64
	query := db.Model(&AttendanceLocation{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("name ASC").Find(&locs).Error; err != nil {
		return nil, 0, err
	}
	return locs, total, nil
}

func (r *Repository) UpdateLocation(ctx context.Context, l *AttendanceLocation) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(l).Error
}

func (r *Repository) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&AttendanceLocation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("location not found")
	}
	return nil
}

// =========================================================================
// Events
// =========================================================================

func (r *Repository) CreateEvent(ctx context.Context, e *AttendanceEvent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(e).Error
}

func (r *Repository) FindEventByID(ctx context.Context, id uuid.UUID) (*AttendanceEvent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var e AttendanceEvent
	if err := db.First(&e, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}
	return &e, nil
}

func (r *Repository) ListEvents(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]AttendanceEvent, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var events []AttendanceEvent
	var total int64
	query := db.Model(&AttendanceEvent{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("event_time_utc DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

// FindLastEventForEmployee returns the employee's most recent attendance
// event (by event_time_utc), used to detect duplicate/out-of-sequence
// check-in or check-out submissions. Returns gorm.ErrRecordNotFound wrapped
// when the employee has no events yet.
func (r *Repository) FindLastEventForEmployee(ctx context.Context, employeeID uuid.UUID) (*AttendanceEvent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var e AttendanceEvent
	if err := db.Where("employee_id = ?", employeeID).Order("event_time_utc DESC").First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, e *AttendanceEvent) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(e).Error
}

// =========================================================================
// Sessions
// =========================================================================

func (r *Repository) FindSessionByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, workDate string) (*AttendanceSession, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s AttendanceSession
	if err := db.Where("employee_id = ? AND work_date = ?", employeeID, workDate).First(&s).Error; err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return &s, nil
}

func (r *Repository) ListSessions(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]AttendanceSession, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var sessions []AttendanceSession
	var total int64
	query := db.Model(&AttendanceSession{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("work_date DESC").Find(&sessions).Error; err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

// =========================================================================
// Overtime Requests
// =========================================================================

func (r *Repository) CreateOvertimeRequest(ctx context.Context, o *AttendanceOvertimeRequest) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(o).Error
}

func (r *Repository) FindOvertimeRequestByID(ctx context.Context, id uuid.UUID) (*AttendanceOvertimeRequest, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var o AttendanceOvertimeRequest
	if err := db.First(&o, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("overtime request not found: %w", err)
	}
	return &o, nil
}

func (r *Repository) ListOvertimeRequests(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]AttendanceOvertimeRequest, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var requests []AttendanceOvertimeRequest
	var total int64
	query := db.Model(&AttendanceOvertimeRequest{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

func (r *Repository) UpdateOvertimeRequest(ctx context.Context, o *AttendanceOvertimeRequest) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(o).Error
}

// =========================================================================
// Exempt Positions
// =========================================================================

func (r *Repository) CreateExemptPosition(ctx context.Context, p *AttendanceExemptPosition) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(p).Error
}

func (r *Repository) FindExemptPositionByID(ctx context.Context, id uuid.UUID) (*AttendanceExemptPosition, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p AttendanceExemptPosition
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("exempt position not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) FindExemptPositionByOrgID(ctx context.Context, orgID uuid.UUID) (*AttendanceExemptPosition, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var p AttendanceExemptPosition
	if err := db.Where("organization_id = ?", orgID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("exempt position not found: %w", err)
	}
	return &p, nil
}

func (r *Repository) ListExemptPositions(ctx context.Context, page, perPage int) ([]AttendanceExemptPosition, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var positions []AttendanceExemptPosition
	var total int64
	query := db.Model(&AttendanceExemptPosition{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Find(&positions).Error; err != nil {
		return nil, 0, err
	}
	return positions, total, nil
}

func (r *Repository) UpdateExemptPosition(ctx context.Context, p *AttendanceExemptPosition) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(p).Error
}

func (r *Repository) DeleteExemptPosition(ctx context.Context, id uuid.UUID) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	result := db.Where("id = ?", id).Delete(&AttendanceExemptPosition{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("exempt position not found")
	}
	return nil
}
