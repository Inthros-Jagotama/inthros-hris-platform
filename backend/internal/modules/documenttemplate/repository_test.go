package documenttemplate

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestRepositoryCreateAndGetByID(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()

	tpl := &DocumentTemplate{ID: uuidStr(), Name: "PKWT", Code: "PKWT01", DocumentType: DocumentTypeContractAgreement, Status: StatusInactive}
	if err := repo.Create(ctx, tpl); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := repo.GetByID(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Code != "PKWT01" {
		t.Fatalf("expected code PKWT01, got %s", got.Code)
	}
}

func TestRepositoryGetByIDNotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)

	_, err := repo.GetByID(context.Background(), uuidStr())
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestRepositoryFindActiveByType(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	createTestTemplate(db, "INACTIVE1", DocumentTypeContractAgreement, StatusInactive)
	active := createTestTemplate(db, "ACTIVE1", DocumentTypeContractAgreement, StatusActive)

	got, err := repo.FindActiveByType(context.Background(), DocumentTypeContractAgreement)
	if err != nil {
		t.Fatalf("FindActiveByType: %v", err)
	}
	if got.ID != active.ID {
		t.Fatalf("expected active template %s, got %s", active.ID, got.ID)
	}
}

// TestRepositoryFindActiveByTypeAndMovement: template SK spesifik per jenis
// movement dipilih, dan fallback ke template umum (movement_type='') bila tidak
// ada template spesifik.
func TestRepositoryFindActiveByTypeAndMovement(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()

	generic := createTestTemplateWithMovement(db, "SK_UMUM", DocumentTypeMovementSK, StatusActive, "")
	promotion := createTestTemplateWithMovement(db, "SK_PROMOSI", DocumentTypeMovementSK, StatusActive, "promotion")

	// Jenis yang punya template spesifik → template spesifik.
	got, err := repo.FindActiveByTypeAndMovement(ctx, DocumentTypeMovementSK, "promotion")
	if err != nil {
		t.Fatalf("FindActiveByTypeAndMovement(promotion): %v", err)
	}
	if got.ID != promotion.ID {
		t.Fatalf("expected promotion template %s, got %s", promotion.ID, got.ID)
	}

	// Jenis tanpa template spesifik → fallback template umum.
	got2, err := repo.FindActiveByTypeAndMovement(ctx, DocumentTypeMovementSK, "mutation")
	if err != nil {
		t.Fatalf("FindActiveByTypeAndMovement(mutation): %v", err)
	}
	if got2.ID != generic.ID {
		t.Fatalf("expected generic template %s for mutation, got %s", generic.ID, got2.ID)
	}

	// Tanpa movement type → template umum.
	got3, err := repo.FindActiveByTypeAndMovement(ctx, DocumentTypeMovementSK, "")
	if err != nil {
		t.Fatalf("FindActiveByTypeAndMovement: %v", err)
	}
	if got3.ID != generic.ID {
		t.Fatalf("expected generic template, got %s", got3.ID)
	}

	// Document type non-movement tetap memakai movement_type=''.
	contract := createTestTemplate(db, "PKWT", DocumentTypeContractAgreement, StatusActive)
	got4, err := repo.FindActiveByTypeAndMovement(ctx, DocumentTypeContractAgreement, "")
	if err != nil {
		t.Fatalf("FindActiveByTypeAndMovement(CONTRACT): %v", err)
	}
	if got4.ID != contract.ID {
		t.Fatalf("expected contract template, got %s", got4.ID)
	}
}

func TestRepositoryListPaginationAndSearch(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	for i := 0; i < 5; i++ {
		createTestTemplate(db, uuidStr()[:8], DocumentTypeContractAgreement, StatusInactive)
	}
	createTestTemplate(db, "FINDME", DocumentTypeMovementSK, StatusInactive)

	items, total, err := repo.List(context.Background(), 1, 3, "", "", "", "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 6 || len(items) != 3 {
		t.Fatalf("expected total=6 len=3, got total=%d len=%d", total, len(items))
	}

	filtered, ftotal, err := repo.List(context.Background(), 1, 10, DocumentTypeMovementSK, "", "", "")
	if err != nil {
		t.Fatalf("List filtered: %v", err)
	}
	if ftotal != 1 || filtered[0].Code != "FINDME" {
		t.Fatalf("expected 1 result FINDME, got total=%d", ftotal)
	}

	searched, stotal, err := repo.List(context.Background(), 1, 10, "", "", "", "findme")
	if err != nil {
		t.Fatalf("List search: %v", err)
	}
	if stotal != 1 || searched[0].Code != "FINDME" {
		t.Fatalf("expected search to find FINDME case-insensitively, got total=%d", stotal)
	}
}

func TestRepositorySoftDeleteExcludesFromList(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()
	tpl := createTestTemplate(db, "TOBEDEL", DocumentTypeContractAgreement, StatusInactive)

	if err := repo.SoftDelete(ctx, tpl.ID); err != nil {
		t.Fatalf("SoftDelete: %v", err)
	}
	_, err := repo.GetByID(ctx, tpl.ID)
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound after soft delete, got %v", err)
	}
	items, total, _ := repo.List(ctx, 1, 10, "", "", "", "")
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected soft-deleted row excluded from List, got total=%d", total)
	}
}

func TestRepositoryVersionsCreateListNextNumber(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()
	tpl := createTestTemplate(db, "VERTEST", DocumentTypeContractAgreement, StatusInactive)

	err := repo.WithTx(ctx, func(tx *gorm.DB) error {
		next, nerr := repo.NextVersionNumber(ctx, tx, tpl.ID)
		if nerr != nil {
			return nerr
		}
		if next != 1 {
			t.Fatalf("expected first version number 1, got %d", next)
		}
		v := &DocumentTemplateVersion{ID: uuidStr(), TemplateID: tpl.ID, Version: next, Content: "<p>v1</p>", PaperSize: "A4", Orientation: "portrait"}
		return repo.CreateVersion(ctx, tx, v)
	})
	if err != nil {
		t.Fatalf("WithTx create version 1: %v", err)
	}

	err = repo.WithTx(ctx, func(tx *gorm.DB) error {
		next, nerr := repo.NextVersionNumber(ctx, tx, tpl.ID)
		if nerr != nil {
			return nerr
		}
		if next != 2 {
			t.Fatalf("expected second version number 2, got %d", next)
		}
		v := &DocumentTemplateVersion{ID: uuidStr(), TemplateID: tpl.ID, Version: next, Content: "<p>v2</p>", PaperSize: "A4", Orientation: "portrait"}
		return repo.CreateVersion(ctx, tx, v)
	})
	if err != nil {
		t.Fatalf("WithTx create version 2: %v", err)
	}

	versions, err := repo.ListVersions(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
}
