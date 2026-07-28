package training

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Training & Development Management"
const ModuleSlug = "training"
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

	return &trainingModule{
		handler: handler,
		logger:  logger,
	}
}

type trainingModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *trainingModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "End-to-end training and development management including course catalog, session scheduling, participant registration, attendance tracking, materials, evaluations, and certificate issuance",
		IsCore:      false,
		DependsOn:   []string{"employee", "organization", "setting"},
		Permissions: []string{
			"training.view",
			"training.create",
			"training.update",
			"training.delete",
			"training.enroll",
		},
		Menus: []module.Menu{
			{
				Name:  "Training",
				Icon:  "book-open",
				Route: "/admin/training",
				Children: []module.Menu{
					{Name: "Courses", Icon: "book", Route: "/admin/training/courses"},
					{Name: "Sessions", Icon: "calendar", Route: "/admin/training/sessions"},
					{Name: "Participants", Icon: "users", Route: "/admin/training/participants"},
					{Name: "Categories", Icon: "tag", Route: "/admin/training/categories"},
					{Name: "Certificates", Icon: "award", Route: "/admin/training/certificates"},
					{Name: "Evaluations", Icon: "star", Route: "/admin/training/evaluations"},
				},
			},
		},
	}
}

func (m *trainingModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *trainingModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&TrainingCategory{},
		&TrainingCourse{},
		&TrainingSession{},
		&TrainingParticipant{},
		&TrainingMaterial{},
		&TrainingEvaluation{},
		&TrainingCertificate{},
	)
}

func (m *trainingModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *trainingModule) Permissions() []string {
	return m.Info().Permissions
}
