package employeemovement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestService_CreateMovementDocument verifies plan §12.15: a document's
// metadata is persisted with the correct document_type/file_name/file_url.
func TestService_CreateMovementDocument(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(svc.repo, employeeID)

	resp, err := svc.CreateMovementDocument(ctx(), created.ID.String(), CreateMovementDocumentRequest{
		DocumentType: string(MovementDocumentTypePromotionSK),
		FileName:     "SK-Promosi-2026.pdf",
		FileURL:      "/uploads/attachments/abc123.pdf",
	})
	if err != nil {
		t.Fatalf("CreateMovementDocument failed: %v", err)
	}
	if resp.MovementID != created.ID.String() {
		t.Errorf("expected movement_id %s, got %s", created.ID.String(), resp.MovementID)
	}
	if resp.DocumentType != string(MovementDocumentTypePromotionSK) {
		t.Errorf("expected document_type PROMOTION_SK, got '%s'", resp.DocumentType)
	}
	if resp.FileName != "SK-Promosi-2026.pdf" {
		t.Errorf("expected file_name 'SK-Promosi-2026.pdf', got '%s'", resp.FileName)
	}
	if resp.FileURL != "/uploads/attachments/abc123.pdf" {
		t.Errorf("expected file_url '/uploads/attachments/abc123.pdf', got '%s'", resp.FileURL)
	}
	if resp.ID == "" {
		t.Error("expected document id to be set")
	}
}

// TestService_CreateMovementDocument_MovementNotFound verifies the service
// rejects adding a document to a movement that doesn't exist.
func TestService_CreateMovementDocument_MovementNotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	_, err := svc.CreateMovementDocument(ctx(), uuid.New().String(), CreateMovementDocumentRequest{
		DocumentType: string(MovementDocumentTypeOther),
		FileName:     "lampiran.pdf",
		FileURL:      "/uploads/attachments/lampiran.pdf",
	})
	if err == nil {
		t.Fatal("expected error when movement does not exist, got nil")
	}
}

// TestService_ListMovementDocuments verifies documents are listed newest-first
// with pagination (multiple documents per movement, plan §12.15).
func TestService_ListMovementDocuments(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(svc.repo, employeeID)

	for i := 0; i < 3; i++ {
		if _, err := svc.CreateMovementDocument(ctx(), created.ID.String(), CreateMovementDocumentRequest{
			DocumentType: string(MovementDocumentTypePromotionSK),
			FileName:     "SK-Promosi-" + string(rune('A'+i)) + ".pdf",
			FileURL:      "/uploads/attachments/doc-" + string(rune('A'+i)) + ".pdf",
		}); err != nil {
			t.Fatalf("CreateMovementDocument %d failed: %v", i, err)
		}
		// Jeda kecil agar created_at berbeda (SQLite menyimpan mikro-detik;
		// tanpa jeda, dua insert berurutan bisa berbagi timestamp yang sama
		// dan urutan "terbaru dulu" jatuh ke tie-break UUID acak).
		time.Sleep(2 * time.Millisecond)
	}

	page1, err := svc.ListMovementDocuments(ctx(), created.ID.String(), 1, 2)
	if err != nil {
		t.Fatalf("ListMovementDocuments failed: %v", err)
	}
	if page1.Total != 3 {
		t.Fatalf("expected total 3, got %d", page1.Total)
	}
	if page1.TotalPages != 2 {
		t.Fatalf("expected total_pages 2, got %d", page1.TotalPages)
	}
	items := page1.Data.([]MovementDocumentResponse)
	if len(items) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(items))
	}
	if items[0].FileName != "SK-Promosi-C.pdf" {
		t.Errorf("expected newest document first (C), got '%s'", items[0].FileName)
	}
}

// TestService_DeleteMovementDocument verifies deleting a document removes only
// that document (not the others on the same movement).
func TestService_DeleteMovementDocument(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(svc.repo, employeeID)

	doc1, err := svc.CreateMovementDocument(ctx(), created.ID.String(), CreateMovementDocumentRequest{
		DocumentType: string(MovementDocumentTypePromotionSK),
		FileName:     "SK-A.pdf",
		FileURL:      "/uploads/attachments/a.pdf",
	})
	if err != nil {
		t.Fatalf("CreateMovementDocument 1 failed: %v", err)
	}
	doc2, err := svc.CreateMovementDocument(ctx(), created.ID.String(), CreateMovementDocumentRequest{
		DocumentType: string(MovementDocumentTypeMutationSK),
		FileName:     "SK-B.pdf",
		FileURL:      "/uploads/attachments/b.pdf",
	})
	if err != nil {
		t.Fatalf("CreateMovementDocument 2 failed: %v", err)
	}

	if err := svc.DeleteMovementDocument(ctx(), doc1.ID); err != nil {
		t.Fatalf("DeleteMovementDocument failed: %v", err)
	}

	// doc1 gone, doc2 remains.
	if _, err := svc.repo.FindMovementDocumentByID(ctx(), uuid.MustParse(doc1.ID)); err == nil {
		t.Error("expected doc1 to be deleted")
	}
	remaining, err := svc.repo.FindMovementDocumentByID(ctx(), uuid.MustParse(doc2.ID))
	if err != nil {
		t.Fatalf("expected doc2 to remain: %v", err)
	}
	if remaining.FileName != "SK-B.pdf" {
		t.Errorf("expected doc2 file_name 'SK-B.pdf', got '%s'", remaining.FileName)
	}
}

// TestService_DeleteMovementDocument_NotFound verifies deleting a non-existent
// document returns an error.
func TestService_DeleteMovementDocument_NotFound(t *testing.T) {
	svc, _, cleanup := newTestService()
	defer cleanup()

	if err := svc.DeleteMovementDocument(ctx(), uuid.New().String()); err == nil {
		t.Fatal("expected error when document does not exist, got nil")
	}
}

// TestHandler_MovementDocuments verifies the documents endpoints end-to-end:
// POST create → GET list → DELETE remove.
func TestHandler_MovementDocuments(t *testing.T) {
	router, repo, cleanup := setupTestRouter()
	defer cleanup()

	employeeID := uuid.New()
	created := createTestMovement(repo, employeeID)

	// POST create
	body := `{
		"document_type": "MUTATION_SK",
		"file_name": "SK-Mutasi.pdf",
		"file_url": "/uploads/attachments/mutasi.pdf"
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tenant/employee-movements/movements/"+created.ID.String()+"/documents", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}

	var createdResp struct {
		Data MovementDocumentResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &createdResp); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}
	if createdResp.Data.DocumentType != string(MovementDocumentTypeMutationSK) {
		t.Errorf("expected document_type MUTATION_SK, got '%s'", createdResp.Data.DocumentType)
	}

	// GET list
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/tenant/employee-movements/movements/"+created.ID.String()+"/documents", nil)
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w2.Code, w2.Body.String())
	}
	var listResp struct {
		Success bool                       `json:"success"`
		Total   int64                      `json:"total"`
		Data    []MovementDocumentResponse `json:"data"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to unmarshal list response: %v", err)
	}
	if !listResp.Success || listResp.Total != 1 || len(listResp.Data) != 1 {
		t.Fatalf("expected 1 document in list, got success=%v total=%d len=%d", listResp.Success, listResp.Total, len(listResp.Data))
	}

	// DELETE
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("DELETE", "/api/v1/tenant/employee-movements/movements/"+created.ID.String()+"/documents/"+createdResp.Data.ID, nil)
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on delete, got %d: %s", w3.Code, w3.Body.String())
	}
}
