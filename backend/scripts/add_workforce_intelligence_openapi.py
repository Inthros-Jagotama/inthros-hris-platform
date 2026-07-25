#!/usr/bin/env python3
"""
Add Workforce Intelligence & Strategic Planning module (68 endpoints, ~25 schemas)
to the existing openapi.json document.
"""

import json, sys, os

OPENAPI_PATH = os.path.join(os.path.dirname(__file__), '..', 'internal', 'pkg', 'docs', 'openapi.json')

with open(OPENAPI_PATH, 'r', encoding='utf-8') as f:
    spec = json.load(f)

TAG = "Tenant: Workforce Intelligence & Strategic Planning"
TAG_DESC = "Workforce Intelligence & Strategic Workforce Planning — strategic analytics layer for headcount planning, forecasting, gap analysis, KPI monitoring, workforce analytics (headcount, attendance, leave, overtime, payroll, performance, learning, recruitment, movement), capacity planning, cost analytics, risk monitoring, executive dashboards, scenario simulation, organization health scoring, and people analytics (training-vs-performance, overtime-vs-productivity, etc.). Read-only analytics module aggregating data from all operational HR modules."

# ─────────────────────────────────────────────────────────────────────────────
# 1. Add tag (if not already present)
# ─────────────────────────────────────────────────────────────────────────────
existing_tags = [t['name'] for t in spec['tags']]
if TAG not in existing_tags:
    # Insert before Health tag or at end
    spec['tags'].append({"name": TAG, "description": TAG_DESC})

# ─────────────────────────────────────────────────────────────────────────────
# 2. Schemas
# ─────────────────────────────────────────────────────────────────────────────
if 'components' not in spec:
    spec['components'] = {}
if 'schemas' not in spec['components']:
    spec['components']['schemas'] = {}

schemas = {
    "CreateHeadcountPlanRequest": {
        "type": "object",
        "required": ["period", "organization_id", "planned_hc"],
        "properties": {
            "period": {"type": "string", "maxLength": 7, "example": "2026-01"},
            "organization_id": {"type": "string", "format": "uuid"},
            "planned_hc": {"type": "integer", "example": 150},
            "snapshot_date": {"type": "string", "format": "date"}
        }
    },
    "UpdateHeadcountPlanRequest": {
        "type": "object",
        "properties": {
            "planned_hc": {"type": "integer"},
            "snapshot_date": {"type": "string", "format": "date"}
        }
    },
    "HeadcountPlanResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "period": {"type": "string"},
            "organization_id": {"type": "string", "format": "uuid"},
            "planned_hc": {"type": "integer"},
            "actual_hc": {"type": "integer"},
            "variance": {"type": "integer"},
            "snapshot_date": {"type": "string", "format": "date"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "CreateForecastRequest": {
        "type": "object",
        "required": ["period", "organization_id", "forecast_type", "headcount"],
        "properties": {
            "period": {"type": "string", "maxLength": 7, "example": "2026-01"},
            "organization_id": {"type": "string", "format": "uuid"},
            "forecast_type": {"type": "string", "enum": ["DEMAND", "SUPPLY", "HIRING"]},
            "headcount": {"type": "integer"},
            "confidence_level": {"type": "number", "format": "float"}
        }
    },
    "UpdateForecastRequest": {
        "type": "object",
        "properties": {
            "headcount": {"type": "integer"},
            "confidence_level": {"type": "number", "format": "float"}
        }
    },
    "ForecastResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "period": {"type": "string"},
            "organization_id": {"type": "string", "format": "uuid"},
            "forecast_type": {"type": "string"},
            "headcount": {"type": "integer"},
            "confidence_level": {"type": "number", "format": "float"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "GapAnalysisResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "supply": {"type": "integer"},
            "demand": {"type": "integer"},
            "gap": {"type": "integer"},
            "status": {"type": "string", "enum": ["SURPLUS", "SHORTAGE", "OPTIMAL"]},
            "departments": {
                "type": "array",
                "items": {
                    "type": "object",
                    "properties": {
                        "organization_id": {"type": "string", "format": "uuid"},
                        "organization_name": {"type": "string"},
                        "supply": {"type": "integer"},
                        "demand": {"type": "integer"},
                        "gap": {"type": "integer"},
                        "status": {"type": "string"}
                    }
                }
            }
        }
    },
    "ProjectionResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "current_hc": {"type": "integer"},
            "projected_hc": {"type": "integer"},
            "hiring_needed": {"type": "integer"},
            "retirement_count": {"type": "integer"},
            "growth_rate": {"type": "number", "format": "float"}
        }
    },
    "KPIResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "period": {"type": "string"},
            "kpi_code": {"type": "string"},
            "kpi_name": {"type": "string"},
            "value": {"type": "number", "format": "float"},
            "target": {"type": "number", "format": "float"},
            "unit": {"type": "string"},
            "dimension": {"type": "string"},
            "dimension_id": {"type": "string"},
            "snapshot_at": {"type": "string", "format": "date"},
            "created_at": {"type": "string", "format": "date-time"}
        }
    },
    "KPISummaryResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "total_kpis": {"type": "integer"},
            "on_target": {"type": "integer"},
            "below_target": {"type": "integer"},
            "kpis": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/KPIResponse"}
            }
        }
    },
    "HeadcountAnalytics": {
        "type": "object",
        "properties": {
            "total_hc": {"type": "integer"},
            "active_hc": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_employment_type": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_gender": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "AttendanceAnalytics": {
        "type": "object",
        "properties": {
            "avg_attendance_rate": {"type": "number", "format": "float"},
            "avg_late_rate": {"type": "number", "format": "float"},
            "avg_absent_rate": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "LeaveAnalytics": {
        "type": "object",
        "properties": {
            "avg_utilization": {"type": "number", "format": "float"},
            "total_days_taken": {"type": "integer"},
            "by_type": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "OvertimeAnalytics": {
        "type": "object",
        "properties": {
            "avg_ot_hours": {"type": "number", "format": "float"},
            "total_ot_cost": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "PayrollAnalytics": {
        "type": "object",
        "properties": {
            "total_payroll": {"type": "number", "format": "float"},
            "avg_salary": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_grade": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "PerformanceAnalytics": {
        "type": "object",
        "properties": {
            "avg_score": {"type": "number", "format": "float"},
            "top_performer_pct": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "distribution": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "LearningAnalytics": {
        "type": "object",
        "properties": {
            "completion_rate": {"type": "number", "format": "float"},
            "avg_score": {"type": "number", "format": "float"},
            "total_hours": {"type": "number", "format": "float"},
            "by_course": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "RecruitmentAnalytics": {
        "type": "object",
        "properties": {
            "time_to_hire": {"type": "number", "format": "float"},
            "cost_per_hire": {"type": "number", "format": "float"},
            "by_source": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "funnel": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "MovementAnalytics": {
        "type": "object",
        "properties": {
            "promotion_count": {"type": "integer"},
            "mutation_count": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_type": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "CapacityDashboardResponse": {
        "type": "object",
        "properties": {
            "utilization_rate": {"type": "number", "format": "float"},
            "available_hc": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "bottlenecks": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/Bottleneck"}
            }
        }
    },
    "Bottleneck": {
        "type": "object",
        "properties": {
            "department_id": {"type": "string", "format": "uuid"},
            "department_name": {"type": "string"},
            "utilization": {"type": "number", "format": "float"},
            "severity": {"type": "string", "enum": ["WARNING", "CRITICAL"]}
        }
    },
    "CostSummaryResponse": {
        "type": "object",
        "properties": {
            "total_payroll": {"type": "number", "format": "float"},
            "total_benefit": {"type": "number", "format": "float"},
            "total_labor": {"type": "number", "format": "float"},
            "cost_per_employee": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "budget_vs_actual": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "RiskDashboardResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "total_risks": {"type": "integer"},
            "high_risks": {"type": "integer"},
            "critical_risks": {"type": "integer"},
            "indicators": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/RiskResponse"}
            }
        }
    },
    "RiskResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "risk_code": {"type": "string"},
            "risk_name": {"type": "string"},
            "risk_level": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "CRITICAL"]},
            "score": {"type": "number", "format": "float"},
            "threshold": {"type": "number", "format": "float"},
            "department_id": {"type": "string", "format": "uuid"},
            "recommendation": {"type": "string"},
            "snapshot_at": {"type": "string", "format": "date"}
        }
    },
    "UpdateRiskRequest": {
        "type": "object",
        "properties": {
            "risk_level": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "CRITICAL"]},
            "recommendation": {"type": "string"}
        }
    },
    "ExecutiveSummaryResponse": {
        "type": "object",
        "properties": {
            "total_hc": {"type": "integer"},
            "hc_growth": {"type": "number", "format": "float"},
            "attrition_rate": {"type": "number", "format": "float"},
            "avg_cost": {"type": "number", "format": "float"},
            "utilization_rate": {"type": "number", "format": "float"},
            "health_score": {"type": "number", "format": "float"},
            "period": {"type": "string"}
        }
    },
    "HiringProgressResponse": {
        "type": "object",
        "properties": {
            "planned": {"type": "integer"},
            "in_progress": {"type": "integer"},
            "completed": {"type": "integer"},
            "total": {"type": "integer"}
        }
    },
    "CreateScenarioRequest": {
        "type": "object",
        "required": ["name", "scenario_type", "parameters"],
        "properties": {
            "name": {"type": "string", "maxLength": 150},
            "description": {"type": "string"},
            "scenario_type": {"type": "string", "enum": ["NEW_BRANCH", "REORG", "GROWTH", "REDUCTION", "RETIREMENT", "BUDGET"]},
            "parameters": {"type": "object"}
        }
    },
    "UpdateScenarioRequest": {
        "type": "object",
        "properties": {
            "name": {"type": "string", "maxLength": 150},
            "description": {"type": "string"},
            "parameters": {"type": "object"}
        }
    },
    "ScenarioResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "description": {"type": "string"},
            "scenario_type": {"type": "string"},
            "parameters": {"type": "object"},
            "results": {"type": "object"},
            "status": {"type": "string"},
            "created_by": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "HealthScoreResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "period": {"type": "string"},
            "organization_id": {"type": "string", "format": "uuid"},
            "score": {"type": "number", "format": "float"},
            "span_of_control": {"type": "number", "format": "float"},
            "manager_ratio": {"type": "number", "format": "float"},
            "promotion_rate": {"type": "number", "format": "float"},
            "internal_hiring_rate": {"type": "number", "format": "float"},
            "succession_coverage": {"type": "number", "format": "float"},
            "stability_ratio": {"type": "number", "format": "float"},
            "components": {"type": "object"},
            "snapshot_at": {"type": "string", "format": "date"},
            "created_at": {"type": "string", "format": "date-time"}
        }
    },
    "CorrelationResponse": {
        "type": "object",
        "properties": {
            "analysis_type": {"type": "string"},
            "correlation": {"type": "number", "format": "float"},
            "strength": {"type": "string", "enum": ["STRONG", "MODERATE", "WEAK", "NONE"]},
            "data_points": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "insight": {"type": "string"}
        }
    },
    "CapacityForecastResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "projected_utilization": {"type": "number", "format": "float"},
            "current_capacity": {"type": "integer"},
            "projected_needed": {"type": "integer"},
            "gap": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "PayrollCostResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "total_salary": {"type": "number", "format": "float"},
            "total_allowance": {"type": "number", "format": "float"},
            "total_deduction": {"type": "number", "format": "float"},
            "total_bpjs": {"type": "number", "format": "float"},
            "by_grade": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_component": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "CostPerEmployeeResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "avg_cost_per_employee": {"type": "number", "format": "float"},
            "median_cost": {"type": "number", "format": "float"},
            "min_cost": {"type": "number", "format": "float"},
            "max_cost": {"type": "number", "format": "float"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_grade": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "ExecutiveTrendResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "current": {"type": "number", "format": "float"},
            "change": {"type": "number", "format": "float"}
        }
    },
    "ExecutiveCapacityResponse": {
        "type": "object",
        "properties": {
            "utilization_rate": {"type": "number", "format": "float"},
            "available_hc": {"type": "integer"},
            "active_dept_count": {"type": "integer"},
            "bottleneck_count": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "ExecutiveRiskOverviewResponse": {
        "type": "object",
        "properties": {
            "total_risks": {"type": "integer"},
            "high_risk_count": {"type": "integer"},
            "critical_count": {"type": "integer"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "by_category": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "ExecutiveHealthScoreResponse": {
        "type": "object",
        "properties": {
            "score": {"type": "number", "format": "float"},
            "span_of_control": {"type": "number", "format": "float"},
            "manager_ratio": {"type": "number", "format": "float"},
            "internal_hiring_rate": {"type": "number", "format": "float"},
            "succession_coverage": {"type": "number", "format": "float"},
            "status": {"type": "string", "enum": ["HEALTHY", "WARNING", "CRITICAL"]},
            "components": {"type": "object"}
        }
    },
    "RiskDetailResponse": {
        "type": "object",
        "properties": {
            "risk_code": {"type": "string"},
            "risk_name": {"type": "string"},
            "risk_level": {"type": "string"},
            "value": {"type": "number", "format": "float"},
            "threshold": {"type": "number", "format": "float"},
            "exceeded_by": {"type": "number", "format": "float"},
            "affected_departments": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "trend": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}},
            "recommendations": {"type": "array", "items": {"type": "string"}}
        }
    },
    "SpanOfControlResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "avg_ratio": {"type": "number", "format": "float"},
            "healthy_range": {"type": "string"},
            "status": {"type": "string"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "SuccessionReadinessResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "roles_with_successors": {"type": "integer"},
            "total_key_roles": {"type": "integer"},
            "coverage_rate": {"type": "number", "format": "float"},
            "status": {"type": "string"},
            "by_department": {"type": "array", "items": {"$ref": "#/components/schemas/DataPoint"}}
        }
    },
    "DataPoint": {
        "type": "object",
        "properties": {
            "label": {"type": "string"},
            "value": {"type": "number", "format": "float"}
        }
    },
    "PaginatedResponseWF": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {"type": "array", "items": {"type": "object"}},
            "page": {"type": "integer"},
            "per_page": {"type": "integer"},
            "total": {"type": "integer"},
            "total_pages": {"type": "integer"}
        }
    }
}

# Add only schemas that don't already exist
for name, schema in schemas.items():
    if name not in spec['components']['schemas']:
        spec['components']['schemas'][name] = schema

# ─────────────────────────────────────────────────────────────────────────────
# 3. Paths
# ─────────────────────────────────────────────────────────────────────────────
PREFIX = "/api/v1/tenant/workforce-intelligence"
TAGS = [TAG]
SECURITY = [{"bearerAuth": []}]
QUERY_PAGE = {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}}
QUERY_PER_PAGE = {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}
QUERY_PERIOD = {"name": "period", "in": "query", "schema": {"type": "string"}}

def path_id():
    return {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}

def std_list_params(extra=None):
    params = [QUERY_PAGE, QUERY_PER_PAGE]
    if extra:
        params.extend(extra)
    return params

def ok_response(ref=None):
    r = {"200": {"description": "Success"}}
    if ref:
        r["200"]["content"] = {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}
    return r

def created_response(ref=None):
    r = {"201": {"description": "Created"}}
    if ref:
        r["201"]["content"] = {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}
    return r

def ok_wrapper(ref=None):
    r = {"200": {"description": "Success"}}
    if ref:
        r["200"]["content"] = {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}
    return r

def ok_array(ref=None):
    r = {"200": {"description": "Success — list of items"}}
    if ref:
        r["200"]["content"] = {"application/json": {"schema": {"type": "array", "items": {"$ref": f"#/components/schemas/{ref}"}}}}
    return r

def request_body(ref):
    return {"required": True, "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}}

# Build all paths
paths = {}

# ── Headcount Plans (5 routes) ──
paths[f"{PREFIX}/planning/headcounts"] = {
    "get": {
        "tags": TAGS, "summary": "List headcount plans (planned vs actual HC per period)",
        "operationId": "wfListHeadcountPlans",
        "parameters": std_list_params([QUERY_PERIOD, {"name": "organization_id", "in": "query", "schema": {"type": "string", "format": "uuid"}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of headcount plans. Filter by period (YYYY-MM) and/or organization unit. Each record shows planned headcount, actual headcount, and variance."
    },
    "post": {
        "tags": TAGS, "summary": "Create headcount plan",
        "operationId": "wfCreateHeadcountPlan",
        "security": SECURITY,
        "requestBody": request_body("CreateHeadcountPlanRequest"),
        "responses": created_response("HeadcountPlanResponse"),
        "description": "Create a new headcount plan for a specific period and organization. Records the planned headcount target for workforce planning purposes."
    }
}

paths[f"{PREFIX}/planning/headcounts/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get headcount plan by ID",
        "operationId": "wfGetHeadcountPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("HeadcountPlanResponse"),
        "description": "Get detailed information about a specific headcount plan including planned vs actual headcount and variance."
    },
    "put": {
        "tags": TAGS, "summary": "Update headcount plan",
        "operationId": "wfUpdateHeadcountPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateHeadcountPlanRequest"),
        "responses": ok_wrapper("HeadcountPlanResponse"),
        "description": "Update an existing headcount plan's planned headcount or snapshot date."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete headcount plan",
        "operationId": "wfDeleteHeadcountPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Headcount plan deleted"}},
        "description": "Remove a headcount plan record from the system."
    }
}

# ── Forecasts (5 routes) ──
paths[f"{PREFIX}/planning/forecasts"] = {
    "get": {
        "tags": TAGS, "summary": "List workforce forecasts",
        "operationId": "wfListForecasts",
        "parameters": std_list_params([QUERY_PERIOD, {"name": "organization_id", "in": "query", "schema": {"type": "string", "format": "uuid"}}, {"name": "forecast_type", "in": "query", "schema": {"type": "string", "enum": ["DEMAND", "SUPPLY", "HIRING"]}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of workforce forecasts. Filter by period, organization, and forecast type (DEMAND, SUPPLY, or HIRING)."
    },
    "post": {
        "tags": TAGS, "summary": "Create workforce forecast",
        "operationId": "wfCreateForecast",
        "security": SECURITY,
        "requestBody": request_body("CreateForecastRequest"),
        "responses": created_response("ForecastResponse"),
        "description": "Create a new workforce forecast. Supports DEMAND (required headcount), SUPPLY (available headcount), and HIRING (gap to fill) forecast types."
    }
}

paths[f"{PREFIX}/planning/forecasts/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get forecast by ID",
        "operationId": "wfGetForecast",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("ForecastResponse"),
        "description": "Get detailed information about a specific workforce forecast including headcount projection and confidence level."
    },
    "put": {
        "tags": TAGS, "summary": "Update forecast",
        "operationId": "wfUpdateForecast",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateForecastRequest"),
        "responses": ok_wrapper("ForecastResponse"),
        "description": "Update an existing workforce forecast's headcount and confidence level."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete forecast",
        "operationId": "wfDeleteForecast",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Forecast deleted"}},
        "description": "Remove a workforce forecast record."
    }
}

# ── Gap Analysis & Projections (2 routes) ──
paths[f"{PREFIX}/planning/gap-analysis"] = {
    "get": {
        "tags": TAGS, "summary": "Workforce gap analysis (supply vs demand)",
        "operationId": "wfGetGapAnalysis",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("GapAnalysisResponse"),
        "description": "Analyze workforce gaps by comparing supply (available HC) vs demand (required HC). Returns SURPLUS, SHORTAGE, or OPTIMAL status per department and overall."
    }
}

paths[f"{PREFIX}/planning/projections"] = {
    "get": {
        "tags": TAGS, "summary": "Workforce projections (hiring, retirement, growth)",
        "operationId": "wfGetProjections",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("ProjectionResponse"),
        "description": "Get workforce projections including projected headcount, hiring needs, retirement counts, and growth rates for strategic workforce planning."
    }
}

# ── KPIs (3 routes) ──
paths[f"{PREFIX}/kpi"] = {
    "get": {
        "tags": TAGS, "summary": "List KPIs",
        "operationId": "wfListKPIs",
        "parameters": std_list_params([QUERY_PERIOD, {"name": "dimension", "in": "query", "schema": {"type": "string"}}, {"name": "kpi_code", "in": "query", "schema": {"type": "string"}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of workforce KPIs. Filter by period, dimension, or KPI code. Each KPI includes value, target, unit, and snapshot date."
    }
}

paths[f"{PREFIX}/kpi/summary"] = {
    "get": {
        "tags": TAGS, "summary": "KPI summary (on-target vs below-target)",
        "operationId": "wfGetKPISummary",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("KPISummaryResponse"),
        "description": "Get a summary of KPIs showing total count, on-target count, and below-target count for a given period."
    }
}

paths[f"{PREFIX}/kpi/{{code}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get KPI by code",
        "operationId": "wfGetKPIByCode",
        "parameters": [{"name": "code", "in": "path", "required": True, "schema": {"type": "string"}}, QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("KPIResponse"),
        "description": "Get a specific KPI by its unique code (e.g., attrition_rate, turnover_rate, retention_rate, etc.)."
    }
}

# ── Analytics Dashboards (9 routes) ──
ANALYTICS = [
    ("headcount", "Headcount analytics dashboard", "HeadcountAnalytics", "Analyze workforce composition: total HC, active HC, distribution by department, employment type, gender, and headcount trend over time."),
    ("attendance", "Attendance analytics dashboard", "AttendanceAnalytics", "Analyze attendance metrics: average attendance rate, late rate, absentee rate, trend, and department breakdown."),
    ("leave", "Leave analytics dashboard", "LeaveAnalytics", "Analyze leave utilization: average utilization rate, total days taken, breakdown by leave type and department."),
    ("overtime", "Overtime analytics dashboard", "OvertimeAnalytics", "Analyze overtime patterns: average OT hours, total OT cost, department breakdown, and trend over time."),
    ("payroll", "Payroll analytics dashboard", "PayrollAnalytics", "Analyze payroll metrics: total payroll cost, average salary, breakdown by department and grade, with trend analysis."),
    ("performance", "Performance analytics dashboard", "PerformanceAnalytics", "Analyze employee performance: average score, top performer percentage, department breakdown, and score distribution."),
    ("learning", "Learning analytics dashboard", "LearningAnalytics", "Analyze learning and development: completion rate, average score, total training hours, and breakdown by course."),
    ("recruitment", "Recruitment analytics dashboard", "RecruitmentAnalytics", "Analyze recruitment efficiency: time to hire, cost per hire, source breakdown, and recruitment funnel metrics."),
    ("movement", "Movement analytics dashboard", "MovementAnalytics", "Analyze employee movement: promotion and mutation counts with breakdown by department and movement type.")
]

for slug, summary, schema_ref, desc in ANALYTICS:
    paths[f"{PREFIX}/analytics/{slug}"] = {
        "get": {
            "tags": TAGS,
            "summary": summary,
            "operationId": f"wfAnalytics{slug.capitalize()}",
            "parameters": [QUERY_PERIOD, {"name": "department_id", "in": "query", "schema": {"type": "string", "format": "uuid"}}],
            "security": SECURITY,
            "responses": ok_wrapper(schema_ref),
            "description": desc
        }
    }

# ── Capacity (4 routes) ──
paths[f"{PREFIX}/capacity/dashboard"] = {
    "get": {
        "tags": TAGS, "summary": "Capacity dashboard",
        "operationId": "wfGetCapacityDashboard",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("CapacityDashboardResponse"),
        "description": "Get workforce capacity dashboard: overall utilization rate, available headcount, department breakdown, and bottleneck identification."
    }
}

paths[f"{PREFIX}/capacity/utilization"] = {
    "get": {
        "tags": TAGS, "summary": "Resource utilization rate",
        "operationId": "wfGetUtilization",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_array("DataPoint"),
        "description": "Get workforce utilization rate data points. Shows how effectively the workforce is being utilized."
    }
}

paths[f"{PREFIX}/capacity/forecast"] = {
    "get": {
        "tags": TAGS, "summary": "Capacity forecast",
        "operationId": "wfGetCapacityForecast",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("CapacityForecastResponse"),
        "description": "Get projected capacity forecast: future utilization, current vs projected needed headcount, and capacity gap analysis by department."
    }
}

paths[f"{PREFIX}/capacity/bottlenecks"] = {
    "get": {
        "tags": TAGS, "summary": "Bottleneck analysis",
        "operationId": "wfGetBottlenecks",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_array("Bottleneck"),
        "description": "Identify capacity bottlenecks across departments. Flags departments with WARNING or CRITICAL utilization levels."
    }
}

# ── Cost (5 routes) ──
paths[f"{PREFIX}/cost/summary"] = {
    "get": {
        "tags": TAGS, "summary": "Cost summary dashboard",
        "operationId": "wfGetCostSummary",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("CostSummaryResponse"),
        "description": "Get workforce cost summary: total payroll, total benefit, total labor cost, cost per employee, department breakdown, and budget vs actual comparison."
    }
}

paths[f"{PREFIX}/cost/payroll"] = {
    "get": {
        "tags": TAGS, "summary": "Payroll cost breakdown",
        "operationId": "wfGetPayrollCostBreakdown",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("PayrollCostResponse"),
        "description": "Get detailed payroll cost breakdown: total salary, allowances, deductions, BPJS contributions, by grade and component."
    }
}

paths[f"{PREFIX}/cost/per-employee"] = {
    "get": {
        "tags": TAGS, "summary": "Cost per employee analysis",
        "operationId": "wfGetCostPerEmployee",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("CostPerEmployeeResponse"),
        "description": "Get cost per employee metrics: average, median, minimum, and maximum cost per employee with breakdown by department and grade."
    }
}

paths[f"{PREFIX}/cost/per-department"] = {
    "get": {
        "tags": TAGS, "summary": "Cost by department",
        "operationId": "wfGetCostByDepartment",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_array("DataPoint"),
        "description": "Get workforce cost broken down by department for comparison and budget allocation analysis."
    }
}

paths[f"{PREFIX}/cost/budget-vs-actual"] = {
    "get": {
        "tags": TAGS, "summary": "Budget vs actual cost comparison",
        "operationId": "wfGetBudgetVsActual",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_array("DataPoint"),
        "description": "Get budget vs actual workforce cost comparison. Shows budget targets versus actual spending for cost control analysis."
    }
}

# ── Risk (8 routes) ──
paths[f"{PREFIX}/risk/dashboard"] = {
    "get": {
        "tags": TAGS, "summary": "Risk dashboard",
        "operationId": "wfGetRiskDashboard",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("RiskDashboardResponse"),
        "description": "Get risk dashboard overview: total risks, high/critical risk counts, and list of active risk indicators with their levels and scores."
    }
}

paths[f"{PREFIX}/risk/indicators"] = {
    "get": {
        "tags": TAGS, "summary": "List risk indicators",
        "operationId": "wfListRiskIndicators",
        "parameters": std_list_params([QUERY_PERIOD, {"name": "risk_level", "in": "query", "schema": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "CRITICAL"]}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of risk indicators. Filter by period and/or risk level (LOW, MEDIUM, HIGH, CRITICAL)."
    }
}

paths[f"{PREFIX}/risk/indicators/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get risk indicator by ID",
        "operationId": "wfGetRiskIndicator",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("RiskResponse"),
        "description": "Get detailed information about a specific risk indicator including current score, threshold, and recommendations."
    },
    "put": {
        "tags": TAGS, "summary": "Update risk indicator",
        "operationId": "wfUpdateRiskIndicator",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateRiskRequest"),
        "responses": ok_wrapper("RiskResponse"),
        "description": "Update a risk indicator's level and/or recommendation. Used to acknowledge and document mitigation actions."
    }
}

# Risk detail types (4 routes)
RISK_TYPES = [
    ("high-turnover", "High turnover risk analysis", "Analyze high turnover risk: current turnover rate vs threshold, affected departments, trend, and recommended interventions."),
    ("retirement", "Retirement risk analysis", "Analyze retirement risk: upcoming retirements, impacted roles, succession gaps, and recommendations for knowledge transfer."),
    ("contract-expiry", "Contract expiration risk analysis", "Analyze contract expiry risk: upcoming contract/PKWT expirations, departments affected, and renewal recommendations."),
    ("high-absenteeism", "High absenteeism risk analysis", "Analyze high absenteeism risk: current absenteeism rate vs threshold, affected departments, trend, and intervention recommendations.")
]

for slug, summary, desc in RISK_TYPES:
    # Convert kebab-case to camelCase operationId
    parts = slug.split('-')
    op_id = 'wfRisk' + parts[0].capitalize() + ''.join(p.capitalize() for p in parts[1:])
    paths[f"{PREFIX}/risk/{slug}"] = {
        "get": {
            "tags": TAGS,
            "summary": summary,
            "operationId": op_id,
            "parameters": [QUERY_PERIOD],
            "security": SECURITY,
            "responses": ok_wrapper("RiskDetailResponse"),
            "description": desc
        }
    }

# ── Executive Dashboard (8 routes) ──
paths[f"{PREFIX}/executive/summary"] = {
    "get": {
        "tags": TAGS, "summary": "Executive workforce summary",
        "operationId": "wfGetExecutiveSummary",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("ExecutiveSummaryResponse"),
        "description": "Executive dashboard summary: total HC, HC growth, attrition rate, average cost, utilization rate, and overall health score in a single view."
    }
}

EXECUTIVE = [
    ("growth", "Executive workforce growth trend", "ExecutiveTrendResponse", "Executive-level workforce growth trend analysis showing HC change over time."),
    ("cost-trend", "Executive cost trend", "ExecutiveTrendResponse", "Executive-level workforce cost trend analysis for budget planning and cost control."),
    ("attrition-trend", "Executive attrition trend", "ExecutiveTrendResponse", "Executive-level attrition rate trend analysis to monitor talent retention."),
    ("capacity", "Executive capacity overview", "ExecutiveCapacityResponse", "Executive-level capacity overview: utilization rate, available HC, active departments, and bottleneck count."),
    ("hiring-progress", "Hiring progress tracker", "HiringProgressResponse", "Track hiring progress: planned, in-progress, and completed hires for executive monitoring."),
    ("risk-overview", "Executive risk overview", "ExecutiveRiskOverviewResponse", "Executive-level risk overview: total risks, high and critical counts, by department and category."),
    ("health-score", "Executive health score", "ExecutiveHealthScoreResponse", "Executive-level organization health score with span of control, manager ratio, internal hiring rate, succession coverage, and overall status (HEALTHY/WARNING/CRITICAL).")
]

for slug, summary, schema_ref, desc in EXECUTIVE:
    # Convert kebab-case to camelCase operationId
    parts = slug.split('-')
    op_id = 'wfExecutive' + parts[0].capitalize() + ''.join(p.capitalize() for p in parts[1:])
    paths[f"{PREFIX}/executive/{slug}"] = {
        "get": {
            "tags": TAGS,
            "summary": summary,
            "operationId": op_id,
            "parameters": [QUERY_PERIOD],
            "security": SECURITY,
            "responses": ok_wrapper(schema_ref),
            "description": desc
        }
    }

# ── Scenario Planning (7 routes) ──
paths[f"{PREFIX}/scenarios"] = {
    "get": {
        "tags": TAGS, "summary": "List scenarios",
        "operationId": "wfListScenarios",
        "parameters": std_list_params([{"name": "scenario_type", "in": "query", "schema": {"type": "string", "enum": ["NEW_BRANCH", "REORG", "GROWTH", "REDUCTION", "RETIREMENT", "BUDGET"]}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of saved simulation scenarios. Filter by scenario type (NEW_BRANCH, REORG, GROWTH, REDUCTION, RETIREMENT, BUDGET)."
    },
    "post": {
        "tags": TAGS, "summary": "Create scenario",
        "operationId": "wfCreateScenario",
        "security": SECURITY,
        "requestBody": request_body("CreateScenarioRequest"),
        "responses": created_response("ScenarioResponse"),
        "description": "Create a new scenario for workforce simulation. Supports NEW_BRANCH, REORG, GROWTH, REDUCTION, RETIREMENT, and BUDGET scenario types. Parameters are JSON-defined per type."
    }
}

paths[f"{PREFIX}/scenarios/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get scenario by ID",
        "operationId": "wfGetScenario",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("ScenarioResponse"),
        "description": "Get detailed information about a specific scenario including its parameters, results (if run), and status (DRAFT, RUNNING, COMPLETED)."
    },
    "put": {
        "tags": TAGS, "summary": "Update scenario",
        "operationId": "wfUpdateScenario",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateScenarioRequest"),
        "responses": ok_wrapper("ScenarioResponse"),
        "description": "Update an existing scenario's name, description, and/or parameters."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete scenario",
        "operationId": "wfDeleteScenario",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Scenario deleted"}},
        "description": "Soft-delete a scenario by ID."
    }
}

paths[f"{PREFIX}/scenarios/{{id}}/run"] = {
    "post": {
        "tags": TAGS, "summary": "Run scenario simulation",
        "operationId": "wfRunScenario",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("ScenarioResponse"),
        "description": "Execute a scenario simulation. Runs the scenario's parameters through the simulation engine and stores results in the scenario record."
    }
}

paths[f"{PREFIX}/scenarios/{{id}}/clone"] = {
    "post": {
        "tags": TAGS, "summary": "Clone scenario",
        "operationId": "wfCloneScenario",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": created_response("ScenarioResponse"),
        "description": "Clone an existing scenario as a new DRAFT scenario. Useful for creating variations of a simulation without affecting the original."
    }
}

# ── Organization Health (5 routes) ──
paths[f"{PREFIX}/health/dashboard"] = {
    "get": {
        "tags": TAGS, "summary": "Organization health dashboard",
        "operationId": "wfGetHealthDashboard",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("HealthScoreResponse"),
        "description": "Get organization health dashboard: composite health score with span of control, manager ratio, promotion rate, internal hiring rate, succession coverage, and stability ratio."
    }
}

paths[f"{PREFIX}/health/scores"] = {
    "get": {
        "tags": TAGS, "summary": "List health scores",
        "operationId": "wfListHealthScores",
        "parameters": std_list_params([QUERY_PERIOD, {"name": "organization_id", "in": "query", "schema": {"type": "string", "format": "uuid"}}]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseWF"),
        "description": "Retrieve a paginated list of organization health scores. Filter by period and organization unit."
    }
}

paths[f"{PREFIX}/health/scores/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get health score by ID",
        "operationId": "wfGetHealthScoreByID",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("HealthScoreResponse"),
        "description": "Get detailed health score components for a specific score record."
    }
}

paths[f"{PREFIX}/health/span-of-control"] = {
    "get": {
        "tags": TAGS, "summary": "Span of control analysis",
        "operationId": "wfGetSpanOfControl",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("SpanOfControlResponse"),
        "description": "Analyze span of control: average manager-to-report ratio, healthy range (3:1-7:1), status indicator, and department breakdown."
    }
}

paths[f"{PREFIX}/health/succession"] = {
    "get": {
        "tags": TAGS, "summary": "Succession readiness analysis",
        "operationId": "wfGetSuccessionReadiness",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("SuccessionReadinessResponse"),
        "description": "Analyze succession readiness: roles with identified successors, total key roles, coverage rate, status (HEALTHY/WARNING/CRITICAL), and department breakdown."
    }
}

# ── People Analytics (7 routes) ──
PEOPLE_ANALYTICS = [
    ("training-vs-performance", "Training vs performance correlation", "wfPeopleTrainingVsPerformance", "Analyze correlation between training participation and employee performance scores. Helps evaluate training ROI and effectiveness."),
    ("overtime-vs-productivity", "Overtime vs productivity correlation", "wfPeopleOvertimeVsProductivity", "Analyze correlation between overtime hours and productivity metrics to identify optimal workload levels."),
    ("attendance-vs-performance", "Attendance vs performance correlation", "wfPeopleAttendanceVsPerformance", "Analyze correlation between attendance rate and employee performance to identify attendance-related performance risks."),
    ("compensation-vs-turnover", "Compensation vs turnover correlation", "wfPeopleCompensationVsTurnover", "Analyze correlation between compensation levels and turnover rates to evaluate compensation competitiveness and retention strategies."),
    ("source-vs-retention", "Recruitment source vs retention correlation", "wfPeopleSourceVsRetention", "Analyze correlation between recruitment source (job board, referral, agency, etc.) and employee retention rates to optimize sourcing strategy."),
    ("career-progression", "Career progression vs performance correlation", "wfPeopleCareerProgression", "Analyze correlation between career advancement (promotions, movements) and employee performance trends."),
    ("learning-effectiveness", "Learning effectiveness analysis", "wfPeopleLearningEffectiveness", "Analyze the effectiveness of learning programs by correlating training completion, scores, and subsequent performance improvements.")
]

for slug, summary, operation_id, desc in PEOPLE_ANALYTICS:
    paths[f"{PREFIX}/people-analytics/{slug}"] = {
        "get": {
            "tags": TAGS,
            "summary": summary,
            "operationId": operation_id,
            "parameters": [QUERY_PERIOD],
            "security": SECURITY,
            "responses": ok_wrapper("CorrelationResponse"),
            "description": desc
        }
    }

# ─────────────────────────────────────────────────────────────────────────────
# 4. Merge paths into the spec
# ─────────────────────────────────────────────────────────────────────────────
for path, methods in paths.items():
    if path not in spec['paths']:
        spec['paths'][path] = methods
    else:
        # Merge methods without overwriting existing ones
        for method, endpoint in methods.items():
            if method not in spec['paths'][path]:
                spec['paths'][path][method] = endpoint

# ─────────────────────────────────────────────────────────────────────────────
# 5. Write output
# ─────────────────────────────────────────────────────────────────────────────
spec['info']['description'] = spec['info']['description'].replace(
    "workforce.view, workforce.create, workforce.update, workforce.delete",
    "workforce.view, workforce.create, workforce.update, workforce.delete"
)

# Add workforce permission to the RBAC permissions block if not already present
if "workforce." not in spec['info']['description']:
    spec['info']['description'] = spec['info']['description'].replace(
        "reimbursement.view, reimbursement.create, reimbursement.update, reimbursement.delete, reimbursement.approve\n\n",
        "reimbursement.view, reimbursement.create, reimbursement.update, reimbursement.delete, reimbursement.approve\n\nworkforce.view\n\n"
    )

with open(OPENAPI_PATH, 'w', encoding='utf-8') as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

# ─────────────────────────────────────────────────────────────────────────────
# 6. Summary
# ─────────────────────────────────────────────────────────────────────────────
wf_paths = {p: list(v.keys()) for p, v in spec['paths'].items() if 'workforce-intelligence' in p}
total_endpoints = sum(len(m) for m in wf_paths.values())
total_schemas = len(spec['components']['schemas'])
total_paths = len(spec['paths'])
total_tags = len(spec['tags'])

new_schemas = [n for n in schemas if n in spec['components']['schemas']]

print(f"--- RESULTS ---")
print(f"WFi endpoints added: {total_endpoints}")
print(f"Total unique paths: {len(wf_paths)}")
print(f"New schemas added: {len(new_schemas)}")
print(f"")
print(f"UPDATED STATE:")
print(f"Total paths: {total_paths}")
print(f"Total schemas: {total_schemas}")
print(f"Total tags: {total_tags}")
print(f"")
print(f"New tag: {TAG}")
print(f"")
print(f"=== Path groups ===")
for p in sorted(wf_paths.keys()):
    methods = ", ".join(m.upper() for m in sorted(wf_paths[p]))
    print(f"  {methods} {p}")
