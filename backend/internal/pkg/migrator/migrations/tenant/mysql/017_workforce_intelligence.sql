-- 017_workforce_intelligence.sql
-- Workforce Intelligence & Strategic Workforce Planning Module
-- Tabel untuk workforce analytics, headcount planning, KPI, risk monitoring,
-- scenario simulation, dan organization health scoring

-- =========================================================================
-- Workforce Planning Headcounts (Perencanaan vs Realisasi Headcount)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_planning_headcounts (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    period            CHAR(7)      NOT NULL COMMENT 'Format: 2026-Q1',
    organization_id   CHAR(36)     NOT NULL,
    planned_hc        INT          NOT NULL DEFAULT 0,
    actual_hc         INT          NOT NULL DEFAULT 0,
    snapshot_date     DATE         NOT NULL,
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_wf_plan_period (period),
    INDEX idx_wf_plan_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce Forecasts (Forecast Headcount Demand/Supply/Hiring)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_forecasts (
    id                CHAR(36)     NOT NULL PRIMARY KEY,
    period            CHAR(7)      NOT NULL COMMENT 'Format: 2026-Q1',
    organization_id   CHAR(36)     NOT NULL,
    forecast_type     VARCHAR(30)  NOT NULL COMMENT 'DEMAND / SUPPLY / HIRING',
    headcount         INT          NOT NULL DEFAULT 0,
    confidence_level  DECIMAL(5,2) NOT NULL DEFAULT 0.00,
    parameters        JSON         NULL COMMENT 'Parameter forecasting dalam format JSON',
    created_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at        TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_wf_forecast_period (period),
    INDEX idx_wf_forecast_org (organization_id),
    INDEX idx_wf_forecast_type (forecast_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce KPIs (Pre-computed KPI Snapshots)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_kpis (
    id              CHAR(36)      NOT NULL PRIMARY KEY,
    period          CHAR(7)       NOT NULL COMMENT 'Format: 2026-Q1',
    kpi_code        VARCHAR(50)   NOT NULL COMMENT 'Kode unik KPI: HC_TOTAL, ATTRITION_RATE, etc.',
    kpi_name        VARCHAR(100)  NULL,
    value           DECIMAL(15,2) NOT NULL,
    target          DECIMAL(15,2) NULL,
    unit            VARCHAR(20)   NULL COMMENT 'PCT / COUNT / AMOUNT / RATIO',
    dimension       VARCHAR(30)   NOT NULL DEFAULT 'COMPANY' COMMENT 'COMPANY / DEPARTMENT / BRANCH',
    dimension_id    CHAR(36)      NULL,
    snapshot_at     DATE          NOT NULL,
    created_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_wf_kpi_period (period),
    INDEX idx_wf_kpi_code (kpi_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce Analytics Cache (Pre-computed Analytics Cache)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_analytics_cache (
    id              CHAR(36)      NOT NULL PRIMARY KEY,
    cache_key       VARCHAR(100)  NOT NULL COMMENT 'Unique cache key: hc_trend_2026Q1',
    cache_type      VARCHAR(50)   NOT NULL COMMENT 'HC / ATTENDANCE / PAYROLL / etc.',
    data            JSON          NOT NULL COMMENT 'Aggregated data dalam JSON',
    period          CHAR(7)       NULL COMMENT 'Format: 2026-Q1',
    expires_at      TIMESTAMP(6)  NULL,
    created_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE INDEX idx_wf_cache_key (cache_key),
    INDEX idx_wf_cache_type (cache_type),
    INDEX idx_wf_cache_period (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce Scenarios (Saved Simulation Scenarios)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_scenarios (
    id              CHAR(36)      NOT NULL PRIMARY KEY,
    name            VARCHAR(150)  NOT NULL,
    description     TEXT          NULL,
    scenario_type   VARCHAR(50)   NOT NULL COMMENT 'NEW_BRANCH / REORG / GROWTH / REDUCTION / RETIREMENT / BUDGET',
    parameters      JSON          NOT NULL COMMENT 'Input parameters dalam JSON',
    results         JSON          NULL COMMENT 'Simulation results dalam JSON',
    status          VARCHAR(20)   NOT NULL DEFAULT 'DRAFT' COMMENT 'DRAFT / COMPLETED / ARCHIVED',
    created_by      CHAR(36)      NULL,
    created_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at      TIMESTAMP(6)  NULL DEFAULT NULL,
    INDEX idx_wf_scenario_type (scenario_type),
    INDEX idx_wf_scenario_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce Risk Indicators (Risk Monitoring Records)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_risk_indicators (
    id              CHAR(36)      NOT NULL PRIMARY KEY,
    period          CHAR(7)       NOT NULL COMMENT 'Format: 2026-Q1',
    risk_code       VARCHAR(50)   NOT NULL COMMENT 'Kode unik risk: HIGH_TURNOVER, RETIREMENT, etc.',
    risk_name       VARCHAR(100)  NULL,
    risk_level      VARCHAR(20)   NOT NULL COMMENT 'LOW / MEDIUM / HIGH / CRITICAL',
    score           DECIMAL(10,2) NULL,
    threshold       DECIMAL(10,2) NULL,
    department_id   CHAR(36)      NULL,
    recommendation  TEXT          NULL,
    snapshot_at     DATE          NOT NULL,
    created_at      TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_wf_risk_period (period),
    INDEX idx_wf_risk_code (risk_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =========================================================================
-- Workforce Health Scores (Organization Health Composite)
-- =========================================================================
CREATE TABLE IF NOT EXISTS workforce_health_scores (
    id                  CHAR(36)      NOT NULL PRIMARY KEY,
    period              CHAR(7)       NOT NULL COMMENT 'Format: 2026-Q1',
    organization_id     CHAR(36)      NOT NULL,
    score               DECIMAL(5,2)  NULL COMMENT 'Composite health score 0-100',
    span_of_control     DECIMAL(5,2)  NULL,
    manager_ratio       DECIMAL(5,2)  NULL,
    promotion_rate      DECIMAL(5,2)  NULL,
    internal_hiring_rate DECIMAL(5,2) NULL,
    succession_coverage DECIMAL(5,2)  NULL,
    stability_ratio     DECIMAL(5,2)  NULL,
    components          JSON          NULL COMMENT 'Component breakdown dalam JSON',
    snapshot_at         DATE          NOT NULL,
    created_at          TIMESTAMP(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_wf_health_period (period),
    INDEX idx_wf_health_org (organization_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
