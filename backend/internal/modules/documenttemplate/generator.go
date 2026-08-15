package documenttemplate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateRequest berisi data yang dibutuhkan Generate Document (Phase 5).
// Business module (contract/movement) yang membangun Values — Generator hanya
// bertanggung jawab rendering + PDF + penyimpanan, sesuai prinsip desain plan
// §24: "Template Dokumen tetap di Settings, Generate Document di module bisnis".
type GenerateRequest struct {
	DocumentType string
	// MovementType: untuk MOVEMENT_SK — jenis movement (promotion/mutation/dll)
	// agar template SK per jenis movement dipakai; kosong → template umum.
	MovementType  string
	ReferenceType string // "contract" | "movement"
	ReferenceID   string
	Values        map[string]string // variable yang sudah di-resolve oleh business module
	GeneratedBy   string
}

// CompanyInfo adalah data perusahaan (platform DB) yang diisi Generator untuk
// variable company.name / company.address — business module tidak perlu tahu
// dari mana data company berasal.
type CompanyInfo struct {
	Name    string
	Address string
}

// CompanyProvider mengambil data company aktif dari konteks (company_id di ctx).
// Di-wire dari main.go dengan akses ke platform DB.
type CompanyProvider func(ctx context.Context) (*CompanyInfo, error)

// Generator adalah shared document generator (plan §15/§24): mengambil template
// aktif + versi aktif, resolve variable, konversi DOCX → PDF, simpan ke
// {uploadDir}/generated_documents/, dan catat ke tabel generated_documents.
type Generator struct {
	svc             *Service
	uploadDir       string
	pdf             PDFService
	companyProvider CompanyProvider
}

// NewGenerator membuat generator. uploadDir = root direktori upload (diserve
// publik via /uploads). pdf nil → Generate mengembalikan error engine-unconfigured.
func NewGenerator(svc *Service, uploadDir string, pdf PDFService) *Generator {
	return &Generator{svc: svc, uploadDir: uploadDir, pdf: pdf}
}

// SetCompanyProvider mengatur provider data company (platform DB). Opsional —
// tanpa provider, variable company.* dibiarkan kosong.
func (g *Generator) SetCompanyProvider(p CompanyProvider) {
	g.companyProvider = p
}

var (
	// ErrNoActiveTemplate: tidak ada template ACTIVE untuk document type ini.
	ErrNoActiveTemplate = errors.New("no active template for this document type")
	// ErrPdfEngineNotConfigured: generator dipakai tanpa PDF service.
	ErrPdfEngineNotConfigured = errors.New("pdf engine not configured")
)

// Generate menghasilkan PDF dari template aktif untuk reference (contract /
// movement) dan mencatatnya di generated_documents. Mengembalikan record yang
// sudah di-decorate dengan FileURL.
func (g *Generator) Generate(ctx context.Context, req GenerateRequest) (*GeneratedDocument, error) {
	if g.pdf == nil {
		return nil, ErrPdfEngineNotConfigured
	}
	if req.DocumentType == "" || req.ReferenceType == "" || req.ReferenceID == "" {
		return nil, fmt.Errorf("document_type, reference_type, and reference_id are required")
	}

	values := map[string]string{}
	for k, v := range req.Values {
		values[k] = v
	}
	if g.companyProvider != nil {
		if ci, err := g.companyProvider(ctx); err == nil {
			if ci.Name != "" {
				values["company.name"] = ci.Name
			}
			if ci.Address != "" {
				values["company.address"] = ci.Address
			}
		}
	}

	tpl, err := g.svc.repo.FindActiveByTypeAndMovement(ctx, req.DocumentType, req.MovementType)
	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return nil, ErrNoActiveTemplate
		}
		return nil, err
	}
	if tpl.ActiveVersionID == nil || *tpl.ActiveVersionID == "" {
		return nil, ErrNoActiveTemplate
	}
	v, err := g.svc.repo.GetVersion(ctx, tpl.ID, *tpl.ActiveVersionID)
	if err != nil {
		return nil, err
	}
	if versionFileURL(v.Content) == "" {
		return nil, fmt.Errorf("active template version has no DOCX file to render")
	}
	srcPath := filepath.Join(g.uploadDir, filepath.FromSlash(v.Content))
	if _, err := os.Stat(srcPath); err != nil {
		return nil, fmt.Errorf("template docx file missing on server: %w", err)
	}

	workDir, err := os.MkdirTemp("", "dtpl-generate-")
	if err != nil {
		return nil, fmt.Errorf("failed to create generation workspace: %w", err)
	}
	defer os.RemoveAll(workDir)

	// 1. Resolve variable → DOCX.
	resolved := filepath.Join(workDir, "resolved.docx")
	if err := resolveDocxVariables(srcPath, resolved, values); err != nil {
		return nil, fmt.Errorf("failed to resolve document variables: %w", err)
	}
	// 2. DOCX → PDF.
	outPdf := filepath.Join(workDir, "generated.pdf")
	if err := g.pdf.ConvertDOCXToPDF(ctx, resolved, outPdf); err != nil {
		return nil, err
	}
	// 3. Simpan ke uploads/generated_documents/ (diserve publik via /uploads).
	pdfDir := filepath.Join(g.uploadDir, "generated_documents")
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create generated documents directory: %w", err)
	}
	fileName := sanitizeFileName(fmt.Sprintf("%s_%s.pdf", tpl.Code, req.ReferenceID))
	relPath := filepath.ToSlash(filepath.Join("generated_documents", fileName))
	if err := copyFile(outPdf, filepath.Join(g.uploadDir, filepath.FromSlash(relPath))); err != nil {
		return nil, fmt.Errorf("failed to store generated pdf: %w", err)
	}

	doc := &GeneratedDocument{
		ID:                uuid.New().String(),
		TemplateID:        tpl.ID,
		TemplateVersionID: v.ID,
		DocumentType:      req.DocumentType,
		ReferenceType:     req.ReferenceType,
		ReferenceID:       req.ReferenceID,
		FileName:          fileName,
		FilePath:          relPath,
		MimeType:          "application/pdf",
		GeneratedBy:       stringPtr(req.GeneratedBy),
		GeneratedAt:       time.Now(),
	}
	if err := g.svc.repo.CreateGeneratedDocument(ctx, doc); err != nil {
		return nil, err
	}
	// Audit: Document Generated (plan §19).
	_ = g.svc.repo.CreateAudit(ctx, nil, &DocumentTemplateAudit{
		ID:         uuid.New().String(),
		TemplateID: tpl.ID,
		VersionID:  &v.ID,
		Action:     "DOCUMENT_GENERATED",
		ActorID:    stringPtr(req.GeneratedBy),
	})
	decorateGeneratedFileURL(doc)
	return doc, nil
}

// ListByReference mengembalikan histori dokumen yang digenerate untuk sebuah
// reference (contract/movement), terbaru dulu.
func (g *Generator) ListByReference(ctx context.Context, referenceType, referenceID string, page, perPage int) ([]GeneratedDocument, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	docs, total, err := g.svc.repo.ListGeneratedByReference(ctx, referenceType, referenceID, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	for i := range docs {
		decorateGeneratedFileURL(&docs[i])
	}
	return docs, total, nil
}

// CreateGeneratedDocument menulis record generated_documents.
func (r *Repository) CreateGeneratedDocument(ctx context.Context, doc *GeneratedDocument) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(doc).Error; err != nil {
		return fmt.Errorf("failed to create generated document record: %w", err)
	}
	return nil
}

// ListGeneratedByReference menampilkan histori generated document untuk
// reference tertentu (contract/movement), terbaru dulu.
func (r *Repository) ListGeneratedByReference(ctx context.Context, referenceType, referenceID string, page, perPage int) ([]GeneratedDocument, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Session(&gorm.Session{}).Model(&GeneratedDocument{}).
		Where("reference_type = ? AND reference_id = ?", referenceType, referenceID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count generated documents: %w", err)
	}
	var docs []GeneratedDocument
	offset := (page - 1) * perPage
	if err := query.Order("generated_at DESC").Offset(offset).Limit(perPage).Find(&docs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list generated documents: %w", err)
	}
	return docs, total, nil
}

// decorateGeneratedFileURL mengisi FileURL (response-only) untuk PDF hasil
// generate — path relatif "generated_documents/{file}.pdf" → "/uploads/...".
func decorateGeneratedFileURL(doc *GeneratedDocument) {
	if doc == nil {
		return
	}
	doc.FileURL = "/uploads/" + doc.FilePath
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_")
	return replacer.Replace(name)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
