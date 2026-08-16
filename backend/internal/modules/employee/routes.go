package employee

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Permission khusus halaman Sensitive Field Settings (toggle encrypt-at-rest
// per field, cakupan tenant-wide). Di-seed oleh migrasi tenant 154.
const (
	PermSensitiveFieldsView   = "setting.sensitive-fields.view"
	PermSensitiveFieldsManage = "setting.sensitive-fields.manage"
)

// requireSensitiveFieldSettings menggating endpoint setelan field sensitif.
//
// Middleware RBAC global menurunkan resource dari path, jadi
// /employees/settings/sensitive-fields hanya tercek sebagai employee.view /
// employee.update — terlalu longgar untuk aksi setingkat admin yang mengubah
// enkripsi at-rest seluruh tenant. Gate tambahan di sini memakai permission
// khusus, dengan fallback module-level "setting.*" — mengikuti aturan yang
// sudah berlaku di authz.NewMiddleware, di mana claim level-module mencakup
// semua submenu-nya (jadi platform role seperti company_admin yang membawa
// claim "setting.view"/"setting.update" tetap bisa masuk).
//
// moduleFallback: claim level-module yang setara dengan aksi ini.
func requireSensitiveFieldSettings(permission, moduleFallback string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range c.GetStringSlice("permissions") {
			if p == "*" || p == permission || p == moduleFallback {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You don't have permission to manage sensitive field settings",
				"details": gin.H{"required": permission},
			},
		})
	}
}

// RegisterRoutes mendaftarkan semua endpoint Employee ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/employees
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	emps := rg.Group("/employees")
	{
		// Employee CRUD
		emps.POST("", handler.Create)
		emps.GET("", handler.List)
		emps.GET("/:id", handler.GetByID)
		emps.PUT("/:id", handler.Update)
		emps.PUT("/:id/photo", handler.UploadPhoto)
		emps.DELETE("/:id/photo", handler.DeletePhoto)
		emps.DELETE("/:id", handler.Delete)

		// Addresses
		emps.POST("/:id/addresses", handler.CreateAddress)
		emps.PUT("/:id/addresses/:addressId", handler.UpdateAddress)
		emps.DELETE("/:id/addresses/:addressId", handler.DeleteAddress)

		// Emergency Contacts
		emps.POST("/:id/emergency-contacts", handler.CreateEmergencyContact)
		emps.PUT("/:id/emergency-contacts/:contactId", handler.UpdateEmergencyContact)
		emps.DELETE("/:id/emergency-contacts/:contactId", handler.DeleteEmergencyContact)

		// Families
		emps.POST("/:id/families", handler.CreateFamily)
		emps.PUT("/:id/families/:familyId", handler.UpdateFamily)
		emps.DELETE("/:id/families/:familyId", handler.DeleteFamily)

		// Educations
		emps.POST("/:id/educations", handler.CreateEducation)
		emps.PUT("/:id/educations/:educationId", handler.UpdateEducation)
		emps.DELETE("/:id/educations/:educationId", handler.DeleteEducation)

		// Experiences
		emps.POST("/:id/experiences", handler.CreateExperience)
		emps.PUT("/:id/experiences/:experienceId", handler.UpdateExperience)
		emps.DELETE("/:id/experiences/:experienceId", handler.DeleteExperience)

		// Documents
		emps.POST("/:id/documents", handler.CreateDocument)
		emps.POST("/:id/documents/upload", handler.UploadDocumentFile)
		emps.PUT("/:id/documents/:documentId", handler.UpdateDocument)
		emps.PUT("/:id/documents/:documentId/upload", handler.UpdateDocumentFile)
		emps.DELETE("/:id/documents/:documentId", handler.DeleteDocument)

		// Insurances
		emps.POST("/:id/insurances", handler.CreateInsurance)
		emps.PUT("/:id/insurances/:insuranceId", handler.UpdateInsurance)
		emps.DELETE("/:id/insurances/:insuranceId", handler.DeleteInsurance)

		// Banks
		emps.POST("/:id/banks", handler.CreateBank)
		emps.PUT("/:id/banks/:bankId", handler.UpdateBank)
		emps.DELETE("/:id/banks/:bankId", handler.DeleteBank)

		// Employments
		emps.POST("/:id/employments", handler.CreateEmployment)
		emps.PUT("/:id/employments/:employmentId", handler.UpdateEmployment)
		emps.DELETE("/:id/employments/:employmentId", handler.DeleteEmployment)

		// Sensitive field settings (encryption toggle, admin only)
		emps.GET("/settings/sensitive-fields",
			requireSensitiveFieldSettings(PermSensitiveFieldsView, "setting.view"),
			handler.ListSensitiveFieldSettings)
		emps.PUT("/settings/sensitive-fields/:fieldKey",
			requireSensitiveFieldSettings(PermSensitiveFieldsManage, "setting.update"),
			handler.SetSensitiveFieldEnabled)
	}
}
