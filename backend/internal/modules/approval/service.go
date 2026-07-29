package approval

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// Service untuk business logic Approval Engine.
type Service struct {
	repo   *Repository
	logger *zap.Logger
}

// NewService membuat Service baru.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// =========================================================================
// Approval Flows
// =========================================================================

// CreateFlow membuat alur persetujuan baru.
func (s *Service) CreateFlow(ctx context.Context, req CreateFlowRequest) (*FlowResponse, error) {
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
		FlowID:       flowUUID,
		StepOrder:    maxOrder + 1,
		StepName:     req.StepName,
		ApproverType: ApproverType(req.ApproverType),
		ApprovalMode: ApprovalModeAnyOne,
		AllowReject:  true,
	}

	if req.ApprovalMode != "" {
		step.ApprovalMode = ApprovalMode(req.ApprovalMode)
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
	if req.ApprovalMode != nil {
		step.ApprovalMode = ApprovalMode(*req.ApprovalMode)
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
		CurrentStep: 1,
	}
	instance.CreatedBy = authctx.GetUserID(ctx)

	// Create tasks for the first step
	var tasks []ApprovalTask
	firstStep := flow.Steps[0]
	taskAssignees, err := s.resolveStepAssignees(firstStep)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve step assignees: %w", err)
	}

	for _, assignee := range taskAssignees {
		tasks = append(tasks, ApprovalTask{
			StepOrder:    firstStep.StepOrder,
			AssigneeType: assignee.Type,
			AssigneeID:   assignee.ID,
			Status:       TaskStatusPending,
		})
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

	return s.repo.CancelInstance(ctx, uid)
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
		// Find next step
		nextStep := 0
		for _, step := range instance.Steps {
			if step.StepOrder > currentStep {
				nextStep = step.StepOrder
				break
			}
		}

		if nextStep > 0 {
			// Create tasks for next step
			var nextTasks []ApprovalTask
			var nextStepDef ApprovalFlowStep
			for _, step := range instance.Steps {
				if step.StepOrder == nextStep {
					nextStepDef = step
					break
				}
			}

			taskAssignees, err := s.resolveStepAssignees(nextStepDef)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve next step assignees: %w", err)
			}

			for _, assignee := range taskAssignees {
				nextTasks = append(nextTasks, ApprovalTask{
					StepOrder:    nextStep,
					AssigneeType: assignee.Type,
					AssigneeID:   assignee.ID,
					Status:       TaskStatusPending,
				})
			}

			if err := s.repo.ApproveStep(ctx, instUUID, currentStep, nextStep, nextTasks); err != nil {
				return nil, err
			}
		} else {
			// No more steps, approve instance fully
			if err := s.repo.ApproveStep(ctx, instUUID, currentStep, currentStep, nil); err != nil {
				return nil, err
			}
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

func (s *Service) resolveStepAssignees(step ApprovalFlowStep) ([]taskAssignee, error) {
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
		// For supervisor, we need to resolve the employee's supervisor dynamically
		// For now, this is a placeholder; actual resolution will be implemented
		// when integrating with the Employee module
		return nil, fmt.Errorf("SUPERVISOR approver type requires employee module integration")
	}

	if len(assignees) == 0 {
		return nil, fmt.Errorf("no assignees resolved for step: %s", step.StepName)
	}

	return assignees, nil
}
