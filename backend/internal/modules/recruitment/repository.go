package recruitment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	dbFunc func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbFunc func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbFunc: dbFunc}
}

func (r *Repository) db(ctx context.Context) (*gorm.DB, error) {
	return r.dbFunc(ctx)
}

// =========================================================================
// Job Requisitions
// =========================================================================

func (r *Repository) CreateRequisition(ctx context.Context, req *JobRequisition) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(req).Error
}

func (r *Repository) FindRequisitionByID(ctx context.Context, id uuid.UUID) (*JobRequisition, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var req JobRequisition
	if err := db.WithContext(ctx).First(&req, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("requisition not found")
		}
		return nil, err
	}
	return &req, nil
}

func (r *Repository) ListRequisitions(ctx context.Context, orgID *uuid.UUID, status *string, page, perPage int) ([]JobRequisition, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []JobRequisition
	var total int64
	query := db.WithContext(ctx).Model(&JobRequisition{})
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateRequisition(ctx context.Context, req *JobRequisition) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(req).Error
}

func (r *Repository) DeleteRequisition(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&JobRequisition{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("requisition not found")
	}
	return result.Error
}

// =========================================================================
// Job Offers (G-3)
// =========================================================================

func (r *Repository) CreateOffer(ctx context.Context, o *JobOffer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(o).Error
}

func (r *Repository) FindOfferByID(ctx context.Context, id uuid.UUID) (*JobOffer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var o JobOffer
	if err := db.WithContext(ctx).First(&o, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("offer not found")
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) ListOffers(ctx context.Context, applicationID *uuid.UUID, status *string, page, perPage int) ([]JobOffer, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []JobOffer
	var total int64
	query := db.WithContext(ctx).Model(&JobOffer{})
	if applicationID != nil {
		query = query.Where("application_id = ?", *applicationID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateOffer(ctx context.Context, o *JobOffer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(o).Error
}

func (r *Repository) DeleteOffer(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&JobOffer{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("offer not found")
	}
	return result.Error
}

// =========================================================================
// Candidates
// =========================================================================

func (r *Repository) CreateCandidate(ctx context.Context, c *Candidate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCandidateByID(ctx context.Context, id uuid.UUID) (*Candidate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c Candidate
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindCandidateByEmail(ctx context.Context, email string) (*Candidate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c Candidate
	if err := db.WithContext(ctx).Where("email = ?", email).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCandidates(ctx context.Context, page, perPage int, search *string) ([]Candidate, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []Candidate
	var total int64
	query := db.WithContext(ctx).Model(&Candidate{})
	if search != nil && *search != "" {
		s := "%" + *search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", s, s, s)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateCandidate(ctx context.Context, c *Candidate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCandidate(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&Candidate{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate not found")
	}
	return result.Error
}

// =========================================================================
// Job Applications
// =========================================================================

func (r *Repository) CreateApplication(ctx context.Context, a *JobApplication) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindApplicationByID(ctx context.Context, id uuid.UUID) (*JobApplication, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a JobApplication
	if err := db.WithContext(ctx).First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("application not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListApplications(ctx context.Context, requisitionID, candidateID *uuid.UUID, status *string, page, perPage int) ([]JobApplication, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []JobApplication
	var total int64
	query := db.WithContext(ctx).Model(&JobApplication{})
	if requisitionID != nil {
		query = query.Where("requisition_id = ?", *requisitionID)
	}
	if candidateID != nil {
		query = query.Where("candidate_id = ?", *candidateID)
	}
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateApplication(ctx context.Context, a *JobApplication) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(a).Error
}

func (r *Repository) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&JobApplication{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("application not found")
	}
	return result.Error
}

// =========================================================================
// Interviews
// =========================================================================

func (r *Repository) CreateInterview(ctx context.Context, i *Interview) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(i).Error
}

func (r *Repository) FindInterviewByID(ctx context.Context, id uuid.UUID) (*Interview, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var i Interview
	if err := db.WithContext(ctx).First(&i, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("interview not found")
		}
		return nil, err
	}
	return &i, nil
}

func (r *Repository) ListInterviews(ctx context.Context, applicationID *uuid.UUID, interviewerID *uuid.UUID, page, perPage int) ([]Interview, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []Interview
	var total int64
	query := db.WithContext(ctx).Model(&Interview{})
	if applicationID != nil {
		query = query.Where("application_id = ?", *applicationID)
	}
	if interviewerID != nil {
		query = query.Where("interviewer_id = ?", *interviewerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("scheduled_at ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateInterview(ctx context.Context, i *Interview) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(i).Error
}

func (r *Repository) DeleteInterview(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&Interview{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("interview not found")
	}
	return result.Error
}

// =========================================================================
// Onboarding Task Templates
// =========================================================================

func (r *Repository) CreateOnboardingTaskTemplate(ctx context.Context, t *OnboardingTaskTemplate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindOnboardingTaskTemplateByID(ctx context.Context, id uuid.UUID) (*OnboardingTaskTemplate, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t OnboardingTaskTemplate
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("onboarding task template not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListOnboardingTaskTemplates(ctx context.Context, category *string, page, perPage int) ([]OnboardingTaskTemplate, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []OnboardingTaskTemplate
	var total int64
	query := db.WithContext(ctx).Model(&OnboardingTaskTemplate{})
	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("day_offset ASC, name ASC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateOnboardingTaskTemplate(ctx context.Context, t *OnboardingTaskTemplate) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeleteOnboardingTaskTemplate(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&OnboardingTaskTemplate{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("onboarding task template not found")
	}
	return result.Error
}

// =========================================================================
// Employee Onboarding
// =========================================================================

func (r *Repository) CreateEmployeeOnboarding(ctx context.Context, o *EmployeeOnboarding) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(o).Error
}

func (r *Repository) FindEmployeeOnboardingByID(ctx context.Context, id uuid.UUID) (*EmployeeOnboarding, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var o EmployeeOnboarding
	if err := db.WithContext(ctx).First(&o, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee onboarding not found")
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) FindEmployeeOnboardingByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*EmployeeOnboarding, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var o EmployeeOnboarding
	if err := db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&o).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) ListEmployeeOnboardings(ctx context.Context, status *string, page, perPage int) ([]EmployeeOnboarding, int64, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, 0, err
	}
	var list []EmployeeOnboarding
	var total int64
	query := db.WithContext(ctx).Model(&EmployeeOnboarding{})
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("start_date ASC, created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *Repository) UpdateEmployeeOnboarding(ctx context.Context, o *EmployeeOnboarding) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(o).Error
}

func (r *Repository) DeleteEmployeeOnboarding(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&EmployeeOnboarding{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee onboarding not found")
	}
	return result.Error
}

// =========================================================================
// Onboarding Task Items
// =========================================================================

func (r *Repository) CreateOnboardingTaskItem(ctx context.Context, t *OnboardingTaskItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(t).Error
}

func (r *Repository) FindOnboardingTaskItemByID(ctx context.Context, id uuid.UUID) (*OnboardingTaskItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var t OnboardingTaskItem
	if err := db.WithContext(ctx).First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("onboarding task item not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListOnboardingTaskItems(ctx context.Context, onboardingID uuid.UUID) ([]OnboardingTaskItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []OnboardingTaskItem
	if err := db.WithContext(ctx).Where("employee_onboarding_id = ?", onboardingID).
		Order("due_date ASC, created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateOnboardingTaskItem(ctx context.Context, t *OnboardingTaskItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(t).Error
}

func (r *Repository) DeleteOnboardingTaskItem(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&OnboardingTaskItem{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("onboarding task item not found")
	}
	return result.Error
}

// =========================================================================
// Recruitment Stages (G-5 — master, seeded)
// =========================================================================

func (r *Repository) CreateStage(ctx context.Context, s *RecruitmentStage) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindStageByCode(ctx context.Context, code string) (*RecruitmentStage, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s RecruitmentStage
	if err := db.WithContext(ctx).Where("code = ?", code).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("recruitment stage not found: %s", code)
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListStages(ctx context.Context) ([]RecruitmentStage, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []RecruitmentStage
	if err := db.WithContext(ctx).Order("sort_order ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Application Stage History (G-5)
// =========================================================================

func (r *Repository) CreateStageHistory(ctx context.Context, h *ApplicationStageHistory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(h).Error
}

func (r *Repository) ListStageHistoryByApplication(ctx context.Context, applicationID uuid.UUID) ([]ApplicationStageHistory, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []ApplicationStageHistory
	if err := db.WithContext(ctx).Where("application_id = ?", applicationID).Order("changed_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Candidate Educations (G-6)
// =========================================================================

func (r *Repository) CreateCandidateEducation(ctx context.Context, e *CandidateEducation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindCandidateEducationByID(ctx context.Context, id uuid.UUID) (*CandidateEducation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e CandidateEducation
	if err := db.WithContext(ctx).Preload("EducationMajor").First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate education not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListCandidateEducations(ctx context.Context, candidateID uuid.UUID) ([]CandidateEducation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateEducation
	if err := db.WithContext(ctx).Preload("EducationMajor").Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateEducation(ctx context.Context, e *CandidateEducation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	// Omit associations so a Save never writes to setting.Education/EducationMajor
	// (owned by another module) and never re-derives a FK from a stale preloaded
	// association (defense in depth alongside the service-layer nil-out).
	return db.WithContext(ctx).Omit(clause.Associations).Save(e).Error
}

func (r *Repository) DeleteCandidateEducation(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateEducation{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate education not found")
	}
	return result.Error
}

// =========================================================================
// Candidate Work Experiences (G-6)
// =========================================================================

func (r *Repository) CreateCandidateWorkExperience(ctx context.Context, e *CandidateWorkExperience) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindCandidateWorkExperienceByID(ctx context.Context, id uuid.UUID) (*CandidateWorkExperience, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e CandidateWorkExperience
	if err := db.WithContext(ctx).First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate work experience not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListCandidateWorkExperiences(ctx context.Context, candidateID uuid.UUID) ([]CandidateWorkExperience, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateWorkExperience
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("start_date DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateWorkExperience(ctx context.Context, e *CandidateWorkExperience) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(e).Error
}

func (r *Repository) DeleteCandidateWorkExperience(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateWorkExperience{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate work experience not found")
	}
	return result.Error
}
