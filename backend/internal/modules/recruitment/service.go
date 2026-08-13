package recruitment

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

var ErrInvalidStatusTransition = errors.New("invalid status transition")

// applicationStageOrder define urutan progresi status non-terminal (G-5).
// Status terminal (ACCEPTED/REJECTED/WITHDRAWN) sengaja tidak masuk sini —
// diperlakukan khusus di isValidStatusTransition.
var applicationStageOrder = map[CandidateStatus]int{
	CandStatusNew:         1,
	CandStatusScreened:    2,
	CandStatusShortlisted: 3,
	CandStatusInterviewed: 4,
	CandStatusOffered:     5,
}

func isTerminalStatus(s CandidateStatus) bool {
	return s == CandStatusAccepted || s == CandStatusRejected || s == CandStatusWithdrawn
}

// isValidStatusTransition menegakkan state machine G-5:
//   - from == to                         → true (no-op, ditangani caller)
//   - from terminal                      → false (tidak ada transisi keluar)
//   - to ∈ {ACCEPTED, REJECTED, WITHDRAWN} dan from non-terminal → true
//   - from & to sama-sama non-terminal    → true hanya jika order[to] >= order[from]
func isValidStatusTransition(from, to CandidateStatus) bool {
	if from == to {
		return true
	}
	if isTerminalStatus(from) {
		return false
	}
	if isTerminalStatus(to) {
		return true
	}
	fromOrder, fromOK := applicationStageOrder[from]
	toOrder, toOK := applicationStageOrder[to]
	if !fromOK || !toOK {
		return false
	}
	return toOrder >= fromOrder
}

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

// InternalCandidateProvider adalah interface narrow yang dipakai Recruitment
// untuk membaca employee internal yang eligible bagi sebuah target position
// dari Career Intelligence (plan S-4 — internal candidate via career path).
// Recruitment TIDAK menentukan career eligibility — CI yang menghitung;
// Recruitment hanya memproses aplikasi internal. Implementasi di-wire di
// cmd/server/main.go melalui adapter (careerintelligence.Service). Bila
// provider nil, list eligible dikembalikan kosong (fail-safe).
type InternalCandidateProvider interface {
	// EligibleEmployeesForPosition mengembalikan daftar employee internal yang
	// memegang source step dari career path menuju targetPositionID.
	EligibleEmployeesForPosition(ctx context.Context, targetPositionID uuid.UUID) ([]InternalCandidate, error)
}

// InternalCandidate adalah employee internal yang eligible untuk posisi target
// (hasil perhitungan CI — Recruitment hanya meneruskan).
type InternalCandidate struct {
	EmployeeID          string
	Name                string
	CurrentPositionID   string
	CurrentPositionName string
	SourceStepSequence  int
	TargetPositionID    string
	TargetPositionName  string
	PathID              string
	PathName            string
}

// SuccessionGapProvider adalah interface narrow yang dipakai Recruitment untuk
// memvalidasi fallback external recruitment (plan S-5 — succession plan →
// fallback external recruitment). Recruitment TIDAK menghitung readiness
// succession sendiri — CI yang menandai posisi kunci tanpa successor siap;
// Recruitment hanya membaca hasilnya untuk requisition reason_type=
// SUCCESSION_GAP. Implementasi di-wire di cmd/server/main.go melalui adapter
// (careerintelligence.Service). Bila provider nil, requisition tetap bisa
// dibuat tanpa error (fail-safe) — kolom succession_position_id tetap tersimpan.
type SuccessionGapProvider interface {
	// SuccessionGapForPosition mengembalikan true bila posisi adalah posisi
	// kunci tanpa successor siap (readiness READY_NOW) sehingga membutuhkan
	// fallback external recruitment.
	SuccessionGapForPosition(ctx context.Context, positionID uuid.UUID) (bool, error)
}

// TrainingHandoffProvider adalah interface narrow yang dipakai Recruitment
// untuk meneruskan kebutuhan training/development ke Training module setelah
// onboarding selesai (plan S-7 — onboarding → training handoff). Recruitment
// TIDAK mengeksekusi training — ia hanya menghasilkan kebutuhan (handoff);
// Training tetap source of truth. Implementasi di-wire di cmd/server/main.go
// melalui adapter (training.Service). Bila provider nil, onboarding tetap
// bisa diselesaikan tanpa error (fail-safe) — handoff hanya di-log.
type TrainingHandoffProvider interface {
	// CreateOnboardingNeed membuat training need source ONBOARDING untuk
	// employee yang baru menyelesaikan onboarding.
	CreateOnboardingNeed(ctx context.Context, employeeID, onboardingID uuid.UUID, reason string) error
}

// EmployeeProvider adalah interface narrow yang dipakai Recruitment untuk
// membuat employee baru dari offer eksternal yang diterima (plan G-4).
// Recruitment TIDAK membuat employee sendiri — ia menyerahkan data hire ke
// Employee module via adapter (employee.Service) dan menyimpan referensi
// employee.recruited_from_application_id. Best-effort: kegagalan provider
// tidak menggagalkan accept offer (di-log).
type EmployeeProvider interface {
	// CreateHiredEmployee membuat employee baru dari kandidat eksternal yang
	// menerima offer; mengembalikan employee UUID.
	CreateHiredEmployee(ctx context.Context, in EmployeeHireInput) (string, error)
}

// EmployeeHireInput membawa data minimal hire eksternal (G-4) — Recruitment
// tidak mengirim field employee lain; Employee module yang melengkapi.
type EmployeeHireInput struct {
	ApplicationID  string
	CandidateID    string
	Name           string
	Email          string
	Phone          string
	StartDate      string
	EmploymentType string
	OrganizationID string
	PositionID     string
	Salary         float64
}

// MovementProvider adalah interface narrow yang dipakai Recruitment untuk
// meneruskan hasil seleksi internal (offer diterima) ke Employee Movement
// (plan G-4 — internal hire → promotion/mutation), bukan employee baru.
type MovementProvider interface {
	// CreateHiredMovement membuat movement internal hire untuk employee yang
	// menang seleksi (candidate_type=INTERNAL).
	CreateHiredMovement(ctx context.Context, in MovementHireInput) error
}

// MovementHireInput membawa data minimal movement internal hire (G-4).
type MovementHireInput struct {
	EmployeeID     string
	ApplicationID  string
	OrganizationID string
	PositionID     string
	EffectiveDate  string
	Reason         string
}

// ApprovalEngine adalah interface narrow yang dipakai Recruitment untuk
// merutekan requisition melalui Central Approval (plan G-1). Recruitment TIDAK
// mengimplementasikan alur approval — ia hanya membuat instance & membaca
// status; Approval module tetap source of truth. Implementasi di-wire di
// cmd/server/main.go melalui adapter (approval.Service) — same
// narrow-interface-plus-adapter pattern payroll/leave/employeemovement.
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error)
	// GetActiveFlowIDForModule lets a requisition submission auto-resolve which
	// flow to route through when the client doesn't supply flow_id explicitly
	// (same pattern leave/employeemovement uses) — without this, a requisition
	// submitted without a flow_id stays in draft and never reaches the
	// Approval module.
	GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
}

type Service struct {
	repo                    *Repository
	logger                  *zap.Logger
	gapProvider             WorkforceGapProvider
	internalProvider        InternalCandidateProvider
	successionProvider      SuccessionGapProvider
	trainingHandoffProvider TrainingHandoffProvider
	approvalEngine          ApprovalEngine
	// G-4: handoff offer diterima → Employee (eksternal) / Employee Movement
	// (internal). Recruitment TIDAK membuat employee/movement sendiri — modul
	// terkait yang mengeksekusi via interface narrow (pola S-1..S-7).
	employeeProvider EmployeeProvider
	movementProvider MovementProvider
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

// SetInternalCandidateProvider wires the Career Intelligence module into this
// service (plan S-4) so recruiters can read internal candidates eligible for a
// target position (via career path) without Recruitment computing eligibility.
func (s *Service) SetInternalCandidateProvider(p InternalCandidateProvider) {
	s.internalProvider = p
}

// SetSuccessionGapProvider wires the Career Intelligence module into this
// service (plan S-5) so requisitions with reason_type=SUCCESSION_GAP can be
// validated as fallback external recruitment for a key position without a
// ready successor.
func (s *Service) SetSuccessionGapProvider(p SuccessionGapProvider) {
	s.successionProvider = p
}

// SetTrainingHandoffProvider wires the Training module into this service
// (plan S-7) so onboarding completion forwards a training need (handoff) —
// Training tetap source of truth; Recruitment hanya menghasilkan kebutuhan.
func (s *Service) SetTrainingHandoffProvider(p TrainingHandoffProvider) {
	s.trainingHandoffProvider = p
}

// SetApprovalEngine wires the central approval module into this service
// (plan G-1) so requisition submissions are routed through it (single approval
// path — manual approve dihapus, keputusan plan G-1/G-5).
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SetEmployeeProvider wires the Employee module into this service (plan G-4)
// so an accepted external offer creates an employee with the
// recruited_from_application_id reference (Employee → Application →
// Requisition → Position traceability). Best-effort: bila provider nil,
// accept offer tetap berhasil — employee creation hanya di-log.
func (s *Service) SetEmployeeProvider(p EmployeeProvider) {
	s.employeeProvider = p
}

// SetMovementProvider wires the Employee Movement module into this service
// (plan G-4) so an accepted internal offer forwards the selection result to a
// movement (promotion/mutation) instead of creating a new employee.
// Best-effort: bila provider nil, accept offer tetap berhasil.
func (s *Service) SetMovementProvider(p MovementProvider) {
	s.movementProvider = p
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
		OrganizationID:   orgID,
		Title:            req.Title,
		Department:       req.Department,
		EmploymentType:   req.EmploymentType,
		Location:         req.Location,
		MinSalary:        req.MinSalary,
		MaxSalary:        req.MaxSalary,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		SlotsAvailable:   1,
	}
	// G-2: priority default MEDIUM bila client tidak mengirim.
	if req.Priority != nil {
		r.Priority = RequisitionPriority(*req.Priority)
	} else {
		r.Priority = ReqPriorityMedium
	}
	// G-2: position_id (referensi master position, opsional).
	if req.PositionID != nil && *req.PositionID != "" {
		if uid, perr := uuid.Parse(*req.PositionID); perr == nil {
			r.PositionID = &uid
		}
	}
	// G-2: requisition_number auto-generated bila client kosong
	// (format REQ-YYYYMM-XXXX, pola nomor dokumen TRN-* training).
	if req.RequisitionNumber != "" {
		r.RequisitionNumber = req.RequisitionNumber
	} else {
		r.RequisitionNumber = generateRequisitionNumber()
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
	if req.SuccessionPositionID != nil {
		uid, _ := uuid.Parse(*req.SuccessionPositionID)
		r.SuccessionPositionID = &uid
	}
	// S-1: auto-resolve hiring need dari Workforce Intelligence ketika
	// requisition dibuat dengan reason WORKFORCE_GAP dan slots tidak
	// ditentukan eksplisit. Gagal resolve = tetap lanjut dengan default.
	s.resolveWorkforceGapSlots(ctx, r, req.SlotsAvailable)
	// S-5: validasi fallback external recruitment untuk reason SUCCESSION_GAP.
	s.validateSuccessionGapFallback(ctx, r)
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
	if req.RequisitionNumber != nil {
		r.RequisitionNumber = *req.RequisitionNumber
	}
	if req.Priority != nil {
		r.Priority = RequisitionPriority(*req.Priority)
	}
	if req.PositionID != nil {
		if *req.PositionID == "" {
			r.PositionID = nil
		} else if uid, perr := uuid.Parse(*req.PositionID); perr == nil {
			// Invalid non-empty UUID diabaikan (jangan simpan uuid.Nil) — pola
			// sama dengan Create path.
			r.PositionID = &uid
		}
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
	prevStatus := r.Status
	if req.Status != nil {
		r.Status = RequisitionStatus(*req.Status)
	}
	// G-2: opened_at diset otomatis saat requisition bertransisi menjadi OPEN
	// (dari non-OPEN). Approval APPROVED juga memicu ini via callback (lihat
	// HandleApprovalStatusChange).
	if r.Status == ReqStatusOpen && prevStatus != ReqStatusOpen {
		now := time.Now().UnixNano()
		r.OpenedAt = &now
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
	if req.SuccessionPositionID != nil {
		if *req.SuccessionPositionID == "" {
			r.SuccessionPositionID = nil
		} else {
			uid, _ := uuid.Parse(*req.SuccessionPositionID)
			r.SuccessionPositionID = &uid
		}
	}
	// S-1: resolve hiring need HANYA saat reason bertransisi menjadi
	// WORKFORCE_GAP (bukan sudah WORKFORCE_GAP) — sehingga update field lain
	// (title, status, dll.) pada requisition yang sudah tertaut gap tidak
	// menimpa slots_available yang tersimpan.
	if req.ReasonType != nil && *req.ReasonType == string(ReqReasonWorkforceGap) && prevReason != string(ReqReasonWorkforceGap) {
		s.resolveWorkforceGapSlots(ctx, r, req.SlotsAvailable)
	}
	// S-5: validasi fallback external recruitment saat reason bertransisi
	// menjadi SUCCESSION_GAP (pola sama dengan resolve S-1).
	if req.ReasonType != nil && *req.ReasonType == string(ReqReasonSuccessionGap) && prevReason != string(ReqReasonSuccessionGap) {
		s.validateSuccessionGapFallback(ctx, r)
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
// Job Offers (G-3)
// =========================================================================

func (s *Service) CreateOffer(ctx context.Context, req CreateOfferRequest) (*OfferResponse, error) {
	appID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	// Validasi aplikasi ada (kandidat + requisition).
	if _, err := s.repo.FindApplicationByID(ctx, appID); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}
	o := &JobOffer{
		ApplicationID:  appID,
		EmploymentType: req.EmploymentType,
		Salary:         req.Salary,
		Allowances:     req.Allowances,
		Benefits:       req.Benefits,
		StartDate:      req.StartDate,
		ExpiryDate:     req.ExpiryDate,
		Status:         OfferStatusDraft,
	}
	o.OfferNumber = generateOfferNumber()
	if err := s.repo.CreateOffer(ctx, o); err != nil {
		return nil, err
	}
	s.logger.Info("Offer created", zap.String("id", o.ID.String()), zap.String("application_id", o.ApplicationID.String()))
	return offerToResponse(o), nil
}

func (s *Service) GetOfferByID(ctx context.Context, id string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return offerToResponse(o), nil
}

func (s *Service) ListOffers(ctx context.Context, applicationID, status *string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	var appUUID *uuid.UUID
	if applicationID != nil && *applicationID != "" {
		uid, _ := uuid.Parse(*applicationID)
		appUUID = &uid
	}
	list, total, err := s.repo.ListOffers(ctx, appUUID, status, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]OfferResponse, 0, len(list))
	for _, o := range list {
		responses = append(responses, *offerToResponse(&o))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateOffer(ctx context.Context, id string, req UpdateOfferRequest) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	// Hanya draft yang bisa diedit (setelah submit/sent terkunci).
	if o.Status != OfferStatusDraft {
		return nil, fmt.Errorf("only draft offers can be updated, current status: %s", o.Status)
	}
	if req.EmploymentType != nil {
		o.EmploymentType = *req.EmploymentType
	}
	if req.Salary != nil {
		o.Salary = *req.Salary
	}
	if req.Allowances != nil {
		o.Allowances = *req.Allowances
	}
	if req.Benefits != nil {
		o.Benefits = *req.Benefits
	}
	if req.StartDate != nil {
		o.StartDate = *req.StartDate
	}
	if req.ExpiryDate != nil {
		o.ExpiryDate = *req.ExpiryDate
	}
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}
	return offerToResponse(o), nil
}

func (s *Service) DeleteOffer(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return err
	}
	if o.Status != OfferStatusDraft {
		return fmt.Errorf("only draft offers can be deleted, current status: %s", o.Status)
	}
	return s.repo.DeleteOffer(ctx, uid)
}

// SubmitOffer mengirim offer draft ke Central Approval (plan G-3):
// DRAFT → PENDING_APPROVAL + approval instance dibuat & approval_instance_id
// disimpan. Auto-resolve flow aktif untuk modul "recruitment_offer" bila
// client tidak mengirim flow_id (pola G-1 requisition).
func (s *Service) SubmitOffer(ctx context.Context, id, flowID string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if o.Status != OfferStatusDraft {
		return nil, fmt.Errorf("only draft offers can be submitted, current status: %s", o.Status)
	}
	if s.approvalEngine == nil {
		return nil, fmt.Errorf("approval engine not configured")
	}

	// Auto-resolve the active flow when no flow_id is supplied (G-3).
	if flowID == "" {
		if resolved, resolveErr := s.approvalEngine.GetActiveFlowIDForModule(ctx, "recruitment_offer"); resolveErr == nil {
			flowID = resolved
		}
	}
	if flowID == "" {
		return nil, fmt.Errorf("approval flow not configured: provide flow_id or activate an approval flow for module recruitment_offer")
	}

	instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "recruitment_offer", o.ID.String(), flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to create approval instance: %w", err)
	}
	if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
		o.ApprovalInstanceID = &parsedInstanceID
	}

	o.Status = OfferStatusPendingApproval
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}

	s.logger.Info("Offer submitted for approval",
		zap.String("offer_id", o.ID.String()),
		zap.String("instance_id", instanceID),
	)
	return offerToResponse(o), nil
}

// HandleOfferApprovalStatusChange adalah push-callback dari Central Approval
// (plan G-3): APPROVED/REJECTED/CANCELLED atas offer di-propagasi ke status
// offer. Hanya offer PENDING_APPROVAL yang diproses (idempotent).
func (s *Service) HandleOfferApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	o, err := s.repo.FindOfferByID(ctx, documentID)
	if err != nil {
		return err
	}
	if o.Status != OfferStatusPendingApproval {
		return nil
	}

	switch status {
	case "APPROVED":
		o.Status = OfferStatusApproved
	case "REJECTED":
		o.Status = OfferStatusRejected
		// Catatan: note reject TIDAK dipersist ke kolom Benefits (teks benefit
		// asli) — note approval tersimpan di instance Approval module.
	case "CANCELLED":
		o.Status = OfferStatusWithdrawn
	default:
		return nil
	}

	s.logger.Info("Offer status updated via approval status handler",
		zap.String("offer_id", o.ID.String()),
		zap.String("approval_status", status),
	)
	return s.repo.UpdateOffer(ctx, o)
}

// SendOffer mengirim offer yang sudah APPROVED ke kandidat: APPROVED → SENT
// + sent_at. Offer yang belum approved tidak bisa dikirim.
func (s *Service) SendOffer(ctx context.Context, id string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if o.Status != OfferStatusApproved {
		return nil, fmt.Errorf("only approved offers can be sent, current status: %s", o.Status)
	}
	now := time.Now().UnixNano()
	o.Status = OfferStatusSent
	o.SentAt = &now
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}
	s.logger.Info("Offer sent", zap.String("offer_id", o.ID.String()))
	return offerToResponse(o), nil
}

// AcceptOffer mencatat penerimaan kandidat: SENT → ACCEPTED + accepted_at.
// Offer yang sudah melewati expiry_date tidak bisa diterima (business rule
// G-3). Penerimaan otomatis memajukan aplikasi ke ACCEPTED dan menaikkan
// slots_filled requisition (pola yang sama dengan UpdateApplicationStatus).
func (s *Service) AcceptOffer(ctx context.Context, id string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if o.Status != OfferStatusSent {
		return nil, fmt.Errorf("only sent offers can be accepted, current status: %s", o.Status)
	}
	// Business rule: offer expired tidak dapat diterima.
	if isOfferExpired(o.ExpiryDate) {
		// Tandai expired + tolak aksi.
		o.Status = OfferStatusExpired
		_ = s.repo.UpdateOffer(ctx, o)
		return nil, fmt.Errorf("offer has expired and cannot be accepted")
	}
	now := time.Now().UnixNano()
	o.Status = OfferStatusAccepted
	o.AcceptedAt = &now
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}

	// Majukan aplikasi ke ACCEPTED + slots_filled requisition (pola G-1
	// UpdateApplicationStatus ACCEPTED — single source of truth). Guard
	// transisi status: hanya increment slots_filled bila aplikasi belum
	// ACCEPTED, supaya accept offer kedua di aplikasi yang sama (atau
	// accept setelah ACCEPTED manual) tidak double-count satu kandidat.
	if a, findErr := s.repo.FindApplicationByID(ctx, o.ApplicationID); findErr == nil && a != nil {
		wasAccepted := a.Status == CandStatusAccepted
		transitionErr := s.transitionApplicationStatus(ctx, a, CandStatusAccepted, nil, "")
		if transitionErr != nil {
			s.logger.Warn("offer accepted but application transition failed",
				zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.Error(transitionErr))
		}
		if err := s.repo.UpdateApplication(ctx, a); err != nil {
			s.logger.Warn("offer accepted but application update failed",
				zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.Error(err))
		}
		// didTransition: satu-satunya sinyal yang boleh dipakai untuk
		// men-trigger efek samping "aplikasi baru saja jadi ACCEPTED".
		// wasAccepted saja tidak cukup — bila transitionApplicationStatus
		// gagal (mis. aplikasi sudah REJECTED/WITHDRAWN, state machine
		// menolak), a.Status TIDAK berubah, tapi wasAccepted (dihitung
		// sebelum pemanggilan) tetap false. Tanpa guard transitionErr,
		// slots_filled dan handoffHiredEmployee bisa jalan untuk aplikasi
		// yang statusnya sebenarnya masih REJECTED/WITHDRAWN.
		didTransition := transitionErr == nil && !wasAccepted
		req, _ := s.repo.FindRequisitionByID(ctx, a.RequisitionID)
		if didTransition && req != nil {
			req.SlotsFilled++
			if req.SlotsFilled >= req.SlotsAvailable {
				req.Status = ReqStatusFilled
			}
			if err := s.repo.UpdateRequisition(ctx, req); err != nil {
				s.logger.Warn("offer accepted but requisition slots update failed",
					zap.String("requisition_id", req.ID.String()), zap.Error(err))
			}
		}

		// G-4: handoff Recruitment → Employee / Employee Movement. External
		// (candidate_type EXTERNAL) → Employee module membuat employee baru
		// dengan referensi recruited_from_application_id. Internal
		// (candidate_type INTERNAL ber-employee_id) → diteruskan ke Employee
		// Movement (promotion/mutation), bukan employee baru. Guard transisi
		// status: hanya dijalankan saat aplikasi BARU menjadi ACCEPTED, supaya
		// accept offer kedua di aplikasi yang sama tidak membuat employee/
		// movement duplikat (idempotensi yang sama dengan slots_filled G-3).
		// Best-effort: kegagalan downstream TIDAK menggagalkan accept offer.
		if didTransition {
			s.handoffHiredEmployee(ctx, o, a, req)
		}
	}

	s.logger.Info("Offer accepted", zap.String("offer_id", o.ID.String()))
	return offerToResponse(o), nil
}

// RejectOffer mencatat penolakan kandidat: SENT → REJECTED + rejected_at.
func (s *Service) RejectOffer(ctx context.Context, id string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if o.Status != OfferStatusSent {
		return nil, fmt.Errorf("only sent offers can be rejected, current status: %s", o.Status)
	}
	now := time.Now().UnixNano()
	o.Status = OfferStatusRejected
	o.RejectedAt = &now
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}
	s.logger.Info("Offer rejected by candidate", zap.String("offer_id", o.ID.String()))
	return offerToResponse(o), nil
}

// WithdrawOffer menarik kembali offer: DRAFT/APPROVED → WITHDRAWN (recruiter
// membatalkan sebelum/sesudah approval tapi belum dikirim).
func (s *Service) WithdrawOffer(ctx context.Context, id string) (*OfferResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	o, err := s.repo.FindOfferByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if o.Status != OfferStatusDraft && o.Status != OfferStatusApproved && o.Status != OfferStatusPendingApproval {
		return nil, fmt.Errorf("offer cannot be withdrawn from status: %s", o.Status)
	}
	o.Status = OfferStatusWithdrawn
	if err := s.repo.UpdateOffer(ctx, o); err != nil {
		return nil, err
	}
	s.logger.Info("Offer withdrawn", zap.String("offer_id", o.ID.String()))
	return offerToResponse(o), nil
}

// SubmitRequisition mengirim requisition draft ke Central Approval (plan G-1):
// DRAFT → SUBMITTED + approval instance dibuat & approval_instance_id
// disimpan. Auto-resolve flow aktif untuk modul "recruitment" bila client
// tidak mengirim flow_id (pola employeemovement G-3).
func (s *Service) SubmitRequisition(ctx context.Context, id, flowID string) (*RequisitionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindRequisitionByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if r.Status != ReqStatusDraft {
		return nil, fmt.Errorf("only draft requisitions can be submitted, current status: %s", r.Status)
	}
	if s.approvalEngine == nil {
		return nil, fmt.Errorf("approval engine not configured")
	}

	// Auto-resolve the active flow when no flow_id is supplied (G-1).
	if flowID == "" {
		if resolved, resolveErr := s.approvalEngine.GetActiveFlowIDForModule(ctx, "recruitment"); resolveErr == nil {
			flowID = resolved
		}
	}
	if flowID == "" {
		return nil, fmt.Errorf("approval flow not configured: provide flow_id or activate an approval flow for module recruitment")
	}

	instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, "recruitment", r.ID.String(), flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to create approval instance: %w", err)
	}
	if parsedInstanceID, parseErr := uuid.Parse(instanceID); parseErr == nil {
		r.ApprovalInstanceID = &parsedInstanceID
	}

	r.Status = ReqStatusSubmitted
	if err := s.repo.UpdateRequisition(ctx, r); err != nil {
		return nil, err
	}

	s.logger.Info("Requisition submitted for approval",
		zap.String("requisition_id", r.ID.String()),
		zap.String("instance_id", instanceID),
	)
	return requisitionToResponse(r), nil
}

// HandleApprovalStatusChange adalah push-callback dari Central Approval (plan
// G-1): keputusan APPROVED/REJECTED/CANCELLED atas requisition di-propagasi
// ke status requisition. Hanya requisition berstatus SUBMITTED yang diproses
// (idempotent — callback ganda tidak menimpa status final).
func (s *Service) HandleApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	r, err := s.repo.FindRequisitionByID(ctx, documentID)
	if err != nil {
		return err
	}
	if r.Status != ReqStatusSubmitted {
		return nil
	}

	switch status {
	case "APPROVED":
		r.Status = ReqStatusOpen
		// G-2: opened_at diset saat requisition dibuka oleh approval.
		now := time.Now().UnixNano()
		r.OpenedAt = &now
	case "REJECTED":
		r.Status = ReqStatusRejected
		// Catatan: note reject TIDAK dipersist ke kolom Requirements — itu teks
		// persyaratan job yang asli (menimpanya = korupsi data). Note approval
		// sudah tersimpan di instance Approval module (source of truth).
	case "CANCELLED":
		r.Status = ReqStatusCancelled
	default:
		return nil
	}

	s.logger.Info("Requisition status updated via approval status handler",
		zap.String("requisition_id", r.ID.String()),
		zap.String("approval_status", status),
	)
	return s.repo.UpdateRequisition(ctx, r)
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
	if req.CandidateNumber != nil && *req.CandidateNumber != "" {
		c.CandidateNumber = *req.CandidateNumber
	} else {
		c.CandidateNumber = generateCandidateNumber()
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
	// G-4: candidate_type + employee_id (internal hire).
	if err := applyCandidateTypeFields(c, req.CandidateType, req.EmployeeID); err != nil {
		return nil, err
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
	if req.CandidateNumber != nil {
		c.CandidateNumber = *req.CandidateNumber
	}
	// G-4: candidate_type + employee_id (internal hire).
	if err := applyCandidateTypeFields(c, req.CandidateType, req.EmployeeID); err != nil {
		return nil, err
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
	if newStage, stageErr := s.repo.FindStageByCode(ctx, string(CandStatusNew)); stageErr == nil {
		if err := s.repo.CreateStageHistory(ctx, &ApplicationStageHistory{
			ApplicationID: a.ID,
			FromStageID:   nil,
			ToStageID:     newStage.ID,
			ChangedAt:     time.Now().UnixNano(),
		}); err != nil {
			s.logger.Warn("failed to write initial stage history", zap.String("application_id", a.ID.String()), zap.Error(err))
		}
	} else {
		s.logger.Warn("failed to look up NEW stage for initial history", zap.String("application_id", a.ID.String()), zap.Error(stageErr))
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

// transitionApplicationStatus memvalidasi transisi (state machine G-5),
// menulis baris job_application_stage_histories, dan meng-update field
// status + timestamp stage pada a (in-memory) — caller bertanggung jawab
// memanggil repo.UpdateApplication untuk menyimpannya. Dipakai oleh
// UpdateApplicationStatus (manual) dan AcceptOffer (otomatis, G-3/G-4) agar
// tidak ada perubahan status yang lolos tanpa histori.
func (s *Service) transitionApplicationStatus(ctx context.Context, a *JobApplication, newStatus CandidateStatus, changedBy *uuid.UUID, notes string) error {
	from := a.Status
	if !isValidStatusTransition(from, newStatus) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, from, newStatus)
	}
	if from == newStatus {
		return nil // no-op — idempotent, tidak menulis history baru
	}

	fromStage, err := s.repo.FindStageByCode(ctx, string(from))
	if err != nil {
		return fmt.Errorf("recruitment stage lookup failed for %s: %w", from, err)
	}
	toStage, err := s.repo.FindStageByCode(ctx, string(newStatus))
	if err != nil {
		return fmt.Errorf("recruitment stage lookup failed for %s: %w", newStatus, err)
	}

	now := time.Now().UnixNano()
	hist := &ApplicationStageHistory{
		ApplicationID: a.ID,
		FromStageID:   &fromStage.ID,
		ToStageID:     toStage.ID,
		ChangedBy:     changedBy,
		Notes:         notes,
		ChangedAt:     now,
	}
	if err := s.repo.CreateStageHistory(ctx, hist); err != nil {
		return fmt.Errorf("failed to write stage history: %w", err)
	}

	a.Status = newStatus
	switch newStatus {
	case CandStatusScreened:
		a.ScreenedAt = &now
	case CandStatusShortlisted:
		a.ShortlistedAt = &now
	case CandStatusOffered:
		a.OfferedAt = &now
	case CandStatusAccepted:
		a.AcceptedAt = &now
	case CandStatusRejected:
		a.RejectedAt = &now
	case CandStatusWithdrawn:
		a.WithdrawnAt = &now
	}
	return nil
}

func (s *Service) UpdateApplicationStatus(ctx context.Context, id, status, reason, notes string, changedBy *uuid.UUID) (*ApplicationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindApplicationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	wasAlreadyAccepted := a.Status == CandStatusAccepted

	if err := s.transitionApplicationStatus(ctx, a, CandidateStatus(status), changedBy, notes); err != nil {
		return nil, err
	}
	if reason != "" {
		a.RejectionReason = reason
	}
	if notes != "" {
		a.Notes = notes
	}

	// ACCEPTED: pertahankan efek samping slots_filled existing (di luar
	// tanggung jawab transitionApplicationStatus, yang hanya urus status +
	// history). Guard wasAlreadyAccepted agar panggilan ACCEPTED berulang
	// (no-op transition) tidak double-increment slots_filled.
	if CandidateStatus(status) == CandStatusAccepted && !wasAlreadyAccepted {
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
	}

	if err := s.repo.UpdateApplication(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Application status updated", zap.String("id", a.ID.String()), zap.String("status", string(a.Status)))
	return applicationToResponse(a), nil
}

func (s *Service) GetApplicationHistory(ctx context.Context, applicationID string) ([]StageHistoryResponse, error) {
	appUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application id: %w", err)
	}
	if _, err := s.repo.FindApplicationByID(ctx, appUID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListStageHistoryByApplication(ctx, appUID)
	if err != nil {
		return nil, err
	}
	stages, err := s.repo.ListStages(ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[uuid.UUID]RecruitmentStage, len(stages))
	for _, st := range stages {
		byID[st.ID] = st
	}

	out := make([]StageHistoryResponse, 0, len(rows))
	for _, r := range rows {
		resp := StageHistoryResponse{
			ID:        r.ID.String(),
			ToStage:   StageRef{Code: byID[r.ToStageID].Code, Name: byID[r.ToStageID].Name},
			Notes:     r.Notes,
			ChangedAt: r.ChangedAt,
		}
		if r.FromStageID != nil {
			if st, ok := byID[*r.FromStageID]; ok {
				resp.FromStage = &StageRef{Code: st.Code, Name: st.Name}
			}
		}
		if r.ChangedBy != nil {
			resp.ChangedBy = r.ChangedBy.String()
		}
		out = append(out, resp)
	}
	return out, nil
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
	prevStatus := o.Status
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
	// S-7: handoff kebutuhan training HANYA saat onboarding BERTRANSIISI menjadi
	// COMPLETED (prevStatus != COMPLETED) — update berulang status COMPLETED
	// (mis. update notes pada onboarding yang sudah selesai) tidak membuat
	// TrainingNeed duplikat. Fail-safe: provider nil/error tidak menggagalkan
	// penyelesaian onboarding.
	if req.Status != nil && *req.Status == "COMPLETED" && prevStatus != "COMPLETED" {
		s.handoffOnboardingTraining(ctx, o)
	}
	return onboardingToResponse(o), nil
}

// handoffOnboardingTraining (S-7) meneruskan kebutuhan training ke Training
// module via interface narrow setelah onboarding selesai. Recruitment hanya
// menghasilkan kebutuhan (handoff) — tidak mengeksekusi training. Reason
// dikirim kosong; Training module yang mengisi default (satu sumber kebenaran
// string di training.CreateOnboardingNeed — hindari drift antar-modul).
func (s *Service) handoffOnboardingTraining(ctx context.Context, o *EmployeeOnboarding) {
	if s.trainingHandoffProvider == nil {
		s.logger.Warn("training handoff provider not wired; onboarding completed without training need",
			zap.String("onboarding_id", o.ID.String()),
			zap.String("employee_id", o.EmployeeID.String()))
		return
	}
	if err := s.trainingHandoffProvider.CreateOnboardingNeed(ctx, o.EmployeeID, o.ID, ""); err != nil {
		s.logger.Warn("training handoff failed; onboarding remains completed",
			zap.String("onboarding_id", o.ID.String()),
			zap.String("employee_id", o.EmployeeID.String()),
			zap.Error(err))
		return
	}
	s.logger.Info("Training handoff created after onboarding completion",
		zap.String("onboarding_id", o.ID.String()),
		zap.String("employee_id", o.EmployeeID.String()))
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
		Name:                 req.Name,
		Description:          req.Description,
		IsCompleted:          false,
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
// Succession Gap fallback (S-5 — CI → Recruitment)
// =========================================================================

// validateSuccessionGapFallback memvalidasi requisition reason SUCCESSION_GAP
// terhadap Career Intelligence: apakah succession_position_id benar-benar
// posisi kunci tanpa successor siap. Fail-safe: provider nil atau error tidak
// menggagalkan operasi — hanya dicatat di log; requisition tetap tersimpan
// dengan succession_position_id sebagai referensi fallback.
func (s *Service) validateSuccessionGapFallback(ctx context.Context, r *JobRequisition) {
	if r.ReasonType != string(ReqReasonSuccessionGap) {
		return
	}
	// Catatan: di create-path r.ID masih uuid.Nil (BeforeCreate berjalan di dalam
	// repo.CreateRequisition), jadi log memakai organization_id + position
	// (pola sama resolveWorkforceGapSlots S-1).
	if r.SuccessionPositionID == nil {
		s.logger.Warn("SUCCESSION_GAP requisition without succession_position_id; fallback not linked",
			zap.String("organization_id", r.OrganizationID.String()))
		return
	}
	if s.successionProvider == nil {
		s.logger.Warn("succession gap provider not wired; fallback unvalidated",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.String("succession_position_id", r.SuccessionPositionID.String()))
		return
	}
	isGap, err := s.successionProvider.SuccessionGapForPosition(ctx, *r.SuccessionPositionID)
	if err != nil {
		s.logger.Warn("succession gap provider failed; fallback unvalidated",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.String("succession_position_id", r.SuccessionPositionID.String()),
			zap.Error(err))
		return
	}
	if isGap {
		s.logger.Info("Succession gap fallback validated — external recruitment for key position",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.String("succession_position_id", r.SuccessionPositionID.String()))
	} else {
		s.logger.Warn("SUCCESSION_GAP requisition but position has a ready successor; fallback may not be needed",
			zap.String("organization_id", r.OrganizationID.String()),
			zap.String("succession_position_id", r.SuccessionPositionID.String()))
	}
}

// =========================================================================
// Internal Candidate (S-4 — CI → Recruitment)
// =========================================================================

// GetEligibleInternalCandidates mengembalikan employee internal yang eligible
// untuk sebuah target position, dibaca dari Career Intelligence via interface
// narrow (plan S-4). Recruitment tidak menghitung eligibility sendiri. Bila
// provider belum di-wire, dikembalikan list kosong (fail-safe).
func (s *Service) GetEligibleInternalCandidates(ctx context.Context, targetPositionID string) ([]InternalCandidate, error) {
	uid, err := uuid.Parse(targetPositionID)
	if err != nil {
		return nil, fmt.Errorf("invalid position_id: %w", err)
	}
	if s.internalProvider == nil {
		s.logger.Warn("internal candidate provider not wired; returning empty list",
			zap.String("target_position_id", targetPositionID))
		return []InternalCandidate{}, nil
	}
	return s.internalProvider.EligibleEmployeesForPosition(ctx, uid)
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

// generateRequisitionNumber (G-2) membuat nomor requisition otomatis dengan
// format REQ-YYYYMM-XXXXXXXX (pola nomor dokumen TRN-* training: prefix + 8
// karakter acak UUID yang di-upper — konsisten dengan training.CertificateNo).
// 8 hex char = 32 bit, cukup untuk uniqueness praktis per bulan tanpa
// sequence table.
func generateRequisitionNumber() string {
	return fmt.Sprintf("REQ-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))
}

// handoffHiredEmployee (G-4) meneruskan hasil hire dari offer yang diterima:
// kandidat EXTERNAL → Employee module membuat employee baru (referensi
// recruited_from_application_id); kandidat INTERNAL ber-employee_id →
// Employee Movement (promotion/mutation). Best-effort — provider nil atau
// error downstream hanya di-log, accept offer tidak pernah gagal karenanya.
// Dipanggil hanya saat aplikasi bertransisi ke ACCEPTED (guard di AcceptOffer)
// supaya idempotent terhadap accept offer kedua di aplikasi yang sama.
func (s *Service) handoffHiredEmployee(ctx context.Context, o *JobOffer, a *JobApplication, req *JobRequisition) {
	cand, err := s.repo.FindCandidateByID(ctx, a.CandidateID)
	if err != nil || cand == nil {
		s.logger.Warn("offer accepted but candidate lookup failed; skipping employee handoff",
			zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.Error(err))
		return
	}

	orgID, posID := "", ""
	if req != nil {
		orgID = req.OrganizationID.String()
		if req.PositionID != nil {
			posID = req.PositionID.String()
		}
	}

	if cand.CandidateType == "INTERNAL" && cand.EmployeeID != nil {
		// Internal hire — diteruskan ke Employee Movement, bukan employee baru.
		if s.movementProvider == nil {
			s.logger.Warn("movement provider not wired; skipping internal movement handoff",
				zap.String("application_id", a.ID.String()))
			return
		}
		if err := s.movementProvider.CreateHiredMovement(ctx, MovementHireInput{
			EmployeeID:     cand.EmployeeID.String(),
			ApplicationID:  a.ID.String(),
			OrganizationID: orgID,
			PositionID:     posID,
			EffectiveDate:  o.StartDate,
			Reason:         "internal hire via accepted recruitment offer",
		}); err != nil {
			s.logger.Warn("offer accepted but internal movement creation failed",
				zap.String("offer_id", o.ID.String()), zap.String("employee_id", cand.EmployeeID.String()), zap.Error(err))
			return
		}
		s.logger.Info("Internal movement created from accepted offer",
			zap.String("offer_id", o.ID.String()), zap.String("employee_id", cand.EmployeeID.String()))
		return
	}

	// External hire — Employee module membuat employee baru.
	if s.employeeProvider == nil {
		s.logger.Warn("employee provider not wired; skipping employee creation",
			zap.String("application_id", a.ID.String()))
		return
	}
	employeeID, err := s.employeeProvider.CreateHiredEmployee(ctx, EmployeeHireInput{
		ApplicationID:  a.ID.String(),
		CandidateID:    cand.ID.String(),
		Name:           strings.TrimSpace(cand.FirstName + " " + cand.LastName),
		Email:          cand.Email,
		Phone:          cand.Phone,
		StartDate:      o.StartDate,
		EmploymentType: o.EmploymentType,
		OrganizationID: orgID,
		PositionID:     posID,
		Salary:         o.Salary,
	})
	if err != nil {
		s.logger.Warn("offer accepted but employee creation failed",
			zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.Error(err))
		return
	}
	s.logger.Info("Employee created from accepted offer",
		zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.String("employee_id", employeeID))
}

// applyCandidateTypeFields (G-4) mengatur candidate_type (default EXTERNAL)
// dan employee_id (referensi employee untuk kandidat INTERNAL) pada kandidat.
// Bila type diubah ke EXTERNAL, referensi employee_id dibersihkan — jalur
// handoff offer (G-4) ditentukan oleh candidate_type, bukan employee_id.
func applyCandidateTypeFields(c *Candidate, candidateType, employeeID *string) error {
	if candidateType != nil && *candidateType != "" {
		c.CandidateType = *candidateType
		if *candidateType == "EXTERNAL" {
			c.EmployeeID = nil
		}
	}
	if c.CandidateType == "" {
		c.CandidateType = "EXTERNAL"
	}
	if employeeID != nil && *employeeID != "" {
		id, err := uuid.Parse(*employeeID)
		if err != nil {
			return fmt.Errorf("invalid employee_id: %w", err)
		}
		c.EmployeeID = &id
	}
	return nil
}

// generateOfferNumber (G-3) membuat nomor offer otomatis dengan format
// OFF-YYYYMM-XXXXXXXX (pola sama generateRequisitionNumber G-2).
func generateOfferNumber() string {
	return fmt.Sprintf("OFF-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))
}

// generateCandidateNumber (G-6) membuat nomor kandidat otomatis dengan
// format CAND-YYYYMM-XXXXXXXX (pola sama generateRequisitionNumber G-2 /
// generateOfferNumber G-3).
func generateCandidateNumber() string {
	return fmt.Sprintf("CAND-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))
}

// isOfferExpired (G-3) membandingkan expiry_date (YYYY-MM-DD) dengan hari ini.
// Expiry kosong = tidak ada batas (tidak pernah expired).
func isOfferExpired(expiryDate string) bool {
	if expiryDate == "" {
		return false
	}
	expiry, err := time.Parse("2006-01-02", expiryDate)
	if err != nil {
		// Format tak dikenal — jangan blokir accept (fail-open).
		return false
	}
	today := time.Now().Format("2006-01-02")
	todayParsed, _ := time.Parse("2006-01-02", today)
	return expiry.Before(todayParsed)
}

// =========================================================================
// Response converters
// =========================================================================

func requisitionToResponse(r *JobRequisition) *RequisitionResponse {
	resp := &RequisitionResponse{
		ID:                r.ID.String(),
		OrganizationID:    r.OrganizationID.String(),
		Title:             r.Title,
		RequisitionNumber: r.RequisitionNumber,
		Priority:          string(r.Priority),
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
	if r.PositionID != nil {
		resp.PositionID = r.PositionID.String()
	}
	// GORM read-back default:0 membuat OpenedAt pointer non-nil bernilai 0 pada
	// draft — perlakukan 0 sebagai belum dibuka (omit), konsisten dengan
	// ClosedAt yang json:"-".
	if r.OpenedAt != nil && *r.OpenedAt != 0 {
		resp.OpenedAt = r.OpenedAt
	}
	if r.RequestedBy != nil {
		resp.RequestedBy = r.RequestedBy.String()
	}
	if r.ApprovedBy != nil {
		resp.ApprovedBy = r.ApprovedBy.String()
	}
	if r.ApprovalInstanceID != nil {
		resp.ApprovalInstanceID = r.ApprovalInstanceID.String()
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
	if r.SuccessionPositionID != nil {
		resp.SuccessionPositionID = r.SuccessionPositionID.String()
	}
	if r.TargetStartDate != nil {
		resp.TargetStartDate = *r.TargetStartDate
	}
	return resp
}

func offerToResponse(o *JobOffer) *OfferResponse {
	resp := &OfferResponse{
		ID:             o.ID.String(),
		ApplicationID:  o.ApplicationID.String(),
		OfferNumber:    o.OfferNumber,
		EmploymentType: o.EmploymentType,
		Salary:         o.Salary,
		Allowances:     o.Allowances,
		Benefits:       o.Benefits,
		StartDate:      o.StartDate,
		ExpiryDate:     o.ExpiryDate,
		Status:         string(o.Status),
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
	if o.ApprovalInstanceID != nil {
		resp.ApprovalInstanceID = o.ApprovalInstanceID.String()
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
		// G-4: jenis kandidat + referensi employee (internal hire).
		CandidateType:   c.CandidateType,
		CandidateNumber: c.CandidateNumber,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
	if c.EmployeeID != nil {
		resp.EmployeeID = c.EmployeeID.String()
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

// =========================================================================
// Candidate Educations (G-6)
// =========================================================================

func (s *Service) CreateCandidateEducation(ctx context.Context, candidateID string, req CreateCandidateEducationRequest) (*CandidateEducationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	e := &CandidateEducation{
		CandidateID:     candUUID,
		InstitutionName: req.InstitutionName,
		GPA:             req.GPA,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		IsHighest:       req.IsHighest,
		Notes:           req.Notes,
	}
	if req.EducationID != nil && *req.EducationID != "" {
		eduID, err := uuid.Parse(*req.EducationID)
		if err != nil {
			return nil, fmt.Errorf("invalid education_id: %w", err)
		}
		e.EducationID = &eduID
	}
	if req.EducationMajorID != nil && *req.EducationMajorID != "" {
		majorID, err := uuid.Parse(*req.EducationMajorID)
		if err != nil {
			return nil, fmt.Errorf("invalid education_major_id: %w", err)
		}
		e.EducationMajorID = &majorID
	}
	if req.Major != nil {
		e.Major = req.Major
	}

	if err := s.repo.CreateCandidateEducation(ctx, e); err != nil {
		return nil, err
	}
	created, err := s.repo.FindCandidateEducationByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	return candidateEducationToResponse(created), nil
}

func (s *Service) ListCandidateEducations(ctx context.Context, candidateID string) ([]CandidateEducationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateEducations(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateEducationResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateEducationToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateEducation(ctx context.Context, id string, req UpdateCandidateEducationRequest) (*CandidateEducationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindCandidateEducationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.EducationID != nil {
		if *req.EducationID == "" {
			e.EducationID = nil
			e.Education = nil // drop preloaded assoc so GORM doesn't resurrect the FK from it on Save
		} else {
			eduID, err := uuid.Parse(*req.EducationID)
			if err != nil {
				return nil, fmt.Errorf("invalid education_id: %w", err)
			}
			e.EducationID = &eduID
			e.Education = nil // also drop here — let the re-fetch after save reload it fresh, don't risk a stale/mismatched pointer
		}
	}
	if req.InstitutionName != nil {
		e.InstitutionName = *req.InstitutionName
	}
	if req.EducationMajorID != nil {
		if *req.EducationMajorID == "" {
			e.EducationMajorID = nil
			e.EducationMajor = nil // drop preloaded assoc so GORM doesn't resurrect the FK from it on Save
		} else {
			majorID, err := uuid.Parse(*req.EducationMajorID)
			if err != nil {
				return nil, fmt.Errorf("invalid education_major_id: %w", err)
			}
			e.EducationMajorID = &majorID
			e.EducationMajor = nil // also drop here — let the re-fetch after save reload it fresh, don't risk a stale/mismatched pointer
		}
	}
	if req.Major != nil {
		e.Major = req.Major
	}
	if req.GPA != nil {
		e.GPA = req.GPA
	}
	if req.StartYear != nil {
		e.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		e.EndYear = req.EndYear
	}
	if req.IsHighest != nil {
		e.IsHighest = *req.IsHighest
	}
	if req.Notes != nil {
		e.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateEducation(ctx, e); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindCandidateEducationByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	return candidateEducationToResponse(updated), nil
}

func (s *Service) DeleteCandidateEducation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateEducation(ctx, uid)
}

func candidateEducationToResponse(e *CandidateEducation) *CandidateEducationResponse {
	resp := &CandidateEducationResponse{
		ID:              e.ID.String(),
		CandidateID:     e.CandidateID.String(),
		InstitutionName: e.InstitutionName,
		IsHighest:       e.IsHighest,
		Notes:           e.Notes,
	}
	if e.EducationID != nil {
		resp.EducationID = e.EducationID.String()
	}
	if e.EducationMajorID != nil {
		resp.EducationMajorID = e.EducationMajorID.String()
	}
	if e.EducationMajor != nil {
		resp.MajorName = e.EducationMajor.Name
	}
	if e.Major != nil {
		resp.Major = *e.Major
	}
	if e.GPA != nil {
		resp.GPA = *e.GPA
	}
	if e.StartYear != nil {
		resp.StartYear = *e.StartYear
	}
	if e.EndYear != nil {
		resp.EndYear = *e.EndYear
	}
	return resp
}

// =========================================================================
// Candidate Work Experiences (G-6)
// =========================================================================

func (s *Service) CreateCandidateWorkExperience(ctx context.Context, candidateID string, req CreateCandidateWorkExperienceRequest) (*CandidateWorkExperienceResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	e := &CandidateWorkExperience{
		CandidateID:    candUUID,
		CompanyName:    req.CompanyName,
		JobTitle:       req.JobTitle,
		EmploymentType: req.EmploymentType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsCurrent:      req.IsCurrent,
		Description:    req.Description,
	}
	if err := s.repo.CreateCandidateWorkExperience(ctx, e); err != nil {
		return nil, err
	}
	return candidateWorkExperienceToResponse(e), nil
}

func (s *Service) ListCandidateWorkExperiences(ctx context.Context, candidateID string) ([]CandidateWorkExperienceResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateWorkExperiences(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateWorkExperienceResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateWorkExperienceToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateWorkExperience(ctx context.Context, id string, req UpdateCandidateWorkExperienceRequest) (*CandidateWorkExperienceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindCandidateWorkExperienceByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CompanyName != nil {
		e.CompanyName = *req.CompanyName
	}
	if req.JobTitle != nil {
		e.JobTitle = *req.JobTitle
	}
	if req.EmploymentType != nil {
		e.EmploymentType = req.EmploymentType
	}
	if req.StartDate != nil {
		e.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		e.EndDate = req.EndDate
	}
	if req.IsCurrent != nil {
		e.IsCurrent = *req.IsCurrent
	}
	if req.Description != nil {
		e.Description = *req.Description
	}
	if err := s.repo.UpdateCandidateWorkExperience(ctx, e); err != nil {
		return nil, err
	}
	return candidateWorkExperienceToResponse(e), nil
}

func (s *Service) DeleteCandidateWorkExperience(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateWorkExperience(ctx, uid)
}

func candidateWorkExperienceToResponse(e *CandidateWorkExperience) *CandidateWorkExperienceResponse {
	resp := &CandidateWorkExperienceResponse{
		ID:          e.ID.String(),
		CandidateID: e.CandidateID.String(),
		CompanyName: e.CompanyName,
		JobTitle:    e.JobTitle,
		StartDate:   e.StartDate,
		IsCurrent:   e.IsCurrent,
		Description: e.Description,
	}
	if e.EmploymentType != nil {
		resp.EmploymentType = *e.EmploymentType
	}
	if e.EndDate != nil {
		resp.EndDate = *e.EndDate
	}
	return resp
}

// =========================================================================
// Candidate Skills (G-6)
// =========================================================================

func (s *Service) CreateCandidateSkill(ctx context.Context, candidateID string, req CreateCandidateSkillRequest) (*CandidateSkillResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	compUUID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	if _, err := s.repo.FindCompetencyByID(ctx, compUUID); err != nil {
		return nil, fmt.Errorf("competency not found: %w", err)
	}

	sk := &CandidateSkill{
		CandidateID:  candUUID,
		CompetencyID: compUUID,
		Level:        req.Level,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateCandidateSkill(ctx, sk); err != nil {
		return nil, err
	}
	created, err := s.repo.FindCandidateSkillByID(ctx, sk.ID)
	if err != nil {
		return nil, err
	}
	return candidateSkillToResponse(created), nil
}

func (s *Service) ListCandidateSkills(ctx context.Context, candidateID string) ([]CandidateSkillResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateSkills(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateSkillResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateSkillToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateSkill(ctx context.Context, id string, req UpdateCandidateSkillRequest) (*CandidateSkillResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sk, err := s.repo.FindCandidateSkillByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Level != nil {
		sk.Level = req.Level
	}
	if req.Notes != nil {
		sk.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateSkill(ctx, sk); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindCandidateSkillByID(ctx, sk.ID)
	if err != nil {
		return nil, err
	}
	return candidateSkillToResponse(updated), nil
}

func (s *Service) DeleteCandidateSkill(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateSkill(ctx, uid)
}

func candidateSkillToResponse(s *CandidateSkill) *CandidateSkillResponse {
	resp := &CandidateSkillResponse{
		ID:           s.ID.String(),
		CandidateID:  s.CandidateID.String(),
		CompetencyID: s.CompetencyID.String(),
		Notes:        s.Notes,
	}
	if s.Competency != nil {
		resp.CompetencyName = s.Competency.Name
	}
	if s.Level != nil {
		resp.Level = *s.Level
	}
	return resp
}

// =========================================================================
// Candidate Certifications (G-6)
// =========================================================================

func (s *Service) CreateCandidateCertification(ctx context.Context, candidateID string, req CreateCandidateCertificationRequest) (*CandidateCertificationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	c := &CandidateCertification{
		CandidateID:         candUUID,
		Name:                req.Name,
		IssuingOrganization: req.IssuingOrganization,
		IssueDate:           req.IssueDate,
		ExpiryDate:          req.ExpiryDate,
		CredentialURL:       req.CredentialURL,
		Notes:               req.Notes,
	}
	if err := s.repo.CreateCandidateCertification(ctx, c); err != nil {
		return nil, err
	}
	return candidateCertificationToResponse(c), nil
}

func (s *Service) ListCandidateCertifications(ctx context.Context, candidateID string) ([]CandidateCertificationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateCertifications(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateCertificationResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateCertificationToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateCertification(ctx context.Context, id string, req UpdateCandidateCertificationRequest) (*CandidateCertificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCandidateCertificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.IssuingOrganization != nil {
		c.IssuingOrganization = req.IssuingOrganization
	}
	if req.IssueDate != nil {
		c.IssueDate = req.IssueDate
	}
	if req.ExpiryDate != nil {
		c.ExpiryDate = req.ExpiryDate
	}
	if req.CredentialURL != nil {
		c.CredentialURL = req.CredentialURL
	}
	if req.Notes != nil {
		c.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateCertification(ctx, c); err != nil {
		return nil, err
	}
	return candidateCertificationToResponse(c), nil
}

func (s *Service) DeleteCandidateCertification(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateCertification(ctx, uid)
}

func candidateCertificationToResponse(c *CandidateCertification) *CandidateCertificationResponse {
	resp := &CandidateCertificationResponse{
		ID:          c.ID.String(),
		CandidateID: c.CandidateID.String(),
		Name:        c.Name,
		Notes:       c.Notes,
	}
	if c.IssuingOrganization != nil {
		resp.IssuingOrganization = *c.IssuingOrganization
	}
	if c.IssueDate != nil {
		resp.IssueDate = *c.IssueDate
	}
	if c.ExpiryDate != nil {
		resp.ExpiryDate = *c.ExpiryDate
	}
	if c.CredentialURL != nil {
		resp.CredentialURL = *c.CredentialURL
	}
	return resp
}

// =========================================================================
// Candidate Documents (G-6)
// =========================================================================

func (s *Service) CreateCandidateDocument(ctx context.Context, candidateID string, req CreateCandidateDocumentRequest) (*CandidateDocumentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	docType := req.DocumentType
	if docType == "" {
		docType = "OTHER"
	}
	d := &CandidateDocument{
		CandidateID:  candUUID,
		DocumentType: docType,
		Name:         req.Name,
		FileURL:      req.FileURL,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateCandidateDocument(ctx, d); err != nil {
		return nil, err
	}
	return candidateDocumentToResponse(d), nil
}

func (s *Service) ListCandidateDocuments(ctx context.Context, candidateID string) ([]CandidateDocumentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateDocuments(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateDocumentResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateDocumentToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateDocument(ctx context.Context, id string, req UpdateCandidateDocumentRequest) (*CandidateDocumentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	d, err := s.repo.FindCandidateDocumentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.DocumentType != nil {
		d.DocumentType = *req.DocumentType
	}
	if req.Name != nil {
		d.Name = *req.Name
	}
	if req.FileURL != nil {
		d.FileURL = *req.FileURL
	}
	if req.Notes != nil {
		d.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateDocument(ctx, d); err != nil {
		return nil, err
	}
	return candidateDocumentToResponse(d), nil
}

func (s *Service) DeleteCandidateDocument(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateDocument(ctx, uid)
}

func candidateDocumentToResponse(d *CandidateDocument) *CandidateDocumentResponse {
	return &CandidateDocumentResponse{
		ID:           d.ID.String(),
		CandidateID:  d.CandidateID.String(),
		DocumentType: d.DocumentType,
		Name:         d.Name,
		FileURL:      d.FileURL,
		Notes:        d.Notes,
	}
}

// =========================================================================
// Candidate Consents (G-6) — append-only, no Update/Delete
// =========================================================================

func (s *Service) CreateCandidateConsent(ctx context.Context, candidateID string, req CreateCandidateConsentRequest, changedBy *uuid.UUID) (*CandidateConsentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	c := &CandidateConsent{
		CandidateID: candUUID,
		Action:      req.Action,
		Notes:       req.Notes,
		ChangedBy:   changedBy,
		ChangedAt:   time.Now().UnixNano(),
	}
	if err := s.repo.CreateCandidateConsent(ctx, c); err != nil {
		return nil, err
	}
	return candidateConsentToResponse(c), nil
}

func (s *Service) ListCandidateConsents(ctx context.Context, candidateID string) ([]CandidateConsentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateConsents(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateConsentResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateConsentToResponse(&list[i]))
	}
	return out, nil
}

func candidateConsentToResponse(c *CandidateConsent) *CandidateConsentResponse {
	resp := &CandidateConsentResponse{
		ID:          c.ID.String(),
		CandidateID: c.CandidateID.String(),
		Action:      c.Action,
		Notes:       c.Notes,
		ChangedAt:   c.ChangedAt,
	}
	if c.ChangedBy != nil {
		resp.ChangedBy = c.ChangedBy.String()
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
		Name:                 t.Name,
		Description:          t.Description,
		IsCompleted:          t.IsCompleted,
		Notes:                t.Notes,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
	if t.TemplateID != nil {
		resp.TemplateID = t.TemplateID.String()
	}
	if t.AssignedTo != nil {
		resp.AssignedTo = t.AssignedTo.String()
	}
	return resp
}

// =========================================================================
// Application Screenings (G-7 sub-project 1)
// =========================================================================

func (s *Service) CreateApplicationScreening(ctx context.Context, applicationID string, req CreateApplicationScreeningRequest) (*ApplicationScreeningResponse, error) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	if _, err := s.repo.FindApplicationByID(ctx, appUUID); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	result := req.Result
	if result == "" {
		result = "HOLD"
	}
	var screenedBy *uuid.UUID
	if req.ScreenedBy != "" {
		id, err := uuid.Parse(req.ScreenedBy)
		if err != nil {
			return nil, fmt.Errorf("invalid screened_by: %w", err)
		}
		screenedBy = &id
	}
	sc := &ApplicationScreening{
		ApplicationID: appUUID,
		ScreenedBy:    screenedBy,
		ScreenedAt:    req.ScreenedAt,
		Score:         req.Score,
		Result:        result,
		Notes:         req.Notes,
	}
	if err := s.repo.CreateApplicationScreening(ctx, sc); err != nil {
		return nil, err
	}
	return applicationScreeningToResponse(sc), nil
}

func (s *Service) ListApplicationScreenings(ctx context.Context, applicationID string) ([]ApplicationScreeningResponse, error) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	list, err := s.repo.ListApplicationScreenings(ctx, appUUID)
	if err != nil {
		return nil, err
	}
	out := make([]ApplicationScreeningResponse, 0, len(list))
	for i := range list {
		out = append(out, *applicationScreeningToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateApplicationScreening(ctx context.Context, id string, req UpdateApplicationScreeningRequest) (*ApplicationScreeningResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindApplicationScreeningByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.ScreenedBy != nil {
		if *req.ScreenedBy == "" {
			sc.ScreenedBy = nil
		} else {
			id, err := uuid.Parse(*req.ScreenedBy)
			if err != nil {
				return nil, fmt.Errorf("invalid screened_by: %w", err)
			}
			sc.ScreenedBy = &id
		}
	}
	if req.ScreenedAt != nil {
		sc.ScreenedAt = *req.ScreenedAt
	}
	if req.Score != nil {
		sc.Score = req.Score
	}
	if req.Result != nil {
		sc.Result = *req.Result
	}
	if req.Notes != nil {
		sc.Notes = *req.Notes
	}
	if err := s.repo.UpdateApplicationScreening(ctx, sc); err != nil {
		return nil, err
	}
	return applicationScreeningToResponse(sc), nil
}

func (s *Service) DeleteApplicationScreening(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteApplicationScreening(ctx, uid)
}

// =========================================================================
// Recruitment Assessments + Participants (G-7 sub-project 2)
// =========================================================================

func (s *Service) CreateAssessment(ctx context.Context, req CreateAssessmentRequest) (*AssessmentResponse, error) {
	var reqUUID *uuid.UUID
	if req.RequisitionID != "" {
		id, err := uuid.Parse(req.RequisitionID)
		if err != nil {
			return nil, fmt.Errorf("invalid requisition_id: %w", err)
		}
		if _, err := s.repo.FindRequisitionByID(ctx, id); err != nil {
			return nil, fmt.Errorf("requisition not found: %w", err)
		}
		reqUUID = &id
	}
	assessType := req.Type
	if assessType == "" {
		assessType = "OTHER"
	}
	var meetingLink *string
	if req.MeetingLink != "" {
		meetingLink = &req.MeetingLink
	}
	a := &RecruitmentAssessment{
		RequisitionID: reqUUID,
		Name:          req.Name,
		Type:          assessType,
		ScheduledAt:   req.ScheduledAt,
		Location:      req.Location,
		MeetingLink:   meetingLink,
		Notes:         req.Notes,
	}
	if err := s.repo.CreateAssessment(ctx, a); err != nil {
		return nil, err
	}
	return assessmentToResponse(a), nil
}

func (s *Service) ListAssessments(ctx context.Context) ([]AssessmentResponse, error) {
	list, err := s.repo.ListAssessments(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]AssessmentResponse, 0, len(list))
	for i := range list {
		out = append(out, *assessmentToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) GetAssessmentByID(ctx context.Context, id string) (*AssessmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindAssessmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return assessmentToResponse(a), nil
}

func (s *Service) UpdateAssessment(ctx context.Context, id string, req UpdateAssessmentRequest) (*AssessmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindAssessmentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		a.Name = *req.Name
	}
	if req.Type != nil {
		a.Type = *req.Type
	}
	if req.ScheduledAt != nil {
		a.ScheduledAt = *req.ScheduledAt
	}
	if req.Location != nil {
		a.Location = *req.Location
	}
	if req.MeetingLink != nil {
		a.MeetingLink = req.MeetingLink
	}
	if req.Notes != nil {
		a.Notes = *req.Notes
	}
	if err := s.repo.UpdateAssessment(ctx, a); err != nil {
		return nil, err
	}
	return assessmentToResponse(a), nil
}

func (s *Service) DeleteAssessment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteAssessment(ctx, uid)
}

func (s *Service) AddAssessmentParticipant(ctx context.Context, assessmentID string, req AddAssessmentParticipantRequest) (*AssessmentParticipantResponse, error) {
	assessUUID, err := uuid.Parse(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment_id: %w", err)
	}
	if _, err := s.repo.FindAssessmentByID(ctx, assessUUID); err != nil {
		return nil, err
	}
	appUUID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	if _, err := s.repo.FindApplicationByID(ctx, appUUID); err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}
	p := &AssessmentParticipant{
		AssessmentID:  assessUUID,
		ApplicationID: appUUID,
		Status:        "INVITED",
	}
	if err := s.repo.CreateAssessmentParticipant(ctx, p); err != nil {
		return nil, err
	}
	return assessmentParticipantToResponse(p), nil
}

func (s *Service) ListAssessmentParticipants(ctx context.Context, assessmentID string) ([]AssessmentParticipantResponse, error) {
	assessUUID, err := uuid.Parse(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid assessment_id: %w", err)
	}
	list, err := s.repo.ListAssessmentParticipants(ctx, assessUUID)
	if err != nil {
		return nil, err
	}
	out := make([]AssessmentParticipantResponse, 0, len(list))
	for i := range list {
		out = append(out, *assessmentParticipantToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateAssessmentParticipant(ctx context.Context, id string, req UpdateAssessmentParticipantRequest) (*AssessmentParticipantResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	p, err := s.repo.FindAssessmentParticipantByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.Score != nil {
		p.Score = req.Score
	}
	if req.Result != nil {
		p.Result = req.Result
	}
	if req.Recommendation != nil {
		p.Recommendation = *req.Recommendation
	}
	if err := s.repo.UpdateAssessmentParticipant(ctx, p); err != nil {
		return nil, err
	}
	return assessmentParticipantToResponse(p), nil
}

func (s *Service) DeleteAssessmentParticipant(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteAssessmentParticipant(ctx, uid)
}

// =========================================================================
// Job Requisition Requirements + Competencies (G-9 sub-project 1)
// =========================================================================

func (s *Service) CreateRequisitionRequirement(ctx context.Context, requisitionID string, req CreateRequisitionRequirementRequest) (*RequisitionRequirementResponse, error) {
	reqUUID, err := uuid.Parse(requisitionID)
	if err != nil {
		return nil, fmt.Errorf("invalid requisition_id: %w", err)
	}
	if _, err := s.repo.FindRequisitionByID(ctx, reqUUID); err != nil {
		return nil, fmt.Errorf("requisition not found: %w", err)
	}
	isRequired := true
	if req.IsRequired != nil {
		isRequired = *req.IsRequired
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	r := &JobRequisitionRequirement{
		RequisitionID:   reqUUID,
		RequirementType: req.RequirementType,
		Name:            req.Name,
		Description:     req.Description,
		MinimumValue:    req.MinimumValue,
		MaximumValue:    req.MaximumValue,
		IsRequired:      isRequired,
		SortOrder:       sortOrder,
	}
	if err := s.repo.CreateRequisitionRequirement(ctx, r); err != nil {
		return nil, err
	}
	return requisitionRequirementToResponse(r), nil
}

func (s *Service) ListRequisitionRequirements(ctx context.Context, requisitionID string) ([]RequisitionRequirementResponse, error) {
	reqUUID, err := uuid.Parse(requisitionID)
	if err != nil {
		return nil, fmt.Errorf("invalid requisition_id: %w", err)
	}
	list, err := s.repo.ListRequisitionRequirements(ctx, reqUUID)
	if err != nil {
		return nil, err
	}
	out := make([]RequisitionRequirementResponse, 0, len(list))
	for i := range list {
		out = append(out, *requisitionRequirementToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateRequisitionRequirement(ctx context.Context, id string, req UpdateRequisitionRequirementRequest) (*RequisitionRequirementResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	r, err := s.repo.FindRequisitionRequirementByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.RequirementType != nil {
		r.RequirementType = *req.RequirementType
	}
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.Description != nil {
		r.Description = *req.Description
	}
	if req.MinimumValue != nil {
		r.MinimumValue = req.MinimumValue
	}
	if req.MaximumValue != nil {
		r.MaximumValue = req.MaximumValue
	}
	if req.IsRequired != nil {
		r.IsRequired = *req.IsRequired
	}
	if req.SortOrder != nil {
		r.SortOrder = *req.SortOrder
	}
	if err := s.repo.UpdateRequisitionRequirement(ctx, r); err != nil {
		return nil, err
	}
	return requisitionRequirementToResponse(r), nil
}

func (s *Service) DeleteRequisitionRequirement(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteRequisitionRequirement(ctx, uid)
}

func (s *Service) CreateRequisitionCompetency(ctx context.Context, requisitionID string, req CreateRequisitionCompetencyRequest) (*RequisitionCompetencyResponse, error) {
	reqUUID, err := uuid.Parse(requisitionID)
	if err != nil {
		return nil, fmt.Errorf("invalid requisition_id: %w", err)
	}
	if _, err := s.repo.FindRequisitionByID(ctx, reqUUID); err != nil {
		return nil, fmt.Errorf("requisition not found: %w", err)
	}
	compUUID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	if _, err := s.repo.FindCompetencyByID(ctx, compUUID); err != nil {
		return nil, fmt.Errorf("competency not found: %w", err)
	}
	isRequired := true
	if req.IsRequired != nil {
		isRequired = *req.IsRequired
	}
	c := &JobRequisitionCompetency{
		RequisitionID: reqUUID,
		CompetencyID:  compUUID,
		RequiredLevel: req.RequiredLevel,
		IsRequired:    isRequired,
		Weight:        req.Weight,
	}
	if err := s.repo.CreateRequisitionCompetency(ctx, c); err != nil {
		return nil, err
	}
	return requisitionCompetencyToResponse(c), nil
}

func (s *Service) ListRequisitionCompetencies(ctx context.Context, requisitionID string) ([]RequisitionCompetencyResponse, error) {
	reqUUID, err := uuid.Parse(requisitionID)
	if err != nil {
		return nil, fmt.Errorf("invalid requisition_id: %w", err)
	}
	list, err := s.repo.ListRequisitionCompetencies(ctx, reqUUID)
	if err != nil {
		return nil, err
	}
	out := make([]RequisitionCompetencyResponse, 0, len(list))
	for i := range list {
		out = append(out, *requisitionCompetencyToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateRequisitionCompetency(ctx context.Context, id string, req UpdateRequisitionCompetencyRequest) (*RequisitionCompetencyResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindRequisitionCompetencyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.RequiredLevel != nil {
		c.RequiredLevel = req.RequiredLevel
	}
	if req.IsRequired != nil {
		c.IsRequired = *req.IsRequired
	}
	if req.Weight != nil {
		c.Weight = req.Weight
	}
	if err := s.repo.UpdateRequisitionCompetency(ctx, c); err != nil {
		return nil, err
	}
	return requisitionCompetencyToResponse(c), nil
}

func (s *Service) DeleteRequisitionCompetency(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteRequisitionCompetency(ctx, uid)
}

func requisitionRequirementToResponse(r *JobRequisitionRequirement) *RequisitionRequirementResponse {
	resp := &RequisitionRequirementResponse{
		ID:              r.ID.String(),
		RequisitionID:   r.RequisitionID.String(),
		RequirementType: r.RequirementType,
		Name:            r.Name,
		Description:     r.Description,
		IsRequired:      r.IsRequired,
		SortOrder:       r.SortOrder,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
	if r.MinimumValue != nil {
		resp.MinimumValue = *r.MinimumValue
	}
	if r.MaximumValue != nil {
		resp.MaximumValue = *r.MaximumValue
	}
	return resp
}

func requisitionCompetencyToResponse(c *JobRequisitionCompetency) *RequisitionCompetencyResponse {
	resp := &RequisitionCompetencyResponse{
		ID:            c.ID.String(),
		RequisitionID: c.RequisitionID.String(),
		CompetencyID:  c.CompetencyID.String(),
		IsRequired:    c.IsRequired,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
	if c.RequiredLevel != nil {
		resp.RequiredLevel = *c.RequiredLevel
	}
	if c.Weight != nil {
		resp.Weight = *c.Weight
	}
	return resp
}

// =========================================================================
// Candidate Match Score (G-9 sub-project 2)
// =========================================================================
// Advisory-only: dihitung on-the-fly, tidak disimpan/dipersist, tidak
// menentukan keputusan apapun secara otomatis — recruiter tetap yang
// memutuskan. Formula: Σ(weight × min(candidate_level/required_level, 1))
// / Σ(weight) × 100, dibatasi ke job_requisition_competencies (bukan
// requirement/education/experience/assessment/interview — keputusan
// brainstorming: skala berbeda-beda dan sulit dibandingkan apple-to-apple).

func (s *Service) GetCandidateMatchScore(ctx context.Context, applicationID string) (*MatchScoreResponse, error) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application_id: %w", err)
	}
	app, err := s.repo.FindApplicationByID(ctx, appUUID)
	if err != nil {
		return nil, err
	}

	resp := &MatchScoreResponse{
		ApplicationID: app.ID.String(),
		CandidateID:   app.CandidateID.String(),
		RequisitionID: app.RequisitionID.String(),
		Breakdown:     []MatchScoreCompetencyBreakdown{},
	}

	reqCompetencies, err := s.repo.ListRequisitionCompetencies(ctx, app.RequisitionID)
	if err != nil {
		return nil, err
	}
	if len(reqCompetencies) == 0 {
		resp.Note = "requisition has no competencies defined; nothing to match"
		return resp, nil
	}

	candSkills, err := s.repo.ListCandidateSkills(ctx, app.CandidateID)
	if err != nil {
		return nil, err
	}
	skillByCompetency := make(map[uuid.UUID]int, len(candSkills))
	for _, sk := range candSkills {
		if sk.Level != nil {
			skillByCompetency[sk.CompetencyID] = *sk.Level
		}
	}

	var weightedSum, weightTotal float64
	for _, rc := range reqCompetencies {
		requiredLevel := 1
		if rc.RequiredLevel != nil {
			requiredLevel = *rc.RequiredLevel
		}
		weight := 1.0
		if rc.Weight != nil {
			weight = *rc.Weight
		}
		candidateLevel := skillByCompetency[rc.CompetencyID]

		ratio := 0.0
		if requiredLevel > 0 {
			ratio = float64(candidateLevel) / float64(requiredLevel)
		}
		if ratio > 1 {
			ratio = 1
		}
		contribution := weight * ratio

		weightedSum += contribution
		weightTotal += weight

		compName := ""
		if rc.Competency != nil {
			compName = rc.Competency.Name
		}
		resp.Breakdown = append(resp.Breakdown, MatchScoreCompetencyBreakdown{
			CompetencyID:   rc.CompetencyID.String(),
			CompetencyName: compName,
			RequiredLevel:  requiredLevel,
			CandidateLevel: candidateLevel,
			Weight:         weight,
			Contribution:   contribution,
		})
	}

	if weightTotal > 0 {
		score := weightedSum / weightTotal * 100
		resp.Score = &score
	}
	return resp, nil
}

// =========================================================================
// Interviewers + Scorecard Items (G-8)
// =========================================================================

func (s *Service) AddInterviewer(ctx context.Context, interviewID string, req AddInterviewerRequest) (*InterviewerResponse, error) {
	intID, err := uuid.Parse(interviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid interview_id: %w", err)
	}
	if _, err := s.repo.FindInterviewByID(ctx, intID); err != nil {
		return nil, err
	}
	empID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee_id: %w", err)
	}
	i := &Interviewer{
		InterviewID: intID,
		EmployeeID:  empID,
		Role:        req.Role,
	}
	if err := s.repo.CreateInterviewer(ctx, i); err != nil {
		return nil, err
	}
	return interviewerToResponse(i), nil
}

func (s *Service) ListInterviewers(ctx context.Context, interviewID string) ([]InterviewerResponse, error) {
	intID, err := uuid.Parse(interviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid interview_id: %w", err)
	}
	list, err := s.repo.ListInterviewers(ctx, intID)
	if err != nil {
		return nil, err
	}
	out := make([]InterviewerResponse, 0, len(list))
	for i := range list {
		out = append(out, *interviewerToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) RemoveInterviewer(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteInterviewer(ctx, uid)
}

func (s *Service) AddScorecardItem(ctx context.Context, interviewID string, req AddScorecardItemRequest) (*ScorecardItemResponse, error) {
	intID, err := uuid.Parse(interviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid interview_id: %w", err)
	}
	if _, err := s.repo.FindInterviewByID(ctx, intID); err != nil {
		return nil, err
	}
	item := &InterviewScorecardItem{
		InterviewID: intID,
		Criterion:   req.Criterion,
		Weight:      req.Weight,
		Score:       req.Score,
		Notes:       req.Notes,
	}
	if err := s.repo.CreateScorecardItem(ctx, item); err != nil {
		return nil, err
	}
	return scorecardItemToResponse(item), nil
}

func (s *Service) ListScorecardItems(ctx context.Context, interviewID string) ([]ScorecardItemResponse, error) {
	intID, err := uuid.Parse(interviewID)
	if err != nil {
		return nil, fmt.Errorf("invalid interview_id: %w", err)
	}
	list, err := s.repo.ListScorecardItems(ctx, intID)
	if err != nil {
		return nil, err
	}
	out := make([]ScorecardItemResponse, 0, len(list))
	for i := range list {
		out = append(out, *scorecardItemToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateScorecardItem(ctx context.Context, id string, req UpdateScorecardItemRequest) (*ScorecardItemResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	item, err := s.repo.FindScorecardItemByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Criterion != nil {
		item.Criterion = *req.Criterion
	}
	if req.Weight != nil {
		item.Weight = *req.Weight
	}
	if req.Score != nil {
		item.Score = req.Score
	}
	if req.Notes != nil {
		item.Notes = *req.Notes
	}
	if err := s.repo.UpdateScorecardItem(ctx, item); err != nil {
		return nil, err
	}
	return scorecardItemToResponse(item), nil
}

func (s *Service) DeleteScorecardItem(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteScorecardItem(ctx, uid)
}

// CompleteInterview menutup interview (status -> COMPLETED, completed_at)
// dan menghitung Interview.Score sebagai weighted average dari
// interview_scorecard_items (Σ(score×weight)/Σ(weight); item tanpa score
// dilewati). Bila tidak ada scorecard item berskor, Score tidak diubah
// (recruiter tetap bisa mengisi manual seperti alur existing).
func (s *Service) CompleteInterview(ctx context.Context, id string) (*InterviewResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	iv, err := s.repo.FindInterviewByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListScorecardItems(ctx, uid)
	if err != nil {
		return nil, err
	}
	var weightedSum, weightTotal float64
	for _, item := range items {
		if item.Score == nil {
			continue
		}
		weightedSum += *item.Score * item.Weight
		weightTotal += item.Weight
	}
	if weightTotal > 0 {
		aggregate := weightedSum / weightTotal
		iv.Score = &aggregate
	}
	iv.Status = IntStatusCompleted
	now := time.Now().UnixNano()
	iv.CompletedAt = &now
	if err := s.repo.UpdateInterview(ctx, iv); err != nil {
		return nil, err
	}
	return interviewToResponse(iv), nil
}

func interviewerToResponse(i *Interviewer) *InterviewerResponse {
	return &InterviewerResponse{
		ID:          i.ID.String(),
		InterviewID: i.InterviewID.String(),
		EmployeeID:  i.EmployeeID.String(),
		Role:        i.Role,
		CreatedAt:   i.CreatedAt,
	}
}

func scorecardItemToResponse(s *InterviewScorecardItem) *ScorecardItemResponse {
	resp := &ScorecardItemResponse{
		ID:          s.ID.String(),
		InterviewID: s.InterviewID.String(),
		Criterion:   s.Criterion,
		Weight:      s.Weight,
		Notes:       s.Notes,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
	if s.Score != nil {
		resp.Score = *s.Score
	}
	return resp
}

func assessmentToResponse(a *RecruitmentAssessment) *AssessmentResponse {
	resp := &AssessmentResponse{
		ID:          a.ID.String(),
		Name:        a.Name,
		Type:        a.Type,
		ScheduledAt: a.ScheduledAt,
		Location:    a.Location,
		Notes:       a.Notes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	if a.RequisitionID != nil {
		resp.RequisitionID = a.RequisitionID.String()
	}
	if a.MeetingLink != nil {
		resp.MeetingLink = *a.MeetingLink
	}
	return resp
}

func assessmentParticipantToResponse(p *AssessmentParticipant) *AssessmentParticipantResponse {
	resp := &AssessmentParticipantResponse{
		ID:             p.ID.String(),
		AssessmentID:   p.AssessmentID.String(),
		ApplicationID:  p.ApplicationID.String(),
		Status:         p.Status,
		Recommendation: p.Recommendation,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	if p.Score != nil {
		resp.Score = *p.Score
	}
	if p.Result != nil {
		resp.Result = *p.Result
	}
	return resp
}

func applicationScreeningToResponse(sc *ApplicationScreening) *ApplicationScreeningResponse {
	resp := &ApplicationScreeningResponse{
		ID:            sc.ID.String(),
		ApplicationID: sc.ApplicationID.String(),
		ScreenedAt:    sc.ScreenedAt,
		Result:        sc.Result,
		Notes:         sc.Notes,
		CreatedAt:     sc.CreatedAt,
		UpdatedAt:     sc.UpdatedAt,
	}
	if sc.ScreenedBy != nil {
		resp.ScreenedBy = sc.ScreenedBy.String()
	}
	if sc.Score != nil {
		resp.Score = *sc.Score
	}
	return resp
}
