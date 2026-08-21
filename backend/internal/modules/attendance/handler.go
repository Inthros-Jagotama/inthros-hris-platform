package attendance

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/modules/approval"
	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Company Settings
// =========================================================================

func (h *Handler) UpsertCompanySetting(c *gin.Context) {
	var req CreateCompanySettingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpsertCompanySetting(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetCompanySetting(c *gin.Context) {
	resp, err := h.service.GetCompanySetting(c.Request.Context())
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetMyTimezone mengembalikan zona waktu efektif untuk user yang login,
// dipakai frontend untuk menampilkan jam/tanggal berjalan.
func (h *Handler) GetMyTimezone(c *gin.Context) {
	tz, err := h.service.GetMyTimezone(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"timezone": tz})
}

// =========================================================================
// Company Shifts
// =========================================================================

func (h *Handler) CreateShift(c *gin.Context) {
	var req CreateCompanyShiftRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateShift(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetShiftByID(c *gin.Context) {
	resp, err := h.service.GetShiftByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListShifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListShifts(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateShift(c *gin.Context) {
	var req UpdateCompanyShiftRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateShift(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteShift(c *gin.Context) {
	if err := h.service.DeleteShift(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Employee Shifts
// =========================================================================

func (h *Handler) CreateEmployeeShift(c *gin.Context) {
	var req CreateEmployeeShiftRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployeeShift(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetEmployeeShiftByID(c *gin.Context) {
	resp, err := h.service.GetEmployeeShiftByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListEmployeeShifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	employeeID := c.Query("employee_id")
	resp, err := h.service.ListEmployeeShifts(c.Request.Context(), &employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateEmployeeShift(c *gin.Context) {
	var req UpdateEmployeeShiftRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmployeeShift(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeShift(c *gin.Context) {
	if err := h.service.DeleteEmployeeShift(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Locations
// =========================================================================

func (h *Handler) CreateLocation(c *gin.Context) {
	var req CreateLocationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateLocation(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetLocationByID(c *gin.Context) {
	resp, err := h.service.GetLocationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListLocations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListLocations(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateLocation(c *gin.Context) {
	var req UpdateLocationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateLocation(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteLocation(c *gin.Context) {
	if err := h.service.DeleteLocation(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =========================================================================
// Events
// =========================================================================

func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEvent(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetEventByID(c *gin.Context) {
	resp, err := h.service.GetEventByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	employeeID := c.Query("employee_id")
	resp, err := h.service.ListEvents(c.Request.Context(), &employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Sessions
// =========================================================================

func (h *Handler) GetSession(c *gin.Context) {
	employeeID := c.Query("employee_id")
	workDate := c.Query("work_date")
	if employeeID == "" || workDate == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "attendance.query_params_required")
		return
	}
	resp, err := h.service.GetSession(c.Request.Context(), employeeID, workDate)
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListSessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	employeeID := c.Query("employee_id")
	resp, err := h.service.ListSessions(c.Request.Context(), &employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEmployeeCalendar(c *gin.Context) {
	employeeID := c.Query("employee_id")
	from := c.Query("from")
	to := c.Query("to")
	if employeeID == "" || from == "" || to == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "attendance.query_params_required")
		return
	}
	resp, err := h.service.GetEmployeeCalendar(c.Request.Context(), employeeID, from, to)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetAttendanceReport(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "attendance.query_params_required")
		return
	}
	resp, err := h.service.GetAttendanceReport(c.Request.Context(), from, to)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) GetEmployeeSummary(c *gin.Context) {
	employeeID := c.Query("employee_id")
	from := c.Query("from")
	to := c.Query("to")
	if employeeID == "" || from == "" || to == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "attendance.query_params_required")
		return
	}
	resp, err := h.service.GetEmployeeSummary(c.Request.Context(), employeeID, from, to)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetAttendanceStats — ringkasan absensi seluruh karyawan dalam rentang
// tanggal (mode HR dashboard). Default: bulan berjalan. Di-gate oleh
// middleware requireAttendanceReport("view") di routes.go.
func (h *Handler) GetAttendanceStats(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	now := time.Now()
	if from == "" {
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if to == "" {
		to = now.Format("2006-01-02")
	}
	resp, err := h.service.GetAttendanceStats(c.Request.Context(), from, to)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetOvertimeTrend — tren lembur per minggu (mode HR dashboard). Default:
// bulan berjalan. Sama seperti /stats/summary, digate attendance.report.view.
func (h *Handler) GetOvertimeTrend(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	now := time.Now()
	if from == "" {
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if to == "" {
		to = now.Format("2006-01-02")
	}
	resp, err := h.service.GetOvertimeTrend(c.Request.Context(), from, to)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =========================================================================
// Overtime Requests
// =========================================================================

func (h *Handler) CreateOvertimeRequest(c *gin.Context) {
	var req CreateOvertimeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateOvertimeRequest(c.Request.Context(), req)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their overtime didn't reach an
		// approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// GET /overtime-requests/assignable-employees — opsi dropdown alur ASSIGNED
// (§32b): karyawan di organisasi bawahan efektif dari user yang login.
func (h *Handler) ListAssignableEmployees(c *gin.Context) {
	resp, err := h.service.ListAssignableEmployees(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /overtime-requests/assign — alur ASSIGNED (§32b): atasan menugaskan
// lembur ke bawahan. Langsung WAITING_ACTUAL + notifikasi ke bawahan.
func (h *Handler) AssignOvertimeRequest(c *gin.Context) {
	var req AssignOvertimeRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.AssignOvertimeRequest(c.Request.Context(), req)
	if err != nil {
		writeOvertimeError(c, err)
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// POST /overtime-requests/:id/actual — submit isian aktual lembur (§32b).
func (h *Handler) SubmitOvertimeActual(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid overtime request id")
		return
	}
	var req SubmitOvertimeActualRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.SubmitActualOvertime(c.Request.Context(), id, req)
	if err != nil {
		writeOvertimeError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// POST /overtime-requests/:id/cancel — batal sebelum isian aktual (§32b).
func (h *Handler) CancelOvertimeRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid overtime request id")
		return
	}
	if err := h.service.CancelOvertimeRequest(c.Request.Context(), id); err != nil {
		writeOvertimeError(c, err)
		return
	}
	httputil.OKJSON(c, gin.H{"id": c.Param("id"), "status": "CANCELLED"})
}

// writeOvertimeError memetakan error layanan alur dua-tahap lembur (§32b) ke
// status HTTP yang tepat.
func writeOvertimeError(c *gin.Context, err error) {
	// Approval routing/assignee failures surface bilingually before the
	// regular error classification.
	if approval.EmitRoutingError(c, err) {
		return
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		httputil.NotFound(c, err.Error())
	case errors.Is(err, ErrOvertimeInvalidState):
		httputil.Conflict(c, "INVALID_STATUS", err.Error())
	case errors.Is(err, ErrOvertimeNotOwner):
		httputil.ErrorRaw(c, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, ErrOvertimeNotAssignable):
		httputil.ErrorRaw(c, http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, ErrOvertimeInvalidActualRange):
		httputil.BadRequest(c, err.Error())
	default:
		httputil.InternalError(c, err.Error())
	}
}

func (h *Handler) GetOvertimeRequestByID(c *gin.Context) {
	resp, err := h.service.GetOvertimeRequestByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListOvertimeRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	employeeID := c.Query("employee_id")
	resp, err := h.service.ListOvertimeRequests(c.Request.Context(), &employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Correction Requests
// =========================================================================

func (h *Handler) CreateCorrectionRequest(c *gin.Context) {
	var req CreateCorrectionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateCorrectionRequest(c.Request.Context(), req)
	if err != nil {
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why their correction didn't reach
		// an approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetCorrectionRequestByID(c *gin.Context) {
	resp, err := h.service.GetCorrectionRequestByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListCorrectionRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	employeeID := c.Query("employee_id")
	resp, err := h.service.ListCorrectionRequests(c.Request.Context(), &employeeID, page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Exempt Positions
// =========================================================================

func (h *Handler) CreateExemptPosition(c *gin.Context) {
	var req CreateExemptPositionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateExemptPosition(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetExemptPositionByID(c *gin.Context) {
	resp, err := h.service.GetExemptPositionByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListExemptPositions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListExemptPositions(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateExemptPosition(c *gin.Context) {
	var req UpdateExemptPositionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateExemptPosition(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteExemptPosition(c *gin.Context) {
	if err := h.service.DeleteExemptPosition(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
