package employeemovement

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

// seedCareerPathPositions membuat tabel positions + mengisi beberapa posisi
// untuk test career path. setupTestDB hanya meng-automigrate model module;
// tabel referensi positions dibuat manual (pola sama seedCareerReferenceTables)
// karena GetPositionNamesByIDs / validasi posisi JOIN ke tabel itu.
func seedCareerPathPositions(t *testing.T, repo *Repository) []uuid.UUID {
	t.Helper()
	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	if err := db.Exec("CREATE TABLE IF NOT EXISTS positions (id CHAR(36) PRIMARY KEY, title VARCHAR(255))").Error; err != nil {
		t.Fatalf("failed to create positions table: %v", err)
	}

	posIDs := make([]uuid.UUID, 0, 4)
	for _, title := range []string{"Staff", "Senior Staff", "Supervisor", "Manager"} {
		id := uuid.New()
		if err := db.Exec("INSERT INTO positions (id, title) VALUES (?, ?)", id.String(), title).Error; err != nil {
			t.Fatalf("failed to seed position %s: %v", title, err)
		}
		posIDs = append(posIDs, id)
	}
	return posIDs
}

func careerPathStepsReq(posIDs []uuid.UUID) []CreateCareerPathStepRequest {
	return []CreateCareerPathStepRequest{
		{PositionID: posIDs[0].String(), Sequence: 1, MinimumServiceMonths: intPtr(0), Requirements: strPtr("Lulusan S1")},
		{PositionID: posIDs[1].String(), Sequence: 2, MinimumServiceMonths: intPtr(24), Requirements: strPtr("Performa >= 80")},
		{PositionID: posIDs[2].String(), Sequence: 3, MinimumServiceMonths: intPtr(36)},
	}
}

// TestService_CreateCareerPath verifies plan §12.9: career path + steps
// tersimpan, response membawa steps terurut sequence dengan nama posisi
// (enrichment, bukan UUID mentah).
func TestService_CreateCareerPath(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	resp, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name:        "Teknologi — Staf ke Supervisor",
		Description: strPtr("Jenjang umum divisi Teknologi"),
		Steps:       careerPathStepsReq(posIDs),
	})
	if err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}
	if resp.ID == "" {
		t.Error("expected career path id to be set")
	}
	if resp.Name != "Teknologi — Staf ke Supervisor" {
		t.Errorf("expected name, got '%s'", resp.Name)
	}
	if !resp.IsActive {
		t.Error("expected is_active true by default")
	}
	if resp.Description == nil || *resp.Description != "Jenjang umum divisi Teknologi" {
		t.Errorf("expected description preserved, got %v", resp.Description)
	}
	if len(resp.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(resp.Steps))
	}
	// Steps harus terurut sequence ascending.
	for i, st := range resp.Steps {
		if st.Sequence != i+1 {
			t.Errorf("expected sequence %d at index %d, got %d", i+1, i, st.Sequence)
		}
		if st.PositionID == "" {
			t.Errorf("step %d: expected position_id set", i)
		}
		if st.ID == "" {
			t.Errorf("step %d: expected step id set", i)
		}
	}
	// Enrichment nama posisi (batch query positions).
	if resp.Steps[0].PositionName != "Staff" {
		t.Errorf("expected position_name 'Staff', got '%s'", resp.Steps[0].PositionName)
	}
	if resp.Steps[1].PositionName != "Senior Staff" {
		t.Errorf("expected position_name 'Senior Staff', got '%s'", resp.Steps[1].PositionName)
	}
	if resp.Steps[1].MinimumServiceMonths == nil || *resp.Steps[1].MinimumServiceMonths != 24 {
		t.Errorf("expected minimum_service_months 24 on step 2, got %v", resp.Steps[1].MinimumServiceMonths)
	}
}

// TestService_CreateCareerPath_DuplicateSequence verifies a path with two
// steps sharing the same sequence is rejected (plan §12.9 — sequence unik
// per path).
func TestService_CreateCareerPath_DuplicateSequence(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	_, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name: "Jalur Duplikat",
		Steps: []CreateCareerPathStepRequest{
			{PositionID: posIDs[0].String(), Sequence: 1},
			{PositionID: posIDs[1].String(), Sequence: 1},
		},
	})
	if err == nil {
		t.Fatal("expected error on duplicate sequence, got nil")
	}
	var ve *MovementValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected MovementValidationError, got %T: %v", err, err)
	}
}

// TestService_CreateCareerPath_PositionNotFound verifies steps referencing a
// position that does not exist are rejected — validasi eksistensi posisi
// (bukan hanya format UUID) di service layer.
func TestService_CreateCareerPath_PositionNotFound(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	_, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name: "Jalur Satu Posisi Palsu",
		Steps: []CreateCareerPathStepRequest{
			{PositionID: posIDs[0].String(), Sequence: 1},
			{PositionID: uuid.New().String(), Sequence: 2},
		},
	})
	if err == nil {
		t.Fatal("expected error on non-existent position, got nil")
	}
	var ve *MovementValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected MovementValidationError, got %T: %v", err, err)
	}
}

// TestService_ListCareerPaths_Pagination verifies list + keyword filter +
// pagination, tiap item membawa steps-nya.
func TestService_ListCareerPaths_Pagination(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	// Prefix unik agar tidak bentrok dengan path test lain (DB SQLite shared
	// antar-test — pola sama flakiness audit/document test).
	names := []string{"Paginate-Alpha", "Paginate-Beta", "Paginate-Gamma"}
	for _, name := range names {
		if _, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
			Name:  name,
			Steps: careerPathStepsReq(posIDs),
		}); err != nil {
			t.Fatalf("CreateCareerPath %s failed: %v", name, err)
		}
	}

	page1, err := svc.ListCareerPaths(ctx(), 1, 2, "Paginate")
	if err != nil {
		t.Fatalf("ListCareerPaths failed: %v", err)
	}
	if page1.Total != 3 {
		t.Fatalf("expected total 3, got %d", page1.Total)
	}
	if page1.TotalPages != 2 {
		t.Fatalf("expected total_pages 2, got %d", page1.TotalPages)
	}
	items := page1.Data.([]CareerPathResponse)
	if len(items) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(items))
	}
	if items[0].Name != "Paginate-Alpha" {
		t.Errorf("expected alphabetical order 'Paginate-Alpha' first, got '%s'", items[0].Name)
	}
	if len(items[0].Steps) != 3 {
		t.Errorf("expected 3 steps on item, got %d", len(items[0].Steps))
	}

	// Filter keyword lebih sempit.
	filtered, err := svc.ListCareerPaths(ctx(), 1, 20, "Beta")
	if err != nil {
		t.Fatalf("ListCareerPaths keyword failed: %v", err)
	}
	if filtered.Total != 1 {
		t.Fatalf("expected 1 match for 'Beta', got %d", filtered.Total)
	}
	fItems := filtered.Data.([]CareerPathResponse)
	if len(fItems) != 1 || fItems[0].Name != "Paginate-Beta" {
		t.Errorf("expected only 'Paginate-Beta', got %+v", fItems)
	}
}

// TestService_GetCareerPathByID verifies detail path returns steps ordered by
// sequence with position names.
func TestService_GetCareerPathByID(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	created, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name:  "Jalur Spesifik",
		Steps: careerPathStepsReq(posIDs),
	})
	if err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}

	got, err := svc.GetCareerPathByID(ctx(), created.ID)
	if err != nil {
		t.Fatalf("GetCareerPathByID failed: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("expected id %s, got %s", created.ID, got.ID)
	}
	if len(got.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(got.Steps))
	}
	if got.Steps[0].PositionName != "Staff" {
		t.Errorf("expected step 1 position_name 'Staff', got '%s'", got.Steps[0].PositionName)
	}
}

// TestService_GetCareerPathByID_NotFound verifies a missing id returns a
// not-found error (mapped to 404 by the handler).
func TestService_GetCareerPathByID_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if _, err := svc.GetCareerPathByID(ctx(), uuid.New().String()); err == nil {
		t.Fatal("expected not found error, got nil")
	}
}

// TestService_UpdateCareerPath_FullReplace verifies PUT replaces the whole
// step list: remove old steps, add new ones, keep header updates.
func TestService_UpdateCareerPath_FullReplace(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	created, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name:        "Jalur Lama",
		Description: strPtr("Sebelum update"),
		Steps:       careerPathStepsReq(posIDs),
	})
	if err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}

	// Update: ganti nama, nonaktifkan, dan ganti steps menjadi 2 langkah saja.
	updated, err := svc.UpdateCareerPath(ctx(), created.ID, UpdateCareerPathRequest{
		Name:        strPtr("Jalur Baru"),
		Description: strPtr("Sesudah update"),
		IsActive:    boolPtr(false),
		Steps: []CreateCareerPathStepRequest{
			{PositionID: posIDs[1].String(), Sequence: 1, MinimumServiceMonths: intPtr(12)},
			{PositionID: posIDs[3].String(), Sequence: 2},
		},
	})
	if err != nil {
		t.Fatalf("UpdateCareerPath failed: %v", err)
	}
	if updated.Name != "Jalur Baru" {
		t.Errorf("expected updated name 'Jalur Baru', got '%s'", updated.Name)
	}
	if updated.IsActive {
		t.Error("expected is_active false after update")
	}
	if len(updated.Steps) != 2 {
		t.Fatalf("expected 2 steps after replace, got %d", len(updated.Steps))
	}
	if updated.Steps[0].PositionName != "Senior Staff" {
		t.Errorf("expected first step 'Senior Staff', got '%s'", updated.Steps[0].PositionName)
	}
	if updated.Steps[1].PositionName != "Manager" {
		t.Errorf("expected second step 'Manager', got '%s'", updated.Steps[1].PositionName)
	}

	// Verifikasi persistensi: detail ulang harus sama.
	got, err := svc.GetCareerPathByID(ctx(), created.ID)
	if err != nil {
		t.Fatalf("GetCareerPathByID after update failed: %v", err)
	}
	if got.Name != "Jalur Baru" || len(got.Steps) != 2 {
		t.Errorf("expected persisted update (name='Jalur Baru', 2 steps), got name='%s' steps=%d", got.Name, len(got.Steps))
	}
}

// TestService_DeleteCareerPath verifies delete removes the path and its steps.
func TestService_DeleteCareerPath(t *testing.T) {
	svc, repo, cleanup := newTestService()
	defer cleanup()

	posIDs := seedCareerPathPositions(t, repo)

	created, err := svc.CreateCareerPath(ctx(), CreateCareerPathRequest{
		Name:  "Jalur Dihapus",
		Steps: careerPathStepsReq(posIDs),
	})
	if err != nil {
		t.Fatalf("CreateCareerPath failed: %v", err)
	}

	if err := svc.DeleteCareerPath(ctx(), created.ID); err != nil {
		t.Fatalf("DeleteCareerPath failed: %v", err)
	}

	if _, err := svc.GetCareerPathByID(ctx(), created.ID); err == nil {
		t.Fatal("expected not found after delete, got nil")
	}

	// Steps juga harus terhapus (DeleteCareerPath menghapus keduanya).
	db, err := repo.getDB(ctx())
	if err != nil {
		t.Fatalf("failed to get test db: %v", err)
	}
	var stepCount int64
	if err := db.Model(&CareerPathStep{}).Where("career_path_id = ?", created.ID).Count(&stepCount).Error; err != nil {
		t.Fatalf("failed to count remaining steps: %v", err)
	}
	if stepCount != 0 {
		t.Errorf("expected 0 steps after delete, got %d", stepCount)
	}
}

// TestService_DeleteCareerPath_NotFound verifies deleting a missing path
// returns a not-found error.
func TestService_DeleteCareerPath_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if err := svc.DeleteCareerPath(ctx(), uuid.New().String()); err == nil {
		t.Fatal("expected not found error, got nil")
	}
}

// boolPtr returns a pointer to the given bool.
func boolPtr(b bool) *bool {
	return &b
}
