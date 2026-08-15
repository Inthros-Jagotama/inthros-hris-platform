package payroll

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Payroll ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	payroll := rg.Group("/payroll")
	{
		// Formula Engine
		formula := payroll.Group("/formula")
		{
			formula.POST("/validate", handler.ValidateFormula)
			formula.GET("/variables", handler.ListFormulaVariables)
		}

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
			ebp.GET("", handler.ListEmployeeBankProfiles)
			ebp.GET("/:id", handler.GetEmployeeBankProfileByID)
			ebp.PUT("/:id", handler.UpdateEmployeeBankProfile)
			ebp.DELETE("/:id", handler.DeleteEmployeeBankProfile)
		}

		// Employee BPJS Profiles
		ebjs := payroll.Group("/employee-bpjs-profiles")
		{
			ebjs.POST("", handler.CreateEmployeeBpjsProfile)
			ebjs.GET("", handler.ListEmployeeBpjsProfiles)
			ebjs.GET("/:id", handler.GetEmployeeBpjsProfileByID)
			ebjs.PUT("/:id", handler.UpdateEmployeeBpjsProfile)
			ebjs.DELETE("/:id", handler.DeleteEmployeeBpjsProfile)
		}

		// Employee Tax Profiles
		etp := payroll.Group("/employee-tax-profiles")
		{
			etp.POST("", handler.CreateEmployeeTaxProfile)
			etp.GET("", handler.ListEmployeeTaxProfiles)
			etp.GET("/:id", handler.GetEmployeeTaxProfileByID)
			etp.PUT("/:id", handler.UpdateEmployeeTaxProfile)
			etp.DELETE("/:id", handler.DeleteEmployeeTaxProfile)
		}

		// Salary Structure — Grade Components
		sgc := payroll.Group("/salary-grade-components")
		{
			sgc.POST("", handler.CreateSalaryGradeComponent)
			sgc.GET("", handler.ListSalaryGradeComponents)
			sgc.GET("/:id", handler.GetSalaryGradeComponentByID)
			sgc.PUT("/:id", handler.UpdateSalaryGradeComponent)
			sgc.DELETE("/:id", handler.DeleteSalaryGradeComponent)
		}

		// Salary Structure — Employee Components (override)
		sec := payroll.Group("/salary-employee-components")
		{
			sec.POST("", handler.CreateSalaryEmployeeComponent)
			sec.GET("", handler.ListSalaryEmployeeComponents)
			sec.GET("/:id", handler.GetSalaryEmployeeComponentByID)
			sec.PUT("/:id", handler.UpdateSalaryEmployeeComponent)
			sec.DELETE("/:id", handler.DeleteSalaryEmployeeComponent)
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
			brc.GET("", handler.ListBpjsRateComponents)
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

		// PPh21 Tax Brackets
		ptb := payroll.Group("/pph21-tax-brackets")
		{
			ptb.POST("", handler.CreatePph21TaxBracket)
			ptb.GET("", handler.ListPph21TaxBrackets)
			ptb.PUT("/:id", handler.UpdatePph21TaxBracket)
			ptb.DELETE("/:id", handler.DeletePph21TaxBracket)
		}

		// Payroll Runs
		pr := payroll.Group("/runs")
		{
			pr.POST("", handler.CreatePayrollRun)
			pr.GET("", handler.ListPayrollRuns)
			pr.GET("/:id", handler.GetPayrollRunByID)
			pr.POST("/:id/calculate", handler.CalculatePayrollRun)
			pr.GET("/:id/employees", handler.ListPayrollRunEmployees)
			pr.GET("/:id/items", handler.ListPayrollRunItems)
			pr.PUT("/:id/status", handler.UpdatePayrollRunStatus)
			pr.GET("/:id/approval", handler.CheckPayrollRunApproval)
			pr.POST("/:id/payslips", handler.GeneratePayslips)
			pr.GET("/:id/payslips", handler.ListPayslipsByRun)
			pr.POST("/:id/payments", handler.CreatePaymentBatch)
			pr.GET("/:id/payments", handler.ListPaymentsByRun)
			pr.GET("/:id/payments/export", handler.ExportPaymentsCSV)
			pr.GET("/:id/reports/summary", handler.GetPayrollSummaryReport)
			pr.GET("/:id/reports/detail", handler.GetPayrollDetailReport)
			pr.GET("/:id/reports/bpjs", handler.GetBpjsReport)
			pr.GET("/:id/reports/tax", handler.GetTaxReport)
			pr.GET("/:id/reports/bank", handler.GetBankTransferReport)
			pr.GET("/:id/dashboard", handler.GetPayrollDashboard)
		}

		// Payslips
		payslips := payroll.Group("/payslips")
		{
			payslips.GET("/:id", handler.GetPayslipByID)
			payslips.GET("/:id/html", handler.GetPayslipHTML)
			payslips.POST("/:id/publish", handler.PublishPayslip)
			payslips.POST("/:id/cancel", handler.CancelPayslip)
		}

		// Payments
		payments := payroll.Group("/payments")
		{
			payments.GET("/:id", handler.GetPaymentByID)
			payments.POST("/:id/status", handler.UpdatePaymentStatus)
		}
	}
}
