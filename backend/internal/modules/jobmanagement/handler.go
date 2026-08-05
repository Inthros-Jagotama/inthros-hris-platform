package jobmanagement

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

// =========================================================================
// Job Titles (9.1)
// =========================================================================

func (h *Handler) CreateJobTitle(c *gin.Context) {
	var req CreateJobTitleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobTitle(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobTitleByID(c *gin.Context) {
	resp, err := h.service.GetJobTitleByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobTitles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobTitles(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobTitle(c *gin.Context) {
	var req UpdateJobTitleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobTitle(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobTitle(c *gin.Context) {
	if err := h.service.DeleteJobTitle(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Title Subs (9.2)
// =========================================================================

func (h *Handler) CreateJobTitleSub(c *gin.Context) {
	var req CreateJobTitleSubRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobTitleSub(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobTitleSubByID(c *gin.Context) {
	resp, err := h.service.GetJobTitleSubByID(c.Request.Context(), c.Param("subId"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobTitleSubs(c *gin.Context) {
	resp, err := h.service.ListJobTitleSubs(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateJobTitleSub(c *gin.Context) {
	var req UpdateJobTitleSubRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobTitleSub(c.Request.Context(), c.Param("subId"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobTitleSub(c *gin.Context) {
	if err := h.service.DeleteJobTitleSub(c.Request.Context(), c.Param("subId")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Values (9.3)
// =========================================================================

func (h *Handler) CreateJobValue(c *gin.Context) {
	var req CreateJobValueRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobValue(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobValueByID(c *gin.Context) {
	resp, err := h.service.GetJobValueByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ListJobValuesTree — tree data job management values (type_group → type → options)
func (h *Handler) ListJobValuesTree(c *gin.Context) {
	resp, err := h.service.ListJobValuesTree(c.Request.Context())
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListJobValues(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobValues(c.Request.Context(), page, perPage, c.Query("type"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListJobValueClusters(c *gin.Context) {
	resp, err := h.service.ListJobValueClusters(c.Request.Context(), c.Param("type"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateJobValueClusters(c *gin.Context) {
	var req UpdateJobValueClustersRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobValueClusters(c.Request.Context(), c.Param("type"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateJobValue(c *gin.Context) {
	var req UpdateJobValueRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobValue(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobValue(c *gin.Context) {
	if err := h.service.DeleteJobValue(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Objectives (9.4)
// =========================================================================

func (h *Handler) CreateJobObjective(c *gin.Context) {
	var req CreateJobObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobObjective(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobObjectiveByID(c *gin.Context) {
	resp, err := h.service.GetJobObjectiveByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobObjectives(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobObjectives(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobObjective(c *gin.Context) {
	var req UpdateJobObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobObjective(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobObjective(c *gin.Context) {
	if err := h.service.DeleteJobObjective(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Identifications (9.5)
// =========================================================================

func (h *Handler) CreateJobIdentification(c *gin.Context) {
	var req CreateJobIdentificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobIdentification(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobIdentificationByID(c *gin.Context) {
	resp, err := h.service.GetJobIdentificationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobIdentifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobIdentifications(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobIdentification(c *gin.Context) {
	var req UpdateJobIdentificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobIdentification(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobIdentification(c *gin.Context) {
	if err := h.service.DeleteJobIdentification(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Responsibilities (9.6)
// =========================================================================

func (h *Handler) CreateJobResponsibility(c *gin.Context) {
	var req CreateJobResponsibilityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobResponsibility(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobResponsibilityByID(c *gin.Context) {
	resp, err := h.service.GetJobResponsibilityByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobResponsibilities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobResponsibilities(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobResponsibility(c *gin.Context) {
	var req UpdateJobResponsibilityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobResponsibility(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobResponsibility(c *gin.Context) {
	if err := h.service.DeleteJobResponsibility(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Education Experiences (9.7)
// =========================================================================

func (h *Handler) CreateJobEducationExperience(c *gin.Context) {
	var req CreateJobEducationExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobEducationExperience(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobEducationExperienceByID(c *gin.Context) {
	resp, err := h.service.GetJobEducationExperienceByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobEducationExperiences(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobEducationExperiences(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobEducationExperience(c *gin.Context) {
	var req UpdateJobEducationExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobEducationExperience(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobEducationExperience(c *gin.Context) {
	if err := h.service.DeleteJobEducationExperience(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job HR Authorities (9.8)
// =========================================================================

func (h *Handler) CreateJobHRAuthority(c *gin.Context) {
	var req CreateJobHRAuthorityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobHRAuthority(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobHRAuthorityByID(c *gin.Context) {
	resp, err := h.service.GetJobHRAuthorityByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobHRAuthorities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobHRAuthorities(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobHRAuthority(c *gin.Context) {
	var req UpdateJobHRAuthorityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobHRAuthority(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobHRAuthority(c *gin.Context) {
	if err := h.service.DeleteJobHRAuthority(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Operational Authorities (9.9)
// =========================================================================

func (h *Handler) CreateJobOperationalAuthority(c *gin.Context) {
	var req CreateJobOperationalAuthorityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobOperationalAuthority(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobOperationalAuthorityByID(c *gin.Context) {
	resp, err := h.service.GetJobOperationalAuthorityByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobOperationalAuthorities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobOperationalAuthorities(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobOperationalAuthority(c *gin.Context) {
	var req UpdateJobOperationalAuthorityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobOperationalAuthority(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobOperationalAuthority(c *gin.Context) {
	if err := h.service.DeleteJobOperationalAuthority(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Working Activities (9.10)
// =========================================================================

func (h *Handler) CreateJobWorkingActivity(c *gin.Context) {
	var req CreateJobWorkingActivityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobWorkingActivity(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobWorkingActivityByID(c *gin.Context) {
	resp, err := h.service.GetJobWorkingActivityByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobWorkingActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobWorkingActivities(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobWorkingActivity(c *gin.Context) {
	var req UpdateJobWorkingActivityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobWorkingActivity(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobWorkingActivity(c *gin.Context) {
	if err := h.service.DeleteJobWorkingActivity(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Working Risks (9.11)
// =========================================================================

func (h *Handler) CreateJobWorkingRisk(c *gin.Context) {
	var req CreateJobWorkingRiskRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobWorkingRisk(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobWorkingRiskByID(c *gin.Context) {
	resp, err := h.service.GetJobWorkingRiskByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobWorkingRisks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobWorkingRisks(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobWorkingRisk(c *gin.Context) {
	var req UpdateJobWorkingRiskRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobWorkingRisk(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobWorkingRisk(c *gin.Context) {
	if err := h.service.DeleteJobWorkingRisk(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Relationships (9.12)
// =========================================================================

func (h *Handler) CreateJobRelationship(c *gin.Context) {
	var req CreateJobRelationshipRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobRelationship(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobRelationshipByID(c *gin.Context) {
	resp, err := h.service.GetJobRelationshipByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobRelationships(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobRelationships(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobRelationship(c *gin.Context) {
	var req UpdateJobRelationshipRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobRelationship(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobRelationship(c *gin.Context) {
	if err := h.service.DeleteJobRelationship(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Relationship Details (9.12b)
// =========================================================================

func (h *Handler) CreateJobRelationshipDetail(c *gin.Context) {
	var req CreateJobRelationshipDetailRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobRelationshipDetail(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobRelationshipDetailByID(c *gin.Context) {
	resp, err := h.service.GetJobRelationshipDetailByID(c.Request.Context(), c.Param("detailId"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobRelationshipDetails(c *gin.Context) {
	resp, err := h.service.ListJobRelationshipDetails(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateJobRelationshipDetail(c *gin.Context) {
	var req UpdateJobRelationshipDetailRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobRelationshipDetail(c.Request.Context(), c.Param("detailId"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobRelationshipDetail(c *gin.Context) {
	if err := h.service.DeleteJobRelationshipDetail(c.Request.Context(), c.Param("detailId")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Subordinate Controls (9.13)
// =========================================================================

func (h *Handler) CreateJobSubordinateControl(c *gin.Context) {
	var req CreateJobSubordinateControlRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobSubordinateControl(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobSubordinateControlByID(c *gin.Context) {
	resp, err := h.service.GetJobSubordinateControlByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobSubordinateControls(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobSubordinateControls(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobSubordinateControl(c *gin.Context) {
	var req UpdateJobSubordinateControlRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobSubordinateControl(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobSubordinateControl(c *gin.Context) {
	if err := h.service.DeleteJobSubordinateControl(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Assets (9.14)
// =========================================================================

func (h *Handler) CreateJobAsset(c *gin.Context) {
	var req CreateJobAssetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobAsset(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobAssetByID(c *gin.Context) {
	resp, err := h.service.GetJobAssetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobAssets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobAssets(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobAsset(c *gin.Context) {
	var req UpdateJobAssetRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobAsset(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobAsset(c *gin.Context) {
	if err := h.service.DeleteJobAsset(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Financials (9.15)
// =========================================================================

func (h *Handler) CreateJobFinancial(c *gin.Context) {
	var req CreateJobFinancialRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobFinancial(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobFinancialByID(c *gin.Context) {
	resp, err := h.service.GetJobFinancialByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobFinancials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobFinancials(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobFinancial(c *gin.Context) {
	var req UpdateJobFinancialRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobFinancial(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobFinancial(c *gin.Context) {
	if err := h.service.DeleteJobFinancial(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Potency Competencies (9.16)
// =========================================================================

func (h *Handler) CreateJobPotencyCompetency(c *gin.Context) {
	var req CreateJobPotencyCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobPotencyCompetency(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobPotencyCompetencyByID(c *gin.Context) {
	resp, err := h.service.GetJobPotencyCompetencyByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobPotencyCompetencies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobPotencyCompetencies(c.Request.Context(), page, perPage, c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateJobPotencyCompetency(c *gin.Context) {
	var req UpdateJobPotencyCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobPotencyCompetency(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobPotencyCompetency(c *gin.Context) {
	if err := h.service.DeleteJobPotencyCompetency(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Job Scores (9.17)
// =========================================================================

func (h *Handler) UpsertJobScore(c *gin.Context) {
	var req UpdateJobScoreRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpsertJobScore(c.Request.Context(), c.Param("orgId"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetJobScoreByOrganization(c *gin.Context) {
	resp, err := h.service.GetJobScoreByOrganization(c.Request.Context(), c.Param("orgId"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobScores(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListJobScores(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Job Competency Groups (9.18)
// =========================================================================

func (h *Handler) CreateJobCompetencyGroup(c *gin.Context) {
	var req CreateJobCompetencyGroupRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateJobCompetencyGroup(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetJobCompetencyGroupByID(c *gin.Context) {
	resp, err := h.service.GetJobCompetencyGroupByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListJobCompetencyGroups(c *gin.Context) {
	resp, err := h.service.ListJobCompetencyGroups(c.Request.Context(), c.Query("organization_id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateJobCompetencyGroup(c *gin.Context) {
	var req UpdateJobCompetencyGroupRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateJobCompetencyGroup(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteJobCompetencyGroup(c *gin.Context) {
	if err := h.service.DeleteJobCompetencyGroup(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
