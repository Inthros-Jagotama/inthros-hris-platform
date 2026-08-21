package training

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
// Training Categories
// =========================================================================

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateTrainingCategoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCategories(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListCategories(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCategory(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCourse(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSession(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteParticipant(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteParticipant(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteMaterial(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteEvaluation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEvaluation(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCertificate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCertificate(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
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
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListProviders(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListProviders(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProvider(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
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
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteTrainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteTrainer(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
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
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListSessionTrainers(c *gin.Context) {
	sessionID := c.Param("id")
	items, err := h.svc.ListSessionTrainers(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) RemoveSessionTrainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RemoveSessionTrainer(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func (h *Handler) MarkAttendance(c *gin.Context) {
	sessionID := c.Param("id")
	var reqs []MarkTrainingAttendanceRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if len(reqs) == 0 {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "training.attendance_items_required")
		return
	}
	resp, err := h.svc.MarkAttendance(c.Request.Context(), sessionID, reqs)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListAssessmentsBySession(c *gin.Context) {
	sessionID := c.Param("id")
	items, err := h.svc.ListAssessmentsBySession(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
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

// =========================================================================
// Training Plans (P1-BE — plan §16)
// =========================================================================

func (h *Handler) CreatePlan(c *gin.Context) {
	var req CreateTrainingPlanRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePlan(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPlans(c *gin.Context) {
	page, perPage := parsePagination(c)
	var year *int
	if y := c.Query("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = &v
		}
	}
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	resp, err := h.svc.ListPlans(c.Request.Context(), year, statusPtr, page, perPage)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetPlanByID(c.Request.Context(), id)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingPlanRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePlan(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePlan(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Plan Items (P1-BE — plan §16)
// =========================================================================

func (h *Handler) CreatePlanItem(c *gin.Context) {
	planID := c.Param("id")
	var req CreateTrainingPlanItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreatePlanItem(c.Request.Context(), planID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPlanItems(c *gin.Context) {
	planID := c.Param("id")
	resp, err := h.svc.ListPlanItems(c.Request.Context(), planID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdatePlanItem(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingPlanItemRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdatePlanItem(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeletePlanItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeletePlanItem(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Needs (P1-BE — plan §17)
// =========================================================================

func (h *Handler) CreateNeed(c *gin.Context) {
	var req CreateTrainingNeedRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateNeed(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListNeeds(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	var empPtr *string
	if empID != "" {
		empPtr = &empID
	}
	courseID := c.Query("course_id")
	var coursePtr *string
	if courseID != "" {
		coursePtr = &courseID
	}
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	resp, err := h.svc.ListNeeds(c.Request.Context(), empPtr, coursePtr, statusPtr, page, perPage)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetNeedByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetNeedByID(c.Request.Context(), id)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateNeed(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingNeedRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateNeed(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Training Requests (P1-BE — plan §15, Central Approval)
// =========================================================================

func (h *Handler) CreateRequest(c *gin.Context) {
	var req CreateTrainingRequestRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateRequest(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListRequests(c *gin.Context) {
	page, perPage := parsePagination(c)
	empID := c.Query("employee_id")
	var empPtr *string
	if empID != "" {
		empPtr = &empID
	}
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}
	resp, err := h.svc.ListRequests(c.Request.Context(), empPtr, statusPtr, page, perPage)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRequestByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetRequestByID(c.Request.Context(), id)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	id := c.Param("id")
	var req SubmitTrainingRequestRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.SubmitRequest(c.Request.Context(), id, req.FlowID)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their submission didn't reach
		// an approver instead of a raw error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) CancelRequest(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.CancelRequest(c.Request.Context(), id)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Course Objectives (P1-BE — plan §8)
// =========================================================================

func (h *Handler) CreateCourseObjective(c *gin.Context) {
	courseID := c.Param("id")
	var req CreateCourseObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCourseObjective(c.Request.Context(), courseID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCourseObjectives(c *gin.Context) {
	courseID := c.Param("id")
	resp, err := h.svc.ListCourseObjectives(c.Request.Context(), courseID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateCourseObjective(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCourseObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCourseObjective(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCourseObjective(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCourseObjective(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Course Competencies (P1-BE — plan §9)
// =========================================================================

func (h *Handler) CreateCourseCompetency(c *gin.Context) {
	courseID := c.Param("id")
	var req CreateCourseCompetencyRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCourseCompetency(c.Request.Context(), courseID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCourseCompetencies(c *gin.Context) {
	courseID := c.Param("id")
	resp, err := h.svc.ListCourseCompetencies(c.Request.Context(), courseID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCourseCompetency(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCourseCompetency(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Course Prerequisites (P1-BE — plan §10)
// =========================================================================

func (h *Handler) CreateCoursePrerequisite(c *gin.Context) {
	courseID := c.Param("id")
	var req CreateCoursePrerequisiteRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCoursePrerequisite(c.Request.Context(), courseID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCoursePrerequisites(c *gin.Context) {
	courseID := c.Param("id")
	resp, err := h.svc.ListCoursePrerequisites(c.Request.Context(), courseID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteCoursePrerequisite(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCoursePrerequisite(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Mandatories (P1-BE — plan §25)
// =========================================================================

func (h *Handler) CreateMandatory(c *gin.Context) {
	var req CreateTrainingMandatoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateMandatory(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListMandatories(c *gin.Context) {
	page, perPage := parsePagination(c)
	courseID := c.Query("course_id")
	var coursePtr *string
	if courseID != "" {
		coursePtr = &courseID
	}
	resp, err := h.svc.ListMandatories(c.Request.Context(), coursePtr, page, perPage)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMandatoryByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetMandatoryByID(c.Request.Context(), id)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateMandatory(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingMandatoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateMandatory(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteNeed(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteNeed(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) DeleteMandatory(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteMandatory(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Session Costs (P1-BE — plan §26)
// =========================================================================

func (h *Handler) CreateSessionCost(c *gin.Context) {
	sessionID := c.Param("id")
	var req CreateTrainingSessionCostRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateSessionCost(c.Request.Context(), sessionID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListSessionCosts(c *gin.Context) {
	sessionID := c.Param("id")
	resp, err := h.svc.ListSessionCosts(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateSessionCost(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTrainingSessionCostRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateSessionCost(c.Request.Context(), id, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteSessionCost(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteSessionCost(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Training Documents (P1-BE — plan §27)
// =========================================================================

func (h *Handler) CreateDocument(c *gin.Context) {
	sessionID := c.Param("id")
	var req CreateTrainingDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateDocument(c.Request.Context(), sessionID, req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListDocuments(c *gin.Context) {
	sessionID := c.Param("id")
	resp, err := h.svc.ListDocuments(c.Request.Context(), sessionID)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteDocument(c.Request.Context(), id); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Evaluation Forms — handler P2 (plan §22)
// =========================================================================

func (h *Handler) CreateEvaluationForm(c *gin.Context) {
	var req CreateEvaluationFormRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEvaluationForm(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetEvaluationForm(c *gin.Context) {
	resp, err := h.svc.GetEvaluationFormByID(c.Request.Context(), c.Param("form_id"))
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListEvaluationForms(c *gin.Context) {
	page, perPage := parsePagination(c)
	var sessionID *string
	if v := c.Query("session_id"); v != "" {
		sessionID = &v
	}
	resp, err := h.svc.ListEvaluationForms(c.Request.Context(), sessionID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEvaluationForm(c *gin.Context) {
	var req UpdateEvaluationFormRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationForm(c.Request.Context(), c.Param("form_id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEvaluationForm(c *gin.Context) {
	if err := h.svc.DeleteEvaluationForm(c.Request.Context(), c.Param("form_id")); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// GetEvaluationFormBySession — form + questions untuk session.
func (h *Handler) GetEvaluationFormBySession(c *gin.Context) {
	resp, err := h.svc.GetEvaluationFormBySession(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Evaluation Questions — handler P2 (plan §22)
// =========================================================================

func (h *Handler) CreateEvaluationQuestion(c *gin.Context) {
	var req CreateEvaluationQuestionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEvaluationQuestion(c.Request.Context(), c.Param("form_id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEvaluationQuestions(c *gin.Context) {
	resp, err := h.svc.ListEvaluationQuestions(c.Request.Context(), c.Param("form_id"))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEvaluationQuestion(c *gin.Context) {
	var req UpdateEvaluationQuestionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEvaluationQuestion(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEvaluationQuestion(c *gin.Context) {
	if err := h.svc.DeleteEvaluationQuestion(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Evaluation Answers — handler P2 (plan §22)
// =========================================================================

func (h *Handler) SubmitEvaluationAnswers(c *gin.Context) {
	var req SubmitEvaluationAnswersRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.SubmitEvaluationAnswers(c.Request.Context(), c.Param("form_id"), c.Param("participant_id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListEvaluationAnswers(c *gin.Context) {
	var questionID, participantID *string
	if v := c.Query("question_id"); v != "" {
		questionID = &v
	}
	if v := c.Query("participant_id"); v != "" {
		participantID = &v
	}
	resp, err := h.svc.ListEvaluationAnswers(c.Request.Context(), questionID, participantID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Effectiveness Assessments — handler P2 (plan §23)
// =========================================================================

func (h *Handler) CreateEffectivenessAssessment(c *gin.Context) {
	var req CreateEffectivenessAssessmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateEffectivenessAssessment(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetEffectivenessAssessment(c *gin.Context) {
	resp, err := h.svc.GetEffectivenessAssessmentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListEffectivenessAssessments(c *gin.Context) {
	page, perPage := parsePagination(c)
	var participantID *string
	if v := c.Query("participant_id"); v != "" {
		participantID = &v
	}
	resp, err := h.svc.ListEffectivenessAssessments(c.Request.Context(), participantID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEffectivenessAssessment(c *gin.Context) {
	var req UpdateEffectivenessAssessmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateEffectivenessAssessment(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEffectivenessAssessment(c *gin.Context) {
	if err := h.svc.DeleteEffectivenessAssessment(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Certifications — handler P2 (plan §24)
// =========================================================================

func (h *Handler) CreateCertification(c *gin.Context) {
	var req CreateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCertification(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetCertification(c *gin.Context) {
	resp, err := h.svc.GetCertificationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListCertifications(c *gin.Context) {
	page, perPage := parsePagination(c)
	var isActive *bool
	if v := c.Query("is_active"); v != "" {
		b := v == "true" || v == "1"
		isActive = &b
	}
	resp, err := h.svc.ListCertifications(c.Request.Context(), isActive, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	// Pakai c.JSON langsung (bukan httputil.SuccessJSON) agar response berbentuk
	// envelope PaginatedResponse standar { success, data: [...], page, per_page, ... }
	// — konsisten dengan seluruh handler list training lain (ListCategories,
	// ListCourses, ListCertificates, dll). Sebelumnya SuccessJSON membungkus
	// envelope dua kali sehingga data menjadi objek (bukan array) dan memicu
	// "Maximum recursive updates exceeded" pada DataTable frontend.
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateCertification(c *gin.Context) {
	var req UpdateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCertification(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCertification(c *gin.Context) {
	if err := h.svc.DeleteCertification(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Certificates — generate dari completion (P2-BE — plan §24)
// =========================================================================

func (h *Handler) GenerateCertificate(c *gin.Context) {
	var req GenerateCertificateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.GenerateCertificate(c.Request.Context(), c.Param("participant_id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateCertificateFile(c *gin.Context) {
	var req UpdateCertificateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCertificateFile(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorRaw(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Reports & History — handler P2 (plan §38)
// =========================================================================

func (h *Handler) GetTrainingHistory(c *gin.Context) {
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "training.employee_id_required")
		return
	}
	resp, err := h.svc.GetTrainingHistory(c.Request.Context(), employeeID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetParticipationReport(c *gin.Context) {
	var sessionStatus *string
	if v := c.Query("session_status"); v != "" {
		sessionStatus = &v
	}
	resp, err := h.svc.GetParticipationReport(c.Request.Context(), sessionStatus)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetCostReport(c *gin.Context) {
	resp, err := h.svc.GetCostReport(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetComplianceReport(c *gin.Context) {
	resp, err := h.svc.GetComplianceReport(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetDashboardReport(c *gin.Context) {
	resp, err := h.svc.GetDashboardReport(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
