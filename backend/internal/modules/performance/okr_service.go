package performance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OKRService interface {
	// OKR Templates
	CreateTemplate(ctx context.Context, db *gorm.DB, userID uuid.UUID, req *CreateOKRTemplateRequest) (*OKRTemplateResponse, error)
	GetTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error)
	GetTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error)
	ListTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplateResponse, int64, error)
	UpdateTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID, req *UpdateOKRTemplateRequest) (*OKRTemplateResponse, error)
	DeleteTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID) error
	DuplicateTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID) (*OKRTemplateResponse, error)

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
	CreateEvaluationKeyResult(db *gorm.DB, req *CreateOKREvaluationKeyResultRequest) (*OKREvaluationDetailResponse, error)
	UpdateEvaluationKeyResultTarget(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationKeyResultTargetRequest) (*OKREvaluationDetailResponse, error)
	DeleteEvaluationKeyResult(db *gorm.DB, id uuid.UUID) error
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
	SubmitKeyResults(ctx context.Context, db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error)
	ApproveKeyResults(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error)
	RejectKeyResults(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID, notes string) (*OKREvaluationResponse, error)
	SubmitEvaluation(ctx context.Context, db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error)
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

	// Cascading objective creation (top-down through the org hierarchy)
	GetObjectiveCreationScope(db *gorm.DB, userID uuid.UUID) (*OKRObjectiveScopeResponse, error)

	// Central approval module integration
	SetApprovalEngine(engine ApprovalEngine)
	SetNotifier(n Notifier)
	HandleKeyResultApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error
	HandleAssessmentApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error
}

// ApprovalModuleOKRKeyResult/OKRAssessment are separate approval module
// slugs (not just "okr") because Key Result approval and final assessment
// approval are two independent approval instances on the same evaluation,
// potentially configured with different flows/approvers — mirrors
// ApprovalModuleKPITarget/KPIRealization.
const (
	ApprovalModuleOKRKeyResult  = "okr_key_result"
	ApprovalModuleOKRAssessment = "okr_assessment"
)

// OKRTemplateStatusPublished marks an OKR template as published (status=1),
// after which it can no longer be edited or deleted. 0 = draft.
const OKRTemplateStatusPublished = 1

type okrServiceImpl struct {
	repo           OKRRepository
	dbResolver     TenantDBFunc
	approvalEngine ApprovalEngine
	notifier       Notifier
}

func NewOKRService(repo OKRRepository, dbResolver TenantDBFunc) OKRService {
	return &okrServiceImpl{repo: repo, dbResolver: dbResolver}
}

// SetApprovalEngine wires the central approval module into the OKR service —
// mirrors Service.SetApprovalEngine for KPI.
func (s *okrServiceImpl) SetApprovalEngine(engine ApprovalEngine) {
	s.approvalEngine = engine
}

// SetNotifier wires the notification module into the OKR service so Key
// Result and assessment approval outcomes can be relayed to the requesting
// employee (docs/module-notification-plan.md §8, Phase 5) — mirrors
// Service.SetNotifier for KPI.
func (s *okrServiceImpl) SetNotifier(n Notifier) {
	s.notifier = n
}

// notifyEvaluationOutcome notifies employeeID of an OKR evaluation's
// approval outcome. Best-effort: if the notifier isn't wired, the employee
// has no linked user account, or Notify fails, this moves on rather than
// failing the approval itself — the OKR service has no logger, so failures
// are intentionally silent (the approval outcome itself is already
// persisted before this runs, mirroring leave.Service.notifyLeaveOutcome's
// best-effort discipline). The ctx must carry company_id (it comes from the
// approval module's push callback) — notification.Service resolves the
// tenant DB from it.
func (s *okrServiceImpl) notifyEvaluationOutcome(ctx context.Context, db *gorm.DB, employeeID uuid.UUID, notifType, referenceType string, referenceID uuid.UUID) {
	if s.notifier == nil {
		return
	}
	userID, err := s.repo.FindUserIDByEmployeeID(db, employeeID)
	if err != nil || userID == nil {
		return
	}
	_ = s.notifier.Notify(ctx, *userID, notifType, nil, referenceType, referenceID)
}

// =========================================================================
// OKR Templates
// =========================================================================

// resolveObjectiveEligibility determines whether ownOrgID is currently
// allowed to create an Objective for a subordinate: true if it's the top of
// the hierarchy (no parent — seeds the cascade), otherwise true only if it
// has already received its own Objective (an active Template with
// Objectives targeting it).
func (s *okrServiceImpl) resolveObjectiveEligibility(db *gorm.DB, ownOrgID uuid.UUID) (bool, string, error) {
	parent, err := s.repo.GetOrganizationParentID(db, ownOrgID)
	if err != nil {
		return false, "", err
	}
	if parent == nil {
		return true, "", nil
	}
	hasOwn, err := s.repo.HasActiveTemplateWithObjectives(db, ownOrgID)
	if err != nil {
		return false, "", err
	}
	if !hasOwn {
		return false, "okr.objective_scope_ineligible_no_own_objective", nil
	}
	return true, "", nil
}

// GetObjectiveCreationScope resolves whether the calling employee can
// currently create an Objective (a Template) for a subordinate Organization,
// and which Organizations qualify as their effective subordinates (direct
// children, walking down through vacant Organizations — see
// GetEffectiveChildOrganizationIDs).
func (s *okrServiceImpl) GetObjectiveCreationScope(db *gorm.DB, userID uuid.UUID) (*OKRObjectiveScopeResponse, error) {
	_, ownOrgID, err := s.repo.GetCurrentEmployeeContext(db, userID)
	if err != nil {
		return nil, err
	}
	if ownOrgID == nil {
		return &OKRObjectiveScopeResponse{
			Eligible:                 false,
			IneligibleReasonKey:      "okr.objective_scope_ineligible_no_position",
			SubordinateOrganizations: []OrganizationOptionResponse{},
		}, nil
	}

	eligible, reasonKey, err := s.resolveObjectiveEligibility(db, *ownOrgID)
	if err != nil {
		return nil, err
	}

	subordinateIDs, err := s.repo.GetEffectiveChildOrganizationIDs(db, *ownOrgID)
	if err != nil {
		return nil, err
	}
	names, err := s.repo.GetOrganizationNamesByIDs(db, subordinateIDs)
	if err != nil {
		return nil, err
	}
	ownName, err := s.repo.GetOrganizationName(db, *ownOrgID)
	if err != nil {
		return nil, err
	}

	subordinates := make([]OrganizationOptionResponse, 0, len(subordinateIDs))
	for _, id := range subordinateIDs {
		subordinates = append(subordinates, OrganizationOptionResponse{ID: id.String(), Name: names[id.String()]})
	}

	return &OKRObjectiveScopeResponse{
		OrganizationID:           ownOrgID.String(),
		OrganizationName:         ownName,
		Eligible:                 eligible,
		IneligibleReasonKey:      reasonKey,
		SubordinateOrganizations: subordinates,
	}, nil
}

func (s *okrServiceImpl) CreateTemplate(ctx context.Context, db *gorm.DB, userID uuid.UUID, req *CreateOKRTemplateRequest) (*OKRTemplateResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, errors.New("invalid organization_id")
	}

	_, ownOrgID, err := s.repo.GetCurrentEmployeeContext(db, userID)
	if err != nil {
		return nil, err
	}
	if ownOrgID == nil {
		return nil, errors.New("you don't have an active position, so no subordinate organization can be resolved")
	}

	eligible, _, err := s.resolveObjectiveEligibility(db, *ownOrgID)
	if err != nil {
		return nil, err
	}
	if !eligible {
		return nil, errors.New("you must have your own objective before you can create one for a subordinate")
	}

	subordinateIDs, err := s.repo.GetEffectiveChildOrganizationIDs(db, *ownOrgID)
	if err != nil {
		return nil, err
	}
	isSubordinate := false
	for _, id := range subordinateIDs {
		if id == orgID {
			isSubordinate = true
			break
		}
	}
	if !isSubordinate {
		return nil, errors.New("this organization is not one of your effective subordinates")
	}

	template := &OKRTemplate{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		EffectiveDate:  req.EffectiveDate,
		ExpiredDate:    req.ExpiredDate,
		CreatedBy:      userID,
		// Organisasi pembuat (bukan org template — template dibuat utk org bawahan).
		// Otorisasi edit/hapus membandingkan org user saat login dengan kolom ini.
		CreatedByOrgID: ownOrgID,
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

	s.notifyTemplateCreated(ctx, db, template)
	return s.templateToResponse(template), nil
}

// notifyTemplateCreated notifies every employee currently occupying the
// template's Organization that a new OKR template is available for them.
// Best-effort: a notification failure never fails template creation, and
// recipients without a linked user account are skipped. The ctx must carry
// company_id (request context) — notification.Service resolves the tenant DB
// from it; a bare context.Background() would silently drop every notification.
func (s *okrServiceImpl) notifyTemplateCreated(ctx context.Context, db *gorm.DB, t *OKRTemplate) {
	if s.notifier == nil {
		return
	}
	userIDs, err := s.repo.GetUserIDsByOrganization(db, t.OrganizationID)
	if err != nil {
		return
	}
	for _, uid := range userIDs {
		_ = s.notifier.Notify(ctx, uid, "OKR_TEMPLATE_CREATED", []string{t.Name}, "okr_template", t.ID)
	}
}

func (s *okrServiceImpl) GetTemplateByID(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateByID(db, id)
	if err != nil {
		return nil, err
	}
	resp := s.templateToResponse(template)
	if err := s.enrichTemplateResponses(db, []*OKRTemplateResponse{resp}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *okrServiceImpl) GetTemplateWithObjectives(db *gorm.DB, id uuid.UUID) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateWithObjectives(db, id)
	if err != nil {
		return nil, err
	}
	resp := s.templateToResponseWithObjectives(template)
	if err := s.enrichTemplateResponses(db, []*OKRTemplateResponse{resp}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *okrServiceImpl) ListTemplates(db *gorm.DB, orgID *uuid.UUID, periodID *uuid.UUID, status *int, page, perPage int) ([]OKRTemplateResponse, int64, error) {
	templates, total, err := s.repo.ListOKRTemplates(db, orgID, periodID, status, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]OKRTemplateResponse, len(templates))
	ptrs := make([]*OKRTemplateResponse, len(templates))
	for i, t := range templates {
		responses[i] = *s.templateToResponse(&t)
		ptrs[i] = &responses[i]
	}
	if err := s.enrichTemplateResponses(db, ptrs); err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}

// enrichTemplateResponses mengisi OrganizationName & PeriodCode yang tidak
// tersedia di model OKRTemplate (butuh lookup lintas tabel organizations dan
// performance_periods) — pola sama dengan KPI enrichTemplateResponses.
func (s *okrServiceImpl) enrichTemplateResponses(db *gorm.DB, responses []*OKRTemplateResponse) error {
	if len(responses) == 0 {
		return nil
	}

	orgIDSet := make(map[uuid.UUID]struct{})
	periodIDSet := make(map[uuid.UUID]struct{})
	for _, r := range responses {
		if r.OrganizationID != "" {
			if id, err := uuid.Parse(r.OrganizationID); err == nil {
				orgIDSet[id] = struct{}{}
			}
		}
		if r.PeriodID != "" {
			if id, err := uuid.Parse(r.PeriodID); err == nil {
				periodIDSet[id] = struct{}{}
			}
		}
	}

	orgIDs := make([]uuid.UUID, 0, len(orgIDSet))
	for id := range orgIDSet {
		orgIDs = append(orgIDs, id)
	}
	periodIDs := make([]uuid.UUID, 0, len(periodIDSet))
	for id := range periodIDSet {
		periodIDs = append(periodIDs, id)
	}

	orgNames, err := s.repo.GetOrganizationNamesByIDs(db, orgIDs)
	if err != nil {
		return err
	}
	periodCodes, err := s.repo.GetPeriodCodesByIDs(db, periodIDs)
	if err != nil {
		return err
	}

	for _, r := range responses {
		r.OrganizationName = orgNames[r.OrganizationID]
		r.PeriodCode = periodCodes[r.PeriodID]
	}
	return nil
}

func (s *okrServiceImpl) UpdateTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID, req *UpdateOKRTemplateRequest) (*OKRTemplateResponse, error) {
	template, err := s.repo.GetOKRTemplateByID(db, id)
	if err != nil {
		return nil, err
	}

	if template.Status == OKRTemplateStatusPublished {
		return nil, errors.New("published OKR templates cannot be modified")
	}
	if err := s.authorizeTemplateOrg(db, template.CreatedByOrgID, userID); err != nil {
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

func (s *okrServiceImpl) DeleteTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID) error {
	template, err := s.repo.GetOKRTemplateByID(db, id)
	if err != nil {
		return err
	}
	if template.Status == OKRTemplateStatusPublished {
		return errors.New("published OKR templates cannot be deleted")
	}
	if err := s.authorizeTemplateOrg(db, template.CreatedByOrgID, userID); err != nil {
		return err
	}
	return s.repo.DeleteOKRTemplate(db, id)
}

// authorizeTemplateOrg memastikan user yang sedang login masih berada di
// organisasi PEMBUAT template (created_by_org_id) sebelum boleh
// mengubah/menghapusnya. Perbandingan STRICT ke created_by_org_id — tanpa
// fallback ke organization_id template. Template tanpa created_by_org_id
// tercatat (legacy) tidak bisa dikelola sama sekali. Aturan berbasis
// ORGANISASI: karyawan yang sudah pindah organisasi tidak lagi bisa
// mengubah/menghapus template meskipun dia yang membuatnya.
func (s *okrServiceImpl) authorizeTemplateOrg(db *gorm.DB, creatorOrgID *uuid.UUID, userID uuid.UUID) error {
	if creatorOrgID == nil {
		return errors.New("this OKR template has no recorded creator organization and cannot be managed")
	}
	_, orgID, err := s.repo.GetCurrentEmployeeContext(db, userID)
	if err != nil {
		return err
	}
	if orgID == nil || *orgID != *creatorOrgID {
		return errors.New("only members of the organization that created this OKR template can modify it")
	}
	return nil
}

func (s *okrServiceImpl) DuplicateTemplate(db *gorm.DB, id uuid.UUID, userID uuid.UUID) (*OKRTemplateResponse, error) {
	original, err := s.repo.GetOKRTemplateWithObjectives(db, id)
	if err != nil {
		return nil, err
	}
	// Duplikasi mengikuti aturan otorisasi yang sama dengan edit/hapus — hanya
	// anggota organisasi pembuat yang boleh menduplikasi template.
	if err := s.authorizeTemplateOrg(db, original.CreatedByOrgID, userID); err != nil {
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
	// Duplikat adalah template baru — organisasi pembuatnya mengikuti user yang
	// melakukan duplikasi (bukan menyalin original), konsisten dengan create.
	// Best-effort: jika gagal resolve, nil → fallback ke organization_id.
	if _, creatorOrg, err := s.repo.GetCurrentEmployeeContext(db, userID); err == nil && creatorOrg != nil {
		newTemplate.CreatedByOrgID = creatorOrg
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

	// Objectives are "received" implicitly via TemplateID (the frontend reads
	// GET /okr/templates/:id/objectives to know which Objectives to propose
	// Key Results under) — Key Results themselves are NOT copied from the
	// Template here. The employee proposes their own Key Results per
	// Objective while the evaluation is DRAFT (see CreateEvaluationKeyResult),
	// mirroring KPI's employee-authored Program items.
	_ = template

	return s.GetEvaluationWithDetails(db, evaluation.ID)
}

// =========================================================================
// Employee-Proposed Key Results (DRAFT phase)
// =========================================================================

// CreateEvaluationKeyResult adds an employee-proposed Key Result under one
// of the evaluation's snapshotted Objectives, only while the evaluation is
// DRAFT. Total Key Result weight within a single Objective cannot exceed
// 100%. Mirrors performance.Service.CreateProgramItem for KPI.
func (s *okrServiceImpl) CreateEvaluationKeyResult(db *gorm.DB, req *CreateOKREvaluationKeyResultRequest) (*OKREvaluationDetailResponse, error) {
	evalID, err := uuid.Parse(req.EvaluationID)
	if err != nil {
		return nil, errors.New("invalid evaluation_id")
	}
	objID, err := uuid.Parse(req.ObjectiveID)
	if err != nil {
		return nil, errors.New("invalid objective_id")
	}

	evaluation, err := s.repo.GetOKREvaluationByID(db, evalID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusDraft {
		return nil, fmt.Errorf("key results can only be proposed while the evaluation is in DRAFT status, current: %s", evaluation.Status)
	}

	existing, err := s.repo.ListOKREvaluationDetailsByEvaluationID(db, evalID)
	if err != nil {
		return nil, err
	}
	existingWeight := 0.0
	sortOrder := 0
	for _, d := range existing {
		if d.ObjectiveID != nil && *d.ObjectiveID == objID {
			existingWeight += d.KeyResultWeight
			sortOrder++
		}
	}
	if existingWeight+req.Weight > 100 {
		return nil, fmt.Errorf("total weight of key results for this objective cannot exceed 100%%, currently at %.2f%%", existingWeight)
	}

	targetType := TargetType(req.TargetType)
	if targetType == "" {
		targetType = TargetTypeNumber
	}
	formulaType := FormulaType(req.FormulaType)
	if formulaType == "" {
		formulaType = FormulaTypeHigherBetter
	}

	detail := &OKREvaluationDetail{
		EvaluationID:    evalID,
		ObjectiveID:     &objID,
		ObjectiveTitle:  req.ObjectiveTitle,
		KeyResultTitle:  req.Title,
		ObjectiveWeight: req.ObjectiveWeight,
		KeyResultWeight: req.Weight,
		TargetValue:     req.TargetValue,
		TargetType:      targetType,
		Unit:            req.Unit,
		FormulaType:     formulaType,
		SortOrder:       sortOrder,
	}
	if err := s.repo.CreateOKREvaluationDetail(db, detail); err != nil {
		return nil, err
	}
	return s.evaluationDetailToResponse(detail), nil
}

// UpdateEvaluationKeyResultTarget edits an employee-proposed Key Result,
// only while the evaluation is DRAFT.
func (s *okrServiceImpl) UpdateEvaluationKeyResultTarget(db *gorm.DB, id uuid.UUID, req *UpdateOKREvaluationKeyResultTargetRequest) (*OKREvaluationDetailResponse, error) {
	detail, err := s.repo.GetOKREvaluationDetailByID(db, id)
	if err != nil {
		return nil, err
	}
	evaluation, err := s.repo.GetOKREvaluationByID(db, detail.EvaluationID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusDraft {
		return nil, fmt.Errorf("key result target can only be edited while the evaluation is in DRAFT status, current: %s", evaluation.Status)
	}

	if req.Weight != nil && *req.Weight != detail.KeyResultWeight {
		others, err := s.repo.ListOKREvaluationDetailsByEvaluationID(db, detail.EvaluationID)
		if err != nil {
			return nil, err
		}
		otherWeight := 0.0
		for _, d := range others {
			if d.ID != detail.ID && d.ObjectiveID != nil && detail.ObjectiveID != nil && *d.ObjectiveID == *detail.ObjectiveID {
				otherWeight += d.KeyResultWeight
			}
		}
		if otherWeight+*req.Weight > 100 {
			return nil, fmt.Errorf("total weight of key results for this objective cannot exceed 100%%, other key results already total %.2f%%", otherWeight)
		}
	}

	if req.Title != nil {
		detail.KeyResultTitle = *req.Title
	}
	if req.TargetType != nil {
		detail.TargetType = TargetType(*req.TargetType)
	}
	if req.TargetValue != nil {
		detail.TargetValue = *req.TargetValue
	}
	if req.Unit != nil {
		detail.Unit = req.Unit
	}
	if req.FormulaType != nil {
		detail.FormulaType = FormulaType(*req.FormulaType)
	}
	if req.Weight != nil {
		detail.KeyResultWeight = *req.Weight
	}

	if err := s.repo.UpdateOKREvaluationDetail(db, detail); err != nil {
		return nil, err
	}
	return s.evaluationDetailToResponse(detail), nil
}

// DeleteEvaluationKeyResult removes an employee-proposed Key Result, only
// while the evaluation is DRAFT.
func (s *okrServiceImpl) DeleteEvaluationKeyResult(db *gorm.DB, id uuid.UUID) error {
	detail, err := s.repo.GetOKREvaluationDetailByID(db, id)
	if err != nil {
		return err
	}
	evaluation, err := s.repo.GetOKREvaluationByID(db, detail.EvaluationID)
	if err != nil {
		return err
	}
	if evaluation.Status != OKRStatusDraft {
		return fmt.Errorf("key results can only be removed while the evaluation is in DRAFT status, current: %s", evaluation.Status)
	}
	return s.repo.DeleteOKREvaluationDetail(db, id)
}

func (s *okrServiceImpl) GetEvaluationByID(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, id)
	if err != nil {
		return nil, err
	}
	resp := s.evaluationToResponse(evaluation)
	if err := s.enrichEvaluationResponses(db, []*OKREvaluationResponse{resp}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *okrServiceImpl) GetEvaluationWithDetails(db *gorm.DB, id uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationWithDetails(db, id)
	if err != nil {
		return nil, err
	}
	resp := s.evaluationToResponseWithDetails(evaluation)
	if err := s.enrichEvaluationResponses(db, []*OKREvaluationResponse{resp}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *okrServiceImpl) ListEvaluations(db *gorm.DB, employeeID, orgID, periodID *uuid.UUID, status *string, page, perPage int) ([]OKREvaluationResponse, int64, error) {
	evaluations, total, err := s.repo.ListOKREvaluations(db, employeeID, orgID, periodID, status, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]OKREvaluationResponse, len(evaluations))
	refs := make([]*OKREvaluationResponse, len(evaluations))
	for i, e := range evaluations {
		responses[i] = *s.evaluationToResponse(&e)
		refs[i] = &responses[i]
	}
	if err := s.enrichEvaluationResponses(db, refs); err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}

// enrichEvaluationResponses batch-populates EmployeeName and
// OrganizationName (Organization = Position in this platform's convention,
// so this doubles as the "jabatan" shown on the evaluation detail page) on
// a set of already-built responses — mirrors the equivalent KPI evaluation
// enrichment in performance/service.go.
func (s *okrServiceImpl) enrichEvaluationResponses(db *gorm.DB, responses []*OKREvaluationResponse) error {
	if len(responses) == 0 {
		return nil
	}

	empIDSet := make(map[uuid.UUID]struct{})
	orgIDSet := make(map[uuid.UUID]struct{})
	periodIDSet := make(map[uuid.UUID]struct{})
	for _, r := range responses {
		if id, err := uuid.Parse(r.EmployeeID); err == nil {
			empIDSet[id] = struct{}{}
		}
		if id, err := uuid.Parse(r.OrganizationID); err == nil {
			orgIDSet[id] = struct{}{}
		}
		if id, err := uuid.Parse(r.PeriodID); err == nil {
			periodIDSet[id] = struct{}{}
		}
	}

	empIDs := make([]uuid.UUID, 0, len(empIDSet))
	for id := range empIDSet {
		empIDs = append(empIDs, id)
	}
	orgIDs := make([]uuid.UUID, 0, len(orgIDSet))
	for id := range orgIDSet {
		orgIDs = append(orgIDs, id)
	}
	periodIDs := make([]uuid.UUID, 0, len(periodIDSet))
	for id := range periodIDSet {
		periodIDs = append(periodIDs, id)
	}

	empNames, err := s.repo.GetEmployeeNamesByIDs(db, empIDs)
	if err != nil {
		return err
	}
	orgNames, err := s.repo.GetOrganizationNamesByIDs(db, orgIDs)
	if err != nil {
		return err
	}
	periodCodes, err := s.repo.GetPeriodCodesByIDs(db, periodIDs)
	if err != nil {
		return err
	}

	for _, r := range responses {
		if name, ok := empNames[r.EmployeeID]; ok {
			r.EmployeeName = name
		}
		if name, ok := orgNames[r.OrganizationID]; ok {
			r.OrganizationName = name
		}
		if code, ok := periodCodes[r.PeriodID]; ok {
			r.PeriodCode = code
		}
	}
	return nil
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
	evaluation, err := s.repo.GetOKREvaluationByID(db, detail.EvaluationID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusKRApproved {
		return nil, fmt.Errorf("actual can only be filled once key results are approved, current status: %s", evaluation.Status)
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
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return err
	}
	if evaluation.Status != OKRStatusKRApproved {
		return fmt.Errorf("actuals can only be filled once key results are approved, current status: %s", evaluation.Status)
	}

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
// Workflow — Key Result proposal phase (DRAFT -> KR_SUBMITTED -> KR_APPROVED)
// =========================================================================

// SubmitKeyResults moves an evaluation from DRAFT to KR_SUBMITTED — the
// employee's proposed Key Results are sent to their supervisor for review.
// Requires every snapshotted Objective to have at least one Key Result with
// weight > 0. Mirrors performance.Service.SubmitTarget for KPI.
func (s *okrServiceImpl) SubmitKeyResults(ctx context.Context, db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusDraft {
		return nil, fmt.Errorf("key results can only be submitted from DRAFT status, current: %s", evaluation.Status)
	}

	details, err := s.repo.ListOKREvaluationDetailsByEvaluationID(db, evaluationID)
	if err != nil {
		return nil, err
	}
	if len(details) == 0 {
		return nil, errors.New("at least one key result must be proposed before submitting")
	}
	weightByObjective := make(map[uuid.UUID]float64)
	for _, d := range details {
		if d.ObjectiveID != nil {
			weightByObjective[*d.ObjectiveID] += d.KeyResultWeight
		}
	}
	for _, w := range weightByObjective {
		if w <= 0 {
			return nil, errors.New("each objective must have at least one key result with weight greater than 0")
		}
	}

	now := time.Now()
	evaluation.Status = OKRStatusKRSubmitted
	evaluation.KRSubmittedAt = &now

	// Route through the central approval module when a flow is configured
	// for this module; ApproveKeyResults/RejectKeyResults remain the manual
	// fallback only when no flow is configured at all. If a flow IS
	// configured but fails to resolve, the whole submission must fail
	// rather than silently degrading to manual mode — mirrors
	// performance.Service.SubmitTarget for KPI.
	if s.approvalEngine != nil {
		if flowID, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, ApprovalModuleOKRKeyResult); err == nil {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, ApprovalModuleOKRKeyResult, evaluation.ID.String(), flowID)
			if err != nil {
				return nil, fmt.Errorf("failed to route key results for approval: %w", err)
			}
			parsedInstanceID, err := uuid.Parse(instanceID)
			if err != nil {
				return nil, fmt.Errorf("invalid approval instance id returned: %w", err)
			}
			evaluation.KRApprovalInstanceID = &parsedInstanceID
		}
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// ApproveKeyResults moves KR_SUBMITTED -> KR_APPROVED ("OKR Active") — the
// manual-fallback path used when no central approval flow is configured.
func (s *okrServiceImpl) ApproveKeyResults(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusKRSubmitted {
		return nil, fmt.Errorf("only key-result-submitted evaluations can be approved, current: %s", evaluation.Status)
	}

	evaluation.Status = OKRStatusKRApproved

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// RejectKeyResults reverts KR_SUBMITTED -> DRAFT so the employee can revise
// their proposed Key Results.
func (s *okrServiceImpl) RejectKeyResults(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID, notes string) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}
	if evaluation.Status != OKRStatusKRSubmitted {
		return nil, fmt.Errorf("only key-result-submitted evaluations can be rejected, current: %s", evaluation.Status)
	}

	evaluation.Status = OKRStatusDraft
	if notes != "" {
		evaluation.ReviewerNotes = &notes
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// =========================================================================
// Workflow — Assessment phase (KR_APPROVED -> SUBMITTED -> COMPLETED)
// =========================================================================

// SubmitEvaluation moves an evaluation from KR_APPROVED ("OKR Active", after
// check-in and self-assessment actuals are filled) to SUBMITTED — the
// employee's final self-assessment sent to their supervisor.
func (s *okrServiceImpl) SubmitEvaluation(ctx context.Context, db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusKRApproved {
		return nil, fmt.Errorf("evaluation can only be submitted once its key results are approved, current: %s", evaluation.Status)
	}

	now := time.Now()
	evaluation.Status = OKRStatusSubmitted
	evaluation.SubmittedAt = &now
	evaluation.SubmittedBy = &userID

	// Same hard-fail-if-configured-but-unresolvable semantics as SubmitKeyResults.
	if s.approvalEngine != nil {
		if flowID, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, ApprovalModuleOKRAssessment); err == nil {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, ApprovalModuleOKRAssessment, evaluation.ID.String(), flowID)
			if err != nil {
				return nil, fmt.Errorf("failed to route assessment for approval: %w", err)
			}
			parsedInstanceID, err := uuid.Parse(instanceID)
			if err != nil {
				return nil, fmt.Errorf("invalid approval instance id returned: %w", err)
			}
			evaluation.AssessmentApprovalInstanceID = &parsedInstanceID
		}
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// ApproveEvaluation changes status from SUBMITTED directly to COMPLETED —
// this is the final approval step for the assessment, so there's no
// separate manual "Complete" action needed afterward. Mirrors
// performance.Service.ApproveEvaluation for KPI.
func (s *okrServiceImpl) ApproveEvaluation(db *gorm.DB, evaluationID uuid.UUID, userID uuid.UUID) (*OKREvaluationResponse, error) {
	evaluation, err := s.repo.GetOKREvaluationByID(db, evaluationID)
	if err != nil {
		return nil, err
	}

	if evaluation.Status != OKRStatusSubmitted {
		return nil, errors.New("only submitted evaluations can be approved")
	}

	now := time.Now()
	evaluation.Status = OKRStatusCompleted
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

	// Reverts to KR_APPROVED (not all the way to DRAFT) — the Key Results
	// were already approved separately, the employee only needs to revise
	// their self-assessment actuals.
	evaluation.Status = OKRStatusKRApproved
	if notes != "" {
		evaluation.ReviewerNotes = &notes
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return nil, err
	}

	return s.evaluationToResponse(evaluation), nil
}

// HandleKeyResultApprovalStatusChange is invoked by the approval module's
// push-based status callback when a Key Result approval instance reaches a
// final state, so the evaluation's own status updates itself. Mirrors
// performance.Service.HandleTargetApprovalStatusChange for KPI.
func (s *okrServiceImpl) HandleKeyResultApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	db, err := s.dbResolver(ctx)
	if err != nil {
		return err
	}
	evaluation, err := s.repo.GetOKREvaluationByID(db, documentID)
	if err != nil {
		return err
	}
	if evaluation.Status != OKRStatusKRSubmitted {
		return nil
	}

	switch status {
	case "APPROVED":
		evaluation.Status = OKRStatusKRApproved
	case "REJECTED", "CANCELLED":
		evaluation.Status = OKRStatusDraft
		evaluation.KRSubmittedAt = nil
		if note != "" {
			evaluation.ReviewerNotes = &note
		}
	default:
		return nil
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return err
	}
	switch status {
	case "APPROVED":
		s.notifyEvaluationOutcome(ctx, db, evaluation.EmployeeID, "OKR_KEY_RESULT_APPROVED", "okr_key_result", evaluation.ID)
	case "REJECTED", "CANCELLED":
		s.notifyEvaluationOutcome(ctx, db, evaluation.EmployeeID, "OKR_KEY_RESULT_REJECTED", "okr_key_result", evaluation.ID)
	}
	return nil
}

// HandleAssessmentApprovalStatusChange is invoked by the approval module's
// push-based status callback when an assessment approval instance reaches a
// final state. Final approval completes the evaluation directly. Mirrors
// performance.Service.HandleRealizationApprovalStatusChange for KPI.
func (s *okrServiceImpl) HandleAssessmentApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	db, err := s.dbResolver(ctx)
	if err != nil {
		return err
	}
	evaluation, err := s.repo.GetOKREvaluationByID(db, documentID)
	if err != nil {
		return err
	}
	if evaluation.Status != OKRStatusSubmitted {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		evaluation.Status = OKRStatusCompleted
		evaluation.ApprovedAt = &now
	case "REJECTED", "CANCELLED":
		evaluation.Status = OKRStatusKRApproved
		evaluation.SubmittedAt = nil
		if note != "" {
			evaluation.ReviewerNotes = &note
		}
	default:
		return nil
	}

	if err := s.repo.UpdateOKREvaluation(db, evaluation); err != nil {
		return err
	}
	switch status {
	case "APPROVED":
		s.notifyEvaluationOutcome(ctx, db, evaluation.EmployeeID, "OKR_ASSESSMENT_APPROVED", "okr_assessment", evaluation.ID)
	case "REJECTED", "CANCELLED":
		s.notifyEvaluationOutcome(ctx, db, evaluation.EmployeeID, "OKR_ASSESSMENT_REJECTED", "okr_assessment", evaluation.ID)
	}
	return nil
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
	if t.CreatedBy != uuid.Nil {
		resp.CreatedBy = t.CreatedBy.String()
	}
	if t.CreatedByOrgID != nil {
		resp.CreatedByOrgID = t.CreatedByOrgID.String()
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
	if e.KRApprovalInstanceID != nil {
		resp.KRApprovalInstanceID = e.KRApprovalInstanceID.String()
	}
	if e.AssessmentApprovalInstanceID != nil {
		resp.AssessmentApprovalInstanceID = e.AssessmentApprovalInstanceID.String()
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
