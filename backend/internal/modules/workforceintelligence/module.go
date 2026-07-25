package workforceintelligence

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Workforce Intelligence & Strategic Planning"
const ModuleSlug = "workforce-intelligence"
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

// NewModule creates a new Workforce Intelligence module.
// This module is a strategic analytics layer that reads data from
// all other tenant modules. It does NOT replace any operational module.
func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	handler := NewHandler(svc)

	return &wfModule{
		handler: handler,
		logger:  logger,
	}
}

type wfModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *wfModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Strategic analytics layer for workforce planning, intelligence, analytics, capacity, cost, risk, executive dashboards, scenario planning, and organization health. Reads data from all tenant modules.",
		DependsOn: []string{
			"organization",
			"employee",
			"attendance",
			"leave",
			"payroll",
			"performance",
			"competency",
			"training",
			"recruitment",
			"employeemovement",
		},
	}
}

func (m *wfModule) RegisterRoutes(router *gin.RouterGroup) {
	RegisterRoutes(router, m.handler)
}

func (m *wfModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&WorkforcePlanningHeadcount{},
		&WorkforceForecast{},
		&WorkforceKPI{},
		&WorkforceAnalyticsCache{},
		&WorkforceScenario{},
		&WorkforceRiskIndicator{},
		&WorkforceHealthScore{},
	)
}

func (m *wfModule) Seed(db *gorm.DB) error {
	// No seed data needed — module is data-driven from other modules
	return nil
}

func (m *wfModule) Permissions() []string {
	return []string{
		"workforce-intelligence.view",
		"workforce-intelligence.planning.manage",
		"workforce-intelligence.planning.view",
		"workforce-intelligence.analytics.view",
		"workforce-intelligence.capacity.view",
		"workforce-intelligence.cost.view",
		"workforce-intelligence.risk.view",
		"workforce-intelligence.executive.view",
		"workforce-intelligence.scenario.manage",
		"workforce-intelligence.scenario.view",
		"workforce-intelligence.health.view",
		"workforce-intelligence.people-analytics.view",
	}
}
