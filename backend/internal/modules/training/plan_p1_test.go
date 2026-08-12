package training

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// =========================================================================
// fakeApprovalEngine — test double untuk training.ApprovalEngine (pola leave).
// =========================================================================

type fakeApprovalEngine struct {
	createCalls []struct {
		module     string
		documentID string
		flowID     string
	}
	instanceID    string
	createErr     error
	activeFlowID  string
	activeFlowErr error
}

func (f *fakeApprovalEngine) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	f.createCalls = append(f.createCalls, struct {
		module     string
		documentID string
		flowID     string
	}{module, documentID, flowID})
	if f.createErr != nil {
		return "", f.createErr
	}
	if f.instanceID == "" {
		f.instanceID = uuid.New().String()
	}
	return f.instanceID, nil
}

func (f *fakeApprovalEngine) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	return "PENDING", nil
}

func (f *fakeApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	if f.activeFlowErr != nil {
		return "", f.activeFlowErr
	}
	if f.activeFlowID == "" {
		return "", fmt.Errorf("no active flow configured")
	}
	return f.activeFlowID, nil
}

// =========================================================================
// Training Plan (P1-BE — plan §16)
// =========================================================================

func TestService_CreatePlan_DefaultsToDraft(t *testing.T) {
	svc := testSvc(t)

	resp, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{
		Code: "TP-2026-001",
		Name: "Annual Development Plan 2026",
		Year: 2026,
	})
	if err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected status DRAFT, got %s", resp.Status)
	}
	if resp.ID == "" {
		t.Error("expected non-empty id")
	}
}

func TestService_ListPlans_FiltersByYearAndStatus(t *testing.T) {
	svc := testSvc(t)

	if _, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{Code: "TP-2025-001", Name: "Plan 2025", Year: 2025}); err != nil {
		t.Fatalf("seed plan 2025 failed: %v", err)
	}
	if _, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{Code: "TP-2026-001", Name: "Plan 2026", Year: 2026}); err != nil {
		t.Fatalf("seed plan 2026 failed: %v", err)
	}

	year := 2026
	resp, err := svc.ListPlans(testCtx(), &year, nil, 1, 20)
	if err != nil {
		t.Fatalf("ListPlans failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 plan for year 2026, got %d", resp.Total)
	}
}

func TestService_UpdatePlan_ChangesStatus(t *testing.T) {
	svc := testSvc(t)

	created, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{Code: "TP-2026-002", Name: "Q1 Plan", Year: 2026})
	if err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}

	status := "ACTIVE"
	updated, err := svc.UpdatePlan(testCtx(), created.ID, UpdateTrainingPlanRequest{Status: &status})
	if err != nil {
		t.Fatalf("UpdatePlan failed: %v", err)
	}
	if updated.Status != "ACTIVE" {
		t.Errorf("expected status ACTIVE, got %s", updated.Status)
	}
}

func TestService_DeletePlan_SoftDelete(t *testing.T) {
	svc := testSvc(t)

	created, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{Code: "TP-2026-003", Name: "To Delete", Year: 2026})
	if err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}
	if err := svc.DeletePlan(testCtx(), created.ID); err != nil {
		t.Fatalf("DeletePlan failed: %v", err)
	}
	if _, err := svc.GetPlanByID(testCtx(), created.ID); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestService_CreatePlanItem_RequiresCourse(t *testing.T) {
	svc := testSvc(t)

	plan, err := svc.CreatePlan(testCtx(), CreateTrainingPlanRequest{Code: "TP-2026-004", Name: "With Items", Year: 2026})
	if err != nil {
		t.Fatalf("CreatePlan failed: %v", err)
	}
	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	item, err := svc.CreatePlanItem(testCtx(), plan.ID, CreateTrainingPlanItemRequest{
		CourseID:      courseID,
		TargetDate:    strP("2026-09-01"),
		Priority:      strP("HIGH"),
	})
	if err != nil {
		t.Fatalf("CreatePlanItem failed: %v", err)
	}
	if item.TrainingPlanID != plan.ID {
		t.Errorf("expected plan id %s, got %s", plan.ID, item.TrainingPlanID)
	}

	// Course tidak ada → error
	if _, err := svc.CreatePlanItem(testCtx(), plan.ID, CreateTrainingPlanItemRequest{CourseID: uuid.New().String()}); err == nil {
		t.Error("expected error for unknown course, got nil")
	}
}

// =========================================================================
// Training Need (P1-BE — plan §17)
// =========================================================================

func TestService_CreateNeed_DefaultsOpenAndManual(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	resp, err := svc.CreateNeed(testCtx(), CreateTrainingNeedRequest{
		CourseID: strP(courseID),
		Reason:   strP("skill gap analysis"),
	})
	if err != nil {
		t.Fatalf("CreateNeed failed: %v", err)
	}
	if resp.Status != "OPEN" {
		t.Errorf("expected status OPEN, got %s", resp.Status)
	}
	if resp.SourceType != "MANUAL" {
		t.Errorf("expected source_type MANUAL, got %s", resp.SourceType)
	}
}

func TestService_UpdateNeed_ToPlanned(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateNeed(testCtx(), CreateTrainingNeedRequest{CourseID: strP(courseID)})
	if err != nil {
		t.Fatalf("CreateNeed failed: %v", err)
	}

	status := "PLANNED"
	updated, err := svc.UpdateNeed(testCtx(), created.ID, UpdateTrainingNeedRequest{Status: &status})
	if err != nil {
		t.Fatalf("UpdateNeed failed: %v", err)
	}
	if updated.Status != "PLANNED" {
		t.Errorf("expected status PLANNED, got %s", updated.Status)
	}
}

// =========================================================================
// Training Request + Central Approval (P1-BE — plan §15/§45)
// =========================================================================

func TestService_CreateRequest_DefaultsDraft(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	resp, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}
	if resp.Status != "DRAFT" {
		t.Errorf("expected status DRAFT, got %s", resp.Status)
	}
}

func TestService_SubmitRequest_WithFlowID_CreatesApprovalInstance(t *testing.T) {
	svc := testSvc(t)
	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}

	flowID := uuid.New().String()
	resp, err := svc.SubmitRequest(testCtx(), created.ID, &flowID)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}
	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got %s", resp.Status)
	}
	if resp.ApprovalInstanceID == "" {
		t.Error("expected approval_instance_id to be set")
	}
	if len(fake.createCalls) != 1 {
		t.Fatalf("expected 1 CreateApprovalInstance call, got %d", len(fake.createCalls))
	}
	if fake.createCalls[0].module != "training_request" {
		t.Errorf("expected module training_request, got %s", fake.createCalls[0].module)
	}
}

func TestService_SubmitRequest_NoFlowID_AutoResolvesActiveFlow(t *testing.T) {
	svc := testSvc(t)
	activeFlowID := uuid.New().String()
	fake := &fakeApprovalEngine{activeFlowID: activeFlowID}
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}

	resp, err := svc.SubmitRequest(testCtx(), created.ID, nil)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}
	if resp.Status != "PENDING_APPROVAL" {
		t.Errorf("expected status PENDING_APPROVAL, got %s", resp.Status)
	}
	if len(fake.createCalls) != 1 || fake.createCalls[0].flowID != activeFlowID {
		t.Errorf("expected auto-resolved flow_id %s, got %+v", activeFlowID, fake.createCalls)
	}
}

func TestService_SubmitRequest_NoFlowAndNoActiveFlow_SkipsApproval(t *testing.T) {
	svc := testSvc(t)
	fake := &fakeApprovalEngine{} // no active flow
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}

	resp, err := svc.SubmitRequest(testCtx(), created.ID, nil)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}
	if resp.Status != "SUBMITTED" {
		t.Errorf("expected status SUBMITTED, got %s", resp.Status)
	}
	if len(fake.createCalls) != 0 {
		t.Errorf("expected no approval calls, got %d", len(fake.createCalls))
	}
}

func TestService_HandleApprovalStatusChange_Approved_AutoEnrolls(t *testing.T) {
	svc := testSvc(t)
	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)
	empID := uuid.New().String()

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    empID,
		CourseID:      courseID,
		SessionID:     &sessionID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}
	flowID := uuid.New().String()
	submitted, err := svc.SubmitRequest(testCtx(), created.ID, &flowID)
	if err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}
	if submitted.Status != "PENDING_APPROVAL" {
		t.Fatalf("expected PENDING_APPROVAL, got %s", submitted.Status)
	}

	uid, _ := uuid.Parse(created.ID)
	if err := svc.HandleApprovalStatusChange(testCtx(), uid, "APPROVED", "approved by HR"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetRequestByID(testCtx(), created.ID)
	if err != nil {
		t.Fatalf("GetRequestByID failed: %v", err)
	}
	if updated.Status != "APPROVED" {
		t.Errorf("expected status APPROVED, got %s", updated.Status)
	}
	if updated.SupervisorNote == "" {
		t.Error("expected supervisor_note to be persisted")
	}

	// Auto-enroll: participant untuk employee + session harus ada.
	participants, err := svc.ListParticipants(testCtx(), &sessionID, nil, 1, 50)
	if err != nil {
		t.Fatalf("ListParticipants failed: %v", err)
	}
	found := false
	if rows, ok := participants.Data.([]TrainingParticipantResponse); ok {
		for _, p := range rows {
			if p.EmployeeID == empID {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected auto-enrolled participant for approved request")
	}
}

func TestService_HandleApprovalStatusChange_Rejected(t *testing.T) {
	svc := testSvc(t)
	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}
	flowID := uuid.New().String()
	if _, err := svc.SubmitRequest(testCtx(), created.ID, &flowID); err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}

	uid, _ := uuid.Parse(created.ID)
	note := "budget not approved"
	if err := svc.HandleApprovalStatusChange(testCtx(), uid, "REJECTED", note); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	updated, err := svc.GetRequestByID(testCtx(), created.ID)
	if err != nil {
		t.Fatalf("GetRequestByID failed: %v", err)
	}
	if updated.Status != "REJECTED" {
		t.Errorf("expected status REJECTED, got %s", updated.Status)
	}
	if updated.SupervisorNote != note {
		t.Errorf("expected supervisor note %q, got %q", note, updated.SupervisorNote)
	}
}

func TestService_CancelRequest(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}

	resp, err := svc.CancelRequest(testCtx(), created.ID)
	if err != nil {
		t.Fatalf("CancelRequest failed: %v", err)
	}
	if resp.Status != "CANCELLED" {
		t.Errorf("expected status CANCELLED, got %s", resp.Status)
	}

	// Request CANCELLED tidak bisa di-submit ulang (hanya DRAFT/REJECTED).
	flowID := uuid.New().String()
	if _, err := svc.SubmitRequest(testCtx(), created.ID, &flowID); err == nil {
		t.Fatal("expected error re-submitting cancelled request, got nil")
	}
}

// TestService_CancelRequest_Approved_Fails — request yang sudah APPROVED tidak
// boleh dibatalkan (guard CancelRequest: status APPROVED/CANCELLED ditolak).
func TestService_CancelRequest_Approved_Fails(t *testing.T) {
	svc := testSvc(t)
	fake := &fakeApprovalEngine{}
	svc.SetApprovalEngine(fake)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}
	flowID := uuid.New().String()
	if _, err := svc.SubmitRequest(testCtx(), created.ID, &flowID); err != nil {
		t.Fatalf("SubmitRequest failed: %v", err)
	}
	uid, _ := uuid.Parse(created.ID)
	if err := svc.HandleApprovalStatusChange(testCtx(), uid, "APPROVED", "ok"); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}

	if _, err := svc.CancelRequest(testCtx(), created.ID); err == nil {
		t.Fatal("expected error cancelling approved request, got nil")
	}
}

func TestService_DeleteNeed_SoftDelete(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateNeed(testCtx(), CreateTrainingNeedRequest{CourseID: strP(courseID)})
	if err != nil {
		t.Fatalf("CreateNeed failed: %v", err)
	}
	if err := svc.DeleteNeed(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteNeed failed: %v", err)
	}
	if _, err := svc.GetNeedByID(testCtx(), created.ID); err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestService_HandleApprovalStatusChange_NotPending_NoOp(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateRequest(testCtx(), CreateTrainingRequestRequest{
		EmployeeID:    uuid.New().String(),
		CourseID:      courseID,
		RequestedDate: "2026-08-20",
	})
	if err != nil {
		t.Fatalf("CreateRequest failed: %v", err)
	}

	uid, _ := uuid.Parse(created.ID)
	if err := svc.HandleApprovalStatusChange(testCtx(), uid, "APPROVED", ""); err != nil {
		t.Fatalf("HandleApprovalStatusChange failed: %v", err)
	}
	updated, err := svc.GetRequestByID(testCtx(), created.ID)
	if err != nil {
		t.Fatalf("GetRequestByID failed: %v", err)
	}
	if updated.Status != "DRAFT" {
		t.Errorf("expected status to remain DRAFT, got %s", updated.Status)
	}
}

// =========================================================================
// Helpers lokal (nama unik — strPtr sudah dipakai service.go sebagai
// converter dereference, bukan pembuat pointer)
// =========================================================================

func strP(s string) *string { return &s }

func intP(i int) *int { return &i }

// =========================================================================
// Course sub-resources (P1-BE — plan §8/§9/§10)
// =========================================================================

func TestService_CourseObjectives_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateCourseObjective(testCtx(), courseID, CreateCourseObjectiveRequest{
		Objective: "Understand Go concurrency",
	})
	if err != nil {
		t.Fatalf("CreateCourseObjective failed: %v", err)
	}
	if created.Objective != "Understand Go concurrency" {
		t.Errorf("unexpected objective: %s", created.Objective)
	}

	list, err := svc.ListCourseObjectives(testCtx(), courseID)
	if err != nil {
		t.Fatalf("ListCourseObjectives failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 objective, got %d", len(list))
	}

	updated, err := svc.UpdateCourseObjective(testCtx(), created.ID, UpdateCourseObjectiveRequest{
		Objective: strP("Master Go concurrency patterns"),
	})
	if err != nil {
		t.Fatalf("UpdateCourseObjective failed: %v", err)
	}
	if updated.Objective != "Master Go concurrency patterns" {
		t.Errorf("unexpected updated objective: %s", updated.Objective)
	}

	if err := svc.DeleteCourseObjective(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteCourseObjective failed: %v", err)
	}
	list, _ = svc.ListCourseObjectives(testCtx(), courseID)
	if len(list) != 0 {
		t.Errorf("expected 0 objectives after delete, got %d", len(list))
	}
}

func TestService_CourseCompetency_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateCourseCompetency(testCtx(), courseID, CreateCourseCompetencyRequest{
		CompetencyID: uuid.New().String(),
		TargetLevel:  intP(4),
	})
	if err != nil {
		t.Fatalf("CreateCourseCompetency failed: %v", err)
	}
	if created.CourseID != courseID {
		t.Errorf("expected course id %s, got %s", courseID, created.CourseID)
	}

	list, err := svc.ListCourseCompetencies(testCtx(), courseID)
	if err != nil {
		t.Fatalf("ListCourseCompetencies failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 competency, got %d", len(list))
	}

	if err := svc.DeleteCourseCompetency(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteCourseCompetency failed: %v", err)
	}
	list, _ = svc.ListCourseCompetencies(testCtx(), courseID)
	if len(list) != 0 {
		t.Errorf("expected 0 competencies after delete, got %d", len(list))
	}
}

func TestService_SessionCost_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	amount := 2500000.0
	created, err := svc.CreateSessionCost(testCtx(), sessionID, CreateTrainingSessionCostRequest{
		CostType: "TRAINER",
		Amount:   &amount,
	})
	if err != nil {
		t.Fatalf("CreateSessionCost failed: %v", err)
	}
	if created.SessionID != sessionID {
		t.Errorf("expected session id %s, got %s", sessionID, created.SessionID)
	}

	list, err := svc.ListSessionCosts(testCtx(), sessionID)
	if err != nil {
		t.Fatalf("ListSessionCosts failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 cost, got %d", len(list))
	}

	newAmount := 3000000.0
	updated, err := svc.UpdateSessionCost(testCtx(), created.ID, UpdateTrainingSessionCostRequest{Amount: &newAmount})
	if err != nil {
		t.Fatalf("UpdateSessionCost failed: %v", err)
	}
	if updated.Amount != newAmount {
		t.Errorf("expected amount %v, got %v", newAmount, updated.Amount)
	}

	if err := svc.DeleteSessionCost(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteSessionCost failed: %v", err)
	}
	list, _ = svc.ListSessionCosts(testCtx(), sessionID)
	if len(list) != 0 {
		t.Errorf("expected 0 costs after delete, got %d", len(list))
	}
}

func TestService_Document_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)
	sessionID := seedSession(t, svc, courseID)

	created, err := svc.CreateDocument(testCtx(), sessionID, CreateTrainingDocumentRequest{
		DocumentType: "PROPOSAL",
		FileURL:      "/uploads/training/proposal.pdf",
	})
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}
	if created.SessionID != sessionID {
		t.Errorf("expected session id %s, got %s", sessionID, created.SessionID)
	}

	list, err := svc.ListDocuments(testCtx(), sessionID)
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 document, got %d", len(list))
	}

	if err := svc.DeleteDocument(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteDocument failed: %v", err)
	}
	list, _ = svc.ListDocuments(testCtx(), sessionID)
	if len(list) != 0 {
		t.Errorf("expected 0 documents after delete, got %d", len(list))
	}
}

func TestService_CoursePrerequisite_RequiresCourse(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	preq, err := svc.CreateCoursePrerequisite(testCtx(), courseID, CreateCoursePrerequisiteRequest{
		PrerequisiteType: "COURSE",
		PrerequisiteID:   strP(uuid.New().String()),
	})
	if err != nil {
		t.Fatalf("CreateCoursePrerequisite failed: %v", err)
	}
	if preq.CourseID != courseID {
		t.Errorf("expected course id %s, got %s", courseID, preq.CourseID)
	}

	list, err := svc.ListCoursePrerequisites(testCtx(), courseID)
	if err != nil {
		t.Fatalf("ListCoursePrerequisites failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 prerequisite, got %d", len(list))
	}
}

func TestService_Mandatory_CRUD(t *testing.T) {
	svc := testSvc(t)

	catID := seedCategory(t, svc)
	courseID := seedCourse(t, svc, catID)

	created, err := svc.CreateMandatory(testCtx(), CreateTrainingMandatoryRequest{
		CourseID: courseID,
		DueDays:  intP(90),
	})
	if err != nil {
		t.Fatalf("CreateMandatory failed: %v", err)
	}
	if created.CourseID != courseID {
		t.Errorf("unexpected course id: %s", created.CourseID)
	}

	resp, err := svc.ListMandatories(testCtx(), &courseID, 1, 20)
	if err != nil {
		t.Fatalf("ListMandatories failed: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("expected 1 mandatory, got %d", resp.Total)
	}

	if err := svc.DeleteMandatory(testCtx(), created.ID); err != nil {
		t.Fatalf("DeleteMandatory failed: %v", err)
	}
	resp, _ = svc.ListMandatories(testCtx(), &courseID, 1, 20)
	if resp.Total != 0 {
		t.Errorf("expected 0 mandatories after delete, got %d", resp.Total)
	}
}

// =========================================================================
// Onboarding → Training Handoff (S-7) Service Tests
// =========================================================================

func TestService_CreateOnboardingNeed_SetsSourceAndReason(t *testing.T) {
	svc := testSvc(t)
	ctx := testCtx()

	empID := uuid.New()
	onbID := uuid.New()
	resp, err := svc.CreateOnboardingNeed(ctx, empID, onbID, "")
	if err != nil {
		t.Fatalf("CreateOnboardingNeed failed: %v", err)
	}
	if resp.SourceType != string(NeedSourceOnboarding) {
		t.Errorf("expected source_type ONBOARDING, got %s", resp.SourceType)
	}
	if resp.SourceID != onbID.String() {
		t.Errorf("expected source_id %s (onboarding id), got %s", onbID.String(), resp.SourceID)
	}
	if resp.EmployeeID != empID.String() {
		t.Errorf("expected employee_id %s, got %s", empID.String(), resp.EmployeeID)
	}
	if resp.Status != "OPEN" {
		t.Errorf("expected status OPEN, got %s", resp.Status)
	}
	if resp.Reason == "" {
		t.Error("expected default reason to be set")
	}

	// Verifikasi tersimpan di DB (bukan hanya response in-memory).
	found, err := svc.GetNeedByID(ctx, resp.ID)
	if err != nil {
		t.Fatalf("GetNeedByID failed: %v", err)
	}
	if found.SourceType != string(NeedSourceOnboarding) || found.SourceID != onbID.String() {
		t.Errorf("expected persisted need source ONBOARDING + source_id, got %+v", found)
	}
}

func TestService_CreateOnboardingNeed_CustomReason(t *testing.T) {
	svc := testSvc(t)
	ctx := testCtx()

	resp, err := svc.CreateOnboardingNeed(ctx, uuid.New(), uuid.New(), "IT orientation required")
	if err != nil {
		t.Fatalf("CreateOnboardingNeed failed: %v", err)
	}
	if resp.Reason != "IT orientation required" {
		t.Errorf("expected custom reason, got %s", resp.Reason)
	}
}

