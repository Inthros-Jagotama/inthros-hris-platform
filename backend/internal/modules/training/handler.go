package training

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Training Categories
// =========================================================================

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateTrainingCategoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCategories(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListCategories(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training category not found" {
			httputil.NotFound(c, "Training category not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingCategoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training category not found" {
			httputil.NotFound(c, "Training category not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Courses
// =========================================================================

func (h *Handler) CreateCourse(c *gin.Context) {
	var req CreateTrainingCourseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCourse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCourses(c *gin.Context) {
	page, perPage := parsePagination(c)
	catID := c.Query("category_id")
	var catPtr *string
	if catID != "" {
		catPtr = &catID
	}
	resp, err := h.svc.ListCourses(c.Request.Context(), catPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training course not found" {
			httputil.NotFound(c, "Training course not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingCourseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCourse(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training course not found" {
			httputil.NotFound(c, "Training course not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCourse(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Sessions
// =========================================================================

func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateTrainingSessionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateSession(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListSessions(c *gin.Context) {
	page, perPage := parsePagination(c)
	courseID := c.Query("course_id")
	status := c.Query("status")
	var coursePtr *string
	if courseID != "" {
		coursePtr = &courseID
	}
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	resp, err := h.svc.ListSessions(c.Request.Context(), coursePtr, statusPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSessionByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetSessionByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training session not found" {
			httputil.NotFound(c, "Training session not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateSession(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingSessionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateSession(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training session not found" {
			httputil.NotFound(c, "Training session not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateSessionStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSessionStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateSessionStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Participants
// =========================================================================

func (h *Handler) CreateParticipant(c *gin.Context) {
	var req CreateTrainingParticipantRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateParticipant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListParticipants(c *gin.Context) {
	page, perPage := parsePagination(c)
	sessionID := c.Query("session_id")
	employeeID := c.Query("employee_id")
	var sessPtr, empPtr *string
	if sessionID != "" {
		sessPtr = &sessionID
	}
	if employeeID != "" {
		empPtr = &employeeID
	}
	resp, err := h.svc.ListParticipants(c.Request.Context(), sessPtr, empPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetParticipantByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetParticipantByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training participant not found" {
			httputil.NotFound(c, "Training participant not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateParticipant(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingParticipantRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateParticipant(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteParticipant(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteParticipant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Materials
// =========================================================================

func (h *Handler) CreateMaterial(c *gin.Context) {
	var req CreateTrainingMaterialRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateMaterial(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListMaterials(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": "session_id is required"}})
		return
	}
	items, err := h.svc.ListMaterials(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) UpdateMaterial(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingMaterialRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateMaterial(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training material not found" {
			httputil.NotFound(c, "Training material not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteMaterial(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Evaluations
// =========================================================================

func (h *Handler) CreateEvaluation(c *gin.Context) {
	var req CreateTrainingEvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEvaluation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEvaluations(c *gin.Context) {
	page, perPage := parsePagination(c)
	sessionID := c.Query("session_id")
	employeeID := c.Query("employee_id")
	var sessPtr, empPtr *string
	if sessionID != "" {
		sessPtr = &sessionID
	}
	if employeeID != "" {
		empPtr = &employeeID
	}
	resp, err := h.svc.ListEvaluations(c.Request.Context(), sessPtr, empPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEvaluationByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetEvaluationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training evaluation not found" {
			httputil.NotFound(c, "Training evaluation not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateEvaluation(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingEvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluation(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training evaluation not found" {
			httputil.NotFound(c, "Training evaluation not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteEvaluation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEvaluation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Certificates
// =========================================================================

func (h *Handler) CreateCertificate(c *gin.Context) {
	var req CreateTrainingCertificateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCertificate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCertificates(c *gin.Context) {
	page, perPage := parsePagination(c)
	participantID := c.Query("participant_id")
	var partPtr *string
	if participantID != "" {
		partPtr = &participantID
	}
	resp, err := h.svc.ListCertificates(c.Request.Context(), partPtr, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCertificateByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetCertificateByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training certificate not found" {
			httputil.NotFound(c, "Training certificate not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateCertificate(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingCertificateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCertificate(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training certificate not found" {
			httputil.NotFound(c, "Training certificate not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCertificate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCertificate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
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
