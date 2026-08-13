package recruitment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/inthros/hris-platform/internal/modules/competency"
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

// =========================================================================
// Candidate Skills (G-6)
// =========================================================================

func (r *Repository) CreateCandidateSkill(ctx context.Context, s *CandidateSkill) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindCandidateSkillByID(ctx context.Context, id uuid.UUID) (*CandidateSkill, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s CandidateSkill
	if err := db.WithContext(ctx).Preload("Competency").First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate skill not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListCandidateSkills(ctx context.Context, candidateID uuid.UUID) ([]CandidateSkill, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateSkill
	if err := db.WithContext(ctx).Preload("Competency").Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateSkill(ctx context.Context, s *CandidateSkill) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	// Omit associations so a Save never re-derives CompetencyID from a stale
	// preloaded Competency pointer (FindCandidateSkillByID preloads Competency).
	return db.WithContext(ctx).Omit(clause.Associations).Save(s).Error
}

func (r *Repository) DeleteCandidateSkill(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateSkill{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate skill not found")
	}
	return result.Error
}

func (r *Repository) FindCompetencyByID(ctx context.Context, id uuid.UUID) (*competency.Competency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c competency.Competency
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("competency not found")
		}
		return nil, err
	}
	return &c, nil
}

// =========================================================================
// Candidate Certifications (G-6)
// =========================================================================

func (r *Repository) CreateCandidateCertification(ctx context.Context, c *CandidateCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCandidateCertificationByID(ctx context.Context, id uuid.UUID) (*CandidateCertification, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c CandidateCertification
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate certification not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCandidateCertifications(ctx context.Context, candidateID uuid.UUID) ([]CandidateCertification, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateCertification
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateCertification(ctx context.Context, c *CandidateCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCandidateCertification(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateCertification{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate certification not found")
	}
	return result.Error
}

// =========================================================================
// Candidate Documents (G-6)
// =========================================================================

func (r *Repository) CreateCandidateDocument(ctx context.Context, d *CandidateDocument) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(d).Error
}

func (r *Repository) FindCandidateDocumentByID(ctx context.Context, id uuid.UUID) (*CandidateDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var d CandidateDocument
	if err := db.WithContext(ctx).First(&d, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate document not found")
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repository) ListCandidateDocuments(ctx context.Context, candidateID uuid.UUID) ([]CandidateDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateDocument
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateDocument(ctx context.Context, d *CandidateDocument) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(d).Error
}

func (r *Repository) DeleteCandidateDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateDocument{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate document not found")
	}
	return result.Error
}

// =========================================================================
// Candidate Consents (G-6) — append-only, no Update/Delete
// =========================================================================

func (r *Repository) CreateCandidateConsent(ctx context.Context, c *CandidateConsent) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) ListCandidateConsents(ctx context.Context, candidateID uuid.UUID) ([]CandidateConsent, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateConsent
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("changed_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Application Screenings (G-7 sub-project 1)
// =========================================================================

func (r *Repository) CreateApplicationScreening(ctx context.Context, s *ApplicationScreening) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindApplicationScreeningByID(ctx context.Context, id uuid.UUID) (*ApplicationScreening, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s ApplicationScreening
	if err := db.WithContext(ctx).First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("application screening not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListApplicationScreenings(ctx context.Context, applicationID uuid.UUID) ([]ApplicationScreening, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []ApplicationScreening
	if err := db.WithContext(ctx).Where("application_id = ?", applicationID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateApplicationScreening(ctx context.Context, s *ApplicationScreening) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(s).Error
}

func (r *Repository) DeleteApplicationScreening(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&ApplicationScreening{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("application screening not found")
	}
	return result.Error
}

// =========================================================================
// Recruitment Assessments + Participants (G-7 sub-project 2)
// =========================================================================

func (r *Repository) CreateAssessment(ctx context.Context, a *RecruitmentAssessment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(a).Error
}

func (r *Repository) FindAssessmentByID(ctx context.Context, id uuid.UUID) (*RecruitmentAssessment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var a RecruitmentAssessment
	if err := db.WithContext(ctx).First(&a, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListAssessments(ctx context.Context) ([]RecruitmentAssessment, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []RecruitmentAssessment
	if err := db.WithContext(ctx).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateAssessment(ctx context.Context, a *RecruitmentAssessment) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(a).Error
}

func (r *Repository) DeleteAssessment(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&RecruitmentAssessment{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("assessment not found")
	}
	return result.Error
}

func (r *Repository) CreateAssessmentParticipant(ctx context.Context, p *AssessmentParticipant) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(p).Error
}

func (r *Repository) FindAssessmentParticipantByID(ctx context.Context, id uuid.UUID) (*AssessmentParticipant, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var p AssessmentParticipant
	if err := db.WithContext(ctx).First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment participant not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListAssessmentParticipants(ctx context.Context, assessmentID uuid.UUID) ([]AssessmentParticipant, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []AssessmentParticipant
	if err := db.WithContext(ctx).Where("assessment_id = ?", assessmentID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateAssessmentParticipant(ctx context.Context, p *AssessmentParticipant) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(p).Error
}

func (r *Repository) DeleteAssessmentParticipant(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&AssessmentParticipant{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("assessment participant not found")
	}
	return result.Error
}

// =========================================================================
// Job Requisition Requirements + Competencies (G-9 sub-project 1)
// =========================================================================

func (r *Repository) CreateRequisitionRequirement(ctx context.Context, req *JobRequisitionRequirement) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(req).Error
}

func (r *Repository) FindRequisitionRequirementByID(ctx context.Context, id uuid.UUID) (*JobRequisitionRequirement, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var req JobRequisitionRequirement
	if err := db.WithContext(ctx).First(&req, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("requisition requirement not found")
		}
		return nil, err
	}
	return &req, nil
}

func (r *Repository) ListRequisitionRequirements(ctx context.Context, requisitionID uuid.UUID) ([]JobRequisitionRequirement, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []JobRequisitionRequirement
	if err := db.WithContext(ctx).Where("requisition_id = ?", requisitionID).Order("sort_order ASC, created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateRequisitionRequirement(ctx context.Context, req *JobRequisitionRequirement) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(req).Error
}

func (r *Repository) DeleteRequisitionRequirement(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&JobRequisitionRequirement{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("requisition requirement not found")
	}
	return result.Error
}

func (r *Repository) CreateRequisitionCompetency(ctx context.Context, c *JobRequisitionCompetency) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindRequisitionCompetencyByID(ctx context.Context, id uuid.UUID) (*JobRequisitionCompetency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c JobRequisitionCompetency
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("requisition competency not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListRequisitionCompetencies(ctx context.Context, requisitionID uuid.UUID) ([]JobRequisitionCompetency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []JobRequisitionCompetency
	if err := db.WithContext(ctx).Where("requisition_id = ?", requisitionID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateRequisitionCompetency(ctx context.Context, c *JobRequisitionCompetency) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteRequisitionCompetency(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&JobRequisitionCompetency{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("requisition competency not found")
	}
	return result.Error
}

// =========================================================================
// Interviewers + Scorecard Items (G-8)
// =========================================================================

func (r *Repository) CreateInterviewer(ctx context.Context, i *Interviewer) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(i).Error
}

func (r *Repository) ListInterviewers(ctx context.Context, interviewID uuid.UUID) ([]Interviewer, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []Interviewer
	if err := db.WithContext(ctx).Where("interview_id = ?", interviewID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) DeleteInterviewer(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&Interviewer{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("interviewer not found")
	}
	return result.Error
}

func (r *Repository) CreateScorecardItem(ctx context.Context, s *InterviewScorecardItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindScorecardItemByID(ctx context.Context, id uuid.UUID) (*InterviewScorecardItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s InterviewScorecardItem
	if err := db.WithContext(ctx).First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("scorecard item not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListScorecardItems(ctx context.Context, interviewID uuid.UUID) ([]InterviewScorecardItem, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []InterviewScorecardItem
	if err := db.WithContext(ctx).Where("interview_id = ?", interviewID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateScorecardItem(ctx context.Context, s *InterviewScorecardItem) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(s).Error
}

func (r *Repository) DeleteScorecardItem(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&InterviewScorecardItem{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("scorecard item not found")
	}
	return result.Error
}
