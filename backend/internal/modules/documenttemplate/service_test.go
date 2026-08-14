package documenttemplate

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestServiceCreateRejectsInvalidDocumentType(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())

	_, err := svc.Create(context.Background(), "X", "CODEX", "NOT_A_TYPE", "", "actor-1")
	var invalidErr *InvalidDocumentTypeError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("expected InvalidDocumentTypeError, got %v", err)
	}
}

func TestServiceCreateAutoGeneratesCodeWhenEmpty(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())
	ctx := context.Background()

	tpl, err := svc.Create(ctx, "Auto Code", "", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("Create with empty code: %v", err)
	}
	if tpl.Code == "" {
		t.Fatal("expected auto-generated code, got empty")
	}
	expectedPrefix := "TMPL-" + DocumentTypeContractAgreement + "-"
	if len(tpl.Code) <= len(expectedPrefix) || tpl.Code[:len(expectedPrefix)] != expectedPrefix {
		t.Fatalf("expected code with prefix %q, got %q", expectedPrefix, tpl.Code)
	}

	// Kode otomatis harus tetap unik untuk dua template berbeda
	second, err := svc.Create(ctx, "Auto Code 2", "", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("second Create with empty code: %v", err)
	}
	if second.Code == tpl.Code {
		t.Fatalf("expected distinct auto-generated codes, both %q", tpl.Code)
	}
}

func TestServiceCreateRejectsDuplicateCode(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())
	ctx := context.Background()

	if _, err := svc.Create(ctx, "First", "DUPCODE", DocumentTypeContractAgreement, "", "actor-1"); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	_, err := svc.Create(ctx, "Second", "DUPCODE", DocumentTypeContractAgreement, "", "actor-1")
	var dupErr *DuplicateCodeError
	if !errors.As(err, &dupErr) {
		t.Fatalf("expected DuplicateCodeError, got %v", err)
	}
}

func TestServiceActivateDeactivatesPreviousActive(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())
	ctx := context.Background()

	first, err := svc.Create(ctx, "First", "ACT1", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := svc.Activate(ctx, first.ID, "actor-1"); err != nil {
		t.Fatalf("activate first: %v", err)
	}

	second, err := svc.Create(ctx, "Second", "ACT2", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := svc.Activate(ctx, second.ID, "actor-1"); err != nil {
		t.Fatalf("activate second: %v", err)
	}

	gotFirst, err := svc.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("get first: %v", err)
	}
	if gotFirst.Status != StatusInactive {
		t.Fatalf("expected first template to be deactivated, got status=%s", gotFirst.Status)
	}
	gotSecond, err := svc.GetByID(ctx, second.ID)
	if err != nil {
		t.Fatalf("get second: %v", err)
	}
	if gotSecond.Status != StatusActive {
		t.Fatalf("expected second template to be active, got status=%s", gotSecond.Status)
	}

	var activeCount int64
	if err := db.Model(&DocumentTemplate{}).
		Where("type = ? AND status = ?", DocumentTypeContractAgreement, StatusActive).
		Count(&activeCount).Error; err != nil {
		t.Fatalf("count active: %v", err)
	}
	if activeCount != 1 {
		t.Fatalf("expected exactly 1 ACTIVE template of type, got %d", activeCount)
	}
}

func TestServiceActivateSelfHealsMultipleActiveRows(t *testing.T) {
	// Simulates a pre-existing invariant violation (e.g. from the race this
	// fix closes, or from data predating the migration): two ACTIVE rows of
	// the same type exist already. Activating a third must deactivate BOTH.
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()

	a := createTestTemplate(db, "MULTIACT1", DocumentTypeContractAgreement, StatusActive)
	b := createTestTemplate(db, "MULTIACT2", DocumentTypeContractAgreement, StatusActive)
	c := createTestTemplate(db, "MULTIACT3", DocumentTypeContractAgreement, StatusInactive)

	if _, err := svc.Activate(ctx, c.ID, "actor-1"); err != nil {
		t.Fatalf("activate c: %v", err)
	}

	var activeCount int64
	if err := db.Model(&DocumentTemplate{}).
		Where("type = ? AND status = ?", DocumentTypeContractAgreement, StatusActive).
		Count(&activeCount).Error; err != nil {
		t.Fatalf("count active: %v", err)
	}
	if activeCount != 1 {
		t.Fatalf("expected self-heal to exactly 1 ACTIVE row, got %d", activeCount)
	}

	gotA, _ := svc.GetByID(ctx, a.ID)
	gotB, _ := svc.GetByID(ctx, b.ID)
	if gotA.Status != StatusInactive || gotB.Status != StatusInactive {
		t.Fatalf("expected both pre-existing active rows deactivated, got a=%s b=%s", gotA.Status, gotB.Status)
	}
}

func TestServiceCreateVersionIncrementsAndSetsActiveVersion(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	tpl, err := svc.Create(ctx, "Vers", "VERSVC", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	v1, err := svc.CreateVersion(ctx, tpl.ID, "<p>v1</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "", "actor-1")
	if err != nil {
		t.Fatalf("CreateVersion v1: %v", err)
	}
	if v1.Version != 1 {
		t.Fatalf("expected version 1, got %d", v1.Version)
	}
	updated, err := svc.GetByID(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("GetByID after version: %v", err)
	}
	if updated.ActiveVersionID == nil || *updated.ActiveVersionID != v1.ID {
		t.Fatalf("expected template.active_version_id to point at v1, got %v", updated.ActiveVersionID)
	}

	v2, err := svc.CreateVersion(ctx, tpl.ID, "<p>v2</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "", "actor-1")
	if err != nil {
		t.Fatalf("CreateVersion v2: %v", err)
	}
	if v2.Version != 2 {
		t.Fatalf("expected version 2, got %d", v2.Version)
	}
}

func TestServiceUpdatePartialDoesNotBlankName(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()

	tpl, err := svc.Create(ctx, "Original Name", "UPD1", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	newDesc := "only description changed"
	updated, err := svc.Update(ctx, tpl.ID, nil, &newDesc, "actor-1")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Original Name" {
		t.Fatalf("expected name preserved, got %q", updated.Name)
	}
	if updated.Description == nil || *updated.Description != newDesc {
		t.Fatalf("expected description updated, got %v", updated.Description)
	}
}

func TestServiceCreateVersionRejectsNonexistentTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())

	_, err := svc.CreateVersion(context.Background(), uuidStr(), "<p>x</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "", "actor-1")
	if !errors.Is(err, ErrTemplateNotFound) {
		t.Fatalf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestServiceDeleteThenRecreateWithSameCode(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()

	tpl, err := svc.Create(ctx, "PKWT", "PKWT01", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := svc.Delete(ctx, tpl.ID, "actor-1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	recreated, err := svc.Create(ctx, "PKWT v2", "PKWT01", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("expected code reuse after soft delete to succeed, got: %v", err)
	}
	if recreated.Code != "PKWT01" {
		t.Fatalf("expected code PKWT01, got %s", recreated.Code)
	}
}

