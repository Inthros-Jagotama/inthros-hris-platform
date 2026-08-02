package setting

import (
	"context"
	"errors"
	"testing"
)

func strPtr2(s string) *string { return &s }

// ── Competency Service Tests ──

func TestService_CreateCompetency(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	resp, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{
		Name:       "Leadership (LEA)",
		Field:      strPtr2("Manajerial"),
		Cluster:    strPtr2("Manajerial"),
		Definition: strPtr2("Kemampuan mengarahkan SDM."),
	})
	if err != nil {
		t.Fatalf("CreateCompetency failed: %v", err)
	}
	if resp.ID == "" || resp.Name != "Leadership (LEA)" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Field == nil || *resp.Field != "Manajerial" {
		t.Fatalf("expected field Manajerial, got %v", resp.Field)
	}
}

func TestService_CreateCompetency_RequiredName(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	// name kosong → validasi binding tidak dijalankan di service, tapi GORM
	// not null akan gagal di DB — pastikan error dikembalikan, bukan panic.
	_, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestService_CreateCompetency_DuplicateName(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	req := CreateCompetencyRequest{Name: "Tenacity (TEN)", Field: strPtr2("Potential")}
	if _, err := svc.CreateCompetency(context.Background(), req); err != nil {
		t.Fatalf("first create failed: %v", err)
	}
	_, err := svc.CreateCompetency(context.Background(), req)
	if err == nil {
		t.Fatal("expected duplicate name error")
	}
	var dupErr *DuplicateCodeError
	if !errors.As(err, &dupErr) {
		t.Fatalf("expected DuplicateCodeError, got %v", err)
	}
}

func TestService_ListCompetencies_Pagination(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	for _, name := range []string{"Alpha", "Beta", "Gamma"} {
		if _, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{Name: name}); err != nil {
			t.Fatalf("create %s failed: %v", name, err)
		}
	}

	resp, err := svc.ListCompetencies(context.Background(), 1, 2, "")
	if err != nil {
		t.Fatalf("ListCompetencies failed: %v", err)
	}
	if resp.Total != 3 {
		t.Fatalf("expected total 3, got %d", resp.Total)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(resp.Data))
	}
	if resp.TotalPages != 2 {
		t.Fatalf("expected 2 total pages, got %d", resp.TotalPages)
	}
	// Urutan name ASC
	if resp.Data[0].Name != "Alpha" || resp.Data[1].Name != "Beta" {
		t.Fatalf("expected sorted names, got %+v", resp.Data)
	}
}

func TestService_ListCompetencies_Search(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	for _, name := range []string{"Leadership (LEA)", "Tenacity (TEN)", "Integrity (INT)"} {
		if _, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{Name: name}); err != nil {
			t.Fatalf("create %s failed: %v", name, err)
		}
	}

	// Search by name
	resp, err := svc.ListCompetencies(context.Background(), 1, 10, "lead")
	if err != nil {
		t.Fatalf("ListCompetencies failed: %v", err)
	}
	if resp.Total != 1 || len(resp.Data) != 1 || resp.Data[0].Name != "Leadership (LEA)" {
		t.Fatalf("expected 1 leadership result, got total=%d data=%+v", resp.Total, resp.Data)
	}

	// Search tanpa match
	resp2, err := svc.ListCompetencies(context.Background(), 1, 10, "tidak-ada")
	if err != nil {
		t.Fatalf("ListCompetencies failed: %v", err)
	}
	if resp2.Total != 0 || len(resp2.Data) != 0 {
		t.Fatalf("expected 0 results, got total=%d data=%+v", resp2.Total, resp2.Data)
	}
}

func TestService_UpdateCompetency(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	created, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{
		Name:    "Before",
		Field:   strPtr2("Core"),
		Cluster: strPtr2("Core"),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	updated, err := svc.UpdateCompetency(context.Background(), created.ID, UpdateCompetencyRequest{
		Name:       strPtr2("After"),
		Definition: strPtr2("Def baru"),
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "After" {
		t.Fatalf("expected name After, got %s", updated.Name)
	}
	if updated.Definition == nil || *updated.Definition != "Def baru" {
		t.Fatalf("expected definition updated, got %v", updated.Definition)
	}
	// Field tidak diubah
	if updated.Field == nil || *updated.Field != "Core" {
		t.Fatalf("expected field Core preserved, got %v", updated.Field)
	}
}

func TestService_UpdateCompetency_DuplicateName(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if _, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{Name: "A"}); err != nil {
		t.Fatalf("create A failed: %v", err)
	}
	b, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{Name: "B"})
	if err != nil {
		t.Fatalf("create B failed: %v", err)
	}

	_, err = svc.UpdateCompetency(context.Background(), b.ID, UpdateCompetencyRequest{Name: strPtr2("A")})
	if err == nil {
		t.Fatal("expected duplicate name error on update")
	}
	var dupErr *DuplicateCodeError
	if !errors.As(err, &dupErr) {
		t.Fatalf("expected DuplicateCodeError, got %v", err)
	}
}

func TestService_DeleteCompetency(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	created, err := svc.CreateCompetency(context.Background(), CreateCompetencyRequest{Name: "To Delete"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := svc.DeleteCompetency(context.Background(), created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := svc.GetCompetencyByID(context.Background(), created.ID); err == nil {
		t.Fatal("expected not found after delete")
	}
}
