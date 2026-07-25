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

		// Training Sessions
		trn.POST("/sessions", handler.CreateSession)
		trn.GET("/sessions", handler.ListSessions)
		trn.GET("/sessions/:id", handler.GetSessionByID)
		trn.PUT("/sessions/:id", handler.UpdateSession)
		trn.PUT("/sessions/:id/status", handler.UpdateSessionStatus)
		trn.DELETE("/sessions/:id", handler.DeleteSession)

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
