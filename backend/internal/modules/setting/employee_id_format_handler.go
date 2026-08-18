package setting

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

// EmployeeIDFormatHandler exposes CRUD/preview endpoints for the employee_id
// generation-mode setting, backed by EmployeeIDFormatService.
type EmployeeIDFormatHandler struct {
	svc *EmployeeIDFormatService
}

func NewEmployeeIDFormatHandler(svc *EmployeeIDFormatService) *EmployeeIDFormatHandler {
	return &EmployeeIDFormatHandler{svc: svc}
}

// GET /api/v1/tenant/settings/employee-id-format — permission setting.employee_id_format.view
func (h *EmployeeIDFormatHandler) Get(c *gin.Context) {
	setting, err := h.svc.GetSetting(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, setting)
}

type updateEmployeeIDFormatRequest struct {
	GenerationMode string `json:"generation_mode" binding:"required"`
	FormatTemplate string `json:"format_template" binding:"required,max=255"`
	ResetPeriod    string `json:"reset_period" binding:"required"`
}

// PUT /api/v1/tenant/settings/employee-id-format — permission setting.employee_id_format.update
func (h *EmployeeIDFormatHandler) Update(c *gin.Context) {
	var req updateEmployeeIDFormatRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	updated, err := h.svc.UpdateSetting(c.Request.Context(), req.GenerationMode, req.FormatTemplate, req.ResetPeriod)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidEmployeeIDGenerationMode), errors.Is(err, ErrInvalidEmployeeIDResetPeriod):
			httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		default:
			httputil.InternalError(c, err.Error())
		}
		return
	}
	httputil.SuccessJSON(c, updated)
}

// GET /api/v1/tenant/settings/employee-id-format/preview — permission setting.employee_id_format.view
func (h *EmployeeIDFormatHandler) Preview(c *gin.Context) {
	preview, err := h.svc.Preview(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"preview": preview})
}
