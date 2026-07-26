package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/auth"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const (
	ModuleName    = "Platform Users & Auth"
	ModuleSlug    = "platform-users"
	ModuleVersion = "1.0.0"
)

// NewModule membuat instance baru Platform User Module.
// Membuat repository dan service baru secara internal.
func NewModule(dbManager *database.Manager, authManager *auth.Manager, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	repo := NewRepository(dbManager.PlatformDB())
	svc := NewService(repo, authManager, logger)
	return newModuleWithService(svc, logger, authMW, rbacMW)
}

// NewModuleWithService membuat instance baru Platform User Module
// dengan service yang sudah ada. Digunakan dari main.go ketika
// service perlu dibagikan ke module lain (seperti Company module).
func NewModuleWithService(svc *Service, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	return newModuleWithService(svc, logger, authMW, rbacMW)
}

func newModuleWithService(svc *Service, logger *zap.Logger, authMW, rbacMW gin.HandlerFunc) module.Module {
	handler := NewHandler(svc)

	return &userModule{
		handler: handler,
		service: svc,
		logger:  logger,
		authMW:  authMW,
		rbacMW:  rbacMW,
	}
}

type userModule struct {
	handler *Handler
	service *Service
	logger  *zap.Logger
	authMW  gin.HandlerFunc
	rbacMW  gin.HandlerFunc
}

// SetRBACMiddleware memperbarui RBAC middleware setelah inisialisasi.
func (m *userModule) SetRBACMiddleware(mw gin.HandlerFunc) {
	m.rbacMW = mw
}

func (m *userModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Platform user authentication, authorization, and user management",
		IsCore:      true,
		DependsOn:   []string{},
		Permissions: []string{
			"user.view",
			"user.create",
			"user.update",
			"user.delete",
		},
		Menus: []module.Menu{
			{
				Name:  "Platform Users",
				Icon:  "users-cog",
				Route: "/admin/platform/users",
			},
		},
	}
}

func (m *userModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler, m.authMW, m.rbacMW)
}

func (m *userModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PlatformUser{})
}

func (m *userModule) Seed(db *gorm.DB) error {
	return m.service.EnsureSeed(db)
}

func (m *userModule) Permissions() []string {
	return m.Info().Permissions
}
