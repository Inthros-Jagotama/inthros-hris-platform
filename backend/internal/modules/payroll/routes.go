package payroll

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Payroll ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	payroll := rg.Group("/payroll")
	{
		// Salary Components
		sc := payroll.Group("/salary-components")
		{
			sc.POST("", handler.CreateSalaryComponent)
			sc.GET("", handler.ListSalaryComponents)
			sc.GET("/:id", handler.GetSalaryComponentByID)
			sc.PUT("/:id", handler.UpdateSalaryComponent)
			sc.DELETE("/:id", handler.DeleteSalaryComponent)
		}

		// Payroll Periods
		pp := payroll.Group("/periods")
		{
			pp.POST("", handler.CreatePayrollPeriod)
			pp.GET("", handler.ListPayrollPeriods)
			pp.PUT("/:id", handler.UpdatePayrollPeriod)
		}

		// Employee Payroll Profiles
		epp := payroll.Group("/employee-payroll-profiles")
		{
			epp.POST("", handler.CreateEmployeePayrollProfile)
			epp.GET("", handler.ListEmployeePayrollProfiles)
			epp.GET("/:id", handler.GetEmployeePayrollProfileByID)
			epp.DELETE("/:id", handler.DeleteEmployeePayrollProfile)
		}

		// Employee Bank Profiles
		ebp := payroll.Group("/employee-bank-profiles")
		{
			ebp.POST("", handler.CreateEmployeeBankProfile)
			ebp.GET("/:id", handler.GetEmployeeBankProfileByID)
			ebp.PUT("/:id", handler.UpdateEmployeeBankProfile)
			ebp.DELETE("/:id", handler.DeleteEmployeeBankProfile)
		}

		// Employee BPJS Profiles
		ebjs := payroll.Group("/employee-bpjs-profiles")
		{
			ebjs.POST("", handler.CreateEmployeeBpjsProfile)
			ebjs.GET("/:id", handler.GetEmployeeBpjsProfileByID)
			ebjs.PUT("/:id", handler.UpdateEmployeeBpjsProfile)
			ebjs.DELETE("/:id", handler.DeleteEmployeeBpjsProfile)
		}

		// Employee Tax Profiles
		etp := payroll.Group("/employee-tax-profiles")
		{
			etp.POST("", handler.CreateEmployeeTaxProfile)
			etp.GET("/:id", handler.GetEmployeeTaxProfileByID)
			etp.PUT("/:id", handler.UpdateEmployeeTaxProfile)
			etp.DELETE("/:id", handler.DeleteEmployeeTaxProfile)
		}

		// BPJS Settings
		bs := payroll.Group("/bpjs-settings")
		{
			bs.POST("", handler.CreateBpjsSetting)
			bs.GET("", handler.ListBpjsSettings)
			bs.GET("/:id", handler.GetBpjsSettingByID)
			bs.PUT("/:id", handler.UpdateBpjsSetting)
			bs.DELETE("/:id", handler.DeleteBpjsSetting)
		}

		// BPJS Rate Components
		brc := payroll.Group("/bpjs-rate-components")
		{
			brc.POST("", handler.CreateBpjsRateComponent)
			brc.GET("/:id", handler.GetBpjsRateComponentByID)
			brc.PUT("/:id", handler.UpdateBpjsRateComponent)
			brc.DELETE("/:id", handler.DeleteBpjsRateComponent)
		}

		// PPh21 Settings
		ps := payroll.Group("/pph21-settings")
		{
			ps.POST("", handler.CreatePph21Setting)
			ps.GET("", handler.ListPph21Settings)
			ps.GET("/:id", handler.GetPph21SettingByID)
			ps.PUT("/:id", handler.UpdatePph21Setting)
			ps.DELETE("/:id", handler.DeletePph21Setting)
		}

		// PPh21 PTKP Rates
		ppr := payroll.Group("/pph21-ptkp-rates")
		{
			ppr.POST("", handler.CreatePph21PtkpRate)
			ppr.GET("", handler.ListPph21PtkpRates)
		}

		// PPh21 Tax Brackets
		ptb := payroll.Group("/pph21-tax-brackets")
		{
			ptb.POST("", handler.CreatePph21TaxBracket)
			ptb.GET("", handler.ListPph21TaxBrackets)
		}

		// Payroll Runs
		pr := payroll.Group("/runs")
		{
			pr.POST("", handler.CreatePayrollRun)
			pr.GET("", handler.ListPayrollRuns)
			pr.GET("/:id", handler.GetPayrollRunByID)
			pr.PUT("/:id/status", handler.UpdatePayrollRunStatus)
			pr.GET("/:id/approval", handler.CheckPayrollRunApproval)
		}
	}
}
