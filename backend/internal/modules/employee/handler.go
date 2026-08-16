package employee

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// POST /api/v1/tenant/employees
func (h *Handler) Create(c *gin.Context) {
	var req CreateEmployeeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		// Kuota employee on-premise tercapai → 403 dengan kode khusus,
		// bukan 500 agar client bisa menangkap error bisnis ini.
		if errors.Is(err, ErrQuotaExceeded) {
			httputil.ErrorRaw(c, http.StatusForbidden, "QUOTA_EXCEEDED", err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/employees
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	status := c.Query("status")
	organizationID := c.Query("organization_id")

	resp, err := h.service.List(c.Request.Context(), page, perPage, search, status, organizationID)
	if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /api/v1/tenant/employees/:id
func (h *Handler) GetByID(c *gin.Context) {
	resp, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
			httputil.NotFound(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

// PUT /api/v1/tenant/employees/:id
func (h *Handler) Update(c *gin.Context) {
	var req UpdateEmployeeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	resp, err := h.service.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.SuccessJSON(c, resp)
}

// DELETE /api/v1/tenant/employees/:id
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// PUT /api/v1/tenant/employees/:id/photo — Upload profile picture
func (h *Handler) UploadPhoto(c *gin.Context) {
	employeeID := c.Param("id")

	file, err := c.FormFile("photo")
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.photo_required")
		return
	}

	// Validate file type (only images)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.image_type_only")
		return
	}

	// Max 2MB
	if file.Size > 2*1024*1024 {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.file_max_2mb")
		return
	}

	// Build filename: employee_id + ext
	filename := employeeID + ext
	uploadDir := "uploads/employees"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		httputil.InternalError(c, "Failed to create upload directory")
		return
	}

	destPath := uploadDir + "/" + filename
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		httputil.InternalError(c, "Failed to save uploaded file")
		return
	}

	// Update employee profile_picture in database
	photoURL := "/" + destPath
	if err := h.service.UpdatePhoto(c.Request.Context(), employeeID, photoURL); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, gin.H{"profile_picture": photoURL})
}

// DELETE /api/v1/tenant/employees/:id/photo — Delete profile picture
func (h *Handler) DeletePhoto(c *gin.Context) {
	employeeID := c.Param("id")

	if err := h.service.DeletePhoto(c.Request.Context(), employeeID); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	// Also delete file from disk
	// Walk known extensions to find and delete the file
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif", ".webp"} {
		filePath := "uploads/employees/" + employeeID + ext
		if err := os.Remove(filePath); err == nil {
			break
		}
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Addresses
// =========================================================================

func (h *Handler) CreateAddress(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateAddressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateAddress(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	var req UpdateAddressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateAddress(c.Request.Context(), c.Param("id"), c.Param("addressId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteAddress(c *gin.Context) {
	if err := h.service.DeleteAddress(c.Request.Context(), c.Param("id"), c.Param("addressId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Emergency Contacts
// =========================================================================

func (h *Handler) CreateEmergencyContact(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateEmergencyContactRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmergencyContact(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateEmergencyContact(c *gin.Context) {
	var req UpdateEmergencyContactRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmergencyContact(c.Request.Context(), c.Param("id"), c.Param("contactId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmergencyContact(c *gin.Context) {
	if err := h.service.DeleteEmergencyContact(c.Request.Context(), c.Param("id"), c.Param("contactId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Families
// =========================================================================

func (h *Handler) CreateFamily(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateFamilyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateFamily(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateFamily(c *gin.Context) {
	var req UpdateFamilyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateFamily(c.Request.Context(), c.Param("id"), c.Param("familyId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteFamily(c *gin.Context) {
	if err := h.service.DeleteFamily(c.Request.Context(), c.Param("id"), c.Param("familyId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Educations
// =========================================================================

func (h *Handler) CreateEducation(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEducation(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateEducation(c *gin.Context) {
	var req UpdateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEducation(c.Request.Context(), c.Param("id"), c.Param("educationId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEducation(c *gin.Context) {
	if err := h.service.DeleteEducation(c.Request.Context(), c.Param("id"), c.Param("educationId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Experiences
// =========================================================================

func (h *Handler) CreateExperience(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateExperience(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateExperience(c *gin.Context) {
	var req UpdateExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateExperience(c.Request.Context(), c.Param("id"), c.Param("experienceId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteExperience(c *gin.Context) {
	if err := h.service.DeleteExperience(c.Request.Context(), c.Param("id"), c.Param("experienceId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Documents
// =========================================================================

func (h *Handler) CreateDocument(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateDocument(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// POST /api/v1/tenant/employees/:id/documents/upload — Upload document file
func (h *Handler) UploadDocumentFile(c *gin.Context) {
	employeeID := c.Param("id")

	name := c.PostForm("name")
	if name == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.name_required")
		return
	}

	var note *string
	if n := c.PostForm("note"); n != "" {
		note = &n
	}

	file, err := c.FormFile("file")
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.document_required")
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true,
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".txt": true,
	}
	if !allowedExts[ext] {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.file_type_not_allowed")
		return
	}

	// Max 10MB
	if file.Size > 10*1024*1024 {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.file_max_10mb")
		return
	}

	// Build unique filename: employeeID_timestamp + ext
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)
	filename := employeeID + "_" + timestamp + ext
	uploadDir := "uploads/documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		httputil.InternalError(c, "Failed to create upload directory")
		return
	}

	destPath := uploadDir + "/" + filename
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		httputil.InternalError(c, "Failed to save uploaded file")
		return
	}

	filePath := "/" + destPath

	empUID, err := uuid.Parse(employeeID)
	if err != nil {
		os.Remove(destPath)
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.invalid_id")
		return
	}

	// Create document record in DB
	doc := &EmployeeDocument{
		EmployeeID: &empUID,
		Name:       name,
		File:       filePath,
		Note:       note,
		CreatedBy:  authctx.GetUserID(c.Request.Context()),
		UpdatedBy:  authctx.GetUserID(c.Request.Context()),
	}

	if err := h.service.CreateDocumentRecord(c.Request.Context(), doc); err != nil {
		// Clean up file if DB insert fails
		os.Remove(destPath)
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, toDocumentResponse(doc), "success.created")
}

// PUT /api/v1/tenant/employees/:id/documents/:documentId/upload — Replace document file
func (h *Handler) UpdateDocumentFile(c *gin.Context) {
	employeeID := c.Param("id")
	documentID := c.Param("documentId")

	name := c.PostForm("name")
	var note *string
	if n := c.PostForm("note"); n != "" {
		note = &n
	}

	file, err := c.FormFile("file")
	if err != nil {
		// No new file — update metadata only via existing endpoint
		var req UpdateDocumentRequest
		if name != "" {
			req.Name = &name
		}
		req.Note = note
		resp, err := h.service.UpdateDocument(c.Request.Context(), employeeID, documentID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
			return
		}
		httputil.SuccessJSON(c, resp)
		return
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true,
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".txt": true,
	}
	if !allowedExts[ext] {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.file_type_not_allowed_short")
		return
	}
	if file.Size > 10*1024*1024 {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "employee.file_max_10mb")
		return
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)
	filename := employeeID + "_" + timestamp + ext
	uploadDir := "uploads/documents"
	os.MkdirAll(uploadDir, 0755)
	destPath := uploadDir + "/" + filename
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		httputil.InternalError(c, "Failed to save uploaded file")
		return
	}

	filePath := "/" + destPath

	if err := h.service.UpdateDocumentFile(c.Request.Context(), documentID, name, filePath, note); err != nil {
		os.Remove(destPath)
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, gin.H{"file": filePath})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	var req UpdateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateDocument(c.Request.Context(), c.Param("id"), c.Param("documentId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	if err := h.service.DeleteDocument(c.Request.Context(), c.Param("id"), c.Param("documentId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Insurances
// =========================================================================

func (h *Handler) CreateInsurance(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateInsuranceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateInsurance(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateInsurance(c *gin.Context) {
	var req UpdateInsuranceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateInsurance(c.Request.Context(), c.Param("id"), c.Param("insuranceId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteInsurance(c *gin.Context) {
	if err := h.service.DeleteInsurance(c.Request.Context(), c.Param("id"), c.Param("insuranceId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Banks
// =========================================================================

func (h *Handler) CreateBank(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateBankRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateBank(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateBank(c *gin.Context) {
	var req UpdateBankRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateBank(c.Request.Context(), c.Param("id"), c.Param("bankId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteBank(c *gin.Context) {
	if err := h.service.DeleteBank(c.Request.Context(), c.Param("id"), c.Param("bankId")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sub-module Handlers: Employments
// =========================================================================

func (h *Handler) CreateEmployment(c *gin.Context) {
	employeeID := c.Param("id")
	var req CreateEmploymentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployment(c.Request.Context(), employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateEmployment(c *gin.Context) {
	var req UpdateEmploymentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmployment(c.Request.Context(), c.Param("id"), c.Param("employmentId"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployment(c *gin.Context) {
	if err := h.service.DeleteEmployment(c.Request.Context(), c.Param("id"), c.Param("employmentId")); err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Sensitive Field Settings
// =========================================================================

// ListSensitiveFieldSettings menampilkan daftar field sensitif beserta
// status toggle enkripsinya. GET /employees/settings/sensitive-fields
func (h *Handler) ListSensitiveFieldSettings(c *gin.Context) {
	settings, err := h.service.ListSensitiveFieldSettings(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, settings)
}

type setSensitiveFieldEnabledRequest struct {
	IsEncryptionEnabled bool `json:"is_encryption_enabled"`
}

// SetSensitiveFieldEnabled mengubah toggle enkripsi satu field.
// PUT /employees/settings/sensitive-fields/:fieldKey
func (h *Handler) SetSensitiveFieldEnabled(c *gin.Context) {
	fieldKey := c.Param("fieldKey")
	var req setSensitiveFieldEnabledRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	if err := h.service.SetSensitiveFieldEnabled(c.Request.Context(), fieldKey, req.IsEncryptionEnabled); err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"field_key": fieldKey, "is_encryption_enabled": req.IsEncryptionEnabled})
}
