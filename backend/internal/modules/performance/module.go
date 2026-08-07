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

	// OKR
	okrRepo := NewOKRRepository()
	okrSvc := NewOKRService(okrRepo)
	okrHandler := NewOKRHandler(okrSvc, resolver)

	return &perfModule{
		handler:    handler,
		okrHandler: okrHandler,
		logger:     logger,
	}
}

// NewModuleWithService mounts the performance module's routes using an
// already-constructed Service, so callers (e.g. main.go) can register a
// push-based approval status handler against it before it's wrapped.
func NewModuleWithService(dbManager *database.Manager, logger *zap.Logger, svc *Service) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	handler := NewHandler(svc)

	okrRepo := NewOKRRepository()
	okrSvc := NewOKRService(okrRepo)
	okrHandler := NewOKRHandler(okrSvc, resolver)

	return &perfModule{
		handler:    handler,
		okrHandler: okrHandler,
		logger:     logger,
	}
}

type perfModule struct {
	handler    *Handler
	okrHandler *OKRHandler
	logger     *zap.Logger
}

func (m *perfModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Balanced Scorecard performance management with KPI/OKR target setting, evaluation periods, perspective-based scoring, and 360 review capabilities",
		IsCore:      false,
		DependsOn:   []string{"organization", "employee", "jobmanagement", "competency", "setting"},
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
	RegisterOKRRoutes(rg, m.okrHandler)
}

func (m *perfModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// KPI Models
		&PerformancePeriod{},
		&PerformancePerspective{},
		&PerformanceTemplate{},
		&PerformanceIndicator{},
		&PerformanceEvaluation{},
		&PerformanceEvaluationDetail{},
		&PerformanceTarget{},
		// OKR Models
		&OKRTemplate{},
		&OKRObjective{},
		&OKRKeyResult{},
		&OKREvaluation{},
		&OKREvaluationDetail{},
		&OKRProgress{},
		&OKRComment{},
		&OKRAttachment{},
		// Phase 5 - Scoring Configuration
		&PerformanceComponent{},
		&PerformanceOrganizationComponent{},
		&PerformanceEvaluationComponent{},
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

	// Seed default performance components if empty
	var compCount int64
	if err := db.Model(&PerformanceComponent{}).Count(&compCount).Error; err != nil {
		return err
	}
	if compCount == 0 {
		cd1 := "Key Performance Indicator score"
		cd2 := "Work program achievement score (manual input)"
		cd3 := "Average final score of direct-report Organizations"
		components := []PerformanceComponent{
			{Code: "KPI", Name: "KPI", Description: &cd1, SortOrder: 1, IsActive: true},
			{Code: "WORK_PROGRAM", Name: "Work Program", Description: &cd2, SortOrder: 2, IsActive: true},
			{Code: "SUBORDINATE", Name: "Subordinate KPI", Description: &cd3, SortOrder: 3, IsActive: true},
		}
		for _, c := range components {
			if err := db.Create(&c).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *perfModule) Permissions() []string {
	return m.Info().Permissions
}
