package careerintelligence

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Career Intelligence & Talent Management endpoints.
func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	ci := r.Group("/career-intelligence")
	{
		// Talent Maps — static paths BEFORE dynamic :id routes (Gin constraint)
		ci.GET("/talent-maps", handler.ListTalentMaps)
		ci.POST("/talent-maps", handler.CreateTalentMap)
		ci.POST("/talent-maps/generate", handler.GenerateTalentMap)
		ci.GET("/talent-maps/settings", handler.GetTalentMapSettings)
		ci.PUT("/talent-maps/settings", handler.UpdateTalentMapSettings)
		ci.GET("/talent-maps/grid", handler.GetTalentGrid)
		ci.GET("/talent-maps/employee/:employeeId", handler.GetEmployeeTalentProfile)
		ci.GET("/employees/:employeeId/training-profile", handler.GetEmployeeTrainingProfile)
		ci.GET("/employees/:employeeId/training-recommendations", handler.GetTrainingRecommendations)
		ci.GET("/talent-maps/:id", handler.GetTalentMapByID)
		ci.PUT("/talent-maps/:id", handler.UpdateTalentMap)
		ci.DELETE("/talent-maps/:id", handler.DeleteTalentMap)

		// Career Interests (3 endpoints)
		ci.GET("/interests", handler.ListCareerInterests)
		ci.POST("/interests", handler.CreateCareerInterest)
		ci.GET("/interests/employee/:employeeId", handler.GetEmployeeCareerInterests)

		// Career Paths — strategical planning (enhancement plan §12.9).
		// Ladder-style (name + steps[]) — endpoint utama FE Career Paths.
		ci.GET("/paths", handler.ListCareerPaths)
		ci.POST("/paths", handler.CreateCareerPathLadder)
		// Static path SEBELUM dynamic :id (Gin constraint).
		ci.GET("/paths/gap-analysis", handler.GetGapAnalysis)
		ci.GET("/paths/:id", handler.GetCareerPathByID)
		// S-4: internal candidate eligibility — Gin :id match satu segmen,
		// jadi /paths/:id/eligible-employees tidak konflik dengan /paths/:id.
		ci.GET("/paths/:id/eligible-employees", handler.GetEligibleEmployees)
		ci.PUT("/paths/:id", handler.UpdateCareerPath)
		ci.DELETE("/paths/:id", handler.DeleteCareerPath)

		// Succession Plans (5 endpoints + S-5 gaps)
		ci.GET("/successions", handler.ListSuccessionPlans)
		ci.POST("/successions", handler.CreateSuccessionPlan)
		// S-5: succession gap (posisi kunci tanpa successor siap) — static
		// path SEBELUM dynamic :id (Gin constraint).
		ci.GET("/successions/gaps", handler.GetSuccessionGaps)
		ci.GET("/successions/:id", handler.GetSuccessionPlanByID)
		ci.PUT("/successions/:id", handler.UpdateSuccessionPlan)
		ci.DELETE("/successions/:id", handler.DeleteSuccessionPlan)
	}
}
