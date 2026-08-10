package careerintelligence

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Career Intelligence & Talent Management"
const ModuleSlug = "career-intelligence"
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

// NewModule creates a new Career Intelligence module.
// This module manages talent mapping, career interests, career paths,
// gap analysis, and succession planning.
func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	return &ciModule{
		handler: handler,
		logger:  logger,
	}
}

type ciModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *ciModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Career Intelligence & Talent Management — talent mapping (9-box grid), career interests and aspirations, career path templates, gap analysis, and succession planning for key positions.",
		DependsOn: []string{
			"organization",
			"employee",
			"jobmanagement",
			"competency",
		},
	}
}

func (m *ciModule) RegisterRoutes(router *gin.RouterGroup) {
	RegisterRoutes(router, m.handler)
}

func (m *ciModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&CareerTalentMap{},
		&CareerInterest{},
		&CareerPath{},
		&CareerPathStep{},
		&CareerSuccessionPlan{},
	)
}

func (m *ciModule) Seed(db *gorm.DB) error {
	// No seed data needed — module is data-driven from employee assessments
	return nil
}

func (m *ciModule) Permissions() []string {
	return []string{
		"career-intelligence.view",
		"career-intelligence.talent-map.manage",
		"career-intelligence.talent-map.view",
		"career-intelligence.interest.manage",
		"career-intelligence.interest.view",
		"career-intelligence.path.manage",
		"career-intelligence.path.view",
		"career-intelligence.succession.manage",
		"career-intelligence.succession.view",
		"career-intelligence.gap-analysis.view",
	}
}
