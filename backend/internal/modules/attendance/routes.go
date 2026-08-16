package attendance

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Attendance ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/attendance
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	att := rg.Group("/attendance")
	{
		// Company Settings
		att.GET("/settings", handler.GetCompanySetting)
		att.PUT("/settings", handler.UpsertCompanySetting)

		// Company Shifts
		att.POST("/shifts", handler.CreateShift)
		att.GET("/shifts", handler.ListShifts)
		att.GET("/shifts/:id", handler.GetShiftByID)
		att.PUT("/shifts/:id", handler.UpdateShift)
		att.DELETE("/shifts/:id", handler.DeleteShift)

		// Employee Shifts
		att.POST("/employee-shifts", handler.CreateEmployeeShift)
		att.GET("/employee-shifts", handler.ListEmployeeShifts)
		att.GET("/employee-shifts/:id", handler.GetEmployeeShiftByID)
		att.PUT("/employee-shifts/:id", handler.UpdateEmployeeShift)
		att.DELETE("/employee-shifts/:id", handler.DeleteEmployeeShift)

		// Locations (Geofence)
		att.POST("/locations", handler.CreateLocation)
		att.GET("/locations", handler.ListLocations)
		att.GET("/locations/:id", handler.GetLocationByID)
		att.PUT("/locations/:id", handler.UpdateLocation)
		att.DELETE("/locations/:id", handler.DeleteLocation)

		// Events (Check-in / Check-out)
		att.POST("/events", handler.CreateEvent)
		att.GET("/events", handler.ListEvents)
		att.GET("/events/:id", handler.GetEventByID)

		// Sessions (Daily Work Sessions)
		att.GET("/sessions", handler.ListSessions)
		att.GET("/sessions/detail", handler.GetSession)
		att.GET("/calendar", handler.GetEmployeeCalendar)
		att.GET("/summary", handler.GetEmployeeSummary)
		att.GET("/reports/sessions", handler.GetAttendanceReport)

		// Overtime Requests (§32b dua-alur: SELF & ASSIGNED → isian aktual)
		att.POST("/overtime-requests", handler.CreateOvertimeRequest)
		att.GET("/overtime-requests/assignable-employees", handler.ListAssignableEmployees)
		att.POST("/overtime-requests/assign", handler.AssignOvertimeRequest)
		att.POST("/overtime-requests/:id/actual", handler.SubmitOvertimeActual)
		att.POST("/overtime-requests/:id/cancel", handler.CancelOvertimeRequest)
		att.GET("/overtime-requests", handler.ListOvertimeRequests)
		att.GET("/overtime-requests/:id", handler.GetOvertimeRequestByID)

		// Correction Requests
		att.POST("/corrections", handler.CreateCorrectionRequest)
		att.GET("/corrections", handler.ListCorrectionRequests)
		att.GET("/corrections/:id", handler.GetCorrectionRequestByID)

		// Exempt Positions
		att.POST("/exempt-positions", handler.CreateExemptPosition)
		att.GET("/exempt-positions", handler.ListExemptPositions)
		att.GET("/exempt-positions/:id", handler.GetExemptPositionByID)
		att.PUT("/exempt-positions/:id", handler.UpdateExemptPosition)
		att.DELETE("/exempt-positions/:id", handler.DeleteExemptPosition)

		// Business Travel (Phase 1-3: Travel CRUD + Approval integration.
		// Funding/Expense/Settlement/Refund/Reimbursement belum diimplementasikan —
		// lihat docs/module-attendance-business-travel-development-plan.md §54.7 urutan kerja).
		att.POST("/business-travels", handler.CreateBusinessTravel)
		att.GET("/business-travels", handler.ListBusinessTravels)
		att.GET("/business-travels/:id", handler.GetBusinessTravelByID)
		att.PUT("/business-travels/:id", handler.UpdateBusinessTravel)
		att.POST("/business-travels/:id/submit", handler.SubmitBusinessTravel)
		att.POST("/business-travels/:id/cancel", handler.CancelBusinessTravel)
		att.POST("/business-travels/:id/participants", handler.AddBusinessTravelParticipant)
		att.GET("/business-travels/:id/participants", handler.ListBusinessTravelParticipants)
		att.POST("/business-travels/:id/destinations", handler.AddBusinessTravelDestination)
		att.GET("/business-travels/:id/destinations", handler.ListBusinessTravelDestinations)
		att.POST("/business-travels/:id/activities", handler.AddBusinessTravelActivity)
		att.GET("/business-travels/:id/activities", handler.ListBusinessTravelActivities)
		att.POST("/business-travels/:id/schedules", handler.AddBusinessTravelSchedule)
		att.GET("/business-travels/:id/schedules", handler.ListBusinessTravelSchedules)

		// Funding (§14-17: hanya bisa dibuat setelah travel APPROVED)
		att.POST("/business-travel-funding-methods", handler.CreateFundingMethod)
		att.GET("/business-travel-funding-methods", handler.ListFundingMethods)
		att.POST("/business-travels/:id/fundings", handler.CreateFunding)
		att.GET("/business-travels/:id/fundings", handler.ListFundings)
		att.PUT("/business-travels/:id/fundings/:fundingId", handler.UpdateFunding)
		att.POST("/business-travels/:id/fundings/:fundingId/confirm", handler.ConfirmFunding)
		att.POST("/business-travels/:id/fundings/:fundingId/documents", handler.AddFundingDocument)

		// Actual Expense (§21-22: dicatat pasca perjalanan, funding method per
		// item boleh berbeda dari funding awal travel — mixed funding §33)
		att.POST("/business-travel-expense-categories", handler.CreateExpenseCategory)
		att.GET("/business-travel-expense-categories", handler.ListExpenseCategories)
		att.POST("/business-travels/:id/expenses", handler.CreateExpense)
		att.GET("/business-travels/:id/expenses", handler.ListExpenses)
		att.PUT("/business-travels/:id/expenses/:expenseId", handler.UpdateExpense)
		att.DELETE("/business-travels/:id/expenses/:expenseId", handler.DeleteExpense)
		att.POST("/business-travels/:id/expenses/:expenseId/documents", handler.AddExpenseDocument)

		// Travel Documents (§23: travel order, ticket, boarding pass, hotel,
		// travel report, dst — level travel, terpisah dari funding/expense document)
		att.POST("/business-travels/:id/documents", handler.AddTravelDocument)
		att.GET("/business-travels/:id/documents", handler.ListTravelDocuments)
		att.DELETE("/business-travels/:id/documents/:documentId", handler.DeleteTravelDocument)

		// Settlement (§24-33: hanya bisa dibuat setelah travel COMPLETED;
		// approval terpisah dari Travel Approval — module slug
		// "business_travel_settlement")
		att.POST("/business-travels/:id/settlements", handler.CreateSettlement)
		att.GET("/business-travels/:id/settlements", handler.ListSettlements)
		att.GET("/business-travels/:id/settlements/:settlementId", handler.GetSettlementByID)
		att.POST("/business-travels/:id/settlements/:settlementId/submit", handler.SubmitSettlement)

		// Refund (§35): dibuat otomatis oleh HandleSettlementApprovalStatusChange
		att.GET("/business-travels/:id/refunds", handler.ListRefunds)
		att.POST("/business-travels/:id/refunds/:refundId/confirm", handler.ConfirmRefund)

		// Reimbursement (§36, §54.7: cek subscription module Reimbursement
		// sebelum diproses — lihat Service.ProcessTravelReimbursement)
		att.GET("/business-travels/:id/reimbursements", handler.ListTravelReimbursements)
		att.POST("/business-travels/:id/reimbursements/:reimbursementId/approve", handler.ApproveTravelReimbursement)
		att.POST("/business-travels/:id/reimbursements/:reimbursementId/process", handler.ProcessTravelReimbursement)
		att.POST("/business-travels/:id/reimbursements/:reimbursementId/pay", handler.PayTravelReimbursement)
	}
}
