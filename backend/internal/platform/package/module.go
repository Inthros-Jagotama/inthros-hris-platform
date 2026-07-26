package pkgmgr

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const (
	ModuleName    = "Package Management"
	ModuleSlug    = "package-management"
	ModuleVersion = "1.0.0"
)

// NewModule membuat instance baru Package Management Module.
// Service dibuat secara internal.
func NewModule(dbManager *database.Manager, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, logger)
	return NewModuleWithService(svc, logger, authMW, rbacMW)
}

// NewModuleWithService membuat instance Package Module dengan service yang sudah ada.
// Digunakan ketika service perlu di-share dengan komponen lain (misal tenant route).
func NewModuleWithService(svc *Service, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	handler := NewHandler(svc)

	return &pkgModule{
		handler: handler,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type pkgModule struct {
	handler *Handler
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

// SetRBACMiddleware memperbarui RBAC middleware setelah inisialisasi.
func (m *pkgModule) SetRBACMiddleware(mw gin.HandlerFunc) {
	m.rbacMW = mw
}

func (m *pkgModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage module packages with pricing, mandatory modules, and dependency validation",
		IsCore:      true,
		DependsOn:   []string{"module-management"},
		Permissions: []string{
			"package.view",
			"package.create",
			"package.update",
			"package.delete",
			"package.publish",
		},
		Menus: []module.Menu{
			{
				Name:  "Packages",
				Icon:  "box",
				Route: "/admin/platform/packages",
			},
		},
	}
}

func (m *pkgModule) RegisterRoutes(rg *gin.RouterGroup) {
	// Public endpoint /api/v1/public/packages didaftarkan langsung di main.go
	// pada root router agar tidak memerlukan auth middleware.
	// Di sini publicRG = nil karena public route sudah di-handle main.go.
	RegisterRoutes(rg, nil, m.handler, m.authMW, m.rbacMW)
}

func (m *pkgModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Package{}, &PackageModule{})
}

func (m *pkgModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *pkgModule) Permissions() []string {
	return m.Info().Permissions
}
