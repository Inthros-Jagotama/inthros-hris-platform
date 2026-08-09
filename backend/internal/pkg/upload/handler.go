// Package upload menyediakan endpoint upload file generik untuk lampiran
// modul tenant (docs/module-attendance-plan.md §32b.4) — konsumen pertama:
// lampiran isian aktual lembur. Berbeda dari upload terikat-resource di modul
// employee (foto/dokumen), handler ini tidak tahu konteks entitas apa pun:
// hanya menyimpan file dan mengembalikan URL publik.
package upload

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

const (
	// maxFileSize membatasi ukuran file upload (10 MB).
	maxFileSize = 10 * 1024 * 1024
	// subDir adalah sub-direktori di dalam upload dir untuk lampiran.
	subDir = "attachments"
)

// allowedExts adalah daftar ekstensi file yang diizinkan.
var allowedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".txt": true, ".csv": true,
}

// Handler melayani POST /api/v1/tenant/uploads.
type Handler struct {
	uploadDir string
}

// NewHandler membuat Handler. uploadDir kosong → default "uploads".
func NewHandler(uploadDir string) *Handler {
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	return &Handler{uploadDir: uploadDir}
}

// Upload menerima multipart field "file", menyimpannya ke
// {uploadDir}/attachments/{uuid}{ext}, lalu mengembalikan URL publik
// "/uploads/attachments/{uuid}{ext}" (diserve oleh router r.Static("/uploads")).
func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "upload.file_required")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "upload.type_not_allowed")
		return
	}
	if file.Size > maxFileSize {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "upload.file_too_large")
		return
	}

	dir := filepath.Join(h.uploadDir, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		httputil.InternalError(c, "Failed to create upload directory")
		return
	}

	filename := uuid.NewString() + ext
	destPath := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		httputil.InternalError(c, "Failed to save uploaded file")
		return
	}

	// URL publik sesuai router.Static("/uploads", uploadDir) —
	// mis. uploadDir="uploads" → "/uploads/attachments/{uuid}.pdf".
	url := "/" + filepath.ToSlash(filepath.Join(h.uploadDir, subDir, filename))
	httputil.SuccessJSON(c, gin.H{"url": url})
}
