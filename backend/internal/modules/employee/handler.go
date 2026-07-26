package employee

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
			httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /api/v1/tenant/employees
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	resp, err := h.service.List(c.Request.Context(), page, perPage)
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
