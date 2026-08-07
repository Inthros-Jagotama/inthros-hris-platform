package approval

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// ModuleSubscriptionChecker abstracts checking which modules a tenant has
// subscribed to. Implemented via an adapter wrapping modulemgmt.Service in
// main.go (same narrow-interface-plus-adapter pattern payroll already uses
// for its ApprovalEngine dependency) — approval must not import modulemgmt
// directly, and modulemgmt lives in the platform DB, not the tenant DB this
// package otherwise talks to.
type ModuleSubscriptionChecker interface {
	IsModuleActive(companyID, moduleSlug string) (bool, error)
	ListActiveModules(companyID string) ([]string, error)
}

// StatusChangeHandler is invoked when an approval instance belonging to a
// given module reaches a final state (APPROVED/REJECTED/CANCELLED), so the
// owning module can update its own record immediately (push) instead of a
// human/caller having to poll or manually call back with the new status.
type StatusChangeHandler func(ctx context.Context, documentID uuid.UUID, status InstanceStatus, note string) error

// Service untuk business logic Approval Engine.
type Service struct {
	repo           *Repository
	logger         *zap.Logger
	moduleChecker  ModuleSubscriptionChecker
	statusHandlers map[string]StatusChangeHandler
}

// NewService membuat Service baru.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:           repo,
		logger:         logger,
		statusHandlers: make(map[string]StatusChangeHandler),
	}
}

// RegisterStatusHandler mendaftarkan callback yang dipanggil setiap kali
// instance milik `module` mencapai status final. Memanggil ulang untuk module
// yang sama menimpa handler sebelumnya (satu handler per module).
func (s *Service) RegisterStatusHandler(module string, handler StatusChangeHandler) {
	s.statusHandlers[module] = handler
}

// notifyStatusChange memanggil status handler yang terdaftar untuk module
// milik instance (jika ada). Kegagalan handler di-log, tidak mem-fail-kan
// operasi approval itu sendiri — status di approval sudah final dan
// tersimpan; kegagalan sinkronisasi ke module konsumen adalah masalah
// terpisah yang tidak boleh membuat approval terlihat gagal ke user.
func (s *Service) notifyStatusChange(ctx context.Context, instance *ApprovalInstance, note string) {
	handler, ok := s.statusHandlers[instance.Module]
	if !ok {
		return
	}
	if err := handler(ctx, instance.DocumentID, instance.Status, note); err != nil {
		s.logger.Error("status change handler failed",
			zap.String("module", instance.Module),
			zap.String("document_id", instance.DocumentID.String()),
			zap.String("status", string(instance.Status)),
			zap.Error(err),
		)
	}
}

// SetModuleChecker menyuntikkan module-subscription checker (opsional). Jika
// tidak pernah dipanggil, validasi subscription pada CreateFlow/UpdateFlow
// dilewati (mis. pada test) — hanya divalidasi ketika benar-benar diwire di
// main.go.
func (s *Service) SetModuleChecker(mc ModuleSubscriptionChecker) {
	s.moduleChecker = mc
}

// ListAvailableModules mengembalikan slug module yang aktif untuk tenant saat
// ini (dari context) — dipakai frontend flow builder agar module picker hanya
// menampilkan module yang benar-benar disubscribe, bukan free-text.
func (s *Service) ListAvailableModules(ctx context.Context) ([]string, error) {
	if s.moduleChecker == nil {
		return nil, fmt.Errorf("module subscription checker is not configured")
	}
	companyID := authctx.GetCompanyID(ctx)
	if companyID == "" {
		return nil, fmt.Errorf("tenant context not found")
	}
	return s.moduleChecker.ListActiveModules(companyID)
}

// ensureModuleSubscribed menolak operasi jika module belum disubscribe tenant.
// No-op jika moduleChecker belum diwire (backward compatible / test-friendly).
func (s *Service) ensureModuleSubscribed(ctx context.Context, module string) error {
	if s.moduleChecker == nil {
		return nil
	}
	companyID := authctx.GetCompanyID(ctx)
	if companyID == "" {
		return fmt.Errorf("tenant context not found")
	}
	active, err := s.moduleChecker.IsModuleActive(companyID, module)
	if err != nil {
		return fmt.Errorf("failed to check module subscription: %w", err)
	}
	if !active {
		return fmt.Errorf("module %q is not subscribed for this tenant", module)
	}
	return nil
}

// =========================================================================
// Approval Flows
// =========================================================================

// CreateFlow membuat alur persetujuan baru.
func (s *Service) CreateFlow(ctx context.Context, req CreateFlowRequest) (*FlowResponse, error) {
	if err := s.ensureModuleSubscribed(ctx, req.Module); err != nil {
		return nil, err
	}

	flow := &ApprovalFlow{
		Module:   req.Module,
		Name:     req.Name,
		Version:  req.Version,
		IsActive: true,
	}
	if flow.Version == 0 {
		flow.Version = 1
	}
	if err := s.repo.CreateFlow(ctx, flow); err != nil {
		return nil, err
	}

	s.logger.Info("Approval flow created",
		zap.String("flow_id", flow.ID.String()),
		zap.String("module", flow.Module),
		zap.String("name", flow.Name),
	)

	response := flow.ToResponse()
	return &response, nil
}

// GetFlowByID mengembalikan alur persetujuan berdasarkan ID.
func (s *Service) GetFlowByID(ctx context.Context, id string) (*FlowResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid flow id: %w", err)
	}

	flow, err := s.repo.FindFlowByIDWithSteps(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := flow.ToResponse()
	return &response, nil
}

// ListFlows mengembalikan daftar semua alur persetujuan dengan pagination.
func (s *Service) ListFlows(ctx context.Context, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	flows, total, err := s.repo.ListFlows(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []FlowResponse
	for _, f := range flows {
		responses = append(responses, f.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateFlow mengupdate alur persetujuan.
func (s *Service) UpdateFlow(ctx context.Context, id string, req UpdateFlowRequest) (*FlowResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid flow id: %w", err)
	}

	flow, err := s.repo.FindFlowByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if req.Module != nil {
		flow.Module = *req.Module
	}
	if req.Name != nil {
		flow.Name = *req.Name
	}
	if req.Version != nil {
		flow.Version = *req.Version
	}
	if req.IsActive != nil {
		flow.IsActive = *req.IsActive
	}

	// Re-validate subscription whenever the module changes or the flow is
	// (re)activated — an already-inactive flow being edited otherwise (e.g.
	// renamed) doesn't need to re-check.
	if req.Module != nil || (req.IsActive != nil && *req.IsActive) {
		if err := s.ensureModuleSubscribed(ctx, flow.Module); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateFlow(ctx, flow); err != nil {
		return nil, err
	}

	response := flow.ToResponse()
	return &response, nil
}

// DeleteFlow menghapus alur persetujuan (soft delete).
func (s *Service) DeleteFlow(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid flow id: %w", err)
	}

	return s.repo.SoftDeleteFlow(ctx, uid)
}

// =========================================================================
// Approval Flow Steps
// =========================================================================

// CreateStep membuat langkah baru dalam alur persetujuan.
func (s *Service) CreateStep(ctx context.Context, flowID string, req CreateStepRequest) (*FlowStepResponse, error) {
	flowUUID, err := uuid.Parse(flowID)
	if err != nil {
		return nil, fmt.Errorf("invalid flow id: %w", err)
	}

	// Verify flow exists
	if _, err := s.repo.FindFlowByID(ctx, flowUUID); err != nil {
		return nil, err
	}

	// Get max step order
	maxOrder, err := s.repo.GetMaxStepOrder(ctx, flowUUID)
	if err != nil {
		return nil, err
	}

	step := &ApprovalFlowStep{
		FlowID:            flowUUID,
		StepOrder:         maxOrder + 1,
		StepName:          req.StepName,
		ApproverType:      ApproverType(req.ApproverType),
		ApprovalMode:      ApprovalModeAnyOne,
		ParticipationType: ParticipationTypeApprover,
		AllowReject:       true,
	}

	if req.HierarchyLevel != nil {
		step.HierarchyLevel = req.HierarchyLevel
	}
	if req.ApprovalMode != "" {
		step.ApprovalMode = ApprovalMode(req.ApprovalMode)
	}
	if req.ParticipationType != "" {
		step.ParticipationType = ParticipationType(req.ParticipationType)
	}
	if req.RoleID != nil {
		if uid, err := uuid.Parse(*req.RoleID); err == nil {
			step.RoleID = &uid
		}
	}
	if req.ApproverUserID != nil {
		if uid, err := uuid.Parse(*req.ApproverUserID); err == nil {
			step.ApproverUserID = &uid
		}
	}
	if req.RequiredApprovals != nil {
		step.RequiredApprovals = req.RequiredApprovals
	}
	if req.AllowReject != nil {
		step.AllowReject = *req.AllowReject
	}
	if req.ConditionsJSON != nil {
		step.ConditionsJSON = req.ConditionsJSON
	}
	if req.SLAHours != nil {
		step.SLAHours = req.SLAHours
	}

	if err := s.repo.CreateStep(ctx, step); err != nil {
		return nil, err
	}

	if orgIDs, err := parseUUIDs(req.OrganizationIDs); err != nil {
		return nil, fmt.Errorf("invalid organization_ids: %w", err)
	} else if len(orgIDs) > 0 {
		if err := s.repo.ReplaceStepOrganizations(ctx, step.ID, orgIDs); err != nil {
			return nil, err
		}
		step.Organizations = make([]ApprovalFlowStepOrganization, 0, len(orgIDs))
		for _, id := range orgIDs {
			step.Organizations = append(step.Organizations, ApprovalFlowStepOrganization{StepID: step.ID, OrganizationID: id})
		}
	}

	response := step.ToResponse()
	return &response, nil
}

// ListStepsByFlow mengembalikan daftar langkah dalam suatu alur.
func (s *Service) ListStepsByFlow(ctx context.Context, flowID string) ([]FlowStepResponse, error) {
	flowUUID, err := uuid.Parse(flowID)
	if err != nil {
		return nil, fmt.Errorf("invalid flow id: %w", err)
	}

	steps, err := s.repo.ListStepsByFlowID(ctx, flowUUID)
	if err != nil {
		return nil, err
	}

	var responses []FlowStepResponse
	for _, step := range steps {
		responses = append(responses, step.ToResponse())
	}

	return responses, nil
}

// UpdateStep mengupdate langkah dalam alur persetujuan.
func (s *Service) UpdateStep(ctx context.Context, flowID, stepID string, req UpdateStepRequest) (*FlowStepResponse, error) {
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return nil, fmt.Errorf("invalid step id: %w", err)
	}

	step, err := s.repo.FindStepByID(ctx, stepUUID)
	if err != nil {
		return nil, err
	}

	if req.StepName != nil {
		step.StepName = *req.StepName
	}
	if req.ApproverType != nil {
		step.ApproverType = ApproverType(*req.ApproverType)
	}
	if req.HierarchyLevel != nil {
		step.HierarchyLevel = req.HierarchyLevel
	}
	if req.ApprovalMode != nil {
		step.ApprovalMode = ApprovalMode(*req.ApprovalMode)
	}
	if req.ParticipationType != nil {
		step.ParticipationType = ParticipationType(*req.ParticipationType)
	}
	if req.RoleID != nil {
		if uid, err := uuid.Parse(*req.RoleID); err == nil {
			step.RoleID = &uid
		}
	}
	if req.ApproverUserID != nil {
		if uid, err := uuid.Parse(*req.ApproverUserID); err == nil {
			step.ApproverUserID = &uid
		}
	}
	if req.RequiredApprovals != nil {
		step.RequiredApprovals = req.RequiredApprovals
	}
	if req.AllowReject != nil {
		step.AllowReject = *req.AllowReject
	}
	if req.ConditionsJSON != nil {
		step.ConditionsJSON = req.ConditionsJSON
	}
	if req.SLAHours != nil {
		step.SLAHours = req.SLAHours
	}

	if err := s.repo.UpdateStep(ctx, step); err != nil {
		return nil, err
	}

	if req.OrganizationIDs != nil {
		orgIDs, err := parseUUIDs(req.OrganizationIDs)
		if err != nil {
			return nil, fmt.Errorf("invalid organization_ids: %w", err)
		}
		if err := s.repo.ReplaceStepOrganizations(ctx, step.ID, orgIDs); err != nil {
			return nil, err
		}
		step.Organizations = make([]ApprovalFlowStepOrganization, 0, len(orgIDs))
		for _, id := range orgIDs {
			step.Organizations = append(step.Organizations, ApprovalFlowStepOrganization{StepID: step.ID, OrganizationID: id})
		}
	} else if orgs, err := s.repo.ListStepOrganizations(ctx, step.ID); err == nil {
		step.Organizations = orgs
	}

	response := step.ToResponse()
	return &response, nil
}

// DeleteStep menghapus langkah dari alur persetujuan.
func (s *Service) DeleteStep(ctx context.Context, flowID, stepID string) error {
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return fmt.Errorf("invalid step id: %w", err)
	}

	return s.repo.SoftDeleteStep(ctx, stepUUID)
}

// =========================================================================
// Approval Instances
// =========================================================================

// CreateInstance membuat instance persetujuan baru.
func (s *Service) CreateInstance(ctx context.Context, req CreateInstanceRequest) (*InstanceResponse, error) {
	flowUUID, err := uuid.Parse(req.FlowID)
	if err != nil {
		return nil, fmt.Errorf("invalid flow id: %w", err)
	}

	docUUID, err := uuid.Parse(req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("invalid document id: %w", err)
	}

	// Load flow with steps
	flow, err := s.repo.FindFlowByIDWithSteps(ctx, flowUUID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	if !flow.IsActive {
		return nil, fmt.Errorf("flow is not active")
	}

	if len(flow.Steps) == 0 {
		return nil, fmt.Errorf("flow has no steps configured")
	}

	// Check if document already has an active instance
	existing, err := s.repo.FindInstanceByDocument(ctx, req.Module, docUUID)
	if err == nil && existing != nil && existing.Status == InstanceStatusPending {
		return nil, fmt.Errorf("document already has a pending approval instance")
	}

	instance := &ApprovalInstance{
		Module:      req.Module,
		DocumentID:  docUUID,
		FlowID:      flowUUID,
		Status:      InstanceStatusPending,
		CurrentStep: flow.Steps[0].StepOrder,
	}
	instance.CreatedBy = authctx.GetUserID(ctx)

	// Walk forward from the first step: WATCHER steps get their tasks
	// auto-completed (informational only) and don't gate progression, so we
	// land on the first real APPROVER step (or the end of the flow).
	var tasks []ApprovalTask
	landedStep, err := s.advanceThroughWatcherSteps(ctx, instance, flow.Steps, 0, &tasks)
	if err != nil {
		return nil, err
	}
	if landedStep != nil {
		instance.CurrentStep = landedStep.StepOrder
	} else {
		// Entire flow was WATCHER-only — nothing left to gate on.
		instance.Status = InstanceStatusApproved
	}

	if err := s.repo.CreateInstanceWithTasks(ctx, instance, tasks); err != nil {
		return nil, err
	}

	s.logger.Info("Approval instance created",
		zap.String("instance_id", instance.ID.String()),
		zap.String("module", req.Module),
		zap.String("document_id", req.DocumentID),
		zap.String("flow_id", req.FlowID),
	)

	if instance.Status == InstanceStatusApproved {
		s.notifyStatusChange(ctx, instance, "")
	}

	// Load and return full instance
	return s.GetInstanceByID(ctx, instance.ID.String())
}

// GetInstanceByID mengembalikan instance persetujuan dengan relasi.
func (s *Service) GetInstanceByID(ctx context.Context, id string) (*InstanceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid instance id: %w", err)
	}

	instance, err := s.repo.FindInstanceByIDWithRelations(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := instance.ToResponse()
	return &response, nil
}

// ListInstances mengembalikan daftar instance dengan filter dan pagination.
func (s *Service) ListInstances(ctx context.Context, page, perPage int, module, status string) (*PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	instances, total, err := s.repo.ListInstances(ctx, page, perPage, module, status)
	if err != nil {
		return nil, err
	}

	var responses []InstanceResponse
	for _, inst := range instances {
		responses = append(responses, inst.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// CancelInstance membatalkan instance persetujuan.
func (s *Service) CancelInstance(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid instance id: %w", err)
	}

	instance, err := s.repo.FindInstanceByID(ctx, uid)
	if err != nil {
		return err
	}

	if err := s.repo.CancelInstance(ctx, uid); err != nil {
		return err
	}

	instance.Status = InstanceStatusCancelled
	s.notifyStatusChange(ctx, instance, "")

	return nil
}

// =========================================================================
// Approval Actions
// =========================================================================

// SubmitAction mengirim aksi (APPROVE/REJECT) pada instance persetujuan.
func (s *Service) SubmitAction(ctx context.Context, instanceID string, userID string, req SubmitActionRequest) (*InstanceResponse, error) {
	instUUID, err := uuid.Parse(instanceID)
	if err != nil {
		return nil, fmt.Errorf("invalid instance id: %w", err)
	}

	actorUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	// Load instance with flow, steps, and tasks
	instance, err := s.repo.FindInstanceByIDWithRelations(ctx, instUUID)
	if err != nil {
		return nil, err
	}

	if instance.Status != InstanceStatusPending {
		return nil, fmt.Errorf("instance is not in PENDING status")
	}

	// Get current step
	currentStep := instance.CurrentStep
	var currentStepDef *ApprovalFlowStep
	for _, step := range instance.Steps {
		if step.StepOrder == currentStep {
			currentStepDef = &step
			break
		}
	}
	if currentStepDef == nil {
		return nil, fmt.Errorf("current step configuration not found")
	}

	// Verify actor has a pending task for this step
	hasTask := false
	for _, task := range instance.Tasks {
		if task.StepOrder == currentStep &&
			task.Status == TaskStatusPending &&
			task.AssigneeID == actorUUID {
			hasTask = true
			break
		}
	}
	if !hasTask {
		return nil, fmt.Errorf("user does not have a pending task for the current step")
	}

	// BUG FIX: Validate reject-permission BEFORE creating action
	// If action is REJECT but step does not allow rejection, reject before any side effects
	if req.Action == "REJECT" && !currentStepDef.AllowReject {
		return nil, fmt.Errorf("rejection is not allowed for this approval step")
	}

	// Create action
	action := &ApprovalAction{
		InstanceID:  instUUID,
		StepOrder:   currentStep,
		ActorUserID: actorUUID,
		Action:      ApprovalActionType(req.Action),
		Note:        req.Note,
	}

	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}

	// Mark corresponding task as done
	for _, task := range instance.Tasks {
		if task.StepOrder == currentStep && task.AssigneeID == actorUUID && task.Status == TaskStatusPending {
			if err := s.repo.UpdateTaskStatus(ctx, task.ID, TaskStatusDone); err != nil {
				return nil, err
			}
			break
		}
	}

	if req.Action == "REJECT" {
		// Reject the entire instance
		if err := s.repo.RejectInstance(ctx, instUUID); err != nil {
			return nil, err
		}

		s.logger.Info("Approval instance rejected",
			zap.String("instance_id", instanceID),
			zap.String("user_id", userID),
		)

		instance.Status = InstanceStatusRejected
		note := ""
		if req.Note != nil {
			note = *req.Note
		}
		s.notifyStatusChange(ctx, instance, note)

		return s.GetInstanceByID(ctx, instanceID)
	}

	// ===== APPROVE flow: check if step conditions are met =====
	tasks, err := s.repo.FindTasksByInstanceID(ctx, instUUID)
	if err != nil {
		return nil, err
	}

	// Count done vs pending for current step
	var doneCount, pendingCount int
	for _, task := range tasks {
		if task.StepOrder == currentStep {
			if task.Status == TaskStatusDone {
				doneCount++
			} else if task.Status == TaskStatusPending {
				pendingCount++
			}
		}
	}

	// Determine if we should proceed based on approval mode
	canProceed := false
	switch currentStepDef.ApprovalMode {
	case ApprovalModeAnyOne:
		// BUG FIX: ANY_ONE — advance as soon as at least ONE approval is done
		canProceed = doneCount >= 1
	case ApprovalModeAll:
		// ALL — advance when all tasks are done (no pending remaining)
		canProceed = pendingCount == 0
	case ApprovalModeNOfM:
		// N_OF_M — advance when doneCount >= requiredApprovals
		if currentStepDef.RequiredApprovals != nil && *currentStepDef.RequiredApprovals > 0 {
			canProceed = doneCount >= *currentStepDef.RequiredApprovals
		} else {
			canProceed = doneCount >= 1
		}
	}

	if canProceed {
		// Find current step's index, then walk forward — skipping WATCHER-only
		// steps (auto-completed, non-blocking) — to the next real gate.
		curIdx := -1
		for i, step := range instance.Steps {
			if step.StepOrder == currentStep {
				curIdx = i
				break
			}
		}

		var nextTasks []ApprovalTask
		var landedStep *ApprovalFlowStep
		if curIdx >= 0 {
			landedStep, err = s.advanceThroughWatcherSteps(ctx, instance, instance.Steps, curIdx+1, &nextTasks)
			if err != nil {
				return nil, err
			}
		}

		if landedStep != nil {
			if err := s.repo.ApproveStep(ctx, instUUID, currentStep, landedStep.StepOrder, nextTasks); err != nil {
				return nil, err
			}
		} else {
			// No more real steps — approve instance fully. Any trailing
			// WATCHER tasks are still persisted for visibility/audit.
			if err := s.repo.ApproveStep(ctx, instUUID, currentStep, currentStep, nextTasks); err != nil {
				return nil, err
			}

			instance.Status = InstanceStatusApproved
			s.notifyStatusChange(ctx, instance, "")
		}
	}

	s.logger.Info("Approval action submitted",
		zap.String("instance_id", instanceID),
		zap.String("user_id", userID),
		zap.String("action", req.Action),
	)

	return s.GetInstanceByID(ctx, instanceID)
}

// =========================================================================
// Approval Tasks (for approver)
// =========================================================================

// ListMyPendingTasks mengembalikan daftar task approval yang pending untuk user.
func (s *Service) ListMyPendingTasks(ctx context.Context, userID string, page, perPage int) (*PaginatedResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	tasks, total, err := s.repo.FindPendingTasksByAssignee(ctx, "USER", uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []TaskResponse
	for _, t := range tasks {
		responses = append(responses, t.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// =========================================================================
// Helper: Resolve step assignees
// =========================================================================

type taskAssignee struct {
	Type string
	ID   uuid.UUID
}

// resolveStepAssignees resolve siapa saja yang harus (atau hanya perlu tahu,
// untuk step WATCHER) menerima task pada sebuah step. USER/ROLE assignee type
// tetap seperti semula; SUPERVISOR/ORGANIZATION/BOTH di-resolve berdasarkan
// hierarki Organization submitter dan/atau Organization eksplisit milik step.
func (s *Service) resolveStepAssignees(ctx context.Context, instance *ApprovalInstance, step ApprovalFlowStep) ([]taskAssignee, error) {
	var assignees []taskAssignee

	switch step.ApproverType {
	case ApproverTypeUser:
		if step.ApproverUserID != nil {
			assignees = append(assignees, taskAssignee{
				Type: "USER",
				ID:   *step.ApproverUserID,
			})
		}
	case ApproverTypeRole:
		if step.RoleID != nil {
			assignees = append(assignees, taskAssignee{
				Type: "ROLE",
				ID:   *step.RoleID,
			})
		}
	case ApproverTypeSupervisor:
		supAssignees, err := s.resolveSupervisorAssignees(ctx, instance, step)
		if err != nil {
			return nil, err
		}
		assignees = append(assignees, supAssignees...)
	case ApproverTypeOrganization:
		orgAssignees, err := s.resolveOrganizationAssignees(ctx, step)
		if err != nil {
			return nil, err
		}
		assignees = append(assignees, orgAssignees...)
	case ApproverTypeBoth:
		supAssignees, err := s.resolveSupervisorAssignees(ctx, instance, step)
		if err != nil {
			return nil, err
		}
		orgAssignees, err := s.resolveOrganizationAssignees(ctx, step)
		if err != nil {
			return nil, err
		}
		assignees = dedupAssignees(append(supAssignees, orgAssignees...))
	}

	if len(assignees) == 0 {
		return nil, fmt.Errorf("no assignees resolved for step: %s", step.StepName)
	}

	return assignees, nil
}

// resolveSupervisorAssignees walks up the submitter's Organization tree
// `step.HierarchyLevel` times (default 1 = atasan langsung) and resolves to
// the platform user(s) currently occupying that ancestor Organization.
func (s *Service) resolveSupervisorAssignees(ctx context.Context, instance *ApprovalInstance, step ApprovalFlowStep) ([]taskAssignee, error) {
	if instance.CreatedBy == nil {
		return nil, fmt.Errorf("cannot resolve supervisor: instance has no submitter")
	}

	orgID, err := s.repo.GetSubmitterOrganizationID(ctx, *instance.CreatedBy)
	if err != nil {
		return nil, err
	}
	if orgID == nil {
		return nil, fmt.Errorf("submitter's organization could not be resolved")
	}

	level := 1
	if step.HierarchyLevel != nil && *step.HierarchyLevel > 0 {
		level = *step.HierarchyLevel
	}

	target := *orgID
	for i := 0; i < level; i++ {
		parent, err := s.repo.GetOrganizationParentID(ctx, target)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("no organization %d level(s) above the submitter", level)
		}
		target = *parent
	}

	userIDs, err := s.repo.GetUserIDsByOrganization(ctx, target)
	if err != nil {
		return nil, err
	}
	assignees := make([]taskAssignee, 0, len(userIDs))
	for _, uid := range userIDs {
		assignees = append(assignees, taskAssignee{Type: "USER", ID: uid})
	}
	return assignees, nil
}

// resolveOrganizationAssignees resolves to the platform user(s) currently
// occupying each Organization explicitly listed on the step.
func (s *Service) resolveOrganizationAssignees(ctx context.Context, step ApprovalFlowStep) ([]taskAssignee, error) {
	var assignees []taskAssignee
	for _, o := range step.Organizations {
		userIDs, err := s.repo.GetUserIDsByOrganization(ctx, o.OrganizationID)
		if err != nil {
			return nil, err
		}
		for _, uid := range userIDs {
			assignees = append(assignees, taskAssignee{Type: "USER", ID: uid})
		}
	}
	return assignees, nil
}

func dedupAssignees(in []taskAssignee) []taskAssignee {
	seen := make(map[string]struct{}, len(in))
	out := make([]taskAssignee, 0, len(in))
	for _, a := range in {
		key := a.Type + ":" + a.ID.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, a)
	}
	return out
}

// advanceThroughWatcherSteps walks forward through steps[fromIndex:], resolving
// assignees and appending tasks for each to tasksOut. WATCHER steps have their
// tasks marked DONE immediately (informational only, they never block
// progression) and the walk continues past them. Returns the first step it
// lands on that actually gates progress (participation_type=APPROVER), or nil
// if every remaining step was WATCHER-only.
func (s *Service) advanceThroughWatcherSteps(ctx context.Context, instance *ApprovalInstance, steps []ApprovalFlowStep, fromIndex int, tasksOut *[]ApprovalTask) (*ApprovalFlowStep, error) {
	for i := fromIndex; i < len(steps); i++ {
		step := steps[i]

		assignees, err := s.resolveStepAssignees(ctx, instance, step)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve assignees for step %q: %w", step.StepName, err)
		}

		status := TaskStatusPending
		if step.ParticipationType == ParticipationTypeWatcher {
			status = TaskStatusDone
		}
		for _, a := range assignees {
			*tasksOut = append(*tasksOut, ApprovalTask{
				StepOrder:    step.StepOrder,
				AssigneeType: a.Type,
				AssigneeID:   a.ID,
				Status:       status,
			})
		}

		if step.ParticipationType != ParticipationTypeWatcher {
			landed := step
			return &landed, nil
		}
	}
	return nil, nil
}

func parseUUIDs(in []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(in))
	for _, str := range in {
		id, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}
