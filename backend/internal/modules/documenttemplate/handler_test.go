package documenttemplate

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, "")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)
	return r
}

func TestHandlerCreateAndList(t *testing.T) {
	r := setupTestRouter(t)

	payload := `{"name":"PKWT","code":"PKWT_H1","document_type":"CONTRACT_AGREEMENT"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-templates", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
	var body struct {
		Data struct {
			Data  []DocumentTemplate `json:"data"`
			Total int64              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listW.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.Total != 1 {
		t.Fatalf("expected 1 template, got %d", body.Data.Total)
	}
}

func TestHandlerCreateRejectsInvalidDocumentType(t *testing.T) {
	r := setupTestRouter(t)
	payload := `{"name":"X","code":"XCODE","document_type":"NOT_REAL"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerUpdateDefaultContentAndVersionDetail(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, "")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	// Create a normal template, create a version, then GET /:id/versions/:versionId
	tpl := createTestTemplate(db, "PKWT_ACTIVE", DocumentTypeContractAgreement, StatusInactive)
	v, err := svc.CreateVersion(context.Background(), tpl.ID, "<p>v1 content</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "", "user-1")
	if err != nil {
		t.Fatalf("create version: %v", err)
	}
	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions/"+v.ID, nil)
	detailW := httptest.NewRecorder()
	r.ServeHTTP(detailW, detailReq)
	if detailW.Code != http.StatusOK {
		t.Fatalf("expected 200 on version detail, got %d: %s", detailW.Code, detailW.Body.String())
	}
	var detailBody struct {
		Data DocumentTemplateVersion `json:"data"`
	}
	if err := json.Unmarshal(detailW.Body.Bytes(), &detailBody); err != nil {
		t.Fatalf("unmarshal version detail: %v", err)
	}
	if detailBody.Data.ID != v.ID || detailBody.Data.Content != "<p>v1 content</p>" {
		t.Fatalf("unexpected version detail: %+v", detailBody.Data)
	}

	// Version detail for a missing version id must 404 (not be shadowed by GET /:id)
	missingReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions/not-a-real-id", nil)
	missingW := httptest.NewRecorder()
	r.ServeHTTP(missingW, missingReq)
	if missingW.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing version, got %d: %s", missingW.Code, missingW.Body.String())
	}
}

func TestHandlerCreateVersionWithDocxUpload(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_DOCX", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	handler := NewHandler(svc, uploadDir)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	// Multipart body: field "file" (.docx valid berisi placeholder) + paper_size/orientation
	docxBytes := makeDocx(t, "Nomor: {{contract.number}} Nama: {{employee.name}}")
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "Perjanjian_PKWT.docx")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(docxBytes); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := mw.WriteField("paper_size", "A4"); err != nil {
		t.Fatalf("write paper_size: %v", err)
	}
	if err := mw.WriteField("orientation", "portrait"); err != nil {
		t.Fatalf("write orientation: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var body struct {
		Data struct {
			Version      DocumentTemplateVersion `json:"version"`
			Placeholders []string               `json:"placeholders"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.Version.FileName == nil || *body.Data.Version.FileName != "Perjanjian_PKWT.docx" {
		t.Fatalf("expected file_name Perjanjian_PKWT.docx, got %v", body.Data.Version.FileName)
	}
	if body.Data.Version.FileURL == "" || !strings.HasPrefix(body.Data.Version.FileURL, "/uploads/document_templates/") {
		t.Fatalf("expected file_url under /uploads/document_templates/, got %q", body.Data.Version.FileURL)
	}
	if !strings.HasPrefix(body.Data.Version.Content, "document_templates/") {
		t.Fatalf("expected content to be a relative file path, got %q", body.Data.Version.Content)
	}

	// Placeholder detection harus menemukan 2 variable terdaftar
	if len(body.Data.Placeholders) != 2 {
		t.Fatalf("expected 2 placeholders detected, got %v", body.Data.Placeholders)
	}

	// File benar-benar tersimpan di disk
	savedPath := filepath.Join(uploadDir, filepath.FromSlash(body.Data.Version.Content))
	if _, err := os.Stat(savedPath); err != nil {
		t.Fatalf("uploaded file not saved on disk: %v", err)
	}

	// ListVersions juga harus membawa file_url
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200 on versions list, got %d", listW.Code)
	}
	var listBody struct {
		Data []DocumentTemplateVersion `json:"data"`
	}
	if err := json.Unmarshal(listW.Body.Bytes(), &listBody); err != nil {
		t.Fatalf("unmarshal versions: %v", err)
	}
	if len(listBody.Data) != 1 || listBody.Data[0].FileURL != body.Data.Version.FileURL {
		t.Fatalf("expected one version with file_url, got %+v", listBody.Data)
	}
}

func TestHandlerCreateVersionRejectsUnknownPlaceholders(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_UNKNOWN", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, t.TempDir())
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	// Placeholder {{signature.author}} tidak terdaftar di registry → harus 400
	docxBytes := makeDocx(t, "Nama: {{employee.name}} TTD: {{signature.author}}")
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "Bad.docx")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(docxBytes); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unknown placeholder, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "signature.author") {
		t.Fatalf("expected error message to name the unknown variable, got: %s", w.Body.String())
	}
}

// mockPDFService menulis file PDF dummy ke outputPath — mensimulasikan hasil
// LibreOffice tanpa membutuhkan binary terinstall di mesin test.
type mockPDFService struct{}

func (mockPDFService) ConvertDOCXToPDF(_ context.Context, _ string, outputPath string) error {
	return os.WriteFile(outputPath, []byte("%PDF-1.4 mock-preview"), 0644)
}

func TestHandlerPreviewWithMockPDF(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_PREVIEW", DocumentTypeContractAgreement, StatusInactive)

	// Simpan file .docx di uploadDir/document_templates/ + buat versi aktif.
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	docxName := uuid.New().String() + ".docx"
	// Relative path pakai forward slash agar cocok dengan versionFileURL()
	// (prefix "document_templates/") dan path resolusi di handler.
	relPath := "document_templates/" + docxName
	if err := os.MkdirAll(filepath.Join(uploadDir, "document_templates"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	docxBytes := makeDocx(t, "Nama: {{employee.name}}")
	if err := os.WriteFile(filepath.Join(uploadDir, filepath.FromSlash(relPath)), docxBytes, 0644); err != nil {
		t.Fatalf("write template docx: %v", err)
	}

	svc := NewService(newTestRepo(db), zap.NewNop())
	v, err := svc.CreateVersion(context.Background(), tpl.ID, relPath, "A4", "portrait", [4]int{20, 20, 20, 20}, "", "user-1")
	if err != nil {
		t.Fatalf("create version: %v", err)
	}
	handler := NewHandler(svc, uploadDir, mockPDFService{})
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/preview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on preview, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			PdfURL   string `json:"pdf_url"`
			FileName string `json:"file_name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.PdfURL == "" || !strings.HasPrefix(body.Data.PdfURL, "/uploads/previews/") {
		t.Fatalf("expected pdf_url under /uploads/previews/, got %q", body.Data.PdfURL)
	}
	// PDF tersimpan di disk
	saved := filepath.Join(uploadDir, "previews", filepath.Base(body.Data.PdfURL))
	if _, err := os.Stat(saved); err != nil {
		t.Fatalf("preview pdf not stored: %v", err)
	}
	// Versi aktif tidak berubah (preview tidak membuat versi baru)
	got, _ := svc.GetByID(context.Background(), tpl.ID)
	if got.ActiveVersionID == nil || *got.ActiveVersionID != v.ID {
		t.Fatalf("preview must not change active version")
	}
}

func TestHandlerPreviewWithoutPDFEngine(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_NOENGINE", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, t.TempDir()) // tanpa PDFService
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/preview", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 without pdf engine, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateVersionRejectsInvalidDocx(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_BADZIP", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, t.TempDir())
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	// File berekstensi .docx tapi bukan zip valid → 400 invalid_docx
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "Broken.docx")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte("this is not a zip")); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid docx, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateVersionRejectsNonDocxUpload(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_BADEXT", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, t.TempDir())
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", "template.txt")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write([]byte("plain text")); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-docx upload, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerCreateVersionRejectsMissingFile(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	tpl := createTestTemplate(db, "PKWT_NOFILE", DocumentTypeContractAgreement, StatusInactive)

	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc, t.TempDir())
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant/settings")
	RegisterRoutes(group, handler)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	if err := mw.WriteField("paper_size", "A4"); err != nil {
		t.Fatalf("write field: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/settings/document-templates/"+tpl.ID+"/versions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing file, got %d: %s", w.Code, w.Body.String())
	}
}
