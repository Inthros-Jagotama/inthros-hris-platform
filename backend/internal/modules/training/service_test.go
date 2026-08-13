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

func TestService_CourseAutoCode(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc) // kategori dengan kode "TECH"

	// Tanpa code → di-generate otomatis {KODE_KATEGORI}-{sekuens}
	resp1, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID: catID,
		Name:       "Course A",
	})
	if err != nil {
		t.Fatalf("failed to create course without code: %v", err)
	}
	if resp1.Code != "TECH-001" {
		t.Fatalf("expected auto code TECH-001, got %s", resp1.Code)
	}

	// Course kedua → sekuens bertambah
	resp2, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID: catID,
		Name:       "Course B",
	})
	if err != nil {
		t.Fatalf("failed to create second course without code: %v", err)
	}
	if resp2.Code != "TECH-002" {
		t.Fatalf("expected auto code TECH-002, got %s", resp2.Code)
	}

	// Code eksplisit tetap dihormati
	resp3, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID: catID,
		Code:       "CUSTOM-01",
		Name:       "Course C",
	})
	if err != nil {
		t.Fatalf("failed to create course with explicit code: %v", err)
	}
	if resp3.Code != "CUSTOM-01" {
		t.Fatalf("expected explicit code CUSTOM-01, got %s", resp3.Code)
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

// =========================================================================
// P0-BE: Providers & Trainers (plan §11–§13)
// =========================================================================

func seedProvider(t *testing.T, svc *Service) string {
	t.Helper()
	typ := "EXTERNAL"
	resp, err := svc.CreateProvider(testCtx(), CreateTrainingProviderRequest{
		Code: "ABC-INST",
		Name: "ABC Training Institute",
		Type: &typ,
	})
	if err != nil {
		t.Fatalf("failed to seed provider: %v", err)
	}
	return resp.ID
}

func seedExternalTrainer(t *testing.T, svc *Service, providerID string) string {
	t.Helper()
	resp, err := svc.CreateTrainer(testCtx(), CreateTrainingTrainerRequest{
		Type:       "EXTERNAL",
		ProviderID: &providerID,
		Name:       "External Trainer A",
	})
	if err != nil {
		t.Fatalf("failed to seed trainer: %v", err)
	}
	return resp.ID
}

func TestService_ProviderLifecycle(t *testing.T) {
	svc := testSvc(t)
	id := seedProvider(t, svc)

	got, err := svc.GetProviderByID(testCtx(), id)
	if err != nil {
		t.Fatalf("GetProviderByID failed: %v", err)
	}
	if got.Type != "EXTERNAL" {
		t.Errorf("expected type EXTERNAL, got %s", got.Type)
	}

	resp, err := svc.ListProviders(testCtx(), 1, 10)
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 provider, got %d", resp.Total)
	}

	name := "ABC Training Institute (Updated)"
	updated, err := svc.UpdateProvider(testCtx(), id, UpdateTrainingProviderRequest{Name: &name})
	if err != nil {
		t.Fatalf("UpdateProvider failed: %v", err)
	}
	if updated.Name != name {
		t.Errorf("expected name %s, got %s", name, updated.Name)
	}
}

func TestService_TrainerRequiresReferenceByType(t *testing.T) {
	svc := testSvc(t)

	// INTERNAL tanpa employee_id → error.
	_, err := svc.CreateTrainer(testCtx(), CreateTrainingTrainerRequest{
		Type: "INTERNAL",
		Name: "Internal Trainer",
	})
	if err == nil {
		t.Fatal("expected error for INTERNAL trainer without employee_id")
	}

	// EXTERNAL tanpa provider_id → error.
	_, err = svc.CreateTrainer(testCtx(), CreateTrainingTrainerRequest{
		Type: "EXTERNAL",
		Name: "External Trainer",
	})
	if err == nil {
		t.Fatal("expected error for EXTERNAL trainer without provider_id")
	}

	// EXTERNAL dengan provider_id → sukses.
	providerID := seedProvider(t, svc)
	trainer, err := svc.CreateTrainer(testCtx(), CreateTrainingTrainerRequest{
		Type:       "EXTERNAL",
		ProviderID: &providerID,
		Name:       "External Trainer OK",
	})
	if err != nil {
		t.Fatalf("CreateTrainer EXTERNAL with provider failed: %v", err)
	}
	if trainer.Type != "EXTERNAL" {
		t.Errorf("expected EXTERNAL, got %s", trainer.Type)
	}
}

func TestService_AddAndListSessionTrainers(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	providerID := seedProvider(t, svc)
	trainerID := seedExternalTrainer(t, svc, providerID)

	added, err := svc.AddSessionTrainer(testCtx(), sessID, AddSessionTrainerRequest{TrainerID: trainerID})
	if err != nil {
		t.Fatalf("AddSessionTrainer failed: %v", err)
	}
	if added.Role != "MAIN" {
		t.Errorf("expected default role MAIN, got %s", added.Role)
	}

	items, err := svc.ListSessionTrainers(testCtx(), sessID)
	if err != nil {
		t.Fatalf("ListSessionTrainers failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 session trainer, got %d", len(items))
	}

	// Duplikasi trainer pada session yang sama → error.
	_, err = svc.AddSessionTrainer(testCtx(), sessID, AddSessionTrainerRequest{TrainerID: trainerID})
	if err == nil {
		t.Fatal("expected error when assigning the same trainer twice")
	}
}

// =========================================================================
// P0-BE: Enrollment & quota (plan §18/§32)
// =========================================================================

func TestService_DuplicateParticipantRejected(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000020"

	if _, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	}); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}
	_, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})
	if err == nil {
		t.Fatal("expected error for duplicate participant")
	}
}

func TestService_ParticipantRegistrationStatus(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000021"

	nom := "NOMINATED"
	p, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID, RegistrationStatus: &nom,
	})
	if err != nil {
		t.Fatalf("CreateParticipant NOMINATED failed: %v", err)
	}
	if p.RegistrationStatus != "NOMINATED" {
		t.Errorf("expected NOMINATED, got %s", p.RegistrationStatus)
	}
	if p.RegisteredAt != "" {
		t.Error("NOMINATED should not set registered_at")
	}
}

// =========================================================================
// P0-BE: Session enhancement (plan §14)
// =========================================================================

func TestService_SessionExternalRequiresProvider(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	provType := "EXTERNAL"
	_, err := svc.CreateSession(testCtx(), CreateTrainingSessionRequest{
		CourseID:     courseID,
		SessionCode:  "CLS-EXT-1",
		TrainerName:  "External Trainer",
		StartDate:    "2026-09-01",
		EndDate:      "2026-09-05",
		ProviderType: &provType,
	})
	if err == nil {
		t.Fatal("expected error for EXTERNAL session without provider_id")
	}

	providerID := seedProvider(t, svc)
	sess, err := svc.CreateSession(testCtx(), CreateTrainingSessionRequest{
		CourseID:     courseID,
		SessionCode:  "CLS-EXT-2",
		TrainerName:  "External Trainer",
		StartDate:    "2026-09-01",
		EndDate:      "2026-09-05",
		ProviderType: &provType,
		ProviderID:   &providerID,
	})
	if err != nil {
		t.Fatalf("CreateSession EXTERNAL with provider failed: %v", err)
	}
	if sess.ProviderType != "EXTERNAL" || sess.ProviderID != providerID {
		t.Errorf("expected EXTERNAL + provider_id, got %s/%s", sess.ProviderType, sess.ProviderID)
	}
}

// =========================================================================
// P0-BE: Attendance (plan §19)
// =========================================================================

func TestService_AttendanceMarkAndList(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000022"

	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	late := "LATE"
	ci := "2026-08-01T08:15:00"
	att, err := svc.MarkAttendance(testCtx(), sessID, []MarkTrainingAttendanceRequest{
		{ParticipantID: part.ID, AttendanceDate: "2026-08-01", CheckIn: &ci, Status: &late},
	})
	if err != nil {
		t.Fatalf("MarkAttendance failed: %v", err)
	}
	if len(att) != 1 || att[0].Status != "LATE" {
		t.Errorf("expected 1 LATE attendance, got %+v", att)
	}

	// Upsert: mark ulang tanggal yang sama tidak membuat baris baru.
	present := "PRESENT"
	att2, err := svc.MarkAttendance(testCtx(), sessID, []MarkTrainingAttendanceRequest{
		{ParticipantID: part.ID, AttendanceDate: "2026-08-01", Status: &present},
	})
	if err != nil {
		t.Fatalf("MarkAttendance upsert failed: %v", err)
	}
	if len(att2) != 1 || att2[0].Status != "PRESENT" || att2[0].ID != att[0].ID {
		t.Errorf("expected upsert same row, got %+v", att2)
	}

	rows, err := svc.ListAttendanceBySession(testCtx(), sessID)
	if err != nil {
		t.Fatalf("ListAttendanceBySession failed: %v", err)
	}
	if len(rows) != 1 {
		t.Errorf("expected 1 attendance row, got %d", len(rows))
	}
	if rows[0].EmployeeID != empID {
		t.Errorf("expected employee %s, got %s", empID, rows[0].EmployeeID)
	}
}

// =========================================================================
// P0-BE: Assessments (plan §21)
// =========================================================================

func TestService_AssessmentLifecycle(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000023"

	part, _ := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})

	// passing > max → error.
	max := 100.0
	passing := 120.0
	_, err := svc.CreateAssessment(testCtx(), CreateTrainingAssessmentRequest{
		SessionID: sessID, Name: "Post Test", MaxScore: &max, PassingScore: &passing,
	})
	if err == nil {
		t.Fatal("expected error when passing_score > max_score")
	}

	// Valid assessment.
	passing = 70.0
	a, err := svc.CreateAssessment(testCtx(), CreateTrainingAssessmentRequest{
		SessionID: sessID, Name: "Post Test", MaxScore: &max, PassingScore: &passing,
	})
	if err != nil {
		t.Fatalf("CreateAssessment failed: %v", err)
	}
	if a.AttemptLimit != 1 {
		t.Errorf("expected default attempt_limit 1, got %d", a.AttemptLimit)
	}

	items, err := svc.ListAssessmentsBySession(testCtx(), sessID)
	if err != nil {
		t.Fatalf("ListAssessmentsBySession failed: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 assessment, got %d", len(items))
	}

	// Score > max → error.
	_, err = svc.SubmitAssessmentResult(testCtx(), a.ID, SubmitAssessmentResultRequest{
		ParticipantID: part.ID, Score: 150,
	})
	if err == nil {
		t.Fatal("expected error when score exceeds max_score")
	}

	// Passed (score >= passing).
	res, err := svc.SubmitAssessmentResult(testCtx(), a.ID, SubmitAssessmentResultRequest{
		ParticipantID: part.ID, Score: 85,
	})
	if err != nil {
		t.Fatalf("SubmitAssessmentResult failed: %v", err)
	}
	if !res.Passed || res.Attempt != 1 {
		t.Errorf("expected passed=true attempt=1, got %+v", res)
	}

	// Attempt limit 1 → attempt kedua ditolak.
	_, err = svc.SubmitAssessmentResult(testCtx(), a.ID, SubmitAssessmentResultRequest{
		ParticipantID: part.ID, Score: 90,
	})
	if err == nil {
		t.Fatal("expected error when attempt limit reached")
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
