package training

import (
	"testing"
)

// =========================================================================
// Evaluation Form & Questions (P2-BE — plan §22)
// =========================================================================

func TestService_CreateEvaluationForm_And_Questions(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	form, err := svc.CreateEvaluationForm(testCtx(), CreateEvaluationFormRequest{
		SessionID: sessionID,
		Name:      "Post-Training Evaluation",
	})
	if err != nil {
		t.Fatalf("CreateEvaluationForm failed: %v", err)
	}
	if form.SessionID != sessionID {
		t.Errorf("expected session id %s, got %s", sessionID, form.SessionID)
	}
	if !form.IsActive {
		t.Error("expected form is_active default true")
	}

	q, err := svc.CreateEvaluationQuestion(testCtx(), form.ID, CreateEvaluationQuestionRequest{
		Question:     "How would you rate the trainer?",
		QuestionType: "RATING",
	})
	if err != nil {
		t.Fatalf("CreateEvaluationQuestion failed: %v", err)
	}
	if q.QuestionType != "RATING" {
		t.Errorf("expected type RATING, got %s", q.QuestionType)
	}

	// Form + questions via session.
	withQ, err := svc.GetEvaluationFormBySession(testCtx(), sessionID)
	if err != nil {
		t.Fatalf("GetEvaluationFormBySession failed: %v", err)
	}
	if withQ.Form.ID != form.ID {
		t.Errorf("expected form id %s, got %s", form.ID, withQ.Form.ID)
	}
	if len(withQ.Questions) != 1 {
		t.Errorf("expected 1 question, got %d", len(withQ.Questions))
	}
}

func TestService_SubmitEvaluationAnswers_Upserts(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	form, err := svc.CreateEvaluationForm(testCtx(), CreateEvaluationFormRequest{SessionID: sessionID, Name: "Eval"})
	if err != nil {
		t.Fatalf("CreateEvaluationForm failed: %v", err)
	}
	q1, err := svc.CreateEvaluationQuestion(testCtx(), form.ID, CreateEvaluationQuestionRequest{Question: "Rating?", QuestionType: "RATING"})
	if err != nil {
		t.Fatalf("CreateEvaluationQuestion failed: %v", err)
	}
	q2, err := svc.CreateEvaluationQuestion(testCtx(), form.ID, CreateEvaluationQuestionRequest{Question: "Comment?", QuestionType: "TEXT"})
	if err != nil {
		t.Fatalf("CreateEvaluationQuestion failed: %v", err)
	}

	// Participant seed.
	empID := "00000000-0000-0000-0000-000000000002"
	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessionID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}
	participantID := part.ID

	answers, err := svc.SubmitEvaluationAnswers(testCtx(), form.ID, participantID, SubmitEvaluationAnswersRequest{
		Answers: []EvaluationAnswerInput{
			{QuestionID: q1.ID, Answer: "5"},
			{QuestionID: q2.ID, Answer: "Very good"},
		},
	})
	if err != nil {
		t.Fatalf("SubmitEvaluationAnswers failed: %v", err)
	}
	if len(answers) != 2 {
		t.Fatalf("expected 2 answers, got %d", len(answers))
	}

	// Re-submit (upsert) — jawaban berubah, tetap 2 baris.
	answers2, err := svc.SubmitEvaluationAnswers(testCtx(), form.ID, participantID, SubmitEvaluationAnswersRequest{
		Answers: []EvaluationAnswerInput{
			{QuestionID: q1.ID, Answer: "4"},
		},
	})
	if err != nil {
		t.Fatalf("re-submit failed: %v", err)
	}
	if len(answers2) != 1 {
		t.Fatalf("expected 1 answer on resubmit, got %d", len(answers2))
	}
	if answers2[0].Answer != "4" {
		t.Errorf("expected updated answer 4, got %s", answers2[0].Answer)
	}

	list, err := svc.ListEvaluationAnswers(testCtx(), nil, strP(participantID))
	if err != nil {
		t.Fatalf("ListEvaluationAnswers failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 answers in list, got %d", len(list))
	}
}

func TestService_SubmitEvaluationAnswers_WrongQuestion_Rejected(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	formA, err := svc.CreateEvaluationForm(testCtx(), CreateEvaluationFormRequest{SessionID: sessionID, Name: "A"})
	if err != nil {
		t.Fatalf("CreateEvaluationForm failed: %v", err)
	}
	formB, err := svc.CreateEvaluationForm(testCtx(), CreateEvaluationFormRequest{SessionID: sessionID, Name: "B"})
	if err != nil {
		t.Fatalf("CreateEvaluationForm failed: %v", err)
	}
	qB, err := svc.CreateEvaluationQuestion(testCtx(), formB.ID, CreateEvaluationQuestionRequest{Question: "Q", QuestionType: "TEXT"})
	if err != nil {
		t.Fatalf("CreateEvaluationQuestion failed: %v", err)
	}

	empID := "00000000-0000-0000-0000-000000000003"
	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessionID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	// Pertanyaan milik form B, dikirim ke form A → ditolak.
	if _, err := svc.SubmitEvaluationAnswers(testCtx(), formA.ID, part.ID, SubmitEvaluationAnswersRequest{
		Answers: []EvaluationAnswerInput{{QuestionID: qB.ID, Answer: "x"}},
	}); err == nil {
		t.Error("expected error for question from different form, got nil")
	}
}

// =========================================================================
// Effectiveness Assessments (P2-BE — plan §23)
// =========================================================================

func TestService_EffectivenessAssessment_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	empID := "00000000-0000-0000-0000-000000000004"
	part, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessionID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	before := 40.0
	after := 75.0
	effect := 35.0
	created, err := svc.CreateEffectivenessAssessment(testCtx(), CreateEffectivenessAssessmentRequest{
		ParticipantID:      part.ID,
		AssessmentDate:     "2026-10-01",
		BeforeScore:        &before,
		AfterScore:         &after,
		EffectivenessScore: &effect,
		Remarks:            strP("improved significantly"),
	})
	if err != nil {
		t.Fatalf("CreateEffectivenessAssessment failed: %v", err)
	}
	if created.BeforeScore == nil || *created.BeforeScore != before {
		t.Errorf("expected before_score %v, got %v", before, created.BeforeScore)
	}

	list, err := svc.ListEffectivenessAssessments(testCtx(), strP(part.ID), 1, 20)
	if err != nil {
		t.Fatalf("ListEffectivenessAssessments failed: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("expected 1 assessment, got %d", list.Total)
	}

	newAfter := 80.0
	updated, err := svc.UpdateEffectivenessAssessment(testCtx(), created.ID, UpdateEffectivenessAssessmentRequest{AfterScore: &newAfter})
	if err != nil {
		t.Fatalf("UpdateEffectivenessAssessment failed: %v", err)
	}
	if updated.AfterScore == nil || *updated.AfterScore != newAfter {
		t.Errorf("expected after_score %v, got %v", newAfter, updated.AfterScore)
	}

	if err := svc.DeleteEffectivenessAssessment(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteEffectivenessAssessment failed: %v", err)
	}
	list, _ = svc.ListEffectivenessAssessments(testCtx(), strP(part.ID), 1, 20)
	if list.Total != 0 {
		t.Errorf("expected 0 assessments after delete, got %d", list.Total)
	}
}

// =========================================================================
// Certifications & Certificate Generation (P2-BE — plan §24)
// =========================================================================

func TestService_Certification_CRUD(t *testing.T) {
	svc := testSvc(t)

	created, err := svc.CreateCertification(testCtx(), CreateCertificationRequest{
		Code: "GCP-ARCH",
		Name: "Google Cloud Architect",
	})
	if err != nil {
		t.Fatalf("CreateCertification failed: %v", err)
	}
	if created.Code != "GCP-ARCH" {
		t.Errorf("unexpected code: %s", created.Code)
	}

	list, err := svc.ListCertifications(testCtx(), nil, 1, 20)
	if err != nil {
		t.Fatalf("ListCertifications failed: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("expected 1 certification, got %d", list.Total)
	}

	active := true
	onlyActive, err := svc.ListCertifications(testCtx(), &active, 1, 20)
	if err != nil {
		t.Fatalf("ListCertifications failed: %v", err)
	}
	if onlyActive.Total != 1 {
		t.Errorf("expected 1 active certification, got %d", onlyActive.Total)
	}

	if err := svc.DeleteCertification(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteCertification failed: %v", err)
	}
	list, _ = svc.ListCertifications(testCtx(), nil, 1, 20)
	if list.Total != 0 {
		t.Errorf("expected 0 certifications after delete, got %d", list.Total)
	}
}

// =========================================================================
// Reports & History smoke tests (P2-BE — plan §38)
// =========================================================================

// TestService_Reports_Smoke — pastikan raw SQL report berjalan di SQLite
// (portabilitas sintaks dasar; dijalankan juga di MySQL/PostgreSQL via migrator).
func TestService_Reports_Smoke(t *testing.T) {
	svc := testSvc(t)

	// Report queries join tabel modul lain (employees/employments/organizations)
	// yang tidak di-AutoMigrate di test DB — buat tabel minimal agar raw SQL
	// berjalan (di production tabel ini ada dari modul employee/organization).
	db, err := svc.repo.dbFunc(testCtx())
	if err != nil {
		t.Fatalf("dbFunc failed: %v", err)
	}
	for _, ddl := range []string{
		`CREATE TABLE IF NOT EXISTS employees (id TEXT PRIMARY KEY, name TEXT)`, // nolint:gosec
		`CREATE TABLE IF NOT EXISTS employments (id TEXT PRIMARY KEY, employee_id TEXT, organization_id TEXT, effective_end_date TEXT)`, // nolint:gosec
		`CREATE TABLE IF NOT EXISTS organizations (id TEXT PRIMARY KEY, nomenclature TEXT)`, // nolint:gosec
	} {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("create helper table failed: %v", err)
		}
	}

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)
	empID := "00000000-0000-0000-0000-000000000006"
	if _, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessionID, EmployeeID: empID,
	}); err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	dash, err := svc.GetDashboardReport(testCtx())
	if err != nil {
		t.Fatalf("GetDashboardReport failed: %v", err)
	}
	if dash.TotalCourses != 1 || dash.TotalSessions != 1 || dash.TotalParticipants != 1 {
		t.Errorf("unexpected dashboard counts: courses=%d sessions=%d participants=%d", dash.TotalCourses, dash.TotalSessions, dash.TotalParticipants)
	}

	history, err := svc.GetTrainingHistory(testCtx(), empID)
	if err != nil {
		t.Fatalf("GetTrainingHistory failed: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("expected 1 history row, got %d", len(history))
	}
	if history[0].CourseID != courseID {
		t.Errorf("unexpected course id: %s", history[0].CourseID)
	}

	parts, err := svc.GetParticipationReport(testCtx(), nil)
	if err != nil {
		t.Fatalf("GetParticipationReport failed: %v", err)
	}
	if len(parts) != 1 {
		t.Errorf("expected 1 participation row, got %d", len(parts))
	}

	costs, err := svc.GetCostReport(testCtx())
	if err != nil {
		t.Fatalf("GetCostReport failed: %v", err)
	}
	if len(costs) != 1 {
		t.Errorf("expected 1 cost row, got %d", len(costs))
	}

	// ComplianceReport: tanpa mandatory, hasil boleh kosong — yang penting tidak error.
	if _, err := svc.GetComplianceReport(testCtx()); err != nil {
		t.Fatalf("GetComplianceReport failed: %v", err)
	}
}

func TestService_GenerateCertificate_OnlyCompleted(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	empID := "00000000-0000-0000-0000-000000000005"
	p, err := svc.CreateParticipant(testCtx(), CreateTrainingParticipantRequest{
		SessionID: sessionID, EmployeeID: empID,
	})
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}

	// Participant belum COMPLETED → generate ditolak.
	if _, err := svc.GenerateCertificate(testCtx(), p.ID, GenerateCertificateRequest{}); err == nil {
		t.Error("expected error for non-completed participant, got nil")
	}

	// Tandai COMPLETED via service UpdateParticipant.
	comp := "COMPLETED"
	if _, err := svc.UpdateParticipant(testCtx(), p.ID, UpdateTrainingParticipantRequest{CompletionStatus: &comp}); err != nil {
		t.Fatalf("mark completed failed: %v", err)
	}

	cert, err := svc.GenerateCertificate(testCtx(), p.ID, GenerateCertificateRequest{})
	if err != nil {
		t.Fatalf("GenerateCertificate failed: %v", err)
	}
	if cert.CertificateNo == "" {
		t.Error("expected certificate_no to be generated")
	}
	if cert.IssuedDate == "" {
		t.Error("expected issued_date to be set")
	}
}
