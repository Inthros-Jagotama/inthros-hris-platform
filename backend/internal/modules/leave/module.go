package leave

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Leave & Time Off"
const ModuleSlug = "leave"
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

	return &leaveModule{
		handler: handler,
		logger:  logger,
	}
}

// NewModuleWithService mounts the leave module's routes using an
// already-constructed Service, so callers (e.g. main.go) can register a
// push-based approval status handler against it before it's wrapped.
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
	handler := NewHandler(svc)

	return &leaveModule{
		handler: handler,
		logger:  logger,
	}
}

type leaveModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *leaveModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage leave & time off including leave types, accrual policies, leave requests, approvals, and employee leave balances",
		IsCore:      false,
		DependsOn:   []string{"employee", "organization"},
		Permissions: []string{
			"leave.view",
			"leave.create",
			"leave.update",
			"leave.delete",
			"leave.approve",
		},
		Menus: []module.Menu{
			{
				Name:  "Leave",
				Icon:  "calendar-off",
				Route: "/admin/leave",
				Children: []module.Menu{
					{Name: "Dashboard", Icon: "bar-chart", Route: "/admin/leave/dashboard"},
					{Name: "Requests", Icon: "list", Route: "/admin/leave/requests"},
					{Name: "Types", Icon: "tag", Route: "/admin/leave/types"},
					{Name: "Balances", Icon: "wallet", Route: "/admin/leave/balances"},
					{Name: "Settings", Icon: "settings", Route: "/admin/leave/settings"},
				},
			},
		},
	}
}

func (m *leaveModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *leaveModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&LeaveType{},
		&LeaveAccrualPolicy{},
		&LeaveReason{},
		&LeaveRequest{},
		&LeaveRequestDetail{},
		&EmployeeLeaveBalance{},
		&LeaveBalanceTransaction{},
	)
}

func (m *leaveModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *leaveModule) Permissions() []string {
	return m.Info().Permissions
}
