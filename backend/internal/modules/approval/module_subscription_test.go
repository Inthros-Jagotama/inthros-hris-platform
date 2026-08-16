package approval

import (
	"context"
	"testing"
)

func ctxAsCompany(companyID string) context.Context {
	return context.WithValue(context.Background(), "company_id", companyID)
}

// fakeModuleChecker is a test double for ModuleSubscriptionChecker.
type fakeModuleChecker struct {
	active map[string]bool // moduleSlug -> active
}

func (f *fakeModuleChecker) IsModuleActive(companyID, moduleSlug string) (bool, error) {
	return f.active[moduleSlug], nil
}

func (f *fakeModuleChecker) ListActiveModules(companyID string) ([]string, error) {
	var slugs []string
	for slug, active := range f.active {
		if active {
			slugs = append(slugs, slug)
		}
	}
	return slugs, nil
}

func newFakeModuleChecker(activeModules ...string) *fakeModuleChecker {
	active := make(map[string]bool, len(activeModules))
	for _, m := range activeModules {
		active[m] = true
	}
	return &fakeModuleChecker{active: active}
}

func TestService_CreateFlow_ModuleNotSubscribed_Rejected(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("leave"))

	_, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "payroll", // not subscribed
		Name:   "Payroll Approval",
	})
	if err == nil {
		t.Fatal("expected CreateFlow to reject an unsubscribed module, got nil error")
	}
}

func TestService_CreateFlow_ModuleSubscribed_Allowed(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("leave"))

	resp, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "leave",
		Name:   "Leave Approval",
	})
	if err != nil {
		t.Fatalf("expected CreateFlow to succeed for a subscribed module, got: %v", err)
	}
	if resp.Module != "leave" {
		t.Errorf("expected module 'leave', got '%s'", resp.Module)
	}
}

func TestService_CreateFlow_RecruitmentOffer_AllowedViaRecruitmentSubscription(t *testing.T) {
	// G-3: flow module recruitment_offer dicek terhadap subscription
	// "recruitment" (alias) — sama seperti performance_kpi_target.
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("recruitment"))

	resp, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "recruitment_offer",
		Name:   "Persetujuan Offer",
	})
	if err != nil {
		t.Fatalf("expected CreateFlow to succeed for recruitment_offer when 'recruitment' is subscribed, got: %v", err)
	}
	if resp.Module != "recruitment_offer" {
		t.Errorf("expected module 'recruitment_offer', got '%s'", resp.Module)
	}
}

func TestService_CreateFlow_TrainingRequest_AllowedViaTrainingSubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("training"))

	resp, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "training_request",
		Name:   "Persetujuan Training",
	})
	if err != nil {
		t.Fatalf("expected CreateFlow to succeed for training_request when 'training' is subscribed, got: %v", err)
	}
	if resp.Module != "training_request" {
		t.Errorf("expected module 'training_request', got '%s'", resp.Module)
	}
}

func TestService_CreateFlow_RecruitmentOffer_RejectedWithoutRecruitmentSubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("leave"))

	_, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "recruitment_offer",
		Name:   "Persetujuan Offer",
	})
	if err == nil {
		t.Fatal("expected CreateFlow to reject recruitment_offer when 'recruitment' is not subscribed")
	}
}

func TestService_UpdateFlow_ReactivateUnsubscribedModule_Rejected(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	// Created while unconfigured (no checker yet) so setup itself isn't blocked.
	flow := createTestFlow(repo, "reimbursement")

	svc.SetModuleChecker(newFakeModuleChecker("leave")) // reimbursement NOT active

	_, err := svc.UpdateFlow(ctxAsCompany("company-1"), flow.ID.String(), UpdateFlowRequest{
		IsActive: boolPtr(true),
	})
	if err == nil {
		t.Fatal("expected UpdateFlow to reject reactivating a flow for an unsubscribed module, got nil error")
	}
}

func TestService_ListAvailableModules_Success(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	// employee ter-subscribe tapi TIDAK terintegrasi approval → tidak muncul.
	svc.SetModuleChecker(newFakeModuleChecker("leave", "reimbursement", "employee"))
	svc.RegisterStatusHandler("leave", nil)
	svc.RegisterStatusHandler("reimbursement", nil)

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	if len(modules) != 2 {
		t.Errorf("expected 2 integrated+subscribed modules, got %d: %v", len(modules), modules)
	}
}

func TestService_CreateFlow_KPISubModule_AllowedViaPerformanceSubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("performance"))

	resp, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "performance_kpi_target",
		Name:   "KPI Target Approval",
	})
	if err != nil {
		t.Fatalf("expected CreateFlow to succeed for performance_kpi_target when 'performance' is subscribed, got: %v", err)
	}
	if resp.Module != "performance_kpi_target" {
		t.Errorf("expected module 'performance_kpi_target', got '%s'", resp.Module)
	}

	if _, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "performance_kpi_realization",
		Name:   "KPI Realization Approval",
	}); err != nil {
		t.Fatalf("expected CreateFlow to succeed for performance_kpi_realization when 'performance' is subscribed, got: %v", err)
	}
}

func TestService_CreateFlow_KPISubModule_RejectedWithoutPerformanceSubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("leave"))

	_, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "performance_kpi_target",
		Name:   "KPI Target Approval",
	})
	if err == nil {
		t.Fatal("expected CreateFlow to reject performance_kpi_target when 'performance' is not subscribed")
	}
}

func TestService_ListAvailableModules_IncludesKPISubModules(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("performance", "leave"))
	svc.RegisterStatusHandler("leave", nil)
	// Sub-checkpoint KPI terintegrasi approval (bukan "performance" itu sendiri).
	for _, m := range subscriptionModuleSubslots["performance"] {
		svc.RegisterStatusHandler(m, nil)
	}

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	found := map[string]bool{}
	for _, m := range modules {
		found[m] = true
	}
	for _, want := range []string{"leave", "performance_kpi_target", "performance_kpi_realization", "okr_key_result", "okr_assessment"} {
		if !found[want] {
			t.Errorf("expected available modules to include %q, got %v", want, modules)
		}
	}
	if found["performance"] {
		t.Errorf("expected base module 'performance' NOT to appear (bukan module flow terintegrasi langsung), got %v", modules)
	}
}

// TestService_ListAvailableModules_OnlyIntegratedModules memastikan module
// yang disubscribe tenant tapi tidak terintegrasi dengan Central Approval
// (mis. employee, organization) tidak dimunculkan di module picker.
func TestService_ListAvailableModules_OnlyIntegratedModules(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("employee", "organization", "leave"))
	svc.RegisterStatusHandler("leave", nil)

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	if len(modules) != 1 || modules[0] != "leave" {
		t.Errorf("expected only 'leave' (terintegrasi + disubscribe), got %v", modules)
	}
}

// TestService_ListAvailableModules_TrainingAndRecruitmentMapping memastikan
// subscription "training"/"recruitment" membuka checkpoint flow module yang
// slug-nya berbeda (training_request / recruitment_offer) dan keduanya hanya
// muncul bila handler-nya terdaftar (terintegrasi).
func TestService_ListAvailableModules_TrainingAndRecruitmentMapping(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("training", "recruitment"))
	svc.RegisterStatusHandler("recruitment", nil)
	svc.RegisterStatusHandler("recruitment_offer", nil)
	svc.RegisterStatusHandler("training_request", nil)

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	found := map[string]bool{}
	for _, m := range modules {
		found[m] = true
	}
	for _, want := range []string{"recruitment", "recruitment_offer", "training_request"} {
		if !found[want] {
			t.Errorf("expected available modules to include %q, got %v", want, modules)
		}
	}
	if found["training"] {
		t.Errorf("expected 'training' NOT to appear (bukan module flow terintegrasi langsung), got %v", modules)
	}
	if len(modules) != 3 {
		t.Errorf("expected exactly 3 modules, got %d: %v", len(modules), modules)
	}
}

// TestService_GetActiveFlowByModule_FallsBackToBaseModule validates that
// if HR only configured a flow under the base "performance" module (not
// the specific "performance_kpi_target"/"performance_kpi_realization"
// sub-checkpoint slugs), that flow is used as a fallback instead of the
// KPI submission silently skipping approval routing entirely.
func TestService_GetActiveFlowByModule_FallsBackToBaseModule(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "performance")

	resp, err := svc.GetActiveFlowByModule(ctxAsCompany("company-1"), "performance_kpi_target")
	if err != nil {
		t.Fatalf("expected fallback to the 'performance' flow, got error: %v", err)
	}
	if resp.ID != flow.ID.String() {
		t.Errorf("expected fallback flow id %s, got %s", flow.ID, resp.ID)
	}
}

// TestService_GetActiveFlowByModule_SpecificFlowTakesPriority validates
// that once a dedicated flow exists for the specific sub-checkpoint slug,
// it's used instead of the base module's flow.
func TestService_GetActiveFlowByModule_SpecificFlowTakesPriority(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	createTestFlow(repo, "performance")
	specific := createTestFlow(repo, "performance_kpi_target")

	resp, err := svc.GetActiveFlowByModule(ctxAsCompany("company-1"), "performance_kpi_target")
	if err != nil {
		t.Fatalf("GetActiveFlowByModule failed: %v", err)
	}
	if resp.ID != specific.ID.String() {
		t.Errorf("expected the specific flow %s to take priority, got %s", specific.ID, resp.ID)
	}
}

func TestService_ListAvailableModules_NoChecker_Error(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err == nil {
		t.Fatal("expected error when module checker is not configured, got nil")
	}
}

// TestService_CreateFlow_Competency360_AllowedViaCompetencySubscription validates
// the Phase 4 wiring: flow module "competency_360_assessment" dicek terhadap
// subscription "competency" (alias) — pola persis performance_kpi_target.
func TestService_CreateFlow_Competency360_AllowedViaCompetencySubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("competency"))

	resp, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "competency_360_assessment",
		Name:   "Finalisasi Assessment 360",
	})
	if err != nil {
		t.Fatalf("expected CreateFlow to succeed for competency_360_assessment when 'competency' is subscribed, got: %v", err)
	}
	if resp.Module != "competency_360_assessment" {
		t.Errorf("expected module 'competency_360_assessment', got '%s'", resp.Module)
	}
}

func TestService_CreateFlow_Competency360_RejectedWithoutCompetencySubscription(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("leave"))

	_, err := svc.CreateFlow(ctxAsCompany("company-1"), CreateFlowRequest{
		Module: "competency_360_assessment",
		Name:   "Finalisasi Assessment 360",
	})
	if err == nil {
		t.Fatal("expected CreateFlow to reject competency_360_assessment when 'competency' is not subscribed")
	}
}

// TestService_ListAvailableModules_IncludesCompetency360SubModule memastikan
// subscription "competency" membuka checkpoint flow competency_360_assessment
// di module picker bila handler status-nya terdaftar (terintegrasi) — pola
// persis TestService_ListAvailableModules_IncludesKPISubModules.
func TestService_ListAvailableModules_IncludesCompetency360SubModule(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("competency"))
	svc.RegisterStatusHandler("competency_360_assessment", nil)

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	found := map[string]bool{}
	for _, m := range modules {
		found[m] = true
	}
	if !found["competency_360_assessment"] {
		t.Errorf("expected available modules to include 'competency_360_assessment', got %v", modules)
	}
	if found["competency"] {
		t.Errorf("expected base module 'competency' NOT to appear (bukan module flow terintegrasi langsung), got %v", modules)
	}
}

// TestService_ListAvailableModules_Competency360HiddenWithoutHandler memastikan
// competency_360_assessment TIDAK muncul di picker bila status handler-nya
// belum terdaftar (module belum terintegrasi Central Approval).
func TestService_ListAvailableModules_Competency360HiddenWithoutHandler(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()
	svc.SetModuleChecker(newFakeModuleChecker("competency"))

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	for _, m := range modules {
		if m == "competency_360_assessment" {
			t.Fatalf("expected competency_360_assessment to be hidden without registered status handler, got %v", modules)
		}
	}
}

// TestService_GetActiveFlowByModule_Competency360_FallsBackToBaseModule validates
// that a flow configured under the base "competency" module is used as fallback
// for the competency_360_assessment checkpoint slug.
func TestService_GetActiveFlowByModule_Competency360_FallsBackToBaseModule(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	flow := createTestFlow(repo, "competency")

	resp, err := svc.GetActiveFlowByModule(ctxAsCompany("company-1"), "competency_360_assessment")
	if err != nil {
		t.Fatalf("expected fallback to the 'competency' flow, got error: %v", err)
	}
	if resp.ID != flow.ID.String() {
		t.Errorf("expected fallback flow id %s, got %s", flow.ID, resp.ID)
	}
}
