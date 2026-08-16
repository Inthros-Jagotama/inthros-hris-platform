package attendance

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// =========================================================================
// Business Travel
// =========================================================================

func (h *Handler) CreateBusinessTravel(c *gin.Context) {
	var req CreateBusinessTravelRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateBusinessTravel(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetBusinessTravelByID(c *gin.Context) {
	resp, err := h.service.GetBusinessTravelByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListBusinessTravels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	requesterID := c.Query("requester_id")
	status := c.Query("status")
	resp, err := h.service.ListBusinessTravels(c.Request.Context(), requesterID, status, page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateBusinessTravel(c *gin.Context) {
	var req UpdateBusinessTravelRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateBusinessTravel(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ErrBusinessTravelInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.NotFound(c, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) SubmitBusinessTravel(c *gin.Context) {
	var req SubmitBusinessTravelRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.SubmitBusinessTravel(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ErrBusinessTravelInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AddBusinessTravelActivity(c *gin.Context) {
	var req CreateBusinessTravelActivityRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AddBusinessTravelActivity(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListBusinessTravelActivities(c *gin.Context) {
	resp, err := h.service.ListBusinessTravelActivities(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AddBusinessTravelSchedule(c *gin.Context) {
	var req CreateBusinessTravelScheduleRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AddBusinessTravelSchedule(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListBusinessTravelSchedules(c *gin.Context) {
	resp, err := h.service.ListBusinessTravelSchedules(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Funding
// =========================================================================

func (h *Handler) CreateFundingMethod(c *gin.Context) {
	var req CreateFundingMethodRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateFundingMethod(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListFundingMethods(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "true") == "true"
	resp, err := h.service.ListFundingMethods(c.Request.Context(), activeOnly)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateFunding(c *gin.Context) {
	var req CreateFundingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateFunding(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ErrBusinessTravelNotApproved) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListFundings(c *gin.Context) {
	resp, err := h.service.ListFundingsByTravel(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateFunding(c *gin.Context) {
	var req UpdateFundingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateFunding(c.Request.Context(), c.Param("fundingId"), req)
	if err != nil {
		if errors.Is(err, ErrFundingInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ConfirmFunding(c *gin.Context) {
	resp, err := h.service.ConfirmFunding(c.Request.Context(), c.Param("fundingId"))
	if err != nil {
		if errors.Is(err, ErrFundingInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) AddFundingDocument(c *gin.Context) {
	var req AddFundingDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AddFundingDocument(c.Request.Context(), c.Param("fundingId"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =========================================================================
// Expense Category & Actual Expense
// =========================================================================

func (h *Handler) CreateExpenseCategory(c *gin.Context) {
	var req CreateExpenseCategoryRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateExpenseCategory(c.Request.Context(), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListExpenseCategories(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "true") == "true"
	resp, err := h.service.ListExpenseCategories(c.Request.Context(), activeOnly)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateExpense(c *gin.Context) {
	var req CreateExpenseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateExpense(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ErrBusinessTravelNotApproved) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListExpenses(c *gin.Context) {
	resp, err := h.service.ListExpensesByTravel(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	var req UpdateExpenseRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateExpense(c.Request.Context(), c.Param("expenseId"), req)
	if err != nil {
		if errors.Is(err, ErrExpenseInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	if err := h.service.DeleteExpense(c.Request.Context(), c.Param("expenseId")); err != nil {
		if errors.Is(err, ErrExpenseInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"deleted": true})
}

func (h *Handler) AddExpenseDocument(c *gin.Context) {
	var req AddExpenseDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AddExpenseDocument(c.Request.Context(), c.Param("expenseId"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) CancelBusinessTravel(c *gin.Context) {
	resp, err := h.service.CancelBusinessTravel(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrBusinessTravelInvalidState) {
			httputil.ErrorSimple(c, http.StatusConflict, err.Error())
			return
		}
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
