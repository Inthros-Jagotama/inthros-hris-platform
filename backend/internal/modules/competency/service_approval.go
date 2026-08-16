package competency

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
		return nil, fmt.Errorf("assessment already finalized")
	}
	if target.Status == "submitted" {
		return nil, fmt.Errorf("assessment already submitted for approval")
	}
	if target.ApprovalInstanceID != nil {
		return nil, fmt.Errorf("assessment already has an active approval instance")
	}

	raters, err := s.repo.FindRatersByTarget(ctx, targetUID)
	if err != nil {
		return nil, err
	}
	if len(raters) == 0 {
		return nil, fmt.Errorf("no raters assigned to this assessment")
	}
	for _, r := range raters {
		if r.Status != string(RaterStatusSubmitted) {
			return nil, fmt.Errorf("not all raters have submitted yet (%s pending)", r.RaterType)
		}
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
	return nil
}
