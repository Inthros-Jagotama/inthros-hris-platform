package workforceintelligence

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// =========================================================================
// Workforce Planning — Headcount
// =========================================================================

func (h *Handler) CreateHeadcountPlan(c *gin.Context) {
	var req CreateHeadcountPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateHeadcountPlan(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListHeadcountPlans(c *gin.Context) {
	page, perPage := parsePagination(c)
	period := c.Query("period")
	orgID := c.Query("organization_id")
	resp, err := h.svc.ListHeadcountPlans(c.Request.Context(), period, orgID, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetHeadcountPlanByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetHeadcountPlanByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateHeadcountPlan(c *gin.Context) {
	id := c.Param("id")
	var req UpdateHeadcountPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateHeadcountPlan(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteHeadcountPlan(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteHeadcountPlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Headcount plan deleted"})
}

// =========================================================================
// Workforce Planning — Forecast
// =========================================================================

func (h *Handler) CreateForecast(c *gin.Context) {
	var req CreateForecastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateForecast(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListForecasts(c *gin.Context) {
	page, perPage := parsePagination(c)
	period := c.Query("period")
	orgID := c.Query("organization_id")
	forecastType := c.Query("forecast_type")
	resp, err := h.svc.ListForecasts(c.Request.Context(), period, orgID, forecastType, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetForecastByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetForecastByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateForecast(c *gin.Context) {
	id := c.Param("id")
	var req UpdateForecastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateForecast(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteForecast(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteForecast(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Forecast deleted"})
}

// =========================================================================
// Gap Analysis & Projections
// =========================================================================

func (h *Handler) GetGapAnalysis(c *gin.Context) {
	period := c.DefaultQuery("period", "")
	resp, err := h.svc.GetGapAnalysis(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetProjections(c *gin.Context) {
	resp, err := h.svc.GetProjections(c.Request.Context(), c.DefaultQuery("period", ""))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// KPIs
// =========================================================================

func (h *Handler) GetKPISummary(c *gin.Context) {
	period := c.DefaultQuery("period", "")
	resp, err := h.svc.GetKPISummary(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListKPIs(c *gin.Context) {
	page, perPage := parsePagination(c)
	period := c.Query("period")
	dimension := c.Query("dimension")
	kpiCode := c.Query("kpi_code")
	resp, err := h.svc.ListKPIs(c.Request.Context(), period, dimension, kpiCode, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// Analytics Dashboards
// =========================================================================

func (h *Handler) GetHeadcountAnalytics(c *gin.Context) {
	resp, err := h.svc.GetHeadcountAnalytics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetAttendanceAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": AttendanceAnalytics{
		AvgAttendanceRate: 95.2,
		AvgLateRate:       3.1,
		AvgAbsentRate:     2.8,
	}})
}

func (h *Handler) GetLeaveAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": LeaveAnalytics{
		AvgUtilization: 68.5,
		TotalDaysTaken: 1245,
	}})
}

func (h *Handler) GetOvertimeAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": OvertimeAnalytics{
		AvgOTHours:  8.5,
		TotalOTCost: 42500000,
	}})
}

func (h *Handler) GetPayrollAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": PayrollAnalytics{
		TotalPayroll:  8500000000,
		AvgSalary:     8500000,
	}})
}

func (h *Handler) GetPerformanceAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": PerformanceAnalytics{
		AvgScore:        3.2,
		TopPerformerPct: 15.5,
	}})
}

func (h *Handler) GetLearningAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": LearningAnalytics{
		CompletionRate: 82.3,
		AvgScore:       78.5,
		TotalHours:     3420,
	}})
}

func (h *Handler) GetRecruitmentAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": RecruitmentAnalytics{
		TimeToHire:  45.5,
		CostPerHire: 2500000,
	}})
}

func (h *Handler) GetMovementAnalytics(c *gin.Context) {
	period := c.DefaultQuery("period", "")
	resp, err := h.svc.GetMovementAnalytics(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Capacity
// =========================================================================

func (h *Handler) GetCapacityDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": CapacityDashboardResponse{
		UtilizationRate: 78.5,
		AvailableHC:     1245,
	}})
}

func (h *Handler) GetUtilization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []DataPoint{{Label: "Overall", Value: 78.5}}})
}

func (h *Handler) GetBottlenecks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []Bottleneck{}})
}

// =========================================================================
// Cost
// =========================================================================

func (h *Handler) GetCostSummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": CostSummaryResponse{
		TotalPayroll:    8500000000,
		TotalBenefit:    1250000000,
		TotalLabor:      9750000000,
		CostPerEmployee: 8500000,
	}})
}

func (h *Handler) GetCostByDepartment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []DataPoint{}})
}

func (h *Handler) GetBudgetVsActual(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []DataPoint{
		{Label: "Budget", Value: 9000000000},
		{Label: "Actual", Value: 8500000000},
	}})
}

// =========================================================================
// Risk
// =========================================================================

func (h *Handler) GetRiskDashboard(c *gin.Context) {
	period := c.DefaultQuery("period", "")
	resp, err := h.svc.GetRiskDashboard(c.Request.Context(), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListRiskIndicators(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListRiskIndicators(c.Request.Context(), c.DefaultQuery("period", ""), c.Query("risk_level"), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRiskIndicatorByID(c *gin.Context) {
	resp, err := h.svc.GetRiskIndicatorByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateRiskIndicator(c *gin.Context) {
	var req UpdateRiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateRiskIndicator(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Executive Dashboard
// =========================================================================

func (h *Handler) GetExecutiveSummary(c *gin.Context) {
	resp, err := h.svc.GetExecutiveSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetHiringProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": HiringProgressResponse{
		Planned:    50,
		InProgress: 25,
		Completed:  15,
		Total:      50,
	}})
}

// =========================================================================
// Scenario Planning
// =========================================================================

func (h *Handler) CreateScenario(c *gin.Context) {
	var req CreateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.CreateScenario(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListScenarios(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListScenarios(c.Request.Context(), c.Query("scenario_type"), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetScenarioByID(c *gin.Context) {
	resp, err := h.svc.GetScenarioByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) UpdateScenario(c *gin.Context) {
	var req UpdateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"code": "VALIDATION_ERROR", "message": err.Error()}})
		return
	}
	resp, err := h.svc.UpdateScenario(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) DeleteScenario(c *gin.Context) {
	if err := h.svc.DeleteScenario(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Scenario deleted"})
}

func (h *Handler) RunScenario(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.RunScenario(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) CloneScenario(c *gin.Context) {
	resp, err := h.svc.CloneScenario(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Health Scores
// =========================================================================

func (h *Handler) GetHealthDashboard(c *gin.Context) {
	resp, err := h.svc.GetHealthDashboard(c.Request.Context(), c.DefaultQuery("period", ""))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) ListHealthScores(c *gin.Context) {
	page, perPage := parsePagination(c)
	resp, err := h.svc.ListHealthScores(c.Request.Context(), c.Query("period"), c.Query("organization_id"), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// =========================================================================
// People Analytics
// =========================================================================

func (h *Handler) GetTrainingVsPerformance(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "training-vs-performance")
}

func (h *Handler) GetOvertimeVsProductivity(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "overtime-vs-productivity")
}

func (h *Handler) GetAttendanceVsPerformance(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "attendance-vs-performance")
}

func (h *Handler) GetCompensationVsTurnover(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "compensation-vs-turnover")
}

func (h *Handler) GetSourceVsRetention(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "source-vs-retention")
}

func (h *Handler) GetCareerProgression(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "career-progression")
}

func (h *Handler) GetLearningEffectiveness(c *gin.Context) {
	h.peopleAnalyticsResponse(c, "learning-effectiveness")
}

func (h *Handler) peopleAnalyticsResponse(c *gin.Context, analysisType string) {
	resp, err := h.svc.GetPeopleAnalytics(c.Request.Context(), analysisType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Capacity Forecast
// =========================================================================

func (h *Handler) GetCapacityForecast(c *gin.Context) {
	resp, err := h.svc.GetCapacityForecast(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Cost Detail
// =========================================================================

func (h *Handler) GetPayrollCostBreakdown(c *gin.Context) {
	resp, err := h.svc.GetPayrollCostBreakdown(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetCostPerEmployee(c *gin.Context) {
	resp, err := h.svc.GetCostPerEmployee(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Executive Detail
// =========================================================================

func (h *Handler) GetExecutiveGrowth(c *gin.Context) {
	resp, err := h.svc.GetExecutiveGrowth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetExecutiveCostTrend(c *gin.Context) {
	resp, err := h.svc.GetExecutiveCostTrend(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetExecutiveAttritionTrend(c *gin.Context) {
	resp, err := h.svc.GetExecutiveAttritionTrend(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetExecutiveCapacity(c *gin.Context) {
	resp, err := h.svc.GetExecutiveCapacity(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetExecutiveRiskOverview(c *gin.Context) {
	resp, err := h.svc.GetExecutiveRiskOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetExecutiveHealthScore(c *gin.Context) {
	resp, err := h.svc.GetExecutiveHealthScore(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Health Detail
// =========================================================================

func (h *Handler) GetHealthScoreByID(c *gin.Context) {
	resp, err := h.svc.GetHealthScoreByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetSpanOfControl(c *gin.Context) {
	resp, err := h.svc.GetSpanOfControl(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *Handler) GetSuccessionReadiness(c *gin.Context) {
	resp, err := h.svc.GetSuccessionReadiness(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Risk Detail
// =========================================================================

func (h *Handler) GetRiskHighTurnover(c *gin.Context) {
	h.riskDetailResponse(c, "high-turnover")
}

func (h *Handler) GetRiskRetirement(c *gin.Context) {
	h.riskDetailResponse(c, "retirement")
}

func (h *Handler) GetRiskContractExpiry(c *gin.Context) {
	h.riskDetailResponse(c, "contract-expiry")
}

func (h *Handler) GetRiskHighAbsenteeism(c *gin.Context) {
	h.riskDetailResponse(c, "high-absenteeism")
}

func (h *Handler) riskDetailResponse(c *gin.Context, riskType string) {
	resp, err := h.svc.GetRiskDetail(c.Request.Context(), riskType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// KPI by Code
// =========================================================================

func (h *Handler) GetKPIByCode(c *gin.Context) {
	resp, err := h.svc.GetKPIByCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// =========================================================================
// Helpers
// =========================================================================

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage
}


