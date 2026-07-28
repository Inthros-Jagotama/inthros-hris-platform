package recruitment

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

const ModuleName = "Recruitment & Onboarding (ATS)"
const ModuleSlug = "recruitment"
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

	return &recModule{
		handler: handler,
		logger:  logger,
	}
}

type recModule struct {
	handler *Handler
	logger  *zap.Logger
}

func (m *recModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "End-to-end recruitment & applicant tracking system with job requisitions, candidate management, interview scheduling, and employee onboarding workflows",
		IsCore:      false,
		DependsOn:   []string{"organization", "employee", "setting"},
		Permissions: []string{
			"recruitment.view",
			"recruitment.create",
			"recruitment.update",
			"recruitment.delete",
			"recruitment.interview",
			"recruitment.onboard",
		},
		Menus: []module.Menu{
			{
				Name:  "Recruitment",
				Icon:  "users",
				Route: "/admin/recruitment",
				Children: []module.Menu{
					{Name: "Requisitions", Icon: "briefcase", Route: "/admin/recruitment/requisitions"},
					{Name: "Candidates", Icon: "user-plus", Route: "/admin/recruitment/candidates"},
					{Name: "Applications", Icon: "file-text", Route: "/admin/recruitment/applications"},
					{Name: "Interviews", Icon: "calendar", Route: "/admin/recruitment/interviews"},
					{Name: "Onboarding", Icon: "log-in", Route: "/admin/recruitment/onboarding"},
				},
			},
		},
	}
}

func (m *recModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *recModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&JobRequisition{},
		&Candidate{},
		&JobApplication{},
		&Interview{},
		&OnboardingTaskTemplate{},
		&EmployeeOnboarding{},
		&OnboardingTaskItem{},
	)
}

func (m *recModule) Seed(db *gorm.DB) error {
	// Seed default onboarding task templates if empty
	var count int64
	if err := db.Model(&OnboardingTaskTemplate{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		defaults := []OnboardingTaskTemplate{
			{Name: "Personal Data Completion", Description: "Complete personal data and emergency contact information", Category: "ADMINISTRATIVE", DayOffset: 0, AssignedRole: "employee", IsMandatory: true},
			{Name: "Employment Contract Signing", Description: "Sign employment contract and non-disclosure agreement", Category: "LEGAL", DayOffset: 0, AssignedRole: "employee", IsMandatory: true},
			{Name: "IT Account Provisioning", Description: "Create email account, system access, and workstation setup", Category: "IT", DayOffset: -7, AssignedRole: "it", IsMandatory: true},
			{Name: "BPJS Registration", Description: "Register employee for BPJS Kesehatan & Ketenagakerjaan", Category: "HR", DayOffset: 0, AssignedRole: "hr", IsMandatory: true},
			{Name: "Bank Account Setup", Description: "Set up payroll bank account", Category: "HR", DayOffset: 0, AssignedRole: "employee", IsMandatory: true},
			{Name: "Orientation & Department Introduction", Description: "Introduction to team members, department overview, and company culture", Category: "ORIENTATION", DayOffset: 1, AssignedRole: "hr", IsMandatory: true},
			{Name: "Policy Review & Acknowledgment", Description: "Review company policies, code of conduct, and employee handbook", Category: "LEGAL", DayOffset: 3, AssignedRole: "employee", IsMandatory: true},
			{Name: "First Week Check-in", Description: "HR check-in to ensure smooth onboarding", Category: "HR", DayOffset: 7, AssignedRole: "hr", IsMandatory: false},
			{Name: "Training Plan Setup", Description: "Create 30-60-90 day training and development plan", Category: "HR", DayOffset: 14, AssignedRole: "manager", IsMandatory: false},
			{Name: "Probation Review Preparation", Description: "Prepare documents and feedback for probation end review", Category: "HR", DayOffset: 90, AssignedRole: "manager", IsMandatory: true},
		}
		for _, t := range defaults {
			if err := db.Create(&t).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *recModule) Permissions() []string {
	return m.Info().Permissions
}
