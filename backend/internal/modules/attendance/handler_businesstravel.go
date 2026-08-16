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
