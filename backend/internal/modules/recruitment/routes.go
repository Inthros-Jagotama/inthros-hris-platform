package recruitment

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Recruitment & Onboarding ke router group tenant.
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	rec := rg.Group("/recruitment")
	{
		// Job Requisitions
		rec.POST("/requisitions", handler.CreateRequisition)
		rec.GET("/requisitions", handler.ListRequisitions)
		rec.GET("/requisitions/:id", handler.GetRequisitionByID)
		rec.PUT("/requisitions/:id", handler.UpdateRequisition)
		rec.DELETE("/requisitions/:id", handler.DeleteRequisition)
		// G-1: submit requisition draft ke Central Approval (single approval path)
		rec.POST("/requisitions/:id/submit", handler.SubmitRequisition)

		// Job Offers (G-3)
		rec.POST("/offers", handler.CreateOffer)
		rec.GET("/offers", handler.ListOffers)
		rec.GET("/offers/:id", handler.GetOfferByID)
		rec.PUT("/offers/:id", handler.UpdateOffer)
		rec.DELETE("/offers/:id", handler.DeleteOffer)
		// G-3: workflow offer — submit ke Central Approval + transisi status
		rec.POST("/offers/:id/submit", handler.SubmitOffer)
		rec.POST("/offers/:id/send", handler.SendOffer)
		rec.POST("/offers/:id/accept", handler.AcceptOffer)
		rec.POST("/offers/:id/reject", handler.RejectOffer)
		rec.POST("/offers/:id/withdraw", handler.WithdrawOffer)

		// Candidates
		rec.POST("/candidates", handler.CreateCandidate)
		rec.GET("/candidates", handler.ListCandidates)
		rec.GET("/candidates/:id", handler.GetCandidateByID)
		rec.PUT("/candidates/:id", handler.UpdateCandidate)
		rec.DELETE("/candidates/:id", handler.DeleteCandidate)

		// S-4: internal candidate eligibility (CI → Recruitment)
		rec.GET("/eligible-internal-candidates", handler.GetEligibleInternalCandidates)

		// Job Applications
		rec.POST("/applications", handler.CreateApplication)
		rec.GET("/applications", handler.ListApplications)
		rec.GET("/applications/:id", handler.GetApplicationByID)
		rec.PUT("/applications/:id/status", handler.UpdateApplicationStatus)
		rec.DELETE("/applications/:id", handler.DeleteApplication)

		// Interviews
		rec.POST("/interviews", handler.CreateInterview)
		rec.GET("/interviews", handler.ListInterviews)
		rec.GET("/interviews/:id", handler.GetInterviewByID)
		rec.PUT("/interviews/:id", handler.UpdateInterview)
		rec.DELETE("/interviews/:id", handler.DeleteInterview)

		// Onboarding Task Templates
		rec.POST("/onboarding-task-templates", handler.CreateOnboardingTaskTemplate)
		rec.GET("/onboarding-task-templates", handler.ListOnboardingTaskTemplates)
		rec.PUT("/onboarding-task-templates/:id", handler.UpdateOnboardingTaskTemplate)
		rec.DELETE("/onboarding-task-templates/:id", handler.DeleteOnboardingTaskTemplate)

		// Employee Onboarding
		rec.POST("/employee-onboardings", handler.CreateEmployeeOnboarding)
		rec.GET("/employee-onboardings", handler.ListEmployeeOnboardings)
		rec.GET("/employee-onboardings/:id", handler.GetEmployeeOnboardingByID)
		rec.PUT("/employee-onboardings/:id", handler.UpdateEmployeeOnboarding)
		rec.DELETE("/employee-onboardings/:id", handler.DeleteEmployeeOnboarding)

		// Onboarding Task Items (nested under onboardings)
		rec.POST("/onboarding-task-items", handler.CreateOnboardingTaskItem)
		rec.GET("/employee-onboardings/:id/task-items", handler.ListOnboardingTaskItems)
		rec.PUT("/onboarding-task-items/:id", handler.UpdateOnboardingTaskItem)
		rec.DELETE("/onboarding-task-items/:id", handler.DeleteOnboardingTaskItem)
	}
}
