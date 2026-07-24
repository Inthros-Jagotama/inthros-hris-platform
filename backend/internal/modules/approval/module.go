package approval

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Approval Engine"
const ModuleSlug = "approval"
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

	return &approvalModule{
		handler: handler,
		logger:  logger,
	}
}

// NewModuleWithService creates an approval module using an existing service instance.
// This allows other modules (like payroll) to share the same approval service instance.
func NewModuleWithService(dbManager *database.Manager, logger *zap.Logger, svc *Service) module.Module {
	handler := NewHandler(svc)

	return &approvalModule{
		handler: handler,
		logger:  logger,
	}
}

type approvalModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *approvalModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Multi-step approval workflow engine supporting flexible approval flows (ANY_ONE, ALL, N_OF_M), configurable per-module, with instance tracking and task management for approvers",
		IsCore:      true,
		DependsOn:   []string{"employee", "organization"},
		Permissions: []string{
			"approval.view",
			"approval.create",
			"approval.update",
			"approval.delete",
			"approval.approve",
		},
		Menus: []module.Menu{
			{
				Name:  "Approval",
				Icon:  "check-circle",
				Route: "/admin/approval",
				Children: []module.Menu{
					{Name: "My Tasks", Icon: "inbox", Route: "/admin/approval/tasks"},
					{Name: "Approval Flows", Icon: "git-branch", Route: "/admin/approval/flows"},
					{Name: "Instances", Icon: "list", Route: "/admin/approval/instances"},
				},
			},
		},
	}
}

func (m *approvalModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *approvalModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ApprovalFlow{},
		&ApprovalFlowStep{},
		&ApprovalInstance{},
		&ApprovalAction{},
		&ApprovalTask{},
	)
}

func (m *approvalModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *approvalModule) Permissions() []string {
	return m.Info().Permissions
}
