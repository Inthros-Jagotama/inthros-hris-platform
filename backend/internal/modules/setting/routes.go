package setting

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	settings := rg.Group("/settings")
	{
		// Zones
		zones := settings.Group("/zones")
		{
			zones.POST("", handler.CreateZone)
			zones.GET("", handler.ListZones)
			zones.GET("/:id", handler.GetZoneByID)
			zones.PUT("/:id", handler.UpdateZone)
			zones.DELETE("/:id", handler.DeleteZone)
		}
		// Provinces
		provinces := settings.Group("/provinces")
		{
			provinces.POST("", handler.CreateProvince)
			provinces.GET("", handler.ListProvinces)
			provinces.GET("/:id", handler.GetProvinceByID)
			provinces.PUT("/:id", handler.UpdateProvince)
			provinces.DELETE("/:id", handler.DeleteProvince)
		}
		// Regencies
		regencies := settings.Group("/regencies")
		{
			regencies.POST("", handler.CreateRegency)
			regencies.GET("", handler.ListRegencies)
			regencies.GET("/:id", handler.GetRegencyByID)
			regencies.PUT("/:id", handler.UpdateRegency)
			regencies.DELETE("/:id", handler.DeleteRegency)
		}
		// Districts
		districts := settings.Group("/districts")
		{
			districts.POST("", handler.CreateDistrict)
			districts.GET("", handler.ListDistricts)
			districts.GET("/:id", handler.GetDistrictByID)
			districts.PUT("/:id", handler.UpdateDistrict)
			districts.DELETE("/:id", handler.DeleteDistrict)
		}
		// Villages
		villages := settings.Group("/villages")
		{
			villages.POST("", handler.CreateVillage)
			villages.GET("", handler.ListVillages)
			villages.GET("/:id", handler.GetVillageByID)
			villages.PUT("/:id", handler.UpdateVillage)
			villages.DELETE("/:id", handler.DeleteVillage)
		}
		// Education
		educations := settings.Group("/educations")
		{
			educations.POST("", handler.CreateEducation)
			educations.GET("", handler.ListEducations)
			educations.GET("/:id", handler.GetEducationByID)
			educations.PUT("/:id", handler.UpdateEducation)
			educations.DELETE("/:id", handler.DeleteEducation)
		}
		// Religion
		religions := settings.Group("/religions")
		{
			religions.POST("", handler.CreateReligion)
			religions.GET("", handler.ListReligions)
			religions.GET("/:id", handler.GetReligionByID)
			religions.PUT("/:id", handler.UpdateReligion)
			religions.DELETE("/:id", handler.DeleteReligion)
		}
		// MaritalStatus
		maritalStatuses := settings.Group("/marital-statuses")
		{
			maritalStatuses.POST("", handler.CreateMaritalStatus)
			maritalStatuses.GET("", handler.ListMaritalStatuses)
			maritalStatuses.GET("/:id", handler.GetMaritalStatusByID)
			maritalStatuses.PUT("/:id", handler.UpdateMaritalStatus)
			maritalStatuses.DELETE("/:id", handler.DeleteMaritalStatus)
		}
		// RelationshipType
		relationshipTypes := settings.Group("/relationship-types")
		{
			relationshipTypes.POST("", handler.CreateRelationshipType)
			relationshipTypes.GET("", handler.ListRelationshipTypes)
			relationshipTypes.GET("/:id", handler.GetRelationshipTypeByID)
			relationshipTypes.PUT("/:id", handler.UpdateRelationshipType)
			relationshipTypes.DELETE("/:id", handler.DeleteRelationshipType)
		}
		// EmploymentStatus
		employmentStatuses := settings.Group("/employment-statuses")
		{
			employmentStatuses.POST("", handler.CreateEmploymentStatus)
			employmentStatuses.GET("", handler.ListEmploymentStatuses)
			employmentStatuses.GET("/:id", handler.GetEmploymentStatusByID)
			employmentStatuses.PUT("/:id", handler.UpdateEmploymentStatus)
			employmentStatuses.DELETE("/:id", handler.DeleteEmploymentStatus)
		}
		// Bank
		banks := settings.Group("/banks")
		{
			banks.POST("", handler.CreateBank)
			banks.GET("", handler.ListBanks)
			banks.GET("/:id", handler.GetBankByID)
			banks.PUT("/:id", handler.UpdateBank)
			banks.DELETE("/:id", handler.DeleteBank)
		}
		// Nationality
		nationalities := settings.Group("/nationalities")
		{
			nationalities.POST("", handler.CreateNationality)
			nationalities.GET("", handler.ListNationalities)
			nationalities.GET("/:id", handler.GetNationalityByID)
			nationalities.PUT("/:id", handler.UpdateNationality)
			nationalities.DELETE("/:id", handler.DeleteNationality)
		}
		// JobFamily
		jobFamilies := settings.Group("/job-families")
		{
			jobFamilies.POST("", handler.CreateJobFamily)
			jobFamilies.GET("", handler.ListJobFamilies)
			jobFamilies.GET("/:id", handler.GetJobFamilyByID)
			jobFamilies.PUT("/:id", handler.UpdateJobFamily)
			jobFamilies.DELETE("/:id", handler.DeleteJobFamily)
		}
		// SalaryGrade
		salaryGrades := settings.Group("/salary-grades")
		{
			salaryGrades.POST("", handler.CreateSalaryGrade)
			salaryGrades.GET("", handler.ListSalaryGrades)
			salaryGrades.GET("/:id", handler.GetSalaryGradeByID)
			salaryGrades.PUT("/:id", handler.UpdateSalaryGrade)
			salaryGrades.DELETE("/:id", handler.DeleteSalaryGrade)
		}
		// TER (Tarif Efektif Rata-rata)
		ters := settings.Group("/ters")
		{
			ters.POST("", handler.CreateTER)
			ters.GET("", handler.ListTERs)
			ters.GET("/:id", handler.GetTERByID)
			ters.PUT("/:id", handler.UpdateTER)
			ters.DELETE("/:id", handler.DeleteTER)
		}
		// PTKP (Penghasilan Tidak Kena Pajak)
		ptkps := settings.Group("/ptkps")
		{
			ptkps.POST("", handler.CreatePTKP)
			ptkps.GET("", handler.ListPTKPs)
			ptkps.GET("/:id", handler.GetPTKPByID)
			ptkps.PUT("/:id", handler.UpdatePTKP)
			ptkps.DELETE("/:id", handler.DeletePTKP)
		}
	}
}
