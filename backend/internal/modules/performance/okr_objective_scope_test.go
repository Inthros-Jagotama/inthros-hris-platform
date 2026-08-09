package performance

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// TestOKRService_TopOfHierarchy_CanCreateForDirectSubordinate verifies that
// an employee at the top of the org hierarchy (no parent Organization) can
// create an Objective for a direct, occupied subordinate without needing an
// Objective of their own first.
func TestOKRService_TopOfHierarchy_CanCreateForDirectSubordinate(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	subOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, subOrgID, &rootOrgID, "Subordinate Org")

	rootEmployeeID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmployeeID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmployeeID, rootUserID)

	subEmployeeID := uuid.New()
	seedOKREmployment(t, db, subEmployeeID, subOrgID)

	scope, err := svc.GetObjectiveCreationScope(db, rootUserID)
	if err != nil {
		t.Fatalf("GetObjectiveCreationScope failed: %v", err)
	}
	if !scope.Eligible {
		t.Fatalf("expected top-of-hierarchy to be eligible, got ineligible: %s", scope.IneligibleReasonKey)
	}
	if len(scope.SubordinateOrganizations) != 1 || scope.SubordinateOrganizations[0].ID != subOrgID.String() {
		t.Fatalf("expected exactly the direct subordinate, got %+v", scope.SubordinateOrganizations)
	}

	if _, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: subOrgID.String(),
		Name:           "Subordinate Objective",
	}); err != nil {
		t.Fatalf("expected CreateTemplate to succeed for a direct subordinate, got: %v", err)
	}
}

// TestOKRService_CannotCreateForOwnOrganization verifies the top-of-hierarchy
// employee cannot create an Objective for themselves — only for subordinates.
func TestOKRService_CannotCreateForOwnOrganization(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")

	rootEmployeeID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmployeeID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmployeeID, rootUserID)

	if _, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: rootOrgID.String(),
		Name:           "Self Objective",
	}); err == nil {
		t.Fatal("expected CreateTemplate to fail when target organization is the caller's own organization")
	}
}

// TestOKRService_CannotCreateForNonSubordinate verifies an employee cannot
// create an Objective for an Organization outside their effective
// subordinate set (e.g. a sibling branch of the hierarchy).
func TestOKRService_CannotCreateForNonSubordinate(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	unrelatedOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, unrelatedOrgID, nil, "Unrelated Org") // separate hierarchy root

	rootEmployeeID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmployeeID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmployeeID, rootUserID)

	if _, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: unrelatedOrgID.String(),
		Name:           "Cross-hierarchy Objective",
	}); err == nil {
		t.Fatal("expected CreateTemplate to fail for an organization outside the caller's effective subordinate set")
	}
}

// TestOKRService_MiddleEmployee_RequiresOwnObjectiveFirst verifies that an
// employee who is NOT at the top of the hierarchy must already have their
// own Objective before they're allowed to create one for a subordinate.
func TestOKRService_MiddleEmployee_RequiresOwnObjectiveFirst(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	middleOrgID := uuid.New()
	leafOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, middleOrgID, &rootOrgID, "Middle Org")
	seedOKROrganization(t, db, leafOrgID, &middleOrgID, "Leaf Org")

	middleEmployeeID := uuid.New()
	middleUserID := uuid.New()
	seedOKREmployment(t, db, middleEmployeeID, middleOrgID)
	seedOKREmployeeAccount(t, db, middleEmployeeID, middleUserID)

	leafEmployeeID := uuid.New()
	seedOKREmployment(t, db, leafEmployeeID, leafOrgID)

	// Middle org has no Objective of its own yet — must be rejected.
	scope, err := svc.GetObjectiveCreationScope(db, middleUserID)
	if err != nil {
		t.Fatalf("GetObjectiveCreationScope failed: %v", err)
	}
	if scope.Eligible {
		t.Fatal("expected middle-of-hierarchy employee without their own Objective to be ineligible")
	}
	if _, err := svc.CreateTemplate(context.Background(), db, middleUserID, &CreateOKRTemplateRequest{
		OrganizationID: leafOrgID.String(),
		Name:           "Leaf Objective",
	}); err == nil {
		t.Fatal("expected CreateTemplate to fail before the middle org has received its own Objective")
	}

	// Give the middle org its own Objective (as if the root had created it).
	rootEmployeeID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmployeeID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmployeeID, rootUserID)
	tmpl, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: middleOrgID.String(),
		Name:           "Middle Objective",
		Status:         intPtr(1),
	})
	if err != nil {
		t.Fatalf("failed to create middle org's own template: %v", err)
	}
	if _, err := svc.CreateObjective(db, &CreateOKRObjectiveRequest{
		TemplateID: tmpl.ID,
		Title:      "Grow the team",
		Weight:     100,
	}); err != nil {
		t.Fatalf("failed to create middle org's own objective: %v", err)
	}

	// Now middle org is eligible and can create one for its own subordinate.
	scope, err = svc.GetObjectiveCreationScope(db, middleUserID)
	if err != nil {
		t.Fatalf("GetObjectiveCreationScope failed: %v", err)
	}
	if !scope.Eligible {
		t.Fatalf("expected middle-of-hierarchy employee to become eligible after receiving their own Objective, reason: %s", scope.IneligibleReasonKey)
	}
	if _, err := svc.CreateTemplate(context.Background(), db, middleUserID, &CreateOKRTemplateRequest{
		OrganizationID: leafOrgID.String(),
		Name:           "Leaf Objective",
	}); err != nil {
		t.Fatalf("expected CreateTemplate to succeed once middle org has its own Objective, got: %v", err)
	}
}

// TestOKRService_WalkDownSkipsVacantDirectChild verifies that when the
// direct-child Organization is vacant (no active employment), the effective
// subordinate resolution walks further down to the next occupied
// Organization instead of stopping.
func TestOKRService_WalkDownSkipsVacantDirectChild(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	rootOrgID := uuid.New()
	vacantOrgID := uuid.New()
	occupiedGrandchildOrgID := uuid.New()
	seedOKROrganization(t, db, rootOrgID, nil, "Root Org")
	seedOKROrganization(t, db, vacantOrgID, &rootOrgID, "Vacant Org")
	seedOKROrganization(t, db, occupiedGrandchildOrgID, &vacantOrgID, "Occupied Grandchild Org")
	// vacantOrgID has no employment row — deliberately vacant.

	rootEmployeeID := uuid.New()
	rootUserID := uuid.New()
	seedOKREmployment(t, db, rootEmployeeID, rootOrgID)
	seedOKREmployeeAccount(t, db, rootEmployeeID, rootUserID)

	grandchildEmployeeID := uuid.New()
	seedOKREmployment(t, db, grandchildEmployeeID, occupiedGrandchildOrgID)

	scope, err := svc.GetObjectiveCreationScope(db, rootUserID)
	if err != nil {
		t.Fatalf("GetObjectiveCreationScope failed: %v", err)
	}
	if len(scope.SubordinateOrganizations) != 1 || scope.SubordinateOrganizations[0].ID != occupiedGrandchildOrgID.String() {
		t.Fatalf("expected walk-down to reach the occupied grandchild, got %+v", scope.SubordinateOrganizations)
	}

	if _, err := svc.CreateTemplate(context.Background(), db, rootUserID, &CreateOKRTemplateRequest{
		OrganizationID: occupiedGrandchildOrgID.String(),
		Name:           "Grandchild Objective",
	}); err != nil {
		t.Fatalf("expected CreateTemplate to succeed for the occupied grandchild reached by walking down, got: %v", err)
	}
}

// TestOKRService_NoPosition_Ineligible verifies a user with no active
// employment record is reported ineligible with the correct reason.
func TestOKRService_NoPosition_Ineligible(t *testing.T) {
	svc, db, cleanup := setupOKRTestDB(t)
	defer cleanup()

	scope, err := svc.GetObjectiveCreationScope(db, uuid.New())
	if err != nil {
		t.Fatalf("GetObjectiveCreationScope failed: %v", err)
	}
	if scope.Eligible {
		t.Fatal("expected a user with no active position to be ineligible")
	}
	if scope.IneligibleReasonKey != "okr.objective_scope_ineligible_no_position" {
		t.Fatalf("expected no_position reason, got %q", scope.IneligibleReasonKey)
	}
}
