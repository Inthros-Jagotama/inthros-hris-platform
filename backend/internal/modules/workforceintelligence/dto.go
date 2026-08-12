package workforceintelligence

import "time"

// =========================================================================
// Pagination
// =========================================================================

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// =========================================================================
// Workforce Planning DTOs
// =========================================================================

type CreateHeadcountPlanRequest struct {
	Period         string `json:"period" binding:"required,max=7"`
	OrganizationID string `json:"organization_id" binding:"required"`
	PlannedHC      int    `json:"planned_hc" binding:"required"`
	SnapshotDate   string `json:"snapshot_date"`
}

type UpdateHeadcountPlanRequest struct {
	PlannedHC    *int    `json:"planned_hc"`
	SnapshotDate *string `json:"snapshot_date"`
}

type HeadcountPlanResponse struct {
	ID             string    `json:"id"`
	Period         string    `json:"period"`
	OrganizationID string    `json:"organization_id"`
	PlannedHC      int       `json:"planned_hc"`
	ActualHC       int       `json:"actual_hc"`
	Variance       int       `json:"variance"`
	SnapshotDate   string    `json:"snapshot_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateForecastRequest struct {
	Period          string  `json:"period" binding:"required,max=7"`
	OrganizationID  string  `json:"organization_id" binding:"required"`
	ForecastType    string  `json:"forecast_type" binding:"required,oneof=DEMAND SUPPLY HIRING"`
	Headcount       int     `json:"headcount" binding:"required"`
	ConfidenceLevel float64 `json:"confidence_level"`
}

type UpdateForecastRequest struct {
	Headcount       *int     `json:"headcount"`
	ConfidenceLevel *float64 `json:"confidence_level"`
}

type ForecastResponse struct {
	ID              string    `json:"id"`
	Period          string    `json:"period"`
	OrganizationID  string    `json:"organization_id"`
	ForecastType    string    `json:"forecast_type"`
	Headcount       int       `json:"headcount"`
	ConfidenceLevel float64   `json:"confidence_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GapAnalysisResponse struct {
	Period          string                `json:"period"`
	Supply          int                   `json:"supply"`
	Demand          int                   `json:"demand"`
	Gap             int                   `json:"gap"`
	Status          string                `json:"status"` // SURPLUS / SHORTAGE / OPTIMAL
	Departments     []DepartmentGap       `json:"departments,omitempty"`
	// S-2: komponen hiring pipeline (Recruitment → WI)
	ExpectedHires   int                   `json:"expected_hires"`    // accepted offers
	OpenPositions   int                   `json:"open_positions"`    // slot requisition aktif
	FilledPositions int                   `json:"filled_positions"`  // slot requisition FILLED
	RemainingGap    int                   `json:"remaining_gap"`     // shortage − expected hires
}

type DepartmentGap struct {
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	Supply           int    `json:"supply"`
	Demand           int    `json:"demand"`
	Gap              int    `json:"gap"`
	Status           string `json:"status"`
}

type ProjectionResponse struct {
	Period       string `json:"period"`
	CurrentHC    int    `json:"current_hc"`
	ProjectedHC  int    `json:"projected_hc"`
	HiringNeeded int    `json:"hiring_needed"`
	RetirementCount int `json:"retirement_count"`
	GrowthRate   float64 `json:"growth_rate"`
}

// =========================================================================
// KPI DTOs
// =========================================================================

type KPIResponse struct {
	ID          string    `json:"id"`
	Period      string    `json:"period"`
	KpiCode     string    `json:"kpi_code"`
	KpiName     string    `json:"kpi_name"`
	Value       float64   `json:"value"`
	Target      *float64  `json:"target,omitempty"`
	Unit        string    `json:"unit"`
	Dimension   string    `json:"dimension"`
	DimensionID string    `json:"dimension_id,omitempty"`
	SnapshotAt  string    `json:"snapshot_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type KPISummaryResponse struct {
	Period      string        `json:"period"`
	TotalKPIs   int           `json:"total_kpis"`
	OnTarget    int           `json:"on_target"`
	BelowTarget int           `json:"below_target"`
	KPIs        []KPIResponse `json:"kpis,omitempty"`
}

// =========================================================================
// Analytics DTOs
// =========================================================================

type AnalyticsDashboardResponse struct {
	Period    string      `json:"period"`
	Summary   interface{} `json:"summary"`
	Trend     []DataPoint `json:"trend,omitempty"`
	Breakdown interface{} `json:"breakdown,omitempty"`
}

type DataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type HeadcountAnalytics struct {
	TotalHC          int            `json:"total_hc"`
	ActiveHC         int            `json:"active_hc"`
	ByDepartment     []DataPoint    `json:"by_department"`
	ByEmploymentType []DataPoint    `json:"by_employment_type"`
	ByGender         []DataPoint    `json:"by_gender"`
	Trend            []DataPoint    `json:"trend"`
}

type AttendanceAnalytics struct {
	AvgAttendanceRate float64       `json:"avg_attendance_rate"`
	AvgLateRate       float64       `json:"avg_late_rate"`
	AvgAbsentRate     float64       `json:"avg_absent_rate"`
	ByDepartment      []DataPoint   `json:"by_department"`
	Trend             []DataPoint   `json:"trend"`
}

type LeaveAnalytics struct {
	AvgUtilization float64       `json:"avg_utilization"`
	TotalDaysTaken int           `json:"total_days_taken"`
	ByType         []DataPoint   `json:"by_type"`
	ByDepartment   []DataPoint   `json:"by_department"`
}

type OvertimeAnalytics struct {
	AvgOTHours    float64       `json:"avg_ot_hours"`
	TotalOTCost   float64       `json:"total_ot_cost"`
	ByDepartment  []DataPoint   `json:"by_department"`
	Trend         []DataPoint   `json:"trend"`
}

type PayrollAnalytics struct {
	TotalPayroll    float64       `json:"total_payroll"`
	AvgSalary       float64       `json:"avg_salary"`
	ByDepartment    []DataPoint   `json:"by_department"`
	ByGrade         []DataPoint   `json:"by_grade"`
	Trend           []DataPoint   `json:"trend"`
}

type PerformanceAnalytics struct {
	AvgScore        float64       `json:"avg_score"`
	TopPerformerPct float64       `json:"top_performer_pct"`
	ByDepartment    []DataPoint   `json:"by_department"`
	Distribution    []DataPoint   `json:"distribution"`
}

type LearningAnalytics struct {
	CompletionRate float64       `json:"completion_rate"`
	AvgScore       float64       `json:"avg_score"`
	TotalHours     float64       `json:"total_hours"`
	ByCourse       []DataPoint   `json:"by_course"`
}

type RecruitmentAnalytics struct {
	TimeToHire    float64       `json:"time_to_hire"`
	TimeToFill    float64       `json:"time_to_fill"`
	OfferAcceptanceRate float64 `json:"offer_acceptance_rate"`
	CandidateMatchScore float64 `json:"candidate_match_score"`
	CostPerHire   float64       `json:"cost_per_hire"`
	BySource      []DataPoint   `json:"by_source"`
	SourceConversion []SourceConversionMetric `json:"source_conversion"`
	Funnel        []DataPoint   `json:"funnel"`
	// S-2: komponen hiring pipeline (Recruitment → WI)
	ExpectedHires   int       `json:"expected_hires"`
	OpenPositions   int       `json:"open_positions"`
	FilledPositions int       `json:"filled_positions"`
	RemainingGap    int       `json:"remaining_gap"`
	Pipeline        []DataPoint `json:"pipeline"`
}

// SourceConversionMetric adalah konversi kandidat→hire per channel rekrutmen
// (S-3): candidates (distinct), hires (aplikasi ACCEPTED), dan conversion rate
// (%).
type SourceConversionMetric struct {
	Source         string  `json:"source"`
	Candidates     int     `json:"candidates"`
	Hires          int     `json:"hires"`
	ConversionRate float64 `json:"conversion_rate"`
}

// =========================================================================
// Quality of Hire DTOs (S-6)
// =========================================================================

// QualityOfHireResponse — metrik agregat kualitas hire (S-6). Komposit dihitung
// dari komponen yang datanya tersedia: interview (interviews.score), probation
// (proxy onboarding COMPLETED), performance (final_score evaluasi selesai),
// retention (employment aktif). RecruitmentMatchScore & AssessmentScore tetap
// placeholder 0 — data kompetensi kandidat (G-9) & assessment belum dikumpulkan
// Recruitment (pola sama S-3).
type QualityOfHireResponse struct {
	OverallScore              float64                    `json:"overall_score"`
	HiresAnalyzed             int                        `json:"hires_analyzed"`
	RecruitmentMatchScore     float64                    `json:"recruitment_match_score"`
	InterviewScore            float64                    `json:"interview_score"`
	AssessmentScore           float64                    `json:"assessment_score"`
	OnboardingCompletionRate  float64                    `json:"onboarding_completion_rate"`
	PerformanceScore          float64                    `json:"performance_score"`
	RetentionRate             float64                    `json:"retention_rate"`
	BySource                  []QualityOfHireBreakdown   `json:"by_source"`
	ByRequisition             []QualityOfHireBreakdown   `json:"by_requisition"`
	ByOrganization            []QualityOfHireBreakdown   `json:"by_organization"`
}

// QualityOfHireBreakdown — agregasi per dimensi breakdown (source/requisition/
// organization): jumlah hire yang punya skor komposit (≥1 komponen berdata) +
// rata-rata skor komposit grup. Hire tanpa data sama sekali tidak muncul di
// breakdown (tetap dihitung di HiresAnalyzed).
type QualityOfHireBreakdown struct {
	Key   string  `json:"key"`
	Hires int     `json:"hires"`
	Score float64 `json:"score"`
}

type MovementAnalytics struct {
	PromotionCount int           `json:"promotion_count"`
	MutationCount  int           `json:"mutation_count"`
	ByDepartment   []DataPoint   `json:"by_department"`
	ByType         []DataPoint   `json:"by_type"`
}

// =========================================================================
// Capacity DTOs
// =========================================================================

type CapacityDashboardResponse struct {
	UtilizationRate float64         `json:"utilization_rate"`
	AvailableHC     int             `json:"available_hc"`
	ByDepartment    []DataPoint     `json:"by_department"`
	Bottlenecks     []Bottleneck    `json:"bottlenecks,omitempty"`
}

type Bottleneck struct {
	DepartmentID   string  `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	Utilization    float64 `json:"utilization"`
	Severity       string  `json:"severity"` // WARNING / CRITICAL
}

// =========================================================================
// Cost DTOs
// =========================================================================

type CostSummaryResponse struct {
	TotalPayroll     float64       `json:"total_payroll"`
	TotalBenefit     float64       `json:"total_benefit"`
	TotalLabor       float64       `json:"total_labor"`
	CostPerEmployee  float64       `json:"cost_per_employee"`
	ByDepartment     []DataPoint   `json:"by_department"`
	BudgetVsActual   []DataPoint   `json:"budget_vs_actual"`
}

// =========================================================================
// Risk DTOs
// =========================================================================

type RiskDashboardResponse struct {
	Period       string         `json:"period"`
	TotalRisks   int            `json:"total_risks"`
	HighRisks    int            `json:"high_risks"`
	CriticalRisks int           `json:"critical_risks"`
	Indicators   []RiskResponse `json:"indicators,omitempty"`
}

type RiskResponse struct {
	ID             string  `json:"id"`
	RiskCode       string  `json:"risk_code"`
	RiskName       string  `json:"risk_name"`
	RiskLevel      string  `json:"risk_level"`
	Score          float64 `json:"score"`
	Threshold      float64 `json:"threshold"`
	DepartmentID   string  `json:"department_id,omitempty"`
	Recommendation string  `json:"recommendation,omitempty"`
	SnapshotAt     string  `json:"snapshot_at"`
}

type UpdateRiskRequest struct {
	RiskLevel      *string `json:"risk_level" binding:"omitempty,oneof=LOW MEDIUM HIGH CRITICAL"`
	Recommendation *string `json:"recommendation"`
}

// =========================================================================
// Executive Dashboard DTOs
// =========================================================================

type ExecutiveSummaryResponse struct {
	TotalHC         int     `json:"total_hc"`
	HCGrowth        float64 `json:"hc_growth"`
	AttritionRate   float64 `json:"attrition_rate"`
	AvgCost         float64 `json:"avg_cost"`
	UtilizationRate float64 `json:"utilization_rate"`
	HealthScore     float64 `json:"health_score"`
	Period          string  `json:"period"`
}

type HiringProgressResponse struct {
	Planned    int `json:"planned"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Total      int `json:"total"`
}

// =========================================================================
// Scenario DTOs
// =========================================================================

type CreateScenarioRequest struct {
	Name         string                 `json:"name" binding:"required,max=150"`
	Description  string                 `json:"description"`
	ScenarioType string                 `json:"scenario_type" binding:"required,oneof=NEW_BRANCH REORG GROWTH REDUCTION RETIREMENT BUDGET"`
	Parameters   map[string]interface{} `json:"parameters" binding:"required"`
}

type UpdateScenarioRequest struct {
	Name        *string                `json:"name" binding:"omitempty,max=150"`
	Description *string                `json:"description"`
	Parameters  *map[string]interface{} `json:"parameters"`
}

type ScenarioResponse struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ScenarioType string                 `json:"scenario_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	Results      map[string]interface{} `json:"results,omitempty"`
	Status       string                 `json:"status"`
	CreatedBy    string                 `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// =========================================================================
// Organization Health DTOs
// =========================================================================

type HealthScoreResponse struct {
	ID                 string               `json:"id"`
	Period             string               `json:"period"`
	OrganizationID     string               `json:"organization_id"`
	Score              float64              `json:"score"`
	SpanOfControl      float64              `json:"span_of_control"`
	ManagerRatio       float64              `json:"manager_ratio"`
	PromotionRate      float64              `json:"promotion_rate"`
	InternalHiringRate float64              `json:"internal_hiring_rate"`
	SuccessionCoverage float64              `json:"succession_coverage"`
	StabilityRatio     float64              `json:"stability_ratio"`
	Components         map[string]interface{} `json:"components,omitempty"`
	SnapshotAt         string               `json:"snapshot_at"`
	CreatedAt          time.Time            `json:"created_at"`
}

// =========================================================================
// People Analytics DTOs
// =========================================================================

type CorrelationResponse struct {
	AnalysisType string       `json:"analysis_type"`
	Correlation  float64      `json:"correlation"`
	Strength     string       `json:"strength"` // STRONG / MODERATE / WEAK / NONE
	DataPoints   []DataPoint  `json:"data_points,omitempty"`
	Insight      string       `json:"insight,omitempty"`
}

// =========================================================================
// Capacity Forecast DTOs
// =========================================================================

type CapacityForecastResponse struct {
	Period          string       `json:"period"`
	ProjectedUtil   float64      `json:"projected_utilization"`
	CurrentCapacity int          `json:"current_capacity"`
	ProjectedNeeded int          `json:"projected_needed"`
	Gap             int          `json:"gap"`
	ByDepartment    []DataPoint  `json:"by_department"`
	Trend           []DataPoint  `json:"trend"`
}

// =========================================================================
// Cost Detail DTOs
// =========================================================================

type PayrollCostResponse struct {
	Period       string      `json:"period"`
	TotalSalary  float64     `json:"total_salary"`
	TotalAllowance float64   `json:"total_allowance"`
	TotalDeduction float64   `json:"total_deduction"`
	TotalBPJS    float64     `json:"total_bpjs"`
	ByGrade      []DataPoint `json:"by_grade"`
	ByComponent  []DataPoint `json:"by_component"`
}

type CostPerEmployeeResponse struct {
	Period           string      `json:"period"`
	AvgCostPerEmployee float64   `json:"avg_cost_per_employee"`
	MedianCost       float64     `json:"median_cost"`
	MinCost          float64     `json:"min_cost"`
	MaxCost          float64     `json:"max_cost"`
	ByDepartment     []DataPoint `json:"by_department"`
	ByGrade          []DataPoint `json:"by_grade"`
}

// =========================================================================
// Executive Detail DTOs
// =========================================================================

type ExecutiveTrendResponse struct {
	Period   string      `json:"period"`
	Trend    []DataPoint `json:"trend"`
	Current  float64     `json:"current"`
	Change   float64     `json:"change"` // absolute or pct change
}

type ExecutiveCapacityResponse struct {
	UtilizationRate  float64      `json:"utilization_rate"`
	AvailableHC      int          `json:"available_hc"`
	ActiveDeptCount  int          `json:"active_dept_count"`
	BottleneckCount  int          `json:"bottleneck_count"`
	ByDepartment     []DataPoint  `json:"by_department"`
}

type ExecutiveRiskOverviewResponse struct {
	TotalRisks      int              `json:"total_risks"`
	HighRiskCount   int              `json:"high_risk_count"`
	CriticalCount   int              `json:"critical_count"`
	ByDepartment    []DataPoint      `json:"by_department"`
	ByCategory      []DataPoint      `json:"by_category"`
}

type ExecutiveHealthScoreResponse struct {
	Score              float64              `json:"score"`
	SpanOfControl      float64              `json:"span_of_control"`
	ManagerRatio       float64              `json:"manager_ratio"`
	InternalHiringRate float64              `json:"internal_hiring_rate"`
	SuccessionCoverage float64              `json:"succession_coverage"`
	Status             string               `json:"status"` // HEALTHY / WARNING / CRITICAL
	Components         map[string]interface{} `json:"components"`
}

// =========================================================================
// Risk Detail DTOs
// =========================================================================

type RiskDetailResponse struct {
	RiskCode    string        `json:"risk_code"`
	RiskName    string        `json:"risk_name"`
	RiskLevel   string        `json:"risk_level"`
	Value       float64       `json:"value"`
	Threshold   float64       `json:"threshold"`
	ExceededBy  float64       `json:"exceeded_by"`
	AffectedDepts []DataPoint `json:"affected_departments"`
	Trend       []DataPoint   `json:"trend"`
	Recommendations []string  `json:"recommendations,omitempty"`
}

// =========================================================================
// Health Detail DTOs
// =========================================================================

type SpanOfControlResponse struct {
	Period    string      `json:"period"`
	AvgRatio  float64     `json:"avg_ratio"`
	HealthyRange string  `json:"healthy_range"` // 3:1 - 7:1
	Status    string      `json:"status"`
	ByDepartment []DataPoint `json:"by_department"`
}

type SuccessionReadinessResponse struct {
	Period           string      `json:"period"`
	RolesWithSuccessors int      `json:"roles_with_successors"`
	TotalKeyRoles    int         `json:"total_key_roles"`
	CoverageRate     float64     `json:"coverage_rate"`
	Status           string      `json:"status"` // HEALTHY / WARNING / CRITICAL
	ByDepartment     []DataPoint `json:"by_department"`
}

// =========================================================================
// Candidate Search DTOs
// =========================================================================

// CandidateSearchPosition adalah posisi kosong: organisasi tanpa employment
// aktif yang berada di bawah Organization Summary berstatus active.
type CandidateSearchPosition struct {
	OrganizationID   string                     `json:"organization_id"`
	OrganizationCode string                     `json:"organization_code"`
	OrganizationName string                     `json:"organization_name"`
	SummaryID        string                     `json:"summary_id"`
	SummaryCode      string                     `json:"summary_code"`
	SummaryDecreeNo  string                     `json:"summary_decree_no"`
	CandidateCount   int                        `json:"candidate_count"`
	Candidates       []CandidateSearchCandidate `json:"candidates"`
	// S-3: integrasi internal candidate eligible (Career Intelligence) —
	// employee internal yang eligible via career path menuju position pada
	// posisi lowong tsb. Kosong bila provider belum di-wire atau tidak ada.
	InternalCandidateCount int                             `json:"internal_candidate_count"`
	InternalCandidates     []CandidateSearchInternalCandidate `json:"internal_candidates"`
}

// CandidateSearchInternalCandidate adalah employee internal yang eligible
// (perhitungan Career Intelligence — S-3/S-4) untuk posisi pada organisasi
// lowong, melengkapi kandidat eksternal dari pool recruitment.
type CandidateSearchInternalCandidate struct {
	EmployeeID          string `json:"employee_id"`
	Name                string `json:"name"`
	CurrentPositionID   string `json:"current_position_id,omitempty"`
	CurrentPositionName string `json:"current_position_name,omitempty"`
	SourceStepSequence  int    `json:"source_step_sequence"`
}

// CandidateSearchCandidate adalah kandidat dari pool recruitment yang melamar
// ke requisition pada organisasi lowong tsb (status aplikasi bukan rejected/withdrawn).
type CandidateSearchCandidate struct {
	ID                string  `json:"id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone,omitempty"`
	CurrentTitle      *string `json:"current_title,omitempty"`
	CurrentCompany    *string `json:"current_company,omitempty"`
	Source            string  `json:"source"`
	ApplicationStatus string  `json:"application_status"`
	RequisitionTitle  string  `json:"requisition_title"`
}
