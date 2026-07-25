package workforceintelligence

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Workforce Intelligence & Strategic Planning endpoints.
func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	wi := r.Group("/workforce-intelligence")
	{
		// Workforce Planning — Headcount
		wi.GET("/planning/headcounts", handler.ListHeadcountPlans)
		wi.POST("/planning/headcounts", handler.CreateHeadcountPlan)
		wi.GET("/planning/headcounts/:id", handler.GetHeadcountPlanByID)
		wi.PUT("/planning/headcounts/:id", handler.UpdateHeadcountPlan)
		wi.DELETE("/planning/headcounts/:id", handler.DeleteHeadcountPlan)

		// Workforce Planning — Forecast
		wi.GET("/planning/forecasts", handler.ListForecasts)
		wi.POST("/planning/forecasts", handler.CreateForecast)
		wi.GET("/planning/forecasts/:id", handler.GetForecastByID)
		wi.PUT("/planning/forecasts/:id", handler.UpdateForecast)
		wi.DELETE("/planning/forecasts/:id", handler.DeleteForecast)

		// Gap Analysis & Projections
		wi.GET("/planning/gap-analysis", handler.GetGapAnalysis)
		wi.GET("/planning/projections", handler.GetProjections)

		// KPIs
		wi.GET("/kpi", handler.ListKPIs)
		wi.GET("/kpi/summary", handler.GetKPISummary)

		// Analytics Dashboards
		wi.GET("/analytics/headcount", handler.GetHeadcountAnalytics)
		wi.GET("/analytics/attendance", handler.GetAttendanceAnalytics)
		wi.GET("/analytics/leave", handler.GetLeaveAnalytics)
		wi.GET("/analytics/overtime", handler.GetOvertimeAnalytics)
		wi.GET("/analytics/payroll", handler.GetPayrollAnalytics)
		wi.GET("/analytics/performance", handler.GetPerformanceAnalytics)
		wi.GET("/analytics/learning", handler.GetLearningAnalytics)
		wi.GET("/analytics/recruitment", handler.GetRecruitmentAnalytics)
		wi.GET("/analytics/movement", handler.GetMovementAnalytics)

		// Capacity
		wi.GET("/capacity/dashboard", handler.GetCapacityDashboard)
		wi.GET("/capacity/utilization", handler.GetUtilization)
		wi.GET("/capacity/forecast", handler.GetCapacityForecast)
		wi.GET("/capacity/bottlenecks", handler.GetBottlenecks)

		// Cost
		wi.GET("/cost/summary", handler.GetCostSummary)
		wi.GET("/cost/payroll", handler.GetPayrollCostBreakdown)
		wi.GET("/cost/per-employee", handler.GetCostPerEmployee)
		wi.GET("/cost/per-department", handler.GetCostByDepartment)
		wi.GET("/cost/budget-vs-actual", handler.GetBudgetVsActual)

		// Risk
		wi.GET("/risk/dashboard", handler.GetRiskDashboard)
		wi.GET("/risk/indicators", handler.ListRiskIndicators)
		wi.GET("/risk/indicators/:id", handler.GetRiskIndicatorByID)
		wi.PUT("/risk/indicators/:id", handler.UpdateRiskIndicator)
		wi.GET("/risk/high-turnover", handler.GetRiskHighTurnover)
		wi.GET("/risk/retirement", handler.GetRiskRetirement)
		wi.GET("/risk/contract-expiry", handler.GetRiskContractExpiry)
		wi.GET("/risk/high-absenteeism", handler.GetRiskHighAbsenteeism)

		// Executive Dashboard
		wi.GET("/executive/summary", handler.GetExecutiveSummary)
		wi.GET("/executive/growth", handler.GetExecutiveGrowth)
		wi.GET("/executive/cost-trend", handler.GetExecutiveCostTrend)
		wi.GET("/executive/attrition-trend", handler.GetExecutiveAttritionTrend)
		wi.GET("/executive/capacity", handler.GetExecutiveCapacity)
		wi.GET("/executive/hiring-progress", handler.GetHiringProgress)
		wi.GET("/executive/risk-overview", handler.GetExecutiveRiskOverview)
		wi.GET("/executive/health-score", handler.GetExecutiveHealthScore)

		// Scenario Planning
		wi.GET("/scenarios", handler.ListScenarios)
		wi.POST("/scenarios", handler.CreateScenario)
		wi.GET("/scenarios/:id", handler.GetScenarioByID)
		wi.PUT("/scenarios/:id", handler.UpdateScenario)
		wi.DELETE("/scenarios/:id", handler.DeleteScenario)
		wi.POST("/scenarios/:id/run", handler.RunScenario)
		wi.POST("/scenarios/:id/clone", handler.CloneScenario)

		// Organization Health
		wi.GET("/health/dashboard", handler.GetHealthDashboard)
		wi.GET("/health/scores", handler.ListHealthScores)
		wi.GET("/health/scores/:id", handler.GetHealthScoreByID)
		wi.GET("/health/span-of-control", handler.GetSpanOfControl)
		wi.GET("/health/succession", handler.GetSuccessionReadiness)

		// People Analytics
		wi.GET("/people-analytics/training-vs-performance", handler.GetTrainingVsPerformance)
		wi.GET("/people-analytics/overtime-vs-productivity", handler.GetOvertimeVsProductivity)
		wi.GET("/people-analytics/attendance-vs-performance", handler.GetAttendanceVsPerformance)
		wi.GET("/people-analytics/compensation-vs-turnover", handler.GetCompensationVsTurnover)
		wi.GET("/people-analytics/source-vs-retention", handler.GetSourceVsRetention)
		wi.GET("/people-analytics/career-progression", handler.GetCareerProgression)
		wi.GET("/people-analytics/learning-effectiveness", handler.GetLearningEffectiveness)

		// KPI by Code
		wi.GET("/kpi/:code", handler.GetKPIByCode)
	}
}
