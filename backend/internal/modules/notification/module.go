package notification

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Notification"
const ModuleSlug = "notification"
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

	return &notificationModule{
		svc:    svc,
		logger: logger,
	}
}

// NewModuleWithService mounts the notification module using an
// already-constructed Service, so callers (e.g. main.go) can wire it into
// other modules as a Notifier before/while mounting it.
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
	return &notificationModule{
		svc:    svc,
		logger: logger,
	}
}

type notificationModule struct {
	svc    *Service
	logger *zap.Logger
}

func (m *notificationModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "In-app notifications delivered to users by other modules, with read/unread tracking",
		IsCore:      false,
		DependsOn:   []string{"useraccount"},
		Permissions: []string{
			"notification.view",
			"notification.manage",
		},
	}
}

// RegisterRoutes is a no-op for now — the REST API (handler + routes) is
// built in Phase 3 per docs/module-notification-plan.md §9.
func (m *notificationModule) RegisterRoutes(rg *gin.RouterGroup) {}

func (m *notificationModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Notification{})
}

func (m *notificationModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *notificationModule) Permissions() []string {
	return m.Info().Permissions
}
