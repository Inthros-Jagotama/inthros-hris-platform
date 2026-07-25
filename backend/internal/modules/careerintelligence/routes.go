package careerintelligence

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Career Intelligence & Talent Management endpoints.
func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	ci := r.Group("/career-intelligence")
	{
		// Talent Maps — static paths BEFORE dynamic :id routes (Gin constraint)
		ci.GET("/talent-maps", handler.ListTalentMaps)
		ci.POST("/talent-maps", handler.CreateTalentMap)
		ci.GET("/talent-maps/grid", handler.GetTalentGrid)
		ci.GET("/talent-maps/employee/:employeeId", handler.GetEmployeeTalentProfile)
		ci.GET("/talent-maps/:id", handler.GetTalentMapByID)
		ci.PUT("/talent-maps/:id", handler.UpdateTalentMap)
		ci.DELETE("/talent-maps/:id", handler.DeleteTalentMap)

		// Career Interests (3 endpoints)
		ci.GET("/interests", handler.ListCareerInterests)
		ci.POST("/interests", handler.CreateCareerInterest)
		ci.GET("/interests/employee/:employeeId", handler.GetEmployeeCareerInterests)

		// Career Paths (3 endpoints + gap analysis)
		ci.GET("/paths", handler.ListCareerPaths)
		ci.POST("/paths", handler.CreateCareerPath)
		ci.DELETE("/paths/:id", handler.DeleteCareerPath)
		ci.GET("/paths/gap-analysis", handler.GetGapAnalysis)

		// Succession Plans (5 endpoints)
		ci.GET("/successions", handler.ListSuccessionPlans)
		ci.POST("/successions", handler.CreateSuccessionPlan)
		ci.GET("/successions/:id", handler.GetSuccessionPlanByID)
		ci.PUT("/successions/:id", handler.UpdateSuccessionPlan)
		ci.DELETE("/successions/:id", handler.DeleteSuccessionPlan)
	}
}
