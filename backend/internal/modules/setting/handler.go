package setting

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
	"github.com/inthros/hris-platform/internal/pkg/numbering"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// handleDupErr checks if the error is a DuplicateCodeError and sends a
// bilingual conflict response. Returns true if the error was handled.
func handleDupErr(c *gin.Context, err error) bool {
	var dupErr *DuplicateCodeError
	if errors.As(err, &dupErr) {
		httputil.ErrorJSON(c, 409, "DUPLICATE_CODE", "setting.duplicate_code", dupErr.Code)
		return true
	}
	return false
}

// ── Company Timezone Handlers ──
func (h *Handler) GetCompanyTimezone(c *gin.Context) {
	tz, err := h.service.GetCompanyTimezone(c.Request.Context())
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"timezone": tz})
}

type updateCompanyTimezoneRequest struct {
	Timezone string `json:"timezone" binding:"required"`
}

func (h *Handler) UpdateCompanyTimezone(c *gin.Context) {
	var req updateCompanyTimezoneRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	if err := h.service.UpdateCompanyTimezone(c.Request.Context(), req.Timezone); err != nil {
		if errors.Is(err, ErrInvalidCompanyTimezone) {
			httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"timezone": req.Timezone})
}

// ── Competency Handlers ──
func (h *Handler) CreateCompetency(c *gin.Context) {
	var req CreateCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateCompetency(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListCompetencies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.DefaultQuery("search", "")
	resp, err := h.service.ListCompetencies(c.Request.Context(), page, perPage, search)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetCompetencyByID(c *gin.Context) {
	resp, err := h.service.GetCompetencyByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateCompetency(c *gin.Context) {
	var req UpdateCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateCompetency(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteCompetency(c *gin.Context) {
	if err := h.service.DeleteCompetency(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Zone Handlers ──
func (h *Handler) CreateZone(c *gin.Context) {
	var req CreateZoneRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateZone(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListZones(c *gin.Context) {
	if c.Query("active_only") == "true" {
		zones, err := h.service.ListZonesActive(c.Request.Context())
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, zones)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListZones(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetZoneByID(c *gin.Context) {
	resp, err := h.service.GetZoneByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateZone(c *gin.Context) {
	var req UpdateZoneRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateZone(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteZone(c *gin.Context) {
	if err := h.service.DeleteZone(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Province Handlers ──
func (h *Handler) CreateProvince(c *gin.Context) {
	var req CreateProvinceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateProvince(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListProvinces(c *gin.Context) {
	if c.Query("all") == "true" {
		provinces, err := h.service.ListAllProvinces(c.Request.Context())
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, provinces)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListProvinces(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetProvinceByID(c *gin.Context) {
	resp, err := h.service.GetProvinceByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateProvince(c *gin.Context) {
	var req UpdateProvinceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateProvince(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteProvince(c *gin.Context) {
	if err := h.service.DeleteProvince(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Regency Handlers ──
func (h *Handler) CreateRegency(c *gin.Context) {
	var req CreateRegencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateRegency(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListRegencies(c *gin.Context) {
	if provID := c.Query("province_id"); provID != "" {
		regencies, err := h.service.ListRegenciesByProvince(c.Request.Context(), provID)
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, regencies)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListRegencies(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetRegencyByID(c *gin.Context) {
	resp, err := h.service.GetRegencyByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateRegency(c *gin.Context) {
	var req UpdateRegencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateRegency(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteRegency(c *gin.Context) {
	if err := h.service.DeleteRegency(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── District Handlers ──
func (h *Handler) CreateDistrict(c *gin.Context) {
	var req CreateDistrictRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateDistrict(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListDistricts(c *gin.Context) {
	if regID := c.Query("regency_id"); regID != "" {
		districts, err := h.service.ListDistrictsByRegency(c.Request.Context(), regID)
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, districts)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListDistricts(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetDistrictByID(c *gin.Context) {
	resp, err := h.service.GetDistrictByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateDistrict(c *gin.Context) {
	var req UpdateDistrictRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateDistrict(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteDistrict(c *gin.Context) {
	if err := h.service.DeleteDistrict(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Village Handlers ──
func (h *Handler) CreateVillage(c *gin.Context) {
	var req CreateVillageRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateVillage(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListVillages(c *gin.Context) {
	if distID := c.Query("district_id"); distID != "" {
		villages, err := h.service.ListVillagesByDistrict(c.Request.Context(), distID)
		if err != nil {
			httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
			return
		}
		httputil.SuccessJSON(c, villages)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListVillages(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetVillageByID(c *gin.Context) {
	resp, err := h.service.GetVillageByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateVillage(c *gin.Context) {
	var req UpdateVillageRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateVillage(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteVillage(c *gin.Context) {
	if err := h.service.DeleteVillage(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
func (h *Handler) GetVillageDetail(c *gin.Context) {
	resp, err := h.service.GetVillageDetail(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SearchVillages(c *gin.Context) {
	query := c.Query("q")
	resp, err := h.service.SearchVillages(c.Request.Context(), query)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ── Education Handlers ──
func (h *Handler) CreateEducation(c *gin.Context) {
	var req CreateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEducation(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListEducations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEducations(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetEducationByID(c *gin.Context) {
	resp, err := h.service.GetEducationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateEducation(c *gin.Context) {
	var req UpdateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEducation(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteEducation(c *gin.Context) {
	if err := h.service.DeleteEducation(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── EducationMajor Handlers ──
func (h *Handler) CreateEducationMajor(c *gin.Context) {
	var req CreateEducationMajorRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEducationMajor(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListEducationMajors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEducationMajors(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetEducationMajorByID(c *gin.Context) {
	resp, err := h.service.GetEducationMajorByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateEducationMajor(c *gin.Context) {
	var req UpdateEducationMajorRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEducationMajor(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteEducationMajor(c *gin.Context) {
	if err := h.service.DeleteEducationMajor(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Religion Handlers ──
func (h *Handler) CreateReligion(c *gin.Context) {
	var req CreateReligionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateReligion(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListReligions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListReligions(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetReligionByID(c *gin.Context) {
	resp, err := h.service.GetReligionByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateReligion(c *gin.Context) {
	var req UpdateReligionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateReligion(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteReligion(c *gin.Context) {
	if err := h.service.DeleteReligion(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── MaritalStatus Handlers ──
func (h *Handler) CreateMaritalStatus(c *gin.Context) {
	var req CreateMaritalStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateMaritalStatus(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListMaritalStatuses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListMaritalStatuses(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetMaritalStatusByID(c *gin.Context) {
	resp, err := h.service.GetMaritalStatusByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateMaritalStatus(c *gin.Context) {
	var req UpdateMaritalStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateMaritalStatus(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteMaritalStatus(c *gin.Context) {
	if err := h.service.DeleteMaritalStatus(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── RelationshipType Handlers ──
func (h *Handler) CreateRelationshipType(c *gin.Context) {
	var req CreateRelationshipTypeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateRelationshipType(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListRelationshipTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListRelationshipTypes(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetRelationshipTypeByID(c *gin.Context) {
	resp, err := h.service.GetRelationshipTypeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateRelationshipType(c *gin.Context) {
	var req UpdateRelationshipTypeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateRelationshipType(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteRelationshipType(c *gin.Context) {
	if err := h.service.DeleteRelationshipType(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── EmploymentStatus Handlers ──
func (h *Handler) CreateEmploymentStatus(c *gin.Context) {
	var req CreateEmploymentStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmploymentStatus(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListEmploymentStatuses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEmploymentStatuses(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetEmploymentStatusByID(c *gin.Context) {
	resp, err := h.service.GetEmploymentStatusByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateEmploymentStatus(c *gin.Context) {
	var req UpdateEmploymentStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmploymentStatus(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteEmploymentStatus(c *gin.Context) {
	if err := h.service.DeleteEmploymentStatus(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Bank Handlers ──
func (h *Handler) CreateBank(c *gin.Context) {
	var req CreateBankRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateBank(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListBanks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListBanks(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetBankByID(c *gin.Context) {
	resp, err := h.service.GetBankByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateBank(c *gin.Context) {
	var req UpdateBankRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateBank(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteBank(c *gin.Context) {
	if err := h.service.DeleteBank(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Nationality Handlers ──
func (h *Handler) CreateNationality(c *gin.Context) {
	var req CreateNationalityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateNationality(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListNationalities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.DefaultQuery("search", "")
	resp, err := h.service.ListNationalities(c.Request.Context(), page, perPage, search)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetNationalityByID(c *gin.Context) {
	resp, err := h.service.GetNationalityByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateNationality(c *gin.Context) {
	var req UpdateNationalityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateNationality(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteNationality(c *gin.Context) {
	if err := h.service.DeleteNationality(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── JobFamily Handlers ──
func (h *Handler) CreateJobFamily(c *gin.Context) {
	var req CreateJobFamilyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobFamily(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListJobFamilies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobFamilies(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetJobFamilyByID(c *gin.Context) {
	resp, err := h.service.GetJobFamilyByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateJobFamily(c *gin.Context) {
	var req UpdateJobFamilyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobFamily(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteJobFamily(c *gin.Context) {
	if err := h.service.DeleteJobFamily(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── TER Handlers ──
func (h *Handler) CreateTER(c *gin.Context) {
	var req CreateTERRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateTER(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListTERs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	resp, err := h.service.ListTERs(c.Request.Context(), page, perPage, search)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetTERByID(c *gin.Context) {
	resp, err := h.service.GetTERByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateTER(c *gin.Context) {
	var req UpdateTERRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateTER(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteTER(c *gin.Context) {
	if err := h.service.DeleteTER(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── PTKP Handlers ──
func (h *Handler) CreatePTKP(c *gin.Context) {
	var req CreatePTKPRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePTKP(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListPTKPs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	resp, err := h.service.ListPTKPs(c.Request.Context(), page, perPage, search)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetPTKPByID(c *gin.Context) {
	resp, err := h.service.GetPTKPByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdatePTKP(c *gin.Context) {
	var req UpdatePTKPRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePTKP(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeletePTKP(c *gin.Context) {
	if err := h.service.DeletePTKP(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Grading Handlers ──
func (h *Handler) CreateGrading(c *gin.Context) {
	var req CreateGradingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateGrading(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListGradings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListGradings(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetGradingByID(c *gin.Context) {
	resp, err := h.service.GetGradingByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateGrading(c *gin.Context) {
	var req UpdateGradingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateGrading(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteGrading(c *gin.Context) {
	if err := h.service.DeleteGrading(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Insurance Handlers ──
func (h *Handler) CreateInsurance(c *gin.Context) {
	var req CreateInsuranceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateInsurance(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListInsurances(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListInsurances(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetInsuranceByID(c *gin.Context) {
	resp, err := h.service.GetInsuranceByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateInsurance(c *gin.Context) {
	var req UpdateInsuranceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateInsurance(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteInsurance(c *gin.Context) {
	if err := h.service.DeleteInsurance(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── CompanyHoliday Handlers ──
func (h *Handler) CreateCompanyHoliday(c *gin.Context) {
	var req CreateCompanyHolidayRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateCompanyHoliday(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListCompanyHolidays(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListCompanyHolidays(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetCompanyHolidayByID(c *gin.Context) {
	resp, err := h.service.GetCompanyHolidayByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateCompanyHoliday(c *gin.Context) {
	var req UpdateCompanyHolidayRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateCompanyHoliday(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteCompanyHoliday(c *gin.Context) {
	if err := h.service.DeleteCompanyHoliday(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── SalaryGrade Handlers ──
func (h *Handler) CreateSalaryGrade(c *gin.Context) {
	var req CreateSalaryGradeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateSalaryGrade(c.Request.Context(), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}
func (h *Handler) ListSalaryGrades(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListSalaryGrades(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}
func (h *Handler) GetSalaryGradeByID(c *gin.Context) {
	resp, err := h.service.GetSalaryGradeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) UpdateSalaryGrade(c *gin.Context) {
	var req UpdateSalaryGradeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateSalaryGrade(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if handleDupErr(c, err) {
			return
		}
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
func (h *Handler) DeleteSalaryGrade(c *gin.Context) {
	if err := h.service.DeleteSalaryGrade(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// ── Document Numbering Handlers ──

// NumberingHandler exposes CRUD/preview endpoints for document numbering
// settings, backed by numbering.Service.
type NumberingHandler struct {
	svc *numbering.Service
}

func NewNumberingHandler(svc *numbering.Service) *NumberingHandler {
	return &NumberingHandler{svc: svc}
}

func (h *NumberingHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, items)
}

type updateNumberingRequest struct {
	FormatTemplate string `json:"format_template" binding:"required,max=255"`
	ResetPeriod    string `json:"reset_period" binding:"required"`
}

func (h *NumberingHandler) Update(c *gin.Context) {
	documentType := c.Param("document_type")
	var req updateNumberingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	updated, err := h.svc.Update(c.Request.Context(), documentType, req.FormatTemplate, req.ResetPeriod)
	if err != nil {
		switch err {
		case numbering.ErrInvalidDocumentType, numbering.ErrInvalidResetPeriod:
			httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		case numbering.ErrSettingNotFound:
			httputil.ErrorSimple(c, http.StatusNotFound, err.Error())
		default:
			httputil.InternalError(c, err.Error())
		}
		return
	}
	httputil.SuccessJSON(c, updated)
}

func (h *NumberingHandler) Preview(c *gin.Context) {
	documentType := c.Param("document_type")
	preview, err := h.svc.Preview(c.Request.Context(), documentType)
	if err != nil {
		switch err {
		case numbering.ErrInvalidDocumentType:
			httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		case numbering.ErrSettingNotFound:
			httputil.ErrorSimple(c, http.StatusNotFound, err.Error())
		default:
			httputil.InternalError(c, err.Error())
		}
		return
	}
	httputil.SuccessJSON(c, gin.H{"preview": preview})
}
