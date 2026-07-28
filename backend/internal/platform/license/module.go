package license

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const (
	ModuleName    = "License Management"
	ModuleSlug    = "license-management"
	ModuleVersion = "1.0.0"
)

// NewModule membuat instance baru License Management Module.
// Service dibuat secara internal (dengan nil moduleManager).
// Untuk auto-activation module, gunakan NewModuleWithService dan inject
// service yang sudah memiliki PackageModuleManager.
func NewModule(dbManager *database.Manager, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, dbManager, nil, logger)
	return NewModuleWithService(svc, logger, authMW, rbacMW)
}

// NewModuleWithService membuat instance License Module dengan service yang sudah ada.
// Digunakan untuk sharing service instance dengan adapter (misal company signup).
func NewModuleWithService(svc *Service, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	handler := NewHandler(svc)

	return &licenseModule{
		handler: handler,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type licenseModule struct {
	handler *Handler
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

// SetRBACMiddleware memperbarui RBAC middleware setelah inisialisasi.
func (m *licenseModule) SetRBACMiddleware(mw gin.HandlerFunc) {
	m.rbacMW = mw
}

func (m *licenseModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage company licenses, plans, and subscription billing",
		IsCore:      true,
		DependsOn:   []string{"company"},
		Permissions: []string{
			"license.view",
			"license.create",
			"license.update",
			"license.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Licenses",
				Icon:  "key",
				Route: "/admin/platform/licenses",
			},
		},
	}
}

func (m *licenseModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler, m.authMW, m.rbacMW)
}

func (m *licenseModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&License{})
}

func (m *licenseModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *licenseModule) Permissions() []string {
	return m.Info().Permissions
}
