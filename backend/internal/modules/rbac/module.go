package rbac

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "RBAC"
const ModuleSlug = "rbac"
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

	return &rbacModule{
		handler: handler,
		logger:  logger,
	}
}

type rbacModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *rbacModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage tenant roles, permissions, and user-role assignments",
		IsCore:      false,
		DependsOn:   []string{"setting"},
		Permissions: []string{
			"rbac.view", "rbac.create", "rbac.update", "rbac.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "RBAC",
				Icon:  "shield",
				Route: "/admin/settings/rbac",
				Children: []module.Menu{
					{Name: "Roles & Permissions", Icon: "shield", Route: "/admin/settings/rbac"},
				},
			},
		},
	}
}

func (m *rbacModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

// Migrate tidak melakukan AutoMigrate — tabel RBAC tenant (roles, permissions,
// role_has_permissions, model_has_roles, users) dibuat oleh migration SQL tenant
// (011_settings.sql & 022_users.sql). Menjalankan AutoMigrate di sini berisiko
// mengubah skema Spatie (composite PK) yang sudah ada.
func (m *rbacModule) Migrate(db *gorm.DB) error {
	return nil
}

func (m *rbacModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *rbacModule) Permissions() []string {
	return m.Info().Permissions
}
