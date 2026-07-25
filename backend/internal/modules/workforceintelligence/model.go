package workforceintelligence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSON type for flexible data storage
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// =========================================================================
// WorkforcePlanningHeadcount — Snapshot of planned vs actual HC
// =========================================================================

type WorkforcePlanningHeadcount struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Period         string     `gorm:"type:char(7);not null;index:idx_wf_plan_period" json:"period"`
	OrganizationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_wf_plan_org" json:"organization_id"`
	PlannedHC      int        `gorm:"not null;default:0" json:"planned_hc"`
	ActualHC       int        `gorm:"not null;default:0" json:"actual_hc"`
	SnapshotDate   time.Time  `gorm:"type:date;not null" json:"snapshot_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (WorkforcePlanningHeadcount) TableName() string { return "workforce_planning_headcounts" }

// =========================================================================
// WorkforceForecast — Forecast data for headcount demand/supply/hiring
// =========================================================================

type WorkforceForecast struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Period          string    `gorm:"type:char(7);not null;index:idx_wf_forecast_period" json:"period"`
	OrganizationID  uuid.UUID `gorm:"type:char(36);not null;index:idx_wf_forecast_org" json:"organization_id"`
	ForecastType    string    `gorm:"type:varchar(30);not null;index:idx_wf_forecast_type" json:"forecast_type"` // DEMAND / SUPPLY / HIRING
	Headcount       int       `gorm:"not null;default:0" json:"headcount"`
	ConfidenceLevel float64   `gorm:"type:decimal(5,2);default:0" json:"confidence_level"`
	Parameters      JSON      `gorm:"type:json" json:"parameters"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (WorkforceForecast) TableName() string { return "workforce_forecasts" }

// =========================================================================
// WorkforceKPI — Pre-computed KPI snapshots per period
// =========================================================================

type WorkforceKPI struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Period      string     `gorm:"type:char(7);not null;index:idx_wf_kpi_period" json:"period"`
	KpiCode     string     `gorm:"type:varchar(50);not null;index:idx_wf_kpi_code" json:"kpi_code"`
	KpiName     string     `gorm:"type:varchar(100)" json:"kpi_name"`
	Value       float64    `gorm:"type:decimal(15,2);not null" json:"value"`
	Target      *float64   `gorm:"type:decimal(15,2)" json:"target,omitempty"`
	Unit        string     `gorm:"type:varchar(20)" json:"unit"` // PCT / COUNT / AMOUNT / RATIO
	Dimension   string     `gorm:"type:varchar(30);not null" json:"dimension"` // COMPANY / DEPARTMENT / BRANCH
	DimensionID *string    `gorm:"type:char(36)" json:"dimension_id,omitempty"`
	SnapshotAt  time.Time  `gorm:"type:date;not null" json:"snapshot_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (WorkforceKPI) TableName() string { return "workforce_kpis" }

// =========================================================================
// WorkforceAnalyticsCache — Aggregated analytics cache (JSON blob)
// =========================================================================

type WorkforceAnalyticsCache struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CacheKey  string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"cache_key"`
	CacheType string     `gorm:"type:varchar(50);not null;index:idx_wf_cache_type" json:"cache_type"` // HC / ATTENDANCE / PAYROLL / etc.
	Data      JSON       `gorm:"type:json;not null" json:"data"`
	Period    string     `gorm:"type:char(7);index:idx_wf_cache_period" json:"period"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (WorkforceAnalyticsCache) TableName() string { return "workforce_analytics_cache" }

// =========================================================================
// WorkforceScenario — Saved simulation scenarios
// =========================================================================

type WorkforceScenario struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(150);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	ScenarioType string         `gorm:"type:varchar(50);not null;index:idx_wf_scenario_type" json:"scenario_type"` // NEW_BRANCH / REORG / GROWTH / REDUCTION
	Parameters   JSON           `gorm:"type:json;not null" json:"parameters"`
	Results      JSON           `gorm:"type:json" json:"results"`
	Status       string         `gorm:"type:varchar(20);default:DRAFT" json:"status"` // DRAFT / COMPLETED / ARCHIVED
	CreatedBy    uuid.UUID      `gorm:"type:char(36)" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_wf_scenario_deleted_at" json:"deleted_at,omitempty"`
}

func (WorkforceScenario) TableName() string { return "workforce_scenarios" }

// =========================================================================
// WorkforceRiskIndicator — Risk monitoring record
// =========================================================================

type WorkforceRiskIndicator struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Period         string     `gorm:"type:char(7);not null;index:idx_wf_risk_period" json:"period"`
	RiskCode       string     `gorm:"type:varchar(50);not null;index:idx_wf_risk_code" json:"risk_code"`
	RiskName       string     `gorm:"type:varchar(100)" json:"risk_name"`
	RiskLevel      string     `gorm:"type:varchar(20);not null" json:"risk_level"` // LOW / MEDIUM / HIGH / CRITICAL
	Score          float64    `gorm:"type:decimal(10,2)" json:"score"`
	Threshold      float64    `gorm:"type:decimal(10,2)" json:"threshold"`
	DepartmentID   *uuid.UUID `gorm:"type:char(36)" json:"department_id,omitempty"`
	Recommendation string     `gorm:"type:text" json:"recommendation"`
	SnapshotAt     time.Time  `gorm:"type:date;not null" json:"snapshot_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (WorkforceRiskIndicator) TableName() string { return "workforce_risk_indicators" }

// =========================================================================
// WorkforceHealthScore — Organization health composite score
// =========================================================================

type WorkforceHealthScore struct {
	ID                 uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	Period             string     `gorm:"type:char(7);not null;index:idx_wf_health_period" json:"period"`
	OrganizationID     uuid.UUID  `gorm:"type:char(36);not null;index:idx_wf_health_org" json:"organization_id"`
	Score              float64    `gorm:"type:decimal(5,2)" json:"score"`
	SpanOfControl      float64    `gorm:"type:decimal(5,2)" json:"span_of_control"`
	ManagerRatio       float64    `gorm:"type:decimal(5,2)" json:"manager_ratio"`
	PromotionRate      float64    `gorm:"type:decimal(5,2)" json:"promotion_rate"`
	InternalHiringRate float64    `gorm:"type:decimal(5,2)" json:"internal_hiring_rate"`
	SuccessionCoverage float64    `gorm:"type:decimal(5,2)" json:"succession_coverage"`
	StabilityRatio     float64    `gorm:"type:decimal(5,2)" json:"stability_ratio"`
	Components         JSON       `gorm:"type:json" json:"components"`
	SnapshotAt         time.Time  `gorm:"type:date;not null" json:"snapshot_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

func (WorkforceHealthScore) TableName() string { return "workforce_health_scores" }
