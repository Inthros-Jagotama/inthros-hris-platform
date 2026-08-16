package competency

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Sentinel error validasi submit-approval — handler memetakannya ke status
// 4xx (bukan 500).
var (
	ErrAssessmentAlreadyFinalized  = errors.New("assessment already finalized")
	ErrAssessmentAlreadySubmitted  = errors.New("assessment already submitted for approval")
	ErrAssessmentHasApproval       = errors.New("assessment already has an active approval instance")
	ErrNoRatersAssigned            = errors.New("no raters assigned to this assessment")
	ErrNotAllRatersSubmitted       = errors.New("not all raters have submitted yet")
	ErrInvalidIndicator            = errors.New("indicator does not belong to this assessment template")
)

// PendingRaterInfo — rater yang belum submit saat submit-approval ditolak.
// Disertakan di response error agar pengguna tahu rater mana yang belum mengisi.
type PendingRaterInfo struct {
	RaterID    string `json:"rater_id"`
	RaterType  string `json:"rater_type"`
	RaterName  string `json:"rater_name,omitempty"`
	EmployeeID string `json:"employee_id,omitempty"`
	Status     string `json:"status"`
}

// ratersPendingError membungkus ErrNotAllRatersSubmitted dengan daftar rater
// yang belum submit — errors.Is tetap mengenali sentinel-nya (Unwrap).
type ratersPendingError struct {
	Pending []PendingRaterInfo
}

func (e *ratersPendingError) Error() string {
	parts := make([]string, 0, len(e.Pending))
	for _, p := range e.Pending {
		name := strings.TrimSpace(p.RaterName)
		if name == "" {
			name = p.RaterType
		}
		parts = append(parts, fmt.Sprintf("%s %s", p.RaterType, name))
	}
	return fmt.Sprintf("%s (%s)", ErrNotAllRatersSubmitted.Error(), strings.Join(parts, ", "))
}

func (e *ratersPendingError) Unwrap() error {
	return ErrNotAllRatersSubmitted
}

const (
	// ApprovalModuleCompetency360Assessment adalah slug approval flow untuk
	// finalisasi hasil assessment 360 (plan generik §13, §34.2). Cukup satu
	// checkpoint — rater submission bukan approval, itu business-state
	// internal modul. Flow di-resolve otomatis; alias ke subscription
	// "competency" di approval/subscriptionModuleAliases (bukan slug
	// subscription asli), jadi tanpa entry itu CreateApprovalInstance selalu
	// menolak dengan "module not subscribed".
	ApprovalModuleCompetency360Assessment = "competency_360_assessment"
)

// ApprovalEngine abstracts the central approval module so assessment
// finalization can be routed through it (same narrow-interface-plus-adapter
// pattern performance/leave/payroll use). Implemented via an adapter wrapping
// approval.Service in main.go. GetActiveFlowIDForModule lets the flow be
// resolved automatically instead of the caller picking a flow_id manually.
type ApprovalEngine interface {
	CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error)
	GetActiveFlowIDForModule(ctx context.Context, module string) (string, error)
}

// SetApprovalEngine wires the central approval module into this service.
func (s *Service) SetApprovalEngine(ae ApprovalEngine) {
	s.approvalEngine = ae
}

// SubmitAssessmentForApproval mengajukan finalisasi sebuah assessment target
// (competency_event_target) ke Central Approval. Prasyarat: seluruh rater
// yang ditugaskan sudah submit dan target belum pernah disubmit. Hasil
// dihitung HANYA setelah approval (lihat HandleAssessmentApprovalStatusChange).
func (s *Service) SubmitAssessmentForApproval(ctx context.Context, targetID string) (*CompetencyEventTargetResponse, error) {
	targetUID, err := uuid.Parse(targetID)
	if err != nil {
		return nil, fmt.Errorf("invalid target id: %w", err)
	}
	target, err := s.repo.FindCompetencyEventTargetByID(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	if target.Status == "finalized" {
		return nil, ErrAssessmentAlreadyFinalized
	}
	if target.Status == "submitted" {
		return nil, ErrAssessmentAlreadySubmitted
	}
	if target.ApprovalInstanceID != nil {
		return nil, ErrAssessmentHasApproval
	}

	raters, err := s.repo.FindRatersByTarget(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	if len(raters) == 0 {
		return nil, ErrNoRatersAssigned
	}

	// Kumpulkan SEMUA rater yang belum submit (bukan berhenti di yang pertama)
	// + resolve nama karyawannya agar response error menyebutkan siapa saja.
	pending := make([]PendingRaterInfo, 0)
	var pendingEmpIDs []uuid.UUID
	for _, r := range raters {
		if r.Status != string(RaterStatusSubmitted) {
			pending = append(pending, PendingRaterInfo{
				RaterID:    r.ID.String(),
				RaterType:  r.RaterType,
				EmployeeID: r.RaterEmployeeID.String(),
				Status:     r.Status,
			})
			pendingEmpIDs = append(pendingEmpIDs, r.RaterEmployeeID)
		}
	}
	if len(pending) > 0 {
		if names, err := s.repo.GetEmployeeNamesByIDs(ctx, pendingEmpIDs); err == nil {
			for i := range pending {
				pending[i].RaterName = names[pending[i].EmployeeID]
			}
		}
		return nil, &ratersPendingError{Pending: pending}
	}

	// Route through Central Approval when a flow is configured. Kalau flow
	// tidak terkonfigurasi sama sekali, submission tetap berjalan dengan
	// status business "submitted" (pola backward-compatible performance) —
	// finalisasi bisa dilakukan jalur manual/status handler bila flow
	// ditambahkan kemudian.
	if s.approvalEngine != nil {
		if flowID, err := s.approvalEngine.GetActiveFlowIDForModule(ctx, ApprovalModuleCompetency360Assessment); err == nil {
			instanceID, err := s.approvalEngine.CreateApprovalInstance(ctx, ApprovalModuleCompetency360Assessment, target.ID.String(), flowID)
			if err != nil {
				return nil, fmt.Errorf("failed to route assessment for approval: %w", err)
			}
			parsed, err := uuid.Parse(instanceID)
			if err != nil {
				return nil, fmt.Errorf("invalid approval instance id returned: %w", err)
			}
			target.ApprovalInstanceID = &parsed
		}
	}
	target.Status = "submitted"
	if err := s.repo.UpdateCompetencyEventTarget(ctx, target); err != nil {
		return nil, err
	}
	response := target.ToResponse()
	return &response, nil
}

// HandleAssessmentApprovalStatusChange di-invoke oleh callback push-based
// Approval Engine saat instance milik module competency_360_assessment
// mencapai status final (APPROVED/REJECTED/CANCELLED) — lihat
// approvalSvc.RegisterStatusHandler di main.go.
func (s *Service) HandleAssessmentApprovalStatusChange(ctx context.Context, documentID uuid.UUID, status string, note string) error {
	target, err := s.repo.FindCompetencyEventTargetByID(ctx, documentID)
	if err != nil {
		return err
	}
	if target.Status != "submitted" {
		return nil
	}

	now := time.Now()
	switch status {
	case "APPROVED":
		target.Status = "finalized"
		target.FinalizedAt = &now
	case "REJECTED", "CANCELLED":
		target.Status = "in_progress"
		target.ApprovalInstanceID = nil
	default:
		return nil
	}

	s.logger.Info("Competency assessment approval status updated",
		zap.String("target_id", target.ID.String()),
		zap.String("approval_status", status),
	)
	if err := s.repo.UpdateCompetencyEventTarget(ctx, target); err != nil {
		return err
	}

	// Hasil HANYA dihitung setelah approval (Phase 5 — §14). Kegagalan
	// kalkulasi tidak boleh membuat callback approval gagal; status target
	// sudah tersimpan, hasil bisa di-recalculate manual via FinalizeTarget.
	if status == "APPROVED" {
		if _, err := s.FinalizeTarget(ctx, target.ID.String()); err != nil {
			s.logger.Warn("assessment approved but calculation failed — result not stored",
				zap.String("target_id", target.ID.String()),
				zap.Error(err),
			)
		}
	}
	return nil
}
