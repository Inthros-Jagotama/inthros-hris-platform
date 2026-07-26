package payroll

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Salary Components
// =============================================================================

func (h *Handler) CreateSalaryComponent(c *gin.Context) {
	var req CreateSalaryComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateSalaryComponent(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListSalaryComponents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListSalaryComponents(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSalaryComponentByID(c *gin.Context) {
	resp, err := h.service.GetSalaryComponentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateSalaryComponent(c *gin.Context) {
	var req UpdateSalaryComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateSalaryComponent(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteSalaryComponent(c *gin.Context) {
	if err := h.service.DeleteSalaryComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =============================================================================
// Payroll Periods
// =============================================================================

func (h *Handler) CreatePayrollPeriod(c *gin.Context) {
	var req CreatePayrollPeriodRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePayrollPeriod(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPayrollPeriods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListPayrollPeriods(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePayrollPeriod(c *gin.Context) {
	var req UpdatePayrollPeriodRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePayrollPeriod(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =============================================================================
// Employee Payroll Profiles
// =============================================================================

func (h *Handler) CreateEmployeePayrollProfile(c *gin.Context) {
	var req CreateEmployeePayrollProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployeePayrollProfile(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetEmployeePayrollProfileByID(c *gin.Context) {
	resp, err := h.service.GetEmployeePayrollProfileByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListEmployeePayrollProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEmployeePayrollProfiles(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEmployeePayrollProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeePayrollProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =============================================================================
// Employee Bank Profiles
// =============================================================================

func (h *Handler) GetEmployeeBankProfileByID(c *gin.Context) {
	resp, err := h.service.GetEmployeeBankProfileByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEmployeeBankProfile(c *gin.Context) {
	var req UpdateEmployeeBankProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmployeeBankProfile(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeBankProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeBankProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateEmployeeBankProfile(c *gin.Context) {
	var req CreateEmployeeBankProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployeeBankProfile(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =============================================================================
// Employee BPJS Profiles
// =============================================================================

func (h *Handler) GetEmployeeBpjsProfileByID(c *gin.Context) {
	resp, err := h.service.GetEmployeeBpjsProfileByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEmployeeBpjsProfile(c *gin.Context) {
	var req UpdateEmployeeBpjsProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmployeeBpjsProfile(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeBpjsProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeBpjsProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateEmployeeBpjsProfile(c *gin.Context) {
	var req CreateEmployeeBpjsProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployeeBpjsProfile(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =============================================================================
// Employee Tax Profiles
// =============================================================================

func (h *Handler) GetEmployeeTaxProfileByID(c *gin.Context) {
	resp, err := h.service.GetEmployeeTaxProfileByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateEmployeeTaxProfile(c *gin.Context) {
	var req UpdateEmployeeTaxProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateEmployeeTaxProfile(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeTaxProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeTaxProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateEmployeeTaxProfile(c *gin.Context) {
	var req CreateEmployeeTaxProfileRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateEmployeeTaxProfile(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =============================================================================
// BPJS Settings
// =============================================================================

func (h *Handler) CreateBpjsSetting(c *gin.Context) {
	var req CreateBpjsSettingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateBpjsSetting(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) GetBpjsSettingByID(c *gin.Context) {
	resp, err := h.service.GetBpjsSettingByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListBpjsSettings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListBpjsSettings(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateBpjsSetting(c *gin.Context) {
	var req UpdateBpjsSettingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateBpjsSetting(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteBpjsSetting(c *gin.Context) {
	if err := h.service.DeleteBpjsSetting(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =============================================================================
// BPJS Rate Components
// =============================================================================

func (h *Handler) GetBpjsRateComponentByID(c *gin.Context) {
	resp, err := h.service.GetBpjsRateComponentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateBpjsRateComponent(c *gin.Context) {
	var req UpdateBpjsRateComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateBpjsRateComponent(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteBpjsRateComponent(c *gin.Context) {
	if err := h.service.DeleteBpjsRateComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateBpjsRateComponent(c *gin.Context) {
	var req CreateBpjsRateComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateBpjsRateComponent(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =============================================================================
// PPh21 Settings
// =============================================================================

func (h *Handler) GetPph21SettingByID(c *gin.Context) {
	resp, err := h.service.GetPph21SettingByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) ListPph21Settings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListPph21Settings(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePph21Setting(c *gin.Context) {
	var req UpdatePph21SettingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePph21Setting(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePph21Setting(c *gin.Context) {
	if err := h.service.DeletePph21Setting(c.Request.Context(), c.Param("id")); err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreatePph21Setting(c *gin.Context) {
	var req CreatePph21SettingRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePph21Setting(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// =============================================================================
// PPh21 PTKP Rates
// =============================================================================

func (h *Handler) CreatePph21PtkpRate(c *gin.Context) {
	var req CreatePph21PtkpRateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePph21PtkpRate(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPph21PtkpRates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListPph21PtkpRates(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =============================================================================
// PPh21 Tax Brackets
// =============================================================================

func (h *Handler) CreatePph21TaxBracket(c *gin.Context) {
	var req CreatePph21TaxBracketRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePph21TaxBracket(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPph21TaxBrackets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListPph21TaxBrackets(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =============================================================================
// Payroll Runs
// =============================================================================

func (h *Handler) CreatePayrollRun(c *gin.Context) {
	var req CreatePayrollRunRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreatePayrollRun(c.Request.Context(), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListPayrollRuns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListPayrollRuns(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPayrollRunByID(c *gin.Context) {
	resp, err := h.service.GetPayrollRunByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdatePayrollRunStatus(c *gin.Context) {
	var req UpdatePayrollRunStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePayrollRunStatus(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.ErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// CheckPayrollRunApproval handles GET /api/v1/tenant/payroll/runs/:id/approval
func (h *Handler) CheckPayrollRunApproval(c *gin.Context) {
	id := c.Param("id")
	instanceID := c.Query("instance_id")
	if instanceID == "" {
		httputil.ErrorRaw(c, http.StatusBadRequest, "VALIDATION_ERROR", "instance_id query parameter is required")
		return
	}

	response, err := h.service.CheckPayrollRunApproval(c.Request.Context(), id, instanceID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "APPROVAL_CHECK_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	httputil.SuccessJSON(c, response)
}
