package reimbursement

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Reimbursement & Claim"
const ModuleSlug = "reimbursement"
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

	return &reimbursementModule{
		handler: handler,
		logger:  logger,
	}
}

// NewModuleWithService mounts the reimbursement module's routes using an
// already-constructed Service, so callers (e.g. main.go) can register a
// push-based approval status handler against it before it's wrapped.
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
	handler := NewHandler(svc)

	return &reimbursementModule{
		handler: handler,
		logger:  logger,
	}
}

type reimbursementModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *reimbursementModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage employee reimbursements & claims including reimbursement types, requests, items, and payment processing",
		IsCore:      false,
		DependsOn:   []string{"employee"},
		Permissions: []string{
			"reimbursement.view",
			"reimbursement.create",
			"reimbursement.update",
			"reimbursement.delete",
			"reimbursement.approve",
		},
		Menus: []module.Menu{
			{
				Name:  "Reimbursement",
				Icon:  "receipt",
				Route: "/admin/reimbursement",
				Children: []module.Menu{
					{Name: "Requests", Icon: "list", Route: "/admin/reimbursement/requests"},
					{Name: "Types", Icon: "tag", Route: "/admin/reimbursement/types"},
					{Name: "Reports", Icon: "bar-chart", Route: "/admin/reimbursement/reports"},
				},
			},
		},
	}
}

func (m *reimbursementModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *reimbursementModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&ReimbursementType{},
		&ReimbursementRequest{},
		&ReimbursementItem{},
	)
}

func (m *reimbursementModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *reimbursementModule) Permissions() []string {
	return m.Info().Permissions
}
