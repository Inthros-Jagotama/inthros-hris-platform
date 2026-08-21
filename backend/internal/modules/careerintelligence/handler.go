package careerintelligence

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *Service
	logger *zap.Logger
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc, logger: svc.logger}
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	return page, perPage
}

func respondError(c *gin.Context, status int, msg string) {
	httputil.ErrorRaw(c, status, http.StatusText(status), msg)
}

func respondSuccess(c *gin.Context, data interface{}) {
	httputil.SuccessJSON(c, data)
}

// =========================================================================
// Talent Map Settings
// =========================================================================

func (h *Handler) GetTalentMapSettings(c *gin.Context) {
	result, err := h.svc.GetTalentMapSettings(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) UpdateTalentMapSettings(c *gin.Context) {
	var req UpdateTalentMapSettingsRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.UpdateTalentMapSettings(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, result)
}

// =========================================================================
// Talent Maps
// =========================================================================

func (h *Handler) GenerateTalentMap(c *gin.Context) {
	var req GenerateTalentMapRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.GenerateTalentMap(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, result, "success.created")
}

func (h *Handler) ListTalentMaps(c *gin.Context) {
	page, perPage := parsePagination(c)
	result, err := h.svc.ListTalentMaps(c.Request.Context(),
		c.Query("period"), c.Query("employee_id"), page, perPage)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateTalentMap(c *gin.Context) {
	var req CreateTalentMapRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.CreateTalentMap(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, result, "success.created")
}

func (h *Handler) GetTalentMapByID(c *gin.Context) {
	result, err := h.svc.GetTalentMapByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) UpdateTalentMap(c *gin.Context) {
	var req UpdateTalentMapRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.UpdateTalentMap(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) DeleteTalentMap(c *gin.Context) {
	if err := h.svc.DeleteTalentMap(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) GetTalentGrid(c *gin.Context) {
	result, err := h.svc.GetTalentGrid(c.Request.Context(), c.Query("period"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) GetEmployeeTalentProfile(c *gin.Context) {
	result, err := h.svc.GetEmployeeTalentProfile(c.Request.Context(), c.Param("employeeId"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

// =========================================================================
// Career Interests
// =========================================================================

func (h *Handler) ListCareerInterests(c *gin.Context) {
	page, perPage := parsePagination(c)
	result, err := h.svc.ListCareerInterests(c.Request.Context(),
		c.Query("employee_id"), page, perPage)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateCareerInterest(c *gin.Context) {
	var req CreateCareerInterestRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.CreateCareerInterest(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, result, "success.created")
}

func (h *Handler) GetEmployeeCareerInterests(c *gin.Context) {
	result, err := h.svc.GetEmployeeCareerInterests(c.Request.Context(), c.Param("employeeId"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

// =========================================================================
// Career Paths
// =========================================================================

func (h *Handler) ListCareerPaths(c *gin.Context) {
	page, perPage := parsePagination(c)
	keyword := c.Query("keyword")
	result, err := h.svc.ListCareerPaths(c.Request.Context(), page, perPage, keyword)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateCareerPathLadder menangani POST /career-intelligence/paths dengan
// body ladder-style (name + steps[]) — endpoint utama FE Career Paths.
func (h *Handler) CreateCareerPathLadder(c *gin.Context) {
	var req CreateCareerPathLadderRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.CreateCareerPathLadder(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, result, "success.created")
}

// GetCareerPathByID menangani GET /career-intelligence/paths/:id
func (h *Handler) GetCareerPathByID(c *gin.Context) {
	result, err := h.svc.GetCareerPathByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(c, result)
}

// UpdateCareerPath menangani PUT /career-intelligence/paths/:id (full-replace)
func (h *Handler) UpdateCareerPath(c *gin.Context) {
	var req UpdateCareerPathRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.UpdateCareerPath(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.UpdatedJSON(c, result, "success.updated")
}

func (h *Handler) DeleteCareerPath(c *gin.Context) {
	if err := h.svc.DeleteCareerPath(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// GetEligibleEmployees menangani GET /career-intelligence/paths/:id/eligible-employees
// (S-4 — internal candidate eligibility via career path).
func (h *Handler) GetEligibleEmployees(c *gin.Context) {
	result, err := h.svc.GetEligibleEmployeesForPath(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) GetGapAnalysis(c *gin.Context) {
	req := GapAnalysisRequest{
		EmployeeID:    c.Query("employee_id"),
		TargetTitleID: c.Query("target_title_id"),
	}
	if req.EmployeeID == "" || req.TargetTitleID == "" {
		respondError(c, http.StatusBadRequest, "employee_id and target_title_id are required")
		return
	}
	result, err := h.svc.GetGapAnalysis(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

// =========================================================================
// Succession Plans
// =========================================================================

func (h *Handler) ListSuccessionPlans(c *gin.Context) {
	page, perPage := parsePagination(c)
	result, err := h.svc.ListSuccessionPlans(c.Request.Context(), page, perPage)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetSuccessionGaps (S-5) menandai posisi kunci tanpa successor siap — basis
// fallback external recruitment.
func (h *Handler) GetSuccessionGaps(c *gin.Context) {
	result, err := h.svc.GetSuccessionGaps(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) CreateSuccessionPlan(c *gin.Context) {
	var req CreateSuccessionPlanRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.CreateSuccessionPlan(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.CreatedJSON(c, result, "success.created")
}

func (h *Handler) GetSuccessionPlanByID(c *gin.Context) {
	result, err := h.svc.GetSuccessionPlanByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) UpdateSuccessionPlan(c *gin.Context) {
	var req UpdateSuccessionPlanRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	result, err := h.svc.UpdateSuccessionPlan(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	respondSuccess(c, result)
}

func (h *Handler) DeleteSuccessionPlan(c *gin.Context) {
	if err := h.svc.DeleteSuccessionPlan(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
