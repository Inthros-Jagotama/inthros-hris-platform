package recruitment

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

// WorkforceGapProvider adalah interface narrow yang dipakai Recruitment untuk
// membaca hiring need dari Workforce Intelligence (plan S-1 — workforce gap →
// requisition). Recruitment TIDAK menghitung gap sendiri; ia hanya membaca
// hasil perhitungan WI. Implementasi di-wire di cmd/server/main.go melalui
// adapter (workforceintelligence.Service). Bila provider nil, requisition
// tetap bisa dibuat dengan slots_available default tanpa error.
type WorkforceGapProvider interface {
	// HiringGapForOrganization mengembalikan jumlah hiring need (shortage)
	// untuk sebuah organisasi. Nilai positif = jumlah slot yang harus
	// di-recruit; 0 = tidak ada shortage.
	HiringGapForOrganization(ctx context.Context, orgID uuid.UUID) (int, error)
}

type Service struct {
	repo        *Repository
	logger      *zap.Logger
	gapProvider WorkforceGapProvider
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// SetWorkforceGapProvider wires the Workforce Intelligence module into this
// service (plan S-1) so requisitions with reason_type=WORKFORCE_GAP can
// auto-resolve slots_available from WI's hiring need.
func (s *Service) SetWorkforceGapProvider(p WorkforceGapProvider) {
	s.gapProvider = p
}

// =========================================================================
// Job Requisitions
// =========================================================================

func (s *Service) CreateRequisition(ctx context.Context, req CreateRequisitionRequest) (*RequisitionResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	r := &JobRequisition{
		OrganizationID:  orgID,
		Title:           req.Title,
		Department:      req.Department,
		EmploymentType:  req.EmploymentType,
		Location:        req.Location,
		MinSalary:       req.MinSalary,
		MaxSalary:       req.MaxSalary,
		Description:     req.Description,
		Requirements:    req.Requirements,
		Responsibilities: req.Responsibilities,
		SlotsAvailable:  1,
	}
	if req.SlotsAvailable != nil {
		r.SlotsAvailable = *req.SlotsAvailable
	}
	if req.RequestedBy != nil {
		uid, _ := uuid.Parse(*req.RequestedBy)
		r.RequestedBy = &uid
	}
	if req.TargetStartDate != nil {
		r.TargetStartDate = req.TargetStartDate
	}
	if req.ReasonType != nil {
		r.ReasonType = *req.ReasonType
	}
	if req.WorkforceGapID != nil {
		uid, _ := uuid.Parse(*req.WorkforceGapID)
		r.WorkforceGapID = &uid
	}
	if req.WorkforcePlanID != nil {
		uid, _ := uuid.Parse(*req.WorkforcePlanID)
		r.WorkforcePlanID = &uid
	}
	// S-1: auto-resolve hiring need dari Workforce Intelligence ketika
	// requisition dibuat dengan reason WORKFORCE_GAP dan slots tidak
	// ditentukan eksplisit. Gagal resolve = tetap lanjut dengan default.
	s.resolveWorkforceGapSlots(ctx, r, req.SlotsAvailable)
	if err := s.repo.CreateRequisition(ctx, r); err != nil {
		return nil, err
	}
	s.logger.Info("Requisition created", zap.String("id", r.ID.String()), zap.String("title", r.Title))
	return requisitionToResponse(r), nil
}

func (s *Service) GetRequisitionByID(ctx context.Context, id string) (*RequisitionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindRequisitionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return requisitionToResponse(r), nil
}

func (s *Service) ListRequisitions(ctx context.Context, orgID, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var orgUUID *uuid.UUID
	if orgID != nil && *orgID != "" {
		uid, _ := uuid.Parse(*orgID)
		orgUUID = &uid
	}
	list, total, err := s.repo.ListRequisitions(ctx, orgUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]RequisitionResponse, 0, len(list))
	for _, r := range list {
		responses = append(responses, *requisitionToResponse(&r))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateRequisition(ctx context.Context, id string, req UpdateRequisitionRequest) (*RequisitionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindRequisitionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		r.Title = *req.Title
	}
	if req.Department != nil {
		r.Department = *req.Department
	}
	if req.EmploymentType != nil {
		r.EmploymentType = *req.EmploymentType
	}
	if req.Location != nil {
		r.Location = *req.Location
	}
	if req.MinSalary != nil {
		r.MinSalary = *req.MinSalary
	}
	if req.MaxSalary != nil {
		r.MaxSalary = *req.MaxSalary
	}
	if req.Description != nil {
		r.Description = *req.Description
	}
	if req.Requirements != nil {
		r.Requirements = *req.Requirements
	}
	if req.Responsibilities != nil {
		r.Responsibilities = *req.Responsibilities
	}
	if req.SlotsAvailable != nil {
		r.SlotsAvailable = *req.SlotsAvailable
	}
	if req.Status != nil {
		r.Status = RequisitionStatus(*req.Status)
	}
	if req.TargetStartDate != nil {
		r.TargetStartDate = req.TargetStartDate
	}
	// Catat reason lama SEBELUM diubah, untuk deteksi transisi ke WORKFORCE_GAP.
	prevReason := r.ReasonType
	if req.ReasonType != nil {
		r.ReasonType = *req.ReasonType
	}
	if req.WorkforceGapID != nil {
		if *req.WorkforceGapID == "" {
			r.WorkforceGapID = nil
		} else {
			uid, _ := uuid.Parse(*req.WorkforceGapID)
			r.WorkforceGapID = &uid
		}
	}
	if req.WorkforcePlanID != nil {
		if *req.WorkforcePlanID == "" {
			r.WorkforcePlanID = nil
		} else {
			uid, _ := uuid.Parse(*req.WorkforcePlanID)
			r.WorkforcePlanID = &uid
		}
	}
	// S-1: resolve hiring need HANYA saat reason bertransisi menjadi
	// WORKFORCE_GAP (bukan sudah WORKFORCE_GAP) — sehingga update field lain
	// (title, status, dll.) pada requisition yang sudah tertaut gap tidak
	// menimpa slots_available yang tersimpan.
	if req.ReasonType != nil && *req.ReasonType == string(ReqReasonWorkforceGap) && prevReason != string(ReqReasonWorkforceGap) {
		s.resolveWorkforceGapSlots(ctx, r, req.SlotsAvailable)
	}
	if err := s.repo.UpdateRequisition(ctx, r); err != nil {
		return nil, err
	}
	return requisitionToResponse(r), nil
}

func (s *Service) DeleteRequisition(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteRequisition(ctx, uid)
}

// =========================================================================
// Candidates
// =========================================================================

func (s *Service) CreateCandidate(ctx context.Context, req CreateCandidateRequest) (*CandidateResponse, error) {
	// Check duplicate email
	existing, err := s.repo.FindCandidateByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("candidate with email %s already exists", req.Email)
	}

	c := &Candidate{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Source:    "direct",
	}
	if req.CurrentCompany != nil {
		c.CurrentCompany = req.CurrentCompany
	}
	if req.CurrentTitle != nil {
		c.CurrentTitle = req.CurrentTitle
	}
	if req.ResumeURL != nil {
		c.ResumeURL = req.ResumeURL
	}
	if req.PortfolioURL != nil {
		c.PortfolioURL = req.PortfolioURL
	}
	if req.LinkedInURL != nil {
		c.LinkedInURL = req.LinkedInURL
	}
	if req.Source != nil {
		c.Source = *req.Source
	}
	if req.Notes != "" {
		c.Notes = req.Notes
	}
	if err := s.repo.CreateCandidate(ctx, c); err != nil {
		return nil, err
	}
	s.logger.Info("Candidate created", zap.String("id", c.ID.String()), zap.String("name", c.FirstName+" "+c.LastName))
	return candidateToResponse(c), nil
}

func (s *Service) GetCandidateByID(ctx context.Context, id string) (*CandidateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCandidateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return candidateToResponse(c), nil
}

func (s *Service) ListCandidates(ctx context.Context, search *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListCandidates(ctx, page, perPage, search)
	if err != nil {
		return nil, err
	}
	responses := make([]CandidateResponse, 0, len(list))
	for _, c := range list {
		responses = append(responses, *candidateToResponse(&c))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateCandidate(ctx context.Context, id string, req UpdateCandidateRequest) (*CandidateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCandidateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.FirstName != nil {
		c.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		c.LastName = *req.LastName
	}
	if req.Email != nil {
		c.Email = *req.Email
	}
	if req.Phone != nil {
		c.Phone = *req.Phone
	}
	if req.Address != nil {
		c.Address = *req.Address
	}
	if req.CurrentCompany != nil {
		c.CurrentCompany = req.CurrentCompany
	}
	if req.CurrentTitle != nil {
		c.CurrentTitle = req.CurrentTitle
	}
	if req.ResumeURL != nil {
		c.ResumeURL = req.ResumeURL
	}
	if req.PortfolioURL != nil {
		c.PortfolioURL = req.PortfolioURL
	}
	if req.LinkedInURL != nil {
		c.LinkedInURL = req.LinkedInURL
	}
	if req.Source != nil {
		c.Source = *req.Source
	}
	if req.Notes != nil {
		c.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidate(ctx, c); err != nil {
		return nil, err
	}
	return candidateToResponse(c), nil
}

func (s *Service) DeleteCandidate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidate(ctx, uid)
}

// =========================================================================
// Job Applications
// =========================================================================

func (s *Service) CreateApplication(ctx context.Context, req CreateApplicationRequest) (*ApplicationResponse, error) {
	reqID, err := uuid.Parse(req.RequisitionID)
	if err != nil {
		return nil, fmt.Errorf("invalid requisition_id: %w", err)
	}
	candID, err := uuid.Parse(req.CandidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}

	// Check requisition exists
	_, err = s.repo.FindRequisitionByID(ctx, reqID)
	if err != nil {
		return nil, fmt.Errorf("requisition not found: %w", err)
	}
	// Check candidate exists
	_, err = s.repo.FindCandidateByID(ctx, candID)
	if err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	a := &JobApplication{
		RequisitionID: reqID,
		CandidateID:   candID,
		Status:        CandStatusNew,
		AppliedAt:     time.Now().UnixNano(),
		Notes:         req.Notes,
	}
	if err := s.repo.CreateApplication(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Job application created", zap.String("id", a.ID.String()))
	return applicationToResponse(a), nil
}

func (s *Service) GetApplicationByID(ctx context.Context, id string) (*ApplicationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindApplicationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return applicationToResponse(a), nil
}

func (s *Service) ListApplications(ctx context.Context, requisitionID, candidateID, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var reqUUID, candUUID *uuid.UUID
	if requisitionID != nil && *requisitionID != "" {
		uid, _ := uuid.Parse(*requisitionID)
		reqUUID = &uid
	}
	if candidateID != nil && *candidateID != "" {
		uid, _ := uuid.Parse(*candidateID)
		candUUID = &uid
	}
	list, total, err := s.repo.ListApplications(ctx, reqUUID, candUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ApplicationResponse, 0, len(list))
	for _, a := range list {
		responses = append(responses, *applicationToResponse(&a))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateApplicationStatus(ctx context.Context, id, status, reason, notes string) (*ApplicationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindApplicationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixNano()
	a.Status = CandidateStatus(status)
	if reason != "" {
		a.RejectionReason = reason
	}
	if notes != "" {
		a.Notes = notes
	}

	switch CandidateStatus(status) {
	case CandStatusScreened:
		a.ScreenedAt = &now
	case CandStatusShortlisted:
		a.ShortlistedAt = &now
	case CandStatusOffered:
		a.OfferedAt = &now
	case CandStatusAccepted:
		a.AcceptedAt = &now
		// Update requisition slots filled
		req, findErr := s.repo.FindRequisitionByID(ctx, a.RequisitionID)
		if findErr == nil && req != nil {
			req.SlotsFilled++
			if req.SlotsFilled >= req.SlotsAvailable {
				req.Status = ReqStatusFilled
			}
			if err := s.repo.UpdateRequisition(ctx, req); err != nil {
				s.logger.Warn("failed to update requisition slots_filled", zap.String("requisition_id", req.ID.String()), zap.Error(err))
			}
		}
	case CandStatusRejected:
		a.RejectedAt = &now
	case CandStatusWithdrawn:
		a.WithdrawnAt = &now
	}

	if err := s.repo.UpdateApplication(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Application status updated", zap.String("id", a.ID.String()), zap.String("status", string(a.Status)))
	return applicationToResponse(a), nil
}

func (s *Service) DeleteApplication(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteApplication(ctx, uid)
}

// =========================================================================
// Interviews
// =========================================================================

func (s *Service) CreateInterview(ctx context.Context, req CreateInterviewRequest) (*InterviewResponse, error) {
	appID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	intvID, err := uuid.Parse(req.InterviewerID)
	if err != nil {
		return nil, fmt.Errorf("invalid interviewer_id: %w", err)
	}

	i := &Interview{
		ApplicationID: appID,
		InterviewerID: intvID,
		Stage:         req.Stage,
		ScheduledAt:   req.ScheduledAt,
		Status:        IntStatusScheduled,
	}
	if i.Stage == "" {
		i.Stage = "FIRST_INTERVIEW"
	}
	if req.DurationMinutes != nil {
		i.DurationMinutes = *req.DurationMinutes
	} else {
		i.DurationMinutes = 60
	}
	if req.Location != "" {
		i.Location = req.Location
	}
	if req.MeetingLink != "" {
		i.MeetingLink = &req.MeetingLink
	}
	if err := s.repo.CreateInterview(ctx, i); err != nil {
		return nil, err
	}
	s.logger.Info("Interview created", zap.String("id", i.ID.String()), zap.String("stage", i.Stage))
	return interviewToResponse(i), nil
}

func (s *Service) GetInterviewByID(ctx context.Context, id string) (*InterviewResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindInterviewByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return interviewToResponse(i), nil
}

func (s *Service) ListInterviews(ctx context.Context, applicationID, interviewerID *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var appUUID, intvUUID *uuid.UUID
	if applicationID != nil && *applicationID != "" {
		uid, _ := uuid.Parse(*applicationID)
		appUUID = &uid
	}
	if interviewerID != nil && *interviewerID != "" {
		uid, _ := uuid.Parse(*interviewerID)
		intvUUID = &uid
	}
	list, total, err := s.repo.ListInterviews(ctx, appUUID, intvUUID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]InterviewResponse, 0, len(list))
	for _, i := range list {
		responses = append(responses, *interviewToResponse(&i))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateInterview(ctx context.Context, id string, req UpdateInterviewRequest) (*InterviewResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	i, err := s.repo.FindInterviewByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.InterviewerID != nil {
		uid, _ := uuid.Parse(*req.InterviewerID)
		i.InterviewerID = uid
	}
	if req.Stage != nil {
		i.Stage = *req.Stage
	}
	if req.ScheduledAt != nil {
		i.ScheduledAt = *req.ScheduledAt
	}
	if req.DurationMinutes != nil {
		i.DurationMinutes = *req.DurationMinutes
	}
	if req.Location != nil {
		i.Location = *req.Location
	}
	if req.MeetingLink != nil {
		i.MeetingLink = req.MeetingLink
	}
	if req.Status != nil {
		i.Status = InterviewStatus(*req.Status)
		if *req.Status == "COMPLETED" {
			now := time.Now().UnixNano()
			i.CompletedAt = &now
		}
	}
	if req.Score != nil {
		i.Score = req.Score
	}
	if req.Feedback != nil {
		i.Feedback = *req.Feedback
	}
	if err := s.repo.UpdateInterview(ctx, i); err != nil {
		return nil, err
	}
	return interviewToResponse(i), nil
}

func (s *Service) DeleteInterview(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteInterview(ctx, uid)
}

// =========================================================================
// Onboarding Task Templates
// =========================================================================

func (s *Service) CreateOnboardingTaskTemplate(ctx context.Context, req CreateOnboardingTaskTemplateRequest) (*OnboardingTaskTemplateResponse, error) {
	t := &OnboardingTaskTemplate{
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		AssignedRole: req.AssignedRole,
		IsMandatory:  true,
	}
	if req.DayOffset != nil {
		t.DayOffset = *req.DayOffset
	}
	if req.IsMandatory != nil {
		t.IsMandatory = *req.IsMandatory
	}
	if err := s.repo.CreateOnboardingTaskTemplate(ctx, t); err != nil {
		return nil, err
	}
	return taskTemplateToResponse(t), nil
}

func (s *Service) ListOnboardingTaskTemplates(ctx context.Context, category *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListOnboardingTaskTemplates(ctx, category, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]OnboardingTaskTemplateResponse, 0, len(list))
	for _, t := range list {
		responses = append(responses, *taskTemplateToResponse(&t))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateOnboardingTaskTemplate(ctx context.Context, id string, req UpdateOnboardingTaskTemplateRequest) (*OnboardingTaskTemplateResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindOnboardingTaskTemplateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.DayOffset != nil {
		t.DayOffset = *req.DayOffset
	}
	if req.AssignedRole != nil {
		t.AssignedRole = *req.AssignedRole
	}
	if req.IsMandatory != nil {
		t.IsMandatory = *req.IsMandatory
	}
	if err := s.repo.UpdateOnboardingTaskTemplate(ctx, t); err != nil {
		return nil, err
	}
	return taskTemplateToResponse(t), nil
}

func (s *Service) DeleteOnboardingTaskTemplate(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteOnboardingTaskTemplate(ctx, uid)
}

// =========================================================================
// Employee Onboarding
// =========================================================================

func (s *Service) CreateEmployeeOnboarding(ctx context.Context, req CreateEmployeeOnboardingRequest) (*EmployeeOnboardingResponse, error) {
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	appID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	o := &EmployeeOnboarding{
		EmployeeID:    empID,
		ApplicationID: appID,
		StartDate:     req.StartDate,
		Status:        "PENDING",
	}
	if req.BuddyID != nil {
		uid, _ := uuid.Parse(*req.BuddyID)
		o.BuddyID = &uid
	}
	if req.Notes != "" {
		o.Notes = req.Notes
	}
	if err := s.repo.CreateEmployeeOnboarding(ctx, o); err != nil {
		return nil, err
	}

	// Auto-create task items from templates
	templates, _, _ := s.repo.ListOnboardingTaskTemplates(ctx, nil, 1, 100)
	for _, t := range templates {
		item := &OnboardingTaskItem{
			EmployeeOnboardingID: o.ID,
			TemplateID:           &t.ID,
			Name:                 t.Name,
			Description:          t.Description,
			AssignedTo:           o.BuddyID,
			IsCompleted:          false,
		}
		s.repo.CreateOnboardingTaskItem(ctx, item)
	}

	s.logger.Info("Employee onboarding created", zap.String("id", o.ID.String()))
	return onboardingToResponse(o), nil
}

func (s *Service) GetEmployeeOnboardingByID(ctx context.Context, id string) (*EmployeeOnboardingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindEmployeeOnboardingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return onboardingToResponse(o), nil
}

func (s *Service) ListEmployeeOnboardings(ctx context.Context, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListEmployeeOnboardings(ctx, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]EmployeeOnboardingResponse, 0, len(list))
	for _, o := range list {
		responses = append(responses, *onboardingToResponse(&o))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateEmployeeOnboarding(ctx context.Context, id string, req UpdateEmployeeOnboardingRequest) (*EmployeeOnboardingResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindEmployeeOnboardingByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.StartDate != nil {
		o.StartDate = *req.StartDate
	}
	if req.BuddyID != nil {
		if *req.BuddyID != "" {
			uid, _ := uuid.Parse(*req.BuddyID)
			o.BuddyID = &uid
		}
	}
	if req.Status != nil {
		o.Status = *req.Status
		if *req.Status == "COMPLETED" {
			now := time.Now().UnixNano()
			o.CompletedAt = &now
		}
	}
	if req.Notes != nil {
		o.Notes = *req.Notes
	}
	if err := s.repo.UpdateEmployeeOnboarding(ctx, o); err != nil {
		return nil, err
	}
	return onboardingToResponse(o), nil
}

func (s *Service) DeleteEmployeeOnboarding(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteEmployeeOnboarding(ctx, uid)
}

// =========================================================================
// Onboarding Task Items
// =========================================================================

func (s *Service) CreateOnboardingTaskItem(ctx context.Context, req CreateOnboardingTaskItemRequest) (*OnboardingTaskItemResponse, error) {
	onbID, err := uuid.Parse(req.EmployeeOnboardingID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_onboarding_id: %w", err)
	}
	t := &OnboardingTaskItem{
		EmployeeOnboardingID: onbID,
		Name:                req.Name,
		Description:         req.Description,
		IsCompleted:         false,
	}
	if req.TemplateID != nil && *req.TemplateID != "" {
		tID, _ := uuid.Parse(*req.TemplateID)
		t.TemplateID = &tID
	}
	if req.AssignedTo != nil && *req.AssignedTo != "" {
		uid, _ := uuid.Parse(*req.AssignedTo)
		t.AssignedTo = &uid
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}
	if err := s.repo.CreateOnboardingTaskItem(ctx, t); err != nil {
		return nil, err
	}
	return taskItemToResponse(t), nil
}

func (s *Service) ListOnboardingTaskItems(ctx context.Context, onboardingID string) ([]OnboardingTaskItemResponse, error) {
	oID, err := uuid.Parse(onboardingID)
	if err != nil {
		return nil, fmt.Errorf("invalid onboarding_id: %w", err)
	}
	list, err := s.repo.ListOnboardingTaskItems(ctx, oID)
	if err != nil {
		return nil, err
	}
	responses := make([]OnboardingTaskItemResponse, 0, len(list))
	for _, t := range list {
		responses = append(responses, *taskItemToResponse(&t))
	}
	return responses, nil
}

func (s *Service) UpdateOnboardingTaskItem(ctx context.Context, id string, req UpdateOnboardingTaskItemRequest) (*OnboardingTaskItemResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	t, err := s.repo.FindOnboardingTaskItemByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.AssignedTo != nil {
		if *req.AssignedTo != "" {
			uid, _ := uuid.Parse(*req.AssignedTo)
			t.AssignedTo = &uid
		} else {
			t.AssignedTo = nil
		}
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}
	if req.IsCompleted != nil {
		t.IsCompleted = *req.IsCompleted
		if *req.IsCompleted {
			now := time.Now().UnixNano()
			t.CompletedAt = &now
		}
	}
	if req.Notes != nil {
		t.Notes = *req.Notes
	}
	if err := s.repo.UpdateOnboardingTaskItem(ctx, t); err != nil {
		return nil, err
	}
	return taskItemToResponse(t), nil
}

func (s *Service) DeleteOnboardingTaskItem(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteOnboardingTaskItem(ctx, uid)
}

// =========================================================================
// Workforce Gap resolution (S-1)
// =========================================================================

// resolveWorkforceGapSlots mengisi slots_available dari hiring need Workforce
// Intelligence ketika requisition (create/update) memakai reason WORKFORCE_GAP
// dan slots tidak ditentukan eksplisit. Fail-safe: provider nil atau error
// tidak menggagalkan operasi — slots tetap bernilai default/eksisting.
func (s *Service) resolveWorkforceGapSlots(ctx context.Context, r *JobRequisition, explicitSlots *int) {
	if r.ReasonType != string(ReqReasonWorkforceGap) {
		return
	}
	if explicitSlots != nil {
		return
	}
	if s.gapProvider == nil {
		return
	}
	need, err := s.gapProvider.HiringGapForOrganization(ctx, r.OrganizationID)
	if err != nil {
		s.logger.Warn("workforce gap provider failed; keeping default slots",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.Error(err))
		return
	}
	if need > 0 {
		r.SlotsAvailable = need
		s.logger.Info("Workforce gap slots auto-resolved",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.Int("slots_available", r.SlotsAvailable))
	}
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

func requisitionToResponse(r *JobRequisition) *RequisitionResponse {
	resp := &RequisitionResponse{
		ID:                r.ID.String(),
		OrganizationID:    r.OrganizationID.String(),
		Title:             r.Title,
		Department:        r.Department,
		EmploymentType:    r.EmploymentType,
		Location:          r.Location,
		MinSalary:         r.MinSalary,
		MaxSalary:         r.MaxSalary,
		Description:       r.Description,
		Requirements:      r.Requirements,
		Responsibilities:  r.Responsibilities,
		SlotsAvailable:    r.SlotsAvailable,
		SlotsFilled:       r.SlotsFilled,
		Status:            string(r.Status),
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
	if r.RequestedBy != nil {
		resp.RequestedBy = r.RequestedBy.String()
	}
	if r.ApprovedBy != nil {
		resp.ApprovedBy = r.ApprovedBy.String()
	}
	if r.ReasonType != "" {
		resp.ReasonType = r.ReasonType
	}
	if r.WorkforceGapID != nil {
		resp.WorkforceGapID = r.WorkforceGapID.String()
	}
	if r.WorkforcePlanID != nil {
		resp.WorkforcePlanID = r.WorkforcePlanID.String()
	}
	if r.TargetStartDate != nil {
		resp.TargetStartDate = *r.TargetStartDate
	}
	return resp
}

func candidateToResponse(c *Candidate) *CandidateResponse {
	resp := &CandidateResponse{
		ID:        c.ID.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		Source:    c.Source,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.CurrentCompany != nil {
		resp.CurrentCompany = *c.CurrentCompany
	}
	if c.CurrentTitle != nil {
		resp.CurrentTitle = *c.CurrentTitle
	}
	if c.ResumeURL != nil {
		resp.ResumeURL = *c.ResumeURL
	}
	if c.PortfolioURL != nil {
		resp.PortfolioURL = *c.PortfolioURL
	}
	if c.LinkedInURL != nil {
		resp.LinkedInURL = *c.LinkedInURL
	}
	return resp
}

func applicationToResponse(a *JobApplication) *ApplicationResponse {
	return &ApplicationResponse{
		ID:              a.ID.String(),
		RequisitionID:   a.RequisitionID.String(),
		CandidateID:     a.CandidateID.String(),
		Status:          string(a.Status),
		RejectionReason: a.RejectionReason,
		Notes:           a.Notes,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func interviewToResponse(i *Interview) *InterviewResponse {
	resp := &InterviewResponse{
		ID:              i.ID.String(),
		ApplicationID:   i.ApplicationID.String(),
		InterviewerID:   i.InterviewerID.String(),
		Stage:           i.Stage,
		DurationMinutes: i.DurationMinutes,
		Location:        i.Location,
		Status:          string(i.Status),
		Feedback:        i.Feedback,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
	}
	if i.MeetingLink != nil {
		resp.MeetingLink = *i.MeetingLink
	}
	if i.Score != nil {
		resp.Score = *i.Score
	}
	return resp
}

func taskTemplateToResponse(t *OnboardingTaskTemplate) *OnboardingTaskTemplateResponse {
	return &OnboardingTaskTemplateResponse{
		ID:           t.ID.String(),
		Name:         t.Name,
		Description:  t.Description,
		Category:     t.Category,
		DayOffset:    t.DayOffset,
		AssignedRole: t.AssignedRole,
		IsMandatory:  t.IsMandatory,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func onboardingToResponse(o *EmployeeOnboarding) *EmployeeOnboardingResponse {
	resp := &EmployeeOnboardingResponse{
		ID:            o.ID.String(),
		EmployeeID:    o.EmployeeID.String(),
		ApplicationID: o.ApplicationID.String(),
		StartDate:     o.StartDate,
		Status:        o.Status,
		Notes:         o.Notes,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
	if o.BuddyID != nil {
		resp.BuddyID = o.BuddyID.String()
	}
	return resp
}

func taskItemToResponse(t *OnboardingTaskItem) *OnboardingTaskItemResponse {
	resp := &OnboardingTaskItemResponse{
		ID:                   t.ID.String(),
		EmployeeOnboardingID: t.EmployeeOnboardingID.String(),
		Name:                t.Name,
		Description:         t.Description,
		IsCompleted:         t.IsCompleted,
		Notes:               t.Notes,
		CreatedAt:           t.CreatedAt,
		UpdatedAt:           t.UpdatedAt,
	}
	if t.TemplateID != nil {
		resp.TemplateID = t.TemplateID.String()
	}
	if t.AssignedTo != nil {
		resp.AssignedTo = t.AssignedTo.String()
	}
	return resp
}
