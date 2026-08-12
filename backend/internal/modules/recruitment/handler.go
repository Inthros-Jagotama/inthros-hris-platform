package recruitment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Job Requisitions
// =========================================================================

func (h *Handler) CreateRequisition(c *gin.Context) {
	var req CreateRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateRequisition(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRequisitions(c *gin.Context) {
	page, perPage := parsePagination(c)
	orgID := c.Query("organization_id")
	status := c.Query("status")
	var orgPtr *string
	if orgID != "" {
		orgPtr = &orgID
	}
	resp, err := h.svc.ListRequisitions(c.Request.Context(), orgPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRequisitionByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetRequisitionByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateRequisition(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateRequisition(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteRequisition(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRequisition(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// SubmitRequisition menangani POST /api/v1/tenant/recruitment/requisitions/:id/submit
// Routes the requisition through the central approval module (plan G-1).
func (h *Handler) SubmitRequisition(c *gin.Context) {
	id := c.Param("id")
	var req SubmitRequisitionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	flowID := ""
	if req.FlowID != nil {
		flowID = *req.FlowID
	}
	resp, err := h.svc.SubmitRequisition(c.Request.Context(), id, flowID)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their submission didn't reach
		// an approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Candidates
// =========================================================================

func (h *Handler) CreateCandidate(c *gin.Context) {
	var req CreateCandidateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidates(c *gin.Context) {
	page, perPage := parsePagination(c)
	search := c.Query("search")
	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}
	resp, err := h.svc.ListCandidates(c.Request.Context(), searchPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCandidateByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCandidateByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidate(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidate(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// GetEligibleInternalCandidates menangani GET /recruitment/eligible-internal-candidates
// (plan S-4) — membaca employee internal yang eligible untuk target position
// dari Career Intelligence (CI menentukan eligibility; Recruitment hanya
// mengeksekusi aplikasi internal).
func (h *Handler) GetEligibleInternalCandidates(c *gin.Context) {
	positionID := c.Query("position_id")
	if positionID == "" {
		httputil.BadRequest(c, "position_id is required")
		return
	}
	result, err := h.svc.GetEligibleInternalCandidates(c.Request.Context(), positionID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, result)
}

// =========================================================================
// Job Applications
// =========================================================================

func (h *Handler) CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateApplication(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListApplications(c *gin.Context) {
	page, perPage := parsePagination(c)
	reqID := c.Query("requisition_id")
	candID := c.Query("candidate_id")
	status := c.Query("status")
	var reqPtr, candPtr *string
	if reqID != "" {
		reqPtr = &reqID
	}
	if candID != "" {
		candPtr = &candID
	}
	resp, err := h.svc.ListApplications(c.Request.Context(), reqPtr, candPtr, &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetApplicationByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetApplicationByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateApplicationStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateApplicationStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateApplicationStatus(c.Request.Context(), id, req.Status, req.RejectionReason, req.Notes)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteApplication(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteApplication(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Interviews
// =========================================================================

func (h *Handler) CreateInterview(c *gin.Context) {
	var req CreateInterviewRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateInterview(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListInterviews(c *gin.Context) {
	page, perPage := parsePagination(c)
	appID := c.Query("application_id")
	intvID := c.Query("interviewer_id")
	var appPtr, intvPtr *string
	if appID != "" {
		appPtr = &appID
	}
	if intvID != "" {
		intvPtr = &intvID
	}
	resp, err := h.svc.ListInterviews(c.Request.Context(), appPtr, intvPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetInterviewByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetInterviewByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateInterview(c *gin.Context) {
	id := c.Param("id")
	var req UpdateInterviewRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateInterview(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteInterview(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteInterview(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Onboarding Task Templates
// =========================================================================

func (h *Handler) CreateOnboardingTaskTemplate(c *gin.Context) {
	var req CreateOnboardingTaskTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateOnboardingTaskTemplate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListOnboardingTaskTemplates(c *gin.Context) {
	page, perPage := parsePagination(c)
	category := c.Query("category")
	var catPtr *string
	if category != "" {
		catPtr = &category
	}
	resp, err := h.svc.ListOnboardingTaskTemplates(c.Request.Context(), catPtr, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateOnboardingTaskTemplate(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOnboardingTaskTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateOnboardingTaskTemplate(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOnboardingTaskTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOnboardingTaskTemplate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Employee Onboarding
// =========================================================================

func (h *Handler) CreateEmployeeOnboarding(c *gin.Context) {
	var req CreateEmployeeOnboardingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEmployeeOnboarding(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEmployeeOnboardings(c *gin.Context) {
	page, perPage := parsePagination(c)
	status := c.Query("status")
	resp, err := h.svc.ListEmployeeOnboardings(c.Request.Context(), &status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEmployeeOnboardingByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEmployeeOnboardingByID(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEmployeeOnboarding(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEmployeeOnboardingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEmployeeOnboarding(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeOnboarding(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEmployeeOnboarding(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Onboarding Task Items
// =========================================================================

func (h *Handler) CreateOnboardingTaskItem(c *gin.Context) {
	var req CreateOnboardingTaskItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateOnboardingTaskItem(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListOnboardingTaskItems(c *gin.Context) {
	onboardingID := c.Param("id")
	items, err := h.svc.ListOnboardingTaskItems(c.Request.Context(), onboardingID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, items)
}

func (h *Handler) UpdateOnboardingTaskItem(c *gin.Context) {
	id := c.Param("id")
	var req UpdateOnboardingTaskItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateOnboardingTaskItem(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteOnboardingTaskItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteOnboardingTaskItem(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Helpers
// =========================================================================

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	perPage := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if v, err := strconv.Atoi(pp); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	return page, perPage
}
