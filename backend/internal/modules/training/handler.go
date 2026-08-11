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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "training.session_id_required")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
			httputil.NotFound(c, "")
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
// Training Providers (P0-BE — plan §11)
// =========================================================================

func (h *Handler) CreateProvider(c *gin.Context) {
	var req CreateTrainingProviderRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateProvider(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListProviders(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListProviders(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProviderByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetProviderByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training provider not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingProviderRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateProvider(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training provider not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProvider(c.Request.Context(), id); err != nil {
		httputil.ErrorJSON(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Trainers (P0-BE — plan §12)
// =========================================================================

func (h *Handler) CreateTrainer(c *gin.Context) {
	var req CreateTrainingTrainerRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateTrainer(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListTrainers(c *gin.Context) {
	page, perPage := parsePagination(c)
	typ := c.Query("type")
	var typPtr *string
	if typ != "" {
		typPtr = &typ
	}
	resp, err := h.svc.ListTrainers(c.Request.Context(), typPtr, page, perPage)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTrainerByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetTrainerByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "training trainer not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateTrainer(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingTrainerRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateTrainer(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training trainer not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteTrainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteTrainer(c.Request.Context(), id); err != nil {
		httputil.ErrorJSON(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Session Trainers (P0-BE — plan §13)
// =========================================================================

func (h *Handler) AddSessionTrainer(c *gin.Context) {
	sessionID := c.Param("id")
	var req AddSessionTrainerRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.AddSessionTrainer(c.Request.Context(), sessionID, req)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListSessionTrainers(c *gin.Context) {
	sessionID := c.Param("id")
	items, err := h.svc.ListSessionTrainers(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) RemoveSessionTrainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RemoveSessionTrainer(c.Request.Context(), id); err != nil {
		httputil.ErrorJSON(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Attendances (P0-BE — plan §19)
// =========================================================================

func (h *Handler) ListAttendanceBySession(c *gin.Context) {
	sessionID := c.Param("id")
	rows, err := h.svc.ListAttendanceBySession(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	var reqs []MarkTrainingAttendanceRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if len(reqs) == 0 {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "attendance items are required")
		return
	}
	resp, err := h.svc.MarkAttendance(c.Request.Context(), sessionID, reqs)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateAttendance(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingAttendanceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateAttendance(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "training attendance not found" {
			httputil.NotFound(c, "")
			return
		}
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Training Assessments (P0-BE — plan §21)
// =========================================================================

func (h *Handler) CreateAssessment(c *gin.Context) {
	var req CreateTrainingAssessmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateAssessment(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListAssessmentsBySession(c *gin.Context) {
	sessionID := c.Param("id")
	items, err := h.svc.ListAssessmentsBySession(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) SubmitAssessmentResult(c *gin.Context) {
	assessmentID := c.Param("id")
	var req SubmitAssessmentResultRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.SubmitAssessmentResult(c.Request.Context(), assessmentID, req)
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
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
