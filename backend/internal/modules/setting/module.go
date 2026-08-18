package setting

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/documenttemplate"
	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
	"github.com/inthros/hris-platform/internal/pkg/numbering"
)

const ModuleName = "Settings"
const ModuleSlug = "setting"
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

// NewModule constructs the setting module. numberingSvc is the shared
// numbering.Service instance owned by main.go (also wired into the
// employeemovement module in Task 5) — the setting module does not
// construct its own, so there is exactly one Service (and therefore one
// consistent DB-transactional sequence state) per process.
// employeeIDFormatSvc is likewise the shared EmployeeIDFormatService
// instance owned by main.go (also wired into the employee module so
// CreateEmployee can resolve employee_id per the configured generation
// mode) — same single-instance rationale as numberingSvc.
// uploadDir adalah root direktori upload (diserve publik via /uploads) —
// dipakai documenttemplate untuk menyimpan file template .docx.
// pdfSvc (opsional) dipakai documenttemplate untuk preview DOCX → PDF.
func NewModule(dbManager *database.Manager, logger *zap.Logger, numberingSvc *numbering.Service, employeeIDFormatSvc *EmployeeIDFormatService, uploadDir string, pdfSvc documenttemplate.PDFService) module.Module {
	resolver := NewTenantDBResolver(dbManager)
	repo := NewRepository(resolver)
	svc := NewService(repo, logger)
	handler := NewHandler(svc)
	numberingHandler := NewNumberingHandler(numberingSvc)
	empIDFormatHandler := NewEmployeeIDFormatHandler(employeeIDFormatSvc)

	// Document templates adalah sub-feature dari Settings (lihat plan
	// docs/module-settngs-fitur-template-dokumen-plan.md) — handler-nya
	// dibangun di sini dan route-nya didaftarkan di bawah /settings/...,
	// bukan sebagai module tenant terpisah.
	dtRepo := documenttemplate.NewRepository(documenttemplate.NewTenantDBResolver(dbManager))
	dtSvc := documenttemplate.NewService(dtRepo, logger)
	dtHandler := documenttemplate.NewHandler(dtSvc, uploadDir, pdfSvc)

	return &settingModule{
		handler:                 handler,
		numberingHandler:        numberingHandler,
		documentTemplateHandler: dtHandler,
		employeeIDFormatHandler: empIDFormatHandler,
		employeeIDFormatSvc:     employeeIDFormatSvc,
		logger:                  logger,
	}
}

type settingModule struct {
	handler                 *Handler
	numberingHandler        *NumberingHandler
	documentTemplateHandler *documenttemplate.Handler
	employeeIDFormatHandler *EmployeeIDFormatHandler
	employeeIDFormatSvc     *EmployeeIDFormatService
	logger                  *zap.Logger
}

// EmployeeIDFormatService returns the tenant employee_id-format service so
// main.go can wire it into the employee module (mirrors the SetQuotaChecker
// / employeemovement numbering wiring pattern used elsewhere in this repo).
func (m *settingModule) EmployeeIDFormatService() *EmployeeIDFormatService {
	return m.employeeIDFormatSvc
}

func (m *settingModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        ModuleName,
		Slug:        ModuleSlug,
		Version:     ModuleVersion,
		Description: "Manage tenant settings including zones, regions, provinces, and administrative areas",
		IsCore:      false,
		DependsOn:   []string{},
		Permissions: []string{
			"setting.zone.view", "setting.zone.create", "setting.zone.update", "setting.zone.delete",
			"setting.province.view", "setting.province.create", "setting.province.update", "setting.province.delete",
			"setting.regency.view", "setting.regency.create", "setting.regency.update", "setting.regency.delete",
			"setting.district.view", "setting.district.create", "setting.district.update", "setting.district.delete",
			"setting.village.view", "setting.village.create", "setting.village.update", "setting.village.delete",
			"setting.education.view", "setting.education.create", "setting.education.update", "setting.education.delete",
			"setting.education_major.view", "setting.education_major.create", "setting.education_major.update", "setting.education_major.delete",
			"setting.religion.view", "setting.religion.create", "setting.religion.update", "setting.religion.delete",
			"setting.marital_status.view", "setting.marital_status.create", "setting.marital_status.update", "setting.marital_status.delete",
			"setting.relationship_type.view", "setting.relationship_type.create", "setting.relationship_type.update", "setting.relationship_type.delete",
			"setting.bank.view", "setting.bank.create", "setting.bank.update", "setting.bank.delete",
			"setting.employment_status.view", "setting.employment_status.create", "setting.employment_status.update", "setting.employment_status.delete",
			"setting.nationality.view", "setting.nationality.create", "setting.nationality.update", "setting.nationality.delete",
			"setting.job_family.view", "setting.job_family.create", "setting.job_family.update", "setting.job_family.delete",
			"setting.grading.view", "setting.grading.create", "setting.grading.update", "setting.grading.delete",
			"setting.salary_grade.view", "setting.salary_grade.create", "setting.salary_grade.update", "setting.salary_grade.delete",
			"setting.ter.view", "setting.ter.create", "setting.ter.update", "setting.ter.delete",
			"setting.ptkp.view", "setting.ptkp.create", "setting.ptkp.update", "setting.ptkp.delete",
			"setting.insurance.view", "setting.insurance.create", "setting.insurance.update", "setting.insurance.delete",
			"setting.company_holiday.view", "setting.company_holiday.create", "setting.company_holiday.update", "setting.company_holiday.delete",
			"setting.competency.view", "setting.competency.create", "setting.competency.update", "setting.competency.delete",
			"setting.document_template.view", "setting.document_template.create", "setting.document_template.update", "setting.document_template.delete",
			"setting.document_template.activate", "setting.document_template.deactivate",
			"setting.document_template.version",
			"setting.employee_id_format.view", "setting.employee_id_format.update",
		},
		Menus: []module.Menu{
			{
				Name:  "Settings",
				Icon:  "cog",
				Route: "/admin/settings",
				Children: []module.Menu{
					{Name: "Zones", Icon: "map-pin", Route: "/admin/settings/zones"},
					{Name: "Provinces", Icon: "map", Route: "/admin/settings/provinces"},
					{Name: "Regencies", Icon: "map", Route: "/admin/settings/regencies"},
					{Name: "Districts", Icon: "map", Route: "/admin/settings/districts"},
					{Name: "Villages", Icon: "map", Route: "/admin/settings/villages"},
					{Name: "Educations", Icon: "graduation-cap", Route: "/admin/settings/educations"},
					{Name: "Education Majors", Icon: "graduation-cap", Route: "/admin/settings/education-majors"},
					{Name: "Religions", Icon: "globe", Route: "/admin/settings/religions"},
					{Name: "Marital Status", Icon: "heart", Route: "/admin/settings/marital-statuses"},
					{Name: "Relationship Types", Icon: "users", Route: "/admin/settings/relationship-types"},
					{Name: "Banks", Icon: "building-columns", Route: "/admin/settings/banks"},
					{Name: "Employment Statuses", Icon: "briefcase", Route: "/admin/settings/employment-statuses"},
					{Name: "Nationalities", Icon: "globe", Route: "/admin/settings/nationalities"},
					{Name: "Job Families", Icon: "briefcase", Route: "/admin/settings/job-families"},
					{Name: "Gradings", Icon: "chart-bar", Route: "/admin/settings/gradings"},
					{Name: "Salary Grades", Icon: "chart-bar", Route: "/admin/settings/salary-grades"},
					{Name: "TER (Tarif Efektif)", Icon: "calculator", Route: "/admin/settings/ters"},
					{Name: "PTKP", Icon: "receipt", Route: "/admin/settings/ptkps"},
					{Name: "Insurances", Icon: "shield", Route: "/admin/settings/insurances"},
					{Name: "Company Holidays", Icon: "calendar", Route: "/admin/settings/company-holidays"},
					{Name: "Competencies", Icon: "star", Route: "/admin/settings/competencies"},
					{Name: "Document Templates", Icon: "file", Route: "/admin/settings/document-templates"},
				},
			},
		},
	}
}

func (m *settingModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutesWithNumbering(rg, m.handler, m.numberingHandler, m.documentTemplateHandler, m.employeeIDFormatHandler)
}

func (m *settingModule) Migrate(db *gorm.DB) error {
	// Catatan: employee_id_format_settings TIDAK di-AutoMigrate di sini —
	// tabel tenant baru wajib lewat migrasi SQL bernomor (lihat
	// internal/pkg/migrator/migrations/tenant/{postgres,mysql}/155_*),
	// AutoMigrate tidak pernah dijalankan untuk modul tenant.
	return db.AutoMigrate(&Zone{}, &Province{}, &Regency{}, &District{}, &Village{}, &Education{}, &EducationMajor{}, &Religion{}, &MaritalStatus{}, &RelationshipType{}, &EmploymentStatus{}, &Bank{}, &Nationality{}, &JobFamily{}, &Grading{}, &SalaryGrade{}, &TER{}, &PTKP{}, &Insurance{}, &CompanyHoliday{}, &Competency{})
}

func (m *settingModule) Seed(db *gorm.DB) error {
	return nil
}

func (m *settingModule) Permissions() []string {
	return m.Info().Permissions
}
