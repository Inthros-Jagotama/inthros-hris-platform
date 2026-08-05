package performance

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Performance Periods
// =========================================================================

func (s *Service) CreatePerformancePeriod(ctx context.Context, req CreatePerformancePeriodRequest) (*PerformancePeriodResponse, error) {
	p := &PerformancePeriod{
		PeriodCode: req.PeriodCode,
		PeriodType: req.PeriodType,
		Year:       req.Year,
	}
	if req.StartDate != nil {
		p.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		p.EndDate = req.EndDate
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if err := s.repo.CreatePerformancePeriod(ctx, p); err != nil {
		return nil, err
	}
	s.logger.Info("Performance period created", zap.String("id", p.ID.String()), zap.String("code", p.PeriodCode))
	return periodToResponse(p), nil
}

func (s *Service) GetPerformancePeriodByID(ctx context.Context, id string) (*PerformancePeriodResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformancePeriodByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return periodToResponse(p), nil
}

func (s *Service) ListPerformancePeriods(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListPerformancePeriods(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformancePeriodResponse, 0, len(list))
	for _, p := range list {
		responses = append(responses, *periodToResponse(&p))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformancePeriod(ctx context.Context, id string, req UpdatePerformancePeriodRequest) (*PerformancePeriodResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformancePeriodByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.PeriodCode != nil {
		p.PeriodCode = *req.PeriodCode
	}
	if req.PeriodType != nil {
		p.PeriodType = *req.PeriodType
	}
	if req.Year != nil {
		p.Year = *req.Year
	}
	if req.StartDate != nil {
		p.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		p.EndDate = req.EndDate
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if err := s.repo.UpdatePerformancePeriod(ctx, p); err != nil {
		return nil, err
	}
	return periodToResponse(p), nil
}

func (s *Service) DeletePerformancePeriod(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformancePeriod(ctx, uid)
}

// =========================================================================
// Performance Perspectives
// =========================================================================

func (s *Service) CreatePerformancePerspective(ctx context.Context, req CreatePerformancePerspectiveRequest) (*PerformancePerspectiveResponse, error) {
	p := &PerformancePerspective{Name: req.Name}
	if req.Description != nil {
		p.Description = req.Description
	}
	if req.SortOrder != nil {
		p.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreatePerformancePerspective(ctx, p); err != nil {
		return nil, err
	}
	return perspectiveToResponse(p), nil
}

func (s *Service) GetPerformancePerspectiveByID(ctx context.Context, id string) (*PerformancePerspectiveResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformancePerspectiveByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return perspectiveToResponse(p), nil
}

func (s *Service) ListPerformancePerspectives(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListPerformancePerspectives(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformancePerspectiveResponse, 0, len(list))
	for _, p := range list {
		responses = append(responses, *perspectiveToResponse(&p))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformancePerspective(ctx context.Context, id string, req UpdatePerformancePerspectiveRequest) (*PerformancePerspectiveResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformancePerspectiveByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = req.Description
	}
	if req.SortOrder != nil {
		p.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdatePerformancePerspective(ctx, p); err != nil {
		return nil, err
	}
	return perspectiveToResponse(p), nil
}

func (s *Service) DeletePerformancePerspective(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformancePerspective(ctx, uid)
}

// =========================================================================
// Performance Templates
// =========================================================================

func (s *Service) CreatePerformanceTemplate(ctx context.Context, req CreatePerformanceTemplateRequest) (*PerformanceTemplateResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	t := &PerformanceTemplate{
		OrganizationID: orgID,
		Name:           req.Name,
	}
	if req.PeriodID != nil && *req.PeriodID != "" {
		periodUID, err := uuid.Parse(*req.PeriodID)
		if err != nil {
			return nil, fmt.Errorf("invalid period_id: %w", err)
		}
		t.PeriodID = &periodUID
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.EffectiveDate != nil {
		t.EffectiveDate = req.EffectiveDate
	}
	if req.ExpiredDate != nil {
		t.ExpiredDate = req.ExpiredDate
	}
	if err := s.repo.CreatePerformanceTemplate(ctx, t); err != nil {
		return nil, err
	}
	return templateToResponse(t), nil
}

func (s *Service) GetPerformanceTemplateByID(ctx context.Context, id string) (*PerformanceTemplateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindPerformanceTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return templateToResponse(t), nil
}

func (s *Service) ListPerformanceTemplates(ctx context.Context, orgID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUUID *uuid.UUID
	if orgID != nil && *orgID != "" {
		uid, err := uuid.Parse(*orgID)
		if err != nil {
			return nil, fmt.Errorf("invalid organization_id: %w", err)
		}
		orgUUID = &uid
	}
	list, total, err := s.repo.ListPerformanceTemplates(ctx, orgUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceTemplateResponse, 0, len(list))
	for _, t := range list {
		responses = append(responses, *templateToResponse(&t))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformanceTemplate(ctx context.Context, id string, req UpdatePerformanceTemplateRequest) (*PerformanceTemplateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindPerformanceTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.PeriodID != nil {
		if *req.PeriodID != "" {
			periodUID, err := uuid.Parse(*req.PeriodID)
			if err != nil {
				return nil, fmt.Errorf("invalid period_id: %w", err)
			}
			t.PeriodID = &periodUID
		} else {
			t.PeriodID = nil
		}
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if req.EffectiveDate != nil {
		t.EffectiveDate = req.EffectiveDate
	}
	if req.ExpiredDate != nil {
		t.ExpiredDate = req.ExpiredDate
	}
	if err := s.repo.UpdatePerformanceTemplate(ctx, t); err != nil {
		return nil, err
	}
	return templateToResponse(t), nil
}

func (s *Service) DeletePerformanceTemplate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceTemplate(ctx, uid)
}

// =========================================================================
// Performance Indicators
// =========================================================================

func (s *Service) CreatePerformanceIndicator(ctx context.Context, req CreatePerformanceIndicatorRequest) (*PerformanceIndicatorResponse, error) {
	tmplID, err := uuid.Parse(req.PerformanceTemplateID)
	if err != nil {
		return nil, fmt.Errorf("invalid performance_template_id: %w", err)
	}
	perspID, err := uuid.Parse(req.PerspectiveID)
	if err != nil {
		return nil, fmt.Errorf("invalid perspective_id: %w", err)
	}
	i := &PerformanceIndicator{
		PerformanceTemplateID: tmplID,
		PerspectiveID:         perspID,
		IndicatorType:         req.IndicatorType,
		Title:                 req.Title,
		Weight:                req.Weight,
		TargetValue:           req.TargetValue,
		FormulaType:           string(FormulaTypeManual),
		MinimumScore:          0,
		MaximumScore:          100,
		TargetType:            string(TargetTypeNumber),
		IsRequired:            true,
	}
	if req.Code != nil {
		i.Code = req.Code
	}
	if req.Description != nil {
		i.Description = req.Description
	}
	if req.UnitOfMeasurement != nil {
		i.UnitOfMeasurement = req.UnitOfMeasurement
	}
	if req.FormulaType != nil {
		i.FormulaType = *req.FormulaType
	}
	if req.MinimumScore != nil {
		i.MinimumScore = *req.MinimumScore
	}
	if req.MaximumScore != nil {
		i.MaximumScore = *req.MaximumScore
	}
	if req.TargetType != nil {
		i.TargetType = *req.TargetType
	}
	if req.IsRequired != nil {
		i.IsRequired = *req.IsRequired
	}
	if req.SortOrder != nil {
		i.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreatePerformanceIndicator(ctx, i); err != nil {
		return nil, err
	}
	return indicatorToResponse(i), nil
}

func (s *Service) GetPerformanceIndicatorByID(ctx context.Context, id string) (*PerformanceIndicatorResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindPerformanceIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return indicatorToResponse(i), nil
}

func (s *Service) ListPerformanceIndicators(ctx context.Context, templateID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	tmplUID, err := uuid.Parse(templateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template_id: %w", err)
	}
	list, total, err := s.repo.ListPerformanceIndicators(ctx, tmplUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceIndicatorResponse, 0, len(list))
	for _, i := range list {
		responses = append(responses, *indicatorToResponse(&i))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformanceIndicator(ctx context.Context, id string, req UpdatePerformanceIndicatorRequest) (*PerformanceIndicatorResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindPerformanceIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		i.Code = req.Code
	}
	if req.IndicatorType != nil {
		i.IndicatorType = *req.IndicatorType
	}
	if req.Title != nil {
		i.Title = *req.Title
	}
	if req.Description != nil {
		i.Description = req.Description
	}
	if req.Weight != nil {
		i.Weight = *req.Weight
	}
	if req.TargetValue != nil {
		i.TargetValue = *req.TargetValue
	}
	if req.UnitOfMeasurement != nil {
		i.UnitOfMeasurement = req.UnitOfMeasurement
	}
	if req.FormulaType != nil {
		i.FormulaType = *req.FormulaType
	}
	if req.MinimumScore != nil {
		i.MinimumScore = *req.MinimumScore
	}
	if req.MaximumScore != nil {
		i.MaximumScore = *req.MaximumScore
	}
	if req.TargetType != nil {
		i.TargetType = *req.TargetType
	}
	if req.IsRequired != nil {
		i.IsRequired = *req.IsRequired
	}
	if req.SortOrder != nil {
		i.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdatePerformanceIndicator(ctx, i); err != nil {
		return nil, err
	}
	return indicatorToResponse(i), nil
}

func (s *Service) DeletePerformanceIndicator(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceIndicator(ctx, uid)
}

// =========================================================================
// Performance Evaluations
// =========================================================================

func (s *Service) CreatePerformanceEvaluation(ctx context.Context, req CreatePerformanceEvaluationRequest) (*PerformanceEvaluationResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	perID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period_id: %w", err)
	}
	tmplID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template_id: %w", err)
	}
	e := &PerformanceEvaluation{
		EmployeeID:     empID,
		OrganizationID: orgID,
		PeriodID:       perID,
		TemplateID:     tmplID,
		Status:         "DRAFT",
	}
	if req.SupervisorID != nil && *req.SupervisorID != "" {
		uid, _ := uuid.Parse(*req.SupervisorID)
		e.SupervisorID = &uid
	}
	if req.Notes != nil {
		e.Notes = req.Notes
	}
	if err := s.repo.CreatePerformanceEvaluation(ctx, e); err != nil {
		return nil, err
	}
	s.logger.Info("Performance evaluation created", zap.String("id", e.ID.String()), zap.String("employee", e.EmployeeID.String()))
	return evaluationToResponse(e), nil
}

func (s *Service) GetPerformanceEvaluationByID(ctx context.Context, id string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return evaluationToResponse(e), nil
}

func (s *Service) ListPerformanceEvaluations(ctx context.Context, employeeID, periodID, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var empUUID, perUUID *uuid.UUID
	if employeeID != nil && *employeeID != "" {
		uid, _ := uuid.Parse(*employeeID)
		empUUID = &uid
	}
	if periodID != nil && *periodID != "" {
		uid, _ := uuid.Parse(*periodID)
		perUUID = &uid
	}
	list, total, err := s.repo.ListPerformanceEvaluations(ctx, empUUID, perUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceEvaluationResponse, 0, len(list))
	for _, e := range list {
		responses = append(responses, *evaluationToResponse(&e))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformanceEvaluation(ctx context.Context, id string, req UpdatePerformanceEvaluationRequest) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.SupervisorID != nil {
		if *req.SupervisorID != "" {
			uid, _ := uuid.Parse(*req.SupervisorID)
			e.SupervisorID = &uid
		} else {
			e.SupervisorID = nil
		}
	}
	if req.Notes != nil {
		e.Notes = req.Notes
	}
	if err := s.repo.UpdatePerformanceEvaluation(ctx, e); err != nil {
		return nil, err
	}
	return evaluationToResponse(e), nil
}

func (s *Service) UpdateEvaluationStatus(ctx context.Context, id, status, notes string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Only allow forward status transitions
	valid := map[string][]string{
		"DRAFT":          {"PLAN_SUBMITTED"},
		"PLAN_SUBMITTED": {"PLAN_APPROVED", "DRAFT"},
		"PLAN_APPROVED":  {"ACTUAL_SUBMITTED"},
		"ACTUAL_SUBMITTED": {"ACTUAL_APPROVED", "DRAFT"},
		"ACTUAL_APPROVED":  {"COMPLETED"},
	}
	allowed, ok := valid[e.Status]
	if !ok {
		return nil, fmt.Errorf("cannot transition from status %s", e.Status)
	}
	canTransition := false
	for _, s := range allowed {
		if s == status {
			canTransition = true
			break
		}
	}
	if !canTransition {
		return nil, fmt.Errorf("cannot transition from %s to %s", e.Status, status)
	}

	e.Status = status
	if notes != "" {
		e.Notes = &notes
	}

	// If completing, calculate final score
	if status == "COMPLETED" || status == "ACTUAL_APPROVED" {
		score, err := s.repo.UpdateEvaluationFinalScore(ctx, uid)
		if err != nil {
			return nil, err
		}
		e.FinalScore = score
	}

	if err := s.repo.UpdatePerformanceEvaluation(ctx, e); err != nil {
		return nil, err
	}
	s.logger.Info("Evaluation status updated", zap.String("id", e.ID.String()), zap.String("status", e.Status))
	return evaluationToResponse(e), nil
}

func (s *Service) DeletePerformanceEvaluation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceEvaluation(ctx, uid)
}

// =========================================================================
// Evaluation Details
// =========================================================================

func (s *Service) CreateEvaluationDetail(ctx context.Context, req CreateEvaluationDetailRequest) (*EvaluationDetailResponse, error) {
	evalID, err := uuid.Parse(req.PerformanceEvaluationID)
	if err != nil {
		return nil, fmt.Errorf("invalid performance_evaluation_id: %w", err)
	}
	perspID, err := uuid.Parse(req.PerspectiveID)
	if err != nil {
		return nil, fmt.Errorf("invalid perspective_id: %w", err)
	}
	d := &PerformanceEvaluationDetail{
		PerformanceEvaluationID: evalID,
		PerspectiveID:           perspID,
		AchievementPercentage:   req.AchievementPercentage,
		Weight:                  req.Weight,
		Target:                  req.Target,
		Actual:                  req.Actual,
		Achievement:             req.Achievement,
		Score:                   req.Score,
	}
	if req.IndicatorID != nil && *req.IndicatorID != "" {
		indUID, err := uuid.Parse(*req.IndicatorID)
		if err != nil {
			return nil, fmt.Errorf("invalid indicator_id: %w", err)
		}
		d.IndicatorID = &indUID
	}
	if req.IndicatorName != nil {
		d.IndicatorName = req.IndicatorName
	}
	if req.Description != nil {
		d.Description = req.Description
	}
	if req.Remarks != nil {
		d.Remarks = req.Remarks
	}
	if err := s.repo.CreateEvaluationDetail(ctx, d); err != nil {
		return nil, err
	}
	return detailToResponse(d), nil
}

func (s *Service) ListEvaluationDetails(ctx context.Context, evalID string) ([]EvaluationDetailResponse, error) {
	eID, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}
	list, err := s.repo.ListEvaluationDetails(ctx, eID)
	if err != nil {
		return nil, err
	}
	responses := make([]EvaluationDetailResponse, 0, len(list))
	for _, d := range list {
		responses = append(responses, *detailToResponse(&d))
	}
	return responses, nil
}

func (s *Service) UpdateEvaluationDetail(ctx context.Context, id string, req UpdateEvaluationDetailRequest) (*EvaluationDetailResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	d, err := s.repo.FindEvaluationDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.AchievementPercentage != nil {
		d.AchievementPercentage = *req.AchievementPercentage
	}
	if req.Weight != nil {
		d.Weight = *req.Weight
	}
	if req.Target != nil {
		d.Target = *req.Target
	}
	if req.Actual != nil {
		d.Actual = *req.Actual
	}
	if req.Achievement != nil {
		d.Achievement = *req.Achievement
	}
	if req.Score != nil {
		d.Score = *req.Score
	}
	if req.Description != nil {
		d.Description = req.Description
	}
	if req.Remarks != nil {
		d.Remarks = req.Remarks
	}
	if err := s.repo.UpdateEvaluationDetail(ctx, d); err != nil {
		return nil, err
	}
	return detailToResponse(d), nil
}

func (s *Service) DeleteEvaluationDetail(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEvaluationDetail(ctx, uid)
}

// =========================================================================
// Performance Targets
// =========================================================================

func (s *Service) CreatePerformanceTarget(ctx context.Context, req CreatePerformanceTargetRequest) (*PerformanceTargetResponse, error) {
	evalID, err := uuid.Parse(req.PerformanceEvaluationID)
	if err != nil {
		return nil, fmt.Errorf("invalid performance_evaluation_id: %w", err)
	}
	indID, err := uuid.Parse(req.IndicatorID)
	if err != nil {
		return nil, fmt.Errorf("invalid indicator_id: %w", err)
	}
	t := &PerformanceTarget{
		PerformanceEvaluationID: evalID,
		IndicatorID:           indID,
		TargetValue:           req.TargetValue,
		Weight:                req.Weight,
	}
	if req.ActualValue != nil {
		t.ActualValue = req.ActualValue
	}
	if req.UnitOfMeasurement != nil {
		t.UnitOfMeasurement = req.UnitOfMeasurement
	}
	if err := s.repo.CreatePerformanceTarget(ctx, t); err != nil {
		return nil, err
	}
	return targetToResponse(t), nil
}

func (s *Service) ListPerformanceTargets(ctx context.Context, evalID string) ([]PerformanceTargetResponse, error) {
	eID, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}
	list, _, err := s.repo.ListPerformanceTargets(ctx, eID)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceTargetResponse, 0, len(list))
	for _, t := range list {
		responses = append(responses, *targetToResponse(&t))
	}
	return responses, nil
}

func (s *Service) UpdatePerformanceTarget(ctx context.Context, id string, req UpdatePerformanceTargetRequest) (*PerformanceTargetResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindPerformanceTargetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.TargetValue != nil {
		t.TargetValue = *req.TargetValue
	}
	if req.ActualValue != nil {
		t.ActualValue = req.ActualValue
	}
	if req.Weight != nil {
		t.Weight = *req.Weight
	}
	if req.UnitOfMeasurement != nil {
		t.UnitOfMeasurement = req.UnitOfMeasurement
	}
	if err := s.repo.UpdatePerformanceTarget(ctx, t); err != nil {
		return nil, err
	}
	return targetToResponse(t), nil
}

func (s *Service) DeletePerformanceTarget(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceTarget(ctx, uid)
}

// =========================================================================
// Helpers
// =========================================================================

func calcTotalPages(total int64, perPage int) int {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	if pages < 1 {
		return 1
	}
	return pages
}

// =========================================================================
// Response converters
// =========================================================================

func periodToResponse(p *PerformancePeriod) *PerformancePeriodResponse {
	r := &PerformancePeriodResponse{
		ID:         p.ID.String(),
		PeriodCode: p.PeriodCode,
		PeriodType: p.PeriodType,
		Year:       p.Year,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
	if p.StartDate != nil {
		r.StartDate = *p.StartDate
	}
	if p.EndDate != nil {
		r.EndDate = *p.EndDate
	}
	return r
}

func perspectiveToResponse(p *PerformancePerspective) *PerformancePerspectiveResponse {
	r := &PerformancePerspectiveResponse{
		ID:        p.ID.String(),
		Name:      p.Name,
		SortOrder: p.SortOrder,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	if p.Description != nil {
		r.Description = *p.Description
	}
	return r
}

func templateToResponse(t *PerformanceTemplate) *PerformanceTemplateResponse {
	r := &PerformanceTemplateResponse{
		ID:             t.ID.String(),
		OrganizationID: t.OrganizationID.String(),
		Name:           t.Name,
		Status:         t.Status,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
	if t.PeriodID != nil {
		r.PeriodID = t.PeriodID.String()
	}
	if t.Description != nil {
		r.Description = *t.Description
	}
	if t.EffectiveDate != nil {
		r.EffectiveDate = *t.EffectiveDate
	}
	if t.ExpiredDate != nil {
		r.ExpiredDate = *t.ExpiredDate
	}
	return r
}

func indicatorToResponse(i *PerformanceIndicator) *PerformanceIndicatorResponse {
	r := &PerformanceIndicatorResponse{
		ID:                    i.ID.String(),
		PerformanceTemplateID: i.PerformanceTemplateID.String(),
		PerspectiveID:         i.PerspectiveID.String(),
		IndicatorType:         i.IndicatorType,
		Title:                 i.Title,
		Weight:                i.Weight,
		TargetValue:           i.TargetValue,
		FormulaType:           i.FormulaType,
		MinimumScore:          i.MinimumScore,
		MaximumScore:          i.MaximumScore,
		TargetType:            i.TargetType,
		IsRequired:            i.IsRequired,
		SortOrder:             i.SortOrder,
		CreatedAt:             i.CreatedAt,
		UpdatedAt:             i.UpdatedAt,
	}
	if i.Code != nil {
		r.Code = *i.Code
	}
	if i.Description != nil {
		r.Description = *i.Description
	}
	if i.UnitOfMeasurement != nil {
		r.UnitOfMeasurement = *i.UnitOfMeasurement
	}
	return r
}

func evaluationToResponse(e *PerformanceEvaluation) *PerformanceEvaluationResponse {
	r := &PerformanceEvaluationResponse{
		ID:             e.ID.String(),
		EmployeeID:     e.EmployeeID.String(),
		OrganizationID: e.OrganizationID.String(),
		PeriodID:       e.PeriodID.String(),
		TemplateID:     e.TemplateID.String(),
		FinalScore:     e.FinalScore,
		Status:         e.Status,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
	if e.SupervisorID != nil {
		r.SupervisorID = e.SupervisorID.String()
	}
	if e.RatingID != nil {
		r.RatingID = e.RatingID.String()
	}
	if e.SubmittedAt != nil {
		r.SubmittedAt = e.SubmittedAt.Format(time.RFC3339)
	}
	if e.ApprovedAt != nil {
		r.ApprovedAt = e.ApprovedAt.Format(time.RFC3339)
	}
	if e.Notes != nil {
		r.Notes = *e.Notes
	}
	return r
}

func detailToResponse(d *PerformanceEvaluationDetail) *EvaluationDetailResponse {
	r := &EvaluationDetailResponse{
		ID:                      d.ID.String(),
		PerformanceEvaluationID: d.PerformanceEvaluationID.String(),
		PerspectiveID:           d.PerspectiveID.String(),
		AchievementPercentage:   d.AchievementPercentage,
		Weight:                  d.Weight,
		Target:                  d.Target,
		Actual:                  d.Actual,
		Achievement:             d.Achievement,
		Score:                   d.Score,
		CreatedAt:               d.CreatedAt,
		UpdatedAt:               d.UpdatedAt,
	}
	if d.IndicatorID != nil {
		r.IndicatorID = d.IndicatorID.String()
	}
	if d.IndicatorName != nil {
		r.IndicatorName = *d.IndicatorName
	}
	if d.Description != nil {
		r.Description = *d.Description
	}
	if d.Remarks != nil {
		r.Remarks = *d.Remarks
	}
	return r
}

func targetToResponse(t *PerformanceTarget) *PerformanceTargetResponse {
	r := &PerformanceTargetResponse{
		ID:                      t.ID.String(),
		PerformanceEvaluationID: t.PerformanceEvaluationID.String(),
		IndicatorID:            t.IndicatorID.String(),
		TargetValue:            t.TargetValue,
		AchievementPercentage:  t.AchievementPercent,
		Weight:                 t.Weight,
		Score:                  t.Score,
		CreatedAt:              t.CreatedAt,
		UpdatedAt:              t.UpdatedAt,
	}
	if t.UnitOfMeasurement != nil {
		r.UnitOfMeasurement = *t.UnitOfMeasurement
	}
	if t.ActualValue != nil {
		r.ActualValue = *t.ActualValue
	}
	return r
}

// =========================================================================
// Performance Progress
// =========================================================================

func (s *Service) CreatePerformanceProgress(ctx context.Context, req CreatePerformanceProgressRequest) (*PerformanceProgressResponse, error) {
	detailID, err := uuid.Parse(req.EvaluationDetailID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_detail_id: %w", err)
	}
	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by: %w", err)
	}
	p := &PerformanceProgress{
		EvaluationDetailID: detailID,
		ProgressDate:       req.ProgressDate,
		ActualValue:        req.ActualValue,
		Achievement:        req.Achievement,
		CreatedBy:          createdBy,
	}
	if req.Notes != nil {
		p.Notes = req.Notes
	}
	if err := s.repo.CreatePerformanceProgress(ctx, p); err != nil {
		return nil, err
	}
	return progressToResponse(p), nil
}

func (s *Service) GetPerformanceProgressByID(ctx context.Context, id string) (*PerformanceProgressResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformanceProgressByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return progressToResponse(p), nil
}

func (s *Service) ListPerformanceProgressByDetailID(ctx context.Context, detailID string) ([]PerformanceProgressResponse, error) {
	dID, err := uuid.Parse(detailID)
	if err != nil {
		return nil, fmt.Errorf("invalid detail_id: %w", err)
	}
	list, err := s.repo.ListPerformanceProgressByDetailID(ctx, dID)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceProgressResponse, 0, len(list))
	for _, p := range list {
		responses = append(responses, *progressToResponse(&p))
	}
	return responses, nil
}

func (s *Service) UpdatePerformanceProgress(ctx context.Context, id string, req UpdatePerformanceProgressRequest) (*PerformanceProgressResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindPerformanceProgressByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.ProgressDate != nil {
		p.ProgressDate = *req.ProgressDate
	}
	if req.ActualValue != nil {
		p.ActualValue = *req.ActualValue
	}
	if req.Achievement != nil {
		p.Achievement = *req.Achievement
	}
	if req.Notes != nil {
		p.Notes = req.Notes
	}
	if err := s.repo.UpdatePerformanceProgress(ctx, p); err != nil {
		return nil, err
	}
	return progressToResponse(p), nil
}

func (s *Service) DeletePerformanceProgress(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceProgress(ctx, uid)
}

// =========================================================================
// Performance Comments
// =========================================================================

func (s *Service) CreatePerformanceComment(ctx context.Context, req CreatePerformanceCommentRequest) (*PerformanceCommentResponse, error) {
	evalID, err := uuid.Parse(req.EvaluationID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid created_by: %w", err)
	}
	c := &PerformanceComment{
		EvaluationID: evalID,
		EmployeeID:   empID,
		Comment:      req.Comment,
		CreatedBy:    createdBy,
	}
	if err := s.repo.CreatePerformanceComment(ctx, c); err != nil {
		return nil, err
	}
	return commentToResponse(c), nil
}

func (s *Service) GetPerformanceCommentByID(ctx context.Context, id string) (*PerformanceCommentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindPerformanceCommentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return commentToResponse(c), nil
}

func (s *Service) ListPerformanceCommentsByEvaluationID(ctx context.Context, evalID string) ([]PerformanceCommentResponse, error) {
	eID, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}
	list, err := s.repo.ListPerformanceCommentsByEvaluationID(ctx, eID)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceCommentResponse, 0, len(list))
	for _, c := range list {
		responses = append(responses, *commentToResponse(&c))
	}
	return responses, nil
}

func (s *Service) UpdatePerformanceComment(ctx context.Context, id string, req UpdatePerformanceCommentRequest) (*PerformanceCommentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindPerformanceCommentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Comment != nil {
		c.Comment = *req.Comment
	}
	if err := s.repo.UpdatePerformanceComment(ctx, c); err != nil {
		return nil, err
	}
	return commentToResponse(c), nil
}

func (s *Service) DeletePerformanceComment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceComment(ctx, uid)
}

// =========================================================================
// Performance Attachments
// =========================================================================

func (s *Service) CreatePerformanceAttachment(ctx context.Context, req CreatePerformanceAttachmentRequest) (*PerformanceAttachmentResponse, error) {
	detailID, err := uuid.Parse(req.EvaluationDetailID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_detail_id: %w", err)
	}
	uploadedBy, err := uuid.Parse(req.UploadedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid uploaded_by: %w", err)
	}
	a := &PerformanceAttachment{
		EvaluationDetailID: detailID,
		FilePath:           req.FilePath,
		FileName:           req.FileName,
		UploadedBy:         uploadedBy,
	}
	if req.FileType != nil {
		a.FileType = req.FileType
	}
	if req.FileSize != nil {
		a.FileSize = req.FileSize
	}
	if req.Description != nil {
		a.Description = req.Description
	}
	if err := s.repo.CreatePerformanceAttachment(ctx, a); err != nil {
		return nil, err
	}
	return attachmentToResponse(a), nil
}

func (s *Service) GetPerformanceAttachmentByID(ctx context.Context, id string) (*PerformanceAttachmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindPerformanceAttachmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return attachmentToResponse(a), nil
}

func (s *Service) ListPerformanceAttachmentsByDetailID(ctx context.Context, detailID string) ([]PerformanceAttachmentResponse, error) {
	dID, err := uuid.Parse(detailID)
	if err != nil {
		return nil, fmt.Errorf("invalid detail_id: %w", err)
	}
	list, err := s.repo.ListPerformanceAttachmentsByDetailID(ctx, dID)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceAttachmentResponse, 0, len(list))
	for _, a := range list {
		responses = append(responses, *attachmentToResponse(&a))
	}
	return responses, nil
}

func (s *Service) UpdatePerformanceAttachment(ctx context.Context, id string, req UpdatePerformanceAttachmentRequest) (*PerformanceAttachmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindPerformanceAttachmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Description != nil {
		a.Description = req.Description
	}
	if err := s.repo.UpdatePerformanceAttachment(ctx, a); err != nil {
		return nil, err
	}
	return attachmentToResponse(a), nil
}

func (s *Service) DeletePerformanceAttachment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceAttachment(ctx, uid)
}

// =========================================================================
// Performance Ratings
// =========================================================================

func (s *Service) CreatePerformanceRating(ctx context.Context, req CreatePerformanceRatingRequest) (*PerformanceRatingResponse, error) {
	rt := &PerformanceRating{
		Code:     req.Code,
		Name:     req.Name,
		MinScore: req.MinScore,
		MaxScore: req.MaxScore,
	}
	if req.Color != nil {
		rt.Color = req.Color
	}
	if req.Description != nil {
		rt.Description = req.Description
	}
	if req.SortOrder != nil {
		rt.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreatePerformanceRating(ctx, rt); err != nil {
		return nil, err
	}
	return ratingToResponse(rt), nil
}

func (s *Service) GetPerformanceRatingByID(ctx context.Context, id string) (*PerformanceRatingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	rt, err := s.repo.FindPerformanceRatingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return ratingToResponse(rt), nil
}

func (s *Service) ListPerformanceRatings(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListPerformanceRatings(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceRatingResponse, 0, len(list))
	for _, rt := range list {
		responses = append(responses, *ratingToResponse(&rt))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformanceRating(ctx context.Context, id string, req UpdatePerformanceRatingRequest) (*PerformanceRatingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	rt, err := s.repo.FindPerformanceRatingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		rt.Code = *req.Code
	}
	if req.Name != nil {
		rt.Name = *req.Name
	}
	if req.MinScore != nil {
		rt.MinScore = *req.MinScore
	}
	if req.MaxScore != nil {
		rt.MaxScore = *req.MaxScore
	}
	if req.Color != nil {
		rt.Color = req.Color
	}
	if req.Description != nil {
		rt.Description = req.Description
	}
	if req.SortOrder != nil {
		rt.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdatePerformanceRating(ctx, rt); err != nil {
		return nil, err
	}
	return ratingToResponse(rt), nil
}

func (s *Service) DeletePerformanceRating(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceRating(ctx, uid)
}

// =========================================================================
// Performance Indicator Formulas
// =========================================================================

func (s *Service) CreatePerformanceIndicatorFormula(ctx context.Context, req CreatePerformanceIndicatorFormulaRequest) (*PerformanceIndicatorFormulaResponse, error) {
	f := &PerformanceIndicatorFormula{
		Code:        req.Code,
		Name:        req.Name,
		FormulaType: req.FormulaType,
	}
	if req.Expression != nil {
		f.Expression = req.Expression
	}
	if req.Description != nil {
		f.Description = req.Description
	}
	if req.SortOrder != nil {
		f.SortOrder = *req.SortOrder
	}
	if err := s.repo.CreatePerformanceIndicatorFormula(ctx, f); err != nil {
		return nil, err
	}
	return formulaToResponse(f), nil
}

func (s *Service) GetPerformanceIndicatorFormulaByID(ctx context.Context, id string) (*PerformanceIndicatorFormulaResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindPerformanceIndicatorFormulaByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return formulaToResponse(f), nil
}

func (s *Service) ListPerformanceIndicatorFormulas(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListPerformanceIndicatorFormulas(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceIndicatorFormulaResponse, 0, len(list))
	for _, f := range list {
		responses = append(responses, *formulaToResponse(&f))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdatePerformanceIndicatorFormula(ctx context.Context, id string, req UpdatePerformanceIndicatorFormulaRequest) (*PerformanceIndicatorFormulaResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindPerformanceIndicatorFormulaByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Code != nil {
		f.Code = *req.Code
	}
	if req.Name != nil {
		f.Name = *req.Name
	}
	if req.FormulaType != nil {
		f.FormulaType = *req.FormulaType
	}
	if req.Expression != nil {
		f.Expression = req.Expression
	}
	if req.Description != nil {
		f.Description = req.Description
	}
	if req.SortOrder != nil {
		f.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdatePerformanceIndicatorFormula(ctx, f); err != nil {
		return nil, err
	}
	return formulaToResponse(f), nil
}

func (s *Service) DeletePerformanceIndicatorFormula(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeletePerformanceIndicatorFormula(ctx, uid)
}

// =========================================================================
// Performance Logs
// =========================================================================

func (s *Service) GetPerformanceLogByID(ctx context.Context, id string) (*PerformanceLogResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	l, err := s.repo.FindPerformanceLogByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return logToResponse(l), nil
}

func (s *Service) ListPerformanceLogs(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListPerformanceLogs(ctx, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceLogResponse, 0, len(list))
	for _, l := range list {
		responses = append(responses, *logToResponse(&l))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) ListPerformanceLogsByEvaluationID(ctx context.Context, evalID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	eID, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}
	list, total, err := s.repo.ListPerformanceLogsByEvaluationID(ctx, eID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]PerformanceLogResponse, 0, len(list))
	for _, l := range list {
		responses = append(responses, *logToResponse(&l))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Phase 2 Response Converters
// =========================================================================

func progressToResponse(p *PerformanceProgress) *PerformanceProgressResponse {
	r := &PerformanceProgressResponse{
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
		r.Notes = *p.Notes
	}
	return r
}

func commentToResponse(c *PerformanceComment) *PerformanceCommentResponse {
	return &PerformanceCommentResponse{
		ID:           c.ID.String(),
		EvaluationID: c.EvaluationID.String(),
		EmployeeID:   c.EmployeeID.String(),
		Comment:      c.Comment,
		CreatedBy:    c.CreatedBy.String(),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func attachmentToResponse(a *PerformanceAttachment) *PerformanceAttachmentResponse {
	r := &PerformanceAttachmentResponse{
		ID:                 a.ID.String(),
		EvaluationDetailID: a.EvaluationDetailID.String(),
		FilePath:           a.FilePath,
		FileName:           a.FileName,
		UploadedBy:         a.UploadedBy.String(),
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
	if a.FileType != nil {
		r.FileType = *a.FileType
	}
	if a.FileSize != nil {
		r.FileSize = *a.FileSize
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

func ratingToResponse(rt *PerformanceRating) *PerformanceRatingResponse {
	r := &PerformanceRatingResponse{
		ID:        rt.ID.String(),
		Code:      rt.Code,
		Name:      rt.Name,
		MinScore:  rt.MinScore,
		MaxScore:  rt.MaxScore,
		SortOrder: rt.SortOrder,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
	}
	if rt.Color != nil {
		r.Color = *rt.Color
	}
	if rt.Description != nil {
		r.Description = *rt.Description
	}
	return r
}

func formulaToResponse(f *PerformanceIndicatorFormula) *PerformanceIndicatorFormulaResponse {
	r := &PerformanceIndicatorFormulaResponse{
		ID:          f.ID.String(),
		Code:        f.Code,
		Name:        f.Name,
		FormulaType: f.FormulaType,
		SortOrder:   f.SortOrder,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
	if f.Expression != nil {
		r.Expression = *f.Expression
	}
	if f.Description != nil {
		r.Description = *f.Description
	}
	return r
}

func logToResponse(l *PerformanceLog) *PerformanceLogResponse {
	r := &PerformanceLogResponse{
		ID:         l.ID.String(),
		EntityType: l.EntityType,
		EntityID:   l.EntityID.String(),
		Action:     l.Action,
		CreatedBy:  l.CreatedBy.String(),
		CreatedAt:  l.CreatedAt,
	}
	if l.EvaluationID != nil {
		r.EvaluationID = l.EvaluationID.String()
	}
	if l.OldValues != nil {
		r.OldValues = *l.OldValues
	}
	if l.NewValues != nil {
		r.NewValues = *l.NewValues
	}
	return r
}

// =========================================================================
// Phase 3 - Business Process Methods
// =========================================================================

// CreateEvaluationWithSnapshot creates an evaluation and snapshots KPIs from template
func (s *Service) CreateEvaluationWithSnapshot(ctx context.Context, req CreateEvaluationWithSnapshotRequest) (*EvaluationWithDetailsResponse, error) {
	// Parse UUIDs
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	periodID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period_id: %w", err)
	}
	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("invalid template_id: %w", err)
	}

	// Get indicators from template
	indicators, err := s.repo.ListIndicatorsByTemplateID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get indicators: %w", err)
	}
	if len(indicators) == 0 {
		return nil, fmt.Errorf("template has no indicators")
	}

	// Create evaluation
	eval := &PerformanceEvaluation{
		EmployeeID:            empID,
		OrganizationID:        orgID,
		PerformancePeriodID:   periodID,
		PerformanceTemplateID: templateID,
		Status:                "DRAFT",
		FinalScore:            0,
	}
	if req.SupervisorID != nil {
		supID, err := uuid.Parse(*req.SupervisorID)
		if err == nil {
			eval.SupervisorID = &supID
		}
	}
	if req.Notes != nil {
		eval.Notes = req.Notes
	}

	// Create evaluation details (snapshot from indicators)
	var details []PerformanceEvaluationDetail
	for _, ind := range indicators {
		detail := PerformanceEvaluationDetail{
			IndicatorID:   &ind.ID,
			IndicatorName: &ind.Name,
			PerspectiveID: ind.PerspectiveID,
			Weight:        ind.Weight,
			Target:        ind.TargetValue,
			Actual:        0,
			Achievement:   0,
			Score:         0,
		}
		if ind.Description != nil {
			detail.Description = ind.Description
		}
		details = append(details, detail)
	}

	// Save evaluation and details
	if err := s.repo.CreateEvaluationWithDetails(ctx, eval, details); err != nil {
		return nil, fmt.Errorf("failed to create evaluation: %w", err)
	}

	// Fetch the created evaluation with details
	return s.GetEvaluationWithDetails(ctx, eval.ID.String())
}

// GetEvaluationWithDetails returns evaluation with all its details
func (s *Service) GetEvaluationWithDetails(ctx context.Context, evalID string) (*EvaluationWithDetailsResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	eval, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	details, err := s.repo.ListEvaluationDetailsByEvaluationID(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := &EvaluationWithDetailsResponse{
		ID:             eval.ID.String(),
		EmployeeID:     eval.EmployeeID.String(),
		OrganizationID: eval.OrganizationID.String(),
		PeriodID:       eval.PerformancePeriodID.String(),
		TemplateID:     eval.PerformanceTemplateID.String(),
		FinalScore:     eval.FinalScore,
		Status:         eval.Status,
		CreatedAt:      eval.CreatedAt,
		UpdatedAt:      eval.UpdatedAt,
	}
	if eval.SupervisorID != nil {
		resp.SupervisorID = eval.SupervisorID.String()
	}
	if eval.RatingID != nil {
		resp.RatingID = eval.RatingID.String()
		// Get rating details
		rating, err := s.repo.FindPerformanceRatingByID(ctx, *eval.RatingID)
		if err == nil && rating != nil {
			resp.RatingName = rating.Name
			if rating.Color != nil {
				resp.RatingColor = *rating.Color
			}
		}
	}
	if eval.SubmittedAt != nil {
		resp.SubmittedAt = eval.SubmittedAt.Format("2006-01-02 15:04:05")
	}
	if eval.ApprovedAt != nil {
		resp.ApprovedAt = eval.ApprovedAt.Format("2006-01-02 15:04:05")
	}
	if eval.Notes != nil {
		resp.Notes = *eval.Notes
	}

	// Convert details
	for _, d := range details {
		resp.Details = append(resp.Details, *detailToResponse(&d))
	}

	return resp, nil
}

// UpdateEvaluationActual updates actual value and calculates achievement/score
func (s *Service) UpdateEvaluationActual(ctx context.Context, detailID string, req UpdateEvaluationActualRequest) (*EvaluationDetailResponse, error) {
	uid, err := uuid.Parse(detailID)
	if err != nil {
		return nil, fmt.Errorf("invalid detail_id: %w", err)
	}

	detail, err := s.repo.FindEvaluationDetailByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Update actual value
	detail.Actual = req.Actual
	if req.Remarks != nil {
		detail.Remarks = req.Remarks
	}

	// Calculate achievement percentage
	if detail.Target > 0 {
		detail.Achievement = (detail.Actual / detail.Target) * 100
		if detail.Achievement > 100 {
			detail.Achievement = 100 // Cap at 100%
		}
	}

	// Calculate score (weight * achievement / 100)
	detail.Score = (detail.Weight * detail.Achievement) / 100

	if err := s.repo.UpdateEvaluationDetail(ctx, detail); err != nil {
		return nil, err
	}

	// Recalculate evaluation final score
	if err := s.RecalculateEvaluationScore(ctx, detail.PerformanceEvaluationID.String()); err != nil {
		return nil, err
	}

	return detailToResponse(detail), nil
}

// BulkUpdateEvaluationActuals updates multiple evaluation details
func (s *Service) BulkUpdateEvaluationActuals(ctx context.Context, req BulkUpdateEvaluationActualRequest) ([]EvaluationDetailResponse, error) {
	var evalID uuid.UUID
	var details []PerformanceEvaluationDetail

	for _, item := range req.Details {
		uid, err := uuid.Parse(item.DetailID)
		if err != nil {
			return nil, fmt.Errorf("invalid detail_id %s: %w", item.DetailID, err)
		}

		detail, err := s.repo.FindEvaluationDetailByID(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("detail %s not found: %w", item.DetailID, err)
		}

		// Track evaluation ID for score recalculation
		evalID = detail.PerformanceEvaluationID

		// Update actual value
		detail.Actual = item.Actual
		if item.Remarks != nil {
			detail.Remarks = item.Remarks
		}

		// Calculate achievement percentage
		if detail.Target > 0 {
			detail.Achievement = (detail.Actual / detail.Target) * 100
			if detail.Achievement > 100 {
				detail.Achievement = 100
			}
		}

		// Calculate score
		detail.Score = (detail.Weight * detail.Achievement) / 100

		details = append(details, *detail)
	}

	if err := s.repo.BulkUpdateEvaluationDetails(ctx, details); err != nil {
		return nil, err
	}

	// Recalculate evaluation final score
	if evalID != uuid.Nil {
		if err := s.RecalculateEvaluationScore(ctx, evalID.String()); err != nil {
			return nil, err
		}
	}

	// Return updated details
	var responses []EvaluationDetailResponse
	for _, d := range details {
		responses = append(responses, *detailToResponse(&d))
	}
	return responses, nil
}

// RecalculateEvaluationScore recalculates final score and assigns rating
func (s *Service) RecalculateEvaluationScore(ctx context.Context, evalID string) error {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return fmt.Errorf("invalid evaluation_id: %w", err)
	}

	// Get sum of all scores
	finalScore, err := s.repo.UpdateEvaluationFinalScore(ctx, uid)
	if err != nil {
		return err
	}

	// Find matching rating
	var ratingID *uuid.UUID
	rating, err := s.repo.FindRatingByScore(ctx, finalScore)
	if err != nil {
		return err
	}
	if rating != nil {
		ratingID = &rating.ID
	}

	// Update evaluation
	return s.repo.UpdateEvaluationWithRating(ctx, uid, finalScore, ratingID)
}

// GetEvaluationProgressSummary returns progress summary for an evaluation
func (s *Service) GetEvaluationProgressSummary(ctx context.Context, evalID string) (*ProgressSummaryResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	total, completed, inProgress, notStarted, avgAchievement, err := s.repo.GetEvaluationProgressSummary(ctx, uid)
	if err != nil {
		return nil, err
	}

	var overallProgress float64
	if total > 0 {
		overallProgress = float64(completed) / float64(total) * 100
	}

	return &ProgressSummaryResponse{
		EvaluationID:       evalID,
		TotalIndicators:    total,
		CompletedCount:     completed,
		InProgressCount:    inProgress,
		NotStartedCount:    notStarted,
		OverallProgress:    overallProgress,
		AverageAchievement: avgAchievement,
	}, nil
}

// SubmitEvaluation changes status from DRAFT to SUBMITTED
func (s *Service) SubmitEvaluation(ctx context.Context, evalID string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	eval, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if eval.Status != "DRAFT" {
		return nil, fmt.Errorf("evaluation can only be submitted from DRAFT status, current: %s", eval.Status)
	}

	// Check if all details have actual values
	details, err := s.repo.ListEvaluationDetailsByEvaluationID(ctx, uid)
	if err != nil {
		return nil, err
	}
	for _, d := range details {
		if d.Actual == 0 {
			return nil, fmt.Errorf("all KPI actual values must be filled before submission")
		}
	}

	// Update status
	now := time.Now()
	eval.Status = "SUBMITTED"
	eval.SubmittedAt = &now

	if err := s.repo.UpdatePerformanceEvaluation(ctx, eval); err != nil {
		return nil, err
	}

	return evalToResponse(eval), nil
}

// ApproveEvaluation changes status from SUBMITTED to APPROVED
func (s *Service) ApproveEvaluation(ctx context.Context, evalID string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	eval, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if eval.Status != "SUBMITTED" {
		return nil, fmt.Errorf("evaluation can only be approved from SUBMITTED status, current: %s", eval.Status)
	}

	// Update status
	now := time.Now()
	eval.Status = "APPROVED"
	eval.ApprovedAt = &now

	if err := s.repo.UpdatePerformanceEvaluation(ctx, eval); err != nil {
		return nil, err
	}

	return evalToResponse(eval), nil
}

// RejectEvaluation changes status from SUBMITTED back to DRAFT for revision
func (s *Service) RejectEvaluation(ctx context.Context, evalID string, notes *string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	eval, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if eval.Status != "SUBMITTED" {
		return nil, fmt.Errorf("evaluation can only be rejected from SUBMITTED status, current: %s", eval.Status)
	}

	// Update status
	eval.Status = "DRAFT"
	eval.SubmittedAt = nil
	if notes != nil {
		eval.Notes = notes
	}

	if err := s.repo.UpdatePerformanceEvaluation(ctx, eval); err != nil {
		return nil, err
	}

	return evalToResponse(eval), nil
}

// CompleteEvaluation changes status from APPROVED to COMPLETED
func (s *Service) CompleteEvaluation(ctx context.Context, evalID string) (*PerformanceEvaluationResponse, error) {
	uid, err := uuid.Parse(evalID)
	if err != nil {
		return nil, fmt.Errorf("invalid evaluation_id: %w", err)
	}

	eval, err := s.repo.FindPerformanceEvaluationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if eval.Status != "APPROVED" {
		return nil, fmt.Errorf("evaluation can only be completed from APPROVED status, current: %s", eval.Status)
	}

	// Update status
	eval.Status = "COMPLETED"

	if err := s.repo.UpdatePerformanceEvaluation(ctx, eval); err != nil {
		return nil, err
	}

	return evalToResponse(eval), nil
}
