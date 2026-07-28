package payroll

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Payroll & Compensation Engine"
const ModuleSlug = "payroll"
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

func NewModule(dbManager *database.Manager, logger *zap.Logger, opts ...interface{}) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	for _, opt := range opts {
		if ae, ok := opt.(ApprovalEngine); ok {
			svc.SetApprovalEngine(ae)
		}
	}
	handler := NewHandler(svc)

	return &payrollModule{
		handler: handler,
		logger:  logger,
	}
}

type payrollModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *payrollModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage payroll structure (salary components, grade components, employee profiles, bank/BPJS/tax profiles, BPJS settings, PPh21 settings), payroll runs, payslips, and tax calculations",
		IsCore:      true,
		DependsOn:   []string{"employee", "organization", "setting"},
		Permissions: []string{
			"payroll.view",
			"payroll.create",
			"payroll.update",
			"payroll.delete",
			"payroll.approve",
			"payroll.run",
		},
		Menus: []module.Menu{
			{
				Name:  "Payroll",
				Icon:  "dollar-sign",
				Route: "/admin/payroll",
				Children: []module.Menu{
					{Name: "Payroll Runs", Icon: "play-circle", Route: "/admin/payroll/runs"},
					{Name: "Periods", Icon: "calendar", Route: "/admin/payroll/periods"},
					{Name: "Salary Components", Icon: "list", Route: "/admin/payroll/salary-components"},
					{Name: "Employee Profiles", Icon: "users", Route: "/admin/payroll/profiles"},
					{Name: "BPJS Settings", Icon: "activity", Route: "/admin/payroll/bpjs-settings"},
					{Name: "PPh21 Settings", Icon: "percent", Route: "/admin/payroll/pph21-settings"},
					{Name: "PTKP Rates", Icon: "trending-up", Route: "/admin/payroll/ptkp-rates"},
					{Name: "Tax Brackets", Icon: "bar-chart", Route: "/admin/payroll/tax-brackets"},
				},
			},
		},
	}
}

func (m *payrollModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *payrollModule) Migrate(db *gorm.DB) error {
	// SQL migration already handled by 006_payroll_structure.sql and 007_payroll_run.sql migrator
	return nil
}

func (m *payrollModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *payrollModule) Permissions() []string {
	return m.Info().Permissions
}
