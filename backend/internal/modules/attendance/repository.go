package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/timezone"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
	platformDB *gorm.DB
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

// NewRepositoryWithPlatformDB is like NewRepository but also wires a
// platform-DB handle, following the pattern established in
// setting.Repository (dbResolver for the tenant DB, platformDB for
// cross-tenant platform tables such as companies). Used by
// ResolveOrganizationTimezone since companies.timezone lives in the
// platform DB, not the tenant DB.
func NewRepositoryWithPlatformDB(dbResolver func(ctx context.Context) (*gorm.DB, error), platformDB *gorm.DB) *Repository {
	return &Repository{dbResolver: dbResolver, platformDB: platformDB}
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
	// Save menulis semua kolom — pertahankan created_at dari baris existing
	// (CreatedAt pada objek request adalah zero time → MySQL menolak '0000-00-00').
	// UpdatedAt diisi otomatis oleh GORM (autoUpdateTime).
	s.CreatedAt = existing.CreatedAt
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

func (r *Repository) FindSessionByID(ctx context.Context, id uuid.UUID) (*AttendanceSession, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var s AttendanceSession
	if err := db.First(&s, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return &s, nil
}

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

// UpsertSession creates or updates a session by its (employee_id, work_date)
// unique key, used by the session calculation engine (session.go).
func (r *Repository) UpsertSession(ctx context.Context, s *AttendanceSession) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(s).Error
}

// FindEmployeeShiftForDate returns the employee's active shift assignment
// covering workDate, if any (most recently effective one wins).
func (r *Repository) FindEmployeeShiftForDate(ctx context.Context, employeeID uuid.UUID, workDate string) (*AttendanceEmployeeShift, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var es AttendanceEmployeeShift
	err = db.Where("employee_id = ? AND effective_date_from <= ?", employeeID, workDate).
		Where("effective_date_to IS NULL OR effective_date_to >= ?", workDate).
		Order("effective_date_from DESC").
		First(&es).Error
	if err != nil {
		return nil, err
	}
	return &es, nil
}

// FindEventsForWorkDate returns the employee's events whose event_time_local
// falls within [workDate 00:00, workDate+1 24:00) local time - a wide-enough
// window to catch a cross-midnight shift's checkout, which lands on the
// following calendar day (§24).
func (r *Repository) FindEventsForWorkDate(ctx context.Context, employeeID uuid.UUID, workDate string) ([]AttendanceEvent, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	start, err := time.Parse("2006-01-02", workDate)
	if err != nil {
		return nil, fmt.Errorf("invalid work date: %w", err)
	}
	end := start.AddDate(0, 0, 2)
	var events []AttendanceEvent
	err = db.Where("employee_id = ? AND event_time_local >= ? AND event_time_local < ?", employeeID, start, end).
		Order("event_time_local ASC").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

// FindSessionsForEmployeeInRange returns an employee's sessions with
// work_date in [fromDate, toDate] inclusive, ordered chronologically - used
// by the Employee Dashboard/Calendar (§36-37, Phase 10).
func (r *Repository) FindSessionsForEmployeeInRange(ctx context.Context, employeeID uuid.UUID, fromDate, toDate string) ([]AttendanceSession, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var sessions []AttendanceSession
	err = db.Where("employee_id = ? AND work_date >= ? AND work_date <= ?", employeeID, fromDate, toDate).
		Order("work_date ASC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// FindSessionsInRange returns every employee's sessions with work_date in
// [fromDate, toDate] inclusive, tenant-wide (no employee filter) - used by
// the attendance reports (§40, Phase 11). Unlike Manager/HR Dashboard (Phase
// 10, deferred), a tenant-wide report doesn't need to resolve "who reports
// to whom" - it's every session in the date range, period.
func (r *Repository) FindSessionsInRange(ctx context.Context, fromDate, toDate string) ([]AttendanceSession, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var sessions []AttendanceSession
	err = db.Where("work_date >= ? AND work_date <= ?", fromDate, toDate).
		Order("work_date ASC, employee_id ASC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
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

// OvertimeStats — agregat pengajuan lembur dalam rentang tanggal (mode HR
// dashboard). Minutes hanya dijumlahkan untuk status APPROVED (calculated).
type OvertimeStats struct {
	Total    int
	Pending  int
	Approved int
	Minutes  int
}

func (r *Repository) CountOvertimeStats(ctx context.Context, fromDate, toDate string) (*OvertimeStats, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []struct {
		Status  string
		Count   int
		Minutes int
	}
	err = db.Model(&AttendanceOvertimeRequest{}).
		Select("status, COUNT(*) AS count, COALESCE(SUM(calculated_minutes), 0) AS minutes").
		Where("work_date >= ? AND work_date <= ?", fromDate, toDate).
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	stats := &OvertimeStats{}
	for _, row := range rows {
		stats.Total += row.Count
		switch OvertimeStatus(row.Status) {
		case OvertimePendingApproval:
			stats.Pending += row.Count
		case OvertimeApproved:
			stats.Approved += row.Count
			stats.Minutes += row.Minutes
		}
	}
	return stats, nil
}

// OvertimeTrendRow — lembur per tanggal (dipakai service untuk mengagregasi
// per minggu; total & menit hanya dari pengajuan APPROVED).
type OvertimeTrendRow struct {
	WorkDate string
	Count    int
	Approved int
	Minutes  int
}

// OvertimeTrendByDate mengembalikan lembur per work_date dalam rentang
// (tenant-wide), diurutkan naik. Service mengagregasinya menjadi per-minggu.
func (r *Repository) OvertimeTrendByDate(ctx context.Context, fromDate, toDate string) ([]OvertimeTrendRow, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []OvertimeTrendRow
	err = db.Model(&AttendanceOvertimeRequest{}).
		Select("work_date, COUNT(*) AS count, COALESCE(SUM(CASE WHEN status = '" + string(OvertimeApproved) + "' THEN 1 ELSE 0 END), 0) AS approved, COALESCE(SUM(CASE WHEN status = '" + string(OvertimeApproved) + "' THEN calculated_minutes ELSE 0 END), 0) AS minutes").
		Where("work_date >= ? AND work_date <= ?", fromDate, toDate).
		Group("work_date").
		Order("work_date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
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

// =========================================================================
// Correction Requests
// =========================================================================

func (r *Repository) CreateCorrectionRequest(ctx context.Context, c *AttendanceCorrectionRequest) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(c).Error
}

func (r *Repository) FindCorrectionRequestByID(ctx context.Context, id uuid.UUID) (*AttendanceCorrectionRequest, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var c AttendanceCorrectionRequest
	if err := db.First(&c, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("correction request not found: %w", err)
	}
	return &c, nil
}

func (r *Repository) ListCorrectionRequests(ctx context.Context, employeeID *uuid.UUID, page, perPage int) ([]AttendanceCorrectionRequest, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	var requests []AttendanceCorrectionRequest
	var total int64
	query := db.Model(&AttendanceCorrectionRequest{})
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

func (r *Repository) UpdateCorrectionRequest(ctx context.Context, c *AttendanceCorrectionRequest) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Save(c).Error
}

// employeeInfo carries the display fields an admin-facing request list needs
// about who submitted it (name, current organization).
type employeeInfo struct {
	Name             string
	OrganizationName string
}

// GetEmployeeInfoByIDs resolves display info (name, current organization
// name) for a batch of employees at once — used to enrich admin-facing
// overtime/correction request lists without an N+1 query per row. Raw
// query — attendance must not import the employee/organization packages
// directly (same reasoning as approval.GetSubmitterInfoByUserIDs).
func (r *Repository) GetEmployeeInfoByIDs(ctx context.Context, employeeIDs []uuid.UUID) (map[string]employeeInfo, error) {
	result := make(map[string]employeeInfo, len(employeeIDs))
	if len(employeeIDs) == 0 {
		return result, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		EmployeeID       string
		Name             string
		OrganizationName string
	}
	var rows []row
	err = db.Table("employees AS e").
		Joins("JOIN employments AS em ON em.employee_id = e.id AND em.effective_end_date IS NULL").
		Joins("JOIN organizations AS o ON o.id = em.organization_id").
		Where("e.id IN ?", employeeIDs).
		Select("e.id AS employee_id, e.name AS name, o.nomenclature AS organization_name").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee info: %w", err)
	}
	for _, rrow := range rows {
		result[rrow.EmployeeID] = employeeInfo{Name: rrow.Name, OrganizationName: rrow.OrganizationName}
	}
	return result, nil
}

// FindUserIDByEmployeeID resolves an employee's platform user_id via
// employee_accounts, mirroring the employee_id -> user_id resolution already
// used by the approval module (GetUserIDsByOrganization) and by
// leave.Repository.FindUserIDByEmployeeID. Returns nil if the employee has
// no linked user account.
func (r *Repository) FindUserIDByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var userIDStrs []string
	err = db.Table("employee_accounts").
		Where("employee_id = ?", employeeID).
		Limit(1).
		Pluck("user_id", &userIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve employee user id: %w", err)
	}
	if len(userIDStrs) == 0 || userIDStrs[0] == "" {
		return nil, nil
	}
	userID, err := uuid.Parse(userIDStrs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	return &userID, nil
}

// FindEmployeeIDByUserID resolves the employee linked to a platform user
// (reverse of FindUserIDByEmployeeID) — used to attribute ASSIGNED overtime
// requests (assigned_by) and to validate ownership when filling the actual
// overtime (§32b dua-alur). Returns nil if the user has no linked employee.
func (r *Repository) FindEmployeeIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var employeeIDStrs []string
	err = db.Table("employee_accounts").
		Where("user_id = ?", userID).
		Limit(1).
		Pluck("employee_id", &employeeIDStrs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user employee id: %w", err)
	}
	if len(employeeIDStrs) == 0 || employeeIDStrs[0] == "" {
		return nil, nil
	}
	employeeID, err := uuid.Parse(employeeIDStrs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}
	return &employeeID, nil
}

// =========================================================================
// Assignable Employees (§32b alur ASSIGNED) — resolusi hirarki organisasi
// =========================================================================
//
// Resolusi "bawahan yang bisa ditugaskan" memakai pola yang sama dengan
// performance.GetEffectiveChildOrganizationIDs: turun dari organisasi user
// yang login, anak langsung yang terisi (ada employment aktif) ikut dipilih;
// anak langsung yang KOSONG dilewati dan pencarian turun ke level berikutnya
// sampai menemukan organisasi terisi — "bawahan langsung, atau turun level
// sampai ada karyawannya". Raw table queries (tanpa import package
// organization/employee — hindari circular dependency, pola sama dengan
// GetEmployeeInfoByIDs).

// FindOrganizationIDByUserID resolves the current (open-ended) organization
// of the employee linked to a platform user — join employee_accounts →
// employments (effective_end_date IS NULL), pola sama dengan
// performance.okrRepositoryImpl.GetCurrentEmployeeContext. Returns nil if the
// user has no linked employee or no active employment.
func (r *Repository) FindOrganizationIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	type row struct {
		OrganizationID string
	}
	var result row
	err = db.Table("employee_accounts AS ea").
		Joins("JOIN employments AS emp ON emp.employee_id = ea.employee_id").
		Where("ea.user_id = ? AND emp.effective_end_date IS NULL", userID).
		Order("emp.effective_date DESC").
		Limit(1).
		Select("emp.organization_id AS organization_id").
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user organization: %w", err)
	}
	if result.OrganizationID == "" {
		return nil, nil
	}
	orgID, err := uuid.Parse(result.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization id: %w", err)
	}
	return &orgID, nil
}

// getCompanyTimezone reads the current tenant's company timezone from the
// platform DB (companies.timezone), keyed by the company ID carried in ctx
// by middleware.TenantRequired (authctx.GetCompanyID). Mirrors
// setting.Repository.FindCompanyTimezone's platform-DB access pattern.
func (r *Repository) getCompanyTimezone(ctx context.Context) (string, error) {
	if r.platformDB == nil {
		return "", fmt.Errorf("platform database not configured")
	}
	companyID := authctx.GetCompanyID(ctx)
	if companyID == "" {
		return "", fmt.Errorf("company id not found in context")
	}
	var c struct {
		Timezone string
	}
	if err := r.platformDB.Table("companies").Select("timezone").Where("id = ?", companyID).First(&c).Error; err != nil {
		return "", fmt.Errorf("failed to resolve company timezone: %w", err)
	}
	return c.Timezone, nil
}

// ResolveOrganizationTimezone mengembalikan zona waktu efektif untuk sebuah
// organization: Zone.timezone (jika di-set) mengalahkan Company.timezone.
func (r *Repository) ResolveOrganizationTimezone(ctx context.Context, organizationID uuid.UUID) (*time.Location, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}

	var zoneTz *string
	err = db.Table("organizations AS o").
		Joins("LEFT JOIN zones AS z ON z.id = o.zone_id").
		Where("o.id = ?", organizationID).
		Select("z.timezone AS zone_tz").
		Scan(&zoneTz).Error
	if err != nil {
		return nil, err
	}

	companyTz, err := r.getCompanyTimezone(ctx)
	if err != nil {
		return nil, err
	}

	return timezone.Resolve(companyTz, zoneTz)
}

// GetChildOrganizationIDs returns the direct children of an organization.
func (r *Repository) GetChildOrganizationIDs(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var ids []uuid.UUID
	if err := db.Table("organizations").
		Where("parent_id = ? AND deleted_at IS NULL", orgID).
		Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf("failed to resolve child organizations: %w", err)
	}
	return ids, nil
}

// IsOrganizationOccupied reports whether an organization currently has an
// active employment (occupied position).
func (r *Repository) IsOrganizationOccupied(ctx context.Context, orgID uuid.UUID) (bool, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return false, err
	}
	var count int64
	if err := db.Table("employments").
		Where("organization_id = ? AND effective_end_date IS NULL", orgID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check organization occupancy: %w", err)
	}
	return count > 0, nil
}

// GetEffectiveChildOrganizationIDs returns the organizations whose employees
// can be assigned overtime by a member of orgID: direct children that are
// occupied, plus — when a direct child is vacant — the effective subordinates
// below that vacant child (walking down past every vacant organization until
// an occupied one is found). Mirrors performance.GetEffectiveChildOrganizationIDs.
func (r *Repository) GetEffectiveChildOrganizationIDs(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error) {
	var effective []uuid.UUID
	frontier := []uuid.UUID{orgID}

	for len(frontier) > 0 {
		var nextFrontier []uuid.UUID
		for _, parent := range frontier {
			children, err := r.GetChildOrganizationIDs(ctx, parent)
			if err != nil {
				return nil, err
			}
			for _, child := range children {
				occupied, err := r.IsOrganizationOccupied(ctx, child)
				if err != nil {
					return nil, err
				}
				if occupied {
					effective = append(effective, child)
					continue
				}
				nextFrontier = append(nextFrontier, child)
			}
		}
		frontier = nextFrontier
	}

	return effective, nil
}

// IsEmployeeInOrganizations reports whether the employee currently occupies
// (active employment) any of the given organizations — server-side guard for
// the ASSIGNED overtime flow so only effective subordinates can be assigned.
// Returns ErrOvertimeNotAssignable (defined in service.go) when the employee
// is outside the assigner's effective subordinate set.
func (r *Repository) IsEmployeeInOrganizations(ctx context.Context, employeeID uuid.UUID, orgIDs []uuid.UUID) error {
	if len(orgIDs) == 0 {
		return ErrOvertimeNotAssignable
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	var count int64
	if err := db.Table("employments").
		Where("employee_id = ? AND organization_id IN ? AND effective_end_date IS NULL", employeeID, orgIDs).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check assignable employee: %w", err)
	}
	if count == 0 {
		return ErrOvertimeNotAssignable
	}
	return nil
}

// assignableEmployee is the display row for the assign dropdown: an employee
// currently occupying one of the assigner's effective subordinate orgs.
type assignableEmployee struct {
	EmployeeID     string
	EmployeeCode   string
	Name           string
	OrganizationID string
	OrgName        string
}

// FindEmployeesByOrganizationIDs resolves employees currently occupying the
// given organizations (active employment), one row per employee, batched —
// used to fill the ASSIGNED overtime dropdown. Raw query without importing
// employee/organization packages (pola sama GetEmployeeInfoByIDs).
func (r *Repository) FindEmployeesByOrganizationIDs(ctx context.Context, orgIDs []uuid.UUID) ([]assignableEmployee, error) {
	if len(orgIDs) == 0 {
		return []assignableEmployee{}, nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var rows []assignableEmployee
	err = db.Table("employments AS em").
		Joins("JOIN employees AS e ON e.id = em.employee_id").
		Joins("JOIN organizations AS o ON o.id = em.organization_id").
		Where("em.organization_id IN ? AND em.effective_end_date IS NULL", orgIDs).
		Select("e.id AS employee_id, e.employee_id AS employee_code, e.name AS name, em.organization_id AS organization_id, o.nomenclature AS org_name").
		Order("e.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to resolve assignable employees: %w", err)
	}
	return rows, nil
}
