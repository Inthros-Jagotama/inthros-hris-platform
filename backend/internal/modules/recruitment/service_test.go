package recruitment

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger := zap.NewNop()
	svc := NewService(repo, logger)
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
