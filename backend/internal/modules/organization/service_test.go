package organization

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TestUpdate_PropagatesFullCodeToDescendants memastikan ketika kode sebuah
// organisasi diubah, full_code seluruh descendants ikut ter-update
// (full_code = chain kode dari root).
func TestUpdate_PropagatesFullCodeToDescendants(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	// Root: code "01" → full_code "01"
	root, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "01",
		Nomenclature:          "Root",
	})
	if err != nil {
		t.Fatalf("create root failed: %v", err)
	}
	if root.FullCode != "01" {
		t.Fatalf("expected root full_code '01', got %q", root.FullCode)
	}

	// Child: code "02" → full_code "0102"
	child, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "02",
		Nomenclature:          "Child",
		ParentID:              strPtrHelper(root.ID),
	})
	if err != nil {
		t.Fatalf("create child failed: %v", err)
	}
	if child.FullCode != "0102" {
		t.Fatalf("expected child full_code '0102', got %q", child.FullCode)
	}

	// Grandchild: code "03" → full_code "010203"
	grand, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "03",
		Nomenclature:          "Grandchild",
		ParentID:              strPtrHelper(child.ID),
	})
	if err != nil {
		t.Fatalf("create grandchild failed: %v", err)
	}
	if grand.FullCode != "010203" {
		t.Fatalf("expected grandchild full_code '010203', got %q", grand.FullCode)
	}

	// ── Update kode root "01" → "09" ──
	newCode := "09"
	updated, err := svc.Update(ctx, root.ID, UpdateOrganizationRequest{Code: &newCode})
	if err != nil {
		t.Fatalf("update root failed: %v", err)
	}
	if updated.FullCode != "09" {
		t.Fatalf("expected updated root full_code '09', got %q", updated.FullCode)
	}

	// ── Verifikasi descendants ikut ter-update ──
	childDB, err := repo.FindByID(ctx, uuid.MustParse(child.ID))
	if err != nil {
		t.Fatalf("fetch child failed: %v", err)
	}
	if childDB.FullCode != "0902" {
		t.Fatalf("expected child full_code '0902' after update, got %q", childDB.FullCode)
	}
	// Parent id anak tidak berubah (hanya kode root yang diubah)
	if childDB.ParentID == nil || childDB.ParentID.String() != root.ID {
		t.Fatalf("expected child parent_id unchanged (root), got %v", childDB.ParentID)
	}

	grandDB, err := repo.FindByID(ctx, uuid.MustParse(grand.ID))
	if err != nil {
		t.Fatalf("fetch grandchild failed: %v", err)
	}
	if grandDB.FullCode != "090203" {
		t.Fatalf("expected grandchild full_code '090203' after update, got %q", grandDB.FullCode)
	}
}

// TestUpdate_MoveToNewParent memastikan ketika sebuah organisasi dipindahkan ke
// parent lain, full_code & level-nya mengikuti parent baru dan seluruh descendants
// ikut ter-update (full_code prefix + level).
func TestUpdate_MoveToNewParent(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	// Root A: "A" → "A"
	rootA, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root A",
	})
	if err != nil {
		t.Fatalf("create root A failed: %v", err)
	}
	// Root B: "B" → "B"
	rootB, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "B",
		Nomenclature:          "Root B",
	})
	if err != nil {
		t.Fatalf("create root B failed: %v", err)
	}
	// Child of A: code "1" → "A1" level 1
	child, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "Child of A",
		ParentID:              strPtrHelper(rootA.ID),
	})
	if err != nil {
		t.Fatalf("create child failed: %v", err)
	}
	// Grandchild of A: code "2" → "A12" level 2
	grand, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "2",
		Nomenclature:          "Grandchild",
		ParentID:              strPtrHelper(child.ID),
	})
	if err != nil {
		t.Fatalf("create grandchild failed: %v", err)
	}

	// ── Pindahkan child dari Root A ke Root B ──
	updated, err := svc.Update(ctx, child.ID, UpdateOrganizationRequest{ParentID: strPtrHelper(rootB.ID)})
	if err != nil {
		t.Fatalf("move child failed: %v", err)
	}
	if updated.FullCode != "B1" {
		t.Fatalf("expected child full_code 'B1' after move, got %q", updated.FullCode)
	}
	if updated.Level != 1 {
		t.Fatalf("expected child level 1 after move, got %d", updated.Level)
	}

	// Parent id child ikut ter-update (bug: sebelumnya ParentID tidak berubah)
	childDB, err := repo.FindByID(ctx, uuid.MustParse(child.ID))
	if err != nil {
		t.Fatalf("fetch child failed: %v", err)
	}
	if childDB.ParentID == nil || childDB.ParentID.String() != rootB.ID {
		t.Fatalf("expected child parent_id to be root B after move, got %v", childDB.ParentID)
	}

	// Descendants ikut ter-update
	grandDB, err := repo.FindByID(ctx, uuid.MustParse(grand.ID))
	if err != nil {
		t.Fatalf("fetch grandchild failed: %v", err)
	}
	if grandDB.FullCode != "B12" {
		t.Fatalf("expected grandchild full_code 'B12' after move, got %q", grandDB.FullCode)
	}
	if grandDB.Level != 2 {
		t.Fatalf("expected grandchild level 2 after move, got %d", grandDB.Level)
	}
}

// TestUpdate_MoveToRoot memastikan pindah ke root (parent dihapus) mengubah
// full_code menjadi kode-nya sendiri dan level 0, serta descendants ikut ter-update.
func TestUpdate_MoveToRoot(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	root, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root",
	})
	if err != nil {
		t.Fatalf("create root failed: %v", err)
	}
	child, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "Child",
		ParentID:              strPtrHelper(root.ID),
	})
	if err != nil {
		t.Fatalf("create child failed: %v", err)
	}
	grand, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "2",
		Nomenclature:          "Grandchild",
		ParentID:              strPtrHelper(child.ID),
	})
	if err != nil {
		t.Fatalf("create grandchild failed: %v", err)
	}

	// ── Pindahkan child ke root (parent_id = "") ──
	updated, err := svc.Update(ctx, child.ID, UpdateOrganizationRequest{ParentID: strPtrHelper("")})
	if err != nil {
		t.Fatalf("move to root failed: %v", err)
	}
	if updated.FullCode != "1" {
		t.Fatalf("expected child full_code '1' after move to root, got %q", updated.FullCode)
	}
	if updated.Level != 0 {
		t.Fatalf("expected child level 0 after move to root, got %d", updated.Level)
	}

	// Parent id harus nil (pindah ke root)
	childDB, err := repo.FindByID(ctx, uuid.MustParse(child.ID))
	if err != nil {
		t.Fatalf("fetch child failed: %v", err)
	}
	if childDB.ParentID != nil {
		t.Fatalf("expected child parent_id nil after move to root, got %v", childDB.ParentID)
	}

	grandDB, err := repo.FindByID(ctx, uuid.MustParse(grand.ID))
	if err != nil {
		t.Fatalf("fetch grandchild failed: %v", err)
	}
	if grandDB.FullCode != "12" {
		t.Fatalf("expected grandchild full_code '12' after move to root, got %q", grandDB.FullCode)
	}
	if grandDB.Level != 1 {
		t.Fatalf("expected grandchild level 1 after move to root, got %d", grandDB.Level)
	}
}

// TestUpdate_MoveDeeper_LevelDelta memastikan saat organisasi dipindahkan ke parent
// yang lebih dalam, level seluruh descendants bergeser dengan delta yang sama.
func TestUpdate_MoveDeeper_LevelDelta(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	rootA, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root A",
	})
	if err != nil {
		t.Fatalf("create root A failed: %v", err)
	}
	rootB, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "B",
		Nomenclature:          "Root B",
	})
	if err != nil {
		t.Fatalf("create root B failed: %v", err)
	}
	// Target parent: C = "B1" level 1, dengan anak D = "B12" level 2
	c, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "C",
		ParentID:              strPtrHelper(rootB.ID),
	})
	if err != nil {
		t.Fatalf("create C failed: %v", err)
	}
	if _, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "2",
		Nomenclature:          "D",
		ParentID:              strPtrHelper(c.ID),
	}); err != nil {
		t.Fatalf("create D failed: %v", err)
	}
	// X = "A1" level 1, Y = "A12" level 2 (di bawah A)
	x, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "X",
		ParentID:              strPtrHelper(rootA.ID),
	})
	if err != nil {
		t.Fatalf("create X failed: %v", err)
	}
	y, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "2",
		Nomenclature:          "Y",
		ParentID:              strPtrHelper(x.ID),
	})
	if err != nil {
		t.Fatalf("create Y failed: %v", err)
	}

	// ── Pindahkan X (level 1) ke bawah C (level 1) → X jadi level 2, Y jadi level 3 ──
	updated, err := svc.Update(ctx, x.ID, UpdateOrganizationRequest{ParentID: strPtrHelper(c.ID)})
	if err != nil {
		t.Fatalf("move X deeper failed: %v", err)
	}
	if updated.FullCode != "B11" {
		t.Fatalf("expected X full_code 'B11' after deeper move, got %q", updated.FullCode)
	}
	if updated.Level != 2 {
		t.Fatalf("expected X level 2 after deeper move, got %d", updated.Level)
	}

	// Parent id X harus C setelah pindah
	xDB, err := repo.FindByID(ctx, uuid.MustParse(x.ID))
	if err != nil {
		t.Fatalf("fetch X failed: %v", err)
	}
	if xDB.ParentID == nil || xDB.ParentID.String() != c.ID {
		t.Fatalf("expected X parent_id to be C after deeper move, got %v", xDB.ParentID)
	}

	yDB, err := repo.FindByID(ctx, uuid.MustParse(y.ID))
	if err != nil {
		t.Fatalf("fetch Y failed: %v", err)
	}
	if yDB.FullCode != "B112" {
		t.Fatalf("expected Y full_code 'B112' after deeper move, got %q", yDB.FullCode)
	}
	if yDB.Level != 3 {
		t.Fatalf("expected Y level 3 after deeper move, got %d", yDB.Level)
	}
}

// TestUpdate_RejectDescendantCollision memastikan perpindahan ditolak bila
// full_code baru salah satu descendant bentrok dengan organisasi lain di
// summary yang sama (data tidak konsisten — defense-in-depth).
func TestUpdate_RejectDescendantCollision(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()
	summaryUUID := uuid.MustParse(summaryID)

	rootA, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root A",
	})
	if err != nil {
		t.Fatalf("create root A failed: %v", err)
	}
	rootB, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "B",
		Nomenclature:          "Root B",
	})
	if err != nil {
		t.Fatalf("create root B failed: %v", err)
	}
	// X = "A1" (level 1) dengan anak Y = "A1X" (level 2)
	x, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "X",
		ParentID:              strPtrHelper(rootA.ID),
	})
	if err != nil {
		t.Fatalf("create X failed: %v", err)
	}
	if _, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "X",
		Nomenclature:          "Y",
		ParentID:              strPtrHelper(x.ID),
	}); err != nil {
		t.Fatalf("create Y failed: %v", err)
	}

	// Sisipkan org tak konsisten langsung via repo: full_code "B1X" (tanpa ancestor "B1").
	fakeID := uuid.New()
	if err := repo.Create(ctx, &Organization{
		ID:                    fakeID,
		Code:                  "1X",
		FullCode:              "B1X",
		Level:                 2,
		Nomenclature:          "Fake",
		OrganizationSummaryID: &summaryUUID,
	}); err != nil {
		t.Fatalf("insert fake org failed: %v", err)
	}

	// Pindahkan X ke bawah B: X baru "B1", anak Y baru "B1X" → bentrok dengan fake.
	if _, err := svc.Update(ctx, x.ID, UpdateOrganizationRequest{ParentID: strPtrHelper(rootB.ID)}); err == nil {
		t.Fatalf("expected descendant collision to be rejected, got nil error")
	}
}

// TestUpdate_RejectCycleMove memastikan organisasi tidak bisa dipindahkan ke
// bawah salah satu descendants-nya sendiri (anti-cycle).
func TestUpdate_RejectCycleMove(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	root, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root",
	})
	if err != nil {
		t.Fatalf("create root failed: %v", err)
	}
	child, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "1",
		Nomenclature:          "Child",
		ParentID:              strPtrHelper(root.ID),
	})
	if err != nil {
		t.Fatalf("create child failed: %v", err)
	}
	grand, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "2",
		Nomenclature:          "Grandchild",
		ParentID:              strPtrHelper(child.ID),
	})
	if err != nil {
		t.Fatalf("create grandchild failed: %v", err)
	}

	// ── Coba pindahkan root ke bawah grandchild → harus ditolak ──
	if _, err := svc.Update(ctx, root.ID, UpdateOrganizationRequest{ParentID: strPtrHelper(grand.ID)}); err == nil {
		t.Fatalf("expected cycle move to be rejected, got nil error")
	}
	// ── Coba pindahkan child ke bawah grandchild-nya sendiri → harus ditolak ──
	if _, err := svc.Update(ctx, child.ID, UpdateOrganizationRequest{ParentID: strPtrHelper(grand.ID)}); err == nil {
		t.Fatalf("expected cycle move to be rejected, got nil error")
	}
}

// TestGetTree_AllLevels memastikan tree yang dikembalikan GetTree tidak dibatasi
// kedalamannya — rantai 5 level harus muncul lengkap (sebelumnya dibatasi 3 level
// oleh preload GORM di FindTree).
func TestGetTree_AllLevels(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	// Bangun rantai 5 level (semua kode "1"): full_code 1, 11, 111, 1111, 11111
	var parentID *string
	lastID := ""
	for i := 0; i < 5; i++ {
		req := CreateOrganizationRequest{
			OrganizationSummaryID: summaryID,
			Code:                  "1",
			Nomenclature:          fmt.Sprintf("Level %d", i+1),
		}
		if parentID != nil {
			req.ParentID = parentID
		}
		org, err := svc.Create(ctx, req)
		if err != nil {
			t.Fatalf("create level %d failed: %v", i+1, err)
		}
		lastID = org.ID
		p := org.ID
		parentID = &p
	}

	tree, err := svc.GetTree(ctx, summaryID)
	if err != nil {
		t.Fatalf("GetTree failed: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("expected 1 root, got %d", len(tree))
	}

	// Telusuri rantai hingga kedalaman 5
	maxDepth := 0
	var walk func(resp OrganizationResponse, d int)
	walk = func(resp OrganizationResponse, d int) {
		if d > maxDepth {
			maxDepth = d
		}
		for _, c := range resp.Children {
			walk(c, d+1)
		}
	}
	walk(tree[0], 1)
	if maxDepth != 5 {
		t.Fatalf("expected tree depth 5 (all levels), got %d", maxDepth)
	}

	// Verifikasi full_code leaf = "11111"
	leaf := tree[0]
	for leaf.Children != nil {
		if len(leaf.Children) == 0 {
			break
		}
		leaf = leaf.Children[0]
	}
	if leaf.FullCode != "11111" {
		t.Fatalf("expected leaf full_code '11111', got %q", leaf.FullCode)
	}
	if leaf.ID != lastID {
		t.Fatalf("expected leaf to be the last created org")
	}
}

// TestGetTree_SiblingOrder memastikan urutan sibling di dalam tree mengikuti
// sort_order (lalu full_code). FindTree menempelkan anak secara reverse-append
// (urutan sibling terbalik) lalu mengandalkan sortChildrenRecursive untuk
// mengembalikan urutan yang benar — test ini mengunci perilaku tersebut.
func TestGetTree_SiblingOrder(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()
	summaryUUID := uuid.MustParse(summaryID)

	root, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "A",
		Nomenclature:          "Root",
	})
	if err != nil {
		t.Fatalf("create root failed: %v", err)
	}
	rootUUID := uuid.MustParse(root.ID)

	// Sisipkan 3 anak langsung via repo dengan sort_order berbeda (terbalik).
	children := []Organization{
		{ID: uuid.New(), Code: "1", FullCode: "A1", Nomenclature: "Sort3", Level: 1, SortOrder: 3, OrganizationSummaryID: &summaryUUID, ParentID: &rootUUID},
		{ID: uuid.New(), Code: "2", FullCode: "A2", Nomenclature: "Sort1", Level: 1, SortOrder: 1, OrganizationSummaryID: &summaryUUID, ParentID: &rootUUID},
		{ID: uuid.New(), Code: "3", FullCode: "A3", Nomenclature: "Sort2", Level: 1, SortOrder: 2, OrganizationSummaryID: &summaryUUID, ParentID: &rootUUID},
	}
	for i := range children {
		if err := repo.Create(ctx, &children[i]); err != nil {
			t.Fatalf("create child %d failed: %v", i, err)
		}
	}

	tree, err := svc.GetTree(ctx, summaryID)
	if err != nil {
		t.Fatalf("GetTree failed: %v", err)
	}
	if len(tree) != 1 {
		t.Fatalf("expected 1 root, got %d", len(tree))
	}
	got := tree[0].Children
	if len(got) != 3 {
		t.Fatalf("expected 3 children, got %d", len(got))
	}
	expected := []string{"Sort1", "Sort2", "Sort3"}
	for i, e := range expected {
		if got[i].Nomenclature != e {
			t.Fatalf("expected child[%d] = %s, got %s", i, e, got[i].Nomenclature)
		}
	}
}

// TestUpdate_SameCode_NoFullCodeChange memastikan update tanpa perubahan code
// tidak mengubah full_code descendants.
func TestUpdate_SameCode_NoFullCodeChange(t *testing.T) {
	repo := newTestRepo(t)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	summaryID := uuid.New().String()

	root, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "01",
		Nomenclature:          "Root",
	})
	if err != nil {
		t.Fatalf("create root failed: %v", err)
	}
	child, err := svc.Create(ctx, CreateOrganizationRequest{
		OrganizationSummaryID: summaryID,
		Code:                  "02",
		Nomenclature:          "Child",
		ParentID:              strPtrHelper(root.ID),
	})
	if err != nil {
		t.Fatalf("create child failed: %v", err)
	}

	// Update hanya nomenclature (code tetap)
	newName := "Root (updated)"
	if _, err := svc.Update(ctx, root.ID, UpdateOrganizationRequest{Nomenclature: &newName}); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	childDB, err := repo.FindByID(ctx, uuid.MustParse(child.ID))
	if err != nil {
		t.Fatalf("fetch child failed: %v", err)
	}
	if childDB.FullCode != "0102" {
		t.Fatalf("expected child full_code unchanged '0102', got %q", childDB.FullCode)
	}
}
