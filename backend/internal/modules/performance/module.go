package performance

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Performance Management"
const ModuleSlug = "performance"
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

	return &perfModule{
		handler: handler,
		logger:  logger,
	}
}

type perfModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *perfModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Balanced Scorecard performance management with KPI/OKR target setting, evaluation periods, perspective-based scoring, and 360 review capabilities",
		IsCore:      false,
		DependsOn:   []string{"organization", "employee", "jobmanagement", "competency"},
		Permissions: []string{
			"performance.view",
			"performance.create",
			"performance.update",
			"performance.delete",
			"performance.approve",
		},
		Menus: []module.Menu{
			{
				Name:  "Performance",
				Icon:  "bar-chart-2",
				Route: "/admin/performance",
				Children: []module.Menu{
					{Name: "Evaluations", Icon: "clipboard", Route: "/admin/performance/evaluations"},
					{Name: "Templates", Icon: "file-text", Route: "/admin/performance/templates"},
					{Name: "Indicators", Icon: "target", Route: "/admin/performance/indicators"},
					{Name: "Periods", Icon: "calendar", Route: "/admin/performance/periods"},
					{Name: "Perspectives", Icon: "pie-chart", Route: "/admin/performance/perspectives"},
				},
			},
		},
	}
}

func (m *perfModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *perfModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&PerformancePeriod{},
		&PerformancePerspective{},
		&PerformanceTemplate{},
		&PerformanceIndicator{},
		&PerformanceEvaluation{},
		&PerformanceEvaluationDetail{},
		&PerformanceTarget{},
	)
}

func (m *perfModule) Seed(db *gorm.DB) error {
	// Seed default perspectives if empty
	var count int64
	if err := db.Model(&PerformancePerspective{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
	d1 := "Financial performance and profitability metrics"
	d2 := "Customer satisfaction and market perspective"
	d3 := "Internal business process efficiency"
	d4 := "Employee development and organizational growth"
	defaults := []PerformancePerspective{
		{Name: "Financial", Description: &d1, SortOrder: 1},
		{Name: "Customer", Description: &d2, SortOrder: 2},
		{Name: "Internal Process", Description: &d3, SortOrder: 3},
		{Name: "Learning & Growth", Description: &d4, SortOrder: 4},
	}
		for _, p := range defaults {
			if err := db.Create(&p).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *perfModule) Permissions() []string {
	return m.Info().Permissions
}
