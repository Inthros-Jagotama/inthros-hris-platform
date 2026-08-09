package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// fakeNotifier is the same capturing double used by leave/attendance
// notifier integration tests — records every Notify call so the test can
// assert recipient, type, and reference without touching the real
// notification module.
type fakeNotifier struct {
	calls []fakeNotifyCall
	err   error
}

type fakeNotifyCall struct {
	recipientUserID uuid.UUID
	notifType       string
	referenceType   string
	referenceID     uuid.UUID
}

func (f *fakeNotifier) Notify(ctx context.Context, recipientUserID uuid.UUID, notifType string, params []string, referenceType string, referenceID uuid.UUID) error {
	f.calls = append(f.calls, fakeNotifyCall{recipientUserID, notifType, referenceType, referenceID})
	return f.err
}

// seedEmployeeAccount links an employee to a platform user in the raw
// employee_accounts table (same shape as my_kpi_context_test.go's insert).
func seedEmployeeAccount(t *testing.T, db *gorm.DB, employeeID, userID uuid.UUID) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO employee_accounts (id, employee_id, user_id) VALUES (?, ?, ?)",
		uuid.New().String(), employeeID.String(), userID.String(),
	).Error; err != nil {
		t.Fatalf("failed to seed employee account: %v", err)
	}
}

// seedKPIMember links an employee to an existing Organization (current,
// open-ended employment) plus a platform user account — used to add extra
// members to an organization that seedMyKPIContextEmployee already created.
func seedKPIMember(t *testing.T, db *gorm.DB, employeeID, userID, orgID uuid.UUID) {
	t.Helper()
	if err := db.Exec(
		"INSERT INTO employments (id, employee_id, organization_id, effective_date, effective_end_date) VALUES (?, ?, ?, ?, NULL)",
		uuid.New().String(), employeeID.String(), orgID.String(), "2026-01-01",
	).Error; err != nil {
		t.Fatalf("failed to seed employment: %v", err)
	}
	seedEmployeeAccount(t, db, employeeID, userID)
}

// seedKPIEvaluation inserts a minimal PerformanceEvaluation directly (the
// approval status handlers only read/update the header row; period/template/
// organization are not dereferenced by them).
func seedKPIEvaluation(t *testing.T, db *gorm.DB, employeeID uuid.UUID, status string) *PerformanceEvaluation {
	t.Helper()
	eval := &PerformanceEvaluation{
		EmployeeID:     employeeID,
		OrganizationID: uuid.New(),
		PeriodID:       uuid.New(),
		TemplateID:     uuid.New(),
		Status:         status,
	}
	if err := db.Create(eval).Error; err != nil {
		t.Fatalf("failed to seed KPI evaluation: %v", err)
	}
	return eval
}

// =========================================================================
// KPI — Target approval
// =========================================================================

func TestService_HandleTargetApprovalStatusChange_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	seedEmployeeAccount(t, db, empID, userID)
	eval := seedKPIEvaluation(t, db, empID, "TARGET_SUBMITTED")

	if err := svc.HandleTargetApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "KPI_TARGET_APPROVED" {
		t.Errorf("expected notif type KPI_TARGET_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "performance_kpi_target" || call.referenceID != eval.ID {
		t.Errorf("expected reference performance_kpi_target/%s, got %s/%s", eval.ID, call.referenceType, call.referenceID)
	}
}

func TestService_HandleTargetApprovalStatusChange_NotifiesEmployeeOnRejection(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	seedEmployeeAccount(t, db, empID, uuid.New())
	eval := seedKPIEvaluation(t, db, empID, "TARGET_SUBMITTED")

	if err := svc.HandleTargetApprovalStatusChange(ctx, eval.ID, "REJECTED", "please revise"); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "KPI_TARGET_REJECTED" {
		t.Errorf("expected notif type KPI_TARGET_REJECTED, got %s", notifier.calls[0].notifType)
	}
}

func TestService_HandleTargetApprovalStatusChange_SkipsNotifyWhenNoLinkedUserAccount(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New() // no employee_accounts row created
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	eval := seedKPIEvaluation(t, db, empID, "TARGET_SUBMITTED")

	if err := svc.HandleTargetApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleTargetApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls when employee has no linked user account, got %d", len(notifier.calls))
	}
}

// =========================================================================
// KPI — Realization approval
// =========================================================================

func TestService_HandleRealizationApprovalStatusChange_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	seedEmployeeAccount(t, db, empID, userID)
	eval := seedKPIEvaluation(t, db, empID, "SUBMITTED")

	if err := svc.HandleRealizationApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleRealizationApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "KPI_REALIZATION_APPROVED" {
		t.Errorf("expected notif type KPI_REALIZATION_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "performance_kpi_realization" || call.referenceID != eval.ID {
		t.Errorf("expected reference performance_kpi_realization/%s, got %s/%s", eval.ID, call.referenceType, call.referenceID)
	}
}

func TestService_HandleRealizationApprovalStatusChange_NotifiesEmployeeOnRejection(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	seedEmployeeAccount(t, db, empID, uuid.New())
	eval := seedKPIEvaluation(t, db, empID, "SUBMITTED")

	if err := svc.HandleRealizationApprovalStatusChange(ctx, eval.ID, "REJECTED", "not accurate"); err != nil {
		t.Fatalf("HandleRealizationApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "KPI_REALIZATION_REJECTED" {
		t.Errorf("expected notif type KPI_REALIZATION_REJECTED, got %s", notifier.calls[0].notifType)
	}
}

func TestService_HandleRealizationApprovalStatusChange_NotifyFailureDoesNotFailApproval(t *testing.T) {
	svc, resolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{err: context.DeadlineExceeded}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	db, err := resolver(ctx)
	if err != nil {
		t.Fatalf("resolver failed: %v", err)
	}
	seedEmployeeAccount(t, db, empID, uuid.New())
	eval := seedKPIEvaluation(t, db, empID, "SUBMITTED")

	if err := svc.HandleRealizationApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("expected HandleRealizationApprovalStatusChange to succeed despite notify failure, got: %v", err)
	}
	if len(notifier.calls) != 1 {
		t.Fatalf("expected the notifier to still be called once despite it returning an error, got %d calls", len(notifier.calls))
	}
}

// =========================================================================
// OKR — Key Result approval
// =========================================================================

// seedOKREvaluation inserts a minimal OKREvaluation header directly.
func seedOKREvaluation(t *testing.T, db *gorm.DB, employeeID uuid.UUID, status OKREvaluationStatus) *OKREvaluation {
	t.Helper()
	eval := &OKREvaluation{
		EmployeeID:     employeeID,
		OrganizationID: uuid.New(),
		PeriodID:       uuid.New(),
		Status:         status,
	}
	if err := db.Create(eval).Error; err != nil {
		t.Fatalf("failed to seed OKR evaluation: %v", err)
	}
	return eval
}

func TestOKR_HandleKeyResultApprovalStatusChange_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	seedOKREmployeeAccount(t, db, empID, userID)
	eval := seedOKREvaluation(t, db, empID, OKRStatusKRSubmitted)

	if err := svc.HandleKeyResultApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "OKR_KEY_RESULT_APPROVED" {
		t.Errorf("expected notif type OKR_KEY_RESULT_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "okr_key_result" || call.referenceID != eval.ID {
		t.Errorf("expected reference okr_key_result/%s, got %s/%s", eval.ID, call.referenceType, call.referenceID)
	}
}

func TestOKR_HandleKeyResultApprovalStatusChange_NotifiesEmployeeOnRejection(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	seedOKREmployeeAccount(t, db, empID, uuid.New())
	eval := seedOKREvaluation(t, db, empID, OKRStatusKRSubmitted)

	if err := svc.HandleKeyResultApprovalStatusChange(ctx, eval.ID, "REJECTED", "revise KR"); err != nil {
		t.Fatalf("HandleKeyResultApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "OKR_KEY_RESULT_REJECTED" {
		t.Errorf("expected notif type OKR_KEY_RESULT_REJECTED, got %s", notifier.calls[0].notifType)
	}
}

// =========================================================================
// OKR — Assessment approval
// =========================================================================

func TestOKR_HandleAssessmentApprovalStatusChange_NotifiesEmployeeOnApproval(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	userID := uuid.New()
	seedOKREmployeeAccount(t, db, empID, userID)
	eval := seedOKREvaluation(t, db, empID, OKRStatusSubmitted)

	if err := svc.HandleAssessmentApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleAssessmentApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.recipientUserID != userID {
		t.Errorf("expected recipient %s, got %s", userID, call.recipientUserID)
	}
	if call.notifType != "OKR_ASSESSMENT_APPROVED" {
		t.Errorf("expected notif type OKR_ASSESSMENT_APPROVED, got %s", call.notifType)
	}
	if call.referenceType != "okr_assessment" || call.referenceID != eval.ID {
		t.Errorf("expected reference okr_assessment/%s, got %s/%s", eval.ID, call.referenceType, call.referenceID)
	}
}

func TestOKR_HandleAssessmentApprovalStatusChange_NotifiesEmployeeOnRejection(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New()
	seedOKREmployeeAccount(t, db, empID, uuid.New())
	eval := seedOKREvaluation(t, db, empID, OKRStatusSubmitted)

	if err := svc.HandleAssessmentApprovalStatusChange(ctx, eval.ID, "REJECTED", "revise assessment"); err != nil {
		t.Fatalf("HandleAssessmentApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(notifier.calls))
	}
	if notifier.calls[0].notifType != "OKR_ASSESSMENT_REJECTED" {
		t.Errorf("expected notif type OKR_ASSESSMENT_REJECTED, got %s", notifier.calls[0].notifType)
	}
}

func TestOKR_HandleAssessmentApprovalStatusChange_SkipsNotifyWhenNoLinkedUserAccount(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()
	ctx := context.Background()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	empID := uuid.New() // no employee_accounts row created
	eval := seedOKREvaluation(t, db, empID, OKRStatusSubmitted)

	if err := svc.HandleAssessmentApprovalStatusChange(ctx, eval.ID, "APPROVED", ""); err != nil {
		t.Fatalf("HandleAssessmentApprovalStatusChange failed: %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notify calls when employee has no linked user account, got %d", len(notifier.calls))
	}
}

// =========================================================================
// Template creation — notify all employees in the template's organization
// =========================================================================

func TestService_CreatePerformanceTemplate_NotifiesOrganizationMembers(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	// Two more employees occupying the same org (both get notified).
	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}
	member1EmpID := uuid.New()
	member1UserID := uuid.New()
	seedKPIMember(t, db, member1EmpID, member1UserID, orgID)
	member2EmpID := uuid.New()
	member2UserID := uuid.New()
	seedKPIMember(t, db, member2EmpID, member2UserID, orgID)

	resp, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Org BSC Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}
	if resp.CreatedBy != creatorID.String() {
		t.Errorf("expected created_by %s, got %s", creatorID, resp.CreatedBy)
	}

	// creator + 2 members = 3 notifications, all with KPI_TEMPLATE_CREATED.
	if len(notifier.calls) != 3 {
		t.Fatalf("expected 3 notify calls (creator + 2 members), got %d", len(notifier.calls))
	}
	recipients := map[uuid.UUID]bool{}
	for _, call := range notifier.calls {
		if call.notifType != "KPI_TEMPLATE_CREATED" {
			t.Errorf("expected notif type KPI_TEMPLATE_CREATED, got %s", call.notifType)
		}
		if call.referenceType != "performance_kpi_template" {
			t.Errorf("expected reference_type performance_kpi_template, got %s", call.referenceType)
		}
		recipients[call.recipientUserID] = true
	}
	for _, expected := range []uuid.UUID{creatorID, member1UserID, member2UserID} {
		if !recipients[expected] {
			t.Errorf("expected notification for user %s, missing", expected)
		}
	}
}

func TestOKR_CreateTemplate_NotifiesOrganizationMembers(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)

	rootOrgID := uuid.New()
	subOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, subOrgID, &rootOrgID, "Subordinate Org")

	rootEmpID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmpID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmpID, rootUserID)

	// Two employees occupying the subordinate org (both get notified).
	subEmp1 := uuid.New()
	subUser1 := uuid.New()
	seedOKREmployment(t, db, subEmp1, subOrgID)
	seedOKREmployeeAccount(t, db, subEmp1, subUser1)
	subEmp2 := uuid.New()
	subUser2 := uuid.New()
	seedOKREmployment(t, db, subEmp2, subOrgID)
	seedOKREmployeeAccount(t, db, subEmp2, subUser2)

	resp, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: subOrgID.String(),
		Name:           "Subordinate OKR Template",
	})
	if err != nil {
		t.Fatalf("CreateTemplate failed: %v", err)
	}
	if resp.CreatedBy != rootUserID.String() {
		t.Errorf("expected created_by %s, got %s", rootUserID, resp.CreatedBy)
	}

	// The two subordinate-org members get notified (creator lives in a
	// different org, so not included).
	if len(notifier.calls) != 2 {
		t.Fatalf("expected 2 notify calls, got %d", len(notifier.calls))
	}
	recipients := map[uuid.UUID]bool{}
	for _, call := range notifier.calls {
		if call.notifType != "OKR_TEMPLATE_CREATED" {
			t.Errorf("expected notif type OKR_TEMPLATE_CREATED, got %s", call.notifType)
		}
		if call.referenceType != "okr_template" {
			t.Errorf("expected reference_type okr_template, got %s", call.referenceType)
		}
		recipients[call.recipientUserID] = true
	}
	for _, expected := range []uuid.UUID{subUser1, subUser2} {
		if !recipients[expected] {
			t.Errorf("expected notification for user %s, missing", expected)
		}
	}
}

// =========================================================================
// Template edit/delete — org membership based: only users currently in the
// template's organization may modify it (never when PUBLISHED), regardless
// of who created it.
// =========================================================================

func TestService_UpdatePerformanceTemplate_RejectsOtherOrg(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Org Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	// A user who has since moved to another org cannot edit the template,
	// even though they created it.
	otherOrgID := uuid.MustParse(createTestOrgID())
	movedID := uuid.New()
	movedEmpID := uuid.New()
	seedMyKPIContextEmployee(t, dbResolver, movedID, movedEmpID, otherOrgID, "Other Org")

	name := "Hijacked"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(movedID), created.ID, UpdatePerformanceTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected update by user in a different org to be rejected")
	}
}

func TestService_UpdatePerformanceTemplate_NoCreatorOrg_Denied(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	userID := uuid.New()
	empID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, userID, empID, orgID, "Org")
	ctx := ctxWithUser(userID)

	// Template legacy: dibuat tanpa user context sehingga created_by_org_id NULL.
	// Di bawah aturan strict, template ini tidak bisa dikelola oleh siapa pun —
	// termasuk user dari organisasi yang sama dengan template.
	legacy, err := svc.CreatePerformanceTemplate(context.Background(), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Legacy Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	name := "Renamed"
	if _, err := svc.UpdatePerformanceTemplate(ctx, legacy.ID, UpdatePerformanceTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected legacy template (no created_by_org_id) update to be rejected")
	}
}

func TestService_UpdatePerformanceTemplate_SameOrgMemberAllowed(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Org Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	// A different member of the SAME org may edit the template.
	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}
	memberID := uuid.New()
	memberEmpID := uuid.New()
	seedKPIMember(t, db, memberEmpID, memberID, orgID)

	name := "Renamed by Member"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(memberID), created.ID, UpdatePerformanceTemplateRequest{Name: &name}); err != nil {
		t.Fatalf("expected same-org member update to succeed: %v", err)
	}
}

func TestService_UpdatePerformanceTemplate_CreatorOrgControls(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	// Creator works in creatorOrg but creates a template FOR targetOrg.
	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	creatorOrgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, creatorOrgID, "Creator Org")

	targetOrgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, uuid.New(), uuid.New(), targetOrgID, "Target Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: targetOrgID.String(),
		Name:           "For Target Org",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}
	// Organisasi pembuat (bukan org template) harus tersimpan di respons.
	if created.CreatedByOrgID != creatorOrgID.String() {
		t.Errorf("expected created_by_org_id %s, got %s", creatorOrgID, created.CreatedByOrgID)
	}

	db, err := dbResolver(context.Background())
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	// Member of the TARGET org (template's org) cannot manage it — hanya
	// anggota organisasi pembuat yang boleh.
	targetMemberID := uuid.New()
	targetMemberEmpID := uuid.New()
	seedKPIMember(t, db, targetMemberEmpID, targetMemberID, targetOrgID)
	name := "Hijacked"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(targetMemberID), created.ID, UpdatePerformanceTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected template-org member (not creator org) to be rejected")
	}

	// Another member of the CREATOR org can manage it.
	peerID := uuid.New()
	peerEmpID := uuid.New()
	seedKPIMember(t, db, peerEmpID, peerID, creatorOrgID)
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(peerID), created.ID, UpdatePerformanceTemplateRequest{Name: &name}); err != nil {
		t.Fatalf("expected creator-org member update to succeed: %v", err)
	}
}

func TestService_UpdatePerformanceTemplate_RejectsPublished(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Publishable Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	published := "PUBLISHED"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(creatorID), created.ID, UpdatePerformanceTemplateRequest{Status: &published}); err != nil {
		t.Fatalf("publishing should succeed for an org member: %v", err)
	}

	name := "Renamed"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(creatorID), created.ID, UpdatePerformanceTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected update of a PUBLISHED template to be rejected")
	}
}

func TestService_DeletePerformanceTemplate_RejectsOtherOrg(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Deletable Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	otherOrgID := uuid.MustParse(createTestOrgID())
	movedID := uuid.New()
	movedEmpID := uuid.New()
	seedMyKPIContextEmployee(t, dbResolver, movedID, movedEmpID, otherOrgID, "Other Org")

	if err := svc.DeletePerformanceTemplate(ctxWithUser(movedID), created.ID); err == nil {
		t.Fatal("expected delete by user in a different org to be rejected")
	}
}

func TestService_DeletePerformanceTemplate_RejectsPublished(t *testing.T) {
	svc, dbResolver, cleanup := setupMyKPIContextTestDB(t)
	defer cleanup()

	creatorID := uuid.New()
	creatorEmpID := uuid.New()
	orgID := uuid.MustParse(createTestOrgID())
	seedMyKPIContextEmployee(t, dbResolver, creatorID, creatorEmpID, orgID, "Org")

	created, err := svc.CreatePerformanceTemplate(ctxWithUser(creatorID), CreatePerformanceTemplateRequest{
		OrganizationID: orgID.String(),
		Name:           "Published Template",
	})
	if err != nil {
		t.Fatalf("CreatePerformanceTemplate failed: %v", err)
	}

	published := "PUBLISHED"
	if _, err := svc.UpdatePerformanceTemplate(ctxWithUser(creatorID), created.ID, UpdatePerformanceTemplateRequest{Status: &published}); err != nil {
		t.Fatalf("publishing should succeed for an org member: %v", err)
	}

	if err := svc.DeletePerformanceTemplate(ctxWithUser(creatorID), created.ID); err == nil {
		t.Fatal("expected delete of a PUBLISHED template to be rejected")
	}
}

func TestOKR_ListTemplates_EnrichesOrgAndPeriod(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Marketing Dept")

	period := &PerformancePeriod{PeriodCode: "Q2-2026", PeriodType: "QUARTERLY", Year: 2026}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("failed to seed period: %v", err)
	}

	tpl := &OKRTemplate{OrganizationID: orgID, PeriodID: &period.ID, Name: "Template A", Status: 0}
	if err := db.Create(tpl).Error; err != nil {
		t.Fatalf("failed to seed OKR template: %v", err)
	}

	list, total, err := svc.ListTemplates(db, nil, nil, nil, 1, 50)
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected 1 template, got total=%d len=%d", total, len(list))
	}
	if list[0].OrganizationName != "Marketing Dept" {
		t.Errorf("expected organization_name 'Marketing Dept', got '%s'", list[0].OrganizationName)
	}
	if list[0].PeriodCode != "Q2-2026" {
		t.Errorf("expected period_code 'Q2-2026', got '%s'", list[0].PeriodCode)
	}
}

func TestOKR_UpdateTemplate_RejectsOtherOrg(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	tpl := seedOKRTemplateWithOwner(t, db, orgID, 0, uuid.New())

	// User currently in a different org cannot edit the template.
	otherOrgID := uuid.New()
	seedOKROrganization(t, db, otherOrgID, nil, "Other Org")
	otherEmpID := uuid.New()
	otherUserID := uuid.New()
	seedOKREmployment(t, db, otherEmpID, otherOrgID)
	seedOKREmployeeAccount(t, db, otherEmpID, otherUserID)

	name := "Hijacked"
	if _, err := svc.UpdateTemplate(db, tpl.ID, otherUserID, &UpdateOKRTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected update by user in a different org to be rejected")
	}
}

func TestOKR_UpdateTemplate_NoCreatorOrg_Denied(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	memberEmpID := uuid.New()
	memberUserID := uuid.New()
	seedOKREmployment(t, db, memberEmpID, orgID)
	seedOKREmployeeAccount(t, db, memberEmpID, memberUserID)

	// Template legacy tanpa created_by_org_id — strict rule: ditolak walau user
	// berada di organisasi yang sama dengan template.
	tpl := &OKRTemplate{OrganizationID: orgID, Name: "Legacy", Status: 0, CreatedBy: uuid.Nil}
	if err := db.Create(tpl).Error; err != nil {
		t.Fatalf("failed to seed legacy OKR template: %v", err)
	}

	name := "Renamed"
	if _, err := svc.UpdateTemplate(db, tpl.ID, memberUserID, &UpdateOKRTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected legacy template (no created_by_org_id) update to be rejected")
	}
}

func TestOKR_UpdateTemplate_SameOrgMemberAllowed(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	tpl := seedOKRTemplateWithOwner(t, db, orgID, 0, uuid.New())

	// A different member of the SAME org may edit the template.
	memberEmpID := uuid.New()
	memberUserID := uuid.New()
	seedOKREmployment(t, db, memberEmpID, orgID)
	seedOKREmployeeAccount(t, db, memberEmpID, memberUserID)

	name := "Renamed by Member"
	if _, err := svc.UpdateTemplate(db, tpl.ID, memberUserID, &UpdateOKRTemplateRequest{Name: &name}); err != nil {
		t.Fatalf("expected same-org member update to succeed: %v", err)
	}
}

func TestOKR_UpdateTemplate_CreatorOrgControls(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	subOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, subOrgID, &rootOrgID, "Sub Org")

	// Creator works in root org; subOrg must be occupied so it counts as an
	// effective subordinate for CreateTemplate.
	subEmpID := uuid.New()
	subUserID := uuid.New()
	seedOKREmployment(t, db, subEmpID, subOrgID)
	seedOKREmployeeAccount(t, db, subEmpID, subUserID)

	rootEmpID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmpID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmpID, rootUserID)

	// Creator (root org) creates a template FOR the subordinate org.
	resp, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: subOrgID.String(),
		Name:           "Sub Template",
	})
	if err != nil {
		t.Fatalf("CreateTemplate failed: %v", err)
	}
	// Organisasi pembuat (bukan org template) harus tersimpan di respons.
	if resp.CreatedByOrgID != rootOrgID.String() {
		t.Errorf("expected created_by_org_id %s, got %s", rootOrgID, resp.CreatedByOrgID)
	}

	// Member of the subordinate org (template's org) cannot manage it.
	name := "Hijacked"
	if _, err := svc.UpdateTemplate(db, uuid.MustParse(resp.ID), subUserID, &UpdateOKRTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected subordinate-org member (not creator org) to be rejected")
	}

	// Another member of the CREATOR org can manage it.
	peerEmpID := uuid.New()
	peerUserID := uuid.New()
	seedOKREmployment(t, db, peerEmpID, rootOrgID)
	seedOKREmployeeAccount(t, db, peerEmpID, peerUserID)
	if _, err := svc.UpdateTemplate(db, uuid.MustParse(resp.ID), peerUserID, &UpdateOKRTemplateRequest{Name: &name}); err != nil {
		t.Fatalf("expected creator-org member update to succeed: %v", err)
	}
}

func TestOKR_UpdateTemplate_RejectsPublished(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	memberEmpID := uuid.New()
	memberUserID := uuid.New()
	seedOKREmployment(t, db, memberEmpID, orgID)
	seedOKREmployeeAccount(t, db, memberEmpID, memberUserID)

	tpl := seedOKRTemplateWithOwner(t, db, orgID, OKRTemplateStatusPublished, memberUserID)

	name := "Renamed"
	if _, err := svc.UpdateTemplate(db, tpl.ID, memberUserID, &UpdateOKRTemplateRequest{Name: &name}); err == nil {
		t.Fatal("expected update of a published OKR template to be rejected")
	}
}

func TestOKR_DeleteTemplate_RejectsOtherOrg(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	tpl := seedOKRTemplateWithOwner(t, db, orgID, 0, uuid.New())

	otherOrgID := uuid.New()
	seedOKROrganization(t, db, otherOrgID, nil, "Other Org")
	otherEmpID := uuid.New()
	otherUserID := uuid.New()
	seedOKREmployment(t, db, otherEmpID, otherOrgID)
	seedOKREmployeeAccount(t, db, otherEmpID, otherUserID)

	if err := svc.DeleteTemplate(db, tpl.ID, otherUserID); err == nil {
		t.Fatal("expected delete by user in a different org to be rejected")
	}
}

func TestOKR_DeleteTemplate_RejectsPublished(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	orgID := uuid.New()
	seedOKROrganization(t, db, orgID, nil, "Org")
	memberEmpID := uuid.New()
	memberUserID := uuid.New()
	seedOKREmployment(t, db, memberEmpID, orgID)
	seedOKREmployeeAccount(t, db, memberEmpID, memberUserID)

	tpl := seedOKRTemplateWithOwner(t, db, orgID, OKRTemplateStatusPublished, memberUserID)

	if err := svc.DeleteTemplate(db, tpl.ID, memberUserID); err == nil {
		t.Fatal("expected delete of a published OKR template to be rejected")
	}
}

// seedOKRTemplateWithOwner inserts an OKRTemplate with the given status and
// owner directly, bypassing CreateTemplate's hierarchy guards. created_by_org_id
// diisi = orgID (organisasi pembuat == organisasi template) agar otorisasi
// strict berbasis created_by_org_id tetap berlaku di test.
func seedOKRTemplateWithOwner(t *testing.T, db *gorm.DB, orgID uuid.UUID, status int, ownerID uuid.UUID) *OKRTemplate {
	t.Helper()
	creatorOrg := orgID
	tpl := &OKRTemplate{
		OrganizationID: orgID,
		Name:           "Template",
		Status:         status,
		CreatedBy:      ownerID,
		CreatedByOrgID: &creatorOrg,
	}
	if err := db.Create(tpl).Error; err != nil {
		t.Fatalf("failed to seed OKR template: %v", err)
	}
	return tpl
}
