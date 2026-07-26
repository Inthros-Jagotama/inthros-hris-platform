package company

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

// ModuleName adalah identitas modul.
const ModuleName = "Company Management"
const ModuleSlug = "company"
const ModuleVersion = "1.0.0"

// NewModule membuat instance baru Company Module untuk registrasi.
// userCreator adalah implementasi dari interface UserCreator (user.Service),
// digunakan untuk membuat company_admin user saat create company.
func NewModule(dbManager *database.Manager, userCreator UserCreator, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, dbManager, userCreator, logger)
	handler := NewHandler(svc)

	return &companyModule{
		handler: handler,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type companyModule struct {
	handler *Handler
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

// SetRBACMiddleware memperbarui RBAC middleware setelah inisialisasi.
// Dipanggil dari main.go setelah RBAC enforcer siap.
func (m *companyModule) SetRBACMiddleware(mw gin.HandlerFunc) {
	m.rbacMW = mw
}

func (m *companyModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage companies/tenants and their lifecycle",
		IsCore:      true,
		DependsOn:   []string{},
		Permissions: []string{
			"company.view",
			"company.create",
			"company.update",
			"company.delete",
			"company.suspend",
			"company.activate",
		},
		Menus: []module.Menu{
			{
				Name:  "Companies",
				Icon:  "building",
				Route: "/admin/platform/companies",
			},
		},
	}
}

func (m *companyModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler, m.authMW, m.rbacMW)
}

func (m *companyModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Company{}, &TenantConnection{})
}

func (m *companyModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *companyModule) Permissions() []string {
	return m.Info().Permissions
}
