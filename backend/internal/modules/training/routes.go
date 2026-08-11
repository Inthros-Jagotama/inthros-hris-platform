package training

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	trn := rg.Group("/trainings")
	{
		// Training Categories
		trn.POST("/categories", handler.CreateCategory)
		trn.GET("/categories", handler.ListCategories)
		trn.GET("/categories/:id", handler.GetCategoryByID)
		trn.PUT("/categories/:id", handler.UpdateCategory)
		trn.DELETE("/categories/:id", handler.DeleteCategory)

		// Training Courses
		trn.POST("/courses", handler.CreateCourse)
		trn.GET("/courses", handler.ListCourses)
		trn.GET("/courses/:id", handler.GetCourseByID)
		trn.PUT("/courses/:id", handler.UpdateCourse)
		trn.DELETE("/courses/:id", handler.DeleteCourse)

		// Training Providers (P0-BE — plan §11)
		trn.POST("/providers", handler.CreateProvider)
		trn.GET("/providers", handler.ListProviders)
		trn.GET("/providers/:id", handler.GetProviderByID)
		trn.PUT("/providers/:id", handler.UpdateProvider)
		trn.DELETE("/providers/:id", handler.DeleteProvider)

		// Training Trainers (P0-BE — plan §12)
		trn.POST("/trainers", handler.CreateTrainer)
		trn.GET("/trainers", handler.ListTrainers)
		trn.GET("/trainers/:id", handler.GetTrainerByID)
		trn.PUT("/trainers/:id", handler.UpdateTrainer)
		trn.DELETE("/trainers/:id", handler.DeleteTrainer)

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

		// Training Certificates
		trn.POST("/certificates", handler.CreateCertificate)
		trn.GET("/certificates", handler.ListCertificates)
		trn.GET("/certificates/:id", handler.GetCertificateByID)
		trn.PUT("/certificates/:id", handler.UpdateCertificate)
		trn.DELETE("/certificates/:id", handler.DeleteCertificate)
	}
}
