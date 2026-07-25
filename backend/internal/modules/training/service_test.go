package training

import (
	"testing"
)

func TestService_CreateCategory_Validation(t *testing.T) {
	svc := testSvc(t)

	// Test with minimal required fields
	resp, err := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "FIN",
		Name: "Finance",
	})
	if err != nil {
		t.Fatalf("CreateCategory with valid input failed: %v", err)
	}
	if resp.Name != "Finance" {
		t.Errorf("expected name=Finance, got %s", resp.Name)
	}
}

func TestService_CourseRequiresCategory(t *testing.T) {
	svc := testSvc(t)

	// Attempt to create course with non-existent category
	_, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID: "00000000-0000-0000-0000-000000000000",
		Code:       "INVALID",
		Name:       "Invalid Course",
	})
	if err == nil {
		t.Fatal("expected error for non-existent category")
	}
}

func TestService_SessionStatusFlow(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)

	// SCHEDULED -> IN_PROGRESS
	resp, err := svc.UpdateSessionStatus(testCtx(), sessID, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("status transition failed: %v", err)
	}
	if resp.Status != "IN_PROGRESS" {
		t.Errorf("expected IN_PROGRESS, got %s", resp.Status)
	}

	// IN_PROGRESS -> COMPLETED
	resp, err = svc.UpdateSessionStatus(testCtx(), sessID, "COMPLETED")
	if err != nil {
		t.Fatalf("status transition failed: %v", err)
	}
	if resp.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", resp.Status)
	}
}

func TestService_ParticipantCannotRegisterCancelledSession(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)

	// Cancel session first
	_, err := svc.UpdateSessionStatus(testCtx(), sessID, "CANCELLED")
	if err != nil {
		t.Fatalf("cancel session failed: %v", err)
	}

	// Try to register participant
	_, err = svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID:  sessID,
		EmployeeID: "00000000-0000-0000-0000-000000000001",
	})
	if err == nil {
		t.Fatal("expected error registering to cancelled session")
	}
}

func TestService_UpdateParticipantScore(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000002"

	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	// Update score
	score := 85.5
	updated, err := svc.UpdateParticipant(testCtx(), part.ID, UpdateTrainingParticipantRequest{
		Score: &score,
	})
	if err != nil {
		t.Fatalf("UpdateParticipant score failed: %v", err)
	}
	if updated.Score != score {
		t.Errorf("expected score=%.1f, got %.1f", score, updated.Score)
	}
	if updated.CompletedAt == "" {
		t.Error("expected CompletedAt to be set after score entry")
	}
}

func TestService_CreateAndListEvaluations(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000003"

	// Register participant
	svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})

	// Submit evaluation
	fb := "Great training!"
	_, err := svc.CreateEvaluation(testCtx(), CreateTrainingEvaluationRequest{
		SessionID: sessID, EmployeeID: empID, Rating: 4, Feedback: &fb,
	})
	if err != nil {
		t.Fatalf("CreateEvaluation failed: %v", err)
	}

	// List evaluations for session
	sessIDStr := sessID
	resp, err := svc.ListEvaluations(testCtx(), &sessIDStr, nil, 1, 10)
	if err != nil {
		t.Fatalf("ListEvaluations failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 evaluation, got %d", resp.Total)
	}
}

func TestService_CertificateLifecycle(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000004"

	part, _ := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})

	// Issue certificate
	cert, err := svc.CreateCertificate(testCtx(), CreateTrainingCertificateRequest{
		ParticipantID: part.ID,
		CertificateNo: "CERT-TEST-001",
		IssuedDate:    "2026-08-05",
	})
	if err != nil {
		t.Fatalf("CreateCertificate failed: %v", err)
	}

	// Update expiry
	expiry := "2027-08-05"
	updated, err := svc.UpdateCertificate(testCtx(), cert.ID, UpdateTrainingCertificateRequest{
		ExpiryDate: &expiry,
	})
	if err != nil {
		t.Fatalf("UpdateCertificate failed: %v", err)
	}
	if updated.ExpiryDate != expiry {
		t.Errorf("expected expiry=%s, got %s", expiry, updated.ExpiryDate)
	}
}

func TestService_MaterialManagement(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)

	// Create materials
	url1 := "https://example.com/slide1.pdf"
	url2 := "https://example.com/slide2.pdf"
	sort1 := 1
	m1, err := svc.CreateMaterial(testCtx(), CreateTrainingMaterialRequest{
		SessionID: sessID,
		Title:     "Slide 1",
		FileURL:   &url1,
		SortOrder: &sort1,
	})
	if err != nil {
		t.Fatalf("CreateMaterial failed: %v", err)
	}
	svc.CreateMaterial(testCtx(), CreateTrainingMaterialRequest{
		SessionID: sessID,
		Title:     "Slide 2",
		FileURL:   &url2,
	})

	// List materials
	mats, err := svc.ListMaterials(testCtx(), sessID)
	if err != nil {
		t.Fatalf("ListMaterials failed: %v", err)
	}
	if len(mats) != 2 {
		t.Errorf("expected 2 materials, got %d", len(mats))
	}

	// Delete material
	if err := svc.DeleteMaterial(testCtx(), m1.ID); err != nil {
		t.Fatalf("DeleteMaterial failed: %v", err)
	}
	mats, _ = svc.ListMaterials(testCtx(), sessID)
	if len(mats) != 1 {
		t.Errorf("expected 1 material after delete, got %d", len(mats))
	}
}

func TestService_PaginationDefaults(t *testing.T) {
	svc := testSvc(t)
	desc := "Page test"
	// Create 5 categories
	for i := 0; i < 5; i++ {
		code := string(rune('A' + i))
		svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
			Code: string(rune('A' + i)), Name: "Category " + string(rune('A' + i)), Description: &desc,
		})
		_ = code
	}

	// First page with per_page=2
	resp, err := svc.ListCategories(testCtx(), 1, 2)
	if err != nil {
		t.Fatalf("ListCategories failed: %v", err)
	}
	if resp.Total != 5 {
		t.Errorf("expected total=5, got %d", resp.Total)
	}
	if resp.TotalPages != 3 {
		t.Errorf("expected total_pages=3, got %d", resp.TotalPages)
	}
}
