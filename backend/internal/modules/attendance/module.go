package attendance

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Time & Attendance"
const ModuleSlug = "attendance"
const ModuleVersion = "1.0.0"

type TenantDBFunc func(ctx context.Context) (*gorm.DB, error)

func NewTenantDBResolver(dbManager *database.Manager) TenantDBFunc {
	return func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
}

func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	return &attModule{
		handler: handler,
		logger:  logger,
	}
}

// NewModuleWithService mounts the attendance module's routes using an
// already-constructed Service, so callers (e.g. main.go) can register a
// push-based approval status handler against it before it's wrapped.
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
	handler := NewHandler(svc)

	return &attModule{
		handler: handler,
		logger:  logger,
	}
}

type attModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *attModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage time & attendance including shifts, check-in/check-out events, geofence locations, overtime requests, and daily work sessions",
		IsCore:      false,
		DependsOn:   []string{"employee", "organization"},
		Permissions: []string{
			"attendance.view",
			"attendance.create",
			"attendance.update",
			"attendance.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Attendance",
				Icon:  "clock",
				Route: "/admin/attendance",
				Children: []module.Menu{
					{Name: "Dashboard", Icon: "bar-chart", Route: "/admin/attendance/dashboard"},
					{Name: "Shifts", Icon: "calendar", Route: "/admin/attendance/shifts"},
					{Name: "Schedules", Icon: "clock", Route: "/admin/attendance/schedules"},
					{Name: "Events", Icon: "list", Route: "/admin/attendance/events"},
					{Name: "Overtime", Icon: "clock", Route: "/admin/attendance/overtime"},
					{Name: "Locations", Icon: "map-pin", Route: "/admin/attendance/locations"},
					{Name: "Settings", Icon: "settings", Route: "/admin/attendance/settings"},
				},
			},
		},
	}
}

func (m *attModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *attModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&AttendanceCompanySetting{},
		&AttendanceCompanyShift{},
		&AttendanceEmployeeShift{},
		&AttendanceLocation{},
		&AttendanceDeviceCapture{},
		&AttendanceFaceCapture{},
		&AttendanceEvent{},
		&AttendanceSession{},
		&AttendanceOvertimeRequest{},
		&AttendanceExemptPosition{},
		&AttendanceCorrectionRequest{},
	)
}

func (m *attModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *attModule) Permissions() []string {
	return m.Info().Permissions
}
