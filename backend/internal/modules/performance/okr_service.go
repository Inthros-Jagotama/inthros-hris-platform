package performance

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OKRService interface {
	// OKR Templates
	CreateTemplate(db *gorm.DB, req *CreateOKRTemplateRequest) (*OKRTemplateResponse, error)
	GetTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error)
	GetTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error)
	ListTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplateResponse, int64, error)
	UpdateTemplate(db *gorm.DB, id uuid.UUID, req *UpdateOKRTemplateRequest) (*OKRTemplateResponse, error)
	DeleteTemplate(db *gorm.DB, id uuid.UUID) error
	DuplicateTemplate(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error)

	// OKR Objectives
	CreateObjective(db *gorm.DB, req *CreateOKRObjectiveRequest) (*OKRObjectiveResponse, error)
	GetObjectiveByID(db *gorm.DB, id uuid.UUID) (*OKRObjectiveResponse, error)
	ListObjectivesByTemplateID(db *gorm.DB, templateID uuid.UUID) ([]OKRObjectiveResponse, error)
	UpdateObjective(db *gorm.DB, id uuid.UUID, req *UpdateOKRObjectiveRequest) (*OKRObjectiveResponse, error)
	DeleteObjective(db *gorm.DB, id uuid.UUID) error

	// OKR Key Results
	CreateKeyResult(db *gorm.DB, req *CreateOKRKeyResultRequest) (*OKRKeyResultResponse, error)
	GetKeyResultByID(db *gorm.DB, id uuid.UUID) (*OKRKeyResultResponse, error)
	ListKeyResultsByObjectiveID(db *gorm.DB, objectiveID uuid.UUID) ([]OKRKeyResultResponse, error)
	UpdateKeyResult(db *gorm.DB, id uuid.UUID, req *UpdateOKRKeyResultRequest) (*OKRKeyResultResponse, error)
	DeleteKeyResult(db *gorm.DB, id uuid.UUID) error

	// OKR Evaluations
	CreateEvaluationWithSnapshot(db *gorm.DB, req *CreateOKREvaluationRequest) (*OKREvaluationResponse, error)
	GetEvaluationByID(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error)
	GetEvaluationWithDetails(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error)
	ListEvaluations(db *gorm.DB, employeeID, orgID, periodID *uuid.UUID, status *string, page, perPage int) ([]OKREvaluationResponse, int64, error)
	UpdateEvaluation(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationRequest) (*OKREvaluationResponse, error)
	DeleteEvaluation(db *gorm.DB, id uuid.UUID) error

	// Evaluation Detail & Score
	UpdateEvaluationDetailActual(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationDetailRequest) (*OKREvaluationDetailResponse, error)
	BulkUpdateEvaluationActuals(db *gorm.DB, evaluationID uuid.UUID, req *OKRBulkUpdateActualsRequest) error
	RecalculateEvaluationScore(db *gorm.DB, evaluationID uuid.UUID) (*OKREvaluationResponse, error)

	// Workflow
	SubmitEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error)
	ApproveEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error)
	RejectEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID, notes string) (*OKREvaluationResponse, error)
	CompleteEvaluation(db *gorm.DB, evaluationID uuid.UUID) (*OKREvaluationResponse, error)

	// OKR Progress
	CreateProgress(db *gorm.DB, req *CreateOKRProgressRequest, userID uuid.UUID) (*OKRProgressResponse, error)
	GetProgressByID(db *gorm.DB, id uuid.UUID) (*OKRProgressResponse, error)
	ListProgressByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRProgressResponse, error)
	UpdateProgress(db *gorm.DB, id uuid.UUID, req *UpdateOKRProgressRequest) (*OKRProgressResponse, error)
	DeleteProgress(db *gorm.DB, id uuid.UUID) error

	// OKR Comments
	CreateComment(db *gorm.DB, req *CreateOKRCommentRequest, userID uuid.UUID) (*OKRCommentResponse, error)
	ListCommentsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKRCommentResponse, error)
	UpdateComment(db *gorm.DB, id uuid.UUID, req *UpdateOKRCommentRequest) (*OKRCommentResponse, error)
	DeleteComment(db *gorm.DB, id uuid.UUID) error

	// OKR Attachments
	CreateAttachment(db *gorm.DB, req *CreateOKRAttachmentRequest, userID uuid.UUID) (*OKRAttachmentResponse, error)
	ListAttachmentsByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRAttachmentResponse, error)
	DeleteAttachment(db *gorm.DB, id uuid.UUID) error

	// Dashboard
	GetHRDashboard(db *gorm.DB, periodID *uuid.UUID) (*OKRDashboardHRResponse, error)

	// My Context (self-assessment)
	GetMyOKRContext(db *gorm.DB, userID uuid.UUID) (*MyOKRContextResponse, error)
}

type okrServiceImpl struct {
	repo OKRRepository
}

func NewOKRService(repo OKRRepository) OKRService {
	return &okrServiceImpl{repo: repo}
}

// =========================================================================
// OKR Templates
// =========================================================================

func (s *okrServiceImpl) CreateTemplate(db *gorm.DB, req *CreateOKRTemplateRequest) (*OKRTemplateResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, errors.New("invalid organization_id")
	}

	template := &OKRTemplate{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		EffectiveDate:  req.EffectiveDate,
		ExpiredDate:    req.ExpiredDate,
	}

	if req.PeriodID != nil {
		periodID, err := uuid.Parse(*req.PeriodID)
		if err != nil {
			return nil, errors.New("invalid period_id")
		}
		template.PeriodID = &periodID
	}

	if req.Status != nil {
		template.Status = *req.Status
	}

	if err := s.repo.CreateOKRTemplate(db, template); err != nil {
		return nil, err
	}

	return s.templateToResponse(template), nil
}

func (s *okrServiceImpl) GetTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateByID(db, id)
	if err != nil {
		return nil, err
	}
	return s.templateToResponse(template), nil
}

func (s *okrServiceImpl) GetTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateWithObjectives(db, id)
	if err != nil {
		return nil, err
	}
	return s.templateToResponseWithObjectives(template), nil
}

func (s *okrServiceImpl) ListTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplateResponse, int64, error) {
	templates, total, err := s.repo.ListOKRTemplates(db, orgID, periodID, status, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]OKRTemplateResponse, len(templates))
	for i, t := range templates {
		responses[i] = *s.templateToResponse(&t)
	}

	return responses, total, nil
}

func (s *okrServiceImpl) UpdateTemplate(db *gorm.DB, id uuid.UUID, req *UpdateOKRTemplateRequest) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateByID(db, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = req.Description
	}
	if req.Status != nil {
		template.Status = *req.Status
	}
	if req.EffectiveDate != nil {
		template.EffectiveDate = req.EffectiveDate
	}
	if req.ExpiredDate != nil {
		template.ExpiredDate = req.ExpiredDate
	}
	if req.PeriodID != nil {
		periodID, err := uuid.Parse(*req.PeriodID)
		if err != nil {
			return nil, errors.New("invalid period_id")
		}
		template.PeriodID = &periodID
	}

	if err := s.repo.UpdateOKRTemplate(db, template); err != nil {
		return nil, err
	}

	return s.templateToResponse(template), nil
}

func (s *okrServiceImpl) DeleteTemplate(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRTemplate(db, id)
}

func (s *okrServiceImpl) DuplicateTemplate(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error) {
	original, err := s.repo.GetOKRTemplateWithObjectives(db, id)
	if err != nil {
		return nil, err
	}

	newTemplate := &OKRTemplate{
		OrganizationID: original.OrganizationID,
		PeriodID:       original.PeriodID,
		Name:           original.Name + " (Copy)",
		Description:    original.Description,
		Status:         0,
		EffectiveDate:  original.EffectiveDate,
		ExpiredDate:    original.ExpiredDate,
	}

	if err := s.repo.CreateOKRTemplate(db, newTemplate); err != nil {
		return nil, err
	}

	for _, obj := range original.Objectives {
		newObj := &OKRObjective{
			TemplateID:  newTemplate.ID,
			Code:        obj.Code,
			Title:       obj.Title,
			Description: obj.Description,
			Weight:      obj.Weight,
			SortOrder:   obj.SortOrder,
		}
		if err := s.repo.CreateOKRObjective(db, newObj); err != nil {
			return nil, err
		}

		for _, kr := range obj.KeyResults {
			newKR := &OKRKeyResult{
				ObjectiveID:  newObj.ID,
				Code:         kr.Code,
				Title:        kr.Title,
				Description:  kr.Description,
				TargetType:   kr.TargetType,
				TargetValue:  kr.TargetValue,
				Unit:         kr.Unit,
				FormulaType:  kr.FormulaType,
				Weight:       kr.Weight,
				MinimumScore: kr.MinimumScore,
				MaximumScore: kr.MaximumScore,
				SortOrder:    kr.SortOrder,
				IsRequired:   kr.IsRequired,
			}
			if err := s.repo.CreateOKRKeyResult(db, newKR); err != nil {
				return nil, err
			}
		}
	}

	return s.GetTemplateWithObjectives(db, newTemplate.ID)
}

// =========================================================================
// OKR Objectives
// =========================================================================

func (s *okrServiceImpl) CreateObjective(db *gorm.DB, req *CreateOKRObjectiveRequest) (*OKRObjectiveResponse, error) {
	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, errors.New("invalid template_id")
	}

	objective := &OKRObjective{
		TemplateID:  templateID,
		Code:        req.Code,
		Title:       req.Title,
		Description: req.Description,
		Weight:      req.Weight,
	}

	if req.SortOrder != nil {
		objective.SortOrder = *req.SortOrder
	}

	if err := s.repo.CreateOKRObjective(db, objective); err != nil {
		return nil, err
	}

	return s.objectiveToResponse(objective), nil
}

func (s *okrServiceImpl) GetObjectiveByID(db *gorm.DB, id uuid.UUID) (*OKRObjectiveResponse, error) {
	objective, err := s.repo.GetOKRObjectiveWithKeyResults(db, id)
	if err != nil {
		return nil, err
	}
	return s.objectiveToResponseWithKeyResults(objective), nil
}

func (s *okrServiceImpl) ListObjectivesByTemplateID(db *gorm.DB, templateID uuid.UUID) ([]OKRObjectiveResponse, error) {
	objectives, err := s.repo.ListOKRObjectivesByTemplateID(db, templateID)
	if err != nil {
		return nil, err
	}

	responses := make([]OKRObjectiveResponse, len(objectives))
	for i, o := range objectives {
		responses[i] = *s.objectiveToResponseWithKeyResults(&o)
	}

	return responses, nil
}

func (s *okrServiceImpl) UpdateObjective(db *gorm.DB, id uuid.UUID, req *UpdateOKRObjectiveRequest) (*OKRObjectiveResponse, error) {
	objective, err := s.repo.GetOKRObjectiveByID(db, id)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		objective.Code = req.Code
	}
	if req.Title != nil {
		objective.Title = *req.Title
	}
	if req.Description != nil {
		objective.Description = req.Description
	}
	if req.Weight != nil {
		objective.Weight = *req.Weight
	}
	if req.SortOrder != nil {
		objective.SortOrder = *req.SortOrder
	}

	if err := s.repo.UpdateOKRObjective(db, objective); err != nil {
		return nil, err
	}

	return s.objectiveToResponse(objective), nil
}

func (s *okrServiceImpl) DeleteObjective(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRObjective(db, id)
}

// =========================================================================
// OKR Key Results
// =========================================================================

func (s *okrServiceImpl) CreateKeyResult(db *gorm.DB, req *CreateOKRKeyResultRequest) (*OKRKeyResultResponse, error) {
	objectiveID, err := uuid.Parse(req.ObjectiveID)
	if err != nil {
		return nil, errors.New("invalid objective_id")
	}

	keyResult := &OKRKeyResult{
		ObjectiveID: objectiveID,
		Code:        req.Code,
		Title:       req.Title,
		Description: req.Description,
		TargetValue: req.TargetValue,
		Unit:        req.Unit,
		Weight:      req.Weight,
	}

	if req.TargetType != nil {
		keyResult.TargetType = TargetType(*req.TargetType)
	}
	if req.FormulaType != nil {
		keyResult.FormulaType = FormulaType(*req.FormulaType)
	}
	if req.MinimumScore != nil {
		keyResult.MinimumScore = *req.MinimumScore
	}
	if req.MaximumScore != nil {
		keyResult.MaximumScore = *req.MaximumScore
	}
	if req.SortOrder != nil {
		keyResult.SortOrder = *req.SortOrder
	}
	if req.IsRequired != nil {
		keyResult.IsRequired = *req.IsRequired
	}

	if err := s.repo.CreateOKRKeyResult(db, keyResult); err != nil {
		return nil, err
	}

	return s.keyResultToResponse(keyResult), nil
}

func (s *okrServiceImpl) GetKeyResultByID(db *gorm.DB, id uuid.UUID) (*OKRKeyResultResponse, error) {
	keyResult, err := s.repo.GetOKRKeyResultByID(db, id)
	if err != nil {
		return nil, err
	}
	return s.keyResultToResponse(keyResult), nil
}

func (s *okrServiceImpl) ListKeyResultsByObjectiveID(db *gorm.DB, objectiveID uuid.UUID) ([]OKRKeyResultResponse, error) {
	keyResults, err := s.repo.ListOKRKeyResultsByObjectiveID(db, objectiveID)
	if err != nil {
		return nil, err
	}

	responses := make([]OKRKeyResultResponse, len(keyResults))
	for i, kr := range keyResults {
		responses[i] = *s.keyResultToResponse(&kr)
	}

	return responses, nil
}

func (s *okrServiceImpl) UpdateKeyResult(db *gorm.DB, id uuid.UUID, req *UpdateOKRKeyResultRequest) (*OKRKeyResultResponse, error) {
	keyResult, err := s.repo.GetOKRKeyResultByID(db, id)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		keyResult.Code = req.Code
	}
	if req.Title != nil {
		keyResult.Title = *req.Title
	}
	if req.Description != nil {
		keyResult.Description = req.Description
	}
	if req.TargetType != nil {
		keyResult.TargetType = TargetType(*req.TargetType)
	}
	if req.TargetValue != nil {
		keyResult.TargetValue = *req.TargetValue
	}
	if req.Unit != nil {
		keyResult.Unit = req.Unit
	}
	if req.FormulaType != nil {
		keyResult.FormulaType = FormulaType(*req.FormulaType)
	}
	if req.Weight != nil {
		keyResult.Weight = *req.Weight
	}
	if req.MinimumScore != nil {
		keyResult.MinimumScore = *req.MinimumScore
	}
	if req.MaximumScore != nil {
		keyResult.MaximumScore = *req.MaximumScore
	}
	if req.SortOrder != nil {
		keyResult.SortOrder = *req.SortOrder
	}
	if req.IsRequired != nil {
		keyResult.IsRequired = *req.IsRequired
	}

	if err := s.repo.UpdateOKRKeyResult(db, keyResult); err != nil {
		return nil, err
	}

	return s.keyResultToResponse(keyResult), nil
}

func (s *okrServiceImpl) DeleteKeyResult(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRKeyResult(db, id)
}

// =========================================================================
// OKR Evaluations
// =========================================================================

func (s *okrServiceImpl) CreateEvaluationWithSnapshot(db *gorm.DB, req *CreateOKREvaluationRequest) (*OKREvaluationResponse, error) {
	employeeID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, errors.New("invalid employee_id")
	}
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, errors.New("invalid organization_id")
	}
	periodID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		return nil, errors.New("invalid period_id")
	}
	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, errors.New("invalid template_id")
	}

	template, err := s.repo.GetOKRTemplateWithObjectives(db, templateID)
	if err != nil {
		return nil, errors.New("template not found")
	}

	evaluation := &OKREvaluation{
		EmployeeID:     employeeID,
		OrganizationID: orgID,
		PeriodID:       periodID,
		TemplateID:     &templateID,
		Status:         OKRStatusDraft,
	}

	if err := s.repo.CreateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	var details []OKREvaluationDetail
	sortOrder := 0
	for _, obj := range template.Objectives {
		for _, kr := range obj.KeyResults {
			detail := OKREvaluationDetail{
				EvaluationID:    evaluation.ID,
				ObjectiveID:     &obj.ID,
				KeyResultID:     &kr.ID,
				ObjectiveTitle:  obj.Title,
				KeyResultTitle:  kr.Title,
				ObjectiveWeight: obj.Weight,
				KeyResultWeight: kr.Weight,
				TargetValue:     kr.TargetValue,
				TargetType:      kr.TargetType,
				Unit:            kr.Unit,
				FormulaType:     kr.FormulaType,
				SortOrder:       sortOrder,
			}
			details = append(details, detail)
			sortOrder++
		}
	}

	if err := s.repo.CreateOKREvaluationDetailsBatch(db, details); err != nil {
		return nil, err
	}

	return s.GetEvaluationWithDetails(db, evaluation.ID)
}

func (s *okrServiceImpl) GetEvaluationByID(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, id)
	if err != nil {
		return nil, err
	}
	return s.evaluationToResponse(evaluation), nil
}

func (s *okrServiceImpl) GetEvaluationWithDetails(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationWithDetails(db, id)
	if err != nil {
		return nil, err
	}
	return s.evaluationToResponseWithDetails(evaluation), nil
}

func (s *okrServiceImpl) ListEvaluations(db *gorm.DB, employeeID, orgID, periodID *uuid.UUID, status *string, page, perPage int) ([]OKREvaluationResponse, int64, error) {
	evaluations, total, err := s.repo.ListOKREvaluations(db, employeeID, orgID, periodID, status, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]OKREvaluationResponse, len(evaluations))
	for i, e := range evaluations {
		responses[i] = *s.evaluationToResponse(&e)
	}

	return responses, total, nil
}

func (s *okrServiceImpl) UpdateEvaluation(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationRequest) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		evaluation.Status = OKREvaluationStatus(*req.Status)
	}
	if req.ReviewerNotes != nil {
		evaluation.ReviewerNotes = req.ReviewerNotes
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

func (s *okrServiceImpl) DeleteEvaluation(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKREvaluation(db, id)
}

// =========================================================================
// Evaluation Detail & Score
// =========================================================================

func (s *okrServiceImpl) UpdateEvaluationDetailActual(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationDetailRequest) (*OKREvaluationDetailResponse, error) {
	detail, err := s.repo.GetOKREvaluationDetailByID(db, id)
	if err != nil {
		return nil, err
	}

	detail.ActualValue = req.ActualValue
	if req.Remarks != nil {
		detail.Remarks = req.Remarks
	}

	detail.Achievement = s.calculateAchievement(detail.ActualValue, detail.TargetValue, detail.FormulaType)
	detail.Score = (detail.KeyResultWeight * detail.Achievement) / 100

	if err := s.repo.UpdateOKREvaluationDetail(db, detail); err != nil {
		return nil, err
	}

	return s.evaluationDetailToResponse(detail), nil
}

func (s *okrServiceImpl) BulkUpdateEvaluationActuals(db *gorm.DB, evaluationID uuid.UUID, req *OKRBulkUpdateActualsRequest) error {
	for _, item := range req.Details {
		detailID, err := uuid.Parse(item.ID)
		if err != nil {
			continue
		}

		detail, err := s.repo.GetOKREvaluationDetailByID(db, detailID)
		if err != nil {
			continue
		}

		detail.ActualValue = item.ActualValue
		detail.Achievement = s.calculateAchievement(detail.ActualValue, detail.TargetValue, detail.FormulaType)
		detail.Score = (detail.KeyResultWeight * detail.Achievement) / 100

		s.repo.UpdateOKREvaluationDetail(db, detail)
	}

	return nil
}

func (s *okrServiceImpl) RecalculateEvaluationScore(db *gorm.DB, evaluationID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationWithDetails(db, evaluationID)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	for i := range evaluation.Details {
		detail := &evaluation.Details[i]
		detail.Achievement = s.calculateAchievement(detail.ActualValue, detail.TargetValue, detail.FormulaType)
		detail.Score = (detail.KeyResultWeight * detail.ObjectiveWeight * detail.Achievement) / 10000
		totalScore += detail.Score
		s.repo.UpdateOKREvaluationDetail(db, detail)
	}

	evaluation.FinalScore = totalScore
	s.repo.UpdateOKREvaluation(db, evaluation)

	return s.evaluationToResponseWithDetails(evaluation), nil
}

func (s *okrServiceImpl) calculateAchievement(actual, target float64, formulaType FormulaType) float64 {
	if target == 0 {
		return 0
	}

	switch formulaType {
	case FormulaTypeHigherBetter:
		return (actual / target) * 100
	case FormulaTypeLowerBetter:
		if actual == 0 {
			return 100
		}
		return (target / actual) * 100
	default:
		return (actual / target) * 100
	}
}

// =========================================================================
// Workflow
// =========================================================================

func (s *okrServiceImpl) SubmitEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusDraft {
		return nil, errors.New("only draft evaluations can be submitted")
	}

	now := time.Now()
	evaluation.Status = OKRStatusSubmitted
	evaluation.SubmittedAt = &now
	evaluation.SubmittedBy = &userID

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

func (s *okrServiceImpl) ApproveEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusSubmitted {
		return nil, errors.New("only submitted evaluations can be approved")
	}

	now := time.Now()
	evaluation.Status = OKRStatusApproved
	evaluation.ApprovedAt = &now
	evaluation.ApprovedBy = &userID

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

func (s *okrServiceImpl) RejectEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID, notes string) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusSubmitted {
		return nil, errors.New("only submitted evaluations can be rejected")
	}

	evaluation.Status = OKRStatusDraft
	evaluation.ReviewerNotes = &notes

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

func (s *okrServiceImpl) CompleteEvaluation(db *gorm.DB, evaluationID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusApproved {
		return nil, errors.New("only approved evaluations can be completed")
	}

	evaluation.Status = OKRStatusCompleted

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// =========================================================================
// OKR Progress
// =========================================================================

func (s *okrServiceImpl) CreateProgress(db *gorm.DB, req *CreateOKRProgressRequest, userID uuid.UUID) (*OKRProgressResponse, error) {
	detailID, err := uuid.Parse(req.EvaluationDetailID)
	if err != nil {
		return nil, errors.New("invalid evaluation_detail_id")
	}

	detail, err := s.repo.GetOKREvaluationDetailByID(db, detailID)
	if err != nil {
		return nil, err
	}

	achievement := s.calculateAchievement(req.ActualValue, detail.TargetValue, detail.FormulaType)

	progress := &OKRProgress{
		EvaluationDetailID: detailID,
		ProgressDate:       req.ProgressDate,
		ActualValue:        req.ActualValue,
		Achievement:        achievement,
		Notes:              req.Notes,
		CreatedBy:          userID,
	}

	if err := s.repo.CreateOKRProgress(db, progress); err != nil {
		return nil, err
	}

	return s.progressToResponse(progress), nil
}

func (s *okrServiceImpl) GetProgressByID(db *gorm.DB, id uuid.UUID) (*OKRProgressResponse, error) {
	progress, err := s.repo.GetOKRProgressByID(db, id)
	if err != nil {
		return nil, err
	}
	return s.progressToResponse(progress), nil
}

func (s *okrServiceImpl) ListProgressByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRProgressResponse, error) {
	progressList, err := s.repo.ListOKRProgressByDetailID(db, detailID)
	if err != nil {
		return nil, err
	}

	responses := make([]OKRProgressResponse, len(progressList))
	for i, p := range progressList {
		responses[i] = *s.progressToResponse(&p)
	}

	return responses, nil
}

func (s *okrServiceImpl) UpdateProgress(db *gorm.DB, id uuid.UUID, req *UpdateOKRProgressRequest) (*OKRProgressResponse, error) {
	progress, err := s.repo.GetOKRProgressByID(db, id)
	if err != nil {
		return nil, err
	}

	if req.ProgressDate != nil {
		progress.ProgressDate = *req.ProgressDate
	}
	if req.ActualValue != nil {
		progress.ActualValue = *req.ActualValue
		detail, _ := s.repo.GetOKREvaluationDetailByID(db, progress.EvaluationDetailID)
		if detail != nil {
			progress.Achievement = s.calculateAchievement(*req.ActualValue, detail.TargetValue, detail.FormulaType)
		}
	}
	if req.Notes != nil {
		progress.Notes = req.Notes
	}

	if err := s.repo.UpdateOKRProgress(db, progress); err != nil {
		return nil, err
	}

	return s.progressToResponse(progress), nil
}

func (s *okrServiceImpl) DeleteProgress(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRProgress(db, id)
}

// =========================================================================
// OKR Comments
// =========================================================================

func (s *okrServiceImpl) CreateComment(db *gorm.DB, req *CreateOKRCommentRequest, userID uuid.UUID) (*OKRCommentResponse, error) {
	evaluationID, err := uuid.Parse(req.EvaluationID)
	if err != nil {
		return nil, errors.New("invalid evaluation_id")
	}

	comment := &OKRComment{
		EvaluationID: evaluationID,
		Comment:      req.Comment,
		CreatedBy:    userID,
	}

	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent_id")
		}
		comment.ParentID = &parentID
	}

	if err := s.repo.CreateOKRComment(db, comment); err != nil {
		return nil, err
	}

	return s.commentToResponse(comment), nil
}

func (s *okrServiceImpl) ListCommentsByEvaluationID(db *gorm.DB, evaluationID uuid.UUID) ([]OKRCommentResponse, error) {
	comments, err := s.repo.ListOKRCommentsByEvaluationID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	responses := make([]OKRCommentResponse, len(comments))
	for i, c := range comments {
		responses[i] = *s.commentToResponseWithReplies(&c)
	}

	return responses, nil
}

func (s *okrServiceImpl) UpdateComment(db *gorm.DB, id uuid.UUID, req *UpdateOKRCommentRequest) (*OKRCommentResponse, error) {
	comment, err := s.repo.GetOKRCommentByID(db, id)
	if err != nil {
		return nil, err
	}

	comment.Comment = req.Comment

	if err := s.repo.UpdateOKRComment(db, comment); err != nil {
		return nil, err
	}

	return s.commentToResponse(comment), nil
}

func (s *okrServiceImpl) DeleteComment(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRComment(db, id)
}

// =========================================================================
// OKR Attachments
// =========================================================================

func (s *okrServiceImpl) CreateAttachment(db *gorm.DB, req *CreateOKRAttachmentRequest, userID uuid.UUID) (*OKRAttachmentResponse, error) {
	detailID, err := uuid.Parse(req.EvaluationDetailID)
	if err != nil {
		return nil, errors.New("invalid evaluation_detail_id")
	}

	attachment := &OKRAttachment{
		EvaluationDetailID: detailID,
		FilePath:           req.FilePath,
		FileName:           req.FileName,
		FileType:           req.FileType,
		FileSize:           req.FileSize,
		Description:        req.Description,
		UploadedBy:         userID,
	}

	if err := s.repo.CreateOKRAttachment(db, attachment); err != nil {
		return nil, err
	}

	return s.attachmentToResponse(attachment), nil
}

func (s *okrServiceImpl) ListAttachmentsByDetailID(db *gorm.DB, detailID uuid.UUID) ([]OKRAttachmentResponse, error) {
	attachments, err := s.repo.ListOKRAttachmentsByDetailID(db, detailID)
	if err != nil {
		return nil, err
	}

	responses := make([]OKRAttachmentResponse, len(attachments))
	for i, a := range attachments {
		responses[i] = *s.attachmentToResponse(&a)
	}

	return responses, nil
}

func (s *okrServiceImpl) DeleteAttachment(db *gorm.DB, id uuid.UUID) error {
	return s.repo.DeleteOKRAttachment(db, id)
}

// =========================================================================
// Dashboard
// =========================================================================

func (s *okrServiceImpl) GetHRDashboard(db *gorm.DB, periodID *uuid.UUID) (*OKRDashboardHRResponse, error) {
	return s.repo.GetOKRHRDashboardStats(db, periodID)
}

// =========================================================================
// My Context (self-assessment)
// =========================================================================

// GetMyOKRContext resolves the calling user's current Organization (posisi
// jabatan terakhir), then returns the ACTIVE (status=1) OKR templates
// configured for that Organization — mirrors GetMyKPIContext for the KPI
// module. Self-assessment: an employee can only start filling an OKR
// evaluation once a template exists for their current position.
func (s *okrServiceImpl) GetMyOKRContext(db *gorm.DB, userID uuid.UUID) (*MyOKRContextResponse, error) {
	empID, orgID, err := s.repo.GetCurrentEmployeeContext(db, userID)
	if err != nil {
		return nil, err
	}
	if empID == nil || orgID == nil {
		return &MyOKRContextResponse{HasPosition: false, Templates: []OKRTemplateResponse{}}, nil
	}

	activeStatus := 1
	templates, _, err := s.repo.ListOKRTemplates(db, orgID, nil, &activeStatus, 1, maxPerPage)
	if err != nil {
		return nil, err
	}
	responses := make([]OKRTemplateResponse, len(templates))
	for i, t := range templates {
		responses[i] = *s.templateToResponse(&t)
	}

	orgName, err := s.repo.GetOrganizationName(db, *orgID)
	if err != nil {
		return nil, err
	}

	return &MyOKRContextResponse{
		HasPosition:      true,
		EmployeeID:       empID.String(),
		OrganizationID:   orgID.String(),
		OrganizationName: orgName,
		Templates:        responses,
	}, nil
}

// =========================================================================
// Response Converters
// =========================================================================

func (s *okrServiceImpl) templateToResponse(t *OKRTemplate) *OKRTemplateResponse {
	resp := &OKRTemplateResponse{
		ID:             t.ID.String(),
		OrganizationID: t.OrganizationID.String(),
		Name:           t.Name,
		Status:         t.Status,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
	if t.PeriodID != nil {
		resp.PeriodID = t.PeriodID.String()
	}
	if t.Description != nil {
		resp.Description = *t.Description
	}
	if t.EffectiveDate != nil {
		resp.EffectiveDate = *t.EffectiveDate
	}
	if t.ExpiredDate != nil {
		resp.ExpiredDate = *t.ExpiredDate
	}
	return resp
}

func (s *okrServiceImpl) templateToResponseWithObjectives(t *OKRTemplate) *OKRTemplateResponse {
	resp := s.templateToResponse(t)
	resp.ObjectiveCount = len(t.Objectives)
	resp.Objectives = make([]OKRObjectiveResponse, len(t.Objectives))
	for i, o := range t.Objectives {
		resp.Objectives[i] = *s.objectiveToResponseWithKeyResults(&o)
	}
	return resp
}

func (s *okrServiceImpl) objectiveToResponse(o *OKRObjective) *OKRObjectiveResponse {
	resp := &OKRObjectiveResponse{
		ID:         o.ID.String(),
		TemplateID: o.TemplateID.String(),
		Title:      o.Title,
		Weight:     o.Weight,
		SortOrder:  o.SortOrder,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
	if o.Code != nil {
		resp.Code = *o.Code
	}
	if o.Description != nil {
		resp.Description = *o.Description
	}
	return resp
}

func (s *okrServiceImpl) objectiveToResponseWithKeyResults(o *OKRObjective) *OKRObjectiveResponse {
	resp := s.objectiveToResponse(o)
	resp.KeyResults = make([]OKRKeyResultResponse, len(o.KeyResults))
	for i, kr := range o.KeyResults {
		resp.KeyResults[i] = *s.keyResultToResponse(&kr)
	}
	return resp
}

func (s *okrServiceImpl) keyResultToResponse(kr *OKRKeyResult) *OKRKeyResultResponse {
	resp := &OKRKeyResultResponse{
		ID:           kr.ID.String(),
		ObjectiveID:  kr.ObjectiveID.String(),
		Title:        kr.Title,
		TargetType:   string(kr.TargetType),
		TargetValue:  kr.TargetValue,
		FormulaType:  string(kr.FormulaType),
		Weight:       kr.Weight,
		MinimumScore: kr.MinimumScore,
		MaximumScore: kr.MaximumScore,
		SortOrder:    kr.SortOrder,
		IsRequired:   kr.IsRequired,
		CreatedAt:    kr.CreatedAt,
		UpdatedAt:    kr.UpdatedAt,
	}
	if kr.Code != nil {
		resp.Code = *kr.Code
	}
	if kr.Description != nil {
		resp.Description = *kr.Description
	}
	if kr.Unit != nil {
		resp.Unit = *kr.Unit
	}
	return resp
}

func (s *okrServiceImpl) evaluationToResponse(e *OKREvaluation) *OKREvaluationResponse {
	resp := &OKREvaluationResponse{
		ID:             e.ID.String(),
		EmployeeID:     e.EmployeeID.String(),
		OrganizationID: e.OrganizationID.String(),
		PeriodID:       e.PeriodID.String(),
		Status:         string(e.Status),
		FinalScore:     e.FinalScore,
		SubmittedAt:    e.SubmittedAt,
		ApprovedAt:     e.ApprovedAt,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
	if e.TemplateID != nil {
		resp.TemplateID = e.TemplateID.String()
	}
	if e.RatingID != nil {
		resp.RatingID = e.RatingID.String()
	}
	if e.ReviewerNotes != nil {
		resp.ReviewerNotes = *e.ReviewerNotes
	}
	return resp
}

func (s *okrServiceImpl) evaluationToResponseWithDetails(e *OKREvaluation) *OKREvaluationResponse {
	resp := s.evaluationToResponse(e)
	resp.Details = make([]OKREvaluationDetailResponse, len(e.Details))
	for i, d := range e.Details {
		resp.Details[i] = *s.evaluationDetailToResponse(&d)
	}
	return resp
}

func (s *okrServiceImpl) evaluationDetailToResponse(d *OKREvaluationDetail) *OKREvaluationDetailResponse {
	resp := &OKREvaluationDetailResponse{
		ID:              d.ID.String(),
		EvaluationID:    d.EvaluationID.String(),
		ObjectiveTitle:  d.ObjectiveTitle,
		KeyResultTitle:  d.KeyResultTitle,
		ObjectiveWeight: d.ObjectiveWeight,
		KeyResultWeight: d.KeyResultWeight,
		TargetValue:     d.TargetValue,
		TargetType:      string(d.TargetType),
		FormulaType:     string(d.FormulaType),
		ActualValue:     d.ActualValue,
		Achievement:     d.Achievement,
		Score:           d.Score,
		SortOrder:       d.SortOrder,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
	if d.ObjectiveID != nil {
		resp.ObjectiveID = d.ObjectiveID.String()
	}
	if d.KeyResultID != nil {
		resp.KeyResultID = d.KeyResultID.String()
	}
	if d.Unit != nil {
		resp.Unit = *d.Unit
	}
	if d.Remarks != nil {
		resp.Remarks = *d.Remarks
	}
	return resp
}

func (s *okrServiceImpl) progressToResponse(p *OKRProgress) *OKRProgressResponse {
	resp := &OKRProgressResponse{
		ID:                 p.ID.String(),
		EvaluationDetailID: p.EvaluationDetailID.String(),
		ProgressDate:       p.ProgressDate,
		ActualValue:        p.ActualValue,
		Achievement:        p.Achievement,
		CreatedBy:          p.CreatedBy.String(),
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
	if p.Notes != nil {
		resp.Notes = *p.Notes
	}
	return resp
}

func (s *okrServiceImpl) commentToResponse(c *OKRComment) *OKRCommentResponse {
	resp := &OKRCommentResponse{
		ID:           c.ID.String(),
		EvaluationID: c.EvaluationID.String(),
		Comment:      c.Comment,
		CreatedBy:    c.CreatedBy.String(),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
	if c.ParentID != nil {
		resp.ParentID = c.ParentID.String()
	}
	return resp
}

func (s *okrServiceImpl) commentToResponseWithReplies(c *OKRComment) *OKRCommentResponse {
	resp := s.commentToResponse(c)
	resp.Replies = make([]OKRCommentResponse, len(c.Replies))
	for i, r := range c.Replies {
		resp.Replies[i] = *s.commentToResponse(&r)
	}
	return resp
}

func (s *okrServiceImpl) attachmentToResponse(a *OKRAttachment) *OKRAttachmentResponse {
	resp := &OKRAttachmentResponse{
		ID:                 a.ID.String(),
		EvaluationDetailID: a.EvaluationDetailID.String(),
		FilePath:           a.FilePath,
		FileName:           a.FileName,
		UploadedBy:         a.UploadedBy.String(),
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
	if a.FileType != nil {
		resp.FileType = *a.FileType
	}
	if a.FileSize != nil {
		resp.FileSize = *a.FileSize
	}
	if a.Description != nil {
		resp.Description = *a.Description
	}
	return resp
}
