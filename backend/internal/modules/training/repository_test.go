package training

import (
	"testing"

	"github.com/google/uuid"
)

func TestRepository_GetEmployeeTrainingSummary(t *testing.T) {
	repo := testRepo(t)
	svc := NewService(repo, testLogger())
	ctx := testCtx()

	categoryID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, categoryID)
	sessionID := seedSession(t, svc, courseID)

	db, err := repo.db(ctx)
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}

	competencyID := uuid.New()
	targetLevel := 4
	if err := db.Create(&TrainingCourseCompetency{
		CourseID:     uuid.MustParse(courseID),
		CompetencyID: competencyID,
		TargetLevel:  &targetLevel,
	}).Error; err != nil {
		t.Fatalf("failed to seed course competency: %v", err)
	}

	employeeID := uuid.New()
	participant := &TrainingParticipant{
		SessionID:        uuid.MustParse(sessionID),
		EmployeeID:       employeeID,
		CompletionStatus: CompletionCompleted,
		Score:            92,
	}
	if err := db.Create(participant).Error; err != nil {
		t.Fatalf("failed to seed participant: %v", err)
	}
	if err := db.Create(&TrainingCertificate{
		ParticipantID: participant.ID,
		CertificateNo: "CERT-001",
		IssuedDate:    "2026-01-15",
	}).Error; err != nil {
		t.Fatalf("failed to seed certificate: %v", err)
	}

	summary, err := repo.GetEmployeeTrainingSummary(ctx, employeeID)
	if err != nil {
		t.Fatalf("GetEmployeeTrainingSummary failed: %v", err)
	}
	if summary.TotalTraining != 1 {
		t.Errorf("expected total_training 1, got %d", summary.TotalTraining)
	}
	if summary.Completed != 1 {
		t.Errorf("expected completed 1, got %d", summary.Completed)
	}
	if summary.Failed != 0 {
		t.Errorf("expected failed 0, got %d", summary.Failed)
	}
	if summary.TrainingHours != 8 {
		t.Errorf("expected training_hours 8, got %.2f", summary.TrainingHours)
	}
	if summary.AverageScore != 92 {
		t.Errorf("expected average_score 92, got %.2f", summary.AverageScore)
	}
	if summary.CertificationCount != 1 {
		t.Errorf("expected certification_count 1, got %d", summary.CertificationCount)
	}
	if summary.CompetencyTrainingCount != 1 {
		t.Errorf("expected competency_training_count 1, got %d", summary.CompetencyTrainingCount)
	}
}

func TestRepository_CreateCategory(t *testing.T) {
	svc := testSvc(t)
	desc := "IT training category"
	resp, err := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code:        "IT",
		Name:        "Information Technology",
		Description: &desc,
	})
	if err != nil {
		t.Fatalf("CreateCategory failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if resp.Code != "IT" {
		t.Errorf("expected code=IT, got %s", resp.Code)
	}
	if resp.Name != "Information Technology" {
		t.Errorf("expected name=Information Technology, got %s", resp.Name)
	}
	if !resp.IsActive {
		t.Error("expected IsActive=true by default")
	}
}

func TestRepository_ListCategories(t *testing.T) {
	svc := testSvc(t)
	desc := "Test"
	svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{Code: "A", Name: "Alpha", Description: &desc})
	svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{Code: "B", Name: "Beta", Description: &desc})

	resp, err := svc.ListCategories(testCtx(), 1, 10)
	if err != nil {
		t.Fatalf("ListCategories failed: %v", err)
	}
	if err := validatePagination(resp, 2); err != nil {
		t.Error(err)
	}
}

func TestRepository_GetCategoryByID_NotFound(t *testing.T) {
	repo := testRepo(t)
	_, err := repo.FindCategoryByID(testCtx(), uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent category")
	}
}

func TestRepository_UpdateCategory(t *testing.T) {
	svc := testSvc(t)
	desc := "Original"
	created, _ := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "OLD", Name: "Old Name", Description: &desc,
	})
	newName := "Updated Name"
	updated, err := svc.UpdateCategory(testCtx(), created.ID, UpdateTrainingCategoryRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateCategory failed: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name=Updated Name, got %s", updated.Name)
	}
}

func TestRepository_DeleteCategory(t *testing.T) {
	svc := testSvc(t)
	desc := "To delete"
	created, _ := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "DEL", Name: "Delete Me", Description: &desc,
	})
	if err := svc.DeleteCategory(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteCategory failed: %v", err)
	}
	// Verify deleted
	_, err := svc.GetCategoryByID(testCtx(), created.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestRepository_CreateCourse(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	dur := 16.0
	resp, err := svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{
		CategoryID:   catID,
		Code:         "PYTHON-101",
		Name:         "Python Basics",
		DurationHour: &dur,
	})
	if err != nil {
		t.Fatalf("CreateCourse failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if resp.CategoryID != catID {
		t.Errorf("expected category_id=%s, got %s", catID, resp.CategoryID)
	}
}

func TestRepository_ListCourses_ByCategory(t *testing.T) {
	svc := testSvc(t)
	catID1 := seedCategory(t, svc)
	// Create second category
	desc := "Business"
	cat2, _ := svc.CreateCategory(testCtx(), CreateTrainingCategoryRequest{
		Code: "BUS", Name: "Business", Description: &desc,
	})
	svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{CategoryID: catID1, Code: "GOLANG-101", Name: "Golang"})
	svc.CreateCourse(testCtx(), CreateTrainingCourseRequest{CategoryID: cat2.ID, Code: "MGMT-101", Name: "Management"})

	// Filter by catID1
	resp, err := svc.ListCourses(testCtx(), &catID1, 1, 10)
	if err != nil {
		t.Fatalf("ListCourses failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 course in category, got %d", resp.Total)
	}
}

func TestRepository_CreateSession(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	resp, err := svc.CreateSession(testCtx(), CreateTrainingSessionRequest{
		CourseID:    courseID,
		SessionCode: "CLS-001",
		TrainerName: "Jane Smith",
		StartDate:   "2026-08-01",
		EndDate:     "2026-08-03",
		MaxQuota:    20,
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if resp.Status != "SCHEDULED" {
		t.Errorf("expected status=SCHEDULED, got %s", resp.Status)
	}
}

func TestRepository_UpdateSessionStatus(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)

	resp, err := svc.UpdateSessionStatus(testCtx(), sessID, "IN_PROGRESS")
	if err != nil {
		t.Fatalf("UpdateSessionStatus failed: %v", err)
	}
	if resp.Status != "IN_PROGRESS" {
		t.Errorf("expected status=IN_PROGRESS, got %s", resp.Status)
	}
}

func TestRepository_CreateParticipant(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := uuid.New().String()

	resp, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID:  sessID,
		EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if resp.AttendanceStatus != "PRESENT" {
		t.Errorf("expected status=PRESENT, got %s", resp.AttendanceStatus)
	}
}

func TestRepository_QuotaFull(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	// Max quota = 1
	sess, err := svc.CreateSession(testCtx(), CreateTrainingSessionRequest{
		CourseID:    courseID,
		SessionCode: "CLS-QUOTA",
		TrainerName: "Quota Test",
		StartDate:   "2026-08-01",
		EndDate:     "2026-08-03",
		MaxQuota:    1,
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	// Register first participant (should succeed)
	_, err = svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID:  sess.ID,
		EmployeeID: uuid.New().String(),
	})
	if err != nil {
		t.Fatalf("first participant should succeed: %v", err)
	}
	// Register second participant (should fail - quota full)
	_, err = svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID:  sess.ID,
		EmployeeID: uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected error for quota full, got nil")
	}
}

func TestRepository_CreateEvaluation(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := uuid.New().String()

	// Register first
	svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})

	resp, err := svc.CreateEvaluation(testCtx(), CreateTrainingEvaluationRequest{
		SessionID:  sessID,
		EmployeeID: empID,
		Rating:     5,
	})
	if err != nil {
		t.Fatalf("CreateEvaluation failed: %v", err)
	}
	if resp.Rating != 5 {
		t.Errorf("expected rating=5, got %d", resp.Rating)
	}
}

func TestRepository_CreateCertificate(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessID := seedSession(t, svc, courseID)
	empID := uuid.New().String()

	// Register participant
	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	resp, err := svc.CreateCertificate(testCtx(), CreateTrainingCertificateRequest{
		ParticipantID: part.ID,
		CertificateNo: "CERT-001",
		IssuedDate:    "2026-08-05",
	})
	if err != nil {
		t.Fatalf("CreateCertificate failed: %v", err)
	}
	if resp.CertificateNo != "CERT-001" {
		t.Errorf("expected cert_no=CERT-001, got %s", resp.CertificateNo)
	}
}

func TestRepository_ListSessions_ByStatus(t *testing.T) {
	svc := testSvc(t)
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	seedSession(t, svc, courseID)

	status := "SCHEDULED"
	resp, err := svc.ListSessions(testCtx(), nil, &status, 1, 10)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if resp.Total < 1 {
		t.Errorf("expected at least 1 SCHEDULED session, got %d", resp.Total)
	}
}
