package training

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// requireTrainingSettings menggating endpoint yang MENGUBAH (create/update/
// delete) master training (categories, courses, providers, trainers beserta
// sub-resource course) — supaya HANYA permission submenu "training.settings.
// <action>" yang berlaku (bukan module-level "training.update"). Middleware
// RBAC global (authz.NewMiddleware) menganggap module-level otomatis mencakup
// semua submenu-nya — cocok untuk kebanyakan resource, tapi tidak untuk aksi
// admin-config ini, yang harus benar-benar terpisah dari sekadar melihat
// training sehari-hari ("training.view", yang dimiliki hampir semua role
// termasuk Employee default). Pola sama seperti requireLeaveSettings /
// requireAttendanceSettings.
//
// GET (list & detail) SENGAJA TIDAK dibatasi middleware ini — dipakai halaman
// Training utama (stat) dan dropdown halaman lain (needs, participants,
// requests, sessions) yang harus tetap bisa diakses siapa pun dengan
// "training.view".
func requireTrainingSettings(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := "training.settings." + action
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
				"message": "You don't have permission to manage training settings",
				"details": gin.H{"required": required},
			},
		})
	}
}

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	trn := rg.Group("/trainings")
	{
		// Training Categories (admin-config)
		trn.POST("/categories", requireTrainingSettings("create"), handler.CreateCategory)
		trn.GET("/categories", handler.ListCategories)
		trn.GET("/categories/:id", handler.GetCategoryByID)
		trn.PUT("/categories/:id", requireTrainingSettings("update"), handler.UpdateCategory)
		trn.DELETE("/categories/:id", requireTrainingSettings("delete"), handler.DeleteCategory)

		// Training Courses (admin-config)
		trn.POST("/courses", requireTrainingSettings("create"), handler.CreateCourse)
		trn.GET("/courses", handler.ListCourses)
		trn.GET("/courses/:id", handler.GetCourseByID)
		trn.PUT("/courses/:id", requireTrainingSettings("update"), handler.UpdateCourse)
		trn.DELETE("/courses/:id", requireTrainingSettings("delete"), handler.DeleteCourse)

		// Course sub-resources (P1-BE — plan §8/§9/§10) — mutasi ikut admin-config
		trn.GET("/courses/:id/objectives", handler.ListCourseObjectives)
		trn.POST("/courses/:id/objectives", requireTrainingSettings("create"), handler.CreateCourseObjective)
		trn.PUT("/course-objectives/:id", requireTrainingSettings("update"), handler.UpdateCourseObjective)
		trn.DELETE("/course-objectives/:id", requireTrainingSettings("delete"), handler.DeleteCourseObjective)

		trn.GET("/courses/:id/competencies", handler.ListCourseCompetencies)
		trn.POST("/courses/:id/competencies", requireTrainingSettings("create"), handler.CreateCourseCompetency)
		trn.DELETE("/course-competencies/:id", requireTrainingSettings("delete"), handler.DeleteCourseCompetency)

		trn.GET("/courses/:id/prerequisites", handler.ListCoursePrerequisites)
		trn.POST("/courses/:id/prerequisites", requireTrainingSettings("create"), handler.CreateCoursePrerequisite)
		trn.DELETE("/course-prerequisites/:id", requireTrainingSettings("delete"), handler.DeleteCoursePrerequisite)

		// Training Providers (P0-BE — plan §11) (admin-config)
		trn.POST("/providers", requireTrainingSettings("create"), handler.CreateProvider)
		trn.GET("/providers", handler.ListProviders)
		trn.GET("/providers/:id", handler.GetProviderByID)
		trn.PUT("/providers/:id", requireTrainingSettings("update"), handler.UpdateProvider)
		trn.DELETE("/providers/:id", requireTrainingSettings("delete"), handler.DeleteProvider)

		// Training Trainers (P0-BE — plan §12) (admin-config)
		trn.POST("/trainers", requireTrainingSettings("create"), handler.CreateTrainer)
		trn.GET("/trainers", handler.ListTrainers)
		trn.GET("/trainers/:id", handler.GetTrainerByID)
		trn.PUT("/trainers/:id", requireTrainingSettings("update"), handler.UpdateTrainer)
		trn.DELETE("/trainers/:id", requireTrainingSettings("delete"), handler.DeleteTrainer)

		// Training Sessions
		trn.POST("/sessions", handler.CreateSession)
		trn.GET("/sessions", handler.ListSessions)
		trn.GET("/sessions/:id", handler.GetSessionByID)
		trn.PUT("/sessions/:id", handler.UpdateSession)
		trn.PUT("/sessions/:id/status", handler.UpdateSessionStatus)
		trn.DELETE("/sessions/:id", handler.DeleteSession)

		// Session Trainers (P0-BE — plan §13)
		trn.GET("/sessions/:id/trainers", handler.ListSessionTrainers)
		trn.POST("/sessions/:id/trainers", handler.AddSessionTrainer)
		trn.DELETE("/session-trainers/:id", handler.RemoveSessionTrainer)

		// Attendance (P0-BE — plan §19)
		trn.GET("/sessions/:id/attendance", handler.ListAttendanceBySession)
		trn.POST("/sessions/:id/attendance", handler.MarkAttendance)
		trn.PUT("/attendances/:id", handler.UpdateAttendance)

		// Assessments (P0-BE — plan §21)
		trn.GET("/sessions/:id/assessments", handler.ListAssessmentsBySession)
		trn.POST("/sessions/:id/assessments", handler.CreateAssessment)
		trn.POST("/assessments/:id/results", handler.SubmitAssessmentResult)

		// Session Costs (P1-BE — plan §26)
		trn.GET("/sessions/:id/costs", handler.ListSessionCosts)
		trn.POST("/sessions/:id/costs", handler.CreateSessionCost)
		trn.PUT("/session-costs/:id", handler.UpdateSessionCost)
		trn.DELETE("/session-costs/:id", handler.DeleteSessionCost)

		// Session Documents (P1-BE — plan §27)
		trn.GET("/sessions/:id/documents", handler.ListDocuments)
		trn.POST("/sessions/:id/documents", handler.CreateDocument)
		trn.DELETE("/documents/:id", handler.DeleteDocument)

		// Training Participants
		trn.POST("/participants", handler.CreateParticipant)
		trn.GET("/participants", handler.ListParticipants)
		trn.GET("/participants/:id", handler.GetParticipantByID)
		trn.PUT("/participants/:id", handler.UpdateParticipant)
		trn.DELETE("/participants/:id", handler.DeleteParticipant)

		// Training Materials
		trn.POST("/materials", handler.CreateMaterial)
		trn.GET("/materials", handler.ListMaterials)
		trn.PUT("/materials/:id", handler.UpdateMaterial)
		trn.DELETE("/materials/:id", handler.DeleteMaterial)

		// Training Evaluations
		trn.POST("/evaluations", handler.CreateEvaluation)
		trn.GET("/evaluations", handler.ListEvaluations)
		trn.GET("/evaluations/:id", handler.GetEvaluationByID)
		trn.PUT("/evaluations/:id", handler.UpdateEvaluation)
		trn.DELETE("/evaluations/:id", handler.DeleteEvaluation)

		// Training Certificates (P2-BE: generate + update file URL)
		trn.POST("/certificates", handler.CreateCertificate)
		trn.GET("/certificates", handler.ListCertificates)
		trn.GET("/certificates/:id", handler.GetCertificateByID)
		trn.PUT("/certificates/:id", handler.UpdateCertificateFile)
		trn.DELETE("/certificates/:id", handler.DeleteCertificate)
		trn.POST("/participants/:id/certificate", handler.GenerateCertificate)

		// Training Plans (P1-BE — plan §16)
		trn.POST("/plans", handler.CreatePlan)
		trn.GET("/plans", handler.ListPlans)
		trn.GET("/plans/:id", handler.GetPlanByID)
		trn.PUT("/plans/:id", handler.UpdatePlan)
		trn.DELETE("/plans/:id", handler.DeletePlan)

		// Plan items (P1-BE — plan §16)
		trn.GET("/plans/:id/items", handler.ListPlanItems)
		trn.POST("/plans/:id/items", handler.CreatePlanItem)
		trn.PUT("/plan-items/:id", handler.UpdatePlanItem)
		trn.DELETE("/plan-items/:id", handler.DeletePlanItem)

		// Training Needs (P1-BE — plan §17)
		trn.POST("/needs", handler.CreateNeed)
		trn.GET("/needs", handler.ListNeeds)
		trn.GET("/needs/:id", handler.GetNeedByID)
		trn.PUT("/needs/:id", handler.UpdateNeed)
		trn.DELETE("/needs/:id", handler.DeleteNeed)

		// Training Requests (P1-BE — plan §15, Central Approval)
		trn.POST("/requests", handler.CreateRequest)
		trn.GET("/requests", handler.ListRequests)
		trn.GET("/requests/:id", handler.GetRequestByID)
		trn.POST("/requests/:id/submit", handler.SubmitRequest)
		trn.POST("/requests/:id/cancel", handler.CancelRequest)

		// Training Mandatories (P1-BE — plan §25)
		trn.POST("/mandatories", handler.CreateMandatory)
		trn.GET("/mandatories", handler.ListMandatories)
		trn.GET("/mandatories/:id", handler.GetMandatoryByID)
		trn.PUT("/mandatories/:id", handler.UpdateMandatory)
		trn.DELETE("/mandatories/:id", handler.DeleteMandatory)

		// ── P2-BE: Evaluation Forms (plan §22) ──
		// Route statis + sub-resource session form SEBELUM parameterized flat.
		trn.GET("/evaluation-forms", handler.ListEvaluationForms)
		trn.POST("/evaluation-forms", handler.CreateEvaluationForm)
		trn.GET("/evaluation-forms/:form_id", handler.GetEvaluationForm)
		trn.PUT("/evaluation-forms/:form_id", handler.UpdateEvaluationForm)
		trn.DELETE("/evaluation-forms/:form_id", handler.DeleteEvaluationForm)
		trn.GET("/sessions/:id/evaluation-form", handler.GetEvaluationFormBySession)

		// Evaluation Questions (P2-BE — plan §22)
		trn.GET("/evaluation-forms/:form_id/questions", handler.ListEvaluationQuestions)
		trn.POST("/evaluation-forms/:form_id/questions", handler.CreateEvaluationQuestion)
		trn.PUT("/evaluation-questions/:id", handler.UpdateEvaluationQuestion)
		trn.DELETE("/evaluation-questions/:id", handler.DeleteEvaluationQuestion)

		// Evaluation Answers (P2-BE — plan §22)
		trn.GET("/evaluation-answers", handler.ListEvaluationAnswers)
		trn.POST("/evaluation-forms/:form_id/participants/:participant_id/answers", handler.SubmitEvaluationAnswers)

		// Effectiveness Assessments (P2-BE — plan §23)
		trn.POST("/effectiveness", handler.CreateEffectivenessAssessment)
		trn.GET("/effectiveness", handler.ListEffectivenessAssessments)
		trn.GET("/effectiveness/:id", handler.GetEffectivenessAssessment)
		trn.PUT("/effectiveness/:id", handler.UpdateEffectivenessAssessment)
		trn.DELETE("/effectiveness/:id", handler.DeleteEffectivenessAssessment)

		// Certifications master (P2-BE — plan §24)
		trn.POST("/certifications", handler.CreateCertification)
		trn.GET("/certifications", handler.ListCertifications)
		trn.GET("/certifications/:id", handler.GetCertification)
		trn.PUT("/certifications/:id", handler.UpdateCertification)
		trn.DELETE("/certifications/:id", handler.DeleteCertification)

		// ── P2-BE: Reports & History (route statis SEBELUM /:id — konstrain Gin) ──
		trn.GET("/history", handler.GetTrainingHistory)
		trn.GET("/reports/participation", handler.GetParticipationReport)
		trn.GET("/reports/cost", handler.GetCostReport)
		trn.GET("/reports/compliance", handler.GetComplianceReport)
		trn.GET("/reports/dashboard", handler.GetDashboardReport)
	}
}
