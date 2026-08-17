package attendance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// requireAttendanceSettings menggating endpoint yang MENGUBAH (create/update/
// delete) konfigurasi absensi (company settings, shifts, employee-shifts,
// locations, exempt-positions) beserta GET-nya — supaya HANYA permission
// submenu "attendance.settings.<action>" yang berlaku (bukan module-level
// "attendance.view"). Middleware RBAC global (authz.NewMiddleware) menganggap
// module-level otomatis mencakup semua submenu-nya — cocok untuk kebanyakan
// resource, tapi tidak untuk aksi admin-config ini, yang harus benar-benar
// terpisah dari sekadar melihat/check-in absensi sehari-hari ("attendance.view",
// yang dimiliki hampir semua role termasuk Employee default). Pola sama seperti
// requireLeaveSettings di modul leave.
//
// GET /settings SENGAJA TIDAK dibatasi middleware ini — dipakai halaman
// Attendance utama (aturan check-in saat hari libur) yang harus tetap bisa
// diakses siapa pun dengan "attendance.view".
func requireAttendanceSettings(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := "attendance.settings." + action
		for _, p := range c.GetStringSlice("permissions") {
			if p == "*" || p == required {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You don't have permission to manage attendance settings",
				"details": gin.H{"required": required},
			},
		})
	}
}

// requireAttendanceReport menggating endpoint laporan absensi (GET
// /reports/sessions) dengan permission "attendance.report.view" secara EKSAK —
// module-level "attendance.view" tidak boleh mencakupnya, pola sama seperti
// requireAttendanceSettings. Endpoint operations (events/sessions) SENGAJA
// tidak digate di sini: GET /events & /sessions/detail juga dipakai halaman
// Attendance utama & halaman Koreksi (alur self-service check-in/koreksi).
func requireAttendanceReport(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := "attendance.report." + action
		for _, p := range c.GetStringSlice("permissions") {
			if p == "*" || p == required {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You don't have permission to view attendance reports",
				"details": gin.H{"required": required},
			},
		})
	}
}

// RegisterRoutes mendaftarkan semua endpoint Attendance ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/attendance
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	att := rg.Group("/attendance")
	{
		// Company Settings (GET tetap terbuka — dipakai halaman Attendance utama)
		att.GET("/settings", handler.GetCompanySetting)
		att.PUT("/settings", requireAttendanceSettings("update"), handler.UpsertCompanySetting)

		// Company Shifts (admin-config)
		att.POST("/shifts", requireAttendanceSettings("create"), handler.CreateShift)
		att.GET("/shifts", requireAttendanceSettings("view"), handler.ListShifts)
		att.GET("/shifts/:id", requireAttendanceSettings("view"), handler.GetShiftByID)
		att.PUT("/shifts/:id", requireAttendanceSettings("update"), handler.UpdateShift)
		att.DELETE("/shifts/:id", requireAttendanceSettings("delete"), handler.DeleteShift)

		// Employee Shifts (admin-config)
		att.POST("/employee-shifts", requireAttendanceSettings("create"), handler.CreateEmployeeShift)
		att.GET("/employee-shifts", requireAttendanceSettings("view"), handler.ListEmployeeShifts)
		att.GET("/employee-shifts/:id", requireAttendanceSettings("view"), handler.GetEmployeeShiftByID)
		att.PUT("/employee-shifts/:id", requireAttendanceSettings("update"), handler.UpdateEmployeeShift)
		att.DELETE("/employee-shifts/:id", requireAttendanceSettings("delete"), handler.DeleteEmployeeShift)

		// Locations (Geofence) (admin-config)
		att.POST("/locations", requireAttendanceSettings("create"), handler.CreateLocation)
		att.GET("/locations", requireAttendanceSettings("view"), handler.ListLocations)
		att.GET("/locations/:id", requireAttendanceSettings("view"), handler.GetLocationByID)
		att.PUT("/locations/:id", requireAttendanceSettings("update"), handler.UpdateLocation)
		att.DELETE("/locations/:id", requireAttendanceSettings("delete"), handler.DeleteLocation)

		// Events (Check-in / Check-out)
		att.POST("/events", handler.CreateEvent)
		att.GET("/events", handler.ListEvents)
		att.GET("/events/:id", handler.GetEventByID)

		// Sessions (Daily Work Sessions)
		att.GET("/sessions", handler.ListSessions)
		att.GET("/sessions/detail", handler.GetSession)
		att.GET("/calendar", handler.GetEmployeeCalendar)
		att.GET("/summary", handler.GetEmployeeSummary)
		att.GET("/reports/sessions", requireAttendanceReport("view"), handler.GetAttendanceReport)

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

		// Exempt Positions (admin-config)
		att.POST("/exempt-positions", requireAttendanceSettings("create"), handler.CreateExemptPosition)
		att.GET("/exempt-positions", requireAttendanceSettings("view"), handler.ListExemptPositions)
		att.GET("/exempt-positions/:id", requireAttendanceSettings("view"), handler.GetExemptPositionByID)
		att.PUT("/exempt-positions/:id", requireAttendanceSettings("update"), handler.UpdateExemptPosition)
		att.DELETE("/exempt-positions/:id", requireAttendanceSettings("delete"), handler.DeleteExemptPosition)

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
		att.PUT("/business-travels/:id/participants/:participantId", handler.UpdateBusinessTravelParticipant)
		att.DELETE("/business-travels/:id/participants/:participantId", handler.DeleteBusinessTravelParticipant)
		att.POST("/business-travels/:id/destinations", handler.AddBusinessTravelDestination)
		att.GET("/business-travels/:id/destinations", handler.ListBusinessTravelDestinations)
		att.PUT("/business-travels/:id/destinations/:destinationId", handler.UpdateBusinessTravelDestination)
		att.DELETE("/business-travels/:id/destinations/:destinationId", handler.DeleteBusinessTravelDestination)
		att.POST("/business-travels/:id/activities", handler.AddBusinessTravelActivity)
		att.GET("/business-travels/:id/activities", handler.ListBusinessTravelActivities)
		att.PUT("/business-travels/:id/activities/:activityId", handler.UpdateBusinessTravelActivity)
		att.DELETE("/business-travels/:id/activities/:activityId", handler.DeleteBusinessTravelActivity)
		att.POST("/business-travels/:id/schedules", handler.AddBusinessTravelSchedule)
		att.GET("/business-travels/:id/schedules", handler.ListBusinessTravelSchedules)
		att.PUT("/business-travels/:id/schedules/:scheduleId", handler.UpdateBusinessTravelSchedule)
		att.DELETE("/business-travels/:id/schedules/:scheduleId", handler.DeleteBusinessTravelSchedule)

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
		// Flat lookup by settlement ID alone, for callers that only know the
		// approval instance's document_id (the settlement ID) and not its
		// parent travel — e.g. the Approvals module's task detail popup.
		att.GET("/business-travel-settlements/:settlementId", handler.GetSettlementByID)
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
