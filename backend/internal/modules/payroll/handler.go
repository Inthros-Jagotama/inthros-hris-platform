package payroll

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/modules/approval"
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

// handleSalaryComponentError memetakan ValidationError (dari formula/reference
// validation) ke HTTP 400, selain error umum ke 500.
func (h *Handler) handleSalaryComponentError(c *gin.Context, err error) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		httputil.BadRequest(c, ve.Error())
		return
	}
	httputil.InternalError(c, err.Error())
}

func (h *Handler) CreateSalaryComponent(c *gin.Context) {
	var req CreateSalaryComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateSalaryComponent(c.Request.Context(), req)
	if err != nil {
		h.handleSalaryComponentError(c, err)
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

// handleServiceError memetakan error service generik: ValidationError → 400,
// lainnya → 500 (dipakai endpoint payslip/payment).
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		httputil.BadRequest(c, ve.Error())
		return
	}
	httputil.InternalError(c, err.Error())
}

func (h *Handler) ListSalaryComponents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListSalaryComponents(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
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
		h.handleSalaryComponentError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteSalaryComponent(c *gin.Context) {
	if err := h.service.DeleteSalaryComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =============================================================================
// Salary Structure — Grade Components
// =============================================================================

func (h *Handler) ListSalaryGradeComponents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListSalaryGradeComponents(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSalaryGradeComponentByID(c *gin.Context) {
	resp, err := h.service.GetSalaryGradeComponentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateSalaryGradeComponent(c *gin.Context) {
	var req CreateSalaryGradeComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateSalaryGradeComponent(c.Request.Context(), req)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateSalaryGradeComponent(c *gin.Context) {
	var req UpdateSalaryGradeComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateSalaryGradeComponent(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteSalaryGradeComponent(c *gin.Context) {
	if err := h.service.DeleteSalaryGradeComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

// =============================================================================
// Salary Structure — Employee Components (override)
// =============================================================================

func (h *Handler) ListSalaryEmployeeComponents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListSalaryEmployeeComponents(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSalaryEmployeeComponentByID(c *gin.Context) {
	resp, err := h.service.GetSalaryEmployeeComponentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) CreateSalaryEmployeeComponent(c *gin.Context) {
	var req CreateSalaryEmployeeComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.CreateSalaryEmployeeComponent(c.Request.Context(), req)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) UpdateSalaryEmployeeComponent(c *gin.Context) {
	var req UpdateSalaryEmployeeComponentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdateSalaryEmployeeComponent(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteSalaryEmployeeComponent(c *gin.Context) {
	if err := h.service.DeleteSalaryEmployeeComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteEmployeePayrollProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeePayrollProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeBankProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeBankProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeBpjsProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeBpjsProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteEmployeeTaxProfile(c *gin.Context) {
	if err := h.service.DeleteEmployeeTaxProfile(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) ListEmployeeBankProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEmployeeBankProfiles(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListEmployeeBpjsProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEmployeeBpjsProfiles(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListEmployeeTaxProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	resp, err := h.service.ListEmployeeTaxProfiles(c.Request.Context(), page, perPage)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteBpjsSetting(c *gin.Context) {
	if err := h.service.DeleteBpjsSetting(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteBpjsRateComponent(c *gin.Context) {
	if err := h.service.DeleteBpjsRateComponent(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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

// ListBpjsRateComponents mengembalikan rate component sebuah setting BPJS.
// Endpoint: GET /payroll/bpjs-rate-components?bpjs_setting_id=...
func (h *Handler) ListBpjsRateComponents(c *gin.Context) {
	settingID := c.Query("bpjs_setting_id")
	if settingID == "" {
		httputil.BadRequest(c, "bpjs_setting_id is required")
		return
	}
	resp, err := h.service.ListBpjsRateComponentsBySettingID(c.Request.Context(), settingID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
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
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePph21Setting(c *gin.Context) {
	if err := h.service.DeletePph21Setting(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
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
		httputil.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdatePph21TaxBracket(c *gin.Context) {
	var req UpdatePph21TaxBracketRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePph21TaxBracket(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeletePph21TaxBracket(c *gin.Context) {
	if err := h.service.DeletePph21TaxBracket(c.Request.Context(), c.Param("id")); err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{"deleted": true})
}

// =============================================================================
// Formula Engine
// =============================================================================

// ValidateFormulaRequest body untuk POST /payroll/formula/validate.
type ValidateFormulaRequest struct {
	Formula string `json:"formula" binding:"required"`
}

// ValidateFormula memvalidasi sintaks formula dan mengembalikan variabel yang
// direferensikan. Dipakai frontend untuk validasi saat user mengetik formula.
func (h *Handler) ValidateFormula(c *gin.Context) {
	var req ValidateFormulaRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	vars, err := h.service.ValidateFormula(c.Request.Context(), req.Formula)
	if err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, gin.H{
		"valid":     true,
		"formula":   req.Formula,
		"variables": vars,
	})
}

// ListFormulaVariables mengembalikan daftar variabel built-in formula engine.
func (h *Handler) ListFormulaVariables(c *gin.Context) {
	vars := h.service.ListFormulaVariables(c.Request.Context())
	httputil.SuccessJSON(c, vars)
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
		httputil.InternalError(c, err.Error())
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

// CalculatePayrollRun mengeksekusi kalkulasi + snapshot untuk sebuah run.
// Endpoint: POST /payroll/runs/:id/calculate
func (h *Handler) CalculatePayrollRun(c *gin.Context) {
	resp, err := h.service.CalculatePayrollRun(c.Request.Context(), c.Param("id"))
	if err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			httputil.BadRequest(c, ve.Error())
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ListPayrollRunEmployees mengembalikan snapshot employee sebuah run.
// Endpoint: GET /payroll/runs/:id/employees
func (h *Handler) ListPayrollRunEmployees(c *gin.Context) {
	resp, err := h.service.ListPayrollRunEmployees(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ListPayrollRunItems mengembalikan snapshot item (detail komponen) sebuah run.
// Endpoint: GET /payroll/runs/:id/items
func (h *Handler) ListPayrollRunItems(c *gin.Context) {
	resp, err := h.service.ListPayrollRunItems(c.Request.Context(), c.Param("id"))
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =============================================================================
// Payslips
// =============================================================================

// GeneratePayslips membuat payslip dari run yang sudah dihitung.
// Endpoint: POST /payroll/runs/:id/payslips
func (h *Handler) GeneratePayslips(c *gin.Context) {
	resp, err := h.service.GeneratePayslips(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ListPayslipsByRun mengembalikan payslip sebuah run.
// Endpoint: GET /payroll/runs/:id/payslips
func (h *Handler) ListPayslipsByRun(c *gin.Context) {
	resp, err := h.service.ListPayslipsByRun(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetPayslipByID mengambil detail payslip.
// Endpoint: GET /payroll/payslips/:id
func (h *Handler) GetPayslipByID(c *gin.Context) {
	resp, err := h.service.GetPayslipByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// PublishPayslip mempublikasikan payslip.
// Endpoint: POST /payroll/payslips/:id/publish
func (h *Handler) PublishPayslip(c *gin.Context) {
	resp, err := h.service.PublishPayslip(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// CancelPayslip membatalkan payslip.
// Endpoint: POST /payroll/payslips/:id/cancel
func (h *Handler) CancelPayslip(c *gin.Context) {
	resp, err := h.service.CancelPayslip(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetPayslipHTML merender payslip sebagai HTML.
// Endpoint: GET /payroll/payslips/:id/html
func (h *Handler) GetPayslipHTML(c *gin.Context) {
	html, err := h.service.GetPayslipHTML(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// =============================================================================
// Payments
// =============================================================================

// CreatePaymentBatch membuat batch pembayaran dari run.
// Endpoint: POST /payroll/runs/:id/payments
func (h *Handler) CreatePaymentBatch(c *gin.Context) {
	resp, err := h.service.CreatePaymentBatch(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ListPaymentsByRun mengembalikan payment sebuah run.
// Endpoint: GET /payroll/runs/:id/payments
func (h *Handler) ListPaymentsByRun(c *gin.Context) {
	resp, err := h.service.ListPaymentsByRun(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// ExportPaymentsCSV mengekspor batch sebagai file bank transfer CSV.
// Endpoint: GET /payroll/runs/:id/payments/export
func (h *Handler) ExportPaymentsCSV(c *gin.Context) {
	csvOut, err := h.service.ExportPaymentsCSV(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=payments.csv")
	c.String(http.StatusOK, csvOut)
}

// GetPaymentByID mengambil detail payment.
// Endpoint: GET /payroll/payments/:id
func (h *Handler) GetPaymentByID(c *gin.Context) {
	resp, err := h.service.GetPaymentByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// UpdatePaymentStatus memindahkan status payment.
// Endpoint: POST /payroll/payments/:id/status
func (h *Handler) UpdatePaymentStatus(c *gin.Context) {
	var req UpdatePaymentStatusRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.service.UpdatePaymentStatus(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// =============================================================================
// Reports & Dashboard
// =============================================================================

// GetPayrollSummaryReport — ringkasan run.
// Endpoint: GET /payroll/runs/:id/reports/summary
func (h *Handler) GetPayrollSummaryReport(c *gin.Context) {
	resp, err := h.service.GetPayrollSummaryReport(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetPayrollDetailReport — rincian per employee per komponen.
// Endpoint: GET /payroll/runs/:id/reports/detail
func (h *Handler) GetPayrollDetailReport(c *gin.Context) {
	resp, err := h.service.GetPayrollDetailReport(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBpjsReport — laporan BPJS per employee.
// Endpoint: GET /payroll/runs/:id/reports/bpjs
func (h *Handler) GetBpjsReport(c *gin.Context) {
	resp, err := h.service.GetBpjsReport(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetTaxReport — laporan pajak per employee.
// Endpoint: GET /payroll/runs/:id/reports/tax
func (h *Handler) GetTaxReport(c *gin.Context) {
	resp, err := h.service.GetTaxReport(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetBankTransferReport — laporan bank transfer dari payment batch.
// Endpoint: GET /payroll/runs/:id/reports/bank
func (h *Handler) GetBankTransferReport(c *gin.Context) {
	resp, err := h.service.GetBankTransferReport(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, resp)
}

// GetPayrollDashboard — agregat dashboard run.
// Endpoint: GET /payroll/runs/:id/dashboard
func (h *Handler) GetPayrollDashboard(c *gin.Context) {
	resp, err := h.service.GetPayrollDashboard(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
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
		// Approval routing/assignee failures (approval.RoutingError) get a
		// bilingual 400 so the user sees why the run didn't reach an
		// approver instead of a raw internal error.
		if approval.EmitRoutingError(c, err) {
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

// CheckPayrollRunApproval handles GET /api/v1/tenant/payroll/runs/:id/approval
func (h *Handler) CheckPayrollRunApproval(c *gin.Context) {
	id := c.Param("id")
	instanceID := c.Query("instance_id")
	if instanceID == "" {
		httputil.ErrorJSON(c, http.StatusBadRequest, "VALIDATION_ERROR", "payroll.instance_param_required")
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
