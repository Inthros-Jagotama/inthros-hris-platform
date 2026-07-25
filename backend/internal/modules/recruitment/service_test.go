package recruitment

import (
	"context"
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
