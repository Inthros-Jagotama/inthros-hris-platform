package attendance

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// businessTravelReportDateRange resolves the optional from/to query params,
// defaulting to first-of-current-month .. today (same convention as
// GetAttendanceStats), plus the optional status filter.
func businessTravelReportDateRange(c *gin.Context) (from, to, status string) {
	from = c.Query("from")
	to = c.Query("to")
	now := time.Now()
	if from == "" {
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if to == "" {
		to = now.Format("2006-01-02")
	}
	status = c.Query("status")
	return from, to, status
}

// GetBusinessTravelReport — daftar perjalanan dinas (rentang start_date).
// Di-gate oleh requireBusinessTravelReport() di routes.go.
func (h *Handler) GetBusinessTravelReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBusinessTravelFundingReport — laporan funding (rentang funding_date).
func (h *Handler) GetBusinessTravelFundingReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelFundingReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBusinessTravelAdvanceReport — advance vs actual expense per settlement.
func (h *Handler) GetBusinessTravelAdvanceReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelAdvanceReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBusinessTravelReimbursementReport — klaim reimbursement travel.
func (h *Handler) GetBusinessTravelReimbursementReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelReimbursementReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBusinessTravelRefundReport — refund kelebihan advance.
func (h *Handler) GetBusinessTravelRefundReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelRefundReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBusinessTravelCostReport — total biaya per perjalanan dinas.
func (h *Handler) GetBusinessTravelCostReport(c *gin.Context) {
	from, to, status := businessTravelReportDateRange(c)
	resp, err := h.service.GetBusinessTravelCostReport(c.Request.Context(), from, to, status)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
