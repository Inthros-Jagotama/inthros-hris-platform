package recruitment

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/modules/competency"
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

func TestRepository_CreateAndFindStageByCode(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	stage := &RecruitmentStage{Code: "NEW", Name: "New Application", SortOrder: 1}
	if err := repo.CreateStage(ctx, stage); err != nil {
		t.Fatalf("CreateStage failed: %v", err)
	}

	found, err := repo.FindStageByCode(ctx, "NEW")
	if err != nil {
		t.Fatalf("FindStageByCode failed: %v", err)
	}
	if found.ID != stage.ID {
		t.Errorf("expected id %s, got %s", stage.ID, found.ID)
	}
}

func TestRepository_ListStages(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	repo.CreateStage(ctx, &RecruitmentStage{Code: "NEW", Name: "New", SortOrder: 1})
	repo.CreateStage(ctx, &RecruitmentStage{Code: "SCREENED", Name: "Screened", SortOrder: 2})

	list, err := repo.ListStages(ctx)
	if err != nil {
		t.Fatalf("ListStages failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 stages, got %d", len(list))
	}
}

func TestRepository_CreateAndListStageHistory(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	req := &JobRequisition{OrganizationID: uuid.New(), Title: "Engineer"}
	repo.CreateRequisition(ctx, req)
	cand := &Candidate{FirstName: "A", LastName: "B", Email: "ab@test.com"}
	repo.CreateCandidate(ctx, cand)
	app := &JobApplication{RequisitionID: req.ID, CandidateID: cand.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, app)

	newStage := &RecruitmentStage{Code: "NEW", Name: "New", SortOrder: 1}
	repo.CreateStage(ctx, newStage)
	screenedStage := &RecruitmentStage{Code: "SCREENED", Name: "Screened", SortOrder: 2}
	repo.CreateStage(ctx, screenedStage)

	h := &ApplicationStageHistory{
		ApplicationID: app.ID,
		FromStageID:   &newStage.ID,
		ToStageID:     screenedStage.ID,
		ChangedAt:     1000,
	}
	if err := repo.CreateStageHistory(ctx, h); err != nil {
		t.Fatalf("CreateStageHistory failed: %v", err)
	}

	list, err := repo.ListStageHistoryByApplication(ctx, app.ID)
	if err != nil {
		t.Fatalf("ListStageHistoryByApplication failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 history row, got %d", len(list))
	}
	if list[0].ToStageID != screenedStage.ID {
		t.Errorf("expected to_stage_id %s, got %s", screenedStage.ID, list[0].ToStageID)
	}
}

// =========================================================================
// Candidate Education / Work Experience Repository Tests (G-6)
// =========================================================================

func TestRepository_CreateAndFindCandidateEducation(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Edu", LastName: "Test", Email: "edutest@test.com"}
	repo.CreateCandidate(ctx, cand)

	edu := &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Universitas Test"}
	if err := repo.CreateCandidateEducation(ctx, edu); err != nil {
		t.Fatalf("CreateCandidateEducation failed: %v", err)
	}

	found, err := repo.FindCandidateEducationByID(ctx, edu.ID)
	if err != nil {
		t.Fatalf("FindCandidateEducationByID failed: %v", err)
	}
	if found.InstitutionName != "Universitas Test" {
		t.Errorf("expected institution 'Universitas Test', got %s", found.InstitutionName)
	}
}

func TestRepository_ListCandidateEducations(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Edu", Email: "listedu@test.com"}
	repo.CreateCandidate(ctx, cand)
	repo.CreateCandidateEducation(ctx, &CandidateEducation{CandidateID: cand.ID, InstitutionName: "SMA 1"})
	repo.CreateCandidateEducation(ctx, &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Universitas A"})

	list, err := repo.ListCandidateEducations(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateEducations failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 education rows, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateEducation(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Edu", Email: "updedu@test.com"}
	repo.CreateCandidate(ctx, cand)
	edu := &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Original"}
	repo.CreateCandidateEducation(ctx, edu)

	edu.InstitutionName = "Updated"
	if err := repo.UpdateCandidateEducation(ctx, edu); err != nil {
		t.Fatalf("UpdateCandidateEducation failed: %v", err)
	}
	found, _ := repo.FindCandidateEducationByID(ctx, edu.ID)
	if found.InstitutionName != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.InstitutionName)
	}

	if err := repo.DeleteCandidateEducation(ctx, edu.ID); err != nil {
		t.Fatalf("DeleteCandidateEducation failed: %v", err)
	}
	if _, err := repo.FindCandidateEducationByID(ctx, edu.ID); err == nil {
		t.Error("expected error finding deleted education, got nil")
	}
}

func TestRepository_CreateAndListCandidateWorkExperience(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Exp", LastName: "Test", Email: "exptest@test.com"}
	repo.CreateCandidate(ctx, cand)

	exp := &CandidateWorkExperience{CandidateID: cand.ID, CompanyName: "Acme Corp", JobTitle: "Engineer", StartDate: "2020-01-01"}
	if err := repo.CreateCandidateWorkExperience(ctx, exp); err != nil {
		t.Fatalf("CreateCandidateWorkExperience failed: %v", err)
	}

	list, err := repo.ListCandidateWorkExperiences(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateWorkExperiences failed: %v", err)
	}
	if len(list) != 1 || list[0].CompanyName != "Acme Corp" {
		t.Errorf("expected 1 row with company 'Acme Corp', got %+v", list)
	}
}

func TestRepository_CreateAndFindCandidateSkill(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "Skill", LastName: "Test", Email: "skilltest@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp := &competency.Competency{Name: "Go Programming"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	skill := &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp.ID}
	if err := repo.CreateCandidateSkill(ctx, skill); err != nil {
		t.Fatalf("CreateCandidateSkill failed: %v", err)
	}

	found, err := repo.FindCandidateSkillByID(ctx, skill.ID)
	if err != nil {
		t.Fatalf("FindCandidateSkillByID failed: %v", err)
	}
	if found.CompetencyID != comp.ID {
		t.Errorf("expected competency_id %s, got %s", comp.ID, found.CompetencyID)
	}
}

func TestRepository_ListCandidateSkills(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Skill", Email: "listskill@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp1 := &competency.Competency{Name: "Go"}
	comp2 := &competency.Competency{Name: "SQL"}
	if err := db.Create(comp1).Error; err != nil {
		t.Fatalf("failed to seed competency comp1: %v", err)
	}
	if err := db.Create(comp2).Error; err != nil {
		t.Fatalf("failed to seed competency comp2: %v", err)
	}

	repo.CreateCandidateSkill(ctx, &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp1.ID})
	repo.CreateCandidateSkill(ctx, &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp2.ID})

	list, err := repo.ListCandidateSkills(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateSkills failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 skills, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateSkill(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Skill", Email: "updskill@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}
	skill := &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp.ID}
	repo.CreateCandidateSkill(ctx, skill)

	level := 4
	skill.Level = &level
	if err := repo.UpdateCandidateSkill(ctx, skill); err != nil {
		t.Fatalf("UpdateCandidateSkill failed: %v", err)
	}
	found, _ := repo.FindCandidateSkillByID(ctx, skill.ID)
	if found.Level == nil || *found.Level != level {
		t.Errorf("expected level %d, got %v", level, found.Level)
	}

	if err := repo.DeleteCandidateSkill(ctx, skill.ID); err != nil {
		t.Fatalf("DeleteCandidateSkill failed: %v", err)
	}
	if _, err := repo.FindCandidateSkillByID(ctx, skill.ID); err == nil {
		t.Error("expected error finding deleted skill, got nil")
	}
}

func TestRepository_CreateAndListCandidateCertification(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Cert", LastName: "Test", Email: "certtest@test.com"}
	repo.CreateCandidate(ctx, cand)

	cert := &CandidateCertification{CandidateID: cand.ID, Name: "AWS Certified Solutions Architect"}
	if err := repo.CreateCandidateCertification(ctx, cert); err != nil {
		t.Fatalf("CreateCandidateCertification failed: %v", err)
	}

	list, err := repo.ListCandidateCertifications(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateCertifications failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "AWS Certified Solutions Architect" {
		t.Errorf("expected 1 cert named 'AWS Certified Solutions Architect', got %+v", list)
	}
}

func TestRepository_CreateAndFindCandidateDocument(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Doc", LastName: "Test", Email: "doctest@test.com"}
	repo.CreateCandidate(ctx, cand)

	doc := &CandidateDocument{CandidateID: cand.ID, DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/uploads/attachments/abc.pdf"}
	if err := repo.CreateCandidateDocument(ctx, doc); err != nil {
		t.Fatalf("CreateCandidateDocument failed: %v", err)
	}

	found, err := repo.FindCandidateDocumentByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("FindCandidateDocumentByID failed: %v", err)
	}
	if found.Name != "resume.pdf" {
		t.Errorf("expected name 'resume.pdf', got %s", found.Name)
	}
}

func TestRepository_ListCandidateDocuments(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Doc", Email: "listdoc@test.com"}
	repo.CreateCandidate(ctx, cand)
	repo.CreateCandidateDocument(ctx, &CandidateDocument{CandidateID: cand.ID, DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/u/a.pdf"})
	repo.CreateCandidateDocument(ctx, &CandidateDocument{CandidateID: cand.ID, DocumentType: "PORTFOLIO", Name: "portfolio.pdf", FileURL: "/u/b.pdf"})

	list, err := repo.ListCandidateDocuments(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateDocuments failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 documents, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateDocument(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Doc", Email: "upddoc@test.com"}
	repo.CreateCandidate(ctx, cand)
	doc := &CandidateDocument{CandidateID: cand.ID, DocumentType: "OTHER", Name: "Original", FileURL: "/u/c.pdf"}
	repo.CreateCandidateDocument(ctx, doc)

	doc.Name = "Updated"
	if err := repo.UpdateCandidateDocument(ctx, doc); err != nil {
		t.Fatalf("UpdateCandidateDocument failed: %v", err)
	}
	found, _ := repo.FindCandidateDocumentByID(ctx, doc.ID)
	if found.Name != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.Name)
	}

	if err := repo.DeleteCandidateDocument(ctx, doc.ID); err != nil {
		t.Fatalf("DeleteCandidateDocument failed: %v", err)
	}
	if _, err := repo.FindCandidateDocumentByID(ctx, doc.ID); err == nil {
		t.Error("expected error finding deleted document, got nil")
	}
}

// =========================================================================
// Candidate Consents (G-6) — append-only, no Update/Delete
// =========================================================================

func TestRepository_CreateAndListCandidateConsents(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Consent", LastName: "Test", Email: "consenttest@test.com"}
	repo.CreateCandidate(ctx, cand)

	consent := &CandidateConsent{CandidateID: cand.ID, Action: "GRANTED", ChangedAt: 1000}
	if err := repo.CreateCandidateConsent(ctx, consent); err != nil {
		t.Fatalf("CreateCandidateConsent failed: %v", err)
	}

	list, err := repo.ListCandidateConsents(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateConsents failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 consent entry, got %d", len(list))
	}
	if list[0].Action != "GRANTED" {
		t.Errorf("expected action 'GRANTED', got %s", list[0].Action)
	}
}

func TestRepository_ListCandidateConsents_OrderedByChangedAt(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Order", LastName: "Consent", Email: "orderconsent@test.com"}
	repo.CreateCandidate(ctx, cand)

	repo.CreateCandidateConsent(ctx, &CandidateConsent{CandidateID: cand.ID, Action: "GRANTED", ChangedAt: 2000})
	repo.CreateCandidateConsent(ctx, &CandidateConsent{CandidateID: cand.ID, Action: "REVOKED", ChangedAt: 1000})

	list, err := repo.ListCandidateConsents(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateConsents failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	if list[0].ChangedAt != 1000 || list[1].ChangedAt != 2000 {
		t.Errorf("expected ascending changed_at order [1000, 2000], got [%d, %d]", list[0].ChangedAt, list[1].ChangedAt)
	}
}

func TestRepository_CreateAndFindApplicationScreening(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Scr", LastName: "Test", Email: "scr@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a)

	sc := &ApplicationScreening{ApplicationID: a.ID, Result: "PASS"}
	if err := repo.CreateApplicationScreening(ctx, sc); err != nil {
		t.Fatalf("CreateApplicationScreening failed: %v", err)
	}

	found, err := repo.FindApplicationScreeningByID(ctx, sc.ID)
	if err != nil {
		t.Fatalf("FindApplicationScreeningByID failed: %v", err)
	}
	if found.Result != "PASS" {
		t.Errorf("expected result PASS, got %s", found.Result)
	}
}

func TestRepository_ListApplicationScreenings(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "List", LastName: "Scr", Email: "listscr@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a)

	repo.CreateApplicationScreening(ctx, &ApplicationScreening{ApplicationID: a.ID, Result: "HOLD"})
	repo.CreateApplicationScreening(ctx, &ApplicationScreening{ApplicationID: a.ID, Result: "PASS"})

	list, err := repo.ListApplicationScreenings(ctx, a.ID)
	if err != nil {
		t.Fatalf("ListApplicationScreenings failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 screenings, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteApplicationScreening(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Upd", LastName: "Scr", Email: "updscr@test.com"}
	repo.CreateCandidate(ctx, c)
	a := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, a)
	sc := &ApplicationScreening{ApplicationID: a.ID, Result: "HOLD"}
	repo.CreateApplicationScreening(ctx, sc)

	sc.Result = "FAIL"
	if err := repo.UpdateApplicationScreening(ctx, sc); err != nil {
		t.Fatalf("UpdateApplicationScreening failed: %v", err)
	}
	found, _ := repo.FindApplicationScreeningByID(ctx, sc.ID)
	if found.Result != "FAIL" {
		t.Errorf("expected FAIL, got %s", found.Result)
	}

	if err := repo.DeleteApplicationScreening(ctx, sc.ID); err != nil {
		t.Fatalf("DeleteApplicationScreening failed: %v", err)
	}
	if _, err := repo.FindApplicationScreeningByID(ctx, sc.ID); err == nil {
		t.Error("expected error finding deleted screening, got nil")
	}
}

func TestRepository_CreateAndFindAssessment(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	a := &RecruitmentAssessment{Name: "Technical Test Batch March", Type: "TECHNICAL"}
	if err := repo.CreateAssessment(ctx, a); err != nil {
		t.Fatalf("CreateAssessment failed: %v", err)
	}

	found, err := repo.FindAssessmentByID(ctx, a.ID)
	if err != nil {
		t.Fatalf("FindAssessmentByID failed: %v", err)
	}
	if found.Name != "Technical Test Batch March" {
		t.Errorf("expected name 'Technical Test Batch March', got %s", found.Name)
	}
}

func TestRepository_ListAssessments(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	repo.CreateAssessment(ctx, &RecruitmentAssessment{Name: "Assessment A", Type: "CODING"})
	repo.CreateAssessment(ctx, &RecruitmentAssessment{Name: "Assessment B", Type: "COGNITIVE"})

	list, err := repo.ListAssessments(ctx)
	if err != nil {
		t.Fatalf("ListAssessments failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 assessments, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteAssessment(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	a := &RecruitmentAssessment{Name: "Original", Type: "OTHER"}
	repo.CreateAssessment(ctx, a)

	a.Name = "Updated"
	if err := repo.UpdateAssessment(ctx, a); err != nil {
		t.Fatalf("UpdateAssessment failed: %v", err)
	}
	found, _ := repo.FindAssessmentByID(ctx, a.ID)
	if found.Name != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.Name)
	}

	if err := repo.DeleteAssessment(ctx, a.ID); err != nil {
		t.Fatalf("DeleteAssessment failed: %v", err)
	}
	if _, err := repo.FindAssessmentByID(ctx, a.ID); err == nil {
		t.Error("expected error finding deleted assessment, got nil")
	}
}

func TestRepository_CreateAndListAssessmentParticipants(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Part", LastName: "Test", Email: "part@test.com"}
	repo.CreateCandidate(ctx, c)
	app := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, app)

	assess := &RecruitmentAssessment{Name: "Batch", Type: "CODING"}
	repo.CreateAssessment(ctx, assess)

	p := &AssessmentParticipant{AssessmentID: assess.ID, ApplicationID: app.ID}
	if err := repo.CreateAssessmentParticipant(ctx, p); err != nil {
		t.Fatalf("CreateAssessmentParticipant failed: %v", err)
	}

	list, err := repo.ListAssessmentParticipants(ctx, assess.ID)
	if err != nil {
		t.Fatalf("ListAssessmentParticipants failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 participant, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteAssessmentParticipant(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Upd", LastName: "Part", Email: "updpart@test.com"}
	repo.CreateCandidate(ctx, c)
	app := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, app)
	assess := &RecruitmentAssessment{Name: "Batch", Type: "CODING"}
	repo.CreateAssessment(ctx, assess)
	p := &AssessmentParticipant{AssessmentID: assess.ID, ApplicationID: app.ID}
	repo.CreateAssessmentParticipant(ctx, p)

	p.Status = "COMPLETED"
	if err := repo.UpdateAssessmentParticipant(ctx, p); err != nil {
		t.Fatalf("UpdateAssessmentParticipant failed: %v", err)
	}
	found, _ := repo.FindAssessmentParticipantByID(ctx, p.ID)
	if found.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", found.Status)
	}

	if err := repo.DeleteAssessmentParticipant(ctx, p.ID); err != nil {
		t.Fatalf("DeleteAssessmentParticipant failed: %v", err)
	}
	if _, err := repo.FindAssessmentParticipantByID(ctx, p.ID); err == nil {
		t.Error("expected error finding deleted participant, got nil")
	}
}

func newTestInterviewForScorecard(t *testing.T, repo *Repository, ctx context.Context, email string) *Interview {
	t.Helper()
	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	c := &Candidate{FirstName: "Iv", LastName: "Test", Email: email}
	repo.CreateCandidate(ctx, c)
	app := &JobApplication{RequisitionID: r.ID, CandidateID: c.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, app)
	iv := &Interview{ApplicationID: app.ID, InterviewerID: uuid.New(), Stage: "FIRST_INTERVIEW", Status: IntStatusScheduled}
	repo.CreateInterview(ctx, iv)
	return iv
}

func TestRepository_CreateAndListInterviewers(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	iv := newTestInterviewForScorecard(t, repo, ctx, "ivrepo1@test.com")

	interviewer := &Interviewer{InterviewID: iv.ID, EmployeeID: uuid.New(), Role: "HR"}
	if err := repo.CreateInterviewer(ctx, interviewer); err != nil {
		t.Fatalf("CreateInterviewer failed: %v", err)
	}

	list, err := repo.ListInterviewers(ctx, iv.ID)
	if err != nil {
		t.Fatalf("ListInterviewers failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 interviewer, got %d", len(list))
	}
}

func TestRepository_DeleteInterviewer(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	iv := newTestInterviewForScorecard(t, repo, ctx, "ivrepo2@test.com")
	interviewer := &Interviewer{InterviewID: iv.ID, EmployeeID: uuid.New()}
	repo.CreateInterviewer(ctx, interviewer)

	if err := repo.DeleteInterviewer(ctx, interviewer.ID); err != nil {
		t.Fatalf("DeleteInterviewer failed: %v", err)
	}
	list, _ := repo.ListInterviewers(ctx, iv.ID)
	if len(list) != 0 {
		t.Errorf("expected 0 interviewers after delete, got %d", len(list))
	}
}

func TestRepository_CreateAndFindScorecardItem(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	iv := newTestInterviewForScorecard(t, repo, ctx, "ivrepo3@test.com")
	item := &InterviewScorecardItem{InterviewID: iv.ID, Criterion: "Technical Skill", Weight: 30}
	if err := repo.CreateScorecardItem(ctx, item); err != nil {
		t.Fatalf("CreateScorecardItem failed: %v", err)
	}

	found, err := repo.FindScorecardItemByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("FindScorecardItemByID failed: %v", err)
	}
	if found.Criterion != "Technical Skill" {
		t.Errorf("expected 'Technical Skill', got %s", found.Criterion)
	}
}

func TestRepository_ListScorecardItems(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	iv := newTestInterviewForScorecard(t, repo, ctx, "ivrepo4@test.com")
	repo.CreateScorecardItem(ctx, &InterviewScorecardItem{InterviewID: iv.ID, Criterion: "Technical Skill", Weight: 30})
	repo.CreateScorecardItem(ctx, &InterviewScorecardItem{InterviewID: iv.ID, Criterion: "Communication", Weight: 20})

	list, err := repo.ListScorecardItems(ctx, iv.ID)
	if err != nil {
		t.Fatalf("ListScorecardItems failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteScorecardItem(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	iv := newTestInterviewForScorecard(t, repo, ctx, "ivrepo5@test.com")
	item := &InterviewScorecardItem{InterviewID: iv.ID, Criterion: "Original", Weight: 10}
	repo.CreateScorecardItem(ctx, item)

	item.Criterion = "Updated"
	if err := repo.UpdateScorecardItem(ctx, item); err != nil {
		t.Fatalf("UpdateScorecardItem failed: %v", err)
	}
	found, _ := repo.FindScorecardItemByID(ctx, item.ID)
	if found.Criterion != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.Criterion)
	}

	if err := repo.DeleteScorecardItem(ctx, item.ID); err != nil {
		t.Fatalf("DeleteScorecardItem failed: %v", err)
	}
	if _, err := repo.FindScorecardItemByID(ctx, item.ID); err == nil {
		t.Error("expected error finding deleted item, got nil")
	}
}

func TestRepository_CreateAndFindRequisitionRequirement(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)

	req := &JobRequisitionRequirement{RequisitionID: r.ID, RequirementType: "EXPERIENCE_YEARS", Name: "Min Experience"}
	if err := repo.CreateRequisitionRequirement(ctx, req); err != nil {
		t.Fatalf("CreateRequisitionRequirement failed: %v", err)
	}

	found, err := repo.FindRequisitionRequirementByID(ctx, req.ID)
	if err != nil {
		t.Fatalf("FindRequisitionRequirementByID failed: %v", err)
	}
	if found.Name != "Min Experience" {
		t.Errorf("expected 'Min Experience', got %s", found.Name)
	}
}

func TestRepository_ListRequisitionRequirements(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	repo.CreateRequisitionRequirement(ctx, &JobRequisitionRequirement{RequisitionID: r.ID, RequirementType: "EDUCATION", Name: "S1"})
	repo.CreateRequisitionRequirement(ctx, &JobRequisitionRequirement{RequisitionID: r.ID, RequirementType: "LANGUAGE", Name: "English"})

	list, err := repo.ListRequisitionRequirements(ctx, r.ID)
	if err != nil {
		t.Fatalf("ListRequisitionRequirements failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteRequisitionRequirement(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	req := &JobRequisitionRequirement{RequisitionID: r.ID, RequirementType: "OTHER", Name: "Original"}
	repo.CreateRequisitionRequirement(ctx, req)

	req.Name = "Updated"
	if err := repo.UpdateRequisitionRequirement(ctx, req); err != nil {
		t.Fatalf("UpdateRequisitionRequirement failed: %v", err)
	}
	found, _ := repo.FindRequisitionRequirementByID(ctx, req.ID)
	if found.Name != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.Name)
	}

	if err := repo.DeleteRequisitionRequirement(ctx, req.ID); err != nil {
		t.Fatalf("DeleteRequisitionRequirement failed: %v", err)
	}
	if _, err := repo.FindRequisitionRequirementByID(ctx, req.ID); err == nil {
		t.Error("expected error finding deleted requirement, got nil")
	}
}

func TestRepository_CreateAndFindRequisitionCompetency(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	comp := &competency.Competency{Name: "Go Programming"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	rc := &JobRequisitionCompetency{RequisitionID: r.ID, CompetencyID: comp.ID}
	if err := repo.CreateRequisitionCompetency(ctx, rc); err != nil {
		t.Fatalf("CreateRequisitionCompetency failed: %v", err)
	}

	found, err := repo.FindRequisitionCompetencyByID(ctx, rc.ID)
	if err != nil {
		t.Fatalf("FindRequisitionCompetencyByID failed: %v", err)
	}
	if found.CompetencyID != comp.ID {
		t.Errorf("expected competency_id %s, got %s", comp.ID, found.CompetencyID)
	}
}

func TestRepository_ListRequisitionCompetencies(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	comp1 := &competency.Competency{Name: "Go"}
	comp2 := &competency.Competency{Name: "SQL"}
	db.Create(comp1)
	db.Create(comp2)
	repo.CreateRequisitionCompetency(ctx, &JobRequisitionCompetency{RequisitionID: r.ID, CompetencyID: comp1.ID})
	repo.CreateRequisitionCompetency(ctx, &JobRequisitionCompetency{RequisitionID: r.ID, CompetencyID: comp2.ID})

	list, err := repo.ListRequisitionCompetencies(ctx, r.ID)
	if err != nil {
		t.Fatalf("ListRequisitionCompetencies failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteRequisitionCompetency(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	r := &JobRequisition{OrganizationID: uuid.New(), Title: "Req", Status: ReqStatusOpen}
	repo.CreateRequisition(ctx, r)
	comp := &competency.Competency{Name: "Go"}
	db.Create(comp)
	rc := &JobRequisitionCompetency{RequisitionID: r.ID, CompetencyID: comp.ID}
	repo.CreateRequisitionCompetency(ctx, rc)

	level := 4
	rc.RequiredLevel = &level
	if err := repo.UpdateRequisitionCompetency(ctx, rc); err != nil {
		t.Fatalf("UpdateRequisitionCompetency failed: %v", err)
	}
	found, _ := repo.FindRequisitionCompetencyByID(ctx, rc.ID)
	if found.RequiredLevel == nil || *found.RequiredLevel != 4 {
		t.Errorf("expected required_level 4, got %v", found.RequiredLevel)
	}

	if err := repo.DeleteRequisitionCompetency(ctx, rc.ID); err != nil {
		t.Fatalf("DeleteRequisitionCompetency failed: %v", err)
	}
	if _, err := repo.FindRequisitionCompetencyByID(ctx, rc.ID); err == nil {
		t.Error("expected error finding deleted competency, got nil")
	}
}
