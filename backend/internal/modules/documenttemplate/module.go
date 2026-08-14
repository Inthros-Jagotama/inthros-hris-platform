package documenttemplate

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

type documentTemplateModule struct {
	handler *Handler
}

func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	repo := NewRepository(NewTenantDBResolver(dbManager))
	svc := NewService(repo, logger)
	return &documentTemplateModule{handler: NewHandler(svc)}
}

func (m *documentTemplateModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        "Document Template",
		Slug:        "documenttemplate",
		Version:     "1.0.0",
		Description: "Manage reusable document templates for contracts and movement SKs",
		IsCore:      false,
		Permissions: []string{
			"documenttemplate.template.view",
			"documenttemplate.template.create",
			"documenttemplate.template.update",
			"documenttemplate.template.delete",
			"documenttemplate.template.activate",
			"documenttemplate.template.deactivate",
			"documenttemplate.template.set_default",
			"documenttemplate.template.version",
			"documenttemplate.generated.view",
		},
	}
}

func (m *documentTemplateModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *documentTemplateModule) Migrate(db *gorm.DB) error {
	// Tenant schema for this module is owned entirely by the versioned SQL
	// migration (110_document_templates.sql) — AutoMigrate does not run
	// against tenant databases in this codebase. This method is a required
	// module.Module interface stub, not the real migration path.
	return nil
}

func (m *documentTemplateModule) Seed(db *gorm.DB) error {
	// Default/reference templates are seeded by the SQL migration itself.
	return nil
}

func (m *documentTemplateModule) Permissions() []string {
	return m.Info().Permissions
}
