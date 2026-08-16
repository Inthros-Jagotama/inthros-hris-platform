package payroll

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/inthros/hris-platform/internal/modules/approval"
)

// =============================================================================
// Service Tests (using real SQLite repository)
// =============================================================================

func newTestService() (*Service, func()) {
	_, dbResolver, cleanup := setupTestDB()
	repo := NewRepository(dbResolver)
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	return svc, func() { cleanup(); logger.Sync() }
}

// =============================================================================
// Salary Component Service Tests
// =============================================================================

func TestService_CreateSalaryComponent_Defaults(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req := CreateSalaryComponentRequest{
		Code:          "BONUS",
		Name:          "THR Bonus",
		ComponentType: "EARNING",
	}

	resp, err := svc.CreateSalaryComponent(ctx, req)
	if err != nil {
		t.Fatalf("CreateSalaryComponent failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.IsTaxable != true {
		t.Errorf("expected IsTaxable true, got false")
	}
	if resp.DisplayOrder != 100 {
		t.Errorf("expected DisplayOrder 100, got %d", resp.DisplayOrder)
	}
}

func TestService_CreateSalaryComponent_WithOptionalFields(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	isTaxable := false
	displayOrder := 50
	status := "INACTIVE"
	req := CreateSalaryComponentRequest{
		Code:          "DEDUCT-TEST",
		Name:          "Test Deduction",
		ComponentType: "DEDUCTION",
		IsTaxable:     &isTaxable,
		DisplayOrder:  &displayOrder,
		Status:        &status,
	}

	resp, err := svc.CreateSalaryComponent(ctx, req)
	if err != nil {
		t.Fatalf("CreateSalaryComponent failed: %v", err)
	}

	if resp.IsTaxable != false {
		t.Errorf("expected IsTaxable false, got true")
	}
	if resp.DisplayOrder != 50 {
		t.Errorf("expected DisplayOrder 50, got %d", resp.DisplayOrder)
	}
	if resp.Status != "INACTIVE" {
		t.Errorf("expected status 'INACTIVE', got '%s'", resp.Status)
	}
}

func TestService_GetSalaryComponentByID_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "BASIC-SVC",
		Name:          "Basic Salary Service",
		ComponentType: "EARNING",
	})

	found, err := svc.GetSalaryComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetSalaryComponentByID failed: %v", err)
	}
	if found.Name != "Basic Salary Service" {
		t.Errorf("expected name 'Basic Salary Service', got '%s'", found.Name)
	}
}

func TestService_GetSalaryComponentByID_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetSalaryComponentByID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestService_GetSalaryComponentByID_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.GetSalaryComponentByID(ctx, uuid.New().String())
	if err == nil {
		t.Fatal("expected error for non-existent salary component")
	}
}

func TestService_ListSalaryComponents_Pagination(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
			Code:          fmt.Sprintf("COMP-LST-%03d", i+1),
			Name:          fmt.Sprintf("Component %d", i+1),
			ComponentType: "EARNING",
		})
	}

	resp, err := svc.ListSalaryComponents(ctx, 0, 0)
	if err != nil {
		t.Fatalf("ListSalaryComponents failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Page)
	}
	if resp.PerPage != 20 {
		t.Errorf("expected per_page 20 (default), got %d", resp.PerPage)
	}
	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
	if !resp.Success {
		t.Error("expected Success true")
	}
}

func TestService_UpdateSalaryComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "UPD-TEST",
		Name:          "Before Update",
		ComponentType: "EARNING",
	})

	newName := "After Update"
	updated, err := svc.UpdateSalaryComponent(ctx, created.ID, UpdateSalaryComponentRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateSalaryComponent failed: %v", err)
	}

	if updated.Name != "After Update" {
		t.Errorf("expected name 'After Update', got '%s'", updated.Name)
	}
}

func TestService_DeleteSalaryComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:          "DEL-TEST",
		Name:          "To Delete",
		ComponentType: "EARNING",
	})

	if err := svc.DeleteSalaryComponent(ctx, created.ID); err != nil {
		t.Fatalf("DeleteSalaryComponent failed: %v", err)
	}

	_, err := svc.GetSalaryComponentByID(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after deleting salary component")
	}
}

// =============================================================================
// Payroll Period Service Tests
// =============================================================================

func TestService_CreatePayrollPeriod_GeneratesCode(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear:  2026,
		PeriodMonth: 2,
		StartDate:   "2026-02-01",
		EndDate:     "2026-02-28",
		AsOfDate:    "2026-02-28",
	})
	if err != nil {
		t.Fatalf("CreatePayrollPeriod failed: %v", err)
	}

	if resp.PeriodCode != "202602" {
		t.Errorf("expected period_code '202602', got '%s'", resp.PeriodCode)
	}
	if resp.Status != "OPEN" {
		t.Errorf("expected default status 'OPEN', got '%s'", resp.Status)
	}
}

// =============================================================================
// BPJS Setting Service Tests
// =============================================================================

func TestService_CreateBpjsSetting_Defaults(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode:        "BPJS-TEST",
		SettingName:        "Test BPJS",
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreateBpjsSetting failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.BaseSource != "BPJS_BASE_COMPONENTS" {
		t.Errorf("expected base_source 'BPJS_BASE_COMPONENTS', got '%s'", resp.BaseSource)
	}
	if resp.RoundingMode != "ROUND" {
		t.Errorf("expected rounding_mode 'ROUND', got '%s'", resp.RoundingMode)
	}
}

func TestService_UpdateBpjsSetting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode:        "BPJS-UPD",
		SettingName:        "Before Update",
		EffectiveStartDate: "2026-01-01",
	})

	newName := "After Update"
	updated, err := svc.UpdateBpjsSetting(ctx, created.ID, UpdateBpjsSettingRequest{
		SettingName: &newName,
	})
	if err != nil {
		t.Fatalf("UpdateBpjsSetting failed: %v", err)
	}

	if updated.SettingName != "After Update" {
		t.Errorf("expected 'After Update', got '%s'", updated.SettingName)
	}
}

func TestService_ListBpjsSettings(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-L1", SettingName: "BPJS 1", EffectiveStartDate: "2026-01-01",
	})
	svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-L2", SettingName: "BPJS 2", EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.ListBpjsSettings(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListBpjsSettings failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestService_DeleteBpjsSetting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-DEL", SettingName: "To Delete", EffectiveStartDate: "2026-01-01",
	})

	if err := svc.DeleteBpjsSetting(ctx, created.ID); err != nil {
		t.Fatalf("DeleteBpjsSetting failed: %v", err)
	}
}

// =============================================================================
// BPJS Rate Component Service Tests
// =============================================================================

func TestService_CreateBpjsRateComponent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	setting, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-RC", SettingName: "Rate Test", EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.CreateBpjsRateComponent(ctx, CreateBpjsRateComponentRequest{
		BpjsSettingID:      setting.ID,
		RateCode:           "JHT-EMP",
		RateName:           "JHT Employee",
		BpjsProgram:        "JHT",
		PaidBy:             "EMPLOYEE",
		RatePercent:        2.0,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreateBpjsRateComponent failed: %v", err)
	}

	if resp.RatePercent != 2.0 {
		t.Errorf("expected rate 2.0, got %f", resp.RatePercent)
	}
}

func TestService_GetBpjsRateComponentByID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	setting, _ := svc.CreateBpjsSetting(ctx, CreateBpjsSettingRequest{
		SettingCode: "BPJS-GET", SettingName: "Get Test", EffectiveStartDate: "2026-01-01",
	})

	created, _ := svc.CreateBpjsRateComponent(ctx, CreateBpjsRateComponentRequest{
		BpjsSettingID:      setting.ID,
		RateCode:           "GET-RATE",
		RateName:           "Get Rate",
		BpjsProgram:        "JKK",
		PaidBy:             "EMPLOYER",
		RatePercent:        0.24,
		EffectiveStartDate: "2026-01-01",
	})

	found, err := svc.GetBpjsRateComponentByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetBpjsRateComponentByID failed: %v", err)
	}
	if found.RateCode != "GET-RATE" {
		t.Errorf("expected 'GET-RATE', got '%s'", found.RateCode)
	}
}

// =============================================================================
// PPh21 Settings Service Tests
// =============================================================================

func TestService_CreatePph21Setting(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePph21Setting(ctx, CreatePph21SettingRequest{
		SettingCode:        "PPH-TEST",
		SettingName:        "Test PPh21",
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreatePph21Setting failed: %v", err)
	}

	if resp.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", resp.Status)
	}
	if resp.AnnualizationMonths != 12 {
		t.Errorf("expected annualization_months 12, got %d", resp.AnnualizationMonths)
	}
}

// =============================================================================
// Payroll Run Service Tests
// =============================================================================

func TestService_CreatePayrollRun(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 3,
		StartDate: "2026-03-01", EndDate: "2026-03-31", AsOfDate: "2026-03-31",
	})

	resp, err := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-MAR-2026",
		RunType:         "REGULAR",
	})
	if err != nil {
		t.Fatalf("CreatePayrollRun failed: %v", err)
	}

	if resp.Status != "DRAFT" {
		t.Errorf("expected default status 'DRAFT', got '%s'", resp.Status)
	}
	if resp.RunCode != "RUN-MAR-2026" {
		t.Errorf("expected run_code 'RUN-MAR-2026', got '%s'", resp.RunCode)
	}
}

// mockApprovalEngine is a minimal ApprovalEngine stub for payroll service tests.
type mockApprovalEngine struct {
	instanceID string
	createErr  error
	activeFlow string
}

func (m *mockApprovalEngine) CreateApprovalInstance(ctx context.Context, module, documentID, flowID string) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	return m.instanceID, nil
}

func (m *mockApprovalEngine) GetApprovalInstanceStatus(ctx context.Context, instanceID string) (string, error) {
	return "APPROVED", nil
}

// GetActiveFlowIDForModule returns the configured active flow; empty + nil
// error means "no active flow" (caller continues without approval).
func (m *mockApprovalEngine) GetActiveFlowIDForModule(ctx context.Context, module string) (string, error) {
	return m.activeFlow, nil
}

func TestService_UpdatePayrollRunStatus_NoApprovalEngine(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})

	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
	})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus failed: %v", err)
	}

	// Without an approval engine (and no flow_id), requesting CALCULATED
	// auto-advances the run straight to REVIEWED — there is no approval step.
	if updated.Status != "REVIEWED" {
		t.Errorf("expected status 'REVIEWED' without approval engine, got '%s'", updated.Status)
	}
}

func TestService_UpdatePayrollRunStatus_WithApprovalEngine(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetApprovalEngine(&mockApprovalEngine{instanceID: uuid.New().String()})

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})

	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	flowID := uuid.New().String()
	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
		FlowID: &flowID,
	})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus failed: %v", err)
	}

	// With an approval engine + flow, the run lands on CALCULATED
	// and is linked to a created approval instance.
	if updated.Status != "CALCULATED" {
		t.Errorf("expected status 'CALCULATED' with approval engine, got '%s'", updated.Status)
	}
	if updated.ApprovalInstanceID == "" {
		t.Error("expected approval_instance_id to be set with approval engine")
	}
}

// fakeNotifier captures Notify calls for payroll run outcome notifications.
type fakeNotifier struct {
	calls []fakeNotifyCall
}
type fakeNotifyCall struct {
	recipientUserID uuid.UUID
	notifType       string
	referenceType   string
	referenceID     uuid.UUID
}

func (f *fakeNotifier) Notify(_ context.Context, recipientUserID uuid.UUID, notifType string, _ []string, referenceType string, referenceID uuid.UUID) error {
	f.calls = append(f.calls, fakeNotifyCall{recipientUserID, notifType, referenceType, referenceID})
	return nil
}

// TestService_UpdatePayrollRunStatus_AutoResolvesActiveFlow verifies that a
// CALCULATED run with an approval engine but no explicit flow_id still creates
// an approval instance via GetActiveFlowIDForModule (same auto-resolve pattern
// as KPI/requisitions).
func TestService_UpdatePayrollRunStatus_AutoResolvesActiveFlow(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	activeFlow := uuid.New().String()
	svc.SetApprovalEngine(&mockApprovalEngine{instanceID: uuid.New().String(), activeFlow: activeFlow})

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})
	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
	})
	if err != nil {
		t.Fatalf("UpdatePayrollRunStatus failed: %v", err)
	}
	if updated.Status != "CALCULATED" {
		t.Errorf("expected status 'CALCULATED' with auto-resolved flow, got '%s'", updated.Status)
	}
	if updated.ApprovalInstanceID == "" {
		t.Error("expected approval_instance_id to be set via auto-resolved flow")
	}
}

// TestService_CheckPayrollRunApproval_NotifiesCreator verifies the run creator
// is notified when the approval instance reaches APPROVED.
func TestService_CheckPayrollRunApproval_NotifiesCreator(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	creatorID := uuid.New()
	notifier := &fakeNotifier{}
	svc.SetNotifier(notifier)
	svc.SetApprovalEngine(&mockApprovalEngine{instanceID: uuid.New().String()})

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})
	resp, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})
	runID, _ := uuid.Parse(resp.ID)
	runModel, err := svc.repo.FindPayrollRunByID(ctx, runID)
	if err != nil {
		t.Fatalf("find run: %v", err)
	}
	runModel.CreatedBy = &creatorID
	if err := svc.repo.UpdatePayrollRun(ctx, runModel); err != nil {
		t.Fatalf("update run creator: %v", err)
	}

	// Bring the run to CALCULATED + create an approval instance.
	flowID := uuid.New().String()
	if _, err := svc.UpdatePayrollRunStatus(ctx, runID.String(), UpdatePayrollRunStatusRequest{Status: "CALCULATED", FlowID: &flowID}); err != nil {
		t.Fatalf("UpdatePayrollRunStatus failed: %v", err)
	}
	fresh, _ := svc.repo.FindPayrollRunByID(ctx, runID)

	_, err = svc.CheckPayrollRunApproval(ctx, fresh.ID.String(), fresh.ApprovalInstanceID.String())
	if err != nil {
		t.Fatalf("CheckPayrollRunApproval failed: %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifier.calls))
	}
	call := notifier.calls[0]
	if call.notifType != "PAYROLL_APPROVED" {
		t.Errorf("expected notif type PAYROLL_APPROVED, got %s", call.notifType)
	}
	if call.recipientUserID != creatorID {
		t.Errorf("expected recipient %s, got %s", creatorID, call.recipientUserID)
	}
	if call.referenceType != "payroll_run" {
		t.Errorf("expected reference type 'payroll_run', got %s", call.referenceType)
	}
}

func TestService_ListPayrollRuns(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 5,
		StartDate: "2026-05-01", EndDate: "2026-05-31", AsOfDate: "2026-05-31",
	})

	svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID, RunCode: "RUN-MAY-01",
	})
	svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID, RunCode: "RUN-MAY-02",
	})

	resp, err := svc.ListPayrollRuns(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPayrollRuns failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// =============================================================================
// PPh21 Tax Bracket Service Tests
// =============================================================================

func TestService_CreatePph21TaxBracket(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	resp, err := svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder:       1,
		LowerBound:         0,
		UpperBound:         float64Ptr(60000000),
		RatePercent:        5,
		EffectiveStartDate: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("CreatePph21TaxBracket failed: %v", err)
	}

	if resp.RatePercent != 5 {
		t.Errorf("expected rate 5, got %f", resp.RatePercent)
	}
}

func TestService_ListPph21TaxBrackets(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder: 1, LowerBound: 0, UpperBound: float64Ptr(60000000),
		RatePercent: 5, EffectiveStartDate: "2026-01-01",
	})
	svc.CreatePph21TaxBracket(ctx, CreatePph21TaxBracketRequest{
		BracketOrder: 2, LowerBound: 60000000, UpperBound: float64Ptr(250000000),
		RatePercent: 15, EffectiveStartDate: "2026-01-01",
	})

	resp, err := svc.ListPph21TaxBrackets(ctx, 1, 10)
	if err != nil {
		t.Fatalf("ListPph21TaxBrackets failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

// TestService_UpdatePayrollRunStatus_ApprovalRoutingErrorFailsLoudly guards
// the fail-loudly policy: when the approval engine rejects routing, the
// RoutingError is propagated (instead of silently advancing to REVIEWED) so
// the handler can show a bilingual message.
func TestService_UpdatePayrollRunStatus_ApprovalRoutingErrorFailsLoudly(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetApprovalEngine(&mockApprovalEngine{
		createErr: &approval.RoutingError{
			Key:    "approval.no_assignees",
			Params: []string{"Persetujuan Supervisor"},
		},
	})

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})

	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	flowID := uuid.New().String()
	_, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
		FlowID: &flowID,
	})
	if err == nil {
		t.Fatal("expected error when approval routing fails")
	}
	var re *approval.RoutingError
	if !errors.As(err, &re) {
		t.Fatalf("expected approval.RoutingError, got: %v", err)
	}
}

// =============================================================================
// Formula Engine Service Tests
// =============================================================================

func TestService_CreateSalaryComponent_FormulaTypeValid(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	formula := "BPJS_WAGE * 2%"
	resp, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "JHT-EMP-F",
		Name:            "JHT Employee Formula",
		ComponentType:   "DEDUCTION",
		CalculationType: CalculationTypeFormula,
		Formula:         &formula,
	})
	if err != nil {
		t.Fatalf("CreateSalaryComponent with valid formula: %v", err)
	}
	if resp.CalculationType != "FORMULA" {
		t.Errorf("expected calculation_type FORMULA, got %q", resp.CalculationType)
	}
	if resp.Formula != formula {
		t.Errorf("expected formula %q persisted, got %q", formula, resp.Formula)
	}
}

func TestService_CreateSalaryComponent_FormulaTypeInvalidSyntax(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	formula := "BASIC +"
	_, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "BAD-FORMULA",
		Name:            "Bad Formula",
		ComponentType:   "EARNING",
		CalculationType: CalculationTypeFormula,
		Formula:         &formula,
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for invalid formula syntax, got %v", err)
	}
	if !strings.Contains(ve.Message, "formula") {
		t.Errorf("expected message to mention formula, got %q", ve.Message)
	}
}

func TestService_CreateSalaryComponent_FormulaTypeMissingFormula(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "NO-FORMULA",
		Name:            "No Formula",
		ComponentType:   "EARNING",
		CalculationType: CalculationTypeFormula,
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError when formula missing, got %v", err)
	}
}

func TestService_CreateSalaryComponent_PercentageTypeValid(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	formula := "BASIC * 1%"
	resp, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "PCT-COMP",
		Name:            "Percent Component",
		ComponentType:   "EARNING",
		CalculationType: CalculationTypePercentage,
		Formula:         &formula,
	})
	if err != nil {
		t.Fatalf("CreateSalaryComponent with percentage: %v", err)
	}
	if resp.CalculationType != "PERCENTAGE" {
		t.Errorf("expected calculation_type PERCENTAGE, got %q", resp.CalculationType)
	}
}

func TestService_CreateSalaryComponent_ReferenceTypeValid(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	source, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code: "PERF-RESULT", Name: "Performance Result", ComponentType: "EARNING",
	})

	resp, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:                  "PERF-BONUS",
		Name:                  "Performance Bonus",
		ComponentType:         "EARNING",
		CalculationType:       CalculationTypeReference,
		ReferenceComponentID:  &source.ID,
	})
	if err != nil {
		t.Fatalf("CreateSalaryComponent with reference: %v", err)
	}
	if resp.ReferenceComponentID != source.ID {
		t.Errorf("expected reference_component_id %s, got %s", source.ID, resp.ReferenceComponentID)
	}
}

func TestService_CreateSalaryComponent_ReferenceTypeMissingRef(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "REF-NO-TARGET",
		Name:            "Reference Without Target",
		ComponentType:   "EARNING",
		CalculationType: CalculationTypeReference,
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError when reference missing, got %v", err)
	}
}

func TestService_CreateSalaryComponent_ReferenceTypeNotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	missing := uuid.New().String()
	_, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:                 "REF-GHOST",
		Name:                 "Reference To Ghost",
		ComponentType:        "EARNING",
		CalculationType:      CalculationTypeReference,
		ReferenceComponentID: &missing,
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError when reference target missing, got %v", err)
	}
}

func TestService_CreateSalaryComponent_UnknownCalculationType(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code:            "WEIRD",
		Name:            "Weird Type",
		ComponentType:   "EARNING",
		CalculationType: "WHATEVER",
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for unknown calculation type, got %v", err)
	}
}

func TestService_UpdateSalaryComponent_FormulaValidation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	created, _ := svc.CreateSalaryComponent(ctx, CreateSalaryComponentRequest{
		Code: "UPD-FORMULA", Name: "Update Formula", ComponentType: "EARNING",
	})

	calcType := CalculationTypeFormula
	badFormula := "BASIC /"
	_, err := svc.UpdateSalaryComponent(ctx, created.ID, UpdateSalaryComponentRequest{
		CalculationType: &calcType,
		Formula:         &badFormula,
	})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError on update with bad formula, got %v", err)
	}

	goodFormula := "BASIC + ALLOWANCE"
	updated, err := svc.UpdateSalaryComponent(ctx, created.ID, UpdateSalaryComponentRequest{
		CalculationType: &calcType,
		Formula:         &goodFormula,
	})
	if err != nil {
		t.Fatalf("UpdateSalaryComponent with valid formula: %v", err)
	}
	if updated.Formula != goodFormula {
		t.Errorf("expected formula %q, got %q", goodFormula, updated.Formula)
	}
}

func TestService_ValidateFormula(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	vars, err := svc.ValidateFormula(ctx, "BASIC + OVERTIME_HOURS * OVERTIME_RATE")
	if err != nil {
		t.Fatalf("ValidateFormula: %v", err)
	}
	if len(vars) != 3 {
		t.Errorf("expected 3 referenced variables, got %v", vars)
	}

	if _, err := svc.ValidateFormula(ctx, "BASIC +"); err == nil {
		t.Error("expected error for invalid formula")
	}
}

func TestService_ListFormulaVariables(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	vars := svc.ListFormulaVariables(ctx)
	if len(vars) == 0 {
		t.Fatal("expected non-empty formula variables")
	}
	names := map[string]bool{}
	for _, v := range vars {
		names[v.Name] = true
	}
	for _, expected := range []string{"GROSS", "BPJS_WAGE", "NET_SALARY"} {
		if !names[expected] {
			t.Errorf("expected built-in variable %s in registry", expected)
		}
	}
}

// TestService_UpdatePayrollRunStatus_NonRoutingApprovalErrorStillSwallowed
// keeps the best-effort contract for non-routing approval failures: the run
// advances to REVIEWED without an approval step.
func TestService_UpdatePayrollRunStatus_NonRoutingApprovalErrorStillSwallowed(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetApprovalEngine(&mockApprovalEngine{createErr: errors.New("database connection lost")})

	period, _ := svc.CreatePayrollPeriod(ctx, CreatePayrollPeriodRequest{
		PeriodYear: 2026, PeriodMonth: 4,
		StartDate: "2026-04-01", EndDate: "2026-04-30", AsOfDate: "2026-04-30",
	})

	run, _ := svc.CreatePayrollRun(ctx, CreatePayrollRunRequest{
		PayrollPeriodID: period.ID,
		RunCode:         "RUN-APR-2026",
	})

	flowID := uuid.New().String()
	updated, err := svc.UpdatePayrollRunStatus(ctx, run.ID, UpdatePayrollRunStatusRequest{
		Status: "CALCULATED",
		FlowID: &flowID,
	})
	if err != nil {
		t.Fatalf("should not fail on non-routing approval error: %v", err)
	}
	if updated.Status != "REVIEWED" {
		t.Errorf("expected status REVIEWED without approval, got '%s'", updated.Status)
	}
}
