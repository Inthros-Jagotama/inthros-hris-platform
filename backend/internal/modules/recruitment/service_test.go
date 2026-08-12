package recruitment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func newTestService() (*Service, func()) {
	db, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger := zap.NewNop()
	svc := NewService(repo, logger)
	seedDefaultRecruitmentStages(db)

	return svc, func() { cleanup() }
}

// =========================================================================
// Requisition Service Tests
// =========================================================================

func TestService_CreateRequisition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
		Department:     "Engineering",
		EmploymentType: "FULL_TIME",
		SlotsAvailable: intPtr(3),
	}

	resp, err := svc.CreateRequisition(ctx, req)
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.Title != "Software Engineer" {
		t.Errorf("expected 'Software Engineer', got '%s'", resp.Title)
	}
	if resp.SlotsAvailable != 3 {
		t.Errorf("expected 3 slots, got %d", resp.SlotsAvailable)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected default status 'DRAFT', got '%s'", resp.Status)
	}
}

func TestService_GetRequisitionByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	_, err := svc.GetRequisitionByID(context.Background(), "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetRequisitionByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	_, err := svc.GetRequisitionByID(context.Background(), uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent requisition")
	}
}

func TestService_UpdateRequisition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Old Title",
	})

	newTitle := "New Title"
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.Title != "New Title" {
		t.Errorf("expected 'New Title', got '%s'", updated.Title)
	}
}

// =========================================================================
// Succession Gap → Fallback External Recruitment (S-5 strategic layer)
// =========================================================================

type fakeSuccessionProvider struct {
	isGap bool
	err   error
}

func (f fakeSuccessionProvider) SuccessionGapForPosition(ctx context.Context, positionID uuid.UUID) (bool, error) {
	return f.isGap, f.err
}

func TestService_CreateRequisition_SuccessionGap_Validated(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetSuccessionGapProvider(fakeSuccessionProvider{isGap: true})
	ctx := context.Background()

	reason := string(ReqReasonSuccessionGap)
	posID := createTestUUID()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:      createTestOrgID(),
		Title:               "CTO (fallback external)",
		ReasonType:          &reason,
		SuccessionPositionID: &posID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.ReasonType != reason {
		t.Errorf("expected reason_type %q, got %q", reason, resp.ReasonType)
	}
	if resp.SuccessionPositionID != posID {
		t.Errorf("expected succession_position_id %s, got %s", posID, resp.SuccessionPositionID)
	}

	// Round-trip: re-read dari repo memastikan kolom benar-benar tersimpan di
	// DB (bukan hanya muncul di response dari objek in-memory).
	persisted, err := svc.GetRequisitionByID(ctx, resp.ID)
	if err != nil {
		t.Fatalf("GetRequisitionByID failed: %v", err)
	}
	if persisted.SuccessionPositionID != posID {
		t.Errorf("expected succession_position_id %s persisted in DB, got %s", posID, persisted.SuccessionPositionID)
	}
	if persisted.ReasonType != reason {
		t.Errorf("expected reason_type %q persisted in DB, got %q", reason, persisted.ReasonType)
	}
}

func TestService_CreateRequisition_SuccessionGap_NoProviderFallsBack(t *testing.T) {
	// Provider tidak di-wire (nil) — requisition SUCCESSION_GAP tetap dibuat
	// tanpa error, referensi succession_position_id tetap tersimpan (fail-safe).
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	reason := string(ReqReasonSuccessionGap)
	posID := createTestUUID()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:      createTestOrgID(),
		Title:               "CTO (fallback external)",
		ReasonType:          &reason,
		SuccessionPositionID: &posID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.SuccessionPositionID != posID {
		t.Errorf("expected succession_position_id preserved without provider, got %s", resp.SuccessionPositionID)
	}
}

func TestService_CreateRequisition_SuccessionGap_ProviderErrorKeepsGoing(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetSuccessionGapProvider(fakeSuccessionProvider{isGap: false, err: fmt.Errorf("ci unavailable")})
	ctx := context.Background()

	reason := string(ReqReasonSuccessionGap)
	posID := createTestUUID()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:      createTestOrgID(),
		Title:               "CTO (fallback external)",
		ReasonType:          &reason,
		SuccessionPositionID: &posID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed on provider error: %v", err)
	}
	if resp.SuccessionPositionID != posID {
		t.Errorf("expected succession_position_id preserved on provider error, got %s", resp.SuccessionPositionID)
	}
}

func TestService_UpdateRequisition_ChangeReasonToSuccessionGap(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetSuccessionGapProvider(fakeSuccessionProvider{isGap: true})
	ctx := context.Background()

	created, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "CTO",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}

	reason := string(ReqReasonSuccessionGap)
	posID := createTestUUID()
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		ReasonType:          &reason,
		SuccessionPositionID: &posID,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.ReasonType != reason {
		t.Errorf("expected reason_type %q, got %q", reason, updated.ReasonType)
	}
	if updated.SuccessionPositionID != posID {
		t.Errorf("expected succession_position_id %s, got %s", posID, updated.SuccessionPositionID)
	}
}

// =========================================================================
// Workforce Gap → Requisition (S-1 strategic layer)
// =========================================================================

type fakeGapProvider struct {
	need int
	err  error
}

func (f fakeGapProvider) HiringGapForOrganization(ctx context.Context, orgID uuid.UUID) (int, error) {
	return f.need, f.err
}

func TestService_CreateRequisition_WorkforceGap_AutoResolveSlots(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetWorkforceGapProvider(fakeGapProvider{need: 3})
	ctx := context.Background()

	reason := string(ReqReasonWorkforceGap)
	gapID := createTestUUID()
	planID := createTestUUID()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:  createTestOrgID(),
		Title:           "Ops Staff",
		ReasonType:      &reason,
		WorkforceGapID:  &gapID,
		WorkforcePlanID: &planID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.ReasonType != reason {
		t.Errorf("expected reason_type %q, got %q", reason, resp.ReasonType)
	}
	if resp.SlotsAvailable != 3 {
		t.Errorf("expected slots auto-resolved to 3, got %d", resp.SlotsAvailable)
	}
	if resp.WorkforceGapID != gapID {
		t.Errorf("expected workforce_gap_id %s, got %s", gapID, resp.WorkforceGapID)
	}
	if resp.WorkforcePlanID != planID {
		t.Errorf("expected workforce_plan_id %s, got %s", planID, resp.WorkforcePlanID)
	}
}

func TestService_CreateRequisition_WorkforceGap_NoProviderFallsBack(t *testing.T) {
	// Provider tidak di-wire (nil) — requisition WORKFORCE_GAP tetap dibuat
	// dengan slots default, tidak error (fail-safe plan S-1).
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	reason := string(ReqReasonWorkforceGap)
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Ops Staff",
		ReasonType:     &reason,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.SlotsAvailable != 1 {
		t.Errorf("expected default slots 1 without provider, got %d", resp.SlotsAvailable)
	}
}

func TestService_CreateRequisition_WorkforceGap_ProviderErrorKeepsDefault(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetWorkforceGapProvider(fakeGapProvider{need: 0, err: fmt.Errorf("wi unavailable")})
	ctx := context.Background()

	reason := string(ReqReasonWorkforceGap)
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Ops Staff",
		ReasonType:     &reason,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.SlotsAvailable != 1 {
		t.Errorf("expected default slots 1 on provider error, got %d", resp.SlotsAvailable)
	}
}

func TestService_CreateRequisition_WorkforceGap_ExplicitSlotsWins(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetWorkforceGapProvider(fakeGapProvider{need: 5})
	ctx := context.Background()

	reason := string(ReqReasonWorkforceGap)
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:  createTestOrgID(),
		Title:           "Ops Staff",
		ReasonType:      &reason,
		SlotsAvailable:  intPtr(2),
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.SlotsAvailable != 2 {
		t.Errorf("expected explicit slots 2 to win, got %d", resp.SlotsAvailable)
	}
}

func TestService_UpdateRequisition_ChangeReasonToWorkforceGap(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Ops Staff",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if created.ReasonType != "" {
		t.Errorf("expected empty reason_type by default, got %q", created.ReasonType)
	}

	svc.SetWorkforceGapProvider(fakeGapProvider{need: 4})
	reason := string(ReqReasonWorkforceGap)
	gapID := createTestUUID()
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		ReasonType:     &reason,
		WorkforceGapID: &gapID,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.ReasonType != reason {
		t.Errorf("expected reason_type %q, got %q", reason, updated.ReasonType)
	}
	if updated.SlotsAvailable != 4 {
		t.Errorf("expected slots auto-resolved to 4 on reason change, got %d", updated.SlotsAvailable)
	}
	if updated.WorkforceGapID != gapID {
		t.Errorf("expected workforce_gap_id %s, got %s", gapID, updated.WorkforceGapID)
	}
}

func TestService_UpdateRequisition_UnrelatedUpdateKeepsGapSlots(t *testing.T) {
	// Regression: requisition yang SUDAH reason_type=WORKFORCE_GAP (slots
	// auto-resolved 3) lalu di-update field lain (title) tanpa slots —
	// slots_available harus TETAP 3, bukan di-resolve ulang dari provider.
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetWorkforceGapProvider(fakeGapProvider{need: 3})
	reason := string(ReqReasonWorkforceGap)
	created, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Ops Staff",
		ReasonType:     &reason,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if created.SlotsAvailable != 3 {
		t.Fatalf("expected auto-resolved slots 3, got %d", created.SlotsAvailable)
	}

	// Provider berubah jadi hiring need 9 — update title TIDAK boleh
	// menimpa slots 3.
	svc.SetWorkforceGapProvider(fakeGapProvider{need: 9})
	newTitle := "Ops Staff (Updated)"
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("expected title updated, got %q", updated.Title)
	}
	if updated.SlotsAvailable != 3 {
		t.Errorf("expected slots to stay 3 on unrelated update, got %d", updated.SlotsAvailable)
	}
}

// =========================================================================
// Internal Candidate (S-4 — CI → Recruitment)
// =========================================================================

type fakeInternalCandidateProvider struct {
	candidates []InternalCandidate
	err        error
}

func (f fakeInternalCandidateProvider) EligibleEmployeesForPosition(ctx context.Context, targetPositionID uuid.UUID) ([]InternalCandidate, error) {
	return f.candidates, f.err
}

func TestService_GetEligibleInternalCandidates_NoProviderReturnsEmpty(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	result, err := svc.GetEligibleInternalCandidates(ctx, uuid.New().String())
	if err != nil {
		t.Fatalf("GetEligibleInternalCandidates failed: %v", err)
	}
	if result == nil || len(result) != 0 {
		t.Errorf("expected empty list without provider, got %v", result)
	}
}

func TestService_GetEligibleInternalCandidates_InvalidPosition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	_, err := svc.GetEligibleInternalCandidates(context.Background(), "bad-uuid")
	if err == nil {
		t.Fatal("expected error for invalid position_id")
	}
}

func TestService_GetEligibleInternalCandidates_WithProvider(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetInternalCandidateProvider(fakeInternalCandidateProvider{
		candidates: []InternalCandidate{
			{EmployeeID: uuid.New().String(), Name: "Andi", CurrentPositionName: "Staff IT"},
		},
	})
	ctx := context.Background()

	result, err := svc.GetEligibleInternalCandidates(ctx, uuid.New().String())
	if err != nil {
		t.Fatalf("GetEligibleInternalCandidates failed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 candidate from provider, got %d", len(result))
	}
	if result[0].Name != "Andi" {
		t.Errorf("expected candidate Andi, got %s", result[0].Name)
	}
}

func TestService_UpdateRequisition_ClearWorkforceGapID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	gapID := createTestUUID()
	created, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Ops Staff",
		WorkforceGapID: &gapID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}

	empty := ""
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		WorkforceGapID: &empty,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.WorkforceGapID != "" {
		t.Errorf("expected cleared workforce_gap_id, got %q", updated.WorkforceGapID)
	}
}

func TestService_DeleteRequisition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "To Delete",
	})
	if err := svc.DeleteRequisition(ctx, created.ID); err != nil {
		t.Fatalf("DeleteRequisition failed: %v", err)
	}
}

// =========================================================================
// Candidate Service Tests
// =========================================================================

func TestService_CreateCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateCandidateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "08123456789",
	}

	resp, err := svc.CreateCandidate(ctx, req)
	if err != nil {
		t.Fatalf("CreateCandidate failed: %v", err)
	}
	if resp.FirstName != "John" {
		t.Errorf("expected 'John', got '%s'", resp.FirstName)
	}
	if resp.Email != "john@example.com" {
		t.Errorf("expected 'john@example.com', got '%s'", resp.Email)
	}
}

func TestService_CreateCandidate_DuplicateEmail(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "dup@example.com",
	})
	_, err := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Jane", LastName: "Doe", Email: "dup@example.com",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestService_UpdateCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@example.com",
	})

	newPhone := "08999999999"
	updated, err := svc.UpdateCandidate(ctx, created.ID, UpdateCandidateRequest{
		Phone: &newPhone,
	})
	if err != nil {
		t.Fatalf("UpdateCandidate failed: %v", err)
	}
	if updated.Phone != "08999999999" {
		t.Errorf("expected '08999999999', got '%s'", updated.Phone)
	}
}

func TestService_DeleteCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "del@example.com",
	})
	if err := svc.DeleteCandidate(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCandidate failed: %v", err)
	}
}

// =========================================================================
// Application Service Tests
// =========================================================================

func TestService_CreateApplication(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@test.com",
	})

	resp, err := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID,
		CandidateID:   cand.ID,
	})
	if err != nil {
		t.Fatalf("CreateApplication failed: %v", err)
	}
	if resp.Status != "NEW" {
		t.Errorf("expected status 'NEW', got '%s'", resp.Status)
	}
}

func TestService_CreateApplication_InvalidRequisition(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	_, err := svc.CreateApplication(context.Background(), CreateApplicationRequest{
		RequisitionID: uuid.New().String(),
		CandidateID:   uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected error for non-existent requisition")
	}
}

func TestService_UpdateApplicationStatus(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@test.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})

	updated, err := svc.UpdateApplicationStatus(ctx, app.ID, "SHORTLISTED", "", "")
	if err != nil {
		t.Fatalf("UpdateApplicationStatus failed: %v", err)
	}
	if updated.Status != "SHORTLISTED" {
		t.Errorf("expected 'SHORTLISTED', got '%s'", updated.Status)
	}
}

func TestService_UpdateApplicationStatus_ForwardJumpAllowed(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "fwd@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	// NEW -> OFFERED direct jump must remain allowed (state machine allows
	// forward jumps between non-terminal stages).
	updated, err := svc.UpdateApplicationStatus(ctx, app.ID, "OFFERED", "", "")
	if err != nil {
		t.Fatalf("expected forward jump to succeed, got error: %v", err)
	}
	if updated.Status != "OFFERED" {
		t.Errorf("expected OFFERED, got %s", updated.Status)
	}
}

func TestService_UpdateApplicationStatus_BackwardRejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "back@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "SHORTLISTED", "", ""); err != nil {
		t.Fatalf("setup transition failed: %v", err)
	}

	_, err := svc.UpdateApplicationStatus(ctx, app.ID, "NEW", "", "")
	if err == nil {
		t.Fatal("expected error for backward transition SHORTLISTED -> NEW, got nil")
	}
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition, got: %v", err)
	}
}

func TestService_UpdateApplicationStatus_FromTerminalRejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "term@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "REJECTED", "", ""); err != nil {
		t.Fatalf("setup transition failed: %v", err)
	}

	_, err := svc.UpdateApplicationStatus(ctx, app.ID, "SCREENED", "", "")
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition from terminal status, got: %v", err)
	}
}

func TestService_UpdateApplicationStatus_SameStatusNoop(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "noop@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "NEW", "", ""); err != nil {
		t.Fatalf("same-status transition should be a no-op, got error: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	// Only the initial NEW history row from CreateApplication — the no-op
	// NEW->NEW call must not add a second row.
	if len(hist) != 1 {
		t.Errorf("expected 1 history row (initial only), got %d", len(hist))
	}
}

func TestService_CreateApplication_WritesInitialHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "init@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	if len(hist) != 1 {
		t.Fatalf("expected 1 initial history row, got %d", len(hist))
	}
	if hist[0].FromStage != nil {
		t.Errorf("expected initial history from_stage nil, got %v", hist[0].FromStage)
	}
	if hist[0].ToStage.Code != "NEW" {
		t.Errorf("expected initial history to_stage NEW, got %s", hist[0].ToStage.Code)
	}
}

func TestService_UpdateApplicationStatus_WritesHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "wh@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "SCREENED", "", "moved to screening"); err != nil {
		t.Fatalf("UpdateApplicationStatus failed: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	if len(hist) != 2 {
		t.Fatalf("expected 2 history rows (initial + transition), got %d", len(hist))
	}
	last := hist[len(hist)-1]
	if last.FromStage == nil || last.FromStage.Code != "NEW" || last.ToStage.Code != "SCREENED" {
		t.Errorf("expected NEW->SCREENED, got from=%v to=%s", last.FromStage, last.ToStage.Code)
	}
	if last.Notes != "moved to screening" {
		t.Errorf("expected notes preserved, got %q", last.Notes)
	}
}

func TestService_DeleteApplication(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@test.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	if err := svc.DeleteApplication(ctx, app.ID); err != nil {
		t.Fatalf("DeleteApplication failed: %v", err)
	}
}

// =========================================================================
// Interview Service Tests
// =========================================================================

func TestService_CreateInterview(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@test.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})

	resp, err := svc.CreateInterview(ctx, CreateInterviewRequest{
		ApplicationID: app.ID,
		InterviewerID: createTestUUID(),
		Stage:         "TECHNICAL",
		ScheduledAt:   1760000000,
	})
	if err != nil {
		t.Fatalf("CreateInterview failed: %v", err)
	}
	if resp.Stage != "TECHNICAL" {
		t.Errorf("expected stage 'TECHNICAL', got '%s'", resp.Stage)
	}
	if resp.Status != "SCHEDULED" {
		t.Errorf("expected status 'SCHEDULED', got '%s'", resp.Status)
	}
}

func TestService_UpdateInterview_Complete(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@test.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	iv, _ := svc.CreateInterview(ctx, CreateInterviewRequest{
		ApplicationID: app.ID, InterviewerID: createTestUUID(),
		Stage: "HR", ScheduledAt: 1760000000,
	})

	completed := "COMPLETED"
	score := 85.5
	feedback := "Great candidate"
	updated, err := svc.UpdateInterview(ctx, iv.ID, UpdateInterviewRequest{
		Status: &completed, Score: &score, Feedback: &feedback,
	})
	if err != nil {
		t.Fatalf("UpdateInterview failed: %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Errorf("expected 'COMPLETED', got '%s'", updated.Status)
	}
	if updated.Score != 85.5 {
		t.Errorf("expected score 85.5, got %f", updated.Score)
	}
}

// =========================================================================
// Onboarding Service Tests
// =========================================================================

func TestService_CreateOnboardingTaskTemplate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateOnboardingTaskTemplate(ctx, CreateOnboardingTaskTemplateRequest{
		Name:        "IT Account Setup",
		Description: "Create email and system accounts",
		Category:    "IT",
		DayOffset:   intPtr(-7),
	})
	if err != nil {
		t.Fatalf("CreateOnboardingTaskTemplate failed: %v", err)
	}
	if resp.Name != "IT Account Setup" {
		t.Errorf("expected 'IT Account Setup', got '%s'", resp.Name)
	}
	if resp.IsMandatory != true {
		t.Errorf("expected is_mandatory true, got %v", resp.IsMandatory)
	}
}

func TestService_CreateEmployeeOnboarding(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	// Seed some task templates first
	svc.CreateOnboardingTaskTemplate(ctx, CreateOnboardingTaskTemplateRequest{
		Name: "Contract Signing", Category: "LEGAL", DayOffset: intPtr(0),
	})
	svc.CreateOnboardingTaskTemplate(ctx, CreateOnboardingTaskTemplateRequest{
		Name: "IT Setup", Category: "IT", DayOffset: intPtr(-7),
	})

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@onboard.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	// Accept candidate
	svc.UpdateApplicationStatus(ctx, app.ID, "ACCEPTED", "", "")

	resp, err := svc.CreateEmployeeOnboarding(ctx, CreateEmployeeOnboardingRequest{
		EmployeeID:    createTestUUID(),
		ApplicationID: app.ID,
		StartDate:     "2026-08-01",
	})
	if err != nil {
		t.Fatalf("CreateEmployeeOnboarding failed: %v", err)
	}
	if resp.Status != "PENDING" {
		t.Errorf("expected status 'PENDING', got '%s'", resp.Status)
	}
	if resp.StartDate != "2026-08-01" {
		t.Errorf("expected start_date '2026-08-01', got '%s'", resp.StartDate)
	}
}

func TestService_CreateOnboardingTaskItem(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@item.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	onb, _ := svc.CreateEmployeeOnboarding(ctx, CreateEmployeeOnboardingRequest{
		EmployeeID:    createTestUUID(),
		ApplicationID: app.ID,
		StartDate:     "2026-08-01",
	})

	resp, err := svc.CreateOnboardingTaskItem(ctx, CreateOnboardingTaskItemRequest{
		EmployeeOnboardingID: onb.ID,
		Name:                "Custom Task",
		Description:         "A custom onboarding task",
	})
	if err != nil {
		t.Fatalf("CreateOnboardingTaskItem failed: %v", err)
	}
	if resp.Name != "Custom Task" {
		t.Errorf("expected 'Custom Task', got '%s'", resp.Name)
	}
}

func TestService_UpdateOnboardingTaskItem_Complete(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@complete.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	onb, _ := svc.CreateEmployeeOnboarding(ctx, CreateEmployeeOnboardingRequest{
		EmployeeID:    createTestUUID(),
		ApplicationID: app.ID,
		StartDate:     "2026-08-01",
	})
	item, _ := svc.CreateOnboardingTaskItem(ctx, CreateOnboardingTaskItemRequest{
		EmployeeOnboardingID: onb.ID,
		Name: "Task to Complete",
	})

	completed := true
	updated, err := svc.UpdateOnboardingTaskItem(ctx, item.ID, UpdateOnboardingTaskItemRequest{
		IsCompleted: &completed,
	})
	if err != nil {
		t.Fatalf("UpdateOnboardingTaskItem failed: %v", err)
	}
	if !updated.IsCompleted {
		t.Error("expected is_completed true")
	}
}

func TestService_ListCandidates_Search(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Alice", LastName: "Smith", Email: "alice@test.com",
	})
	svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Bob", LastName: "Jones", Email: "bob@test.com",
	})
	svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Charlie", LastName: "Brown", Email: "charlie@test.com",
	})

	search := "bob"
	resp, err := svc.ListCandidates(ctx, &search, 1, 10)
	if err != nil {
		t.Fatalf("ListCandidates search failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 result for 'bob', got %d", resp.Total)
	}
}

// =========================================================================
// Onboarding → Training Handoff (S-7 strategic layer)
// =========================================================================

type fakeTrainingHandoff struct {
	calls    int
	empID    uuid.UUID
	onbID    uuid.UUID
	err      error
}

func (f *fakeTrainingHandoff) CreateOnboardingNeed(ctx context.Context, employeeID, onboardingID uuid.UUID, reason string) error {
	f.calls++
	f.empID = employeeID
	f.onbID = onboardingID
	return f.err
}

// seedOnboarding mengembalikan EmployeeOnboarding PENDING untuk test handoff.
func seedOnboarding(t *testing.T, svc *Service, ctx context.Context) *EmployeeOnboardingResponse {
	t.Helper()
	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Engineer",
	})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: "john@handoff.com",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	onb, err := svc.CreateEmployeeOnboarding(ctx, CreateEmployeeOnboardingRequest{
		EmployeeID:    createTestUUID(),
		ApplicationID: app.ID,
		StartDate:     "2026-08-01",
	})
	if err != nil {
		t.Fatalf("CreateEmployeeOnboarding failed: %v", err)
	}
	return onb
}

func TestService_OnboardingComplete_TriggersTrainingHandoff(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeTrainingHandoff{}
	svc.SetTrainingHandoffProvider(provider)
	ctx := context.Background()

	onb := seedOnboarding(t, svc, ctx)
	status := "COMPLETED"
	updated, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status})
	if err != nil {
		t.Fatalf("UpdateEmployeeOnboarding failed: %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED, got %s", updated.Status)
	}
	if provider.calls != 1 {
		t.Fatalf("expected 1 handoff call, got %d", provider.calls)
	}
	if provider.onbID.String() != onb.ID || provider.empID.String() != onb.EmployeeID {
		t.Errorf("expected handoff with onboarding %s + employee %s, got %s/%s",
			onb.ID, onb.EmployeeID, provider.onbID, provider.empID)
	}
}

func TestService_OnboardingRepeatedCompleted_NoDuplicateHandoff(t *testing.T) {
	// Handoff hanya saat BERTRANSISI ke COMPLETED — update kedua dengan status
	// yang sama (mis. update notes) tidak membuat TrainingNeed duplikat.
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeTrainingHandoff{}
	svc.SetTrainingHandoffProvider(provider)
	ctx := context.Background()

	onb := seedOnboarding(t, svc, ctx)
	status := "COMPLETED"
	notes := "updated note"
	if _, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status}); err != nil {
		t.Fatalf("first complete failed: %v", err)
	}
	if provider.calls != 1 {
		t.Fatalf("expected 1 handoff after first completion, got %d", provider.calls)
	}
	// Update ulang dengan status COMPLETED (notes berubah) — tidak boleh handoff lagi.
	if _, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status, Notes: &notes}); err != nil {
		t.Fatalf("second complete update failed: %v", err)
	}
	if provider.calls != 1 {
		t.Errorf("expected still 1 handoff on repeated COMPLETED update, got %d", provider.calls)
	}
}

func TestService_OnboardingNotCompleted_NoHandoff(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	provider := &fakeTrainingHandoff{}
	svc.SetTrainingHandoffProvider(provider)
	ctx := context.Background()

	onb := seedOnboarding(t, svc, ctx)
	status := "IN_PROGRESS"
	if _, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status}); err != nil {
		t.Fatalf("UpdateEmployeeOnboarding failed: %v", err)
	}
	if provider.calls != 0 {
		t.Errorf("expected no handoff for non-completed status, got %d calls", provider.calls)
	}
}

func TestService_OnboardingComplete_NoProviderFallsBack(t *testing.T) {
	// Provider tidak di-wire (nil) — onboarding tetap selesai tanpa error;
	// handoff hanya di-log (fail-safe plan S-7).
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	onb := seedOnboarding(t, svc, ctx)
	status := "COMPLETED"
	updated, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status})
	if err != nil {
		t.Fatalf("UpdateEmployeeOnboarding failed without provider: %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Errorf("expected status COMPLETED without provider, got %s", updated.Status)
	}
}

func TestService_OnboardingComplete_ProviderErrorKeepsCompleted(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetTrainingHandoffProvider(&fakeTrainingHandoff{err: fmt.Errorf("training unavailable")})
	ctx := context.Background()

	onb := seedOnboarding(t, svc, ctx)
	status := "COMPLETED"
	updated, err := svc.UpdateEmployeeOnboarding(ctx, onb.ID, UpdateEmployeeOnboardingRequest{Status: &status})
	if err != nil {
		t.Fatalf("UpdateEmployeeOnboarding failed on handoff error: %v", err)
	}
	if updated.Status != "COMPLETED" {
		t.Errorf("expected onboarding stays COMPLETED on handoff error, got %s", updated.Status)
	}
}

// =========================================================================
// Module Approval Integration (G-1 — requisition → Central Approval)
// =========================================================================

type fakeApprovalEngine struct {
	instanceID    string
	flowID        string
	createErr     error
	createdModule string
	createdDocID  string
	createdFlowID string
}

func (f *fakeApprovalEngine) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	f.createdModule = module
	f.createdDocID = documentID
	f.createdFlowID = flowID
	if f.createErr != nil {
		return "", f.createErr
	}
	return f.instanceID, nil
}

func (f *fakeApprovalEngine) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	return "PENDING", nil
}

func (f *fakeApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	if f.flowID == "" {
		return "", fmt.Errorf("no active flow for module %s", module)
	}
	return f.flowID, nil
}

// seedDraftRequisition membuat requisition DRAFT untuk test approval (G-1).
func seedDraftRequisition(t *testing.T, svc *Service, ctx context.Context) *RequisitionResponse {
	t.Helper()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	return resp
}

func TestService_SubmitRequisition_CreatesInstance(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	instanceID := uuid.New().String()
	engine := &fakeApprovalEngine{instanceID: instanceID, flowID: uuid.New().String()}
	svc.SetApprovalEngine(engine)
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	resp, err := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitRequisition failed: %v", err)
	}
	if resp.Status != string(ReqStatusSubmitted) {
		t.Errorf("expected status SUBMITTED, got %s", resp.Status)
	}
	if resp.ApprovalInstanceID != instanceID {
		t.Errorf("expected approval_instance_id %s, got %s", instanceID, resp.ApprovalInstanceID)
	}
	if engine.createdModule != "recruitment" {
		t.Errorf("expected module 'recruitment', got %q", engine.createdModule)
	}
	if engine.createdDocID != req.ID {
		t.Errorf("expected document_id %s, got %s", req.ID, engine.createdDocID)
	}
}

func TestService_SubmitRequisition_AutoResolveFlow(t *testing.T) {
	// flowID tidak dikirim — flow aktif modul recruitment di-auto-resolve
	// (pola employeemovement G-3).
	svc, cleanup := newTestService()
	defer cleanup()
	flowID := uuid.New().String()
	engine := &fakeApprovalEngine{instanceID: uuid.New().String(), flowID: flowID}
	svc.SetApprovalEngine(engine)
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	resp, err := svc.SubmitRequisition(ctx, req.ID, "")
	if err != nil {
		t.Fatalf("SubmitRequisition failed: %v", err)
	}
	if resp.Status != string(ReqStatusSubmitted) {
		t.Errorf("expected status SUBMITTED, got %s", resp.Status)
	}
	if engine.createdFlowID != flowID {
		t.Errorf("expected auto-resolved flow_id %s, got %s", flowID, engine.createdFlowID)
	}
}

func TestService_SubmitRequisition_NotDraft(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	// Pindahkan ke OPEN dulu (simulasi non-draft)
	open := string(ReqStatusOpen)
	if _, err := svc.UpdateRequisition(ctx, req.ID, UpdateRequisitionRequest{Status: &open}); err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	_, err := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())
	if err == nil {
		t.Fatal("expected error submitting non-draft requisition")
	}
}

func TestService_SubmitRequisition_NoEngine(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	_, err := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())
	if err == nil {
		t.Fatal("expected error when approval engine not wired")
	}
}

func TestService_SubmitRequisition_NoFlow(t *testing.T) {
	// Tidak ada flow aktif (GetActiveFlowIDForModule error) dan client tidak
	// mengirim flow_id — submit ditolak dengan pesan flow not configured.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	_, err := svc.SubmitRequisition(ctx, req.ID, "")
	if err == nil {
		t.Fatal("expected error when no approval flow configured")
	}
}

func TestService_HandleApprovalStatusChange_ApprovedOpens(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	submitted, err := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitRequisition failed: %v", err)
	}

	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, err := svc.GetRequisitionByID(ctx, submitted.ID)
	if err != nil {
		t.Fatalf("GetRequisitionByID failed: %v", err)
	}
	if persisted.Status != string(ReqStatusOpen) {
		t.Errorf("expected status OPEN after approval, got %s", persisted.Status)
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	submitted, _ := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())

	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "REJECTED", "budget not approved"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetRequisitionByID(ctx, submitted.ID)
	if persisted.Status != string(ReqStatusRejected) {
		t.Errorf("expected status REJECTED, got %s", persisted.Status)
	}
}

func TestService_HandleApprovalStatusChange_Cancelled(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	submitted, _ := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())

	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "CANCELLED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetRequisitionByID(ctx, submitted.ID)
	if persisted.Status != string(ReqStatusCancelled) {
		t.Errorf("expected status CANCELLED, got %s", persisted.Status)
	}
}

func TestService_HandleApprovalStatusChange_NotSubmittedIsNoop(t *testing.T) {
	// Callback untuk requisition yang belum SUBMITTED (mis. DRAFT/OPEN)
	// tidak mengubah status — idempotent, callback ganda aman.
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(req.ID), "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetRequisitionByID(ctx, req.ID)
	if persisted.Status != string(ReqStatusDraft) {
		t.Errorf("expected status unchanged DRAFT, got %s", persisted.Status)
	}
}

func TestService_HandleApprovalStatusChange_UnknownStatusIsNoop(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	submitted, _ := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())

	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "WEIRD_STATUS", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetRequisitionByID(ctx, submitted.ID)
	if persisted.Status != string(ReqStatusSubmitted) {
		t.Errorf("expected status unchanged SUBMITTED, got %s", persisted.Status)
	}
}

// =========================================================================
// Requisition Enhancement (G-2 — number, priority, position, opened_at)
// =========================================================================

func TestService_CreateRequisition_AutoGeneratesNumber(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.RequisitionNumber == "" {
		t.Fatal("expected auto-generated requisition_number")
	}
	// Format REQ-YYYYMM-XXXXXXXX (19 char)
	if len(resp.RequisitionNumber) != 19 || resp.RequisitionNumber[:4] != "REQ-" {
		t.Errorf("unexpected requisition_number format: %q", resp.RequisitionNumber)
	}
}

func TestService_CreateRequisition_RespectsExplicitNumberAndPriority(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	explicit := "REQ-202601-ABCD1234"
	priority := string(ReqPriorityUrgent)
	posID := createTestUUID()
	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID:    createTestOrgID(),
		Title:             "Software Engineer",
		RequisitionNumber: explicit,
		Priority:          &priority,
		PositionID:        &posID,
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.RequisitionNumber != explicit {
		t.Errorf("expected requisition_number %q, got %q", explicit, resp.RequisitionNumber)
	}
	if resp.Priority != priority {
		t.Errorf("expected priority %q, got %q", priority, resp.Priority)
	}
	if resp.PositionID != posID {
		t.Errorf("expected position_id %s, got %s", posID, resp.PositionID)
	}
}

func TestService_CreateRequisition_DefaultPriorityMedium(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if resp.Priority != string(ReqPriorityMedium) {
		t.Errorf("expected default priority MEDIUM, got %q", resp.Priority)
	}
}

func TestService_UpdateRequisition_SetsOpenedAtOnOpen(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
	})
	if created.OpenedAt != nil {
		t.Fatal("expected opened_at nil for draft")
	}

	open := string(ReqStatusOpen)
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{Status: &open})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.Status != string(ReqStatusOpen) {
		t.Errorf("expected status OPEN, got %s", updated.Status)
	}
	if updated.OpenedAt == nil || *updated.OpenedAt == 0 {
		t.Error("expected opened_at set when requisition becomes OPEN")
	}

	// Round-trip: persisted di DB
	persisted, _ := svc.GetRequisitionByID(ctx, created.ID)
	if persisted.OpenedAt == nil {
		t.Error("expected opened_at persisted in DB")
	}
}

func TestService_ApprovalApproved_SetsOpenedAt(t *testing.T) {
	// G-1 + G-2: approval APPROVED membuka requisition → opened_at diset otomatis.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	req := seedDraftRequisition(t, svc, ctx)
	submitted, err := svc.SubmitRequisition(ctx, req.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitRequisition failed: %v", err)
	}
	if submitted.OpenedAt != nil {
		t.Fatal("expected opened_at nil while submitted")
	}

	if err := svc.HandleApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetRequisitionByID(ctx, submitted.ID)
	if persisted.Status != string(ReqStatusOpen) {
		t.Errorf("expected status OPEN after approval, got %s", persisted.Status)
	}
	if persisted.OpenedAt == nil || *persisted.OpenedAt == 0 {
		t.Error("expected opened_at set by approval APPROVED")
	}
}

func TestService_UpdateRequisition_ClearPositionID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	posID := createTestUUID()
	created, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
		PositionID:     &posID,
	})
	if created.PositionID != posID {
		t.Fatalf("expected position_id %s, got %s", posID, created.PositionID)
	}

	empty := ""
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{PositionID: &empty})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.PositionID != "" {
		t.Errorf("expected cleared position_id, got %q", updated.PositionID)
	}
}

func TestService_UpdateRequisition_SetPriorityAndNumber(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(),
		Title:          "Software Engineer",
	})

	priority := string(ReqPriorityHigh)
	number := "REQ-202602-12345678"
	updated, err := svc.UpdateRequisition(ctx, created.ID, UpdateRequisitionRequest{
		Priority:          &priority,
		RequisitionNumber: &number,
	})
	if err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	if updated.Priority != priority {
		t.Errorf("expected priority %q, got %q", priority, updated.Priority)
	}
	if updated.RequisitionNumber != number {
		t.Errorf("expected requisition_number %q, got %q", number, updated.RequisitionNumber)
	}
}

// =========================================================================
// G-3 - Job Offer management
// =========================================================================

// seedOffer membuat requisition + candidate + application, lalu offer DRAFT
// dengan expiry date yang dikirim pemanggil (untuk test expired).
func seedOffer(t *testing.T, svc *Service, ctx context.Context, expiry string) *OfferResponse {
	return seedOfferWithEmail(t, svc, ctx, expiry, "john.offer@test.com")
}

func seedOfferWithEmail(t *testing.T, svc *Service, ctx context.Context, expiry, email string) *OfferResponse {
	t.Helper()
	req, err := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Software Engineer",
	})
	if err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	cand, err := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "John", LastName: "Doe", Email: email,
	})
	if err != nil {
		t.Fatalf("CreateCandidate failed: %v", err)
	}
	app, err := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: cand.ID,
	})
	if err != nil {
		t.Fatalf("CreateApplication failed: %v", err)
	}
	offer, err := svc.CreateOffer(ctx, CreateOfferRequest{
		ApplicationID: app.ID, EmploymentType: "FULL_TIME",
		Salary: 10000000, StartDate: "2026-09-01", ExpiryDate: expiry,
	})
	if err != nil {
		t.Fatalf("CreateOffer failed: %v", err)
	}
	return offer
}

func seedDraftOffer(t *testing.T, svc *Service, ctx context.Context) *OfferResponse {
	t.Helper()
	return seedOffer(t, svc, ctx, time.Now().AddDate(0, 0, 30).Format("2006-01-02"))
}

// approveAndSendOffer memajukan offer DRAFT -> APPROVED -> SENT (G-3).
func approveAndSendOffer(t *testing.T, svc *Service, ctx context.Context, offerID string) string {
	t.Helper()
	submitted, err := svc.SubmitOffer(ctx, offerID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitOffer failed: %v", err)
	}
	if err := svc.HandleOfferApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "APPROVED", ""); err != nil {
		t.Fatalf("HandleOfferApprovalStatusChange failed: %v", err)
	}
	if _, err := svc.SendOffer(ctx, submitted.ID); err != nil {
		t.Fatalf("SendOffer failed: %v", err)
	}
	return submitted.ID
}

func TestService_CreateOffer(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	if offer.Status != string(OfferStatusDraft) {
		t.Errorf("expected status DRAFT, got %s", offer.Status)
	}
	if offer.OfferNumber == "" {
		t.Error("expected auto-generated offer_number")
	}
	if !strings.HasPrefix(offer.OfferNumber, "OFF-") {
		t.Errorf("expected OFF- prefix, got %s", offer.OfferNumber)
	}
	if len(offer.OfferNumber) != 19 {
		t.Errorf("expected offer_number length 19, got %d (%s)", len(offer.OfferNumber), offer.OfferNumber)
	}
}

func TestService_CreateOffer_InvalidApplication(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	_, err := svc.CreateOffer(context.Background(), CreateOfferRequest{
		ApplicationID: uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected error for invalid application_id")
	}
}

func TestService_SubmitOffer_CreatesInstance(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	instanceID := uuid.New().String()
	engine := &fakeApprovalEngine{instanceID: instanceID, flowID: uuid.New().String()}
	svc.SetApprovalEngine(engine)
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	resp, err := svc.SubmitOffer(ctx, offer.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitOffer failed: %v", err)
	}
	if resp.Status != string(OfferStatusPendingApproval) {
		t.Errorf("expected status PENDING_APPROVAL, got %s", resp.Status)
	}
	if resp.ApprovalInstanceID != instanceID {
		t.Errorf("expected approval_instance_id %s, got %s", instanceID, resp.ApprovalInstanceID)
	}
	if engine.createdModule != "recruitment_offer" {
		t.Errorf("expected module 'recruitment_offer', got %q", engine.createdModule)
	}
	if engine.createdDocID != offer.ID {
		t.Errorf("expected document_id %s, got %s", offer.ID, engine.createdDocID)
	}
}

func TestService_SubmitOffer_AutoResolveFlow(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	flowID := uuid.New().String()
	engine := &fakeApprovalEngine{instanceID: uuid.New().String(), flowID: flowID}
	svc.SetApprovalEngine(engine)
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	resp, err := svc.SubmitOffer(ctx, offer.ID, "")
	if err != nil {
		t.Fatalf("SubmitOffer auto-resolve failed: %v", err)
	}
	if resp.Status != string(OfferStatusPendingApproval) {
		t.Errorf("expected status PENDING_APPROVAL, got %s", resp.Status)
	}
	if engine.createdFlowID != flowID {
		t.Errorf("expected auto-resolved flow_id %s, got %s", flowID, engine.createdFlowID)
	}
}

func TestService_HandleOfferApprovalStatusChange_Rejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	submitted, err := svc.SubmitOffer(ctx, offer.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitOffer failed: %v", err)
	}
	if err := svc.HandleOfferApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "REJECTED", "salary too high"); err != nil {
		t.Fatalf("HandleOfferApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetOfferByID(ctx, submitted.ID)
	if persisted.Status != string(OfferStatusRejected) {
		t.Errorf("expected status REJECTED after approval reject, got %s", persisted.Status)
	}
}

func TestService_OfferApproval_Sent_Accepted(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	submitted, err := svc.SubmitOffer(ctx, offer.ID, uuid.New().String())
	if err != nil {
		t.Fatalf("SubmitOffer failed: %v", err)
	}
	if err := svc.HandleOfferApprovalStatusChange(ctx, uuid.MustParse(submitted.ID), "APPROVED", ""); err != nil {
		t.Fatalf("HandleOfferApprovalStatusChange failed: %v", err)
	}
	persisted, _ := svc.GetOfferByID(ctx, submitted.ID)
	if persisted.Status != string(OfferStatusApproved) {
		t.Errorf("expected status APPROVED after approval, got %s", persisted.Status)
	}

	sent, err := svc.SendOffer(ctx, submitted.ID)
	if err != nil {
		t.Fatalf("SendOffer failed: %v", err)
	}
	if sent.Status != string(OfferStatusSent) {
		t.Errorf("expected status SENT, got %s", sent.Status)
	}

	accepted, err := svc.AcceptOffer(ctx, submitted.ID)
	if err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}
	if accepted.Status != string(OfferStatusAccepted) {
		t.Errorf("expected status ACCEPTED, got %s", accepted.Status)
	}

	app, err := svc.GetApplicationByID(ctx, accepted.ApplicationID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}
	if app.Status != string(CandStatusAccepted) {
		t.Errorf("expected application ACCEPTED, got %s", app.Status)
	}
	req, err := svc.GetRequisitionByID(ctx, app.RequisitionID)
	if err != nil {
		t.Fatalf("GetRequisitionByID failed: %v", err)
	}
	if req.SlotsFilled < 1 {
		t.Errorf("expected slots_filled >= 1 after offer accept, got %d", req.SlotsFilled)
	}
}

func TestService_AcceptOffer_NoDoubleIncrementSlotsFilled(t *testing.T) {
	// Idempotensi (G-3): aplikasi sudah ACCEPTED (jalur manual
	// UpdateApplicationStatus) lalu AcceptOffer — slots_filled TIDAK boleh
	// naik dua kali untuk satu kandidat.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	app, err := svc.GetApplicationByID(ctx, offer.ApplicationID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}

	// Aplikasi di-ACCEPTED manual dulu (jalur G-1 — sudah increment).
	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "ACCEPTED", "", ""); err != nil {
		t.Fatalf("UpdateApplicationStatus failed: %v", err)
	}
	reqBefore, _ := svc.GetRequisitionByID(ctx, app.RequisitionID)
	slotsBefore := reqBefore.SlotsFilled

	// Offer tetap bisa di-accept (status offer SENT) — tidak boleh menambah slot.
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)
	accepted, err := svc.AcceptOffer(ctx, offerID)
	if err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}
	if accepted.Status != string(OfferStatusAccepted) {
		t.Errorf("expected status ACCEPTED, got %s", accepted.Status)
	}
	reqAfter, _ := svc.GetRequisitionByID(ctx, app.RequisitionID)
	if reqAfter.SlotsFilled != slotsBefore {
		t.Errorf("expected slots_filled unchanged (%d), got %d", slotsBefore, reqAfter.SlotsFilled)
	}
}

func TestService_UpdateApplicationStatus_RepeatedAcceptedNoDoubleIncrement(t *testing.T) {
	// Idempotensi (G-5): panggilan manual PUT .../status ACCEPTED yang kedua
	// kali (aplikasi sudah ACCEPTED) adalah no-op transition dan TIDAK boleh
	// menambah slots_filled lagi.
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "repeat-accept@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "ACCEPTED", "", ""); err != nil {
		t.Fatalf("first ACCEPTED transition failed: %v", err)
	}
	reqAfterFirst, _ := svc.GetRequisitionByID(ctx, req.ID)
	slotsAfterFirst := reqAfterFirst.SlotsFilled
	if slotsAfterFirst != 1 {
		t.Fatalf("expected slots_filled=1 after first ACCEPTED call, got %d", slotsAfterFirst)
	}

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "ACCEPTED", "", ""); err != nil {
		t.Fatalf("second (no-op) ACCEPTED transition failed: %v", err)
	}
	reqAfterSecond, _ := svc.GetRequisitionByID(ctx, req.ID)
	if reqAfterSecond.SlotsFilled != slotsAfterFirst {
		t.Errorf("expected slots_filled unchanged after repeated ACCEPTED call (%d), got %d", slotsAfterFirst, reqAfterSecond.SlotsFilled)
	}
}

func TestService_AcceptOffer_WritesHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})

	// seedDraftOffer + approveAndSendOffer are existing test helpers in this
	// file (used by TestService_AcceptOffer_NoDoubleIncrementSlotsFilled)
	// that create a requisition+candidate+application+draft offer and drive
	// it to SENT via the approval flow — reuse them instead of re-deriving
	// the submit/approve/send sequence.
	offer := seedDraftOffer(t, svc, ctx)
	app, err := svc.GetApplicationByID(ctx, offer.ApplicationID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	if _, err := svc.AcceptOffer(ctx, offerID); err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	last := hist[len(hist)-1]
	if last.ToStage.Code != "ACCEPTED" {
		t.Errorf("expected last history to_stage ACCEPTED, got %s", last.ToStage.Code)
	}
}

func TestService_AcceptOffer_TransitionFails_SkipsSideEffects(t *testing.T) {
	// Regression (code review finding): if transitionApplicationStatus fails
	// inside AcceptOffer (e.g. application already in a terminal status
	// other than ACCEPTED, so the state machine rejects the -> ACCEPTED
	// transition), the slots_filled increment and the G-4 employee/movement
	// handoff must NOT run — wasAccepted alone is not a sufficient guard,
	// since it's computed before the (possibly failing) transition attempt.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	hire := &fakeEmployeeProvider{respID: uuid.New().String()}
	mov := &fakeMovementProvider{}
	svc.SetEmployeeProvider(hire)
	svc.SetMovementProvider(mov)
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	app, err := svc.GetApplicationByID(ctx, offer.ApplicationID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}

	// Force the application into a terminal status that is NOT ACCEPTED.
	// isValidStatusTransition rejects any transition out of a terminal
	// status, so transitionApplicationStatus(..., CandStatusAccepted, ...)
	// will fail inside AcceptOffer below.
	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "REJECTED", "", ""); err != nil {
		t.Fatalf("setup REJECTED transition failed: %v", err)
	}
	reqBefore, _ := svc.GetRequisitionByID(ctx, app.RequisitionID)
	slotsBefore := reqBefore.SlotsFilled

	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	// AcceptOffer itself must still succeed (offer state machine only cares
	// about OfferStatusSent) — it's the application-side side effects that
	// must be skipped.
	accepted, err := svc.AcceptOffer(ctx, offerID)
	if err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}
	if accepted.Status != string(OfferStatusAccepted) {
		t.Errorf("expected offer status ACCEPTED, got %s", accepted.Status)
	}

	appAfter, err := svc.GetApplicationByID(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}
	if appAfter.Status != "REJECTED" {
		t.Errorf("expected application status to remain REJECTED, got %s", appAfter.Status)
	}

	reqAfter, _ := svc.GetRequisitionByID(ctx, app.RequisitionID)
	if reqAfter.SlotsFilled != slotsBefore {
		t.Errorf("expected slots_filled unchanged (%d), got %d", slotsBefore, reqAfter.SlotsFilled)
	}

	if hire.called {
		t.Error("expected employee provider NOT called when application transition failed")
	}
	if mov.called {
		t.Error("expected movement provider NOT called when application transition failed")
	}
}

func TestService_AcceptOffer_Expired(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	past := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	offer := seedOffer(t, svc, ctx, past)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	_, err := svc.AcceptOffer(ctx, offerID)
	if err == nil {
		t.Fatal("expected error accepting expired offer")
	}
	persisted, _ := svc.GetOfferByID(ctx, offerID)
	if persisted.Status != string(OfferStatusExpired) {
		t.Errorf("expected status EXPIRED after expired accept attempt, got %s", persisted.Status)
	}
}

func TestService_RejectOffer(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	rejected, err := svc.RejectOffer(ctx, offerID)
	if err != nil {
		t.Fatalf("RejectOffer failed: %v", err)
	}
	if rejected.Status != string(OfferStatusRejected) {
		t.Errorf("expected status REJECTED, got %s", rejected.Status)
	}
}

func TestService_WithdrawOffer(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	withdrawn, err := svc.WithdrawOffer(ctx, offer.ID)
	if err != nil {
		t.Fatalf("WithdrawOffer failed: %v", err)
	}
	if withdrawn.Status != string(OfferStatusWithdrawn) {
		t.Errorf("expected status WITHDRAWN, got %s", withdrawn.Status)
	}
}

func TestService_UpdateOffer_OnlyDraft(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	_, err := svc.UpdateOffer(ctx, offerID, UpdateOfferRequest{})
	if err == nil {
		t.Fatal("expected error updating non-draft offer")
	}
}

func TestService_DeleteOffer_OnlyDraft(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	if err := svc.DeleteOffer(ctx, offerID); err == nil {
		t.Fatal("expected error deleting non-draft offer")
	}

	// Draft masih bisa dihapus.
	other := seedOfferWithEmail(t, svc, ctx, time.Now().AddDate(0, 0, 30).Format("2006-01-02"), "john.other@test.com")
	if err := svc.DeleteOffer(ctx, other.ID); err != nil {
		t.Errorf("expected draft offer deletable, got error: %v", err)
	}
}

// =========================================================================
// G-4 — Recruitment → Employee / Employee Movement
// =========================================================================

// fakeEmployeeProvider mencatat panggilan CreateHiredEmployee (G-4).
type fakeEmployeeProvider struct {
	callCount int
	called    bool
	input     EmployeeHireInput
	respID    string
	errResp   error
}

func (f *fakeEmployeeProvider) CreateHiredEmployee(ctx context.Context, in EmployeeHireInput) (string, error) {
	f.callCount++
	f.called = true
	f.input = in
	if f.errResp != nil {
		return "", f.errResp
	}
	return f.respID, nil
}

// fakeMovementProvider mencatat panggilan CreateHiredMovement (G-4).
type fakeMovementProvider struct {
	called  bool
	input   MovementHireInput
	errResp error
}

func (f *fakeMovementProvider) CreateHiredMovement(ctx context.Context, in MovementHireInput) error {
	f.called = true
	f.input = in
	return f.errResp
}

func TestService_AcceptOffer_External_CreatesEmployee(t *testing.T) {
	// G-4: offer eksternal diterima → Employee module dipanggil (bukan movement),
	// dengan application_id & data kandidat yang benar.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	hire := &fakeEmployeeProvider{respID: uuid.New().String()}
	mov := &fakeMovementProvider{}
	svc.SetEmployeeProvider(hire)
	svc.SetMovementProvider(mov)
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	if _, err := svc.AcceptOffer(ctx, offerID); err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}
	if !hire.called {
		t.Fatal("expected employee provider called for external hire")
	}
	if mov.called {
		t.Fatal("expected movement provider NOT called for external hire")
	}
	if hire.input.ApplicationID != offer.ApplicationID {
		t.Errorf("expected application_id %s, got %s", offer.ApplicationID, hire.input.ApplicationID)
	}
	if hire.input.Name == "" {
		t.Error("expected candidate name in hire input")
	}
	if hire.input.CandidateID == "" {
		t.Error("expected candidate_id in hire input")
	}
}

func TestService_AcceptOffer_Internal_CreatesMovement(t *testing.T) {
	// G-4: kandidat INTERNAL (employee_id terisi) → Employee Movement dipanggil
	// (bukan employee baru), dengan employee_id yang benar.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	hire := &fakeEmployeeProvider{respID: uuid.New().String()}
	mov := &fakeMovementProvider{}
	svc.SetEmployeeProvider(hire)
	svc.SetMovementProvider(mov)
	ctx := context.Background()

	internalEmpID := uuid.New().String()
	internalCand, err := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Budi", LastName: "Internal", Email: "budi.internal@test.com",
		CandidateType: strPtr("INTERNAL"),
		EmployeeID:    &internalEmpID,
	})
	if err != nil {
		t.Fatalf("CreateCandidate internal failed: %v", err)
	}
	if internalCand.CandidateType != "INTERNAL" {
		t.Errorf("expected candidate_type INTERNAL, got %s", internalCand.CandidateType)
	}
	if internalCand.EmployeeID != internalEmpID {
		t.Errorf("expected employee_id %s, got %s", internalEmpID, internalCand.EmployeeID)
	}

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "Internal Engineer",
	})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{
		RequisitionID: req.ID, CandidateID: internalCand.ID,
	})
	offer, err := svc.CreateOffer(ctx, CreateOfferRequest{
		ApplicationID: app.ID, EmploymentType: "FULL_TIME",
		StartDate: time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
	})
	if err != nil {
		t.Fatalf("CreateOffer failed: %v", err)
	}
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	if _, err := svc.AcceptOffer(ctx, offerID); err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}
	if !mov.called {
		t.Fatal("expected movement provider called for internal hire")
	}
	if hire.called {
		t.Fatal("expected employee provider NOT called for internal hire")
	}
	if mov.input.EmployeeID != internalEmpID {
		t.Errorf("expected movement employee_id %s, got %s", internalEmpID, mov.input.EmployeeID)
	}
	if mov.input.ApplicationID != app.ID {
		t.Errorf("expected movement application_id %s, got %s", app.ID, mov.input.ApplicationID)
	}
}

func TestService_AcceptOffer_NoProviders_StillSucceeds(t *testing.T) {
	// G-4 fail-safe: provider tidak di-wire → accept offer tetap berhasil.
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	accepted, err := svc.AcceptOffer(ctx, offerID)
	if err != nil {
		t.Fatalf("AcceptOffer failed without providers: %v", err)
	}
	if accepted.Status != string(OfferStatusAccepted) {
		t.Errorf("expected ACCEPTED, got %s", accepted.Status)
	}
}

func TestService_Candidate_InternalCreateAndUpdate(t *testing.T) {
	// G-4: kandidat bisa dibuat internal lalu di-update (candidate_type + employee_id).
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	empID := uuid.New().String()
	created, err := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Sari", LastName: "Internal", Email: "sari.internal@test.com",
		CandidateType: strPtr("INTERNAL"),
		EmployeeID:    &empID,
	})
	if err != nil {
		t.Fatalf("CreateCandidate failed: %v", err)
	}
	if created.CandidateType != "INTERNAL" {
		t.Errorf("expected INTERNAL, got %s", created.CandidateType)
	}

	// Update ke EXTERNAL membersihkan referensi employee_id (jalur handoff G-4
	// ditentukan candidate_type — referensi tidak bermakna untuk eksternal).
	external := "EXTERNAL"
	updated, err := svc.UpdateCandidate(ctx, created.ID, UpdateCandidateRequest{
		CandidateType: &external,
	})
	if err != nil {
		t.Fatalf("UpdateCandidate failed: %v", err)
	}
	if updated.CandidateType != "EXTERNAL" {
		t.Errorf("expected EXTERNAL after update, got %s", updated.CandidateType)
	}
	if updated.EmployeeID != "" {
		t.Errorf("expected employee_id cleared on EXTERNAL, got %s", updated.EmployeeID)
	}

	// Default create = EXTERNAL (tanpa candidate_type).
	def, err := svc.CreateCandidate(ctx, CreateCandidateRequest{
		FirstName: "Dewi", LastName: "Default", Email: "dewi.default@test.com",
	})
	if err != nil {
		t.Fatalf("CreateCandidate default failed: %v", err)
	}
	if def.CandidateType != "EXTERNAL" {
		t.Errorf("expected default EXTERNAL, got %s", def.CandidateType)
	}
}

func TestService_AcceptOffer_SecondOffer_NoDuplicateHandoff(t *testing.T) {
	// G-4 idempotensi: aplikasi sudah ACCEPTED (offer pertama) — offer kedua di
	// aplikasi yang sama yang di-accept TIDAK memanggil employee provider lagi
	// (guard transisi status, mirror bug slots_filled G-3).
	svc, cleanup := newTestService()
	defer cleanup()
	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})
	hire := &fakeEmployeeProvider{respID: uuid.New().String()}
	svc.SetEmployeeProvider(hire)
	ctx := context.Background()

	offer := seedDraftOffer(t, svc, ctx)
	firstID := approveAndSendOffer(t, svc, ctx, offer.ID)
	if _, err := svc.AcceptOffer(ctx, firstID); err != nil {
		t.Fatalf("first AcceptOffer failed: %v", err)
	}
	if hire.called != true {
		t.Fatal("expected employee provider called on first accept")
	}

	// Offer kedua di aplikasi yang sama (aplikasi sudah ACCEPTED).
	second, err := svc.CreateOffer(ctx, CreateOfferRequest{
		ApplicationID: offer.ApplicationID, EmploymentType: "FULL_TIME",
		StartDate: time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
	})
	if err != nil {
		t.Fatalf("CreateOffer second failed: %v", err)
	}
	secondID := approveAndSendOffer(t, svc, ctx, second.ID)
	if _, err := svc.AcceptOffer(ctx, secondID); err != nil {
		t.Fatalf("second AcceptOffer failed: %v", err)
	}

	if hire.callCount != 1 {
		t.Fatalf("expected employee provider called exactly once, got %d", hire.callCount)
	}
}
