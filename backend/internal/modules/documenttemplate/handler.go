package documenttemplate

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

const (
	// maxTemplateFileSize membatasi ukuran file template .docx (10 MB).
	maxTemplateFileSize = 10 * 1024 * 1024
	// templateUploadSubDir adalah sub-direktori upload untuk template DOCX.
	templateUploadSubDir = "document_templates"
)

type Handler struct {
	svc       *Service
	uploadDir string
	pdf       PDFService
}

// NewHandler membuat handler. pdfServices opsional (variadic) — diisi saat
// preview PDF diaktifkan; tanpa PDFService, endpoint preview mengembalikan 503.
func NewHandler(svc *Service, uploadDir string, pdfServices ...PDFService) *Handler {
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	h := &Handler{svc: svc, uploadDir: uploadDir}
	if len(pdfServices) > 0 {
		h.pdf = pdfServices[0]
	}
	return h
}

func actorID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	var invalidTypeErr *InvalidDocumentTypeError
	if errors.As(err, &invalidTypeErr) {
		httputil.ErrorSimple(c, http.StatusBadRequest, invalidTypeErr.Error())
		return
	}
	var invalidMovementErr *InvalidMovementTypeError
	if errors.As(err, &invalidMovementErr) {
		httputil.ErrorSimple(c, http.StatusBadRequest, invalidMovementErr.Error())
		return
	}
	var movementNotAllowedErr *MovementTypeNotAllowedError
	if errors.As(err, &movementNotAllowedErr) {
		httputil.ErrorSimple(c, http.StatusBadRequest, movementNotAllowedErr.Error())
		return
	}
	var dupCodeErr *DuplicateCodeError
	if errors.As(err, &dupCodeErr) {
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_CODE", "documenttemplate.duplicate_code", dupCodeErr.Code)
		return
	}
	var dupActiveErr *DuplicateActiveTemplateError
	if errors.As(err, &dupActiveErr) {
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_ACTIVE", "documenttemplate.duplicate_active", dupActiveErr.DocumentType)
		return
	}
	if errors.Is(err, ErrTemplateNotFound) || errors.Is(err, ErrVersionNotFound) {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.InternalError(c, err.Error())
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	items, total, err := h.svc.List(c.Request.Context(), page, perPage, c.Query("document_type"), c.Query("movement_type"), c.Query("status"), c.Query("search"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, TemplateListResponse{Data: items, Total: total, Page: page})
}

func (h *Handler) ListVariables(c *gin.Context) {
	httputil.SuccessJSON(c, VariableRegistry())
}

// ListMovementTypes mengembalikan daftar jenis movement yang tersedia untuk
// template SK Movement (nilai + label bilingual dipakai frontend).
func (h *Handler) ListMovementTypes(c *gin.Context) {
	httputil.SuccessJSON(c, MovementTypeOptions())
}

// copyFile menyalin file (util kecil untuk memindahkan PDF hasil konversi ke
// direktori uploads agar diserve publik).
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// Preview menghasilkan PDF dari versi aktif template menggunakan data contoh,
// memakai pipeline yang sama dengan Generate Document (DOCX → LibreOffice →
// PDF). Hasil disimpan di {uploadDir}/previews/ dan dikembalikan sebagai
// { pdf_url, file_name } untuk ditampilkan di PDF viewer Settings (Phase 4).
func (h *Handler) Preview(c *gin.Context) {
	if h.pdf == nil {
		httputil.ErrorJSON(c, http.StatusServiceUnavailable, "PDF_ENGINE_NOT_CONFIGURED", "documenttemplate.pdf_engine_unavailable")
		return
	}

	tpl, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	if tpl.ActiveVersionID == nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.no_active_version")
		return
	}
	v, err := h.svc.GetVersion(c.Request.Context(), tpl.ID, *tpl.ActiveVersionID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	if versionFileURL(v.Content) == "" {
		// Konten bukan path file DOCX (mis. versi HTML lama) — tidak bisa dipreview.
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.no_template_file")
		return
	}
	srcPath := filepath.Join(h.uploadDir, filepath.FromSlash(v.Content))
	if _, err := os.Stat(srcPath); err != nil {
		httputil.ErrorJSON(c, http.StatusNotFound, "FILE_NOT_FOUND", "documenttemplate.template_file_missing")
		return
	}

	workDir, err := os.MkdirTemp("", "dtpl-preview-")
	if err != nil {
		httputil.InternalError(c, "Failed to create preview workspace")
		return
	}
	defer os.RemoveAll(workDir)

	// 1. Salin + resolve variabel dengan data contoh.
	resolved := filepath.Join(workDir, "resolved.docx")
	if err := resolveDocxVariables(srcPath, resolved, sampleData()); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	// 2. Konversi DOCX → PDF (LibreOffice headless).
	outPdf := filepath.Join(workDir, "preview.pdf")
	if err := h.pdf.ConvertDOCXToPDF(c.Request.Context(), resolved, outPdf); err != nil {
		// LibreOffice tidak terinstall/terkonfigurasi → 503 dengan pesan jelas
		// (bukan INTERNAL_ERROR mentah).
		if strings.Contains(err.Error(), "LibreOffice not installed or not found") ||
			strings.Contains(err.Error(), "configured libreoffice binary") {
			httputil.ErrorJSON(c, http.StatusServiceUnavailable, "PDF_ENGINE_NOT_CONFIGURED", "documenttemplate.libreoffice_missing")
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	// 3. Simpan hasil ke {uploadDir}/previews/ (diserve publik via /uploads).
	pdfDir := filepath.Join(h.uploadDir, "previews")
	if err := os.MkdirAll(pdfDir, 0755); err != nil {
		httputil.InternalError(c, "Failed to create preview output directory")
		return
	}
	pdfName := uuid.NewString() + ".pdf"
	if err := copyFile(outPdf, filepath.Join(pdfDir, pdfName)); err != nil {
		httputil.InternalError(c, "Failed to store preview pdf")
		return
	}

	pdfURL := "/uploads/previews/" + pdfName
	httputil.SuccessJSON(c, gin.H{
		"pdf_url":   pdfURL,
		"file_name": fmt.Sprintf("preview_%s.pdf", tpl.Code),
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	tpl, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, tpl)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.Create(c.Request.Context(), req.Name, req.Code, req.DocumentType, req.MovementType, req.Description, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, tpl, "documenttemplate.created")
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.Update(c.Request.Context(), c.Param("id"), req.Name, req.Description, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.updated")
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), actorID(c)); err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.DeletedJSON(c, "documenttemplate.deleted")
}

func (h *Handler) Activate(c *gin.Context) {
	tpl, err := h.svc.Activate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.activated")
}

func (h *Handler) Deactivate(c *gin.Context) {
	tpl, err := h.svc.Deactivate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.deactivated")
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, versions)
}

func (h *Handler) GetVersionByID(c *gin.Context) {
	v, err := h.svc.GetVersion(c.Request.Context(), c.Param("id"), c.Param("versionId"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, v)
}

// parseMarginField membaca satu field margin multipart (default 20).
func parseMarginField(c *gin.Context, field string) int {
	v, _ := strconv.Atoi(c.PostForm(field))
	if v == 0 {
		return 20
	}
	return v
}

func (h *Handler) CreateVersion(c *gin.Context) {
	content := ""
	fileName := ""
	paperSize := "A4"
	orientation := "portrait"
	margins := [4]int{20, 20, 20, 20}
	var placeholders []string

	if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
		// Mode DOCX (plan baru): template diupload sebagai file .docx, disimpan
		// ke {uploadDir}/document_templates/, content = path relatif file.
		file, err := c.FormFile("file")
		if err != nil {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.file_required")
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".docx" {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.file_type_not_allowed")
			return
		}
		if file.Size > maxTemplateFileSize {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.file_too_large")
			return
		}

		dir := filepath.Join(h.uploadDir, templateUploadSubDir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			httputil.InternalError(c, "Failed to create upload directory")
			return
		}
		destName := uuid.NewString() + ext
		destPath := filepath.Join(dir, destName)
		if err := c.SaveUploadedFile(file, destPath); err != nil {
			httputil.InternalError(c, "Failed to save uploaded template file")
			return
		}
		content = filepath.ToSlash(filepath.Join(templateUploadSubDir, destName))
		fileName = file.Filename

		// Placeholder detection + variable validation (Phase 3): baca file yang
		// barusan disimpan, deteksi {{key}} di XML word/, dan pastikan semua key
		// terdaftar di registry. Unknown placeholder → tolak (400) dengan daftar.
		data, err := os.ReadFile(destPath)
		if err != nil {
			httputil.InternalError(c, "Failed to read uploaded template file")
			return
		}
		pl, perr := extractDocxPlaceholders(data)
		if perr != nil {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.invalid_docx")
			return
		}
		placeholders = pl
		if unknown := unknownPlaceholders(placeholders); len(unknown) > 0 {
			httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "documenttemplate.unknown_variables", strings.Join(unknown, ", "))
			return
		}
		if ps := c.PostForm("paper_size"); ps != "" {
			paperSize = ps
		}
		if or := c.PostForm("orientation"); or != "" {
			orientation = or
		}
		margins = [4]int{
			parseMarginField(c, "margin_top"),
			parseMarginField(c, "margin_right"),
			parseMarginField(c, "margin_bottom"),
			parseMarginField(c, "margin_left"),
		}
	} else {
		// Mode JSON (backward compat dengan versi lama / konten HTML):
		// tetap didukung agar konsumen lama tidak rusak selama transisi.
		var req CreateVersionRequest
		if !httputil.BindAndValidate(c, &req) {
			return
		}
		content = req.Content
		if req.PaperSize != "" {
			paperSize = req.PaperSize
		}
		if req.Orientation != "" {
			orientation = req.Orientation
		}
		margins = [4]int{req.MarginTop, req.MarginRight, req.MarginBottom, req.MarginLeft}
		for i, m := range margins {
			if m == 0 {
				margins[i] = 20
			}
		}
	}

	v, err := h.svc.CreateVersion(c.Request.Context(), c.Param("id"), content, paperSize, orientation, margins, fileName, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	// Pada mode multipart sudah dihitung placeholders; mode JSON (backward
	// compat) tidak punya file sehingga placeholders kosong.
	resp := gin.H{"version": v}
	if placeholders != nil {
		resp["placeholders"] = placeholders
	}
	httputil.CreatedJSON(c, resp, "documenttemplate.version_created")
}
