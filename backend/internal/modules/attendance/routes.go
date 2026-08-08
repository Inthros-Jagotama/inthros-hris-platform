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

		// Overtime Requests
		att.POST("/overtime-requests", handler.CreateOvertimeRequest)
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
	}
}
