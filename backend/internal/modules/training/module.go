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
	return NewModuleWithService(logger, svc)
}

// NewModuleWithService constructs the module from an already-built service
// (pola leave/payroll/reimbursement) — dipakai di main.go agar status handler
// approval (training_request) bisa didaftarkan sebelum module di-mount.
func NewModuleWithService(logger *zap.Logger, svc *Service) module.Module {
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
		Description: "End-to-end training and development management including course catalog, provider/trainer masters, session scheduling, enrollment, attendance, assessments, materials, evaluations, and certificate issuance",
		IsCore:      false,
		DependsOn:   []string{"employee", "organization", "setting", "competency", "approval"},
		Permissions: []string{
			"training.view",
			"training.create",
			"training.update",
			"training.delete",
			"training.enroll",
			"training.course.manage",
			"training.session.manage",
			"training.participant.manage",
			"training.attendance.manage",
			"training.assessment.manage",
			"training.evaluation.manage",
			"training.certificate.manage",
			"training.plan.manage",
			"training.request.create",
			"training.request.approve",
			"training.report.view",
		},
		Menus: []module.Menu{
			{
				Name:  "Training",
				Icon:  "book-open",
				Route: "/training",
				Children: []module.Menu{
					{Name: "Courses", Icon: "book", Route: "/training/courses"},
					{Name: "Categories", Icon: "tag", Route: "/training/categories"},
					{Name: "Providers", Icon: "building", Route: "/training/providers"},
					{Name: "Trainers", Icon: "user", Route: "/training/trainers"},
					{Name: "Sessions", Icon: "calendar", Route: "/training/sessions"},
					{Name: "Participants", Icon: "users", Route: "/training/participants"},
					{Name: "Certificates", Icon: "award", Route: "/training/certificates"},
					{Name: "Evaluations", Icon: "star", Route: "/training/evaluations"},
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
		// P0-BE (plan §42 P0-BE)
		&TrainingProvider{},
		&TrainingTrainer{},
		&TrainingSessionTrainer{},
		&TrainingAttendance{},
		&TrainingAssessment{},
		&TrainingAssessmentResult{},
		// P1-BE (plan §42 P1-BE)
		&TrainingPlan{},
		&TrainingPlanItem{},
		&TrainingNeed{},
		&TrainingRequest{},
		&TrainingCourseObjective{},
		&TrainingCourseCompetency{},
		&TrainingCoursePrerequisite{},
		&TrainingMandatory{},
		&TrainingSessionCost{},
		&TrainingDocument{},
		// P2-BE (plan §42 P2-BE)
		&TrainingEvaluationForm{},
		&TrainingEvaluationQuestion{},
		&TrainingEvaluationAnswer{},
		&TrainingEffectivenessAssessment{},
		&TrainingCertification{},
	)
}

func (m *trainingModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *trainingModule) Permissions() []string {
	return m.Info().Permissions
}
