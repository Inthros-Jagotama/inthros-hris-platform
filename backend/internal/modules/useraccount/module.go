package useraccount

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "User Account"
const ModuleSlug = "useraccount"
const ModuleVersion = "1.0.0"

// TenantDBFunc adalah resolver tenant database (pola sama dengan module lain).
type TenantDBFunc func(ctx context.Context) (*gorm.DB, error)

// NewTenantDBResolver membuat resolver dari company_id di context.
func NewTenantDBResolver(dbManager *database.Manager) TenantDBFunc {
	return func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
}

// NewModule membuat module user-account. mailer dipakai untuk kirim email setup,
// authManager untuk JWT tenant login, platformDB untuk resolve company by slug.
func NewModule(dbManager *database.Manager, authManager *auth.Manager, mailer Mailer, logger *zap.Logger) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver, dbManager.PlatformDB())
	svc := NewService(repo, authManager, mailer, logger)
	handler := NewHandler(svc)

	return &userAccountModule{
		handler: handler,
		logger:  logger,
	}
}

type userAccountModule struct {
	handler *Handler
	logger  *zap.Logger
}

// PublicHandler diekspos agar main.go bisa mendaftarkan route publik
// (set-password via link email) di luar middleware auth tenant.
func (m *userAccountModule) PublicHandler() *Handler {
	return m.handler
}

func (m *userAccountModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage employee login accounts (user tenant + setup password via email link)",
		IsCore:      true,
		DependsOn:   []string{"employee", "setting", "rbac"},
		Permissions: []string{
			"useraccount.create",
			"useraccount.view",
			"useraccount.update",
		},
		Menus: nil,
	}
}

func (m *userAccountModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *userAccountModule) Migrate(db *gorm.DB) error {
	// Tabel employee_accounts dibuat oleh SQL migration 023_user_accounts.sql.
	return nil
}

func (m *userAccountModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *userAccountModule) Permissions() []string {
	return m.Info().Permissions
}
