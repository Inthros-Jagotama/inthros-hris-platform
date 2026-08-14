package documenttemplate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// mockGenPDFService menulis PDF dummy — mensimulasikan hasil LibreOffice tanpa
// membutuhkan binary terinstall (sama seperti mockPDFService di handler_test).
type mockGenPDFService struct{}

func (mockGenPDFService) ConvertDOCXToPDF(_ context.Context, _ string, outputPath string) error {
	return os.WriteFile(outputPath, []byte("%PDF-1.4 mock-generated"), 0644)
}

// setupActiveDocxTemplate membuat template ACTIVE dengan versi DOCX yang punya
// placeholder, mengembalikan template + uploadDir (file docx tersimpan di disk).
func setupActiveDocxTemplate(t *testing.T, repo *Repository, uploadDir, docType, body string) *DocumentTemplate {
	t.Helper()
	ctx := context.Background()
	db, _ := repo.getDB(ctx)
	tpl := createTestTemplate(db, "", docType, StatusActive)
	tpl, _ = repo.GetByID(ctx, tpl.ID)

	// Simpan file docx ke uploadDir/document_templates/ + buat versi.
	docxName := uuid.NewString() + ".docx"
	relPath := "document_templates/" + docxName
	if err := os.MkdirAll(filepath.Join(uploadDir, "document_templates"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, filepath.FromSlash(relPath)), makeDocx(t, body), 0644); err != nil {
		t.Fatalf("write docx: %v", err)
	}
	svc := NewService(repo, zap.NewNop())
	v, err := svc.CreateVersion(ctx, tpl.ID, relPath, "A4", "portrait", [4]int{20, 20, 20, 20}, "template.docx", "user-1")
	if err != nil {
		t.Fatalf("create version: %v", err)
	}
	// Activate agar template punya versi aktif (CreateVersion sudah set
	// active_version_id, tapi pastikan status ACTIVE — createTestTemplate sudah).
	_ = v
	return tpl
}

func TestGeneratorGenerateCreatesPDFAndRecord(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	uploadDir := filepath.Join(t.TempDir(), "uploads")

	tpl := setupActiveDocxTemplate(t, repo, uploadDir, DocumentTypeContractAgreement, "Nama: {{employee.name}} Kontrak: {{contract.number}}")
	if tpl == nil {
		t.Fatal("setup failed")
	}

	gen := NewGenerator(NewService(repo, zap.NewNop()), uploadDir, mockGenPDFService{})
	gen.SetCompanyProvider(func(ctx context.Context) (*CompanyInfo, error) {
		return &CompanyInfo{Name: "PT Maju", Address: "Jl. Merdeka"}, nil
	})

	refID := uuid.NewString()
	doc, err := gen.Generate(context.Background(), GenerateRequest{
		DocumentType:  DocumentTypeContractAgreement,
		ReferenceType: "contract",
		ReferenceID:   refID,
		Values:        map[string]string{"employee.name": "Asep", "contract.number": "CTR-001"},
		GeneratedBy:   "user-1",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if doc.FileURL == "" || !strings.HasPrefix(doc.FileURL, "/uploads/generated_documents/") {
		t.Fatalf("expected file_url under /uploads/generated_documents/, got %q", doc.FileURL)
	}
	if !strings.Contains(doc.FileName, ".pdf") {
		t.Fatalf("expected pdf file name, got %q", doc.FileName)
	}
	// PDF tersimpan di disk.
	saved := filepath.Join(uploadDir, "generated_documents", doc.FileName)
	if _, err := os.Stat(saved); err != nil {
		t.Fatalf("generated pdf not stored: %v", err)
	}

	// Record tercatat di generated_documents + histori listable.
	docs, total, err := gen.ListByReference(context.Background(), "contract", refID, 1, 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(docs) != 1 {
		t.Fatalf("expected 1 generated document, got total=%d len=%d", total, len(docs))
	}
	if docs[0].ID != doc.ID || docs[0].FileURL != doc.FileURL {
		t.Fatalf("history record mismatch: %+v", docs[0])
	}
}

func TestGeneratorGenerateNoActiveTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })

	gen := NewGenerator(NewService(repo, zap.NewNop()), filepath.Join(t.TempDir(), "uploads"), mockGenPDFService{})
	_, err := gen.Generate(context.Background(), GenerateRequest{
		DocumentType:  DocumentTypeContractAgreement,
		ReferenceType: "contract",
		ReferenceID:   uuid.NewString(),
	})
	if err == nil {
		t.Fatal("expected error when no active template exists")
	}
	if !strings.Contains(err.Error(), "no active template") {
		t.Fatalf("expected no-active-template error, got: %v", err)
	}
}

func TestGeneratorGenerateWithoutPDFEngine(t *testing.T) {
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	setupActiveDocxTemplate(t, repo, uploadDir, DocumentTypeContractAgreement, "Nama: {{employee.name}}")

	gen := NewGenerator(NewService(repo, zap.NewNop()), uploadDir, nil)
	_, err := gen.Generate(context.Background(), GenerateRequest{
		DocumentType:  DocumentTypeContractAgreement,
		ReferenceType: "contract",
		ReferenceID:   uuid.NewString(),
	})
	if err == nil {
		t.Fatal("expected error when pdf engine is not configured")
	}
	if err != ErrPdfEngineNotConfigured {
		t.Fatalf("expected ErrPdfEngineNotConfigured, got: %v", err)
	}
}
