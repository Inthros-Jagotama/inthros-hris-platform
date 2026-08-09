package employeemovement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
)

// ApprovalEngine abstracts the central approval module so employee
// movements are approved through it (single approval path — manual approve
// dihapus, keputusan plan §11.5 / G-5). Implemented via an adapter wrapping
// approval.Service in main.go (same narrow-interface-plus-adapter pattern
// payroll/leave/reimbursement already use).
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
	// GetActiveFlowIDForModule lets a movement submission auto-resolve which
	// flow to route through when the client doesn't supply flow_id explicitly
	// (same pattern leave/attendance uses) — without this, a movement
	// submitted without a flow_id stays in draft and never reaches the
	// Approval module.
	GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
}

// Notifier abstracts the notification module so employeemovement can notify
// the employee of their movement's approval/execution outcome
// (docs/module-movement-plan.md §7 — same pattern attendance/leave).
// notification.Service satisfies this structurally.
type Notifier interface {
	Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error
}

// CareerEmployment is the data needed to create an employment record from a
// movement's to_* fields. Defined here (instead of importing the employee
// module) so employeemovement only depends on a narrow interface.
//
// ID is only populated by FindCurrentEmployment (so CloseEmployment knows
// which record to close); it is empty when building a new employment.
type CareerEmployment struct {
	ID                   uuid.UUID
	OrganizationID       *uuid.UUID
	PositionID           *uuid.UUID
	EmploymentStatusID   *uuid.UUID
	DecisionLetterNumber string
	DecisionLetterDate   string
	EffectiveDate        string
}

// MovementValidationError is returned by the service when movement fields
// violate the per-type business rules (plan G-7). Handlers map it to a 400
// Bad Request response so the FE gets a field-level message instead of a 500.
type MovementValidationError struct {
	Message string
}

func (e *MovementValidationError) Error() string {
	return e.Message
}

// validateMovementFields enforces per-type required fields (plan G-7):
//   - mutation            → wajib to_organization_id (dan/atau to_position_id)
//   - promotion / demotion → wajib to_position_id
//   - status_change       → wajib to_employment_status_id
//   - contract_extension  → wajib merujuk kontrak aktif (dicek via repo)
//   - offboarding / retirement → boleh tanpa to_* (tanpa validasi)
//
// `hasActiveContract` is only consulted for contract_extension and may be nil
// (skip) when the caller already knows the employee has no contract.
func validateMovementFields(movementType MovementType, toOrganizationID, toPositionID, toEmploymentStatusID *uuid.UUID, hasActiveContract bool) error {
	switch movementType {
	case MovementTypeMutation:
		if toOrganizationID == nil && toPositionID == nil {
			return &MovementValidationError{Message: "movement type 'mutation' requires to_organization_id or to_position_id"}
		}
	case MovementTypePromotion, MovementTypeDemotion:
		if toPositionID == nil {
			return &MovementValidationError{Message: fmt.Sprintf("movement type '%s' requires to_position_id", movementType)}
		}
	case MovementTypeStatusChange:
		if toEmploymentStatusID == nil {
			return &MovementValidationError{Message: "movement type 'status_change' requires to_employment_status_id"}
		}
	case MovementTypeContractExtension:
		if !hasActiveContract {
			return &MovementValidationError{Message: "movement type 'contract_extension' requires an active employee contract"}
		}
	}
	return nil
}

// CareerExecutor abstracts the employee module's employment + employee
// status changes so ExecuteMovement can push the real HR data change
// (create new employment, close the previous one, mark offboarding /
// retirement employees inactive). Implemented via an adapter wrapping
// employee.Service in main.go (same narrow-interface-plus-adapter pattern
// as ApprovalEngine / AttendanceSessionUpdater).
type CareerExecutor interface {
	// FindCurrentEmployment returns the employee's currently active employment
	// (most recent with no effective_end_date), or nil if none.
	FindCurrentEmployment(ctx context.Context, employeeID uuid.UUID) (*CareerEmployment, error)
	// CloseEmployment sets the employment's effective_end_date to the day
	// before effectiveDate (so the new employment can take over).
	CloseEmployment(ctx context.Context, employmentID uuid.UUID, effectiveDate string) error
	// CreateEmployment persists a new employment record and returns its ID.
	CreateEmployment(ctx context.Context, employeeID uuid.UUID, data CareerEmployment) (uuid.UUID, error)
	// SetEmployeeInactive marks an offboarded/retired employee as inactive.
	SetEmployeeInactive(ctx context.Context, employeeID uuid.UUID) error
}

// Service untuk business logic Employee Movement & Career Management.
type Service struct {
	repo           *Repository
	logger         *zap.Logger
	approvalEngine ApprovalEngine
	careerExecutor CareerExecutor
	notifier       Notifier
}

// NewService membuat Service baru.
func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SetCareerExecutor wires the employee module (employment + employee status
// changes) into this service so ExecuteMovement touches real HR data.
func (s *Service) SetCareerExecutor(ce CareerExecutor) {
	s.careerExecutor = ce
}

// SetNotifier wires the notification module into this service so
// HandleApprovalStatusChange and ExecuteMovement can notify the employee of
// the movement outcome (docs/module-movement-plan.md §7).
func (s *Service) SetNotifier(n Notifier) {
	s.notifier = n
}

// notifyMovementOutcome best-effort notification to the movement's employee
// (mirrors attendance.notifyRequestOutcome / leave.notifyLeaveOutcome).
func (s *Service) notifyMovementOutcome(ctx context.Context, employeeID uuid.UUID, notifType, referenceType string, referenceID uuid.UUID) {
	if s.notifier == nil {
		return
	}
	userID, err := s.repo.FindUserIDByEmployeeID(ctx, employeeID)
	if err != nil {
		s.logger.Warn("Failed to resolve employee user id for movement notification",
			zap.String("reference_type", referenceType),
			zap.String("reference_id", referenceID.String()),
			zap.Error(err),
		)
		return
	}
	if userID == nil {
		return
	}
	if err := s.notifier.Notify(ctx, *userID, notifType, nil, referenceType, referenceID); err != nil {
		s.logger.Warn("Failed to send movement notification",
			zap.String("notif_type", notifType),
			zap.String("reference_type", referenceType),
			zap.String("reference_id", referenceID.String()),
			zap.Error(err),
		)
	}
}

// validateMovement checks the per-type required fields before persisting or
// updating a movement. contract_extension also verifies the employee has an
// active contract (plan G-7). It is called both on create and after an update
// (movement_type and to_* fields may have changed).
func (s *Service) validateMovement(ctx context.Context, m *EmployeeMovement) error {
	hasActiveContract := false
	if m.MovementType == MovementTypeContractExtension {
		has, err := s.repo.HasActiveContractByEmployeeID(ctx, m.EmployeeID)
		if err != nil {
			return err
		}
		hasActiveContract = has
	}
	return validateMovementFields(m.MovementType, m.ToOrganizationID, m.ToPositionID, m.ToEmploymentStatusID, hasActiveContract)
}

// =========================================================================
// Employee Movement
// =========================================================================

// CreateMovement membuat pergerakan karyawan baru.
func (s *Service) CreateMovement(ctx context.Context, req CreateMovementRequest) (*MovementResponse, error) {
	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	movement := &EmployeeMovement{
		CreatedBy:            authctx.GetUserID(ctx),
		UpdatedBy:            authctx.GetUserID(ctx),
		EmployeeID:           employeeUUID,
		MovementType:         MovementType(req.MovementType),
		DecisionLetterNumber: req.DecisionLetterNumber,
		DecisionLetterDate:   req.DecisionLetterDate,
		EffectiveDate:        req.EffectiveDate,
		Status:               MovementStatusDraft,
	}

	if req.FromEmploymentID != nil {
		if uid, err := uuid.Parse(*req.FromEmploymentID); err == nil {
			movement.FromEmploymentID = &uid
		}
	}
	if req.ToEmploymentID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentID); err == nil {
			movement.ToEmploymentID = &uid
		}
	}
	if req.FromOrganizationID != nil {
		if uid, err := uuid.Parse(*req.FromOrganizationID); err == nil {
			movement.FromOrganizationID = &uid
		}
	}
	if req.ToOrganizationID != nil {
		if uid, err := uuid.Parse(*req.ToOrganizationID); err == nil {
			movement.ToOrganizationID = &uid
		}
	}
	if req.FromPositionID != nil {
		if uid, err := uuid.Parse(*req.FromPositionID); err == nil {
			movement.FromPositionID = &uid
		}
	}
	if req.ToPositionID != nil {
		if uid, err := uuid.Parse(*req.ToPositionID); err == nil {
			movement.ToPositionID = &uid
		}
	}
	if req.FromEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.FromEmploymentStatusID); err == nil {
			movement.FromEmploymentStatusID = &uid
		}
	}
	if req.ToEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentStatusID); err == nil {
			movement.ToEmploymentStatusID = &uid
		}
	}
	if req.Reason != nil {
		movement.Reason = req.Reason
	}
	if req.Notes != nil {
		movement.Notes = req.Notes
	}

	// Business validation per movement type (plan G-7).
	if err := s.validateMovement(ctx, movement); err != nil {
		return nil, err
	}

	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement created",
		zap.String("employee_id", req.EmployeeID),
		zap.String("movement_type", req.MovementType),
		zap.String("movement_id", movement.ID.String()),
	)

	responses := []MovementResponse{movement.ToResponse()}
	s.enrichMovementResponses(ctx, responses)
	return &responses[0], nil
}

// =========================================================================
// Enrichment helpers (plan G-4) — fill display names on responses via batch
// JOINs so the frontend does not need to resolve UUIDs one-by-one.
// =========================================================================

// collectUUIDStrings parses non-empty string ids into a deduped uuid slice,
// silently skipping values that are not valid UUIDs.
func collectUUIDStrings(ids []string) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	var result []uuid.UUID
	for _, id := range ids {
		if id == "" {
			continue
		}
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		if _, ok := seen[uid]; !ok {
			seen[uid] = struct{}{}
			result = append(result, uid)
		}
	}
	return result
}

// fillEmployeeNames copies resolved employee info (name + code) onto the given
// employee ids, if present in the info map.
func fillEmployeeNames(info map[string]employeeRefInfo, employeeID string, resp *MovementResponse) {
	if info != nil {
		if emp, ok := info[employeeID]; ok {
			resp.EmployeeName = emp.Name
			resp.EmployeeCode = emp.Code
		}
	}
}

// fillContractEmployeeNames copies resolved employee info (name + code) onto
// the given contract response, if present in the info map.
func fillContractEmployeeNames(info map[string]employeeRefInfo, employeeID string, resp *ContractResponse) {
	if info != nil {
		if emp, ok := info[employeeID]; ok {
			resp.EmployeeName = emp.Name
			resp.EmployeeCode = emp.Code
		}
	}
}

// enrichMovementResponses fills employee/organization/position/status display
// names on movement responses (single or list) with batch queries (G-4).
func (s *Service) enrichMovementResponses(ctx context.Context, responses []MovementResponse) {
	if len(responses) == 0 {
		return
	}

	// Collect distinct ids per reference table.
	empIDs := make(map[uuid.UUID]struct{}, len(responses))
	var orgIDs, posIDs, statusIDs []string
	for i := range responses {
		r := &responses[i]
		eid, err := uuid.Parse(r.EmployeeID)
		if err == nil {
			empIDs[eid] = struct{}{}
		}
		if r.FromOrganizationID != nil {
			orgIDs = append(orgIDs, *r.FromOrganizationID)
		}
		if r.ToOrganizationID != nil {
			orgIDs = append(orgIDs, *r.ToOrganizationID)
		}
		if r.FromPositionID != nil {
			posIDs = append(posIDs, *r.FromPositionID)
		}
		if r.ToPositionID != nil {
			posIDs = append(posIDs, *r.ToPositionID)
		}
		if r.FromEmploymentStatusID != nil {
			statusIDs = append(statusIDs, *r.FromEmploymentStatusID)
		}
		if r.ToEmploymentStatusID != nil {
			statusIDs = append(statusIDs, *r.ToEmploymentStatusID)
		}
	}

	empList := make([]uuid.UUID, 0, len(empIDs))
	for id := range empIDs {
		empList = append(empList, id)
	}
	orgList := collectUUIDStrings(orgIDs)
	posList := collectUUIDStrings(posIDs)
	statusList := collectUUIDStrings(statusIDs)

	if empInfo, err := s.repo.GetEmployeeInfoByIDs(ctx, empList); err == nil {
		for i := range responses {
			fillEmployeeNames(empInfo, responses[i].EmployeeID, &responses[i])
		}
	} else {
		s.logger.Warn("failed to resolve employee info for movements", zap.Error(err))
	}

	if names, err := s.repo.GetOrganizationNamesByIDs(ctx, orgList); err == nil {
		for i := range responses {
			if responses[i].FromOrganizationID != nil {
				responses[i].FromOrganizationName = names[*responses[i].FromOrganizationID]
			}
			if responses[i].ToOrganizationID != nil {
				responses[i].ToOrganizationName = names[*responses[i].ToOrganizationID]
			}
		}
	} else {
		s.logger.Warn("failed to resolve organization names for movements", zap.Error(err))
	}

	if names, err := s.repo.GetPositionNamesByIDs(ctx, posList); err == nil {
		for i := range responses {
			if responses[i].FromPositionID != nil {
				responses[i].FromPositionName = names[*responses[i].FromPositionID]
			}
			if responses[i].ToPositionID != nil {
				responses[i].ToPositionName = names[*responses[i].ToPositionID]
			}
		}
	} else {
		s.logger.Warn("failed to resolve position names for movements", zap.Error(err))
	}

	if names, err := s.repo.GetEmploymentStatusNamesByIDs(ctx, statusList); err == nil {
		for i := range responses {
			if responses[i].FromEmploymentStatusID != nil {
				responses[i].FromEmploymentStatusName = names[*responses[i].FromEmploymentStatusID]
			}
			if responses[i].ToEmploymentStatusID != nil {
				responses[i].ToEmploymentStatusName = names[*responses[i].ToEmploymentStatusID]
			}
		}
	} else {
		s.logger.Warn("failed to resolve employment status names for movements", zap.Error(err))
	}
}

// enrichContractResponses fills employee name/code and previous contract number
// on contract responses (single or list) with batch queries (G-4).
func (s *Service) enrichContractResponses(ctx context.Context, responses []ContractResponse) {
	if len(responses) == 0 {
		return
	}
	empIDs := make(map[uuid.UUID]struct{}, len(responses))
	var prevIDs []string
	for i := range responses {
		r := &responses[i]
		eid, err := uuid.Parse(r.EmployeeID)
		if err == nil {
			empIDs[eid] = struct{}{}
		}
		if r.PreviousContractID != nil {
			prevIDs = append(prevIDs, *r.PreviousContractID)
		}
	}

	empList := make([]uuid.UUID, 0, len(empIDs))
	for id := range empIDs {
		empList = append(empList, id)
	}
	prevList := collectUUIDStrings(prevIDs)

	if empInfo, err := s.repo.GetEmployeeInfoByIDs(ctx, empList); err == nil {
		for i := range responses {
			fillContractEmployeeNames(empInfo, responses[i].EmployeeID, &responses[i])
		}
	} else {
		s.logger.Warn("failed to resolve employee info for contracts", zap.Error(err))
	}

	if numbers, err := s.repo.GetContractNumbersByIDs(ctx, prevList); err == nil {
		for i := range responses {
			if responses[i].PreviousContractID != nil {
				responses[i].PreviousContractNumber = numbers[*responses[i].PreviousContractID]
			}
		}
	} else {
		s.logger.Warn("failed to resolve previous contract numbers", zap.Error(err))
	}
}

// GetMovementByID mengembalikan pergerakan berdasarkan ID.
func (s *Service) GetMovementByID(ctx context.Context, id string) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	responses := []MovementResponse{movement.ToResponse()}
	s.enrichMovementResponses(ctx, responses)
	return &responses[0], nil
}

// ListMovementsByEmployee mengembalikan daftar pergerakan untuk seorang karyawan.
func (s *Service) ListMovementsByEmployee(ctx context.Context, employeeID string, page, perPage int) (*PaginatedMovementResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	movements, total, err := s.repo.FindMovementsByEmployeeID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []MovementResponse
	for _, m := range movements {
		responses = append(responses, m.ToResponse())
	}
	s.enrichMovementResponses(ctx, responses)

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListMovements mengembalikan daftar semua pergerakan dengan pagination.
// Optional filters: movementType, status, search (by decision letter number or
// employee name/code) — dipakai halaman FE Movements (langkah 9 plan).
func (s *Service) ListMovements(ctx context.Context, page, perPage int, movementType, status, search string) (*PaginatedMovementResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	movements, total, err := s.repo.ListMovements(ctx, page, perPage, movementType, status, search)
	if err != nil {
		return nil, err
	}

	var responses []MovementResponse
	for _, m := range movements {
		responses = append(responses, m.ToResponse())
	}
	s.enrichMovementResponses(ctx, responses)

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateMovement mengupdate pergerakan.
func (s *Service) UpdateMovement(ctx context.Context, id string, req UpdateMovementRequest) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	movement.UpdatedBy = authctx.GetUserID(ctx)

	if movement.Status != MovementStatusDraft {
		return nil, fmt.Errorf("cannot update movement with status '%s', only draft movements can be updated", movement.Status)
	}

	if req.MovementType != nil {
		movement.MovementType = MovementType(*req.MovementType)
	}
	if req.ToOrganizationID != nil {
		if uid, err := uuid.Parse(*req.ToOrganizationID); err == nil {
			movement.ToOrganizationID = &uid
		}
	}
	if req.ToPositionID != nil {
		if uid, err := uuid.Parse(*req.ToPositionID); err == nil {
			movement.ToPositionID = &uid
		}
	}
	if req.ToEmploymentStatusID != nil {
		if uid, err := uuid.Parse(*req.ToEmploymentStatusID); err == nil {
			movement.ToEmploymentStatusID = &uid
		}
	}
	if req.Reason != nil {
		movement.Reason = req.Reason
	}
	if req.DecisionLetterNumber != nil {
		movement.DecisionLetterNumber = *req.DecisionLetterNumber
	}
	if req.DecisionLetterDate != nil {
		movement.DecisionLetterDate = *req.DecisionLetterDate
	}
	if req.EffectiveDate != nil {
		movement.EffectiveDate = *req.EffectiveDate
	}
	if req.Status != nil {
		movement.Status = MovementStatus(*req.Status)
	}
	if req.Notes != nil {
		movement.Notes = req.Notes
	}

	// Business validation per movement type (plan G-7) — the movement_type
	// may have changed in this request, so re-validate with the effective value.
	if err := s.validateMovement(ctx, movement); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return nil, err
	}

	responses := []MovementResponse{movement.ToResponse()}
	s.enrichMovementResponses(ctx, responses)
	return &responses[0], nil
}

// DeleteMovement menghapus pergerakan (hanya draft).
func (s *Service) DeleteMovement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return err
	}

	if movement.Status != MovementStatusDraft {
		return fmt.Errorf("cannot delete movement with status '%s', only draft movements can be deleted", movement.Status)
	}

	return s.repo.DeleteMovement(ctx, uid)
}

// SubmitMovement routes a draft movement through the central approval
// module — the single approval path (manual approve dihapus, G-5). If the
// client doesn't supply flow_id, the active flow for module
// "employeemovement" is auto-resolved (G-3 — same pattern leave/attendance).
func (s *Service) SubmitMovement(ctx context.Context, id string, req SubmitMovementRequest) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}
	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if movement.Status != MovementStatusDraft {
		return nil, fmt.Errorf("only draft movements can be submitted, current status: %s", movement.Status)
	}
	if s.approvalEngine == nil {
		return nil, fmt.Errorf("approval engine not configured")
	}

	// Auto-resolve the active flow when no flow_id is supplied (G-3).
	flowID := ""
	if req.FlowID != nil && *req.FlowID != "" {
		flowID = *req.FlowID
	} else if resolved, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, "employeemovement"); err == nil {
		flowID = resolved
	}
	if flowID == "" {
		return nil, fmt.Errorf("approval flow not configured: provide flow_id or activate an approval flow for module employeemovement")
	}

	instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "employeemovement", movement.ID.String(), flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to create approval instance: %w", err)
	}
	if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
		movement.ApprovalInstanceID = &parsedInstanceID
	}
	movement.Status = MovementStatusPendingApproval
	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement submitted for approval",
		zap.String("movement_id", movement.ID.String()),
		zap.String("instance_id", instanceID),
	)

	responses := []MovementResponse{movement.ToResponse()}
	s.enrichMovementResponses(ctx, responses)
	return &responses[0], nil
}

// HandleApprovalStatusChange is invoked by the approval module's push-based
// status callback when a movement's approval instance reaches a final
// state, so the movement's own status field updates itself. Since G-5 the
// approval module is the single approval path. REJECTED maps to the
// dedicated "rejected" status (keputusan plan §11.4) — not cancelled.
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	movement, err := s.repo.FindMovementByID(ctx, documentID)
	if err != nil {
		return err
	}
	if movement.Status != MovementStatusPendingApproval {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		movement.Status = MovementStatusApproved
		movement.ApprovedAt = &now
	case "REJECTED":
		movement.Status = MovementStatusRejected
		if note != "" {
			movement.Notes = &note
		}
	case "CANCELLED":
		movement.Status = MovementStatusCancelled
	default:
		return nil
	}

	s.logger.Info("Employee movement status updated via approval status handler",
		zap.String("movement_id", movement.ID.String()),
		zap.String("approval_status", status),
	)
	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return err
	}

	// Best-effort notification to the movement's employee (plan §7).
	switch movement.Status {
	case MovementStatusApproved:
		s.notifyMovementOutcome(ctx, movement.EmployeeID, "MOVEMENT_APPROVED", "employeemovement", movement.ID)
	case MovementStatusRejected, MovementStatusCancelled:
		s.notifyMovementOutcome(ctx, movement.EmployeeID, "MOVEMENT_REJECTED", "employeemovement", movement.ID)
	}
	return nil
}

// movementCreatesEmployment reports whether the movement type should create
// a new employment record when executed. contract_extension only extends the
// contract; offboarding/retirement close the employment without a new one.
func movementCreatesEmployment(t MovementType) bool {
	switch t {
	case MovementTypePromotion, MovementTypeDemotion, MovementTypeMutation, MovementTypeStatusChange, MovementTypeOther:
		return true
	default:
		return false
	}
}

// movementDeactivatesEmployee reports whether the movement marks the
// employee as inactive (offboarding / retirement — keputusan plan §11.3).
func movementDeactivatesEmployee(t MovementType) bool {
	return t == MovementTypeOffboarding || t == MovementTypeRetirement
}

// dayBefore returns the date one day before the given date. Accepts both
// plain YYYY-MM-DD and RFC3339 timestamps (MySQL returns DATETIME values for
// DATE columns), so movement execution is robust regardless of driver.
func dayBefore(date string) (string, error) {
	// Normalize: strip the time portion when an RFC3339 value is stored.
	if len(date) >= 10 {
		if _, err := time.Parse("2006-01-02", date[:10]); err == nil {
			date = date[:10]
		}
	}
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid effective_date %q: %w", date, err)
	}
	return d.AddDate(0, 0, -1).Format("2006-01-02"), nil
}

// ExecuteMovement mengeksekusi pergerakan. Selain mengubah status movement
// menjadi executed, transaksi HR data juga dijalankan (G-1):
//   - promotion/demotion/mutation/status_change/other → buat employment baru
//     (to_* + effective_date), tutup employment aktif lama (effective_end_date
//     = effective_date - 1), simpan to_employment_id di movement.
//   - offboarding/retirement → tutup employment aktif lama dan tandai
//     employee non-aktif (keputusan §11.3).
//   - contract_extension → tanpa perubahan employment.
//
// effective_date boleh di masa depan (keputusan §11.2): employment baru
// disimpan dengan tanggal tsb, employment lama tetap aktif sampai sehari
// sebelumnya.
func (s *Service) ExecuteMovement(ctx context.Context, id string, executedBy string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	executorUUID, err := uuid.Parse(executedBy)
	if err != nil {
		return fmt.Errorf("invalid executor id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return err
	}
	if movement.Status != MovementStatusApproved {
		return fmt.Errorf("movement not found or not in approved status")
	}

	// --- HR data transaction (G-1) ---
	if movementCreatesEmployment(movement.MovementType) || movementDeactivatesEmployee(movement.MovementType) {
		if s.careerExecutor == nil {
			return fmt.Errorf("career executor not configured: cannot execute movement type '%s'", movement.MovementType)
		}

		current, err := s.careerExecutor.FindCurrentEmployment(ctx, movement.EmployeeID)
		if err != nil {
			return fmt.Errorf("failed to find current employment: %w", err)
		}

		// Tutup employment aktif lama (effective_end_date = effective_date - 1).
		if current != nil {
			endDate, err := dayBefore(movement.EffectiveDate)
			if err != nil {
				return err
			}
			if err := s.careerExecutor.CloseEmployment(ctx, current.ID, endDate); err != nil {
				return fmt.Errorf("failed to close previous employment: %w", err)
			}
		}

		if movementCreatesEmployment(movement.MovementType) {
			data := CareerEmployment{
				OrganizationID:       movement.ToOrganizationID,
				PositionID:           movement.ToPositionID,
				EmploymentStatusID:   movement.ToEmploymentStatusID,
				DecisionLetterNumber: movement.DecisionLetterNumber,
				DecisionLetterDate:   movement.DecisionLetterDate,
				EffectiveDate:        movement.EffectiveDate,
			}
			newEmploymentID, err := s.careerExecutor.CreateEmployment(ctx, movement.EmployeeID, data)
			if err != nil {
				return fmt.Errorf("failed to create new employment: %w", err)
			}
			movement.ToEmploymentID = &newEmploymentID
		}

		if movementDeactivatesEmployee(movement.MovementType) {
			if err := s.careerExecutor.SetEmployeeInactive(ctx, movement.EmployeeID); err != nil {
				return fmt.Errorf("failed to mark employee inactive: %w", err)
			}
		}
	}

	if err := s.repo.ExecuteMovement(ctx, uid, executorUUID, movement.ToEmploymentID); err != nil {
		return err
	}
	s.notifyMovementOutcome(ctx, movement.EmployeeID, "MOVEMENT_EXECUTED", "employeemovement", movement.ID)
	return nil
}

// CancelMovement membatalkan pergerakan.
func (s *Service) CancelMovement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid movement id: %w", err)
	}

	return s.repo.CancelMovement(ctx, uid)
}

// =========================================================================
// Employee Contract
// =========================================================================

// CreateContract membuat kontrak karyawan baru.
func (s *Service) CreateContract(ctx context.Context, req CreateContractRequest) (*ContractResponse, error) {
	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	contract := &EmployeeContract{
		CreatedBy:      authctx.GetUserID(ctx),
		UpdatedBy:      authctx.GetUserID(ctx),
		EmployeeID:     employeeUUID,
		ContractNumber: req.ContractNumber,
		ContractType:   ContractType(req.ContractType),
		StartDate:      req.StartDate,
		Status:         ContractStatusActive,
	}

	if req.EndDate != nil {
		contract.EndDate = req.EndDate
	}
	if req.PreviousContractID != nil {
		if uid, err := uuid.Parse(*req.PreviousContractID); err == nil {
			contract.PreviousContractID = &uid
		}
	}
	if req.DecisionLetterNumber != nil {
		contract.DecisionLetterNumber = req.DecisionLetterNumber
	}
	if req.Notes != nil {
		contract.Notes = req.Notes
	}
	if req.DocumentURL != nil {
		contract.DocumentURL = req.DocumentURL
	}

	// Jika ada previous_contract_id, gunakan ExtendContract flow — extension
	// count dihitung berantai dari kontrak sebelumnya (G-6).
	if contract.PreviousContractID != nil {
		if err := s.repo.ExtendContract(ctx, contract, *contract.PreviousContractID); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.CreateContract(ctx, contract); err != nil {
			return nil, err
		}
	}

	s.logger.Info("Employee contract created",
		zap.String("employee_id", req.EmployeeID),
		zap.String("contract_number", req.ContractNumber),
		zap.String("contract_type", req.ContractType),
	)

	responses := []ContractResponse{contract.ToResponse()}
	s.enrichContractResponses(ctx, responses)
	return &responses[0], nil
}

// GetContractByID mengembalikan kontrak berdasarkan ID.
func (s *Service) GetContractByID(ctx context.Context, id string) (*ContractResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}

	contract, err := s.repo.FindContractByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	responses := []ContractResponse{contract.ToResponse()}
	s.enrichContractResponses(ctx, responses)
	return &responses[0], nil
}

// ListContractsByEmployee mengembalikan daftar kontrak untuk seorang karyawan.
func (s *Service) ListContractsByEmployee(ctx context.Context, employeeID string, page, perPage int) (*PaginatedContractResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	contracts, total, err := s.repo.FindContractsByEmployeeID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []ContractResponse
	for _, c := range contracts {
		responses = append(responses, c.ToResponse())
	}
	s.enrichContractResponses(ctx, responses)

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedContractResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListContracts mengembalikan daftar semua kontrak dengan pagination.
func (s *Service) ListContracts(ctx context.Context, page, perPage int) (*PaginatedContractResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	contracts, total, err := s.repo.ListContracts(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []ContractResponse
	for _, c := range contracts {
		responses = append(responses, c.ToResponse())
	}
	s.enrichContractResponses(ctx, responses)

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedContractResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateContract mengupdate kontrak.
func (s *Service) UpdateContract(ctx context.Context, id string, req UpdateContractRequest) (*ContractResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}

	contract, err := s.repo.FindContractByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	contract.UpdatedBy = authctx.GetUserID(ctx)

	if req.ContractNumber != nil {
		contract.ContractNumber = *req.ContractNumber
	}
	if req.ContractType != nil {
		contract.ContractType = ContractType(*req.ContractType)
	}
	if req.EndDate != nil {
		contract.EndDate = req.EndDate
	}
	if req.DecisionLetterNumber != nil {
		contract.DecisionLetterNumber = req.DecisionLetterNumber
	}
	if req.Notes != nil {
		contract.Notes = req.Notes
	}
	if req.DocumentURL != nil {
		contract.DocumentURL = req.DocumentURL
	}
	if req.Status != nil {
		contract.Status = ContractStatus(*req.Status)
	}

	if err := s.repo.UpdateContract(ctx, contract); err != nil {
		return nil, err
	}

	responses := []ContractResponse{contract.ToResponse()}
	s.enrichContractResponses(ctx, responses)
	return &responses[0], nil
}

// DeleteContract menghapus kontrak.
func (s *Service) DeleteContract(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid contract id: %w", err)
	}

	return s.repo.DeleteContract(ctx, uid)
}
