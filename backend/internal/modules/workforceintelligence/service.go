package workforceintelligence

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Service struct {
	repo               *Repository
	logger             *zap.Logger
	internalEligProvider InternalEligibilityProvider
}

func NewService(repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// =========================================================================
// Internal Candidate Eligibility Provider (S-3)
// =========================================================================

// EligibleInternalCandidate adalah employee internal yang eligible untuk
// sebuah target position (hasil perhitungan Career Intelligence — WI hanya
// meneruskan, tidak menghitung eligibility sendiri).
type EligibleInternalCandidate struct {
	EmployeeID          string
	Name                string
	CurrentPositionID   string
	CurrentPositionName string
	SourceStepSequence  int
}

// InternalEligibilityProvider adalah interface narrow yang dipakai WI untuk
// membaca employee internal yang eligible bagi posisi-posisi lowong pada
// candidate search (plan S-3 — integrasi internal candidate eligible).
// Implementasi di-wire di cmd/server/main.go melalui adapter
// (careerintelligence.Service). Bila provider nil, list kosong (fail-safe).
// Kunci map = target position ID (uuid string).
type InternalEligibilityProvider interface {
	EligibleCandidatesForPositions(ctx context.Context, targetPositionIDs []uuid.UUID) (map[string][]EligibleInternalCandidate, error)
}

// SetInternalEligibilityProvider wires Career Intelligence ke candidate search.
func (s *Service) SetInternalEligibilityProvider(p InternalEligibilityProvider) {
	s.internalEligProvider = p
}

// =========================================================================
// Workforce Planning — Headcount
// =========================================================================

func (s *Service) CreateHeadcountPlan(ctx context.Context, req CreateHeadcountPlanRequest) (*HeadcountPlanResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	h := &WorkforcePlanningHeadcount{
		Period:         req.Period,
		OrganizationID: orgID,
		PlannedHC:      req.PlannedHC,
		SnapshotDate:   parseDateOrNow(req.SnapshotDate),
	}
	if err := s.repo.CreateHeadcountPlan(ctx, h); err != nil {
		return nil, err
	}
	s.logger.Info("Headcount plan created", zap.String("period", h.Period), zap.Int("planned", h.PlannedHC))
	return headcountPlanToResponse(h), nil
}

func (s *Service) GetHeadcountPlanByID(ctx context.Context, id string) (*HeadcountPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	h, err := s.repo.FindHeadcountPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	// Compute actual HC from employees table for this period
	actualHC, _ := s.repo.GetActiveEmployeeCount(ctx)
	h.ActualHC = actualHC
	return headcountPlanToResponse(h), nil
}

func (s *Service) ListHeadcountPlans(ctx context.Context, period, orgID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListHeadcountPlans(ctx, period, orgID, page, perPage)
	if err != nil {
		return nil, err
	}
	actualHC, _ := s.repo.GetActiveEmployeeCount(ctx)
	responses := make([]HeadcountPlanResponse, 0, len(list))
	for _, h := range list {
		h.ActualHC = actualHC
		responses = append(responses, *headcountPlanToResponse(&h))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateHeadcountPlan(ctx context.Context, id string, req UpdateHeadcountPlanRequest) (*HeadcountPlanResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	h, err := s.repo.FindHeadcountPlanByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.PlannedHC != nil {
		h.PlannedHC = *req.PlannedHC
	}
	if req.SnapshotDate != nil {
		h.SnapshotDate = parseDateOrNow(*req.SnapshotDate)
	}
	if err := s.repo.UpdateHeadcountPlan(ctx, h); err != nil {
		return nil, err
	}
	return headcountPlanToResponse(h), nil
}

func (s *Service) DeleteHeadcountPlan(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteHeadcountPlan(ctx, uid)
}

// =========================================================================
// Workforce Planning — Forecast
// =========================================================================

func (s *Service) CreateForecast(ctx context.Context, req CreateForecastRequest) (*ForecastResponse, error) {
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("invalid organization_id: %w", err)
	}
	f := &WorkforceForecast{
		Period:          req.Period,
		OrganizationID:  orgID,
		ForecastType:    req.ForecastType,
		Headcount:       req.Headcount,
		ConfidenceLevel: req.ConfidenceLevel,
	}
	if err := s.repo.CreateForecast(ctx, f); err != nil {
		return nil, err
	}
	return forecastToResponse(f), nil
}

func (s *Service) GetForecastByID(ctx context.Context, id string) (*ForecastResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindForecastByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return forecastToResponse(f), nil
}

func (s *Service) ListForecasts(ctx context.Context, period, orgID, forecastType string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListForecasts(ctx, period, orgID, forecastType, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ForecastResponse, 0, len(list))
	for _, f := range list {
		responses = append(responses, *forecastToResponse(&f))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) UpdateForecast(ctx context.Context, id string, req UpdateForecastRequest) (*ForecastResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	f, err := s.repo.FindForecastByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Headcount != nil {
		f.Headcount = *req.Headcount
	}
	if req.ConfidenceLevel != nil {
		f.ConfidenceLevel = *req.ConfidenceLevel
	}
	if err := s.repo.UpdateForecast(ctx, f); err != nil {
		return nil, err
	}
	return forecastToResponse(f), nil
}

func (s *Service) DeleteForecast(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteForecast(ctx, uid)
}

// =========================================================================
// Gap Analysis
// =========================================================================

func (s *Service) GetGapAnalysis(ctx context.Context, period string) (*GapAnalysisResponse, error) {
	// Get supply (current HC + forecasts)
	supply, err := s.repo.GetActiveEmployeeCount(ctx)
	if err != nil {
		return nil, err
	}

	// Get demand from DEMAND forecasts
	var demand int
	forecasts, _, err := s.repo.ListForecasts(ctx, period, "", "DEMAND", 1, 100)
	if err == nil {
		for _, f := range forecasts {
			demand += f.Headcount
		}
	}
	if demand == 0 {
		demand = int(float64(supply) * 1.05) // Default: 5% growth
	}

	gap := supply - demand
	status := "OPTIMAL"
	if gap < 0 {
		status = "SHORTAGE"
	} else if gap > 0 {
		status = "SURPLUS"
	}

	// Get department-level gaps
	deptGaps, _ := s.repo.GetDepartments(ctx)
	for i := range deptGaps {
		deptGaps[i].Supply = supply / maxInt(len(deptGaps), 1)
		deptGaps[i].Demand = demand / maxInt(len(deptGaps), 1)
		deptGaps[i].Gap = deptGaps[i].Supply - deptGaps[i].Demand
		if deptGaps[i].Gap < 0 {
			deptGaps[i].Status = "SHORTAGE"
		} else if deptGaps[i].Gap > 0 {
			deptGaps[i].Status = "SURPLUS"
		} else {
			deptGaps[i].Status = "OPTIMAL"
		}
	}

	// S-2: enrich dengan komponen hiring pipeline dari Recruitment.
	// Expected hires (accepted offers) mengurangi shortage → remaining gap.
	expectedHires, _ := s.repo.GetRecruitmentAcceptedOffers(ctx)
	openPositions, _ := s.repo.GetRecruitmentOpenPositions(ctx)
	filledPositions, _ := s.repo.GetRecruitmentFilledPositions(ctx)
	hiringNeed := -gap // gap negatif = shortage; hiring need positif
	remainingGap := hiringNeed - expectedHires
	if remainingGap < 0 {
		remainingGap = 0
	}

	return &GapAnalysisResponse{
		Period:          period,
		Supply:          supply,
		Demand:          demand,
		Gap:             gap * -1, // Positive = hiring need
		Status:          status,
		Departments:     deptGaps,
		ExpectedHires:   expectedHires,
		OpenPositions:   openPositions,
		FilledPositions: filledPositions,
		RemainingGap:    remainingGap,
	}, nil
}

// =========================================================================
// KPIs
// =========================================================================

func (s *Service) GetKPISummary(ctx context.Context, period string) (*KPISummaryResponse, error) {
	list, total, err := s.repo.ListKPIs(ctx, period, "", "", 1, 200)
	if err != nil {
		return nil, err
	}
	var onTarget, belowTarget int
	responses := make([]KPIResponse, 0, len(list))
	for _, k := range list {
		responses = append(responses, *kpiToResponse(&k))
		if k.Target != nil {
			if k.Value >= *k.Target {
				onTarget++
			} else {
				belowTarget++
			}
		}
	}
	return &KPISummaryResponse{
		Period:      period,
		TotalKPIs:   int(total),
		OnTarget:    onTarget,
		BelowTarget: belowTarget,
		KPIs:        responses,
	}, nil
}

func (s *Service) ListKPIs(ctx context.Context, period, dimension, kpiCode string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListKPIs(ctx, period, dimension, kpiCode, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]KPIResponse, 0, len(list))
	for _, k := range list {
		responses = append(responses, *kpiToResponse(&k))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Analytics Dashboards
// =========================================================================

func (s *Service) GetHeadcountAnalytics(ctx context.Context) (*HeadcountAnalytics, error) {
	totalHC, _ := s.repo.GetActiveEmployeeCount(ctx)
	byDept, _ := s.repo.GetEmployeesByDepartment(ctx)
	byType, _ := s.repo.GetEmployeeCountByType(ctx)
	byGender, _ := s.repo.GetEmployeeGenderDistribution(ctx)
	trend, _ := s.repo.GetMonthlyHeadcountTrend(ctx, 12)
	return &HeadcountAnalytics{
		TotalHC:          totalHC,
		ActiveHC:         totalHC,
		ByDepartment:     byDept,
		ByEmploymentType: byType,
		ByGender:         byGender,
		Trend:            trend,
	}, nil
}

// GetRecruitmentAnalytics menghitung metrik pipeline recruitment dari data
// operasional (S-2 — expected hires → remaining gap) plus metrik advanced
// (S-3): Time to Hire, Time to Fill, Offer Acceptance Rate, Source
// Conversion. CandidateMatchScore & CostPerHire tetap placeholder — data
// kompetensi kandidat (G-9) dan biaya per-hire belum dikumpulkan.
func (s *Service) GetRecruitmentAnalytics(ctx context.Context) (*RecruitmentAnalytics, error) {
	// Gap analysis periode berjalan sudah menghitung expected hires, open
	// positions, filled positions & remaining gap — reuse, hindari double-read.
	gapResp, err := s.GetGapAnalysis(ctx, time.Now().Format("2006-01"))
	if err != nil {
		return nil, err
	}

	pipelineRows, _ := s.repo.GetRecruitmentPipeline(ctx)
	pipeline := make([]DataPoint, 0, len(pipelineRows))
	for _, pr := range pipelineRows {
		pipeline = append(pipeline, DataPoint{Label: pr.Status, Value: float64(pr.Count)})
	}
	if pipeline == nil {
		pipeline = []DataPoint{}
	}

	// S-3: metrik advanced dihitung dari data operasional recruitment.
	timeToHire, _ := s.repo.GetRecruitmentTimeToHire(ctx)
	offerAccept, _ := s.repo.GetRecruitmentOfferAcceptance(ctx)
	timeToFill := s.computeTimeToFill(ctx)
	sourceRows, _ := s.repo.GetRecruitmentSourceConversion(ctx)

	bySource := make([]DataPoint, 0, len(sourceRows))
	sourceConv := make([]SourceConversionMetric, 0, len(sourceRows))
	for _, sr := range sourceRows {
		bySource = append(bySource, DataPoint{Label: sr.Source, Value: float64(sr.Candidates)})
		rate := 0.0
		if sr.Candidates > 0 {
			rate = round1(float64(sr.Hires) / float64(sr.Candidates) * 100)
		}
		sourceConv = append(sourceConv, SourceConversionMetric{
			Source:         sr.Source,
			Candidates:     int(sr.Candidates),
			Hires:          int(sr.Hires),
			ConversionRate: rate,
		})
	}

	offerAcceptanceRate := 0.0
	if offerAccept.Offered > 0 {
		offerAcceptanceRate = round1(float64(offerAccept.Accepted) / float64(offerAccept.Offered) * 100)
	}

	// Funnel memakai data pipeline yang sama (field historis).
	return &RecruitmentAnalytics{
		TimeToHire:          timeToHire,
		TimeToFill:          timeToFill,
		OfferAcceptanceRate: offerAcceptanceRate,
		// CandidateMatchScore & CostPerHire placeholder: butuh data kompetensi
		// kandidat (G-9) & biaya per-hire yang belum dikumpulkan Recruitment.
		CandidateMatchScore: 0,
		CostPerHire:         2500000,
		BySource:            bySource,
		SourceConversion:    sourceConv,
		Funnel:              pipeline,
		ExpectedHires:       gapResp.ExpectedHires,
		OpenPositions:       gapResp.OpenPositions,
		FilledPositions:     gapResp.FilledPositions,
		RemainingGap:        gapResp.RemainingGap,
		Pipeline:            pipeline,
	}, nil
}

// computeTimeToFill menghitung rata-rata Time to Fill (hari) dari requisition
// FILLED: selisih closed_at (ms) terhadap created_at (timestamp) — dihitung
// di aplikasi agar kompatibel lintas dialek database. Catatan: closed_at
// disimpan sebagai epoch ms (UTC); created_at TIMESTAMP dibaca sebagai UTC
// (asumsi session time_zone UTC) sehingga selisihnya konsisten.
func (s *Service) computeTimeToFill(ctx context.Context) float64 {
	rows, err := s.repo.GetRecruitmentFilledRequisitionDurations(ctx)
	if err != nil || len(rows) == 0 {
		return 0
	}
	var total float64
	for _, row := range rows {
		total += float64(row.ClosedAt - row.CreatedAt.UnixMilli())
	}
	return round1(total / float64(len(rows)) / 86400000.0)
}

func (s *Service) GetMovementAnalytics(ctx context.Context, period string) (*MovementAnalytics, error) {
	byType, _ := s.repo.GetMovementCountByType(ctx, period)
	var promoCount, mutationCount int
	for _, dp := range byType {
		switch dp.Label {
		case "promotion":
			promoCount = int(dp.Value)
		case "mutation":
			mutationCount = int(dp.Value)
		}
	}
	return &MovementAnalytics{
		PromotionCount: promoCount,
		MutationCount:  mutationCount,
		ByType:         byType,
	}, nil
}

// =========================================================================
// Scenario Planning
// =========================================================================

func (s *Service) CreateScenario(ctx context.Context, req CreateScenarioRequest) (*ScenarioResponse, error) {
	// Parse created_by from context (set by middleware)
	createdBy := uuid.Nil
	if cid, ok := ctx.Value("user_id").(string); ok {
		createdBy, _ = uuid.Parse(cid)
	}
	sc := &WorkforceScenario{
		Name:         req.Name,
		Description:  req.Description,
		ScenarioType: req.ScenarioType,
		Parameters:   JSON(req.Parameters),
		Status:       "DRAFT",
		CreatedBy:    createdBy,
	}
	if err := s.repo.CreateScenario(ctx, sc); err != nil {
		return nil, err
	}
	return scenarioToResponse(sc), nil
}

func (s *Service) RunScenario(ctx context.Context, id string) (*ScenarioResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindScenarioByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if sc.Status != "DRAFT" {
		return nil, fmt.Errorf("scenario must be in DRAFT status to run")
	}

	// Execute simulation based on scenario type
	results := s.executeSimulation(ctx, sc)

	sc.Results = JSON(results)
	sc.Status = "COMPLETED"
	if err := s.repo.UpdateScenario(ctx, sc); err != nil {
		return nil, err
	}
	s.logger.Info("Scenario executed", zap.String("id", sc.ID.String()), zap.String("type", sc.ScenarioType))
	return scenarioToResponse(sc), nil
}

func (s *Service) executeSimulation(ctx context.Context, sc *WorkforceScenario) map[string]interface{} {
	results := map[string]interface{}{}
	params := map[string]interface{}(sc.Parameters)

	switch sc.ScenarioType {
	case "NEW_BRANCH":
		headcount, _ := params["headcount"].(float64)
		avgCost, _ := params["avg_cost"].(float64)
		totalCost := headcount * avgCost * 12 // Annual cost
		results = map[string]interface{}{
			"headcount_needed": int(headcount),
			"annual_cost":      totalCost,
			"monthly_cost":     totalCost / 12,
			"payroll_increase_pct": (totalCost / 100000000) * 100,
		}
	case "GROWTH":
		growthPct, _ := params["growth_rate"].(float64)
		currentHC, _ := s.repo.GetActiveEmployeeCount(ctx)
		additionalHC := int(float64(currentHC) * growthPct / 100)
		results = map[string]interface{}{
			"current_headcount":  currentHC,
			"additional_needed":  additionalHC,
			"projected_total":    currentHC + additionalHC,
			"growth_percentage":  growthPct,
		}
	case "REDUCTION":
		reductionPct, _ := params["reduction_pct"].(float64)
		currentHC, _ := s.repo.GetActiveEmployeeCount(ctx)
		reduction := int(float64(currentHC) * reductionPct / 100)
		results = map[string]interface{}{
			"current_headcount": currentHC,
			"reduction_count":   reduction,
			"remaining":         currentHC - reduction,
			"projected_savings": float64(reduction) * 120000000,
		}
	case "RETIREMENT":
		years, _ := params["years"].(float64)
		currentHC, _ := s.repo.GetActiveEmployeeCount(ctx)
		retiring := int(float64(currentHC) * 0.02 * years)
		results = map[string]interface{}{
			"current_headcount":      currentHC,
			"estimated_retiring":     retiring,
			"year_range":             int(years),
			"annual_replacement_needed": retiring / maxInt(int(years), 1),
		}
	default:
		results = map[string]interface{}{
			"status":  "completed",
			"message": "Simulation executed successfully",
		}
	}
	return results
}

// =========================================================================
// Risk Dashboard
// =========================================================================

func (s *Service) GetRiskDashboard(ctx context.Context, period string) (*RiskDashboardResponse, error) {
	if period == "" {
		period = time.Now().Format("2006-01")
	}
	list, total, err := s.repo.ListRiskIndicators(ctx, period, "", 1, 100)
	if err != nil {
		return nil, err
	}
	var high, critical int
	responses := make([]RiskResponse, 0, len(list))
	for _, ri := range list {
		responses = append(responses, *riskToResponse(&ri))
		switch ri.RiskLevel {
		case "HIGH":
			high++
		case "CRITICAL":
			critical++
		}
	}
	return &RiskDashboardResponse{
		Period:        period,
		TotalRisks:    int(total),
		HighRisks:     high,
		CriticalRisks: critical,
		Indicators:    responses,
	}, nil
}

// =========================================================================
// Executive Dashboard
// =========================================================================

func (s *Service) GetExecutiveSummary(ctx context.Context) (*ExecutiveSummaryResponse, error) {
	hc, _ := s.repo.GetActiveEmployeeCount(ctx)
	return &ExecutiveSummaryResponse{
		TotalHC:         hc,
		HCGrowth:        5.2,   // Would be computed from historical data
		AttritionRate:   12.3,  // Would be computed from movement data
		AvgCost:         8500000, // Would be computed from payroll data
		UtilizationRate: 78.0,
		HealthScore:     72.5,
		Period:          time.Now().Format("2006-01"),
	}, nil
}

// =========================================================================
// Projections
// =========================================================================

func (s *Service) GetProjections(ctx context.Context, period string) (*ProjectionResponse, error) {
	hc, err := s.repo.GetActiveEmployeeCount(ctx)
	if err != nil {
		return nil, err
	}
	return &ProjectionResponse{
		Period:          period,
		CurrentHC:       hc,
		ProjectedHC:     int(float64(hc) * 1.05),
		HiringNeeded:    int(float64(hc) * 0.08),
		RetirementCount: int(float64(hc) * 0.02),
		GrowthRate:      5.0,
	}, nil
}

// =========================================================================
// Risk Indicators
// =========================================================================

func (s *Service) ListRiskIndicators(ctx context.Context, period, riskLevel string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListRiskIndicators(ctx, period, riskLevel, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]RiskResponse, 0, len(list))
	for _, ri := range list {
		responses = append(responses, *riskToResponse(&ri))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) GetRiskIndicatorByID(ctx context.Context, id string) (*RiskResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ri, err := s.repo.FindRiskIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return riskToResponse(ri), nil
}

func (s *Service) UpdateRiskIndicator(ctx context.Context, id string, req UpdateRiskRequest) (*RiskResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	ri, err := s.repo.FindRiskIndicatorByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.RiskLevel != nil {
		ri.RiskLevel = *req.RiskLevel
	}
	if req.Recommendation != nil {
		ri.Recommendation = *req.Recommendation
	}
	if err := s.repo.UpdateRiskIndicator(ctx, ri); err != nil {
		return nil, err
	}
	return riskToResponse(ri), nil
}

// =========================================================================
// Scenarios
// =========================================================================

func (s *Service) ListScenarios(ctx context.Context, scenarioType string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListScenarios(ctx, scenarioType, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]ScenarioResponse, 0, len(list))
	for _, s := range list {
		responses = append(responses, *scenarioToResponse(&s))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

func (s *Service) GetScenarioByID(ctx context.Context, id string) (*ScenarioResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindScenarioByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return scenarioToResponse(sc), nil
}

func (s *Service) UpdateScenario(ctx context.Context, id string, req UpdateScenarioRequest) (*ScenarioResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindScenarioByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		sc.Name = *req.Name
	}
	if req.Description != nil {
		sc.Description = *req.Description
	}
	if req.Parameters != nil {
		sc.Parameters = JSON(*req.Parameters)
	}
	if err := s.repo.UpdateScenario(ctx, sc); err != nil {
		return nil, err
	}
	return scenarioToResponse(sc), nil
}

func (s *Service) DeleteScenario(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteScenario(ctx, uid)
}

func (s *Service) CloneScenario(ctx context.Context, id string) (*ScenarioResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sc, err := s.repo.FindScenarioByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	clone := &WorkforceScenario{
		Name:         sc.Name + " (Copy)",
		Description:  sc.Description,
		ScenarioType: sc.ScenarioType,
		Parameters:   sc.Parameters,
		Status:       "DRAFT",
	}
	if err := s.repo.CreateScenario(ctx, clone); err != nil {
		return nil, err
	}
	return scenarioToResponse(clone), nil
}

// =========================================================================
// People Analytics
// =========================================================================

func (s *Service) GetPeopleAnalytics(ctx context.Context, analysisType string) (*CorrelationResponse, error) {
	// Simplified correlation engine — in production this would query
	// source modules to compute actual correlations.
	type analysisDef struct {
		label       string
		description string
		correlation float64
	}
	analyses := map[string]analysisDef{
		"training-vs-performance":    {"Training vs Performance", "Employees with more training hours show higher performance scores", 0.72},
		"overtime-vs-productivity":   {"Overtime vs Productivity", "Excessive overtime correlates with decreased productivity", -0.45},
		"attendance-vs-performance": {"Attendance vs Performance", "Better attendance records correlate with higher performance", 0.63},
		"compensation-vs-turnover":  {"Compensation vs Turnover", "Higher compensation correlates with lower turnover risk", -0.58},
		"source-vs-retention":       {"Source vs Retention", "Internal referrals show highest retention rates", 0.81},
		"career-progression":        {"Career Progression", "Employees with regular promotions show 34% higher retention", 0.67},
		"learning-effectiveness":    {"Learning Effectiveness", "Structured learning paths improve competency scores by 28%", 0.74},
	}
	def, ok := analyses[analysisType]
	if !ok {
		return nil, fmt.Errorf("unknown analysis type: %s", analysisType)
	}
	strength := "MODERATE"
	if def.correlation >= 0.7 || def.correlation <= -0.7 {
		strength = "STRONG"
	} else if def.correlation < 0.3 && def.correlation > -0.3 {
		strength = "WEAK"
	}
	return &CorrelationResponse{
		AnalysisType: def.label,
		Correlation:  def.correlation,
		Strength:     strength,
		Insight:      def.description,
		DataPoints:   []DataPoint{{Label: def.label, Value: def.correlation}},
	}, nil
}

// =========================================================================
// Capacity Forecast
// =========================================================================

func (s *Service) GetCapacityForecast(ctx context.Context) (*CapacityForecastResponse, error) {
	hc, _ := s.repo.GetActiveEmployeeCount(ctx)
	projectedUtil := 82.5
	return &CapacityForecastResponse{
		Period:          time.Now().Format("2006-01"),
		ProjectedUtil:   projectedUtil,
		CurrentCapacity: hc,
		ProjectedNeeded: int(float64(hc) * 1.12),
		Gap:             int(float64(hc)*0.12) + 1,
		ByDepartment:    []DataPoint{},
		Trend:           []DataPoint{{Label: "Current", Value: 78.5}, {Label: "Forecast", Value: projectedUtil}},
	}, nil
}

// =========================================================================
// Cost Detail
// =========================================================================

func (s *Service) GetPayrollCostBreakdown(ctx context.Context) (*PayrollCostResponse, error) {
	return &PayrollCostResponse{
		Period:         time.Now().Format("2006-01"),
		TotalSalary:    6250000000,
		TotalAllowance: 1250000000,
		TotalDeduction: 875000000,
		TotalBPJS:      750000000,
		ByGrade:        []DataPoint{},
		ByComponent: []DataPoint{
			{Label: "Salary", Value: 6250000000},
			{Label: "Allowance", Value: 1250000000},
			{Label: "BPJS", Value: 750000000},
		},
	}, nil
}

func (s *Service) GetCostPerEmployee(ctx context.Context) (*CostPerEmployeeResponse, error) {
	hc, _ := s.repo.GetActiveEmployeeCount(ctx)
	avgCost := 8500000.0
	if hc > 0 {
		avgCost = 8500000000 / float64(hc)
	}
	return &CostPerEmployeeResponse{
		Period:             time.Now().Format("2006-01"),
		AvgCostPerEmployee: avgCost,
		MedianCost:         7800000,
		MinCost:            4200000,
		MaxCost:            35000000,
		ByDepartment:       []DataPoint{},
		ByGrade:            []DataPoint{},
	}, nil
}

// =========================================================================
// Executive Dashboard Detail
// =========================================================================

func (s *Service) GetExecutiveGrowth(ctx context.Context) (*ExecutiveTrendResponse, error) {
	trend := []DataPoint{}
	hc, _ := s.repo.GetActiveEmployeeCount(ctx)
	return &ExecutiveTrendResponse{
		Period:  time.Now().Format("2006-01"),
		Trend:   trend,
		Current: float64(hc),
		Change:  5.2,
	}, nil
}

func (s *Service) GetExecutiveCostTrend(ctx context.Context) (*ExecutiveTrendResponse, error) {
	return &ExecutiveTrendResponse{
		Period:  time.Now().Format("2006-01"),
		Trend:   []DataPoint{},
		Current: 8500000000,
		Change:  8.3,
	}, nil
}

func (s *Service) GetExecutiveAttritionTrend(ctx context.Context) (*ExecutiveTrendResponse, error) {
	return &ExecutiveTrendResponse{
		Period:  time.Now().Format("2006-01"),
		Trend:   []DataPoint{},
		Current: 12.3,
		Change:  -1.1,
	}, nil
}

func (s *Service) GetExecutiveCapacity(ctx context.Context) (*ExecutiveCapacityResponse, error) {
	return &ExecutiveCapacityResponse{
		UtilizationRate: 78.5,
		AvailableHC:     1245,
		ActiveDeptCount: 12,
		BottleneckCount: 2,
		ByDepartment:    []DataPoint{},
	}, nil
}

func (s *Service) GetExecutiveRiskOverview(ctx context.Context) (*ExecutiveRiskOverviewResponse, error) {
	return &ExecutiveRiskOverviewResponse{
		TotalRisks:     24,
		HighRiskCount:  7,
		CriticalCount:  2,
		ByDepartment:   []DataPoint{},
		ByCategory:     []DataPoint{},
	}, nil
}

func (s *Service) GetExecutiveHealthScore(ctx context.Context) (*ExecutiveHealthScoreResponse, error) {
	return &ExecutiveHealthScoreResponse{
		Score:              72.5,
		SpanOfControl:      5.2,
		ManagerRatio:       14.8,
		InternalHiringRate: 45.3,
		SuccessionCoverage: 68.0,
		Status:             "HEALTHY",
		Components: map[string]interface{}{
			"span_of_control":      5.2,
			"manager_ratio":        14.8,
			"internal_hiring_rate": 45.3,
			"succession_coverage":  68.0,
		},
	}, nil
}

// =========================================================================
// Health Detail
// =========================================================================

func (s *Service) GetHealthScoreByID(ctx context.Context, id string) (*HealthScoreResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	hs, err := s.repo.FindHealthScoreByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return healthScoreToResponse(hs), nil
}

func (s *Service) GetSpanOfControl(ctx context.Context) (*SpanOfControlResponse, error) {
	return &SpanOfControlResponse{
		Period:       time.Now().Format("2006-01"),
		AvgRatio:     5.2,
		HealthyRange: "3:1 - 7:1",
		Status:       "HEALTHY",
		ByDepartment: []DataPoint{},
	}, nil
}

func (s *Service) GetSuccessionReadiness(ctx context.Context) (*SuccessionReadinessResponse, error) {
	return &SuccessionReadinessResponse{
		Period:             time.Now().Format("2006-01"),
		RolesWithSuccessors: 34,
		TotalKeyRoles:      50,
		CoverageRate:       68.0,
		Status:             "WARNING",
		ByDepartment:       []DataPoint{},
	}, nil
}

// =========================================================================
// Risk Detail
// =========================================================================

func (s *Service) GetRiskDetail(ctx context.Context, riskType string) (*RiskDetailResponse, error) {
	riskDefs := map[string]*RiskDetailResponse{
		"high-turnover": {
			RiskCode: "HIGH_TURNOVER", RiskName: "High Turnover Risk", RiskLevel: "HIGH",
			Value: 18.5, Threshold: 15.0, ExceededBy: 3.5,
			AffectedDepts: []DataPoint{
				{Label: "Sales", Value: 24.5},
				{Label: "Operations", Value: 18.2},
			},
			Recommendations: []string{"Conduct exit interviews", "Review compensation competitiveness", "Implement retention programs"},
		},
		"retirement": {
			RiskCode: "RETIREMENT", RiskName: "Retirement Risk", RiskLevel: "MEDIUM",
			Value: 8.5, Threshold: 10.0, ExceededBy: -1.5,
			AffectedDepts: []DataPoint{
				{Label: "Finance", Value: 12.0},
				{Label: "HR", Value: 10.5},
			},
			Recommendations: []string{"Identify critical knowledge transfer needs", "Plan succession pipeline for retiring roles", "Offer phased retirement options"},
		},
		"contract-expiry": {
			RiskCode: "CONTRACT_EXPIRY", RiskName: "Contract Expiration Risk", RiskLevel: "HIGH",
			Value: 42, Threshold: 30, ExceededBy: 12,
			AffectedDepts: []DataPoint{
				{Label: "IT", Value: 55},
				{Label: "Marketing", Value: 38},
			},
			Recommendations: []string{"Review and plan contract renewals", "Identify critical roles on expiring contracts", "Start recruitment pipeline early"},
		},
		"high-absenteeism": {
			RiskCode: "HIGH_ABSENTEEISM", RiskName: "High Absenteeism Risk", RiskLevel: "MEDIUM",
			Value: 12.5, Threshold: 10.0, ExceededBy: 2.5,
			AffectedDepts: []DataPoint{
				{Label: "Manufacturing", Value: 15.2},
				{Label: "Logistics", Value: 13.8},
			},
			Recommendations: []string{"Investigate root causes of absenteeism", "Implement wellness programs", "Review shift scheduling"},
		},
	}
	def, ok := riskDefs[riskType]
	if !ok {
		return nil, fmt.Errorf("unknown risk type: %s", riskType)
	}
	return def, nil
}

// =========================================================================
// KPI by Code
// =========================================================================

func (s *Service) GetKPIByCode(ctx context.Context, code string) (*KPIResponse, error) {
	list, _, err := s.repo.ListKPIs(ctx, "", "", code, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("KPI not found: %s", code)
	}
	return kpiToResponse(&list[0]), nil
}

// =========================================================================
// Health Scores
// =========================================================================

func (s *Service) GetHealthDashboard(ctx context.Context, period string) ([]HealthScoreResponse, error) {
	list, _, err := s.repo.ListHealthScores(ctx, period, "", 1, 10)
	if err != nil {
		return nil, err
	}
	responses := make([]HealthScoreResponse, 0, len(list))
	for _, hs := range list {
		responses = append(responses, *healthScoreToResponse(&hs))
	}
	return responses, nil
}

func (s *Service) ListHealthScores(ctx context.Context, period, orgID string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}
	list, total, err := s.repo.ListHealthScores(ctx, period, orgID, page, perPage)
	if err != nil {
		return nil, err
	}
	responses := make([]HealthScoreResponse, 0, len(list))
	for _, hs := range list {
		responses = append(responses, *healthScoreToResponse(&hs))
	}
	return &PaginatedResponse{
		Success: true, Data: responses, Page: page, PerPage: perPage,
		Total: total, TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// =========================================================================
// Helpers
// =========================================================================

func calcTotalPages(total int64, perPage int) int {
	pages := int(math.Ceil(float64(total) / float64(perPage)))
	if pages < 1 {
		return 1
	}
	return pages
}

func parseDateOrNow(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}
	return t
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// =========================================================================
// Response converters
// =========================================================================

func headcountPlanToResponse(h *WorkforcePlanningHeadcount) *HeadcountPlanResponse {
	variance := h.PlannedHC - h.ActualHC
	return &HeadcountPlanResponse{
		ID:             h.ID.String(),
		Period:         h.Period,
		OrganizationID: h.OrganizationID.String(),
		PlannedHC:      h.PlannedHC,
		ActualHC:       h.ActualHC,
		Variance:       variance,
		SnapshotDate:   h.SnapshotDate.Format("2006-01-02"),
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
	}
}

func forecastToResponse(f *WorkforceForecast) *ForecastResponse {
	return &ForecastResponse{
		ID:              f.ID.String(),
		Period:          f.Period,
		OrganizationID:  f.OrganizationID.String(),
		ForecastType:    f.ForecastType,
		Headcount:       f.Headcount,
		ConfidenceLevel: f.ConfidenceLevel,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}

func kpiToResponse(k *WorkforceKPI) *KPIResponse {
	r := &KPIResponse{
		ID:         k.ID.String(),
		Period:     k.Period,
		KpiCode:    k.KpiCode,
		KpiName:    k.KpiName,
		Value:      k.Value,
		Target:     k.Target,
		Unit:       k.Unit,
		Dimension:  k.Dimension,
		SnapshotAt: k.SnapshotAt.Format("2006-01-02"),
		CreatedAt:  k.CreatedAt,
	}
	if k.DimensionID != nil {
		r.DimensionID = *k.DimensionID
	}
	return r
}

func scenarioToResponse(s *WorkforceScenario) *ScenarioResponse {
	r := &ScenarioResponse{
		ID:           s.ID.String(),
		Name:         s.Name,
		Description:  s.Description,
		ScenarioType: s.ScenarioType,
		Parameters:   map[string]interface{}(s.Parameters),
		Results:      map[string]interface{}(s.Results),
		Status:       s.Status,
		CreatedBy:    s.CreatedBy.String(),
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
	if r.Parameters == nil {
		r.Parameters = map[string]interface{}{}
	}
	return r
}

func riskToResponse(ri *WorkforceRiskIndicator) *RiskResponse {
	r := &RiskResponse{
		ID:             ri.ID.String(),
		RiskCode:       ri.RiskCode,
		RiskName:       ri.RiskName,
		RiskLevel:      ri.RiskLevel,
		Score:          ri.Score,
		Threshold:      ri.Threshold,
		Recommendation: ri.Recommendation,
		SnapshotAt:     ri.SnapshotAt.Format("2006-01-02"),
	}
	if ri.DepartmentID != nil {
		r.DepartmentID = ri.DepartmentID.String()
	}
	return r
}

func healthScoreToResponse(hs *WorkforceHealthScore) *HealthScoreResponse {
	r := &HealthScoreResponse{
		ID:                 hs.ID.String(),
		Period:             hs.Period,
		OrganizationID:     hs.OrganizationID.String(),
		Score:              hs.Score,
		SpanOfControl:      hs.SpanOfControl,
		ManagerRatio:       hs.ManagerRatio,
		PromotionRate:      hs.PromotionRate,
		InternalHiringRate: hs.InternalHiringRate,
		SuccessionCoverage: hs.SuccessionCoverage,
		StabilityRatio:     hs.StabilityRatio,
		SnapshotAt:         hs.SnapshotAt.Format("2006-01-02"),
		CreatedAt:          hs.CreatedAt,
	}
	if hs.Components != nil {
		r.Components = map[string]interface{}(hs.Components)
	}
	return r
}

// =========================================================================
// Candidate Search
// =========================================================================

// CandidateSearch mencari posisi kosong (organisasi tanpa employment aktif di
// bawah Organization Summary active) beserta kandidat recruitment yang melamar
// ke requisition pada posisi tsb. S-3: mendukung filter posisi dan integrasi
// internal candidate eligible (Career Intelligence via narrow provider).
func (s *Service) CandidateSearch(ctx context.Context, search, position string, page, perPage int) (*PaginatedResponse, error) {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 || perPage > maxPerPage {
		perPage = defaultPerPage
	}

	var searchPtr, positionPtr *string
	if search != "" {
		searchPtr = &search
	}
	if position != "" {
		positionPtr = &position
	}

	rows, total, err := s.repo.CandidateSearchVacantOrgs(ctx, searchPtr, positionPtr, page, perPage)
	if err != nil {
		return nil, err
	}

	orgIDs := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		if uid, err := uuid.Parse(r.OrganizationID); err == nil {
			orgIDs = append(orgIDs, uid)
		}
	}

	// S-3: internal candidate eligible per posisi lowong (Career Intelligence).
	internalByOrg := s.resolveEligibleInternalCandidates(ctx, orgIDs)

	candRows, err := s.repo.CandidateSearchCandidatesByOrgIDs(ctx, orgIDs)
	if err != nil {
		return nil, err
	}

	candsByOrg := make(map[string][]CandidateSearchCandidate, len(rows))
	for _, cr := range candRows {
		// Dedupe per kandidat (satu kandidat bisa melamar ke beberapa requisition
		// di organisasi yang sama) — ambil aplikasi paling baru.
		cands := candsByOrg[cr.OrganizationID]
		dup := false
		for _, existing := range cands {
			if existing.ID == cr.ID {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		candsByOrg[cr.OrganizationID] = append(cands, CandidateSearchCandidate{
			ID:                cr.ID,
			FirstName:         cr.FirstName,
			LastName:          cr.LastName,
			Email:             cr.Email,
			Phone:             cr.Phone,
			CurrentTitle:      cr.CurrentTitle,
			CurrentCompany:    cr.CurrentCompany,
			Source:            cr.Source,
			ApplicationStatus: cr.ApplicationStatus,
			RequisitionTitle:  cr.RequisitionTitle,
		})
	}

	responses := make([]CandidateSearchPosition, 0, len(rows))
	for _, r := range rows {
		cands := candsByOrg[r.OrganizationID]
		if cands == nil {
			cands = []CandidateSearchCandidate{}
		}
		internals := internalByOrg[r.OrganizationID]
		if internals == nil {
			internals = []CandidateSearchInternalCandidate{}
		}
		responses = append(responses, CandidateSearchPosition{
			OrganizationID:        r.OrganizationID,
			OrganizationCode:      r.OrganizationCode,
			OrganizationName:      r.OrganizationName,
			SummaryID:             r.SummaryID,
			SummaryCode:           r.SummaryCode,
			SummaryDecreeNo:       r.SummaryDecreeNo,
			CandidateCount:        len(cands),
			Candidates:            cands,
			InternalCandidateCount: len(internals),
			InternalCandidates:    internals,
		})
	}

	return &PaginatedResponse{
		Success:    true,
		Data:       responses,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calcTotalPages(total, perPage),
	}, nil
}

// resolveEligibleInternalCandidates meminta Career Intelligence (via narrow
// provider) employee internal yang eligible untuk position-position pada
// organisasi lowong; hasil dikelompokkan per organization_id. Fail-safe:
// provider nil / error / tanpa positions → peta kosong.
func (s *Service) resolveEligibleInternalCandidates(ctx context.Context, orgIDs []uuid.UUID) map[string][]CandidateSearchInternalCandidate {
	result := map[string][]CandidateSearchInternalCandidate{}
	if s.internalEligProvider == nil || len(orgIDs) == 0 {
		return result
	}

	posRows, err := s.repo.CandidateSearchPositionsByOrgIDs(ctx, orgIDs)
	if err != nil {
		s.logger.Warn("candidate-search: failed to load positions for internal candidates",
			zap.Error(err))
		return result
	}
	if len(posRows) == 0 {
		return result
	}

	posIDs := make([]uuid.UUID, 0, len(posRows))
	for _, pr := range posRows {
		if uid, err := uuid.Parse(pr.ID); err == nil {
			posIDs = append(posIDs, uid)
		}
	}

	elig, err := s.internalEligProvider.EligibleCandidatesForPositions(ctx, posIDs)
	if err != nil {
		s.logger.Warn("candidate-search: internal eligibility provider failed; returning without internal candidates",
			zap.Error(err))
		return result
	}

	orgByPos := make(map[string]string, len(posRows))
	for _, pr := range posRows {
		orgByPos[pr.ID] = pr.OrganizationID
	}
	// Dedupe per employee per org: satu employee bisa eligible via beberapa
	// position/path pada org yang sama — konsisten dengan dedupe kandidat
	// eksternal di CandidateSearch.
	seen := map[string]map[string]bool{} // orgID -> employeeID
	for posID, cands := range elig {
		orgID := orgByPos[posID]
		if seen[orgID] == nil {
			seen[orgID] = map[string]bool{}
		}
		for _, c := range cands {
			if seen[orgID][c.EmployeeID] {
				continue
			}
			seen[orgID][c.EmployeeID] = true
			result[orgID] = append(result[orgID], CandidateSearchInternalCandidate{
				EmployeeID:          c.EmployeeID,
				Name:                c.Name,
				CurrentPositionID:   c.CurrentPositionID,
				CurrentPositionName: c.CurrentPositionName,
				SourceStepSequence:  c.SourceStepSequence,
			})
		}
	}
	return result
}
