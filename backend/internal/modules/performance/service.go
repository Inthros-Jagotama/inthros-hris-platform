package performance

import (
	"context"
	"fmt"
	"math"

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
	if req.Description != nil {
		t.Description = req.Description
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
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = req.Description
	}
	if req.Status != nil {
		t.Status = *req.Status
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
		PerspectiveID:        perspID,
		IndicatorType:        req.IndicatorType,
		Title:                req.Title,
		Weight:               req.Weight,
		TargetValue:          req.TargetValue,
	}
	if req.Description != nil {
		i.Description = req.Description
	}
	if req.UnitOfMeasurement != nil {
		i.UnitOfMeasurement = req.UnitOfMeasurement
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
		PerspectiveID:          perspID,
		AchievementPercentage:  req.AchievementPercentage,
		Weight:                 req.Weight,
		Score:                  req.Score,
	}
	if req.Description != nil {
		d.Description = req.Description
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
	if req.Score != nil {
		d.Score = *req.Score
	}
	if req.Description != nil {
		d.Description = req.Description
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
	if t.Description != nil {
		r.Description = *t.Description
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
		SortOrder:             i.SortOrder,
		CreatedAt:             i.CreatedAt,
		UpdatedAt:             i.UpdatedAt,
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
	if e.Notes != nil {
		r.Notes = *e.Notes
	}
	return r
}

func detailToResponse(d *PerformanceEvaluationDetail) *EvaluationDetailResponse {
	r := &EvaluationDetailResponse{
		ID:                      d.ID.String(),
		PerformanceEvaluationID: d.PerformanceEvaluationID.String(),
		PerspectiveID:          d.PerspectiveID.String(),
		AchievementPercentage:  d.AchievementPercentage,
		Weight:                 d.Weight,
		Score:                  d.Score,
		CreatedAt:              d.CreatedAt,
		UpdatedAt:              d.UpdatedAt,
	}
	if d.Description != nil {
		r.Description = *d.Description
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
