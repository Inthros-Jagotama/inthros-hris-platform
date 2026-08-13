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
		// G-9 sub-project 1: requisition requirements + competencies (fondasi candidate matching)
		rec.POST("/requisitions/:id/requirements", handler.CreateRequisitionRequirement)
		rec.GET("/requisitions/:id/requirements", handler.ListRequisitionRequirements)
		rec.PUT("/requirements/:id", handler.UpdateRequisitionRequirement)
		rec.DELETE("/requirements/:id", handler.DeleteRequisitionRequirement)
		rec.POST("/requisitions/:id/competencies", handler.CreateRequisitionCompetency)
		rec.GET("/requisitions/:id/competencies", handler.ListRequisitionCompetencies)
		rec.PUT("/requisition-competencies/:id", handler.UpdateRequisitionCompetency)
		rec.DELETE("/requisition-competencies/:id", handler.DeleteRequisitionCompetency)

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
		rec.POST("/candidates/:id/educations", handler.CreateCandidateEducation)
		rec.GET("/candidates/:id/educations", handler.ListCandidateEducations)
		rec.PUT("/educations/:id", handler.UpdateCandidateEducation)
		rec.DELETE("/educations/:id", handler.DeleteCandidateEducation)
		rec.POST("/candidates/:id/work-experiences", handler.CreateCandidateWorkExperience)
		rec.GET("/candidates/:id/work-experiences", handler.ListCandidateWorkExperiences)
		rec.PUT("/work-experiences/:id", handler.UpdateCandidateWorkExperience)
		rec.DELETE("/work-experiences/:id", handler.DeleteCandidateWorkExperience)
		rec.POST("/candidates/:id/skills", handler.CreateCandidateSkill)
		rec.GET("/candidates/:id/skills", handler.ListCandidateSkills)
		rec.PUT("/skills/:id", handler.UpdateCandidateSkill)
		rec.DELETE("/skills/:id", handler.DeleteCandidateSkill)
		rec.POST("/candidates/:id/certifications", handler.CreateCandidateCertification)
		rec.GET("/candidates/:id/certifications", handler.ListCandidateCertifications)
		rec.PUT("/certifications/:id", handler.UpdateCandidateCertification)
		rec.DELETE("/certifications/:id", handler.DeleteCandidateCertification)
		rec.POST("/candidates/:id/documents", handler.CreateCandidateDocument)
		rec.GET("/candidates/:id/documents", handler.ListCandidateDocuments)
		rec.PUT("/documents/:id", handler.UpdateCandidateDocument)
		rec.DELETE("/documents/:id", handler.DeleteCandidateDocument)
		rec.POST("/candidates/:id/consents", handler.CreateCandidateConsent)
		rec.GET("/candidates/:id/consents", handler.ListCandidateConsents)

		// S-4: internal candidate eligibility (CI → Recruitment)
		rec.GET("/eligible-internal-candidates", handler.GetEligibleInternalCandidates)

		// Job Applications
		rec.POST("/applications", handler.CreateApplication)
		rec.GET("/applications", handler.ListApplications)
		rec.GET("/applications/:id", handler.GetApplicationByID)
		rec.PUT("/applications/:id/status", handler.UpdateApplicationStatus)
		rec.GET("/applications/:id/history", handler.GetApplicationHistory)
		// G-9 sub-project 2: candidate match score (advisory, computed on-the-fly)
		rec.GET("/applications/:id/match-score", handler.GetCandidateMatchScore)
		rec.DELETE("/applications/:id", handler.DeleteApplication)
		// G-7 sub-project 1: screening (many-per-application, no auto status transition)
		rec.POST("/applications/:id/screenings", handler.CreateApplicationScreening)
		rec.GET("/applications/:id/screenings", handler.ListApplicationScreenings)
		rec.PUT("/screenings/:id", handler.UpdateApplicationScreening)
		rec.DELETE("/screenings/:id", handler.DeleteApplicationScreening)

			// G-7 sub-project 2: assessments (batch session) + participants
			rec.POST("/assessments", handler.CreateAssessment)
			rec.GET("/assessments", handler.ListAssessments)
			rec.GET("/assessments/:id", handler.GetAssessmentByID)
			rec.PUT("/assessments/:id", handler.UpdateAssessment)
			rec.DELETE("/assessments/:id", handler.DeleteAssessment)
			rec.POST("/assessments/:id/participants", handler.AddAssessmentParticipant)
			rec.GET("/assessments/:id/participants", handler.ListAssessmentParticipants)
			rec.PUT("/assessment-participants/:id", handler.UpdateAssessmentParticipant)
			rec.DELETE("/assessment-participants/:id", handler.DeleteAssessmentParticipant)

		// Interviews
		rec.POST("/interviews", handler.CreateInterview)
		rec.GET("/interviews", handler.ListInterviews)
		rec.GET("/interviews/:id", handler.GetInterviewByID)
		rec.PUT("/interviews/:id", handler.UpdateInterview)
		rec.DELETE("/interviews/:id", handler.DeleteInterview)
		// G-8: multi-interviewer + scorecard
		rec.POST("/interviews/:id/interviewers", handler.AddInterviewer)
		rec.GET("/interviews/:id/interviewers", handler.ListInterviewers)
		rec.DELETE("/interviewers/:id", handler.RemoveInterviewer)
		rec.POST("/interviews/:id/scorecard-items", handler.AddScorecardItem)
		rec.GET("/interviews/:id/scorecard-items", handler.ListScorecardItems)
		rec.PUT("/scorecard-items/:id", handler.UpdateScorecardItem)
		rec.DELETE("/scorecard-items/:id", handler.DeleteScorecardItem)
		rec.POST("/interviews/:id/complete", handler.CompleteInterview)

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
