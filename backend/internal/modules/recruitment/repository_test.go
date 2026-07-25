package recruitment

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func newTestRepository() (*Repository, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	return repo, func() { cleanup() }
}

// =========================================================================
// Job Requisition Repository Tests
// =========================================================================

func TestRepo_CreateRequisition(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{
		OrganizationID: uuid.New(),
		Title:          "Software Engineer",
		Status:         ReqStatusDraft,
	}
	if err := repo.CreateRequisition(ctx, r); err != nil {
		t.Fatalf("CreateRequisition failed: %v", err)
	}
	if r.ID == uuid.Nil {
		t.Fatal("expected ID to be generated")
	}
}

func TestRepo_FindRequisitionByID_NotFound(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	_, err := repo.FindRequisitionByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent requisition")
	}
}

func TestRepo_ListRequisitions(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	orgID := uuid.New()
	otherOrgID := uuid.New()

	for i := 0; i < 3; i++ {
		repo.CreateRequisition(ctx, &JobRequisition{
			OrganizationID: orgID, Title: "Req A", Status: ReqStatusOpen,
		})
	}
	repo.CreateRequisition(ctx, &JobRequisition{
		OrganizationID: otherOrgID, Title: "Req B", Status: ReqStatusOpen,
	})

	// Filter by org
	list, total, err := repo.ListRequisitions(ctx, &orgID, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListRequisitions failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 for org, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 results, got %d", len(list))
	}
}

func TestRepo_UpdateRequisition(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	r := &JobRequisition{
		OrganizationID: uuid.New(), Title: "Original", Status: ReqStatusDraft,
	}
	repo.CreateRequisition(ctx, r)

	r.Title = "Updated"
	if err := repo.UpdateRequisition(ctx, r); err != nil {
		t.Fatalf("UpdateRequisition failed: %v", err)
	}
	found, _ := repo.FindRequisitionByID(ctx, r.ID)
	if found.Title != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", found.Title)
	}
}

func TestRepo_DeleteRequisition(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	r := &JobRequisition{
		OrganizationID: uuid.New(), Title: "To Delete", Status: ReqStatusDraft,
	}
	repo.CreateRequisition(ctx, r)

	if err := repo.DeleteRequisition(ctx, r.ID); err != nil {
		t.Fatalf("DeleteRequisition failed: %v", err)
	}
	// Verify deleted
	_, err := repo.FindRequisitionByID(ctx, r.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestRepo_DeleteRequisition_NotFound(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	err := repo.DeleteRequisition(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error deleting non-existent requisition")
	}
}

// =========================================================================
// Candidate Repository Tests
// =========================================================================

func TestRepo_CreateCandidate(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	c := &Candidate{
		FirstName: "John", LastName: "Doe", Email: "john@example.com",
	}
	if err := repo.CreateCandidate(context.Background(), c); err != nil {
		t.Fatalf("CreateCandidate failed: %v", err)
	}
}

func TestRepo_FindCandidateByEmail_NotFound(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	found, err := repo.FindCandidateByEmail(context.Background(), "nonexist@test.com")
	if err != nil {
		t.Fatalf("FindCandidateByEmail failed: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil for non-existent email")
	}
}

func TestRepo_FindCandidateByEmail_Found(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	repo.CreateCandidate(ctx, &Candidate{
		FirstName: "Jane", LastName: "Doe", Email: "jane@example.com",
	})
	found, err := repo.FindCandidateByEmail(ctx, "jane@example.com")
	if err != nil {
		t.Fatalf("FindCandidateByEmail failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected candidate, got nil")
	}
}

func TestRepo_ListCandidates_Search(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	repo.CreateCandidate(ctx, &Candidate{FirstName: "Alice", LastName: "Wonder", Email: "alice@test.com"})
	repo.CreateCandidate(ctx, &Candidate{FirstName: "Bob", LastName: "Builder", Email: "bob@test.com"})

	search := "alice"
	list, total, err := repo.ListCandidates(ctx, 1, 10, &search)
	if err != nil {
		t.Fatalf("ListCandidates failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 result for search 'alice', got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item, got %d", len(list))
	}
}

func TestRepo_UpdateCandidate(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	c := &Candidate{FirstName: "Old", LastName: "Name", Email: "old@test.com"}
	repo.CreateCandidate(ctx, c)

	c.FirstName = "New"
	repo.UpdateCandidate(ctx, c)
	found, _ := repo.FindCandidateByID(ctx, c.ID)
	if found.FirstName != "New" {
		t.Errorf("expected 'New', got '%s'", found.FirstName)
	}
}

func TestRepo_DeleteCandidate(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	c := &Candidate{FirstName: "Delete", LastName: "Me", Email: "del@test.com"}
	repo.CreateCandidate(ctx, c)

	if err := repo.DeleteCandidate(ctx, c.ID); err != nil {
		t.Fatalf("DeleteCandidate failed: %v", err)
	}
}

// =========================================================================
// Job Application Repository Tests
// =========================================================================

func TestRepo_CreateApplication(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "App", LastName: "Test", Email: "app@test.com"}
	repo.CreateCandidate(ctx, c)

	a := &JobApplication{
		RequisitionID: r.ID,
		CandidateID:   c.ID,
		Status:        CandStatusNew,
	}
	if err := repo.CreateApplication(ctx, a); err != nil {
		t.Fatalf("CreateApplication failed: %v", err)
	}
}

func TestRepo_ListApplications_FilterByStatus(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "App2", LastName: "Test", Email: "app2@test.com"}
	repo.CreateCandidate(ctx, c)

	status := CandStatusNew
	repo.CreateApplication(ctx, &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: status})

	filteredStatus := "NEW"
	list, total, err := repo.ListApplications(ctx, nil, nil, &filteredStatus, 1, 10)
	if err != nil {
		t.Fatalf("ListApplications failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 application with status NEW, got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item, got %d", len(list))
	}
}

func TestRepo_DeleteApplication(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "AppDel", LastName: "Test", Email: "appdel@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a)

	if err := repo.DeleteApplication(ctx, a.ID); err != nil {
		t.Fatalf("DeleteApplication failed: %v", err)
	}
}

// =========================================================================
// Interview Repository Tests
// =========================================================================

func TestRepo_CreateInterview(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Int", LastName: "Test", Email: "int@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a)

	i := &Interview{
		ApplicationID: a.ID,
		InterviewerID: uuid.New(),
		Stage:         "TECHNICAL",
		Status:        IntStatusScheduled,
	}
	if err := repo.CreateInterview(ctx, i); err != nil {
		t.Fatalf("CreateInterview failed: %v", err)
	}
}

func TestRepo_ListInterviews_FilterByApp(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c1 := &Candidate{FirstName: "Int1", LastName: "A", Email: "int1@test.com"}
	repo.CreateCandidate(ctx, c1)
	c2 := &Candidate{FirstName: "Int2", LastName: "B", Email: "int2@test.com"}
	repo.CreateCandidate(ctx, c2)
	a1 := &JobApplication{RequisitionID: r.ID, CandidateID: c1.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a1)
	a2 := &JobApplication{RequisitionID: r.ID, CandidateID: c2.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a2)

	repo.CreateInterview(ctx, &Interview{ApplicationID: a1.ID, InterviewerID: uuid.New(), Stage: "HR", Status: IntStatusScheduled})
	repo.CreateInterview(ctx, &Interview{ApplicationID: a1.ID, InterviewerID: uuid.New(), Stage: "TECHNICAL", Status: IntStatusScheduled})
	repo.CreateInterview(ctx, &Interview{ApplicationID: a2.ID, InterviewerID: uuid.New(), Stage: "HR", Status: IntStatusScheduled})

	list, total, err := repo.ListInterviews(ctx, &a1.ID, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListInterviews failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 interviews for a1, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

// =========================================================================
// Onboarding Task Template Repository Tests
// =========================================================================

func TestRepo_CreateOnboardingTaskTemplate(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	tmpl := &OnboardingTaskTemplate{
		Name: "IT Setup", Category: "IT", DayOffset: -7, IsMandatory: true,
	}
	if err := repo.CreateOnboardingTaskTemplate(context.Background(), tmpl); err != nil {
		t.Fatalf("CreateOnboardingTaskTemplate failed: %v", err)
	}
}

func TestRepo_ListOnboardingTaskTemplates_FilterByCategory(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	repo.CreateOnboardingTaskTemplate(ctx, &OnboardingTaskTemplate{Name: "IT Setup", Category: "IT", IsMandatory: true})
	repo.CreateOnboardingTaskTemplate(ctx, &OnboardingTaskTemplate{Name: "HR Setup", Category: "HR", IsMandatory: true})
	repo.CreateOnboardingTaskTemplate(ctx, &OnboardingTaskTemplate{Name: "Legal Setup", Category: "LEGAL", IsMandatory: true})

	cat := "IT"
	list, total, err := repo.ListOnboardingTaskTemplates(ctx, &cat, 1, 10)
	if err != nil {
		t.Fatalf("ListOnboardingTaskTemplates failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected 1 IT template, got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item, got %d", len(list))
	}
}

// =========================================================================
// Employee Onboarding Repository Tests
// =========================================================================

func TestRepo_CreateEmployeeOnboarding(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Onb", LastName: "Test", Email: "onb@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)

	o := &EmployeeOnboarding{
		EmployeeID:    uuid.New(),
		ApplicationID: a.ID,
		StartDate:     "2026-08-01",
		Status:        "PENDING",
	}
	if err := repo.CreateEmployeeOnboarding(ctx, o); err != nil {
		t.Fatalf("CreateEmployeeOnboarding failed: %v", err)
	}
}

func TestRepo_FindEmployeeOnboardingByEmployeeID(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()
	empID := uuid.New()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Onb2", LastName: "Test", Email: "onb2@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)

	repo.CreateEmployeeOnboarding(ctx, &EmployeeOnboarding{
		EmployeeID: empID, ApplicationID: a.ID, StartDate: "2026-08-01", Status: "PENDING",
	})

	found, err := repo.FindEmployeeOnboardingByEmployeeID(ctx, empID)
	if err != nil {
		t.Fatalf("FindEmployeeOnboardingByEmployeeID failed: %v", err)
	}
	if found == nil {
		t.Fatal("expected onboarding, got nil")
	}
}

func TestRepo_ListEmployeeOnboardings_FilterByStatus(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	createOnboarding := func(status string) {
		r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
		repo.CreateRequisition(ctx, r)
		c := &Candidate{FirstName: "Onb3", LastName: status, Email: "onb3_" + status + "@test.com"}
		repo.CreateCandidate(ctx, c)
		a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
		repo.CreateApplication(ctx, a)
		repo.CreateEmployeeOnboarding(ctx, &EmployeeOnboarding{
			EmployeeID: uuid.New(), ApplicationID: a.ID, StartDate: "2026-08-01", Status: status,
		})
	}
	createOnboarding("PENDING")
	createOnboarding("PENDING")
	createOnboarding("COMPLETED")

	status := "PENDING"
	list, total, err := repo.ListEmployeeOnboardings(ctx, &status, 1, 10)
	if err != nil {
		t.Fatalf("ListEmployeeOnboardings failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 PENDING onboardings, got %d", total)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

// =========================================================================
// Onboarding Task Item Repository Tests
// =========================================================================

func TestRepo_CreateOnboardingTaskItem(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	// Create prerequisite
	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Item", LastName: "Test", Email: "item@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)
	o := &EmployeeOnboarding{EmployeeID: uuid.New(), ApplicationID: a.ID, StartDate: "2026-08-01", Status: "PENDING"}
	repo.CreateEmployeeOnboarding(ctx, o)

	item := &OnboardingTaskItem{
		EmployeeOnboardingID: o.ID,
		Name:                "Task 1",
		IsCompleted:         false,
	}
	if err := repo.CreateOnboardingTaskItem(ctx, item); err != nil {
		t.Fatalf("CreateOnboardingTaskItem failed: %v", err)
	}
}

func TestRepo_ListOnboardingTaskItems(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Items2", LastName: "Test", Email: "items2@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)
	o := &EmployeeOnboarding{EmployeeID: uuid.New(), ApplicationID: a.ID, StartDate: "2026-08-01", Status: "PENDING"}
	repo.CreateEmployeeOnboarding(ctx, o)

	repo.CreateOnboardingTaskItem(ctx, &OnboardingTaskItem{EmployeeOnboardingID: o.ID, Name: "Task 1", IsCompleted: false})
	repo.CreateOnboardingTaskItem(ctx, &OnboardingTaskItem{EmployeeOnboardingID: o.ID, Name: "Task 2", IsCompleted: false})

	items, err := repo.ListOnboardingTaskItems(ctx, o.ID)
	if err != nil {
		t.Fatalf("ListOnboardingTaskItems failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestRepo_UpdateOnboardingTaskItem(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "ItemUpd", LastName: "Test", Email: "itemupd@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)
	o := &EmployeeOnboarding{EmployeeID: uuid.New(), ApplicationID: a.ID, StartDate: "2026-08-01", Status: "PENDING"}
	repo.CreateEmployeeOnboarding(ctx, o)
	item := &OnboardingTaskItem{EmployeeOnboardingID: o.ID, Name: "Initial", IsCompleted: false}
	repo.CreateOnboardingTaskItem(ctx, item)

	item.Name = "Updated"
	repo.UpdateOnboardingTaskItem(ctx, item)
	found, _ := repo.FindOnboardingTaskItemByID(ctx, item.ID)
	if found.Name != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", found.Name)
	}
}

func TestRepo_DeleteOnboardingTaskItem(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusDraft}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "ItemDel", LastName: "Test", Email: "itemdel@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusAccepted}
	repo.CreateApplication(ctx, a)
	o := &EmployeeOnboarding{EmployeeID: uuid.New(), ApplicationID: a.ID, StartDate: "2026-08-01", Status: "PENDING"}
	repo.CreateEmployeeOnboarding(ctx, o)
	item := &OnboardingTaskItem{EmployeeOnboardingID: o.ID, Name: "To Delete", IsCompleted: false}
	repo.CreateOnboardingTaskItem(ctx, item)

	if err := repo.DeleteOnboardingTaskItem(ctx, item.ID); err != nil {
		t.Fatalf("DeleteOnboardingTaskItem failed: %v", err)
	}
}

func TestRepo_DeleteOnboardingTaskItem_NotFound(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	err := repo.DeleteOnboardingTaskItem(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error deleting non-existent task item")
	}
}
