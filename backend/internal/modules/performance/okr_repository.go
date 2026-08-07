package performance

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OKRRepository interface {
	// OKR Templates
	CreateOKRTemplate(db *gorm.DB, template *OKRTemplate) error
	GetOKRTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplate, error)
	GetOKRTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplate, error)
	ListOKRTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplate, int64, error)
	UpdateOKRTemplate(db *gorm.DB, template *OKRTemplate) error
	DeleteOKRTemplate(db *gorm.DB, id uuid.UUID) error

	// OKR Objectives
	CreateOKRObjective(db *gorm.DB, objective *OKRObjective) error
	GetOKRObjectiveByID(db *gorm.DB, id uuid.UUID) (*OKRObjective, error)
	GetOKRObjectiveWithKeyResults(db *gorm.DB, id uuid.UUID) (*OKRObjective, error)
	ListOKRObjectivesByTemplateID(db *gorm.DB, templateID uuid.UUID) ([]OKRObjective, error)
	UpdateOKRObjective(db *gorm.DB, objective *OKRObjective) error
	DeleteOKRObjective(db *gorm.DB, id uuid.UUID) error

	// OKR Key Results
	CreateOKRKeyResult(db *gorm.DB, keyResult *OKRKeyResult) error
	GetOKRKeyResultByID(db *gorm.DB, id uuid.UUID) (*OKRKeyResult, error)
	ListOKRKeyResultsByObjectiveID(db *gorm.DB, objectiveID uuid.UUID) ([]OKRKeyResult, error)
	UpdateOKRKeyResult(db *gorm.DB, keyResult *OKRKeyResult) error
	DeleteOKRKeyResult(db *gorm.DB, id uuid.UUID) error

	// OKR Evaluations
	CreateOKREvaluation(db *gorm.DB, evaluation *OKREvaluation) error
	GetOKREvaluationByID(db *gorm.DB, id uuid.UUID) (*OKREvaluation, error)
	GetOKREvaluationWithDetails(db *gorm.DB, id uuid.UUID) (*OKREvaluation, error)
	ListOKREvaluations(db *gorm.DB, employeeID, orgID, periodID *uuid.UUID, status *string, page, perPage int) ([]OKREvaluation, int64, error)
	UpdateOKREvaluation(db *gorm.DB, evaluation *OKREvaluation) error
	DeleteOKREvaluation(db *gorm.DB, id uuid.UUID) error

	// OKR Evaluation Details
	CreateOKREvaluationDetail(db *gorm.DB, detail *OKREvaluationDetail) error
	CreateOKREvaluationDetailsBatch(db *gorm.DB, details []OKREvaluationDetail) error
	GetOKREvaluationDetailByID(db *gorm.DB, id uuid.UUID) (*OKREvaluationDetail, error)
	ListOKREvaluationDetailsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKREvaluationDetail, error)
	UpdateOKREvaluationDetail(db *gorm.DB, detail *OKREvaluationDetail) error

	// OKR Progress
	CreateOKRProgress(db *gorm.DB, progress *OKRProgress) error
	GetOKRProgressByID(db *gorm.DB, id uuid.UUID) (*OKRProgress, error)
	ListOKRProgressByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRProgress, error)
	UpdateOKRProgress(db *gorm.DB, progress *OKRProgress) error
	DeleteOKRProgress(db *gorm.DB, id uuid.UUID) error

	// OKR Comments
	CreateOKRComment(db *gorm.DB, comment *OKRComment) error
	GetOKRCommentByID(db *gorm.DB, id uuid.UUID) (*OKRComment, error)
	ListOKRCommentsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKRComment, error)
	UpdateOKRComment(db *gorm.DB, comment *OKRComment) error
	DeleteOKRComment(db *gorm.DB, id uuid.UUID) error

	// OKR Attachments
	CreateOKRAttachment(db *gorm.DB, attachment *OKRAttachment) error
	GetOKRAttachmentByID(db *gorm.DB, id uuid.UUID) (*OKRAttachment, error)
	ListOKRAttachmentsByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRAttachment, error)
	DeleteOKRAttachment(db *gorm.DB, id uuid.UUID) error

	// Dashboard
	GetOKRHRDashboardStats(db *gorm.DB, periodID *uuid.UUID) (*OKRDashboardHRResponse, error)

	// My Context (self-assessment)
	GetCurrentEmployeeContext(db *gorm.DB, userID uuid.UUID) (employeeID *uuid.UUID, organizationID *uuid.UUID, err error)
	GetOrganizationName(db *gorm.DB, orgID uuid.UUID) (string, error)
}

type okrRepositoryImpl struct{}

func NewOKRRepository() OKRRepository {
	return &okrRepositoryImpl{}
}

// =========================================================================
// OKR Templates
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRTemplate(db *gorm.DB, template *OKRTemplate) error {
	return db.Create(template).Error
}

func (r *okrRepositoryImpl) GetOKRTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplate, error) {
	var template OKRTemplate
	if err := db.First(&template, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *okrRepositoryImpl) GetOKRTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplate, error) {
	var template OKRTemplate
	if err := db.Preload("Objectives", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Preload("Objectives.KeyResults", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&template, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *okrRepositoryImpl) ListOKRTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplate, int64, error) {
	var templates []OKRTemplate
	var total int64

	query := db.Model(&OKRTemplate{})
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	if periodID != nil {
		query = query.Where("period_id = ?", *periodID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

func (r *okrRepositoryImpl) UpdateOKRTemplate(db *gorm.DB, template *OKRTemplate) error {
	return db.Save(template).Error
}

func (r *okrRepositoryImpl) DeleteOKRTemplate(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRTemplate{}, "id = ?", id).Error
}

// =========================================================================
// OKR Objectives
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRObjective(db *gorm.DB, objective *OKRObjective) error {
	return db.Create(objective).Error
}

func (r *okrRepositoryImpl) GetOKRObjectiveByID(db *gorm.DB, id uuid.UUID) (*OKRObjective, error) {
	var objective OKRObjective
	if err := db.First(&objective, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &objective, nil
}

func (r *okrRepositoryImpl) GetOKRObjectiveWithKeyResults(db *gorm.DB, id uuid.UUID) (*OKRObjective, error) {
	var objective OKRObjective
	if err := db.Preload("KeyResults", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&objective, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &objective, nil
}

func (r *okrRepositoryImpl) ListOKRObjectivesByTemplateID(db *gorm.DB, templateID uuid.UUID) ([]OKRObjective, error) {
	var objectives []OKRObjective
	if err := db.Preload("KeyResults", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("template_id = ?", templateID).Order("sort_order ASC").Find(&objectives).Error; err != nil {
		return nil, err
	}
	return objectives, nil
}

func (r *okrRepositoryImpl) UpdateOKRObjective(db *gorm.DB, objective *OKRObjective) error {
	return db.Save(objective).Error
}

func (r *okrRepositoryImpl) DeleteOKRObjective(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRObjective{}, "id = ?", id).Error
}

// =========================================================================
// OKR Key Results
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRKeyResult(db *gorm.DB, keyResult *OKRKeyResult) error {
	return db.Create(keyResult).Error
}

func (r *okrRepositoryImpl) GetOKRKeyResultByID(db *gorm.DB, id uuid.UUID) (*OKRKeyResult, error) {
	var keyResult OKRKeyResult
	if err := db.First(&keyResult, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &keyResult, nil
}

func (r *okrRepositoryImpl) ListOKRKeyResultsByObjectiveID(db *gorm.DB, objectiveID uuid.UUID) ([]OKRKeyResult, error) {
	var keyResults []OKRKeyResult
	if err := db.Where("objective_id = ?", objectiveID).Order("sort_order ASC").Find(&keyResults).Error; err != nil {
		return nil, err
	}
	return keyResults, nil
}

func (r *okrRepositoryImpl) UpdateOKRKeyResult(db *gorm.DB, keyResult *OKRKeyResult) error {
	return db.Save(keyResult).Error
}

func (r *okrRepositoryImpl) DeleteOKRKeyResult(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRKeyResult{}, "id = ?", id).Error
}

// =========================================================================
// OKR Evaluations
// =========================================================================

func (r *okrRepositoryImpl) CreateOKREvaluation(db *gorm.DB, evaluation *OKREvaluation) error {
	return db.Create(evaluation).Error
}

func (r *okrRepositoryImpl) GetOKREvaluationByID(db *gorm.DB, id uuid.UUID) (*OKREvaluation, error) {
	var evaluation OKREvaluation
	if err := db.First(&evaluation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (r *okrRepositoryImpl) GetOKREvaluationWithDetails(db *gorm.DB, id uuid.UUID) (*OKREvaluation, error) {
	var evaluation OKREvaluation
	if err := db.Preload("Details", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).First(&evaluation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &evaluation, nil
}

func (r *okrRepositoryImpl) ListOKREvaluations(db *gorm.DB, employeeID, orgID, periodID *uuid.UUID, status *string, page, perPage int) ([]OKREvaluation, int64, error) {
	var evaluations []OKREvaluation
	var total int64

	query := db.Model(&OKREvaluation{})
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	if periodID != nil {
		query = query.Where("period_id = ?", *periodID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&evaluations).Error; err != nil {
		return nil, 0, err
	}

	return evaluations, total, nil
}

func (r *okrRepositoryImpl) UpdateOKREvaluation(db *gorm.DB, evaluation *OKREvaluation) error {
	return db.Save(evaluation).Error
}

func (r *okrRepositoryImpl) DeleteOKREvaluation(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKREvaluation{}, "id = ?", id).Error
}

// =========================================================================
// OKR Evaluation Details
// =========================================================================

func (r *okrRepositoryImpl) CreateOKREvaluationDetail(db *gorm.DB, detail *OKREvaluationDetail) error {
	return db.Create(detail).Error
}

func (r *okrRepositoryImpl) CreateOKREvaluationDetailsBatch(db *gorm.DB, details []OKREvaluationDetail) error {
	if len(details) == 0 {
		return nil
	}
	return db.Create(&details).Error
}

func (r *okrRepositoryImpl) GetOKREvaluationDetailByID(db *gorm.DB, id uuid.UUID) (*OKREvaluationDetail, error) {
	var detail OKREvaluationDetail
	if err := db.First(&detail, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *okrRepositoryImpl) ListOKREvaluationDetailsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKREvaluationDetail, error) {
	var details []OKREvaluationDetail
	if err := db.Where("evaluation_id = ?", evaluationID).Order("sort_order ASC").Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

func (r *okrRepositoryImpl) UpdateOKREvaluationDetail(db *gorm.DB, detail *OKREvaluationDetail) error {
	return db.Save(detail).Error
}

// =========================================================================
// OKR Progress
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRProgress(db *gorm.DB, progress *OKRProgress) error {
	return db.Create(progress).Error
}

func (r *okrRepositoryImpl) GetOKRProgressByID(db *gorm.DB, id uuid.UUID) (*OKRProgress, error) {
	var progress OKRProgress
	if err := db.First(&progress, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *okrRepositoryImpl) ListOKRProgressByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRProgress, error) {
	var progressList []OKRProgress
	if err := db.Where("evaluation_detail_id = ?", detailID).Order("progress_date DESC").Find(&progressList).Error; err != nil {
		return nil, err
	}
	return progressList, nil
}

func (r *okrRepositoryImpl) UpdateOKRProgress(db *gorm.DB, progress *OKRProgress) error {
	return db.Save(progress).Error
}

func (r *okrRepositoryImpl) DeleteOKRProgress(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRProgress{}, "id = ?", id).Error
}

// =========================================================================
// OKR Comments
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRComment(db *gorm.DB, comment *OKRComment) error {
	return db.Create(comment).Error
}

func (r *okrRepositoryImpl) GetOKRCommentByID(db *gorm.DB, id uuid.UUID) (*OKRComment, error) {
	var comment OKRComment
	if err := db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *okrRepositoryImpl) ListOKRCommentsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKRComment, error) {
	var comments []OKRComment
	if err := db.Preload("Replies").Where("evaluation_id = ? AND parent_id IS NULL", evaluationID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *okrRepositoryImpl) UpdateOKRComment(db *gorm.DB, comment *OKRComment) error {
	return db.Save(comment).Error
}

func (r *okrRepositoryImpl) DeleteOKRComment(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRComment{}, "id = ?", id).Error
}

// =========================================================================
// OKR Attachments
// =========================================================================

func (r *okrRepositoryImpl) CreateOKRAttachment(db *gorm.DB, attachment *OKRAttachment) error {
	return db.Create(attachment).Error
}

func (r *okrRepositoryImpl) GetOKRAttachmentByID(db *gorm.DB, id uuid.UUID) (*OKRAttachment, error) {
	var attachment OKRAttachment
	if err := db.First(&attachment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *okrRepositoryImpl) ListOKRAttachmentsByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRAttachment, error) {
	var attachments []OKRAttachment
	if err := db.Where("evaluation_detail_id = ?", detailID).Order("created_at DESC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *okrRepositoryImpl) DeleteOKRAttachment(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&OKRAttachment{}, "id = ?", id).Error
}

// =========================================================================
// Dashboard
// =========================================================================

func (r *okrRepositoryImpl) GetOKRHRDashboardStats(db *gorm.DB, periodID *uuid.UUID) (*OKRDashboardHRResponse, error) {
	var result OKRDashboardHRResponse

	query := db.Model(&OKREvaluation{})
	if periodID != nil {
		query = query.Where("period_id = ?", *periodID)
	}

	// Count by status
	var statusCounts []struct {
		Status string
		Count  int
	}
	if err := query.Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range statusCounts {
		result.TotalEvaluations += sc.Count
		switch sc.Status {
		case string(OKRStatusCompleted):
			result.CompletedCount = sc.Count
		case string(OKRStatusApproved):
			result.ApprovedCount = sc.Count
		case string(OKRStatusSubmitted):
			result.SubmittedCount = sc.Count
		case string(OKRStatusDraft):
			result.DraftCount = sc.Count
		}
	}

	// Average score
	var avgScore float64
	db.Model(&OKREvaluation{}).Where("status IN ?", []string{string(OKRStatusCompleted), string(OKRStatusApproved)}).Select("COALESCE(AVG(final_score), 0)").Scan(&avgScore)
	result.AverageScore = avgScore

	return &result, nil
}

// =========================================================================
// My Context (self-assessment)
// =========================================================================

// GetCurrentEmployeeContext resolves the current (open-ended) employee_id and
// organization_id for a logged-in user — same join pattern as the KPI
// self-assessment context resolver, duplicated here because the OKR
// repository takes *gorm.DB directly rather than a context-based resolver.
func (r *okrRepositoryImpl) GetCurrentEmployeeContext(db *gorm.DB, userID uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {
	type row struct {
		EmployeeID     string
		OrganizationID string
	}
	var result row
	err := db.Table("employee_accounts AS ea").
		Joins("JOIN employments AS emp ON emp.employee_id = ea.employee_id").
		Where("ea.user_id = ? AND emp.effective_end_date IS NULL", userID).
		Order("emp.effective_date DESC").
		Limit(1).
		Select("ea.employee_id AS employee_id, emp.organization_id AS organization_id").
		Scan(&result).Error
	if err != nil {
		return nil, nil, err
	}
	if result.EmployeeID == "" || result.OrganizationID == "" {
		return nil, nil, nil
	}
	empID, err := uuid.Parse(result.EmployeeID)
	if err != nil {
		return nil, nil, err
	}
	orgID, err := uuid.Parse(result.OrganizationID)
	if err != nil {
		return nil, nil, err
	}
	return &empID, &orgID, nil
}

// GetOrganizationName mengambil nomenclature Organization via raw table
// query (tanpa import package organization, hindari circular dependency).
func (r *okrRepositoryImpl) GetOrganizationName(db *gorm.DB, orgID uuid.UUID) (string, error) {
	var name string
	if err := db.Table("organizations").
		Where("id = ?", orgID).
		Pluck("nomenclature", &name).Error; err != nil {
		return "", err
	}
	return name, nil
}
