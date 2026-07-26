package pkgmgr

import (
	"testing"
)

// =========================================================================
// CreatePackage Tests
// =========================================================================

func TestCreatePackage_Success_WithoutModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	req := CreatePackageRequest{
		Name:        "Basic HR",
		Slug:        "basic-hr",
		Description: "Paket dasar HR",
		Price:       1000000,
	}

	resp, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	if resp.Name != req.Name {
		t.Errorf("expected name '%s', got '%s'", req.Name, resp.Name)
	}
	if resp.Slug != req.Slug {
		t.Errorf("expected slug '%s', got '%s'", req.Slug, resp.Slug)
	}
	if resp.Price != req.Price {
		t.Errorf("expected price %f, got %f", req.Price, resp.Price)
	}
	if resp.Status != "draft" {
		t.Errorf("expected status 'draft', got '%s'", resp.Status)
	}
	if resp.ModuleCount != 0 {
		t.Errorf("expected module_count 0, got %d", resp.ModuleCount)
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestCreatePackage_Success_WithModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	req := CreatePackageRequest{
		Name:  "HR Pro",
		Slug:  "hr-pro",
		Price: 5000000,
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), IsMandatory: true, SortOrder: 1},
			{ModuleID: empID.String(), IsMandatory: true, SortOrder: 2},
		},
	}

	resp, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	if resp.ModuleCount != 2 {
		t.Errorf("expected module_count 2, got %d", resp.ModuleCount)
	}
	if len(resp.Modules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(resp.Modules))
	}

	// Verify modules are sorted by sort_order
	if resp.Modules[0].ModuleSlug != "organization" {
		t.Errorf("expected first module 'organization', got '%s'", resp.Modules[0].ModuleSlug)
	}
	if resp.Modules[1].ModuleSlug != "employee" {
		t.Errorf("expected second module 'employee', got '%s'", resp.Modules[1].ModuleSlug)
	}
}

func TestCreatePackage_Success_DependenciesFulfilled(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")
	payID := findModuleUUIDBySlug("payroll")
	leaveID := findModuleUUIDBySlug("leave")

	req := CreatePackageRequest{
		Name:  "HR Complete",
		Slug:  "hr-complete",
		Price: 10000000,
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), IsMandatory: true, SortOrder: 1},
			{ModuleID: empID.String(), IsMandatory: true, SortOrder: 2},
			{ModuleID: payID.String(), IsMandatory: false, SortOrder: 3},
			{ModuleID: leaveID.String(), IsMandatory: false, SortOrder: 4},
		},
	}

	resp, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage with fulfilled dependencies failed: %v", err)
	}

	if resp.ModuleCount != 4 {
		t.Errorf("expected 4 modules, got %d", resp.ModuleCount)
	}
}

func TestCreatePackage_Error_DuplicateSlug(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	req := CreatePackageRequest{
		Name: "First",
		Slug: "same-slug",
	}

	_, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("First CreatePackage failed: %v", err)
	}

	req2 := CreatePackageRequest{
		Name: "Second",
		Slug: "same-slug",
	}

	_, err = svc.CreatePackage(req2)
	if err == nil {
		t.Fatal("expected error for duplicate slug, got nil")
	}
}

func TestCreatePackage_Error_InvalidModuleID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	req := CreatePackageRequest{
		Name: "Invalid",
		Slug: "invalid-module",
		Modules: []PackageModuleInput{
			{ModuleID: "not-a-uuid"},
		},
	}

	_, err := svc.CreatePackage(req)
	if err == nil {
		t.Fatal("expected error for invalid module_id, got nil")
	}
}

func TestCreatePackage_Error_ModuleNotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	req := CreatePackageRequest{
		Name: "Module Not Found",
		Slug: "module-not-found",
		Modules: []PackageModuleInput{
			{ModuleID: uuidStr()},
		},
	}

	_, err := svc.CreatePackage(req)
	if err == nil {
		t.Fatal("expected error for non-existent module, got nil")
	}
}

func TestCreatePackage_Error_MissingDependency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	// "employee" depends on "organization" — only include "employee"
	empID := findModuleUUIDBySlug("employee")

	req := CreatePackageRequest{
		Name:  "Missing Deps",
		Slug:  "missing-deps",
		Price: 3000000,
		Modules: []PackageModuleInput{
			{ModuleID: empID.String(), IsMandatory: true},
		},
	}

	_, err := svc.CreatePackage(req)
	if err == nil {
		t.Fatal("expected error for missing dependency, got nil")
	}
}

func TestCreatePackage_Error_MissingNestedDependency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	// "performance" depends on "organization,employee,job-management,competency"
	// Only include "organization" and "competency" — missing "employee" and "job-management"
	orgID := findModuleUUIDBySlug("organization")
	compID := findModuleUUIDBySlug("competency")
	perfID := findModuleUUIDBySlug("performance")

	req := CreatePackageRequest{
		Name:  "Missing Nested Deps",
		Slug:  "missing-nested-deps",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: compID.String()},
			{ModuleID: perfID.String()},
		},
	}

	_, err := svc.CreatePackage(req)
	if err == nil {
		t.Fatal("expected error for nested missing dependencies (employee, job-management), got nil")
	}
}

func TestCreatePackage_Error_DuplicateModule(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")

	req := CreatePackageRequest{
		Name:  "Duplicate Module",
		Slug:  "duplicate-module",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: orgID.String()},
		},
	}

	_, err := svc.CreatePackage(req)
	if err == nil {
		t.Fatal("expected error for duplicate module, got nil")
	}
}

// =========================================================================
// UpdatePackage Tests
// =========================================================================

func TestUpdatePackage_Success_UpdateName(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	// Create initial package
	createReq := CreatePackageRequest{
		Name:  "Original",
		Slug:  "original-pkg",
		Price: 1000000,
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), SortOrder: 1},
			{ModuleID: empID.String(), SortOrder: 2},
		},
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	// Update package name
	newName := "Updated Name"
	updateReq := UpdatePackageRequest{
		Name: &newName,
	}

	updated, err := svc.UpdatePackage(created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdatePackage failed: %v", err)
	}

	if updated.Name != newName {
		t.Errorf("expected name '%s', got '%s'", newName, updated.Name)
	}
	if updated.ModuleCount != 2 {
		t.Errorf("expected module_count 2 (unchanged), got %d", updated.ModuleCount)
	}
}

func TestUpdatePackage_Success_UpdateModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")
	leaveID := findModuleUUIDBySlug("leave")
	reimbID := findModuleUUIDBySlug("reimbursement")

	// Create initial package with 2 modules
	createReq := CreatePackageRequest{
		Name:  "Growing Pkg",
		Slug:  "growing-pkg",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), SortOrder: 1},
			{ModuleID: empID.String(), SortOrder: 2},
		},
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	// Update: replace with 4 modules (all dependencies fulfilled)
	updateReq := UpdatePackageRequest{
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), SortOrder: 1},
			{ModuleID: empID.String(), SortOrder: 2},
			{ModuleID: leaveID.String(), SortOrder: 3},
			{ModuleID: reimbID.String(), SortOrder: 4},
		},
	}

	updated, err := svc.UpdatePackage(created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdatePackage with new modules failed: %v", err)
	}

	if updated.ModuleCount != 4 {
		t.Errorf("expected 4 modules, got %d", updated.ModuleCount)
	}
}

func TestUpdatePackage_Error_MissingDependency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	// Create initial package
	createReq := CreatePackageRequest{
		Name:  "Base",
		Slug:  "base-pkg",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: empID.String()},
		},
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	// Update: replace modules — "payroll" depends on "employee,organization" but we'll only include "payroll"
	payID := findModuleUUIDBySlug("payroll")
	updateReq := UpdatePackageRequest{
		Modules: []PackageModuleInput{
			{ModuleID: payID.String()},
		},
	}

	_, err = svc.UpdatePackage(created.ID, updateReq)
	if err == nil {
		t.Fatal("expected error for missing dependency on update, got nil")
	}
}

func TestUpdatePackage_Error_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	_, err := svc.UpdatePackage("not-a-uuid", UpdatePackageRequest{})
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
}

func TestUpdatePackage_Error_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	_, err := svc.UpdatePackage(uuidStr(), UpdatePackageRequest{})
	if err == nil {
		t.Fatal("expected error for non-existent package, got nil")
	}
}

// =========================================================================
// PublishPackage Tests
// =========================================================================

func TestPublishPackage_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	// Create package with modules (dependencies fulfilled)
	createReq := CreatePackageRequest{
		Name:  "Ready Package",
		Slug:  "ready-pkg",
		Price: 5000000,
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String(), IsMandatory: true, SortOrder: 1},
			{ModuleID: empID.String(), IsMandatory: true, SortOrder: 2},
		},
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	// Publish
	published, err := svc.PublishPackage(created.ID)
	if err != nil {
		t.Fatalf("PublishPackage failed: %v", err)
	}

	if published.Status != "published" {
		t.Errorf("expected status 'published', got '%s'", published.Status)
	}
	if !published.IsPublic {
		t.Error("expected is_public to be true after publish")
	}

	// Verify it appears in published list
	publishedList, err := svc.ListPublishedPackages("")
	if err != nil {
		t.Fatalf("ListPublishedPackages failed: %v", err)
	}
	if len(publishedList) != 1 {
		t.Errorf("expected 1 published package, got %d", len(publishedList))
	}
}

func TestPublishPackage_Error_NoModules(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	// Create package without modules
	createReq := CreatePackageRequest{
		Name: "Empty Package",
		Slug: "empty-pkg",
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	// Try to publish
	_, err = svc.PublishPackage(created.ID)
	if err == nil {
		t.Fatal("expected error for publishing package without modules, got nil")
	}
}

func TestPublishPackage_Error_MissingDependency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	empID := findModuleUUIDBySlug("employee")
	payID := findModuleUUIDBySlug("payroll")

	// Buat package langsung via repo untuk bypass validasi service
	// (payroll depends_on employee,organization — organization sengaja tidak disertakan)
	pkg := &Package{
		Name:   "Incomplete Package",
		Slug:   "incomplete-pkg",
		Price:  3000000,
		Status: string(PackageDraft),
		Modules: []PackageModule{
			{
				ModuleID:    empID,
				ModuleName:  "Employee Management",
				ModuleSlug:  "employee",
				IsMandatory: false,
				SortOrder:   1,
			},
			{
				ModuleID:    payID,
				ModuleName:  "Payroll Management",
				ModuleSlug:  "payroll",
				IsMandatory: false,
				SortOrder:   2,
			},
		},
	}
	if err := svc.repo.Create(pkg); err != nil {
		t.Fatalf("repo.Create failed: %v", err)
	}

	// Try to publish — should fail because payroll dependency "organization" is missing
	_, err := svc.PublishPackage(pkg.ID.String())
	if err == nil {
		t.Fatal("expected error for publishing package with missing dependency, got nil")
	}
}

func TestUnpublishPackage_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	// Create and publish
	createReq := CreatePackageRequest{
		Name:  "Publish Then Unpublish",
		Slug:  "pub-unpub",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: empID.String()},
		},
	}
	created, err := svc.CreatePackage(createReq)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	published, err := svc.PublishPackage(created.ID)
	if err != nil {
		t.Fatalf("PublishPackage failed: %v", err)
	}

	// Unpublish
	unpublished, err := svc.UnpublishPackage(published.ID)
	if err != nil {
		t.Fatalf("UnpublishPackage failed: %v", err)
	}

	if unpublished.Status != "draft" {
		t.Errorf("expected status 'draft', got '%s'", unpublished.Status)
	}
}

// =========================================================================
// ValidatePackageDependencies Tests
// =========================================================================

func TestValidatePackageDependencies_AllResolved(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")
	payID := findModuleUUIDBySlug("payroll")
	leaveID := findModuleUUIDBySlug("leave")

	req := CreatePackageRequest{
		Name:  "Full Package",
		Slug:  "full-pkg",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: empID.String()},
			{ModuleID: payID.String()},
			{ModuleID: leaveID.String()},
		},
	}
	created, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	deps, err := svc.ValidatePackageDependencies(created.ID)
	if err != nil {
		t.Fatalf("ValidatePackageDependencies failed: %v", err)
	}

	if len(deps) != 4 {
		t.Errorf("expected 4 dependency results, got %d", len(deps))
	}

	for _, d := range deps {
		if !d.Resolved {
			t.Errorf("expected module '%s' to be resolved, but it's not (depends_on: %s)",
				d.ModuleName, d.DependsOn)
		}
	}
}

func TestValidatePackageDependencies_SomeUnresolved(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	perfID := findModuleUUIDBySlug("performance")

	// Buat package langsung via repo untuk bypass validasi service
	// performance depends_on organization,employee,job-management,competency
	// — hanya organization yang disertakan, sisanya tidak
	pkg := &Package{
		Name:   "Partial Package",
		Slug:   "partial-pkg",
		Status: string(PackageDraft),
		Modules: []PackageModule{
			{
				ModuleID:    orgID,
				ModuleName:  "Organization Management",
				ModuleSlug:  "organization",
				IsMandatory: true,
				SortOrder:   1,
			},
			{
				ModuleID:    perfID,
				ModuleName:  "Performance Management",
				ModuleSlug:  "performance",
				IsMandatory: false,
				SortOrder:   2,
			},
		},
	}
	if err := svc.repo.Create(pkg); err != nil {
		t.Fatalf("repo.Create failed: %v", err)
	}

	deps, err := svc.ValidatePackageDependencies(pkg.ID.String())
	if err != nil {
		t.Fatalf("ValidatePackageDependencies failed: %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("expected 2 dependency results, got %d", len(deps))
	}

	// Organization Management should be resolved (no dependencies)
	orgResolved := false
	perfResolved := true // should be false
	for _, d := range deps {
		if d.ModuleName == "Organization Management" && d.Resolved {
			orgResolved = true
		}
		if d.ModuleName == "Performance Management" && !d.Resolved {
			perfResolved = false
		}
	}
	if !orgResolved {
		t.Error("expected 'Organization Management' to be resolved")
	}
	if perfResolved {
		t.Error("expected 'Performance Management' to be unresolved (missing employee, job-management, competency)")
	}
}

// =========================================================================
// GetPackage & ListPackages Tests
// =========================================================================

func TestGetPackage_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")

	req := CreatePackageRequest{
		Name:  "Get Me",
		Slug:  "get-me",
		Price: 2000000,
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
		},
	}
	created, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	fetched, err := svc.GetPackage(created.ID)
	if err != nil {
		t.Fatalf("GetPackage failed: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID '%s', got '%s'", created.ID, fetched.ID)
	}
	if fetched.Name != req.Name {
		t.Errorf("expected name '%s', got '%s'", req.Name, fetched.Name)
	}
}

func TestGetPackage_Error_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetPackage("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
}

func TestGetPackage_Error_NotFound(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	_, err := svc.GetPackage(uuidStr())
	if err == nil {
		t.Fatal("expected error for non-existent package, got nil")
	}
}

func TestListPackages_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	// Create 2 packages
	_, err := svc.CreatePackage(CreatePackageRequest{Name: "Pkg A", Slug: "pkg-a"})
	if err != nil {
		t.Fatalf("CreatePackage A failed: %v", err)
	}
	_, err = svc.CreatePackage(CreatePackageRequest{Name: "Pkg B", Slug: "pkg-b"})
	if err != nil {
		t.Fatalf("CreatePackage B failed: %v", err)
	}

	result, err := svc.ListPackages(1, 10, "")
	if err != nil {
		t.Fatalf("ListPackages failed: %v", err)
	}

	if !result.Success {
		t.Error("expected success to be true")
	}
	if len(result.Data.([]PackageResponse)) != 2 {
		t.Errorf("expected 2 packages, got %d", len(result.Data.([]PackageResponse)))
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
}

func TestListPublishedPackages_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	orgID := findModuleUUIDBySlug("organization")
	empID := findModuleUUIDBySlug("employee")

	// Create and publish one package
	req := CreatePackageRequest{
		Name:  "Published Pkg",
		Slug:  "published-pkg",
		Modules: []PackageModuleInput{
			{ModuleID: orgID.String()},
			{ModuleID: empID.String()},
		},
	}
	created, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}
	_, err = svc.PublishPackage(created.ID)
	if err != nil {
		t.Fatalf("PublishPackage failed: %v", err)
	}

	// Create another draft package (should not appear in published list)
	_, err = svc.CreatePackage(CreatePackageRequest{Name: "Draft Pkg", Slug: "draft-pkg"})
	if err != nil {
		t.Fatalf("CreatePackage draft failed: %v", err)
	}

	published, err := svc.ListPublishedPackages("")
	if err != nil {
		t.Fatalf("ListPublishedPackages failed: %v", err)
	}

	if len(published) != 1 {
		t.Errorf("expected 1 published package, got %d", len(published))
	}
}

// =========================================================================
// DeletePackage Tests
// =========================================================================

func TestDeletePackage_Success(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	req := CreatePackageRequest{
		Name: "Delete Me",
		Slug: "delete-me",
	}
	created, err := svc.CreatePackage(req)
	if err != nil {
		t.Fatalf("CreatePackage failed: %v", err)
	}

	if err := svc.DeletePackage(created.ID); err != nil {
		t.Fatalf("DeletePackage failed: %v", err)
	}

	// Verify it's soft-deleted
	_, err = svc.GetPackage(created.ID)
	if err == nil {
		t.Fatal("expected error for deleted package, got nil")
	}
}

func TestDeletePackage_Error_InvalidUUID(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()

	err := svc.DeletePackage("not-a-uuid")
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}
}
