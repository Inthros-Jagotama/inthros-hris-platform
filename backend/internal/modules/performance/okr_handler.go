package performance

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type OKRHandler struct {
	service  OKRService
	resolver TenantDBFunc
}

func NewOKRHandler(service OKRService, resolver TenantDBFunc) *OKRHandler {
	return &OKRHandler{service: service, resolver: resolver}
}

func (h *OKRHandler) getTenantDB(c *gin.Context) *gorm.DB {
	db, err := h.resolver(c.Request.Context())
	if err != nil {
		return nil
	}
	return db
}

func (h *OKRHandler) getUserID(c *gin.Context) uuid.UUID {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		return uuid.Nil
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}
	return userID
}

func (h *OKRHandler) GetMyOKRContext(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "failed to resolve tenant database")
		return
	}
	resp, err := h.service.GetMyOKRContext(db, h.getUserID(c))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func okrParsePagination(c *gin.Context) (int, int) {
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
// OKR Templates
// =========================================================================

func (h *OKRHandler) CreateTemplate(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.CreateTemplate(db, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) GetTemplateByID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetTemplateWithObjectives(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ListTemplates(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	page, perPage := okrParsePagination(c)

	var orgID, periodID *uuid.UUID
	var status *int

	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if id, err := uuid.Parse(orgIDStr); err == nil {
			orgID = &id
		}
	}
	if periodIDStr := c.Query("period_id"); periodIDStr != "" {
		if id, err := uuid.Parse(periodIDStr); err == nil {
			periodID = &id
		}
	}
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	results, total, err := h.service.ListTemplates(db, orgID, periodID, status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"meta": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

func (h *OKRHandler) UpdateTemplate(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKRTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateTemplate(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteTemplate(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteTemplate(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

func (h *OKRHandler) DuplicateTemplate(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.DuplicateTemplate(db, id)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

// =========================================================================
// OKR Objectives
// =========================================================================

func (h *OKRHandler) CreateObjective(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.CreateObjective(db, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) GetObjectiveByID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetObjectiveByID(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ListObjectivesByTemplateID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	templateID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	results, err := h.service.ListObjectivesByTemplateID(db, templateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, results)
}

func (h *OKRHandler) UpdateObjective(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKRObjectiveRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateObjective(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteObjective(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteObjective(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// OKR Key Results
// =========================================================================

func (h *OKRHandler) CreateKeyResult(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRKeyResultRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.CreateKeyResult(db, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) GetKeyResultByID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetKeyResultByID(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ListKeyResultsByObjectiveID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	objectiveID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	results, err := h.service.ListKeyResultsByObjectiveID(db, objectiveID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, results)
}

func (h *OKRHandler) UpdateKeyResult(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKRKeyResultRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateKeyResult(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteKeyResult(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteKeyResult(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// OKR Evaluations
// =========================================================================

func (h *OKRHandler) CreateEvaluationWithSnapshot(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKREvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.CreateEvaluationWithSnapshot(db, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) GetEvaluationByID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetEvaluationByID(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) GetEvaluationWithDetails(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetEvaluationWithDetails(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ListEvaluations(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	page, perPage := okrParsePagination(c)

	var employeeID, orgID, periodID *uuid.UUID
	var status *string

	if empIDStr := c.Query("employee_id"); empIDStr != "" {
		if id, err := uuid.Parse(empIDStr); err == nil {
			employeeID = &id
		}
	}
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if id, err := uuid.Parse(orgIDStr); err == nil {
			orgID = &id
		}
	}
	if periodIDStr := c.Query("period_id"); periodIDStr != "" {
		if id, err := uuid.Parse(periodIDStr); err == nil {
			periodID = &id
		}
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	results, total, err := h.service.ListEvaluations(db, employeeID, orgID, periodID, status, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"meta": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

func (h *OKRHandler) UpdateEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKREvaluationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateEvaluation(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteEvaluation(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Evaluation Detail & Score
// =========================================================================

func (h *OKRHandler) UpdateEvaluationDetailActual(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKREvaluationDetailRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateEvaluationDetailActual(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) BulkUpdateEvaluationActuals(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req OKRBulkUpdateActualsRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	if err := h.service.BulkUpdateEvaluationActuals(db, evaluationID, &req); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.MessageJSON(c, "success.actuals_updated")
}

func (h *OKRHandler) RecalculateEvaluationScore(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.RecalculateEvaluationScore(db, evaluationID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

// =========================================================================
// Workflow
// =========================================================================

func (h *OKRHandler) SubmitEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	userID := h.getUserID(c)

	result, err := h.service.SubmitEvaluation(db, evaluationID, userID)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ApproveEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	userID := h.getUserID(c)

	result, err := h.service.ApproveEvaluation(db, evaluationID, userID)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) RejectEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	userID := h.getUserID(c)

	var req struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&req)

	result, err := h.service.RejectEvaluation(db, evaluationID, userID, req.Notes)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) CompleteEvaluation(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.CompleteEvaluation(db, evaluationID)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

// =========================================================================
// OKR Progress
// =========================================================================

func (h *OKRHandler) CreateProgress(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRProgressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	userID := h.getUserID(c)

	result, err := h.service.CreateProgress(db, &req, userID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) GetProgressByID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	result, err := h.service.GetProgressByID(db, id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) ListProgressByDetailID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	detailID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	results, err := h.service.ListProgressByDetailID(db, detailID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, results)
}

func (h *OKRHandler) UpdateProgress(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKRProgressRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateProgress(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteProgress(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteProgress(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// OKR Comments
// =========================================================================

func (h *OKRHandler) CreateComment(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRCommentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	userID := h.getUserID(c)

	result, err := h.service.CreateComment(db, &req, userID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) ListCommentsByEvaluationID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	evaluationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	results, err := h.service.ListCommentsByEvaluationID(db, evaluationID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, results)
}

func (h *OKRHandler) UpdateComment(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	var req UpdateOKRCommentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	result, err := h.service.UpdateComment(db, id, &req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}

func (h *OKRHandler) DeleteComment(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteComment(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// OKR Attachments
// =========================================================================

func (h *OKRHandler) CreateAttachment(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var req CreateOKRAttachmentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}

	userID := h.getUserID(c)

	result, err := h.service.CreateAttachment(db, &req, userID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.CreatedJSON(c, result, "success.created")
}

func (h *OKRHandler) ListAttachmentsByDetailID(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	detailID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	results, err := h.service.ListAttachmentsByDetailID(db, detailID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, results)
}

func (h *OKRHandler) DeleteAttachment(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", "error.invalid_id")
		return
	}

	if err := h.service.DeleteAttachment(db, id); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Dashboard
// =========================================================================

func (h *OKRHandler) GetHRDashboard(c *gin.Context) {
	db := h.getTenantDB(c)
	if db == nil {
		httputil.InternalError(c, "Database connection not found")
		return
	}

	var periodID *uuid.UUID
	if periodIDStr := c.Query("period_id"); periodIDStr != "" {
		if id, err := uuid.Parse(periodIDStr); err == nil {
			periodID = &id
		}
	}

	result, err := h.service.GetHRDashboard(db, periodID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}

	httputil.SuccessJSON(c, result)
}
