package competency

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Competency Management"
const ModuleSlug = "competency"
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

	return &compModule{
		handler: handler,
		logger:  logger,
	}
}

// NewModuleWithService mounts the competency module's routes using an
// already-constructed Service, so callers (e.g. main.go) can register a
// push-based approval status handler against it before it's wrapped
// (pola sama employeemovement/attendance — §34.2).
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
	handler := NewHandler(svc)

	return &compModule{
		handler: handler,
		logger:  logger,
	}
}

type compModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *compModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage competency dictionaries, values, events, scoring, and gap analysis",
		IsCore:      true,
		DependsOn:   []string{"organization", "employee"},
		Permissions: []string{
			"competency.view",
			"competency.create",
			"competency.update",
			"competency.delete",
			// Competency 360 (§23 plan generik)
			"competency_360.view",
			"competency_360.manage",
			"competency_360.create",
			"competency_360.assign_rater",
			"competency_360.assess",
			"competency_360.approve",
			"competency_360.view_result",
			"competency_360.view_report",
		},
		Menus: []module.Menu{
			{
				Name:  "Competency",
				Icon:  "award",
				Route: "/admin/competency",
				Children: []module.Menu{
					{Name: "Competencies", Icon: "book", Route: "/admin/competency/competencies"},
					{Name: "Values", Icon: "sliders", Route: "/admin/competency/values"},
					{Name: "Events", Icon: "calendar", Route: "/admin/competency/events"},
					{Name: "Scores", Icon: "bar-chart", Route: "/admin/competency/scores"},
					{Name: "Assessment Templates", Icon: "clone", Route: "/admin/competency/templates"},
					{Name: "Rater Assignment", Icon: "users", Route: "/admin/competency/raters"},
					{Name: "My Assessment", Icon: "clipboard-list", Route: "/admin/competency/my-assessments"},
				},
			},
		},
	}
}

func (m *compModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *compModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Competency{},
		&CompetenceValue{},
		&CompetencyValue{},
		&CompetencyEvent{},
		&CompetencyEventTarget{},
		&CompetencyScore{},
		&CompetencyScoreDetail{},
		// Competency 360
		&CompetencyRatingScale{},
		&CompetencyRatingScaleItem{},
		&CompetencyAssessmentTemplate{},
		&CompetencyAssessmentTemplateCompetency{},
		&CompetencyAssessmentTemplateRaterType{},
		&CompetencyIndicator{},
		&CompetencyAssessmentTemplateIndicator{},
		&CompetencyAssessmentRater{},
		&CompetencyAssessmentResponse{},
	)
}

func (m *compModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *compModule) Permissions() []string {
	return m.Info().Permissions
}
