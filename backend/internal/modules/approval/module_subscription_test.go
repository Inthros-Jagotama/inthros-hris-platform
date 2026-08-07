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
	svc.SetModuleChecker(newFakeModuleChecker("leave", "reimbursement"))

	modules, err := svc.ListAvailableModules(ctxAsCompany("company-1"))
	if err != nil {
		t.Fatalf("ListAvailableModules failed: %v", err)
	}
	if len(modules) != 2 {
		t.Errorf("expected 2 active modules, got %d: %v", len(modules), modules)
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
