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
}

func TestServiceActivateRejectsDefaultTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	def := createTestTemplate(db, "DEF1", DocumentTypeContractAgreement, StatusReference, true)

	_, err := svc.Activate(context.Background(), def.ID, "actor-1")
	var immErr *ReferenceTemplateImmutableError
	if !errors.As(err, &immErr) {
		t.Fatalf("expected ReferenceTemplateImmutableError, got %v", err)
	}
}

func TestServiceCreateFromDefaultCopiesContent(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	content := "<p>default content</p>"
	def := createTestTemplate(db, "DEFCOPY", DocumentTypeContractAgreement, StatusReference, true)
	def.Content = &content
	if err := db.Save(def).Error; err != nil {
		t.Fatalf("seed default content: %v", err)
	}

	copied, err := svc.CreateFromDefault(context.Background(), DocumentTypeContractAgreement, "New From Default", "NEWCODE", "actor-1")
	if err != nil {
		t.Fatalf("CreateFromDefault: %v", err)
	}
	if copied.Content == nil || *copied.Content != content {
		t.Fatalf("expected copied content %q, got %v", content, copied.Content)
	}
	if copied.IsDefault {
		t.Fatalf("copied template must not itself be default")
	}
	if copied.Status != StatusInactive {
		t.Fatalf("copied template must start INACTIVE, got %s", copied.Status)
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

	v1, err := svc.CreateVersion(ctx, tpl.ID, "<p>v1</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "actor-1")
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

	v2, err := svc.CreateVersion(ctx, tpl.ID, "<p>v2</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "actor-1")
	if err != nil {
		t.Fatalf("CreateVersion v2: %v", err)
	}
	if v2.Version != 2 {
		t.Fatalf("expected version 2, got %d", v2.Version)
	}
}

func TestServiceDeleteRejectsDefaultTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	def := createTestTemplate(db, "DEFDEL", DocumentTypeMovementSK, StatusReference, true)

	err := svc.Delete(context.Background(), def.ID, "actor-1")
	var immErr *ReferenceTemplateImmutableError
	if !errors.As(err, &immErr) {
		t.Fatalf("expected ReferenceTemplateImmutableError, got %v", err)
	}
}
