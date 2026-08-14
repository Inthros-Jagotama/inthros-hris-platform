package employeemovement

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/documenttemplate"
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

// PerformanceProvider abstracts the performance module so employeemovement
// can read the latest completed evaluation's final score as eligibility input
// (plan §12.10: Movement hanya membaca hasil final, tidak menghitung KPI/OKR).
// Implemented via adapter wrapping performance.Service in main.go.
type PerformanceProvider interface {
	LatestFinalScore(ctx context.Context, employeeID uuid.UUID) (float64, bool, error)
}

// CompetencyProvider abstracts the competency module so employeemovement can
// read the latest assessment score as eligibility input (plan §12.10/§12.11:
// Movement hanya membaca hasil final dari competency management).
// Implemented via adapter wrapping competency.Repository in main.go.
type CompetencyProvider interface {
	LatestScore(ctx context.Context, employeeID uuid.UUID) (float64, bool, error)
}

// OKRProvider abstracts the performance (OKR) module so employeemovement can
// read the latest completed OKR evaluation's final score as eligibility input
// (plan §12.10/§12.11 — KPI, OKR, dan Competency menjadi input eligibility;
// movement tidak menghitung skor sendiri). Implemented via adapter wrapping
// performance.OKRRepository in main.go.
type OKRProvider interface {
	LatestScore(ctx context.Context, employeeID uuid.UUID) (float64, bool, error)
}

// MovementConflictError is returned when a movement cannot be created or
// executed because it conflicts with the current employment state — the
// target position is already occupied by another employee (plan §12.3) or the
// effective date overlaps an existing employment period of the same employee
// (plan §12.4). Handlers map it to a 409 Conflict response so the FE shows a
// meaningful message instead of a generic 500.
type MovementConflictError struct {
	Message string
}

func (e *MovementConflictError) Error() string {
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
	FindCurrentEmployment(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID) (*CareerEmployment, error)
	// CloseEmployment sets the employment's effective_end_date to the day
	// before effectiveDate (so the new employment can take over).
	CloseEmployment(ctx context.Context, tx *gorm.DB, employmentID uuid.UUID, effectiveDate string) error
	// CreateEmployment persists a new employment record and returns its ID.
	CreateEmployment(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID, data CareerEmployment) (uuid.UUID, error)
	// SetEmployeeInactive marks an offboarded/retired employee as inactive.
	SetEmployeeInactive(ctx context.Context, tx *gorm.DB, employeeID uuid.UUID) error
}

// NumberingGenerator generates the next document number for a document type
// (see backend/internal/pkg/numbering for the concrete implementation).
type NumberingGenerator interface {
	Generate(ctx context.Context, documentType string) (string, error)
}

// GeneratedDocumentRef adalah hasil Generate Document (Phase 5) yang
// dikembalikan oleh DocumentGenerator — dipakai untuk menampilkan histori
// dokumen di UI contract/movement. Field hanya yang dibutuhkan module ini.
type GeneratedDocumentRef struct {
	ID           string    `json:"id"`
	TemplateID   string    `json:"template_id"`
	DocumentType string    `json:"document_type"`
	ReferenceType string   `json:"reference_type"`
	ReferenceID  string    `json:"reference_id"`
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	GeneratedBy  *string   `json:"generated_by,omitempty"`
	GeneratedAt  time.Time `json:"generated_at"`
}

// DocumentGenerator mengabstraksi shared document generator (documenttemplate)
// agar module ini hanya bergantung pada interface sempit (pola yang sama dengan
// ApprovalEngine/Notifier/CareerExecutor). Implementasinya di-wire dari main.go.
type DocumentGenerator interface {
	Generate(ctx context.Context, req DocumentGenerateRequest) (*GeneratedDocumentRef, error)
	ListByReference(ctx context.Context, referenceType, referenceID string, page, perPage int) ([]GeneratedDocumentRef, int64, error)
}

// DocumentGenerateRequest berisi data yang dibutuhkan Generate Document —
// Values adalah map variable ter-resolve (contract/movement/employee/company).
type DocumentGenerateRequest struct {
	DocumentType  string
	ReferenceType string
	ReferenceID   string
	Values        map[string]string
	GeneratedBy   string
}

// CareerExecutor methods all receive a *gorm.DB transaction (tx) opened by
// ExecuteMovementTx: every HR data change runs on the caller's transaction so
// movement execution is atomic (plan §12.2) — if any step fails, the whole
// execution rolls back (employment unchanged, movement stays approved).
// Implementations (the adapter in cmd/server/main.go) run their repository
// calls through the Tx variants of the employee repository.

// Service untuk business logic Employee Movement & Career Management.
type Service struct {
	repo                *Repository
	logger              *zap.Logger
	approvalEngine      ApprovalEngine
	careerExecutor      CareerExecutor
	notifier            Notifier
	performanceProvider PerformanceProvider
	competencyProvider  CompetencyProvider
	okrProvider         OKRProvider
	numberingService    NumberingGenerator
	docGenerator        DocumentGenerator
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

func (s *Service) SetPerformanceProvider(p PerformanceProvider) {
	s.performanceProvider = p
}

func (s *Service) SetCompetencyProvider(c CompetencyProvider) {
	s.competencyProvider = c
}

func (s *Service) SetOKRProvider(o OKRProvider) {
	s.okrProvider = o
}

// SetDocumentGenerator wires the shared document generator (documenttemplate)
// into this service so GenerateMovementDocument / GenerateContractDocument can
// produce PDF dari template aktif (plan §16/§17).
func (s *Service) SetDocumentGenerator(g DocumentGenerator) {
	s.docGenerator = g
}

// SetNumberingService wires the document numbering package so
// CreateMovement/CreateContract can auto-generate a number when the caller
// leaves it blank.
func (s *Service) SetNumberingService(ns NumberingGenerator) {
	s.numberingService = ns
}

// movementAuditJSON membungkus movement menjadi JSON string untuk kolom
// old_data/new_data pada audit trail (plan §12.6). Nilai dikembalikan sebagai
// *string agar selaras dengan kolom JSON (pola OrganizationHistory /
// payroll before_json/after_json). nil → nil (tidak ada snapshot).
func movementAuditJSON(m *EmployeeMovement) *string {
	if m == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	s := string(b)
	return &s
}

// statusPtr returns a pointer to the string form of a movement status, for
// the audit trail's old_status/new_status columns.
func statusPtr(s MovementStatus) *string {
	str := string(s)
	return &str
}

// recordAudit mencatat satu perubahan lifecycle movement ke tabel audit trail
// (plan §12.6: CREATED, UPDATED, SUBMITTED, APPROVED, REJECTED, CANCELLED,
// EXECUTED). Best-effort: kegagalan menyimpan audit hanya di-log dan TIDAK
// menggagalkan operasi movement utama (pola organization.captureHistory).
// actedBy nil untuk aksi yang tidak membawa konteks user (mis. callback).
func (s *Service) recordAudit(ctx context.Context, movementID uuid.UUID, action MovementAuditAction, oldStatus, newStatus *string, oldData, newData *string, reason string, actedBy *uuid.UUID) {
	audit := &EmployeeMovementAudit{
		MovementID: movementID,
		Action:     action,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		OldData:    oldData,
		NewData:    newData,
		ActedBy:    actedBy,
	}
	if reason != "" {
		audit.Reason = &reason
	}
	if err := s.repo.CreateAudit(ctx, audit); err != nil {
		s.logger.Warn("Failed to record movement audit",
			zap.String("movement_id", movementID.String()),
			zap.String("action", string(action)),
			zap.Error(err),
		)
	}
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
//
// Beyond G-7, when the movement carries a target position it also runs the
// position conflict check (plan §12.3): the target position must not be
// occupied by another employee's open employment at the effective date. The
// same check is repeated atomically inside ExecuteMovement — this early check
// just gives the HR user immediate feedback while drafting.
func (s *Service) validateMovement(ctx context.Context, m *EmployeeMovement) error {
	hasActiveContract := false
	if m.MovementType == MovementTypeContractExtension {
		has, err := s.repo.HasActiveContractByEmployeeID(ctx, m.EmployeeID)
		if err != nil {
			return err
		}
		hasActiveContract = has
	}
	if err := validateMovementFields(m.MovementType, m.ToOrganizationID, m.ToPositionID, m.ToEmploymentStatusID, hasActiveContract); err != nil {
		return err
	}

	if m.ToPositionID != nil && m.MovementType != MovementTypeContractExtension {
		occupied, err := s.repo.PositionConflict(ctx, nil, *m.ToPositionID, m.EmployeeID, m.EffectiveDate)
		if err != nil {
			return fmt.Errorf("failed to check target position conflict: %w", err)
		}
		if occupied {
			return &MovementConflictError{Message: "target position is already occupied by another employee at the effective date"}
		}
	}
	return nil
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

	// Auto-generate the decision letter number when the caller leaves it
	// blank (document numbering settings feature). "employee_movement" must
	// stay in sync with numbering.DocumentTypeEmployeeMovement.
	decisionLetterNumber := req.DecisionLetterNumber
	if decisionLetterNumber == "" && s.numberingService != nil {
		generated, err := s.numberingService.Generate(ctx, "employee_movement")
		if err != nil {
			return nil, fmt.Errorf("failed to generate decision letter number: %w", err)
		}
		decisionLetterNumber = generated
	}

	movement := &EmployeeMovement{
		CreatedBy:            authctx.GetUserID(ctx),
		UpdatedBy:            authctx.GetUserID(ctx),
		EmployeeID:           employeeUUID,
		MovementType:         MovementType(req.MovementType),
		DecisionLetterNumber: decisionLetterNumber,
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

	// Persist snapshot names (plan §12.5) so history keeps the master-data
	// names as they are at creation time.
	s.fillMovementSnapshot(ctx, movement)

	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement created",
		zap.String("employee_id", req.EmployeeID),
		zap.String("movement_type", req.MovementType),
		zap.String("movement_id", movement.ID.String()),
	)

	// Audit trail (plan §12.6): CREATED, snapshot awal movement.
	s.recordAudit(ctx, movement.ID, MovementAuditActionCreated, nil, statusPtr(movement.Status), nil, movementAuditJSON(movement), "", movement.CreatedBy)

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

// fillMovementSnapshot resolves the display names for the movement's
// from_*/to_* references and stores them on the movement row itself
// (plan §12.5). The snapshot is persisted at create/update time so history
// keeps the master-data names as they were then — renaming an organization,
// position or employment status later does not rewrite the past.
//
// Resolution is best-effort: when a reference id can't be resolved the name
// is left empty, and response-time enrichment (G-4) fills it from live master
// data as a fallback.
func (s *Service) fillMovementSnapshot(ctx context.Context, m *EmployeeMovement) {
	var orgIDs, posIDs, statusIDs []uuid.UUID
	if m.FromOrganizationID != nil {
		orgIDs = append(orgIDs, *m.FromOrganizationID)
	}
	if m.ToOrganizationID != nil {
		orgIDs = append(orgIDs, *m.ToOrganizationID)
	}
	if m.FromPositionID != nil {
		posIDs = append(posIDs, *m.FromPositionID)
	}
	if m.ToPositionID != nil {
		posIDs = append(posIDs, *m.ToPositionID)
	}
	if m.FromEmploymentStatusID != nil {
		statusIDs = append(statusIDs, *m.FromEmploymentStatusID)
	}
	if m.ToEmploymentStatusID != nil {
		statusIDs = append(statusIDs, *m.ToEmploymentStatusID)
	}

	if names, err := s.repo.GetOrganizationNamesByIDs(ctx, orgIDs); err == nil {
		if m.FromOrganizationID != nil {
			m.FromOrganizationName = names[m.FromOrganizationID.String()]
		}
		if m.ToOrganizationID != nil {
			m.ToOrganizationName = names[m.ToOrganizationID.String()]
		}
	} else {
		s.logger.Warn("failed to resolve organization names for movement snapshot", zap.Error(err))
	}

	if names, err := s.repo.GetPositionNamesByIDs(ctx, posIDs); err == nil {
		if m.FromPositionID != nil {
			m.FromPositionName = names[m.FromPositionID.String()]
		}
		if m.ToPositionID != nil {
			m.ToPositionName = names[m.ToPositionID.String()]
		}
	} else {
		s.logger.Warn("failed to resolve position names for movement snapshot", zap.Error(err))
	}

	if names, err := s.repo.GetEmploymentStatusNamesByIDs(ctx, statusIDs); err == nil {
		if m.FromEmploymentStatusID != nil {
			m.FromEmploymentStatusName = names[m.FromEmploymentStatusID.String()]
		}
		if m.ToEmploymentStatusID != nil {
			m.ToEmploymentStatusName = names[m.ToEmploymentStatusID.String()]
		}
	} else {
		s.logger.Warn("failed to resolve employment status names for movement snapshot", zap.Error(err))
	}
}

// enrichMovementResponses fills employee/organization/position/status display
// names on movement responses (single or list) with batch queries (G-4).
//
// Snapshot-aware (plan §12.5): names already persisted as movement snapshots
// (ToResponse carries them from the row) are preserved — live master data is
// only used as a fallback for rows created before the snapshot migration.
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

	// Snapshot-aware (plan §12.5): rows with a persisted snapshot already
	// carry their names, so the batch query is only needed when some field is
	// still empty (legacy movements created before migration 083).
	needsOrgNames, needsPosNames, needsStatusNames := false, false, false
	for i := range responses {
		if (responses[i].FromOrganizationID != nil && responses[i].FromOrganizationName == "") ||
			(responses[i].ToOrganizationID != nil && responses[i].ToOrganizationName == "") {
			needsOrgNames = true
		}
		if (responses[i].FromPositionID != nil && responses[i].FromPositionName == "") ||
			(responses[i].ToPositionID != nil && responses[i].ToPositionName == "") {
			needsPosNames = true
		}
		if (responses[i].FromEmploymentStatusID != nil && responses[i].FromEmploymentStatusName == "") ||
			(responses[i].ToEmploymentStatusID != nil && responses[i].ToEmploymentStatusName == "") {
			needsStatusNames = true
		}
	}

	if needsOrgNames {
		if names, err := s.repo.GetOrganizationNamesByIDs(ctx, orgList); err == nil {
			for i := range responses {
				if responses[i].FromOrganizationID != nil && responses[i].FromOrganizationName == "" {
					responses[i].FromOrganizationName = names[*responses[i].FromOrganizationID]
				}
				if responses[i].ToOrganizationID != nil && responses[i].ToOrganizationName == "" {
					responses[i].ToOrganizationName = names[*responses[i].ToOrganizationID]
				}
			}
		} else {
			s.logger.Warn("failed to resolve organization names for movements", zap.Error(err))
		}
	}

	if needsPosNames {
		if names, err := s.repo.GetPositionNamesByIDs(ctx, posList); err == nil {
			for i := range responses {
				if responses[i].FromPositionID != nil && responses[i].FromPositionName == "" {
					responses[i].FromPositionName = names[*responses[i].FromPositionID]
				}
				if responses[i].ToPositionID != nil && responses[i].ToPositionName == "" {
					responses[i].ToPositionName = names[*responses[i].ToPositionID]
				}
			}
		} else {
			s.logger.Warn("failed to resolve position names for movements", zap.Error(err))
		}
	}

	if needsStatusNames {
		if names, err := s.repo.GetEmploymentStatusNamesByIDs(ctx, statusList); err == nil {
			for i := range responses {
				if responses[i].FromEmploymentStatusID != nil && responses[i].FromEmploymentStatusName == "" {
					responses[i].FromEmploymentStatusName = names[*responses[i].FromEmploymentStatusID]
				}
				if responses[i].ToEmploymentStatusID != nil && responses[i].ToEmploymentStatusName == "" {
					responses[i].ToEmploymentStatusName = names[*responses[i].ToEmploymentStatusID]
				}
			}
		} else {
			s.logger.Warn("failed to resolve employment status names for movements", zap.Error(err))
		}
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

	// Audit trail (plan §12.6): simpan snapshot SEBELUM perubahan (masih state
	// lama — updated_by belum ditimpa user saat ini).
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

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

	// Refresh snapshot names (plan §12.5) after the to_* fields changed.
	s.fillMovementSnapshot(ctx, movement)

	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return nil, err
	}

	// Audit trail (plan §12.6): UPDATED, snapshot sebelum & sesudah.
	s.recordAudit(ctx, movement.ID, MovementAuditActionUpdated, oldStatus, statusPtr(movement.Status), oldData, movementAuditJSON(movement), "", movement.UpdatedBy)

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

	// Audit trail (plan §12.6): snapshot sebelum status diubah jadi pending
	// (approval_instance_id sudah set karena memang bagian transisi submit).
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

	movement.Status = MovementStatusPendingApproval
	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement submitted for approval",
		zap.String("movement_id", movement.ID.String()),
		zap.String("instance_id", instanceID),
	)

	// Audit trail (plan §12.6): SUBMITTED — acted_by = user yang submit
	// (bukan updated_by lama dari create).
	s.recordAudit(ctx, movement.ID, MovementAuditActionSubmitted, oldStatus, statusPtr(movement.Status), oldData, movementAuditJSON(movement), "", authctx.GetUserID(ctx))

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

	// Audit trail (plan §12.6): snapshot sebelum status diubah (masih
	// pending_approval).
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

	now := time.Now()
	auditAction := MovementAuditActionUpdated
	switch status {
	case "APPROVED":
		movement.Status = MovementStatusApproved
		movement.ApprovedAt = &now
		auditAction = MovementAuditActionApproved
	case "REJECTED":
		movement.Status = MovementStatusRejected
		if note != "" {
			movement.Notes = &note
		}
		auditAction = MovementAuditActionRejected
	case "CANCELLED":
		movement.Status = MovementStatusCancelled
		auditAction = MovementAuditActionCancelled
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

	// Audit trail (plan §12.6): APPROVED / REJECTED / CANCELLED — acted_by
	// kosong (aksi datang dari push-callback Central Approval, bukan user).
	s.recordAudit(ctx, movement.ID, auditAction, oldStatus, statusPtr(movement.Status), oldData, movementAuditJSON(movement), note, nil)

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

// normalizeDate memotong nilai tanggal menjadi YYYY-MM-DD. Driver database
// (MySQL/SQLite) dapat mengembalikan DATETIME/RFC3339 untuk kolom DATE — sama
// seperti yang sudah diantisipasi dayBefore — sehingga career timeline selalu
// menampilkan tanggal bersih.
func normalizeDate(date string) string {
	if len(date) >= 10 {
		if _, err := time.Parse("2006-01-02", date[:10]); err == nil {
			return date[:10]
		}
	}
	return date
}

// dayBefore returns the date one day before the given date. Accepts both
// plain YYYY-MM-DD and RFC3339 timestamps (MySQL returns DATETIME values for
// DATE columns), so movement execution is robust regardless of driver.
func dayBefore(date string) (string, error) {
	// Normalize: strip the time portion when an RFC3339 value is stored.
	date = normalizeDate(date)
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid effective_date %q: %w", date, err)
	}
	return d.AddDate(0, 0, -1).Format("2006-01-02"), nil
}

// ExecuteMovement mengeksekusi pergerakan secara ATOMIC (enhancement plan
// §12.2). Selain mengubah status movement menjadi executed, transaksi HR data
// juga dijalankan dalam SATU database transaction (G-1 + P0):
//   - promotion/demotion/mutation/status_change/other → buat employment baru
//     (to_* + effective_date), tutup employment aktif lama (effective_end_date
//     = effective_date - 1), simpan to_employment_id di movement.
//   - offboarding/retirement → tutup employment aktif lama dan tandai
//     employee non-aktif (keputusan §11.3).
//   - contract_extension → tanpa perubahan employment.
//
// Sebelum perubahan HR dijalankan, dua konflik divalidasi di dalam transaksi:
//   - §12.3 position conflict — target position tidak boleh terisi employment
//     terbuka employee lain pada effective date.
//   - §12.4 effective date conflict — employee tidak boleh memiliki
//     employment terbuka yang mulai pada/ setelah effective date (mencegah
//     overlap / backdate).
//
// Jika salah satu langkah gagal, seluruh transaksi di-ROLLBACK: employment
// lama utuh dan movement tetap berstatus approved (bisa di-retry HR).
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

	createsEmployment := movementCreatesEmployment(movement.MovementType)
	deactivatesEmployee := movementDeactivatesEmployee(movement.MovementType)

	// Audit trail (plan §12.6): snapshot & status sebelum eksekusi.
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

	// contract_extension → tanpa perubahan HR data; update status cukup
	// atomic dengan satu query (tanpa transaksi lintas modul).
	if !createsEmployment && !deactivatesEmployee {
		if err := s.repo.ExecuteMovement(ctx, uid, executorUUID, nil); err != nil {
			return err
		}
		after, err := s.repo.FindMovementByID(ctx, uid)
		if err == nil {
			s.recordAudit(ctx, movement.ID, MovementAuditActionExecuted, oldStatus, statusPtr(after.Status), oldData, movementAuditJSON(after), "", &executorUUID)
		} else {
			s.logger.Warn("failed to reload movement after execution for audit", zap.String("movement_id", movement.ID.String()), zap.Error(err))
		}
		return nil
	}

	if s.careerExecutor == nil {
		return fmt.Errorf("career executor not configured: cannot execute movement type '%s'", movement.MovementType)
	}

	// Transaksi atomik: conflict detection + perubahan HR data + update status
	// movement dijalankan pada satu transaksi (plan §12.2).
	err = s.repo.ExecuteMovementTx(ctx, uid, executorUUID, func(tx *gorm.DB) (*uuid.UUID, error) {
		if createsEmployment && movement.ToPositionID != nil {
			// §12.3 — target position conflict (hard check, diulang di execute
			// agar tidak hanya bergantung pada validasi frontend/validasi draft).
			occupied, err := s.repo.PositionConflict(ctx, tx, *movement.ToPositionID, movement.EmployeeID, movement.EffectiveDate)
			if err != nil {
				return nil, fmt.Errorf("failed to check target position conflict: %w", err)
			}
			if occupied {
				return nil, &MovementConflictError{Message: "target position is already occupied by another employee at the effective date"}
			}
		}

		// §12.4 — effective date conflict: tidak boleh ada employment terbuka
		// (open-ended) yang dimulai pada/ setelah effective date. Dicek juga
		// untuk offboarding/retirement agar menutup employment future-dated
		// tidak menghasilkan periode invalid (end sebelum start-nya sendiri).
		if createsEmployment || deactivatesEmployee {
			overlap, err := s.repo.EmploymentEffectiveDateConflict(ctx, tx, movement.EmployeeID, movement.EffectiveDate)
			if err != nil {
				return nil, fmt.Errorf("failed to check employment period conflict: %w", err)
			}
			if overlap {
				return nil, &MovementConflictError{Message: "movement effective date overlaps an existing employment period of the employee"}
			}
		}

		current, err := s.careerExecutor.FindCurrentEmployment(ctx, tx, movement.EmployeeID)
		if err != nil {
			return nil, fmt.Errorf("failed to find current employment: %w", err)
		}

		// Tutup employment aktif lama (effective_end_date = effective_date - 1).
		if current != nil {
			endDate, err := dayBefore(movement.EffectiveDate)
			if err != nil {
				return nil, err
			}
			if err := s.careerExecutor.CloseEmployment(ctx, tx, current.ID, endDate); err != nil {
				return nil, fmt.Errorf("failed to close previous employment: %w", err)
			}
		}

		if createsEmployment {
			data := CareerEmployment{
				OrganizationID:       movement.ToOrganizationID,
				PositionID:           movement.ToPositionID,
				EmploymentStatusID:   movement.ToEmploymentStatusID,
				DecisionLetterNumber: movement.DecisionLetterNumber,
				DecisionLetterDate:   movement.DecisionLetterDate,
				EffectiveDate:        movement.EffectiveDate,
			}
			newEmploymentID, err := s.careerExecutor.CreateEmployment(ctx, tx, movement.EmployeeID, data)
			if err != nil {
				return nil, fmt.Errorf("failed to create new employment: %w", err)
			}
			return &newEmploymentID, nil
		}

		if deactivatesEmployee {
			if err := s.careerExecutor.SetEmployeeInactive(ctx, tx, movement.EmployeeID); err != nil {
				return nil, fmt.Errorf("failed to mark employee inactive: %w", err)
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	// Audit trail (plan §12.6): EXECUTED — reload untuk menangkap state
	// akhir (status executed + to_employment_id hasil eksekusi).
	after, err := s.repo.FindMovementByID(ctx, uid)
	if err == nil {
		s.recordAudit(ctx, movement.ID, MovementAuditActionExecuted, oldStatus, statusPtr(after.Status), oldData, movementAuditJSON(after), "", &executorUUID)
	} else {
		s.logger.Warn("failed to reload movement after execution for audit", zap.String("movement_id", movement.ID.String()), zap.Error(err))
	}

	s.notifyMovementOutcome(ctx, movement.EmployeeID, "MOVEMENT_EXECUTED", "employeemovement", movement.ID)
	return nil
}

// CancelMovement membatalkan pergerakan. Per plan §12.16, perilaku bergantung
// pada status saat ini:
//
//	draft     → dibatalkan langsung oleh HR (repo.CancelMovement draft-only).
//	approved  → pembatalan menjadi Cancellation Request: approval instance
//	            baru dibuat di module "employeemovement_cancellation",
//	            movement masuk status cancellation_pending; hasil akhir
//	            (cancelled / kembali approved) diputuskan Central Approval
//	            dan ditangani HandleCancellationStatusChange.
//	lain      → error (executed/rejected/cancelled/pending tidak bisa dibatalkan).
//
// Flow (opsional) di-auto-resolve dari module "employeemovement_cancellation"
// bila req.FlowID kosong (pola G-3 sama seperti SubmitMovement).
func (s *Service) CancelMovement(ctx context.Context, id string, req CancelMovementRequest) (*MovementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Audit trail (plan §12.6): snapshot sebelum dibatalkan.
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

	switch movement.Status {
	case MovementStatusDraft:
		if err := s.repo.CancelMovement(ctx, uid); err != nil {
			return nil, err
		}
		// Audit trail (plan §12.6): CANCELLED — newData mencerminkan status
		// akhir (repo meng-update DB langsung, jadi salinan di memori
		// disesuaikan); acted_by = user yang cancel (authctx).
		after := *movement
		after.Status = MovementStatusCancelled
		s.recordAudit(ctx, movement.ID, MovementAuditActionCancelled, oldStatus, statusPtr(after.Status), oldData, movementAuditJSON(&after), "", authctx.GetUserID(ctx))

		responses := []MovementResponse{after.ToResponse()}
		s.enrichMovementResponses(ctx, responses)
		return &responses[0], nil

	case MovementStatusApproved:
		// Cancellation Request melalui Central Approval Module (plan §12.16).
		if s.approvalEngine == nil {
			return nil, fmt.Errorf("approval engine not configured")
		}
		flowID := ""
		if req.FlowID != nil && *req.FlowID != "" {
			flowID = *req.FlowID
		} else if resolved, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, "employeemovement_cancellation"); err == nil {
			flowID = resolved
		}
		if flowID == "" {
			return nil, fmt.Errorf("approval flow not configured: provide flow_id or activate an approval flow for module employeemovement_cancellation")
		}

		instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "employeemovement_cancellation", movement.ID.String(), flowID)
		if err != nil {
			return nil, fmt.Errorf("failed to create cancellation approval instance: %w", err)
		}
		parsedInstanceID, parseErr := uuid.Parse(instanceID)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid cancellation approval instance id: %w", parseErr)
		}
		if err := s.repo.SetCancellationRequested(ctx, uid, parsedInstanceID, req.Reason); err != nil {
			return nil, err
		}

		s.logger.Info("Employee movement cancellation requested for approval",
			zap.String("movement_id", movement.ID.String()),
			zap.String("instance_id", instanceID),
		)

		// Audit trail (plan §12.6): CANCELLATION_REQUESTED — newData mencerminkan
		// status akhir (repo meng-update DB langsung, salinan di memori disesuaikan);
		// acted_by = user yang mengajukan pembatalan (authctx).
		after := *movement
		after.Status = MovementStatusCancellationPending
		after.CancellationApprovalInstanceID = &parsedInstanceID
		if req.Reason != nil {
			after.Notes = req.Reason
		}
		s.recordAudit(ctx, movement.ID, MovementAuditActionCancellationRequested, oldStatus, statusPtr(after.Status), oldData, movementAuditJSON(&after), "", authctx.GetUserID(ctx))

		responses := []MovementResponse{after.ToResponse()}
		s.enrichMovementResponses(ctx, responses)
		return &responses[0], nil

	default:
		return nil, fmt.Errorf("movement cannot be cancelled in status '%s'", movement.Status)
	}
}

// HandleCancellationStatusChange is invoked by the approval module's push-based
// status callback when a Cancellation Request instance (module
// "employeemovement_cancellation") reaches a final state, so the movement's
// own status follows the cancellation decision (plan §12.16):
//
//	APPROVED → movement dibatalkan (status cancelled)
//	REJECTED → pembatalan ditolak, movement kembali approved
//	CANCELLED → request dibatalkan, movement kembali approved
//
// Callback hanya diproses ketika movement sedang dalam status
// cancellation_pending; callback basi (untuk movement yang sudah berubah)
// diabaikan.
func (s *Service) HandleCancellationStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	movement, err := s.repo.FindMovementByID(ctx, documentID)
	if err != nil {
		return err
	}
	if movement.Status != MovementStatusCancellationPending {
		return nil
	}

	// Audit trail (plan §12.6): snapshot sebelum status diubah (masih
	// cancellation_pending).
	oldData := movementAuditJSON(movement)
	oldStatus := statusPtr(movement.Status)

	auditAction := MovementAuditActionUpdated
	switch status {
	case "APPROVED":
		movement.Status = MovementStatusCancelled
		auditAction = MovementAuditActionCancelled
	case "REJECTED", "CANCELLED":
		movement.Status = MovementStatusApproved
		if note != "" {
			movement.Notes = &note
		}
		auditAction = MovementAuditActionCancellationRejected
	default:
		return nil
	}

	s.logger.Info("Employee movement cancellation status updated via approval status handler",
		zap.String("movement_id", movement.ID.String()),
		zap.String("cancellation_approval_status", status),
	)
	if err := s.repo.UpdateMovement(ctx, movement); err != nil {
		return err
	}

	// Audit trail (plan §12.6): acted_by kosong (aksi datang dari push-callback
	// Central Approval, bukan user).
	s.recordAudit(ctx, movement.ID, auditAction, oldStatus, statusPtr(movement.Status), oldData, movementAuditJSON(movement), note, nil)

	// Best-effort notification ke employee pemilik movement.
	switch movement.Status {
	case MovementStatusCancelled:
		s.notifyMovementOutcome(ctx, movement.EmployeeID, "MOVEMENT_REJECTED", "employeemovement", movement.ID)
	}
	return nil
}

// ListMovementAudits mengembalikan audit trail satu movement (plan §12.6),
// terurut acted_at DESC (baru dulu) dengan pagination.
func (s *Service) ListMovementAudits(ctx context.Context, id string, page, perPage int) (*PaginatedMovementAuditResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	audits, total, err := s.repo.ListAuditsByMovementID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []MovementAuditResponse
	for _, a := range audits {
		responses = append(responses, a.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementAuditResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// =========================================================================
// Contract Expiry Management (plan §12.13)
// =========================================================================

// contractExpiryReminderDays adalah jadwal reminder kontrak (plan §12.13):
// H-30, H-14, H-7, H-1 sebelum end_date.
var contractExpiryReminderDays = []int{30, 14, 7, 1}

// contractExpiryHRPermission adalah permission yang menandai user sebagai HR
// untuk penerima reminder kontrak (plan §12.13: "Notification dikirim kepada
// HR").
const contractExpiryHRPermission = "employeemovement.view"

// notifyContractEvent mengirim notifikasi kontrak kepada employee pemilik
// kontrak (via akun user terhubung) DAN seluruh user HR (yang punya
// permission contractExpiryHRPermission). Best-effort: kegagalan hanya di-log
// (pola sama notifyMovementOutcome). params diisi mis. nomor kontrak + jumlah
// hari tersisa.
func (s *Service) notifyContractEvent(ctx context.Context, contract EmployeeContract, notifType string, params []string) {
	if s.notifier == nil {
		return
	}

	// Kumpulkan penerima dalam satu set user id (dedup) — employee pemilik
	// kontrak ditambah seluruh user HR (yang punya permission
	// employeemovement.view). Tanpa dedup, user yang merupakan employee
	// sekaligus HR (mis. staf HR yang juga punya kontrak) akan menerima
	// notifikasi dua kali.
	recipients := make(map[uuid.UUID]struct{})

	if userID, err := s.repo.FindUserIDByEmployeeID(ctx, contract.EmployeeID); err == nil && userID != nil {
		recipients[*userID] = struct{}{}
	} else if err != nil {
		s.logger.Warn("failed to resolve employee user id for contract notification",
			zap.String("contract_id", contract.ID.String()),
			zap.Error(err),
		)
	}

	hrUserIDs, err := s.repo.FindUserIDsWithPermission(ctx, contractExpiryHRPermission)
	if err != nil {
		s.logger.Warn("failed to resolve HR users for contract notification",
			zap.String("contract_id", contract.ID.String()),
			zap.Error(err),
		)
		return
	}
	for _, hrID := range hrUserIDs {
		recipients[hrID] = struct{}{}
	}

	for recipientID := range recipients {
		if err := s.notifier.Notify(ctx, recipientID, notifType, params, "employeemovement", contract.ID); err != nil {
			s.logger.Warn("failed to send contract notification",
				zap.String("contract_id", contract.ID.String()),
				zap.String("recipient", recipientID.String()),
				zap.String("notif_type", notifType),
				zap.Error(err),
			)
		}
	}
}

// addDays mengembalikan tanggal YYYY-MM-DD yang berjarak `days` dari `date`
// (days negatif = mundur). Dipakai menghitung target reminder H-N.
func addDays(date string, days int) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.AddDate(0, 0, days).Format("2006-01-02")
}

// ProcessContractExpiration menjalankan proses harian manajemen kedaluwarsa
// kontrak (plan §12.13):
//
//  1. MARK EXPIRED — kontrak status=active yang end_date < hari ini dipindah
//     ke status expired, lalu employee & HR dinotifikasi CONTRACT_EXPIRED.
//  2. REMINDER — untuk tiap jadwal H-30 / H-14 / H-7 / H-1, kontrak active
//     yang berakhir tepat N hari lagi dinotifikasi CONTRACT_EXPIRING (ke
//     employee pemilik + user HR).
//
// Method ini murni per-tenant (context membawa company_id) dan sengaja TANPA
// tahu jadwal — scheduler (goroutine + time.Ticker di main.go) yang memanggil
// method ini harian. Best-effort per kontrak: kegagalan notifikasi satu
// kontrak tidak menggagalkan sisa proses.
func (s *Service) ProcessContractExpiration(ctx context.Context) error {
	today := time.Now().Format("2006-01-02")

	// 1. Mark expired: kontrak active yang sudah lewat end_date.
	expired, err := s.repo.FindContractsExpiredBefore(ctx, today)
	if err != nil {
		return err
	}
	if len(expired) > 0 {
		ids := make([]uuid.UUID, 0, len(expired))
		for _, c := range expired {
			ids = append(ids, c.ID)
		}
		if err := s.repo.MarkContractsExpired(ctx, ids); err != nil {
			return err
		}
		s.logger.Info("Contracts marked expired", zap.Int("count", len(expired)))
		for _, c := range expired {
			s.notifyContractEvent(ctx, c, "CONTRACT_EXPIRED", []string{c.ContractNumber})
		}
	}

	// 2. Reminder H-30 / H-14 / H-7 / H-1.
	for _, days := range contractExpiryReminderDays {
		target := addDays(today, days)
		contracts, err := s.repo.FindContractsExpiringOn(ctx, target)
		if err != nil {
			s.logger.Warn("failed to list contracts for expiry reminder", zap.Int("days", days), zap.Error(err))
			continue
		}
		for _, c := range contracts {
			s.logger.Info("Contract expiry reminder",
				zap.String("contract_id", c.ID.String()),
				zap.String("contract_number", c.ContractNumber),
				zap.Int("days_left", days),
			)
			s.notifyContractEvent(ctx, c, "CONTRACT_EXPIRING", []string{c.ContractNumber, addDays(today, days)})
		}
	}
	return nil
}

// =========================================================================
// Movement Documents (plan §12.15)
// =========================================================================

// CreateMovementDocument menambahkan metadata dokumen ke sebuah movement.
// Alur upload: file fisik di-upload dulu lewat endpoint upload generik
// (POST /api/v1/tenant/uploads) yang mengembalikan file_url; service ini
// hanya memvalidasi movement ada + menyimpan metadata (document_type,
// file_name, file_url).
func (s *Service) CreateMovementDocument(ctx context.Context, movementID string, req CreateMovementDocumentRequest) (*MovementDocumentResponse, error) {
	uid, err := uuid.Parse(movementID)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	// Pastikan movement benar-benar ada — selain jadi FK guard, ini memberi
	// error yang jelas (404/500 bukan FK violation cryptic).
	if _, err := s.repo.FindMovementByID(ctx, uid); err != nil {
		return nil, err
	}

	doc := &EmployeeMovementDocument{
		MovementID:   uid,
		DocumentType: MovementDocumentType(req.DocumentType),
		FileName:     req.FileName,
		FileURL:      req.FileURL,
		UploadedBy:   authctx.GetUserID(ctx),
	}
	if err := s.repo.CreateMovementDocument(ctx, doc); err != nil {
		return nil, err
	}

	s.logger.Info("Employee movement document created",
		zap.String("movement_id", movementID),
		zap.String("document_type", req.DocumentType),
		zap.String("file_name", req.FileName),
	)

	response := doc.ToResponse()
	return &response, nil
}

// ListMovementDocuments mengembalikan daftar dokumen sebuah movement
// (plan §12.15), terurut created_at DESC dengan pagination.
func (s *Service) ListMovementDocuments(ctx context.Context, movementID string, page, perPage int) (*PaginatedMovementDocumentResponse, error) {
	uid, err := uuid.Parse(movementID)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	documents, total, err := s.repo.ListDocumentsByMovementID(ctx, uid, page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []MovementDocumentResponse
	for _, d := range documents {
		responses = append(responses, d.ToResponse())
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedMovementDocumentResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// DeleteMovementDocument menghapus metadata dokumen movement.
func (s *Service) DeleteMovementDocument(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid document id: %w", err)
	}
	return s.repo.DeleteMovementDocument(ctx, uid)
}

// =========================================================================
// Career History (plan §12.8) — read model
// =========================================================================

// movementFromToLabel membangun label perpindahan "dari → ke" dari snapshot
// names movement (plan §12.5) — mis. "IT → Finance" atau "Staff → Supervisor".
// Setiap sisi memakai posisi bila ada, fallback organisasi, dan dilewati bila
// keduanya kosong (mis. movement tanpa from_* saat dibuat dari nol).
func movementFromToLabel(m *EmployeeMovement) string {
	from := ""
	if m.FromPositionName != "" {
		from = m.FromPositionName
	} else if m.FromOrganizationName != "" {
		from = m.FromOrganizationName
	}
	to := ""
	if m.ToPositionName != "" {
		to = m.ToPositionName
	} else if m.ToOrganizationName != "" {
		to = m.ToOrganizationName
	}
	switch {
	case from != "" && to != "":
		return from + " → " + to
	case to != "":
		return to
	case from != "":
		return from
	default:
		return ""
	}
}

// GetCareerHistory menyusun timeline karier seorang karyawan (plan §12.8) dari
// tiga sumber transaksional — TANPA tabel duplikasi:
//
//	employee_movements  → event MOVEMENT (hanya status executed)
//	employments         → event JOINED (employment pertama) + current position
//	employee_contracts  → event CONTRACT
//
// Timeline terurut kronologis ascending (test §17.6: Join → Promotion →
// Mutation → ... → Current Position). Nama organisasi/posisi/status
// employment di-resolve via batch query (G-4); movement memakai snapshot names
// (plan §12.5) sehingga histori tidak berubah walau master data berganti.
func (s *Service) GetCareerHistory(ctx context.Context, employeeID string) (*CareerHistoryResponse, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	employments, err := s.repo.FindEmploymentsByEmployeeID(ctx, uid)
	if err != nil {
		return nil, err
	}
	movements, err := s.repo.FindExecutedMovementsByEmployeeID(ctx, uid)
	if err != nil {
		return nil, err
	}
	contracts, err := s.repo.FindAllContractsByEmployeeID(ctx, uid)
	if err != nil {
		return nil, err
	}

	data := CareerHistoryData{
		EmployeeID: uid.String(),
		Timeline:   []CareerTimelineEntry{},
	}

	// Employee display info (G-4) — best-effort.
	if empInfo, err := s.repo.GetEmployeeInfoByIDs(ctx, []uuid.UUID{uid}); err == nil {
		if emp, ok := empInfo[uid.String()]; ok {
			data.EmployeeName = emp.Name
			data.EmployeeCode = emp.Code
		}
	} else {
		s.logger.Warn("failed to resolve employee info for career history", zap.String("employee_id", employeeID), zap.Error(err))
	}

	// Resolve org/posisi/status names untuk employment (movement sudah punya
	// snapshot sendiri).
	var orgIDs, posIDs, statusIDs []uuid.UUID
	for _, e := range employments {
		if e.OrganizationID != nil {
			orgIDs = append(orgIDs, *e.OrganizationID)
		}
		if e.PositionID != nil {
			posIDs = append(posIDs, *e.PositionID)
		}
		if e.EmploymentStatusID != nil {
			statusIDs = append(statusIDs, *e.EmploymentStatusID)
		}
	}
	orgNames := map[string]string{}
	posNames := map[string]string{}
	statusNames := map[string]string{}
	if n, err := s.repo.GetOrganizationNamesByIDs(ctx, orgIDs); err == nil {
		orgNames = n
	} else {
		s.logger.Warn("failed to resolve organization names for career history", zap.Error(err))
	}
	if n, err := s.repo.GetPositionNamesByIDs(ctx, posIDs); err == nil {
		posNames = n
	} else {
		s.logger.Warn("failed to resolve position names for career history", zap.Error(err))
	}
	if n, err := s.repo.GetEmploymentStatusNamesByIDs(ctx, statusIDs); err == nil {
		statusNames = n
	} else {
		s.logger.Warn("failed to resolve employment status names for career history", zap.Error(err))
	}

	// JOINED — dari employment pertama (effective_date terawal).
	if len(employments) > 0 {
		first := employments[0]
		title := ""
		if first.PositionID != nil {
			title = posNames[first.PositionID.String()]
		}
		if title == "" && first.OrganizationID != nil {
			title = orgNames[first.OrganizationID.String()]
		}
		eid := first.ID.String()
		data.Timeline = append(data.Timeline, CareerTimelineEntry{
			Date:         normalizeDate(first.EffectiveDate),
			EventType:    "JOINED",
			Title:        title,
			EmploymentID: &eid,
		})
	}

	// MOVEMENT — dari movement yang sudah dieksekusi, urut effective_date ASC.
	for i := range movements {
		m := &movements[i]
		label := movementFromToLabel(m)
		mt := string(m.MovementType)
		mid := m.ID.String()
		entry := CareerTimelineEntry{
			Date:         normalizeDate(m.EffectiveDate),
			EventType:    "MOVEMENT",
			Title:        mt,
			MovementType: &mt,
			MovementID:   &mid,
		}
		if label != "" {
			entry.Description = &label
		}
		data.Timeline = append(data.Timeline, entry)
	}

	// CONTRACT — dari employee_contracts, urut start_date ASC.
	for i := range contracts {
		c := &contracts[i]
		ct := string(c.ContractType)
		cid := c.ID.String()
		entry := CareerTimelineEntry{
			Date:         normalizeDate(c.StartDate),
			EventType:    "CONTRACT",
			Title:        c.ContractNumber,
			ContractType: &ct,
			ContractID:   &cid,
		}
		if c.EndDate != nil {
			desc := normalizeDate(c.StartDate) + " – " + normalizeDate(*c.EndDate)
			entry.Description = &desc
		}
		data.Timeline = append(data.Timeline, entry)
	}

	// Urut kronologis (tanggal ASC). Untuk tanggal sama: JOINED → MOVEMENT →
	// CONTRACT agar urutan "Join dulu, lalu transaksi, lalu kontrak" stabil.
	sort.SliceStable(data.Timeline, func(i, j int) bool {
		if data.Timeline[i].Date != data.Timeline[j].Date {
			return data.Timeline[i].Date < data.Timeline[j].Date
		}
		return careerEventPriority(data.Timeline[i].EventType) < careerEventPriority(data.Timeline[j].EventType)
	})

	// Current position — employment terbuka terakhir (effective_end_date NULL)
	// bila ada; fallback employment terakhir menurut effective_date.
	if pos := s.currentPosition(employments, orgNames, posNames, statusNames); pos != nil {
		data.CurrentPosition = pos
	}

	return &CareerHistoryResponse{Success: true, Data: data}, nil
}

// careerEventPriority mengurutkan event dengan tanggal sama: JOINED pertama,
// lalu MOVEMENT, lalu CONTRACT.
func careerEventPriority(eventType string) int {
	switch eventType {
	case "JOINED":
		return 0
	case "MOVEMENT":
		return 1
	default:
		return 2
	}
}

// currentPosition memilih employment aktif terakhir karyawan: yang masih
// terbuka (effective_end_date NULL) dengan effective_date terbesar; jika
// tidak ada yang terbuka, pakai employment terakhir menurut effective_date.
// Mengembalikan nil bila karyawan belum memiliki employment sama sekali.
func (s *Service) currentPosition(employments []careerEmploymentRow, orgNames, posNames, statusNames map[string]string) *CareerPositionInfo {
	if len(employments) == 0 {
		return nil
	}
	best := 0
	for i := range employments {
		e := &employments[i]
		// Bandingkan tanggal yang sudah dinormalisasi (YYYY-MM-DD) agar driver
		// yang mengembalikan DATETIME/RFC3339 tidak memengaruhi urutan.
		currentBest := &employments[best]
		if e.EffectiveEndDate == nil && currentBest.EffectiveEndDate != nil {
			// Employment terbuka selalu menang atas yang sudah ditutup.
			best = i
		} else if (e.EffectiveEndDate == nil) == (currentBest.EffectiveEndDate == nil) {
			// Sama-sama terbuka atau sama-sama tertutup → pilih tanggal terbesar.
			if normalizeDate(e.EffectiveDate) > normalizeDate(currentBest.EffectiveDate) {
				best = i
			}
		}
	}
	e := &employments[best]
	info := &CareerPositionInfo{
		EmploymentID:  e.ID.String(),
		EffectiveDate: normalizeDate(e.EffectiveDate),
	}
	if e.OrganizationID != nil {
		info.OrganizationID = e.OrganizationID.String()
		info.OrganizationName = orgNames[e.OrganizationID.String()]
	}
	if e.PositionID != nil {
		info.PositionID = e.PositionID.String()
		info.PositionName = posNames[e.PositionID.String()]
	}
	if e.EmploymentStatusID != nil {
		info.EmploymentStatusID = e.EmploymentStatusID.String()
		info.EmploymentStatusName = statusNames[e.EmploymentStatusID.String()]
	}
	return info
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

	// Auto-generate the contract number when the caller leaves it blank
	// (document numbering settings feature). "employee_contract" must stay
	// in sync with numbering.DocumentTypeEmployeeContract.
	contractNumber := req.ContractNumber
	if contractNumber == "" && s.numberingService != nil {
		generated, err := s.numberingService.Generate(ctx, "employee_contract")
		if err != nil {
			return nil, fmt.Errorf("failed to generate contract number: %w", err)
		}
		contractNumber = generated
	}

	contract := &EmployeeContract{
		CreatedBy:      authctx.GetUserID(ctx),
		UpdatedBy:      authctx.GetUserID(ctx),
		EmployeeID:     employeeUUID,
		ContractNumber: contractNumber,
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
		zap.String("contract_number", contractNumber),
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

// =========================================================================
// Generate Document (plan §16/§17) — business module menyiapkan data, shared
// Document Generator (documenttemplate) yang melakukan rendering + PDF.
// =========================================================================

// employeeValues memetakan profil karyawan ke variable {{employee.*}} yang
// terdaftar di VariableRegistry. Key dengan nilai kosong tetap disertakan agar
// placeholder ter-replace menjadi kosong (bukan literal {{...}}). Position &
// organization TIDAK dimasukkan di sini — keduanya kontekstual (movement
// memakai posisi/org tujuan; contract memakai employment aktif).
func employeeValues(emp *EmployeeProfileData) map[string]string {
	return map[string]string{
		"employee.employee_id":      emp.EmployeeID,
		"employee.name":             emp.Name,
		"employee.nik":              emp.NIK,
		"employee.family_id":        emp.FamilyID,
		"employee.mother_name":      emp.MotherName,
		"employee.gender":           emp.Gender,
		"employee.dob":              emp.DOB,
		"employee.pob":              emp.POB,
		"employee.nationality_type": emp.NationalityType,
		"employee.nationality_id":   emp.NationalityID,
		"employee.passport":         emp.Passport,
		"employee.phone_number":     emp.PhoneNumber,
		"employee.email":            emp.Email,
		"employee.linkedin":         emp.LinkedIn,
		"employee.instagram":        emp.Instagram,
		"employee.religion":         emp.Religion,
		"employee.marital_status":   emp.MaritalStatus,
		"employee.status":           emp.Status,
		"employee.join_date":        emp.JoinDate,
	}
}

// GenerateMovementDocument menghasilkan PDF SK Movement dari template aktif
// MOVEMENT_SK. Movement harus berstatus approved/executed (SK yang diterbitkan).
func (s *Service) GenerateMovementDocument(ctx context.Context, movementID, actorID string) (*GeneratedDocumentRef, error) {
	if s.docGenerator == nil {
		return nil, fmt.Errorf("document generator is not configured")
	}
	uid, err := uuid.Parse(movementID)
	if err != nil {
		return nil, fmt.Errorf("invalid movement id: %w", err)
	}
	movement, err := s.repo.FindMovementByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if movement.Status != MovementStatusApproved && movement.Status != MovementStatusExecuted {
		return nil, fmt.Errorf("document can only be generated for approved or executed movements")
	}

	emp, err := s.repo.GetEmployeeProfile(ctx, movement.EmployeeID)
	if err != nil {
		return nil, err
	}
	values := map[string]string{
		"employee.position":          movement.ToPositionName,
		"employee.organization":      movement.ToOrganizationName,
		"movement.number":            movement.DecisionLetterNumber,
		"movement.effective_date":    movement.EffectiveDate,
		"movement.previous_position": movement.FromPositionName,
		"movement.new_position":      movement.ToPositionName,
	}
	for k, v := range employeeValues(emp) {
		values[k] = v
	}
	return s.docGenerator.Generate(ctx, DocumentGenerateRequest{
		DocumentType:  documenttemplate.DocumentTypeMovementSK,
		ReferenceType: "movement",
		ReferenceID:   movementID,
		Values:        values,
		GeneratedBy:   actorID,
	})
}

// GenerateContractDocument menghasilkan PDF Perjanjian Kerja dari template aktif
// CONTRACT_AGREEMENT.
func (s *Service) GenerateContractDocument(ctx context.Context, contractID, actorID string) (*GeneratedDocumentRef, error) {
	if s.docGenerator == nil {
		return nil, fmt.Errorf("document generator is not configured")
	}
	uid, err := uuid.Parse(contractID)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}
	contract, err := s.repo.FindContractByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	emp, err := s.repo.GetEmployeeProfile(ctx, contract.EmployeeID)
	if err != nil {
		return nil, err
	}
	values := map[string]string{
		"employee.position":     emp.Position,
		"employee.organization": emp.Organization,
		"contract.number":       contract.ContractNumber,
		"contract.start_date":   contract.StartDate,
	}
	if contract.EndDate != nil {
		values["contract.end_date"] = *contract.EndDate
	}
	for k, v := range employeeValues(emp) {
		values[k] = v
	}
	return s.docGenerator.Generate(ctx, DocumentGenerateRequest{
		DocumentType:  documenttemplate.DocumentTypeContractAgreement,
		ReferenceType: "contract",
		ReferenceID:   contractID,
		Values:        values,
		GeneratedBy:   actorID,
	})
}

// ListGeneratedDocuments menampilkan histori dokumen yang digenerate untuk
// sebuah reference (movement/contract), terbaru dulu.
func (s *Service) ListGeneratedDocuments(ctx context.Context, referenceType, referenceID string, page, perPage int) ([]GeneratedDocumentRef, int64, error) {
	if s.docGenerator == nil {
		return nil, 0, fmt.Errorf("document generator is not configured")
	}
	return s.docGenerator.ListByReference(ctx, referenceType, referenceID, page, perPage)
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
func (s *Service) ListContracts(ctx context.Context, page, perPage int, status, search string) (*PaginatedContractResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	contracts, total, err := s.repo.ListContracts(ctx, page, perPage, status, search)
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

// GetMovementReport mengembalikan Movement Report (plan §12.17): jumlah
// movement per tipe (promosi/demosi/mutasi/dll) dan per status, dengan filter
// opsional periode (effective_date), organisasi, posisi, employee, tipe, dan
// status. Semua filter opsional — kosong berarti semua data.
func (s *Service) GetMovementReport(ctx context.Context, dateFrom, dateTo, organizationID, positionID, employeeID, movementType, status string) (*MovementReportResponse, error) {
	// Validasi periode: rentang terbalik (date_from > date_to) tidak bermakna —
	// tolak dengan 400 (MovementValidationError) alih-alih mengembalikan 0 baris.
	if dateFrom != "" && dateTo != "" && dateFrom > dateTo {
		return nil, &MovementValidationError{Message: "date_from cannot be after date_to"}
	}

	f := MovementReportFilter{
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		MovementType: movementType,
		Status:       status,
	}
	if organizationID != "" {
		uid, err := uuid.Parse(organizationID)
		if err != nil {
			return nil, &MovementValidationError{Message: "invalid organization_id"}
		}
		f.OrganizationID = &uid
	}
	if positionID != "" {
		uid, err := uuid.Parse(positionID)
		if err != nil {
			return nil, &MovementValidationError{Message: "invalid position_id"}
		}
		f.PositionID = &uid
	}
	if employeeID != "" {
		uid, err := uuid.Parse(employeeID)
		if err != nil {
			return nil, &MovementValidationError{Message: "invalid employee_id"}
		}
		f.EmployeeID = &uid
	}

	byType, err := s.repo.CountMovementsByType(ctx, f)
	if err != nil {
		return nil, err
	}
	byStatus, err := s.repo.CountMovementsByStatus(ctx, f)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, c := range byType {
		total += c
	}

	return &MovementReportResponse{
		Success: true,
		Data: MovementReportData{
			Total:    total,
			ByType:   byType,
			ByStatus: byStatus,
		},
	}, nil
}

// GetContractReport mengembalikan Contract Report (plan §12.17): jumlah
// kontrak per status (active/expired/extended/terminated) plus jumlah kontrak
// aktif yang berakhir dalam 30 hari ke depan (expiring — plan §12.18).
func (s *Service) GetContractReport(ctx context.Context) (*ContractReportResponse, error) {
	byStatus, err := s.repo.CountContractsByStatus(ctx)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, c := range byStatus {
		total += c
	}

	today := time.Now().Format("2006-01-02")
	in30Days := addDays(today, 30)
	expiring, err := s.repo.CountExpiringContracts(ctx, today, in30Days)
	if err != nil {
		return nil, err
	}

	return &ContractReportResponse{
		Success: true,
		Data: ContractReportData{
			Total:    total,
			ByStatus: byStatus,
			Expiring: expiring,
		},
	}, nil
}

// GetHRDashboard mengembalikan data untuk kartu HR Dashboard (plan §12.18)
// dalam satu panggilan:
//   - movement_by_type: jumlah movement per tipe (semua status)
//   - pending_approval: jumlah movement berstatus pending_approval
//   - effective_this_month: jumlah movement dengan effective_date di bulan berjalan
//   - contracts: ringkasan kontrak (active / expiring < 30 hari / expired)
//
// Memakai agregasi yang sama dengan Movement/Contract Report (§12.17) sehingga
// kartu dashboard dan halaman report selalu konsisten.
func (s *Service) GetHRDashboard(ctx context.Context) (*HRDashboardResponse, error) {
	// Movement per tipe — tanpa filter (semua status).
	byType, err := s.repo.CountMovementsByType(ctx, MovementReportFilter{})
	if err != nil {
		return nil, err
	}

	// Pending approval.
	pendingStatus, err := s.repo.CountMovementsByStatus(ctx, MovementReportFilter{Status: string(MovementStatusPendingApproval)})
	if err != nil {
		return nil, err
	}
	pendingApproval := pendingStatus[string(MovementStatusPendingApproval)]

	// Effective this month — rentang hari pertama s/d terakhir bulan berjalan.
	now := time.Now()
	monthStart := now.Format("2006-01") + "-01"
	monthEnd := lastDayOfMonth(now)
	effectiveThisMonth, err := s.repo.CountMovementsEffectiveBetween(ctx, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	// Ringkasan kontrak.
	contractByStatus, err := s.repo.CountContractsByStatus(ctx)
	if err != nil {
		return nil, err
	}
	today := now.Format("2006-01-02")
	expiring, err := s.repo.CountExpiringContracts(ctx, today, addDays(today, 30))
	if err != nil {
		return nil, err
	}

	return &HRDashboardResponse{
		Success: true,
		Data: HRDashboardData{
			MovementByType:     byType,
			PendingApproval:    pendingApproval,
			EffectiveThisMonth: effectiveThisMonth,
			Contracts: ContractSummaryData{
				Active:   contractByStatus[string(ContractStatusActive)],
				Expiring: expiring,
				Expired:  contractByStatus[string(ContractStatusExpired)],
			},
		},
	}, nil
}

// lastDayOfMonth mengembalikan tanggal terakhir bulan dari t (YYYY-MM-DD),
// dipakai menghitung rentang "effective this month" pada HR Dashboard.
func lastDayOfMonth(t time.Time) string {
	firstOfNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	last := firstOfNext.AddDate(0, 0, -1)
	return last.Format("2006-01-02")
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

// =========================================================================
// Promotion Eligibility (plan §12.10/§12.11)
// =========================================================================

const (
	// Threshold default untuk eligibility. Sesuai contoh plan §12.10/§12.11:
	// Minimum Service >= 24 months, Performance (KPI) >= 80, OKR >= 80,
	// Competency >= 80. Threshold masa kerja dapat di-override oleh
	// career_path_steps (minimum_service_months step berikutnya) pada
	// promotion-eligibility. Kebijakan nil-data: rule tanpa data
	// (Available=false) tidak memblokir eligible (pragmatis).
	eligibilityDefaultMinServiceMonths = 24
	eligibilityMinPerformanceScore     = 80.0
	eligibilityMinCompetencyScore      = 80.0
	eligibilityMinOKRScore             = 80.0
)

// GetMovementEligibility mengembalikan ringkasan eligibility umum seorang
// karyawan: masa kerja (bulan), posisi sekarang, skor performa & kompetensi
// terakhir (jika tersedia), career path (jika posisi termasuk path), dan
// hasil rule default (tenure >= 24, performance >= 80, competency >= 80).
// Provider nil / data tidak tersedia → rule not met dengan detail.
func (s *Service) GetMovementEligibility(ctx context.Context, employeeID string) (*MovementEligibilityResponse, error) {
	data, err := s.buildEligibility(ctx, employeeID, false)
	if err != nil {
		return nil, err
	}
	return &MovementEligibilityResponse{Success: true, Data: *data}, nil
}

// GetPromotionEligibility mengembalikan eligibility khusus promosi: rule
// tenure menggunakan minimum_service_months dari step berikutnya dalam career
// path (jika employee berada dalam path yang aktif), fallback 24 bulan.
// Menyertakan info target posisi promosi berikutnya.
func (s *Service) GetPromotionEligibility(ctx context.Context, employeeID string) (*PromotionEligibilityResponse, error) {
	data, err := s.buildEligibility(ctx, employeeID, true)
	if err != nil {
		return nil, err
	}
	promo := PromotionEligibilityData{
		EmployeeID:           data.EmployeeID,
		EmployeeName:         data.EmployeeName,
		EmployeeCode:         data.EmployeeCode,
		TenureMonths:         data.TenureMonths,
		CurrentPosition:      data.CurrentPosition,
		MinimumServiceMonths: eligibilityDefaultMinServiceMonths,
		PerformanceScore:     data.PerformanceScore,
		CompetencyScore:      data.CompetencyScore,
		OKRScore:             data.OKRScore,
	}
	// Ambil info next step + override minimum_service dari career path
	// (untuk promotion, tenure target = min_service_months step berikutnya).
	for _, r := range data.Rules {
		if r.Code == "career_path" && r.Met && data.CareerPathID != nil {
			// Cari step berikutnya — data sudah dihitung di buildEligibility.
			s.findPromotionNextStep(ctx, &promo, *data.CareerPathID, data.CurrentPosition)
			break
		}
	}

	promo.Rules = s.evaluatePromotionRules(data.TenureMonths, promo.MinimumServiceMonths,
		data.PerformanceScore, data.CompetencyScore, data.OKRScore)
	promo.Eligible = true
	for _, r := range promo.Rules {
		if r.Available && !r.Met {
			promo.Eligible = false
			break
		}
	}

	return &PromotionEligibilityResponse{Success: true, Data: promo}, nil
}

// buildEligibility adalah helper bersama untuk kedua endpoint eligibility.
// promotion=false: rule default 24/80/80, tanpa target next step.
// promotion=true: juga menghitung career path + next step info.
func (s *Service) buildEligibility(ctx context.Context, employeeID string, promotion bool) (*MovementEligibilityData, error) {
	uid, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}

	employments, err := s.repo.FindEmploymentsByEmployeeID(ctx, uid)
	if err != nil {
		return nil, err
	}

	data := MovementEligibilityData{
		EmployeeID: uid.String(),
	}

	// Employee display info (G-4) — best-effort.
	if empInfo, err := s.repo.GetEmployeeInfoByIDs(ctx, []uuid.UUID{uid}); err == nil {
		if emp, ok := empInfo[uid.String()]; ok {
			data.EmployeeName = emp.Name
			data.EmployeeCode = emp.Code
		}
	} else {
		s.logger.Warn("failed to resolve employee info for eligibility", zap.String("employee_id", employeeID), zap.Error(err))
	}

	// Posisi sekarang + tenure.
	var orgIDs, posIDs, statusIDs []uuid.UUID
	for _, e := range employments {
		if e.OrganizationID != nil {
			orgIDs = append(orgIDs, *e.OrganizationID)
		}
		if e.PositionID != nil {
			posIDs = append(posIDs, *e.PositionID)
		}
		if e.EmploymentStatusID != nil {
			statusIDs = append(statusIDs, *e.EmploymentStatusID)
		}
	}

	orgNames := map[string]string{}
	posNames := map[string]string{}
	statusNames := map[string]string{}
	if n, err := s.repo.GetOrganizationNamesByIDs(ctx, orgIDs); err == nil {
		orgNames = n
	} else {
		s.logger.Warn("failed to resolve org names for eligibility", zap.Error(err))
	}
	if n, err := s.repo.GetPositionNamesByIDs(ctx, posIDs); err == nil {
		posNames = n
	} else {
		s.logger.Warn("failed to resolve position names for eligibility", zap.Error(err))
	}
	if n, err := s.repo.GetEmploymentStatusNamesByIDs(ctx, statusIDs); err == nil {
		statusNames = n
	} else {
		s.logger.Warn("failed to resolve employment status names for eligibility", zap.Error(err))
	}

	data.CurrentPosition = s.currentPosition(employments, orgNames, posNames, statusNames)
	data.TenureMonths = computeTenureMonths(employments)

	// Career path — apakah posisi sekarang termasuk path aktif?
	if data.CurrentPosition != nil && data.CurrentPosition.PositionID != "" {
		if curPosID, err := uuid.Parse(data.CurrentPosition.PositionID); err == nil {
			if steps, err := s.repo.FindCareerPathStepsByPositionID(ctx, curPosID); err == nil && len(steps) > 0 {
				// Seluruh steps dari path pertama yang cocok.
				pathID := steps[0].CareerPathID
				if pathInfo, err := s.repo.FindCareerPathByID(ctx, pathID); err == nil {
					s := pathInfo.ID.String()
					data.CareerPathID = &s
					s2 := pathInfo.Name
					data.CareerPathName = &s2
				}
			}
		}
	}

	// Skor performa (KPI), kompetensi, dan OKR — best-effort melalui provider.
	// Movement hanya membaca hasil final, tidak menghitung KPI/OKR sendiri
	// (plan §12.11 Integration Principle).
	if s.performanceProvider != nil {
		if score, found, err := s.performanceProvider.LatestFinalScore(ctx, uid); err == nil && found {
			data.PerformanceScore = &score
		} else if err != nil {
			s.logger.Warn("failed to get performance score for eligibility", zap.String("employee_id", employeeID), zap.Error(err))
		}
	}
	if s.competencyProvider != nil {
		if score, found, err := s.competencyProvider.LatestScore(ctx, uid); err == nil && found {
			data.CompetencyScore = &score
		} else if err != nil {
			s.logger.Warn("failed to get competency score for eligibility", zap.String("employee_id", employeeID), zap.Error(err))
		}
	}
	if s.okrProvider != nil {
		if score, found, err := s.okrProvider.LatestScore(ctx, uid); err == nil && found {
			data.OKRScore = &score
		} else if err != nil {
			s.logger.Warn("failed to get okr score for eligibility", zap.String("employee_id", employeeID), zap.Error(err))
		}
	}

	// Evaluasi rule default.
	data.Rules = s.evaluateDefaultRules(data.TenureMonths, data.PerformanceScore, data.CompetencyScore, data.OKRScore)

	// Tambah rule career_path jika posisi ada dalam path aktif.
	if data.CareerPathID != nil && data.CareerPathName != nil {
		met := true
		detail := fmt.Sprintf("Posisi saat ini ada dalam career path '%s'", *data.CareerPathName)
		data.Rules = append(data.Rules, EligibilityRuleResult{
			Code: "career_path", Label: "Career Path",
			Met: met, Detail: &detail,
			Available: true,
		})
	}

	// Eligible dihitung hanya dari rule yang datanya tersedia (Available=true).
	// Rule tanpa data (mis. tenant belum menjalankan OKR/performance/competency)
	// dilaporkan tapi tidak memblokir — kebijakan pragmatis.
	data.Eligible = true
	for _, r := range data.Rules {
		if r.Available && !r.Met {
			data.Eligible = false
			break
		}
	}

	return &data, nil
} // evaluateDefaultRules membangun rule untuk movement-eligibility: minimum
// service (default 24), performance (>= 80), competency (>= 80), OKR (>= 80).
func (s *Service) evaluateDefaultRules(tenureMonths int, performanceScore, competencyScore, okrScore *float64) []EligibilityRuleResult {
	rules := make([]EligibilityRuleResult, 0, 4)

	tenureReq := fmt.Sprintf("%d bulan", eligibilityDefaultMinServiceMonths)
	tenureActual := fmt.Sprintf("%d bulan", tenureMonths)
	rules = append(rules, EligibilityRuleResult{
		Code: "minimum_service", Label: "Minimum Masa Kerja",
		Met:    tenureMonths >= eligibilityDefaultMinServiceMonths,
		Actual: &tenureActual, Required: &tenureReq,
		Available: true,
	})

	rules = append(rules, buildScoreRule("performance", "Skor Kinerja (KPI)", performanceScore, eligibilityMinPerformanceScore))
	rules = append(rules, buildScoreRule("competency", "Skor Kompetensi", competencyScore, eligibilityMinCompetencyScore))
	rules = append(rules, buildScoreRule("okr", "Skor OKR", okrScore, eligibilityMinOKRScore))

	return rules
}

// evaluatePromotionRules membangun rule untuk promotion-eligibility:
// minimum service (dari career path atau default 24), performance, competency,
// dan OKR.
func (s *Service) evaluatePromotionRules(tenureMonths, minService int, performanceScore, competencyScore, okrScore *float64) []EligibilityRuleResult {
	rules := make([]EligibilityRuleResult, 0, 4)

	tenureReq := fmt.Sprintf("%d bulan", minService)
	tenureActual := fmt.Sprintf("%d bulan", tenureMonths)
	rules = append(rules, EligibilityRuleResult{
		Code: "minimum_service", Label: "Minimum Masa Kerja (Target Promosi)",
		Met:    tenureMonths >= minService,
		Actual: &tenureActual, Required: &tenureReq,
		Available: true,
	})

	rules = append(rules, buildScoreRule("performance", "Skor Kinerja (KPI)", performanceScore, eligibilityMinPerformanceScore))
	rules = append(rules, buildScoreRule("competency", "Skor Kompetensi", competencyScore, eligibilityMinCompetencyScore))
	rules = append(rules, buildScoreRule("okr", "Skor OKR", okrScore, eligibilityMinOKRScore))

	return rules
}

// findPromotionNextStep mencari step berikutnya dalam career path setelah
// posisi employee sekarang dan mengisi field promotion di PromotionEligibilityData.
func (s *Service) findPromotionNextStep(ctx context.Context, promo *PromotionEligibilityData, careerPathID string, currentPos *CareerPositionInfo) {
	pid, err := uuid.Parse(careerPathID)
	if err != nil {
		return
	}
	steps, err := s.repo.ListCareerPathStepsByPathID(ctx, pid)
	if err != nil || len(steps) == 0 {
		return
	}
	// Cari indeks step yang match current position.
	idx := -1
	for i, st := range steps {
		if currentPos != nil && st.PositionID.String() == currentPos.PositionID {
			idx = i
			break
		}
	}
	if idx < 0 || idx >= len(steps)-1 {
		// Posisi tidak ditemukan di path atau sudah di step terakhir.
		return
	}
	nextStep := steps[idx+1]

	// Enrich nama posisi target.
	posNames := map[string]string{}
	if n, err := s.repo.GetPositionNamesByIDs(ctx, []uuid.UUID{nextStep.PositionID}); err == nil {
		posNames = n
	} else {
		s.logger.Warn("failed to resolve next step position name", zap.Error(err))
	}

	pn := nextStep.PositionID.String()
	promo.NextPositionID = &pn
	if name, ok := posNames[pn]; ok {
		promo.NextPositionName = &name
	}
	seq := nextStep.Sequence
	promo.NextPositionSeq = &seq
	promo.MinimumServiceMonths = eligibilityDefaultMinServiceMonths
	if nextStep.MinimumServiceMonths != nil && *nextStep.MinimumServiceMonths > 0 {
		promo.MinimumServiceMonths = *nextStep.MinimumServiceMonths
	}
}

// computeTenureMonths menghitung masa kerja (dalam bulan) dari employment
// pertama (effective_date terawal) hingga sekarang.
func computeTenureMonths(employments []careerEmploymentRow) int {
	if len(employments) == 0 {
		return 0
	}
	start := employments[0].EffectiveDate
	startDate, err := time.Parse("2006-01-02", normalizeDate(start))
	if err != nil {
		return 0
	}
	now := time.Now()
	months := int(now.Sub(startDate).Hours() / (24 * 30.44))
	if months < 0 {
		return 0
	}
	return months
} // buildScoreRule membuat rule result untuk satu nilai skor (perform/kompetensi/
// OKR) terhadap threshold. Jika nil (data tidak tersedia), rule dilaporkan
// dengan Available=false (met=false + detail) dan TIDAK ikut menentukan
// `eligible` — kebijakan pragmatis: tenant yang belum menjalankan modul
// terkait tidak membuat karyawannya otomatis ineligible.
func buildScoreRule(code, label string, score *float64, threshold float64) EligibilityRuleResult {
	if score == nil {
		msg := "data belum tersedia — rule dilewati dari perhitungan eligible"
		return EligibilityRuleResult{
			Code: code, Label: label,
			Met: false, Detail: &msg,
			Available: false,
		}
	}
	actual := fmt.Sprintf("%.0f", *score)
	required := fmt.Sprintf("%.0f", threshold)
	return EligibilityRuleResult{
		Code: code, Label: label,
		Met:    *score >= threshold,
		Actual: &actual, Required: &required,
		Available: true,
	}
}

// =========================================================================
// Career Path CRUD dipindahkan ke modul Career Intelligence (keputusan user
// 2026-08-10 — pemisahan transactional vs strategical). Modul ini hanya
// MEMBACA career_paths/career_path_steps untuk promotion eligibility
// (FindCareerPathStepsByPositionID / findPromotionNextStep di atas).
// =========================================================================
